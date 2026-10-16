// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func initJwtKey() {
	file := define.AppRoot + `data/jwt_key.conf`
	if tool.IsFile(file) {
		key, err := tool.ReadFile(file)
		if err != nil {
			logs.Error(`read JWT key file failed: %v`, err)
		} else {
			key = strings.TrimSpace(key)
			if len(key) == 64 {
				if _, err := hex.DecodeString(key); err == nil {
					define.JwtKey = key
					return
				}
			}
			logs.Error(`invalid JWT key file, generating a new key`)
		}
	}

	key := make([]byte, 32)
	rand.Read(key)
	define.JwtKey = hex.EncodeToString(key)
	if err := os.WriteFile(file, []byte(define.JwtKey), 0600); err != nil {
		logs.Error(`write JWT key file failed: %v`, err)
	}
}
