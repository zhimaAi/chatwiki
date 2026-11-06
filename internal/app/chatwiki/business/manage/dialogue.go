// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func GetDialogueList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//format check
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	//data check
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//get params
	minId := cast.ToUint(c.PostForm(`min_id`))
	size := max(1, cast.ToInt(c.PostForm(`size`)))
	m := msql.Model(`chat_ai_dialogue`, define.Postgres).Where(`is_background`, `1`).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, robot[`id`])
	if minId > 0 {
		m.Where(`id`, `<`, cast.ToString(minId))
	}
	list, err := m.Limit(size).Order(`id desc`).Field(`id,openid,subject,create_time`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetAnswerSource(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	messageId := cast.ToInt(c.Query(`message_id`))
	fileId := cast.ToInt(c.Query(`file_id`))
	if messageId <= 0 || fileId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	access, ok := getAnswerSourceAccess(c, messageId)
	if !ok {
		return
	}
	list, err := answerSourceQuery(access.AdminUserID, messageId, fileId).
		Order(`c.id`).
		Field(`
			c.paragraph_id as id,
			c.word_total,
			c.similarity,
			c.title,
			c.type,
			c.content,
			c.question,
			c.answer,
			c.images,
			case when f.file_ext != 'pdf' then 0 else d.page_num end as page_num
		`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(list) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	formatedList := make([]map[string]any, 0, len(list))
	for _, item := range list {
		formatedItem := make(map[string]any)
		for k, v := range item {
			formatedItem[k] = v
		}
		var images []string
		_ = json.Unmarshal([]byte(item[`images`]), &images)
		formatedItem[`images`] = images
		if page := cast.ToInt(item[`page_num`]); page > 0 {
			previewToken, signErr := common.SignLibFilePagePreview(common.LibFilePagePreview{
				AdminUserID: access.AdminUserID, RobotID: access.RobotID, OpenID: access.OpenID,
				MessageID: messageId, FileID: fileId, ParagraphID: cast.ToInt(item[`id`]),
				Page: page, AdminAccess: access.AdminAccess,
			}, time.Now())
			if signErr != nil {
				logs.Error(signErr.Error())
				common.FmtError(c, `sys_err`)
				return
			}
			formatedItem[`preview_url`] = "/chat/getLibRawFileOnePage?preview_token=" + url.QueryEscape(previewToken)
		}
		formatedList = append(formatedList, formatedItem)
	}
	c.String(http.StatusOK, lib_web.FmtJson(formatedList, nil))
}
