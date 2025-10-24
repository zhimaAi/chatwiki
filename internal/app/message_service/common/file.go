// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"

	"github.com/zhimaAi/go_tools/tool"
)

func WriteFileByString(objectKey, content string) (string, error) {
	if err := tool.WriteFile(define.UploadDir+objectKey, content); err != nil {
		return ``, err
	}
	return define.LocalUploadPrefix + objectKey, nil
}
