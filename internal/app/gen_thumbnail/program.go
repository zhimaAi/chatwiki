// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package gen_thumbnail

import (
	"chatwiki/internal/app/gen_thumbnail/define"
	"chatwiki/internal/app/gen_thumbnail/initialize"
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
