// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package i18n

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
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
