// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/plugin/define"
	"fmt"
	"github.com/zhimaAi/go_tools/msql"
)

func initPostgres() {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		define.Config.Postgres["host"], define.Config.Postgres["port"],
		define.Config.Postgres["user"], define.Config.Postgres["password"],
		define.Config.Postgres["dbname"], define.Config.Postgres["sslmode"])
	if err := msql.RegisterDataBase(define.Postgres, conn, define.Postgres); err != nil {
		panic(`Postgres Register Error`)
	}
	if err := msql.SetDebug(define.Postgres, true); err != nil {
		panic(`Postgres SetDebug Error`)
	}
}
