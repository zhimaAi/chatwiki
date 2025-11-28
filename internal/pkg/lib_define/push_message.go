// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
