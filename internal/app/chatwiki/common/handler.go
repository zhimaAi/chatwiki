// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func GetAzureHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `azure`,
			EndPoint:   config[`api_endpoint`],
			APIVersion: config[`api_version`],
			APIKey:     config[`api_key`],
			Model:      useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetAzureSupplierHandler(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:       `azure`,
			EndPoint:   config[`api_endpoint`],
			APIVersion: config[`api_version`],
			APIKey:     config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetClaudeHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `claude`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetClaudeSupplierHandler(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `claude`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetGeminiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `gemini`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetGeminiSupplierHandler(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `gemini`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetYiyanHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `baidu`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetYiyanSupplierHandler(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:      `baidu`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetTongyiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:              `ali`,
			EndPoint:          config[`api_endpoint`],
			APIKey:            config[`api_key`],
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		},
		config: config,
	}
	return handler, nil
}

func GetTongyiSupplierHandler(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `ali`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetBaaiHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `baai`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetBaaiSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `baai`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetCohereHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `cohere`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetCohereSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `cohere`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetOllamaHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `ollama`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetOllamaSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `ollama`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetXinferenceHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `xinference`,
			EndPoint:   config[`api_endpoint`],
			APIKey:     config[`api_key`],
			APIVersion: config["api_version"],
			Model:      useModel,
		},
	}
	return handler, nil
}

func GetXinferenceSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:       `xinference`,
			EndPoint:   config[`api_endpoint`],
			APIKey:     config[`api_key`],
			APIVersion: config["api_version"],
		},
	}
	return handler, nil
}

func GetDeepseekHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `deepseek`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetDeepseekSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `deepseek`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetJinaHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `jina`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetJinaSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `jina`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetLingYiWanWuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `lingyiwanwu`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetLingYiWanWuSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `lingyiwanwu`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetMoonShotHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `moonshot`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetMoonShotSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `moonshot`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetBaichuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `baichuan`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetBaichuanSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `baichuan`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetZhipuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `zhipu`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetZhipuSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `zhipu`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetOpenAIHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `openai`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetOpenAISupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `openai`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetOpenAIAgentHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `openaiAgent`,
			APIKey:     config[`api_key`],
			EndPoint:   config[`api_endpoint`],
			APIVersion: config["api_version"],
			Model:      useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetOpenAIAgentSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:       `openaiAgent`,
			APIKey:     config[`api_key`],
			EndPoint:   config[`api_endpoint`],
			APIVersion: config["api_version"],
		},
		config: config,
	}
	return handler, nil
}

func GetSparkHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `spark`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			APPID:     config[`app_id`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetSparkSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:      `spark`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			APPID:     config[`app_id`],
		},
		config: config,
	}
	return handler, nil
}

func GetHunyuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `hunyuan`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetHunyuanSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:      `hunyuan`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
		},
		config: config,
	}
	return handler, nil
}

func GetDoubaoHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:              `doubao`,
			EndPoint:          config[`api_endpoint`],
			APIKey:            config[`api_key`],
			SecretKey:         config[`secret_key`],
			Region:            config[`region`],
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		},
		config: config,
	}
	return handler, nil
}

func GetDoubaoSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:      `doubao`,
			EndPoint:  config[`api_endpoint`],
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
		},
		config: config,
	}
	return handler, nil
}

func GetMinimaxHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:     `minimax`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
			Model:    useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetMinimaxSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:     `minimax`,
			EndPoint: config[`api_endpoint`],
			APIKey:   config[`api_key`],
		},
		config: config,
	}
	return handler, nil
}

func GetSiliconFlowHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	endpoint := config[`api_endpoint`]
	if len(endpoint) == 0 {
		endpoint = `https://api.siliconflow.cn`
	}
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			EndPoint:          endpoint,
			Corp:              ModelSiliconFlow,
			APIKey:            config[`api_key`],
			APIVersion:        `v1`,
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		},
		config: config,
	}
	return handler, nil
}

func GetSiliconFlowSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	endpoint := config[`api_endpoint`]
	if len(endpoint) == 0 {
		endpoint = `https://api.siliconflow.cn`
	}
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			EndPoint:   endpoint,
			Corp:       ModelSiliconFlow,
			APIKey:     config[`api_key`],
			APIVersion: `v1`,
		},
		config: config,
	}
	return handler, nil
}

func Get302AiHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:              `302ai`,
			APIKey:            config[`api_key`],
			SecretKey:         config[`secret_key`],
			Region:            config[`region`],
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		},
		config: config,
	}
	return handler, nil
}

func Get302AiSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			Corp:      `302ai`,
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
		},
		config: config,
	}
	return handler, nil
}

func GetOpenRouterHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			EndPoint:          config[`api_endpoint`],
			Corp:              `openrouter`,
			APIKey:            config[`api_key`],
			APIVersion:        `v1`,
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
		},
		config: config,
	}
	return handler, nil
}

func GetOpenRouterSupplierHandle(modelInfo ModelInfo, config msql.Params) (*SupplierHandler, error) {
	handler := &SupplierHandler{
		Meta: adaptor.Meta{
			EndPoint:   config[`api_endpoint`],
			Corp:       `openrouter`,
			APIKey:     config[`api_key`],
			APIVersion: `v1`,
		},
		config: config,
	}
	return handler, nil
}
