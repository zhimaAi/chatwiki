// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/client_side_build/define"
	"runtime"

	"github.com/Unknwon/goconfig"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initConfig() {
	file := `configs/client_side_build/config_pro.ini`
	define.Env = "production"
	if define.IsDev {
		file = `configs/client_side_build/config_dev.ini`
		define.Env = "development"
		if tool.InArrayString(runtime.GOOS, []string{`windows`, `darwin`}) {
			file = `configs/client_side_build/config_loc.ini`
			define.Env = "local"
		}
	}
	logs.Info(`read config file:%s`, file)
	config, err := goconfig.LoadConfigFile(file)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config file error`)
	}
	define.Config.NumCPU, err = config.GetSection(`num_cpu`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config num_cpu error`)
	}
	define.Config.NsqLookup, err = config.GetSection(`nsqlookup`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config nsqlookup error`)
	}
	define.Config.Nsqd, err = config.GetSection(`nsqd`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config nsqd error`)
	}
}
