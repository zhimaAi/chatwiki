// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/message_service/define"
	"flag"

	"github.com/zhimaAi/go_tools/logs"
)

const usage = `false is pro,true is dev,default run dev`

func Initialize() {
	//get run env
	flag.BoolVar(&define.IsDev, `IsDev`, true, usage)
	flag.Parse()
	if define.IsDev {
		logs.SetTerminal(true)
		logs.Info(`current run env:dev`)
	} else {
		logs.Info(`current run env:pro`)
	}
	//initialize config
	initConfig()
	logs.Info(`initialize config finish`)
	//initialize cpu
	initNumCPU()
	logs.Info(`initialize cpu finish`)
	//initialize redis
	initRedis()
	logs.Info(`initialize redis finish`)
	//initialize postgres
	initPostgres()
	logs.Info(`initialize postgres finish`)
	//initialize gin
	initGin()
	logs.Info(`initialize gin finish`)
}
