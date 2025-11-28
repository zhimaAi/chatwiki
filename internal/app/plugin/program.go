// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package plugin

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/app/plugin/initialize"
	"chatwiki/internal/pkg/lib_web"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"net/http"
)

func Run() {
	initialize.Initialize()

	go lib_web.WebRun(define.WebService)
	go func() {
		err := http.ListenAndServe(":55559", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()

	initialize.StartPhpGoridge()

	dbPluginList, err := msql.Model(`plugin_config`, define.Postgres).Where(`has_loaded`, cast.ToString(true)).Select()
	if err != nil {
		logs.Error(err.Error())
		panic(err)
	}
	for _, dbPlugin := range dbPluginList {
		err = define.PhpPlugin.LoadPhpPlugin(dbPlugin[`name`], define.Version)
		if err != nil {
			logs.Error(err.Error())
			_, err = msql.Model(`plugin_config`, define.Postgres).Where(`id`, cast.ToString(dbPlugin[`id`])).Update(msql.Datas{
				`has_loaded`: false,
			})
			if err != nil {
				logs.Error(err.Error())
			}
		}
	}
}

func Stop() {
	lib_web.Shutdown(define.WebService)
	initialize.StopPhpGoridge()
}
