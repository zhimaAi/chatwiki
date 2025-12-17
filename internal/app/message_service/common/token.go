// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/message_service/define"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func GetWeComApp(corpId string, secret string) (*work.Work, error) {
	return work.NewWork(&work.UserConfig{
		CorpID: corpId, Secret: secret,
		OAuth:     work.OAuth{Callback: `https://xxx.xxx`},
		HttpDebug: false, Debug: define.IsDev,
	})
}

func GetCorpAccessToken(appInfo msql.Params) string {
	WeComApp, err := GetWeComApp(appInfo[`app_id`], appInfo[`app_secret`])
	if err != nil {
		logs.Error(err.Error())
		return ``
	}
	resToken, err := WeComApp.AccessToken.GetToken(false)
	if err != nil {
		logs.Error(err.Error())
		return ``
	}
	return resToken.AccessToken
}
