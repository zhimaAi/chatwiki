// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetPromptLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	m := msql.Model(`chat_ai_prompt_library_group`, define.Postgres)
	list, err := m.Where(`admin_user_id`, cast.ToString(userId)).Field(`id,group_name,group_desc`).Order(`id`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SavePromptLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	groupName := strings.TrimSpace(c.PostForm(`group_name`))
	groupDesc := strings.TrimSpace(c.PostForm(`group_desc`))
	//check required
	if id < 0 || len(groupName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_prompt_library_group`, define.Postgres)
	if id > 0 {
		groupId, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if cast.ToUint(groupId) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}
	//database dispose
	var err error
	data := msql.Datas{
		`group_name`:  groupName,
		`group_desc`:  groupDesc,
		`update_time`: tool.Time2Int(),
	}
	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		data[`admin_user_id`] = userId
		data[`create_time`] = data[`update_time`]
		id, err = m.Insert(data, `id`)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
}

func DeletePromptLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_prompt_library_group`, define.Postgres)
	groupId, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToUint(groupId) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//dispose relation data
	_, err = msql.Model(`chat_ai_prompt_library_items`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
		Where(`group_id`, groupId).Update(msql.Datas{`group_id`: `0`, `update_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//database dispose
	_, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetPromptLibraryItems(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	groupId := cast.ToInt(c.DefaultQuery(`group_id`, `-1`))
	page := max(1, cast.ToInt(c.DefaultQuery(`page`, `1`)))
	size := max(1, cast.ToInt(c.DefaultQuery(`size`, `10`)))
	m := msql.Model(`chat_ai_prompt_library_items`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId))
	if groupId >= 0 {
		m.Where(`group_id`, cast.ToString(groupId))
	}
	list, err := m.Order(`id DESC`).Limit(size*(page-1), size).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for i, item := range list {
		delete(list[i], `admin_user_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		if cast.ToInt(item[`prompt_type`]) == define.PromptTypeStruct {
			item[`prompt_struct`], _ = common.CheckPromptConfig(define.PromptTypeStruct, item[`prompt_struct`])
			list[i][`markdown`] = common.BuildPromptStruct(define.PromptTypeStruct, ``, item[`prompt_struct`])
		}
	}
	data := map[string]any{`list`: list, `page`: page, `size`: size, `has_more`: len(list) == size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func SavePromptLibraryItems(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	title := strings.TrimSpace(c.PostForm(`title`))
	groupId := cast.ToUint(c.PostForm(`group_id`))
	promptType := cast.ToInt(c.PostForm(`prompt_type`))
	prompt := strings.TrimSpace(c.PostForm(`prompt`))
	promptStruct := strings.TrimSpace(c.PostForm(`prompt_struct`))
	//check required
	if id < 0 || len(title) == 0 || (promptType == define.PromptTypeCustom && len(prompt) == 0) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	var err error
	if promptStruct, err = common.CheckPromptConfig(promptType, promptStruct); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//data check
	m := msql.Model(`chat_ai_prompt_library_items`, define.Postgres)
	if id > 0 {
		item, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Field(`id`).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(item) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}
	if groupId > 0 {
		group, err := msql.Model(`chat_ai_prompt_library_group`, define.Postgres).Where(`id`, cast.ToString(groupId)).
			Where(`admin_user_id`, cast.ToString(userId)).Field(`id`).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(group) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_id`))))
			return
		}
	}
	//database dispose
	data := msql.Datas{
		`title`:       title,
		`group_id`:    groupId,
		`prompt_type`: promptType,
		`update_time`: tool.Time2Int(),
	}
	switch promptType {
	case define.PromptTypeCustom:
		data[`prompt`] = prompt
	case define.PromptTypeStruct:
		data[`prompt_struct`] = promptStruct
	}
	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		data[`admin_user_id`] = userId
		data[`create_time`] = data[`update_time`]
		id, err = m.Insert(data, `id`)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
}

func DeletePromptLibraryItems(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_prompt_library_items`, define.Postgres)
	item, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(item) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	_, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func CreatePromptByLlm(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	demand := strings.TrimSpace(c.Query(`demand`))
	modelConfigId := cast.ToInt(c.Query(`model_config_id`))
	useModel := strings.TrimSpace(c.Query(`use_model`))
	if len(demand) == 0 || modelConfigId <= 0 || len(useModel) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//check model_config_id and use_model
	if ok := common.CheckModelIsValid(adminUserId, modelConfigId, useModel, common.Llm); !ok {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`使用的LLM模型选择错误`)))
		return
	}
	promptStruct, err := common.CreatePromptByAi(demand, adminUserId, modelConfigId, useModel)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`promptStruct`: promptStruct, `markdown`: common.BuildPromptStruct(define.PromptTypeStruct, ``, promptStruct)}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func MovePromptLibraryItems(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	groupId := cast.ToUint(c.PostForm(`group_id`))
	//check required
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_prompt_library_items`, define.Postgres)
	item, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(item) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if groupId > 0 {
		group, err := msql.Model(`chat_ai_prompt_library_group`, define.Postgres).Where(`id`, cast.ToString(groupId)).
			Where(`admin_user_id`, cast.ToString(userId)).Field(`id`).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(group) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_id`))))
			return
		}
	}
	//database dispose
	data := msql.Datas{`group_id`: groupId, `update_time`: tool.Time2Int()}
	_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
}
