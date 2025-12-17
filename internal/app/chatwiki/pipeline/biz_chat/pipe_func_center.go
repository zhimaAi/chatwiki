// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"

	"github.com/gin-contrib/sse"
)

// CheckKeywordReply 关键词检测处理
func CheckKeywordReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	var keywordReplyList []common.ReplyContent
	keywordReplyList, in.keywordSkipAI, _ = common.BuildKeywordReplyMessage(in.params)
	if len(keywordReplyList) > 0 {
		out.replyContentList = append(out.replyContentList, keywordReplyList...)
	}
	return pipeline.PipeContinue
}

// CheckReceivedMessageReply 收到消息回复处理
func CheckReceivedMessageReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.keywordSkipAI {
		return pipeline.PipeContinue
	}
	receivedMessageReplyList, _ := common.BuildReceivedMessageReply(in.params, lib_define.MsgTypeText)
	if len(receivedMessageReplyList) > 0 {
		out.replyContentList = append(out.replyContentList, receivedMessageReplyList...)
	}
	return pipeline.PipeContinue
}

// PushReplyContentList 推送回复内容列表
func PushReplyContentList(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(out.replyContentList) > 0 {
		in.Stream(sse.Event{Event: `reply_content_list`, Data: out.replyContentList})
	}
	return pipeline.PipeContinue
}
