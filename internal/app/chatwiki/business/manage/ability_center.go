// Copyright © 2016- 2024 Sesame Network Technology all right reserved
package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
)

type SaveAbilityReq struct {
	AbilityType  string `form:"ability_type" json:"ability_type" binding:"required"`
	SwitchStatus int    `form:"switch_status" json:"switch_status"`
}

// SaveUserAbility 保存用户功能开关状态
func SaveUserAbility(c *gin.Context) {
	var (
		err error
		req SaveAbilityReq
	)

	// 获取参数
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 检查开关状态参数是否合法 (0:关闭, 1:开启)
	if req.SwitchStatus != define.SwitchOff && req.SwitchStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 保存用户功能开关状态
	err = common.SaveUserAbility(adminUserId, req.AbilityType, req.SwitchStatus)
	if err != nil {
		logs.Error("SaveUserAbility error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetAbilityList 获取用户功能列表
func GetAbilityList(c *gin.Context) {
	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取功能列表
	list, err := common.GetAbilityList(adminUserId)
	if err != nil {
		logs.Error("GetAbilityList error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, list)
}
