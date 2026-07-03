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

// WhatsappInboundDedup msgid dedup: Aliyun ChatApp retries pushing the same message
// (more so under multi-IP/slow responses); SetNX ensures each MsgId is processed only once
// to avoid duplicate replies.
const WhatsappInboundDedup = `chatwiki.whatsapp.inbound.dedup.`
