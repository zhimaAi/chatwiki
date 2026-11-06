// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
)

// Keep the edition-specific entry at the same path and with the same signature;
// this edition does not register a human customer service tool.
func newManualTransferTool(in *ChatInParam, workflowTools ...chat.Tool) einotool.BaseTool {
	return nil
}
