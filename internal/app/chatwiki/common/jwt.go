// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var (
	jwtClient = &lib_web.JwtToken{}
)

func newClient() *lib_web.JwtToken {
	return lib_web.NewTokenClient(define.JwtTtl, define.JwtKey)
}
func GetToken(userId, userName, parentId any) (jwt.MapClaims, error) {
	return newClient().GetToken(userId, userName, parentId)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	return newClient().ParseToken(tokenString)
}

func GetLoginUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := ParseToken(token)
	if err != nil {
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		return userId
	}
	return userId
}

func GetAdminUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := ParseToken(token)
	if err != nil {
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		return userId
	}
	if cast.ToInt(data["parent_id"]) <= 0 {
		return cast.ToInt(data["user_id"])
	}
	return cast.ToInt(data["parent_id"])
}
