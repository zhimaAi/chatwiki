// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const NoRoute = `NoRoute`
const NoMethod = `NoMethod`

type Action = func(c *gin.Context)
type Route map[string]map[string]Action

// InitGin 初始化Gin服务器，现在支持添加额外的中间件
func InitGin(port string, route Route, auth Action, extraMiddlewares ...gin.HandlerFunc) *http.Server {
	handler := gin.Default()
	handler.Use(gin.Recovery())

	// 应用额外的中间件
	for _, middleware := range extraMiddlewares {
		handler.Use(middleware)
	}

	routeMap := make(map[string]string)
	for method, routers := range route {
		for path, action := range routers {
			switch method {
			case http.MethodGet, http.MethodPost, http.MethodPut,
				http.MethodPatch, http.MethodHead, http.MethodOptions,
				http.MethodDelete, http.MethodConnect, http.MethodTrace:
				routeMap[path] = method
				// check is need auth
				if _, ok := NoAuthRouteMap[path]; ok {
					handler.Handle(method, path, action)
				} else {
					handler.Handle(method, path, auth, action)
				}
			case NoRoute:
				handler.NoRoute(action)
			case NoMethod:
				fallthrough
			default:
				handler.HandleMethodNotAllowed = true
				handler.NoMethod(action)
			}
		}
	}

	return &http.Server{Addr: `:` + port, Handler: handler}
}
