// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package i18n

import (
	"chatwiki/internal/app/chatwiki/define"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/beego/i18n"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func init() {
	for _, lang := range define.Langs {
		filePath := fmt.Sprintf(`%si18n/locale_%s.ini`, define.AppRoot, lang)
		if err := i18n.SetMessage(lang, filePath); err != nil {
			logs.Error(err.Error())
		}
	}
}

func Show(lang, message string, args ...any) string {
	if !tool.InArrayString(lang, define.Langs) {
		lang = define.Langs[0]
	}
	content := i18n.Tr(lang, message, args...)
	//replace special symbol
	content = strings.ReplaceAll(content, `\r`, "\r")
	content = strings.ReplaceAll(content, `\n`, "\n")
	content = strings.ReplaceAll(content, `\f`, "\f")
	return content
}

const PlaceholderRegexStr = `\[\[ZM--\s*(\w+)\s*--ZM\]\]`
const PlaceholderPrefix = "[[ZM--"
const PlaceholderSuffix = "--ZM]]"

// PlaceholderRegex 正则表达式用于匹配 [[ZM--变量名--ZM]] 格式的占位符
var PlaceholderRegex = regexp.MustCompile(PlaceholderRegexStr)

// ReplacePlaceholdersInString 替换字符串中的占位符
func ReplacePlaceholdersInString(text string, lang string) string {
	// 查找所有匹配的占位符
	matches := PlaceholderRegex.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return text
	}

	// 从后往前替换，防止索引偏移
	for i := len(matches) - 1; i >= 0; i-- {
		start, end := matches[i][0], matches[i][1]
		fullMatch := text[start:end]                                                                            // 完整的占位符，例如 [[ZM--variable_name--ZM]]
		variableName := strings.TrimPrefix(strings.TrimSuffix(fullMatch, PlaceholderSuffix), PlaceholderPrefix) // 提取变量名部分 (去掉 [[ZM-- 和 --ZM]])

		// 去除变量名前后的空格
		variableName = strings.TrimSpace(variableName)
		// 从i18n获取对应的翻译文本
		translatedText := Show(lang, variableName)

		// 确保翻译文本经过JSON转义，避免破坏JSON格式
		escapedText, err := json.Marshal(translatedText)
		if err != nil {
			// 如果转义失败，使用原文本作为备选方案
			logs.Info("JSON escape failed for translation: %s, using raw text", translatedText)
		} else {
			// 移除首尾的引号（json.Marshal会添加引号）
			escapedText = escapedText[1 : len(escapedText)-1]
			translatedText = string(escapedText)
		}

		// 替换占位符
		text = text[:start] + translatedText + text[end:]
	}

	return text
}
