// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/app/chatwiki/route"
	"chatwiki/internal/pkg/lib_web"
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

// go:embed static/*
var staticFiles embed.FS

func initGin() {
	if !define.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}
	port := define.Config.WebService[`port`]
	define.WebService = lib_web.InitGin(port, route.Route, middlewares.CasbinAuth())

	// Setup upload directory and MCP endpoint with multi-tenant support
	if handler, ok := define.WebService.Handler.(*gin.Engine); ok {
		handler.StaticFS(`/upload`, http.Dir(define.UploadDir))
		handler.LoadHTMLGlob(define.GetTemplatesPath() + `*.html`)
		handler.StaticFS(`/open/static`, http.Dir(define.GetTemplatesStaticPath()))

		// Register MCP endpoints with Bearer Token auth and dynamic tool loading
		mcpAuth := business.MCPAuthMiddleware()
		handler.POST(`/mcp`, mcpAuth, business.HandleMCPRequest)
		handler.GET(`/mcp`, mcpAuth, business.HandleMCPRequest)
		handler.DELETE(`/mcp`, mcpAuth, business.HandleMCPRequest)
	}
}
