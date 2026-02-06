// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	wechatCommon "chatwiki/internal/pkg/wechat/common"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type RobotSmartMenuCacheBuildHandler struct {
	ID int
}

func (h *RobotSmartMenuCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`robot_smart_menu_%d`, h.ID)
}
func (h *RobotSmartMenuCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where(`id`, cast.ToString(h.ID)).Find()
	if err != nil {
		return nil, err
	}
	// Convert
	if len(data) > 0 {
		//
		var smartMenuContent = make([]lib_define.SmartMenuContent, 0)
		_ = tool.JsonDecodeUseNumber(data[`menu_content`], &smartMenuContent)
		result := lib_define.SmartMenu{
			ID:              cast.ToInt(data[`id`]),
			AdminUserID:     cast.ToInt(data[`admin_user_id`]),
			RobotID:         cast.ToInt(data[`robot_id`]),
			MenuTitle:       data[`menu_title`],
			MenuDescription: data[`menu_description`],
			MenuContent:     smartMenuContent,
			CreateTime:      cast.ToInt(data[`create_time`]),
			UpdateTime:      cast.ToInt(data[`update_time`]),
		}
		return result, err
	} else {
		return nil, nil
	}
}

// GetSmartMenuInfo Get smart menu information
func GetSmartMenuInfo(id int) (lib_define.SmartMenu, error) {
	result := lib_define.SmartMenu{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// FormatReplyListToDb Format reply list
// Note that this formatting may modify the array length, so check if the length is 0 after calling
func FormatReplyListToDb(replyList []ReplyContent, sendSource string) []ReplyContent {
	if len(replyList) == 0 {
		return replyList
	}
	newReplyList := make([]ReplyContent, 0)
	for i := range replyList {
		// Mark source
		replyList[i].SendSource = sendSource
		replyType := replyList[i].ReplyType
		if replyType == `` && replyList[i].Type != `` {
			replyType = replyList[i].Type
		}
		// Smart menu
		if replyType == ReplyTypeSmartMenu {
			// Get smart menu content
			smartMenu, err := GetSmartMenuInfo(cast.ToInt(replyList[i].SmartMenuId))
			if err == nil && smartMenu.ID > 0 {
				// Replace date
				smartMenu.MenuDescription = wechatCommon.ReplaceDate(smartMenu.MenuDescription)
				replyList[i].SmartMenu = smartMenu
			} else {
				continue
			}
		}
		newReplyList = append(newReplyList, replyList[i])
	}
	return newReplyList
}

// SaveSmartMenu Save smart menu (create or update)
func SaveSmartMenu(id, adminUserID, robotID int, menuTitle, menuDescription string, menuContent []lib_define.SmartMenuContent) (int64, error) {
	menuContentJson, _ := tool.JsonEncode(menuContent)

	data := msql.Datas{
		"admin_user_id":    adminUserID,
		"robot_id":         robotID,
		"menu_title":       menuTitle,
		"menu_description": menuDescription,
		"menu_content":     menuContentJson,
		"update_time":      tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// Create new record
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Insert(data, "id")
	} else {
		// Update existing record
		_, err = msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id})
	}
	return newId, err
}

// DeleteSmartMenu Delete smart menu
func DeleteSmartMenu(id, robotID int) error {
	_, err := msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Delete()
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id})
	}
	return err
}

// GetSmartMenu Get single smart menu
func GetSmartMenu(id int, robotID int) (msql.Params, error) {
	return msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Find()
}

// GetSmartMenuListWithFilter Get smart menu list (with filters and pagination)
func GetSmartMenuListWithFilter(robotID int, menuTitle string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("robot_id", cast.ToString(robotID))

	// Fuzzy query by menu title
	if menuTitle != "" {
		model.Where("menu_title LIKE ?", "%"+menuTitle+"%")
	}

	// Add pagination
	list, total, err := model.Order("id DESC").Paginate(page, size)
	return list, total, err
}
