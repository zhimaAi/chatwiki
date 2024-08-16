// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"net/http"
	"strings"
	"time"
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
