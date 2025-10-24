package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetWechatKefuOpenid(adminUserId int, message map[string]any) string {
	appid := strings.TrimSpace(cast.ToString(message[`appid`]))
	externalUserid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
	openKfid := strings.TrimSpace(cast.ToString(message[`open_kfid`]))
	//generate a custom openid
	openid := tool.MD5(fmt.Sprintf(`%s:%s:%s`, appid, externalUserid, openKfid))
	if e, o := GetExternalUserInfo(openid); len(e) == 0 || len(o) == 0 {
		//save the mapping relationship to the database
		data := msql.Datas{
			`admin_user_id`:   adminUserId,
			`openid`:          openid,
			`app_id`:          appid,
			`external_userid`: externalUserid,
			`open_kfid`:       openKfid,
			`create_time`:     tool.Time2Int(),
			`update_time`:     tool.Time2Int(),
		}
		_, err := msql.Model(`wechat_kefu_openid_map`, define.Postgres).Insert(data)
		if err != nil {
			logs.Error(err.Error())
		}
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &WechatKefuOpenidCacheBuildHandler{Openid: openid})
	}
	return openid
}

type WechatKefuOpenidCacheBuildHandler struct{ Openid string }

func (h *WechatKefuOpenidCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.wechat_kefu_openid_map.%s`, h.Openid)
}
func (h *WechatKefuOpenidCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`wechat_kefu_openid_map`, define.Postgres).Where(`openid`, h.Openid).Find()
}

func GetExternalUserInfo(openid string) (externalUserid, openKfid string) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &WechatKefuOpenidCacheBuildHandler{Openid: openid}, &result, time.Hour)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	return result[`external_userid`], result[`open_kfid`]
}
