// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

type TriggerConfigCacheBuildHandler struct {
	AdminUserId int
	TriggerType string
}

func (h TriggerConfigCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.trigger.config.%d.%s`, h.AdminUserId, h.TriggerType)
}
func (h TriggerConfigCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`trigger_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(h.AdminUserId)).
		Where(`trigger_type`, h.TriggerType).Find()
}

type TriggerOfficialCacheBuildHandler struct {
	AppId   string
	MsgType string
}

func (h TriggerOfficialCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(define.RedisPrefixOfficialTrigger, h.AppId, h.MsgType)
}

func (h TriggerOfficialCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`find_key`, fmt.Sprintf(`%s.%s`, h.AppId, h.MsgType)).
		Select()
}

type TriggerFindKeyCacheBuildHandler struct {
	AdminUserId string
	FindKey     string
}

func (h TriggerFindKeyCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(define.RedisPrefixFindKeyTrigger, h.AdminUserId, h.FindKey)
}

func (h TriggerFindKeyCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`admin_user_id`, h.AdminUserId).
		Where(`find_key`, h.FindKey).
		Find()
}
