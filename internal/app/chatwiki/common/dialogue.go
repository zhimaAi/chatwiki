// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"errors"
	"fmt"

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
	m := msql.Model(`chat_ai_dialogue`, define.Postgres)
	if isBackground != 1 { //not background create
		dialogueId, _ := m.Where(`robot_id`, chatBaseParam.Robot[`id`]).Where(`openid`, chatBaseParam.Openid).
			Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).Order(`id DESC`).Value(`id`)
		if cast.ToUint(dialogueId) > 0 {
			return cast.ToInt(dialogueId), nil //use the old first
		}
	}
	id, err := m.Insert(msql.Datas{
		`admin_user_id`: chatBaseParam.AdminUserId,
		`robot_id`:      chatBaseParam.Robot[`id`],
		`openid`:        chatBaseParam.Openid,
		`is_background`: isBackground,
		`subject`:       MbSubstr(question, 0, 1000),
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}, `id`)
	if err == nil { //clear cached data
		lib_redis.DelCacheData(define.Redis, &DialogueCacheBuildHandler{DialogueId: int(id)})
		lib_redis.DelCacheData(define.Redis, &WeChatDialogueCacheBuildHandler{
			AdminUserId: chatBaseParam.AdminUserId, RobotId: cast.ToInt(chatBaseParam.Robot[`id`]), Openid: chatBaseParam.Openid})
	}
	return int(id), err
}

func SessionCacheKey(dialogueId int) string {
	return fmt.Sprintf(`chatwiki.get_session.by_dialogue.%d`, dialogueId)
}

func GetSessionId(params *define.ChatRequestParam, dialogueId int) (int, error) {
	cacheKey := SessionCacheKey(dialogueId)
	sessionId, err := define.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if sessionId := cast.ToInt(sessionId); sessionId > 0 {
		return sessionId, nil
	}
	//create new session
	var appId string
	if params.ChatBaseParam != nil && len(params.AppInfo) > 0 {
		appId = params.AppInfo[`app_id`]
	}
	id, err := msql.Model(`chat_ai_session`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`:     params.ChatBaseParam.AdminUserId,
		`app_type`:          params.ChatBaseParam.AppType,
		`app_id`:            appId,
		`dialogue_id`:       dialogueId,
		`robot_id`:          params.ChatBaseParam.Robot[`id`],
		`openid`:            params.ChatBaseParam.Openid,
		`last_chat_time`:    tool.Time2Int(),
		`last_chat_message`: MbSubstr(params.Question, 0, 1000),
		`create_time`:       tool.Time2Int(),
		`update_time`:       tool.Time2Int(),
	}, `id`)
	if err != nil {
		return 0, err
	}
	//write cache
	_, err = define.Redis.Set(context.Background(), cacheKey, id, GetSessionTtl()).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
	//create new receiver
	createNewReceiver(params, id)
	return int(id), nil
}

func UpLastChat(dialogueId, sessionId int, lastChat msql.Datas, isCustomer int) {
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
	_, err = define.Redis.Expire(context.Background(), SessionCacheKey(dialogueId), GetSessionTtl()).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
	//update receiver
	updateReceiver(sessionId, lastChat, isCustomer)
}
