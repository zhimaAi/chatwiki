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
	"github.com/zhimaAi/go_tools/msql"
)

type SaveRobotAbilityReq struct {
	RobotId       int    `form:"robot_id" json:"robot_id" binding:"required"`
	AbilityType   string `form:"ability_type" json:"ability_type" binding:"required"`
	SwitchStatus  int    `form:"switch_status" json:"switch_status"`
	FixedMenu     int    `form:"fixed_menu" json:"fixed_menu"`
	AiReplyStatus int    `form:"ai_reply_status" json:"ai_reply_status"`
}

type SaveRobotAbilitySwitchStatusReq struct {
	RobotId      int    `form:"robot_id" json:"robot_id" binding:"required"`
	AbilityType  string `form:"ability_type" json:"ability_type" binding:"required"`
	SwitchStatus int    `form:"switch_status" json:"switch_status"`
}

type SaveRobotAbilityFixedMenuReq struct {
	RobotId     int    `form:"robot_id" json:"robot_id" binding:"required"`
	AbilityType string `form:"ability_type" json:"ability_type" binding:"required"`
	FixedMenu   int    `form:"fixed_menu" json:"fixed_menu"`
}

type SaveRobotAbilityAiReplyStatusReq struct {
	RobotId       int    `form:"robot_id" json:"robot_id" binding:"required"`
	AbilityType   string `form:"ability_type" json:"ability_type" binding:"required"`
	AiReplyStatus int    `form:"ai_reply_status" json:"ai_reply_status"`
}

// SaveRobotAbility 保存机器人功能开关状态
func SaveRobotAbility(c *gin.Context) {
	var (
		err error
		req SaveRobotAbilityReq
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

	// 检查固定菜单参数是否合法 (0:关闭, 1:开启)
	if req.FixedMenu != define.SwitchOff && req.FixedMenu != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `fixed_menu`)
		return
	}

	// 检查AI回复参数是否合法 (0:关闭, 1:开启)
	if req.AiReplyStatus != define.SwitchOff && req.AiReplyStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `ai_reply_status`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where("id", cast.ToString(req.RobotId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error("Get robot info error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(robotInfo) == 0 {
		common.FmtError(c, `robot_not_exist`)
		return
	}

	// 保存机器人功能开关状态
	err = common.SaveRobotAbility(adminUserId, req.RobotId, req.AbilityType, req.SwitchStatus, req.FixedMenu, req.AiReplyStatus)
	if err != nil {
		logs.Error("SaveRobotAbility error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// SaveRobotAbilitySwitchStatus 保存机器人功能开关状态
func SaveRobotAbilitySwitchStatus(c *gin.Context) {
	var (
		err error
		req SaveRobotAbilitySwitchStatusReq
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

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotId, adminUserId) {
		return
	}

	// 保存机器人功能开关状态
	err = common.SaveRobotAbilitySwitchStatus(adminUserId, req.RobotId, req.AbilityType, req.SwitchStatus)
	if err != nil {
		logs.Error("SaveRobotAbilitySwitchStatus error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

func checkRobotByAdminUserId(c *gin.Context, robotId, adminUserId int) bool {
	// 检查机器人是否存在且属于当前用户
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where("id", cast.ToString(robotId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error("Get robot info error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return true
	}

	if len(robotInfo) == 0 {
		common.FmtError(c, `robot_not_exist`)
		return true
	}
	return false
}

// SaveRobotAbilityFixedMenu 保存机器人功能固定菜单状态
func SaveRobotAbilityFixedMenu(c *gin.Context) {
	var (
		err error
		req SaveRobotAbilityFixedMenuReq
	)

	// 获取参数
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 检查固定菜单参数是否合法 (0:关闭, 1:开启)
	if req.FixedMenu != define.SwitchOff && req.FixedMenu != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `fixed_menu`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotId, adminUserId) {
		return
	}

	// 保存机器人功能固定菜单状态
	err = common.SaveRobotAbilityFixedMenu(adminUserId, req.RobotId, req.AbilityType, req.FixedMenu)
	if err != nil {
		logs.Error("SaveRobotAbilityFixedMenu error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// SaveRobotAbilityAiReplyStatus 保存机器人功能AI回复状态
func SaveRobotAbilityAiReplyStatus(c *gin.Context) {
	var (
		err error
		req SaveRobotAbilityAiReplyStatusReq
	)

	// 获取参数
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 检查AI回复参数是否合法 (0:关闭, 1:开启)
	if req.AiReplyStatus != define.SwitchOff && req.AiReplyStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `ai_reply_status`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotId, adminUserId) {
		return
	}

	// 保存机器人功能AI回复状态
	err = common.SaveRobotAbilityAiReplyStatus(adminUserId, req.RobotId, req.AbilityType, req.AiReplyStatus)
	if err != nil {
		logs.Error("SaveRobotAbilityAiReplyStatus error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetRobotAbilityList 根据robot_id获取机器人功能列表
func GetRobotAbilityList(c *gin.Context) {
	// 获取参数
	robotId := cast.ToInt(c.Query("robot_id"))
	if robotId <= 0 {
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
	if checkRobotByAdminUserId(c, robotId, adminUserId) {
		return
	}

	// 获取机器人功能列表
	list, err := common.GetRobotAbilityList(robotId, adminUserId, ``)
	if err != nil {
		logs.Error("GetRobotAbilityByRobotId error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, list)
}

// GetRobotSpecifyAbilityConfig 获取指定机器人功能配置
func GetRobotSpecifyAbilityConfig(c *gin.Context) {
	// 获取参数
	robotId := cast.ToInt(c.Query("robot_id"))
	abilityType := cast.ToString(c.Query("ability_type"))
	if robotId <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}
	if abilityType == `` {
		common.FmtError(c, `param_lack`, `ability_type`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, robotId, adminUserId) {
		return
	}

	// 获取机器人功能列表
	list, err := common.GetRobotAbilityList(robotId, adminUserId, abilityType)
	if err != nil {
		logs.Error("GetRobotAbilityByRobotId error: %s", err.Error())
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
