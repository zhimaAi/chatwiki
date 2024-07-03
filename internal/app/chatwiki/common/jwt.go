// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"

	"github.com/dgrijalva/jwt-go/v4"
)

var (
	client = &lib_web.JwtToken{}
)

func newClient() *lib_web.JwtToken {
	client = lib_web.NewTokenClient(define.JwtTtl, define.JwtKey)
	return client
}
func GetToken(userId, userName, parentId any) (jwt.MapClaims, error) {
	return newClient().GetToken(userId, userName, parentId)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	return newClient().ParseToken(tokenString)
}
