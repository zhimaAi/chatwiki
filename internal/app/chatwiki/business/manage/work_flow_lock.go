// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

const workFlowDraftLockLease = 120 * time.Second

type workFlowDraftLockOwner struct {
	LoginUserId   int    `json:"login_user_id"`
	RobotKey      string `json:"robot_key"`
	RemoteAddr    string `json:"remote_addr"`
	UniIdentifier string `json:"uni_identifier"`
	LeaseToken    string `json:"lease_token"`
	UserAgent     string `json:"user_agent"`
}

type workFlowDraftLockResult struct {
	Acquired bool
	LockTtl  int64
	ExpireAt int64
	Owner    workFlowDraftLockOwner
}

// The scripts calculate an absolute deadline from Redis TIME. A heartbeat always
// replaces the deadline with now + lease; it never adds time to the current TTL.
const acquireWorkFlowDraftLockScript = `
local now = redis.call('TIME')
local nowMs = (tonumber(now[1]) * 1000) + math.floor(tonumber(now[2]) / 1000)
local leaseMs = tonumber(ARGV[5])
local expireAt = nowMs + leaseMs
local current = redis.call('GET', KEYS[1])

if not current then
    redis.call('SET', KEYS[1], ARGV[1])
    redis.call('PEXPIREAT', KEYS[1], expireAt)
    return {1, ARGV[1], leaseMs, expireAt}
end

local decodedOk, decoded = pcall(cjson.decode, current)
if decodedOk
    and tostring(decoded.login_user_id) == ARGV[2]
    and tostring(decoded.uni_identifier) == ARGV[3] then
    redis.call('SET', KEYS[1], ARGV[1])
    redis.call('PEXPIREAT', KEYS[1], expireAt)
    return {1, ARGV[1], leaseMs, expireAt}
end

local ttl = redis.call('PTTL', KEYS[1])
local currentExpireAt = 0
if ttl > 0 then
    currentExpireAt = nowMs + ttl
end
return {0, current, ttl, currentExpireAt}
`

const renewWorkFlowDraftLockScript = `
local now = redis.call('TIME')
local nowMs = (tonumber(now[1]) * 1000) + math.floor(tonumber(now[2]) / 1000)
local leaseMs = tonumber(ARGV[4])
local expireAt = nowMs + leaseMs
local current = redis.call('GET', KEYS[1])

if not current then
    return {0, '', 0, 0}
end

local decodedOk, decoded = pcall(cjson.decode, current)
if decodedOk
    and tostring(decoded.login_user_id) == ARGV[1]
    and tostring(decoded.uni_identifier) == ARGV[2]
    and tostring(decoded.lease_token) == ARGV[3] then
    redis.call('PEXPIREAT', KEYS[1], expireAt)
    return {1, current, leaseMs, expireAt}
end

local ttl = redis.call('PTTL', KEYS[1])
local currentExpireAt = 0
if ttl > 0 then
    currentExpireAt = nowMs + ttl
end
return {0, current, ttl, currentExpireAt}
`

const releaseWorkFlowDraftLockScript = `
local current = redis.call('GET', KEYS[1])
if not current then
    return 1
end

local decodedOk, decoded = pcall(cjson.decode, current)
if decodedOk
    and tostring(decoded.login_user_id) == ARGV[1]
    and tostring(decoded.uni_identifier) == ARGV[2]
    and tostring(decoded.lease_token) == ARGV[3] then
    redis.call('DEL', KEYS[1])
    return 1
end
return 0
`

func getWorkFlowDraftLockRedisKey(adminUserId int, robotKey string) string {
	// Preserve the existing key bytes. The old %s formatting of an int produced
	// "%!s(int=<id>)" and changing it during rollout would create split locks.
	legacyAdminUserId := fmt.Sprintf("%%!s(int=%d)", adminUserId)
	lockIdentity := fmt.Sprintf("user_id:%s,robot:%s,", legacyAdminUserId, robotKey)
	return define.LockPreKey + ".draft_lock." + tool.MD5(lockIdentity)
}

func buildWorkFlowDraftLockOwner(c *gin.Context, robotKey, uniIdentifier, leaseToken, userAgent string) workFlowDraftLockOwner {
	return workFlowDraftLockOwner{
		LoginUserId:   getLoginUserId(c),
		RobotKey:      robotKey,
		RemoteAddr:    lib_web.GetClientIP(c),
		UniIdentifier: strings.TrimSpace(uniIdentifier),
		LeaseToken:    strings.TrimSpace(leaseToken),
		UserAgent:     strings.TrimSpace(userAgent),
	}
}

func validateWorkFlowDraftLockOwner(owner workFlowDraftLockOwner) error {
	if owner.LoginUserId <= 0 || owner.UniIdentifier == "" || owner.LeaseToken == "" {
		return errors.New("invalid workflow draft lock owner")
	}
	return nil
}

func acquireWorkFlowDraftLock(ctx context.Context, adminUserId int, owner workFlowDraftLockOwner) (workFlowDraftLockResult, error) {
	result := workFlowDraftLockResult{}
	if err := validateWorkFlowDraftLockOwner(owner); err != nil {
		return result, err
	}
	lockValue, err := json.Marshal(owner)
	if err != nil {
		return result, err
	}
	redisResult, err := define.Redis.Eval(
		ctx,
		acquireWorkFlowDraftLockScript,
		[]string{getWorkFlowDraftLockRedisKey(adminUserId, owner.RobotKey)},
		string(lockValue),
		strconv.Itoa(owner.LoginUserId),
		owner.UniIdentifier,
		owner.LeaseToken,
		workFlowDraftLockLease.Milliseconds(),
	).Result()
	if err != nil {
		return result, err
	}
	return parseWorkFlowDraftLockResult(redisResult)
}

func renewWorkFlowDraftLock(ctx context.Context, adminUserId int, owner workFlowDraftLockOwner) (workFlowDraftLockResult, error) {
	result := workFlowDraftLockResult{}
	if err := validateWorkFlowDraftLockOwner(owner); err != nil {
		return result, err
	}
	redisResult, err := define.Redis.Eval(
		ctx,
		renewWorkFlowDraftLockScript,
		[]string{getWorkFlowDraftLockRedisKey(adminUserId, owner.RobotKey)},
		strconv.Itoa(owner.LoginUserId),
		owner.UniIdentifier,
		owner.LeaseToken,
		workFlowDraftLockLease.Milliseconds(),
	).Result()
	if err != nil {
		return result, err
	}
	return parseWorkFlowDraftLockResult(redisResult)
}

func releaseWorkFlowDraftLock(ctx context.Context, adminUserId int, owner workFlowDraftLockOwner) (bool, error) {
	if err := validateWorkFlowDraftLockOwner(owner); err != nil {
		return false, err
	}
	result, err := define.Redis.Eval(
		ctx,
		releaseWorkFlowDraftLockScript,
		[]string{getWorkFlowDraftLockRedisKey(adminUserId, owner.RobotKey)},
		strconv.Itoa(owner.LoginUserId),
		owner.UniIdentifier,
		owner.LeaseToken,
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func parseWorkFlowDraftLockResult(redisResult any) (workFlowDraftLockResult, error) {
	result := workFlowDraftLockResult{}
	values, ok := redisResult.([]interface{})
	if !ok || len(values) != 4 {
		return result, errors.New("invalid workflow draft lock response")
	}
	result.Acquired = cast.ToInt(values[0]) == 1
	ttlMillis := cast.ToInt64(values[2])
	if ttlMillis > 0 {
		result.LockTtl = (ttlMillis + 999) / 1000
	}
	expireAtMillis := cast.ToInt64(values[3])
	if expireAtMillis > 0 {
		result.ExpireAt = expireAtMillis / 1000
	}
	storedValue := cast.ToString(values[1])
	if storedValue != "" {
		_ = json.Unmarshal([]byte(storedValue), &result.Owner)
	}
	return result, nil
}

func getWorkFlowDraftLockResultData(result workFlowDraftLockResult) map[string]any {
	data := map[string]any{
		"lock_res":       result.Acquired,
		"is_self":        result.Acquired,
		"lock_ttl":       result.LockTtl,
		"lock_expire_at": result.ExpireAt,
	}
	if !result.Acquired {
		data["lock_conflict"] = 1
		data["remote_addr"] = result.Owner.RemoteAddr
		data["user_agent"] = result.Owner.UserAgent
		staffUserName := lib_define.UnknownUser
		userInfo, err := GetUserInfo(cast.ToString(result.Owner.LoginUserId))
		if err == nil && userInfo["user_name"] != "" {
			staffUserName = userInfo["user_name"]
		}
		data["login_user_name"] = staffUserName
	}
	return data
}

func responseWorkFlowDraftLockConflict(c *gin.Context, result workFlowDraftLockResult) {
	data := getWorkFlowDraftLockResultData(result)
	c.String(http.StatusOK, lib_web.FmtJson(
		data,
		errors.New(i18n.Show(common.GetLang(c), `no_draft_edit_permission`, data["login_user_name"], result.Owner.RemoteAddr)),
	))
}

func checkAndRenewWorkFlowDraftLock(c *gin.Context, adminUserId int, owner workFlowDraftLockOwner) bool {
	result, err := renewWorkFlowDraftLock(c.Request.Context(), adminUserId, owner)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return false
	}
	if !result.Acquired {
		responseWorkFlowDraftLockConflict(c, result)
		return false
	}
	return true
}

func HeartbeatDraftKey(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	owner := buildWorkFlowDraftLockOwner(
		c,
		robotKey,
		c.PostForm(`uni_identifier`),
		c.PostForm(`lease_token`),
		c.PostForm(`user_agent`),
	)
	if err := validateWorkFlowDraftLockOwner(owner); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `lease_token`))))
		return
	}
	result, err := renewWorkFlowDraftLock(c.Request.Context(), adminUserId, owner)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(getWorkFlowDraftLockResultData(result), nil))
}

func ReleaseDraftKey(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	owner := buildWorkFlowDraftLockOwner(
		c,
		robotKey,
		c.PostForm(`uni_identifier`),
		c.PostForm(`lease_token`),
		c.PostForm(`user_agent`),
	)
	if err := validateWorkFlowDraftLockOwner(owner); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `lease_token`))))
		return
	}
	released, err := releaseWorkFlowDraftLock(c.Request.Context(), adminUserId, owner)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{"released": released}, nil))
}
