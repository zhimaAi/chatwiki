// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
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
// websocket notifications. dialogueId/sessionId default to 0 => no chat history
// (stateless call), consistent with how WorkflowNode invokes a sub-workflow.
func runClawbotForWorkflow(robotInfo msql.Params, question string, lang string) (string, error) {
	if len(lang) == 0 {
		lang = define.LangZhCn
	}
	chatBaseParam := &define.ChatBaseParam{
		AppType:     lib_define.AppOpenApi,
		Openid:      ``,
		AdminUserId: cast.ToInt(robotInfo[`admin_user_id`]),
		Robot:       robotInfo,
	}
	chatParams := &define.ChatRequestParam{
		ChatBaseParam: chatBaseParam,
		Lang:          lang,
		Question:      question,
	}
	// A non-nil channel with a draining goroutine is required: the error path of
	// doApplicationTypeClaw (pipe_clawbot.go) calls SendDefaultUnknownQuestionPrompt,
	// which writes to chanStream directly (bypassing in.Stream). A nil channel would
	// block forever on send. The goroutine discards events so nothing reaches the
	// C-side, preserving the "no trace" behaviour. useStream is false, so the normal
	// path never writes to the channel (in.Stream is a no-op when useStream is false).
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
		useStream:  false,
		chanStream: chanStream,
		monitor:    common.NewMonitor(chatParams),
	}
	out := &ChatOutParam{
		debugLog: make([]any, 0),
	}
	if err := doApplicationTypeClaw(in, out); err != nil {
		return ``, err
	}
	if out.Error != nil {
		return ``, out.Error
	}
	return out.content, nil
}
