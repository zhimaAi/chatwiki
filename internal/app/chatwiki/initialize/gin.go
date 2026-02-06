// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/app/chatwiki/route"
	"chatwiki/internal/pkg/lib_web"
	"embed"
	"net/http"

	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// go:embed static/*
var staticFiles embed.FS

func initGin() {
	if !define.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}
	port := define.Config.WebService[`port`]
	define.WebService = lib_web.InitGin(port, route.Route, middlewares.CasbinAuth(), middlewares.I18nPlaceholderMiddleware()) // add i18n placeholder middleware

	// Setup upload directory and MCP endpoint with multi-tenant support
	if handler, ok := define.WebService.Handler.(*gin.Engine); ok {
		handler.StaticFS(`/upload`, http.Dir(define.UploadDir))
		handler.LoadHTMLGlob(define.GetTemplatesPath() + `*.html`)
		handler.StaticFS(`/open/static`, http.Dir(define.GetTemplatesStaticPath()))

		// Register MCP endpoints with Bearer Token auth and dynamic tool loading
		mcpAuth := business.MCPAuthMiddleware()
		handler.Any(`/mcp`, mcpAuth, business.HandleMCPRequest)

		// Proxy Plugin Api
		target, _ := url.Parse(define.Config.Plugin[`endpoint`])
		proxy := httputil.NewSingleHostReverseProxy(target)
		handler.Any("/manage/plugin/*path", middlewares.CasbinAuth(), func(c *gin.Context) {
			adminUserId := common.GetAdminUserId(c)
			c.Request.Header.Set("admin_user_id", cast.ToString(adminUserId))
			c.Request.URL.Scheme = target.Scheme
			c.Request.URL.Host = target.Host
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}
