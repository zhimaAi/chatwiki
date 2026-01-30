// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

const DefaultUseModelFile = define.AppRoot + `data/default_use_model.json`

type DefaultUseModelConfig struct {
	ModelDefine string `json:"model_define"`
	ModelName   string `json:"model_name"` //仅用于展示
	UseModelConfig
}

type ModelInputSupport struct {
	InputText     uint `json:"input_text,omitzero"`     //文本
	InputVoice    uint `json:"input_voice,omitzero"`    //语音
	InputImage    uint `json:"input_image,omitzero"`    //图片
	InputVideo    uint `json:"input_video,omitzero"`    //视频
	InputDocument uint `json:"input_document,omitzero"` //文档
}

type ModelOutputSupport struct {
	OutputText  uint `json:"output_text,omitzero"`  //文本
	OutputVoice uint `json:"output_voice,omitzero"` //语音
	OutputImage uint `json:"output_image,omitzero"` //图片
	OutputVideo uint `json:"output_video,omitzero"` //视频
}

type ImageGeneration struct {
	ImageSizes          string `json:"image_sizes"`            //支持的图片比例
	ImageMax            string `json:"image_max"`              //最大生成图片数量
	ImageWatermark      string `json:"image_watermark"`        //是否添加水印，1添加，0不添加
	ImageOptimizePrompt string `json:"image_optimize_prompt"`  //是否开启优化提示词，1开启，0不开启
	ImageInputsImageMax string `json:"image_inputs_image_max"` //输入图片时最多输入几张图片
}

type UseModelConfig struct {
	Id                  uint   `json:"id,omitzero"`
	ModelType           string `json:"model_type,omitzero"`
	UseModelName        string `json:"use_model_name,omitzero"`
	ShowModelName       string `json:"show_model_name,omitzero"`
	ThinkingType        uint   `json:"thinking_type,omitzero"`         //深度思考选项:0不支持,1支持,2可选
	FunctionCall        uint   `json:"function_call,omitzero"`         //是否支持function call
	ModelInputSupport                                                  //支持的输入类型
	ModelOutputSupport                                                 //支持的输出类型
	VectorDimensionList string `json:"vector_dimension_list,omitzero"` //向量维度列表(英文逗号分割)
	ImageGeneration     string `json:"image_generation,omitzero"`      //图片生成配置
}

func (useModel *UseModelConfig) ToDatas() (data msql.Datas, err error) {
	var jsonStr string
	if jsonStr, err = tool.JsonEncode(useModel); err != nil {
		return
	}
	err = tool.JsonDecodeUseNumber(jsonStr, &data)
	delete(data, `id`) //释放主键id字段
	return
}

func (useModel *UseModelConfig) ToSave(lang string, adminUserId, modelConfigId int) error {
	m := msql.Model(`chat_ai_model_list`, define.Postgres)
	if useModel.Id > 0 { //编辑可用模型
		id, _ := m.Where(`id`, cast.ToString(useModel.Id)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`model_config_id`, cast.ToString(modelConfigId)).Value(`id`)
		if cast.ToUint(id) == 0 {
			return errors.New(i18n.Show(lang, `edit_model_info_not_exist`, useModel.Id))
		}
	}
	data, err := useModel.ToDatas()
	if err != nil {
		return err
	}
	//校验模型唯一性
	if useModel.Id > 0 { //编辑可用模型时,剔除这一条记录
		m.Where(`id`, `<>`, cast.ToString(useModel.Id))
	}
	id, err := m.Where(`model_config_id`, cast.ToString(modelConfigId)).
		Where(`model_type`, useModel.ModelType).
		Where(`use_model_name`, useModel.UseModelName).Value(`id`)
	if err != nil {
		return fmt.Errorf(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	if cast.ToUint(id) > 0 {
		return errors.New(i18n.Show(lang, `model_already_exist_under_config`, useModel.UseModelName, useModel.ModelType))
	}
	//保存模型数据
	if useModel.Id == 0 { //新增可用模型
		data[`admin_user_id`] = adminUserId
		data[`model_config_id`] = modelConfigId
		data[`create_time`] = tool.Time2Int()
		data[`update_time`] = tool.Time2Int()
		_, err = m.Insert(data)
	} else { //编辑可用模型
		data[`update_time`] = tool.Time2Int()
		_, err = m.Where(`id`, cast.ToString(useModel.Id)).Update(data)
	}
	if err != nil {
		return fmt.Errorf(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &ModelListCacheBuildHandler{ModelConfigId: modelConfigId})
	//set use guide finish
	switch useModel.ModelType {
	case Llm:
		_ = SetStepFinish(adminUserId, define.StepSetLlm)
	case TextEmbedding:
		_ = SetStepFinish(adminUserId, define.StepSetEmbedding)
	}
	return nil
}

func LoadUseModelConfig(params msql.Params, modelSupplier string) UseModelConfig {
	useModel := UseModelConfig{
		Id:            cast.ToUint(params[`id`]),
		ModelType:     params[`model_type`],
		UseModelName:  params[`use_model_name`],
		ShowModelName: params[`show_model_name`],
	}
	if useModel.ModelType == Llm {
		if tool.InArrayString(modelSupplier, []string{ModelChatWiki, ModelAliyunTongyi, ModelDoubao, ModelSiliconFlow}) {
			useModel.ThinkingType = cast.ToUint(params[`thinking_type`])
		}
		useModel.FunctionCall = cast.ToUint(cast.ToBool(params[`function_call`]))
		useModel.ModelInputSupport = ModelInputSupport{
			InputText:     cast.ToUint(cast.ToBool(params[`input_text`])),
			InputVoice:    cast.ToUint(cast.ToBool(params[`input_voice`])),
			InputImage:    cast.ToUint(cast.ToBool(params[`input_image`])),
			InputVideo:    cast.ToUint(cast.ToBool(params[`input_video`])),
			InputDocument: cast.ToUint(cast.ToBool(params[`input_document`])),
		}
		useModel.ModelOutputSupport = ModelOutputSupport{
			OutputText:  cast.ToUint(cast.ToBool(params[`output_text`])),
			OutputVoice: cast.ToUint(cast.ToBool(params[`output_voice`])),
			OutputImage: cast.ToUint(cast.ToBool(params[`output_image`])),
			OutputVideo: cast.ToUint(cast.ToBool(params[`output_video`])),
		}
	}
	if useModel.ModelType == TextEmbedding {
		useModel.VectorDimensionList = params[`vector_dimension_list`]
		useModel.ModelInputSupport = ModelInputSupport{
			InputText: cast.ToUint(cast.ToBool(params[`input_text`])),
		}
	}
	if useModel.ModelType == Image {
		useModel.ImageGeneration = cast.ToString(params[`image_generation`])
		if useModel.ImageGeneration == `` {
			useModel.ImageGeneration = `{}`
		}
		useModel.ModelInputSupport = ModelInputSupport{
			InputText:  cast.ToUint(cast.ToBool(params[`input_text`])),
			InputImage: cast.ToUint(cast.ToBool(params[`input_image`])),
		}
	}
	if useModel.ModelType == Tts {
		useModel.ModelInputSupport = ModelInputSupport{
			InputText: cast.ToUint(cast.ToBool(params[`input_text`])),
		}
		useModel.ModelOutputSupport = ModelOutputSupport{
			OutputVoice: cast.ToUint(cast.ToBool(params[`output_voice`])),
		}
	}
	return useModel
}

func GetDefaultUseModel(modelDefine string) (useModels []UseModelConfig) {
	useModels = make([]UseModelConfig, 0)
	if tool.InArrayString(modelDefine, []string{ModelChatWiki, ModelOpenAIAgent, ModelAzureOpenAI, ModelOllama, ModelXnference, ModelDoubao}) {
		return //没有默认的模型
	}
	jsonStr, err := tool.ReadFile(DefaultUseModelFile)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	models := make([]DefaultUseModelConfig, 0)
	if err = tool.JsonDecodeUseNumber(jsonStr, &models); err != nil {
		logs.Error(err.Error())
		return
	}
	for _, model := range models {
		if model.ModelDefine == modelDefine {
			useModels = append(useModels, model.UseModelConfig)
		}
	}
	return
}

func AutoAddDefaultUseModel(lang string, adminUserId, modelConfigId int, modelDefine string) {
	for _, useModel := range GetDefaultUseModel(modelDefine) {
		if err := useModel.ToSave(lang, adminUserId, modelConfigId); err != nil {
			logs.Error(err.Error())
		}
	}
}

func ConfigurationTest(meta adaptor.Meta, modelType string) (err error) {
	client := &adaptor.Adaptor{}
	client.Init(meta)
	switch modelType {
	case Llm:
		messages := []adaptor.ZhimaChatCompletionMessage{{Role: `user`, Content: `configuration test`}}
		req := adaptor.ZhimaChatCompletionRequest{Messages: messages, MaxToken: 100, Temperature: 0.1}
		_, err = client.CreateChatCompletion(req)
	case TextEmbedding:
		req := adaptor.ZhimaEmbeddingRequest{Input: `configuration test`}
		_, err = client.CreateEmbeddings(req)
	case Rerank:
		req := &adaptor.ZhimaRerankReq{Passages: []string{`chatwiki`, `ChatWiki`}, Query: `ChatWiki?`, TopK: 1}
		_, err = client.CreateRerank(req)
	case Image:
		req := &adaptor.ZhimaImageGenerationReq{Prompt: "Generate a test image", Size: tea.String("2k"), ResponseFormat: tea.String("b64_json")}
		_, err = client.CreateImageGenerate(req)
		if err != nil {
			if strings.Contains(err.Error(), `prompt cannot be empty`) {
				err = nil
			}
		}
	default:
		return fmt.Errorf(`not support model_type :%s`, modelType)
	}
	return
}
