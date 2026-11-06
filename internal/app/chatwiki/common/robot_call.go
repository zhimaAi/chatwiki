// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"maps"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

// BuildRobotCallParams keeps the supplied entry dialogue while isolating the target
// robot's configuration and request-specific overrides.
func BuildRobotCallParams(parent *define.ChatRequestParam, robot msql.Params, question string, global map[string]any, dialogueId int) *define.ChatRequestParam {
	return &define.ChatRequestParam{
		ChatBaseParam: &define.ChatBaseParam{
			AppType: parent.AppType, AppInfo: maps.Clone(parent.AppInfo),
			Openid: parent.Openid, AdminUserId: parent.AdminUserId,
			Robot: maps.Clone(robot), Customer: maps.Clone(parent.Customer), RelUserId: parent.RelUserId,
		},
		Lang: parent.Lang, Question: question, DialogueId: dialogueId, WorkFlowGlobal: maps.Clone(global),
		WechatappAppid: parent.WechatappAppid,
		StopCtx:        parent.StopCtx, IsClose: parent.IsClose, ManualTransfer: parent.ManualTransfer,
	}
}

// ChatHistoryRobotID matches the entry robot that owns the inherited conversation.
func ChatHistoryRobotID(params *define.ChatRequestParam) int {
	robot := params.ManualTransfer.SourceRobot()
	if len(robot) == 0 {
		robot = params.Robot
	}
	return cast.ToInt(robot[`id`])
}
