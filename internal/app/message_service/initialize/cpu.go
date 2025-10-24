// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/message_service/define"
	"runtime"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initNumCPU() {
	ratio, err := tool.String2Float64(define.Config.NumCPU["ratio"])
	if err != nil {
		logs.Error(err.Error())
		panic(`num_cpu ratio error`)
	}
	maximum := runtime.NumCPU()
	set := tool.Floor(float64(maximum) * ratio)
	old := runtime.GOMAXPROCS(set)
	logs.Info(`cur use cpus:` + tool.Int2String(old))
	logs.Info(`max use cpus:` + tool.Int2String(maximum))
	logs.Info(`set use cpus:` + tool.Int2String(set))
}
