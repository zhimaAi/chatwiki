// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package websocket

import (
	"chatwiki/internal/app/websocket/business"
	"chatwiki/internal/app/websocket/common"
	"chatwiki/internal/app/websocket/define"
	"chatwiki/internal/app/websocket/initialize"
	"chatwiki/internal/pkg/lib_define"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/mq"
)

func Run() {
	//initialize
	initialize.Initialize()
	//consumer handle
	define.ConsumerHandle = mq.NewConsumerHandle().SetHostAndPort(define.Config.NsqLookup[`host`], cast.ToUint(define.Config.NsqLookup[`port`]))
	//web start
	go business.InitWs()
	//pprof api
	go func() {
		err := http.ListenAndServe(":55560", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	//consumer start
	common.RunTask(lib_define.WsMessagePushTopic, lib_define.WsMessagePushChannel, 2, business.WsMessagePush)
}

func Stop() {
	define.ConsumerHandle.Stop()
}
