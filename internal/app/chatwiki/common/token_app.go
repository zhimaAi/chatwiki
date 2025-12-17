// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func TokenAppAllowUse(adminUserId, robotId int, tokenAppType string) bool {
	configCache := TokenAppLimitConfigCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: tokenAppType,
		RobotId:      robotId,
	}
	config := msql.Params{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &configCache, &config, time.Hour*30)
	if err != nil {
		logs.Error(err.Error())
		return true
	}
	switchStatus := cast.ToInt(config[`switch_status`])
	if switchStatus == 0 {
		return true
	}
	useCache := TokenAppUseCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: tokenAppType,
		RobotId:      robotId,
		DateYmd:      ``,
	}
	useToken, err := lib_redis.GetCacheIncrWithBuild(define.Redis, &useCache, time.Hour*30)
	if err != nil {
		logs.Error(err.Error())
		return true
	}
	maxToken := cast.ToInt64(config[`max_token`])
	if maxToken <= cast.ToInt64(useToken) {
		return false
	}
	return true
}

func TokenAppIsSwitchOpen(adminUserId, robotId int, tokenAppType string) bool {
	configCache := TokenAppLimitConfigCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: tokenAppType,
		RobotId:      robotId,
	}
	config := msql.Params{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &configCache, &config, time.Hour*30)
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	switchStatus := cast.ToInt(config[`switch_status`])
	if switchStatus == 0 {
		return false
	}
	return true
}
