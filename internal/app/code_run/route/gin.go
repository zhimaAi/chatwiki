// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/code_run/business"
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
	/*code_run API*/
	Route[http.MethodPost][`/javaScript`] = business.JavaScript
}
