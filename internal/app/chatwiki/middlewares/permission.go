// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package middlewares

import (
	"chatwiki/internal/app/chatwiki/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func CasbinAuth() gin.HandlerFunc {
	return func(request *gin.Context) {
		// get params
		obj := strings.TrimRight(request.Request.URL.RequestURI(), "/")
		before, _, _ := strings.Cut(obj, "?")
		if before != "" && strings.ContainsAny(before, "/") {
			obj = before
		}
		// get user info
		data, err := common.ParseToken(request.GetHeader(`token`))
		if err != nil || data == nil {
			common.FmtError(request, "user_no_login")
			return
		}
	}
}
