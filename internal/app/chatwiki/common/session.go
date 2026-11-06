// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"math"
	"os"
	"time"

	"github.com/spf13/cast"
)

const (
	DefaultSessionSecond     = 24 * 60 * 60
	MinSessionSecond         = 6 * 60
	MaxSessionSecond         = 48 * 60 * 60
	sessionExpireHourEpsilon = 1e-9
)

// GetSessionSecond returns the configured inactivity timeout for an account.
// The optional account ID keeps callers that do not have account context
// compatible with the historical environment-variable fallback.
func GetSessionSecond(adminUserIds ...int) int {
	if sessionSecond, ok := devSessionSecond(); ok {
		return sessionSecond
	}
	if len(adminUserIds) > 0 && adminUserIds[0] > 0 {
		config := GetAdminConfig(adminUserIds[0])
		return normalizeSessionSecond(cast.ToInt(config[`session_expire_seconds`]))
	}
	return DefaultSessionSecond
}

// GetSessionTtl only keeps the session id cached. Dev still uses session_second.
// Otherwise the key lives for the longest allowed setting, so a 48 hour
// timeout is not cut off early. The account setting ends the session by last_chat_time.
func GetSessionTtl() time.Duration {
	if sessionSecond, ok := devSessionSecond(); ok {
		return time.Duration(sessionSecond) * time.Second
	}
	return MaxSessionSecond * time.Second
}

func SessionExpireSecondsFromHours(hours float64) (int, bool) {
	if math.IsNaN(hours) || math.IsInf(hours, 0) ||
		hours < 0.1-sessionExpireHourEpsilon || hours > 48+sessionExpireHourEpsilon {
		return 0, false
	}
	seconds := int(math.Round(hours * 60 * 60))
	if seconds < MinSessionSecond || seconds > MaxSessionSecond {
		return 0, false
	}
	return seconds, true
}

// SessionRemainingSeconds is the time left before a session ends.
// Zero means the session is already due, including the exact timeout boundary.
func SessionRemainingSeconds(lastChatTime, now, sessionSeconds int) int {
	if lastChatTime <= 0 || sessionSeconds <= 0 {
		return 0
	}
	remaining := lastChatTime + sessionSeconds - now
	if remaining <= 0 {
		return 0
	}
	return remaining
}

func devSessionSecond() (int, bool) {
	if !define.IsDev {
		return 0, false
	}
	sessionSecond := cast.ToInt(os.Getenv(`session_second`))
	if sessionSecond <= 0 {
		return 0, false
	}
	return sessionSecond, true
}

func normalizeSessionSecond(sessionSecond int) int {
	if sessionSecond >= MinSessionSecond && sessionSecond <= MaxSessionSecond {
		return sessionSecond
	}
	return DefaultSessionSecond
}
