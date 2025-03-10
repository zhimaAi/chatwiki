// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/app/chatwiki/route"
	"chatwiki/internal/pkg/lib_web"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initGin() {
	if !define.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}
	port := define.Config.WebService[`port`]
	define.WebService = lib_web.InitGin(port, route.Route, middlewares.CasbinAuth())
	//set the upload directory
	if handler, ok := define.WebService.Handler.(*gin.Engine); ok {
		handler.StaticFS(`/upload`, http.Dir(define.UploadDir))
	}
}
