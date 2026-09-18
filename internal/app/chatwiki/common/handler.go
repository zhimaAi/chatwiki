// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	llm "github.com/zhimaAi/llm_adaptor/v2"
)

const (
	legacyOpenAIEmbeddingModelLarge = "text-embedding-3-large"
	legacyAliEmbeddingModelV3       = "text-embedding-v3"
	legacyEmbeddingDimension1024    = 1024
	legacyEmbeddingDimension1536    = 1536
)

func legacyEmbeddingDimensions(provider llm.Provider, useModel string) *int {
	switch provider {
	case llm.ProviderAli:
		if useModel == legacyAliEmbeddingModelV3 {
			return tea.Int(legacyEmbeddingDimension1024)
		}
		return tea.Int(legacyEmbeddingDimension1536)
	case llm.ProviderOpenAI, llm.ProviderOpenAIAgent, llm.ProviderBaichuan, llm.ProviderZhipu:
		if useModel == legacyOpenAIEmbeddingModelLarge {
			return tea.Int(legacyEmbeddingDimension1536)
		}
	}
	return nil
}

func newModelCallHandler(modelInfo ModelInfo, config msql.Params, useModel string, clientConfig llm.ClientConfig) (*ModelCallHandler, error) {
	client, err := llm.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	return &ModelCallHandler{
		Client:              client,
		Model:               useModel,
		EmbeddingDimensions: legacyEmbeddingDimensions(clientConfig.Provider, useModel),
		ChoosableThinking:   tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		config:              config,
	}, nil
}

func Get302AiHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.Provider302AI,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetOpenRouterHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenRouter,
		BaseURL:  ResolveOpenRouterEndpoint(),
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetOrcaRouterHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenAI,
		BaseURL:  ResolveOrcaRouterEndpoint(),
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetDeepseekHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderDeepSeek,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetGeminiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderGemini,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetOpenAIHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenAI,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetDoubaoHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderDoubao,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetSiliconFlowHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderSiliconFlow,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetTongyiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderAli,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetOpenAIAgentHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider:   llm.ProviderOpenAIAgent,
		BaseURL:    config[`api_endpoint`],
		APIVersion: config[`api_version`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetAzureHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderAzure,
		BaseURL:  config[`api_endpoint`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetClaudeHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderClaude,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetYiyanHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBaidu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetBaaiHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBAAI,
		BaseURL:  config[`api_endpoint`],
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetCohereHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderCohere,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetOllamaHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOllama,
		BaseURL:  config[`api_endpoint`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetXinferenceHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider:   llm.ProviderXinference,
		BaseURL:    config[`api_endpoint`],
		APIVersion: config[`api_version`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetJinaHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderJina,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetLingYiWanWuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderLingYiWanWu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetMoonShotHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderMoonshot,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetSparkHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderSpark,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetHunyuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderHunyuan,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetBaichuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBaichuan,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetZhipuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderZhipu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func GetMinimaxHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderMiniMax,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newModelCallHandler(modelInfo, config, useModel, clientConfig)
}

func ResolveOpenRouterEndpoint() string {
	return ``
}

func ResolveOrcaRouterEndpoint() string {
	// OrcaRouter exposes an OpenAI-compatible gateway at the /v1 namespace.
	// Returning the full endpoint here keeps the provider self-contained so the
	// OpenAI-compatible client does not need any provider-specific default.
	return DefaultOrcaRouterEndpoint + `/v1`
}
