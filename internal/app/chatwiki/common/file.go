// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/zhimaAi/go_tools/tool"
)

func WriteFileByString(objectKey, content string) (string, error) {
	if err := tool.WriteFile(define.UploadDir+objectKey, content); err != nil {
		return ``, err
	}
	return define.LocalUploadPrefix + objectKey, nil
}

func GetFileByLink(link string) string {
	if len(link) >= len(define.LocalUploadPrefix) {
		return define.UploadDir + link[len(define.LocalUploadPrefix):]
	}
	return ``
}
