// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type CacheBuildHandler interface {
	GetCacheKey() string
	GetCacheData() (any, error)
}

func GetCacheWithBuild(cache *redis.Client, handler CacheBuildHandler, result any, ttl time.Duration) error {
	cacheKey := handler.GetCacheKey()
	if len(cacheKey) == 0 {
		return errors.New(`cache key is empty`)
	}
	jsonStr, err := cache.Get(context.Background(), cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return err
	}
	if len(jsonStr) == 0 { //cache not created
		data, err := handler.GetCacheData()
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		jsonStr, err = tool.JsonEncode(data)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if data == nil || len(jsonStr) == 0 || jsonStr == `{}` || jsonStr == `[]` { //no data
			jsonStr = `-1`
		}
		_, err = cache.Set(context.Background(), cacheKey, jsonStr, ttl).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			logs.Error(err.Error())
			return err
		}
	}
	if jsonStr == `-1` { //no data
		return nil
	}
	return tool.JsonDecodeUseNumber(jsonStr, result)
}

func DelCacheData(cache *redis.Client, handler CacheBuildHandler) {
	cacheKey := handler.GetCacheKey()
	if len(cacheKey) == 0 {
		return
	}
	_, err := cache.Del(context.Background(), cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
}

func AddLock(cache *redis.Client, key string, ttl time.Duration) bool {
	ok, err := cache.SetNX(context.Background(), key, time.Now().Unix(), ttl).Result()
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	return ok
}

func UnLock(cache *redis.Client, key string) {
	_, err := cache.Del(context.Background(), key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
}
