// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package user_domain_service

import (
	"chatwiki/internal/app/user_domain_service/define"
	"chatwiki/internal/app/user_domain_service/initialize"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	_ "net/http/pprof"

	"github.com/zhimaAi/go_tools/logs"
)

func Run() {
	//initialize
	initialize.Initialize()
	//web start
	go lib_web.WebRun(define.WebService)
	//pprof api
	go func() {
		err := http.ListenAndServe(":55565", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
}

func Stop() {
	lib_web.Shutdown(define.WebService)
}
