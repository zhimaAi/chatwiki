// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/pkg/lib_define"
	"flag"
	"github.com/zhimaAi/go_tools/logs"
)

const usage = `false is pro,true is dev,default run dev`

func Initialize() {
	flag.BoolVar(&define.IsDev, `IsDev`, true, usage)
	lib_define.IsDev = define.IsDev //pkg
	flag.Parse()

	if define.IsDev {
		logs.SetTerminal(true)
		logs.Info("current run env:dev")
	} else {
		logs.Info("current run env:pro")
	}

	logs.Info("init config")
	initConfig()

	logs.Info("init postgres")
	initPostgres()

	logs.Info("init redis")
	initRedis()

	initGin()
	logs.Info(`initialize gin finish`)

	logs.Info("init php_plugin_manager")
	initPhpPluginPool()
}
