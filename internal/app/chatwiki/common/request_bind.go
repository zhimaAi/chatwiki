// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestParamsBind only support POST and GET
func RequestParamsBind(req any, c *gin.Context, method ...string) error {
	requestMethod := c.Request.Method
	if len(method) > 0 {
		requestMethod = method[0]
	}
	if requestMethod == http.MethodPost {
		return c.ShouldBind(req)
	} else {
		return c.ShouldBindQuery(req)
	}
}
