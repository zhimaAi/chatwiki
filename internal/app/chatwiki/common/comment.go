// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
)

func GetCommentExecTypeMap(lang string) map[int]string {
	return map[int]string{
		define.CommentExecTypeDelete: i18n.Show(lang, `comment_exec_type_delete`),
		define.CommentExecTypeReply:  i18n.Show(lang, `comment_exec_type_reply`),
		define.CommentExecTypeTop:    i18n.Show(lang, `comment_exec_type_top`),
	}
}
