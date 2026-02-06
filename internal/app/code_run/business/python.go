// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/code_run/common"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func Python(c *gin.Context) {
	body, err := GetCodeRunBody(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	result, code, signal := ``, http.StatusOK, make(chan struct{})
	go asyncRunPython(body, &result, &err, signal)
	select {
	case <-c.Request.Context().Done():
		code = 499
	case <-signal:
	}
	c.String(code, lib_web.FmtJson(result, err))
}

func asyncRunPython(body lib_define.CodeRunBody, result *string, err *error, signal chan struct{}) {
	batchNo, over := tool.Random(20), false
	logs.Debug(`Python-%s:body:%s`, batchNo, tool.JsonEncodeNoError(body))
	time.AfterFunc(time.Second*10, func() {
		if !over {
			logs.Debug(`Python-%s:body:%s`, batchNo, tool.JsonEncodeNoError(body))
		}
	})
	temp := time.Now()
	*result, *err = common.RunPython(body.MainFunc, batchNo, body.Params)
	over = true //标记已完成
	logs.Debug(`Python-%s:result:%s,err:%v`, batchNo, *result, *err)
	logs.Debug(`Python-%s:time:%v`, batchNo, time.Now().Sub(temp).Milliseconds())
	signal <- struct{}{}
	close(signal)
}
