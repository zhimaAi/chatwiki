// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package work_flow

import (
	"fmt"
	"strings"

	"github.com/zhimaAi/go_tools/tool"
)

// FindKeyIsUse 寻找key是否在传入的节点中使用
func FindKeyIsUse(nodeList []WorkFlowNode, findKey string) bool {
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
		case NodeTypeLlm:
			if strings.Contains(node.NodeParams.Llm.Prompt, findKey) {
				return true
			}
			if node.NodeParams.Llm.QuestionValue == findKey {
				return true
			}
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
		case NodeTypeParamsExtractor:
			if node.NodeParams.ParamsExtractor.QuestionValue == findKey {
				return true
			}
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
		}
	}
	return false
}

// FindNodeByUseKey 寻找key使用的节点
func FindNodeByUseKey(nodeList []WorkFlowNode, findKey string) *WorkFlowNode {
	findKey = strings.TrimPrefix(findKey, `【`)
	findKey = strings.TrimSuffix(findKey, `】`)
	fmt.Println(fmt.Sprintf(`开始查找节点 %s`, findKey))
	for _, node := range nodeList {
		if node.NodeType == NodeTypeStart && (findKey == `global.question` || node.NodeKey == `global.openid`) {
			return &node
		}
		keys := node.GetVariables()
		if tool.InArray(findKey, keys) {
			return &node
		}
	}
	return nil
}
