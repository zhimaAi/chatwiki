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

// PlaceholderRegex regular expression to match placeholders in [[ZM--variable_name--ZM]] format
var PlaceholderRegex = regexp.MustCompile(PlaceholderRegexStr)

// ReplacePlaceholdersInString replace placeholders in string
func ReplacePlaceholdersInString(text string, lang string) string {
	// find all matching placeholders
	matches := PlaceholderRegex.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return text
	}

	// replace from back to front to prevent index offset
	for i := len(matches) - 1; i >= 0; i-- {
		start, end := matches[i][0], matches[i][1]
		fullMatch := text[start:end]                                                                            // full placeholder, e.g. [[ZM--variable_name--ZM]]
		variableName := strings.TrimPrefix(strings.TrimSuffix(fullMatch, PlaceholderSuffix), PlaceholderPrefix) // extract variable name (remove [[ZM-- and --ZM]])

		// trim spaces around variable name
		variableName = strings.TrimSpace(variableName)
		// get translated text from i18n
		translatedText := Show(lang, variableName)

		// ensure translated text is JSON escaped to avoid breaking JSON format
		escapedText, err := json.Marshal(translatedText)
		if err != nil {
			// if escape fails, use raw text as fallback
			logs.Info("JSON escape failed for translation: %s, using raw text", translatedText)
		} else {
			// remove quotes at both ends (json.Marshal adds quotes)
			escapedText = escapedText[1 : len(escapedText)-1]
			translatedText = string(escapedText)
		}

		// replace placeholder
		text = text[:start] + translatedText + text[end:]
	}

	return text
}
