// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetChatVariables(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm("robot_key"))
	robotInfo, err := common.GetRobotInfo(robotKey)
	if err != nil || len(robotInfo) == 0 {
		common.FmtError(c, "no_data")
		return
	}
	m := msql.Model("chat_ai_variables", define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where(`robot_key`, robotKey)

	list, err := m.Order("id asc").Select()
	if err != nil {
		logs.Error("GetChatVariables error: ", err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	common.FmtOk(c, list)
}

func CreateChatVariable(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm("robot_key"))
	robotInfo, err := common.GetRobotInfo(robotKey)
	if err != nil || len(robotInfo) == 0 {
		common.FmtError(c, "no_data")
		return
	}
	variable := common.ChatVariable{}
	variable.VariableType = c.PostForm("variable_type")
	variable.VariableName = c.PostForm("variable_name")
	variable.VariableKey = c.PostForm("variable_key")
	variable.DefaultValue = c.PostForm("default_value")
	variable.MustInput = c.PostForm("must_input")
	variable.MaxInputLength = c.PostForm("max_input_length")
	variable.Options = c.PostForm("options")
	variable.Id = c.PostForm("id")
	variable.RobotKey = robotKey
	variable.AdminUserId = cast.ToString(adminUserId)
	variable.RobotId = robotInfo[`id`]
	if err := variable.VerifyParams(common.GetLang(c)); err != nil {
		common.FmtError(c, err.Error())
		return
	}
	var id int64
	if cast.ToInt(variable.Id) > 0 {
		id = cast.ToInt64(variable.Id)
		variable.UpdateTime = cast.ToString(time.Now().Unix())
		upDataStr, err := tool.JsonEncode(variable)
		if err != nil {
			common.FmtError(c, "param_err", `options`)
			return
		}
		upData := make(map[string]any)
		err = tool.JsonDecode(upDataStr, &upData)
		if err != nil {
			common.FmtError(c, "param_err", `options`)
			return
		}
		delete(upData, `id`)
		delete(upData, `admin_user_id`)
		delete(upData, `create_time`)
		delete(upData, `robot_id`)
		delete(upData, `robot_key`)
		delete(upData, `value`)

		_, err = msql.Model("chat_ai_variables", define.Postgres).
			Where("id", variable.Id).Update(upData)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, "sys_err")
			return
		}
	} else {
		variable.UpdateTime = cast.ToString(time.Now().Unix())
		variable.CreateTime = variable.UpdateTime
		upDataStr, err := tool.JsonEncode(variable)
		if err != nil {
			common.FmtError(c, "param_err", `options`)
			return
		}
		insertData := make(map[string]any)
		err = tool.JsonDecode(upDataStr, &insertData)
		if err != nil {
			common.FmtError(c, "param_err", `options`)
			return
		}
		delete(insertData, `value`)
		delete(insertData, `id`)
		id, err = msql.Model("chat_ai_variables", define.Postgres).
			Insert(insertData, `id`)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, "sys_err")
			return
		}
	}
	common.FmtOk(c, map[string]any{"id": id})
}

func DeleteChatVariable(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	id := strings.TrimSpace(c.PostForm("id"))
	if cast.ToInt(id) <= 0 {
		common.FmtError(c, "param_err", "id")
		return
	}
	existData, err := msql.Model("chat_ai_variables", define.Postgres).
		Where("id", id).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error("DeleteVariable find error: ", err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	if len(existData) == 0 {
		common.FmtError(c, "no_data")
		return
	}
	_, err = msql.Model("chat_ai_variables", define.Postgres).
		Where("id", id).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Delete()
	if err != nil {
		logs.Error("DeleteVariable delete error: ", err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	common.FmtOk(c, nil)
}

func TakeVariablePlaceholders(content string) []string {
	if len(content) == 0 {
		return []string{}
	}
	reg := regexp.MustCompile(`【chat_variable:([A-Za-z_]{1,10})】`)
	matches := reg.FindAllStringSubmatch(content, -1)
	var placeholders []string
	for _, match := range matches {
		if len(match) >= 2 {
			placeholders = append(placeholders, match[1])
		}
	}
	return placeholders
}

func GetChatRobotVariables(dialogueId int, chatBaseParam *define.ChatBaseParam) (data map[string]any) {
	robotInfo := chatBaseParam.Robot
	var fillVariables []map[string]any
	var waitVariables []msql.Params
	data = map[string]any{
		`need_fill_variable`: false,
		`fill_variables`:     fillVariables,
		`wait_variables`:     waitVariables,
		`dialogue_id`:        0,
		`session_id`:         0,
	}
	var err error
	//dialog id
	if dialogueId == 0 {
		dialogueId, err = common.GetDialogueIdNoCreate(chatBaseParam)
		if err != nil {
			logs.Error(err.Error())
			return
		}
	}
	data[`dialogue_id`] = dialogueId
	//session id
	var sessionId = 0
	if dialogueId == 0 {
		sessionId = 0
	} else {
		sessionId, err = common.GetSessionIdNoCreate(dialogueId)
	}
	var sessionInfo = make(msql.Params)
	if sessionId > 0 {
		sessionInfo, err = msql.Model(`chat_ai_session`, define.Postgres).
			Where(`id`, cast.ToString(sessionId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			return
		}
		//fill variables
		if len(sessionInfo) > 0 && len(sessionInfo[`chat_prompt_variables`]) > 0 {
			err = tool.JsonDecode(sessionInfo[`chat_prompt_variables`], &fillVariables)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			data[`fill_variables`] = fillVariables
		}
	}
	data[`session_id`] = sessionId

	//wait variables
	waitVariables, err = msql.Model(`chat_ai_variables`, define.Postgres).
		Where(`robot_key`, robotInfo[`robot_key`]).Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	data[`wait_variables`] = waitVariables

	//need fill variable
	if len(waitVariables) > 0 {
		data[`need_fill_variable`] = true
	}
	if len(fillVariables) > 0 {
		data[`need_fill_variable`] = false
	}
	return
}
