// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"log"
	netURL "net/url"
	"regexp"
	"runtime"
	"strings"
	"time"
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
	time.Sleep(3 * time.Second)
	err := (*page).Locator("#zoomable-container div.surface").Focus()
	if err != nil {
		return err
	}
	t, err := (*page).Content()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	log.Println(t)

	modifier := "Control"
	if php2go.InArray(runtime.GOOS, []string{`windows`, `darwin`}) {
		modifier = "Meta"
	}
	_ = (*page).Keyboard().Press(fmt.Sprintf(`%s+KeyA`, modifier))
	_ = (*page).Keyboard().Press(fmt.Sprintf(`%s+KeyC`, modifier))
	time.Sleep(2 * time.Second)
	r, err := (*page).Evaluate("navigator.clipboard.readText()")
	if err != nil {
		return err
	}
	content, ok := r.(string)
	if !ok {
		return errors.New("could not convert content to string")
	}
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
		return err
	}
	return nil
}

func PageProcessForKDoc(page *playwright.Page) error {
	err := (*page).Locator("#workspace").Focus()
	if err != nil {
		return err
	}
	modifier := "Control"
	if php2go.InArray(runtime.GOOS, []string{`windows`, `darwin`}) {
		modifier = "Meta"
	}
	_ = (*page).Keyboard().Press(fmt.Sprintf(`%s+KeyA`, modifier))
	_ = (*page).Keyboard().Press(fmt.Sprintf(`%s+KeyC`, modifier))

	r, err := (*page).Evaluate("navigator.clipboard.readText()")
	if err != nil {
		return err
	}
	content, ok := r.(string)
	if !ok {
		return errors.New("could not convert content to string")
	}
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", escapeJSContent(content)))
	if err != nil {
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
