// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := cast.ToInt(c.Query(`library_id`))
	groupType := cast.ToInt(c.DefaultQuery(`group_type`, `0`))
	if libraryId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	m := msql.Model(`chat_ai_library_group`, define.Postgres)
	wheres := [][]string{{`admin_user_id`, cast.ToString(userId)}, {`library_id`, cast.ToString(libraryId)}}
	list, err := m.Where2(wheres).Where(`group_type`, cast.ToString(groupType)).Field(`id,group_name`).Order(`sort desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list = append([]msql.Params{{`id`: `0`, `group_name`: `未分组`}}, list...)

	switch groupType {
	case define.LibraryGroupTypeQA:
		//统计数据
		stats, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where2(wheres).Where(`isolated`, `false`).Where(`delete_time`, `0`).
			Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for i, params := range list {
			list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
		}
	case define.LibraryGroupTypeFile:
		//统计数据
		stats, err := msql.Model(`chat_ai_library_file`, define.Postgres).
			Where2(wheres).Where(`delete_time`, `0`).
			Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for i, params := range list {
			list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	groupType := cast.ToInt(c.DefaultPostForm(`group_type`, `0`))
	groupName := strings.TrimSpace(c.PostForm(`group_name`))
	//check required
	if id < 0 || libraryId <= 0 || len(groupName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if utf8.RuneCountInString(groupName) > 50 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_name`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_library_group`, define.Postgres)
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
		`library_id`:  libraryId,
		`group_name`:  groupName,
		`update_time`: tool.Time2Int(),
	}
	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		data[`admin_user_id`] = userId
		if groupType != define.LibraryGroupTypeQA {
			sort, _ := m.Where(`admin_user_id`, cast.ToString(userId)).Where(`group_type`, cast.ToString(groupType)).Max(`sort`)
			data[`sort`] = cast.ToInt(sort) + 1
		}
		data[`create_time`] = data[`update_time`]
		data[`group_type`] = groupType
		id, err = m.Insert(data, `id`)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(id, nil))
}

func DeleteLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	groupType := cast.ToInt(c.DefaultPostForm(`group_type`, `0`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_library_group`, define.Postgres)
	libraryId, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`library_id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToUint(libraryId) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//dispose relation data
	relationModel := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	if groupType == define.LibraryGroupTypeFile {
		relationModel = msql.Model(`chat_ai_library_file`, define.Postgres)
	}
	_, err = relationModel.Where(`admin_user_id`, cast.ToString(userId)).
		Where(`library_id`, libraryId).Where(`group_id`, cast.ToString(id)).
		Update(msql.Datas{`group_id`: `0`, `update_time`: tool.Time2Int()})
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

func GetLibraryListGroup(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	typ := cast.ToString(c.Query(`type`))
	m := msql.Model(`chat_ai_library_list_group`, define.Postgres)
	wheres := [][]string{{`admin_user_id`, cast.ToString(adminUserId)}}
	list, err := m.Where2(wheres).Field(`id,group_name`).Order(`sort desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	list = append([]msql.Params{{`id`: `0`, `group_name`: `未分组`}}, list...)
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
		permissionData, _ := common.GetAllPermissionManage(adminUserId, cast.ToString(userId), define.IdentityTypeUser, define.ObjectTypeLibrary)
		for _, permission := range permissionData {
			managedRobotIdList = append(managedRobotIdList, cast.ToString(permission[`object_id`]))
		}
		//wheres = append(wheres, []string{`id`, `in`, strings.Join(managedRobotIdList, `,`)})
	}
	if typ == "" {
		typ = fmt.Sprintf(`%v,%v`, define.GeneralLibraryType, define.QALibraryType)
	} else {
		typ = fmt.Sprintf(`%v`, cast.ToInt(typ))
	}
	//统计数据
	stats, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where2(wheres).
		Where(`type`, `in`, typ).
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

func SaveLibraryListGroup(c *gin.Context) {
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
	m := msql.Model(`chat_ai_library_list_group`, define.Postgres)
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
		sort, _ := m.Where(`admin_user_id`, cast.ToString(userId)).Max(`sort`)
		data[`admin_user_id`] = userId
		data[`sort`] = cast.ToInt(sort) + 1
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

func DeleteLibraryListGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	m := msql.Model(`chat_ai_library_list_group`, define.Postgres)
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
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
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

func SortLibararyListGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	sortGroup := strings.TrimSpace(c.PostForm(`sort_group`))
	sortSense := cast.ToInt(c.DefaultPostForm(`sense`, `0`))
	if len(sortGroup) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	groupList := []msql.Datas{}
	tool.JsonDecode(sortGroup, &groupList)
	m := msql.Model(`chat_ai_library_list_group`, define.Postgres)
	if sortSense == define.LibraryGroupTypeFile {
		m = msql.Model(`chat_ai_library_group`, define.Postgres)
	}
	for _, item := range groupList {
		if cast.ToInt(item[`id`]) <= 0 {
			continue
		}
		_, err := m.Where(`admin_user_id`, cast.ToString(userId)).
			Where(`id`, cast.ToString(item[`id`])).
			Update(msql.Datas{`sort`: cast.ToInt(item[`sort`]), `update_time`: tool.Time2Int()})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func RelationLibraryGroup(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	libraryId := cast.ToInt64(c.PostForm(`library_id`))
	fileId := cast.ToInt64(c.PostForm(`file_id`))
	groupId := cast.ToInt64(c.PostForm(`group_id`))
	sense := cast.ToInt(c.DefaultPostForm(`sense`, `0`))
	//check required
	if libraryId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	//检查机器人是否存在
	libraryInfo, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`id`, cast.ToString(libraryId)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToUint(libraryInfo) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	//更新机器人分组
	if sense == define.LibraryGroupTypeFile {
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, cast.ToString(fileId)).
			Where(`admin_user_id`, cast.ToString(userId)).
			Update(msql.Datas{`group_id`: cast.ToString(groupId), `update_time`: tool.Time2Int()})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		_, err = msql.Model(`chat_ai_library`, define.Postgres).
			Where(`id`, cast.ToString(libraryId)).
			Where(`admin_user_id`, cast.ToString(userId)).
			Update(msql.Datas{`group_id`: cast.ToString(groupId), `update_time`: tool.Time2Int()})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
