// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/llm/adaptor"
	"errors"
	"github.com/zhimaAi/go_tools/msql"
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

func GetVector2000(adminUserId int, openid string, robot msql.Params, library msql.Params, file msql.Params, modelConfigId int, useModel, input string) (string, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return ``, err
	}

	return handler.GetVector2000(adminUserId, openid, robot, library, file, input)
}

func RequestChatStream(adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, chanStream chan sse.Event, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	return handler.RequestChatStream(adminUserId, openid, robot, appType, messages, chanStream, temperature, maxToken)
}

func RequestChat(adminUserId int, openid string, robot msql.Params, appType string, modelConfigId int, useModel string, messages []adaptor.ZhimaChatCompletionMessage, temperature float32, maxToken int) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return adaptor.ZhimaChatCompletionResponse{}, 0, err
	}
	return handler.RequestChat(adminUserId, openid, robot, appType, messages, temperature, maxToken)
}
