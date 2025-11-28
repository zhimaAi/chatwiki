// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/plugin/business"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
)

var Route lib_web.Route

func init() {
	//step1:initialize
	Route = make(map[string]map[string]lib_web.Action)
	Route[http.MethodGet] = make(map[string]lib_web.Action)
	Route[http.MethodPost] = make(map[string]lib_web.Action)
	Route[lib_web.NoMethod] = make(map[string]lib_web.Action)
	Route[lib_web.NoRoute] = make(map[string]lib_web.Action)
	//step2:define route
	Route[http.MethodGet][`/ping`] = business.Ping   //ping<-->pong
	Route[lib_web.NoMethod][`/`] = business.NoMethod //NoMethod
	Route[lib_web.NoRoute][`/`] = business.NoRoute   //NoMethod

	//local plugins
	Route[http.MethodGet][`/manage/plugin/local-plugins/list`] = business.GetLocalPluginList
	Route[http.MethodGet][`/manage/plugin/local-plugins/detail`] = business.GetLocalPluginDetail
	Route[http.MethodPost][`/manage/plugin/local-plugins/destroy`] = business.DestroyLocalPlugin
	Route[http.MethodPost][`/manage/plugin/local-plugins/load`] = business.LoadLocalPlugin
	Route[http.MethodPost][`/manage/plugin/local-plugins/unload`] = business.UnloadLocalPlugin
	Route[http.MethodPost][`/manage/plugin/local-plugins/run`] = business.RunLocalPluginLambda

	//remote plugins
	Route[http.MethodGet][`/manage/plugin/remote-plugins/list`] = business.GetRemotePluginList
	Route[http.MethodGet][`/manage/plugin/remote-plugins/detail`] = business.GetRemotePluginDetail
	Route[http.MethodPost][`/manage/plugin/remote-plugins/download`] = business.DownloadRemotePlugin

}
