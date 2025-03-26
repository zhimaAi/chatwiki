// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"errors"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
	"strings"

	"github.com/zhimaAi/go_tools/msql"
)

func GetAzureHandler(config msql.Params, _ string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `azure`,
			EndPoint:   config[`api_endpoint`],
			APIVersion: config[`api_version`],
			APIKey:     config[`api_key`],
			Model:      config[`deployment_name`],
		},
		config: config,
	}
	return handler, nil
}

func GetClaudeHandler(config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `claude`,
			APIKey:     config[`api_key`],
			APIVersion: config[`api_version`],
			Model:      useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetGeminiHandler(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetYiyanHandler(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func CheckYiyanFancCall(modelInfo ModelInfo, config msql.Params, useModel string) error {
	if config[`secret_key`] == "" {
		useModel = strings.ToLower(useModel)
	}
	if tool.InArrayString(useModel, modelInfo.SupportedFunctionCallList) {
		return nil
	}
	return errors.New(`model is not support`)
}

func GetTongyiHandler(config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:   `ali`,
			APIKey: config[`api_key`],
			Model:  useModel,
		},
		config: config,
	}
	return handler, nil
}

func GetBaaiHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetCohereHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetOllamaHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
	if useModel == "默认" && cast.ToString(config["deployment_name"]) != "" {
		useModel = cast.ToString(config["deployment_name"])
	}
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

func GetDeepseekHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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
func GetJinaHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetLingYiWanWuHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetMoonShotHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetBaichuanHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetZhipuHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetOpenAIHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetOpenAIAgentHandle(config msql.Params, _ string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:       `openaiAgent`,
			APIKey:     config[`api_key`],
			EndPoint:   config[`api_endpoint`],
			APIVersion: config["api_version"],
			Model:      config[`deployment_name`],
		},
		config: config,
	}
	return handler, nil
}

func GetSparkHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetHunyuanHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetDoubaoHandle(config msql.Params, _ string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			Corp:      `doubao`,
			APIKey:    config[`api_key`],
			SecretKey: config[`secret_key`],
			Region:    config[`region`],
			Model:     config[`deployment_name`],
		},
		config: config,
	}
	return handler, nil
}

func GetMinimaxHandle(config msql.Params, useModel string) (*ModelCallHandler, error) {
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

func GetSiliconFlow(config msql.Params, useModel string) (*ModelCallHandler, error) {
	handler := &ModelCallHandler{
		Meta: adaptor.Meta{
			EndPoint:   `https://api.siliconflow.cn`,
			Corp:       ModelSiliconFlow,
			APIKey:     config[`api_key`],
			APIVersion: `v1`,
			Model:      useModel,
		},
		config: config,
	}
	return handler, nil
}
