// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

const LockPreKey = `chatwiki.op_lock.`
