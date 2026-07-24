// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"github.com/cloudwego/eino/schema"
	"github.com/zhimaAi/go_tools/tool"
)

func newSkillDebugLog(logType string, content any, message *schema.Message) map[string]any {
	debugLog := map[string]any{
		`type`:        logType,
		`content`:     content,
		`create_time`: tool.Time2Int(),
	}
	if message == nil || message.ResponseMeta == nil || message.ResponseMeta.Usage == nil {
		return debugLog
	}
	debugLog[`prompt_tokens`] = message.ResponseMeta.Usage.PromptTokens
	debugLog[`completion_tokens`] = message.ResponseMeta.Usage.CompletionTokens
	return debugLog
}
