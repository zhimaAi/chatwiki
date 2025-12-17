// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_web

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
)

func GetRequestIP(c *gin.Context) string {
	reqIP := c.RemoteIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func GetPublicIp() string {
	var ipInfo struct {
		IP string `json:"ip"`
	}
	url := `https://api64.ipify.org/?format=json`
	err := curl.Get(url).SetTimeout(time.Second*3, time.Second*3).ToJSON(&ipInfo)
	if err != nil {
		logs.Error(err.Error())
	}
	if len(ipInfo.IP) == 0 {
		return `127.0.0.1`
	}
	return ipInfo.IP
}

func GetClientIP(c *gin.Context) string {
	// get from X-Forwarded-For
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// get first ip
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// get from X-Real-IP
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// if not exists get from request
	return GetRequestIP(c)
}

// 获取请求的域名
func GetRequestDomain(c *gin.Context) string {
	host := c.Request.Host
	// 如果域名中包含端口号，去除端口号部分
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	return host
}
