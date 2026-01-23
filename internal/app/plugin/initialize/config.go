// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/pkg/lib_web"
	"fmt"
	"runtime"

	"github.com/zhimaAi/go_tools/logs"

	"github.com/Unknwon/goconfig"
	"github.com/zhimaAi/go_tools/tool"
)

func initConfig() {
	file := `configs/plugin/config_pro.ini`
	define.Env = "production"
	if define.IsDev {
		file = `configs/plugin/config_dev.ini`
		define.Env = "development"
		if tool.InArrayString(runtime.GOOS, []string{`windows`, `darwin`}) {
			file = `configs/plugin/config_loc.ini`
			define.Env = "local"
		}
	}
	logs.Info(`read config file: %s`, file)
	config, err := goconfig.LoadConfigFile(file)
	if err != nil {
		panic(`read config file error`)
	}
	define.Config.NumCPU, err = config.GetSection(`num_cpu`)
	if err != nil {
		panic(`read config num_cpu error`)
	}
	define.Config.Postgres, err = config.GetSection(`postgres`)
	if err != nil {
		panic(`read config postgres error`)
	}
	define.Config.Redis, err = config.GetSection(`redis`)
	if err != nil {
		panic(`read config redis error`)
	}
	define.Config.Xiaokefu, err = config.GetSection(`xiaokefu`)
	if err != nil {
		panic(`read config xiaokefu error`)
	}
	define.Config.WebService, err = config.GetSection(`webservice`)
	if err != nil {
		logs.Error(err.Error())
		panic(`read config webservice error`)
	}
	ip := lib_web.GetPublicIp()
	if len(define.Config.WebService[`wechat_article_crawler_host`]) == 0 {
		define.Config.WebService[`wechat_article_crawler_host`] = fmt.Sprintf(`http://%s:18086`, ip)
	}
	define.Config.RpcService, err = config.GetSection(`rpcservice`)
	if err != nil {
		panic(`read config rpcservice error`)
	}
	define.DefaultPhpEnv = []string{
		"CRAWLER_HOST=" + define.Config.WebService[`crawler`],
		"WECHAT_ARTICLE_CRAWLER_HOST=" + define.Config.WebService[`wechat_article_crawler_host`],
		"WECHAT_ARTICLE_CRAWLER_API_TOKEN=" + define.Config.WebService[`wechat_article_crawler_api_token`],
	}
}
