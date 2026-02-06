// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

type SaveAbilityReq struct {
	AbilityType  string `form:"ability_type" json:"ability_type" binding:"required"`
	SwitchStatus int    `form:"switch_status" json:"switch_status"`
}

// SaveUserAbility saves the user's feature switch status
func SaveUserAbility(c *gin.Context) {
	var (
		err error
		req SaveAbilityReq
	)

	//Get params
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	//Validate switch status (0: off, 1: on)
	if req.SwitchStatus != define.SwitchOff && req.SwitchStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}

	//Get current user
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	//Save user's feature switch status
	_, err = common.SaveUserAbility(adminUserId, req.AbilityType, req.SwitchStatus)
	if err != nil {
		logs.Error("SaveUserAbility error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetAbilityList gets the user's feature list
func GetAbilityList(c *gin.Context) {
	//Get current user
	abilityType := cast.ToString(c.Query("ability_type"))
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	//Get feature list
	list, err := common.GetAbilityList(adminUserId, abilityType)
	if err != nil {
		logs.Error("GetAbilityList error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, list)
}

// GetSpecifyAbilityConfig gets the user's feature config
func GetSpecifyAbilityConfig(c *gin.Context) {
	//Get params
	abilityType := cast.ToString(c.Query("ability_type"))
	if abilityType == `` {
		common.FmtError(c, `param_lack`, `ability_type`)
		return
	}

	//Get current user
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	//Get robot ability list
	list, err := common.GetAbilityList(adminUserId, abilityType)
	if err != nil {
		logs.Error("GetAbilityList error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(list) > 0 {
		common.FmtOk(c, list[0])
	} else {
		logs.Error(`no specified robot ability`)
		common.FmtError(c, `sys_err`)
	}
}
