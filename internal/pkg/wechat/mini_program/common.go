// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package mini_program

import (
	"chatwiki/internal/pkg/lib_define"
	"context"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	response2 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/spf13/cast"
)

// Copy From app.CustomerServiceMessage.Send
func miniCustomMessageSend(app *miniProgram.MiniProgram, push *lib_define.PushMessage, ctx context.Context, toUser string, msgType string, msg interface{}) (*response2.ResponseMiniProgram, error) {
	data := object.HashMap{
		`touser`:  toUser,
		`msgtype`: msgType,
		msgType:   msg,
	}
	if push != nil && len(push.Robot) > 0 && cast.ToBool(push.Robot[`show_ai_msg_mini`]) {
		data[`aimsgcontext`] = object.HashMap{`is_ai_msg`: 1} //消息下方增加灰色 wording “内容由第三方AI生成”
	}
	result := &response2.ResponseMiniProgram{}
	_, err := app.CustomerServiceMessage.BaseClient.HttpPostJson(ctx, `cgi-bin/message/custom/send`, &data, nil, nil, result)
	return result, err
}
