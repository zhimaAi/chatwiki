// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package official_account

import (
	"chatwiki/internal/pkg/lib_define"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/spf13/cast"
)

func gzhBuildSendMsgAppends(customer string, push *lib_define.PushMessage) *object.HashMap {
	appends := object.HashMap{`touser`: customer}
	if push != nil && len(push.Robot) > 0 && cast.ToBool(push.Robot[`show_ai_msg_gzh`]) {
		appends[`aimsgcontext`] = object.HashMap{`is_ai_msg`: 1} //消息下方增加灰色 wording “内容由第三方AI生成”
	}
	return &appends
}
