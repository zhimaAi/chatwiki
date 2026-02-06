// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/gen_thumbnail/define"
	"runtime"

	"github.com/Unknwon/goconfig"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initConfig() {
	file := `configs/gen_thumbnail/config_pro.ini`
	define.Env = "production"
	if define.IsDev {
		file = `configs/gen_thumbnail/config_dev.ini`
		define.Env = "development"
		if tool.InArrayString(runtime.GOOS, []string{`windows`, `darwin`}) {
			file = `configs/gen_thumbnail/config_loc.ini`
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
