// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"strings"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func CheckChatRequestPermission(c *gin.Context) {
	params := getChatRequestParam(c)
	var externalConfigH5 map[string]any
	from := strings.TrimSpace(c.PostForm(`from`))
	if params.Error == nil {
		if cast.ToInt(params.Robot[`sensitive_words_switch`]) == define.SwitchOn && params.Question != "" {
			if ok, wordsArr := common.CheckSensitiveWords(params.Question, cast.ToInt(params.Robot[`admin_user_id`]), cast.ToInt(params.Robot[`id`])); ok {
				common.FmtOk(c, map[string]any{`code`: define.ErrorCodeContainsSensitiveWords, `words`: wordsArr})
				return
			}
		}
		if params.Robot[`external_config_h5`] != "" && from != lib_define.AppYunPc {
			err := tool.JsonDecode(params.Robot[`external_config_h5`], &externalConfigH5)
			if err != nil {
				logs.Error(err.Error())
				common.FmtError(c, `sys_err`)
				return
			}
			if cast.ToInt(externalConfigH5[`accessRestrictionsType`]) == define.AccessRestrictionsTypeLogin && common.GetAdminUserId(c) == 0 {
				common.FmtOk(c, map[string]any{`code`: define.ErrorCodeNeedLogin})
				return
			}
		}
	}
	common.FmtOk(c, nil)
}
