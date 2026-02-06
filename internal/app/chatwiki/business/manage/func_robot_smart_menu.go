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

// SmartMenuRequest General request struct
type SmartMenuRequest struct {
	ID              int    `form:"id" json:"id"`
	RobotID         int    `form:"robot_id" json:"robot_id" binding:"required"`
	MenuTitle       string `form:"menu_title" json:"menu_title"`
	MenuDescription string `form:"menu_description" json:"menu_description"`
	MenuContent     string `form:"menu_content" json:"menu_content"`
}

// SmartMenuListFilterRequest List filter request struct
type SmartMenuListFilterRequest struct {
	RobotID   int    `json:"robot_id" form:"robot_id" binding:"required"`
	MenuTitle string `json:"menu_title" form:"menu_title"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// SaveSmartMenu Save smart menu (create or update)
func SaveSmartMenu(c *gin.Context) {
	var req SmartMenuRequest

	// Get params
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Parse MenuContent
	var menuContent []lib_define.SmartMenuContent
	if req.MenuContent != "" {
		if err := tool.JsonDecodeUseNumber(req.MenuContent, &menuContent); err != nil {
			common.FmtError(c, `param_err`, "invalid menu_content format")
			return
		}
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the robot exists and belongs to current user
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// If updating, check whether the menu exists and belongs to current user and robot
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

	// Save menu (create or update)
	id, err := common.SaveSmartMenu(
		req.ID,
		adminUserId,
		req.RobotID,
		req.MenuTitle,
		req.MenuDescription,
		menuContent, // use parsed value
	)

	if err != nil {
		logs.Error("SaveSmartMenu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteSmartMenu Delete smart menu
func DeleteSmartMenu(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	robotID := cast.ToInt(c.PostForm("robot_id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the robot exists and belongs to current user
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// Check whether the menu exists and belongs to current user and robot
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

	// Delete menu
	err = common.DeleteSmartMenu(id, robotID)
	if err != nil {
		logs.Error("DeleteSmartMenu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetSmartMenu Get a single smart menu
func GetSmartMenu(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	robotID := cast.ToInt(c.Query("robot_id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get menu info
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

	// Check permissions
	if cast.ToInt(menuInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// Build response
	var menuContent []lib_define.SmartMenuContent

	// Parse JSON
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

// GetSmartMenuList Get smart menu list
func GetSmartMenuList(c *gin.Context) {
	var req SmartMenuListFilterRequest

	// Get params
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the robot exists and belongs to current user
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// Get menu list
	list, total, err := common.GetSmartMenuListWithFilter(req.RobotID, req.MenuTitle, req.Page, req.Size)
	if err != nil {
		logs.Error("GetSmartMenuListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// Build response
	var result []map[string]interface{}
	for _, item := range list {
		var menuContent []lib_define.SmartMenuContent

		// Parse JSON
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

	// Return paginated data
	response := map[string]interface{}{
		"list":  result,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}
