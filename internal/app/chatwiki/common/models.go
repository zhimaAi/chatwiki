// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"errors"
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
	ModelDefine               string        `json:"model_define"`
	ModelName                 string        `json:"model_name"`
	ModelIconUrl              string        `json:"model_icon_url"`
	Introduce                 string        `json:"introduce"`
	IsOffline                 bool          `json:"is_offline"`
	SupportList               []string      `json:"support_list"`
	SupportedType             []string      `json:"supported_type"`
	SupportedFunctionCallList []string      `json:"supported_function_call_list"`
	ConfigParams              []string      `json:"config_params"`
	HistoryConfigParams       []string      `json:"history_config_params"`
	ConfigList                []msql.Params `json:"config_list"`
	ApiVersions               []string      `json:"api_versions"`
	LlmModelList              []string      `json:"llm_model_list"`
	ChoosableThinkingModels   []string      `json:"choosable_thinking_models"`
	NetworkSearchModelList    []string      `json:"network_search_model_list"`
	VectorModelList           []string      `json:"vector_model_list"`
	RerankModelList           []string      `json:"rerank_model_list"`
	HelpLinks                 string        `json:"help_links"`
	CallHandlerFunc           HandlerFunc   `json:"-"`
	CheckAllowRequest         BeforeFunc    `json:"-"`
	CheckFancCallRequest      BeforeFunc    `json:"-"`
	TokenUseReport            AfterFunc     `json:"-"`
}

const (
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
	Rerank        = "RERANK"
	MaxContent    = 10000
)

func GetModelList() []ModelInfo {
	//模型过滤处理
	list := make([]ModelInfo, 0)
	for _, info := range modelList {
		if tool.InArrayString(info.ModelDefine, []string{}) {
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
		if info.SupportedFunctionCallList == nil {
			list[i].SupportedFunctionCallList = make([]string, 0)
		}
		if info.ConfigParams == nil {
			list[i].ConfigParams = make([]string, 0)
		}
		if info.HistoryConfigParams == nil {
			list[i].HistoryConfigParams = make([]string, 0)
		}
		if info.ConfigList == nil {
			list[i].ConfigList = make([]msql.Params, 0)
		}
		if info.ApiVersions == nil {
			list[i].ApiVersions = make([]string, 0)
		}
		if info.LlmModelList == nil {
			list[i].LlmModelList = make([]string, 0)
		}
		if info.ChoosableThinkingModels == nil {
			list[i].ChoosableThinkingModels = make([]string, 0)
		}
		if info.NetworkSearchModelList == nil {
			list[i].NetworkSearchModelList = make([]string, 0)
		}
		if info.VectorModelList == nil {
			list[i].VectorModelList = make([]string, 0)
		}
		if info.RerankModelList == nil {
			list[i].RerankModelList = make([]string, 0)
		}
	}
	return list
}

var modelList = [...]ModelInfo{
	{
		ModelDefine:   ModelOpenAI,
		ModelName:     `OpenAI`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
		Introduce:     `基于OpenAI官方提供的API`,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		SupportedFunctionCallList: []string{
			`gpt-4o`,
			`gpt-4o-mini`,
			`gpt-4-turbo`,
			`gpt-4-turbo-2024-04-09`,
			`gpt-4-turbo-preview`,
		},
		ConfigParams: []string{`api_key`},
		LlmModelList: []string{
			`gpt-4o`,
			`gpt-4o-mini`,
			`gpt-4-turbo`,
			`gpt-4-turbo-2024-04-09`,
			`gpt-4-turbo-preview`,
			`gpt-4-0125-preview`,
			`gpt-4-1106-preview`,
			`gpt-4-vision-preview`,
			`gpt-4-1106-vision-preview`,
			`gpt-4`,
			`gpt-4-0613`,
			`gpt-4-32k`,
			`gpt-4-32k-0613`,
			`gpt-3.5-turbo-0125`,
			`gpt-3.5-turbo`,
			`gpt-3.5-turbo-1106`,
		},
		VectorModelList: []string{
			`text-embedding-3-large`,
			`text-embedding-3-small`,
			`text-embedding-ada-002`,
		},
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
		ConfigParams:    []string{`model_type`, `deployment_name`, `api_endpoint`, `api_key`, `api_version`},
		ApiVersions:     []string{"v1", `v3`},
		LlmModelList:    []string{"默认"},
		VectorModelList: []string{"默认"},
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
		ConfigParams:  []string{`model_type`, `deployment_name`, `api_endpoint`, `api_key`, `api_version`},
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
		LlmModelList:    []string{`默认`},
		VectorModelList: []string{`默认`},
		HelpLinks:       `https://azure.microsoft.com/en-us/products/ai-services/openai-service`,
		CallHandlerFunc: GetAzureHandler,
	},
	{
		ModelDefine:   ModelAnthropicClaude,
		ModelName:     `Anthropic Claude`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAnthropicClaude + `.png`,
		Introduce:     `Anthropic出品的Claude模型`,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`, `api_version`},
		ApiVersions:   []string{`2023-06-01`},
		LlmModelList: []string{
			`claude-3-opus-20240229`,
			`claude-3-sonnet-20240229`,
			`claude-3-haiku-20240307`,
			`claude-3-5-sonnet-20240620`,
		},
		VectorModelList: []string{`voyage-2`, `voyage-large-2`, `voyage-code-2`},
		HelpLinks:       `https://claude.ai/`,
		CallHandlerFunc: GetClaudeHandler,
	},
	{
		ModelDefine:   ModelGoogleGemini,
		ModelName:     `Google Gemini`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelGoogleGemini + `.png`,
		Introduce:     `基于Google提供的Gemini API`,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`},
		LlmModelList: []string{
			`gemini-1.0-pro`,
			`gemini-1.5-flash`,
			`gemini-1.5-pro`,
			`gemini-pro`,
			`gemini-pro-vision`,
		},
		VectorModelList: []string{
			`text-embedding-004`,
			`embedding-001`,
		},
		HelpLinks:       `https://ai.google.dev/`,
		CallHandlerFunc: GetGeminiHandler,
	},
	{
		ModelDefine:   ModelBaiduYiyan,
		ModelName:     `文心一言`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelBaiduYiyan + `.png`,
		Introduce:     `基于百度千帆大模型平台提供的文心一言API`,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		SupportedFunctionCallList: []string{
			//v2 chat
			`ernie-x1-turbo-32k`,
			`deepseek-v3`,
			`deepseek-r1-250528`,
			`deepseek-r1`,
			//v1 chat
			`ERNIE-4.0-8K`,
			`ERNIE-4.0-Turbo-8K`,
			`ERNIE-3.5-8K`,
		},
		ConfigParams:        []string{`api_key`},
		HistoryConfigParams: []string{`secret_key`},
		LlmModelList: []string{
			// v2 chat
			`ernie-4.5-turbo-32k`,
			`ernie-4.5-turbo-128k`,
			`ernie-4.0-8k`,
			`ernie-speed-128k`,
			`ernie-x1-turbo-32k`,
			`deepseek-v3`,
			`deepseek-r1-250528`,
			`deepseek-r1`,
			//v1 chat
			`ERNIE-4.0-8K`,
			`ERNIE-4.0-Turbo-8K`,
			`ERNIE-4.0-8K-Preview`,
			`ERNIE-4.0-8K-Latest`,
			`ERNIE-3.5-8K`,
			`ERNIE-3.5-128K`,
		},
		NetworkSearchModelList: []string{
			`ernie-4.5-turbo-32k`,
			`ernie-4.5-turbo-128k`,
			`ernie-4.0-8k`,
			`ernie-x1-turbo-32k`,
			`deepseek-v3`,
			`deepseek-r1`,
		},
		VectorModelList: []string{
			`embedding-v1`,
			`bge-large-zh`,
			`bge-large-en`,
			`tao-8k`,
		},
		HelpLinks:            `https://cloud.baidu.com/`,
		CallHandlerFunc:      GetYiyanHandler,
		CheckFancCallRequest: CheckYiyanFancCall,
	},
	{
		ModelDefine:   ModelAliyunTongyi,
		ModelName:     `通义千问`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelAliyunTongyi + `.png`,
		Introduce:     `基于阿里云提供的通义千问API`,
		SupportList:   []string{Llm, TextEmbedding, Tts, Rerank},
		SupportedType: []string{Llm, TextEmbedding, Rerank},
		SupportedFunctionCallList: []string{
			`qwen-plus`,
			`qwen-turbo`,
			`qwen3-235b-a22b`,
			`qwen-max`,
			`qwen-long`,
			`deepseek-v3`,
			`deepseek-r1`,
			`deepseek-r1-0528`,
		},
		ConfigParams: []string{`api_key`},
		LlmModelList: []string{
			`qwen-plus`,
			`qwen-turbo`,
			`qwen3-235b-a22b`,
			`qwen-max`,
			`qwen-long`,
			`deepseek-v3`,
			`deepseek-r1`,
			`deepseek-r1-0528`,
			`Moonshot-Kimi-K2-Instruct`,
		},
		ChoosableThinkingModels: []string{
			`qwen-plus`,
			`qwen-turbo`,
			`qwen3-235b-a22b`,
		},
		NetworkSearchModelList: []string{
			`qwen-plus`,
			`qwen-turbo`,
			`qwen3-235b-a22b`,
			`qwen-max`,
			`Moonshot-Kimi-K2-Instruct`,
		},
		VectorModelList: []string{
			`text-embedding-v4`,
			`text-embedding-v3`,
			`text-embedding-v2`,
			`text-embedding-v1`,
		},
		RerankModelList: []string{
			`gte-rerank-v2`,
		},
		HelpLinks:       `https://dashscope.aliyun.com/?spm=a2c4g.11186623.nav-dropdown-menu-0.142.6d1b46c1EeV28g&scm=20140722.X_data-37f0c4e3bf04683d35bc._.V_1`,
		CallHandlerFunc: GetTongyiHandler,
	},
	{
		ModelDefine:   ModelBaai,
		ModelName:     `BGE`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelBaai + `.png`,
		Introduce:     `由北京智源人工智能研究院研发的本地模型，包含bge-rerank-base、bge-m3模型，支持嵌入和rerank。使用bge系列模型，无需消耗token，但是本地模型运行需要硬件支持，请确保服务器有足够的内存（至少8G内存）和用于计算的GPU`,
		SupportList:   []string{TextEmbedding, Rerank},
		SupportedType: []string{TextEmbedding, Rerank},
		ConfigParams:  []string{`api_endpoint`},
		VectorModelList: []string{
			"bge-m3",
		},
		RerankModelList: []string{
			`bge-reranker-base-onnx-o4`,
			"bge-m3",
		},
		HelpLinks:       `https://www.baidu.com/`,
		CallHandlerFunc: GetBaaiHandle,
	},
	{
		ModelDefine:   ModelCohere,
		ModelName:     `Cohere`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelCohere + `.png`,
		Introduce:     `cohere提供的模型，包含Command、Command R、Command R+等`,
		SupportList:   []string{Llm, TextEmbedding, Rerank},
		SupportedType: []string{Llm, TextEmbedding, Rerank},
		ConfigParams:  []string{`api_key`},
		LlmModelList: []string{
			`command-r-plus`,
			`command-r`,
			`command`,
			`command-nightly`,
			`command-light`,
			`command-light-nightly`,
		},
		VectorModelList: []string{
			`embed-english-v3.0`,
			`embed-english-light-v3.0`,
			`embed-multilingual-v3.0`,
			`embed-multilingual-light-v3.0`,
			`embed-english-v2.0`,
			`embed-english-light-v2.0`,
			`embed-multilingual-v2.0`,
		},
		RerankModelList: []string{
			"rerank-english-v3.0",
			"rerank-multilingual-v3.0",
			"rerank-english-v2.0",
			"rerank-multilingual-v2.0",
		},
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
		ConfigParams:    []string{`model_type`, `deployment_name`, `api_endpoint`},
		LlmModelList:    []string{"默认"},
		VectorModelList: []string{"默认"},
		HelpLinks:       `https://www.ollama.com/`,
		CallHandlerFunc: GetOllamaHandle,
	},
	{
		ModelDefine:     ModelXnference,
		ModelName:       `xorbitsai inference`,
		ModelIconUrl:    define.LocalUploadPrefix + `model_icon/` + ModelXnference + `.png`,
		Introduce:       `Xorbits Inference(Xinference)是一个开源平台,用于简化各种AI模型的运行和集成,借助Xinference,您可以使用任何开源LLM,嵌入模型和多模态模型在本地服务器中部署`,
		IsOffline:       true,
		SupportList:     []string{Llm, TextEmbedding, Rerank},
		SupportedType:   []string{Llm, TextEmbedding, Rerank},
		ConfigParams:    []string{`model_type`, `deployment_name`, `api_version`, `api_endpoint`},
		ApiVersions:     []string{"v1"},
		LlmModelList:    []string{"默认"},
		VectorModelList: []string{"默认"},
		RerankModelList: []string{"默认"},
		HelpLinks:       `https://baidu.com/`,
		CallHandlerFunc: GetXinferenceHandle,
	},
	{
		ModelDefine:               ModelDeepseek,
		ModelName:                 `DeepSeek`,
		ModelIconUrl:              define.LocalUploadPrefix + `model_icon/` + ModelDeepseek + `.png`,
		Introduce:                 `由DeepSeek提供的大模型API`,
		SupportList:               []string{Llm},
		SupportedType:             []string{Llm},
		SupportedFunctionCallList: []string{`deepseek-chat`},
		ConfigParams:              []string{`api_key`},
		LlmModelList: []string{
			`deepseek-chat`,
			`deepseek-reasoner`,
		},
		HelpLinks:       `https://www.deepseek.com/`,
		CallHandlerFunc: GetDeepseekHandle,
	},
	{
		ModelDefine:   ModelJina,
		ModelName:     `Jina`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelJina + `.png`,
		Introduce:     `有Jina提供的嵌入和Rerank模型，`,
		SupportList:   []string{TextEmbedding, Rerank},
		SupportedType: []string{TextEmbedding, Rerank},
		ConfigParams:  []string{`api_key`},
		VectorModelList: []string{
			`jina-embeddings-v2-base-en`,
			`jina-embeddings-v2-base-zh`,
			`jina-embeddings-v2-base-de`,
			`jina-embeddings-v2-base-es`,
			`jina-colbert-v1-en`,
			`jina-embeddings-v2-base-code`,
		},
		RerankModelList: []string{
			"jina-reranker-v1-base-en",
			"jina-reranker-v1-turbo-en",
			"jina-reranker-v1-tiny-en",
			"jina-colbert-v1-en",
		},
		HelpLinks:       `https://jina.ai/`,
		CallHandlerFunc: GetJinaHandle,
	},
	{
		ModelDefine:               ModelLingYiWanWu,
		ModelName:                 `零一万物`,
		ModelIconUrl:              define.LocalUploadPrefix + `model_icon/` + ModelLingYiWanWu + `.png`,
		Introduce:                 `基于零一万物提供的零一大模型API`,
		SupportList:               []string{Llm},
		SupportedType:             []string{Llm},
		SupportedFunctionCallList: []string{`yi-large-fc`},
		ConfigParams:              []string{`api_key`},
		LlmModelList: []string{
			`yi-large`,
			`yi-large-fc`,
		},
		HelpLinks:       `https://platform.lingyiwanwu.com/`,
		CallHandlerFunc: GetLingYiWanWuHandle,
	},
	{
		ModelDefine:               ModelMoonShot,
		ModelName:                 `月之暗面`,
		ModelIconUrl:              define.LocalUploadPrefix + `model_icon/` + ModelMoonShot + `.png`,
		Introduce:                 `基于月之暗面提供的Kimi API`,
		SupportList:               []string{Llm},
		SupportedType:             []string{Llm},
		SupportedFunctionCallList: []string{`moonshot-v1-8k`},
		ConfigParams:              []string{`api_key`},
		LlmModelList: []string{
			`moonshot-v1-8k`,
		},
		HelpLinks:       `https://www.moonshot.cn/`,
		CallHandlerFunc: GetMoonShotHandle,
	},
	{
		ModelDefine:   ModelSpark,
		ModelName:     `讯飞星火`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelSpark + `.png`,
		Introduce:     `基于科大讯飞提供的讯飞星火大模型API`,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`app_id`, `api_key`, `secret_key`},
		LlmModelList: []string{
			`Spark Lite`,
			`Spark V2.0`,
			`Spark Pro`,
			`Spark Max`,
			`Spark4.0 Ultra`,
		},
		HelpLinks:       `https://xinghuo.xfyun.cn/sparkapi`,
		CallHandlerFunc: GetSparkHandle,
	},
	{
		ModelDefine:   ModelHunyuan,
		ModelName:     `腾讯混元`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelHunyuan + `.png`,
		Introduce:     `hunyuan`,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`, `secret_key`},
		LlmModelList: []string{
			`hunyuan-lite`,
			`hunyuan-functioncall`,
			`hunyuan-standard`,
			`hunyuan-standard-256K`,
			`hunyuan-pro`,
		},
		VectorModelList: []string{
			`默认`,
		},
		HelpLinks:       `https://cloud.tencent.com/product/hunyuan`,
		CallHandlerFunc: GetHunyuanHandle,
	},
	{
		ModelDefine:   ModelDoubao,
		ModelName:     `火山引擎`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelDoubao + `.png`,
		Introduce:     `基于火山引擎提供的豆包大模型API`,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`model_type`, `deployment_name`, `show_model_name`, `api_key`, `thinking_type`, `region`},
		HistoryConfigParams: []string{
			`secret_key`,
		},
		LlmModelList:    []string{`默认`},
		VectorModelList: []string{`默认`},
		HelpLinks:       `https://www.volcengine.com/product/doubao`,
		CallHandlerFunc: GetDoubaoHandle,
	},
	{
		ModelDefine:               ModelBaichuan,
		ModelName:                 `百川智能`,
		ModelIconUrl:              define.LocalUploadPrefix + `model_icon/` + ModelBaichuan + `.png`,
		Introduce:                 `基于百川智能提供的百川大模型API`,
		SupportList:               []string{Llm, TextEmbedding},
		SupportedType:             []string{Llm, TextEmbedding},
		SupportedFunctionCallList: []string{`Baichuan4`, `Baichuan3-Turbo`, `Baichuan3-Turbo-128k`, `Baichuan2-Turbo`, `Baichuan2-Turbo-192k`},
		ConfigParams:              []string{`api_key`},
		LlmModelList: []string{
			`Baichuan4`,
			`Baichuan3-Turbo`,
			`Baichuan3-Turbo-128k`,
			`Baichuan2-Turbo`,
			`Baichuan2-Turbo-192k`,
		},
		VectorModelList: []string{
			`Baichuan-Text-Embedding`,
		},
		HelpLinks:       `https://platform.baichuan-ai.com`,
		CallHandlerFunc: GetBaichuanHandle,
	},
	{
		ModelDefine:               ModelZhipu,
		ModelName:                 `智谱`,
		ModelIconUrl:              define.LocalUploadPrefix + `model_icon/` + ModelZhipu + `.png`,
		Introduce:                 `领先的认知大模型AI开放平台`,
		SupportList:               []string{Llm, TextEmbedding},
		SupportedType:             []string{Llm, TextEmbedding},
		SupportedFunctionCallList: []string{`glm-4-0520`, `glm-4`, `glm-4-air`, `glm-4-airx`, `glm-4-flash`},
		ConfigParams:              []string{`api_key`},
		LlmModelList: []string{
			`glm-4-0520`,
			`glm-4`,
			`glm-4-air`,
			`glm-4-airx`,
			`glm-4-flash`,
		},
		VectorModelList: []string{
			`embedding-2`,
		},
		HelpLinks:       `https://open.bigmodel.cn/`,
		CallHandlerFunc: GetZhipuHandle,
	},
	{
		ModelDefine:   ModelMinimax,
		ModelName:     `minimax`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelMinimax + `.png`,
		Introduce:     `MiniMax 成立于 2021 年 12 月，是领先的通用人工智能科技公司，致力于与用户共创智能。MiniMax 自主研发多模态、万亿参数的 MoE 大模型，并基于大模型推出海螺AI、星野等原生应用。MiniMax API 开放平台提供安全、灵活、可靠的 API 服务，助力企业和开发者快速搭建 AI 应用。`,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`},
		LlmModelList: []string{
			`abab6.5s-chat`,
			`abab6.5g-chat`,
			`abab6.5t-chat`,
			`abab5.5s-chat`,
			`abab5.5-chat`,
		},
		HelpLinks:       `https://www.minimaxi.com/`,
		CallHandlerFunc: GetMinimaxHandle,
	},
	{
		ModelDefine:   ModelSiliconFlow,
		ModelName:     `硅基流动`,
		ModelIconUrl:  define.LocalUploadPrefix + `model_icon/` + ModelSiliconFlow + `.png`,
		Introduce:     `支持通义千问，mata-lama，google-gemma，bge-m3等开源模型，可以免部署、低成本使用`,
		SupportList:   []string{Llm, TextEmbedding, Rerank},
		SupportedType: []string{Llm, TextEmbedding, Rerank},
		SupportedFunctionCallList: []string{
			`deepseek-ai/DeepSeek-R1`,
			`Pro/deepseek-ai/DeepSeek-R1`,
			`deepseek-ai/DeepSeek-V3`,
			`Pro/deepseek-ai/DeepSeek-V3`,
			`Pro/moonshotai/Kimi-K2-Instruct`,
			`Qwen/Qwen3-235B-A22B`,
		},
		ConfigParams: []string{`api_key`},
		LlmModelList: []string{
			`deepseek-ai/DeepSeek-R1`,
			`Pro/deepseek-ai/DeepSeek-R1`,
			`deepseek-ai/DeepSeek-V3`,
			`Pro/deepseek-ai/DeepSeek-V3`,
			`Pro/moonshotai/Kimi-K2-Instruct`,
			`baidu/ERNIE-4.5-300B-A47B`,
			`Qwen/Qwen3-235B-A22B`,
		},
		ChoosableThinkingModels: []string{
			`Qwen/Qwen3-235B-A22B`,
		},
		VectorModelList: []string{
			`Qwen/Qwen3-Embedding-8B`,
			`Qwen/Qwen3-Embedding-4B`,
			`Qwen/Qwen3-Embedding-0.6B`,
			`BAAI/bge-m3`,
			`Pro/BAAI/bge-m3`,
		},
		RerankModelList: []string{
			`Qwen/Qwen3-Reranker-8B`,
			`Qwen/Qwen3-Reranker-4B`,
			`Qwen/Qwen3-Reranker-0.6B`,
			`BAAI/bge-reranker-v2-m3`,
			`Pro/BAAI/bge-reranker-v2-m3`,
		},
		HelpLinks:       `https://siliconflow.cn/zh-cn/`,
		CallHandlerFunc: GetSiliconFlow,
	},
}

func GetModelInfoByDefine(modelDefine string) (ModelInfo, bool) {
	for _, info := range GetModelList() {
		if info.ModelDefine == modelDefine {
			return info, true
		}
	}
	return ModelInfo{}, false
}

func IsMultiConfModel(defineName string) bool {
	return tool.InArrayString(defineName, []string{ModelAzureOpenAI, ModelOllama, ModelXnference, ModelOpenAIAgent, ModelDoubao})
}

func GetModelCallHandler(modelConfigId int, useModel string) (*ModelCallHandler, error) {
	if modelConfigId <= 0 {
		return nil, errors.New(`model config id is empty`)
	}
	config, err := GetModelConfigInfo(modelConfigId, 0)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	if len(config) == 0 {
		return nil, errors.New(`model config is empty`)
	}
	modelInfo, ok := GetModelInfoByDefine(config[`model_define`])
	if !ok {
		return nil, errors.New(`model define invalid`)
	}
	if modelInfo.CheckAllowRequest != nil { //check allow request
		if err = modelInfo.CheckAllowRequest(modelInfo, config, useModel); err != nil {
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
	handler, err := GetModelCallHandler(modelConfigId, useModel)
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
	handler, err := GetModelCallHandler(modelConfigId, useModel)
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
	handler, err := GetModelCallHandler(modelConfigId, useModel)
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
	handler, err := GetModelCallHandler(modelConfigId, useModel)
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
	go func() {
		err := LlmLogRequest("Text Embedding", adminUserId, openid, robot, library, h.config, lib_define.AppYunH5, fileInfo, h.Meta.Model, res.PromptToken, res.CompletionToken, req, res)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
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
	go func() {
		err := LlmLogRequest("RERANK", adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
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

	go func() {
		library := msql.Params{}
		if appType == "" && openid == "" {
			library, robot = robot, library
		}
		err = LlmLogRequest("LLM", adminUserId, openid, robot, library, h.config, appType, msql.Params{}, h.Meta.Model, totalResponse.PromptToken, totalResponse.CompletionToken, req, totalResponse)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
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
	go func() {
		err := LlmLogRequest("LLM", adminUserId, openid, robot, msql.Params{}, h.config, appType, msql.Params{}, h.Meta.Model, resp.PromptToken, resp.CompletionToken, req, resp)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
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
	config, err := GetModelConfigInfo(modelConfigId, userId)
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	modelInfo, _ := GetModelInfoByDefine(config[`model_define`])
	modelList := modelInfo.LlmModelList
	if modelType == TextEmbedding {
		modelList = modelInfo.VectorModelList
	}
	if !tool.InArrayString(useModel, modelList) && !IsMultiConfModel(config["model_define"]) {
		return false
	}
	return true
}

func CheckModelIsDeepSeek(model string) bool {
	modelLower := strings.ToLower(model)
	return strings.Contains(modelLower, `deepseek-r1`) ||
		strings.Contains(modelLower, `deepseek-reasoner`)
}

func CheckSupportFuncCall(adminUserId, modelConfigId int, useModel string) error {
	config, err := GetModelConfigInfo(modelConfigId, adminUserId)
	if err != nil {
		return err
	}
	if len(config) == 0 || !tool.InArrayString(Llm, strings.Split(config[`model_types`], `,`)) {
		return errors.New(`模型配置ID参数错误`)
	}
	modelInfo, _ := GetModelInfoByDefine(config[`model_define`])
	if !tool.InArrayString(useModel, modelInfo.LlmModelList) && !IsMultiConfModel(config["model_define"]) {
		return errors.New(`使用模型名称参数错误`)
	}
	if len(modelInfo.SupportedFunctionCallList) == 0 {
		return errors.New(`模型服务商不支持func call`)
	}
	if modelInfo.CheckFancCallRequest != nil {
		if err = modelInfo.CheckFancCallRequest(modelInfo, config, useModel); err != nil {
			return errors.New(`使用模型不支持func call`)
		}
	} else if !tool.InArrayString(useModel, modelInfo.SupportedFunctionCallList) {
		return errors.New(`使用模型不支持func call`)
	}
	return nil
}
