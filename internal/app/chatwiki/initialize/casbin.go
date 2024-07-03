// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/casbin"
	"fmt"

	"github.com/zhimaAi/go_tools/logs"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func initCasbin() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		define.Config.Postgres["host"], define.Config.Postgres["port"],
		define.Config.Postgres["user"], define.Config.Postgres["password"],
		define.Config.Postgres["dbname"], define.Config.Postgres["sslmode"])
	Db, err = gorm.Open(pgDriver.New(pgDriver.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		logs.Error(fmt.Sprintf("connection postgre fail :%s", err))
		return
	}
	if err := casbin.Handler.Init(Db); err != nil {
		logs.Error(fmt.Sprintf("init casbin fail :%s", err))
		return
	}
	middlewares.Init()
}
