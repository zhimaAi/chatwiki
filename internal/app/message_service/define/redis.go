// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

const WechatappWxkfMessageLasttime = `chatwiki.wxkf.message.lasttime.`
const WechatappWxkfSyncMsgToken = `chatwiki.wxkf.sync.msg.token.`
const WechatappWxkfMessageRunning = `chatwiki.wxkf.message.running.`
const WechatappWxkfMessageCursor = `chatwiki.wxkf.message.cursor.`
