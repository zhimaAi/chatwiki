// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func UnknownIssueStats(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	startDay := cast.ToUint(c.DefaultQuery(`start_day`, cast.ToString(tool.GetYmdBeforeDay(6))))
	endDay := cast.ToUint(c.DefaultQuery(`end_day`, tool.Date(`Ymd`)))
	startDay, endDay = tool.GetYmd(tool.GetTimestamp(startDay)), tool.GetYmd(tool.GetTimestamp(endDay))
	page := max(1, cast.ToInt(c.DefaultQuery(`page`, `1`)))
	size := max(1, cast.ToInt(c.DefaultQuery(`size`, `20`)))
	if robotId <= 0 || startDay < 20250101 || startDay > endDay {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`robot_id`, cast.ToString(robotId)).
		Where(`stats_day`, `between`, fmt.Sprintf(`%d,%d`, startDay, endDay))
	list, total, err := m.Order(`id desc`).
		Field(`question,stats_day,trigger_total`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for i, item := range list {
		list[i][`show_date`] = tool.Date(`Y-m-d`, tool.GetTimestamp(cast.ToUint(item[`stats_day`])))
	}
	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetUnknownIssueSummary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	startDay := cast.ToUint(c.Query(`start_day`))
	endDay := cast.ToUint(c.Query(`end_day`))
	page := max(1, cast.ToInt(c.DefaultQuery(`page`, `1`)))
	size := max(1, cast.ToInt(c.DefaultQuery(`size`, `20`)))
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_unknown_issue_summary`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`robot_id`, cast.ToString(robotId))
	if startDay > 0 && endDay > 0 {
		startDay, endDay = tool.GetYmd(tool.GetTimestamp(startDay)), tool.GetYmd(tool.GetTimestamp(endDay))
		if startDay < 20250101 || startDay > endDay {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `trigger_day`))))
			return
		}
		m.Where(`trigger_day`, `between`, fmt.Sprintf(`%d,%d`, startDay, endDay))
	}
	export := strings.TrimSpace(c.Query(`export`))
	if tool.InArrayString(export, []string{`docx`, `xlsx`}) {
		list, err := m.Field(`question,unknown_list,answer,images`).Order(`id desc`).Select()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		filepath, err := common.ExportUnknownIssueSummary(common.GetLang(c), list, export)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		c.FileAttachment(filepath, i18n.Show(common.GetLang(c), `unknown_issue_export_filename`, tool.Date(`Y-m-d-H-i-s`), export))
		return
	}
	list, total, err := m.Field(`id,question,unknown_total,unknown_list,trigger_day,answer,images,to_library_id,to_library_name`).
		Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for i, item := range list {
		list[i][`show_date`] = tool.Date(`Y-m-d`, tool.GetTimestamp(cast.ToUint(item[`trigger_day`])))
		if cast.ToUint(item[`unknown_total`]) == 0 {
			list[i][`unknown_list`], list[i][`unknown_total`] = tool.JsonEncodeNoError([]string{item[`question`]}), `1`
		}
	}
	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func UnknownIssueSummaryAnswer(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	question := strings.TrimSpace(c.PostForm(`question`))
	unknownList := common.DisposeStringList(c.PostForm(`unknown_list`))
	answer := strings.TrimSpace(c.PostForm(`answer`))
	images := c.PostFormArray(`images`)
	//check required
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if len(question) < 1 || len(question) > common.MaxContent {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
		return
	}
	if len(answer) < 1 || len(answer) > common.MaxContent {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
		return
	}
	jsonImages, err := common.CheckLibraryImage(images)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `images`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_unknown_issue_summary`, define.Postgres)
	summary, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Field(`question,robot_id`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(summary) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{
		`question`:      question,
		`unknown_list`:  tool.JsonEncodeNoError(unknownList),
		`unknown_total`: len(unknownList),
		`answer`:        answer,
		`images`:        jsonImages,
		`update_time`:   tool.Time2Int(),
	}
	if summary[`question`] != question { //question content changed; regenerate embedding
		robot, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, summary[`robot_id`]).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(robot) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		if cast.ToUint(robot[`unknown_summary_status`]) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `unknown_issue_summary_switch_disabled`))))
			return
		}
		embedding, err := common.GetVector2000(common.GetLang(c), adminUserId, robot[`admin_user_id`], robot, msql.Params{}, msql.Params{},
			cast.ToInt(robot[`unknown_summary_model_config_id`]), robot[`unknown_summary_use_model`], question)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		data[`embedding`] = embedding
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UnknownIssueSummaryImport(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	toLibraryId := cast.ToInt(c.PostForm(`to_library_id`))
	//check required
	if id <= 0 || toLibraryId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_unknown_issue_summary`, define.Postgres)
	summary, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(summary) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	library, err := common.GetLibraryInfo(toLibraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(library) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{
		`to_library_id`:   toLibraryId,
		`to_library_name`: library[`library_name`],
		`update_time`:     tool.Time2Int(),
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetUnknownIssueChatContext(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	statsId := cast.ToInt(c.Query(`stats_id`))
	robotId := cast.ToInt(c.Query(`robot_id`))
	triggerDay := cast.ToInt(c.Query(`trigger_day`))
	question := strings.TrimSpace(c.Query(`question`))
	window := cast.ToInt(c.DefaultQuery(`window`, `15`))
	if window <= 0 {
		window = 15
	}
	if window > 50 {
		window = 50
	}
	if statsId <= 0 && (robotId <= 0 || triggerDay <= 0 || len(question) == 0) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId))
	if statsId > 0 {
		m.Where(`id`, cast.ToString(statsId))
	} else {
		m.Where(`robot_id`, cast.ToString(robotId)).Where(`stats_day`, cast.ToString(triggerDay)).Where(`question`, question)
	}
	stats, err := m.Field(`id,robot_id,stats_day,question,sample_openid,sample_rel_user_id,sample_dialogue_id,sample_session_id,sample_message_id,last_dialogue_id,last_session_id,last_message_id,last_trigger_time`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(stats) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	robotId = cast.ToInt(stats[`robot_id`])
	dialogueId := cast.ToInt(stats[`last_dialogue_id`])
	sessionId := cast.ToInt(stats[`last_session_id`])
	messageId := cast.ToInt(stats[`last_message_id`])
	openid := stats[`sample_openid`]
	relUserId := cast.ToInt(stats[`sample_rel_user_id`])
	usedFallback := false
	if dialogueId <= 0 || messageId <= 0 {
		startTs := tool.GetTimestamp(cast.ToUint(stats[`stats_day`]))
		endTs := startTs + 86400
		lastMessage, queryErr := msql.Model(`chat_ai_message`, define.Postgres).Alias(`m`).
			Join(`chat_ai_session s`, `m.session_id=s.id`, `left`).
			Where(`m.admin_user_id`, cast.ToString(adminUserId)).
			Where(`m.robot_id`, cast.ToString(robotId)).
			Where(`m.is_customer`, `1`).
			Where(`m.content`, stats[`question`]).
			Where(`m.create_time`, `between`, fmt.Sprintf(`%d,%d`, startTs, endTs)).
			Order(`m.id desc`).
			Field(`m.id,m.dialogue_id,m.session_id,m.openid,m.create_time,s.rel_user_id`).
			Find()
		if queryErr != nil {
			logs.Error(queryErr.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(lastMessage) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		usedFallback = true
		dialogueId = cast.ToInt(lastMessage[`dialogue_id`])
		sessionId = cast.ToInt(lastMessage[`session_id`])
		messageId = cast.ToInt(lastMessage[`id`])
		openid = lastMessage[`openid`]
		relUserId = cast.ToInt(lastMessage[`rel_user_id`])
		upData := msql.Datas{
			`last_dialogue_id`:  dialogueId,
			`last_session_id`:   sessionId,
			`last_message_id`:   messageId,
			`last_trigger_time`: cast.ToInt(lastMessage[`create_time`]),
			`update_time`:       tool.Time2Int(),
		}
		if cast.ToInt(stats[`sample_message_id`]) <= 0 {
			upData[`sample_openid`] = openid
			upData[`sample_rel_user_id`] = relUserId
			upData[`sample_dialogue_id`] = dialogueId
			upData[`sample_session_id`] = sessionId
			upData[`sample_message_id`] = messageId
		}
		if _, queryErr = msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres).
			Where(`id`, stats[`id`]).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(upData); queryErr != nil {
			logs.Error(queryErr.Error())
		}
	}
	anchorMessage, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`id`, cast.ToString(messageId)).
		Field(`id,dialogue_id`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(anchorMessage) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if dialogueId <= 0 {
		dialogueId = cast.ToInt(anchorMessage[`dialogue_id`])
	}
	fields := `id,admin_user_id,robot_id,openid,dialogue_id,session_id,is_customer,msg_type,content,nickname,name,avatar,received_message_type,media_id_to_oss_url,menu_json,quote_file,reply_content_list,create_time,update_time`
	beforeList, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`id`, `<`, cast.ToString(messageId)).
		Order(`id desc`).Limit(window).Field(fields).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for i, j := 0, len(beforeList)-1; i < j; i, j = i+1, j-1 {
		beforeList[i], beforeList[j] = beforeList[j], beforeList[i]
	}
	anchorList, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`id`, cast.ToString(messageId)).
		Field(fields).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	afterList, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`id`, `>`, cast.ToString(messageId)).
		Order(`id`).Limit(window).Field(fields).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list := append(beforeList, anchorList...)
	list = append(list, afterList...)
	robotInfo, _ := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,robot_name,robot_avatar,robot_key`).
		Find()
	customerInfo := msql.Params{}
	if len(openid) > 0 {
		customerInfo, _ = common.GetCustomerInfo(openid, adminUserId)
	}
	data := map[string]any{
		`stats_id`:          stats[`id`],
		`question`:          stats[`question`],
		`trigger_day`:       stats[`stats_day`],
		`dialogue_id`:       dialogueId,
		`session_id`:        sessionId,
		`openid`:            openid,
		`rel_user_id`:       relUserId,
		`anchor_message_id`: messageId,
		`window`:            window,
		`used_fallback`:     usedFallback,
		`robot`:             robotInfo,
		`customer`:          customerInfo,
		`list`:              list,
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}
