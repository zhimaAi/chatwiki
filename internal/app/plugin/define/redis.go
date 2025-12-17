// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

const RedisKeyLoadedLocalPluginList = "loaded_local_plugin_list"
