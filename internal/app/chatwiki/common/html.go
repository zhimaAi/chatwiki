// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// ProcessHTMLImages 把html中的网络图片下载转换成本地图片
func ProcessHTMLImages(content string, userId int) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	var wg sync.WaitGroup

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		if isHTTPURL(src) {
			wg.Add(1)
			go func(sel *goquery.Selection, url string) {
				defer wg.Done()
				if err := handleNetworkImage(sel, url, userId); err != nil {
					logs.Error(err.Error())
				}
			}(s, src)
			return
		}
		if strings.HasPrefix(src, "data:image") {
			handleBase64Image(s, src, userId)
			return
		}
	})

	wg.Wait()

	out, err := doc.Html()
	if err != nil {
		return "", err
	}

	if !utf8.ValidString(out) {
		out = tool.Convert(out, "gbk", "utf-8")
	}

	return out, nil
}

// handleNetworkImage 下载网络图片并替换成本地图片
func handleNetworkImage(sel *goquery.Selection, src string, userId int) error {
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 判断 MIME，决定扩展名
	mimeType := resp.Header.Get("Content-Type")
	ext := mimeToExt(mimeType)
	if ext == "" {
		ext = filepath.Ext(src)
		if ext == "" {
			ext = ".png"
		}
	}

	key := fmt.Sprintf("chat_ai/%d/%s/%s/%s%s", userId, "library_image", tool.Date("Ym"), tool.MD5(string(data)), ext)
	imgURL, err := WriteFileByString(key, string(data))
	if err != nil {
		return err
	}

	sel.ReplaceWithHtml(fmt.Sprintf("<b>{{!!%s!!}}</b>", imgURL))
	return nil
}

// handleBase64Image 把base64格式图片下载到本地
func handleBase64Image(sel *goquery.Selection, src string, userId int) error {
	parts := strings.Split(src, ";")
	if len(parts) < 2 {
		return fmt.Errorf("invalid base64 image")
	}

	format := strings.TrimPrefix(parts[0], "data:image/")
	if format == "svg+xml" {
		return nil // 不处理 svg
	}

	b64 := strings.TrimPrefix(parts[1], "base64,")
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	// webp → png
	if format == "webp" {
		data, err = ConvertWebPToPNG(data)
		if err != nil {
			return err
		}
		format = "png"
	}

	key := fmt.Sprintf("chat_ai/%d/%s/%s/%s.%s", userId, "library_image", tool.Date("Ym"), tool.MD5(string(data)), format)
	imgURL, err := WriteFileByString(key, string(data))
	if err != nil {
		return err
	}

	sel.ReplaceWithHtml(fmt.Sprintf("<b>{{!!%s!!}}</b>", imgURL))
	return nil
}

func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func mimeToExt(m string) string {
	if m == "" {
		return ""
	}
	exts, _ := mime.ExtensionsByType(m)
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}
