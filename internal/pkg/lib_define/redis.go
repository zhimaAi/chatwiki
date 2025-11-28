// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_define

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

// 注:修改版本号时应用于:chatwiki,message_service两个项目
const (
	RedisPrefixRobotInfo           = `chatwiki.robot_info.v20251107.%s`
	RedisPrefixAppInfo             = `chatwiki.app_info.v20251107.%s.%s`
	RedisPrefixPassiveSubscribe    = `chatwiki.passive_subscribe.v20251107`
	RedisPrefixMediaUpload         = `chatwiki.media_upload.v20251107.%s.%s`
	RedisPrefixDingtalkAccessToken = `chatwiki.dingtalk.access_token.v20251106.%s.%s`
)
