// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"time"

	"github.com/zhimaAi/go_tools/logs"
)

func AddJobs(topic, message string, delay ...time.Duration) error {
	topic = define.Env + `_` + topic
	return define.ProducerHandle.AddJobs(topic, message, delay...)
}

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
