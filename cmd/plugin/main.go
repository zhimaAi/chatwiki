// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package main

import (
	"chatwiki/internal/app/plugin"
	"chatwiki/internal/app/plugin/define"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhimaAi/go_tools/logs"
)

func main() {
	logs.SetLogsDir(define.AppRoot + `logs`)
	plugin.Run()
	sc := make(chan os.Signal, 1)
	sl := []os.Signal{
		syscall.SIGHUP,  // hangup
		syscall.SIGINT,  // interrupt
		syscall.SIGTERM, // terminated
	}
	signal.Notify(sc, sl...)
	_ = <-sc
	plugin.Stop()
}
