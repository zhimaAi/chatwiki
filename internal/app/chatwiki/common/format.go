// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
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
