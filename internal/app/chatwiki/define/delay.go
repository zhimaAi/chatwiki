// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"time"
)

var DelayTicker = time.NewTicker(time.Second)

const (
	OfficialAccountBatchSendDelayTask       = 1 //公众号消息群发延时任务
	OfficialAccountBatchSendSyncCommentTask = 2 //公众号消息评论同步延时任务
)

type BaseDelayTask struct {
	Type uint `json:"type"`
}

type DelayTaskEvent struct {
	BaseDelayTask
	AdminUserId int `json:"admin_user_id"`
	TaskId      int `json:"task_id"`
}
