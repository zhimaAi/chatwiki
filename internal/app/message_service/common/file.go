// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func IsUrl(s string) bool {
	return strings.HasPrefix(s, `https://`) || strings.HasPrefix(s, `http://`)
}

func DownloadFile(fileUrl, temFile string) error {
	request := curl.Get(fileUrl)
	response, err := request.Response()
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(`StatusCode:%d`, response.StatusCode))
	}
	err = request.ToFile(temFile)
	if err != nil && tool.IsFile(temFile) {
		_ = os.Remove(temFile) //delete error file
	}
	return err
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
