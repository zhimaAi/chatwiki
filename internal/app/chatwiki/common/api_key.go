// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"fmt"
	"strings"

	"chatwiki/internal/app/chatwiki/define"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetAuthorizationToken(robotKey string) string {
	return BuildMessageId(tool.MD5(tool.Random(20)), robotKey, tool.Time2Int())
}

func ParseAuthorizationToken(c *gin.Context) (msql.Params, error) {
	// get user info
	headers := c.GetHeader(`Authorization`)
	parts := strings.SplitN(headers, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer" && len(parts[1]) == 72) {
		return nil, fmt.Errorf("open_apikey_format_err")
	}
	token := strings.TrimSpace(parts[1])
	robotKey := ParseMessageId(token)
	if robotKey == "" {
		return nil, fmt.Errorf("open_apikey_failed")
	}
	// robot check
	if !CheckRobotKey(robotKey) {
		return nil, fmt.Errorf("open_apikey_failed")
	}
	// token expired check
	keyData, err := GetRobotApikeyInfo(robotKey)
	if err != nil {
		return nil, fmt.Errorf("sys_err")
	}
	if len(keyData) == 0 {
		return nil, fmt.Errorf("open_apikey_failed")
	}
	flag := false
	for _, item := range keyData {
		if cast.ToInt(item["status"]) != define.SwitchOn || cast.ToInt(item["expire_time"]) > 0 && cast.ToInt(item["expire_time"]) <= tool.Time2Int() {
			continue
		}
		if strings.TrimSpace(cast.ToString(item["key"])) == strings.TrimSpace(token) {
			flag = true
			break
		}
	}
	if !flag {
		return nil, fmt.Errorf("open_apikey_failed")
	}
	return msql.Params{
		"robot_key": robotKey,
	}, nil
}

func BuildMessageId(name, id string, createTime int) string {
	str := fmt.Sprintf("%s_%v_%d", name, id, createTime)
	return tool.Base64Encode(str)
}

func BuildOpenAiMsgId() string {
	str := fmt.Sprintf("%v_%d", tool.Random(20), tool.Time2Int())
	return "chatcmpl-" + tool.MD5(str)
}

func ParseMessageId(id string) string {
	data, err := tool.Base64Decode(id)
	if err != nil {
		return ""
	}
	message := strings.SplitN(data, "_", 3)
	if len(message) >= 3 {
		return cast.ToString(message[1])
	}
	return ""
}

func BuildOpenId(id string) string {
	str := fmt.Sprintf("%v_%d_%v", id, tool.Time2Int(), tool.MD5(tool.Random(20)))
	return tool.Base64Encode(str)
}
