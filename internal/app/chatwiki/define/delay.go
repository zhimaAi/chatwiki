// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"time"
)

var DelayTicker = time.NewTicker(time.Second)

const (
	OfficialAccountBatchSendDelayTask       = 1 // official account message mass sending delay task
	OfficialAccountBatchSendSyncCommentTask = 2 // official account message comment synchronization delay task
)

const (
	ReplayHisCommentSwitch1 = 1
	ReplayHisCommentSwitch2 = 2
)

type BaseDelayTask struct {
	Type uint `json:"type"`
}

type DelayTaskEvent struct {
	BaseDelayTask
	AdminUserId            int `json:"admin_user_id"`
	TaskId                 int `json:"task_id"`
	ReplayHisCommentSwitch int `json:"replay_his_comment_switch"`
}
