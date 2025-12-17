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

func GetClaudeHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `claude`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetGeminiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `gemini`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetYiyanHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `baidu`,
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetTongyiHandler(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:              `ali`,
			APIKey:            config[`api_key`],
			Model:             useModel,
			ChoosableThinking: tool.InArrayString(useModel, modelInfo.GetChoosableThinkingModels()),
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

func GetDeepseekHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `deepseek`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}
func GetJinaHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `jina`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetLingYiWanWuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `lingyiwanwu`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetMoonShotHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `moonshot`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetBaichuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `baichuan`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetZhipuHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `zhipu`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetOpenAIHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `openai`,
			APIKey: config[`api_key`],
			Model:  useModel,
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

func GetSparkHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `spark`,
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			APPID:     config[`app_id`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetHunyuanHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `hunyuan`,
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
			Model:     useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetDoubaoHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:              `doubao`,
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

func GetMinimaxHandle(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `minimax`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetSiliconFlow(modelInfo ModelInfo, config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			EndPoint:          `https://api.siliconflow.cn`,
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
