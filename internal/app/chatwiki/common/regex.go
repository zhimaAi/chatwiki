// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import "regexp"

func TakeUrls(str *string) []string {
	// 标准HTTP/HTTPS链接
	regex := regexp.MustCompile(`https?://(?:[-\w]+\.)+[-\w]+(?:/[-\w\./\?%&=]*)?`)
	return regex.FindAllString(*str, -1)
}
