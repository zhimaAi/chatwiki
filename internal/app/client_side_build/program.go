// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package client_side_build

import (
	"chatwiki/internal/app/client_side_build/business"
	"chatwiki/internal/app/client_side_build/common"
	"chatwiki/internal/app/client_side_build/define"
	"chatwiki/internal/app/client_side_build/initialize"
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
	//pprof api
	go func() {
		err := http.ListenAndServe(":55559", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	//consumer start
	common.RunTask(lib_define.ClientSideBuildTopic, lib_define.ClientsidebuildchannelWin, 1, business.ClientsidebuildWin)
}

func Stop() {
	define.ConsumerHandle.Stop()
}
