// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// GetPermissionManageList
func GetPermissionManageList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// Get identity type params
	identityId := cast.ToString(c.Query("identity_id"))
	identityType := cast.ToInt(c.Query("identity_type"))
	objectType := cast.ToInt(c.Query("object_type"))
	tab := cast.ToInt(c.Query("tab")) // 1: member department 2: parent department
	if identityType <= 0 || objectType <= 0 || objectType > 3 {
		common.FmtError(c, "param_lack")
		return
	}
	var list = make([]msql.Params, 0)
	if tool.InArrayInt(tab, []int{1, 2}) {
		list, _ = common.GetPartentPermissionList(adminUserId, identityId, identityType, objectType)
	} else {
		list, _ = common.GetPermissionManageList(adminUserId, identityId, identityType, objectType)
	}
	if len(list) == 0 {
		common.FmtOk(c, list)
		return
	}
	objectIds := []string{}
	for _, item := range list {
		if cast.ToInt(item[`object_id`]) == define.ObjectTypeAll {
			continue
		}
		objectIds = append(objectIds, item[`object_id`])
	}

	// relation data
	relationData, err := getRelationData(adminUserId, objectType, strings.Join(objectIds, `,`))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	for _, item := range list {
		item[`name`] = ``
		item[`avatar`] = ``
		if data, ok := relationData[item[`object_id`]]; ok {
			item[`name`] = data[`name`]
			item[`avatar`] = data[`avatar`]
		}
	}
	common.FmtOk(c, list)
}

// GetPartnerManageList Collaborators
func GetPartnerManageList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	// Get params
	objectId := cast.ToInt(c.Query("object_id"))
	objectType := cast.ToInt(c.Query("object_type"))
	if objectId <= 0 || objectType <= 0 || objectType > 3 {
		common.FmtError(c, "param_lack")
		return
	}

	// Build query conditions
	m := msql.Model("permission_manage", define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where(`object_id`, `in`, fmt.Sprintf(`%v,%v`, objectId, define.ObjectTypeAll)).
		Where("object_type", cast.ToString(objectType))
	// Get permission list
	list, err := m.Order("id desc").Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	if len(list) == 0 {
		common.FmtOk(c, list)
		return
	}
	indentityUser, identityDepartment := []string{}, []string{}
	for _, item := range list {
		if cast.ToInt(item[`identity_type`]) == define.IdentityTypeUser {
			indentityUser = append(indentityUser, item[`identity_id`])
		} else if cast.ToInt(item[`identity_type`]) == define.IdentityTypeDepartment {
			identityDepartment = append(identityDepartment, item[`identity_id`])
		}
	}
	// relation data
	relationUserData, err := getIdentityRelationData(adminUserId, define.IdentityTypeUser, strings.Join(indentityUser, `,`))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	relationDepartmentData, err := getIdentityRelationData(adminUserId, define.IdentityTypeDepartment, strings.Join(identityDepartment, `,`))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	var (
		topData   = make([]msql.Params, 1)
		otherList = make([]msql.Params, 0)
		creator   int
	)
	relationData, _ := getRelationData(adminUserId, objectType, cast.ToString(objectId))
	if data, ok := relationData[cast.ToString(objectId)]; ok {
		creator = cast.ToInt(data[`creator`])
	}
	for _, item := range list {
		item[`name`] = ``
		item[`user_name`] = ``
		item[`avatar`] = ``
		item[`role_type`] = ``
		item[`is_creator`] = `0`
		if cast.ToInt(item[`identity_type`]) == define.IdentityTypeUser {
			if data, ok := relationUserData[item[`identity_id`]]; ok {
				item[`name`] = data[`name`]
				item[`avatar`] = data[`avatar`]
				item[`user_name`] = data[`user_name`]
				item[`role_type`] = data[`role_type`]
				if common.CheckIsAdminOrCreator(cast.ToInt(item[`identity_id`]), adminUserId, 0) {
					topData[0] = item
				} else if common.CheckIsAdminOrCreator(cast.ToInt(item[`identity_id`]), 0, creator) {
					item[`is_creator`] = `1`
					if len(topData) == 1 {
						topData = append(topData, item)
					}
				} else {
					otherList = append(otherList, item)
				}
			}
		} else if data, ok := relationDepartmentData[item[`identity_id`]]; ok {
			item[`name`] = data[`name`]
			item[`avatar`] = data[`avatar`]
			otherList = append(otherList, item)
		}
	}
	common.FmtOk(c, append(topData, otherList...))
}

func getRelationData(adminUserId, objectType int, objectIds string) (map[string]msql.Params, error) {
	if len(objectIds) == 0 {
		return nil, nil
	}
	var (
		relationData []msql.Params
		err          error
	)
	if objectType == define.ObjectTypeRobot {
		relationData, err = msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, objectIds).
			Field(`id,robot_name as name,robot_avatar as avatar,creator`).
			Select()
	} else if objectType == define.ObjectTypeLibrary {
		relationData, err = msql.Model(`chat_ai_library`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, objectIds).
			Where(`type`, `in`, fmt.Sprintf(`%v,%v`, cast.ToString(define.GeneralLibraryType), cast.ToString(define.QALibraryType))).
			Field(`id,library_name as name,creator`).
			Select()
	} else if objectType == define.ObjectTypeForm {
		relationData, err = msql.Model(`form`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, `in`, objectIds).
			Where(`delete_time`, `0`).
			Field(`id,name,creator`).
			Select()
	}
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return formatRelationData(relationData), nil
}

func getIdentityRelationData(adminUserId, identityType int, identityIds string) (map[string]msql.Params, error) {
	if len(identityIds) == 0 {
		return nil, nil
	}
	var (
		relationData []msql.Params
		err          error
	)
	if identityType == define.IdentityTypeUser {
		relationData, err = msql.Model(define.TableUser, define.Postgres).
			Where(`id`, `in`, identityIds).
			Field(`id,user_roles,user_name,nick_name as name,avatar`).
			Select()
		roles, _ := msql.Model(define.TableRole, define.Postgres).Select()
		roleMap := make(map[string]msql.Params)
		for _, role := range roles {
			roleMap[role["id"]] = role
		}
		for _, user := range relationData {
			user["role_name"] = ""
			user["role_type"] = ""
			if role, ok := roleMap[user["user_roles"]]; ok {
				user["role_name"] = role[`name`]
				user["role_type"] = role[`role_type`]
			}
		}
	} else if identityType == define.IdentityTypeDepartment {
		relationData, err = msql.Model(`department`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, identityIds).
			Field(`id,department_name as name`).
			Select()
	}
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return formatRelationData(relationData), nil
}

func formatRelationData(list []msql.Params) map[string]msql.Params {
	var data = make(map[string]msql.Params)
	for _, item := range list {
		data[item[`id`]] = item
	}
	return data
}

// BatchSavePermissionManage Save permission management config
func BatchSavePermissionManage(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	type PermissionManage struct {
		ObjectId      int `json:"object_id"`
		ObjectType    int `json:"object_type"`
		OperateRights int `json:"operate_rights"`
	}
	var objectArray []PermissionManage
	// Get params
	identityIds := strings.TrimSpace(c.PostForm("identity_ids"))
	identityType := cast.ToInt(c.PostForm("identity_type"))
	objectType := cast.ToInt(c.PostForm("object_type"))
	err := tool.JsonDecode(c.PostForm("object_array"), &objectArray)
	// Validate params
	if len(identityIds) == 0 || identityType <= 0 || err != nil {
		common.FmtError(c, "param_lack")
		return
	}
	var (
		objectIds = []string{}
	)
	for _, item := range objectArray {
		objectIds = append(objectIds, cast.ToString(item.ObjectId))
	}
	m := msql.Model("permission_manage", define.Postgres)
	if cast.ToInt(identityType) == define.IdentityTypeUser {
		// check is creator
		allPermissions, err := m.Where("identity_id", `in`, identityIds).
			Where("identity_type", cast.ToString(identityType)).
			Where("object_type", cast.ToString(objectType)).ColumnArr(`object_id`)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		var exceptObjectIds = []string{}
		for _, item := range allPermissions {
			if cast.ToInt(item) > 0 && (!tool.InArrayString(item, objectIds) || len(objectIds) == 0) {
				exceptObjectIds = append(exceptObjectIds, item)
			}
		}
		getRelationUserData, _ := getIdentityRelationData(adminUserId, identityType, identityIds)
		if len(exceptObjectIds) > 0 {
			relationData, _ := getRelationData(adminUserId, objectType, strings.Join(exceptObjectIds, ","))
			for _, item := range relationData {
				if tool.InArrayString(item[`creator`], strings.Split(identityIds, `,`)) {
					userData := getRelationUserData[item[`creator`]]
					common.FmtError(c, `permission_creator_no_operate`, item[`name`], userData[`name`])
					return
				}
			}
		}
		if len(objectIds) > 0 {
			relationData, _ := getRelationData(adminUserId, objectType, strings.Join(objectIds, ","))
			for _, item := range objectArray {
				creatorData, ok := relationData[cast.ToString(item.ObjectId)]
				if !ok {
					continue
				}
				if tool.InArrayString(creatorData[`creator`], strings.Split(identityIds, `,`)) && item.OperateRights < define.PermissionManageRights {
					userData := getRelationUserData[creatorData[`creator`]]
					common.FmtError(c, `permission_creator_no_operate`, creatorData[`name`], userData[`name`])
					return
				}
			}
		}
	}
	// Begin transaction
	m.Begin()
	_, err = m.Where("identity_id", `in`, identityIds).
		Where("identity_type", cast.ToString(identityType)).
		Where("object_type", cast.ToString(objectType)).Delete()
	if err != nil {
		logs.Error(err.Error())
		m.Rollback()
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	for _, user := range strings.Split(identityIds, ",") {
		if cast.ToInt(user) == 0 {
			m.Rollback()
			common.FmtError(c, "param_lack")
			return
		}
		for _, item := range objectArray {
			if item.ObjectId <= 0 || item.ObjectType <= 0 {
				m.Rollback()
				common.FmtError(c, "param_lack")
				return
			}
			// Validate operation rights value
			validRights := []int{define.PermissionQueryRights, define.PermissionManageRights, define.PermissionEditRights}
			if !tool.InArray(item.OperateRights, validRights) {
				m.Rollback()
				common.FmtError(c, "param_invalid", "operate_rights")
				return
			}
			// check permission
			if permission := common.CheckOperateRights(adminUserId, userId, define.IdentityTypeUser, item.ObjectId, item.ObjectType, define.PermissionManageRights); !permission {
				m.Rollback()
				common.FmtError(c, "auth_no_permission")
				return
			}
			// Insert new permission
			_, err = common.SavePermissionManage(0, cast.ToInt64(userId), msql.Datas{
				"admin_user_id":  adminUserId,
				"identity_id":    user,
				"object_id":      item.ObjectId,
				"operate_rights": item.OperateRights,
				"creator":        getLoginUserId(c),
				"identity_type":  identityType,
				"object_type":    item.ObjectType,
			})
			if err != nil {
				m.Rollback()
				logs.Error(err.Error())
				common.FmtError(c, "sys_err", err.Error())
				return
			}
		}
	}
	m.Commit()
	common.FmtOk(c, nil)
}

// SavePermissionManage Save permission management config
func SavePermissionManage(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	type PermissionManage struct {
		ObjectId      int `json:"object_id"`
		ObjectType    int `json:"object_type"`
		OperateRights int `json:"operate_rights"`
	}
	var objectArray []PermissionManage
	// Get params
	identityIds := strings.TrimSpace(c.PostForm("identity_ids"))
	identityType := cast.ToInt(c.PostForm("identity_type"))
	err := tool.JsonDecode(c.PostForm("object_array"), &objectArray)

	// Validate params
	if len(identityIds) == 0 || identityType <= 0 || len(objectArray) <= 0 || err != nil {
		common.FmtError(c, "param_lack")
		return
	}
	// Begin transaction
	m := msql.Model("permission_manage", define.Postgres)
	m.Begin()
	for _, user := range strings.Split(identityIds, ",") {
		if cast.ToInt(user) == 0 {
			m.Rollback()
			common.FmtError(c, "param_lack")
			return
		}
		for _, item := range objectArray {
			if item.ObjectId <= 0 || item.ObjectType <= 0 {
				m.Rollback()
				common.FmtError(c, "param_lack")
				return
			}
			// Validate operation rights value
			validRights := []int{define.PermissionQueryRights, define.PermissionManageRights, define.PermissionEditRights}
			if !tool.InArray(item.OperateRights, validRights) {
				m.Rollback()
				common.FmtError(c, "param_invalid", "operate_rights")
				return
			}
			// check permission
			permission := common.CheckOperateRights(adminUserId, userId, define.IdentityTypeUser, item.ObjectId, item.ObjectType, define.PermissionManageRights)
			if !permission {
				m.Rollback()
				common.FmtError(c, "auth_no_permission")
				return
			}
			// Check whether permission record exists
			existingPerms, err := m.Where("admin_user_id", cast.ToString(adminUserId)).
				Where("identity_id", user).
				Where("identity_type", cast.ToString(identityType)).
				Where("object_id", cast.ToString(item.ObjectId)).
				Where("object_type", cast.ToString(item.ObjectType)).
				Value(`id`)
			if err != nil {
				m.Rollback()
				logs.Error(err.Error())
				common.FmtError(c, "sys_err", err.Error())
				return
			}
			if len(existingPerms) > 0 {
				// Update existing permission
				_, err = common.SavePermissionManage(cast.ToInt64(existingPerms), cast.ToInt64(adminUserId), msql.Datas{
					"operate_rights": item.OperateRights,
					"update_time":    tool.Time2Int(),
				})
			} else {
				// Insert new permission
				_, err = common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
					"admin_user_id":  adminUserId,
					"identity_id":    user,
					"object_id":      item.ObjectId,
					"operate_rights": item.OperateRights,
					"creator":        getLoginUserId(c),
					"identity_type":  identityType,
					"object_type":    item.ObjectType,
				})
			}
			if err != nil {
				m.Rollback()
				logs.Error(err.Error())
				common.FmtError(c, "sys_err", err.Error())
				return
			}
		}
	}
	m.Commit()
	common.FmtOk(c, nil)
}

func DeletePermissionMange(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	identityId := cast.ToInt(c.PostForm("identity_id"))
	identityType := cast.ToInt(c.PostForm("identity_type"))
	objectId := cast.ToInt(c.PostForm("object_id"))
	objectType := cast.ToInt(c.PostForm("object_type"))
	// Validate params
	if identityId == 0 || identityType <= 0 || objectId == 0 || objectType == 0 {
		common.FmtError(c, "param_lack")
		return
	}
	// delete
	m := msql.Model("permission_manage", define.Postgres)
	if _, err := m.Where("admin_user_id", cast.ToString(adminUserId)).
		Where("identity_id", cast.ToString(identityId)).
		Where("identity_type", cast.ToString(identityType)).
		Where("object_id", cast.ToString(objectId)).
		Where("object_type", cast.ToString(objectType)).
		Delete(); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	common.FmtOk(c, nil)
}

func AddDefaultPermissionManage(adminUserId, userId, objectId, objectType int) {
	if adminUserId == userId {
		return
	}
	common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
		`admin_user_id`:  adminUserId,
		`identity_id`:    userId,
		`object_id`:      objectId,
		`operate_rights`: define.PermissionManageRights,
		`creator`:        adminUserId,
		`identity_type`:  define.IdentityTypeUser,
		`object_type`:    objectType,
	})
}
