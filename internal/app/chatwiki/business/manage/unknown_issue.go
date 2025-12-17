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
		filepath, err := common.ExportUnknownIssueSummary(list, export)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		c.FileAttachment(filepath, fmt.Sprintf(`未知问题总结导出%s.%s`, tool.Date(`Y-m-d-H-i-s`), export))
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
	if summary[`question`] != question { //更换了问题内容,重新转向量
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
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`未知问题总结开关未开启`)))
			return
		}
		embedding, err := common.GetVector2000(adminUserId, robot[`admin_user_id`], robot, msql.Params{}, msql.Params{},
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
