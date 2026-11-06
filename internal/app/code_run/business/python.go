// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/code_run/common"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"context"
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
	type outcome struct {
		result string
		err    error
	}
	done := make(chan outcome, 1)
	ctx := c.Request.Context()
	go func() {
		result, err := runPython(ctx, body)
		done <- outcome{result, err}
	}()
	select {
	case <-c.Request.Context().Done():
		c.String(499, lib_web.FmtJson(nil, c.Request.Context().Err()))
	case out := <-done:
		c.String(http.StatusOK, lib_web.FmtJson(out.result, out.err))
	}
}

func runPython(ctx context.Context, body lib_define.CodeRunBody) (string, error) {
	batchNo := tool.Random(20)
	logs.Debug("Python-%s:body:%s", batchNo, tool.JsonEncodeNoError(body))
	start := time.Now()
	result, err := common.RunPythonContext(ctx, body.MainFunc, batchNo, body.Params)
	logs.Debug("Python-%s:result:%s,err:%v", batchNo, result, err)
	logs.Debug("Python-%s:time:%v", batchNo, time.Since(start).Milliseconds())
	return result, err
}
