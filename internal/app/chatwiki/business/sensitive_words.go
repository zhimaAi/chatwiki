// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
)

func CheckSensitiveWords(c *gin.Context) {
	params := getChatRequestParam(c)
	if params.Error == nil && cast.ToInt(params.Robot[`sensitive_words_switch`]) == define.SwitchOn {
		if ok, wordsArr := common.CheckSensitiveWords(params.Question, cast.ToInt(params.Robot[`admin_user_id`]), cast.ToInt(params.Robot[`id`])); ok {
			common.FmtOk(c, map[string]any{`code`: define.ErrorCodeContainsSensitiveWords, `words`: wordsArr})
			return
		}
	}
	common.FmtOk(c, nil)
}
