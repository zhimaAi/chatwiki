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
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type ModelCallHandler struct {
	modelInfo *ModelInfo
	adaptor.Meta
	config msql.Params
}

type HandlerFunc func(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error)
type BeforeFunc func(info ModelInfo, config msql.Params, useModel string) error
type AfterFunc func(config msql.Params, useModel string, promptToken, completionToken int, robot msql.Params)

type ModelInfo struct {
	ModelDefine            string           `json:"model_define"`
	ModelName              string           `json:"model_name"`
	ModelIconUrl           string           `json:"model_icon_url"`
	Introduce              string           `json:"introduce"`
	SupportList            []string         `json:"support_list"`
	SupportedType          []string         `json:"supported_type"`
	ConfigParams           []string         `json:"config_params"`
	HistoryConfigParams    []string         `json:"history_config_params"`
	ApiVersions            []string         `json:"api_versions"`
	NetworkSearchModelList []string         `json:"network_search_model_list"`
	HelpLinks              string           `json:"help_links"`
	CallHandlerFunc        HandlerFunc      `json:"-"`
	CheckAllowRequest      BeforeFunc       `json:"-"`
	TokenUseReport         AfterFunc        `json:"-"`
	ConfigInfo             msql.Params      `json:"config_info"`
	UseModelConfigs        []UseModelConfig `json:"use_model_configs"`
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
			continue //获取支持function call,不支持的跳过
		}
		if modelType == Llm && choosableThinking && useModel.ThinkingType != 2 { //深度思考选项:0不支持,1支持,2可选
			continue //获取支持配置可选Thinking,不支持的跳过
		}
		models = append(models, useModel.UseModelName)
	}
	return models
}

// GetFunctionCallModels 获取支持function call模型列表
func (modelInfo *ModelInfo) GetFunctionCallModels() []string {
	return modelInfo.GetModelList(Llm, true, false)
}

// GetChoosableThinkingModels 获取支持配置可选Thinking的模型列表
func (modelInfo *ModelInfo) GetChoosableThinkingModels() []string {
	return modelInfo.GetModelList(Llm, false, true)
}

// GetLlmModelList 获取大语言模型列表
func (modelInfo *ModelInfo) GetLlmModelList() []string {
	return modelInfo.GetModelList(Llm, false, false)
}

// GetVectorModelList 获取嵌入模型列表
func (modelInfo *ModelInfo) GetVectorModelList() []string {
	return modelInfo.GetModelList(TextEmbedding, false, false)
}

// GetRerankModelList 获取重排序模型列表
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

// GetModelNameByDefine 获取指定模型的服务商名称
func GetModelNameByDefine(modelDefine string) string {
	if modelConfig, exist := GetModelConfigByDefine(modelDefine); exist {
		return modelConfig.ModelName
	}
	return fmt.Sprintf(`未知(%s)`, modelDefine)
}

// GetModelConfigByDefine 获取指定模型的基础定义
func GetModelConfigByDefine(modelDefine string) (modelConfig ModelInfo, exist bool) {
	for _, info := range GetModelConfigList() {
		if info.ModelDefine == modelDefine {
			return info, true
		}
	}
	return
}

// GetModelInfoByConfig 获取用户配置的模型完整信息
func GetModelInfoByConfig(adminUserId, modelConfigId int) (_ ModelInfo, exist bool) {
	config, err := GetModelConfigInfo(modelConfigId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
	}
	if len(config) == 0 {
		return
	}
	modelConfig, exist := GetModelConfigByDefine(config[`model_define`])
	if !exist {
		return
	}
	modelInfo := modelConfig //拷贝一个新的
	//填充配置信息
	modelInfo.ConfigInfo = config
	//新旧数据兼容处理
	historyConfigParams := make([]string, 0)
	for _, item := range modelInfo.HistoryConfigParams {
		if data, ok := config[item]; ok && len(data) > 0 {
			historyConfigParams = append(historyConfigParams, item)
		}
	}
	modelInfo.HistoryConfigParams = historyConfigParams
	//填充可使用模型数据
	if config[`model_define`] != ModelChatWiki {
		if useModelList, err := GetModelListInfo(modelConfigId); err != nil {
			logs.Error(err.Error())
		} else {
			modelInfo.SetUseModelConfigs(useModelList)
		}
	}
	return modelInfo, true
}

// GetModelConfigList 获取全部模型的基础定义
func GetModelConfigList() []ModelInfo {
	//模型过滤处理
	list := make([]ModelInfo, 0)
	for _, info := range modelConfigList {
		if !define.IsDev && tool.InArrayString(info.ModelDefine, []string{}) {
			continue
		}
		list = append(list, info)
	}
	//添加自定义模型
	//list = append(list, ModelInfo{ModelDefine: `DIY MODEL`})
	//零值处理
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

var modelConfigList = [...]ModelInfo{
	{
		ModelDefine:     ModelOpenAI,
		ModelName:       `OpenAI`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
		Introduce:       `基于OpenAI官方提供的API`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://openai.com/`,
		CallHandlerFunc: GetOpenAIHandle,
	},
	{
		ModelDefine:     ModelOpenAIAgent,
		ModelName:       `其他兼容OpenAI API的模型服务商`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
		Introduce:       `支持添加其他兼容OpenAi API的模型服务商，比如api2d、oneapi等`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_endpoint`, `api_key`, `api_version`},
		ApiVersions:     []string{"v1", `v3`},
		HelpLinks:       `https://openai.com/`,
		CallHandlerFunc: GetOpenAIAgentHandle,
	},
	{
		ModelDefine:   ModelAzureOpenAI,
		ModelName:     `Azure OpenAI Service`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAzureOpenAI + `.png`,
		Introduce:     `Microsoft Azure提供的OpenAI API服务`,
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
		HelpLinks:       `https://azure.microsoft.com/en-us/products/ai-services/openai-service`,
		CallHandlerFunc: GetAzureHandler,
	},
	{
		ModelDefine:     ModelAnthropicClaude,
		ModelName:       `Anthropic Claude`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelAnthropicClaude + `.png`,
		Introduce:       `Anthropic出品的Claude模型`,
		SupportList:     []string{Llm},
		SupportedType:   []string{Llm},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://claude.ai/`,
		CallHandlerFunc: GetClaudeHandler,
	},
	{
		ModelDefine:     ModelGoogleGemini,
		ModelName:       `Google Gemini`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelGoogleGemini + `.png`,
		Introduce:       `基于Google提供的Gemini API`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://ai.google.dev/`,
		CallHandlerFunc: GetGeminiHandler,
	},
	{
		ModelDefine:         ModelBaiduYiyan,
		ModelName:           `文心一言`,
		ModelIconUrl:        define.LocalUploadPrefix + `model_icon/` + ModelBaiduYiyan + `.png`,
		Introduce:           `基于百度千帆大模型平台提供的文心一言API`,
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
		HelpLinks:       `https://cloud.baidu.com/`,
		CallHandlerFunc: GetYiyanHandler,
	},
	{
		ModelDefine:   ModelAliyunTongyi,
		ModelName:     `通义千问`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAliyunTongyi + `.png`,
		Introduce:     `基于阿里云提供的通义千问API`,
		SupportList:   []string{Llm, TextEmbedding, Tts, Rerank},
		SupportedType: []string{Llm, TextEmbedding, Rerank},
		ConfigParams:  []string{`api_key`},
		NetworkSearchModelList: []string{
			`qwen-plus`,
			`qwen-turbo`,
			`qwen3-235b-a22b`,
			`qwen-max`,
			`Moonshot-Kimi-K2-Instruct`,
		},
		HelpLinks:       `https://dashscope.aliyun.com/?spm=a2c4g.11186623.nav-dropdown-menu-0.142.6d1b46c1EeV28g&scm=20140722.X_data-37f0c4e3bf04683d35bc._.V_1`,
		CallHandlerFunc: GetTongyiHandler,
	},
	{
		ModelDefine:     ModelBaai,
		ModelName:       `BGE`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelBaai + `.png`,
		Introduce:       `由北京智源人工智能研究院研发的本地模型，包含bge-rerank-base、bge-m3模型，支持嵌入和rerank。使用bge系列模型，无需消耗token，但是本地模型运行需要硬件支持，请确保服务器有足够的内存（至少8G内存）和用于计算的GPU`,
		SupportList:     []string{TextEmbedding, Rerank},
		SupportedType:   []string{TextEmbedding, Rerank},
		ConfigParams:    []string{`api_endpoint`},
		HelpLinks:       `https://www.baidu.com/`,
		CallHandlerFunc: GetBaaiHandle,
	},
	{
		ModelDefine:     ModelCohere,
		ModelName:       `Cohere`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelCohere + `.png`,
		Introduce:       `cohere提供的模型，包含Command、Command R、Command R+等`,
		SupportList:     []string{Llm, TextEmbedding, Rerank},
		SupportedType:   []string{Llm, TextEmbedding, Rerank},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://cohere.com/`,
		CallHandlerFunc: GetCohereHandle,
	},
	{
		ModelDefine:     ModelOllama,
		ModelName:       `Ollama`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelOllama + `.png`,
		Introduce:       `Ollama是一个轻量级的简单易用的本地大模型运行框架,通过Ollama可以在本地服务器构建和运营大语言模型(比如Llama3等).ChatWiki支持使用Ollama部署LLM的型和Text Embedding模型`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_endpoint`},
		HelpLinks:       `https://www.ollama.com/`,
		CallHandlerFunc: GetOllamaHandle,
	},
	{
		ModelDefine:     ModelXnference,
		ModelName:       `xorbitsai inference`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelXnference + `.png`,
		Introduce:       `Xorbits Inference(Xinference)是一个开源平台,用于简化各种AI模型的运行和集成,借助Xinference,您可以使用任何开源LLM,嵌入模型和多模态模型在本地服务器中部署`,
		SupportList:     []string{Llm, TextEmbedding, Rerank},
		SupportedType:   []string{Llm, TextEmbedding, Rerank},
		ConfigParams:    []string{`api_version`, `api_endpoint`},
		ApiVersions:     []string{"v1"},
		HelpLinks:       `https://baidu.com/`,
		CallHandlerFunc: GetXinferenceHandle,
	},
	{
		ModelDefine:     ModelDeepseek,
		ModelName:       `DeepSeek`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelDeepseek + `.png`,
		Introduce:       `由DeepSeek提供的大模型API`,
		SupportList:     []string{Llm},
		SupportedType:   []string{Llm},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://www.deepseek.com/`,
		CallHandlerFunc: GetDeepseekHandle,
	},
	{
		ModelDefine:     ModelJina,
		ModelName:       `Jina`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelJina + `.png`,
		Introduce:       `有Jina提供的嵌入和Rerank模型，`,
		SupportList:     []string{TextEmbedding, Rerank},
		SupportedType:   []string{TextEmbedding, Rerank},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://jina.ai/`,
		CallHandlerFunc: GetJinaHandle,
	},
	{
		ModelDefine:     ModelLingYiWanWu,
		ModelName:       `零一万物`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelLingYiWanWu + `.png`,
		Introduce:       `基于零一万物提供的零一大模型API`,
		SupportList:     []string{Llm},
		SupportedType:   []string{Llm},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://platform.lingyiwanwu.com/`,
		CallHandlerFunc: GetLingYiWanWuHandle,
	},
	{
		ModelDefine:     ModelMoonShot,
		ModelName:       `月之暗面`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelMoonShot + `.png`,
		Introduce:       `基于月之暗面提供的Kimi API`,
		SupportList:     []string{Llm},
		SupportedType:   []string{Llm},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://www.moonshot.cn/`,
		CallHandlerFunc: GetMoonShotHandle,
	},
	{
		ModelDefine:     ModelSpark,
		ModelName:       `讯飞星火`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelSpark + `.png`,
		Introduce:       `基于科大讯飞提供的讯飞星火大模型API`,
		SupportList:     []string{Llm},
		SupportedType:   []string{Llm},
		ConfigParams:    []string{`app_id`, `api_key`, `secret_key`},
		HelpLinks:       `https://xinghuo.xfyun.cn/sparkapi`,
		CallHandlerFunc: GetSparkHandle,
	},
	{
		ModelDefine:     ModelHunyuan,
		ModelName:       `腾讯混元`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelHunyuan + `.png`,
		Introduce:       `腾讯混元大模型由腾讯公司全链路自研`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_key`, `secret_key`},
		HelpLinks:       `https://cloud.tencent.com/product/hunyuan`,
		CallHandlerFunc: GetHunyuanHandle,
	},
	{
		ModelDefine:         ModelDoubao,
		ModelName:           `火山引擎`,
		ModelIconUrl:        define.LocalUploadPrefix + `model_icon/` + ModelDoubao + `.png`,
		Introduce:           `基于火山引擎提供的豆包大模型API`,
		SupportList:         []string{Llm, TextEmbedding, Image},
		SupportedType:       []string{Llm, TextEmbedding, Image},
		ConfigParams:        []string{`api_key`, `region`},
		HistoryConfigParams: []string{`secret_key`},
		HelpLinks:           `https://www.volcengine.com/product/doubao`,
		CallHandlerFunc:     GetDoubaoHandle,
	},
	{
		ModelDefine:     ModelBaichuan,
		ModelName:       `百川智能`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelBaichuan + `.png`,
		Introduce:       `基于百川智能提供的百川大模型API`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://platform.baichuan-ai.com`,
		CallHandlerFunc: GetBaichuanHandle,
	},
	{
		ModelDefine:     ModelZhipu,
		ModelName:       `智谱`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelZhipu + `.png`,
		Introduce:       `领先的认知大模型AI开放平台`,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://open.bigmodel.cn/`,
		CallHandlerFunc: GetZhipuHandle,
	},
	{
		ModelDefine:     ModelMinimax,
		ModelName:       `minimax`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelMinimax + `.png`,
		Introduce:       `MiniMax 成立于 2021 年 12 月，是领先的通用人工智能科技公司，致力于与用户共创智能。MiniMax 自主研发多模态、万亿参数的 MoE 大模型，并基于大模型推出海螺AI、星野等原生应用。MiniMax API 开放平台提供安全、灵活、可靠的 API 服务，助力企业和开发者快速搭建 AI 应用。`,
		SupportList:     []string{Llm, Tts},
		SupportedType:   []string{Llm, Tts},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://www.minimaxi.com/`,
		CallHandlerFunc: GetMinimaxHandle,
	},
	{
		ModelDefine:     ModelSiliconFlow,
		ModelName:       `硅基流动`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelSiliconFlow + `.png`,
		Introduce:       `支持通义千问，mata-lama，google-gemma，bge-m3等开源模型，可以免部署、低成本使用`,
		SupportList:     []string{Llm, TextEmbedding, Rerank},
		SupportedType:   []string{Llm, TextEmbedding, Rerank},
		ConfigParams:    []string{`api_key`},
		HelpLinks:       `https://siliconflow.cn/zh-cn/`,
		CallHandlerFunc: GetSiliconFlow,
	},
}

func CompatibleUseModelOldData(config msql.Params, useModel string) string {
	if len(config[`deployment_name`]) > 0 && tool.InArrayString(useModel, []string{`默认`, config[`show_model_name`]}) {
		useModel = config[`deployment_name`]
	}
	return useModel
}

func GetModelCallHandler(adminUserId, modelConfigId int, useModel string, robot msql.Params) (*ModelCallHandler, error) {
	modelInfo, ok := GetModelInfoByConfig(adminUserId, modelConfigId)
	if !ok {
		return nil, errors.New(`模型配置ID参数错误`)
	}
	//校验使用的模型是否有效
	var isValid bool
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) //兼容旧数据
	for i := range modelInfo.UseModelConfigs {
		if modelInfo.UseModelConfigs[i].UseModelName == useModel {
			isValid = true
			break
		}
	}
	if !isValid {
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
	return handler, nil
}

func GetVector2000(adminUserId int, openid string, robot msql.Params, library msql.Params, file msql.Params, modelConfigId int, useModel, input string) (string, error) {
	handler, err := GetModelCallHandler(adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return ``, err
	}
	res, err := handler.GetVector2000(adminUserId, openid, robot, library, file, input)
	if err != nil {
		return ``, err
	}
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.PromptToken, res.CompletionToken, robot)
	}
	return tool.JsonEncode(res.Result)
}

func RequestChatStream(adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, chanStream chan sse.Event, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChatStream(adminUserId, openid, robot, appType, messages, functionTools, chanStream, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, robot)
	}
	return chatResp, requestTime, err
}

func RequestSearchStream(adminUserId int, modelConfigId int, useModel string, library msql.Params, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, chanStream chan sse.Event, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(adminUserId, modelConfigId, useModel, nil)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChatStream(adminUserId, "", library, "", messages, functionTools, chanStream, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, msql.Params{})
	}
	return chatResp, requestTime, err
}

func RequestChat(adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, functionTools []adaptor.FunctionTool, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	chatResp, requestTime, err := handler.RequestChat(adminUserId, openid, robot, appType, messages, functionTools, temperature, maxToken)
	if err == nil && handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, chatResp.PromptToken, chatResp.CompletionToken, robot)
	}
	return chatResp, requestTime, err
}

func (h *ModelCallHandler) GetVector2000(adminUserId int, openid string, robot msql.Params, library msql.Params, fileInfo msql.Params, input string) (adaptor.ZhimaEmbeddingResponse, error) {
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
	err = LlmLogRequest("Text Embedding", adminUserId, openid, robot, library, h.config, lib_define.AppYunH5, fileInfo, h.Meta.Model, res.PromptToken, res.CompletionToken, req, res)
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

func (h *ModelCallHandler) RequestRerank(adminUserId int, openid, appType string, robot msql.Params, params *adaptor.ZhimaRerankReq) (adaptor.ZhimaRerankResp, error) {
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
	err = LlmLogRequest("RERANK", adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
	if err != nil {
		logs.Error(err.Error())
	}
	//}()
	return res, nil
}

func (h *ModelCallHandler) RequestChatStream(
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
	err = LlmLogRequest("LLM", adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
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
	err = LlmLogRequest("LLM", adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, resp.PromptToken, resp.CompletionToken, req, resp)
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
	modelInfo, exist := GetModelInfoByConfig(userId, modelConfigId)
	if !exist {
		return false
	}
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) //兼容旧数据
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

func CheckSupportFuncCall(adminUserId, modelConfigId int, useModel string) error {
	modelInfo, exist := GetModelInfoByConfig(adminUserId, modelConfigId)
	if !exist {
		return errors.New(`模型配置ID参数错误`)
	}
	useModel = CompatibleUseModelOldData(modelInfo.ConfigInfo, useModel) //兼容旧数据
	if !tool.InArrayString(useModel, modelInfo.GetLlmModelList()) {
		return errors.New(`使用模型名称参数错误`)
	}
	if !tool.InArrayString(useModel, modelInfo.GetFunctionCallModels()) {
		return errors.New(`使用模型不支持func call`)
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
			modelInfo, ok := GetModelInfoByConfig(adminUserId, cast.ToInt(config[`id`]))
			if !ok {
				continue
			}
			//过滤掉非当前检索的模型列表
			useModels := make([]UseModelConfig, 0)
			for _, useModel := range modelInfo.UseModelConfigs {
				if useModel.ModelType == modelType {
					useModels = append(useModels, useModel)
				}
			}
			if len(useModels) == 0 {
				continue //过滤掉空数据模型服务商
			}
			modelInfo.UseModelConfigs = useModels
			list = append(list, modelInfo)
		}
	}
	return list, nil
}
