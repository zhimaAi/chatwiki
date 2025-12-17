// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"

	"github.com/gin-contrib/sse"
	"github.com/zhimaAi/go_tools/msql"
)

func CloseReceiverFromAppOpenApi(params *define.ChatRequestParam) {

}

func IsManualReplyPauseRobotReply(params *define.ChatRequestParam, sessionId int, chanStream chan sse.Event) bool {
	return false
}

func IsKeywordSwitchManual(params *define.ChatRequestParam, sessionId, dialogueId int, chanStream chan sse.Event) (msql.Params, bool) {
	return nil, false
}

func IsIntentionSwitchManual(params *define.ChatRequestParam, sessionId, dialogueId int, monitor *common.Monitor, chanStream chan sse.Event) (msql.Params, bool) {
	return nil, false
}

func IsRepetitionSwitchManual(params *define.ChatRequestParam, sessionId, dialogueId int, curMsgId int64, chanStream chan sse.Event) (msql.Params, bool) {
	return nil, false
}

func IsUnknownSwitchManual(params *define.ChatRequestParam, sessionId, dialogueId int, list []msql.Params, chanStream chan sse.Event) (msql.Params, bool) {
	return nil, false
}

func AdditionQuoteLib(params *define.ChatRequestParam, list []msql.Params, message *msql.Datas) {

}
