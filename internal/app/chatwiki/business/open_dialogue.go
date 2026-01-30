// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func OpenGetDialogueList(c *gin.Context) {
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
	openid := c.PostForm(`openid`)
	if len(openid) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))))
		return
	}
	m := msql.Model(`chat_ai_dialogue`, define.Postgres).Where(`is_background`, `0`).
		Where(`admin_user_id`, robot[`admin_user_id`]).
		Where(`openid`, openid).Where(`robot_id`, robot[`id`])
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

func OpenDeleteDialogue(c *gin.Context) {
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
	ids := c.PostForm(`ids`)
	openid := c.PostForm(`openid`)
	if len(openid) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))))
		return
	}
	if ids == `-1` {
		_, err := msql.Model(`chat_ai_dialogue`, define.Postgres).
			Where(`is_background`, `0`).
			Where(`admin_user_id`, robot[`admin_user_id`]).
			Where(`openid`, openid).Where(`robot_id`, robot[`id`]).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		_, err := msql.Model(`chat_ai_dialogue`, define.Postgres).
			Where(`is_background`, `0`).
			Where(`admin_user_id`, robot[`admin_user_id`]).
			Where(`openid`, openid).Where(`robot_id`, robot[`id`]).
			Where(`id`, `in`, ids).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func OpenEditDialogue(c *gin.Context) {
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
	id := c.PostForm(`id`)
	openid := c.PostForm(`openid`)
	subject := c.PostForm(`subject`)
	if len(subject) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `subject`))))
		return
	}
	if len(openid) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))))
		return
	}
	_, err = msql.Model(`chat_ai_dialogue`, define.Postgres).
		Where(`is_background`, `0`).
		Where(`admin_user_id`, robot[`admin_user_id`]).
		Where(`openid`, openid).Where(`id`, id).Where(`robot_id`, robot[`id`]).Update(msql.Datas{
		`subject`:     subject,
		`update_time`: time.Now().Unix(),
	})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
