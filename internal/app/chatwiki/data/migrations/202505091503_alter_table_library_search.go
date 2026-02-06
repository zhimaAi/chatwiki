// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		id, err := msql.Model(`goose_db_version`, define.Postgres).Where(`version_id`, `20250509`).Value(`id`)
		if err != nil {
			return err
		}
		if cast.ToInt(id) > 0 {
			return nil // already upgraded
		}
		_, err = db.Exec(`ALTER TABLE "user_search_config" ADD COLUMN "prompt_type" int4 NOT NULL DEFAULT 0;
								ALTER TABLE "user_search_config" ADD COLUMN "prompt" varchar(500) NOT NULL DEFAULT '';
								COMMENT ON COLUMN "user_search_config"."prompt" IS '提示词';`)
		if err != nil {
			return err
		}
		return nil
	}, nil)
}
