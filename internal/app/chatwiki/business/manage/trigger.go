// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
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

func GetTriggerOfficialMessage(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	list := []map[string]any{
		{
			`msg_type`: define.TriggerOfficialMessage,
			`desc`:     i18n.Show(common.GetLang(c), `trigger_gzh_message`),
			`fields`:   work_flow.GetMessage(common.GetLang(c)),
		},
		{
			`msg_type`: define.TriggerOfficialSubscribeUnScribe,
			`desc`:     i18n.Show(common.GetLang(c), `trigger_gzh_subscribe`),
			`fields`:   work_flow.GetSubscribeUnsubscribe(common.GetLang(c)),
		},
		{
			`msg_type`: define.TriggerOfficialQrCodeScan,
			`desc`:     i18n.Show(common.GetLang(c), `trigger_gzh_scan`),
			`fields`:   work_flow.GetQrcodeScan(common.GetLang(c)),
		},
		{
			`msg_type`: define.TriggerOfficialMenuClick,
			`desc`:     i18n.Show(common.GetLang(c), `trigger_gzh_click`),
			`fields`:   work_flow.GetMenuClick(common.GetLang(c)),
		},
	}
	robotKey := cast.ToString(c.Query(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	appList, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`app_type`, lib_define.AppOfficeAccount).
		Order(`id desc`).Field(`app_name,app_id,app_avatar,account_customer_type`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for idx, appInfo := range appList {
		accountIsVerify := lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`])
		appList[idx][`account_is_verify`] = cast.ToString(accountIsVerify)
	}
	common.FmtOk(c, map[string]any{
		`messages`: list,
		`apps`:     appList,
	})
	return
}
