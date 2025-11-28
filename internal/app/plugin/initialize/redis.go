// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/pkg/lib_define"
	"context"
	"github.com/go-redis/redis/v8"
)

func initRedis() {
	server := define.Config.Redis["host"] + `:` + define.Config.Redis["port"]
	option := &redis.Options{Addr: server, Password: define.Config.Redis["password"], DB: 0}
	define.Redis = redis.NewClient(option)
	if _, err := define.Redis.Ping(context.Background()).Result(); err != nil {
		panic(`Redis Ping Error`)
	}
	lib_define.Redis = define.Redis //pkg
}
