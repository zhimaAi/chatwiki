// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type ChatVariable struct {
	Id             string `form:"id" json:"id"`
	AdminUserId    string `form:"admin_user_id" json:"admin_user_id"`
	RobotId        string `form:"robot_id" json:"robot_id"`
	RobotKey       string `form:"robot_key" json:"robot_key"`
	VariableType   string `form:"variable_type" json:"variable_type"`
	VariableKey    string `form:"variable_key" json:"variable_key"`
	VariableName   string `form:"variable_name" json:"variable_name"`
	MaxInputLength string `form:"max_input_length" json:"max_input_length"`
	DefaultValue   string `form:"default_value" json:"default_value"`
	MustInput      string `form:"must_input" json:"must_input"`
	Options        any    `form:"options" json:"options"`
	CreateTime     string `form:"create_time" json:"create_time"`
	UpdateTime     string `form:"update_time" json:"update_time"`
	Value          any    `form:"value" json:"value"`
}

type ChatVariableOption struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

const (
	VariableTypeInputString    = "input_string"
	VariableTypeInputNumber    = "input_number"
	VariableTypeSelectOne      = "select_one"
	VariableTypeCheckboxSwitch = "checkbox_switch"
)

func (h *ChatVariable) VerifyParams(lang string) error {
	if h.VariableType == `` {
		return errors.New(i18n.Show(lang, "param_err", `variable_type`))
	}
	if h.VariableKey == `` {
		return errors.New(i18n.Show(lang, "param_err", `variable_key`))
	}
	if h.VariableName == `` {
		return errors.New(i18n.Show(lang, "param_err", `variable_name`))
	}
	validTypes := []string{VariableTypeSelectOne, VariableTypeInputString, VariableTypeInputNumber, VariableTypeCheckboxSwitch}
	if !tool.InArrayString(h.VariableType, validTypes) {
		return errors.New(i18n.Show(lang, "param_err", `variable_type`))
	}
	keyRegex := regexp.MustCompile(`^[a-zA-Z_]{1,10}$`)
	if !keyRegex.MatchString(h.VariableKey) {
		return errors.New(i18n.Show(lang, "param_err", `variable_key`))
	}
	if utf8.RuneCountInString(h.VariableName) > 10 {
		return errors.New(i18n.Show(lang, "param_err", `variable_name`))
	}
	maxLen := cast.ToInt(h.MaxInputLength)
	if maxLen < 1 || maxLen > 50 {
		return errors.New(i18n.Show(lang, "param_err", `max_input_length`))
	}
	mustInput := cast.ToInt8(h.MustInput)
	if mustInput != 0 && mustInput != 1 {
		return errors.New(i18n.Show(lang, "param_err", `must_input`))
	}
	if h.VariableType == VariableTypeSelectOne {
		var options = make([]ChatVariableOption, 0)
		if err := tool.JsonDecode(cast.ToString(h.Options), &options); err != nil {
			return errors.New(i18n.Show(lang, "param_err", `options`))
		}
		if len(options) == 0 {
			return errors.New(i18n.Show(lang, "param_err", `options`))
		}
		if len(options) > 100 {
			return errors.New(i18n.Show(lang, "param_err", `options`))
		}
		for opKey, opVal := range options {
			if opVal.Key == `` {
				options[opKey].Key = h.RobotId + tool.Random(10)
			}
		}
		h.Options = tool.JsonEncodeNoError(options)
	}
	if cast.ToInt(h.Id) == 0 {
		num, err := msql.Model(`chat_ai_variables`, define.Postgres).Where(`robot_key`, h.RobotKey).Count(`1`)
		if err != nil {
			return errors.New(i18n.Show(lang, "sys_err"))
		}
		if num > 10 {
			return errors.New(i18n.Show(lang, "max_create", num))
		}
		num, err = msql.Model(`chat_ai_variables`, define.Postgres).
			Where(`robot_key`, h.RobotKey).Where(`variable_key`, h.VariableKey).Count(`1`)
		if err != nil {
			return errors.New(i18n.Show(lang, "sys_err"))
		}
		if num > 0 {
			return errors.New(i18n.Show(lang, "data_exists"))
		}
	} else {
		count, err := msql.Model(`chat_ai_variables`, define.Postgres).
			Where(`robot_key`, h.RobotKey).
			Where(`variable_key`, h.VariableKey).
			Where("id", "<>", h.Id).Count(`1`)
		if err != nil {
			return errors.New(i18n.Show(lang, "sys_err"))
		}
		if count > 0 {
			return errors.New(i18n.Show(lang, "data_exists"))
		}
	}
	return nil
}
