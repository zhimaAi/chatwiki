// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/tool"
)

func IsUrl(s string) bool {
	return strings.HasPrefix(s, `https://`) || strings.HasPrefix(s, `http://`)
}

func CheckUsePrivate(objectKey string) bool {
	return regexp.MustCompile(`chat_ai/\d+/library_file/\d+/`).MatchString(objectKey)
}

func GetOssOptions(objectKey string) []oss.Option {
	if CheckUsePrivate(objectKey) {
		return []oss.Option{oss.ObjectACL(oss.ACLPrivate)}
	}
	return nil
}

func DownloadFile(fileUrl, temFile string) error {
	if cast.ToUint(define.Config.OssConfig[`enable`]) > 0 && CheckUsePrivate(fileUrl) {
		if temp, err := url.Parse(fileUrl); err == nil && len(temp.Path) > 0 {
			if err = GetObjectToFile(temp.Path[1:], temFile); err == nil {
				return nil
			}
		}
	}
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

func GetObjectToFile(objectKey, temFile string) error {
	client, err := oss.New(define.Config.OssConfig[`endpoint_internal`], define.Config.OssConfig[`keyid`], define.Config.OssConfig[`secret`])
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(define.Config.OssConfig[`bucket`])
	if err != nil {
		return err
	}
	err = bucket.GetObjectToFile(objectKey, temFile)
	if err != nil && tool.IsFile(temFile) {
		_ = os.Remove(temFile) //delete error file
	}
	return err
}

func PutObjectFromFile(objectKey, filePath string) (string, error) {
	client, err := oss.New(define.Config.OssConfig[`endpoint_internal`], define.Config.OssConfig[`keyid`], define.Config.OssConfig[`secret`])
	if err != nil {
		return ``, err
	}
	bucket, err := client.Bucket(define.Config.OssConfig[`bucket`])
	if err != nil {
		return ``, err
	}
	err = bucket.PutObjectFromFile(objectKey, filePath, GetOssOptions(objectKey)...)
	if err != nil {
		return ``, err
	}
	return "https://" + define.Config.OssConfig[`bucket`] + "." + define.Config.OssConfig[`endpoint`] + "/" + objectKey, nil
}

func PutObjectFromString(objectKey, content string) (string, error) {
	client, err := oss.New(define.Config.OssConfig[`endpoint_internal`], define.Config.OssConfig[`keyid`], define.Config.OssConfig[`secret`])
	if err != nil {
		return ``, err
	}
	bucket, err := client.Bucket(define.Config.OssConfig[`bucket`])
	if err != nil {
		return ``, err
	}
	err = bucket.PutObject(objectKey, strings.NewReader(content), GetOssOptions(objectKey)...)
	if err != nil {
		return ``, err
	}
	return "https://" + define.Config.OssConfig[`bucket`] + "." + define.Config.OssConfig[`endpoint`] + "/" + objectKey, nil
}

// GetLinkByFile is the inverse of GetFileByLink: it restores a local file path back to a
// publicly fetchable URL, for channels whose peer server (e.g. Aliyun ChatApp for WhatsApp)
// pulls the media itself. A local copy is still kept as fallback; on send the local path is
// turned into an api_domain public link. Returns input as-is if it is already a URL or unrecognized.
func GetLinkByFile(file string) string {
	if len(file) == 0 || IsUrl(file) {
		return file
	}
	apiDomain := strings.TrimRight(define.Config.WebService[`api_domain`], `/`)
	//download cache & export file: static/public/xxx -> api_domain + /public/xxx
	if strings.HasPrefix(file, `static/`) {
		return apiDomain + strings.TrimPrefix(file, `static`)
	}
	//upload file: define.UploadDir + key -> api_domain + LocalUploadPrefix + key
	if strings.HasPrefix(file, define.UploadDir) {
		return apiDomain + define.LocalUploadPrefix + file[len(define.UploadDir):]
	}
	return file
}

// GetCDNLinkByFile is the yun-only variant of GetLinkByFile: when OSS is enabled it returns an
// OSS + CDN link (so the peer server can pull the media), otherwise it falls back to GetLinkByFile.
// Kept separate so GetLinkByFile stays identical to main and merges don't conflict.
func GetCDNLinkByFile(file string) string {
	return GetLinkByFile(file)
}

func GetUrlPath(urlStr string) string {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return ``
	}
	return parsedUrl.Path
}
