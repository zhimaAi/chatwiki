// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func StopDelayService() {
	define.DelayTicker.Stop()
}

func AddDelayTask(task any, seconds int64) {
	taskStr, err := tool.JsonEncode(task)
	if err != nil {
		logs.Error(`task:%v,err:%s`, task, err.Error())
		return
	}
	member := &redis.Z{Score: float64(time.Now().Unix() + seconds), Member: taskStr}
	err = define.Redis.ZAdd(context.Background(), define.DelayZset, member).Err()
	if err != nil {
		logs.Error(`ZAdd error:%s/%d/%s`, taskStr, member.Score, err.Error())
	}
}

func DelDelayTask(task any) {
	taskStr, err := tool.JsonEncode(task)
	if err != nil {
		logs.Error(`task:%v,err:%s`, task, err.Error())
		return
	}
	err = define.Redis.ZRem(context.Background(), define.DelayZset, taskStr).Err()
	if err != nil {
		logs.Error(`ZRem error:%s/%s`, taskStr, err.Error())
	}
}
