// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const LlmRequestLogsBaseTableName = `llm_request_logs`

func GetLlmRequestLogsTableName(timestamp uint) string {
	return fmt.Sprintf(`%s_%s`, LlmRequestLogsBaseTableName, tool.Date(`Ym`, timestamp))
}

func SaveLlmRequestLogs(data msql.Datas) error {
	tableName := GetLlmRequestLogsTableName(0)
	_, err := msql.Model(tableName, define.Postgres).Insert(data)
	if err == nil {
		return nil
	}
	var sqlerr *pq.Error
	if errors.As(err, &sqlerr) && sqlerr.Code == `42P01` { // Table does not exist
		// Create new table
		sql := fmt.Sprintf(`CREATE TABLE "%s" (LIKE "%s" INCLUDING ALL)`, tableName, LlmRequestLogsBaseTableName)
		_, err = msql.RawExec(define.Postgres, sql, nil)
		if err != nil {
			logs.Error(err.Error()) // Error creating new table
		}
		// Try to re-insert
		if _, err = msql.Model(tableName, define.Postgres).Insert(data); err == nil {
			return nil
		}
	}
	logs.Error(err.Error())
	return err
}
