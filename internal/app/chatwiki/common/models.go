// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/llm/adaptor"
	"errors"
	"github.com/zhimaAi/go_tools/tool"

	"github.com/gin-contrib/sse"
	"github.com/zhimaAi/go_tools/logs"
)

func GetModelInfoByDefine(modelDefine string) (define.ModelInfo, bool) {
	for _, info := range define.ModelList {
		if info.ModelDefine == modelDefine {
			return info, true
		}
	}
	return define.ModelInfo{}, false
}
func IsMultiConfModel(defineName string) bool {
	return tool.InArrayString(defineName, []string{define.ModelOllama, define.ModelXnference, define.ModelOpenAIAgent})
}
func GetOfflineTypeModelInfos(isOffline bool) []define.ModelInfo {
	var result []define.ModelInfo
	for _, info := range define.ModelList {
		if isOffline && info.IsOffline || !isOffline && !info.IsOffline {
			result = append(result, info)
		}
	}
	return result
}

func GetModelCallHandler(modelConfigId int, useModel string) (*define.ModelCallHandler, error) {
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
	return modelInfo.CallHandlerFunc(config, useModel)
}

func GetVector2000(modelConfigId int, useModel, input string) (string, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return ``, err
	}

	return handler.GetVector2000(input)
}

func GetSimilarity(modelConfigId int, useModel string, query []float64, inputs [][]float64) (string, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return ``, err
	}

	return handler.GetSimilarity(query, inputs)
}

func RequestChatStream(modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, chanStream chan sse.Event, temperature float32, maxToken int) (string, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return ``, err
	}
	return handler.RequestChatStream(messages, chanStream, temperature, maxToken)
}

func RequestChat(modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, temperature float32, maxToken int) (string, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return ``, err
	}
	return handler.RequestChat(messages, temperature, maxToken)
}
