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

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetModelConfigList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	m := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
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
	for _, modelConfig := range common.GetModelConfigList() {
		if len(modelDefine) > 0 && modelConfig.ModelDefine != modelDefine {
			continue //存在检索的时候,不符合跳过
		}
		for _, config := range configs {
			if modelConfig.ModelDefine != config[`model_define`] {
				continue //模型服务商不一致,跳过
			}
			modelInfo, _ := common.GetModelInfoByConfig(adminUserId, cast.ToInt(config[`id`]))
			list = append(list, modelInfo)
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func ShowModelConfigList(c *gin.Context) {
	list := make([]common.ModelInfo, 0)
	for _, modelConfig := range common.GetModelConfigList() {
		if modelConfig.ModelDefine == common.ModelChatWiki {
			continue
		}
		modelConfig.UseModelConfigs = common.GetDefaultUseModel(modelConfig.ModelDefine) //填充默认的模型列表
		list = append(list, modelConfig)
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
	modelConfig, exist := common.GetModelConfigByDefine(modelDefine)
	if !exist {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_define`))))
		return
	}
	data := msql.Datas{
		`admin_user_id`: userId,
		`model_define`:  modelDefine,
		`model_types`:   strings.Join(modelConfig.SupportedType, `,`),
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
		`config_name`:   strings.TrimSpace(c.PostForm(`config_name`)),
	}
	for _, field := range modelConfig.ConfigParams {
		value := strings.TrimSpace(c.PostForm(field))
		if len(value) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
			return
		}
		if field == `api_version` && !tool.InArrayString(value, modelConfig.ApiVersions) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, field))))
			return
		}
		data[field] = value
	}
	//configuration test
	if err := configurationTest(common.ToStringMap(data), modelConfig); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//database dispose
	id, err := msql.Model(`chat_ai_model_config`, define.Postgres).Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//添加默认的可使用模型
	common.AutoAddDefaultUseModel(userId, int(id), modelDefine)
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.ModelConfigCacheBuildHandler{ModelConfigId: int(id)})
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
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
	id := cast.ToInt(c.PostForm(`id`))
	modelInfo, exist := common.GetModelInfoByConfig(userId, id)
	if !exist {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`模型配置ID参数错误`)))
		return
	}
	data := msql.Datas{
		`model_types`: strings.Join(modelInfo.SupportedType, `,`),
		`update_time`: tool.Time2Int(),
		`config_name`: strings.TrimSpace(c.PostForm(`config_name`)),
	}
	for _, field := range modelInfo.ConfigParams {
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
	info := modelInfo.ConfigInfo
	if tool.InArrayString(info[`model_define`], []string{common.ModelBaiduYiyan, common.ModelDoubao}) {
		secretKey := strings.TrimSpace(c.PostForm(`secret_key`))
		data[`secret_key`] = secretKey
	}
	//configuration test
	if err := configurationTest(common.ToStringMap(data, `model_define`, info[`model_define`]), modelInfo); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//database dispose
	_, err := msql.Model(`chat_ai_model_config`, define.Postgres).Where(`id`, cast.ToString(id)).Update(data)
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
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	modelType := strings.TrimSpace(c.Query(`model_type`))
	list, err := common.GetModelConfigOption(adminUserId, modelType, common.GetLang(c))
	c.String(http.StatusOK, lib_web.FmtJson(list, err))
}

func configurationTest(config msql.Params, modelInfo common.ModelInfo) error {
	if modelInfo.ModelDefine == common.ModelChatWiki {
		return errors.New(`自有模型ChatWiki禁止操作`)
	}
	//优先选取默认的模型进行参数测试
	useModels := append(common.GetDefaultUseModel(modelInfo.ModelDefine), modelInfo.UseModelConfigs...)
	if len(useModels) == 0 {
		return nil //没有获取到模型的,不校验
	}
	//调用模型测试
	modelInfo.ConfigInfo = config //替换成当前提交的参数
	handler, err := modelInfo.CallHandlerFunc(modelInfo, modelInfo.ConfigInfo, useModels[0].UseModelName)
	if err != nil {
		return err
	}
	return common.ConfigurationTest(handler.Meta, useModels[0].ModelType)
}
