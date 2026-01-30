// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type Draft struct {
	IsDraft      bool
	NodeMaps     map[string]msql.Params
	StartNodeKey string
	//Passed from the frontend, whether to enable multimodal input
	QuestionMultipleSwitch bool
}

type WorkFlowParams struct {
	*define.ChatRequestParam
	RealRobot          msql.Params
	CurMsgId           int
	DialogueId         int
	SessionId          int
	Draft              Draft
	IsTestLoopNodeRun  bool
	TriggerParams      TriggerParams //Trigger parameters
	IsTestBatchNodeRun bool
	IsFromWorkflow     bool //Whether it comes from workflow
	//Immediate reply output handler
	ImmediatelyReplyHandle func(replyContent common.ReplyContent)
}

type LoopIntermediate struct {
	LoopNodeKey string
	Params      *[]common.LoopField
}

type TriggerParams struct {
	TriggerType    uint                 //Trigger type, default: TriggerTypeChat
	TestParams     map[string]any       //Test params passed to trigger
	TriggerOutputs []TriggerOutputParam //Manually built trigger outputs
}

type WorkFlow struct {
	params           *WorkFlowParams
	nodeLogs         []common.NodeLog
	StartTime        int
	EndTime          int
	context          context.Context
	cancel           context.CancelFunc
	ticker           *time.Ticker //Workflow timeout
	isTimeout        bool
	global           common.SimpleFields
	output           common.SimpleFields
	input            common.SimpleFields
	outputs          map[string]common.SimpleFields
	inputs           map[string]common.SimpleFields
	curNodeKey       string
	runNodeKeys      []string
	curNode          NodeAdapter
	runLogs          []string
	running          bool
	isFinish         bool
	VersionId        int
	LoopIntermediate LoopIntermediate //Intermediate variable; currently only for loop nodes, used in variable assignment
	isStorage        bool
}

func (flow *WorkFlow) Logs(format string, a ...any) {
	msg := fmt.Sprintf(`[%s] %s`, tool.Date(), fmt.Sprintf(format, a...))
	flow.runLogs = append(flow.runLogs, msg)
	if define.IsDev {
		logs.Debug(fmt.Sprintf(`【%v】`, flow.global[`openid`].GetVal(common.TypString))+format, a...) //debug log
	}
}

type LlmCallInfo struct {
	Params      LlmBaseParams                        `json:"params"`
	Messages    []adaptor.ZhimaChatCompletionMessage `json:"messages"`
	ChatResp    adaptor.ZhimaChatCompletionResponse  `json:"chat_resp"`
	RequestTime int64                                `json:"request_time"`
	Error       error                                `json:"error"`
}

func (flow *WorkFlow) LlmCallLogs(info LlmCallInfo) {
	msg := fmt.Sprintf(`[%s] llm call:%s`, tool.Date(), tool.JsonEncodeNoError(info))
	flow.runLogs = append(flow.runLogs, msg)
	if define.IsDev {
		jsonStr, _ := tool.JsonEncodeIndent(info, ``, "\t")
		logs.Debug(fmt.Sprintf("【%v】llm call:\r\n%s", flow.global[`openid`].GetVal(common.TypString), jsonStr))
	}
}

func (flow *WorkFlow) Running() (err error) {
	flow.running = true
	flow.Logs(`Running workflow...`)
	if flow.params.Draft.IsDraft {
		flow.Logs(`Debugging with draft...`)
	} else {
		flow.VersionId = flow.getLastVersionId()
	}
	for {
		var nodeInfo msql.Params
		flow.Logs(`Current running node:%s`, flow.curNodeKey)
		flow.curNode, nodeInfo, err = GetNodeByKey(flow, cast.ToUint(flow.params.RealRobot[`id`]), flow.curNodeKey)
		if err != nil {
			flow.Logs(err.Error())
		}
		if flow.curNode == nil {
			break //exit
		}
		flow.runNodeKeys = append(flow.runNodeKeys, flow.curNodeKey)
		flow.getNodeInputs()
		flow.inputs[flow.curNodeKey] = flow.input
		var nextNodeKey string
		//node run start
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flow.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeLoop {
			flow.output, nextNodeKey, err = LoopNodeRunning(nodeInfo, flow)
		} else if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeBatch {
			flow.output, nextNodeKey, err = BatchNodeRunning(nodeInfo, flow)
		} else {
			flow.output, nextNodeKey, err = flow.curNode.Running(flow)
		}
		flow.outputs[flow.curNodeKey] = flow.output //record variables output by each node

		//add icon
		if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeHttpTool {
			nodeParams := flow.curNode.Params().(CurlNodeParams)
			nodeLog.NodeIcon = nodeParams.ToolInfo.HttpToolAvatar
		}

		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flow.output))
		nodeLog.Input = common.GetFieldsObject(common.GetRecurveFields(flow.input))
		nodeLog.NodeOutput = GetNodeOutput(nodeLog.Output)
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flow.nodeLogs = append(flow.nodeLogs, nodeLog)
		//node run end
		flow.Logs(`Result nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flow.output) > 0 {
			flow.Logs(`Output variables:%s`, tool.JsonEncodeNoError(flow.output))
		}
		if flow.isTimeout || err != nil || len(nextNodeKey) == 0 || flow.isStorage {
			break //end
		}
		//external interruption listener
		select {
		case <-flow.context.Done():
			goto flowExit //note: break cannot exit the for loop!!!
		default: //run next node
			flow.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flow.Logs(`Workflow finished...`)
	flow.cancel() //close context
	flow.EndTime = tool.Time2Int()
	flow.running = false //run finished
	return
}

func (flow *WorkFlow) getLastVersionId() int {
	versionId, err := msql.Model(`work_flow_version`, define.Postgres).
		Where(`robot_id`, cast.ToString(flow.params.Robot[`id`])).
		Where(`admin_user_id`, cast.ToString(flow.params.AdminUserId)).
		Order(`id desc`).Limit(1).Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		return 0
	}
	return cast.ToInt(versionId)
}

func (flow *WorkFlow) Ending() {
	flow.Logs(`Saving workflow run logs...`)
	_, err := msql.Model(`work_flow_logs`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: flow.params.AdminUserId,
		`robot_id`:      flow.params.RealRobot[`id`],
		`openid`:        cast.ToString(flow.global[`openid`].GetVal(common.TypString)),
		`run_node_keys`: strings.Join(flow.runNodeKeys, `,`),
		`run_logs`:      tool.JsonEncodeNoError(flow.runLogs),
		`create_time`:   flow.StartTime, //store start time here
		`update_time`:   flow.EndTime,   //store end time here
		`node_logs`:     tool.JsonEncodeNoError(flow.nodeLogs),
		`version_id`:    flow.VersionId,
		`question`:      cast.ToString(flow.global[`question`].GetVal(common.TypString)),
	})
	if err != nil {
		logs.Error(err.Error())
		flow.Logs(`Failed to save workflow logs:%s`, err.Error())
	}
}

func (flow *WorkFlow) VariableReplace(content string) string {
	//Replace global variables first
	for key, field := range flow.global {
		content = strings.ReplaceAll(content, fmt.Sprintf(`【global.%s】`, key), field.ShowVals())
	}
	//Then replace node output variables
	for nodeKey, output := range flow.outputs {
		for key, field := range output {
			content = strings.ReplaceAll(content, fmt.Sprintf(`【%s.%s】`, nodeKey, key), field.ShowVals())
		}
	}
	//Legacy compatibility for old data
	for key, field := range flow.output {
		content = strings.ReplaceAll(content, fmt.Sprintf(`【%s】`, key), field.ShowVals())
	}
	return regexp.MustCompile(`【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】`).ReplaceAllString(content, ``)
}

func (flow *WorkFlow) VariableReplaceJson(jsonStr string) string {
	//Replace global variables first
	for key, field := range flow.global {
		jsonStr = strings.ReplaceAll(jsonStr, fmt.Sprintf(`【global.%s】`, key), strings.ReplaceAll(field.ShowVals(), `"`, `\"`))
	}
	//Then replace node output variables
	for nodeKey, output := range flow.outputs {
		for key, field := range output {
			jsonStr = strings.ReplaceAll(jsonStr, fmt.Sprintf(`【%s.%s】`, nodeKey, key), strings.ReplaceAll(field.ShowVals(), `"`, `\"`))
		}
	}
	//Legacy compatibility for old data
	for key, field := range flow.output {
		jsonStr = strings.ReplaceAll(jsonStr, fmt.Sprintf(`【%s】`, key), strings.ReplaceAll(field.ShowVals(), `"`, `\"`))
	}
	return regexp.MustCompile(`【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】`).ReplaceAllString(jsonStr, ``)
}

func (flow *WorkFlow) GetVariable(key string) (field common.SimpleField, exist bool) {
	if strings.HasPrefix(key, `global.`) {
		realKey := strings.TrimPrefix(key, `global.`)
		if field, exist = flow.global[realKey]; exist {
			return
		}
	}
	if temp := strings.SplitN(key, `.`, 2); len(temp) == 2 && common.IsMd5Str(temp[0]) {
		if output := flow.outputs[temp[0]]; len(output) > 0 {
			if field, exist = output[temp[1]]; exist {
				return
			}
		}
	}
	if field, exist = flow.output[key]; exist {
		return
	}
	return
}

func SysGlobalVariables() []string { //Fixed values, immutable at runtime
	return []string{}
}

func RunningWorkFlow(params *WorkFlowParams, startNodeKey string) (*WorkFlow, error) {
	ctx, cancel := context.WithCancel(context.Background())
	flow := &WorkFlow{
		params:      params,
		nodeLogs:    make([]common.NodeLog, 0),
		StartTime:   tool.Time2Int(),
		context:     ctx,
		cancel:      cancel,
		ticker:      time.NewTicker(time.Minute * 60),     //DIY
		global:      common.SimpleFields{},                //no system global variables anymore
		outputs:     make(map[string]common.SimpleFields), //record variables output by each node
		inputs:      make(map[string]common.SimpleFields), //input parameters
		curNodeKey:  startNodeKey,                         //start node
		runNodeKeys: make([]string, 0),
		runLogs:     make([]string, 0),
	}
	var err error
	err = WorkFlowRestore(flow)
	if err != nil {
		return flow, err
	}
	go func(flow *WorkFlow) {
		defer flow.ticker.Stop()
		select {
		case <-flow.context.Done():
			return //workflow finished
		case <-flow.ticker.C:
		}
		flow.Logs(`Workflow execution timeout...`)
		flow.isTimeout = true
		flow.cancel()
	}(flow)

	//Inject variables when testing loop node standalone
	if flow.params.IsTestLoopNodeRun {
		err = FlowRunningLoopTest(flow)
	} else if flow.params.IsTestBatchNodeRun {
		err = FlowRunningBatchTest(flow)
	} else {
		err = flow.Running() //run workflow
	}
	if err == nil { //additional validation logic
		if flow.isTimeout {
			err = errors.New(i18n.Show(flow.params.Lang, "workflow_execution_timeout"))
		} else if !flow.isFinish {
			err = errors.New(i18n.Show(flow.params.Lang, "workflow_not_reach_end_node"))
		}
	}
	if flow.isStorage {
		SetWorkFlowStorage(flow)
	} else {
		go flow.Ending() //record runtime logs
	}
	return flow, err //return data
}

func BaseCallWorkFlow(params *WorkFlowParams) (flow *WorkFlow, nodeLogs []common.NodeLog, err error) {
	if len(params.RealRobot) == 0 { //If not provided, use params.Robot
		params.RealRobot = params.Robot
	}
	var startNodeKey string
	if params.Draft.IsDraft {
		startNodeKey = params.Draft.StartNodeKey
		if len(startNodeKey) == 0 {
			err = errors.New(i18n.Show(flow.params.Lang, "workflow_no_start_node_data"))
			return
		}
	} else {
		startNodeKey = params.RealRobot[`start_node_key`]
		if len(startNodeKey) == 0 {
			err = errors.New(i18n.Show(flow.params.Lang, "workflow_no_start_node_data"))
			return
		}
	}
	if params.TriggerParams.TriggerType == 0 {
		params.TriggerParams.TriggerType = TriggerTypeChat //default chat trigger
	}
	flow, err = RunningWorkFlow(params, startNodeKey)
	if flow != nil {
		nodeLogs = flow.nodeLogs
	}
	return
}

func CallWorkFlow(params *WorkFlowParams, debugLog *[]any, monitor *common.Monitor, isSwitchManual *bool) (content string, requestTime int64, libUseTime common.LibUseTime, list []msql.Params, replyContentList []common.ReplyContent, err error) {
	flow, _, err := BaseCallWorkFlow(params)
	if flow != nil && len(flow.nodeLogs) > 0 {
		monitor.NodeLogs = flow.nodeLogs //record monitoring data
	}
	if flow == nil || err != nil {
		return
	}

	content, replyContentList = TakeOutputReply(flow)
	if len(content) == 0 && len(replyContentList) == 0 {
		err = errors.New(i18n.Show(params.Lang, "workflow_no_reply_content"))
		return
	}
	requestTime = cast.ToInt64(flow.output[`special.llm_request_time`].GetVal(common.TypNumber))
	libUseTime, _ = flow.output[`special.lib_use_time`].GetVal().(common.LibUseTime)
	for _, val := range flow.output[`special.lib_paragraph_list`].Vals {
		list = append(list, val.Params)
	}
	*debugLog = append(*debugLog, flow.output[`special.llm_debug_log`].GetVals()...)
	if cast.ToBool(flow.output[`special.is_switch_manual`].GetVal()) {
		*isSwitchManual = true
	}
	return
}

func BuildFunctionTools(lang string, robot msql.Params) ([]adaptor.FunctionTool, bool) {
	if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeChat || len(robot[`work_flow_ids`]) == 0 {
		return nil, false
	}
	//check function call capability
	if err := common.CheckSupportFuncCall(lang, cast.ToInt(robot[`admin_user_id`]), cast.ToInt(robot[`model_config_id`]), robot[`use_model`]); err != nil {
		return nil, false
	}
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	list, err := m.Where(`application_type`, cast.ToString(define.ApplicationTypeFlow)).Where(`start_node_key`, `<>`, ``).
		Where(`admin_user_id`, robot[`admin_user_id`]).Where(`id`, `in`, robot[`work_flow_ids`]).
		Field(`id,robot_key,robot_name,robot_intro,start_node_key,en_name`).Select()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	functionTools := make([]adaptor.FunctionTool, 0)
	for _, item := range list {
		node, _, err := GetNodeByKey(nil, cast.ToUint(item[`id`]), item[`start_node_key`])
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		startNode, _ := node.(*StartNode)
		properties, required := make(map[string]map[string]string), make([]string, 0)
		if startNode != nil && len(startNode.params.DiyGlobal) > 0 {
			for _, param := range startNode.params.DiyGlobal {
				properties[param.Key] = map[string]string{`type`: param.Typ, `description`: param.Desc}
				if param.Required {
					required = append(required, param.Key)
				}
			}
		}
		name := fmt.Sprintf(`work_flow_#_%s`, item[`robot_key`])
		if len(item[`en_name`]) > 0 {
			name = item[`en_name`]
		}
		functionTools = append(functionTools, adaptor.FunctionTool{
			Name:        name,
			Description: fmt.Sprintf(`%s(%s)`, item[`robot_name`], item[`robot_intro`]),
			Parameters: adaptor.Parameters{
				Type:       `object`,
				Properties: properties,
				Required:   required,
			},
		})
	}
	return functionTools, len(functionTools) > 0
}

func ChooseWorkFlowRobot(adminUserId string, functionTools []adaptor.FunctionToolCall) (_ msql.Params, global map[string]any) {
	for _, functionTool := range functionTools {
		robotKey, ok := common.IsWorkFlowFuncCall(adminUserId, functionTool.Name)
		if !ok {
			continue
		}
		robot, err := common.GetRobotInfo(robotKey)
		if err != nil {
			logs.Error(err.Error())
		}
		if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeFlow || len(robot[`start_node_key`]) == 0 {
			continue
		}
		_ = tool.JsonDecodeUseNumber(functionTool.Arguments, &global)
		return robot, global
	}
	return
}

func BuildWorkFlowParams(params define.ChatRequestParam, workFlowRobot msql.Params, workFlowGlobal map[string]any, curMsgId, dialogueId, sessionId int) *WorkFlowParams {
	params.WorkFlowGlobal = workFlowGlobal
	return &WorkFlowParams{ChatRequestParam: &params, RealRobot: workFlowRobot, CurMsgId: curMsgId, DialogueId: dialogueId, SessionId: sessionId}
}

func TakeOutputReply(flow *WorkFlow) (string, []common.ReplyContent) {
	//add reply content list
	var replyContentList []common.ReplyContent
	if len(flow.output) <= 0 {
		return ``, replyContentList
	}
	replyContentListJson := cast.ToString(flow.output[`special.reply_content_list`].GetVal(common.TypString))
	if len(replyContentListJson) > 0 {
		err := tool.JsonDecode(replyContentListJson, &replyContentList)
		if err != nil {
			logs.Warning(`output reply content list json error: %s`, err.Error())
		}
	}
	//finish node to string
	outputSlice := make([]string, 1000)
	for outputKey, _ := range flow.output {
		if !strings.HasPrefix(outputKey, define.FinishReplyPrefixKey) {
			continue
		}
		params := strings.Split(strings.TrimPrefix(outputKey, define.FinishReplyPrefixKey), `_`)
		if len(params) != 2 {
			logs.Warning(`finish node output message format error: %s`, outputKey)
			continue
		}
		// Defined output index
		messageIdx := cast.ToInt(params[1]) - 1
		takeContent := cast.ToString(flow.output[outputKey].GetVal(common.TypString))
		if len(takeContent) == 0 {
			logs.Warning(`finish node output message is empty: %s, %s`, outputKey, takeContent)
			continue
		}
		if params[0] == lib_define.MsgTypeText {
			outputSlice[cast.ToInt(cast.ToString(messageIdx)+`0`)] = takeContent
		} else if params[0] == lib_define.MsgTypeImage {
			urls := common.TakeUrls(&takeContent)
			if len(urls) > 0 {
				for idx, url := range urls {
					outputSlice[cast.ToInt(cast.ToString(messageIdx)+cast.ToString(idx))] = fmt.Sprintf(`![%s](%s)`, url, url)
				}
			}
		} else if params[0] == lib_define.MsgTypeVoice {
			urls := common.TakeUrls(&takeContent)
			if len(urls) > 0 {
				for idx, url := range urls {
					ext := filepath.Ext(url)
					if strings.ToLower(ext) != `.mp3` {
						logs.Warning(`finish node output voice url is not mp3`)
						continue
					}
					outputSlice[cast.ToInt(cast.ToString(messageIdx)+cast.ToString(idx))] = fmt.Sprintf(`!voice[](%s)`, url)
				}
			}
		}
	}
	contents := make([]string, 0)
	for _, content := range outputSlice {
		if len(content) > 0 {
			contents = append(contents, content)
		}
	}
	if len(contents) > 0 {
		return strings.Join(contents, "\n"), replyContentList
	}
	//default reply
	var replyContentNodes = []string{`special.llm_reply_content`, `special.question_optimize_reply_content`, `special.mcp_reply_content`}
	content := ``
	for _, fieldsKey := range replyContentNodes {
		content = cast.ToString(flow.output[fieldsKey].GetVal(common.TypString))
		if len(content) > 0 {
			return content, replyContentList
		}
	}
	return ``, replyContentList
}

func CallHttpTest(workFlowParams *WorkFlowParams, curlNode *WorkFlowNode) (map[string]any, error) {
	flow := &WorkFlow{
		params:      workFlowParams,
		nodeLogs:    make([]common.NodeLog, 0),
		outputs:     make(map[string]common.SimpleFields),
		runNodeKeys: make([]string, 0),
		runLogs:     make([]string, 0),
		global:      map[string]common.SimpleField{},
	}
	curlNodeRun := CurlNode{
		params:      curlNode.NodeParams.Curl,
		nextNodeKey: "",
	}
	fillTestParamsToRunningParams(flow, workFlowParams.TestParams)
	output, _, err := curlNodeRun.Running(flow)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	httpResultJson := curlNodeRun.GetHttpResultJson()
	data := make(map[string]any)
	err = tool.JsonDecode(*httpResultJson, &data)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	recursiveFields := common.GetRecursiveFieldsFromMap(data)
	return map[string]any{
		`output`: output,
		`result`: data,
		`fields`: recursiveFields,
	}, nil
}

func (flow *WorkFlow) getNodeInputs() {
	if flow == nil {
		return
	}
	flow.input = common.SimpleFields{}
	if flow.curNode.Params() == nil {
		return
	}
	var variables []string
	ExtractVariables(flow.curNode.Params(), &variables)
	if len(variables) == 0 {
		return
	}
	var err error
	for _, variable := range variables {
		variable = strings.TrimPrefix(variable, `【`)
		variable = strings.TrimSuffix(variable, `】`)
		var nodeName string
		var variableKey string
		var field common.SimpleField
		var exist bool
		if strings.HasPrefix(variable, `global.`) {
			nodeName = i18n.Show(flow.params.Lang, `start_node`)
			variableKey = strings.TrimPrefix(variable, `global.`)
			field, exist = flow.global[variableKey]
			if !exist {
				continue
			}
		} else {
			params := strings.Split(variable, `.`)
			if len(params) == 0 {
				continue
			}
			nodeKey := params[0]
			if len(nodeKey) == 0 {
				continue
			}
			var info msql.Params
			if flow.params.Draft.IsDraft {
				info = flow.params.Draft.NodeMaps[nodeKey]
			} else {
				info, err = common.GetRobotNode(cast.ToUint(flow.params.Robot[`id`]), nodeKey)
				if err != nil {
					logs.Error(err.Error())
					continue
				}
			}
			if len(info) == 0 {
				continue
			}
			nodeName = cast.ToString(info[`node_name`])
			variableKey = strings.TrimPrefix(variable, nodeKey+`.`)
			output := flow.outputs[nodeKey]
			if len(output) == 0 {
				continue
			}
			field, exist = output[variableKey]
			if !exist {
				continue
			}
		}
		flow.input[nodeName+`/`+variableKey] = common.SimpleField{
			Key:  nodeName + `/` + variableKey,
			Typ:  field.Typ,
			Vals: field.Vals,
		}
	}
	return
}
