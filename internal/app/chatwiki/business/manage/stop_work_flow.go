// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func StopWorkFlow(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	logID := cast.ToInt64(c.PostForm(`log_id`))
	if logID <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `log_id`))))
		return
	}
	logInfo, err := msql.Model(`work_flow_logs`, define.Postgres).
		Where(`id`, cast.ToString(logID)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Field(`id,status`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(logInfo) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	status := cast.ToInt(logInfo[`status`])
	// Already stopped: return success idempotently.
	if status == define.WorkFlowStatusStopped {
		c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`stopped`: true}, nil))
		return
	}
	// Completed and other non-running states cannot be stopped.
	if status != define.WorkFlowStatusRunning {
		c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`stopped`: false}, nil))
		return
	}
	// The database is the source of truth: atomically transition running to stopped.
	affected, err := msql.Model(`work_flow_logs`, define.Postgres).
		Where(`id`, cast.ToString(logID)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`status`, cast.ToString(define.WorkFlowStatusRunning)).
		Update(msql.Datas{
			`status`:      define.WorkFlowStatusStopped,
			`update_time`: tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if affected == 0 {
		// Concurrent race: re-check whether another request already stopped it.
		curStatus, qErr := msql.Model(`work_flow_logs`, define.Postgres).
			Where(`id`, cast.ToString(logID)).
			Where(`admin_user_id`, cast.ToString(userId)).
			Value(`status`)
		if qErr != nil {
			logs.Error(qErr.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`stopped`: cast.ToInt(curStatus) == define.WorkFlowStatusStopped}, nil))
		return
	}
	// Remove paused snapshots so workflows waiting for user input cannot be resumed.
	work_flow.DelWorkFlowStorageByLogID(logID)
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`stopped`: true}, nil))
}
