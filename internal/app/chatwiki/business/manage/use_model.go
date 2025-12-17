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

func SaveUseModelConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	if modelConfigId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	modelInfo, exist := common.GetModelInfoByConfig(adminUserId, modelConfigId)
	if !exist {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`模型配置ID参数错误`)))
		return
	}
	if modelInfo.ModelDefine == common.ModelChatWiki {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`自有模型ChatWiki禁止操作`)))
		return
	}
	params := make(msql.Params)
	for key := range c.Request.PostForm {
		params[key] = strings.TrimSpace(c.PostForm(key)) //参数太多了,套个娃
	}
	useModel := common.LoadUseModelConfig(params, modelInfo.ModelDefine)
	//参数校验
	if !tool.InArrayString(useModel.ModelType, []string{common.Llm, common.TextEmbedding, common.Rerank}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_type`))))
		return
	}
	if len(useModel.UseModelName) == 0 || common.IsContainChinese(useModel.UseModelName) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model_name`))))
		return
	}
	if len(useModel.ShowModelName) == 0 {
		useModel.ShowModelName = useModel.UseModelName //填充默认值
	}
	if !tool.InArrayInt(int(useModel.ThinkingType), []int{0, 1, 2}) { //深度思考选项:0不支持,1支持,2可选
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `thinking_type`))))
		return
	}
	if len(useModel.VectorDimensionList) > 0 && !common.CheckIds(useModel.VectorDimensionList) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `vector_dimension_list`))))
		return
	}
	//调用模型测试
	handler, err := modelInfo.CallHandlerFunc(modelInfo, modelInfo.ConfigInfo, useModel.UseModelName)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if err = common.ConfigurationTest(handler.Meta, useModel.ModelType); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//保存数据
	err = useModel.ToSave(adminUserId, modelConfigId)
	c.String(http.StatusOK, lib_web.FmtJson(nil, err))
}

func DelUseModelConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//校验可用模型是否存在
	m := msql.Model(`chat_ai_model_list`, define.Postgres)
	info, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Field(`model_config_id`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Delete(); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.ModelListCacheBuildHandler{ModelConfigId: cast.ToInt(info[`model_config_id`])})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
