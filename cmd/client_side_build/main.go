// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package main

import (
	"chatwiki/internal/app/client_side_build"
	"chatwiki/internal/app/client_side_build/define"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhimaAi/go_tools/logs"
)

func main() {
	logs.SetLogsDir(define.AppRoot + `logs`)
	client_side_build.Run()

	sc := make(chan os.Signal)
	sl := []os.Signal{
		syscall.SIGHUP,  //hangup
		syscall.SIGINT,  //interrupt
		syscall.SIGTERM, //terminated
	}
	signal.Notify(sc, sl...)
	sig := <-sc
	logs.Info(`sign:` + sig.String())
	client_side_build.Stop()
	logs.Info(`stop finish`)
}
