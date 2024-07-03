// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/zhimaAi/go_tools/logs"
)

func WebRun(ws *http.Server) {
	logs.Info(`webservice listen:%s`, ws.Addr)
	err := ws.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logs.Info(`webservice closed`)
	} else if err != nil {
		logs.Error(err.Error())
		panic(`webservice listening error`)
	}
}

func Shutdown(ws *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ws.Shutdown(ctx); err != nil {
		logs.Error(err.Error())
		panic(`webservice close error`)
	}
	logs.Info(`webservice shutdown`)
}
