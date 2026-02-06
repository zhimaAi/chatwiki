// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
)

// CheckManualReplyPauseRobotReply manual intervention + pause robot reply after manual reply
func CheckManualReplyPauseRobotReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if IsManualReplyPauseRobotReply(in.params, in.sessionId, in.chanStream) {
		out.AiMessage = nil
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckKeywordSwitchManual keyword switch to manual
func CheckKeywordSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if msg, ok := IsKeywordSwitchManual(in.params, in.sessionId, in.dialogueId, in.chanStream); ok {
		out.AiMessage = msg
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckIntentionSwitchManual switch to manual based on user intention
func CheckIntentionSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.AppInfo) > 0 && len(in.params.ReceivedMessageType) > 0 && in.params.ReceivedMessageType != lib_define.MsgTypeText {
		return pipeline.PipeContinue // wechat and other apps, skip intention switch for non-text messages
	}
	if msg, ok := IsIntentionSwitchManual(in.params, in.sessionId, in.dialogueId, in.monitor, in.chanStream); ok {
		out.AiMessage = msg
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}
