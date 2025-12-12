// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func TriggerList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	list, err := work_flow.TriggerList(adminUserId, common.GetLang(c))
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, list)
}

func TriggerSwitch(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	switchStatus := cast.ToInt(c.PostForm("switch_status"))
	if !tool.InArrayInt(switchStatus, []int{1, 0}) {
		common.FmtError(c, `param_err`, `switch_status`)
		return
	}
	id := cast.ToInt(c.PostForm("id"))
	if id == 0 {
		common.FmtError(c, `param_err`, `id`)
		return
	}
	info, err := msql.Model(`trigger_config`, define.Postgres).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	_, err = msql.Model(`trigger_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).Update(msql.Datas{
		"switch_status": switchStatus,
		"update_time":   time.Now().Unix(),
	})
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	lib_redis.DelCacheData(define.Redis, common.TriggerConfigCacheBuildHandler{
		AdminUserId: adminUserId,
		TriggerType: info[`trigger_type`],
	})
	common.FmtOk(c, nil)
}
