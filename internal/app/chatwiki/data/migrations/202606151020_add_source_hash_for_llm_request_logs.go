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
		return common.UnifyRunSql(common.LlmRequestLogsBaseTableName,
			`ALTER TABLE "public"."llm_request_logs" ADD COLUMN "source_hash" varchar(32) NOT NULL DEFAULT '';COMMENT ON COLUMN "public"."llm_request_logs"."source_hash" IS '隐式溯源标识(hash,MD5)'`)
	}, nil)
}
