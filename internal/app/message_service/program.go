// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package message_service

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/app/message_service/initialize"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/mq"
)

func Run() {
	//initialize
	initialize.Initialize()
	//producer handle
	define.ProducerHandle = mq.NewProducerHandle().SetWorkNum(5).SetHostAndPort(define.Config.Nsqd[`host`], cast.ToUint(define.Config.Nsqd[`port`]))
	//web start
	go lib_web.WebRun(define.WebService)
	//pprof api
	go func() {
		err := http.ListenAndServe(":55558", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
}

func Stop() {
	lib_web.Shutdown(define.WebService)
	define.ProducerHandle.Stop()
}
