// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	NodeTypeRemark           = -1 // Remark
	NodeTypeEdges            = 0  // Graph edges
	NodeTypeStart            = 1  // Start node
	NodeTypeTerm             = 2  // Conditional branch
	NodeTypeCate             = 3  // Question classification
	NodeTypeCurl             = 4  // HTTP request
	NodeTypeLibs             = 5  // Knowledge base retrieval
	NodeTypeLlm              = 6  // AI conversation
	NodeTypeFinish           = 7  // End node
	NodeTypeAssign           = 8  // Assignment node
	NodeTypeReply            = 9  // Specified reply
	NodeTypeManual           = 10 // Transfer to human (xkf)
	NodeTypeQuestionOptimize = 11 // Question optimization
	NodeTypeParamsExtractor  = 12 // Parameter extraction
	NodeTypeFormInsert       = 13 // Data form insert
	NodeTypeFormDelete       = 14 // Data form delete
	NodeTypeFormUpdate       = 15 // Data form update
	NodeTypeFormSelect       = 16 // Data form query
	NodeTypeCodeRun          = 17 // Code execution
	NodeTypeMcp              = 20 // MCP
	NodeTypePlugin           = 21 // Plugin
	NodeTypeLoop             = 25 // Loop node
	NodeTypeLoopEnd          = 26 // Loop end node
	NodeTypeLoopStart        = 27 // Loop start node
	NodeTypeBatch            = 30 // Batch processing
	NodeTypeBatchStart       = 31 // Batch processing start
	NodeTypeImageGeneration  = 33 // Image generation
	NodeTypeJsonEncode       = 36 // JSON serialization
	NodeTypeJsonDecode       = 37 // JSON deserialization
	NodeTypeTextToAudio      = 38 // Text to audio
	NodeTypeVoiceClone       = 39 // Voice cloning
	NodeTypeLibraryImport    = 40 // Knowledge base import
	NodeTypeWorkflow         = 41 // Workflow
	NodeTypeImmediatelyReply = 42 // Immediate reply
	NodeTypeQuestion         = 43 // Question
	NodeTypeHttpTool         = 45 // HTTP tool
)

const (
	AddToHEADERS = `HEADERS`
	AddToPARAMS  = `PARAMS`
	AddToBODY    = `BODY`
)

var NodeTypes = [...]int{
	NodeTypeRemark,
	NodeTypeEdges,
	NodeTypeStart,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeHttpTool,
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
	NodeTypeJsonEncode,
	NodeTypeJsonDecode,
	NodeTypeTextToAudio,
	NodeTypeVoiceClone,
	NodeTypeLibraryImport,
	NodeTypeWorkflow,
	NodeTypeImmediatelyReply,
	NodeTypeQuestion,
}

type NodeAdapter interface {
	Params() any
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
	lang := define.LangEnUs
	if flow != nil && flow.params != nil {
		lang = flow.params.Lang
	}
	if len(info) == 0 {
		return nil, info, errors.New(i18n.Show(lang, `node_info_not_exist`, nodeKey))
	}
	nodeType := cast.ToInt(info[`node_type`])
	nodeParams := DisposeNodeParams(nodeType, info[`node_params`], lang)
	switch nodeType {
	case NodeTypeStart:
		return &StartNode{params: nodeParams.Start, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeTerm:
		return &TermNode{params: nodeParams.Term, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeCate:
		return &CateNode{params: nodeParams.Cate, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeCurl:
		return &CurlNode{params: nodeParams.Curl, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeHttpTool:
		return &HttpToolNode{params: nodeParams.Curl, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLibs:
		return &LibsNode{params: nodeParams.Libs, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLlm:
		return &LlmNode{params: nodeParams.Llm, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeFinish:
		return &FinishNode{params: nodeParams.Finish}, info, nil
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
	case NodeTypeJsonEncode:
		return &JsonEncode{params: nodeParams.JsonEncode, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeJsonDecode:
		return &JsonDecode{params: nodeParams.JsonDecode, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeTextToAudio:
		return &TextToAudioNode{params: nodeParams.TextToAudio, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeVoiceClone:
		return &VoiceCloneNode{params: nodeParams.VoiceClone, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeLibraryImport:
		return &LibraryImport{params: nodeParams.LibraryImport, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeWorkflow:
		return &WorkflowNode{params: nodeParams.Workflow, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeImmediatelyReply:
		return &ImmediatelyReplyNode{params: nodeParams.ImmediatelyReply, nextNodeKey: info[`next_node_key`]}, info, nil
	case NodeTypeQuestion:
		return &QuestionNode{params: nodeParams.Question, nextNodeKey: info[`next_node_key`]}, info, nil
	default:
		return nil, info, errors.New(i18n.Show(lang, `node_type_not_supported`, info[`node_type`]))
	}
}

type StartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *StartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing start node logic...`)
	workFlowGlobal := common.RecurveFields{}
	for _, param := range n.params.DiyGlobal { // Assemble from configuration parameters
		field := common.SimpleField{Key: param.Key, Desc: &param.Desc, Typ: param.Typ}
		workFlowGlobal = append(workFlowGlobal, common.RecurveField{SimpleField: field})
	}
	if len(flow.params.WorkFlowGlobal) > 0 { // Extract incoming parameter data
		workFlowGlobal = workFlowGlobal.ExtractionData(flow.params.WorkFlowGlobal)
	}
	for key, field := range common.SimplifyFields(workFlowGlobal) {
		flow.global[key] = field
	}
	if flow.params.Draft.IsDraft { // Draft debugging scenario
		flow.params.Robot[`question_multiple_switch`] = cast.ToString(cast.ToUint(flow.params.Draft.QuestionMultipleSwitch))
	} else {
		if flow.params.IsFromWorkflow { // Workflow calling workflow scenario, temporarily no processing

		} else { // Trigger logic
			var findTrigger bool
			for i, trigger := range n.params.TriggerList {
				if trigger.TriggerSwitch && trigger.TriggerType == flow.params.TriggerParams.TriggerType {
					findTrigger = true
					flow.Logs(`Select trigger (%d):%s`, i+1, trigger.TriggerName)
					trigger.SetGlobalValue(flow) // Fill variable values from trigger
				}
			}
			if !findTrigger { // No enabled trigger
				err = errors.New(i18n.Show(flow.params.Lang, `no_enabled_trigger_for_scenario`))
				return
			}
		}
	}
	output = make(common.SimpleFields) //init
	for key, field := range flow.global {
		output[fmt.Sprintf(`global.%s`, key)] = field
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *StartNode) Params() any {
	return n.params
}

type TermNode struct {
	params      TermNodeParams
	nextNodeKey string
}

func (n *TermNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing condition branch logic...`)
	for _, param := range n.params {
		if param.Verify(flow) {
			nextNodeKey = param.NextNodeKey
			return
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *TermNode) Params() any {
	return n.params
}

type CateNode struct {
	params      CateNodeParams
	nextNodeKey string
}

func (n *CateNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`Executing question classification logic...`)
	debugLog := common.SimpleField{Key: `special.llm_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)

	//part1:prompt
	categorys := make([]string, 0)
	for i, category := range n.params.Categorys {
		categorys = append(categorys, fmt.Sprintf(`%d.%s`, i+1, category.Category))
	}
	// Compatible with frontend input
	if n.params.Role >= 1 {
		n.params.Role -= 1
	}
	// question_multiple_switch
	if cast.ToBool(flow.params.Robot[`question_multiple_switch`]) {
		//when question_multiple_switch is true, we will use system role to classify question
		n.params.Role = define.PromptRoleTypeSystem
	}

	prompt := fmt.Sprintf(define.WorkFlowCateNodePrompt, strings.Join(categorys, "\n"))

	if n.params.Role == define.PromptRoleTypeSystem {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: define.PromptRoleTypeMap[define.PromptRoleTypeSystem], Content: prompt})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	} else {
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: ``}})
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
	//part3:cur_question
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}

	//part3:question,prompt+question
	if n.params.Role == define.PromptRoleTypeUser {
		content := strings.Join([]string{prompt, question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: content}})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `cur_question`, `content`: question}})
	}

	//request chat
	flow.params.Robot[`enable_thinking`] = cast.ToString(cast.ToUint(n.params.EnableThinking))
	chatResp, requestTime, err := common.RequestChat(
		flow.params.Lang, flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	// Output log in advance to avoid missing log due to error below
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `llm_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `llm_error`, `content`: err}})
	output = common.SimpleFields{debugLog.Key: debugLog}
	if err != nil {
		flow.Logs(`llm request failed:` + err.Error())
		nextNodeKey = n.nextNodeKey
		return
	}
	output[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.CompletionToken}}}
	output[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &chatResp.PromptToken}}}
	number, err := cast.ToIntE(chatResp.Result)
	if err != nil || number < 0 || number > len(n.params.Categorys) {
		flow.Logs(`llm returned unexpected result:` + chatResp.Result)
		nextNodeKey = n.nextNodeKey
		return
	}
	if number == 0 {
		flow.Logs(`llm determined not belonging to listed categories:` + chatResp.Result)
		nextNodeKey = n.nextNodeKey
		return
	}
	flow.Logs(`llm determined category is:(%s)%s`, chatResp.Result, n.params.Categorys[number-1].Category)
	nextNodeKey = n.params.Categorys[number-1].NextNodeKey
	return
}

func (n *CateNode) Params() any {
	return n.params
}

type CurlNode struct {
	params         CurlNodeParams
	nextNodeKey    string
	httpResultJson *string
}

func (n *CurlNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing http request logic...`)
	rawurl := n.params.Rawurl
	if len(n.params.Params) > 0 {
		params := make([]string, 0)
		for _, param := range n.params.Params {
			params = append(params, fmt.Sprintf(`%s=%s`, url.QueryEscape(param.Key), url.QueryEscape(flow.VariableReplace(param.Value))))
		}

		// add auth params
		if len(n.params.HttpAuth) > 0 {
			for _, authParam := range n.params.HttpAuth {
				if authParam.AddTo == AddToPARAMS {
					params = append(params, fmt.Sprintf(`%s=%s`, url.QueryEscape(authParam.Key), url.QueryEscape(flow.VariableReplace(authParam.Value))))
				}
			}
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

	// add auth params
	if len(n.params.HttpAuth) > 0 {
		for _, authParam := range n.params.HttpAuth {
			if authParam.AddTo == AddToHEADERS {
				request.Header(authParam.Key, flow.VariableReplace(authParam.Value))
			}
		}
	}
	if n.params.Method != http.MethodGet && len(n.params.Body) > 0 {
		switch n.params.Type {
		case TypeUrlencoded:
			for _, param := range n.params.Body {
				request.Param(param.Key, flow.VariableReplace(param.Value))
			}
			// add auth params to body
			if len(n.params.HttpAuth) > 0 {
				for _, authParam := range n.params.HttpAuth {
					if authParam.AddTo == AddToBODY {
						request.Param(authParam.Key, flow.VariableReplace(authParam.Value))
					}
				}
			}

		case TypeJsonBody:
			// json body can not add auth params
			request.Body(flow.VariableReplaceJson(n.params.BodyRaw)).
				Header("Content-Type", "application/json")

		case TypeNone:
			// add auth params to body
			if len(n.params.HttpAuth) > 0 {
				for _, authParam := range n.params.HttpAuth {
					if authParam.AddTo == AddToBODY {
						request.Param(authParam.Key, flow.VariableReplace(authParam.Value))
					}
				}
			}
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
	n.httpResultJson = tea.String(tool.JsonEncodeNoError(result))
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
	nextNodeKey = n.nextNodeKey
	return
}

func (n *CurlNode) GetHttpResultJson() *string {
	return n.httpResultJson
}

func (n *CurlNode) Params() any {
	return n.params
}

type HttpToolNode struct {
	params         CurlNodeParams
	nextNodeKey    string
	httpResultJson *string
}

func (n *HttpToolNode) GetHttpResultJson() *string {
	return n.httpResultJson
}

func (n *HttpToolNode) Params() any {
	return n.params
}

func (n *HttpToolNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	// Ensure HttpToolNode logic is consistent with CurlNode for easier updates
	curlNode := &CurlNode{
		params:         n.params,
		nextNodeKey:    n.nextNodeKey,
		httpResultJson: n.httpResultJson,
	}
	output, nextNodeKey, err = curlNode.Running(flow)
	n.httpResultJson = curlNode.GetHttpResultJson() // Ensure httpResultJson is updated
	return
}

type LibsNode struct {
	params      LibsNodeParams
	nextNodeKey string
}

func (n *LibsNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing knowledge base retrieval logic...`)
	//guise robot
	robot := flow.params.Robot
	robot[`library_ids`] = n.params.LibraryIds
	robot[`search_type`] = cast.ToString(n.params.SearchType)
	robot[`rrf_weight`] = n.params.RrfWeight
	robot[`top_k`] = cast.ToString(n.params.TopK)
	robot[`similarity`] = cast.ToString(n.params.Similarity)
	robot[`rerank_status`] = cast.ToString(n.params.RerankStatus)
	robot[`rerank_model_config_id`] = cast.ToString(n.params.RerankModelConfigId)
	robot[`rerank_use_model`] = n.params.RerankUseModel
	robot[`meta_search_switch`] = cast.ToString(n.params.MetaSearchSwitch)
	robot[`meta_search_type`] = cast.ToString(n.params.MetaSearchType)
	robot[`meta_search_condition_list`] = n.params.MetaSearchConditionList

	robot[`recall_neighbor_switch`] = cast.ToString(n.params.RecallNeighborSwitch)
	robot[`recall_neighbor_after_num`] = cast.ToString(n.params.RecallNeighborAfterNum)
	robot[`recall_neighbor_before_num`] = cast.ToString(n.params.RecallNeighborBeforeNum)
	//start call
	var openid = cast.ToString(flow.global[`openid`].GetVal(common.TypString))
	var question = cast.ToString(flow.global[`question`].GetVal(common.TypString))
	field, exist := flow.GetVariable(n.params.QuestionValue)
	if exist {
		question = field.ShowVals()
	}
	list, libUseTime, err := common.GetMatchLibraryParagraphList(
		flow.params.Lang, openid, flow.params.AppType, question, []string{},
		n.params.LibraryIds, n.params.TopK.Int(), n.params.Similarity, n.params.SearchType.Int(), robot,
	)
	isBackground := len(flow.params.Customer) > 0 && cast.ToInt(flow.params.Customer[`is_background`]) > 0
	if !isBackground && len(list) == 0 { // Unknown question statistics
		common.SaveUnknownIssueRecord(flow.params.Lang, flow.params.AdminUserId, flow.params.Robot, question)
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
		`special.lib_paragraph_list`: common.SimpleField{Key: `special.lib_paragraph_list`, Typ: common.TypArrParams, Vals: vals, Lang: flow.params.Lang},
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *LibsNode) Params() any {
	return n.params
}

type LlmNode struct {
	params      LlmNodeParams
	nextNodeKey string
}

func (n *LlmNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing AI conversation logic...`)
	debugLog := common.SimpleField{Key: `special.llm_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}

	// Compatible with frontend input
	if n.params.Role >= 1 {
		n.params.Role -= 1
	}
	// question_multiple_switch
	if cast.ToBool(flow.params.Robot[`question_multiple_switch`]) {
		//when question_multiple_switch is true, we will use system role to classify question
		n.params.Role = define.PromptRoleTypeSystem
	}

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
	prompt, libraryContent := common.FormatSystemPrompt(flow.params.Lang, confPrompt, list)
	if n.params.Role == define.PromptRoleTypeUser {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: define.PromptRoleTypeMap[define.PromptRoleTypeSystem], Content: libraryContent})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: libraryContent}})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: define.PromptRoleTypeMap[define.PromptRoleTypeSystem], Content: prompt})
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
	if n.params.Role == define.PromptRoleTypeUser {
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
		flow.params.Lang, flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	// Output log in advance to avoid missing log due to error below
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `llm_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `llm_error`, `content`: err}})
	output = common.SimpleFields{
		`special.lib_use_time`:       libsOutput[`special.lib_use_time`],
		`special.lib_paragraph_list`: libsOutput[`special.lib_paragraph_list`],
		debugLog.Key:                 debugLog,
	}
	if err != nil {
		err = errors.New(i18n.Show(flow.params.Lang, `llm_request_failed`, err.Error()))
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

func (n *LlmNode) Params() any {
	return n.params
}

type FinishNode struct {
	params FinishNodeParams
}

func (n *FinishNode) Running(flow *WorkFlow) (output common.SimpleFields, _ string, _ error) {
	flow.Logs(`Executing finish node logic...`)
	flow.isFinish = true
	output = common.SimpleFields{}
	if n.params.OutType == define.FinishNodeOutTypeMessage && len(n.params.Messages) > 0 {
		output = n.getOutPuts(flow)
	} else if n.params.OutType == define.FinishNodeOutTypeVariable {
		output = n.getOutPuts(flow)
		if len(n.params.Outputs) > 0 {
			output[`special.finish_variables`] = common.SimpleField{
				Key: `special.finish_variables`,
				Typ: common.TypString,
				Vals: []common.Val{
					{
						String: tea.String(tool.JsonEncodeNoError(n.getVariables(flow, n.params.Outputs))),
					},
				},
			}
		}
	} else {
		output = flow.output
	}
	return
}
func (n *FinishNode) getVariables(flow *WorkFlow, variables common.RecurveFields) map[string]any {
	var out = make(map[string]any)
	for _, field := range variables {
		var val any
		switch field.Typ {
		case common.TypString:
			val = flow.VariableReplace(cast.ToString(field.Desc))
			out[field.Key] = val
		case common.TypNumber:
			val = flow.VariableReplace(cast.ToString(field.Desc))
			out[field.Key] = cast.ToInt(val)
		case common.TypBoole:
			valReplace := flow.VariableReplace(cast.ToString(field.Desc))
			if strings.ToLower(valReplace) == `true` {
				val = true
			} else {
				val = false
			}
			out[field.Key] = val
		case common.TypFloat:
			val = flow.VariableReplace(cast.ToString(field.Desc))
			out[field.Key] = cast.ToFloat64(val)
		case common.TypObject:
			if len(field.Subs) > 0 {
				out[field.Key] = n.getVariables(flow, field.Subs)
			}
		case common.TypArrString, common.TypArrNumber, common.TypArrFloat, common.TypArrBoole, common.TypArrObject:
			variable := GetFirstVariable(cast.ToString(field.Desc))
			simpleField, _ := flow.GetVariable(variable)
			anys := simpleField.GetVals()
			out[field.Key] = anys
		}
	}
	return out
}

func (n *FinishNode) getOutPuts(flow *WorkFlow) (output common.SimpleFields) {
	output = common.SimpleFields{}
	if len(n.params.Messages) == 1 && n.params.Messages[0].Type == `text` && n.params.Messages[0].Content == `` {
		output = flow.output
		return
	}
	for idx, message := range n.params.Messages {
		content := flow.VariableReplace(message.Content)
		if content == `` {
			continue
		}
		key := fmt.Sprintf(define.FinishReplyPrefixKey+`%s_%d`, message.Type, idx+1)
		output[key] = common.SimpleField{Key: key, Typ: common.TypString, Vals: []common.Val{{String: &content}}}
	}
	return
}

func (n *FinishNode) Params() any {
	return n.params
}

type AssignNode struct {
	params      AssignNodeParams
	nextNodeKey string
}

func (n *AssignNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing assignment branch logic...`)
	output = common.SimpleFields{}
	for _, param := range n.params {
		variable, _ := strings.CutPrefix(param.Variable, `global.`)
		field, ok := flow.global[variable]
		if !ok || field.Sys {
			continue // Custom variable does not exist
		}
		var data any = flow.VariableReplace(param.Value) // Variable replacement
		if tool.InArrayString(field.Typ, common.TypArrays[:]) {
			var temp []any // Special handling for array type
			for _, item := range strings.Split(cast.ToString(data), `、`) {
				temp = append(temp, item)
			}
			data = temp
		}
		flow.global[variable] = field.SetVals(data) // Assign value to custom global variable
		output[`global_set.`+variable] = flow.global[variable]
	}
	flow.Logs(`Current global values:%s`, tool.JsonEncodeNoError(flow.global))
	// Handling intermediate variables
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
				// Variable replacement, supports arrays, etc.
				var data any = flow.VariableReplace(param.Value)
				if tool.InArrayString(loopParam.Typ, common.TypArrays[:]) {
					var temp []any // Special handling for array type
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

func (n *AssignNode) Params() any {
	return n.params
}

type ReplyNode struct {
	params      ReplyNodeParams
	nextNodeKey string
}

func (n *ReplyNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing specified reply logic...`)
	content := flow.VariableReplace(n.params.Content)
	output = common.SimpleFields{
		`special.llm_reply_content`: common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &content}}},
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *ReplyNode) Params() any {
	return n.params
}

type ManualNode struct {
	params ManualNodeParams
}

func (n *ManualNode) Running(flow *WorkFlow) (_ common.SimpleFields, _ string, err error) {
	flow.Logs(`Executing transfer to human logic...`)
	flow.isFinish = true
	err = errors.New(i18n.Show(flow.params.Lang, `cloud_version_manual_node_only`))
	return
}

func (n *ManualNode) Params() any {
	return n.params
}

type QuestionOptimizeNode struct {
	params      QuestionOptimizeNodeParams
	nextNodeKey string
}

func (n *QuestionOptimizeNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing question optimization logic...`)
	debugLog := common.SimpleField{Key: `special.question_optimize_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	// Compatible with frontend input
	if n.params.Role >= 1 {
		n.params.Role -= 1
	}
	// question_multiple_switch
	if cast.ToBool(flow.params.Robot[`question_multiple_switch`]) {
		//when question_multiple_switch is true, we will use system role to classify question
		n.params.Role = define.PromptRoleTypeSystem
	}

	//part1:prompt
	prompt := define.PromptWorkFlowQuestionOptimize
	userPrompt := flow.VariableReplace(n.params.Prompt)
	if userPrompt != `` {
		prompt += fmt.Sprintf(`\n\n# %s:\n %s`, i18n.Show(flow.params.Lang, `dialogue_background`), userPrompt)
	}

	if n.params.Role == define.PromptRoleTypeSystem {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: define.PromptRoleTypeMap[define.PromptRoleTypeSystem], Content: prompt})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	} else {
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: ``}})
	}

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

	if n.params.Role == define.PromptRoleTypeUser {
		content := strings.Join([]string{prompt, question}, "\n\n")
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
		flow.params.Lang, flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	// Output log in advance to avoid missing log due to error below
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `question_optimize_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `question_optimize_error`, `content`: err}})
	output = common.SimpleFields{
		debugLog.Key: debugLog,
	}
	if err != nil {
		err = errors.New(i18n.Show(flow.params.Lang, `llm_request_failed`, err.Error()))
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

func (n *QuestionOptimizeNode) Params() any {
	return n.params
}

type ParamsExtractorNode struct {
	params      ParamsExtractorNodeParams
	nextNodeKey string
}

func (n *ParamsExtractorNode) Running(flow *WorkFlow) (outputs common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing parameter extraction logic...`)
	debugLog := common.SimpleField{Key: `special.params_extractor_debug_log`, Typ: common.TypArrObject, Vals: []common.Val{}}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	// Compatible with frontend input
	if n.params.Role >= 1 {
		n.params.Role -= 1
	}
	// question_multiple_switch
	if cast.ToBool(flow.params.Robot[`question_multiple_switch`]) {
		//when question_multiple_switch is true, we will use system role to classify question
		n.params.Role = define.PromptRoleTypeSystem
	}

	//part1:prompt
	prompt := define.CompletionGenerateJsonPrompt
	userPrompt := flow.VariableReplace(n.params.Prompt)
	// Get parameter list
	params := tool.JsonEncodeNoError(n.params.Output)
	if len(params) <= 0 {
		params = `[]`
	}

	prompt = fmt.Sprintf(prompt, userPrompt, params)
	if n.params.Role == define.PromptRoleTypeSystem {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: define.PromptRoleTypeMap[define.PromptRoleTypeSystem], Content: prompt})
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: prompt}})
	} else {
		debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `prompt`, `content`: ``}})
	}

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

	if n.params.Role == define.PromptRoleTypeUser {
		content := strings.Join([]string{prompt, question}, "\n\n")
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
		flow.params.Lang, flow.params.AdminUserId, openid, flow.params.Robot, flow.params.AppType,
		n.params.ModelConfigId.Int(), n.params.UseModel, messages, nil, n.params.Temperature, n.params.MaxToken.Int(),
	)
	flow.LlmCallLogs(LlmCallInfo{Params: n.params.LlmBaseParams, Messages: messages, ChatResp: chatResp, RequestTime: requestTime, Error: err})
	chatResp.Result, _ = strings.CutPrefix(chatResp.Result, "```json")
	chatResp.Result, _ = strings.CutSuffix(chatResp.Result, "```")
	// Output log in advance to avoid missing log due to error below
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]string{`type`: `params_extractor_answer`, `content`: chatResp.Result}})
	debugLog.Vals = append(debugLog.Vals, common.Val{Object: map[string]any{`type`: `params_extractor_error`, `content`: err}})
	outputs = common.SimpleFields{
		debugLog.Key: debugLog,
	}
	if err != nil {
		err = errors.New(i18n.Show(flow.params.Lang, `llm_request_failed`, err.Error()))
		return
	}
	var result = make([]map[string]any, 0)
	if err = tool.JsonDecodeUseNumber(chatResp.Result, &result); err != nil {
		err = errors.New(i18n.Show(flow.params.Lang, `llm_response_format_error`, err.Error()))
		return
	}
	mapResult := make(map[string]any)
	for _, item := range result {
		mapResult[cast.ToString(item[`key`])] = item[`vals`]
	}
	output := common.SimplifyFields(n.params.Output.ExtractionData(mapResult)) // Extract data
	for key, out := range output {
		// Filter enum values
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

func (n *ParamsExtractorNode) Params() any {
	return n.params
}

type FormInsertNode struct {
	params      FormInsertNodeParams
	nextNodeKey string
}

func (n *FormInsertNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing form data insertion logic...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int(), flow.params.Lang); err != nil {
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

func (n *FormInsertNode) Params() any {
	return n.params
}

type FormDeleteNode struct {
	params      FormDeleteNodeParams
	nextNodeKey string
}

func (n *FormDeleteNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing form data deletion logic...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int(), flow.params.Lang); err != nil {
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

func (n *FormDeleteNode) Params() any {
	return n.params
}

type FormUpdateNode struct {
	params      FormUpdateNodeParams
	nextNodeKey string
}

func (n *FormUpdateNode) Running(flow *WorkFlow) (_ common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing form data update logic...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int(), flow.params.Lang); err != nil {
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

func (n *FormUpdateNode) Params() any {
	return n.params
}

type FormSelectNode struct {
	params      FormSelectNodeParams
	nextNodeKey string
}

func (n *FormSelectNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing form data query logic...`)
	if err = checkFormId(flow.params.AdminUserId, n.params.FormId.Int(), flow.params.Lang); err != nil {
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
		if order.IsAsc { // Ascending
			orderBy = append(orderBy, fmt.Sprintf(`%s asc`, order.Name))
		} else { // Descending
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
	// Concatenate data
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
	// Assemble output results
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

func (n *FormSelectNode) Params() any {
	return n.params
}

type CodeRunNode struct {
	params      CodeRunNodeParams
	nextNodeKey string
}

func (n *CodeRunNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`Executing code execution logic...`)
	// Start assembling variables
	params := make(map[string]any)
	for _, param := range n.params.Params {
		field, exist := flow.GetVariable(param.Variable)
		if !exist { // Variable does not exist
			params[param.Field] = nil
			continue
		}
		if tool.InArrayString(field.Typ, common.TypArrays[:]) {
			params[param.Field] = field.GetVals()
		} else {
			params[param.Field] = field.GetVal()
		}
	}
	// Start code execution
	data := lib_define.CodeRunBody{MainFunc: n.params.MainFunc, Params: params}
	flow.Logs(`body:%s`, tool.JsonEncodeNoError(data))
	timeout := time.Duration(n.params.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	jsonStr, err := common.RequestCodeRun(ctx, n.params.Language, data)
	flow.Logs(`result:%s,err:%v`, jsonStr, err)
	if err != nil {
		nextNodeKey = n.params.Exception
		return // Code execution exception
	}
	result := make(map[string]any)
	if err = tool.JsonDecodeUseNumber(jsonStr, &result); err != nil {
		flow.Logs(`Result parsing exception:%s`, err.Error())
		nextNodeKey = n.params.Exception
		return // Result parsing exception
	}
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
	nextNodeKey = n.nextNodeKey
	return
}

func (n *CodeRunNode) Params() any {
	return n.params
}

type McpNode struct {
	params      McpNodeParams
	nextNodeKey string
}

func (n *McpNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing MCP logic...`)

	for s, a := range n.params.Arguments {
		n.params.Arguments[s] = flow.VariableReplace(cast.ToString(a))
	}

	// Get configuration from database
	provider, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`id`, cast.ToString(n.params.ProviderId)).
		Find()
	if err != nil {
		return nil, "", err
	}

	// Get tools
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
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `mcp_tool_not_found`))
	}

	// Timeout configuration
	timeout := time.Duration(cast.ToUint(provider["request_timeout"])) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Initialize MCP client
	mcpClient, err := common.NewMcpClient(ctx, cast.ToInt(provider[`client_type`]), provider[`url`], provider[`headers`])
	if err != nil {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `mcp_client_init_failed`, err.Error()))
	}
	result, err := common.CallTool(ctx, mcpClient, selectedTool, n.params.Arguments)
	if err != nil {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `mcp_tool_call_failed`, err.Error()))
	}

	// Build output
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

func (n *McpNode) Params() any {
	return n.params
}

type LoopNode struct {
	params      LoopNodeParams
	nextNodeKey string
}

func (n *LoopNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

func (n *LoopNode) Params() any {
	return n.params
}

type LoopEndNode struct {
	nextNodeKey string
}

func (n *LoopEndNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

func (n *LoopEndNode) Params() any {
	return nil
}

type LoopStartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *LoopStartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing start node logic...`)
	nextNodeKey = n.nextNodeKey
	return
}

func (n *LoopStartNode) Params() any {
	return n.params
}

type PluginNode struct {
	params      PluginNodeParams
	nextNodeKey string
}

func (n *PluginNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	u := define.Config.Plugin[`endpoint`] + "/manage/plugin/local-plugins/run"

	// Handle Map and Slice of arbitrary depth
	var deepReplace func(val any) any
	deepReplace = func(val any) any {
		switch v := val.(type) {
		case map[string]any: // If it is a Map, recursively process each value of the Map
			newMap := make(map[string]any)
			for k, subVal := range v {
				newMap[k] = deepReplace(subVal)
			}
			return newMap
		case []any: // If it is a Slice, recursively process each element of the Slice
			newSlice := make([]any, len(v))
			for i, subVal := range v {
				newSlice[i] = deepReplace(subVal)
			}
			return newSlice
		case string: // Only execute replacement when encountering a string
			return flow.VariableReplace(v)
		default:
			// Other basic types (int, bool, etc.), try to replace after converting to string
			return flow.VariableReplace(cast.ToString(v))
		}
	}

	// Call recursive function for each top-level parameter
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
		output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
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
		output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
		nextNodeKey = n.nextNodeKey
		return
	} else {
		nextNodeKey = n.nextNodeKey
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `plugin_type_not_supported`))
	}
}

func (n *PluginNode) Params() any {
	return n.params
}

type BatchStartNode struct {
	params      StartNodeParams
	nextNodeKey string
}

func (n *BatchStartNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing batch start node logic...`)
	nextNodeKey = n.nextNodeKey
	return
}

func (n *BatchStartNode) Params() any {
	return n.params
}

type BatchNode struct {
	params      BatchNodeParams
	nextNodeKey string
}

func (n *BatchNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	nextNodeKey = n.nextNodeKey
	return
}

func (n *BatchNode) Params() any {
	return n.params
}

type ImageGeneration struct {
	params      ImageGenerationParams
	nextNodeKey string
}

func (n *ImageGeneration) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing image generation logic...`)
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
	res, err := common.RequestImageGenerate(flow.params.Lang, flow.params.AdminUserId, flow.params.Openid, flow.params.Robot, flow.params.AppType,
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
		return // Result parsing exception
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

func (n *ImageGeneration) Params() any {
	return n.params
}

type JsonEncode struct {
	params      JsonEncodeParams
	nextNodeKey string
}

func (n *JsonEncode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing JSON serialization logic...`)
	variable := GetFirstVariable(n.params.InputVariable)
	if variable == `` {
		nextNodeKey = n.nextNodeKey
		return
	}
	output = common.SimpleFields{}
	field, exist := flow.GetVariable(variable)
	if !exist || len(n.params.Outputs) == 0 {
		flow.Logs(`Variable %s does not exist or output is empty`, variable)
	} else {
		out := n.params.Outputs[0]
		if tool.InArray(field.Typ, []string{common.TypObject, common.TypParams}) {
			vals := field.GetVals()
			if len(vals) == 0 {
				nextNodeKey = n.nextNodeKey
				return
			}
			output[out.Key] = common.SimpleField{
				Key:  out.Key,
				Typ:  common.TypString,
				Vals: []common.Val{{String: tea.String(tool.JsonEncodeNoError(vals[0]))}},
			}
		} else if tool.InArray(field.Typ, []string{common.TypArrString, common.TypArrFloat, common.TypArrObject,
			common.TypArrNumber, common.TypArrParams, common.TypArrBoole,
		}) {
			output[out.Key] = common.SimpleField{
				Key:  out.Key,
				Typ:  common.TypString,
				Vals: []common.Val{{String: tea.String(tool.JsonEncodeNoError(field.GetVals()))}},
			}
		} else {
			output[out.Key] = common.SimpleField{
				Key:  out.Key,
				Typ:  common.TypString,
				Vals: []common.Val{{String: tea.String(``)}},
			}
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *JsonEncode) Params() any {
	return n.params
}

type JsonDecode struct {
	params      JsonDecodeParams
	nextNodeKey string
}

func (n *JsonDecode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing JSON deserialization logic...`)
	output = common.SimpleFields{}
	variable := n.params.InputVariable
	if variable == `` {
		nextNodeKey = n.nextNodeKey
		return
	}
	field, exist := flow.GetVariable(variable)
	if !exist {
		flow.Logs(`Variable %s does not exist`, variable)
	} else {
		vals := field.GetVals()
		if len(vals) == 0 {
			nextNodeKey = n.nextNodeKey
			return
		}
		data := ``
		if tool.InArray(field.Typ, []string{common.TypObject, common.TypParams}) {
			data = tool.JsonEncodeNoError(vals[0])
		} else if field.Typ == common.TypString {
			data = cast.ToString(vals[0])
		} else {
			nextNodeKey = n.nextNodeKey
			return
		}
		if data != `` {
			result := make(map[string]any)
			err := tool.JsonDecode(data, &result)
			if err != nil {
				flow.Logs(`JSON deserialization failed:%s,%s`, data, err.Error())
			} else if len(n.params.Outputs) == 1 {
				resOutput := n.params.Outputs[0]
				var resVal any
				resVal = result
				output[resOutput.Key] = common.SimpleField{
					Key:  resOutput.Key,
					Typ:  common.TypObject,
					Vals: []common.Val{{Object: resVal}},
				}
				if len(resOutput.Subs) > 0 {
					resOutput.Subs.SimplifyFieldsDeepExtract(&output, resOutput.Key+`.`, result) // Extract data
				}
			}
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *JsonDecode) Params() any {
	return n.params
}

// TextToAudioNode Text to Speech
type TextToAudioNode struct {
	params      TextToAudioNodeParams
	nextNodeKey string
}

func (n *TextToAudioNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`Executing text to speech logic...`)

	// Replace variables
	text := flow.VariableReplace(n.params.Arguments.Text)
	voiceId := flow.VariableReplace(n.params.Arguments.VoiceSetting.VoiceId)

	// Build parameters
	params := make(map[string]any)
	params[`model`] = n.params.Arguments.UseModel
	params[`text`] = text
	params[`output_format`] = n.params.Arguments.OutputFormat

	// Build voice_setting
	voiceSetting := make(map[string]any)
	voiceSetting[`voice_id`] = voiceId
	if n.params.Arguments.VoiceSetting.Speed > 0 {
		voiceSetting[`speed`] = float64(n.params.Arguments.VoiceSetting.Speed) / 50.0 // Convert to 0.5-2.0 range
	}
	if n.params.Arguments.VoiceSetting.Vol > 0 {
		voiceSetting[`vol`] = n.params.Arguments.VoiceSetting.Vol
	}
	if n.params.Arguments.VoiceSetting.Pitch != 0 {
		voiceSetting[`pitch`] = n.params.Arguments.VoiceSetting.Pitch
	}
	if len(n.params.Arguments.VoiceSetting.Emotion) > 0 {
		voiceSetting[`emotion`] = flow.VariableReplace(n.params.Arguments.VoiceSetting.Emotion)
	}
	params[`voice_setting`] = voiceSetting

	// Build audio_setting
	if n.params.Arguments.AudioSetting.SampleRate > 0 || n.params.Arguments.AudioSetting.Bitrate > 0 ||
		len(n.params.Arguments.AudioSetting.Format) > 0 || n.params.Arguments.AudioSetting.Channel > 0 {
		audioSetting := make(map[string]any)
		if n.params.Arguments.AudioSetting.SampleRate > 0 {
			audioSetting[`sample_rate`] = n.params.Arguments.AudioSetting.SampleRate
		}
		if n.params.Arguments.AudioSetting.Bitrate > 0 {
			audioSetting[`bitrate`] = n.params.Arguments.AudioSetting.Bitrate
		}
		if len(n.params.Arguments.AudioSetting.Format) > 0 {
			audioSetting[`format`] = flow.VariableReplace(n.params.Arguments.AudioSetting.Format)
		}
		if n.params.Arguments.AudioSetting.Channel > 0 {
			audioSetting[`channel`] = n.params.Arguments.AudioSetting.Channel
		}
		params[`audio_setting`] = audioSetting
	}
	voiceModify := make(map[string]any)
	voiceModify[`pitch`] = n.params.Arguments.VoiceModify.Intensity
	voiceModify[`speed`] = n.params.Arguments.VoiceModify.Pitch
	voiceModify[`volume`] = n.params.Arguments.VoiceModify.Timbre
	if len(n.params.Arguments.VoiceModify.SoundEffects) > 0 {
		voiceModify[`sound_effects`] = n.params.Arguments.VoiceModify.SoundEffects
	}

	// Call API
	result, err := common.TtsSpeechT2A(flow.params.Lang, flow.params.AdminUserId, n.params.Arguments.ModelId, n.params.Arguments.UseModel, params)
	if err != nil {
		flow.Logs(`Text to speech failed: %v`, err)
		return nil, "", err
	}

	flow.Logs(`Text to speech success`)

	// Build output
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
	nextNodeKey = n.nextNodeKey
	return
}

func (n *TextToAudioNode) Params() any {
	return n.params
}

// VoiceCloneNode Voice Cloning
type VoiceCloneNode struct {
	params      VoiceCloneNodeParams
	nextNodeKey string
}

func (n *VoiceCloneNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	flow.Logs(`Executing voice cloning logic...`)

	// Get model configuration
	config, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`model_define`, common.ModelMinimax).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `get_model_config_failed`))
	}
	if len(config) == 0 {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `minimax_model_not_configured`))
	}

	// Replace variables
	fileUrl := flow.VariableReplace(n.params.Arguments.FileUrl)
	voiceId := flow.VariableReplace(n.params.Arguments.VoiceId)

	// Download file to temporary location
	tempDir := `/tmp/voice_clone_` + tool.Random(8)
	if err = tool.MkDirAll(tempDir); err != nil {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `create_temp_dir_failed`, err.Error()))
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Download cloned audio file
	cloneAudioPath := tempDir + `/clone_audio.mp3`
	if err = common.DownloadFile(fileUrl, cloneAudioPath); err != nil {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `download_clone_audio_failed`, err.Error()))
	}

	// Upload cloned audio to get file_id
	flow.Logs(`Uploading cloned audio...`)
	uploadResult, err := common.TtsUploadVoiceFile(flow.params.Lang, flow.params.AdminUserId, cast.ToInt(config["id"]), "voice_clone", cloneAudioPath)
	if err != nil {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `upload_clone_audio_failed`, err.Error()))
	}

	var cloneFileID int64
	if fileInfo, ok := uploadResult[`file`].(map[string]any); ok {
		if fileID, ok := fileInfo[`file_id`].(float64); ok {
			cloneFileID = int64(fileID)
		}
	}
	if cloneFileID <= 0 {
		return nil, "", errors.New(i18n.Show(flow.params.Lang, `get_clone_audio_file_id_failed`))
	}

	flow.Logs(`Cloned audio uploaded successfully, file_id: %d`, cloneFileID)

	// Process example audio (optional)
	var promptFileID int64
	if len(n.params.Arguments.ClonePrompt.PromptAudioUrl) > 0 || len(n.params.Arguments.ClonePrompt.PromptText) > 0 {
		promptAudioUrl := flow.VariableReplace(n.params.Arguments.ClonePrompt.PromptAudioUrl)
		promptText := flow.VariableReplace(n.params.Arguments.ClonePrompt.PromptText)

		if len(promptAudioUrl) > 0 && len(promptText) > 0 {
			// Download example audio file
			promptAudioPath := tempDir + `/prompt_audio.mp3`
			if err = common.DownloadFile(promptAudioUrl, promptAudioPath); err != nil {
				return nil, "", errors.New(i18n.Show(flow.params.Lang, `download_prompt_audio_failed`, err.Error()))
			}

			// Upload example audio to get file_id
			flow.Logs(`Uploading example audio...`)
			promptUploadResult, err := common.TtsUploadVoiceFile(flow.params.Lang, flow.params.AdminUserId, cast.ToInt(config["id"]), "prompt_audio", promptAudioPath)
			if err != nil {
				return nil, "", errors.New(i18n.Show(flow.params.Lang, `upload_prompt_audio_failed`, err.Error()))
			}

			if fileInfo, ok := promptUploadResult[`file`].(map[string]any); ok {
				if fileID, ok := fileInfo[`file_id`].(float64); ok {
					promptFileID = int64(fileID)
				}
			}
			if promptFileID <= 0 {
				return nil, "", errors.New(i18n.Show(flow.params.Lang, `get_prompt_audio_file_id_failed`))
			}
			flow.Logs(`Example audio uploaded successfully, file_id: %d`, promptFileID)
		}
	}

	// Build parameters
	params := make(map[string]any)
	params[`file_id`] = cloneFileID
	params[`voice_id`] = voiceId

	// Build clone_prompt (optional)
	if promptFileID > 0 && len(n.params.Arguments.ClonePrompt.PromptText) > 0 {
		clonePrompt := make(map[string]any)
		clonePrompt[`prompt_audio`] = promptFileID
		clonePrompt[`prompt_text`] = flow.VariableReplace(n.params.Arguments.ClonePrompt.PromptText)
		params[`clone_prompt`] = clonePrompt
	}

	// Add audition parameters (optional)
	if len(n.params.Arguments.Text) > 0 {
		params[`text`] = flow.VariableReplace(n.params.Arguments.Text)
	}

	// Add model parameters (optional)
	if len(n.params.Arguments.Model) > 0 {
		params[`model`] = flow.VariableReplace(n.params.Arguments.Model)
	}

	// Add language boost parameters (optional)
	if len(n.params.Arguments.LanguageBoost) > 0 {
		params[`language_boost`] = flow.VariableReplace(n.params.Arguments.LanguageBoost)
	}

	// Add noise reduction parameters (optional)
	if n.params.Arguments.NeedNoiseReduction {
		params[`need_noise_reduction`] = true
	}

	// Add volume normalization parameters (optional)
	if n.params.Arguments.NeedVolumeNormalization {
		params[`need_volume_normalization`] = true
	}

	// Call cloning API
	flow.Logs(`Starting voice cloning...`)
	result, err := common.TtsCloneVoice(flow.params.Lang, flow.params.AdminUserId, cast.ToInt(config["id"]), params)
	if err != nil {
		flow.Logs(`Voice cloning failed: %v`, err)
		return nil, "", err
	}

	flow.Logs(`Voice cloning success`)

	// Build output
	output = common.SimplifyFields(n.params.Output.ExtractionData(result)) // Extract data
	nextNodeKey = n.nextNodeKey
	return
}

func (n *VoiceCloneNode) Params() any {
	return n.params
}

type LibraryImport struct {
	params      LibraryImportParams
	nextNodeKey string
}

func (n *LibraryImport) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing knowledge base import logic...`)
	output = common.SimpleFields{}
	libraryInfo, err := common.GetLibraryInfo(cast.ToInt(n.params.LibraryId), flow.params.AdminUserId)
	if err != nil {
		return
	}
	if len(libraryInfo) == 0 {
		return
	}
	msg, ok := n.toImport(flow, libraryInfo)
	if ok {
		output[`msg`] = common.SimpleField{
			Key:  `msg`,
			Typ:  common.TypString,
			Vals: []common.Val{{String: tea.String(`ok`)}},
		}
	} else {
		output[`msg`] = common.SimpleField{
			Key:  `msg`,
			Typ:  common.TypString,
			Vals: []common.Val{{String: tea.String(msg)}},
		}
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *LibraryImport) toImport(flow *WorkFlow, libraryInfo msql.Params) (msg string, ok bool) {
	token := flow.params.HeaderToken
	if token == `` {
		info, err := msql.Model(define.TableUser, define.Postgres).
			Where(`id`, cast.ToString(flow.params.AdminUserId)).
			Field(`id,user_name,parent_id`).Find()
		if err != nil {
			logs.Error(err.Error())
			msg = i18n.Show(flow.params.Lang, `get_account_info_failed`)
			return
		}
		if len(info) == 0 {
			msg = i18n.Show(flow.params.Lang, `get_account_info_failed`)
			return
		}
		jwtToken, err := common.GetToken(info[`id`], info[`user_name`], info["parent_id"])
		if err != nil {
			logs.Error(err.Error())
			msg = i18n.Show(flow.params.Lang, `get_token_failed`)
			return
		}
		token = cast.ToString(jwtToken[`token`])
	}
	if cast.ToInt(libraryInfo[`type`]) == define.GeneralLibraryType {
		if n.params.ImportType == define.LibraryImportContent {
			msg, ok = n.toImportFile(flow, token, libraryInfo)
		} else if n.params.ImportType == define.LibraryImportUrl {
			normalUrl := flow.VariableReplace(n.params.NormalUrl)
			if len(normalUrl) == 0 {
				msg = i18n.Show(flow.params.Lang, `import_url_empty`)
				return
			}
			urlInfo, err := msql.Model(`chat_ai_library_file`, define.Postgres).
				Where(`library_id`, libraryInfo[`id`]).Where(`doc_url`, normalUrl).Find()
			if err != nil {
				logs.Error(fmt.Sprintf(`Failed to get document info %s`, err.Error()))
				msg = i18n.Show(flow.params.Lang, `get_doc_info_failed`)
				return
			}
			if n.params.NormalUrlRepeatOp == define.LibraryImportRepeatImport {
				msg, ok = n.toImportUrl(flow.params.Lang, token, normalUrl)
			} else if n.params.NormalUrlRepeatOp == define.LibraryImportRepeatNotImport {
				if len(urlInfo) > 0 {
					msg = i18n.Show(flow.params.Lang, `import_url_exist_skip`, normalUrl)
					ok = true
				} else {
					msg, ok = n.toImportUrl(flow.params.Lang, token, normalUrl)
				}
			} else if n.params.NormalUrlRepeatOp == define.LibraryImportRepeatUpdate {
				msg, ok = n.toUpdateUrl(flow.params.Lang, token, urlInfo)
			}
		}
	} else if cast.ToInt(libraryInfo[`type`]) == define.QALibraryType {
		question := flow.VariableReplace(n.params.QaQuestion)
		answer := flow.VariableReplace(n.params.QaAnswer)
		if question == `` || answer == `` {
			msg = i18n.Show(flow.params.Lang, `import_qa_empty`)
			return
		}
		var images []any
		images = []any{}
		imageVariables := GetVariablePlaceholders(n.params.QaImagesVariable)
		for _, imageVariable := range imageVariables {
			imageField, exist := flow.GetVariable(imageVariable)
			if exist {
				images = append(images, imageField.GetVals()...)
			}
		}
		sourceUrls := common.GetUrls(n.params.QaImagesVariable)
		for _, sourceUrl := range sourceUrls {
			images = append(images, sourceUrl)
		}
		imagesNew := make([]any, 0)
		for _, image := range images {
			if common.IsUrl(cast.ToString(image)) {
				imagesNew = append(imagesNew, image)
			}
		}
		var similars []any
		similars = []any{}
		similarVariables := GetVariablePlaceholders(n.params.QaSimilarQuestionVariable)
		for _, similarVariable := range similarVariables {
			similarField, exist := flow.GetVariable(similarVariable)
			if exist {
				similars = append(similars, similarField.GetVals()...)
			}
		}
		similarReplace := RemoveVariablePlaceholders(n.params.QaSimilarQuestionVariable)
		if len(similarReplace) > 0 {
			similars = append(similars, similarReplace)
		}
		qaInfo, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`library_id`, libraryInfo[`id`]).Where(`question`, n.params.QaQuestion).Find()
		if err != nil {
			logs.Error(fmt.Sprintf(`Failed to get QA info %s`, err.Error()))
			msg = i18n.Show(flow.params.Lang, `get_qa_info_failed`)
			return
		}
		if n.params.QaRepeatOp == define.LibraryImportRepeatImport {
			msg, ok = n.toImportQa(flow.params.Lang, token, question, answer, imagesNew, similars, nil)
		} else if n.params.QaRepeatOp == define.LibraryImportRepeatNotImport {
			if len(qaInfo) > 0 {
				msg = i18n.Show(flow.params.Lang, `import_qa_exist_skip`, question)
				ok = true
			} else {
				msg, ok = n.toImportQa(flow.params.Lang, token, question, answer, imagesNew, similars, nil)
			}
		} else if n.params.QaRepeatOp == define.LibraryImportRepeatUpdate {
			if len(qaInfo) > 0 {
				msg, ok = n.toImportQa(flow.params.Lang, token, question, answer, imagesNew, similars, qaInfo)
			} else {
				msg, ok = n.toImportQa(flow.params.Lang, token, question, answer, imagesNew, similars, nil)
			}
		}
	} else {
		msg = i18n.Show(flow.params.Lang, `library_import_type_not_supported`)
	}
	return
}

func (n *LibraryImport) toImportFile(flow *WorkFlow, token string, libraryInfo msql.Params) (msg string, ok bool) {
	content := flow.VariableReplace(n.params.NormalContent)
	if len(content) == 0 {
		msg = i18n.Show(flow.params.Lang, `import_content_empty`)
		return
	}
	title := flow.VariableReplace(n.params.NormalTitle)
	fileName := define.UploadDir + `work_flow/library_import/` + title + `.txt`
	err := tool.WriteFile(fileName, content)
	if err != nil {
		logs.Error(`Failed to create file %s: %v\n`, fileName, err)
		msg = i18n.Show(flow.params.Lang, `create_file_failed`)
		return
	}
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			logs.Error(`Failed to delete file %s: %v\n`, fileName, err)
		}
	}()
	var res lib_web.Response
	err = curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/addLibraryFile`, define.Config.WebService[`port`])).
		Header(`Token`, token).Header(`lang`, flow.params.Lang).
		Param(`library_id`, libraryInfo[`id`]).
		Param(`group_id`, n.params.LibraryGroupId).
		PostFile(`library_files`, fileName).ToJSON(&res)
	if err != nil {
		logs.Error(`Create document request failed: %s %v\n`, res.Msg, err)
		msg = i18n.Show(flow.params.Lang, `create_doc_request_failed`)
		return
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		logs.Error(fmt.Sprintf(`Create document failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(flow.params.Lang, `create_doc_failed`)
		return
	} else {
		fileIds := cast.ToSlice(cast.ToStringMap(res.Data)[`file_ids`])
		if len(fileIds) == 0 || cast.ToInt64(fileIds[0]) <= 0 {
			logs.Error(fmt.Sprintf(`Create default custom document failed:%s`, tool.JsonEncodeNoError(res)))
			msg = i18n.Show(flow.params.Lang, `create_default_custom_doc_failed`)
			return
		} else {
			ok = true
		}
	}
	return
}

func (n *LibraryImport) toImportQa(lang string, token, question, answer string, images []any, similars []any, qaInfo msql.Params) (msg string, ok bool) {
	var res lib_web.Response
	req := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/editParagraph`, define.Config.WebService[`port`])).
		Header(`Token`, token).Header(`lang`, lang).
		Param(`question`, question).
		Param(`answer`, answer).
		Param(`library_id`, n.params.LibraryId)
	if len(images) > 0 {
		for _, item := range images {
			req.Param(`images`, cast.ToString(item))
		}
	}
	if len(qaInfo) > 0 {
		req.Param(`id`, qaInfo[`id`])
	}
	if len(similars) > 0 {
		req.Param(`similar_questions`, cast.ToString(tool.JsonEncodeNoError(similars)))
	} else {
		req.Param(`similar_questions`, `[]`)
	}
	err := req.Param(`group_id`, n.params.LibraryGroupId).ToJSON(&res)
	if err != nil {
		logs.Error(fmt.Sprintf(`Import QA pair request failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(lang, `import_qa_pair_request_failed`)
		return
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		logs.Error(fmt.Sprintf(`Create default custom document failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(lang, `create_default_custom_doc_failed`)
	} else {
		ok = true
	}
	return
}

func (n *LibraryImport) toUpdateUrl(lang string, token string, urlInfo msql.Params) (msg string, ok bool) {
	var res lib_web.Response
	err := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/manualCrawl`, define.Config.WebService[`port`])).
		Header(`Token`, token).Header(`lang`, lang).
		Param(`id`, urlInfo[`id`]).
		ToJSON(&res)
	if err != nil {
		msg = i18n.Show(lang, `update_url_failed`) + fmt.Sprintf(":%s(%v)", res.Msg, err)
		return
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		logs.Error(fmt.Sprintf(`Update failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(lang, `update_url_failed`)
	} else {
		msg = i18n.Show(lang, `update_url_success`)
		ok = true
	}
	return
}

func (n *LibraryImport) toImportUrl(lang string, token, normalUrl string) (msg string, ok bool) {
	var res lib_web.Response
	err := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/addLibraryFile`, define.Config.WebService[`port`])).
		Header(`Token`, token).Header(`lang`, lang).
		Param(`urls`, tool.JsonEncodeNoError([]map[string]string{{`url`: normalUrl}})).
		Param(`doc_auto_renew_frequency`, `1`).
		Param(`doc_auto_renew_minute`, `120`).
		Param(`doc_type`, `2`).
		Param(`library_id`, n.params.LibraryId).
		Param(`group_id`, n.params.LibraryGroupId).ToJSON(&res)
	if err != nil {
		logs.Error(fmt.Sprintf(`Import request failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(lang, `import_request_failed`)
		return
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		logs.Error(fmt.Sprintf(`Import failed:%s(%v)`, res.Msg, err))
		msg = i18n.Show(lang, `import_failed`)
	} else {
		fileIds := cast.ToSlice(cast.ToStringMap(res.Data)[`file_ids`])
		if len(fileIds) == 0 || cast.ToInt64(fileIds[0]) <= 0 {
			logs.Error(fmt.Sprintf(`Import failed:%s`, tool.JsonEncodeNoError(res)))
			msg = i18n.Show(lang, `import_not_success`)
		} else {
			ok = true
		}
	}
	return
}

func (n *LibraryImport) Params() any {
	return n.params
}

// WorkflowNode Workflow as a node
type WorkflowNode struct {
	params      WorkflowNodeParams
	nextNodeKey string
}

func (n *WorkflowNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, _ error) {
	nextNodeKey = n.nextNodeKey
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Where(`id`, cast.ToString(n.params.RobotId)).
		Find()
	if err != nil {
		flow.Logs(err.Error())
		nextNodeKey = n.params.Exception
		return
	}
	if len(robotInfo) == 0 {
		flow.Logs(i18n.Show(flow.params.Lang, `robot_not_exist`))
		nextNodeKey = n.params.Exception
		return
	}

	// Execute the next workflow
	subParams := &WorkFlowParams{
		ChatRequestParam:  flow.params.ChatRequestParam,
		RealRobot:         robotInfo,
		IsTestLoopNodeRun: false,
		TriggerParams:     TriggerParams{},
		IsFromWorkflow:    true,
	}
	subParams.WorkFlowGlobal = make(map[string]any)
	subParams.Robot = robotInfo
	subParams.ImmediatelyReplyHandle = nil // Do not continue nesting (Immediate Reply Node)

	// Variable replacement
	for _, p := range n.params.Params {
		subParams.WorkFlowGlobal[p.Key] = flow.VariableReplace(p.Variable)
	}
	subWorkflow, _, err := BaseCallWorkFlow(subParams)
	if err != nil {
		flow.Logs(err.Error())
		nextNodeKey = n.params.Exception
		return
	}

	// Extract specified reply content
	var replyArr []string
	var result string

	for key, out := range subWorkflow.outputs {
		if key == subWorkflow.curNodeKey { // Special handling for end node
			for specialKey, specialVal := range out {
				if strings.HasPrefix(specialKey, "special.finish_reply_") {
					for _, val := range specialVal.Vals {
						if len(*val.String) > 0 {
							replyArr = append(replyArr, *val.String)
						}
					}
				}
			}
		} else {
			if contentList, ok := out["special.llm_reply_content"]; ok {
				for _, val := range contentList.Vals {
					if len(*val.String) > 0 {
						replyArr = append(replyArr, *val.String)
					}
				}
			}
		}
	}
	result = strings.Join(replyArr, "\n")

	output = common.SimpleFields{}
	output[`data`] = common.SimpleField{
		Key:  `data`,
		Typ:  common.TypString,
		Vals: []common.Val{{String: &result}},
	}
	return
}

type ImmediatelyReplyNode struct {
	params      ImmediatelyReplyNodeParams
	nextNodeKey string
}

func (n *WorkflowNode) Params() any {
	return n.params
}

func (n *ImmediatelyReplyNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`Executing immediate reply logic...`)
	content := flow.VariableReplace(n.params.Content)
	// Immediately output message
	if flow.params.ImmediatelyReplyHandle != nil {
		flow.params.ImmediatelyReplyHandle(ImmediatelyReplyBuildReplyContent(content))
	}
	// Return output variables
	output = common.SimpleFields{
		`special.llm_reply_content`: common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &content}}},
	}
	nextNodeKey = n.nextNodeKey
	return
}

func (n *ImmediatelyReplyNode) Params() any {
	return n.params
}

type QuestionNode struct {
	params      QuestionParams
	nextNodeKey string
}

func (n *QuestionNode) Running(flow *WorkFlow) (output common.SimpleFields, nextNodeKey string, err error) {
	flow.Logs(`run question ...`)
	output = common.SimpleFields{}
	nextNodeKey = n.nextNodeKey
	n.params.AnswerText = flow.VariableReplace(n.params.AnswerText)
	if n.params.AnswerType == define.QuestionAnswerTypeText {
		output = common.SimpleFields{
			`special.llm_reply_content`: common.SimpleField{Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &n.params.AnswerText}}},
		}
	} else {
		replyContentList := n.params.ReplyContentList
		for replyContentKey, replyContent := range replyContentList {
			if replyContent.ReplyType != common.ReplyTypeSmartMenu {
				continue
			}
			replyContentList[replyContentKey].Type = replyContent.ReplyType
			replyContentList[replyContentKey].SmartMenu.MenuDescription = flow.VariableReplace(n.params.AnswerText)
			if len(replyContent.SmartMenu.MenuContent) > 0 {
				newMenuContents := []SmartMenuContent{}
				for _, menuContent := range replyContent.SmartMenu.MenuContent {
					menuContent.Content = flow.VariableReplace(menuContent.Content)
					if menuContent.MenuType != `-1` && menuContent.Content != `` {
						newMenuContents = append(newMenuContents, menuContent)
					}
				}
				replyContentList[replyContentKey].SmartMenu.MenuContent = newMenuContents
			}
		}
		output = common.SimpleFields{
			`special.reply_content_list`: common.SimpleField{Key: `special.reply_content_list`, Typ: common.TypString, Vals: []common.Val{{String: tea.String(tool.JsonEncodeNoError(replyContentList))}}},
		}
	}
	flow.isFinish = true
	flow.isStorage = true
	return
}

func (n *QuestionNode) Params() any {
	return n.params
}
