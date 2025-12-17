// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"os"
	"time"

	"github.com/spf13/cast"
)

func GetSessionSecond() int {
	if define.IsDev {
		sessionSecond := cast.ToInt(os.Getenv(`session_second`))
		if sessionSecond > 0 {
			return sessionSecond
		}
	}
	return 86400 //24hour
}

func GetSessionTtl() time.Duration {
	return time.Duration(GetSessionSecond()) * time.Second
}
