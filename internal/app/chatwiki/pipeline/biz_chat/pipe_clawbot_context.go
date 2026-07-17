// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import "context"

func clawbotRunContext(in *ChatInParam) context.Context {
	if in != nil && in.params != nil && in.params.StopCtx != nil {
		return in.params.StopCtx
	}
	return context.Background()
}
