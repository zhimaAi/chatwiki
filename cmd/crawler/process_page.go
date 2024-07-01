// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	netURL "net/url"
	"regexp"
	"runtime"
	"strings"
)

type PageInfo struct {
	RawHtml    string `json:"html"`
	MainHtml   string `json:"main_html"`
	Screenshot string `json:"screenshot"`
}

type NeedWaitPageFunc func(page *playwright.Page) error
type SpecialPageProcessingFunc func(page *playwright.Page) error

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
	screen := playwright.Size{Width: 1980, Height: 1080}
	viewport := playwright.Size{Width: 1920, Height: 1080}
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Permissions:     permissions,
		AcceptDownloads: &acceptDownloads,
		Screen:          &screen,
		Viewport:        &viewport,
		UserAgent:       &userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create context: %v", err)
	}
	// create a new page
	page, err := context.NewPage()
	//page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	defer func(page *playwright.Page) {
		_ = (*page).Close()
	}(&page)

	// abort requests
	_ = page.Route("**/*", func(route playwright.Route) {
		if php2go.InArray(route.Request().ResourceType(), []string{"image", "media", "font", "websocket", "texttrack", "eventsource"}) {
			_ = route.Abort()
		} else {
			_ = route.Continue()
		}
	})

	// navigate to url, with timeout 15s
	timeout := float64(15000)
	if _, err = page.Goto(parsedURL.String(), playwright.PageGotoOptions{Timeout: &timeout}); err != nil {
		return nil, fmt.Errorf("could not navigate to url: %v", err)
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
	"www.kdocs.cn":    PageProcessForKDoc,
	"feishu.cn":       PageProcessForFeishu,
	"www.jianshu.com": PageProcessForJianshu,
	"yuque.com":       PageProcessForYuque,
}

func PageProcessForTencentDoc(page *playwright.Page) error {
	err := (*page).WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle})
	if err != nil {
		return err
	}

	err = (*page).Locator("#scrollable-content > div.resize-sensor > div.resize-sensor-expand").Focus()
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
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", strings.ReplaceAll(content, "`", "\\`")))
	if err != nil {
		return err
	}
	return nil
}

func PageProcessForKDoc(page *playwright.Page) error {
	err := (*page).Locator("#mainpages").Focus()
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
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", strings.ReplaceAll(content, "`", "\\`")))
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
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", strings.ReplaceAll(content, "`", "\\`")))
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
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", strings.ReplaceAll(content, "`", "\\`")))
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
	_, err = (*page).Evaluate(fmt.Sprintf("document.body.innerHTML = `%s`", strings.ReplaceAll(content, "`", "\\`")))
	if err != nil {
		return err
	}
	return nil
}
