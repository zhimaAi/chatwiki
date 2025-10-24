// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetRobotGroupList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	applicationType := cast.ToInt(c.DefaultQuery(`application_type`, `-1`))
	m := msql.Model(`chat_ai_robot_group`, define.Postgres)
	wheres := [][]string{{`admin_user_id`, cast.ToString(adminUserId)}}
	list, err := m.Where2(wheres).Field(`id,group_name`).Order(`id`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	list = append([]msql.Params{{`id`: `0`, `group_name`: `未分组`}}, list...)
	if applicationType >= 0 {
		wheres = append(wheres, []string{`application_type`, cast.ToString(applicationType)})
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	//check permission
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot}) {
		managedRobotIdList := []string{`0`}
		permissionData, _ := common.GetAllPermissionManage(adminUserId, cast.ToString(userId), define.IdentityTypeUser, define.ObjectTypeRobot)
		for _, permission := range permissionData {
			managedRobotIdList = append(managedRobotIdList, cast.ToString(permission[`object_id`]))
		}
		//wheres = append(wheres, []string{`id`, `in`, strings.Join(managedRobotIdList, `,`)})
	}
	//统计每个分组下的机器人数量
	stats, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where2(wheres).
		Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	for i, params := range list {
		list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveRobotGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	var err error
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	groupName := strings.TrimSpace(c.PostForm(`group_name`))
	//check required
	if id < 0 || len(groupName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	if utf8.RuneCountInString(groupName) > 50 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_name`))))
		return
	}

	//data check
	m := msql.Model(`chat_ai_robot_group`, define.Postgres)
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
	data := msql.Datas{
		`group_name`:  groupName,
		`update_time`: tool.Time2Int(),
	}
	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		data[`admin_user_id`] = userId
		data[`create_time`] = data[`update_time`]
		id, err = m.Insert(data, `id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
}

func DeleteRobotGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	m := msql.Model(`chat_ai_robot_group`, define.Postgres)
	//检查分组是否存在
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

	//将该分组下的机器人移到未分组
	_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
		Where(`group_id`, cast.ToString(id)).
		Update(msql.Datas{`group_id`: `0`, `update_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	//删除分组
	_, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func RelationRobotGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	robotId := cast.ToInt64(c.PostForm(`robot_id`))
	groupId := cast.ToInt64(c.PostForm(`group_id`))
	//check required
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	//检查机器人是否存在
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToUint(robotInfo) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	//更新机器人分组
	_, err = msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Update(msql.Datas{`group_id`: cast.ToString(groupId), `update_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
