// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"unicode"

	"github.com/gin-gonic/gin"
)

func GetLang(c *gin.Context) string {
	lang := c.GetHeader(`lang`)
	if len(lang) == 0 {
		lang = define.LangZhCn
	}
	return lang
}

func IsContainChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts[`Han`], r) {
			return true
		}
	}
	return false
}
