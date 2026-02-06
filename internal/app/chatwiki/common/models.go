// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type SupplierHandler struct {
	modelInfo *ModelInfo
	adaptor.Meta
	config msql.Params
}

type ModelCallHandler struct {
	modelInfo *ModelInfo
	adaptor.Meta
	config msql.Params
	// UseModel corresponding model type information
	CurModelMap map[string]UseModelConfig
}

type HandlerFunc func(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error)
type SupplierHandlerFunc func(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error)
type BeforeFunc func(info ModelInfo, config msql.Params, useModel string) error
type AfterFunc func(config msql.Params, useModel string, promptToken, completionToken int, robot msql.Params, imageNum int)

type ModelInfo struct {
	ModelDefine             string              `json:"model_define"`
	ModelName               string              `json:"model_name"`
	ModelIconUrl            string              `json:"model_icon_url"`
	Introduce               string              `json:"introduce"`
	SupportList             []string            `json:"support_list"`
	SupportedType           []string            `json:"supported_type"`
	ConfigParams            []string            `json:"config_params"`
	HistoryConfigParams     []string            `json:"history_config_params"`
	ApiVersions             []string            `json:"api_versions"`
	NetworkSearchModelList  []string            `json:"network_search_model_list"`
	HelpLinks               string              `json:"help_links"`
	CallHandlerFunc         HandlerFunc         `json:"-"`
	CallSupplierhandlerFunc SupplierHandlerFunc `json:"-"`
	CheckAllowRequest       BeforeFunc          `json:"-"`
	TokenUseReport          AfterFunc           `json:"-"`
	ConfigInfo              msql.Params         `json:"config_info"`
	UseModelConfigs         []UseModelConfig    `json:"use_model_configs"`
}

func (modelInfo *ModelInfo) SetUseModelConfigs(useModelList []msql.Params) {
	useModels := make([]UseModelConfig, len(useModelList))
	for idx, params := range useModelList {
		useModels[idx] = LoadUseModelConfig(params, modelInfo.ModelDefine)
	}
	modelInfo.UseModelConfigs = useModels
}

func (modelInfo *ModelInfo) GetModelList(modelType string, functionCall, choosableThinking bool) []string {
	models := make([]string, 0)
	for _, useModel := range modelInfo.UseModelConfigs {
		if useModel.ModelType != modelType {
			continue
		}
		if modelType == Llm && functionCall && !cast.ToBool(useModel.FunctionCall) {
			continue // Get support for function call, skip if not supported
		}
		if modelType == Llm && choosableThinking && useModel.ThinkingType != 2 { // Deep thinking option: 0 not supported, 1 supported, 2 optional
			continue // Get support for optional Thinking configuration, skip if not supported
		}
		models = append(models, useModel.UseModelName)
	}
	return models
}

// GetFunctionCallModels Get the list of models that support function call
func (modelInfo *ModelInfo) GetFunctionCallModels() []string {
	return modelInfo.GetModelList(Llm, true, false)
}

// GetChoosableThinkingModels Get the list of models that support optional Thinking configuration
func (modelInfo *ModelInfo) GetChoosableThinkingModels() []string {
	return modelInfo.GetModelList(Llm, false, true)
}

// GetLlmModelList Get the list of large language models
func (modelInfo *ModelInfo) GetLlmModelList() []string {
	return modelInfo.GetModelList(Llm, false, false)
}

// GetVectorModelList Get the list of embedding models
func (modelInfo *ModelInfo) GetVectorModelList() []string {
	return modelInfo.GetModelList(TextEmbedding, false, false)
}

// GetRerankModelList Get the list of reranking models
func (modelInfo *ModelInfo) GetRerankModelList() []string {
	return modelInfo.GetModelList(Rerank, false, false)
}

const (
	ModelChatWiki        = `chatwiki` //DIY
	ModelAzureOpenAI     = `azure`
	ModelAnthropicClaude = `claude`
	ModelGoogleGemini    = `gemini`
	ModelBaiduYiyan      = `yiyan`
	ModelAliyunTongyi    = `tongyi`
	ModelBaai            = "baai"
	ModelCohere          = "cohere"
	ModelOllama          = "ollama"
	ModelXnference       = "xinference"
	ModelDeepseek        = "deepseek"
	ModelJina            = "jina"
	ModelLingYiWanWu     = "lingyiwanwu"
	ModelMoonShot        = "moonshot"
	ModelOpenAI          = "openai"
	ModelOpenAIAgent     = "openaiAgent"
	ModelSpark           = "spark"
	ModelHunyuan         = "hunyuan"
	ModelDoubao          = "doubao"
	ModelBaichuan        = "baichuan"
	ModelZhipu           = "zhipu"
	ModelMinimax         = "minimax"
	ModelSiliconFlow     = "siliconflow"
)

const (
	Llm           = `LLM`
	TextEmbedding = `TEXT EMBEDDING`
	Speech2Text   = `SPEECH2TEXT`
	Tts           = `TTS`
	Rerank        = `RERANK`
	Image         = `IMAGE`
	MaxContent    = 10000
)

// GetModelNameByDefine Get the provider name of the specified model
func GetModelNameByDefine(lang string, modelDefine string) string {
	if modelConfig, exist := GetModelConfigByDefine(lang, modelDefine); exist {
		return modelConfig.ModelName
	}
	return fmt.Sprintf(`Unknown(%s)`, modelDefine)
}

// GetModelConfigByDefine Get the base definition of the specified model
func GetModelConfigByDefine(lang string, modelDefine string) (modelConfig ModelInfo, exist bool) {
	for _, info := range GetModelConfigList(lang) {
		if info.ModelDefine == modelDefine {
			return info, true
		}
	}
	return
}

// GetModelInfoByConfig Get the complete model information of the user configuration
func GetModelInfoByConfig(lang string, adminUserId, modelConfigId int) (_ ModelInfo, exist bool) {
	config, err := GetModelConfigInfo(modelConfigId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
	}
	if len(config) == 0 {
		return
	}
	modelConfig, exist := GetModelConfigByDefine(lang, config[`model_define`])
	if !exist {
		return
	}
	modelInfo := modelConfig // Copy a new one
	// Fill configuration information
	modelInfo.ConfigInfo = config
	// Compatibility handling for old and new data
	historyConfigParams := make([]string, 0)
	for _, item := range modelInfo.HistoryConfigParams {
		if data, ok := config[item]; ok && len(data) > 0 {
			historyConfigParams = append(historyConfigParams, item)
		}
	}
	modelInfo.HistoryConfigParams = historyConfigParams
	// Fill available model data
	if config[`model_define`] != ModelChatWiki {
		if useModelList, err := GetModelListInfo(modelConfigId); err != nil {
			logs.Error(err.Error())
		} else {
			modelInfo.SetUseModelConfigs(useModelList)
		}
	}
	return modelInfo, true
}

// GetModelConfigList Get the base definitions of all models
func GetModelConfigList(lang string) []ModelInfo {
	// Model filtering process
	list := make([]ModelInfo, 0)
	for _, info := range getModelConfigList(lang) {
		if !define.IsDev && tool.InArrayString(info.ModelDefine, []string{}) {
			continue
		}
		list = append(list, info)
	}
	// Add custom model
	//list = append(list, ModelInfo{ModelDefine: `DIY MODEL`})
	// Zero value processing
	for i, info := range list {
		if info.SupportList == nil {
			list[i].SupportList = make([]string, 0)
		}
		if info.SupportedType == nil {
			list[i].SupportedType = make([]string, 0)
		}
		if info.ConfigParams == nil {
			list[i].ConfigParams = make([]string, 0)
		}
		if info.HistoryConfigParams == nil {
			list[i].HistoryConfigParams = make([]string, 0)
		}
		if info.ApiVersions == nil {
			list[i].ApiVersions = make([]string, 0)
		}
		if info.NetworkSearchModelList == nil {
			list[i].NetworkSearchModelList = make([]string, 0)
		}
		if info.ConfigInfo == nil {
			list[i].ConfigInfo = make(msql.Params, 0)
		}
		if info.UseModelConfigs == nil {
			list[i].UseModelConfigs = make([]UseModelConfig, 0)
		}
	}
	return list
}

func getModelConfigList(lang string) []ModelInfo {
	return []ModelInfo{
		{
			ModelDefine:             ModelOpenAI,
			ModelName:               `OpenAI`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
			Introduce:               i18n.Show(lang, `model_openai_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://openai.com/`,
			CallHandlerFunc:         GetOpenAIHandle,
			CallSupplierhandlerFunc: GetOpenAISupplierHandle,
		},
		{
			ModelDefine:             ModelOpenAIAgent,
			ModelName:               i18n.Show(lang, `model_openai_agent_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
			Introduce:               i18n.Show(lang, `model_openai_agent_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_endpoint`, `api_key`, `api_version`},
			ApiVersions:             []string{"v1", `v3`},
			HelpLinks:               `https://openai.com/`,
			CallHandlerFunc:         GetOpenAIAgentHandle,
			CallSupplierhandlerFunc: GetOpenAIAgentSupplierHandle,
		},
		{
			ModelDefine:   ModelAzureOpenAI,
			ModelName:     `Azure OpenAI Service`,
			ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAzureOpenAI + `.png`,
			Introduce:     i18n.Show(lang, `model_azure_introduce`),
			SupportList:   []string{Llm, TextEmbedding, Speech2Text, Tts},
			SupportedType: []string{Llm, TextEmbedding},
			ConfigParams:  []string{`api_endpoint`, `api_key`, `api_version`},
			ApiVersions: []string{
				`2023-05-15`,
				`2023-06-01-preview`,
				`2023-10-01-preview`,
				`2024-02-15-preview`,
				`2024-03-01-preview`,
				`2024-04-01-preview`,
				`2024-05-01-preview`,
				`2024-02-01`,
			},
			HelpLinks:               `https://azure.microsoft.com/en-us/products/ai-services/openai-service`,
			CallHandlerFunc:         GetAzureHandler,
			CallSupplierhandlerFunc: GetAzureSupplierHandler,
		},
		{
			ModelDefine:             ModelAnthropicClaude,
			ModelName:               `Anthropic Claude`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelAnthropicClaude + `.png`,
			Introduce:               i18n.Show(lang, `model_claude_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://claude.ai/`,
			CallHandlerFunc:         GetClaudeHandler,
			CallSupplierhandlerFunc: GetClaudeSupplierHandler,
		},
		{
			ModelDefine:             ModelGoogleGemini,
			ModelName:               `Google Gemini`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelGoogleGemini + `.png`,
			Introduce:               i18n.Show(lang, `model_gemini_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://ai.google.dev/`,
			CallHandlerFunc:         GetGeminiHandler,
			CallSupplierhandlerFunc: GetGeminiSupplierHandler,
		},
		{
			ModelDefine:         ModelBaiduYiyan,
			ModelName:           i18n.Show(lang, `model_yiyan_name`),
			ModelIconUrl:        define.LocalUploadPrefix + `model_icon/` + ModelBaiduYiyan + `.png`,
			Introduce:           i18n.Show(lang, `model_yiyan_introduce`),
			SupportList:         []string{Llm, TextEmbedding},
			SupportedType:       []string{Llm, TextEmbedding},
			ConfigParams:        []string{`api_key`},
			HistoryConfigParams: []string{`secret_key`},
			NetworkSearchModelList: []string{
				`ernie-4.5-turbo-32k`,
				`ernie-4.5-turbo-128k`,
				`ernie-4.0-8k`,
				`ernie-x1-turbo-32k`,
				`deepseek-v3`,
				`deepseek-r1`,
			},
			HelpLinks:               `https://cloud.baidu.com/`,
			CallHandlerFunc:         GetYiyanHandler,
			CallSupplierhandlerFunc: GetYiyanSupplierHandler,
		},
		{
			ModelDefine:   ModelAliyunTongyi,
			ModelName:     i18n.Show(lang, `model_tongyi_name`),
			ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAliyunTongyi + `.png`,
			Introduce:     i18n.Show(lang, `model_tongyi_introduce`),
			SupportList:   []string{Llm, TextEmbedding, Tts, Rerank, Image},
			SupportedType: []string{Llm, TextEmbedding, Rerank, Image},
			ConfigParams:  []string{`api_key`},
			NetworkSearchModelList: []string{
				`qwen-plus`,
				`qwen-turbo`,
				`qwen3-235b-a22b`,
				`qwen-max`,
				`Moonshot-Kimi-K2-Instruct`,
			},
			HelpLinks:               `https://dashscope.aliyun.com/?spm=a2c4g.11186623.nav-dropdown-menu-0.142.6d1b46c1EeV28g&scm=20140722.X_data-37f0c4e3bf04683d35bc._.V_1`,
			CallHandlerFunc:         GetTongyiHandler,
			CallSupplierhandlerFunc: GetTongyiSupplierHandler,
		},
		{
			ModelDefine:             ModelBaai,
			ModelName:               `BGE`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelBaai + `.png`,
			Introduce:               i18n.Show(lang, `model_baai_introduce`),
			SupportList:             []string{TextEmbedding, Rerank},
			SupportedType:           []string{TextEmbedding, Rerank},
			ConfigParams:            []string{`api_endpoint`},
			HelpLinks:               `https://www.baidu.com/`,
			CallHandlerFunc:         GetBaaiHandle,
			CallSupplierhandlerFunc: GetBaaiSupplierHandle,
		},
		{
			ModelDefine:             ModelCohere,
			ModelName:               `Cohere`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelCohere + `.png`,
			Introduce:               i18n.Show(lang, `model_cohere_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Rerank},
			SupportedType:           []string{Llm, TextEmbedding, Rerank},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://cohere.com/`,
			CallHandlerFunc:         GetCohereHandle,
			CallSupplierhandlerFunc: GetCohereSupplierHandle,
		},
		{
			ModelDefine:             ModelOllama,
			ModelName:               `Ollama`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelOllama + `.png`,
			Introduce:               i18n.Show(lang, `model_ollama_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_endpoint`},
			HelpLinks:               `https://www.ollama.com/`,
			CallHandlerFunc:         GetOllamaHandle,
			CallSupplierhandlerFunc: GetOllamaSupplierHandle,
		},
		{
			ModelDefine:             ModelXnference,
			ModelName:               `xorbitsai inference`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelXnference + `.png`,
			Introduce:               i18n.Show(lang, `model_xinference_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Rerank},
			SupportedType:           []string{Llm, TextEmbedding, Rerank},
			ConfigParams:            []string{`api_version`, `api_endpoint`},
			ApiVersions:             []string{"v1"},
			HelpLinks:               `https://baidu.com/`,
			CallHandlerFunc:         GetXinferenceHandle,
			CallSupplierhandlerFunc: GetXinferenceSupplierHandle,
		},
		{
			ModelDefine:             ModelDeepseek,
			ModelName:               `DeepSeek`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelDeepseek + `.png`,
			Introduce:               i18n.Show(lang, `model_deepseek_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://www.deepseek.com/`,
			CallHandlerFunc:         GetDeepseekHandle,
			CallSupplierhandlerFunc: GetDeepseekSupplierHandle,
		},
		{
			ModelDefine:             ModelJina,
			ModelName:               `Jina`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelJina + `.png`,
			Introduce:               i18n.Show(lang, `model_jina_introduce`),
			SupportList:             []string{TextEmbedding, Rerank},
			SupportedType:           []string{TextEmbedding, Rerank},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://jina.ai/`,
			CallHandlerFunc:         GetJinaHandle,
			CallSupplierhandlerFunc: GetJinaSupplierHandle,
		},
		{
			ModelDefine:             ModelLingYiWanWu,
			ModelName:               i18n.Show(lang, `model_lingyiwanwu_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelLingYiWanWu + `.png`,
			Introduce:               i18n.Show(lang, `model_lingyiwanwu_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://platform.lingyiwanwu.com/`,
			CallHandlerFunc:         GetLingYiWanWuHandle,
			CallSupplierhandlerFunc: GetLingYiWanWuSupplierHandle,
		},
		{
			ModelDefine:             ModelMoonShot,
			ModelName:               i18n.Show(lang, `model_moonshot_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelMoonShot + `.png`,
			Introduce:               i18n.Show(lang, `model_moonshot_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://www.moonshot.cn/`,
			CallHandlerFunc:         GetMoonShotHandle,
			CallSupplierhandlerFunc: GetMoonShotSupplierHandle,
		},
		{
			ModelDefine:             ModelSpark,
			ModelName:               i18n.Show(lang, `model_spark_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelSpark + `.png`,
			Introduce:               i18n.Show(lang, `model_spark_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`app_id`, `api_key`, `secret_key`},
			HelpLinks:               `https://xinghuo.xfyun.cn/sparkapi`,
			CallHandlerFunc:         GetSparkHandle,
			CallSupplierhandlerFunc: GetSparkSupplierHandle,
		},
		{
			ModelDefine:             ModelHunyuan,
			ModelName:               i18n.Show(lang, `model_hunyuan_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelHunyuan + `.png`,
			Introduce:               i18n.Show(lang, `model_hunyuan_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`, `secret_key`},
			HelpLinks:               `https://cloud.tencent.com/product/hunyuan`,
			CallHandlerFunc:         GetHunyuanHandle,
			CallSupplierhandlerFunc: GetHunyuanSupplierHandle,
		},
		{
			ModelDefine:             ModelDoubao,
			ModelName:               i18n.Show(lang, `model_doubao_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelDoubao + `.png`,
			Introduce:               i18n.Show(lang, `model_doubao_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Image},
			SupportedType:           []string{Llm, TextEmbedding, Image},
			ConfigParams:            []string{`api_key`, `region`},
			HistoryConfigParams:     []string{`secret_key`},
			HelpLinks:               `https://www.volcengine.com/product/doubao`,
			CallHandlerFunc:         GetDoubaoHandle,
			CallSupplierhandlerFunc: GetDoubaoSupplierHandle,
		},
		{
			ModelDefine:             ModelBaichuan,
			ModelName:               i18n.Show(lang, `model_baichuan_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelBaichuan + `.png`,
			Introduce:               i18n.Show(lang, `model_baichuan_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://platform.baichuan-ai.com`,
			CallHandlerFunc:         GetBaichuanHandle,
			CallSupplierhandlerFunc: GetBaichuanSupplierHandle,
		},
		{
			ModelDefine:             ModelZhipu,
			ModelName:               i18n.Show(lang, `model_zhipu_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelZhipu + `.png`,
			Introduce:               i18n.Show(lang, `model_zhipu_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://open.bigmodel.cn/`,
			CallHandlerFunc:         GetZhipuHandle,
			CallSupplierhandlerFunc: GetZhipuSupplierHandle,
		},
		{
			ModelDefine:             ModelMinimax,
			ModelName:               `minimax`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelMinimax + `.png`,
			Introduce:               i18n.Show(lang, `model_minimax_introduce`),
			SupportList:             []string{Llm, Tts},
			SupportedType:           []string{Llm, Tts},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://www.minimaxi.com/`,
			CallHandlerFunc:         GetMinimaxHandle,
			CallSupplierhandlerFunc: GetMinimaxSupplierHandle,
		},
		{
			ModelDefine:             ModelSiliconFlow,
			ModelName:               i18n.Show(lang, `model_siliconflow_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelSiliconFlow + `.png`,
			Introduce:               i18n.Show(lang, `model_siliconflow_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Rerank},
			SupportedType:           []string{Llm, TextEmbedding, Rerank},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://siliconflow.cn/zh-cn/`,
			CallHandlerFunc:         GetSiliconFlowHandle,
			CallSupplierhandlerFunc: GetSiliconFlowSupplierHandle,
		},
	}
}

func CompatibleUseModelOldData(config msql.Params, useModel string) string {
	if len(config[`deployment_name`]) > 0 && tool.InArrayString(useModel, []string{lib_define.DefaultUseModel, config[`show_model_name`]}) {
		useModel = config[`deployment_name`]
	}
	return useModel
}

func GetSupplierCallHandler(lang string, adminUserId, modelConfigId int) (*SupplierHandler, error) {
	modelInfo, ok := GetModelInfoByConfig(lang, adminUserId, modelConfigId)
	if !ok {
		return nil, errors.New(i18n.Show(lang, `model_config_id_invalid`))
	}
	config := modelInfo.ConfigInfo
	logs.Debug("modelInfo", modelInfo)
	handler, err := modelInfo.CallSupplierhandlerFunc(modelInfo, config)
	if err != nil {
		return nil, err
	}
	handler.modelInfo = &modelInfo //save quote
	return handler, nil
}

func GetModelCallHandler(lang string, adminUserId, modelConfigId int, useModel string, robot msql.Params) (*ModelCallHandler, error) {
	modelInfo, ok := GetModelInfoByConfig(lang, adminUserId, modelConfigId)
	if !ok {
		return nil, errors.New(i18n.Show(lang, `model_config_id_invalid`))
	}
	// Validate if the used model is valid
	curModelMap := make(map[string]UseModelConfig)
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) // Compatible with old data
	for i := range modelInfo.UseModelConfigs {
		if modelInfo.UseModelConfigs[i].UseModelName == useModel {
			curModelMap[modelInfo.UseModelConfigs[i].ModelType] = modelInfo.UseModelConfigs[i]
		}
	}
	if len(curModelMap) == 0 {
		return nil, fmt.Errorf(`model(%s) not config`, useModel)
	}
	config := modelInfo.ConfigInfo
	//check token limit
	robotId := 0
	if len(robot) > 0 {
		robotId = cast.ToInt(robot[`id`])
	}
	if !TokenAppAllowUse(cast.ToInt(config[`admin_user_id`]), robotId, GetTokenAppType(robot)) {
		return nil, errors.New(`token usage exceeded`)
	}
	if modelInfo.CheckAllowRequest != nil { //check allow request
		if err := modelInfo.CheckAllowRequest(modelInfo, config, useModel); err != nil {
			return nil, err
		}
	}
	handler, err := modelInfo.CallHandlerFunc(modelInfo, config, useModel)
	if err != nil {
		return nil, err
	}
	handler.modelInfo = &modelInfo //save quote
	handler.CurModelMap = curModelMap
	return handler, nil
}

func GetVector2000(lang string, adminUserId int, openid string, robot msql.Params, library msql.Params, file msql.Params, modelConfigId int, useModel, input string) (string, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return ``, err
	}
	res, err := handler.GetVector2000(lang, adminUserId, openid, robot, library, file, input)
	if err != nil {
		return ``, err
	}
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.PromptToken, res.CompletionToken, robot, 0)
	}
	return tool.JsonEncode(res.Result)
}

func RequestChatStream(lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, chanStream chan sse.Event, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChatStream(lang, adminUserId, openid, robot, appType, messages, functionTools, chanStream, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, robot, 0)
	}
	return chatResp, requestTime, err
}

func RequestSearchStream(lang string, adminUserId int, modelConfigId int, useModel string, library msql.Params, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, chanStream chan sse.Event, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, nil)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChatStream(lang, adminUserId, "", library, "", messages, functionTools, chanStream, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, msql.Params{}, 0)
	}
	return chatResp, requestTime, err
}

func RequestChat(lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChat(lang, adminUserId, openid, robot, appType, messages, functionTools, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, robot, 0)
	}
	return chatResp, requestTime, err
}

func (h *ModelCallHandler) GetVector2000(lang string, adminUserId int, openid string, robot msql.Params, library msql.Params, fileInfo msql.Params, input string) (adaptor.ZhimaEmbeddingResponse, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaEmbeddingRequest{Input: input}
	var res adaptor.ZhimaEmbeddingResponse
	var err error
	maxTryCount := 3
	for i := 0; i < maxTryCount; i++ {
		res, err = client.CreateEmbeddings(req)
		if err != nil {
			logs.Error(err.Error())
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	if err != nil {
		return res, err
	}

	if res.Result == nil {
		return res, errors.New(`get vector return nil`)
	}
	if len(res.Result) < define.VectorDimension {
		res.Result = append(res.Result, make([]float64, define.VectorDimension-len(res.Result))...)
	}
	//go func() {
	err = LlmLogRequest(lang, TextEmbedding, adminUserId, openid, robot, library, h.config, lib_define.AppYunH5, fileInfo, h.Meta.Model, res.PromptToken, res.CompletionToken, req, res)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return res, nil
}

func (h *ModelCallHandler) GetSimilarity(query []float64, inputs [][]float64) (string, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaSimilarityRequest{Model: h.Meta.Model, Query: query, Input: inputs}
	res, err := client.CreateSimilarity(req)
	if err != nil {
		return ``, err
	}
	if res.Result == nil {
		return ``, errors.New(`get vector return nil`)
	}
	return tool.JsonEncode(res.Result)
}

func (h *ModelCallHandler) RequestRerank(lang string, adminUserId int, openid, appType string, robot msql.Params, params *adaptor.ZhimaRerankReq) (adaptor.ZhimaRerankResp, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := &adaptor.ZhimaRerankReq{
		Enable:   params.Enable,
		Query:    params.Query,
		Passages: params.Passages,
		Data:     params.Data,
		TopK:     params.TopK,
	}
	res, err := client.CreateRerank(req)
	if err != nil {
		return res, err
	}
	if res.Data == nil {
		return res, errors.New(`get rerank return nil`)
	}
	result, _ := tool.JsonEncode(res.Data)
	totalResponse := adaptor.ZhimaChatCompletionResponse{
		Result:          result,
		PromptToken:     res.InputToken,
		CompletionToken: res.OutputToken,
	}
	//go func() {
	err = LlmLogRequest(lang, Rerank, adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return res, nil
}

func (h *ModelCallHandler) RequestChatStream(
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	messages []adaptor.ZhimaChatCompletionMessage,
	functionTools []adaptor.FunctionTool,
	chanStream chan sse.Event,
	temperature float32,
	maxToken int,
) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	client := &adaptor.Adaptor{}
	if h.Meta.ChoosableThinking && len(robot) > 0 && cast.ToBool(robot[`enable_thinking`]) {
		h.Meta.EnabledThinking = true
	}
	if h.CurModelMap[Llm].InputImage > 0 && len(robot) > 0 && cast.ToBool(robot[`question_multiple_switch`]) {
		messages = ConvertQuestionMultiple(messages) // Convert to multimodal input structure
	}
	client.Init(h.Meta)
	req := adaptor.ZhimaChatCompletionRequest{
		Messages:      messages,
		MaxToken:      maxToken,
		Temperature:   float64(temperature),
		FunctionTools: functionTools,
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	defer func(stream *adaptor.ZhimaChatCompletionStreamResponse) {
		_ = stream.Close()
	}(stream)

	var totalResponse adaptor.ZhimaChatCompletionResponse
	var content string
	var functionToolCall adaptor.FunctionToolCall
	requestTime := int64(0)
	requestStartTime := time.Now()

	for {
		response, err := stream.Read()
		if requestTime == 0 {
			requestTime = time.Now().Sub(requestStartTime).Milliseconds()
			chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
		}

		totalResponse.PromptToken += response.PromptToken
		totalResponse.CompletionToken += response.CompletionToken

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return adaptor.ZhimaChatCompletionResponse{}, 0, err
		}
		if len(response.FunctionToolCalls) > 0 {
			chunkFunc := response.FunctionToolCalls[0]
			if len(functionToolCall.Name) == 0 {
				functionToolCall.Name = chunkFunc.Name
				functionToolCall.Arguments = chunkFunc.Arguments
			} else {
				functionToolCall.Arguments += chunkFunc.Arguments
			}
		}
		if len(response.ReasoningContent) > 0 {
			if cast.ToInt(robot[`think_switch`]) == define.SwitchOn {
				chanStream <- sse.Event{Event: `reasoning_content`, Data: response.ReasoningContent}
			}
			totalResponse.ReasoningContent += response.ReasoningContent
		}
		if len(response.Result) == 0 && len(functionToolCall.Arguments) == 0 {
			continue
		}

		totalResponse.Result += response.Result
		content += response.Result
		chanStream <- sse.Event{Event: `sending`, Data: response.Result}
	}

	if h.CheckFunctionArguments(functionToolCall, functionTools) && len(totalResponse.Result) == 0 { // only function call response
		totalResponse.Result = `ok`
		totalResponse.IsValidFunctionCall = true
		chanStream <- sse.Event{Event: `sending`, Data: totalResponse.Result}
	}

	//go func() {
	library := msql.Params{}
	if appType == "" && openid == "" {
		library, robot = robot, library
	}
	err = LlmLogRequest(lang, Llm, adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	if len(functionToolCall.Name) > 0 && len(functionToolCall.Arguments) > 0 {
		go func(adminUserId, robotId int, functionToolCall adaptor.FunctionToolCall) {
			err := SaveFormData(adminUserId, robotId, functionToolCall)
			if err != nil {
				logs.Error(err.Error())
			}
		}(adminUserId, cast.ToInt(robot[`id`]), functionToolCall)
	}

	return totalResponse, requestTime, nil
}

func (h *ModelCallHandler) CheckFunctionArguments(functionToolCall adaptor.FunctionToolCall, functionTools []adaptor.FunctionTool) bool {
	for _, functionTool := range functionTools {
		arguments := make(map[string]any)
		err := json.Unmarshal([]byte(functionToolCall.Arguments), &arguments)
		if err != nil {
			logs.Error(err.Error())
			break
		}
		if functionTool.Name == functionToolCall.Name {
			allRequired := true
			for _, requiredArgument := range functionTool.Parameters.Required {
				if _, ok := arguments[requiredArgument]; !ok {
					allRequired = false
					break
				}
			}
			if allRequired {
				return true
			}
		}
	}
	return false
}

func (h *ModelCallHandler) RequestChat(
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	messages []adaptor.ZhimaChatCompletionMessage,
	functionTools []adaptor.FunctionTool,
	temperature float32,
	maxToken int,
) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	client := &adaptor.Adaptor{}
	if h.Meta.ChoosableThinking && len(robot) > 0 && cast.ToBool(robot[`enable_thinking`]) {
		h.Meta.EnabledThinking = true
	}
	if h.CurModelMap[Llm].InputImage > 0 && len(robot) > 0 && cast.ToBool(robot[`question_multiple_switch`]) {
		messages = ConvertQuestionMultiple(messages) // Convert to multimodal input structure
	}
	client.Init(h.Meta)
	req := adaptor.ZhimaChatCompletionRequest{
		Messages:      messages,
		MaxToken:      maxToken,
		Temperature:   float64(temperature),
		FunctionTools: functionTools,
	}
	var functionToolCall adaptor.FunctionToolCall
	requestStartTime := time.Now()
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	if len(resp.FunctionToolCalls) > 0 {
		functionToolCall.Name = resp.FunctionToolCalls[0].Name
		functionToolCall.Arguments = resp.FunctionToolCalls[0].Arguments
	}
	requestTime := time.Now().Sub(requestStartTime).Milliseconds()
	if h.CheckFunctionArguments(functionToolCall, functionTools) && len(resp.Result) == 0 {
		resp.Result = `OK`
	}
	//go func() {
	err = LlmLogRequest(lang, Llm, adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, resp.PromptToken, resp.CompletionToken, req, resp)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	if len(functionToolCall.Name) > 0 && len(functionToolCall.Arguments) > 0 {
		go func(adminUserId, robotId int, functionToolCall adaptor.FunctionToolCall) {
			err := SaveFormData(adminUserId, robotId, functionToolCall)
			if err != nil {
				logs.Error(err.Error())
			}
		}(adminUserId, cast.ToInt(robot[`id`]), functionToolCall)
	}

	return resp, requestTime, nil
}

func CheckModelIsValid(userId, modelConfigId int, useModel, modelType string) bool {
	modelInfo, exist := GetModelInfoByConfig(define.LangEnUs, userId, modelConfigId)
	if !exist {
		return false
	}
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) // Compatible with old data
	switch modelType {
	case Llm:
		return tool.InArrayString(useModel, modelInfo.GetLlmModelList())
	case TextEmbedding:
		return tool.InArrayString(useModel, modelInfo.GetVectorModelList())
	case Rerank:
		return tool.InArrayString(useModel, modelInfo.GetRerankModelList())
	}
	return false
}

func CheckModelIsDeepSeek(model string) bool {
	modelLower := strings.ToLower(model)
	return strings.Contains(modelLower, `deepseek-r1`) ||
		strings.Contains(modelLower, `deepseek-reasoner`)
}

func CheckSupportFuncCall(lang string, adminUserId, modelConfigId int, useModel string) error {
	modelInfo, exist := GetModelInfoByConfig(lang, adminUserId, modelConfigId)
	if !exist {
		return errors.New(i18n.Show(lang, `model_config_id_invalid`))
	}
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) // Compatible with old data
	if !tool.InArrayString(useModel, modelInfo.GetLlmModelList()) {
		return errors.New(i18n.Show(lang, `use_model_name_param_error`))
	}
	if !tool.InArrayString(useModel, modelInfo.GetFunctionCallModels()) {
		return errors.New(i18n.Show(lang, `use_model_not_support_func_call`))
	}
	return nil
}

func GetModelConfigOption(adminUserId int, modelType, lang string) ([]ModelInfo, error) {
	if len(modelType) == 0 {
		return nil, errors.New(i18n.Show(lang, `param_lack`))
	}
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	list := make([]ModelInfo, 0)
	for _, config := range configs {
		if tool.InArrayString(modelType, strings.Split(config[`model_types`], `,`)) {
			modelInfo, ok := GetModelInfoByConfig(lang, adminUserId, cast.ToInt(config[`id`]))
			if !ok {
				continue
			}
			// Filter out models that are not in the current search
			useModels := make([]UseModelConfig, 0)
			for _, useModel := range modelInfo.UseModelConfigs {
				if useModel.ModelType == modelType {
					useModels = append(useModels, useModel)
				}
			}
			if len(useModels) == 0 {
				continue // Filter out empty data model providers
			}
			modelInfo.UseModelConfigs = useModels
			list = append(list, modelInfo)
		}
	}
	return list, nil
}

func (h *ModelCallHandler) RequestImageGenerate(lang string, adminUserId int, openid, appType string, robot msql.Params, params *adaptor.ZhimaImageGenerationReq) (*adaptor.ZhimaImageGenerationResp, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	params.Stream = false
	params.ResponseFormat = tea.String(`b64_json`)
	formatImageGenerateParams(params, h.Meta.Model, h.modelInfo.UseModelConfigs)
	res, err := client.CreateImageGenerate(params)
	if err != nil {
		return nil, err
	}
	if len(res.Datas) == 0 {
		return res, errors.New(`image generate empty`)
	}
	datas := make([]*adaptor.ImageGenerationData, 0)
	for _, data := range res.Datas {
		fileData, err := tool.Base64Decode(data.B64Json)
		if err != nil {
			logs.Error(`image generate base64 decode failed : %s`, err.Error())
			continue
		}
		objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, adminUserId, `image_generation`, tool.Date(`Ym`), tool.MD5(fileData), data.Ext)
		fileLink, err := WriteFileByString(objectKey, fileData)
		if err != nil {
			logs.Error(`image generate save file failed : %s`, err.Error())
			continue
		}
		data.Url = fileLink
		if !IsUrl(fileLink) {
			data.Url = define.Config.WebService[`image_domain`] + fileLink
		}
		data.B64Json = ``
		datas = append(datas, data)
	}
	res.Datas = datas
	err = LlmLogRequest(lang, Image, adminUserId, openid, robot, msql.Params{}, h.config, appType,
		msql.Params{}, h.Meta.Model, res.InputToken, res.OutputToken, params, res)
	if err != nil {
		logs.Error(err.Error())
	}
	return res, nil
}

func (h *ModelCallHandler) RequestImageGenerateStream(
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	params *adaptor.ZhimaImageGenerationReq,
	chanStream chan sse.Event,
) (*adaptor.ZhimaImageGenerationResp, int64, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	params.Stream = true
	params.ResponseFormat = tea.String(`b64_json`)
	formatImageGenerateParams(params, h.Meta.Model, h.modelInfo.UseModelConfigs)
	stream, err := client.CreateImageGenerateStream(params)
	if err != nil {
		return &adaptor.ZhimaImageGenerationResp{}, 0, err
	}
	defer func(stream *adaptor.ZhimaImageGenerationStreamRes) {
		_ = stream.Close()
	}(stream)

	var totalResponse = &adaptor.ZhimaImageGenerationResp{}
	requestTime := int64(0)
	requestStartTime := time.Now()

	for {
		response, err := stream.Read()
		if requestTime == 0 {
			requestTime = time.Now().Sub(requestStartTime).Milliseconds()
			chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
		}
		if err != nil {
			logs.Error(`image generate failed:` + err.Error())
			continue
		}

		totalResponse.InputToken += response.InputToken
		totalResponse.OutputToken += response.OutputToken

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return &adaptor.ZhimaImageGenerationResp{}, 0, err
		}
		if len(response.Datas) > 0 {
			totalResponse.Datas = append(totalResponse.Datas, response.Datas...)
		}
		chanStream <- sse.Event{Event: `sending`, Data: response.Datas}
	}

	//go func() {
	library := msql.Params{}
	if appType == "" && openid == "" {
		library, robot = robot, library
	}
	err = LlmLogRequest(lang, Image, adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.InputToken, totalResponse.OutputToken, params, totalResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return totalResponse, requestTime, nil
}

func RequestImageGenerate(lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, params *adaptor.ZhimaImageGenerationReq) (*adaptor.ZhimaImageGenerationResp, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return &adaptor.ZhimaImageGenerationResp{}, err
	}
	params.Stream = false
	res, err := handler.RequestImageGenerate(lang, adminUserId, openid, appType, robot, params)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.InputToken, res.OutputToken, robot, 0)
	}
	return res, err
}

func formatImageGenerateParams(params *adaptor.ZhimaImageGenerationReq, useModel string, useModelConfigs []UseModelConfig) {
	if len(*params.Image) > 0 && len(useModelConfigs) > 0 {
		for _, modelConfig := range useModelConfigs {
			if modelConfig.UseModelName != useModel {
				continue
			}
			imageGenerate := ImageGeneration{}
			err := tool.JsonDecode(modelConfig.ImageGeneration, &imageGenerate)
			if err != nil {
				logs.Error(err.Error())
				break
			}
			if modelConfig.InputImage != 1 {
				*params.Image = []string{}
				break
			}
			imageInputsMax := cast.ToInt(imageGenerate.ImageInputsImageMax)
			if imageInputsMax > 0 && len(*params.Image) > imageInputsMax {
				*params.Image = (*params.Image)[0:imageInputsMax]
			}
			if *params.Image != nil && len(*params.Image) > 0 {
				base64s := make([]string, 0)
				for _, imageUrl := range *params.Image {
					ext := GetUrlExt(imageUrl)
					if ext == `` {
						logs.Warning(`get url ext failed: %s`, imageUrl)
						continue
					}
					data, err := curl.Get(imageUrl).String()
					if err != nil {
						logs.Error(`get image(%s) failed: %s`, imageUrl, err.Error())
						continue
					}
					base64s = append(base64s, fmt.Sprintf(`data:image/%s;base64,%s`, ext, tool.Base64Encode(data)))
				}
				*params.Image = base64s
			}
			break
		}
	}
}

func (h *SupplierHandler) TtsGetVoiceList(adminUserId int) ([]map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	url := "https://api.minimaxi.com/v1/get_voice"

	var result map[string]any
	var voiceList []map[string]any
	request := curl.Post(url).
		Header("Authorization", "Bearer "+h.APIKey).
		Header("Content-Type", "application/json").
		Param(`voice_type`, `all`)

	if err := request.ToJSON(&result); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check response status
	if baseResp, ok := result["base_resp"].(map[string]any); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok {
			if statusCode != 0 {
				statusMsg, _ := baseResp["status_msg"].(string)
				return nil, fmt.Errorf("API error: %s", statusMsg)
			}
		}
	}

	if systemVoice, ok := result["system_voice"].([]any); ok {
		for _, v := range systemVoice {
			if voiceMap, ok := v.(map[string]any); ok {
				voiceMap["type"] = "system"
				voiceList = append(voiceList, voiceMap)
			}
		}
	}
	if voiceCloning, ok := result["voice_cloning"].([]any); ok {
		for _, v := range voiceCloning {
			if voiceMap, ok := v.(map[string]any); ok {
				voiceMap["type"] = "voice_cloning"
				voiceList = append(voiceList, voiceMap)
			}
		}
	}
	if voiceGeneration, ok := result["voice_generation"].([]any); ok {
		for _, v := range voiceGeneration {
			if voiceMap, ok := v.(map[string]any); ok {
				voiceMap["type"] = "voice_generation"
				voiceList = append(voiceList, voiceMap)
			}
		}
	}

	return voiceList, nil
}

func (h *SupplierHandler) TtsUploadVoiceFile(purpose, filePath string) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}

	url := "https://api.minimaxi.com/v1/files/upload"

	// Verify if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("file not found: " + filePath)
	}

	// Check file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.Size() > 20*1024*1024 { // 20MB
		return nil, errors.New("file size exceeds 20MB limit")
	}

	var result map[string]any
	request := curl.Post(url).
		Header("Authorization", "Bearer "+h.APIKey).
		PostFile("file", filePath).
		Param("purpose", purpose)

	if resp, err := request.Response(); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	if err := request.ToJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check response status
	if baseResp, ok := result["base_resp"].(map[string]any); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok {
			if statusCode != 0 {
				statusMsg, _ := baseResp["status_msg"].(string)
				return nil, fmt.Errorf("API error: %s", statusMsg)
			}
		}
	}

	return result, nil
}

func (h *SupplierHandler) TtsCloneVoice(params map[string]any) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}

	url := "https://api.minimaxi.com/v1/voice_clone"

	var result map[string]any
	request := curl.Post(url).
		Header("Authorization", "Bearer "+h.APIKey).
		Header("Content-Type", "application/json")

	request, err := request.JSONBody(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}
	if err = request.ToJSON(&result); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check response status
	if baseResp, ok := result["base_resp"].(map[string]any); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok {
			if statusCode != 0 {
				statusMsg, _ := baseResp["status_msg"].(string)
				return nil, fmt.Errorf("API error: %s", statusMsg)
			}
		}
	}

	return result, nil
}

func (h *ModelCallHandler) TtsSpeechT2A(params map[string]any) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}

	url := "https://api.minimaxi.com/v1/t2a_v2"

	var result map[string]any
	request := curl.Post(url).
		Header("Authorization", "Bearer "+h.APIKey).
		Header("Content-Type", "application/json")

	request, err := request.JSONBody(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}
	if err := request.ToJSON(&result); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check response status
	if baseResp, ok := result["base_resp"].(map[string]any); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok {
			if statusCode != 0 {
				statusMsg, _ := baseResp["status_msg"].(string)
				return nil, fmt.Errorf("API error: %s", statusMsg)
			}
		}
	}

	return result, nil
}

func TtsGetVoiceList(lang string, adminUserId, modelConfigId int) ([]map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsGetVoiceList(adminUserId)
}

func TtsUploadVoiceFile(lang string, adminUserId, modelConfigId int, perpose, filePath string) (map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsUploadVoiceFile(perpose, filePath)
}

func TtsCloneVoice(lang string, adminUserId, modelConfigId int, params map[string]any) (map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsCloneVoice(params)
}

func TtsSpeechT2A(lang string, adminUserId, modelConfigId int, useModel string, params map[string]any) (map[string]any, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, nil)
	if err != nil {
		return nil, err
	}
	return handler.TtsSpeechT2A(params)
}
