// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
		question := common.GetFirstQuestionByInput(params.Question) //多模态输入特殊处理
		if cast.ToInt(params.Robot[`sensitive_words_switch`]) == define.SwitchOn && len(question) > 0 {
			if ok, wordsArr := common.CheckSensitiveWords(question, cast.ToInt(params.Robot[`admin_user_id`]), cast.ToInt(params.Robot[`id`])); ok {
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
			if cast.ToInt(externalConfigH5[`accessRestrictionsType`]) == define.AccessPermissionTypeLogin {
				var (
					adminUserId = common.GetAdminUserId(c)
					userId      = common.GetLoginUserId(c)
				)
				if adminUserId == 0 {
					common.FmtOk(c, map[string]any{`code`: define.ErrorCodeNeedLogin})
					return
				}
				// check permission
				permission := common.CheckObjectAccessRights(adminUserId, userId, define.IdentityTypeUser, cast.ToInt(params.Robot[`id`]), define.ObjectTypeRobot, define.PermissionQueryRights)
				if !permission {
					common.FmtOk(c, map[string]any{`code`: define.ErrorCodeNeedNoPermissionLogin})
					return
				}
			}
		}
	}
	common.FmtOk(c, nil)
}
