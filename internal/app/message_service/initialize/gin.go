// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/app/message_service/middlewares"
	"chatwiki/internal/app/message_service/route"
	"chatwiki/internal/pkg/lib_web"

	"github.com/gin-gonic/gin"
)

func initGin() {
	if !define.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}
	port := define.Config.WebService[`port`]
	define.WebService = lib_web.InitGin(port, route.Route, middlewares.CasbinAuth())
}
