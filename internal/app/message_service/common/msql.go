// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"time"

	"github.com/zhimaAi/go_tools/msql"
)

type WechatAppCacheBuildHandler struct {
	Field string
	Value string
}

func (h *WechatAppCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.app_info.%s.%s`, h.Field, h.Value)
}
func (h *WechatAppCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(h.Field, h.Value).Find()
}

func GetWechatAppInfo(field, value string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &WechatAppCacheBuildHandler{Field: field, Value: value}, &result, time.Hour)
	return result, err
}
