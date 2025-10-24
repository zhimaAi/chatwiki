// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/code_run/define"
	"runtime"

	"github.com/Unknwon/goconfig"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initConfig() {
	file := `configs/code_run/config_pro.ini`
	define.Env = "production"
	if define.IsDev {
		file = `configs/code_run/config_dev.ini`
		define.Env = "development"
		if tool.InArrayString(runtime.GOOS, []string{`windows`, `darwin`}) {
			file = `configs/code_run/config_loc.ini`
			define.Env = "local"
		}
	}
	logs.Info(`read config file:%s`, file)
	config, err := goconfig.LoadConfigFile(file)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config file error`)
	}
	define.Config.WebService, err = config.GetSection(`webservice`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config webservice error`)
	}
	define.Config.NumCPU, err = config.GetSection(`num_cpu`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config num_cpu error`)
	}
}
