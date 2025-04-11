// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
	return 1800 //30minute
}

func GetSessionTtl() time.Duration {
	return time.Duration(GetSessionSecond()) * time.Second
}
