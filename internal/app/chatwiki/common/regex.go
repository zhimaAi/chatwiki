// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import "regexp"

func TakeUrls(str *string) []string {
	// Standard HTTP/HTTPS links
	regex := regexp.MustCompile(`https?://(?:[-\w]+\.)+[-\w]+(?:/[-\w\./\?%&=]*)?`)
	return regex.FindAllString(*str, -1)
}
