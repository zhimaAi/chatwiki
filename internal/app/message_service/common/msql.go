// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"time"

	"github.com/zhimaAi/go_tools/msql"
)

type RobotCacheBuildHandler struct{ RobotKey string }

func (h *RobotCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(lib_define.RedisPrefixRobotInfo, h.RobotKey)
}
func (h *RobotCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, h.RobotKey).Find()
}

func GetRobotInfo(robotKey string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotCacheBuildHandler{RobotKey: robotKey}, &result, time.Hour)
	return result, err
}

type WechatAppCacheBuildHandler struct {
	Field string
	Value string
}

func (h *WechatAppCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(lib_define.RedisPrefixAppInfo, h.Field, h.Value)
}
func (h *WechatAppCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(h.Field, h.Value).Find()
}

func GetWechatAppInfo(field, value string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &WechatAppCacheBuildHandler{Field: field, Value: value}, &result, time.Hour)
	return result, err
}
