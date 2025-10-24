// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func getTopic(msgType string) string {
	if msgType == `event` {
		return lib_define.PushEvent
	} else {
		return lib_define.PushMessage
	}
}

func PushNSQ(message map[string]interface{}) {
	topic := getTopic(cast.ToString(message[`MsgType`]))
	msg, err := tool.JsonEncode(message)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if err := AddJobs(topic, msg); err != nil {
		logs.Error(err.Error())
		return
	}
	if define.IsDev {
		logs.Debug(`topic:%s,msg:%s`, topic, msg)
	}
}
