// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"
	"fmt"
	"runtime"

	"github.com/Unknwon/goconfig"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initConfig() {
	file := `configs/chatwiki/config_pro.ini`
	define.Env = "production"
	if define.IsDev {
		file = `configs/chatwiki/config_dev.ini`
		define.Env = "development"
		if tool.InArrayString(runtime.GOOS, []string{`windows`, `darwin`}) {
			file = `configs/chatwiki/config_loc.ini`
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
	//the ip and port of the domain name are not configured
	ip := lib_web.GetPublicIp()
	if len(define.Config.WebService[`h5_domain`]) == 0 {
		define.Config.WebService[`h5_domain`] = fmt.Sprintf(`http://%s:18081`, ip)
	}
	if len(define.Config.WebService[`pc_domain`]) == 0 {
		define.Config.WebService[`pc_domain`] = fmt.Sprintf(`http://%s:18082`, ip)
	}
	if len(define.Config.WebService[`ws_domain`]) == 0 {
		define.Config.WebService[`ws_domain`] = fmt.Sprintf(`%s:18083`, ip)
	}
	define.Config.NumCPU, err = config.GetSection(`num_cpu`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config num_cpu error`)
	}
	define.Config.Redis, err = config.GetSection("redis")
	if err != nil {
		logs.Error(err.Error())
		panic(`read config redis error`)
	}
	define.Config.Postgres, err = config.GetSection("postgres")
	if err != nil {
		logs.Error(err.Error())
		panic(`read config postgres error`)
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

	define.Config.Crawler, err = config.GetSection(`crawler`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config crawler error`)
	}
}
