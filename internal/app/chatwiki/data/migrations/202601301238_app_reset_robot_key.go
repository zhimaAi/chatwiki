// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"

	"chatwiki/internal/pkg/lib_redis"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		go func() {
			err := ResetAppEmptyRobotKey()
			if err != nil {
				logs.Other(`migration`, `error:%s`, err.Error())
			}
		}()
		return nil
	}, nil)
}

func ResetAppEmptyRobotKey() error {
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`robot_id`, cast.ToString(0)).
		Where(`robot_key`, `!=`, ``)
	oldAppInfoList, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(oldAppInfoList) > 0 {
		// clear apps with empty robot_key
		_, err = m.Where(`robot_id`, cast.ToString(0)).
			Where(`robot_key`, `!=`, ``).
			Update(msql.Datas{`robot_key`: ``})
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		// clear cache
		for _, appInfo := range oldAppInfoList {
			lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `id`, Value: appInfo[`id`]})
			lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: appInfo[`app_id`]})
			lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: appInfo[`access_key`]})
		}
	}

	return nil
}
