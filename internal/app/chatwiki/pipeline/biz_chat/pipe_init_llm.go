// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/pkg/pipeline"

	"github.com/gin-contrib/sse"
)

// RobotInfoPush 推送机器人信息
func RobotInfoPush(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.Stream(sse.Event{Event: `robot`, Data: in.params.Robot})
	return pipeline.PipeContinue
}
