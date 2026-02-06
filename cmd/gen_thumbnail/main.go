// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package main

import (
	"chatwiki/internal/app/gen_thumbnail"
	"chatwiki/internal/app/gen_thumbnail/define"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhimaAi/go_tools/logs"
)

func main() {
	logs.SetLogsDir(define.AppRoot + `logs`)
	gen_thumbnail.Run()
	sc := make(chan os.Signal)
	sl := []os.Signal{
		syscall.SIGHUP,  //hangup
		syscall.SIGINT,  //interrupt
		syscall.SIGTERM, //terminated
	}
	signal.Notify(sc, sl...)
	sig := <-sc
	logs.Info(`sign:` + sig.String())
	gen_thumbnail.Stop()
	logs.Info(`stop finish`)
}
