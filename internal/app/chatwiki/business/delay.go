// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func StartDelayService() {
	for range define.DelayTicker.C {
		opt := &redis.ZRangeBy{Min: `0`, Max: tool.Time2String()}
		list, err := define.Redis.ZRangeByScore(context.Background(), define.DelayZset, opt).Result()
		if err != nil {
			logs.Error(`ZRangeByScore 错误:%s/%s`, tool.JsonEncodeNoError(opt), err.Error())
		}
		if len(list) == 0 {
			continue //说明此刻没有任务可以跳过此轮
		}
		for _, task := range list {
			go DelayTaskTrigger(task)
		}
		err = define.Redis.ZRemRangeByScore(context.Background(), define.DelayZset, opt.Min, opt.Max).Err()
		if err != nil {
			logs.Error(`ZRemRangeByScore 错误:%s/%s/%s`, opt.Min, opt.Max, err.Error())
		}
	}
}

func DelayTaskTrigger(task string) {
	base := define.BaseDelayTask{}
	if err := tool.JsonDecodeUseNumber(task, &base); err != nil {
		logs.Error(`task:%s,err:%s`, task, err.Error())
		return
	}
	switch base.Type {
	case define.OfficialAccountBatchSendDelayTask:
		taskInfo := define.DelayTaskEvent{}
		if err := tool.JsonDecodeUseNumber(task, &taskInfo); err != nil {
			logs.Error(`task:%s,err:%s`, task, err.Error())
			break
		}

		if err := common.AddJobs(define.OfficialAccountBatchSendTopic, tool.JsonEncodeNoError(taskInfo)); err != nil {
			logs.Error(`NSQ生产异常,走同步逻辑,错误:%s`, err.Error())
			_ = OfficialAccountBatchSend(tool.JsonEncodeNoError(taskInfo))
		}
	case define.OfficialAccountBatchSendSyncCommentTask:
		taskInfo := define.DelayTaskEvent{}
		if err := tool.JsonDecodeUseNumber(task, &taskInfo); err != nil {
			logs.Error(`task:%s,err:%s`, task, err.Error())
			break
		}

		if err := common.AddJobs(define.OfficialAccountCommentSyncTopic, tool.JsonEncodeNoError(taskInfo)); err != nil {
			logs.Error(`NSQ生产异常,走同步逻辑,错误:%s`, err.Error())
			_ = OfficialAccountCommentSync(tool.JsonEncodeNoError(taskInfo))
		}
	default:
		logs.Error(`未知的延时任务类型:%s`, task)
	}
}
