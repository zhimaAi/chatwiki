// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"context"
)

func DoWorkWithContext(ctx context.Context, doworkHandle func()) error {
	channel := make(chan struct{})
	go func() {
		doworkHandle()
		select {
		case channel <- struct{}{}: //发送通知
			close(channel) //关闭管道
		default:
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-channel:
		return nil
	}
}
