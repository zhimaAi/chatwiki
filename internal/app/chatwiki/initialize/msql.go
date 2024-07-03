// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func initPostgres() {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		define.Config.Postgres["host"], define.Config.Postgres["port"],
		define.Config.Postgres["user"], define.Config.Postgres["password"],
		define.Config.Postgres["dbname"], define.Config.Postgres["sslmode"])
	if err := msql.RegisterDataBase(define.Postgres, conn, define.Postgres); err != nil {
		logs.Error(err.Error())
		panic(`Postgres Register Error`)
	}
	if err := msql.SetDebug(define.Postgres, true); err != nil {
		logs.Error(err.Error())
	}
}
