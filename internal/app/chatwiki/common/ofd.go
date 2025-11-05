// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"archive/zip"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/antchfx/xmlquery"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func GetOfdPages(reader *zip.ReadCloser) (pages []string) {
	pages = make([]string, 0)
	doc, err := GetFileByReader(reader, `Doc_0/Document.xml`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	document, err := xmlquery.Parse(doc)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	xmlquery.FindEach(document, `//ofd:Page`, func(_ int, node *xmlquery.Node) {
		page, _ := GetNodeAttr(node.Attr, ``, `BaseLoc`)
		if len(page) > 0 {
			pages = append(pages, page)
		}
	})
	return
}

func GetOfdRels(reader *zip.ReadCloser) (rels map[string]string) {
	rels = make(map[string]string)
	rc, err := GetFileByReader(reader, `Doc_0/DocumentRes.xml`)
	if err != nil {
		return
	}
	document, err := xmlquery.Parse(rc)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	xmlquery.FindEach(document, `//ofd:MultiMedia`, func(_ int, node *xmlquery.Node) {
		id, _ := GetNodeAttr(node.Attr, ``, `ID`)
		sub := xmlquery.FindOne(node, `ofd:MediaFile`)
		if sub == nil {
			return
		}
		mediaFile := sub.InnerText()
		if len(id) > 0 && len(mediaFile) > 0 {
			rels[id] = mediaFile
		}
	})
	return
}

func IsNewLine(minH, maxH *float64, boundary string) bool {
	temp := strings.Fields(boundary)
	if len(boundary) == 0 || len(temp) != 4 {
		return false //错误的boundary
	}
	start, end := cast.ToFloat64(temp[1]), cast.ToFloat64(temp[1])+cast.ToFloat64(temp[3])
	if *minH == 0 && *maxH == 0 { //首行位置
		*minH, *maxH = start, end
		return false
	}
	if *maxH < start || end < *minH { //开启新行
		*minH, *maxH = start, end
		return true
	}
	*minH, *maxH = min(*minH, start), max(*maxH, end)
	return false
}

func OfdInfoExtract(name string, userId int) (result []string, err error) {
	reader, err := zip.OpenReader(name)
	if err != nil {
		return
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	pages := GetOfdPages(reader)
	rels := GetOfdRels(reader)
	var content string
	for idx := range pages {
		if idx > 0 {
			content += "\r\n"
		}
		page, err := GetFileByReader(reader, `Doc_0/`+pages[idx])
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		document, err := xmlquery.Parse(page)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		var minH, maxH float64
		xmlquery.FindEach(document, `//ofd:Layer`, func(_ int, wp *xmlquery.Node) {
			xmlquery.FindEach(wp, `//*`, func(_ int, node *xmlquery.Node) {
				if node.Prefix == `ofd` && node.Data == `TextObject` {
					boundary, _ := GetNodeAttr(node.Attr, ``, `Boundary`)
					if IsNewLine(&minH, &maxH, boundary) && !strings.HasSuffix(content, "\r\n") {
						content += "\r\n"
					}
					sub := xmlquery.FindOne(node, `ofd:TextCode`)
					if sub != nil {
						content += sub.InnerText()
					}
				}
				if node.Prefix == `ofd` && node.Data == `ImageObject` && len(node.Attr) > 0 {
					if id, ok := GetNodeAttr(node.Attr, ``, `ResourceID`); ok && len(rels[id]) > 0 {
						if imgStr, err := GetImgByZip(reader, `Doc_0/Res/`+rels[id], userId); err == nil {
							content, _ = strings.CutSuffix(content, "\r\n")
							content += "\r\n" + imgStr + "\r\n" //图片信息
						} else {
							logs.Error(err.Error())
						}
					}
				}
			})
		})
	}
	result = strings.Split(content, "\r\n")
	return
}

type OfdInfoCacheBuildHandler struct {
	fileUrl string
	userId  int
}

func (h *OfdInfoCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.ofd.info.%d.%s`, h.userId, tool.MD5(h.fileUrl))
}
func (h *OfdInfoCacheBuildHandler) GetCacheData() (any, error) {
	return OfdInfoExtract(GetFileByLink(h.fileUrl), h.userId)
}
func GetOfdInfoInfo(fileUrl string, userId int) (string, error) {
	result := make([]string, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &OfdInfoCacheBuildHandler{fileUrl: fileUrl, userId: userId}, &result, time.Hour*24)
	content := strings.Join(result, "\r\n")
	return content, err
}

func ReadOfd(fileUrl string, userId int) (define.DocSplitItems, int, error) {
	if !LinkExists(fileUrl) {
		return nil, 0, errors.New(`file not exist:` + fileUrl)
	}
	content, err := GetOfdInfoInfo(fileUrl, userId)
	if err != nil {
		return nil, 0, err
	}
	list := define.DocSplitItems{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}
