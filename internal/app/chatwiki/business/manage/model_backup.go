// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetBackupModelConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	lang := common.GetLang(c)

	var backup any
	cfg, err := common.GetBackupModelConfig(adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
		return
	}
	if len(cfg) > 0 && cast.ToInt(cfg[`model_config_id`]) > 0 {
		backupConfigId := cast.ToInt(cfg[`model_config_id`])
		backupUseModel := cfg[`use_model`]
		valid := false
		if modelInfo, ok := common.GetModelInfoByConfig(lang, adminUserId, backupConfigId); ok {
			for i := range modelInfo.UseModelConfigs {
				if modelInfo.UseModelConfigs[i].ModelType == common.Llm && modelInfo.UseModelConfigs[i].UseModelName == backupUseModel {
					backup = map[string]any{
						`model_config_id`: backupConfigId,
						`use_model`:       backupUseModel,
						`corp`:            modelInfo.ModelName,
						`model_icon_url`:  modelInfo.ModelIconUrl,
						`show_model_name`: modelInfo.UseModelConfigs[i].ShowModelName,
					}
					valid = true
					break
				}
			}
		}
		if !valid {
			_ = common.SetBackupModelConfig(adminUserId, 0, ``)
		}
	}

	options := make([]map[string]any, 0)
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
		return
	}
	for _, config := range configs {
		if !tool.InArrayString(common.Llm, strings.Split(config[`model_types`], `,`)) {
			continue
		}
		modelInfo, ok := common.GetModelInfoByConfig(lang, adminUserId, cast.ToInt(config[`id`]))
		if !ok {
			continue
		}
		models := make([]map[string]any, 0)
		for i := range modelInfo.UseModelConfigs {
			if modelInfo.UseModelConfigs[i].ModelType != common.Llm {
				continue
			}
			models = append(models, map[string]any{
				`use_model`:       modelInfo.UseModelConfigs[i].UseModelName,
				`show_model_name`: modelInfo.UseModelConfigs[i].ShowModelName,
			})
		}
		if len(models) == 0 {
			continue
		}
		options = append(options, map[string]any{
			`model_config_id`: cast.ToInt(config[`id`]),
			`corp`:            modelInfo.ModelName,
			`model_icon_url`:  modelInfo.ModelIconUrl,
			`models`:          models,
		})
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`backup`:  backup,
		`options`: options,
	}, nil))
}

func SetBackupModelConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	lang := common.GetLang(c)
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))

	if modelConfigId == 0 {
		if err := common.SetBackupModelConfig(adminUserId, 0, ``); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
			return
		}
		c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
		return
	}

	if len(useModel) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `param_lack`))))
		return
	}
	modelInfo, ok := common.GetModelInfoByConfig(lang, adminUserId, modelConfigId)
	if !ok {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `model_config_id_invalid`))))
		return
	}
	valid := false
	for i := range modelInfo.UseModelConfigs {
		if modelInfo.UseModelConfigs[i].ModelType == common.Llm && modelInfo.UseModelConfigs[i].UseModelName == useModel {
			valid = true
			break
		}
	}
	if !valid {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `param_invalid`, `use_model`))))
		return
	}

	if err := common.SetBackupModelConfig(adminUserId, modelConfigId, useModel); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetModelErrorLogs(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	lang := common.GetLang(c)
	startDate := time.Now().AddDate(0, 0, -29).Format(`2006-01-02`)

	rows, err := msql.Model(`llm_model_error_logs`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date`, `>=`, startDate).
		Field(`date,corp,model,robot_id,robot_name,application_type,count(*) as error_num`).
		Group(`date,corp,model,robot_id,robot_name,application_type`).
		Order(`date desc,error_num desc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
		return
	}

	robotIdSet := make(map[string]struct{})
	for _, row := range rows {
		if cast.ToInt(row[`robot_id`]) > 0 {
			robotIdSet[row[`robot_id`]] = struct{}{}
		}
	}
	existRobotIds := make(map[string]struct{})
	robotKeyMap := make(map[string]string)
	if len(robotIdSet) > 0 {
		ids := make([]string, 0, len(robotIdSet))
		for id := range robotIdSet {
			ids = append(ids, id)
		}
		robotRows, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, strings.Join(ids, `,`)).
			Field(`id,robot_key`).Select()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(lang, `sys_err`))))
			return
		}
		for _, r := range robotRows {
			existRobotIds[r[`id`]] = struct{}{}
			robotKeyMap[r[`id`]] = r[`robot_key`]
		}
	}

	type modelGroup struct {
		date   string
		data   map[string]any
		robots []map[string]any
	}
	groups := make(map[string]*modelGroup)
	order := make([]string, 0, len(rows))
	for _, row := range rows {
		key := row[`date`] + "\x00" + row[`corp`] + "\x00" + row[`model`]
		g, ok := groups[key]
		if !ok {
			g = &modelGroup{
				date: row[`date`],
				data: map[string]any{
					`date`:       common.DbDateToDateFormat(row[`date`]),
					`corp`:       row[`corp`],
					`model`:      row[`model`],
					`model_name`: strings.TrimSpace(row[`corp`] + ` ` + row[`model`]),
					`error_num`:  0,
				},
				robots: make([]map[string]any, 0),
			}
			groups[key] = g
			order = append(order, key)
		}
		errNum := cast.ToInt(row[`error_num`])
		g.data[`error_num`] = g.data[`error_num`].(int) + errNum
		robotId := cast.ToInt(row[`robot_id`])
		robotKey := robotKeyMap[row[`robot_id`]]
		if robotId > 0 {
			if _, exists := existRobotIds[row[`robot_id`]]; !exists {
				robotId = 0
				robotKey = ``
			}
		}
		g.robots = append(g.robots, map[string]any{
			`robot_id`:         robotId,
			`robot_key`:        robotKey,
			`robot_name`:       row[`robot_name`],
			`application_type`: cast.ToInt(row[`application_type`]),
			`error_num`:        errNum,
		})
	}

	sort.SliceStable(order, func(i, j int) bool {
		gi, gj := groups[order[i]], groups[order[j]]
		if gi.date != gj.date {
			return gi.date > gj.date
		}
		return gi.data[`error_num`].(int) > gj.data[`error_num`].(int)
	})

	list := make([]map[string]any, 0, len(order))
	for _, key := range order {
		g := groups[key]
		g.data[`robots`] = g.robots
		list = append(list, g.data)
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`list`: list}, nil))
}
