// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zhimaAi/go_tools/curl"
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

func PutObjectFromFile(objectKey, filePath string) (string, error) {
	client, err := oss.New(define.Config.OssConfig[`endpoint_internal`], define.Config.OssConfig[`keyid`], define.Config.OssConfig[`secret`])
	if err != nil {
		return ``, err
	}
	bucket, err := client.Bucket(define.Config.OssConfig[`bucket`])
	if err != nil {
		return ``, err
	}
	err = bucket.PutObjectFromFile(objectKey, filePath)
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
	err = bucket.PutObject(objectKey, strings.NewReader(content))
	if err != nil {
		return ``, err
	}
	return "https://" + define.Config.OssConfig[`bucket`] + "." + define.Config.OssConfig[`endpoint`] + "/" + objectKey, nil
}
