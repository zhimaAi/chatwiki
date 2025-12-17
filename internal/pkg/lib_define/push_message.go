// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_define

import "github.com/zhimaAi/go_tools/msql"

type PushMessage struct {
	MsgRaw      string
	Message     map[string]any
	AdminUserId int
	CreateTime  int
	Openid      string
	Content     string
	AppInfo     msql.Params
	Robot       msql.Params
	Customer    msql.Params
}
