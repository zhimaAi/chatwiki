// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	messageId := cast.ToInt(c.Query(`message_id`))
	fileId := cast.ToInt(c.Query(`file_id`))
	if messageId <= 0 || fileId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	list, err := msql.Model(`chat_ai_answer_source`, define.Postgres).
		Alias(`c`).
		Join(`chat_ai_library_file_data d`, `c.paragraph_id = d.id`, `left`).
		Join(`chat_ai_library_file f`, `d.file_id = f.id`, `left`).
		Where(`c.admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`c.message_id`, cast.ToString(messageId)).
		Where(`c.file_id`, cast.ToString(fileId)).
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
			case when f.id is null or f.file_ext != 'pdf' then 0 else d.page_num end as page_num
		`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var formatedList []map[string]any
	for _, item := range list {
		formatedItem := make(map[string]any)
		for k, v := range item {
			formatedItem[k] = v
		}
		var images []string
		_ = json.Unmarshal([]byte(item[`images`]), &images)
		formatedItem[`images`] = images
		formatedList = append(formatedList, formatedItem)
	}
	c.String(http.StatusOK, lib_web.FmtJson(formatedList, nil))
}
