// Copyright © 2016- 2024 Sesame Network Technology all right reserved
package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetLibraryRobotInfo(adminUserId, libraryId int) ([]msql.Params, error) {
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,robot_key,robot_name,robot_intro,robot_avatar,application_type,library_ids`).
		Order(`id desc`).
		Select()
	if err != nil {
		return nil, err
	}
	robotData := []msql.Params{}
	if len(robotInfo) == 0 {
		return []msql.Params{}, nil
	}
	for _, item := range robotInfo {
		libraryIds := strings.Split(cast.ToString(item[`library_ids`]), ",")
		if tool.InArrayString(cast.ToString(libraryId), libraryIds) {
			robotData = append(robotData, item)
		}
	}
	return robotData, nil
}

func GetFormRobotInfo(adminUserId, formId int) ([]msql.Params, error) {
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,robot_key,robot_name,robot_intro,robot_avatar,application_type,form_ids`).
		Order(`id desc`).
		Select()
	if err != nil {
		return nil, err
	}
	robotData := []msql.Params{}
	if len(robotInfo) == 0 {
		return []msql.Params{}, nil
	}
	for _, item := range robotInfo {
		formIds := strings.Split(cast.ToString(item[`form_ids`]), ",")
		if tool.InArrayString(cast.ToString(formId), formIds) {
			robotData = append(robotData, item)
		}
	}
	return robotData, nil
}
