// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"chatwiki/internal/app/chatwiki/define"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

const (
	ChatClawTokenStatusActive  = 1
	ChatClawTokenStatusRevoked = 2
)

func GetRequestToken(c *gin.Context) string {
	token := strings.TrimSpace(c.GetHeader(`token`))
	if token == "" {
		token = strings.TrimSpace(c.Query(`token`))
	}
	return token
}

func GetTokenSha256(token string) string {
	if strings.TrimSpace(token) == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func CheckChatClawTokenStatus(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("chatclaw token is empty")
	}
	tokenHash := GetTokenSha256(token)

	var (
		info msql.Params
		err  error
	)
	if tokenHash != "" {
		info, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
			Where("token_hash", tokenHash).
			Order("id desc").
			Find()
		if err != nil {
			info, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
				Where("token", token).
				Order("id desc").
				Find()
			if err != nil {
				return err
			}
		}
	}
	if len(info) == 0 {
		info, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
			Where("token", token).
			Order("id desc").
			Find()
		if err != nil {
			return err
		}
	}
	if len(info) == 0 {
		return errors.New("chatclaw token not found")
	}
	if cast.ToInt(info["status"]) == ChatClawTokenStatusRevoked {
		return errors.New("chatclaw token has been revoked")
	}
	return nil
}

func GetChatClawAuthClaims(c *gin.Context) (jwt.MapClaims, string, error) {
	token := GetRequestToken(c)
	if token == "" {
		return nil, "", errors.New("chatclaw token is empty")
	}
	claims, err := ParseChatClawToken(token)
	if err != nil {
		return nil, "", err
	}
	if cast.ToInt(claims["user_id"]) <= 0 {
		return nil, "", errors.New("chatclaw token user_id invalid")
	}
	if err = CheckChatClawTokenStatus(token); err != nil {
		return nil, "", err
	}
	return claims, token, nil
}
