// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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

// 1发送消息 2跳转网页 3跳转小程序 4人工客服 5推送事件
const (
	OfficialCustomMenuActTypeSendMessage     = 1
	OfficialCustomMenuActTypeJumpURL         = 2
	OfficialCustomMenuActTypeJumpMiniProgram = 3
	OfficialCustomMenuActTypeCustomerService = 4
	OfficialCustomMenuActTypePushEvent       = 5
)

const OfficialCustomMenuMenuLevelOne = 1
const OfficialCustomMenuMenuLevelTwo = 2

// OfficialCustomMenu 聊天机器人智能菜单
type OfficialCustomMenu struct {
	ID            int                        `json:"id" form:"id"`                             // 自增ID
	AdminUserID   int                        `json:"admin_user_id" form:"admin_user_id"`       // 管理员用户ID
	AppID         string                     `json:"appid" form:"appid"`                       // 公众号appid
	SeqID         int                        `json:"seq_id" form:"seq_id"`                     // 排序id
	MenuName      string                     `json:"menu_name" form:"menu_name"`               // 菜单名称
	MenuLevel     int                        `json:"menu_level" form:"menu_level"`             // 菜单层级 1根菜单 2二级菜单
	ParentMenuID  int                        `json:"parent_menu_id" form:"parent_menu_id"`     // 上级节点id
	ChooseActItem int                        `json:"choose_act_item" form:"choose_act_item"`   // 默认选中栏目 0有子节点时无此功能 1发送消息 2跳转网页 3跳转小程序 4人工客服 5推送事件
	ActParams     OfficialCustomMenuActParam `json:"act_params" form:"act_params"`             // 配置json串
	OperUserID    int                        `json:"oper_user_id" form:"oper_user_id"`         // 操作人ID
	TemplateID    int                        `json:"template_id,omitempty" form:"template_id"` // 模版id
	BatchID       int                        `json:"batch_id,omitempty" form:"batch_id"`       // 批量id
	SubMenuList   []OfficialCustomMenu       `json:"sub_menu_list,omitempty" form:"sub_menu_list"`
	CreateTime    int                        `json:"create_time,omitempty"` // 创建时间
	UpdateTime    int                        `json:"update_time,omitempty"` // 更新时间
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

// OfficialCustomMenuCacheBuildHandler 菜单缓存处理器
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

	// 转换
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

// GetOfficialCustomMenuInfo 获取自定义菜单信息
func GetOfficialCustomMenuInfo(id int) (OfficialCustomMenu, error) {
	result := OfficialCustomMenu{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// SaveOfficialCustomMenu 保存自定义菜单（创建或更新）
// 更新菜单而不是删除重建，保持菜单ID不变以维持与微信的关联
func SaveOfficialCustomMenu(adminUserID int, appid string, operUserID int, menuList []OfficialCustomMenu, sendWx bool) (int, error) {
	// 开始事务
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

	// 获取当前数据库中存在的菜单
	existingMenus, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Select()
	if err != nil {
		return 0, err
	}

	// 创建现有菜单ID映射
	existingMenuMap := make(map[int]msql.Params)
	for _, menu := range existingMenus {
		menuID := cast.ToInt(menu["id"])
		existingMenuMap[menuID] = menu
	}

	// 解析主菜单
	var parentMenus []OfficialCustomMenu

	for _, menu := range menuList {
		if menu.MenuLevel == OfficialCustomMenuMenuLevelOne {
			parentMenus = append(parentMenus, menu)
		}
	}

	// 用于跟踪已处理的菜单ID，以便删除不再需要的菜单
	processedMenuIDs := make(map[int]bool)

	// 更新或创建主菜单及对应的子菜单
	parentIDMap := make(map[int]int) // oldID -> newID (对于新建菜单)
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
			// 更新现有菜单
			_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
				Where("id", cast.ToString(menu.ID)).
				Update(data)
			if err != nil {
				return 0, err
			}
			newId = int64(menu.ID)
			processedMenuIDs[menu.ID] = true
		} else {
			// 创建新菜单
			data["parent_menu_id"] = 0
			data["create_time"] = tool.Time2Int()
			newId, err = msql.Model(`func_official_custom_menu`, define.Postgres).Insert(data, "id")
			if err != nil {
				return 0, err
			}
			// 记录新创建的菜单ID映射
			if menu.ID > 0 {
				parentIDMap[menu.ID] = int(newId)
			}
			processedMenuIDs[int(newId)] = true
		}

		// 更新父菜单ID映射，确保即使原ID不存在也能正确映射
		if menu.ID > 0 {
			parentIDMap[menu.ID] = int(newId)
		}
		parentMenus[i].ID = int(newId)

		// 处理该主菜单下的子菜单
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
				// 更新现有的子菜单
				_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
					Where("id", cast.ToString(subMenu.ID)).
					Update(subData)
				if err != nil {
					return 0, err
				}
				subNewId = int64(subMenu.ID)
				processedMenuIDs[subMenu.ID] = true
			} else {
				// 创建新的子菜单
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

	// 删除未被处理的旧菜单（即已被删除的菜单）
	for menuID := range existingMenuMap {
		if !processedMenuIDs[menuID] {
			_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
				Where("id", cast.ToString(menuID)).
				Delete()
			if err != nil {
				return 0, err
			}
			// 清除被删除菜单的缓存
			lib_redis.DelCacheData(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: menuID})
		}
	}

	// 清除所有相关缓存
	for menuID := range processedMenuIDs {
		lib_redis.DelCacheData(define.Redis, &OfficialCustomMenuCacheBuildHandler{ID: menuID})
	}

	// 记录菜单变更历史
	err = saveMenuHistory(adminUserID, appid, operUserID, oldMenus)
	if err != nil {
		return 0, err
	}
	if sendWx {
		//将本地菜单 发送到微信 菜单去
		return SendOfficialCustomMenuToWx(adminUserID, appid)
	}
	//将本地菜单 发送到微信 菜单去
	return 0, nil
}

// SendOfficialCustomMenuToWx 发送自定义菜单
func SendOfficialCustomMenuToWx(adminUserID int, appid string) (int, error) {
	appInfo, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserID)).Where(`app_id`, appid).Find()
	if err != nil {
		return 0, err
	}

	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}

	//获取菜单 转化成request.RequestMenuCreate

	localMenuList, err := GetOfficialCustomMenuList(adminUserID, appid)
	if err != nil {
		return 0, err
	}
	//转换成request.RequestMenuCreate
	menu := LocalMenuListConvertToRequestMenuCreate(localMenuList)
	if len(menu.Buttons) == 0 {
		return 0, errors.New(`菜单为空`)
	}
	//处理
	setMenuCode, err := app.SetMenu(menu)
	if err == nil {
		//保存菜单 开启状态
		_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
			Where(`id`, appInfo[`id`]).
			Where(`admin_user_id`, appInfo[`admin_user_id`]).
			Update(msql.Datas{
				"custom_menu_status": define.SwitchOn,
			})
	}

	return setMenuCode, err
}

// LocalMenuListConvertToRequestMenuCreate 本地菜单列表转换为 request.RequestMenuCreate
func LocalMenuListConvertToRequestMenuCreate(localMenuList []OfficialCustomMenu) request.RequestMenuCreate {
	var menu request.RequestMenuCreate
	buttons := make([]*request.Button, 0)
	for _, localMenu := range localMenuList {

		curButton := request.Button{
			Name: localMenu.MenuName,
		}
		// 不再需要循环 actParams，因为 ActParams 现在是单个对象而不是数组
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

		//处理子菜单
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

// SyncWxMenuToShow 同步微信菜单到展示
func SyncWxMenuToShow(adminUserID int, appid string, operUserID int) ([]OfficialCustomMenu, error) {
	appInfo, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserID)).Where(`app_id`, appid).Find()
	if err != nil {
		return nil, err
	}
	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}

	localMenuList := make([]OfficialCustomMenu, 0)
	WxMenu, err := app.GetMenu()
	if err != nil {
		//报错如果包含 menu no exist
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
	//转化为本地菜单
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
				//获取对应菜单的详情
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
						//获取对应菜单的详情
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
	// 检查是否使用默认菜单
	localMenuList = CheckUseDefaultMenu(adminUserID, appid, localMenuList)
	//转化为本地菜单返回给前端
	return localMenuList, err
}

// saveMenuHistory 保存菜单变更历史
func saveMenuHistory(adminUserID int, appid string, operUserID int, oldMenus []OfficialCustomMenu) error {
	// 获取当前数据库中的菜单（完整的一级和二级菜单结构）
	currentMenuList, err := GetOfficialCustomMenuList(adminUserID, appid)
	if err != nil {
		return err
	}

	// 清理当前菜单列表中的时间字段
	cleanCurrentMenuList := make([]OfficialCustomMenu, len(currentMenuList))
	for i, menu := range currentMenuList {
		cleanMenu := menu
		cleanMenu.CreateTime = 0
		cleanMenu.UpdateTime = 0

		// 清理子菜单的时间字段
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
	// 清理提交的菜单列表中的时间字段
	cleanHistoryMenuList := make([]OfficialCustomMenu, len(oldMenus))
	for i, menu := range historyMenuList {
		cleanMenu := menu
		cleanMenu.CreateTime = 0
		cleanMenu.UpdateTime = 0

		// 清理子菜单的时间字段
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

	// 使用 DeepEqual 比较两个结构体切片
	if !reflect.DeepEqual(cleanHistoryMenuList, cleanCurrentMenuList) {
		// 重新编码当前菜单（包含时间字段）用于存储
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

// GetOfficialCustomMenuList 获取自定义菜单列表
func GetOfficialCustomMenuList(adminUserID int, appid string) ([]OfficialCustomMenu, error) {
	// 获取所有菜单
	menus, err := msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Order("menu_level ASC, seq_id ASC, id ASC").
		Select()
	if err != nil {
		return nil, err
	}

	// 分别存储一级菜单和二级菜单
	parentMenus := make([]OfficialCustomMenu, 0)
	subMenuMap := make(map[int][]OfficialCustomMenu) // parent_menu_id -> []SubMenu

	// 处理所有菜单项
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
			// 将二级菜单按parent_menu_id分组
			subMenuMap[menu.ParentMenuID] = append(subMenuMap[menu.ParentMenuID], menu)
		}
	}

	// 将二级菜单附加到对应的一级菜单
	for i := range parentMenus {
		if subMenus, exists := subMenuMap[parentMenus[i].ID]; exists {
			parentMenus[i].SubMenuList = subMenus
		}
	}
	// 检查是否使用默认菜单
	parentMenus = CheckUseDefaultMenu(adminUserID, appid, parentMenus)

	return parentMenus, nil
}

// CheckUseDefaultMenu 检查是否使用默认菜单
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
			MenuName:      `主菜单`,
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

// DeleteAllOfficialCustomMenuToWx 删除所有公众号自定义菜单
func DeleteAllOfficialCustomMenuToWx(adminUserId int, appid string) string {
	//获取所有公众号
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
			errStr += wxName + `报错：` + err.Error() + `\n`
		}
	}
	return errStr
}

func WxDeleteMenu(appInfo msql.Params) error {
	//捕获报错
	app := &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}
	_, err := app.DeleteMenu()
	if err != nil {
		return err
	}
	//更新公众号菜单开关状态
	_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`id`, appInfo[`id`]).
		Where(`admin_user_id`, appInfo[`admin_user_id`]).
		Update(msql.Datas{
			"custom_menu_status": define.SwitchOff,
		})
	return err
}

// DeleteOfficialCustomMenuToDb 删除自定义菜单
func DeleteOfficialCustomMenuToDb(adminUserID int, appid string) error {
	// 开始事务
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

	// 删除菜单
	_, err = msql.Model(`func_official_custom_menu`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserID)).
		Where("appid", appid).
		Delete()
	if err != nil {
		return err
	}

	// 清除所有相关缓存
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
