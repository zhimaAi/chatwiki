// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

const (
	NodeTypeEdges  = 0 //图的edges
	NodeTypeStart  = 1 //开始节点
	NodeTypeTerm   = 2 //判断分支
	NodeTypeCate   = 3 //问题分类
	NodeTypeCurl   = 4 //http请求
	NodeTypeLibs   = 5 //知识库检索
	NodeTypeLlm    = 6 //AI对话
	NodeTypeFinish = 7 //结束节点
)

var NodeTypes = [...]int{
	NodeTypeEdges,
	NodeTypeStart,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeLibs,
	NodeTypeLlm,
	NodeTypeFinish,
}

type NodeAdapter interface {
	Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error)
}

func GetNodeByKey(robotId uint, nodeKey string) (NodeAdapter, error) {
	info, err := common.GetRobotNode(robotId, nodeKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	if len(info) == 0 {
		return nil, errors.New(`节点信息不存在:` + nodeKey)
	}
	nodeParams := NodeParams{}
	if err = tool.JsonDecodeUseNumber(info[`node_params`], &nodeParams); err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`node_params err:%s/%s`, info[`node_params`], err.Error())
	}
	switch cast.ToInt(info[`node_type`]) {
	case NodeTypeStart:
		return &StartNode{nextNodeKey: info[`next_node_key`]}, nil
	case NodeTypeTerm:
		return &TermNode{params: nodeParams.Term, nextNodeKey: info[`next_node_key`]}, nil
	case NodeTypeCate:
		return &CateNode{params: nodeParams.Cate}, nil
	case NodeTypeCurl:
		return &CurlNode{params: nodeParams.Curl, nextNodeKey: info[`next_node_key`]}, nil
	case NodeTypeLibs:
		return &LibsNode{params: nodeParams.Libs, nextNodeKey: info[`next_node_key`]}, nil
	case NodeTypeLlm:
		return &LlmNode{params: nodeParams.Llm, nextNodeKey: info[`next_node_key`]}, nil
	case NodeTypeFinish:
		return &FinishNode{params: nodeParams}, nil
	default:
		return nil, errors.New(`不支持的节点类型:` + info[`node_type`])
	}
}

type StartNode struct {
	nextNodeKey string
}

func (n *StartNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行开始节点逻辑...`)
	output = SimpleFields{
		`global.question`: flow.global[`question`],
		`global.openid`:   flow.global[`openid`],
	}
	nextNodeKey = n.nextNodeKey
	return
}

type TermNode struct {
	params      TermNodeParams
	nextNodeKey string
}

func (n *TermNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
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
	params CateNodeParams
}

func (n *CateNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行问题分类逻辑...`)
	//系统提示词
	categorys := make([]string, 0)
	for i, category := range n.params.Categorys {
		categorys = append(categorys, fmt.Sprintf(`%d.%s`, i+1, category.Category))
	}
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: fmt.Sprintf(`## 角色
你是一个超过10年的资深客服，能够准确识别用户的情绪和意图，并将其分类。
## 任务
你需要分析用户的对话内容，识别用户意图，判断用户问题属于哪个分类。
## 分类
%s
## 输出格式
- 返回你认为用户问题归属分类的序号，如果你认为没有合适的分类，返回0。
- 你只能返回分类的序号或者0，否则你将受到惩罚。
- 只需要按要求返回，不要附带你的思考过程。`, strings.Join(categorys, "\n"))}}
	//上下文内容
	contextList := common.BuildChatContextPair(flow.params.Openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: flow.params.Question})
	//发起请求
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	number, err := cast.ToIntE(chatResp.Result)
	if err != nil || number < 0 || number > len(n.params.Categorys) {
		err = errors.New(`llm返回超出预期:` + chatResp.Result)
		return
	}
	if number == 0 {
		err = errors.New(`llm判定为不属于列举的分类:` + chatResp.Result)
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

func (n *CurlNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
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
	output = SimplifyFields(n.params.Output.ExtractionData(result)) //提取数据
	nextNodeKey = n.nextNodeKey
	return
}

type LibsNode struct {
	params      LibsNodeParams
	nextNodeKey string
}

func (n *LibsNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
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
	recallStart := time.Now()
	list, err := common.GetMatchLibraryParagraphList(
		flow.params.Openid, flow.params.AppType, flow.params.Question, []string{},
		n.params.LibraryIds, n.params.TopK.Int(), n.params.Similarity, n.params.SearchType.Int(), robot,
	)
	if err != nil {
		return
	}
	recallTime := int(time.Now().Sub(recallStart).Milliseconds())
	vals := make([]Val, 0)
	for _, params := range list {
		vals = append(vals, Val{Params: params})
	}
	output = SimpleFields{
		`special.lib_recall_time`:    SimpleField{Key: `special.lib_recall_time`, Typ: TypNumber, Vals: []Val{{Number: &recallTime}}},
		`special.lib_paragraph_list`: SimpleField{Key: `special.lib_paragraph_list`, Typ: TypArrParams, Vals: vals},
	}
	nextNodeKey = n.nextNodeKey
	return
}

type LlmNode struct {
	params      LlmNodeParams
	nextNodeKey string
}

func (n *LlmNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行AI对话逻辑...`)
	debugLog := SimpleField{Key: `special.llm_debug_log`, Typ: TypArrObject, Vals: []Val{}}
	//part1:prompt
	prompt := flow.VariableReplace(n.params.Prompt) + "\n\n" + define.PromptDefaultAnswerImage
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	debugLog.Vals = append(debugLog.Vals, Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	//part2:library
	for _, val := range flow.output[`special.lib_paragraph_list`].Vals {
		var images []string
		if err = tool.JsonDecodeUseNumber(val.Params[`images`], &images); err != nil {
			logs.Error(err.Error())
		}
		if cast.ToInt(val.Params[`type`]) == define.ParagraphTypeNormal {
			messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: common.EmbTextImages(val.Params[`content`], images)})
			debugLog.Vals = append(debugLog.Vals, Val{Object: map[string]string{`type`: `library`, `content`: common.EmbTextImages(val.Params[`content`], images)}})
		} else {
			messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: "question: " + val.Params[`question`] + "\nanswer: " + common.EmbTextImages(val.Params[`answer`], images)})
			debugLog.Vals = append(debugLog.Vals, Val{Object: map[string]string{`type`: `library`, `content`: "question: " + val.Params[`question`] + "\nanswer: " + common.EmbTextImages(val.Params[`answer`], images)}})
		}
	}
	//part3:context_qa
	contextList := common.BuildChatContextPair(flow.params.Openid, cast.ToInt(flow.params.Robot[`id`]),
		flow.params.DialogueId, flow.params.CurMsgId, n.params.ContextPair.Int())
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		debugLog.Vals = append(debugLog.Vals, Val{Object: map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]}})
	}
	//part4:cur_question
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: flow.params.Question})
	debugLog.Vals = append(debugLog.Vals, Val{Object: map[string]string{`type`: `cur_question`, `content`: flow.params.Question}})
	//append OpenApiContent
	messages = common.BuildOpenApiContent(flow.params.ChatRequestParam, messages)
	//request chat
	chatResp, requestTime, err := common.RequestChat(
		flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	if err != nil {
		err = errors.New(`llm请求失败:` + err.Error())
		return
	}
	llmTime := int(requestTime)
	output = SimpleFields{
		`special.lib_recall_time`:    flow.output[`special.lib_recall_time`],
		`special.lib_paragraph_list`: flow.output[`special.lib_paragraph_list`],
		`special.llm_request_time`:   SimpleField{Key: `special.llm_request_time`, Typ: TypNumber, Vals: []Val{{Number: &llmTime}}},
		`special.llm_reply_content`:  SimpleField{Key: `special.llm_reply_content`, Typ: TypString, Vals: []Val{{String: &chatResp.Result}}},
		debugLog.Key:                 debugLog,
	}
	nextNodeKey = n.nextNodeKey
	return
}

type FinishNode struct {
	params NodeParams
}

func (n *FinishNode) Running(flow *WorkFlow) (output SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`执行结束节点逻辑...`)
	flow.isFinish = true
	return flow.output, ``, nil
}
