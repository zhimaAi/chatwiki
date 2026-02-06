// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/wechat/official_account"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/menu/request"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// 1Send Message 2Jump to Web Page 3Jump to Mini Program 4Customer Service 5Push Event
const (
	OfficialCustomMenuActTypeSendMessage     = 1
	OfficialCustomMenuActTypeJumpURL         = 2
	OfficialCustomMenuActTypeJumpMiniProgram = 3
	OfficialCustomMenuActTypeCustomerService = 4
	OfficialCustomMenuActTypePushEvent       = 5
)

const OfficialCustomMenuMenuLevelOne = 1
const OfficialCustomMenuMenuLevelTwo = 2

// OfficialCustomMenu Intelligent chatbot menu
type OfficialCustomMenu struct {
	ID            int                        `json:"id" form:"id"`                             // Auto-increment ID
	AdminUserID   int                        `json:"admin_user_id" form:"admin_user_id"`       // Admin user ID
	AppID         string                     `json:"appid" form:"appid"`                       // Official account appid
	SeqID         int                        `json:"seq_id" form:"seq_id"`                     // Sequence ID
	MenuName      string                     `json:"menu_name" form:"menu_name"`               // Menu name
	MenuLevel     int                        `json:"menu_level" form:"menu_level"`             // Menu level 1Root menu 2Sub menu
	ParentMenuID  int                        `json:"parent_menu_id" form:"parent_menu_id"`     // Parent node ID
	ChooseActItem int                        `json:"choose_act_item" form:"choose_act_item"`   // Default selected item 0No function when has child nodes 1Send message 2Jump to web page 3Jump to mini program 4Customer service 5Push event
	ActParams     OfficialCustomMenuActParam `json:"act_params" form:"act_params"`             // Configuration JSON string
	OperUserID    int                        `json:"oper_user_id" form:"oper_user_id"`         // Operator user ID
	TemplateID    int                        `json:"template_id,omitempty" form:"template_id"` // Template ID
	BatchID       int                        `json:"batch_id,omitempty" form:"batch_id"`       // Batch ID
	SubMenuList   []OfficialCustomMenu       `json:"sub_menu_list,omitempty" form:"sub_menu_list"`
	CreateTime    int                        `json:"create_time,omitempty"` // Creation time
	UpdateTime    int                        `json:"update_time,omitempty"` // Update time
}

type OfficialCustomMenuActParam struct {
	Item         int            `json:"item"`
	LinkURL      string         `json:"linkUrl,omitempty"`
	OpenSwitch   string         `json:"open_switch,omitempty"`
	ReplyContent []ReplyContent `json:"reply_content,omitempty"`
	ReplyNum     int            `json:"reply_num,omitempty"`
	Appid        string         `json:"appid,omitempty"`
	Pagepath     string         `json:"pagepath,omitempty"`
	StandbyURL   string         `json:"standbyUrl,omitempty"`
	Key          string         `json:"key,omitempty"`
}

// OfficialCustomMenuCacheBuildHandler Menu cache handler
type OfficialCustomMenuCacheBuildHandler struct {
	ID int
}

func (h *OfficialCustomMenuCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`official_custom_menu_%d`, h.ID)
}

func (h *OfficialCustomMenuCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where(`id`, cast.ToString(h.ID)).
		Find()
	if err != nil {
		return nil, err
	}

	// Convert
	if len(data) > 0 {
		var actParams OfficialCustomMenuActParam
		_ = tool.JsonDecodeUseNumber(data[`act_params`], &actParams)

		result := OfficialCustomMenu{
			ID:            cast.ToInt(data[`id`]),
			AdminUserID:   cast.ToInt(data[`admin_user_id`]),
			AppID:         data[`appid`],
			SeqID:         cast.ToInt(data[`seq_id`]),
			MenuName:      data[`menu_name`],
			MenuLevel:     cast.ToInt(data[`menu_level`]),
			ParentMenuID:  cast.ToInt(data[`parent_menu_id`]),
			ChooseActItem: cast.ToInt(data[`choose_act_item`]),
			ActParams:     actParams,
			OperUserID:    cast.ToInt(data[`oper_user_id`]),
			TemplateID:    cast.ToInt(data[`template_id`]),
			BatchID:       cast.ToInt(data[`batch_id`]),
			CreateTime:    cast.ToInt(data[`create_time`]),
			UpdateTime:    cast.ToInt(data[`update_time`]),
		}
		return result, err
	} else {
		return nil, nil
	}
}

// GetOfficialCustomMenuInfo Get custom menu information
func GetOfficialCustomMenuInfo(id int) (OfficialCustomMenu, error) {
	result := OfficialCustomMenu{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// SaveOfficialCustomMenu Save custom menu (create or update)
// Update menu instead of deleting and recreating, keeping menu ID unchanged to maintain association with WeChat
func SaveOfficialCustomMenu(adminUserID int, appid string, operUserID int, menuList []OfficialCustomMenu, sendWx bool) (int, error) {
	// Begin transaction
	tx, err := msql.Begin(define.Postgres)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	//
	oldMenus, err := GetOfficialCustomMenuList(adminUserID, appid)

	// Get menus currently existing in database
	existingMenus, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Select()
	if err != nil {
		return 0, err
	}

	// Create mapping of existing menu IDs
	existingMenuMap := make(map[int]msql.Params)
	for _, menu := range existingMenus {
		menuID := cast.ToInt(menu["id"])
		existingMenuMap[menuID] = menu
	}

	// Parse main menus
	var parentMenus []OfficialCustomMenu

	for _, menu := range menuList {
		if menu.MenuLevel == OfficialCustomMenuMenuLevelOne {
			parentMenus = append(parentMenus, menu)
		}
	}

	// Used to track processed menu IDs, to delete menus that are no longer needed
	processedMenuIDs := make(map[int]bool)

	// Update or create main menus and their corresponding sub-menus
	parentIDMap := make(map[int]int) // oldID -> newID (for newly created menus)
	for i, menu := range parentMenus {
		actParamsJson, _ := tool.JsonEncode(menu.ActParams)

		data := msql.Datas{
			"admin_user_id":   adminUserID,
			"appid":           appid,
			"seq_id":          menu.SeqID,
			"menu_name":       menu.MenuName,
			"menu_level":      menu.MenuLevel,
			"choose_act_item": menu.ChooseActItem,
			"act_params":      actParamsJson,
			"oper_user_id":    operUserID,
			"template_id":     menu.TemplateID,
			"batch_id":        menu.BatchID,
			"update_time":     tool.Time2Int(),
		}

		var newId int64
		if menu.ID > 0 && existingMenuMap[menu.ID] != nil {
			// Update existing menu
			_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
				Where("id", cast.ToString(menu.ID)).
				Update(data)
			if err != nil {
				return 0, err
			}
			newId = int64(menu.ID)
			processedMenuIDs[menu.ID] = true
		} else {
			// Create new menu
			data["parent_menu_id"] = 0
			data["create_time"] = tool.Time2Int()
			newId, err = msql.Model(`func_official_custom_menu`, define.Postgres).Insert(data, "id")
			if err != nil {
				return 0, err
			}
			// Record mapping of newly created menu IDs
			if menu.ID > 0 {
				parentIDMap[menu.ID] = int(newId)
			}
			processedMenuIDs[int(newId)] = true
		}

		// Update parent menu ID mapping, ensuring correct mapping even if original ID does not exist
		if menu.ID > 0 {
			parentIDMap[menu.ID] = int(newId)
		}
		parentMenus[i].ID = int(newId)

		// Process sub-menus under this main menu
		for j := range parentMenus[i].SubMenuList {
			parentMenus[i].SubMenuList[j].MenuLevel = OfficialCustomMenuMenuLevelTwo
			parentMenus[i].SubMenuList[j].ParentMenuID = int(newId)

			subMenu := parentMenus[i].SubMenuList[j]
			subActParamsJson, _ := tool.JsonEncode(subMenu.ActParams)

			subData := msql.Datas{
				"admin_user_id":   adminUserID,
				"appid":           appid,
				"seq_id":          subMenu.SeqID,
				"menu_name":       subMenu.MenuName,
				"menu_level":      subMenu.MenuLevel,
				"parent_menu_id":  subMenu.ParentMenuID,
				"choose_act_item": subMenu.ChooseActItem,
				"act_params":      subActParamsJson,
				"oper_user_id":    operUserID,
				"template_id":     subMenu.TemplateID,
				"batch_id":        subMenu.BatchID,
				"update_time":     tool.Time2Int(),
			}

			var subNewId int64
			if subMenu.ID > 0 && existingMenuMap[subMenu.ID] != nil {
				// Update existing sub-menu
				_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
					Where("id", cast.ToString(subMenu.ID)).
					Update(subData)
				if err != nil {
					return 0, err
				}
				subNewId = int64(subMenu.ID)
				processedMenuIDs[subMenu.ID] = true
			} else {
				// Create new sub-menu
				subData["create_time"] = tool.Time2Int()
				subNewId, err = msql.Model(`func_official_custom_menu`, define.Postgres).Insert(subData, "id")
				if err != nil {
					return 0, err
				}
				processedMenuIDs[int(subNewId)] = true
			}
			parentMenus[i].SubMenuList[j].ID = int(subNewId)
		}
	}

	// Delete unprocessed old menus (i.e. menus that have been deleted)
	for menuID := range existingMenuMap {
		if !processedMenuIDs[menuID] {
			_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
				Where("id", cast.ToString(menuID)).
				Delete()
			if err != nil {
				return 0, err
			}
			// Clear cache of deleted menu
			lib_redis.DelCacheData(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: menuID})
		}
	}

	// Clear all related caches
	for menuID := range processedMenuIDs {
		lib_redis.DelCacheData(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: menuID})
	}

	// Record menu change history
	err = saveMenuHistory(adminUserID, appid, operUserID, oldMenus)
	if err != nil {
		return 0, err
	}
	if sendWx {
		// Send local menu to WeChat
		return SendOfficialCustomMenuToWx(adminUserID, appid)
	}
	// Send local menu to WeChat
	return 0, nil
}

// SendOfficialCustomMenuToWx Send custom menu
func SendOfficialCustomMenuToWx(adminUserID int, appid string) (int, error) {
	appInfo, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserID)).Where(`app_id`, appid).Find()
	if err != nil {
		return 0, err
	}

	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}

	// Get menu and convert to request.RequestMenuCreate

	localMenuList, err := GetOfficialCustomMenuList(adminUserID, appid)
	if err != nil {
		return 0, err
	}
	// Convert to request.RequestMenuCreate
	menu := LocalMenuListConvertToRequestMenuCreate(localMenuList)
	if len(menu.Buttons) == 0 {
		return 0, errors.New(`the menu is empty`)
	}
	// Process
	setMenuCode, err := app.SetMenu(menu)
	if err == nil {
		// Save menu enabled status
		_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
			Where(`id`, appInfo[`id`]).
			Where(`admin_user_id`, appInfo[`admin_user_id`]).
			Update(msql.Datas{
				"custom_menu_status": define.SwitchOn,
			})
	}

	return setMenuCode, err
}

// LocalMenuListConvertToRequestMenuCreate Convert local menu list to request.RequestMenuCreate
func LocalMenuListConvertToRequestMenuCreate(localMenuList []OfficialCustomMenu) request.RequestMenuCreate {
	var menu request.RequestMenuCreate
	buttons := make([]*request.Button, 0)
	for _, localMenu := range localMenuList {

		curButton := request.Button{
			Name: localMenu.MenuName,
		}
		// No longer need to loop through actParams, as ActParams is now a single object rather than an array
		act := localMenu.ActParams
		switch act.Item {
		case OfficialCustomMenuActTypeSendMessage:
			curButton.Type = lib_define.MenuButtonTypeClick
			curButton.Key = cast.ToString(localMenu.ID)
			break
		case OfficialCustomMenuActTypeJumpURL:
			curButton.Type = lib_define.MenuButtonTypeView
			curButton.URL = act.LinkURL
			break
		case OfficialCustomMenuActTypeJumpMiniProgram:
			curButton.Type = lib_define.MenuButtonTypeMiniprogram
			curButton.AppID = act.Appid
			curButton.PagePath = act.Pagepath
			curButton.URL = act.StandbyURL
			break
		case OfficialCustomMenuActTypePushEvent:
			curButton.Type = lib_define.MenuButtonTypeClick
			curButton.Key = act.Key
			break
		}

		// Process sub-menu
		if len(localMenu.SubMenuList) > 0 {
			subButtons := make([]request.SubButton, 0)
			for _, subMenu := range localMenu.SubMenuList {
				subButton := request.SubButton{
					Name: subMenu.MenuName,
				}

				switch subMenu.ActParams.Item {
				case OfficialCustomMenuActTypeSendMessage:
					subButton.Type = lib_define.MenuButtonTypeClick
					subButton.Key = cast.ToString(subMenu.ID)
					break
				case OfficialCustomMenuActTypeJumpURL:
					subButton.Type = lib_define.MenuButtonTypeView
					subButton.URL = subMenu.ActParams.LinkURL
					break
				case OfficialCustomMenuActTypeJumpMiniProgram:
					subButton.Type = lib_define.MenuButtonTypeMiniprogram
					subButton.AppID = subMenu.ActParams.Appid
					subButton.PagePath = subMenu.ActParams.Pagepath
					subButton.URL = subMenu.ActParams.StandbyURL
					break
				case OfficialCustomMenuActTypePushEvent:
					subButton.Type = lib_define.MenuButtonTypeClick
					subButton.Key = subMenu.ActParams.Key
					break
				}
				subButtons = append(subButtons, subButton)
			}
			curButton.SubButtons = subButtons
		}
		buttons = append(buttons, &curButton)
	}
	menu.Buttons = buttons
	return menu
}

// SyncWxMenuToShow Synchronize WeChat menu to display
func SyncWxMenuToShow(adminUserID int, appid string, operUserID int) ([]OfficialCustomMenu, error) {
	appInfo, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserID)).Where(`app_id`, appid).Find()
	if err != nil {
		return nil, err
	}
	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}

	localMenuList := make([]OfficialCustomMenu, 0)
	WxMenu, err := app.GetMenu()
	if err != nil {
		// Error if contains menu no exist
		if strings.Contains(err.Error(), `menu no exist`) {
			localMenuList = CheckUseDefaultMenu(adminUserID, appid, localMenuList)
			return localMenuList, nil
		}
		return nil, err
	}

	if len(WxMenu.Menus.Buttons) == 0 {
		localMenuList = CheckUseDefaultMenu(adminUserID, appid, localMenuList)
		return localMenuList, nil
	}
	// Convert to local menu
	for menuIndex, button := range WxMenu.Menus.Buttons {
		localMenu := OfficialCustomMenu{
			MenuName:     button.Name,
			AppID:        appid,
			SeqID:        menuIndex,
			MenuLevel:    OfficialCustomMenuMenuLevelOne,
			ParentMenuID: 0,
			OperUserID:   operUserID,
			AdminUserID:  adminUserID,
		}
		switch button.Type {
		case lib_define.MenuButtonTypeClick:
			if tool.IsNumeric(button.Key) {
				localMenu.ChooseActItem = OfficialCustomMenuActTypeSendMessage
				localMenu.ActParams.Item = OfficialCustomMenuActTypeSendMessage
				// Get details of corresponding menu
				getMenu, err := GetOfficialCustomMenuInfo(cast.ToInt(button.Key))
				if err == nil && getMenu.ID > 0 {
					localMenu.ID = getMenu.ID
					localMenu.ActParams = getMenu.ActParams
					localMenu.ChooseActItem = getMenu.ChooseActItem
					break
				}
			}
			localMenu.ChooseActItem = OfficialCustomMenuActTypePushEvent
			localMenu.ActParams.Item = OfficialCustomMenuActTypePushEvent
			localMenu.ActParams.Key = button.Key
			break
		case lib_define.MenuButtonTypeView:
			localMenu.ChooseActItem = OfficialCustomMenuActTypeJumpURL
			localMenu.ActParams.Item = OfficialCustomMenuActTypeJumpURL
			localMenu.ActParams.LinkURL = button.URL
			break
		case lib_define.MenuButtonTypeMiniprogram:
			localMenu.ChooseActItem = OfficialCustomMenuActTypeJumpMiniProgram
			localMenu.ActParams.Item = OfficialCustomMenuActTypeJumpMiniProgram
			localMenu.ActParams.Appid = button.AppID
			localMenu.ActParams.Pagepath = button.PagePath
			localMenu.ActParams.StandbyURL = button.URL
			break

		}
		if len(button.SubButtons) > 0 {
			for subMenuIndex, subButton := range button.SubButtons {
				subMenu := OfficialCustomMenu{
					MenuName:     subButton.Name,
					AppID:        appid,
					SeqID:        subMenuIndex,
					MenuLevel:    OfficialCustomMenuMenuLevelTwo,
					ParentMenuID: localMenu.ID,
					OperUserID:   operUserID,
					AdminUserID:  adminUserID,
				}
				switch subButton.Type {
				case lib_define.MenuButtonTypeClick:
					if tool.IsNumeric(subButton.Key) {
						subMenu.ChooseActItem = OfficialCustomMenuActTypeSendMessage
						subMenu.ActParams.Item = OfficialCustomMenuActTypeSendMessage
						// Get details of corresponding menu
						getMenu, err := GetOfficialCustomMenuInfo(cast.ToInt(subButton.Key))
						if err == nil && getMenu.ID > 0 {
							subMenu.ID = getMenu.ID
							if localMenu.ID == 0 {
								localMenu.ParentMenuID = subMenu.ParentMenuID
							}
							subMenu.ActParams = getMenu.ActParams
							subMenu.ChooseActItem = getMenu.ChooseActItem
							break
						}
					}
					subMenu.ChooseActItem = OfficialCustomMenuActTypePushEvent
					subMenu.ActParams.Item = OfficialCustomMenuActTypePushEvent
					subMenu.ActParams.Key = subButton.Key
					break
				case lib_define.MenuButtonTypeView:
					subMenu.ChooseActItem = OfficialCustomMenuActTypeJumpURL
					subMenu.ActParams.Item = OfficialCustomMenuActTypeJumpURL
					subMenu.ActParams.LinkURL = subButton.URL
					break
				case lib_define.MenuButtonTypeMiniprogram:
					subMenu.ChooseActItem = OfficialCustomMenuActTypeJumpMiniProgram
					subMenu.ActParams.Item = OfficialCustomMenuActTypeJumpMiniProgram
					subMenu.ActParams.Appid = subButton.AppID
					subMenu.ActParams.Pagepath = subButton.PagePath
					subMenu.ActParams.StandbyURL = subButton.URL
					break
				}
				localMenu.SubMenuList = append(localMenu.SubMenuList, subMenu)
			}
		}
		localMenuList = append(localMenuList, localMenu)
	}
	// Check if default menu is used
	localMenuList = CheckUseDefaultMenu(adminUserID, appid, localMenuList)
	// Convert to local menu and return to frontend
	return localMenuList, err
}

// saveMenuHistory Save menu change history
func saveMenuHistory(adminUserID int, appid string, operUserID int, oldMenus []OfficialCustomMenu) error {
	// Get menus in current database (complete primary and secondary menu structure)
	currentMenuList, err := GetOfficialCustomMenuList(adminUserID, appid)
	if err != nil {
		return err
	}

	// Clean time fields in current menu list
	cleanCurrentMenuList := make([]OfficialCustomMenu, len(currentMenuList))
	for i, menu := range currentMenuList {
		cleanMenu := menu
		cleanMenu.CreateTime = 0
		cleanMenu.UpdateTime = 0

		// Clean time fields of sub-menus
		cleanSubMenuList := make([]OfficialCustomMenu, len(cleanMenu.SubMenuList))
		for j, subMenu := range cleanMenu.SubMenuList {
			cleanSubMenu := subMenu
			cleanSubMenu.CreateTime = 0
			cleanSubMenu.UpdateTime = 0
			cleanSubMenuList[j] = cleanSubMenu
		}
		cleanMenu.SubMenuList = cleanSubMenuList

		cleanCurrentMenuList[i] = cleanMenu
	}

	//
	historyMenuList := oldMenus
	// Clean time fields in submitted menu list
	cleanHistoryMenuList := make([]OfficialCustomMenu, len(oldMenus))
	for i, menu := range historyMenuList {
		cleanMenu := menu
		cleanMenu.CreateTime = 0
		cleanMenu.UpdateTime = 0

		// Clean time fields of sub-menus
		cleanSubMenuList := make([]OfficialCustomMenu, len(cleanMenu.SubMenuList))
		for j, subMenu := range cleanMenu.SubMenuList {
			cleanSubMenu := subMenu
			cleanSubMenu.CreateTime = 0
			cleanSubMenu.UpdateTime = 0
			cleanSubMenuList[j] = cleanSubMenu
		}
		cleanMenu.SubMenuList = cleanSubMenuList

		cleanHistoryMenuList[i] = cleanMenu
	}

	// Use DeepEqual to compare two struct slices
	if !reflect.DeepEqual(cleanHistoryMenuList, cleanCurrentMenuList) {
		// Re-encode current menu (including time fields) for storage
		historyMenuJson, _ := tool.JsonEncode(cleanHistoryMenuList)
		currentMenuJson, _ := tool.JsonEncode(cleanCurrentMenuList)

		historyData := msql.Datas{
			"admin_user_id": adminUserID,
			"appid":         appid,
			"history_menu":  historyMenuJson,
			"mew_menu":      currentMenuJson,
			"oper_user_id":  operUserID,
			"create_time":   tool.Time2Int(),
			"update_time":   tool.Time2Int(),
		}

		_, err = msql.Model(`func_official_custom_menu_history`, define.Postgres).Insert(historyData)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOfficialCustomMenuList Get custom menu list
func GetOfficialCustomMenuList(adminUserID int, appid string) ([]OfficialCustomMenu, error) {
	// Get all menus
	menus, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Order("menu_level ASC, seq_id ASC, id ASC").
		Select()
	if err != nil {
		return nil, err
	}

	// Store primary and secondary menus separately
	parentMenus := make([]OfficialCustomMenu, 0)
	subMenuMap := make(map[int][]OfficialCustomMenu) // parent_menu_id -> []SubMenu

	// Process all menu items
	for _, item := range menus {
		var actParams OfficialCustomMenuActParam
		_ = tool.JsonDecodeUseNumber(item[`act_params`], &actParams)

		actParams.ReplyContent = FormatReplyListToDb(actParams.ReplyContent, OfficialAbilityCustomMenu)

		menu := OfficialCustomMenu{
			ID:            cast.ToInt(item[`id`]),
			AdminUserID:   cast.ToInt(item[`admin_user_id`]),
			AppID:         item[`appid`],
			SeqID:         cast.ToInt(item[`seq_id`]),
			MenuName:      item[`menu_name`],
			MenuLevel:     cast.ToInt(item[`menu_level`]),
			ParentMenuID:  cast.ToInt(item[`parent_menu_id`]),
			ChooseActItem: cast.ToInt(item[`choose_act_item`]),
			ActParams:     actParams,
			OperUserID:    cast.ToInt(item[`oper_user_id`]),
			TemplateID:    cast.ToInt(item[`template_id`]),
			BatchID:       cast.ToInt(item[`batch_id`]),
			CreateTime:    cast.ToInt(item[`create_time`]),
			UpdateTime:    cast.ToInt(item[`update_time`]),
		}

		if menu.MenuLevel == OfficialCustomMenuMenuLevelOne {
			parentMenus = append(parentMenus, menu)
		} else if menu.MenuLevel == OfficialCustomMenuMenuLevelTwo {
			// Group secondary menus by parent_menu_id
			subMenuMap[menu.ParentMenuID] = append(subMenuMap[menu.ParentMenuID], menu)
		}
	}

	// Attach secondary menus to corresponding primary menu
	for i := range parentMenus {
		if subMenus, exists := subMenuMap[parentMenus[i].ID]; exists {
			parentMenus[i].SubMenuList = subMenus
		}
	}
	// Check if default menu is used
	parentMenus = CheckUseDefaultMenu(adminUserID, appid, parentMenus)

	return parentMenus, nil
}

// CheckUseDefaultMenu Check if default menu is used
func CheckUseDefaultMenu(adminUserID int, appid string, parentMenus []OfficialCustomMenu) []OfficialCustomMenu {
	if len(parentMenus) == 0 {
		sendReplyContent := make([]ReplyContent, 0)
		sendReplyContent = append(sendReplyContent, ReplyContent{
			Type:        `text`,
			Description: ``,
		})
		defaultMenu := OfficialCustomMenu{
			ID:            0,
			AppID:         appid,
			MenuName:      lib_define.MainMenu,
			MenuLevel:     OfficialCustomMenuMenuLevelOne,
			OperUserID:    adminUserID,
			CreateTime:    tool.Time2Int(),
			UpdateTime:    tool.Time2Int(),
			SubMenuList:   []OfficialCustomMenu{},
			ChooseActItem: OfficialCustomMenuActTypeSendMessage,
			ActParams: OfficialCustomMenuActParam{
				Item:         OfficialCustomMenuActTypeSendMessage,
				ReplyContent: sendReplyContent,
			},
		}
		parentMenus = append(parentMenus, defaultMenu)
	}
	return parentMenus
}

// DeleteAllOfficialCustomMenuToWx Delete all official account custom menus
func DeleteAllOfficialCustomMenuToWx(adminUserId int, appid string) string {
	// Get all official accounts
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	m.Where(`app_type`, lib_define.AppOfficeAccount)
	if appid != `` {
		m.Where(`app_id`, appid)
	}
	m.Order(`sort asc`)
	list, err := m.Order(`id desc`).Select()
	if err != nil {
		return err.Error()
	}
	errStr := ``
	for _, appInfo := range list {
		err = WxDeleteMenu(appInfo)
		if err != nil {
			wxName := appInfo[`app_name`]
			errStr += wxName + ` error：` + err.Error() + `\n`
		}
	}
	return errStr
}

func WxDeleteMenu(appInfo msql.Params) error {
	// Capture error
	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}
	_, err := app.DeleteMenu()
	if err != nil {
		return err
	}
	// Update official account menu switch status
	_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`id`, appInfo[`id`]).
		Where(`admin_user_id`, appInfo[`admin_user_id`]).
		Update(msql.Datas{
			"custom_menu_status": define.SwitchOff,
		})
	return err
}

// DeleteOfficialCustomMenuToDb Delete custom menu
func DeleteOfficialCustomMenuToDb(adminUserID int, appid string) error {
	// Begin transaction
	tx, err := msql.Begin(define.Postgres)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Delete menu
	_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Delete()
	if err != nil {
		return err
	}

	// Clear all related caches
	allMenus, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Select()
	if err != nil {
		return err
	}

	for _, menu := range allMenus {
		menuID := cast.ToInt(menu["id"])
		lib_redis.DelCacheData(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: menuID})
	}

	return nil
}
