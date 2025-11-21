// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/wechat"

	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func RefreshAccountVerify(appInfo msql.Params) error {
	app, err := wechat.GetApplication(appInfo)
	if err != nil {
		return err
	}
	basic, _, err := app.GetAccountBasicInfo()
	if err != nil {
		return err
	}
	//database dispose
	data := msql.Datas{
		`account_customer_type`: basic.CustomerType,
		`update_time`:           tool.Time2Int(),
	}
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres)
	_, err = m.Where(`id`, appInfo[`id`]).Update(data)
	if err != nil {
		return err
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `id`, Value: appInfo[`id`]})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: appInfo[`app_id`]})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: appInfo[`access_key`]})
	return nil
}
