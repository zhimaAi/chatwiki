// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package plugin

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/app/plugin/initialize"
	"chatwiki/internal/app/plugin/php"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
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
		err = define.PhpPlugin.LoadPhpPlugin(dbPlugin[`name`], define.Version, define.DefaultPhpEnv)
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

	// 默认开启的插件
	var defaultOpenPluginList []string
	defaultOpenPluginListStr, isExist := define.Config.WebService[`default_open_plugin_list`]
	if isExist {
		defaultOpenPluginList = strings.Split(defaultOpenPluginListStr, ",")
	}
	for _, pluginName := range defaultOpenPluginList {
		// 检查本地是否存在该插件
		_, err := php.GetPluginManifest(pluginName)
		if err != nil {
			logs.Error("plugin %s not found locally: %s", pluginName, err.Error())
			continue
		}

		//// 检查数据库中是否已有该插件的配置记录
		//dbPluginDetail, err := msql.Model(`plugin_config`, define.Postgres).Where(`name`, pluginName).Find()
		//if err != nil {
		//	logs.Error("failed to query plugin_config for %s: %s", pluginName, err.Error())
		//	continue
		//}
		//
		//// 如果数据库中不存在该插件的记录,则创建记录并加载插件
		//if len(dbPluginDetail) == 0 {
		//	// 创建数据库记录
		//	_, err = msql.Model(`plugin_config`, define.Postgres).Insert(msql.Datas{
		//		`name`:       pluginName,
		//		`has_loaded`: true,
		//	})
		//	if err != nil {
		//		logs.Error("failed to insert plugin_config for %s: %s", pluginName, err.Error())
		//		continue
		//	}
		//
		//	// 加载插件
		//	err = define.PhpPlugin.LoadPhpPlugin(pluginName, define.Version, define.DefaultPhpEnv)
		//	if err != nil {
		//		logs.Error("failed to load default plugin %s: %s", pluginName, err.Error())
		//		// 如果加载失败,将数据库记录的 has_loaded 设置为 false
		//		_, err = msql.Model(`plugin_config`, define.Postgres).Where(`name`, pluginName).Update(msql.Datas{
		//			`has_loaded`: false,
		//		})
		//		if err != nil {
		//			logs.Error("failed to update plugin_config for %s: %s", pluginName, err.Error())
		//		}
		//	} else {
		//		logs.Info("successfully loaded default plugin: %s", pluginName)
		//	}
		//}
	}

}

func Stop() {
	lib_web.Shutdown(define.WebService)
	initialize.StopPhpGoridge()
}
