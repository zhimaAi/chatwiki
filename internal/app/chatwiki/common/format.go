// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func FillRobotName(list *[]msql.Params, lang string) error {
	robotIds := make([]string, 0)
	for _, item := range *list {
		if cast.ToInt(item[`robot_id`]) == 0 {
			continue
		}
		boolExist := false
		for _, item2 := range robotIds {
			if item2 == item[`robot_id`] {
				boolExist = true
				break
			}
		}
		if !boolExist {
			robotIds = append(robotIds, item[`robot_id`])
		}
	}
	if len(robotIds) > 0 {
		robots, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, `in`, strings.Join(robotIds, `,`)).Field(`id,robot_name`).Select()
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		for key, item := range *list {
			for _, item2 := range robots {
				if item[`robot_id`] == item2[`id`] {
					(*list)[key][`robot_name`] = item2[`robot_name`]
					break
				}
			}
		}
	}
	for key, item := range *list {
		if cast.ToString(item[`robot_name`]) == `` {
			(*list)[key][`robot_name`] = ``
		}
	}
	return nil
}

func FormatRolesLang(lang string, data []msql.Params) {
	if len(data) == 0 {
		return
	}
	for i, item := range data {
		if item[`parent_id`] != `0` {
			continue
		}
		switch lang {
		case define.LangZhCn:
			if item[`create_name`] != lib_define.System {
				data[i][`create_name`] = lib_define.System
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeRoot && item[`name`] != lib_define.DefaultRoleRoot {
				data[i][`name`] = lib_define.DefaultRoleRoot
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeAdmin && item[`name`] != lib_define.DefaultRoleAdmin {
				data[i][`name`] = lib_define.DefaultRoleAdmin
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeUser && item[`name`] != lib_define.DefaultRoleUser {
				data[i][`name`] = lib_define.DefaultRoleUser
			}
		case define.LangEnUs:
			if item[`create_name`] != lib_define.SystemEnUS {
				data[i][`create_name`] = lib_define.SystemEnUS
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeRoot && item[`name`] != lib_define.DefaultRoleRootEnUS {
				data[i][`name`] = lib_define.DefaultRoleRootEnUS
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeAdmin && item[`name`] != lib_define.DefaultRoleAdminEnUS {
				data[i][`name`] = lib_define.DefaultRoleAdminEnUS
			}
			if cast.ToInt(item[`role_type`]) == define.RoleTypeUser && item[`name`] != lib_define.DefaultRoleUserEnUS {
				data[i][`name`] = lib_define.DefaultRoleUserEnUS
			}
		}
	}
}

func FormatDepartmentsLang(lang string, data []msql.Params) {
	if len(data) == 0 {
		return
	}
	for k, v := range data {
		if cast.ToInt(v[`is_default`]) == define.SwitchOn {
			if lang == define.LangZhCn && v[`department_name`] == lib_define.DefaultDepartmentEnUS {
				data[k][`department_name`] = lib_define.DefaultDepartment
			}
			if lang == define.LangEnUs && v[`department_name`] == lib_define.DefaultDepartment {
				data[k][`department_name`] = lib_define.DefaultDepartmentEnUS
			}
		}
	}
}
