// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func init() {
	goose.AddMigrationNoTxContext(up, nil)
}

func buildAlterSql(tableName string) string {
	return fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "error_msg" TYPE varchar(1000)`, tableName)
}

func up(_ context.Context, db *sql.DB) error {
	_, err := db.Exec(buildAlterSql(common.ChatMonitorBaseTableName))
	if err != nil {
		return err
	}
	curTablename := fmt.Sprintf(`%s_%s`, common.ChatMonitorBaseTableName, tool.Date(`Ym`))
	exist, err := msql.Model(`information_schema.tables`, define.Postgres).Where(`table_name`, curTablename).Value(`1`)
	if err != nil {
		return err
	}
	if cast.ToInt(exist) == 0 {
		return nil
	}
	_, err = db.Exec(buildAlterSql(curTablename))
	if err != nil {
		return err
	}
	return nil
}
