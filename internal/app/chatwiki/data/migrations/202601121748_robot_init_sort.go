// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		go func() {
			err := InitRobotSort()
			if err != nil {
				logs.Other(`migration`, `async process init sort error: %s`, err.Error())
			}
		}()
		return nil
	}, nil)
}

func InitRobotSort() error {
	m := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_name,admin_user_id`).
		Order(`admin_user_id asc,id asc`) //基础排序处理
	list, err := m.Select()
	//排序处理
	if err != nil {
		logs.Other(`migration`, `query all robot sort data failed: %s`, err.Error())
		return err
	}
	logs.Other(`migration`, `start processing robot sort data`)
	sortNum := 1
	lastAdminUserId := 0
	for _, robot := range list {
		if lastAdminUserId != cast.ToInt(robot[`admin_user_id`]) {
			lastAdminUserId = cast.ToInt(robot[`admin_user_id`])
			sortNum = 1
		}
		_, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, robot[`id`]).
			Update(msql.Datas{`sort_num`: sortNum})
		if err != nil {
			logs.Other(`migration`, `update robot sort failed: %s`, err.Error())
			return err
		}
		sortNum++
	}
	logs.Other(`migration`, `completed robot sort processing`)
	return nil
}
