// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func SetMessageLasttimeCache(appid string, lasttime int) {
	if len(appid) == 0 || lasttime == 0 {
		return
	}
	key := define.WechatappWxkfMessageLasttime + appid
	_, err := define.Redis.Set(context.Background(), key, lasttime, 0).Result()
	if err != nil {
		logs.Error(err.Error())
	}
}

func GetMessageLasttimeCache(appid string) int {
	if len(appid) == 0 {
		return 0
	}
	key := define.WechatappWxkfMessageLasttime + appid
	lasttime, err := define.Redis.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return 0
	}
	return cast.ToInt(lasttime)
}

func SetSyncMsgTokenCache(appid, token string) {
	if len(appid) == 0 || len(token) == 0 {
		return
	}
	key := define.WechatappWxkfSyncMsgToken + appid
	_, err := define.Redis.Set(context.Background(), key, token, time.Minute*10).Result()
	if err != nil {
		logs.Error(err.Error())
	}
}

func GetSyncMsgTokenCache(appid string) string {
	if len(appid) == 0 {
		return ``
	}
	key := define.WechatappWxkfSyncMsgToken + appid
	token, err := define.Redis.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return ``
	}
	return token
}

func SetMessageRunningLock(appid string) bool {
	if len(appid) == 0 {
		return false
	}
	key := define.WechatappWxkfMessageRunning + appid
	ret, err := define.Redis.SetNX(context.Background(), key, tool.Time2String(), 3*time.Minute).Result()
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	return ret
}

func KeepMessageRunningLock(appid string) {
	if len(appid) == 0 {
		return
	}
	key := define.WechatappWxkfMessageRunning + appid
	ttl, err := define.Redis.TTL(context.Background(), key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	if ttl < time.Minute {
		_, err := define.Redis.Expire(context.Background(), key, 3*time.Minute).Result()
		if err != nil {
			logs.Error(err.Error())
		}
		keeptime := time.Now().Format("2006-01-02 15:04:05.000")
		fmt.Println(fmt.Sprintf(`%s[keep_lock]appid:%s`, keeptime, appid))
	}
}

func DelMessageRunningLock(appid string) {
	if len(appid) > 0 {
		key := define.WechatappWxkfMessageRunning + appid
		_, err := define.Redis.Del(context.Background(), key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			logs.Error(err.Error())
		}
	}
}

func SetMessageCursorCache(appid, cursor string) {
	if len(appid) == 0 || len(cursor) == 0 {
		return
	}
	key := define.WechatappWxkfMessageCursor + appid
	_, err := define.Redis.Set(context.Background(), key, cursor, 0).Result()
	if err != nil {
		logs.Error(err.Error())
	}
}

func GetMessageCursorCache(appid string) string {
	if len(appid) == 0 {
		return ``
	}
	key := define.WechatappWxkfMessageCursor + appid
	cursor, err := define.Redis.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return ``
	}
	return cursor
}
