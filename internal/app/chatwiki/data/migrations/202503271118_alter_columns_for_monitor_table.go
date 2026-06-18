// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		return common.UnifyRunSql(common.ChatMonitorBaseTableName,
			`ALTER TABLE "public"."chat_ai_chat_monitor" ALTER COLUMN "error_msg" TYPE varchar(1000)`)
	}, nil)
}
