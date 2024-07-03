// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetDialogueId(chatBaseParam *define.ChatBaseParam, question string) (int, error) {
	var isBackground int
	if len(chatBaseParam.Customer) > 0 && cast.ToInt(chatBaseParam.Customer[`is_background`]) > 0 {
		isBackground = 1
	}
	id, err := msql.Model(`chat_ai_dialogue`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: chatBaseParam.AdminUserId,
		`robot_id`:      chatBaseParam.Robot[`id`],
		`openid`:        chatBaseParam.Openid,
		`is_background`: isBackground,
		`subject`:       question,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}, `id`)
	if err == nil { //clear cached data
		lib_redis.DelCacheData(define.Redis, &DialogueCacheBuildHandler{DialogueId: int(id)})
	}
	return int(id), err
}

func sessionCacheKey(dialogueId int) string {
	return fmt.Sprintf(`chatwiki.get_session.by_dialogue.%d`, dialogueId)
}

func GetSessionId(chatBaseParam *define.ChatBaseParam, dialogueId int) (int, error) {
	cacheKey := sessionCacheKey(dialogueId)
	sessionId, err := define.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if sessionId := cast.ToInt(sessionId); sessionId > 0 {
		return sessionId, nil
	}
	//create new session
	id, err := msql.Model(`chat_ai_session`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: chatBaseParam.AdminUserId,
		`dialogue_id`:   dialogueId,
		`robot_id`:      chatBaseParam.Robot[`id`],
		`openid`:        chatBaseParam.Openid,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}, `id`)
	if err != nil {
		return 0, err
	}
	//write cache
	_, err = define.Redis.Set(context.Background(), cacheKey, id, 30*time.Minute).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
	return int(id), nil
}

func UpLastChat(dialogueId, sessionId int, lastChat msql.Datas) {
	if len(lastChat) == 0 {
		return
	}
	lastChat[`update_time`] = tool.Time2Int()
	_, err := msql.Model(`chat_ai_session`, define.Postgres).
		Where(`id`, cast.ToString(sessionId)).Update(lastChat)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//update session_id ttl
	_, err = define.Redis.Expire(context.Background(), sessionCacheKey(dialogueId), 30*time.Minute).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
}
