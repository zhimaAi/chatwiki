// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	llm "github.com/zhimaAi/llm_adaptor/v2"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
	"github.com/zhimaAi/llm_adaptor/v2/embedding"
	"github.com/zhimaAi/llm_adaptor/v2/image"
	"github.com/zhimaAi/llm_adaptor/v2/rerank"
	"github.com/zhimaAi/llm_adaptor/v2/speech"
)

type SupplierHandler struct {
	modelInfo *ModelInfo
	Client    *llm.Client
	config    msql.Params
}

type ModelCallHandler struct {
	modelInfo           *ModelInfo
	Client              *llm.Client
	Model               string
	LogModel            string
	EmbeddingDimensions *int
	ChoosableThinking   bool
	config              msql.Params
	// UseModel corresponding model type information
	CurModelMap map[string]UseModelConfig
}

type ThinkingSwitch bool

const (
	ThinkingDisabled ThinkingSwitch = false
	ThinkingEnabled  ThinkingSwitch = true
)

func ToThinkingSwitch(enableThinking any) ThinkingSwitch {
	if cast.ToBool(cast.ToUint(enableThinking)) {
		return ThinkingEnabled
	}
	return ThinkingDisabled
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
	Weight                  int                 `json:"weight"`
	ApiEndPoint             string              `json:"api_end_point"`
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
	Model302Ai           = "302ai"
	ModelOpenRouter      = "openrouter"
	ModelOrcaRouter      = "orcarouter"
)

const DefaultOpenRouterEndpoint = `https://openrouter.ai/api`
const DefaultOrcaRouterEndpoint = `https://api.orcarouter.ai`

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
			ModelDefine:             Model302Ai,
			ModelName:               `302.AI`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + Model302Ai + `.png`,
			Introduce:               i18n.Show(lang, `model_302ai_introduce`),
			SupportList:             []string{Llm, Image},
			SupportedType:           []string{Llm, Image},
			ConfigParams:            []string{`api_key`},
			HistoryConfigParams:     []string{},
			HelpLinks:               `https://302.ai`,
			CallHandlerFunc:         Get302AiHandle,
			CallSupplierhandlerFunc: Get302AiSupplierHandle,
			ApiEndPoint:             `https://api.302ai.cn`,
		},
		{
			ModelDefine:             ModelOpenRouter,
			ModelName:               `OpenRouter`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelOpenRouter + `.png`,
			Introduce:               i18n.Show(lang, `model_openrouter_introduce`),
			SupportList:             []string{Llm, Image},
			SupportedType:           []string{Llm, Image},
			ConfigParams:            []string{`api_key`},
			HistoryConfigParams:     []string{},
			HelpLinks:               `https://openrouter.ai/`,
			CallHandlerFunc:         GetOpenRouterHandle,
			CallSupplierhandlerFunc: GetOpenRouterSupplierHandle,
			ApiEndPoint:             DefaultOpenRouterEndpoint,
		},
		{
			ModelDefine:             ModelOrcaRouter,
			ModelName:               `OrcaRouter`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelOrcaRouter + `.png`,
			Introduce:               i18n.Show(lang, `model_orcarouter_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HistoryConfigParams:     []string{},
			HelpLinks:               `https://www.orcarouter.ai`,
			CallHandlerFunc:         GetOrcaRouterHandle,
			CallSupplierhandlerFunc: GetOrcaRouterSupplierHandle,
			ApiEndPoint:             DefaultOrcaRouterEndpoint,
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
			ApiEndPoint:             `https://api.deepseek.com`,
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
			ApiEndPoint:             `https://generativelanguage.googleapis.com`,
		},
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
			ApiEndPoint:             `https://api.openai.com`,
		},
		{
			ModelDefine:             ModelDoubao,
			ModelName:               i18n.Show(lang, `model_doubao_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelDoubao + `.png`,
			Introduce:               i18n.Show(lang, `model_doubao_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Image},
			SupportedType:           []string{Llm, TextEmbedding, Image},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://www.volcengine.com/product/doubao`,
			CallHandlerFunc:         GetDoubaoHandle,
			CallSupplierhandlerFunc: GetDoubaoSupplierHandle,
			ApiEndPoint:             `https://ark.cn-beijing.volces.com`,
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
			ApiEndPoint:             `https://api.siliconflow.cn`,
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
			ApiEndPoint:             `https://dashscope.aliyuncs.com`,
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
			ApiEndPoint:             ``,
		},
		{
			ModelDefine:             ModelAzureOpenAI,
			ModelName:               `Azure OpenAI Service`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelAzureOpenAI + `.png`,
			Introduce:               i18n.Show(lang, `model_azure_introduce`),
			SupportList:             []string{Llm, TextEmbedding, Speech2Text, Tts, Image},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_endpoint`, `api_key`},
			HelpLinks:               `https://azure.microsoft.com/en-us/products/ai-services/openai-service`,
			CallHandlerFunc:         GetAzureHandler,
			CallSupplierhandlerFunc: GetAzureSupplierHandler,
			ApiEndPoint:             ``,
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
			ApiEndPoint:             `https://api.anthropic.com`,
		},
		{
			ModelDefine:   ModelBaiduYiyan,
			ModelName:     i18n.Show(lang, `model_yiyan_name`),
			ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelBaiduYiyan + `.png`,
			Introduce:     i18n.Show(lang, `model_yiyan_introduce`),
			SupportList:   []string{Llm, TextEmbedding},
			SupportedType: []string{Llm, TextEmbedding},
			ConfigParams:  []string{`api_key`},
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
			ModelDefine:             ModelBaai,
			ModelName:               `BGE`,
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelBaai + `.png`,
			Introduce:               i18n.Show(lang, `model_baai_introduce`),
			SupportList:             []string{TextEmbedding, Rerank},
			SupportedType:           []string{TextEmbedding, Rerank},
			ConfigParams:            []string{`api_endpoint`},
			HelpLinks:               `https://bge.baai.ac.cn/home`,
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
			ApiEndPoint:             `https://api.cohere.com`,
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
			ApiEndPoint:             ``,
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
			HelpLinks:               `https://xinference.io/zh`,
			CallHandlerFunc:         GetXinferenceHandle,
			CallSupplierhandlerFunc: GetXinferenceSupplierHandle,
			ApiEndPoint:             ``,
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
			ApiEndPoint:             `https://api.jina.ai`,
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
			ApiEndPoint:             `https://api.lingyiwanwu.com`,
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
			ApiEndPoint:             `https://api.moonshot.cn`,
		},
		{
			ModelDefine:             ModelSpark,
			ModelName:               i18n.Show(lang, `model_spark_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelSpark + `.png`,
			Introduce:               i18n.Show(lang, `model_spark_introduce`),
			SupportList:             []string{Llm},
			SupportedType:           []string{Llm},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://xinghuo.xfyun.cn/sparkapi`,
			CallHandlerFunc:         GetSparkHandle,
			CallSupplierhandlerFunc: GetSparkSupplierHandle,
			ApiEndPoint:             `https://spark-api-open.xf-yun.com/v1`,
		},
		{
			ModelDefine:             ModelHunyuan,
			ModelName:               i18n.Show(lang, `model_hunyuan_name`),
			ModelIconUrl:            define.LocalUploadPrefix + `model_icon/` + ModelHunyuan + `.png`,
			Introduce:               i18n.Show(lang, `model_hunyuan_introduce`),
			SupportList:             []string{Llm, TextEmbedding},
			SupportedType:           []string{Llm, TextEmbedding},
			ConfigParams:            []string{`api_key`},
			HelpLinks:               `https://cloud.tencent.com/product/hunyuan`,
			CallHandlerFunc:         GetHunyuanHandle,
			CallSupplierhandlerFunc: GetHunyuanSupplierHandle,
			ApiEndPoint:             `https://tokenhub.tencentmaas.com/v1`,
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
			ApiEndPoint:             `https://api.baichuan-ai.com`,
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
			ApiEndPoint:             `https://open.bigmodel.cn`,
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
			ApiEndPoint:             `https://api.minimaxi.com`,
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
	handler.LogModel = useModel
	return handler, nil
}

func GetVector2000(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, library msql.Params, file msql.Params, modelConfigId int, useModel, input string) (string, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return ``, err
	}
	res, err := handler.GetVector2000(ctx, lang, adminUserId, openid, robot, library, file, input)
	if err != nil {
		return ``, err
	}
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.Usage.PromptTokens, res.Usage.TotalTokens-res.Usage.PromptTokens, robot, 0)
	}
	values, err := res.Data[0].Embedding.Float64s()
	if err != nil {
		return ``, err
	}
	return tool.JsonEncode(values)
}

func requestChatStreamWithState(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []chat.Message, functionTools []chat.Tool, chanStream chan sse.Event, temperature float32, maxToken int, enableThinking ThinkingSwitch) (ChatResponse, int64, bool, ModelErrStage, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return ChatResponse{}, 0, false, precheckErrStage(err), err
	}
	chatResp, requestTime, streamed, stage, err := handler.requestChatStreamWithState(ctx, lang, adminUserId, openid, robot, appType, messages, functionTools, chanStream, temperature, maxToken, enableThinking)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, robot, 0)
	}
	return chatResp, requestTime, streamed, stage, err
}

func RequestChatStream(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []chat.Message, functionTools []chat.Tool, chanStream chan sse.Event, temperature float32, maxToken int, enableThinking ThinkingSwitch) (ChatResponse, int64, error) {
	chatResp, requestTime, streamed, stage, err := requestChatStreamWithState(ctx, lang, adminUserId, openid, robot, appType, modelConfigId, useModel, messages, functionTools, chanStream, temperature, maxToken, enableThinking)
	if err == nil {
		return chatResp, requestTime, nil
	}
	if ctx != nil && ctx.Err() != nil {
		return chatResp, requestTime, err
	}
	if GetTokenAppType(robot) == define.TokenAppTypeOther || stage == ModelErrPrecheck {
		return chatResp, requestTime, err
	}
	logModelError(lang, adminUserId, modelConfigId, useModel, robot, err.Error())
	if streamed {
		return chatResp, requestTime, err
	}
	backupConfigId, backupUseModel, ok := getUsableBackupModel(lang, adminUserId, modelConfigId, useModel, messages, functionTools)
	if !ok {
		return chatResp, requestTime, err
	}
	bResp, bTime, _, bStage, bErr := requestChatStreamWithState(ctx, lang, adminUserId, openid, robot, appType, backupConfigId, backupUseModel, messages, functionTools, chanStream, temperature, maxToken, enableThinking)
	if bErr == nil {
		return bResp, bTime, nil
	}
	if bStage == ModelErrProvider || bStage == ModelErrStreamRead {
		logModelError(lang, adminUserId, backupConfigId, backupUseModel, robot, bErr.Error())
	}
	return chatResp, requestTime, err
}

func RequestSearchStream(ctx context.Context, lang string, adminUserId int, modelConfigId int, useModel string, library msql.Params, messages []chat.Message, functionTools []chat.Tool, chanStream chan sse.Event, temperature float32, maxToken int) (ChatResponse, int64, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, nil)
	if err != nil {
		return ChatResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChatStream(ctx, lang, adminUserId, "", library, "", messages, functionTools, chanStream, temperature, maxToken, ThinkingDisabled)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, msql.Params{}, 0)
	}
	return chatResp, requestTime, err
}

func requestChatWithState(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []chat.Message, functionTools []chat.Tool, temperature float32, maxToken int, enableThinking ThinkingSwitch) (ChatResponse, int64, ModelErrStage, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return ChatResponse{}, 0, precheckErrStage(err), err
	}
	chatResp, requestTime, err := handler.RequestChat(ctx, lang, adminUserId, openid, robot, appType, messages, functionTools, temperature, maxToken, enableThinking)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, robot, 0)
	}
	if err != nil {
		return chatResp, requestTime, ModelErrProvider, err
	}
	return chatResp, requestTime, ModelErrNone, nil
}

func RequestChat(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []chat.Message, functionTools []chat.Tool, temperature float32, maxToken int, enableThinking ThinkingSwitch) (ChatResponse, int64, error) {
	chatResp, requestTime, stage, err := requestChatWithState(ctx, lang, adminUserId, openid, robot, appType, modelConfigId, useModel, messages, functionTools, temperature, maxToken, enableThinking)
	if err == nil {
		return chatResp, requestTime, nil
	}
	if ctx != nil && ctx.Err() != nil {
		return chatResp, requestTime, err
	}
	if GetTokenAppType(robot) == define.TokenAppTypeOther || stage == ModelErrPrecheck {
		return chatResp, requestTime, err
	}
	logModelError(lang, adminUserId, modelConfigId, useModel, robot, err.Error())
	backupConfigId, backupUseModel, ok := getUsableBackupModel(lang, adminUserId, modelConfigId, useModel, messages, functionTools)
	if !ok {
		return chatResp, requestTime, err
	}
	bResp, bTime, bStage, bErr := requestChatWithState(ctx, lang, adminUserId, openid, robot, appType, backupConfigId, backupUseModel, messages, functionTools, temperature, maxToken, enableThinking)
	if bErr == nil {
		return bResp, bTime, nil
	}
	if bStage == ModelErrProvider {
		logModelError(lang, adminUserId, backupConfigId, backupUseModel, robot, bErr.Error())
	}
	return chatResp, requestTime, err
}

func (h *ModelCallHandler) newEmbeddingRequest(input string) *embedding.CreateRequest {
	return &embedding.CreateRequest{
		Model:      h.Model,
		Input:      embedding.Input{Text: tea.String(input)},
		Dimensions: h.EmbeddingDimensions,
	}
}

func (h *ModelCallHandler) GetVector2000(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, library msql.Params, fileInfo msql.Params, input string) (*embedding.CreateResponse, error) {
	ctx = llm.NormalizeContext(ctx)
	req := h.newEmbeddingRequest(input)
	var res *embedding.CreateResponse
	var err error
	maxTryCount := 3
	for i := 0; i < maxTryCount; i++ {
		res, err = h.Client.Embeddings.Create(ctx, req)
		if err != nil {
			logLLMAdaptorError(llmAdaptorAPIEmbeddingCreate, req, res, err)
			select {
			case <-ctx.Done():
				return res, ctx.Err()
			case <-time.After(time.Second):
			}
		} else {
			break
		}
	}
	if err != nil {
		return res, err
	}

	if res == nil || len(res.Data) == 0 {
		return res, errors.New(`get vector return nil`)
	}
	values, err := res.Data[0].Embedding.Float64s()
	if err != nil {
		return res, err
	}
	if len(values) < define.VectorDimension {
		values = append(values, make([]float64, define.VectorDimension-len(values))...)
	}
	res.Data[0].Embedding = embedding.FloatEmbedding(values)
	//go func() {
	err = LlmLogRequest(lang, TextEmbedding, adminUserId, openid, robot, library, h.config, lib_define.AppYunH5, fileInfo, h.LogModel, res.Usage.PromptTokens, res.Usage.TotalTokens-res.Usage.PromptTokens, req, res)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return res, nil
}
func (h *ModelCallHandler) RequestRerank(ctx context.Context, lang string, adminUserId int, openid, appType string, robot msql.Params, req *rerank.CreateRequest) (*rerank.CreateResponse, error) {
	res, err := h.Client.Rerank.Create(ctx, req)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPIRerankCreate, req, res, err)
		return res, err
	}
	if res == nil || res.Results == nil {
		return res, errors.New(`get rerank return nil`)
	}
	inputToken, outputToken := rerankTokens(res)
	//go func() {
	err = LlmLogRequest(lang, Rerank, adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.LogModel, inputToken, outputToken, req, res)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return res, nil
}

func rerankTokens(response *rerank.CreateResponse) (int, int) {
	if response == nil || response.Meta == nil {
		return 0, 0
	}
	inputToken, outputToken := 0, 0
	if response.Meta.Tokens != nil {
		if response.Meta.Tokens.InputTokens != nil {
			inputToken = *response.Meta.Tokens.InputTokens
		}
		if response.Meta.Tokens.OutputTokens != nil {
			outputToken = *response.Meta.Tokens.OutputTokens
		}
	}
	if response.Meta.BilledUnits != nil {
		if inputToken == 0 && response.Meta.BilledUnits.InputTokens != nil {
			inputToken = *response.Meta.BilledUnits.InputTokens
		}
		if outputToken == 0 && response.Meta.BilledUnits.OutputTokens != nil {
			outputToken = *response.Meta.BilledUnits.OutputTokens
		}
	}
	return inputToken, outputToken
}

// AmendFuncToolsPropertiesType 修改函数工具的属性类型(string|number|boolean|object|array)
func AmendFuncToolsPropertiesType(Type string) string {
	switch Type {
	case TypString:
		return TypString
	case TypNumber, TypFloat, `integer`:
		return TypNumber
	case TypBoole, `boolean`:
		return `boolean`
	case TypObject, TypParams:
		return TypObject
	case TypArrString, TypArrNumber, TypArrBoole, TypArrFloat, TypArrObject, TypArrParams, `array`:
		return `array`
	default:
		logs.Warning(`unsupported type:%s`, Type)
		return TypString // fallback type
	}
}

// sendOrAbort sends an SSE event but gives up immediately if ctx is canceled
// (client disconnected), so the streaming loop never blocks on an unbuffered
// channel whose consumer has stopped reading.
func sendOrAbort(ctx context.Context, chanStream chan sse.Event, event sse.Event) {
	if chanStream == nil {
		return
	}
	var done <-chan struct{}
	if ctx != nil {
		done = ctx.Done()
	}
	select {
	case chanStream <- event:
	case <-done:
	}
}

func (h *ModelCallHandler) RequestChatStream(
	ctx context.Context,
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	messages []chat.Message,
	functionTools []chat.Tool,
	chanStream chan sse.Event,
	temperature float32,
	maxToken int,
	enableThinking ThinkingSwitch,
) (ChatResponse, int64, error) {
	chatResp, requestTime, _, _, err := h.requestChatStreamWithState(ctx, lang, adminUserId, openid, robot, appType, messages, functionTools, chanStream, temperature, maxToken, enableThinking)
	return chatResp, requestTime, err
}

func (h *ModelCallHandler) requestChatStreamWithState(
	ctx context.Context,
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	messages []chat.Message,
	functionTools []chat.Tool,
	chanStream chan sse.Event,
	temperature float32,
	maxToken int,
	enableThinking ThinkingSwitch,
) (ChatResponse, int64, bool, ModelErrStage, error) {
	if (h.CurModelMap[Llm].InputImage > 0 || h.CurModelMap[Llm].InputVideo > 0 || h.CurModelMap[Llm].InputVoice > 0) && len(robot) > 0 && cast.ToBool(robot[`question_multiple_switch`]) {
		messages = ConvertQuestionMultiple(messages) // Convert to multimodal input structure
	}
	req := h.chatRequest(messages, functionTools, temperature, maxToken, enableThinking)
	streamRequest := &chat.StreamRequest{CreateRequest: *req}
	stream, err := h.Client.Chat.Stream(ctx, streamRequest)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPIChatStream, streamRequest, nil, err)
		return ChatResponse{}, 0, false, ModelErrProvider, err
	}
	defer stream.Close()

	accumulator := chat.NewAccumulator()
	requestTime := int64(0)
	streamed := false
	requestStartTime := time.Now()

	for {
		response, err := stream.Recv()
		if requestTime == 0 {
			requestTime = time.Since(requestStartTime).Milliseconds()
			sendOrAbort(ctx, chanStream, sse.Event{Event: `request_time`, Data: requestTime})
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logLLMAdaptorError(llmAdaptorAPIChatStream, streamRequest, accumulator.Response(), err)
			return NewChatResponse(accumulator.Response()), requestTime, streamed, ModelErrStreamRead, err
		}
		if err := accumulator.Add(response); err != nil {
			logLLMAdaptorError(llmAdaptorAPIChatStream, streamRequest, accumulator.Response(), err)
			return NewChatResponse(accumulator.Response()), requestTime, streamed, ModelErrStreamRead, err
		}

		// push stream raw events for clawbot and BTS (ChatClawClient) so callers
		// can collect and merge via schema.ConcatMessages
		if len(robot) > 0 && chanStream != nil && (cast.ToInt(robot[`application_type`]) == define.ApplicationTypeClaw || appType == lib_define.ChatClawClient) {
			sendOrAbort(ctx, chanStream, sse.Event{Event: `stream_raw`, Data: response})
			streamed = true
		}
		if len(response.Choices) == 0 {
			continue
		}
		delta := response.Choices[0].Delta
		if len(delta.ReasoningContent) > 0 {
			if cast.ToInt(robot[`think_switch`]) == define.SwitchOn {
				sendOrAbort(ctx, chanStream, sse.Event{Event: `reasoning_content`, Data: delta.ReasoningContent})
				streamed = true
			}
		}
		if delta.Content.Text == nil || *delta.Content.Text == `` {
			continue
		}
		sendOrAbort(ctx, chanStream, sse.Event{Event: `sending`, Data: *delta.Content.Text})
		streamed = true
	}

	providerResponse := accumulator.Response()
	totalResponse := NewChatResponse(providerResponse)
	var functionToolCall chat.ToolCall
	if len(totalResponse.ToolCalls()) > 0 {
		functionToolCall = totalResponse.ToolCalls()[0]
	}
	if h.CheckFunctionArguments(functionToolCall, functionTools) && len(totalResponse.Result()) == 0 {
		totalResponse.SetResult(`OK`)
		totalResponse.IsValidFunctionCall = true
		sendOrAbort(ctx, chanStream, sse.Event{Event: `sending`, Data: totalResponse.Result()})
		streamed = true
	}

	//go func() {
	library := msql.Params{}
	if appType == "" && openid == "" {
		library, robot = robot, library
	}
	err = LlmLogRequest(lang, Llm, adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.LogModel, totalResponse.Usage.PromptTokens, totalResponse.Usage.CompletionTokens, req, providerResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	if len(functionToolCall.Function.Name) > 0 && len(functionToolCall.Function.Arguments) > 0 {
		go func(adminUserId, robotId int, functionToolCall chat.ToolCall) {
			err := SaveFormData(adminUserId, robotId, functionToolCall)
			if err != nil {
				logs.Error(err.Error())
			}
		}(adminUserId, cast.ToInt(robot[`id`]), functionToolCall)
	}

	return totalResponse, requestTime, streamed, ModelErrNone, nil
}

func (h *ModelCallHandler) CheckFunctionArguments(functionToolCall chat.ToolCall, functionTools []chat.Tool) bool {
	if functionToolCall.Function.Name == `` && functionToolCall.Function.Arguments == `` {
		return false // no functiontoolcall was returned
	}
	for _, functionTool := range functionTools {
		arguments := make(map[string]any)
		err := json.Unmarshal([]byte(functionToolCall.Function.Arguments), &arguments)
		if err != nil {
			logs.Error(err.Error())
			break
		}
		if functionTool.Function.Name == functionToolCall.Function.Name {
			var parameters struct {
				Required []string `json:"required"`
			}
			if err := json.Unmarshal(functionTool.Function.Parameters, &parameters); err != nil {
				logs.Error(err.Error())
				continue
			}
			allRequired := true
			for _, requiredArgument := range parameters.Required {
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

func (h *ModelCallHandler) chatRequest(messages []chat.Message, functionTools []chat.Tool, temperature float32, maxToken int, enableThinking ThinkingSwitch) *chat.CreateRequest {
	request := &chat.CreateRequest{
		Model:       h.Model,
		Messages:    append([]chat.Message(nil), messages...),
		Temperature: tea.Float64(float64(temperature)),
		Tools:       append([]chat.Tool(nil), functionTools...),
	}
	if maxToken > 0 {
		request.MaxTokens = tea.Int(maxToken)
	}
	if h.ChoosableThinking {
		request.ReasoningEffort = chat.ReasoningEffortNone
		if enableThinking == ThinkingEnabled {
			request.ReasoningEffort = chat.ReasoningEffortMedium
		}
	}
	return request
}

func (h *ModelCallHandler) RequestChat(
	ctx context.Context,
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	messages []chat.Message,
	functionTools []chat.Tool,
	temperature float32,
	maxToken int,
	enableThinking ThinkingSwitch,
) (ChatResponse, int64, error) {
	if (h.CurModelMap[Llm].InputImage > 0 || h.CurModelMap[Llm].InputVideo > 0 || h.CurModelMap[Llm].InputVoice > 0) && len(robot) > 0 && cast.ToBool(robot[`question_multiple_switch`]) {
		messages = ConvertQuestionMultiple(messages) // Convert to multimodal input structure
	}
	req := h.chatRequest(messages, functionTools, temperature, maxToken, enableThinking)
	var functionToolCall chat.ToolCall
	requestStartTime := time.Now()
	providerResponse, err := h.Client.Chat.Create(ctx, req)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPIChatCreate, req, providerResponse, err)
		return ChatResponse{}, 0, err
	}
	resp := NewChatResponse(providerResponse)
	if len(resp.ToolCalls()) > 0 {
		functionToolCall = resp.ToolCalls()[0]
	}
	requestTime := time.Since(requestStartTime).Milliseconds()
	if h.CheckFunctionArguments(functionToolCall, functionTools) && len(resp.Result()) == 0 {
		resp.SetResult(`OK`)
		resp.IsValidFunctionCall = true
	}
	//go func() {
	err = LlmLogRequest(lang, Llm, adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.LogModel, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, req, providerResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	if len(functionToolCall.Function.Name) > 0 && len(functionToolCall.Function.Arguments) > 0 {
		go func(adminUserId, robotId int, functionToolCall chat.ToolCall) {
			err := SaveFormData(adminUserId, robotId, functionToolCall)
			if err != nil {
				logs.Error(err.Error())
			}
		}(adminUserId, cast.ToInt(robot[`id`]), functionToolCall)
	}

	return resp, requestTime, nil
}

func CheckModelIsValid(userId, modelConfigId int, useModel, modelType string, runtime ...bool) bool {
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

func CheckSupportFuncCall(lang string, adminUserId, modelConfigId int, useModel string, runtime ...bool) error {
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
	// Get weight from model_define_weight table
	weightMap := make(map[string]int)
	weightData, err := msql.Model(`model_define_weight`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Select()
	if err == nil && len(weightData) > 0 {
		for _, weight := range weightData {
			weightMap[weight[`model_config_id`]] = cast.ToInt(weight[`weight`])
		}
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
			if weight, exists := weightMap[modelInfo.ConfigInfo[`id`]]; exists {
				modelInfo.Weight = weight
			} else {
				modelInfo.Weight = 0
			}
			list = append(list, modelInfo)
		}
	}

	return list, nil
}

func (h *ModelCallHandler) RequestImageGenerate(ctx context.Context, lang string, adminUserId int, openid, appType string, robot msql.Params, params *image.GenerateRequest, inputImages []string) (*image.GenerateResponse, error) {
	params.Model = h.Model
	params.ResponseFormat = `b64_json`
	files, closers, err := prepareImageFiles(inputImages, h.Model, h.modelInfo.UseModelConfigs)
	defer closeImageFiles(closers)
	if err != nil {
		return nil, err
	}
	var res *image.GenerateResponse
	api := llmAdaptorAPIImageGenerate
	if len(files) == 0 {
		res, err = h.Client.Images.Generate(ctx, params)
	} else {
		api = llmAdaptorAPIImageEdit
		res, err = h.Client.Images.Edit(ctx, imageEditRequest(params, files))
	}
	if err != nil {
		logLLMAdaptorError(api, imageRequestLog(params, files), res, err)
		return nil, err
	}
	if len(res.Data) == 0 {
		return res, errors.New(`image generate empty`)
	}
	if err := persistImageResponse(adminUserId, res); err != nil {
		return res, err
	}
	err = LlmLogRequest(lang, Image, adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.LogModel, res.Usage.InputTokens, res.Usage.OutputTokens, imageRequestLog(params, files), res)
	if err != nil {
		logs.Error(err.Error())
	}
	return res, nil
}

func (h *ModelCallHandler) RequestImageGenerateStream(
	ctx context.Context,
	lang string,
	adminUserId int,
	openid string,
	robot msql.Params,
	appType string,
	params *image.GenerateRequest,
	inputImages []string,
	chanStream chan sse.Event,
) (*image.GenerateResponse, int64, error) {
	params.Model = h.Model
	params.ResponseFormat = `b64_json`
	files, closers, err := prepareImageFiles(inputImages, h.Model, h.modelInfo.UseModelConfigs)
	defer closeImageFiles(closers)
	if err != nil {
		return &image.GenerateResponse{}, 0, err
	}
	var stream image.Stream
	api := llmAdaptorAPIImageStream
	if len(files) == 0 {
		stream, err = h.Client.Images.Stream(ctx, &image.StreamRequest{GenerateRequest: *params})
	} else {
		api = llmAdaptorAPIImageEditStream
		stream, err = h.Client.Images.EditStream(ctx, &image.EditStreamRequest{EditRequest: *imageEditRequest(params, files)})
	}
	if err != nil {
		logLLMAdaptorError(api, imageRequestLog(params, files), nil, err)
		return &image.GenerateResponse{}, 0, err
	}
	defer func(stream image.Stream) {
		_ = stream.Close()
	}(stream)

	var totalResponse = &image.GenerateResponse{OutputFormat: params.OutputFormat}
	requestTime := int64(0)
	requestStartTime := time.Now()

	for {
		response, err := stream.Recv()
		if requestTime == 0 {
			requestTime = time.Now().Sub(requestStartTime).Milliseconds()
			chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logLLMAdaptorError(api, imageRequestLog(params, files), totalResponse, err)
			return &image.GenerateResponse{}, 0, err
		}

		totalResponse.Usage.InputTokens += response.Usage.InputTokens
		totalResponse.Usage.OutputTokens += response.Usage.OutputTokens
		if response.OutputFormat != `` {
			totalResponse.OutputFormat = response.OutputFormat
		}
		if response.B64JSON != `` {
			chunkResponse := &image.GenerateResponse{Data: []image.Data{{B64JSON: response.B64JSON}}, OutputFormat: totalResponse.OutputFormat}
			if saveErr := persistImageResponse(adminUserId, chunkResponse); saveErr != nil {
				return &image.GenerateResponse{}, 0, saveErr
			}
			totalResponse.Data = append(totalResponse.Data, chunkResponse.Data...)
			chanStream <- sse.Event{Event: `sending`, Data: chunkResponse.Data}
		}
	}

	//go func() {
	library := msql.Params{}
	if appType == "" && openid == "" {
		library, robot = robot, library
	}
	err = LlmLogRequest(lang, Image, adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.LogModel, totalResponse.Usage.InputTokens, totalResponse.Usage.OutputTokens, imageRequestLog(params, files), totalResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return totalResponse, requestTime, nil
}

func RequestImageGenerate(ctx context.Context, lang string, adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, params *image.GenerateRequest, inputImages []string) (*image.GenerateResponse, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return &image.GenerateResponse{}, err
	}
	res, err := handler.RequestImageGenerate(ctx, lang, adminUserId, openid, appType, robot, params, inputImages)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.Usage.InputTokens, res.Usage.OutputTokens, robot, len(res.Data))
	}
	return res, err
}

func imageFileExtension(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case `png`, `webp`:
		return strings.ToLower(strings.TrimSpace(format))
	default:
		return `jpg`
	}
}

func persistImageResponse(adminUserId int, response *image.GenerateResponse) error {
	if response.OutputFormat == `` {
		response.OutputFormat = `jpeg`
	}
	datas := make([]image.Data, 0, len(response.Data))
	for _, data := range response.Data {
		fileData, err := tool.Base64Decode(data.B64JSON)
		if err != nil {
			logs.Error(`image generate base64 decode failed : %s`, err.Error())
			continue
		}
		ext := imageFileExtension(response.OutputFormat)
		objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, adminUserId, `image_generation`, tool.Date(`Ym`), tool.MD5(fileData), ext)
		fileLink, err := WriteFileByString(objectKey, fileData)
		if err != nil {
			logs.Error(`image generate save file failed : %s`, err.Error())
			continue
		}
		data.URL = fileLink
		if !IsUrl(fileLink) {
			data.URL = define.Config.WebService[`image_domain`] + fileLink
		}
		data.B64JSON = ``
		datas = append(datas, data)
	}
	response.Data = datas
	if len(response.Data) == 0 {
		return errors.New(`image generate contains no valid data`)
	}
	return nil
}

func prepareImageFiles(inputImages []string, useModel string, useModelConfigs []UseModelConfig) ([]image.File, []io.Closer, error) {
	links := append([]string(nil), inputImages...)
	for _, modelConfig := range useModelConfigs {
		if modelConfig.UseModelName != useModel {
			continue
		}
		if modelConfig.InputImage != 1 {
			return nil, nil, nil
		}
		imageGenerate := ImageGeneration{}
		if err := tool.JsonDecode(modelConfig.ImageGeneration, &imageGenerate); err != nil {
			return nil, nil, err
		}
		imageInputsMax := cast.ToInt(imageGenerate.ImageInputsImageMax)
		if imageInputsMax > 0 && len(links) > imageInputsMax {
			links = links[:imageInputsMax]
		}
		break
	}
	files := make([]image.File, 0, len(links))
	closers := make([]io.Closer, 0, len(links))
	for index, link := range links {
		var file image.File
		var closer io.Closer
		var err error
		if strings.HasPrefix(strings.TrimSpace(link), `data:`) {
			file, err = imageFileFromDataURL(link, index)
		} else {
			file, closer, err = imageFileFromLink(link)
		}
		if err != nil {
			logs.Error(`prepare input image(%s) failed: %s`, link, err.Error())
			continue
		}
		files = append(files, file)
		if closer != nil {
			closers = append(closers, closer)
		}
	}
	if len(links) > 0 && len(files) == 0 {
		return nil, closers, errors.New(`image edit contains no valid input file`)
	}
	return files, closers, nil
}

func imageFileFromLink(link string) (image.File, io.Closer, error) {
	filePath := GetFileByLink(link)
	if filePath == `` {
		return image.File{}, nil, fmt.Errorf(`get local file by link failed`)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return image.File{}, nil, err
	}
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	if contentType == `` {
		header := make([]byte, 512)
		readSize, readErr := file.Read(header)
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			_ = file.Close()
			return image.File{}, nil, readErr
		}
		if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
			_ = file.Close()
			return image.File{}, nil, seekErr
		}
		contentType = http.DetectContentType(header[:readSize])
	}
	return image.File{Filename: filepath.Base(filePath), ContentType: contentType, Reader: file}, file, nil
}

func imageFileFromDataURL(value string, index int) (image.File, error) {
	comma := strings.IndexByte(value, ',')
	if comma < 0 {
		return image.File{}, errors.New(`invalid image data URL`)
	}
	header := value[5:comma]
	if !strings.Contains(header, `;base64`) {
		return image.File{}, errors.New(`image data URL is not base64 encoded`)
	}
	contentType := strings.TrimSpace(strings.SplitN(header, `;`, 2)[0])
	data, err := base64.StdEncoding.DecodeString(value[comma+1:])
	if err != nil {
		return image.File{}, err
	}
	if len(data) == 0 {
		return image.File{}, errors.New(`image data URL is empty`)
	}
	extension := `.jpg`
	if extensions, extErr := mime.ExtensionsByType(contentType); extErr == nil && len(extensions) > 0 {
		extension = extensions[0]
	}
	filename := fmt.Sprintf(`image_%d%s`, index+1, extension)
	return image.File{Filename: filename, ContentType: contentType, Reader: bytes.NewReader(data)}, nil
}

func imageEditRequest(params *image.GenerateRequest, files []image.File) *image.EditRequest {
	return &image.EditRequest{
		Model: params.Model, Images: append([]image.File(nil), files...), Prompt: params.Prompt, N: params.N,
		Quality: params.Quality, ResponseFormat: params.ResponseFormat, Size: params.Size, User: params.User,
		OutputFormat: params.OutputFormat, ExtraBody: params.ExtraBody,
	}
}

func imageRequestLog(params *image.GenerateRequest, files []image.File) map[string]any {
	fileLogs := make([]map[string]string, 0, len(files))
	for _, file := range files {
		fileLogs = append(fileLogs, map[string]string{`filename`: file.Filename, `content_type`: file.ContentType})
	}
	return map[string]any{
		`model`: params.Model, `prompt`: params.Prompt, `n`: params.N, `quality`: params.Quality,
		`response_format`: params.ResponseFormat, `size`: params.Size, `user`: params.User,
		`output_format`: params.OutputFormat, `extra_body`: params.ExtraBody, `images`: fileLogs,
	}
}

func closeImageFiles(closers []io.Closer) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			logs.Error(`close image file failed: %s`, err.Error())
		}
	}
}

func (h *SupplierHandler) TtsGetVoiceList(ctx context.Context) ([]map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	request := &speech.ListVoicesRequest{VoiceType: speech.VoiceTypeAll}
	response, err := h.Client.Speech.ListVoices(ctx, request)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPISpeechListVoices, request, response, err)
		return nil, err
	}
	voiceList := make([]map[string]any, 0, len(response.SystemVoices)+len(response.ClonedVoices)+len(response.GeneratedVoices))
	appendVoices := func(voices []speech.Voice, voiceType string) {
		for _, voice := range voices {
			item, convertErr := speechResponseMap(voice)
			if convertErr != nil {
				continue
			}
			item["type"] = voiceType
			voiceList = append(voiceList, item)
		}
	}
	appendVoices(response.SystemVoices, "system")
	appendVoices(response.ClonedVoices, "voice_cloning")
	appendVoices(response.GeneratedVoices, "voice_generation")
	return voiceList, nil
}

func (h *SupplierHandler) TtsUploadVoiceFile(ctx context.Context, purpose, filePath string) (*speech.UploadVoiceFileResponse, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	request := &speech.UploadVoiceFileRequest{Purpose: speech.VoiceFilePurpose(purpose), FilePath: filePath}
	response, err := h.Client.Speech.UploadVoiceFile(ctx, request)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPISpeechUploadVoice, request, response, err)
		return nil, err
	}
	return response, nil
}

func (h *SupplierHandler) TtsCloneVoice(ctx context.Context, params map[string]any) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	request, err := decodeCloneVoiceRequest(params)
	if err != nil {
		return nil, err
	}
	response, err := h.Client.Speech.CloneVoice(ctx, request)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPISpeechCloneVoice, request, response, err)
		return nil, err
	}
	return speechResponseMap(response)
}

func (h *SupplierHandler) TtsCloneVoiceFromFiles(ctx context.Context, sourceFilePath, promptFilePath string, params map[string]any) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	cloneRequest, err := decodeCloneVoiceRequest(params)
	if err != nil {
		return nil, err
	}
	request := &speech.CloneVoiceFromFilesRequest{
		SourceFilePath: sourceFilePath,
		PromptFilePath: promptFilePath,
		CloneRequest:   *cloneRequest,
	}
	response, err := h.Client.Speech.CloneVoiceFromFiles(ctx, request)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPISpeechCloneFromFiles, request, response, err)
		return nil, err
	}
	return speechResponseMap(response.Clone)
}

func decodeCloneVoiceRequest(params map[string]any) (*speech.CloneVoiceRequest, error) {
	raw, err := tool.JsonEncode(params)
	if err != nil {
		return nil, fmt.Errorf("encode clone voice request: %w", err)
	}
	request := &speech.CloneVoiceRequest{}
	if err := tool.JsonDecodeUseNumber(raw, request); err != nil {
		return nil, fmt.Errorf("decode clone voice request: %w", err)
	}
	return request, nil
}

func speechResponseMap(response any) (map[string]any, error) {
	raw, err := tool.JsonEncode(response)
	if err != nil {
		return nil, fmt.Errorf("encode speech response: %w", err)
	}
	result := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(raw, &result); err != nil {
		return nil, fmt.Errorf("decode speech response: %w", err)
	}
	return result, nil
}

func (h *ModelCallHandler) TtsSpeechT2A(ctx context.Context, params map[string]any) (map[string]any, error) {
	if h.modelInfo.ModelDefine != ModelMinimax {
		return nil, fmt.Errorf("model not support")
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("encode speech request: %w", err)
	}
	request := &speech.CreateRequest{}
	if err := json.Unmarshal(raw, request); err != nil {
		return nil, fmt.Errorf("decode speech request: %w", err)
	}
	if request.Model == "" {
		request.Model = h.Model
	}
	response, err := h.Client.Speech.Create(ctx, request)
	if err != nil {
		logLLMAdaptorError(llmAdaptorAPISpeechCreate, request, response, err)
		return nil, err
	}
	responseRaw, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("encode speech response: %w", err)
	}
	result := make(map[string]any)
	if err := json.Unmarshal(responseRaw, &result); err != nil {
		return nil, fmt.Errorf("decode speech response: %w", err)
	}
	return result, nil
}

func TtsGetVoiceList(ctx context.Context, lang string, adminUserId, modelConfigId int) ([]map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsGetVoiceList(ctx)
}

func TtsUploadVoiceFile(ctx context.Context, lang string, adminUserId, modelConfigId int, perpose, filePath string) (*speech.UploadVoiceFileResponse, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsUploadVoiceFile(ctx, perpose, filePath)
}

func TtsCloneVoice(ctx context.Context, lang string, adminUserId, modelConfigId int, params map[string]any) (map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsCloneVoice(ctx, params)
}

func TtsCloneVoiceFromFiles(ctx context.Context, lang string, adminUserId, modelConfigId int, sourceFilePath, promptFilePath string, params map[string]any) (map[string]any, error) {
	handler, err := GetSupplierCallHandler(lang, adminUserId, modelConfigId)
	if err != nil {
		return nil, err
	}
	return handler.TtsCloneVoiceFromFiles(ctx, sourceFilePath, promptFilePath, params)
}

func TtsSpeechT2A(ctx context.Context, lang string, adminUserId, modelConfigId int, useModel string, params map[string]any) (map[string]any, error) {
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, nil)
	if err != nil {
		return nil, err
	}
	return handler.TtsSpeechT2A(ctx, params)
}
