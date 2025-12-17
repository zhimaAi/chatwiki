// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package middlewares

import (
	"chatwiki/internal/app/chatwiki/common"
	"net/http"
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
		token := request.GetHeader(`token`)
		if len(token) == 0 {
			token = request.Query(`token`)
		}
		data, err := common.ParseToken(token)
		if err != nil || data == nil {
			common.FmtErrorWithCode(request, http.StatusUnauthorized, "user_no_login")
			return
		}
	}
}
