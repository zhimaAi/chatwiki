// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"sync"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-contrib/sse"
)

// Guard every tool, including the filesystem tools installed by DeepAgent.
// The barrier waits for active callbacks before the caller reads their output.
func clawbotManualTransferContext(parent context.Context, state *define.ManualTransferState) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-state.Done():
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

func clawbotManualTransferTools(ctx context.Context, cancel context.CancelFunc, state *define.ManualTransferState) (compose.ToolMiddleware, func()) {
	var active sync.RWMutex
	enter := func() (func(), error) {
		active.RLock()
		if _, ok := state.Result(); ok {
			active.RUnlock()
			return nil, context.Canceled
		}
		if err := ctx.Err(); err != nil {
			active.RUnlock()
			return nil, err
		}
		return func() {
			if _, ok := state.Result(); ok {
				cancel()
			}
			active.RUnlock()
		}, nil
	}
	middleware := compose.ToolMiddleware{
		Invokable: func(next compose.InvokableToolEndpoint) compose.InvokableToolEndpoint {
			return func(ctx context.Context, input *compose.ToolInput) (*compose.ToolOutput, error) {
				leave, err := enter()
				if err != nil {
					return nil, err
				}
				defer leave()
				return next(ctx, input)
			}
		},
		Streamable: func(next compose.StreamableToolEndpoint) compose.StreamableToolEndpoint {
			return func(ctx context.Context, input *compose.ToolInput) (*compose.StreamToolOutput, error) {
				leave, err := enter()
				if err != nil {
					return nil, err
				}
				defer leave()
				return next(ctx, input)
			}
		},
		EnhancedInvokable: func(next compose.EnhancedInvokableToolEndpoint) compose.EnhancedInvokableToolEndpoint {
			return func(ctx context.Context, input *compose.ToolInput) (*compose.EnhancedInvokableToolOutput, error) {
				leave, err := enter()
				if err != nil {
					return nil, err
				}
				defer leave()
				return next(ctx, input)
			}
		},
		EnhancedStreamable: func(next compose.EnhancedStreamableToolEndpoint) compose.EnhancedStreamableToolEndpoint {
			return func(ctx context.Context, input *compose.ToolInput) (*compose.EnhancedStreamableToolOutput, error) {
				leave, err := enter()
				if err != nil {
					return nil, err
				}
				defer leave()
				return next(ctx, input)
			}
		},
	}
	return middleware, func() { active.Lock(); active.Unlock() }
}

func finishClawbotManualTransfer(in *ChatInParam, out *ChatOutParam) bool {
	content, ok := in.params.ManualTransfer.Result()
	if !ok {
		return false
	}
	in.isSwitchManual = true
	out.Error = nil
	out.content = content
	out.chatResp.SetResult(content)
	in.Stream(sse.Event{Event: `stream_message`, Data: schema.AssistantMessage(content, nil)})
	return true
}
