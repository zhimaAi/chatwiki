// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	netURL "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/playwright-community/playwright-go"
	"github.com/zhimaAi/go_tools/logs"
)

type PageInfo struct {
	RawHtml    string `json:"html"`
	MainHtml   string `json:"main_html"`
	Screenshot string `json:"screenshot"`
}

type NeedWaitPageFunc func(page *playwright.Page) error
type SpecialPageProcessingFunc func(page *playwright.Page) error

const scroll = `
async (args) => {
    const {direction, speed} = args;
    const delay = ms => new Promise(resolve => setTimeout(resolve, ms));
    const scrollHeight = () => document.body.scrollHeight;
    const start = direction === "down" ? 0 : scrollHeight();
    const shouldStop = (position) => direction === "down" ? position > scrollHeight() : position < 0;
    const increment = direction === "down" ? 100 : -100;
    const delayTime = speed === "slow" ? 50 : 10;
    console.error(start, shouldStop(start), increment)
    for (let i = start; !shouldStop(i); i += increment) {
        window.scrollTo(0, i);
        await delay(delayTime);
    }
};
`

func fetchURLContent(parsedURL *netURL.URL) (*PageInfo, error) {
	// check if semaphore is full
	select {
	case concurrent <- struct{}{}:
	default:
		return nil, TooManyRequestsError
	}
	defer func() { <-concurrent }() // release semaphore

	// Create an isolated context
	acceptDownloads := false
	permissions := []string{"clipboard-read", "clipboard-write"}
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Permissions:     permissions,
		AcceptDownloads: &acceptDownloads,
		UserAgent:       &userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create context: %v", err)
	}
	// create a new page
	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	defer func(page *playwright.Page, context *playwright.BrowserContext) {
		_ = (*context).Close()
		_ = (*page).Close()
	}(&page, &context)

	// navigate to url, with timeout 15s
	if _, err = page.Goto(parsedURL.String(), playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateLoad}); err != nil {
		return nil, fmt.Errorf("could not navigate to url: %v", err)
	}

	// scroll to bottom to wait all images loaded
	_, err = page.Evaluate(scroll, map[string]interface{}{"direction": "down", "speed": "fast"})
	if err != nil {
		return nil, fmt.Errorf("could not scroll to bottom: %v", err)
	}

	// wait for element
	requestHostPath := parsedURL.Host + parsedURL.Path
	for baseURL, processor := range needWaitPages {
		matched, _ := regexp.MatchString(baseURL, requestHostPath)
		if matched {
			err := processor(&page)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	// save raw html
	var pageInfo PageInfo
	pageInfo.RawHtml, err = page.Content()
	if err != nil {
		return nil, nil
	}

	// take screenshot
	takeScreenshot(&page, &pageInfo)

	//var mainHtml string
	requestHostPath = parsedURL.Host + parsedURL.Path
	for baseURL, processor := range specialPageProcessors {
		matched, _ := regexp.MatchString(baseURL, requestHostPath)
		if matched {
			err := processor(&page)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	pageInfo.MainHtml, err = page.Content()
	if err != nil {
		return nil, err
	}

	return &pageInfo, nil
}

func takeScreenshot(page *playwright.Page, pageInfo *PageInfo) {
	fullPage := true
	screenshotBytes, err := (*page).Screenshot(playwright.PageScreenshotOptions{
		FullPage: &fullPage,
	})
	if err != nil {
		logs.Error("could not take screenshot: %v", err)
		return
	}

	pageInfo.Screenshot = base64.StdEncoding.EncodeToString(screenshotBytes)
}

var needWaitPages = map[string]NeedWaitPageFunc{
	"channels.weixin.qq.com/shop/learning-center": NeedWaitPageForChannelWeixin,
	"juejin.cn/post": NeedWaitPageForJuejin,
	"yuque.com":      NeedWaitPageForYuque,
}

func NeedWaitPageForChannelWeixin(page *playwright.Page) error {
	state := playwright.WaitForSelectorState("visible")
	timeout := float64(5000)
	err := (*page).Locator("#exeditor-preview").WaitFor(playwright.LocatorWaitForOptions{State: &state, Timeout: &timeout})
	if err != nil {
		return err
	}
	return nil
}

func NeedWaitPageForJuejin(page *playwright.Page) error {
	state := playwright.WaitForSelectorState("visible")
	timeout := float64(10000)
	err := (*page).Locator("article").WaitFor(playwright.LocatorWaitForOptions{State: &state, Timeout: &timeout})
	if err != nil {
		return err
	}
	return nil
}

func NeedWaitPageForYuque(page *playwright.Page) error {
	state := playwright.WaitForSelectorState("visible")
	timeout := float64(10000)
	err := (*page).Locator("#doc-reader-content").WaitFor(playwright.LocatorWaitForOptions{State: &state, Timeout: &timeout})
	if err != nil {
		return err
	}
	return nil
}

var specialPageProcessors = map[string]SpecialPageProcessingFunc{
	"docs.qq.com/doc": PageProcessForTencentDoc,
	"kdocs.cn":        PageProcessForKDoc,
	"feishu.cn":       PageProcessForFeishu,
	"www.jianshu.com": PageProcessForJianshu,
	"yuque.com":       PageProcessForYuque,
}

// escapeJSContent escapes content for safe JavaScript evaluation
func escapeJSContent(content string) string {
	content = strings.ReplaceAll(content, "`", "\\`")
	content = strings.ReplaceAll(content, "${", "\\${")
	return content
}

func PageProcessForTencentDoc(page *playwright.Page) error {
	waitCondition := `
        () => window.pad &&
              window.pad.contextService &&
              window.pad.contextService.context &&
              window.pad.contextService.context.editor &&
              typeof window.pad.contextService.context.editor.run === 'function'
    `
	_, err := (*page).WaitForFunction(waitCondition, nil, playwright.PageWaitForFunctionOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	time.Sleep(3 * time.Second)
	// window.pad.editor.readContext(0) 也可以获取纯文本内容
	r, err := (*page).Evaluate(`
let editor = window.pad.contextService.context.editor
editor.run("selectAll")
editor.clipboardManager.copyInterface.getCopyContent().html
`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	content, ok := r.(string)
	if !ok {
		return errors.New("could not convert content to string")
	}

	// 处理HTML中的图片转为base64
	content = processTencentDocImages(content)

	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
		return err
	}
	return nil
}

// processTencentDocImages 处理HTML中的图片链接，将其转换为base64
func processTencentDocImages(htmlContent string) string {
	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches := imgRegex.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		imgURL := match[1]
		// 如果已经是base64格式，则跳过
		if strings.HasPrefix(imgURL, "data:") {
			continue
		}

		// 下载图片
		imgData, mimeType, err := downloadImage(imgURL)
		if err != nil {
			logs.Error("下载图片失败: %v", err)
			continue
		}

		// 转换为base64
		base64Img := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imgData))

		// 替换原始URL为base64
		oldImgTag := fmt.Sprintf(`src="%s"`, imgURL)
		newImgTag := fmt.Sprintf(`src="%s"`, base64Img)
		htmlContent = strings.Replace(htmlContent, oldImgTag, newImgTag, -1)

		oldImgTag = fmt.Sprintf(`src='%s'`, imgURL)
		newImgTag = fmt.Sprintf(`src='%s'`, base64Img)
		htmlContent = strings.Replace(htmlContent, oldImgTag, newImgTag, -1)
	}

	return htmlContent
}

// downloadImage 下载图片并返回字节数据和MIME类型
func downloadImage(url string) ([]byte, string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// 添加模拟浏览器的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://docs.qq.com/")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	// 发送GET请求下载图片
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("下载图片HTTP状态码错误: %d", resp.StatusCode)
	}

	// 获取MIME类型
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg" // 默认MIME类型
	}

	// 读取图片数据
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return imgData, mimeType, nil
}

func PageProcessForKDoc(page *playwright.Page) error {
	// 解析出文档URL后缀
	u, err := netURL.Parse((*page).URL())
	if err != nil {
		return err
	}
	path := u.Path
	parts := strings.Split(path, "/")
	lastSegment := ""
	if len(parts) > 0 {
		lastSegment = parts[len(parts)-1]
	}

	// 获取下载链接
	downloadUrl, err := GetKDocDownloadUrl(fmt.Sprintf("https://www.kdocs.cn/api/office/file/%s/download", lastSegment))
	if err != nil {
		log.Println(err)
		return err
	}

	// 下载文件
	downloadFilename := fmt.Sprintf("kdoc_%s_%d.docx", lastSegment, time.Now().Unix())
	defer func() {
		if err := os.Remove(downloadFilename); err != nil {
			log.Printf("删除文件 %s 失败: %v\n", downloadFilename, err)
		}
	}()
	err = downloadKdocFile(downloadUrl, downloadFilename)
	if err != nil {
		log.Println(err)
		return err
	}

	// 去读文档内容
	content, err := DocxInfoExtract(downloadFilename)
	if err != nil {
		log.Println(err)
		return err
	}

	// 拼接到html中
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	return nil
}

func PageProcessForFeishu(page *playwright.Page) error {
	_, _ = (*page).Evaluate(`document.querySelector(".doc-info-wrapper").remove()`)
	content, err := (*page).Locator(`div[data-content-editable-root=true]`).InnerHTML()
	if err != nil {
		return err
	}
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
		return err
	}
	return nil
}

func PageProcessForJianshu(page *playwright.Page) error {
	content, err := (*page).Locator(`//*[@id="__next"]/div[1]/div/div[1]/section[1]`).InnerHTML()
	if err != nil {
		return err
	}
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
		return err
	}
	return nil
}

func PageProcessForYuque(page *playwright.Page) error {
	content, err := (*page).Locator(`#content`).InnerHTML()
	if err != nil {
		return err
	}
	content = escapeJSContent(content)
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", content))
	if err != nil {
		return err
	}
	return nil
}

func GetKDocDownloadUrl(apiURL string) (string, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", fmt.Errorf("HTTP 请求返回非 OK 状态 (%v)，且读取错误信息失败: %w", resp.Status, readErr)
		}
		return "", fmt.Errorf("HTTP 请求返回非 OK 状态: %v, 响应体: %s", resp.Status, errorBody)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	type DownloadInfo struct {
		DownloadURL string `json:"download_url"`
		URL         string `json:"url"`
		Fize        int    `json:"fize"`
		Fver        int    `json:"fver"`
		Store       string `json:"store"`
	}
	var downloadInfo DownloadInfo
	err = json.Unmarshal(body, &downloadInfo)
	if err != nil {
		return "", fmt.Errorf("解析 JSON 数据失败: %w", err)
	}

	if downloadInfo.DownloadURL == "" {
		return "", fmt.Errorf("JSON 中未找到 download_url 字段或其值为空")
	}

	return downloadInfo.DownloadURL, nil
}

func downloadKdocFile(url string, filename string) error {
	// 创建一个文件用于保存下载内容
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Println(err)
		}
	}(out)

	// 发送 GET 请求下载文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("发送下载请求失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		errorBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("下载请求返回非 OK 状态 (%v)，且读取错误信息失败: %w", resp.Status, readErr)
		}
		return fmt.Errorf("下载请求返回非 OK 状态: %v, 响应体: %s", resp.Status, errorBody)
	}

	// 将响应体内容复制到文件中
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("复制响应体到文件失败: %w", err)
	}

	return nil
}

func DocxInfoExtract(name string) (string, error) {
	reader, err := zip.OpenReader(name)
	if err != nil {
		return "", err
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	rels := GetDocxRels(reader)
	rc, err := GetFileByReader(reader, `word/document.xml`)
	if err != nil {
		return "", err
	}
	document, err := xmlquery.Parse(rc)
	if err != nil {
		return "", err
	}
	var result []string
	xmlquery.FindEach(document, `//w:p`, func(_ int, wp *xmlquery.Node) {
		var temp string
		xmlquery.FindEach(wp, `//*`, func(_ int, node *xmlquery.Node) {
			if node.Prefix == `w` && node.Data == `t` {
				temp += "<p>" + node.InnerText() + "</p>"
			}
			if node.Prefix == `a` && node.Data == `blip` && len(node.Attr) > 0 {
				if id, ok := GetNodeAttr(node.Attr, `r`, `embed`); ok && len(rels[id]) > 0 {
					if imgStr, err := GetImgByZip(reader, `word/`+rels[id]); err == nil {
						temp += fmt.Sprintf(`<img src="%s"/>`, imgStr) //图片信息
					} else {
						logs.Error(err.Error())
					}
				}
			}
		})
		result = append(result, temp)
	})

	return strings.Join(result, "\r\n"), nil
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

func GetImgByZip(reader *zip.ReadCloser, name string) (imgStr string, err error) {
	rc, err := GetFileByReader(reader, name)
	if err != nil {
		return
	}
	bs, err := io.ReadAll(rc)
	if err != nil {
		return
	}
	base64String := base64.StdEncoding.EncodeToString(bs)
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(name), `.`))
	return fmt.Sprintf("data:image/%s;base64,%s", ext, base64String), nil
}
