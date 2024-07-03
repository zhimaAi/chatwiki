// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_web

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/zhimaAi/go_tools/tool"
)

type JwtToken struct {
	Ttl int    `json:"ttl"`
	Key string `json:"key"`
}

func NewTokenClient(ttl int, key string) *JwtToken {
	return &JwtToken{
		Ttl: ttl,
		Key: key,
	}
}
func (t JwtToken) GetToken(userId, userName, parentId any) (jwt.MapClaims, error) {
	data := jwt.MapClaims{
		"user_id": userId, "user_name": userName, "parent_id": parentId,
		`ttl`: t.Ttl, "exp": tool.Time2Int() + t.Ttl,
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString([]byte(t.Key))
	if err != nil {
		return nil, err
	}
	data[`token`] = tokenString
	return data, nil
}

func (t JwtToken) ParseToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if len(tokenString) == 0 {
		return nil, errors.New(`token is empty`)
	}
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(t.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(`token is invalid`)
	}
}
