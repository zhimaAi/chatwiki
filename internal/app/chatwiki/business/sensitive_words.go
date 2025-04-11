package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
