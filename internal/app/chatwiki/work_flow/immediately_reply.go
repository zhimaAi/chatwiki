// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
)

// ImmediatelyReplyBuildReplyContent 组装立即回复节点的参数成ReplyContent结构
func ImmediatelyReplyBuildReplyContent(content string) common.ReplyContent {
	return common.ReplyContent{
		ReplyType:   common.ReplyTypeText,
		Description: content,
		Status:      `1`,
		Type:        common.ReplyTypeText,
		SendSource:  `work_flow_immediately_reply`,
	}
}
