// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"errors"
	"fmt"
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
	Sys  bool    `json:"sys"`
	Key  string  `json:"key"`
	Desc *string `json:"desc,omitempty"`
	Typ  string  `json:"typ"`
	Vals []Val   `json:"vals,omitempty"`
}

func (field SimpleField) SetVals(data any) SimpleField {
	var datas []any
	if tool.InArrayString(field.Typ, TypArrays[:]) {
		if array, ok := data.([]any); ok {
			datas = append(datas, array...)
		}
	} else {
		datas = append(datas, data)
	}
	if data == nil || len(datas) == 0 {
		field.Vals = nil
		return field
	}
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
	if field.Key == `special.lib_paragraph_list` {
		list := make([]msql.Params, 0)
		for _, val := range field.Vals {
			list = append(list, val.Params)
		}
		_, libraryContent := FormatSystemPrompt(``, list)
		return libraryContent
	}
	temp := make([]string, 0)
	for _, v := range field.GetVals(specifyTyp...) {
		show := fmt.Sprintf(`%v`, v)
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
	*m = MixedInt(cast.ToInt(string(data)))
	return nil
}
