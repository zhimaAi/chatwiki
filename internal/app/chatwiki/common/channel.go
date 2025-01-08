// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/pkg/lib_define"
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
	return list
}
