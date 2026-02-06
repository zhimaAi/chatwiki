// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// GetDepartmentList gets the department list
func GetDepartmentList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	// get tree
	departmentData, _, err := common.GetDepartmentTrees(userId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	for _, v := range departmentData {
		if v.IsDefault == define.SwitchOn {
			count, _ := msql.Model(define.TableUser, define.Postgres).Where("parent_id", cast.ToString(userId)).WhereOr("id", cast.ToString(userId)).Where("is_deleted", define.Normal).Count(`id`)
			v.ChildrenNums = count
			break
		}
	}
	common.FmtOk(c, departmentData)
}

// GetAllDepartment
func GetAllDepartment(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	//Build query conditions
	m := msql.Model(`department`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))
	//Get all department data
	list, err := m.Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	common.FmtOk(c, list)
}

// SaveDepartment saves department info
func SaveDepartment(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	var err error
	//Get params
	id := cast.ToInt64(c.PostForm(`id`))
	pid := cast.ToInt64(c.PostForm(`pid`))
	departmentName := strings.TrimSpace(c.PostForm(`department_name`))
	if len(departmentName) == 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	// check level is max
	if ok, level := common.OverDepartmentLevel(userId, cast.ToInt(pid), cast.ToInt(id)); ok {
		common.FmtError(c, `department_level_max`, cast.ToString(level))
		return
	}
	id, err = common.SaveDepartment(id, cast.ToInt64(userId), msql.Datas{
		`admin_user_id`:   userId,
		`department_name`: departmentName,
		`pid`:             pid,
	})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	common.FmtOk(c, map[string]any{`id`: id})
}

// DeleteDepartment deletes a department
func DeleteDepartment(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	id := cast.ToInt(c.PostForm(`id`))
	newDid := cast.ToInt(c.PostForm(`new_department_id`))
	if id <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	m := msql.Model(`department`, define.Postgres)
	info, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	m.Begin()
	// check department
	var (
		delIds   = []string{cast.ToString(id)}
		children []string
	)
	// get children
	departmentList, _ := common.GetAllDepartmentList(adminUserId)
	common.FindDepartmentChildren(departmentList, cast.ToInt(id), &children)
	if newDid <= 0 {
		delIds = append(delIds, children...)
	}
	// new department
	departmentMbember := msql.Model(`department_member`, define.Postgres)
	if newDid > 0 {
		if tool.InArrayString(cast.ToString(newDid), append(children, cast.ToString(id))) {
			common.FmtError(c, `department_move_err`)
			return
		}
		// check level is max
		if ok, level := common.OverDepartmentLevel(adminUserId, cast.ToInt(newDid), cast.ToInt(id)); ok {
			common.FmtError(c, `department_level_max`, cast.ToString(level))
			return
		}
		// update department
		_, err = m.Where(`pid`, cast.ToString(id)).Update(msql.Datas{
			`pid`:         newDid,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			m.Rollback()
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		// update department member
		_, err = departmentMbember.Where(`department_id`, `in`, cast.ToString(id)).Update(msql.Datas{
			`department_id`: newDid,
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			m.Rollback()
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
	} else {
		// get default department
		var defaultDepartmentId int
		departmentInfo, _ := common.GetDefaultDepartmentInfo(adminUserId)
		if len(departmentInfo) > 0 {
			defaultDepartmentId = cast.ToInt(departmentInfo[`id`])
		}
		// delete department member
		_, err = departmentMbember.Where(`department_id`, `in`, strings.Join(delIds, `,`)).Update(msql.Datas{
			`department_id`: defaultDepartmentId,
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			m.Rollback()
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
	}

	// delete permission
	_, err = msql.Model(`permission_manage`, define.Postgres).
		Where(`identity_id`, `in`, strings.Join(delIds, `,`)).
		Where(`identity_type`, cast.ToString(define.IdentityTypeDepartment)).
		Delete()
	if err != nil {
		m.Rollback()
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	//Delete departments
	_, err = m.Where(`id`, `in`, strings.Join(delIds, `,`)).Delete()
	if err != nil {
		m.Rollback()
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	m.Commit()
	common.FmtOk(c, nil)
}

// BatchUpdateUserDepartment batch-updates user departments
func BatchUpdateUserDepartment(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	//Get params
	departmentIds := strings.TrimSpace(c.PostForm("department_ids"))
	userIds := strings.TrimSpace(c.PostForm("user_ids"))
	if len(departmentIds) <= 0 || len(userIds) == 0 {
		common.FmtError(c, "param_lack")
		return
	}

	//Validate that departments exist
	departmentInfo, err := msql.Model("department", define.Postgres).
		Where("id", `in`, departmentIds).
		Where("admin_user_id", cast.ToString(adminUserId)).
		ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	if len(departmentInfo) != len(strings.Split(departmentIds, ",")) {
		common.FmtError(c, "no_data")
		return
	}

	if err := common.SaveUserDepartmentData(adminUserId, userIds, departmentIds); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err", err.Error())
		return
	}
	common.FmtOk(c, nil)
}
