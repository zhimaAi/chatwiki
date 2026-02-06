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
	ModelName   string `json:"model_name"` // For display only
	UseModelConfig
}

type ModelInputSupport struct {
	InputText     uint `json:"input_text,omitzero"`     // Text
	InputVoice    uint `json:"input_voice,omitzero"`    // Voice
	InputImage    uint `json:"input_image,omitzero"`    // Image
	InputVideo    uint `json:"input_video,omitzero"`    // Video
	InputDocument uint `json:"input_document,omitzero"` // Document
}

type ModelOutputSupport struct {
	OutputText  uint `json:"output_text,omitzero"`  // Text
	OutputVoice uint `json:"output_voice,omitzero"` // Voice
	OutputImage uint `json:"output_image,omitzero"` // Image
	OutputVideo uint `json:"output_video,omitzero"` // Video
}

type ImageGeneration struct {
	ImageSizes          string `json:"image_sizes"`            // Supported image ratios
	ImageMax            string `json:"image_max"`              // Maximum number of generated images
	ImageWatermark      string `json:"image_watermark"`        // Whether to add watermark, 1 to add, 0 not to add
	ImageOptimizePrompt string `json:"image_optimize_prompt"`  // Whether to enable optimized prompt words, 1 to enable, 0 not to enable
	ImageInputsImageMax string `json:"image_inputs_image_max"` // Maximum number of images to input when inputting images
}

type UseModelConfig struct {
	Id                  uint   `json:"id,omitzero"`
	ModelType           string `json:"model_type,omitzero"`
	UseModelName        string `json:"use_model_name,omitzero"`
	ShowModelName       string `json:"show_model_name,omitzero"`
	ThinkingType        uint   `json:"thinking_type,omitzero"` //Deep thinking option: 0 not supported, 1 supported, 2 optional
	FunctionCall        uint   `json:"function_call,omitzero"` //Whether function call is supported
	ModelInputSupport          //Supported input types
	ModelOutputSupport         //Supported output types
	VectorDimensionList string `json:"vector_dimension_list,omitzero"` //Vector dimension list (comma separated)
	ImageGeneration     string `json:"image_generation,omitzero"`      //Image generation configuration
}

func (useModel *UseModelConfig) ToDatas() (data msql.Datas, err error) {
	var jsonStr string
	if jsonStr, err = tool.JsonEncode(useModel); err != nil {
		return
	}
	err = tool.JsonDecodeUseNumber(jsonStr, &data)
	delete(data, `id`) // Release primary key id field
	return
}

func (useModel *UseModelConfig) ToSave(lang string, adminUserId, modelConfigId int) error {
	m := msql.Model(`chat_ai_model_list`, define.Postgres)
	if useModel.Id > 0 { // Edit available model
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
	// Verify model uniqueness
	if useModel.Id > 0 { // When editing an available model, exclude this record
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
	// Save model data
	if useModel.Id == 0 { // Add new available model
		data[`admin_user_id`] = adminUserId
		data[`model_config_id`] = modelConfigId
		data[`create_time`] = tool.Time2Int()
		data[`update_time`] = tool.Time2Int()
		_, err = m.Insert(data)
	} else { // Edit available model
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
		return // No default models
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
