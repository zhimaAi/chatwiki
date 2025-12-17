// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/pkg/pipeline"
)

// CheckManualReplyPauseRobotReply 人工介入+人工回复后暂停机器人回复
func CheckManualReplyPauseRobotReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if IsManualReplyPauseRobotReply(in.params, in.sessionId, in.chanStream) {
		out.AiMessage = nil
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckKeywordSwitchManual 关键词转人工
func CheckKeywordSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if msg, ok := IsKeywordSwitchManual(in.params, in.sessionId, in.dialogueId, in.chanStream); ok {
		out.AiMessage = msg
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckIntentionSwitchManual 根据用户意图转人工
func CheckIntentionSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if msg, ok := IsIntentionSwitchManual(in.params, in.sessionId, in.dialogueId, in.monitor, in.chanStream); ok {
		out.AiMessage = msg
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}
