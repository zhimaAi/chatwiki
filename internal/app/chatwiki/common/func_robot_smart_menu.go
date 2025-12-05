// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
	//转换
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

// GetSmartMenuInfo 获取智能菜单信息
func GetSmartMenuInfo(id int) (lib_define.SmartMenu, error) {
	result := lib_define.SmartMenu{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// FormatReplyListToDb 格式化修改回复列表
// 注意这个格式化可能会修改 数组长度，所以需要调用后判断是否长度为0
func FormatReplyListToDb(replyList []ReplyContent, sendSource string) []ReplyContent {
	if len(replyList) == 0 {
		return replyList
	}
	newReplyList := make([]ReplyContent, 0)
	for i := range replyList {
		//标记来源
		replyList[i].SendSource = sendSource
		replyType := replyList[i].ReplyType
		if replyType == `` && replyList[i].Type != `` {
			replyType = replyList[i].Type
		}
		//智能菜单
		if replyType == ReplyTypeSmartMenu {
			//获取智能菜单内容
			smartMenu, err := GetSmartMenuInfo(cast.ToInt(replyList[i].SmartMenuId))
			if err == nil && smartMenu.ID > 0 {
				//替换时间
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

// SaveSmartMenu 保存智能菜单（创建或更新）
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
		// 创建新记录
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Insert(data, "id")
	} else {
		// 更新现有记录
		_, err = msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id})
	}
	return newId, err
}

// DeleteSmartMenu 删除智能菜单
func DeleteSmartMenu(id, robotID int) error {
	_, err := msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Delete()
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSmartMenuCacheBuildHandler{ID: id})
	}
	return err
}

// GetSmartMenu 获取单个智能菜单
func GetSmartMenu(id int, robotID int) (msql.Params, error) {
	return msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Find()
}

// GetSmartMenuListWithFilter 获取智能菜单列表（带过滤条件和分页）
func GetSmartMenuListWithFilter(robotID int, menuTitle string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(`func_chat_robot_smart_menu`, define.Postgres).Where("robot_id", cast.ToString(robotID))

	// 根据菜单标题模糊查询
	if menuTitle != "" {
		model.Where("menu_title LIKE ?", "%"+menuTitle+"%")
	}

	// 添加分页
	list, total, err := model.Order("id DESC").Paginate(page, size)
	return list, total, err
}
