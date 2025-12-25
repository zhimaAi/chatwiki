// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

const (
	NodeTypeRemark           = -1 //注释
	NodeTypeEdges            = 0  //图的edges
	NodeTypeStart            = 1  //开始节点
	NodeTypeTerm             = 2  //判断分支
	NodeTypeCate             = 3  //问题分类
	NodeTypeCurl             = 4  //http请求
	NodeTypeLibs             = 5  //知识库检索
	NodeTypeLlm              = 6  //AI对话
	NodeTypeFinish           = 7  //结束节点
	NodeTypeAssign           = 8  //赋值节点
	NodeTypeReply            = 9  //指定回复
	NodeTypeManual           = 10 //转人工(xkf)
	NodeTypeQuestionOptimize = 11 // 问题优化
	NodeTypeParamsExtractor  = 12 // 参数提取
	NodeTypeFormInsert       = 13 //数据表单新增
	NodeTypeFormDelete       = 14 //数据表单删除
	NodeTypeFormUpdate       = 15 //数据表单更新
	NodeTypeFormSelect       = 16 //数据表单查询
	NodeTypeCodeRun          = 17 //代码运行
	NodeTypeMcp              = 20 //MCP
	NodeTypePlugin           = 21 //插件
	NodeTypeLoop             = 25 //循环节点
	NodeTypeLoopEnd          = 26 //终止循环节点
	NodeTypeLoopStart        = 27 //开始循环节点
	NodeTypeBatch            = 30 //批处理
	NodeTypeBatchStart       = 31 //批处理开始
	NodeTypeImageGeneration  = 33 //图片生成
)

var NodeTypes = [...]int{
	NodeTypeRemark,
	NodeTypeEdges,
	NodeTypeStart,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeLibs,
	NodeTypeLlm,
	NodeTypeFinish,
	NodeTypeAssign,
	NodeTypeReply,
	NodeTypeManual,
	NodeTypeQuestionOptimize,
	NodeTypeParamsExtractor,
	NodeTypeFormInsert,
	NodeTypeFormDelete,
	NodeTypeFormUpdate,
	NodeTypeFormSelect,
	NodeTypeCodeRun,
	NodeTypeMcp,
	NodeTypeLoop,
	NodeTypeLoopEnd,
	NodeTypeLoopStart,
	NodeTypePlugin,
	NodeTypeBatch,
	NodeTypeBatchStart,
	NodeTypeImageGeneration,
}

type NodeAdapter interface {
	Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error)
}

func GetNodeByKey(flow *WorkFlow, robotId uint, nodeKey string) (NodeAdapter, msql.Params, error) {
	var err error
	var info msql.Params
	if flow != nil && flow.params.Draft.IsDraft {
		info = flow.params.Draft.NodeMaps[nodeKey]
	} else {
		info, err = common.GetRobotNode(robotId, nodeKey)
	}
	if err != nil {
		logs.Error(err.Error())
		return nil, info, err
	}
	if len(info) == 0 {
		return nil, info, errors.New(`节点信息不存在:` + nodeKey)
	}
	nodeType := cast.ToInt(info[`node_type`])
	nodeParams := DisposeNodeParams(nodeType, info[`node_params`])
	switch nodeType {
	case NodeTypeStart:
		return &StartNode{params: nodeParams.Start, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeTerm:
		return &TermNode{params: nodeParams.Term, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeCate:
		return &CateNode{params: nodeParams.Cate, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeCurl:
		return &CurlNode{params: nodeParams.Curl, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLibs:
		return &LibsNode{params: nodeParams.Libs, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLlm:
		return &LlmNode{params: nodeParams.Llm, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFinish:
		return &FinishNode{}, info, nil
	case NodeTypeAssign:
		return &AssignNode{params: nodeParams.Assign, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeReply:
		return &ReplyNode{params: nodeParams.Reply, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeManual:
		return &ManualNode{params: nodeParams.Manual}, info, nil
	case NodeTypeQuestionOptimize:
		return &QuestionOptimizeNode{params: nodeParams.QuestionOptimize, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeParamsExtractor:
		return &ParamsExtractorNode{params: nodeParams.ParamsExtractor, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFormInsert:
		return &FormInsertNode{params: nodeParams.FormInsert, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFormDelete:
		return &FormDeleteNode{params: nodeParams.FormDelete, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFormUpdate:
		return &FormUpdateNode{params: nodeParams.FormUpdate, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFormSelect:
		return &FormSelectNode{params: nodeParams.FormSelect, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeCodeRun:
		return &CodeRunNode{params: nodeParams.CodeRun, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeMcp:
		return &McpNode{params: nodeParams.Mcp, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLoop:
		return &LoopNode{params: nodeParams.Loop, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLoopEnd:
		return &LoopEndNode{nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLoopStart:
		return &LoopStartNode{nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeBatchStart:
		return &BatchStartNode{nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeBatch:
		return &BatchNode{params: nodeParams.Batch, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypePlugin:
		return &PluginNode{params: nodeParams.Plugin, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeImageGeneration:
		return &ImageGeneration{params: nodeParams.ImageGeneration, nextNodeKey: info[`next_node_key`]}, info, nil
	default:
		return nil, info, errors.New(`不支持的节点类型:` + info[`node_type`])
	}
}

type StartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *StartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行开始节点逻辑...`)
	workFlowGlobal := common.RecurveFields{}
	for _, param := range n.params.DiyGlobal { //从配置参数组装
		field := common.SimpleField{Key: param.Key, Desc: &param.Desc, Typ: param.Typ}
		workFlowGlobal = append(workFlowGlobal, common.RecurveField{SimpleField: field})
	}
	if len(flow.params.WorkFlowGlobal) > 0 { //传入参数数据提取
		workFlowGlobal = workFlowGlobal.ExtractionData(flow.params.WorkFlowGlobal)
	}
	for key, field := range common.SimplifyFields(workFlowGlobal) {
		flow.global[key] = field
	}
	if flow.params.Draft.IsDraft { //草稿调试场景
		flow.params.Robot[`question_multiple_switch`] = cast.ToString(cast.ToUint(flow.params.Draft.QuestionMultipleSwitch))
	} else { //触发器逻辑
		var findTrigger bool
		for i, trigger := range n.params.TriggerList {
			if trigger.TriggerSwitch && trigger.TriggerType == flow.params.TriggerParams.TriggerType {
				findTrigger = true
				flow.Logs(`选择使用触发器(%d):%s`, i+1, trigger.TriggerName)
				trigger.SetGlobalValue(flow) //从触发器填充变量值
			}
		}
		if !findTrigger { //没有开启的触发器
			err = errors.New(`当前场景没有对应开启的触发器`)
			return
		}
	}
	output = make(common.SimpleFields) //init
	for key, field := range flow.global {
		output[fmt.Sprintf(`global.%s`, key)] = field
	}
	nextNodeKey = n.nextNodeKey
	return
}

type TermNode struct {
	params      TermNodeParams
	nextNodeKey string
}

func (n *TermNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行判断分支逻辑...`)
	for _, param := range n.params {
		if param.Verify(flow) {
			nextNodeKey = param.NextNodeKey
			return
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

type CateNode struct {
	params      CateNodeParams
	nextNodeKey string
}

func (n *CateNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`执行问题分类逻辑...`)
	debugLog := common.SimpleField{Key: `special.llm_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}
	//part1:prompt
	categorys := make([]string, 0)
	for i, category := range n.params.Categorys {
		categorys = append(categorys, fmt.Sprintf(`%d.%s`, i+1, category.Category))
	}
	prompt := fmt.Sprintf(`## 角色
你是一个超过10年的资深客服，能够准确识别用户的情绪和意图，并将其分类。
## 任务
你需要分析用户的对话内容，识别用户意图，判断用户问题属于哪个分类。
## 分类
%s
## 输出格式
- 返回你认为用户问题归属分类的序号，如果你认为没有合适的分类，返回0。
- 你只能返回分类的序号或者0，否则你将受到惩罚。
- 只需要按要求返回，不要附带你的思考过程。`, strings.Join(categorys, "\n"))
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	//part2:context_qa
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	contextList := common.BuildChatContextPair(openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part3:cur_question
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: question}})
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	//提前给输出日志,避免下面报错丢失日志
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `llm_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `llm_error`, `content`: err}})
	output = common.SimpleFields{debugLog.Key: debugLog}
	if err != nil {
		flow.Logs(`llm请求失败:` + err.Error())
		nextNodeKey = n.nextNodeKey
		return
	}
	output[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.CompletionToken}}}
	output[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.PromptToken}}}
	number, err := cast.ToIntE(chatResp.Result)
	if err != nil || number < 0 || number > len(n.params.Categorys) {
		flow.Logs(`llm返回超出预期:` + chatResp.Result)
		nextNodeKey = n.nextNodeKey
		return
	}
	if number == 0 {
		flow.Logs(`llm判定为不属于列举的分类:` + chatResp.Result)
		nextNodeKey = n.nextNodeKey
		return
	}
	flow.Logs(`llm判定的分类是:(%s)%s`, chatResp.Result, n.params.Categorys[number-1].Category)
	nextNodeKey = n.params.Categorys[number-1].NextNodeKey
	return
}

type CurlNode struct {
	params      CurlNodeParams
	nextNodeKey string
}

func (n *CurlNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行http请求逻辑...`)
	rawurl := n.params.Rawurl
	if len(n.params.Params) > 0 {
		params := make([]string, 0)
		for _, param := range n.params.Params {
			params = append(params, fmt.Sprintf(`%s=%s`, url.QueryEscape(param.Key), url.QueryEscape(flow.VariableReplace(param.Value))))
		}
		if strings.Contains(rawurl, `?`) {
			rawurl += "&" + strings.Join(params, `&`)
		} else {
			rawurl += "?" + strings.Join(params, `&`)
		}
	}
	request := curl.NewRequest(rawurl, n.params.Method)
	for _, header := range n.params.Headers {
		request.Header(header.Key, flow.VariableReplace(header.Value))
	}
	if n.params.Method != http.MethodGet && len(n.params.Body) > 0 {
		switch n.params.Type {
		case TypeUrlencoded:
			for _, param := range n.params.Body {
				request.Param(param.Key, flow.VariableReplace(param.Value))
			}
		case TypeJsonBody:
			request.Body(flow.VariableReplaceJson(n.params.BodyRaw)).
				Header("Content-Type", "application/json")
		}
	}
	if n.params.Timeout > 0 {
		timeout := time.Duration(n.params.Timeout) * time.Second
		request.SetTimeout(timeout, timeout)
	}
	resp, err := request.String()
	flow.Logs(`resp:%s,err:%v`, resp, err)
	if err != nil {
		return
	}
	response, err := request.Response()
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf(`StatusCode:%d`, response.StatusCode))
		return
	}
	result := make(map[string]any)
	err = request.ToJSON(&result)
	if err != nil {
		return
	}
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) //提取数据
	nextNodeKey = n.nextNodeKey
	return
}

type LibsNode struct {
	params      LibsNodeParams
	nextNodeKey string
}

func (n *LibsNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行知识库检索逻辑...`)
	//guise robot
	robot := flow.params.Robot
	robot[`library_ids`] = n.params.LibraryIds
	robot[`search_type`] = cast.ToString(n.params.SearchType)
	robot[`top_k`] = cast.ToString(n.params.TopK)
	robot[`similarity`] = cast.ToString(n.params.Similarity)
	robot[`rerank_status`] = cast.ToString(n.params.RerankStatus)
	robot[`rerank_model_config_id`] = cast.ToString(n.params.RerankModelConfigId)
	robot[`rerank_use_model`] = n.params.RerankUseModel
	//start call
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	list, libUseTime, err := common.GetMatchLibraryParagraphList(
		openid, flow.params.AppType, question, []string{},
		n.params.LibraryIds, n.params.TopK.Int(), n.params.Similarity, n.params.SearchType.Int(), robot,
	)
	isBackground := len(flow.params.Customer) > 0 && cast.ToInt(flow.params.Customer[`is_background`]) > 0
	if !isBackground && len(list) == 0 { //未知问题统计
		common.SaveUnknownIssueRecord(flow.params.AdminUserId, flow.params.Robot, question)
	}
	if err != nil {
		return
	}
	vals := make([]common.Val, 0)
	for _, params := range list {
		vals = append(vals, common.Val{Params: params})
	}
	output = common.SimpleFields{
		`special.lib_use_time`:       common.SimpleField{Key: `special.lib_use_time`, Typ: common.TypObject, Vals: []common.Val{{Object: libUseTime}}},
		`special.lib_paragraph_list`: common.SimpleField{Key: `special.lib_paragraph_list`, Typ: common.TypArrParams, Vals: vals},
	}
	nextNodeKey = n.nextNodeKey
	return
}

type LlmNode struct {
	params      LlmNodeParams
	nextNodeKey string
}

func (n *LlmNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行AI对话逻辑...`)
	debugLog := common.SimpleField{Key: `special.llm_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}
	// check is deep-seek-r1
	isDeepSeek := common.CheckModelIsDeepSeek(n.params.UseModel)
	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	//part1:prompt
	list := make([]msql.Params, 0)
	libsOutput := flow.output
	if len(n.params.LibsNodeKey) > 0 {
		libsOutput = flow.outputs[n.params.LibsNodeKey]
	}
	for _, val := range libsOutput[`special.lib_paragraph_list`].Vals {
		list = append(list, val.Params)
	}
	confPrompt := flow.VariableReplace(n.params.Prompt)
	prompt, libraryContent := common.FormatSystemPrompt(confPrompt, list)
	if isDeepSeek {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: libraryContent})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: libraryContent}})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: prompt})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	}
	//part2:context_qa
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	contextList := common.BuildChatContextPair(openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part3:question,prompt+question
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	if isDeepSeek {
		content := strings.Join([]string{confPrompt, question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: content}})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: question}})
	}
	//append OpenApiContent
	messages = common.BuildOpenApiContent(flow.params.ChatRequestParam, messages)
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	//提前给输出日志,避免下面报错丢失日志
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `llm_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `llm_error`, `content`: err}})
	output = common.SimpleFields{
		`special.lib_use_time`:       libsOutput[`special.lib_use_time`],
		`special.lib_paragraph_list`: libsOutput[`special.lib_paragraph_list`],
		debugLog.Key:                 debugLog,
	}
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	llmTime := int(requestTime)
	output[`special.llm_request_time`] = common.SimpleField{Key: `special.llm_request_time`, Typ: common.TypNumber, Vals: []common.Val{{Number: &llmTime}}}
	output[`special.llm_reply_content`] = common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &chatResp.Result}}}
	output[`special.mcp_reply_content`] = common.SimpleField{Key: `special.mcp_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &chatResp.Result}}}
	output[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.CompletionToken}}}
	output[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.PromptToken}}}
	nextNodeKey = n.nextNodeKey
	return
}

type FinishNode struct{}

func (n *FinishNode) Running(flow *WorkFlow) (output common.SimpleFields, _ string, _ error) {
	flow.Logs(`执行结束节点逻辑...`)
	flow.isFinish = true
	output = flow.output
	return
}

type AssignNode struct {
	params      AssignNodeParams
	nextNodeKey string
}

func (n *AssignNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行赋值分支逻辑...`)
	output = common.SimpleFields{}
	for _, param := range n.params {
		variable, _ := strings.CutPrefix(param.Variable, `global.`)
		field, ok := flow.global[variable]
		if !ok || field.Sys {
			continue //自定义变量不存在
		}
		var data any = flow.VariableReplace(param.Value) //变量替换
		if tool.InArrayString(field.Typ, common.TypArrays[:]) {
			var temp []any //数组类型特殊处理
			for _, item := range strings.Split(cast.ToString(data), `、`) {
				temp = append(temp, item)
			}
			data = temp
		}
		flow.global[variable] = field.SetVals(data) //给自定义全局变量赋值
		output[`global_set.`+variable] = flow.global[variable]
	}
	flow.Logs(`当前global值:%s`, tool.JsonEncodeNoError(flow.global))
	//中间变量的处理
	n.Intermediate(flow, output)
	nextNodeKey = n.nextNodeKey
	return
}

func (n *AssignNode) Intermediate(flow *WorkFlow, output common.SimpleFields) {
	if flow.LoopIntermediate.Params == nil {
		return
	}
	for loopKey, loopParam := range *flow.LoopIntermediate.Params {
		for _, param := range n.params {
			if param.Variable == flow.LoopIntermediate.LoopNodeKey+`.`+loopParam.Key {
				//变量替换 支持数组等
				var data any = flow.VariableReplace(param.Value)
				if tool.InArrayString(loopParam.Typ, common.TypArrays[:]) {
					var temp []any //数组类型特殊处理
					for _, item := range strings.Split(cast.ToString(data), `、`) {
						temp = append(temp, item)
					}
					data = temp
				}
				(*flow.LoopIntermediate.Params)[loopKey].SimpleField = loopParam.SimpleField.SetVals(data)
				output[`loop_intermediate_set.`+loopParam.Key] = (*flow.LoopIntermediate.Params)[loopKey].SimpleField
			}
		}
	}
}

type ReplyNode struct {
	params      ReplyNodeParams
	nextNodeKey string
}

func (n *ReplyNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行指定回复逻辑...`)
	content := flow.VariableReplace(n.params.Content)
	output = common.SimpleFields{
		`special.llm_reply_content`: common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &content}}},
	}
	nextNodeKey = n.nextNodeKey
	return
}

type ManualNode struct {
	params ManualNodeParams
}

func (n *ManualNode) Running(flow *WorkFlow) (_ common.SimpleFields, _ string, err error) {
	flow.Logs(`执行转人工逻辑...`)
	flow.isFinish = true
	err = errors.New(`仅云版支持转人工节点`)
	return
}

type QuestionOptimizeNode struct {
	params      QuestionOptimizeNodeParams
	nextNodeKey string
}

func (n *QuestionOptimizeNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行问题优化逻辑...`)
	debugLog := common.SimpleField{Key: `special.question_optimize_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}
	//part1:prompt
	prompt := define.PromptWorkFlowQuestionOptimize
	userPrompt := flow.VariableReplace(n.params.Prompt)
	if userPrompt != `` {
		prompt += fmt.Sprintf(`\n\n# 对话背景:\n %s`, userPrompt)
	}
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	//part3:context_qa
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	contextList := common.BuildChatContextPair(openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part4:cur_question
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: question}})
	//append OpenApiContent
	messages = common.BuildOpenApiContent(flow.params.ChatRequestParam, messages)
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	//提前给输出日志,避免下面报错丢失日志
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `question_optimize_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `question_optimize_error`, `content`: err}})
	output = common.SimpleFields{
		debugLog.Key: debugLog,
	}
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	llmTime := int(requestTime)
	output[`special.question_optimize_request_time`] = common.SimpleField{Key: `special.question_optimize_request_time`, Typ: common.TypNumber, Vals: []common.Val{{Number: &llmTime}}}
	output[`special.question_optimize_reply_content`] = common.SimpleField{Key: `special.question_optimize_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &chatResp.Result}}}
	output[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.CompletionToken}}}
	output[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.PromptToken}}}
	nextNodeKey = n.nextNodeKey
	return
}

type ParamsExtractorNode struct {
	params      ParamsExtractorNodeParams
	nextNodeKey string
}

func (n *ParamsExtractorNode) Running(flow *WorkFlow) (outputs common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行参数提取逻辑...`)
	debugLog := common.SimpleField{Key: `special.params_extractor_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}
	//part1:prompt
	prompt := define.CompletionGenerateJsonPrompt
	userPrompt := flow.VariableReplace(n.params.Prompt)
	// 获取参数列表
	params := tool.JsonEncodeNoError(n.params.Output)
	if len(params) <= 0 {
		params = `[]`
	}
	prompt = fmt.Sprintf(prompt, userPrompt, params)

	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	//part3:context_qa
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	contextList := common.BuildChatContextPair(openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part4:cur_question
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: question}})
	//append OpenApiContent
	messages = common.BuildOpenApiContent(flow.params.ChatRequestParam, messages)
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	chatResp.Result, _ = strings.CutPrefix(chatResp.Result, "```json")
	chatResp.Result, _ = strings.CutSuffix(chatResp.Result, "```")
	//提前给输出日志,避免下面报错丢失日志
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `params_extractor_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `params_extractor_error`, `content`: err}})
	outputs = common.SimpleFields{
		debugLog.Key: debugLog,
	}
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	var result = make([]map[string]any, 0)
	if err = tool.JsonDecodeUseNumber(chatResp.Result, &result); err != nil {
		err = errors.New(`llm返回数据格式错误:` + err.Error())
		return
	}
	mapResult := make(map[string]any)
	for _, item := range result {
		mapResult[cast.ToString(item[`key`])] = item[`vals`]
	}
	output := common.SimplifyFields(n.params.Output.ExtractionData(mapResult)) //提取数据
	for key, out := range output {
		// 枚举值过滤
		if len(out.Enum) > 0 {
			enumValues := make([]any, 0)
			enums := strings.Split(strings.ReplaceAll(out.Enum, "\n", ","), ",")
			for _, val := range out.GetVals() {
				if tool.InArrayString(cast.ToString(val), enums) {
					enumValues = append(enumValues, val)
				}
			}
			out.Vals = nil
			out = out.SetVals(enumValues)
		}
		if len(out.GetVals()) <= 0 {
			out = out.SetVals(out.Default)
		}
		outputs[key] = out
	}

	llmTime := int(requestTime)
	outputs[`special.params_extractor_request_time`] = common.SimpleField{Key: `special.params_extractor_request_time`, Typ: common.TypNumber, Vals: []common.Val{{Number: &llmTime}}}
	outputs[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.CompletionToken}}}
	outputs[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.PromptToken}}}
	nextNodeKey = n.nextNodeKey
	return
}

type FormInsertNode struct {
	params      FormInsertNodeParams
	nextNodeKey string
}

func (n *FormInsertNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行数据表单新增逻辑...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int()); err != nil {
		return
	}
	entryValues := make(map[string]any)
	for _, field := range n.params.Datas {
		entryValues[field.Name] = flow.VariableReplace(field.Value)
	}
	err = common.SaveFormEntry(flow.params.AdminUserId, n.params.FormId.Int(), 0, entryValues)
	if err != nil {
		return
	}
	nextNodeKey = n.nextNodeKey
	return
}

type FormDeleteNode struct {
	params      FormDeleteNodeParams
	nextNodeKey string
}

func (n *FormDeleteNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行数据表单删除逻辑...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int()); err != nil {
		return
	}
	where := make([]define.FormFilterCondition, len(n.params.Where))
	for idx, condition := range n.params.Where {
		condition.RuleValue1 = flow.VariableReplace(condition.RuleValue1)
		condition.RuleValue2 = flow.VariableReplace(condition.RuleValue2)
		where[idx] = condition
	}
	_, err = msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`form_id`, cast.ToString(n.params.FormId.Int())).Where(`delete_time`, `0`).
		Where(common.BuildFilterConditionSql(n.params.Typ.Int(), where)).
		Update(msql.Datas{`delete_time`: tool.Time2Int()})
	if err != nil {
		return
	}
	nextNodeKey = n.nextNodeKey
	return
}

type FormUpdateNode struct {
	params      FormUpdateNodeParams
	nextNodeKey string
}

func (n *FormUpdateNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行数据表单更新逻辑...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int()); err != nil {
		return
	}
	where := make([]define.FormFilterCondition, len(n.params.Where))
	for idx, condition := range n.params.Where {
		condition.RuleValue1 = flow.VariableReplace(condition.RuleValue1)
		condition.RuleValue2 = flow.VariableReplace(condition.RuleValue2)
		where[idx] = condition
	}
	entryIds, err := msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`form_id`, cast.ToString(n.params.FormId.Int())).Where(`delete_time`, `0`).
		Where(common.BuildFilterConditionSql(n.params.Typ.Int(), where)).ColumnArr(`id`)
	if err != nil {
		return
	}
	if len(entryIds) > 0 {
		errs := make([]string, 0)
		entryValues := make(map[string]any)
		for _, field := range n.params.Datas {
			entryValues[field.Name] = flow.VariableReplace(field.Value)
		}
		for _, entryId := range entryIds {
			err = common.SaveFormEntry(flow.params.AdminUserId, n.params.FormId.Int(), cast.ToInt(entryId), entryValues)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			err = errors.New(strings.Join(errs, "\n"))
			return
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

type FormSelectNode struct {
	params      FormSelectNodeParams
	nextNodeKey string
}

func (n *FormSelectNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行数据表单查询逻辑...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int()); err != nil {
		return
	}
	where := make([]define.FormFilterCondition, len(n.params.Where))
	for idx, condition := range n.params.Where {
		condition.RuleValue1 = flow.VariableReplace(condition.RuleValue1)
		condition.RuleValue2 = flow.VariableReplace(condition.RuleValue2)
		where[idx] = condition
	}
	orderBy := make([]string, 0)
	for _, order := range n.params.Order {
		if order.IsAsc { //升序
			orderBy = append(orderBy, fmt.Sprintf(`%s asc`, order.Name))
		} else { //降序
			orderBy = append(orderBy, fmt.Sprintf(`%s desc`, order.Name))
		}
	}
	entryIds, err := msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`form_id`, cast.ToString(n.params.FormId.Int())).Where(`delete_time`, `0`).
		Where(common.BuildFilterConditionSql(n.params.Typ.Int(), where)).
		Order(strings.Join(orderBy, `,`)).Limit(n.params.Limit.Int()).ColumnArr(`id`)
	if err != nil {
		return
	}
	//拼接数据
	list := make([]map[string]any, len(entryIds))
	if len(entryIds) > 0 {
		fields, ok := msql.Model(`form_field`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
			Where(`form_id`, cast.ToString(n.params.FormId.Int())).ColumnMap(`id,type`, `name`)
		if ok != nil {
			err = ok
			return
		}
		values, ok := msql.Model(`form_field_value`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
			Where(`form_entry_id`, `in`, strings.Join(entryIds, `,`)).
			ColumnMap(`type,string_content,integer_content,number_content,boolean_content`,
				`concat(form_entry_id, '-', form_field_id) uni`)
		if ok != nil {
			err = ok
			return
		}
		for idx, entryId := range entryIds {
			list[idx] = map[string]any{`id`: cast.ToInt(entryId)}
			for _, field := range n.params.Fields {
				item := values[fmt.Sprintf(`%s-%s`, entryId, fields[field.Name][`id`])]
				var value any = item[fmt.Sprintf(`%s_content`, item[`type`])]
				switch fields[field.Name][`type`] {
				case `string`:
					value = cast.ToString(value)
				case `integer`:
					value = cast.ToInt(value)
				case `number`:
					value = cast.ToFloat32(value)
				case `boolean`:
					value = cast.ToBool(value)
				}
				list[idx][field.Name] = value
			}
		}
	}
	//组装输出结果
	vals := make([]common.Val, 0)
	for _, obj := range list {
		vals = append(vals, common.Val{Object: obj})
	}
	total := len(vals)
	output = common.SimpleFields{
		`output_list`: common.SimpleField{Key: `output_list`, Typ: common.TypArrObject, Vals: vals},
		`row_num`:     common.SimpleField{Key: `row_num`, Typ: common.TypNumber, Vals: []common.Val{{Number: &total}}},
	}
	nextNodeKey = n.nextNodeKey
	return
}

type CodeRunNode struct {
	params      CodeRunNodeParams
	nextNodeKey string
}

func (n *CodeRunNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`执行代码运行逻辑...`)
	//开始组装变量
	params := make(map[string]any)
	for _, param := range n.params.Params {
		field, exist := flow.GetVariable(param.Variable)
		if !exist { //变量不存在
			params[param.Field] = nil
			continue
		}
		if tool.InArrayString(field.Typ, common.TypArrays[:]) {
			params[param.Field] = field.GetVals()
		} else {
			params[param.Field] = field.GetVal()
		}
	}
	//开始代码运行
	data := lib_define.CodeRunBody{MainFunc: n.params.MainFunc, Params: params}
	flow.Logs(`body:%s`, tool.JsonEncodeNoError(data))
	timeout := time.Duration(n.params.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	jsonStr, err := common.RequestCodeRun(ctx, `javaScript`, data)
	flow.Logs(`result:%s,err:%v`, jsonStr, err)
	if err != nil {
		nextNodeKey = n.params.Exception
		return //代码运行异常
	}
	result := make(map[string]any)
	if err = tool.JsonDecodeUseNumber(jsonStr, &result); err != nil {
		flow.Logs(`结果解析异常:%s`, err.Error())
		nextNodeKey = n.params.Exception
		return //结果解析异常
	}
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) //提取数据
	nextNodeKey = n.nextNodeKey
	return
}

type McpNode struct {
	params      McpNodeParams
	nextNodeKey string
}

func (n *McpNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行MCP逻辑...`)

	for s, a := range n.params.Arguments {
		n.params.Arguments[s] = flow.VariableReplace(cast.ToString(a))
	}

	// 从数据库获取配置
	provider, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`id`, cast.ToString(n.params.ProviderId)).
		Find()
	if err != nil {
		return nil, "", err
	}

	// 获取工具
	var toolList []mcp.Tool
	err = json.Unmarshal([]byte(provider[`tools`]), &toolList)
	if err != nil {
		return nil, "", err
	}
	var selectedTool mcp.Tool
	found := false
	for _, t := range toolList {
		if t.Name == n.params.ToolName {
			selectedTool = t
			found = true
		}
	}
	if !found {
		return nil, "", errors.New("没有找到对应的工具")
	}

	// 超时配置
	timeout := time.Duration(cast.ToUint(provider["request_timeout"])) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 初始化mcp客户端
	mcpClient, err := common.NewMcpClient(ctx, cast.ToInt(provider[`client_type`]), provider[`url`], provider[`headers`])
	if err != nil {
		return nil, "", fmt.Errorf("mcp客户端初始化失败: %v", err.Error())
	}
	result, err := common.CallTool(ctx, mcpClient, selectedTool, n.params.Arguments)
	if err != nil {
		return nil, "", fmt.Errorf("调用mcp工具出错: %v", err.Error())
	}

	// 构建 output
	output = common.SimpleFields{
		`special.mcp_reply_content`: common.SimpleField{
			Key:  `special.mcp_reply_content`,
			Typ:  common.TypString,
			Vals: []common.Val{{String: &result}},
		},
	}

	nextNodeKey = n.nextNodeKey
	return
}

type LoopNode struct {
	params      LoopNodeParams
	nextNodeKey string
}

func (n *LoopNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

type LoopEndNode struct {
	nextNodeKey string
}

func (n *LoopEndNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

type LoopStartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *LoopStartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行开始节点逻辑...`)
	nextNodeKey = n.nextNodeKey
	return
}

type PluginNode struct {
	params      PluginNodeParams
	nextNodeKey string
}

func (n *PluginNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	u := define.Config.Plugin[`endpoint`] + "/manage/plugin/local-plugins/run"

	// 处理任意深度的 Map 和 Slice
	var deepReplace func(val any) any
	deepReplace = func(val any) any {
		switch v := val.(type) {
		case map[string]any: // 如果是 Map，递归处理 Map 的每一个值
			newMap := make(map[string]any)
			for k, subVal := range v {
				newMap[k] = deepReplace(subVal)
			}
			return newMap
		case []any: // 如果是 Slice，递归处理 Slice 的每一个元素
			newSlice := make([]any, len(v))
			for i, subVal := range v {
				newSlice[i] = deepReplace(subVal)
			}
			return newSlice
		case string: // 只有遇到字符串，才真正执行替换
			return flow.VariableReplace(v)
		default:
			// 其他基本类型（int, bool等），转字符串后尝试替换
			return flow.VariableReplace(cast.ToString(v))
		}
	}

	// 对顶层的每个参数调用递归函数
	for key, value := range n.params.Params {
		n.params.Params[key] = deepReplace(value)
	}

	if n.params.Type == `notice` {
		resp := &lib_web.Response{}
		request := curl.Post(u).Header(`admin_user_id`, cast.ToString(flow.params.AdminUserId))
		request.Param("name", n.params.Name)
		request.Param("action", "default/send-message")
		params, err := json.Marshal(n.params.Params)
		if err != nil {
			return nil, "", err
		}
		err = request.Param("params", string(params)).ToJSON(resp)
		if err != nil {
			return nil, "", err
		}
		if resp.Res != 0 {
			return nil, "", errors.New(resp.Msg)
		}

		result := make(map[string]any)
		err = request.ToJSON(&result)
		if err != nil {
			return
		}
		output = common.SimplifyFields(n.params.Output.ExtractionData(result)) //提取数据
		nextNodeKey = n.nextNodeKey
		return
	} else if n.params.Type == `extension` {
		result := make(map[string]any)
		request := curl.Post(u).Header(`admin_user_id`, cast.ToString(flow.params.AdminUserId))
		request.Param("name", n.params.Name)
		request.Param("action", "default/exec")
		params, err := json.Marshal(n.params.Params)
		if err != nil {
			return nil, "", err
		}
		err = request.Param("params", string(params)).ToJSON(&result)
		if err != nil {
			return nil, "", err
		}
		if cast.ToInt(result["res"]) != 0 {
			return nil, "", errors.New(cast.ToString(result["msg"]))
		}
		output = common.SimplifyFields(n.params.Output.ExtractionData(result)) //提取数据
		nextNodeKey = n.nextNodeKey
		return
	} else {
		nextNodeKey = n.nextNodeKey
		return nil, "", errors.New("暂不支持的插件类型")
	}
}

type BatchStartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *BatchStartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行批处理开始节点逻辑...`)
	nextNodeKey = n.nextNodeKey
	return
}

type BatchNode struct {
	params      BatchNodeParams
	nextNodeKey string
}

func (n *BatchNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

type ImageGeneration struct {
	params      ImageGenerationParams
	nextNodeKey string
}

func (n *ImageGeneration) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行图片生成逻辑...`)
	output = common.SimpleFields{}
	var optimizePromptMode *string
	if n.params.ImageOptimizePrompt == `1` {
		optimizePromptMode = tea.String(`standard`)
	}
	var images = make([]string, 0)
	if len(n.params.InputImages) > 0 {
		for _, image := range n.params.InputImages {
			images = append(images, flow.VariableReplace(image))
		}
	}
	res, err := common.RequestImageGenerate(flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
		cast.ToInt(n.params.ModelConfigId), n.params.UseModel, &adaptor.ZhimaImageGenerationReq{
			Prompt:                    flow.VariableReplace(n.params.Prompt),
			Size:                      &n.params.Size,
			Image:                     &images,
			SequentialImageGeneration: tea.String(`auto`),
			MaxImages:                 cast.ToInt(n.params.ImageNum),
			Stream:                    false,
			ResponseFormat:            tea.String(`b64_json`),
			Watermark:                 tea.Bool(n.params.ImageWatermark == `1`),
			OptimizePromptMode:        optimizePromptMode,
		})
	if err != nil {
		output[`msg`] = common.SimpleField{
			Key:  "msg",
			Typ:  common.TypString,
			Vals: []common.Val{{String: tea.String(err.Error())}},
		}
		nextNodeKey = n.nextNodeKey
		return //结果解析异常
	} else {
		output[`msg`] = common.SimpleField{
			Key:  "msg",
			Typ:  common.TypString,
			Vals: []common.Val{{String: tea.String(`success`)}},
		}
		if res == nil {
			logs.Info(`image generation res is empty %#v`, res)
			return
		}
		for i := 0; i < len(res.Datas); i++ {
			letter := 'a' + rune(i)
			key := fmt.Sprintf(`picture_url_%c`, letter)
			output[key] = common.SimpleField{
				Key:  key,
				Typ:  common.TypString,
				Vals: []common.Val{{String: tea.String(res.Datas[i].Url)}},
			}
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}
