// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_define"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// SmartMenuRequest 通用请求结构
type SmartMenuRequest struct {
	ID              int    `form:"id" json:"id"`
	RobotID         int    `form:"robot_id" json:"robot_id" binding:"required"`
	MenuTitle       string `form:"menu_title" json:"menu_title"`
	MenuDescription string `form:"menu_description" json:"menu_description"`
	MenuContent     string `form:"menu_content" json:"menu_content"`
}

// SmartMenuListFilterRequest 列表过滤请求结构
type SmartMenuListFilterRequest struct {
	RobotID   int    `json:"robot_id" form:"robot_id" binding:"required"`
	MenuTitle string `json:"menu_title" form:"menu_title"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// SaveSmartMenu 保存智能菜单（创建或更新）
func SaveSmartMenu(c *gin.Context) {
	var req SmartMenuRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 解析 MenuContent
	var menuContent []lib_define.SmartMenuContent
	if req.MenuContent != "" {
		if err := tool.JsonDecodeUseNumber(req.MenuContent, &menuContent); err != nil {
			common.FmtError(c, `param_err`, "invalid menu_content format")
			return
		}
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// 如果是更新操作，检查规则是否存在且属于当前用户和机器人
	if req.ID > 0 {
		menuInfo, err := common.GetSmartMenu(req.ID, req.RobotID)
		if err != nil {
			logs.Error("Get smart menu error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}

		if len(menuInfo) == 0 {
			common.FmtError(c, `menu_not_exist`)
			return
		}

		if cast.ToInt(menuInfo["admin_user_id"]) != adminUserId || cast.ToInt(menuInfo["robot_id"]) != req.RobotID {
			common.FmtError(c, `auth_no_permission`)
			return
		}
	}

	// 保存规则（创建或更新）
	id, err := common.SaveSmartMenu(
		req.ID,
		adminUserId,
		req.RobotID,
		req.MenuTitle,
		req.MenuDescription,
		menuContent, // 使用解析后的值
	)

	if err != nil {
		logs.Error("SaveSmartMenu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteSmartMenu 删除智能菜单
func DeleteSmartMenu(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	robotID := cast.ToInt(c.PostForm("robot_id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// 检查菜单是否存在且属于当前用户和机器人
	menuInfo, err := common.GetSmartMenu(id, robotID)
	if err != nil {
		logs.Error("Get smart menu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(menuInfo) == 0 {
		common.FmtError(c, `menu_not_exist`)
		return
	}

	if cast.ToInt(menuInfo["admin_user_id"]) != adminUserId || cast.ToInt(menuInfo["robot_id"]) != robotID {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 删除菜单
	err = common.DeleteSmartMenu(id, robotID)
	if err != nil {
		logs.Error("DeleteSmartMenu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetSmartMenu 获取单个智能菜单
func GetSmartMenu(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	robotID := cast.ToInt(c.Query("robot_id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取菜单信息
	menuInfo, err := common.GetSmartMenu(id, robotID)
	if err != nil {
		logs.Error("Get smart menu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(menuInfo) == 0 {
		common.FmtError(c, `menu_not_exist`)
		return
	}

	// 检查权限
	if cast.ToInt(menuInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 处理返回数据
	var menuContent []lib_define.SmartMenuContent

	// 解析JSON数据
	if menuInfo["menu_content"] != "" {
		_ = tool.JsonDecodeUseNumber(menuInfo["menu_content"], &menuContent)
	}

	result := map[string]interface{}{
		"id":               menuInfo["id"],
		"robot_id":         menuInfo["robot_id"],
		"menu_title":       menuInfo["menu_title"],
		"menu_description": menuInfo["menu_description"],
		"menu_content":     menuContent,
		"create_time":      menuInfo["create_time"],
		"update_time":      menuInfo["update_time"],
	}

	common.FmtOk(c, result)
}

// GetSmartMenuList 获取智能菜单列表
func GetSmartMenuList(c *gin.Context) {
	var req SmartMenuListFilterRequest

	// 获取参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// 获取菜单列表
	list, total, err := common.GetSmartMenuListWithFilter(req.RobotID, req.MenuTitle, req.Page, req.Size)
	if err != nil {
		logs.Error("GetSmartMenuListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// 处理返回数据
	var result []map[string]interface{}
	for _, item := range list {
		var menuContent []lib_define.SmartMenuContent

		// 解析JSON数据
		if item["menu_content"] != "" {
			_ = tool.JsonDecodeUseNumber(item["menu_content"], &menuContent)
		}

		result = append(result, map[string]interface{}{
			"id":               item["id"],
			"robot_id":         item["robot_id"],
			"menu_title":       item["menu_title"],
			"menu_description": item["menu_description"],
			"menu_content":     menuContent,
			"create_time":      item["create_time"],
			"update_time":      item["update_time"],
		})
	}

	// 返回分页数据
	response := map[string]interface{}{
		"list":  result,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}
