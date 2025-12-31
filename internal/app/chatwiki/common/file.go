// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
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

func GetUrlExt(rawURL string) string { //return jpg jpeg ...
	u, err := url.Parse(rawURL)
	if err != nil {
		logs.Error(`get url ext：%s error:%s`, rawURL, err.Error())
		return ``
	}
	ext := path.Ext(u.Path)
	if len(ext) > 1 {
		return strings.ToLower(ext[1:])
	}
	return ``
}

func Mp3ToAmr(adminUserId int, mp3Url string) (amr string) {
	domain := define.Config.WebService[`file_transfer`]
	locFile := ``
	if IsUrl(mp3Url) { //download
		temFile := `static/public/download/` + tool.MD5(mp3Url) + strings.ToLower(filepath.Ext(mp3Url))
		if tool.IsFile(temFile) {
			locFile = temFile //local exist
		} else {
			if err := DownloadFile(mp3Url, temFile); err != nil {
				logs.Error(err.Error())
			}
		}
	} else if strings.HasPrefix(mp3Url, define.LocalUploadPrefix) {
		locFile = define.UploadDir + mp3Url[len(define.LocalUploadPrefix):]
		if !tool.IsFile(locFile) {
			locFile = ``
		}
	} else if strings.HasPrefix(mp3Url, `/public/export/`) {
		locFile = `static` + mp3Url
		if !tool.IsFile(locFile) {
			locFile = ``
		}
	}
	if locFile == `` {
		logs.Warning(`not found local file from link %s`, mp3Url)
		return
	}
	request := curl.Post(domain+`/convert`).PostFile(`mp3`, locFile)
	resp, err := request.Bytes()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	key := fmt.Sprintf("chat_ai/%d/%s/%s/%s%s", adminUserId, "mp3_to_amr", tool.Date("Ym"), tool.MD5(string(resp)), `.amr`)
	amr, err = WriteFileByString(key, string(resp))
	if err != nil {
		logs.Error(err.Error())
		return
	}
	return amr
}
