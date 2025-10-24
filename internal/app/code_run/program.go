// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package code_run

import (
	"chatwiki/internal/app/code_run/define"
	"chatwiki/internal/app/code_run/initialize"
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
		err := http.ListenAndServe(":55559", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
}

func Stop() {
	lib_web.Shutdown(define.WebService)
}
