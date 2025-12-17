// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/pkg/lib_define"

	"github.com/ArtisanCloud/PowerLibs/v3/cache"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var wechatCache *cache.GRedis

func GetWechatCache() *cache.GRedis {
	if wechatCache == nil {
		wechatCache = cache.NewGRedis(&redis.UniversalOptions{
			Addrs:                    []string{lib_define.Redis.Options().Addr},
			Password:                 lib_define.Redis.Options().Password,
			MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
		})
	}
	return wechatCache
}
