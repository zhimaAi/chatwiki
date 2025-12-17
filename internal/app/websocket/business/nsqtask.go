// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"encoding/json"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func WsMessagePush(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	openid := cast.ToString(data[`openid`])
	message, err := json.Marshal(data[`message`])
	if err != nil {
		logs.Error(`message failure:%s/%s`, msg, err.Error())
		return nil
	}
	EventPushChan <- &WsMessage{openid: openid, message: message}
	return nil
}
