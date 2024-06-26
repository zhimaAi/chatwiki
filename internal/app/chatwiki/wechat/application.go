// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package wechat

import (
	"chatwiki/internal/app/chatwiki/wechat/mini_program"
	"chatwiki/internal/app/chatwiki/wechat/official_account"
	"chatwiki/internal/app/chatwiki/wechat/wechat_kefu"
	"chatwiki/internal/pkg/lib_define"
	"errors"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/zhimaAi/go_tools/msql"
)

type ApplicationInterface interface {
	SendText(touser, content string) (int, error)
	GetToken(refresh bool) (*response.ResponseGetToken, int, error)
	SendMsgOnEvent(code, content string) (int, error)
}

func GetApplication(appInfo msql.Params, openKfid string) (ApplicationInterface, error) {
	if len(appInfo) == 0 {
		return nil, errors.New(`app info is empty`)
	}
	switch appInfo[`app_type`] {
	case lib_define.AppOfficeAccount:
		return &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.AppMini:
		return &mini_program.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.AppWechatKefu:
		return &wechat_kefu.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`], OpenKfid: openKfid}, nil
	}
	return nil, errors.New(`app type not support`)
}
