// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
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

const (
	LoopTypeArray  = `array`
	LoopTypeNumber = `number`
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
	//获取字段的vals值
	typ := field.Typ
	if len(specifyTyp) > 0 && len(specifyTyp[0]) > 0 {
		typ = specifyTyp[0] //指定类型
	}
	vals := field.GetVals(typ)
	//复杂类型json形式输出
	if len(vals) > 0 && tool.InArrayString(typ, []string{TypObject, TypParams, TypArrObject, TypArrParams}) {
		if tool.InArrayString(typ, TypArrays[:]) {
			return tool.JsonEncodeNoError(vals)
		}
		return tool.JsonEncodeNoError(vals[0])
	}
	//标准的变量值,统一格式输出
	temp := make([]string, 0)
	for _, v := range vals {
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

func (fields RecurveFields) SimplifyFieldsDeep(simpleFields *SimpleFields, prefix string) {
	for _, field := range fields {
		if field.Typ != TypObject {
			field.SimpleField.Key = prefix + field.SimpleField.Key
			(*simpleFields)[field.SimpleField.Key] = field.SimpleField
		} else {
			(*simpleFields)[prefix+field.Key] = SimpleField{
				Key: prefix + field.Key,
				Typ: field.Typ,
			}
			if len(field.Subs) > 0 {
				field.Subs.SimplifyFieldsDeep(simpleFields, prefix+field.Key+`.`)
			}
		}
	}
}

func (fields RecurveFields) SimplifyFieldsDeepExtract(simpleFields *SimpleFields, prefix string, data map[string]any) {
	for _, field := range fields {
		if field.Typ != TypObject {
			field.SimpleField = field.SimpleField.SetVals(data[field.Key])
			field.SimpleField.Key = prefix + field.SimpleField.Key
			(*simpleFields)[field.SimpleField.Key] = field.SimpleField
		} else {
			(*simpleFields)[prefix+field.Key] = SimpleField{
				Key:  prefix + field.Key,
				Typ:  field.Typ,
				Vals: []Val{{Object: data[field.Key]}},
			}
			if len(field.Subs) > 0 {
				if subResult, ok := data[field.Key].(map[string]any); ok {
					field.Subs.SimplifyFieldsDeepExtract(simpleFields, prefix+field.Key+`.`, subResult)
				}
			}
		}
	}
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

type LoopField struct {
	SimpleField
	Value string `json:"value"` //node_key.key
}

func (loopField *LoopField) ParseValue() (string, string) {
	params := strings.SplitN(loopField.Value, `.`, 2)
	if len(params) != 2 {
		return ``, ``
	}
	return params[0], params[1]
}

func (loopField *LoopField) NodeKey() (nodeKey string) {
	nodeKey, _ = loopField.ParseValue()
	return
}

func (loopField *LoopField) ChooseKey() (key string) {
	_, key = loopField.ParseValue()
	return
}

type LoopTestParams struct {
	NodeKey  string      `json:"node_key"`
	NodeName string      `json:"node_name"`
	Field    SimpleField `json:"field"`
}

type BatchTestParams struct {
	NodeKey  string      `json:"node_key"`
	NodeName string      `json:"node_name"`
	Field    SimpleField `json:"field"`
}

const WorkFlowEnNameRegex = `[a-zA-Z0-9_\-\.]{1,50}`

func CheckEnName(adminUserId, enName string, id string) bool {
	re := regexp.MustCompile(`^` + WorkFlowEnNameRegex + `$`)
	if !re.MatchString(enName) {
		return false
	}
	id, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`en_name`, enName).Where(`id`, `!=`, id).Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	if len(id) == 0 {
		return true
	}
	return false
}

func TakeWorkFlowTestUseToken(flowOutputs []NodeLog) (int, int64) {
	if len(flowOutputs) == 0 {
		return 0, 0
	}
	token := 0
	type outputs struct {
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
		Output    struct {
			LlmResult struct {
				CompletionToken int `json:"completion_token"`
				PromptToken     int `json:"prompt_token"`
			} `json:"llm_result,omitempty"`
		} `json:"output"`
	}
	nodeOutputs := make([]outputs, 0)
	err := tool.JsonDecode(tool.JsonEncodeNoError(flowOutputs), &nodeOutputs)
	if err != nil {
		logs.Error(err.Error())
		return 0, 0
	}
	var useMills int64 = 0
	for _, nodeOutPut := range nodeOutputs {
		token += nodeOutPut.Output.LlmResult.CompletionToken
		token += nodeOutPut.Output.LlmResult.PromptToken
		useMills += nodeOutPut.EndTime - nodeOutPut.StartTime
	}
	return token, useMills
}

func GetRecursiveFieldsFromMap(data map[string]any) RecurveFields {
	defer func() {
		if r := recover(); r != nil {
			logs.Error(`%v`, r)
		}
	}()
	recurveFields := make(RecurveFields, 0)
	if len(data) == 0 {
		return recurveFields
	}
	for mapKey, mapVal := range data {
		if mapVal == nil {
			continue
		}
		valKind := reflect.TypeOf(mapVal).Kind()
		if valKind == reflect.Map {
			if dataMap, ok := mapVal.(map[string]any); ok {
				recurveField := RecurveField{
					SimpleField: SimpleField{
						Key: mapKey,
						Typ: TypObject,
					},
					Subs: GetRecursiveFieldsFromMap(dataMap),
				}
				recurveFields = append(recurveFields, recurveField)
			}
			continue
		}
		switch valKind {
		case reflect.String:
			recurveFields = append(recurveFields, RecurveField{
				SimpleField: SimpleField{Key: mapKey, Typ: TypString},
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			recurveFields = append(recurveFields, RecurveField{
				SimpleField: SimpleField{Key: mapKey, Typ: TypNumber},
			})
		case reflect.Float32, reflect.Float64:
			if valKind == reflect.Float64 {
				if num, ok := mapVal.(float64); ok {
					// 检查float64是否是整数（小数部分为0）
					if num == float64(int64(num)) {
						recurveFields = append(recurveFields, RecurveField{
							SimpleField: SimpleField{Key: mapKey, Typ: TypNumber},
						})
					} else {
						recurveFields = append(recurveFields, RecurveField{
							SimpleField: SimpleField{Key: mapKey, Typ: TypFloat},
						})
					}
				}
			} else {
				recurveFields = append(recurveFields, RecurveField{
					SimpleField: SimpleField{Key: mapKey, Typ: TypFloat},
				})
			}
		case reflect.Bool:
			recurveFields = append(recurveFields, RecurveField{
				SimpleField: SimpleField{Key: mapKey, Typ: TypBoole},
			})
		case reflect.Slice, reflect.Array:
			arrField := RecurveField{
				SimpleField: SimpleField{Key: mapKey},
			}
			arrVal := reflect.ValueOf(mapVal)
			if arrVal.Len() == 0 {
				continue
			} else {
				var elemKind reflect.Kind
				var elemValue any
				for i := 0; i < arrVal.Len(); i++ {
					elem := arrVal.Index(i)
					if elem.IsNil() {
						continue // 跳过nil元素
					}
					// 安全获取元素类型：处理interface{}包装的情况
					elemKind = elem.Kind()
					if elemKind == reflect.Interface {
						// 解包interface{}获取真实类型
						elemKind = elem.Elem().Kind()
					}
					elemValue = elem.Interface()
					break
				}
				switch elemKind {
				case reflect.String:
					arrField.Typ = TypArrString
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					arrField.Typ = TypArrNumber
				case reflect.Float32, reflect.Float64:
					if elemKind == reflect.Float64 {
						if num, ok := elemValue.(float64); ok {
							// 检查float64是否是整数（小数部分为0）
							if num == float64(int64(num)) {
								arrField.Typ = TypArrNumber
							} else {
								arrField.Typ = TypArrFloat
							}
						}
					} else {
						arrField.Typ = TypArrFloat
					}
				case reflect.Bool:
					arrField.Typ = TypArrBoole
				case reflect.Map:
					arrField.Typ = TypArrObject
				default:
					arrField.Typ = TypArrString
				}
			}
			recurveFields = append(recurveFields, arrField)
		default:
			continue
		}
	}
	return recurveFields
}
