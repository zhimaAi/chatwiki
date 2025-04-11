// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"archive/zip"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"github.com/zhimaAi/pdf"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/antchfx/xmlquery"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func GetFileByReader(reader *zip.ReadCloser, name string) (io.ReadCloser, error) {
	for _, file := range reader.File {
		if file.Name == name {
			return file.Open()
		}
	}
	return nil, fmt.Errorf("file not found: %s", name)
}

func GetNodeAttr(attr []xmlquery.Attr, space, Local string) (string, bool) {
	for _, item := range attr {
		if item.Name.Space == space && item.Name.Local == Local {
			return item.Value, true
		}
	}
	return ``, false
}

func GetDocxRels(reader *zip.ReadCloser) (rels map[string]string) {
	rels = make(map[string]string)
	rc, err := GetFileByReader(reader, `word/_rels/document.xml.rels`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	document, err := xmlquery.Parse(rc)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	xmlquery.FindEach(document, `//Relationship`, func(_ int, node *xmlquery.Node) {
		id, _ := GetNodeAttr(node.Attr, ``, `Id`)
		target, _ := GetNodeAttr(node.Attr, ``, `Target`)
		if len(id) > 0 && len(target) > 0 {
			rels[id] = target
		}
	})
	return
}

func GetImgByZip(reader *zip.ReadCloser, name string, userId int) (imgStr string, err error) {
	rc, err := GetFileByReader(reader, name)
	if err != nil {
		return
	}
	bs, err := io.ReadAll(rc)
	if err != nil {
		return
	}
	// save to file
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(name), `.`))
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_image`, tool.Date(`Ym`), tool.MD5(string(bs)), ext)
	imgUrl, err := WriteFileByString(objectKey, string(bs))
	if err != nil {
		logs.Error(err.Error())
		return
	}
	imgStr = fmt.Sprintf(`{{!!%s!!}}`, imgUrl)
	return
}

func DocxInfoExtract(name string, userId int) (result []string, err error) {
	reader, err := zip.OpenReader(name)
	if err != nil {
		return
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	rels := GetDocxRels(reader)
	rc, err := GetFileByReader(reader, `word/document.xml`)
	if err != nil {
		return
	}
	document, err := xmlquery.Parse(rc)
	if err != nil {
		return
	}
	xmlquery.FindEach(document, `//w:p`, func(_ int, wp *xmlquery.Node) {
		var temp string
		xmlquery.FindEach(wp, `//*`, func(_ int, node *xmlquery.Node) {
			if node.Prefix == `w` && node.Data == `t` {
				temp += node.InnerText()
			}
			if node.Prefix == `a` && node.Data == `blip` && len(node.Attr) > 0 {
				if id, ok := GetNodeAttr(node.Attr, `r`, `embed`); ok && len(rels[id]) > 0 {
					if imgStr, err := GetImgByZip(reader, `word/`+rels[id], userId); err == nil {
						temp += imgStr //图片信息
					} else {
						logs.Error(err.Error())
					}
				}
			}
		})
		result = append(result, temp)
	})
	return
}

type DocxInfoCacheBuildHandler struct {
	fileUrl string
	userId  int
}

func (h *DocxInfoCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.docx.info.%d.%s`, h.userId, tool.MD5(h.fileUrl))
}
func (h *DocxInfoCacheBuildHandler) GetCacheData() (any, error) {
	return DocxInfoExtract(GetFileByLink(h.fileUrl), h.userId)
}
func GetDocxInfoInfo(fileUrl string, userId int) (string, error) {
	result := make([]string, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &DocxInfoCacheBuildHandler{fileUrl: fileUrl, userId: userId}, &result, time.Hour*24)
	content := strings.Join(result, "\r\n")
	return content, err
}

func ReadDocx(fileUrl string, userId int) ([]define.DocSplitItem, int, error) {
	if !LinkExists(fileUrl) {
		return nil, 0, errors.New(`file not exist:` + fileUrl)
	}
	content, err := GetDocxInfoInfo(fileUrl, userId)
	if err != nil {
		return nil, 0, err
	}
	list := []define.DocSplitItem{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}

func ReadPdf(pdfUrl, lang string) ([]define.DocSplitItem, int, error) {
	if len(pdfUrl) == 0 {
		return nil, 0, errors.New(`file link cannot be empty`)
	}
	//read pdf
	file, reader, err := pdf.Open(GetFileByLink(pdfUrl))
	if err != nil {
		logs.Error(err.Error())
		return nil, 0, errors.New(i18n.Show(lang, `can_not_parse_pdf`))
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	//paging collection
	list := make([]define.DocSplitItem, 0)
	wordTotal := 0
	for num := 1; num <= reader.NumPage(); num++ {
		p := reader.Page(num)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(content) > 0 {
			wordTotal += utf8.RuneCountInString(content)
			list = append(list, define.DocSplitItem{PageNum: num, Content: content})
		}
	}
	return list, wordTotal, nil
}
