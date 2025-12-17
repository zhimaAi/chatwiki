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

// CustomMenuRequest 通用请求结构
type CustomMenuRequest struct {
	Appid    string `form:"appid" json:"appid" binding:"required"`
	MenuJson string `form:"menu_json" json:"menu_json" binding:"required"`
}

// CustomMenuListFilterRequest 列表过滤请求结构
type CustomMenuListFilterRequest struct {
	Appid string `json:"appid" form:"appid" binding:"required"`
}

// SaveCustomMenu 保存自定义菜单（创建或更新）
func SaveCustomMenu(c *gin.Context) {
	var req CustomMenuRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 解析 MenuJson
	var menuList []common.OfficialCustomMenu
	if req.MenuJson != "" {
		if err := tool.JsonDecodeUseNumber(req.MenuJson, &menuList); err != nil {
			common.FmtError(c, `param_err`, "invalid menu_json format")
			return
		}
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	//功能开关检测
	useAbility := common.CheckUseAbilityByAbilityType(adminUserId, common.OfficialAbilityCustomMenu)
	// 获取操作用户信息
	operUserID := getLoginUserId(c) // 默认使用管理员ID作为操作人ID

	// 保存菜单（创建或更新）
	_, err := common.SaveOfficialCustomMenu(
		adminUserId,
		req.Appid,
		operUserID,
		menuList, // 使用解析后的值
		useAbility,
	)

	if err != nil {
		logs.Error("SaveCustomMenu error: %s", err.Error())
		common.FmtError(c, `official_api_error`, err.Error())
		return
	}

	common.FmtOk(c, map[string]interface{}{"appid": req.Appid})
}

// GetCustomMenuList 获取自定义菜单列表
func GetCustomMenuList(c *gin.Context) {
	var req CustomMenuListFilterRequest

	// 获取参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取菜单列表
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

// CloseWxMenu 关闭微信菜单
func CloseWxMenu(c *gin.Context) {
	appid := cast.ToString(c.PostForm("appid"))
	if appid == `` {
		common.FmtError(c, `param_err`, "appid")
		return
	}
	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取菜单列表
	errorString := common.DeleteAllOfficialCustomMenuToWx(adminUserId, appid)
	common.FmtOk(c, map[string]interface{}{
		"errorString": errorString,
	})
}

// SyncWxMenuToShow 同步微信菜单到展示
func SyncWxMenuToShow(c *gin.Context) {
	appid := cast.ToString(c.Query("appid"))
	if appid == "" {
		common.FmtError(c, `param_err`, "appid is required")
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	// 获取操作用户信息
	operUserID := getLoginUserId(c) // 默认使用管理员ID作为操作人ID

	// 获取菜单列表
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
