// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

const LockPreKey = `chatwiki.op_lock.`

const DelayZset = `chatwiki.delay.zset` // 存储延时的zset

const (
	RedisPrefixOfficialTrigger = `chatwiki.official_trigger.%s.%s`
	RedisPrefixFindKeyTrigger  = `chatwiki.find_key_trigger.%s.%s`
)
