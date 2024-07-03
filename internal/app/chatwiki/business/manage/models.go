// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetModelConfigList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	m := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))
	modelDefine := strings.TrimSpace(c.Query(`model_define`))
	if len(modelDefine) > 0 {
		m.Where(`model_define`, modelDefine)
	}
	configs, err := m.Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list := make([]define.ModelInfo, 0)
	for _, info := range define.ModelList {
		if len(modelDefine) == 0 || info.ModelDefine == modelDefine {
			info.ConfigList = make([]msql.Params, 0)
			list = append(list, info)
		}
	}
	for _, config := range configs {
		for i := range list {
			if list[i].ModelDefine == config[`model_define`] {
				list[i].ConfigList = append(list[i].ConfigList, config)
			}
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func AddModelConfig(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	modelDefine := strings.TrimSpace(c.PostForm(`model_define`))
	if len(modelDefine) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	modelInfo, ok := common.GetModelInfoByDefine(modelDefine)
	if !ok {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_define`))))
		return
	}
	data := msql.Datas{
		`admin_user_id`: userId,
		`model_define`:  modelDefine,
		`model_types`:   strings.Join(modelInfo.SupportedType, `,`),
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	for _, field := range modelInfo.ConfigParams {
		value := strings.TrimSpace(c.PostForm(field))
		if len(value) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
			return
		}
		if field == `model_type` { //special handling parameters
			if !tool.InArrayString(value, modelInfo.SupportedType) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_type`))))
				return
			}
			data[`model_types`] = value
		} else {
			if field == `api_version` && !tool.InArrayString(value, modelInfo.ApiVersions) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
				return
			}
			data[field] = value
		}
	}
	//database dispose
	modelId, err := msql.Model(`chat_ai_model_config`, define.Postgres).Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.ModelConfigCacheBuildHandler{ModelConfigId: int(modelId)})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DelModelConfig(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get config info
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetModelConfigInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//check relation data
	library, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`model_config_id`, cast.ToString(id)).Field(`library_name`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(library) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exist_relation_library`, library[`library_name`]))))
		return
	}
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`model_config_id`, cast.ToString(id)).Field(`robot_name`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exist_relation_robot`, robot[`robot_name`]))))
		return
	}
	robot, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`rerank_status`, `1`).Where(`rerank_model_config_id`, cast.ToString(id)).Field(`robot_name`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exist_relation_robot`, robot[`robot_name`]))))
		return
	}
	//database dispose
	_, err = msql.Model(`chat_ai_model_config`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.ModelConfigCacheBuildHandler{ModelConfigId: id})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func EditModelConfig(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get config info
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetModelConfigInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//get model define
	modelInfo, ok := common.GetModelInfoByDefine(info[`model_define`])
	if !ok {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//database dispose
	data := msql.Datas{`update_time`: tool.Time2Int()}
	for _, field := range modelInfo.ConfigParams {
		if field != `model_type` {
			value := strings.TrimSpace(c.PostForm(field))
			if len(value) == 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
				return
			}
			if field == `api_version` && !tool.InArrayString(value, modelInfo.ApiVersions) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
				return
			}
			data[field] = value
		}
	}
	_, err = msql.Model(`chat_ai_model_config`, define.Postgres).Where(`id`, cast.ToString(id)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.ModelConfigCacheBuildHandler{ModelConfigId: id})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetModelConfigOption(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	modelType := strings.TrimSpace(c.Query(`model_type`))
	isOffline := strings.TrimSpace(c.Query(`is_offline`))
	if len(modelType) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	var modelDefines []string
	if len(isOffline) > 0 {
		models := common.GetOfflineTypeModelInfos(cast.ToBool(isOffline))
		for _, model := range models {
			modelDefines = append(modelDefines, model.ModelDefine)
		}
	}

	t := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))
	if len(modelDefines) > 0 {
		t.Where(`model_define`, `in`, strings.Join(modelDefines, `,`))
	}
	configs, err := t.Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list := make([]map[string]any, 0)
	for _, config := range configs {
		if tool.InArrayString(modelType, strings.Split(config[`model_types`], `,`)) {
			modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
			list = append(list, map[string]any{`model_config`: config, `model_info`: modelInfo})
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}
