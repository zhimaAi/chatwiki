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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetActiveModels(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	models, err := msql.Model(`llm_token_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Field(`distinct(model) as model`).
		Select()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	var result []string
	for _, item := range models {
		result = append(result, item[`model`])
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func StatToken(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	m := msql.Model(`llm_token_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))

	model := strings.TrimSpace(c.Query(`model`))
	if len(model) > 0 {
		m.Where(`model`, model)
	}
	t := strings.TrimSpace(c.Query(`type`))
	if len(t) > 0 {
		m.Where(`type`, t)
	}
	startDate := strings.TrimSpace(c.Query(`start_date`))
	if len(startDate) > 0 {
		m.Where(`date`, `>=`, startDate)
	}
	endDate := strings.TrimSpace(c.Query(`end_date`))
	if len(endDate) > 0 {
		m.Where(`date`, `<=`, endDate)
	}

	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	list, total, err := m.Field(`id,admin_user_id,corp,model,type,prompt_token,completion_token,to_char(date, 'YYYY-MM-DD') AS date`).Order(`date desc, corp asc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func StatTokenApp(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	m := msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))

	tokenAppType := strings.TrimSpace(c.Query(`token_app_type`))
	if tokenAppType != `` {
		m.Where(`token_app_type`, tokenAppType)
	} else {
		m.Where(`token_app_type`, `in`, strings.Join(define.GetTokenAppTypes(), `,`))
	}
	robotId := strings.TrimSpace(c.Query(`robot_id`))
	if cast.ToInt(robotId) > 0 {
		m.Where(`robot_id`, robotId)
	}
	page := max(1, cast.ToInt(c.DefaultQuery(`page`, `1`)))
	size := max(1, cast.ToInt(c.DefaultQuery(`size`, `10`)))
	startDate := strings.TrimSpace(c.Query(`start_date`))
	endDate := strings.TrimSpace(c.Query(`end_date`))
	if len(startDate) == 0 || len(endDate) == 0 { // default to last 7 days if empty
		startDate = time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
		endDate = time.Now().Format(`2006-01-02`)
	}
	checkRet := common.CheckDataRangeDay(startDate, endDate, 365)
	if !checkRet {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `date_range_max`, 365))))
		return
	}
	if len(startDate) > 0 {
		m.Where(`date`, `>=`, startDate)
	}
	if len(endDate) > 0 {
		m.Where(`date`, `<=`, endDate)
	}

	fields := `token_app_type,to_char(date, 'YYYY-MM-DD') AS date,robot_id,prompt_token,completion_token,request_num,(prompt_token + completion_token) as total_token`
	list, total, err := m.Field(fields).Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	err = common.FillRobotName(&list, common.GetLang(c))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func StatTokenAppChart(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	m := msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))

	tokenAppType := strings.TrimSpace(c.Query(`token_app_type`))
	if tokenAppType != `` {
		m.Where(`token_app_type`, tokenAppType)
	} else {
		m.Where(`token_app_type`, `in`, strings.Join(define.GetTokenAppTypes(), `,`))
	}
	robotId := strings.TrimSpace(c.Query(`robot_id`))
	if cast.ToInt(robotId) > 0 {
		m.Where(`robot_id`, robotId)
	}
	startDate := strings.TrimSpace(c.Query(`start_date`))
	endDate := strings.TrimSpace(c.Query(`end_date`))
	if len(startDate) == 0 || len(endDate) == 0 { // default to last 7 days if empty
		startDate = time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
		endDate = time.Now().Format(`2006-01-02`)
	}
	checkRet := common.CheckDataRangeDay(startDate, endDate, 365)
	if !checkRet {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `date_range_max`, 365))))
		return
	}
	if len(startDate) > 0 {
		m.Where(`date`, `>=`, startDate)
	}
	if len(endDate) > 0 {
		m.Where(`date`, `<=`, endDate)
	}

	fields := `to_char(date, 'YYYY-MM-DD') AS date,sum(prompt_token) as prompt_token,sum(completion_token) as completion_token,sum(request_num) as request_num,sum(request_num) as request_num`
	list, err := m.Field(fields).Group(`date`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	dataList := make([]msql.Params, 0)
	err = RangeDate(startDate, endDate, common.GetLang(c), func(date string) {
		boolFind := false
		for _, item := range list {
			if item[`date`] == date {
				item[`total_token`] = cast.ToString(cast.ToInt(item[`prompt_token`]) + cast.ToInt(item[`completion_token`]))
				dataList = append(dataList, item)
				boolFind = true
			}
		}
		if !boolFind {
			dataList = append(dataList, msql.Params{
				`date`:             date,
				`prompt_token`:     `0`,
				`completion_token`: `0`,
				`request_num`:      `0`,
				`total_token`:      `0`,
			})
		}
	})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`list`: dataList}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func RangeDate(startDate, endDate, lang string, call func(string)) error {
	start, err := time.Parse(`2006-01-02`, startDate)
	if err != nil || endDate == "" {
		return errors.New(i18n.Show(lang, `param_invalid`, `start_date`))
	}
	end, err := time.Parse(`2006-01-02`, endDate)
	if err != nil {
		return errors.New(i18n.Show(lang, `param_invalid`, `end_date`))
	}
	if end.Before(start) {
		return errors.New(i18n.Show(lang, `param_invalid`, `start_date,end_date`))
	}
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		call(d.Format(`2006-01-02`))
	}
	return nil
}

func StatAnalyse(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	_type := cast.ToInt(c.Query(`type`))
	robotId := cast.ToInt(c.Query(`robot_id`))
	startDate := strings.TrimSpace(c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`)))
	endDate := strings.TrimSpace(c.DefaultQuery(`end_date`, time.Now().Format(`2006-01-02`)))
	if robotId <= 0 || _type <= 0 || len(startDate) == 0 || len(endDate) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !php2go.InArray(cast.ToInt(_type), []int{common.StatsTypeDailyActiveUser, common.StatsTypeDailyNewUser, common.StatsTypeDailyMsgCount, common.StatsTypeDailyTokenCount}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `type`))))
		return
	}
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `start_date`))))
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date`))))
		return
	}
	condition := fmt.Sprintf(`date_series.date = ds.date and admin_user_id = %d and robot_id = %d and type = %d and ds.date >= '%s' and ds.date <= '%s'`, userId, robotId, _type, startDate, endDate)
	channel := strings.TrimSpace(c.Query(`channel`))
	if len(channel) > 0 {
		condition = condition + fmt.Sprintf(`and app_type = '%s'`, channel)
	}
	m := msql.Model(fmt.Sprintf(`generate_series('%s'::date, '%s'::date, '1 day') AS date_series(date)`, startDate, endDate), define.Postgres).
		Join(`llm_request_daily_stats ds`, condition, `left`).
		Group(`date_series.date`).
		Order(`date asc`).
		Field(fmt.Sprintf(`to_char(date_series.date, 'YYYY-MM-DD') AS date,COALESCE(sum(amount), 0) as amount,%d as type`, _type))

	result, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func StatAiTipAnalyse(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	startDate := strings.TrimSpace(c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`)))
	channel := strings.TrimSpace(c.Query(`channel`))
	endDate := strings.TrimSpace(c.DefaultQuery(`end_date`, time.Now().Format(`2006-01-02`)))
	if robotId <= 0 || len(startDate) == 0 || len(endDate) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `start_date`))))
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date`))))
		return
	}
	data, err := common.StatAiTipAnalyse(userId, robotId, startDate, endDate, common.GetLang(c), channel)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func WorkflowLogs(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	startDate := strings.TrimSpace(c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`)))
	openid := strings.TrimSpace(c.Query(`openid`))
	question := strings.TrimSpace(c.Query(`question`))
	page := max(1, cast.ToInt(c.DefaultQuery(`page`, `1`)))
	size := max(1, cast.ToInt(c.DefaultQuery(`size`, `10`)))
	robotId := cast.ToInt(c.Query(`robot_id`))
	endDate := strings.TrimSpace(c.DefaultQuery(`end_date`, time.Now().Format(`2006-01-02`)))
	if len(startDate) == 0 || len(endDate) == 0 {
		startDate = time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
		endDate = time.Now().Format(`2006-01-02`)
	}
	checkRet := common.CheckDataRangeDay(startDate, endDate, 365)
	if !checkRet {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `date_range_max`, 365))))
		return
	}
	_, checkErr := common.CheckWorkflowRobotById(cast.ToString(userId), cast.ToString(robotId), common.GetLang(c))
	if checkErr != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, checkErr))
		return
	}
	startDateT, err := time.ParseInLocation(`2006-01-02`, startDate, time.Local)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `start_date`))))
		return
	}
	endDateT, err := time.ParseInLocation(`2006-01-02`, endDate, time.Local)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date`))))
		return
	}
	m := msql.Model(`work_flow_logs`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Field(`id,openid,question,node_logs,version_id,create_time`)
	m.Where(`create_time`, `>=`, cast.ToString(startDateT.Unix()))
	m.Where(`create_time`, `<=`, cast.ToString(endDateT.Unix()+86399))
	if openid != `` {
		m.Where(`openid`, openid)
	}
	if question != `` {
		m.Where(`question`, `like`, `%`+question+`%`)
	}
	list, total, err := m.Order(`id desc`).Paginate(page, size)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	versionIds := make([]string, 0)
	for _, item := range list {
		versionId := cast.ToString(item[`version_id`])
		if versionId == `` {
			continue
		}
		if !php2go.InArray(versionId, versionIds) {
			versionIds = append(versionIds, cast.ToString(item[`version_id`]))
		}
	}
	if len(versionIds) > 0 {
		versionInfos, err := msql.Model(`work_flow_version`, define.Postgres).Where(`id`, `in`, strings.Join(versionIds, `,`)).Field(`id,version`).Select()
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		for key, item := range list {
			for _, versionInfo := range versionInfos {
				if cast.ToString(item[`version_id`]) == cast.ToString(versionInfo[`id`]) {
					list[key][`version`] = versionInfo[`version`]
				}
			}
			// Duration
			durationMills := 0
			totalToken := 0
			nodeLogs := make([]define.NodeLogs, 0)
			dErr := tool.JsonDecode(item[`node_logs`], &nodeLogs)
			if dErr != nil {
				logs.Error(dErr.Error())
			} else {
				for key, nodeLog := range nodeLogs {
					durationMills += cast.ToInt(nodeLog.UseTime)
					if len(nodeLogs) != key+1 {
						totalToken += cast.ToInt(nodeLog.Output.LlmResult.PromptToken) + cast.ToInt(nodeLog.Output.LlmResult.CompletionToken)
					}
				}
			}
			list[key][`duration_mills`] = cast.ToString(durationMills)
			list[key][`total_token`] = cast.ToString(totalToken)
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`list`:  list,
		`total`: total,
		`page`:  page,
		`size`:  size,
	}, nil))
}
