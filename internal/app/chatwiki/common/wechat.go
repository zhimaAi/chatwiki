// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/wechat"
	"errors"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"

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
	lib_redis.DelCacheData(define.Redis, &WechatAppCacheBuildHandler{Field: `id`, Value: appInfo[`id`]})
	lib_redis.DelCacheData(define.Redis, &WechatAppCacheBuildHandler{Field: `app_id`, Value: appInfo[`app_id`]})
	lib_redis.DelCacheData(define.Redis, &WechatAppCacheBuildHandler{Field: `access_key`, Value: appInfo[`access_key`]})
	return nil
}

func GetOfficialAccountApp(appid, secret string) (*officialAccount.OfficialAccount, error) {

	app, err := wechat.GetApplication(msql.Params{
		`app_type`:   lib_define.AppOfficeAccount,
		`app_id`:     appid,
		`app_secret`: secret})
	if err != nil {
		return nil, err
	}
	officialApp, ok := app.(wechat.OfficialAccountInterface)
	if !ok {
		return nil, errors.New("公众号初始化失败")
	}

	return officialApp.GetAccountClient()

}
