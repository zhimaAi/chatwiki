// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/zhimaAi/go_tools/tool"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type CrawlerResult struct {
	Page    CrawlerPage         `json:"page"`
	Article readability.Article `json:"article"`
}
type CrawlerPage struct {
	RawHtml    string `json:"html"`
	MainHtml   string `json:"main_html"`
	Screenshot string `json:"screenshot"`
}

func SaveUploadedFile(fileHeader *multipart.FileHeader, limitSize, userId int, saveDir string, allowExt []string) (*define.UploadInfo, error) {
	if fileHeader == nil {
		return nil, errors.New(`file header is nil`)
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(fileHeader.Filename), `.`))
	if !tool.InArrayString(ext, allowExt) {
		return nil, errors.New(ext + ` not allow`)
	}
	if fileHeader.Size > int64(limitSize) {
		return nil, errors.New(`file size too big`)
	}
	reader, err := fileHeader.Open()
	defer func(reader multipart.File) {
		_ = reader.Close()
	}(reader)
	bs, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(bs) == 0 {
		return nil, errors.New(`file content is empty`)
	}
	content := string(bs)
	md5Hash := tool.MD5(content)
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, saveDir, tool.Date(`Ym`), md5Hash, ext)
	link, err := WriteFileByString(objectKey, content)
	if err != nil {
		return nil, err
	}
	return &define.UploadInfo{Name: fileHeader.Filename, Size: fileHeader.Size, Ext: ext, Link: link}, nil
}

func SaveUploadedFileMulti(c *gin.Context, name string, limitSize, userId int, saveDir string, allowExt []string) ([]*define.UploadInfo, []string) {
	uploadInfos := make([]*define.UploadInfo, 0)
	uploadErrors := make([]string, 0)
	if c.Request.MultipartForm == nil || len(c.Request.MultipartForm.File) == 0 {
		return uploadInfos, uploadErrors
	}
	for _, fileHeader := range c.Request.MultipartForm.File[name] {
		uploadInfo, err := SaveUploadedFile(fileHeader, limitSize, userId, saveDir, allowExt)
		if err != nil {
			uploadErrors = append(uploadErrors, err.Error())
			continue
		}
		uploadInfos = append(uploadInfos, uploadInfo)
	}
	return uploadInfos, uploadErrors
}

func SaveUrlPage(userId int, url, saveDir string) (*define.UploadInfo, error) {

	// request crawler
	endPoint := "http://" + define.Config.Crawler[`host`] + ":" + define.Config.Crawler[`port`] + "/content"
	resp, err := common.HttpPost(endPoint, nil, nil, map[string]interface{}{"url": url})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		if errMsg, ok := errResp["error"]; ok {
			return nil, errors.New(errMsg)
		} else {
			return nil, errors.New("request " + url + " error")
		}
	}

	// parse response
	var crawlerResult CrawlerResult
	err = common.HttpDecodeResponse(resp, &crawlerResult)
	if err != nil {
		return nil, err
	}
	if len(crawlerResult.Article.TextContent) == 0 || len(crawlerResult.Article.Title) == 0 {
		return nil, errors.New("parse url " + url + " article failed")
	}

	// save to file
	md5Hash := tool.MD5(crawlerResult.Article.TextContent)
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, saveDir, tool.Date(`Ym`), md5Hash, "html")

	link, err := WriteFileByString(objectKey, crawlerResult.Article.TextContent)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(define.AppRoot + link)
	if err != nil {
		return nil, err
	}

	return &define.UploadInfo{Name: crawlerResult.Article.Title, Size: fileInfo.Size(), Ext: "html", Link: link, Online: true, DocUrl: url}, nil
}
