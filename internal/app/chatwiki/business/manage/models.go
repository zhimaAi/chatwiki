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
	"github.com/zhimaAi/llm_adaptor/adaptor"
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
	list := make([]common.ModelInfo, 0)
	for _, info := range common.GetModelList() {
		if len(modelDefine) == 0 || info.ModelDefine == modelDefine {
			info.ConfigList = make([]msql.Params, 0)
			list = append(list, info)
		}
	}
	for _, config := range configs {
		for i := range list {
			if list[i].ModelDefine == config[`model_define`] {
				historyConfigParams := make([]string, 0)
				for _, item := range list[i].HistoryConfigParams {
					if data, ok := config[item]; ok && len(data) > 0 {
						historyConfigParams = append(historyConfigParams, item)
					}
				}
				list[i].HistoryConfigParams = historyConfigParams
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
		if len(value) == 0 && field != `show_model_name` {
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
	//configuration test
	config := msql.Params{}
	for k, v := range data {
		if str, ok := v.(string); ok {
			config[k] = str
		}
	}
	err := configurationTest(config, modelInfo)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
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
	robot, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(cast.ToString(id) + `=ANY(work_flow_model_config_ids)`).Field(`robot_name`).Find()
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
			if len(value) == 0 && field != `show_model_name` {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
				return
			}
			if field == `api_version` && !tool.InArrayString(value, modelInfo.ApiVersions) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
				return
			}
			data[field] = value
			info[field] = value
		}
	}
	if tool.InArrayString(c.PostForm(`model_define`), []string{common.ModelBaiduYiyan, common.ModelDoubao}) {
		secretKey := strings.TrimSpace(c.PostForm(`secret_key`))
		data[`secret_key`] = secretKey
		info[`secret_key`] = secretKey
	}
	err = configurationTest(info, modelInfo)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
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
	if len(modelType) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list := make([]map[string]any, 0)
	for _, config := range configs {
		if tool.InArrayString(modelType, strings.Split(config[`model_types`], `,`)) {
			modelInfo, ok := common.GetModelInfoByDefine(config[`model_define`])
			if !ok {
				continue
			}
			list = append(list, map[string]any{`model_config`: config, `model_info`: modelInfo})
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func configurationTest(config msql.Params, modelInfo common.ModelInfo) error {
	modelInfo, ok := common.GetModelInfoByDefine(config[`model_define`])
	if !ok {
		return errors.New(`model define invalid`)
	}

	if strings.Contains(config[`model_types`], `LLM`) {
		handler, err := modelInfo.CallHandlerFunc(config, modelInfo.LlmModelList[0])
		if err != nil {
			return err
		}
		client := &adaptor.Adaptor{}
		client.Init(handler.Meta)

		var messages []adaptor.ZhimaChatCompletionMessage
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: `configuration test`})
		req := adaptor.ZhimaChatCompletionRequest{
			Messages:    messages,
			MaxToken:    10,
			Temperature: 0.1,
		}
		r, err := client.CreateChatCompletion(req)
		if err != nil {
			return err
		}
		logs.Info(r.Result)
	} else if strings.Contains(config[`model_types`], `TEXT EMBEDDING`) {
		handler, err := modelInfo.CallHandlerFunc(config, modelInfo.VectorModelList[0])
		if err != nil {
			return err
		}
		client := &adaptor.Adaptor{}
		client.Init(handler.Meta)

		req := adaptor.ZhimaEmbeddingRequest{Input: `configuration test`}
		_, err = client.CreateEmbeddings(req)
		if err != nil {
			return err
		}
	}
	return nil
}
