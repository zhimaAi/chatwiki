// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/user_domain_service/define"
	"chatwiki/internal/pkg/lib_define"
	"flag"

	"github.com/zhimaAi/go_tools/logs"
)

const usage = `false is pro,true is dev,default run dev`

func Initialize() {
	//get run env
	flag.BoolVar(&define.IsDev, `IsDev`, true, usage)
	lib_define.IsDev = define.IsDev //pkg
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
	//initialize gin
	initGin()
	logs.Info(`initialize gin finish`)
}
