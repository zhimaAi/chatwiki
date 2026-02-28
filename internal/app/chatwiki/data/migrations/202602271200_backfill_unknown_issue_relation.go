// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		go func() {
			_ = BackfillUnknownIssueRelationByHistory()
		}()
		return nil
	}, nil)
}

func BackfillUnknownIssueRelationByHistory() error {
	logs.Other(`migration`, `start backfill unknown issue relation data`)
	const batchSize = 500
	lastId := 0
	total := 0
	updated := 0
	notFound := 0
	failed := 0

	for {
		list, err := msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres).
			Where(`id`, `>`, cast.ToString(lastId)).
			Where(`sample_session_id`, cast.ToString(0)).
			Order(`id asc`).
			Limit(batchSize).
			Field(`id,admin_user_id,robot_id,stats_day,question`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if len(list) == 0 {
			break
		}
		for _, row := range list {
			total++
			lastId = cast.ToInt(row[`id`])
			startTs := tool.GetTimestamp(cast.ToUint(row[`stats_day`]))
			endTs := startTs + 86400
			msg, queryErr := msql.Model(`chat_ai_message`, define.Postgres).Alias(`m`).
				Join(`chat_ai_session s`, `m.session_id=s.id`, `left`).
				Where(`m.admin_user_id`, row[`admin_user_id`]).
				Where(`m.robot_id`, row[`robot_id`]).
				Where(`m.is_customer`, `1`).
				Where(`m.content`, row[`question`]).
				Where(`m.create_time`, `between`, fmt.Sprintf(`%d,%d`, startTs, endTs)).
				Order(`m.id desc`).
				Field(`m.id,m.openid,m.dialogue_id,m.session_id,m.create_time,s.rel_user_id`).
				Find()
			if queryErr != nil {
				failed++
				logs.Error(queryErr.Error())
				continue
			}
			if len(msg) == 0 {
				notFound++
				continue
			}
			updateData := msql.Datas{
				`sample_openid`:      msg[`openid`],
				`sample_rel_user_id`: cast.ToInt(msg[`rel_user_id`]),
				`sample_dialogue_id`: cast.ToInt(msg[`dialogue_id`]),
				`sample_session_id`:  cast.ToInt(msg[`session_id`]),
				`sample_message_id`:  cast.ToInt(msg[`id`]),
				`last_dialogue_id`:   cast.ToInt(msg[`dialogue_id`]),
				`last_session_id`:    cast.ToInt(msg[`session_id`]),
				`last_message_id`:    cast.ToInt(msg[`id`]),
				`last_trigger_time`:  cast.ToInt(msg[`create_time`]),
				`update_time`:        tool.Time2Int(),
			}
			if _, queryErr = msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres).
				Where(`id`, row[`id`]).
				Where(`sample_session_id`, cast.ToString(0)).
				Update(updateData); queryErr != nil {
				failed++
				logs.Error(queryErr.Error())
				continue
			}
			updated++
		}
		logs.Other(`migration`, `backfill unknown issue relation progress: total=%d,updated=%d,not_found=%d,failed=%d,last_id=%d`,
			total, updated, notFound, failed, lastId)
	}
	logs.Other(`migration`, `backfill unknown issue relation finished: total=%d,updated=%d,not_found=%d,failed=%d`,
		total, updated, notFound, failed)
	return nil
}
