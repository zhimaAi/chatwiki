// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// CustomMenuRequest generic request structure
type CustomMenuRequest struct {
	Appid    string `form:"appid" json:"appid" binding:"required"`
	MenuJson string `form:"menu_json" json:"menu_json" binding:"required"`
}

// CustomMenuListFilterRequest list filter request structure
type CustomMenuListFilterRequest struct {
	Appid string `json:"appid" form:"appid" binding:"required"`
}

// SaveCustomMenu saves custom menu (create or update)
func SaveCustomMenu(c *gin.Context) {
	var req CustomMenuRequest

	// Get parameters
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Parse MenuJson
	var menuList []common.OfficialCustomMenu
	if req.MenuJson != "" {
		if err := tool.JsonDecodeUseNumber(req.MenuJson, &menuList); err != nil {
			common.FmtError(c, `param_err`, "invalid menu_json format")
			return
		}
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	// Feature switch check
	useAbility := common.CheckUseAbilityByAbilityType(adminUserId, common.OfficialAbilityCustomMenu)
	// Get operator user info
	operUserID := getLoginUserId(c) // Default use admin ID as operator ID

	// Save menu (create or update)
	_, err := common.SaveOfficialCustomMenu(
		adminUserId,
		req.Appid,
		operUserID,
		menuList, // Use parsed values
		useAbility,
	)

	if err != nil {
		logs.Error("SaveCustomMenu error: %s", err.Error())
		common.FmtError(c, `official_api_error`, err.Error())
		return
	}

	common.FmtOk(c, map[string]interface{}{"appid": req.Appid})
}

// GetCustomMenuList gets custom menu list
func GetCustomMenuList(c *gin.Context) {
	var req CustomMenuListFilterRequest

	// Get parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get menu list
	menuList, err := common.GetOfficialCustomMenuList(adminUserId, req.Appid)
	if err != nil {
		logs.Error("GetCustomMenuList error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{
		"list": menuList,
	})
}

// CloseWxMenu closes WeChat menu
func CloseWxMenu(c *gin.Context) {
	appid := cast.ToString(c.PostForm("appid"))
	if appid == `` {
		common.FmtError(c, `param_err`, "appid")
		return
	}
	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get menu list
	errorString := common.DeleteAllOfficialCustomMenuToWx(adminUserId, appid)
	common.FmtOk(c, map[string]interface{}{
		"errorString": errorString,
	})
}

// SyncWxMenuToShow syncs WeChat menu to display
func SyncWxMenuToShow(c *gin.Context) {
	appid := cast.ToString(c.Query("appid"))
	if appid == "" {
		common.FmtError(c, `param_err`, "appid is required")
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	// Get operator user info
	operUserID := getLoginUserId(c) // Default use admin ID as operator ID

	// Get menu list
	menuList, err := common.SyncWxMenuToShow(adminUserId, appid, operUserID)
	if err != nil {
		logs.Error("SyncWxMenuToShow error: %s", err.Error())
		common.FmtError(c, `official_api_error`, err.Error())
		return
	}

	common.FmtOk(c, map[string]interface{}{
		"list": menuList,
	})
}
