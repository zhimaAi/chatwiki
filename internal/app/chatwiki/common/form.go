// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func GetFormInfo(formId, adminUserId int) (msql.Params, error) {
	form, err := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(formId)).
		Find()
	if err != nil {
		return nil, err
	}
	if len(form) == 0 {
		return nil, nil
	} else {
		entryCount, err := msql.Model(`form_entry`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`form_id`, cast.ToString(formId)).
			Where(`delete_time`, `0`).
			Count()
		if err != nil {
			return nil, err
		}
		form[`entry_count`] = cast.ToString(entryCount)
		return form, nil
	}
}

func GetFormEntryList(adminUserId, formId, filterId, page, size int, startTimeStamp, endTimeStamp int64) ([]msql.Params, int, error) {
	formFields, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).
		Select()
	if err != nil {
		return nil, 0, err
	}

	m := msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).
		Where(`delete_time`, `0`)
	if startTimeStamp > 0 && endTimeStamp > 0 {
		m.Where(`create_time`, `between`, fmt.Sprintf(`%d,%d`, startTimeStamp, endTimeStamp))
	}
	if filterId > 0 {
		filter, err := msql.Model(`form_filter`, define.Postgres).
			Where(`form_id`, cast.ToString(formId)).
			Where(`id`, cast.ToString(filterId)).
			Find()
		if err != nil {
			return nil, 0, err
		}
		if len(filter) == 0 {
			return nil, 0, errors.New(`filter not found`)
		}
		filterConditions, err := msql.Model(`form_filter_condition`, define.Postgres).
			Where(`form_filter_id`, cast.ToString(filterId)).
			Select()
		if err != nil {
			return nil, 0, err
		}
		if len(filterConditions) == 0 {
			return nil, 0, errors.New(`filter condition not found`)
		}
		var formatedFilterConditions []define.FormFilterCondition
		for _, filterCondition := range filterConditions {
			formatedFilterConditions = append(formatedFilterConditions, define.FormFilterCondition{
				FormFieldId: cast.ToInt(filterCondition[`form_field_id`]),
				Rule:        cast.ToString(filterCondition[`rule`]),
				RuleValue1:  cast.ToString(filterCondition[`rule_value1`]),
				RuleValue2:  cast.ToString(filterCondition[`rule_value2`]),
			})
		}
		m.Where(BuildFilterConditionSql(cast.ToInt(filter[`type`]), formatedFilterConditions))
	}

	entries, total, err := m.Order(`id desc`).Field(`id`).Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}
	entryIdList := []string{`-1`}
	for _, entry := range entries {
		entryIdList = append(entryIdList, entry[`id`])
	}
	fieldValues, err := msql.Model(`form_field_value`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_entry_id`, `in`, strings.Join(entryIdList, `,`)).
		Select()
	if err != nil {
		return nil, 0, err
	}
	for _, entry := range entries {
		for _, field := range formFields {
			entry[field[`name`]] = ``
			for _, fieldValue := range fieldValues {
				if entry[`id`] == fieldValue[`form_entry_id`] && field[`id`] == fieldValue[`form_field_id`] {
					if fieldValue[`type`] == `string` {
						entry[field[`name`]] = fieldValue[`string_content`]
					} else if fieldValue[`type`] == `integer` {
						entry[field[`name`]] = fieldValue[`integer_content`]
					} else if fieldValue[`type`] == `number` {
						entry[field[`name`]] = fieldValue[`number_content`]
					} else if fieldValue[`type`] == `boolean` {
						entry[field[`name`]] = fieldValue[`boolean_content`]
					}
				}
			}
		}
	}

	return entries, total, nil
}

func GetFormEntryCountByFilter(adminUserId, formId, _type int, filterConditions []define.FormFilterCondition) (int, error) {
	m := msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).
		Where(`delete_time`, `0`).
		Where(BuildFilterConditionSql(_type, filterConditions))
	return m.Order(`id desc`).Count(`id`)
}

func BuildFilterConditionSql(_type int, filterConditions []define.FormFilterCondition) string {
	if len(filterConditions) == 0 {
		return `` // Prevent concatenating incorrect SQL
	}
	var rawSql []string
	for _, filterCondition := range filterConditions {
		sql := `(`
		sql = sql + `form_field_id=` + cast.ToString(filterCondition.FormFieldId) + ` and `
		if filterCondition.Rule == `string_eq` {
			sql = sql + `string_content = '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `string_neq` {
			sql = sql + `string_content <> '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `string_contain` {
			sql = sql + `string_content like '%` + filterCondition.RuleValue1 + `%'`
		} else if filterCondition.Rule == `string_not_contain` {
			sql = sql + `string_content not like '%` + filterCondition.RuleValue1 + `%'`
		} else if filterCondition.Rule == `string_empty` {
			sql = sql + `string_content = ''`
		} else if filterCondition.Rule == `string_not_empty` {
			sql = sql + `string_content <> ''`
		} else if filterCondition.Rule == `integer_gt` {
			sql = sql + `integer_content > ` + filterCondition.RuleValue1
		} else if filterCondition.Rule == `integer_gte` {
			sql = sql + `integer_content >= ` + filterCondition.RuleValue1
		} else if filterCondition.Rule == `integer_lt` {
			sql = sql + `integer_content < ` + filterCondition.RuleValue1
		} else if filterCondition.Rule == `integer_lte` {
			sql = sql + `integer_content <= ` + filterCondition.RuleValue1
		} else if filterCondition.Rule == `integer_eq` {
			sql = sql + `integer_content = ` + filterCondition.RuleValue1
		} else if filterCondition.Rule == `integer_between` {
			sql = sql + `integer_content between ` + filterCondition.RuleValue1 + ` and ` + filterCondition.RuleValue2
		} else if filterCondition.Rule == `number_gt` {
			sql = sql + `number_content > '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `number_gte` {
			sql = sql + `number_content >= '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `number_lt` {
			sql = sql + `number_content < '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `number_lte` {
			sql = sql + `number_content <= '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `number_eq` {
			sql = sql + `number_content = '` + filterCondition.RuleValue1 + `'`
		} else if filterCondition.Rule == `number_between` {
			sql = sql + `number_content between '` + filterCondition.RuleValue1 + `' and '` + filterCondition.RuleValue2 + `'`
		} else if filterCondition.Rule == `boolean_true` {
			sql = sql + `boolean_content = true`
		} else if filterCondition.Rule == `boolean_false` {
			sql = sql + `boolean_content = false`
		}
		if filterCondition.FormFieldId == 0 { // When FormFieldId==0, it represents form_entry_id, int type
			sql = strings.ReplaceAll(sql, `form_field_id=0 and integer_content`, `form_entry_id`)
		}
		sql = sql + `)`
		rawSql = append(rawSql, sql)
	}
	if _type == 1 {
		var existsSql []string
		for _, sql := range rawSql {
			existsSql = append(existsSql, `(exists (select 1 from form_field_value where form_entry_id = form_entry.id and `+sql+`))`)
		}
		return strings.Join(existsSql, ` and `)
	} else {
		return `exists (select 1 from form_field_value where form_entry_id = form_entry.id and (` + strings.Join(rawSql, ` or `) + `))`
	}
}

func UploadFormFile(adminUserId, formId int, uuid, ext string, fieldsMap map[string]msql.Params, uploadInfoMap []map[string]any) (err error) {
	title := uploadInfoMap[0]
	uploadFormErrData := &define.UploadFormFile{
		Total:   len(uploadInfoMap),
		ErrData: make([]map[string]any, 0),
	}
	errData := make([]map[string]any, 0)
	for _, item := range uploadInfoMap {
		uploadFormErrData.Processed++
		var isErr bool
		if ext == `json` {
			title = item
		}
		for field, fieldData := range fieldsMap {
			if _, ok := title[field]; !ok {
				item[`err_msg`] = `missing field:` + field
				errData = append(errData, item)
				isErr = true
				break
			}
			if cast.ToBool(fieldData[`required`]) && checkIsNull(cast.ToString(fieldData[`type`]), item[field]) {
				item[`err_msg`] = `fields required:` + field
				errData = append(errData, item)
				isErr = true
				break
			}
		}
		if isErr {
			continue
		}
		formEntryId := cast.ToInt(item[`id`])
		err = SaveFormEntry(adminUserId, formId, formEntryId, item)
		if err != nil {
			logs.Error(err.Error())
			item[`err_msg`] = err.Error()
			errData = append(errData, item)
			continue
		}
		_ = SetUploadFormFileProc(uuid, uploadFormErrData, time.Duration(3600))
	}
	uploadFormErrData.Finish = true
	uploadFormErrData.ErrData = errData
	_ = SetUploadFormFileProc(uuid, uploadFormErrData, time.Duration(3600))
	return nil
}
func checkIsNull(fieldType string, data any) bool {
	switch fieldType {
	case `string`, `number`:
		return cast.ToString(data) == ""
	case `integer`:
		return cast.ToInt(data) == 0
	case `boolean`:
		return !cast.ToBool(data)
	}
	return true
}
func SaveFormEntry(adminUserId, formId, formEntryId int, entryValues map[string]any) error {
	if formEntryId > 0 {
		formEntry, err := msql.Model(`form_entry`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(formEntryId)).
			Where(`delete_time`, `0`).
			Find()
		if err != nil {
			return err
		}
		if len(formEntry) == 0 {
			return errors.New(`form entry not found`)
		}
	}
	formFields, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).
		Select()
	if err != nil {
		return err
	}
	var fieldValues []msql.Datas
	for _, formField := range formFields {
		content, ok := entryValues[formField[`name`]]
		if !ok && formEntryId > 0 {
			continue // When updating, allow updating only partial fields
		}
		if cast.ToBool(formField[`required`]) && !ok {
			return errors.New(`lack param field ` + formField[`name`])
		}
		if formField[`type`] == `string` {
			fieldValues = append(fieldValues, msql.Datas{
				`admin_user_id`:  adminUserId,
				`form_field_id`:  formField[`id`],
				`type`:           `string`,
				`string_content`: cast.ToString(content),
				`update_time`:    tool.Time2Int(),
			})
		} else if formField[`type`] == `integer` {
			fieldValues = append(fieldValues, msql.Datas{
				`admin_user_id`:   adminUserId,
				`form_field_id`:   formField[`id`],
				`type`:            `integer`,
				`integer_content`: cast.ToInt(content),
				`update_time`:     tool.Time2Int(),
			})
		} else if formField[`type`] == `number` {
			fieldValues = append(fieldValues, msql.Datas{
				`admin_user_id`:  adminUserId,
				`form_field_id`:  formField[`id`],
				`type`:           `number`,
				`number_content`: cast.ToFloat32(content),
				`update_time`:    tool.Time2Int(),
			})
		} else if formField[`type`] == `boolean` {
			fieldValues = append(fieldValues, msql.Datas{
				`admin_user_id`:   adminUserId,
				`form_field_id`:   formField[`id`],
				`type`:            `boolean`,
				`boolean_content`: cast.ToBool(content),
				`update_time`:     tool.Time2Int(),
			})
		}

	}
	m := msql.Model(`form_field_value`, define.Postgres)
	err = m.Begin()
	if err != nil {
		return err
	}
	if formEntryId > 0 {
		_, err = msql.Model(`form_entry`, define.Postgres).
			Where(`id`, cast.ToString(formEntryId)).
			Update(msql.Datas{`update_time`: tool.Time2Int()})
		if err != nil {
			_ = m.Rollback()
			return err
		}
		for _, fieldValue := range fieldValues {
			wheres := [][]string{
				{`admin_user_id`, cast.ToString(adminUserId)},
				{`form_entry_id`, cast.ToString(formEntryId)},
				{`form_field_id`, cast.ToString(fieldValue[`form_field_id`])},
			}
			fieldValueId, err := m.Where2(wheres).Value(`id`)
			if err != nil {
				_ = m.Rollback()
				return err
			}
			if cast.ToUint(fieldValueId) == 0 { // New field, insert if no data
				fieldValue[`form_entry_id`] = formEntryId
				fieldValue[`create_time`] = tool.Time2Int()
				_, err = m.Insert(fieldValue)
			} else { // Update logic
				_, err = m.Where2(wheres).Update(fieldValue)
			}
			if err != nil {
				_ = m.Rollback()
				return err
			}
		}
	} else {
		formEntryId, err := msql.Model(`form_entry`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`form_id`:       formId,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}, `id`)
		if err != nil {
			_ = m.Rollback()
			return err
		}
		for _, fieldValue := range fieldValues {
			fieldValue[`form_entry_id`] = formEntryId
			fieldValue[`create_time`] = tool.Time2Int()
			_, err = m.Insert(fieldValue)
			if err != nil {
				_ = m.Rollback()
				return err
			}
		}
	}
	err = m.Commit()
	if err != nil {
		return err
	}
	return nil
}

func BuildFunctionTools(formIdList []string, adminUserId int) ([]adaptor.FunctionTool, error) {
	var functionTools []adaptor.FunctionTool
	forms, err := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, `in`, strings.Join(formIdList, `,`)).
		Select()
	if err != nil {
		return nil, err
	}
	for _, form := range forms {
		formFields, err := msql.Model(`form_field`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`form_id`, cast.ToString(form[`id`])).
			Select()
		if err != nil {
			return nil, err
		}

		properties := make(map[string]interface{})
		required := make([]string, 0)
		for _, formField := range formFields {
			properties[formField[`name`]] = map[string]interface{}{
				`type`:        formField[`type`],
				`description`: formField[`description`],
			}
			if cast.ToBool(formField[`required`]) {
				required = append(required, formField[`name`])
			}
		}
		functionTool := adaptor.FunctionTool{
			Name:        form[`name`],
			Description: form[`description`] + `(Do not assign default values to any field)`,
			Parameters: adaptor.Parameters{
				Type:       `object`,
				Properties: properties,
				Required:   required,
			},
		}
		functionTools = append(functionTools, functionTool)
	}
	return functionTools, nil
}

func SaveFormData(adminUserId, robotId int, functionToolCall adaptor.FunctionToolCall) error {
	if _, ok := IsWorkFlowFuncCall(cast.ToString(adminUserId), functionToolCall.Name); ok {
		return nil
	}
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		return err
	}
	if len(robot) == 0 {
		return errors.New(`no robot data`)
	}
	forms, err := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, `in`, robot[`form_ids`]).
		Select()
	if err != nil {
		return err
	}
	for _, form := range forms {
		if functionToolCall.Name == form[`name`] {
			entryValues := make(map[string]any)
			err := json.Unmarshal([]byte(functionToolCall.Arguments), &entryValues)
			if err != nil {
				return err
			}
			err = SaveFormEntry(adminUserId, cast.ToInt(form[`id`]), 0, entryValues)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
