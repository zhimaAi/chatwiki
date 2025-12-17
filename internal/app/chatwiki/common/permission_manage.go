// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/casbin"
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func SavePermissionManage(id, adminUserId int64, data msql.Datas) (int64, error) {
	if len(data) <= 0 {
		return id, errors.New(`data empty`)
	}
	data[`update_time`] = tool.Time2Int()
	m := msql.Model(`permission_manage`, define.Postgres)
	if id > 0 {
		// 更新
		info, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`id`)
		if err != nil {
			logs.Error(err.Error())
			return id, err
		}
		if len(info) == 0 {
			return id, err
		}
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		if err != nil {
			logs.Error(err.Error())
			return id, err
		}
	} else {
		// 新增
		data[`create_time`] = tool.Time2Int()
		id, err := m.Insert(data, `id`)
		if err != nil {
			logs.Error(err.Error())
			return id, err
		}
	}
	return id, nil
}

func GetPermissionManageList(adminUserId int, identityIds string, identityType, objectType int) ([]msql.Params, error) {
	m := msql.Model(`permission_manage`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`identity_type`, cast.ToString(identityType))
	if identityIds != `` {
		m.Where(`identity_id`, `in`, cast.ToString(identityIds))
	}
	if objectType > 0 {
		m.Where(`object_type`, cast.ToString(objectType))
	}
	list, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
	}
	return list, err
}

func GetPartentPermissionList(adminUserId int, identityId string, identityType, objectType int) ([]msql.Params, []msql.Params) {
	var (
		departmentIds       = []string{}
		permissionMergeData = []msql.Params{}
		permissionData      = []msql.Params{}
		err                 error
	)
	if identityType == define.IdentityTypeUser {
		departmentIds, err = GetUserDepartmentIds(cast.ToString(identityId))
		if err != nil {
			logs.Error(err.Error())
			return permissionData, permissionMergeData
		}
	} else {
		departmentInfo, err := GetDepartmentInfo(cast.ToInt(identityId))
		if err != nil {
			logs.Error(err.Error())
			return permissionData, permissionMergeData
		}
		if cast.ToInt(departmentInfo[`pid`]) > 0 {
			departmentIds = append(departmentIds, cast.ToString(departmentInfo[`pid`]))
		}
	}

	// 获取上级部门权限
	if len(departmentIds) > 0 {
		departmentPermissionData, err := GetPermissionManageList(adminUserId, strings.Join(departmentIds, `,`), define.IdentityTypeDepartment, objectType)
		if err != nil {
			logs.Error(err.Error())
			return permissionData, permissionMergeData
		}
		permissionData = append(permissionData, departmentPermissionData...)
	}
	permissionMergeData = mergePermissionData(permissionData)
	return permissionData, permissionMergeData
}

func GetAllPermissionManage(adminUserId int, identityId string, identityType, objectType int) ([]msql.Params, []msql.Params) {
	permissionMergeData := []msql.Params{}
	permissionData, err := GetPermissionManageList(adminUserId, cast.ToString(identityId), identityType, objectType)
	if err != nil {
		logs.Error(err.Error())
		return permissionData, permissionMergeData
	}
	permissionData = GetAdminPermissionData(adminUserId, cast.ToInt(identityId), permissionData)
	departmentIds, err := GetUserAllDepartmentIds(adminUserId, cast.ToString(identityId))
	if err != nil {
		logs.Error(err.Error())
		return permissionData, permissionMergeData
	}
	// 获取上级部门权限
	if len(departmentIds) > 0 {
		departmentPermissionData, err := GetPermissionManageList(adminUserId, strings.Join(departmentIds, `,`), define.IdentityTypeDepartment, objectType)
		if err != nil {
			logs.Error(err.Error())
			return permissionData, permissionMergeData
		}
		permissionData = append(permissionData, departmentPermissionData...)
	}
	permissionMergeData = mergePermissionData(permissionData)
	return permissionData, permissionMergeData
}

func CheckOperateRights(adminUserId, identityId, identityType int, objectId, objectType, operateRights int) bool {
	// role
	info := GetUserInfo(identityId)
	rules, err := casbin.Handler.GetPolicyForUser(info["user_roles"])
	if err != nil {
		logs.Error(err.Error())
	}
	rolePermission := make([]string, 0)
	for _, rule := range rules {
		if len(rule) >= 1 {
			if strings.ContainsAny(rule[1], "/") {
				continue
			}
			rolePermission = append(rolePermission, rule[1])
		}
	}
	if tool.InArrayString(`TeamManage`, rolePermission) {
		return true
	}
	if CheckIsAdminOrCreator(cast.ToInt(identityId), adminUserId, 0) {
		return true
	}
	_, mergePermissionData := GetAllPermissionManage(adminUserId, cast.ToString(identityId), identityType, objectType)
	if identityType == define.IdentityTypeDepartment {
		identityId = 0
	}
	return checkOperateRights(mergePermissionData, adminUserId, identityId, objectId, objectType, operateRights)
}

func CheckObjectAccessRights(adminUserId, identityId, identityType int, objectId, objectType, operateRights int) bool {
	// role
	if CheckIsAdminOrCreator(cast.ToInt(identityId), adminUserId, 0) {
		return true
	}
	_, mergePermissionData := GetAllPermissionManage(adminUserId, cast.ToString(identityId), identityType, objectType)
	if identityType == define.IdentityTypeDepartment {
		identityId = 0
	}
	return checkOperateRights(mergePermissionData, adminUserId, identityId, objectId, objectType, operateRights)
}

func CheckIsAdminOrCreator(user, adminUserId, creator int) bool {
	return user == adminUserId || user == creator
}

func GetAdminPermissionData(adminUserId, userId int, permissionData []msql.Params) []msql.Params {
	if adminUserId != userId {
		return permissionData
	}
	addPermission := map[int]int{
		define.ObjectTypeRobot:   define.ObjectTypeAll,
		define.ObjectTypeLibrary: define.ObjectTypeAll,
		define.ObjectTypeForm:    define.ObjectTypeAll,
	}
	for _, item := range permissionData {
		if cast.ToInt(item[`object_id`]) == define.ObjectTypeAll {
			if _, ok := addPermission[cast.ToInt(item[`object_type`])]; ok {
				delete(addPermission, cast.ToInt(item[`object_type`]))
			}
		}
	}
	if len(addPermission) > 0 {
		for item := range addPermission {
			permissionData = append(permissionData, msql.Params{
				`identity_id`:    cast.ToString(userId),
				`identity_type`:  cast.ToString(define.IdentityTypeUser),
				`object_id`:      cast.ToString(define.ObjectTypeAll),
				`object_type`:    cast.ToString(item),
				`operate_rights`: cast.ToString(define.PermissionManageRights),
			})
		}
	}
	return permissionData
}

func mergePermissionData(permissionData []msql.Params) []msql.Params {
	var (
		permissionMergeData []msql.Params
		permissionMap       = make(map[string]msql.Params)
	)
	if len(permissionData) == 0 {
		return permissionMergeData
	}
	for _, item := range permissionData {
		key := item[`object_id`] + `_` + cast.ToString(item[`object_type`])
		if rights, ok := permissionMap[key]; ok {
			if cast.ToInt(rights[`operate_rights`]) >= cast.ToInt(item[`operate_rights`]) {
				continue
			}
		}
		permissionMap[key] = item
	}
	for _, item := range permissionMap {
		permissionMergeData = append(permissionMergeData, item)
	}
	return permissionMergeData
}

func checkOperateRights(permissionData []msql.Params, adminUserId, userId, objectId, objectType, operateRights int) bool {
	if len(permissionData) == 0 {
		return false
	}
	// admin
	if adminUserId == userId {
		return true
	}
	for _, item := range permissionData {
		// all permission
		if cast.ToInt(item[`object_id`]) == define.ObjectTypeAll && cast.ToInt(item[`object_type`]) == objectType {
			return true
		}
		key := item[`object_id`] + `_` + cast.ToString(item[`object_type`])
		operateKey := cast.ToString(objectId) + `_` + cast.ToString(objectType)
		if key == operateKey {
			return cast.ToInt(item[`operate_rights`]) >= cast.ToInt(operateRights)
		}
	}
	return false
}

func GetUserInfo(user int) msql.Params {
	info, err := msql.Model(define.TableUser, define.Postgres).
		Where(`id`, cast.ToString(user)).
		Find()
	if err != nil {
		logs.Error(err.Error())
	}
	return info
}
