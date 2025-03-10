// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const NoRoute = `NoRoute`
const NoMethod = `NoMethod`

type Action = func(c *gin.Context)
type Route map[string]map[string]Action

func InitGin(port string, route Route, auth Action) *http.Server {
	handler := gin.Default()
	handler.Use(gin.Recovery())
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
