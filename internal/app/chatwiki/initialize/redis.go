// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/zhimaAi/go_tools/logs"
)

func initRedis() {
	server := define.Config.Redis["host"] + `:` + define.Config.Redis["port"]
	option := &redis.Options{Addr: server, Password: define.Config.Redis["password"], DB: 0}
	define.Redis = redis.NewClient(option)
	if _, err := define.Redis.Ping(context.Background()).Result(); err != nil {
		logs.Error(err.Error())
		panic(`Redis Ping Error`)
	}
}
