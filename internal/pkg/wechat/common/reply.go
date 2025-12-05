// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/pkg/lib_define"
	"regexp"
	"strings"
	"time"
)

const (
	PlainText       = 0 // 纯文本
	NormalLink      = 1 // 普通链接
	MiniProgramLink = 2 // 小程序链接
)

// LinkInfo 链接信息结构
type LinkInfo struct {
	Type            int
	Text            string
	URL             string
	AppID           string
	MiniProgramPath string
}

// ContentToALabel 反向解析文本内容，判断是纯文本、普通a标签还是小程序a标签并返回对应参数
func ContentToALabel(content string) *LinkInfo {
	// 去除首尾空格
	content = strings.TrimSpace(content)

	// 如果为空，直接返回纯文本
	if content == "" {
		return &LinkInfo{
			Type: PlainText,
			Text: "",
		}
	}

	// 使用正则表达式匹配普通链接
	normalLinkRegex := regexp.MustCompile(`<a\s+href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	normalMatches := normalLinkRegex.FindStringSubmatch(content)

	// 使用正则表达式匹配小程序链接
	miniprogramLinkRegex := regexp.MustCompile(`<a\s+href=["']#["']\s+data-miniprogram-path=["']([^"']+)["']\s+data-miniprogram-appid=["']([^"']+)["'][^>]*>(.*?)</a>`)
	miniprogramMatches := miniprogramLinkRegex.FindStringSubmatch(content)

	// 判断是否为小程序链接
	if len(miniprogramMatches) == 4 {
		return &LinkInfo{
			Type:            MiniProgramLink,
			Text:            miniprogramMatches[3],
			URL:             "#",
			AppID:           miniprogramMatches[2],
			MiniProgramPath: miniprogramMatches[1],
		}
	}

	// 判断是否为普通链接
	if len(normalMatches) == 3 {
		return &LinkInfo{
			Type: NormalLink,
			Text: normalMatches[2],
			URL:  normalMatches[1],
		}
	}

	// 默认为纯文本
	return &LinkInfo{
		Type: PlainText,
		Text: content,
	}
}

// ProcessEscapeSequences 处理字符串中的转义序列
func ProcessEscapeSequences(content string) string {
	// 将字面意义的 \n 替换为实际的换行符
	content = strings.ReplaceAll(content, `\n`, "\n")
	return content
}

func ReplaceDate(content string) string {
	t := time.Now()

	// 定义日期格式映射
	dateType := map[string]string{
		"{{yyyy-MM-dd hh:mm:ss}}": t.Format("2006-01-02 15:04:05"),
		"{{MM-dd hh:mm:ss}}":      t.Format("01-02 15:04:05"),
		"{{yyyy-MM-dd hh:mm}}":    t.Format("2006-01-02 15:04"),
		"{{MM-dd hh:mm}}":         t.Format("01-02 15:04"),
		"{{yyyy-MM-dd}}":          t.Format("2006-01-02"),
		"{{MM-dd}}":               t.Format("01-02"),
		"{{MM月dd日}}":              t.Format("1月2日"),
		`\n`:                      "\n",
	}

	// 替换所有匹配项
	result := content
	for key, value := range dateType {
		result = strings.ReplaceAll(result, key, value)
	}

	return result
}

func WechatFormatSmartMenu2C(smartMenu lib_define.SmartMenu) string {
	// 处理 MenuDescription 中的转义序列
	description := ProcessEscapeSequences(smartMenu.MenuDescription)
	content := description + "\n"
	for _, menuContent := range smartMenu.MenuContent {
		itemContent := ``
		if menuContent.SerialNo != `` {
			itemContent = menuContent.SerialNo + ` `
		}
		if menuContent.MenuType == lib_define.SmartMenuTypeNormal || menuContent.MenuType == `` {
			itemContent += menuContent.Content
			content += itemContent + "\n"
			continue
		}
		if menuContent.ID != `` {
			itemContent += `<a href="weixin://bizmsgmenu?msgmenucontent=` + menuContent.Content + `&msgmenuid=` + menuContent.ID + `">` + menuContent.Content + `</a>`
		} else {
			if menuContent.MenuType == lib_define.SmartMenuTypeKey {
				itemContent += `<a href="weixin://bizmsgmenu?msgmenucontent=` + menuContent.Content + `&msgmenuid=0">` + menuContent.Content + `</a>`
			}
		}
		content += itemContent + "\n"
	}
	return ReplaceDate(content)
}
