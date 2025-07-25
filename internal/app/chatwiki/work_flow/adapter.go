// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

const (
	NodeTypeRemark = -1 //注释
	NodeTypeEdges  = 0  //图的edges
	NodeTypeStart  = 1  //开始节点
	NodeTypeTerm   = 2  //判断分支
	NodeTypeCate   = 3  //问题分类
	NodeTypeCurl   = 4  //http请求
	NodeTypeLibs   = 5  //知识库检索
	NodeTypeLlm    = 6  //AI对话
	NodeTypeFinish = 7  //结束节点
	NodeTypeAssign = 8  //赋值节点
	NodeTypeReply  = 9  //指定回复
	NodeTypeManual = 10 //转人工(xkf)
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
}

type NodeAdapter interface {
	Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error)
}

func GetNodeByKey(robotId uint, nodeKey string) (NodeAdapter, msql.Params, error) {
	info, err := common.GetRobotNode(robotId, nodeKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, info, err
	}
	if len(info) == 0 {
		return nil, info, errors.New(`节点信息不存在:` + nodeKey)
	}
	nodeParams := NodeParams{}
	if err = tool.JsonDecodeUseNumber(info[`node_params`], &nodeParams); err != nil {
		logs.Error(err.Error())
		return nil, info, fmt.Errorf(`node_params err:%s/%s`, info[`node_params`], err.Error())
	}
	switch cast.ToInt(info[`node_type`]) {
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
		if _, ok := flow.global[key]; !ok {
			flow.global[key] = field //给全局变量赋值,不能覆盖系统级参数
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
	contextList := common.BuildChatContextPair(flow.params.Openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part3:cur_question
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: flow.params.Question})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: flow.params.Question}})
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
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
			request.Body(flow.VariableReplace(n.params.BodyRaw)).
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
	list, libUseTime, err := common.GetMatchLibraryParagraphList(
		flow.params.Openid, flow.params.AppType, flow.params.Question, []string{},
		n.params.LibraryIds, n.params.TopK.Int(), n.params.Similarity, n.params.SearchType.Int(), robot,
	)
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
	for _, val := range flow.output[`special.lib_paragraph_list`].Vals {
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
	contextList := common.BuildChatContextPair(flow.params.Openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part3:question,prompt+question
	if isDeepSeek {
		content := strings.Join([]string{confPrompt, flow.params.Question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: content}})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: flow.params.Question})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: flow.params.Question}})
	}
	//append OpenApiContent
	messages = common.BuildOpenApiContent(flow.params.ChatRequestParam, messages)
	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	//提前给输出日志,避免下面报错丢失日志
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `llm_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `llm_error`, `content`: err}})
	output = common.SimpleFields{
		`special.lib_use_time`:       flow.output[`special.lib_use_time`],
		`special.lib_paragraph_list`: flow.output[`special.lib_paragraph_list`],
		debugLog.Key:                 debugLog,
	}
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	llmTime := int(requestTime)
	output[`special.llm_request_time`] = common.SimpleField{Key: `special.llm_request_time`, Typ: common.TypNumber, Vals: []common.Val{{Number: &llmTime}}}
	output[`special.llm_reply_content`] = common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &chatResp.Result}}}
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
	}
	flow.Logs(`当前global值:%s`, tool.JsonEncodeNoError(flow.global))
	nextNodeKey = n.nextNodeKey
	return
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
