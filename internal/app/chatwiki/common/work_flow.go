// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	TypString    = `string`
	TypNumber    = `number`
	TypBoole     = `boole`
	TypFloat     = `float`
	TypObject    = `object`
	TypParams    = `params`
	TypArrString = `array<string>`
	TypArrNumber = `array<number>`
	TypArrBoole  = `array<boole>`
	TypArrFloat  = `array<float>`
	TypArrObject = `array<object>`
	TypArrParams = `array<params>`
)

var TypScalars = [...]string{
	TypString,
	TypNumber,
	TypBoole,
	TypFloat,
	TypObject,
	TypParams,
}

var TypArrays = [...]string{
	TypArrString,
	TypArrNumber,
	TypArrBoole,
	TypArrFloat,
	TypArrObject,
	TypArrParams,
}

type SimpleFields map[string]SimpleField
type SimpleField struct {
	Sys      bool    `json:"sys"`
	Key      string  `json:"key"`
	Desc     *string `json:"desc,omitempty"`
	Typ      string  `json:"typ"`
	Vals     []Val   `json:"vals,omitempty"`
	Required bool    `json:"required"`
	Default  string  `json:"default,omitempty"`
	Enum     string  `json:"enum,omitempty"` //枚举值
}

func (field SimpleField) SetVals(data any) SimpleField {
	var datas []any
	if array, ok := data.([]any); ok {
		datas = append(datas, array...)
	} else {
		datas = append(datas, data)
	}
	if !tool.InArrayString(field.Typ, TypArrays[:]) && len(datas) > 0 {
		datas = datas[:1]
	}
	if data == nil || len(datas) == 0 {
		field.Vals = nil
		return field
	}
	field.Vals = make([]Val, 0) //清空原有的值
	for _, one := range datas {
		switch field.Typ {
		case TypString, TypArrString:
			if temp, err := cast.ToStringE(one); err == nil {
				field.Vals = append(field.Vals, Val{String: &temp})
			}
		case TypNumber, TypArrNumber:
			if temp, err := cast.ToIntE(one); err == nil {
				field.Vals = append(field.Vals, Val{Number: &temp})
			}
		case TypBoole, TypArrBoole:
			if temp, err := cast.ToBoolE(one); err == nil {
				field.Vals = append(field.Vals, Val{Boole: &temp})
			}
		case TypFloat, TypArrFloat:
			if temp, err := cast.ToFloat64E(one); err == nil {
				field.Vals = append(field.Vals, Val{Float: &temp})
			}
		case TypObject, TypArrObject:
			if one != nil {
				field.Vals = append(field.Vals, Val{Object: one})
			}
		case TypParams, TypArrParams:
			if temp, ok := one.(msql.Params); ok {
				field.Vals = append(field.Vals, Val{Params: temp})
			}
		}
	}
	return field
}

func (field SimpleField) GetVals(specifyTyp ...string) []any {
	vals := make([]any, 0)
	typ := field.Typ
	if len(specifyTyp) > 0 && len(specifyTyp[0]) > 0 {
		typ = specifyTyp[0] //指定类型
	}
	for _, val := range field.Vals {
		switch typ {
		case TypString, TypArrString:
			if val.String != nil {
				vals = append(vals, *val.String)
			}
		case TypNumber, TypArrNumber:
			if val.Number != nil {
				vals = append(vals, *val.Number)
			}
		case TypBoole, TypArrBoole:
			if val.Boole != nil {
				vals = append(vals, *val.Boole)
			}
		case TypFloat, TypArrFloat:
			if val.Float != nil {
				vals = append(vals, *val.Float)
			}
		case TypObject, TypArrObject:
			if val.Object != nil {
				vals = append(vals, val.Object)
			}
		case TypParams, TypArrParams:
			if val.Params != nil {
				vals = append(vals, val.Params)
			}
		}
	}
	if len(vals) == 0 || tool.InArrayString(typ, TypArrays[:]) {
		return vals
	}
	return vals[:1]
}

func (field SimpleField) GetVal(specifyTyp ...string) any {
	vals := field.GetVals(specifyTyp...)
	if len(vals) > 0 {
		return vals[0]
	}
	return nil
}

func (field SimpleField) ShowVals(specifyTyp ...string) string {
	//特殊的变量值,使用特定的输出方式
	if field.Key == `special.lib_paragraph_list` {
		list := make([]msql.Params, 0)
		for _, val := range field.Vals {
			list = append(list, val.Params)
		}
		_, libraryContent := FormatSystemPrompt(``, list)
		return libraryContent
	}
	//标准的变量值,统一格式输出
	temp := make([]string, 0)
	typ := field.Typ
	if len(specifyTyp) > 0 && len(specifyTyp[0]) > 0 {
		typ = specifyTyp[0] //指定类型
	}
	for _, v := range field.GetVals(specifyTyp...) {
		show := fmt.Sprintf(`%v`, v)
		if tool.InArrayString(typ, []string{TypObject, TypParams, TypArrObject, TypArrParams}) {
			show = tool.JsonEncodeNoError(v) //复杂类型json形式输出
		}
		if len(show) > 0 { //filter empty
			temp = append(temp, show)
		}
	}
	return strings.Join(temp, `、`)
}

type Val struct {
	String *string     `json:"string,omitempty"`
	Number *int        `json:"number,omitempty"`
	Boole  *bool       `json:"boole,omitempty"`
	Float  *float64    `json:"float,omitempty"`
	Object any         `json:"object,omitempty"`
	Params msql.Params `json:"params,omitempty"`
}

type RecurveFields []RecurveField
type RecurveField struct {
	SimpleField
	Subs RecurveFields `json:"subs,omitempty"`
}

func SimplifyFields(recurveFields RecurveFields) SimpleFields {
	simpleFields := make(SimpleFields)
	for _, field := range recurveFields {
		if field.Typ != TypObject {
			simpleFields[field.Key] = field.SimpleField
		} else {
			for key, simpleField := range SimplifyFields(field.Subs) {
				simpleField.Key = field.Key + `.` + key //统一key值
				simpleFields[field.Key+`.`+key] = simpleField
			}
		}
	}
	return simpleFields
}

func GetRecurveFields(simpleFields SimpleFields) RecurveFields {
	fieldMap := make(map[string]RecurveField)
	for key, field := range simpleFields {
		if before, after, found := strings.Cut(key, `.`); !found {
			fieldMap[key] = RecurveField{SimpleField: field}
		} else { //递归组装object
			if _, ok := fieldMap[before]; !ok {
				objField := field //组建出来的Object
				objField.Key, objField.Typ, objField.Vals = before, TypObject, nil
				fieldMap[before] = RecurveField{SimpleField: objField, Subs: make(RecurveFields, 0)}
			}
			temp := fieldMap[before]
			field.Key = after //取下一级的key值
			for _, recurveField := range GetRecurveFields(SimpleFields{after: field}) {
				var exist bool
				for idx, sub := range temp.Subs {
					if sub.Key == recurveField.Key {
						exist = true
						temp.Subs[idx].Subs = append(sub.Subs, recurveField.Subs...)
						break
					}
				}
				if !exist {
					temp.Subs = append(temp.Subs, recurveField)
				}
			}
			fieldMap[before] = temp
		}
	}
	recurveFields := make(RecurveFields, 0)
	for _, field := range fieldMap {
		recurveFields = append(recurveFields, field)
	}
	return recurveFields
}

func GetFieldsObject(recurveFields RecurveFields) map[string]any {
	result := make(map[string]any)
	for _, field := range recurveFields {
		if field.Typ == TypObject {
			result[field.Key] = GetFieldsObject(field.Subs)
		} else {
			if tool.InArrayString(field.Typ, TypArrays[:]) {
				result[field.Key] = field.GetVals()
			} else {
				result[field.Key] = field.GetVal()
			}
		}
	}
	return result
}

func (fields RecurveFields) ExtractionData(result map[string]any) RecurveFields {
	if result == nil {
		result = make(map[string]any)
	}
	for i, field := range fields {
		if field.Typ == TypObject {
			if subResult, ok := result[field.Key].(map[string]any); ok {
				fields[i].Subs = fields[i].Subs.ExtractionData(subResult)
			}
		} else {
			fields[i].SimpleField = fields[i].SetVals(result[field.Key])
		}
	}
	return fields
}

func (fields RecurveFields) Verify(parent ...string) error {
	maps := map[string]struct{}{}
	for _, field := range fields {
		fullKey := strings.Join(append(parent, field.Key), `.`)
		if !IsVariableName(field.Key) {
			return errors.New(fmt.Sprintf(`字段[%s]不能为空或格式错误`, fullKey))
		}
		if !tool.InArrayString(field.Typ, TypScalars[:]) && !tool.InArrayString(field.Typ, TypArrays[:]) {
			return errors.New(fmt.Sprintf(`字段类型[%s]不在指定列表范围内`, field.Typ))
		}
		if _, ok := maps[field.Key]; ok {
			return errors.New(fmt.Sprintf(`字段[%s]重复定义`, fullKey))
		}
		maps[field.Key] = struct{}{}
		if field.Typ == TypObject {
			if len(field.Subs) == 0 {
				return errors.New(fmt.Sprintf(`字段[%s]没有添加子项`, fullKey))
			}
			if err := field.Subs.Verify(append(parent, field.Key)...); err != nil {
				return err
			}
		} else if len(field.Subs) > 0 {
			return errors.New(fmt.Sprintf(`字段[%s]不允许添加子项`, fullKey))
		}
	}
	return nil
}

type MixedInt int

func (m *MixedInt) Int() int {
	return int(*m)
}

func (m *MixedInt) UnmarshalJSON(data []byte) error {
	*m = MixedInt(cast.ToInt(strings.Trim(string(data), `"`)))
	return nil
}

var verReg = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

func NextWorkFlowVersion(currentVersion string) string {
	if currentVersion == `` {
		return `0.0.1`
	}
	sm := verReg.FindStringSubmatch(currentVersion)
	if sm == nil {
		return ``
	}
	major, _ := strconv.Atoi(sm[1])
	minor, _ := strconv.Atoi(sm[2])
	patch, _ := strconv.Atoi(sm[3])
	patch++
	if patch >= 20 {
		patch = 0
		minor++
		if minor >= 20 {
			minor = 0
			major++
		}
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

var validSeg = regexp.MustCompile(`^(0|[1-9]\d{0,3})$`) // 0 or 1-999

func ValidateWorkFlowVersion(ver string) bool {
	seg := strings.Split(ver, `.`)
	if len(seg) != 3 {
		return false
	}
	for _, s := range seg {
		if !validSeg.MatchString(s) {
			return false
		}
		n, _ := strconv.Atoi(s)
		if n > 1000 {
			return false
		}
	}
	return true
}
