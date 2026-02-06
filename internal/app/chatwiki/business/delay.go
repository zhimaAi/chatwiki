// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
			logs.Error(`zrange by score error:%s/%s`, tool.JsonEncodeNoError(opt), err.Error())
		}
		if len(list) == 0 {
			continue // indicates that there are no tasks at this moment, skip this round
		}
		for _, task := range list {
			go DelayTaskTrigger(task)
		}
		err = define.Redis.ZRemRangeByScore(context.Background(), define.DelayZset, opt.Min, opt.Max).Err()
		if err != nil {
			logs.Error(`zrem range by score error:%s/%s/%s`, opt.Min, opt.Max, err.Error())
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
			logs.Error(`nsq production error, fallback to sync logic:%s`, err.Error())
			_ = OfficialAccountBatchSend(tool.JsonEncodeNoError(taskInfo))
		}
	case define.OfficialAccountBatchSendSyncCommentTask:
		taskInfo := define.DelayTaskEvent{}
		if err := tool.JsonDecodeUseNumber(task, &taskInfo); err != nil {
			logs.Error(`task:%s,err:%s`, task, err.Error())
			break
		}

		if err := common.AddJobs(define.OfficialAccountCommentSyncTopic, tool.JsonEncodeNoError(taskInfo)); err != nil {
			logs.Error(`nsq production error, fallback to sync logic:%s`, err.Error())
			_ = OfficialAccountCommentSync(tool.JsonEncodeNoError(taskInfo))
		}
	default:
		logs.Error(`unknown delay task type:%s`, task)
	}
}
