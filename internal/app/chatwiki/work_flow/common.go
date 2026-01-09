// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// FindKeyIsUse 寻找key是否在传入的节点中使用
func FindKeyIsUse(nodeList []WorkFlowNode, findKey string) bool {
	var isNeedOpenid bool
	for _, node := range nodeList {
		switch node.NodeType {
		case NodeTypeTerm:
			for _, param := range node.NodeParams.Term {
				for _, term := range param.Terms {
					if term.Variable == findKey {
						return true
					}
				}
			}
		case NodeTypeCurl:
			for _, param := range node.NodeParams.Curl.Headers {
				if strings.Contains(param.Value, findKey) {
					return true
				}
			}
			for _, param := range node.NodeParams.Curl.Params {
				if strings.Contains(param.Value, findKey) {
					return true
				}
			}
			switch node.NodeParams.Curl.Type {
			case TypeUrlencoded:
				for _, param := range node.NodeParams.Curl.Body {
					if strings.Contains(param.Value, findKey) {
						return true
					}
				}
			case TypeJsonBody:
				if strings.Contains(node.NodeParams.Curl.BodyRaw, findKey) {
					return true
				}
			}
		case NodeTypeCate:
			if strings.Contains(node.NodeParams.Cate.Prompt, findKey) {
				return true
			}
			if node.NodeParams.Cate.QuestionValue == findKey {
				return true
			}
			if findKey == `global.question` && len(node.NodeParams.Cate.QuestionValue) == 0 {
				return true //为空的时候使用默认值,所有被使用
			}
			isNeedOpenid = true //此节点存在隐藏的openid参数调用
		case NodeTypeLibs:
			if node.NodeParams.Libs.QuestionValue == findKey {
				return true
			}
			if findKey == `global.question` && len(node.NodeParams.Libs.QuestionValue) == 0 {
				return true //为空的时候使用默认值,所有被使用
			}
			isNeedOpenid = true //此节点存在隐藏的openid参数调用
		case NodeTypeLlm:
			if strings.Contains(node.NodeParams.Llm.Prompt, findKey) {
				return true
			}
			if node.NodeParams.Llm.QuestionValue == findKey {
				return true
			}
			if findKey == `global.question` && len(node.NodeParams.Llm.QuestionValue) == 0 {
				return true //为空的时候使用默认值,所有被使用
			}
			isNeedOpenid = true //此节点存在隐藏的openid参数调用
			if len(node.NodeParams.Llm.LibsNodeKey) > 0 {
				variable := fmt.Sprintf(`%s.%s`, node.NodeParams.Llm.LibsNodeKey, `special.lib_paragraph_list`)
				if variable == findKey {
					return true
				}
			}
		case NodeTypeAssign:
			for _, param := range node.NodeParams.Assign {
				if param.Variable == findKey { //自定义变量不存在
					return true
				}
				if param.Value == findKey {
					return true
				}
			}
		case NodeTypeReply:
			if strings.Contains(node.NodeParams.Reply.Content, findKey) {
				return true
			}
		case NodeTypeQuestionOptimize:
			if strings.Contains(node.NodeParams.QuestionOptimize.Prompt, findKey) {
				return true
			}
			if node.NodeParams.QuestionOptimize.QuestionValue == findKey {
				return true
			}
			if findKey == `global.question` && len(node.NodeParams.QuestionOptimize.QuestionValue) == 0 {
				return true //为空的时候使用默认值,所有被使用
			}
			isNeedOpenid = true //此节点存在隐藏的openid参数调用
		case NodeTypeParamsExtractor:
			if node.NodeParams.ParamsExtractor.QuestionValue == findKey {
				return true
			}
			if findKey == `global.question` && len(node.NodeParams.ParamsExtractor.QuestionValue) == 0 {
				return true //为空的时候使用默认值,所有被使用
			}
			isNeedOpenid = true //此节点存在隐藏的openid参数调用
		case NodeTypeFormInsert:
			for _, field := range node.NodeParams.FormInsert.Datas {
				if strings.Contains(field.Value, findKey) {
					return true
				}
			}
		case NodeTypeFormDelete:
			for _, field := range node.NodeParams.FormDelete.Where {
				if field.RuleValue1 == findKey {
					return true
				}
				if strings.Contains(field.RuleValue2, findKey) {
					return true
				}
			}
		case NodeTypeFormUpdate:
			for _, field := range node.NodeParams.FormUpdate.Where {
				if strings.Contains(field.RuleValue1, findKey) {
					return true
				}
				if strings.Contains(field.RuleValue2, findKey) {
					return true
				}
			}
			for _, field := range node.NodeParams.FormUpdate.Datas {
				if strings.Contains(field.Value, findKey) {
					return true
				}
			}
		case NodeTypeFormSelect:
			for _, field := range node.NodeParams.FormSelect.Where {
				if strings.Contains(field.RuleValue1, findKey) {
					return true
				}
				if strings.Contains(field.RuleValue2, findKey) {
					return true
				}
			}
		case NodeTypeCodeRun:
			for _, param := range node.NodeParams.CodeRun.Params {
				if param.Variable == findKey {
					return true
				}
			}
		case NodeTypeMcp:
			if strings.Contains(node.NodeParams.Mcp.ToolName, findKey) {
				return true
			}
		case NodeTypeLoop:
			for _, field := range node.NodeParams.Loop.LoopArrays {
				if field.Value == findKey {
					return true
				}
			}
			for _, field := range node.NodeParams.Loop.Output {
				if field.Value == findKey {
					return true
				}
			}
		case NodeTypeFinish:
			if node.NodeParams.Finish.OutType == define.FinishNodeOutTypeMessage {
				for _, field := range node.NodeParams.Finish.Messages {
					if strings.Contains(field.Content, findKey) {
						return true
					}
				}
			}
		case NodeTypeImageGeneration:
			for _, param := range node.NodeParams.ImageGeneration.InputImages {
				if strings.Contains(param, findKey) {
					return true
				}
			}
			if strings.Contains(node.NodeParams.ImageGeneration.Prompt, findKey) {
				return true
			}
		case NodeTypeJsonEncode:
			if strings.Contains(node.NodeParams.JsonEncode.InputVariable, findKey) {
				return true
			}
		case NodeTypeJsonDecode:
			if strings.Contains(node.NodeParams.JsonDecode.InputVariable, findKey) {
				return true
			}
		case NodeTypeLibraryImport:
			if strings.Contains(node.NodeParams.LibraryImport.NormalContent, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.NormalTitle, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.NormalUrl, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.QaQuestion, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.QaAnswer, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.QaImagesVariable, findKey) {
				return true
			}
			if strings.Contains(node.NodeParams.LibraryImport.QaSimilarQuestionVariable, findKey) {
				return true
			}
		}
	}
	if findKey == `global.openid` && isNeedOpenid {
		return true //存在隐藏的openid参数调用
	}
	return false
}

// FindNodeByUseKey 寻找key使用的节点
func FindNodeByUseKey(nodeList []WorkFlowNode, findKey string) *WorkFlowNode {
	findKey = strings.TrimPrefix(findKey, `【`)
	findKey = strings.TrimSuffix(findKey, `】`)
	for _, node := range nodeList {
		keys := node.GetVariables()
		if tool.InArray(findKey, keys) {
			return &node
		}
	}
	return nil
}

// TakeTestParams 提取测试参数到全局变量
func TakeTestParams(question, openid, value string, workFlow *map[string]any) []any {
	takeData := make([]any, 0)
	testParamData := make([]TestFillVal, 0)
	err := tool.JsonDecodeUseNumber(value, &testParamData)
	if err != nil {
		logs.Error(`TakeTestParams err:%s`, err.Error())
		return takeData
	}
	if len(testParamData) == 0 {
		return takeData
	}
	(*workFlow)[`question`] = cast.ToString(question)
	(*workFlow)[`openid`] = cast.ToString(openid)
	for _, testParam := range testParamData {
		takeData = append(takeData, testParam)
		if !strings.HasPrefix(testParam.Field.Key, `global`) {
			continue
		}
		globalKey := strings.TrimPrefix(testParam.Field.Key, `global.`)
		if strings.Contains(testParam.Field.Typ, `array`) {
			anys, ok := testParam.Field.Vals.([]any)
			if !ok {
				anys = make([]any, 0)
			}
			if testParam.Field.Typ == common.TypArrObject {
				vals := make([]any, 0)
				for _, anyVal := range anys {
					anyMap := make(map[string]any)
					_ = tool.JsonDecode(cast.ToString(anyVal), &anyMap)
					var anyVal any
					anyVal = anyMap
					vals = append(vals, anyVal)
				}
				(*workFlow)[globalKey] = vals
			} else {
				(*workFlow)[globalKey] = anys
			}
		} else {
			(*workFlow)[globalKey] = cast.ToString(testParam.Field.Vals)
		}
	}
	return takeData
}

// fillTestParamsToRunningParams 提取测试参数至运行时output和global变量
func fillTestParamsToRunningParams(flowL *WorkFlow, params []any) {
	flowL.Logs(`执行批处理节点测试，参数注入...`)
	testParams := make([]TestFillVal, 0)
	for _, field := range params { //从配置参数组装
		fieldParse := TestFillVal{}
		err := tool.JsonDecode(tool.JsonEncodeNoError(field), &fieldParse)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		testParams = append(testParams, fieldParse)
	}
	desc := ``
	workFlowGlobal := common.RecurveFields{}
	workFlowGlobal = append(workFlowGlobal, common.RecurveField{SimpleField: common.SimpleField{Key: `question`, Desc: &desc, Typ: common.TypString}})
	workFlowGlobal = append(workFlowGlobal, common.RecurveField{SimpleField: common.SimpleField{Key: `openid`, Desc: &desc, Typ: common.TypString}})
	for _, fieldParse := range testParams { //从配置参数组装
		if strings.HasPrefix(fieldParse.Field.Key, `global`) {
			field := common.SimpleField{Key: strings.TrimPrefix(fieldParse.Field.Key, `global.`), Desc: &desc, Typ: fieldParse.Field.Typ}
			workFlowGlobal = append(workFlowGlobal, common.RecurveField{SimpleField: field})
		}
	}
	if len(flowL.params.WorkFlowGlobal) > 0 { //传入参数数据提取
		workFlowGlobal = workFlowGlobal.ExtractionData(flowL.params.WorkFlowGlobal)
	}
	for key, field := range common.SimplifyFields(workFlowGlobal) {
		flowL.global[key] = field
	}
	if len(testParams) > 0 {
		for _, fieldParse := range testParams {
			parseVals, ok := fieldParse.Field.Vals.([]any)
			if !ok {
				parseVals = make([]any, 0)
			}
			vals := make([]common.Val, 0)
			switch fieldParse.Field.Typ {
			case common.TypArrFloat:
				for _, val := range parseVals {
					fl := cast.ToFloat64(val)
					vals = append(vals, common.Val{
						Float: &fl,
					})
				}
			case common.TypArrObject:
				for _, val := range parseVals {
					object := make(map[string]any)
					_ = tool.JsonDecode(cast.ToString(val), &object)
					vals = append(vals, common.Val{Object: object})
				}
			case common.TypArrBoole:
				for _, val := range parseVals {
					bl := cast.ToBool(val)
					vals = append(vals, common.Val{
						Boole: &bl,
					})
				}
			case common.TypArrNumber:
				for _, val := range parseVals {
					il := cast.ToInt(val)
					vals = append(vals, common.Val{
						Number: &il,
					})
				}
			case common.TypArrString:
				for _, val := range parseVals {
					sl := cast.ToString(val)
					vals = append(vals, common.Val{
						String: &sl,
					})
				}
			case common.TypString:
				val := cast.ToString(fieldParse.Field.Vals)
				vals = append(vals, common.Val{
					String: &val,
				})
			case common.TypNumber:
				val := cast.ToInt(fieldParse.Field.Vals)
				vals = append(vals, common.Val{
					Number: &val,
				})
			case common.TypBoole:
				val := cast.ToBool(fieldParse.Field.Vals)
				vals = append(vals, common.Val{
					Boole: &val,
				})
			case common.TypObject:
				val := make(map[string]any)
				_ = tool.JsonDecode(cast.ToString(fieldParse.Field.Vals), &val)
				vals = append(vals, common.Val{
					Object: &val,
				})
			case common.TypFloat:
				val := cast.ToFloat64(fieldParse.Field.Vals)
				vals = append(vals, common.Val{
					Float: &val,
				})
			}
			if _, ok := flowL.outputs[fieldParse.NodeKey]; !ok {
				flowL.outputs[fieldParse.NodeKey] = make(common.SimpleFields)
			}
			fieldKey := fieldParse.Field.Key
			if !strings.HasPrefix(fieldKey, `global`) {
				fieldKey = strings.Replace(fieldKey, fieldParse.NodeKey+`.`, ``, 1)
			}
			flowL.outputs[fieldParse.NodeKey][fieldKey] = common.SimpleField{
				Key:  fieldKey,
				Typ:  fieldParse.Field.Typ,
				Vals: vals,
			}
		}
	}
	flowL.Logs(`执行循环节点测试，outputs初始化完成，%s %s`, "\n", tool.JsonEncodeNoError(flowL.outputs))
}
