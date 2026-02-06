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
	modelInfo, exist := common.GetModelInfoByConfig(common.GetLang(c), adminUserId, modelConfigId)
	if !exist {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `model_config_id_invalid`))))
		return
	}
	if modelInfo.ModelDefine == common.ModelChatWiki {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `own_model_chatwiki_operation_forbidden`))))
		return
	}
	params := make(msql.Params)
	for key := range c.Request.PostForm {
		params[key] = strings.TrimSpace(c.PostForm(key)) //Too many params; wrap once
	}
	useModel := common.LoadUseModelConfig(params, modelInfo.ModelDefine)
	//Parameter validation
	if !tool.InArrayString(useModel.ModelType, []string{common.Llm, common.TextEmbedding, common.Rerank, common.Image, common.Tts}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_type`))))
		return
	}
	if len(useModel.UseModelName) == 0 || common.IsContainChinese(useModel.UseModelName) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model_name`))))
		return
	}
	if len(useModel.ShowModelName) == 0 {
		useModel.ShowModelName = useModel.UseModelName //Fill default value
	}
	if useModel.ModelType == common.Image {
		if !tool.InArrayString(modelInfo.ModelDefine, []string{common.ModelDoubao, common.ModelAliyunTongyi}) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_define`))))
			return
		}
		imageGeneration := common.ImageGeneration{}
		err := tool.JsonDecode(useModel.ImageGeneration, &imageGeneration)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_generation`))))
			return
		}
		if common.CheckArrayInArray(strings.Split(imageGeneration.ImageSizes, `,`), define.ImageSizes) >= 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_sizes`))))
			return
		}
		if cast.ToInt(imageGeneration.ImageMax) < 1 || cast.ToInt(imageGeneration.ImageMax) > 15 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_max`))))
			return
		}
		if useModel.InputImage > 0 && (cast.ToInt(imageGeneration.ImageInputsImageMax) < 1 || cast.ToInt(imageGeneration.ImageInputsImageMax) > 15) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_inputs_image_max`))))
			return
		}
		if !tool.InArray(imageGeneration.ImageWatermark, []string{`1`, `0`}) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_watermark`))))
			return
		}
		if !tool.InArray(imageGeneration.ImageOptimizePrompt, []string{`1`, `0`}) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `image_optimize_prompt`))))
			return
		}
	} else if useModel.ModelType == common.Tts {
		if modelInfo.ModelDefine != common.ModelMinimax {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_define`))))
			return
		}
	} else {
		if !tool.InArrayInt(int(useModel.ThinkingType), []int{0, 1, 2}) { //Thinking options: 0 unsupported, 1 supported, 2 optional
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `thinking_type`))))
			return
		}
		if len(useModel.VectorDimensionList) > 0 && !common.CheckIds(useModel.VectorDimensionList) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `vector_dimension_list`))))
			return
		}
	}
	//Model call test
	if useModel.ModelType != common.Tts { //TTS models are not tested for now
		handler, err := modelInfo.CallHandlerFunc(modelInfo, modelInfo.ConfigInfo, useModel.UseModelName)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		if err = common.ConfigurationTest(handler.Meta, useModel.ModelType); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		//Save data
		err = useModel.ToSave(common.GetLang(c), adminUserId, modelConfigId)
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	} else {
		err := useModel.ToSave(common.GetLang(c), adminUserId, modelConfigId)
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	}
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
	//Validate whether the available model exists
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

func GetMiniMaxVoiceList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	//Get model config
	config, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`model_define`, common.ModelMinimax).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(config) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `minimax_model_not_exist`))))
		return
	}

	//Call MiniMax API to get the voice list
	result, err := common.TtsGetVoiceList(common.GetLang(c), adminUserId, cast.ToInt(config[`id`]))
	if err != nil {
		logs.Error(`get minimax voice list failed: %v`, err)
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//Return result
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}
