// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"net"
	"net/http"
	"strings"
)

func GetRequestIP(r *http.Request) string {
	ip := r.Header.Get(`X-Forwarded-For`)
	if len(ip) == 0 {
		ip = r.Header.Get(`X-Real-IP`)
	}
	if len(ip) == 0 {
		ip = strings.TrimSpace(strings.Split(r.RemoteAddr, `:`)[0])
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ip
	}
	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		return ip
	}
	return ipv4.String()
}
