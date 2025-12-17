// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/websocket/define"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func RunTask(topic, channel string, workNum uint, callback func(msg string, args ...string) error) {
	topic = define.Env + `_` + topic
	err := define.ConsumerHandle.PushZero(define.Config.Nsqd[`host`]+`:`+define.Config.Nsqd[`port`], topic)
	if err != nil {
		logs.Error(`PushZero Error:%s`, err.Error())
	}
	err = define.ConsumerHandle.Run(topic, channel, workNum, callback)
	if err != nil {
		logs.Error(`Consumer Run Error:%s`, err.Error())
	}
}

func OpenPush(openid, ip string, stime int) {
	logs.Debug(`ClosePush:%s`, tool.JsonEncodeNoError(map[string]any{`openid`: openid, `ip`: ip, `stime`: stime}))
}

func ClosePush(openid, ip string, stime, etime int) {
	logs.Debug(`ClosePush:%s`, tool.JsonEncodeNoError(map[string]any{`openid`: openid, `ip`: ip, `stime`: stime, `etime`: etime}))
}
