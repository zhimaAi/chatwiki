// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"github.com/zhimaAi/go_tools/msql"
	llm "github.com/zhimaAi/llm_adaptor/v2"
)

func newSupplierHandler(config msql.Params, clientConfig llm.ClientConfig) (*SupplierHandler, error) {
	client, err := llm.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	return &SupplierHandler{
		Client: client,
		config: config,
	}, nil
}

func Get302AiSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.Provider302AI,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetOpenRouterSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenRouter,
		BaseURL:  ResolveOpenRouterEndpoint(),
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetOrcaRouterSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenAI,
		BaseURL:  ResolveOrcaRouterEndpoint(),
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetDeepseekSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderDeepSeek,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetGeminiSupplierHandler(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderGemini,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetOpenAISupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOpenAI,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetDoubaoSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderDoubao,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetSiliconFlowSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderSiliconFlow,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetTongyiSupplierHandler(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderAli,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetOpenAIAgentSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider:   llm.ProviderOpenAIAgent,
		BaseURL:    config[`api_endpoint`],
		APIVersion: config[`api_version`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetAzureSupplierHandler(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderAzure,
		BaseURL:  config[`api_endpoint`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetClaudeSupplierHandler(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderClaude,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetYiyanSupplierHandler(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBaidu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetBaaiSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBAAI,
		BaseURL:  config[`api_endpoint`],
	}
	return newSupplierHandler(config, clientConfig)
}

func GetCohereSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderCohere,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetOllamaSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderOllama,
		BaseURL:  config[`api_endpoint`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetXinferenceSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider:   llm.ProviderXinference,
		BaseURL:    config[`api_endpoint`],
		APIVersion: config[`api_version`],
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetJinaSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderJina,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetLingYiWanWuSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderLingYiWanWu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetMoonShotSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderMoonshot,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetSparkSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderSpark,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetHunyuanSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderHunyuan,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetBaichuanSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderBaichuan,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetZhipuSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderZhipu,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}

func GetMinimaxSupplierHandle(_ ModelInfo, config msql.Params) (*SupplierHandler, error) {
	clientConfig := llm.ClientConfig{
		Provider: llm.ProviderMiniMax,
		Credentials: llm.CredentialConfig{
			APIKeys: config[`api_key`],
		},
	}
	return newSupplierHandler(config, clientConfig)
}
