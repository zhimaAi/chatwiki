// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
