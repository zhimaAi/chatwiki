// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func WriteFileByString(objectKey, content string) (string, error) {
	if cast.ToUint(define.Config.OssConfig[`enable`]) > 0 { //put oss
		if link, err := PutObjectFromString(objectKey, content); err == nil {
			return link, nil
		} else {
			logs.Error(err.Error())
		}
	}
	if err := tool.WriteFile(define.UploadDir+objectKey, content); err != nil {
		return ``, err
	}
	return define.LocalUploadPrefix + objectKey, nil
}

func GetFileByLink(link string) string {
	if IsUrl(link) { //download
		temFile := `static/public/download/` + tool.MD5(link) + strings.ToLower(filepath.Ext(link))
		if tool.IsFile(temFile) {
			return temFile //local exist
		}
		if err := DownloadFile(link, temFile); err != nil {
			logs.Error(err.Error())
			return ``
		}
		return temFile
	}
	//upload file
	if strings.HasPrefix(link, define.LocalUploadPrefix) {
		locFile := define.UploadDir + link[len(define.LocalUploadPrefix):]
		if tool.IsFile(locFile) {
			return locFile //local file
		}
	}
	//export file
	if strings.HasPrefix(link, `/public/export/`) {
		locFile := `static` + link
		if tool.IsFile(locFile) {
			return locFile //local file
		}
	}
	return ``
}

func LinkExists(link string) bool {
	return len(GetFileByLink(link)) > 0
}
