// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"chatwiki/internal/app/chatwiki/llm/adaptor"
	"errors"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

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
)

const (
	Llm           = `LLM`
	TextEmbedding = `TEXT EMBEDDING`
	Speech2Text   = `SPEECH2TEXT`
	Tts           = `TTS`
	Rerank        = "RERANK"
	MaxContent    = 5000
)

type ModelCallHandler struct {
	adaptor.Meta
	UseModel string
}

func (h *ModelCallHandler) GetVector2000(input string) (string, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaEmbeddingRequest{Input: input}
	var res adaptor.ZhimaEmbeddingResponse
	var err error
	maxTryCount := 2
	for i := 0; i < maxTryCount; i++ {
		res, err = client.CreateEmbeddings(req)
		if err != nil {
			logs.Error(err.Error())
		} else {
			break
		}
	}
	if err != nil {
		return ``, err
	}

	if res.Result == nil {
		return ``, errors.New(`get vector return nil`)
	}
	if len(res.Result) < VectorDimension {
		res.Result = append(res.Result, make([]float64, VectorDimension-len(res.Result))...)
	}
	return tool.JsonEncode(res.Result)
}

func (h *ModelCallHandler) GetSimilarity(query []float64, inputs [][]float64) (string, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaSimilarityRequest{Model: h.UseModel, Query: query, Input: inputs}
	res, err := client.CreateSimilarity(req)
	if err != nil {
		return ``, err
	}
	if res.Result == nil {
		return ``, errors.New(`get vector return nil`)
	}
	return tool.JsonEncode(res.Result)
}

func (h *ModelCallHandler) RequestRerank(params *adaptor.ZhimaRerankReq) ([]msql.Params, error) {
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
		return nil, err
	}
	if res == nil {
		return nil, errors.New(`get rerank return nil`)
	}
	return res, nil
}

func (h *ModelCallHandler) RequestChatStream(messages []adaptor.ZhimaChatCompletionMessage, chanStream chan sse.Event, temperature float32, maxToken int) (string, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaChatCompletionRequest{
		Messages:    messages,
		MaxToken:    maxToken,
		Temperature: float64(temperature),
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		return ``, err
	}
	defer func(stream *adaptor.ZhimaChatCompletionStreamResponse) {
		_ = stream.Close()
	}(stream)
	var content string
	for {
		response, err := stream.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return ``, err
		}
		if len(response.Result) == 0 {
			continue
		}
		content += response.Result
		chanStream <- sse.Event{Event: `sending`, Data: response.Result}
	}
	return content, nil
}

func (h *ModelCallHandler) RequestChat(messages []adaptor.ZhimaChatCompletionMessage, temperature float32, maxToken int) (string, error) {
	client := &adaptor.Adaptor{}
	client.Init(h.Meta)
	req := adaptor.ZhimaChatCompletionRequest{
		Messages:    messages,
		MaxToken:    maxToken,
		Temperature: float64(temperature),
	}
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		return ``, err
	}
	return resp.Result, nil
}

type HandlerFunc func(config msql.Params, useModel string) (*ModelCallHandler, error)

type ModelInfo struct {
	ModelDefine     string        `json:"model_define"`
	ModelName       string        `json:"model_name"`
	ModelIconUrl    string        `json:"model_icon_url"`
	Introduce       string        `json:"introduce"`
	IsOffline       bool          `json:"is_offline"`
	SupportList     []string      `json:"support_list"`
	SupportedType   []string      `json:"supported_type"`
	ConfigParams    []string      `json:"config_params"`
	ConfigList      []msql.Params `json:"config_list"`
	ApiVersions     []string      `json:"api_versions"`
	LlmModelList    []string      `json:"llm_model_list"`
	VectorModelList []string      `json:"vector_model_list"`
	RerankModelList []string      `json:"rerank_model_list"`
	HelpLinks       string        `json:"help_links"`
	CallHandlerFunc HandlerFunc   `json:"-"`
}

var ModelList = []ModelInfo{
	{
		ModelDefine:   ModelOpenAI,
		ModelName:     `OpenAI`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
		Introduce:     `基于OpenAI官方提供的API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`gpt-4o`,
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
			`gpt-3.5-turbo-instruct`,
		},
		VectorModelList: []string{
			`text-embedding-3-large`,
			`text-embedding-3-small`,
			`text-embedding-ada-002`,
		},
		RerankModelList: []string{},
		HelpLinks:       `https://openai.com/`,
		CallHandlerFunc: GetOpenAIHandle,
	},
	{
		ModelDefine:     ModelOpenAIAgent,
		ModelName:       `其他兼容OpenAI API的模型服务商`,
		ModelIconUrl:    LocalUploadPrefix + `model_icon/` + ModelOpenAI + `.png`,
		Introduce:       `支持添加其他兼容OpenAi API的模型服务商，比如api2d、oneapi等`,
		IsOffline:       false,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`model_type`, `deployment_name`, `api_endpoint`, `api_key`, `api_version`},
		ConfigList:      nil,
		ApiVersions:     []string{"v1"},
		LlmModelList:    []string{"默认"},
		VectorModelList: []string{"默认"},
		RerankModelList: []string{},
		HelpLinks:       `https://openai.com/`,
		CallHandlerFunc: GetOpenAIAgentHandle,
	},
	{
		ModelDefine:   ModelAzureOpenAI,
		ModelName:     `Azure OpenAI Service`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelAzureOpenAI + `.png`,
		Introduce:     `Microsoft Azure提供的OpenAI API服务`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding, Speech2Text, Tts},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`model_type`, `deployment_name`, `api_endpoint`, `api_key`, `api_version`},
		ConfigList:    nil,
		ApiVersions: []string{
			`2023-03-15-preview`,
			`2023-05-15`,
			`2023-06-01-preview`,
			`2023-07-01-preview`,
			`2023-08-01-preview`,
			`2023-09-01-preview`,
			`2023-10-01-preview`,
			`2023-11-01-preview`,
			`2023-12-01-preview`,
			`2024-02-15-preview`,
			`2024-03-01-preview`,
			`2024-04-01-preview`,
			`2024-02-01`,
		},
		LlmModelList:    []string{`默认`},
		VectorModelList: []string{`默认`},
		RerankModelList: []string{},
		HelpLinks:       `https://azure.microsoft.com/en-us/products/ai-services/openai-service`,
		CallHandlerFunc: GetAzureHandler,
	},
	{
		ModelDefine:   ModelAnthropicClaude,
		ModelName:     `Anthropic Claude`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelAnthropicClaude + `.png`,
		Introduce:     `Anthropic出品的Claude模型`,
		IsOffline:     false,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`, `api_version`},
		ConfigList:    nil,
		ApiVersions:   []string{`2023-06-01`},
		LlmModelList: []string{
			`claude-3-opus-20240229`,
			`claude-3-sonnet-20240229`,
			`claude-3-haiku-20240307`,
		},
		VectorModelList: []string{`voyage-2`, `voyage-large-2`, `voyage-code-2`},
		RerankModelList: []string{},
		HelpLinks:       `https://claude.ai/`,
		CallHandlerFunc: GetClaudeHandler,
	},
	{
		ModelDefine:   ModelGoogleGemini,
		ModelName:     `Google Gemini`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelGoogleGemini + `.png`,
		Introduce:     `基于Google提供的Gemini API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
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
		RerankModelList: []string{},
		HelpLinks:       `https://ai.google.dev/`,
		CallHandlerFunc: GetGeminiHandler,
	},
	{
		ModelDefine:   ModelBaiduYiyan,
		ModelName:     `文心一言`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelBaiduYiyan + `.png`,
		Introduce:     `基于百度千帆大模型平台提供的文心一言API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`, `secret_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`ERNIE-4.0-Turbo-8K`,
			`ERNIE-4.0-8K`,
			`ERNIE-4.0-8K-Preemptible`,
			`ERNIE-4.0-8K-Preview`,
			`ERNIE-4.0-8K-Preview-0518`,
			`ERNIE-4.0-8K-Latest`,
			`ERNIE-4.0-8K-0329`,
			`ERNIE-4.0-8K-0104`,
			`ERNIE-4.0-8K-0613`,
			`ERNIE-3.5-8K`,
			`ERNIE-3.5-8K-0205`,
			`ERNIE-3.5-8K-Preview`,
			`ERNIE-3.5-8K-0329`,
			`ERNIE-3.5-128K`,
			`ERNIE-3.5-8K-0613`,
			`ERNIE-Speed-8K`,
			`ERNIE-Speed-128K`,
			`ERNIE-Lite-8K-0922`,
			`ERNIE-Lite-8K-0308`,
		},
		VectorModelList: []string{
			`Embedding-V1`,
			`bge-large-zh`,
			`bge-large-en`,
			`tao-8k`,
		},
		RerankModelList: []string{},
		HelpLinks:       `https://cloud.baidu.com/`,
		CallHandlerFunc: GetYiyanHandler,
	},
	{
		ModelDefine:   ModelAliyunTongyi,
		ModelName:     `通义千问`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelAliyunTongyi + `.png`,
		Introduce:     `基于阿里云提供的通义千问API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding, Tts},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`qwen-turbo`,
			`qwen-plus`,
			`qwen-max`,
			`qwen-max-0428`,
			`qwen-max-0403`,
			`qwen-max-0107`,
			`qwen-max-longcontext`,
		},
		VectorModelList: []string{
			`text-embedding-v1`,
			`text-embedding-v2`,
		},
		RerankModelList: []string{},
		HelpLinks:       `https://dashscope.aliyun.com/?spm=a2c4g.11186623.nav-dropdown-menu-0.142.6d1b46c1EeV28g&scm=20140722.X_data-37f0c4e3bf04683d35bc._.V_1`,
		CallHandlerFunc: GetTongyiHandler,
	},
	{
		ModelDefine:   ModelBaai,
		ModelName:     `BAAI 智源研究院`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelBaai + `.png`,
		Introduce:     `由北京智源人工智能研究院研发的本地模型，包含bge-rerank-base、bge-m3模型，支持嵌入和rerank。使用bge系列模型，无需消耗token，但是本地模型运行需要硬件支持，请确保服务器有足够的内存（至少8G内存）和用于计算的GPU`,
		IsOffline:     true,
		SupportList:   []string{TextEmbedding, Rerank},
		SupportedType: []string{TextEmbedding, Rerank},
		ConfigParams:  []string{`api_endpoint`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList:  []string{},
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
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelCohere + `.png`,
		Introduce:     `cohere提供的模型，包含Command、Command R、Command R+等`,
		IsOffline:     false,
		SupportList:   []string{Rerank},
		SupportedType: []string{Rerank},
		ConfigParams:  []string{`api_endpoint`},
		ConfigList:    nil,
		ApiVersions:   []string{},
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
		ModelIconUrl:    LocalUploadPrefix + `model_icon/` + ModelOllama + `.png`,
		Introduce:       `Ollama是一个轻量级的简单易用的本地大模型运行框架,通过Ollama可以在本地服务器构建和运营大语言模型(比如Llama3等).ChatWiki支持使用Ollama部署LLM的型和Text Embedding模型`,
		IsOffline:       true,
		SupportList:     []string{Llm, TextEmbedding},
		SupportedType:   []string{Llm, TextEmbedding},
		ConfigParams:    []string{`model_type`, `deployment_name`, `api_endpoint`},
		ConfigList:      nil,
		ApiVersions:     []string{},
		LlmModelList:    []string{"默认"},
		VectorModelList: []string{"默认"},
		RerankModelList: []string{},
		HelpLinks:       `https://www.ollama.com/`,
		CallHandlerFunc: GetOllamaHandle,
	},
	//{
	//	ModelDefine:     ModelXnference,
	//	ModelName:       `xorbitsai inference`,
	//	ModelIconUrl:    LocalUploadPrefix + `model_icon/` + ModelXnference + `.png`,
	//	Introduce:       `Xorbits Inference(Xinference)是一个开源平台,用于简化各种AI模型的运行和集成,借助Xinference,您可以使用任何开源LLM,嵌入模型和多模态模型在本地服务器中部署`,
	//	IsOffline:       true,
	//	SupportList:     []string{Llm, TextEmbedding, Rerank},
	//	SupportedType:   []string{Llm, TextEmbedding, Rerank},
	//	ConfigParams:    []string{`model_type`, `deployment_name`, `api_version`, `api_endpoint`},
	//	ConfigList:      nil,
	//	ApiVersions:     []string{"v1"},
	//	LlmModelList:    []string{"默认"},
	//	VectorModelList: []string{"默认"},
	//	RerankModelList: []string{"默认"},
	//	HelpLinks:       `https://baidu.com/`,
	//	CallHandlerFunc: GetXinferenceHandle,
	//},
	{
		ModelDefine:   ModelDeepseek,
		ModelName:     `DeepSeek`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelDeepseek + `.png`,
		Introduce:     `由DeepSeek提供的大模型API`,
		IsOffline:     false,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`deepseek-chat`,
		},
		VectorModelList: []string{},
		RerankModelList: []string{},
		HelpLinks:       `https://www.deepseek.com/`,
		CallHandlerFunc: GetDeepseekHandle,
	},
	{
		ModelDefine:   ModelJina,
		ModelName:     `Jina`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelJina + `.png`,
		Introduce:     `有Jina提供的嵌入和Rerank模型，`,
		IsOffline:     false,
		SupportList:   []string{Llm, Rerank},
		SupportedType: []string{Llm, Rerank},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`jina-embeddings-v2-base-en`,
			`jina-embeddings-v2-base-zh`,
			`jina-embeddings-v2-base-de`,
			`jina-embeddings-v2-base-es`,
			`jina-colbert-v1-en`,
			`jina-embeddings-v2-base-code`,
		},
		VectorModelList: []string{},
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
		ModelDefine:   ModelLingYiWanWu,
		ModelName:     `零一万物`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelLingYiWanWu + `.png`,
		Introduce:     `基于零一万物提供的零一大模型API`,
		IsOffline:     false,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`yi-large`,
		},
		VectorModelList: []string{},
		RerankModelList: []string{},
		HelpLinks:       `https://platform.lingyiwanwu.com/`,
		CallHandlerFunc: GetLingYiWanWuHandle,
	},
	{
		ModelDefine:   ModelMoonShot,
		ModelName:     `月之暗面`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelMoonShot + `.png`,
		Introduce:     `基于月之暗面提供的Kimi API`,
		IsOffline:     false,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`moonshot-v1-8k`,
		},
		VectorModelList: []string{},
		RerankModelList: []string{},
		HelpLinks:       `https://www.moonshot.cn/`,
		CallHandlerFunc: GetMoonShotHandle,
	},
	{
		ModelDefine:   ModelSpark,
		ModelName:     `讯飞星火`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelSpark + `.png`,
		Introduce:     `基于科大讯飞提供的讯飞星火大模型API`,
		IsOffline:     false,
		SupportList:   []string{Llm},
		SupportedType: []string{Llm},
		ConfigParams:  []string{`api_key`, `app_id`, `secret_key`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`generalv3.5`,
			`generalv3`,
			`generalv2`,
			`general`,
		},
		VectorModelList: []string{},
		RerankModelList: []string{},
		HelpLinks:       `https://xinghuo.xfyun.cn/sparkapi`,
		CallHandlerFunc: GetSparkHandle,
	},

	{
		ModelDefine:   ModelHunyuan,
		ModelName:     `腾讯混元`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelHunyuan + `.png`,
		Introduce:     `hunyuan`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`api_key`, `secret_key`, `region`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`hunyuan-lite`,
			`hunyuan-standard`,
			`hunyuan-standard-256K`,
			`hunyuan-pro`,
		},
		VectorModelList: []string{
			`默认`,
		},
		RerankModelList: []string{},
		HelpLinks:       `https://cloud.tencent.com/product/hunyuan`,
		CallHandlerFunc: GetHunyuanHandle,
	},
	{
		ModelDefine:   ModelDoubao,
		ModelName:     `火山引擎`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelDoubao + `.png`,
		Introduce:     `基于火山引擎提供的豆包大模型API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`deployment_name`, `api_key`, `secret_key`, `region`},
		ConfigList:    nil,
		ApiVersions:   []string{},
		LlmModelList: []string{
			`Doubao-lite-4k`,
			`Doubao-lite-32k`,
			`Doubao-lite-128k`,
			`Doubao-pro-4k`,
			`Doubao-pro-32k`,
			`Doubao-pro-128k`,
		},
		VectorModelList: []string{
			`默认`,
		},
		RerankModelList: []string{},
		HelpLinks:       `https://www.volcengine.com/product/doubao`,
		CallHandlerFunc: GetDoubaoHandle,
	},
	{
		ModelDefine:   ModelBaichuan,
		ModelName:     `百川智能`,
		ModelIconUrl:  LocalUploadPrefix + `model_icon/` + ModelBaichuan + `.png`,
		Introduce:     `基于百川智能提供的百川大模型API`,
		IsOffline:     false,
		SupportList:   []string{Llm, TextEmbedding},
		SupportedType: []string{Llm, TextEmbedding},
		ConfigParams:  []string{`deployment_name`, `api_key`, `secret_key`, `region`},
		ConfigList:    nil,
		ApiVersions:   []string{},
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
		RerankModelList: []string{},
		HelpLinks:       `https://platform.baichuan-ai.com`,
		CallHandlerFunc: GetDoubaoHandle,
	},
}
