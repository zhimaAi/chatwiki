// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

const (
	systemModeAIService = "ai_service"
	systemModeAgent     = "agent"
	systemModeAdvanced  = "advanced"
)

type SaveSystemModeReq struct {
	SystemMode string `form:"system_mode" json:"system_mode" binding:"required"`
}

func isValidSystemMode(mode string) bool {
	switch mode {
	case systemModeAIService, systemModeAgent, systemModeAdvanced:
		return true
	default:
		return false
	}
}

func GetSystemMode(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}

	company, err := msql.Model(define.TableCompany, define.Postgres).
		Where(`parent_id`, cast.ToString(adminUserID)).
		Order(`id asc`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(company) == 0 {
		common.FmtError(c, `company_not_exist`)
		return
	}

	mode := company[`system_mode`]
	if !isValidSystemMode(mode) {
		mode = systemModeAdvanced
	}
	selected := 1
	if company[`system_mode_selected`] == `0` {
		selected = 0
	}
	common.FmtOk(c, map[string]any{`system_mode`: mode, `system_mode_selected`: selected})
}

func SaveSystemMode(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	loginUserID := getLoginUserId(c)
	if loginUserID != adminUserID && !common.CheckPermission(loginUserID, define.CompanyManage) {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	var req SaveSystemModeReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_invalid`, `system_mode`)
		return
	}
	mode := strings.TrimSpace(req.SystemMode)
	if !isValidSystemMode(mode) {
		common.FmtError(c, `param_invalid`, `system_mode`)
		return
	}

	affected, err := msql.Model(define.TableCompany, define.Postgres).
		Where(`parent_id`, cast.ToString(adminUserID)).
		Update(msql.Datas{
			`system_mode`:          mode,
			`system_mode_selected`: 1,
			`update_time`:          time.Now().Unix(),
		})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if affected == 0 {
		common.FmtError(c, `company_not_exist`)
		return
	}

	lib_redis.DelCacheData(define.Redis, &CompanyCacheBuildHandler{AdminUserId: adminUserID})
	lib_redis.DelCacheData(define.Redis, &CompanyCacheBuildHandler{AdminUserId: 0})
	common.FmtOk(c, map[string]any{`system_mode`: mode, `system_mode_selected`: 1})
}
