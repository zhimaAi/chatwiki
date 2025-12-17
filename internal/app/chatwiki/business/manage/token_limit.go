// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func TokenLimitCreate(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	maxToken := cast.ToInt(c.PostForm(`max_token`))
	robotId := cast.ToInt(strings.TrimSpace(c.PostForm(`robot_id`)))
	tokenAppType := strings.TrimSpace(c.PostForm(`token_app_type`))
	description := strings.TrimSpace(c.PostForm(`description`))
	if utf8.RuneCountInString(description) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `desc_max`, 500))))
		return
	}
	if tokenAppType == `` || !tool.InArray(tokenAppType, define.GetTokenAppTypes()) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `token_app_type`))))
		return
	}
	if tokenAppType != define.TokenAppTypeOther && robotId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `robot_id`))))
		return
	}
	if robotId > 0 {
		robotDbId, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(robotId)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`id`)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if cast.ToInt(robotDbId) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}
	existConfig, err := msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`token_app_type`, tokenAppType).
		Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(existConfig) > 0 {
		_, err = msql.Model(`llm_token_app_limit`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, cast.ToString(robotId)).
			Where(`token_app_type`, tokenAppType).
			Update(map[string]any{
				`description`: description,
				`max_token`:   maxToken,
				`update_time`: time.Now().Unix(),
			})
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		_, err = msql.Model(`llm_token_app_limit`, define.Postgres).
			Insert(map[string]any{
				`admin_user_id`:  adminUserId,
				`robot_id`:       robotId,
				`token_app_type`: tokenAppType,
				`description`:    description,
				`max_token`:      maxToken,
				`create_time`:    time.Now().Unix(),
				`update_time`:    time.Now().Unix(),
			})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	TokenLimitConfigCacheRemove(adminUserId, msql.Params{
		`token_app_type`: tokenAppType,
		`robot_id`:       cast.ToString(robotId),
	})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func TokenLimitList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	search := strings.TrimSpace(c.PostForm(`search`))
	robotIds := make([]string, 0)
	robotInfos := make([]msql.Params, 0)
	var err error
	m := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	if search != `` {
		m.Where(`robot_name`, `like`, search)
	}
	robotInfos, err = m.Field(`id,id as robot_id,robot_name,application_type`).Order(`id desc`).Select()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for _, robotInfo := range robotInfos {
		robotIds = append(robotIds, cast.ToString(robotInfo[`robot_id`]))
	}
	ml := msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	configs, err := ml.Select()
	if err != nil {
		logs.Error(err.Error() + ` ` + ml.GetLastSql())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotInfos = append(robotInfos, msql.Params{
		`robot_id`:       `0`,
		`robot_name`:     ``,
		`token_app_type`: define.TokenAppTypeOther,
		`is_config`:      `0`,
	})
	todayUseList, err := msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date`, time.Now().Format("2006-01-02")).
		Field(`admin_user_id,token_app_type,date,robot_id,sum(prompt_token + completion_token) as total_token`).
		Group(`admin_user_id,token_app_type,date,robot_id`).Select()
	if err != nil {
		logs.Error(err.Error() + ` ` + msql.Model(`llm_token_app_daily_stats`, define.Postgres).GetLastSql())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for _, robotInfo := range robotInfos {
		robotInfo[`token_app_type`] = common.GetTokenAppType(robotInfo)
		//今日使用token
		for _, todayUse := range todayUseList {
			if cast.ToInt(todayUse[`robot_id`]) == cast.ToInt(robotInfo[`robot_id`]) &&
				todayUse[`token_app_type`] == robotInfo[`token_app_type`] {
				robotInfo[`today_use_token`] = todayUse[`total_token`]
			}
		}
		boolFind := false
		for _, config := range configs {
			if cast.ToInt(config[`robot_id`]) == cast.ToInt(robotInfo[`robot_id`]) &&
				config[`token_app_type`] == robotInfo[`token_app_type`] {
				robotInfo[`switch_status`] = config[`switch_status`]
				robotInfo[`max_token`] = config[`max_token`]
				robotInfo[`description`] = config[`description`]
				robotInfo[`is_config`] = `1`
				robotInfo[`use_token`] = config[`use_token`]
				boolFind = true
				break
			}
		}
		delete(robotInfo, `application_type`)
		delete(robotInfo, `id`)
		if !boolFind {
			robotInfo[`switch_status`] = `0`
			robotInfo[`max_token`] = `0`
			robotInfo[`description`] = ``
			robotInfo[`is_config`] = `0`
			robotInfo[`use_token`] = `0`
		}
	}
	data := map[string]any{`list`: robotInfos}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func TokenLimitSwitch(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	tokenAppType := cast.ToString(c.PostForm(`token_app_type`))
	if tokenAppType == `` || !tool.InArray(tokenAppType, define.GetTokenAppTypes()) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `token_app_type`))))
		return
	}
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	if tokenAppType != define.TokenAppTypeOther && robotId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `robot_id`))))
		return
	}
	if robotId > 0 {
		robotDbId, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(robotId)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`id`)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if cast.ToInt(robotDbId) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}
	tokenLimit, err := msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`token_app_type`, tokenAppType).
		Field(`id,robot_id,token_app_type`).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(tokenLimit) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	switchStatus := strings.TrimSpace(c.PostForm(`switch_status`))
	if !tool.InArray(switchStatus, []string{`0`, `1`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `switch_status`))))
		return
	}
	updateData := map[string]any{
		`switch_status`: switchStatus,
		`update_time`:   time.Now().Unix(),
	}
	if cast.ToInt(switchStatus) == 0 {
		updateData[`use_token`] = 0
	}
	_, err = msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`token_app_type`, tokenAppType).
		Update(updateData)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	TokenLimitConfigCacheRemove(adminUserId, tokenLimit)
	if cast.ToInt(switchStatus) == 0 {
		useCache := common.TokenAppUseCacheBuildHandler{
			AdminUserId:  adminUserId,
			TokenAppType: cast.ToString(tokenLimit[`token_app_type`]),
			RobotId:      cast.ToInt(tokenLimit[`robot_id`]),
			DateYmd:      ``,
		}
		lib_redis.DelCacheData(define.Redis, &useCache)
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func TokenLimitConfigCacheRemove(adminUserId int, tokenLimit msql.Params) {
	configCache := common.TokenAppLimitConfigCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: cast.ToString(tokenLimit[`token_app_type`]),
		RobotId:      cast.ToInt(tokenLimit[`robot_id`]),
	}
	lib_redis.DelCacheData(define.Redis, &configCache)
}
