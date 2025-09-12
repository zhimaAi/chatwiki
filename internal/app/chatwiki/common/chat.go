// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func BuildChatContextPair(openid string, robotId, dialogueId, curMsgId, contextPair int) []map[string]string {
	contextList := make([]map[string]string, 0)
	if contextPair <= 0 {
		return contextList //no context required
	}
	m := msql.Model(`chat_ai_message`, define.Postgres).Where(`openid`, openid).
		Where(`robot_id`, cast.ToString(robotId)).Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`msg_type`, cast.ToString(define.MsgTypeText))
	if curMsgId > 0 { //兼容调试运行获取上下文
		m.Where(`id`, `<`, cast.ToString(curMsgId))
	}
	list, err := m.Order(`id desc`).Field(`id,content,is_customer,is_valid_function_call`).Limit(contextPair * 4).Select()
	if err != nil {
		logs.Error(err.Error())
	}
	if len(list) == 0 {
		return contextList
	}
	//reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	//remove the record before the function call
	for index := len(list) - 1; index >= 0; index-- {
		if cast.ToBool(list[index][`is_valid_function_call`]) {
			if index == len(list)-1 {
				list = []msql.Params{}
			} else {
				list = list[index+1:]
			}
			break
		}
	}
	//foreach
	for i := 0; i < len(list)-1; i++ {
		if cast.ToInt(list[i][`is_customer`]) == define.MsgFromCustomer && cast.ToInt(list[i+1][`is_customer`]) == define.MsgFromRobot {
			contextList = append(contextList, map[string]string{`question`: list[i][`content`], `answer`: list[i+1][`content`]})
			i++ //skip answer
		}
	}
	//cut out
	if len(contextList) > contextPair {
		contextList = contextList[len(contextList)-contextPair:]
	}

	return contextList
}

func BuildOpenApiContent(params *define.ChatRequestParam, messages []adaptor.ZhimaChatCompletionMessage) []adaptor.ZhimaChatCompletionMessage {
	if params.AppType != lib_define.AppOpenApi || len(params.OpenApiContent) == 0 {
		return messages
	}
	var contents = make([]adaptor.ZhimaChatCompletionMessage, 0)
	err := tool.JsonDecode(params.OpenApiContent, &contents)
	if err != nil {
		logs.Error(err.Error())
		return messages
	}
	if len(contents) > 0 {
		messages = append(contents, messages...)
	}
	if define.IsDev {
		logs.Debug("%+v", messages)
	}
	return messages
}
