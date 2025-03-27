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
	goose.AddMigrationNoTxContext(up, nil)
}

func up(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(`ALTER TABLE "chat_ai_chat_monitor" ALTER COLUMN "error_msg" TYPE varchar(1000)`)
	if err != nil {
		return err
	}
	exist, err := msql.Model(`information_schema.tables`, define.Postgres).Where(`table_name`, `chat_ai_chat_monitor_202503`).Value(`1`)
	if err != nil {
		return err
	}
	if cast.ToInt(exist) == 0 {
		return nil
	}
	_, err = db.Exec(`ALTER TABLE "chat_ai_chat_monitor_202503" ALTER COLUMN "error_msg" TYPE varchar(1000)`)
	if err != nil {
		return err
	}
	return nil
}
