// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

type SessionChannel struct {
	AppType string `json:"app_type"`
	AppName string `json:"app_name"`
	AppId   string `json:"app_id"`
}

func GetChannelList(userId int, robotId uint) []SessionChannel {
	list := []SessionChannel{
		{AppType: lib_define.AppYunH5, AppName: `WebAPP`},
		{AppType: lib_define.AppYunPc, AppName: `嵌入网站`},
		{AppType: lib_define.AppOpenApi, AppName: `开放接口`},
	}
	//wechat_app
	if robotId > 0 {
		apps, err := msql.Model(define.TableChatAiWechatApp, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_id`, cast.ToString(robotId)).Field(`app_type,app_name,app_id`).Select()
		if err != nil {
			logs.Error(err.Error())
		}
		for _, app := range apps {
			list = append(list, SessionChannel{AppType: app[`app_type`], AppName: app[`app_name`], AppId: app[`app_id`]})
		}
	}
	return list
}
