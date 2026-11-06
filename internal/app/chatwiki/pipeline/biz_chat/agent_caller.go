// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"

	"github.com/gin-contrib/sse"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	// register the agent caller so the work_flow package can invoke a clawbot
	// synchronously without depending on this package at compile time.
	work_flow.AgentCaller = runClawbotForWorkflow
}

// runClawbotForWorkflow runs a clawbot synchronously and returns its full reply.
// It calls doApplicationTypeClaw directly (not the full DoChatRequest pipeline),
// so it does NOT persist any C-side chat message / dialogue and does NOT push
// websocket notifications. Nested calls retain the entry conversation and
// current message while using the target robot's configuration.
func runClawbotForWorkflow(robotInfo msql.Params, question string, params *work_flow.WorkFlowParams) (string, []common.ReplyContent, bool, error) {
	params.ManualTransfer.SetSource(params.ChatRequestParam, params.SessionId)
	chatParams := common.BuildRobotCallParams(params.ChatRequestParam, robotInfo, question, nil, params.DialogueId)
	if len(chatParams.Lang) == 0 {
		chatParams.Lang = define.LangZhCn
	}

	chanStream := make(chan sse.Event)
	go func() {
		for range chanStream {
			// discard events to prevent the producer from blocking
		}
	}()
	// doApplicationTypeClaw is synchronous; once it returns, no further sends happen,
	// so closing here lets the draining goroutine exit without a send-on-closed panic.
	defer close(chanStream)
	in := &ChatInParam{
		params:     chatParams,
		dialogueId: chatParams.DialogueId,
		sessionId:  params.SessionId,
		useStream:  false,
		chanStream: chanStream,
		monitor:    common.NewMonitor(chatParams),
	}
	out := &ChatOutParam{
		debugLog: make([]any, 0),
		cMsgId:   int64(params.CurMsgId),
	}
	if err := doApplicationTypeClaw(in, out); err != nil {
		if content, transferred := chatParams.ManualTransfer.Result(); transferred {
			return content, out.replyContentList, true, nil
		}
		return ``, out.replyContentList, false, err
	}
	if out.Error != nil {
		return ``, out.replyContentList, false, out.Error
	}
	return out.content, out.replyContentList, in.isSwitchManual, nil
}
