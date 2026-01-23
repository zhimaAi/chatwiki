// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
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
	//由前端传入,是否启用多模态输入
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
	TriggerParams      TriggerParams //触发器参数
	IsTestBatchNodeRun bool
	IsFromWorkflow     bool //是否来自工作流
	//立即回复输出句柄
	ImmediatelyReplyHandle func(replyContent common.ReplyContent)
}

type LoopIntermediate struct {
	LoopNodeKey string
	Params      *[]common.LoopField
}

type TriggerParams struct {
	TriggerType    uint                 //触发方式,默认值:TriggerTypeChat
	TestParams     map[string]any       //测试触发器传入参数
	TriggerOutputs []TriggerOutputParam //自己构建的trigger outputs
}

type WorkFlow struct {
	params           *WorkFlowParams
	nodeLogs         []common.NodeLog
	StartTime        int
	EndTime          int
	context          context.Context
	cancel           context.CancelFunc
	ticker           *time.Ticker //流程超时
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
	LoopIntermediate LoopIntermediate //中间变量 目前只用于循环节点，在变量赋值中使用
}

func (flow *WorkFlow) Logs(format string, a ...any) {
	msg := fmt.Sprintf(`[%s] %s`, tool.Date(), fmt.Sprintf(format, a...))
	flow.runLogs = append(flow.runLogs, msg)
	if define.IsDev {
		logs.Debug(fmt.Sprintf(`【%v】`, flow.global[`openid`].GetVal(common.TypString))+format, a...) //debug日志
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
	msg := fmt.Sprintf(`[%s] llm调用:%s`, tool.Date(), tool.JsonEncodeNoError(info))
	flow.runLogs = append(flow.runLogs, msg)
	if define.IsDev {
		jsonStr, _ := tool.JsonEncodeIndent(info, ``, "\t")
		logs.Debug(fmt.Sprintf("【%v】llm调用:\r\n%s", flow.global[`openid`].GetVal(common.TypString), jsonStr))
	}
}

func (flow *WorkFlow) Running() (err error) {
	flow.running = true
	flow.Logs(`进行工作流...`)
	if flow.params.Draft.IsDraft {
		flow.Logs(`使用草稿调试...`)
	} else {
		flow.VersionId = flow.getLastVersionId()
	}
	for {
		var nodeInfo msql.Params
		flow.Logs(`当前运行节点:%s`, flow.curNodeKey)
		flow.curNode, nodeInfo, err = GetNodeByKey(flow, cast.ToUint(flow.params.RealRobot[`id`]), flow.curNodeKey)
		if err != nil {
			flow.Logs(err.Error())
		}
		if flow.curNode == nil {
			break //退出
		}
		flow.runNodeKeys = append(flow.runNodeKeys, flow.curNodeKey)
		flow.getNodeInputs()
		flow.inputs[flow.curNodeKey] = flow.input
		var nextNodeKey string
		//节点运行开始
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
		flow.outputs[flow.curNodeKey] = flow.output //记录每个节点输出的变量
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flow.output))
		nodeLog.Input = common.GetFieldsObject(common.GetRecurveFields(flow.input))
		nodeLog.NodeOutput = GetNodeOutput(nodeLog.Output)
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flow.nodeLogs = append(flow.nodeLogs, nodeLog)
		//节点运行结束
		flow.Logs(`结果nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flow.output) > 0 {
			flow.Logs(`输出变量:%s`, tool.JsonEncodeNoError(flow.output))
		}
		if flow.isTimeout || err != nil || len(nextNodeKey) == 0 {
			break //结束
		}
		//外部中断监听处理
		select {
		case <-flow.context.Done():
			goto flowExit //特别注意:break不能做到退出for!!!
		default: //执行下一个节点
			flow.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flow.Logs(`工作流结束...`)
	flow.cancel() //关闭上下文
	flow.EndTime = tool.Time2Int()
	flow.running = false //运行结束
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
	flow.Logs(`保存工作流运行日志...`)
	_, err := msql.Model(`work_flow_logs`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: flow.params.AdminUserId,
		`robot_id`:      flow.params.RealRobot[`id`],
		`openid`:        cast.ToString(flow.global[`openid`].GetVal(common.TypString)),
		`run_node_keys`: strings.Join(flow.runNodeKeys, `,`),
		`run_logs`:      tool.JsonEncodeNoError(flow.runLogs),
		`create_time`:   flow.StartTime, //这里放开始时间
		`update_time`:   flow.EndTime,   //这里放结束时间
		`node_logs`:     tool.JsonEncodeNoError(flow.nodeLogs),
		`version_id`:    flow.VersionId,
		`question`:      cast.ToString(flow.global[`question`].GetVal(common.TypString)),
	})
	if err != nil {
		logs.Error(err.Error())
		flow.Logs(`保存工作流日志失败:%s`, err.Error())
	}
}

func (flow *WorkFlow) VariableReplace(content string) string {
	//优先替换全局变量
	for key, field := range flow.global {
		content = strings.ReplaceAll(content, fmt.Sprintf(`【global.%s】`, key), field.ShowVals())
	}
	//再替换节点输出变量
	for nodeKey, output := range flow.outputs {
		for key, field := range output {
			content = strings.ReplaceAll(content, fmt.Sprintf(`【%s.%s】`, nodeKey, key), field.ShowVals())
		}
	}
	//这个变成了旧数据兼容
	for key, field := range flow.output {
		content = strings.ReplaceAll(content, fmt.Sprintf(`【%s】`, key), field.ShowVals())
	}
	return regexp.MustCompile(`【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】`).ReplaceAllString(content, ``)
}

func (flow *WorkFlow) VariableReplaceJson(jsonStr string) string {
	//优先替换全局变量
	for key, field := range flow.global {
		jsonStr = strings.ReplaceAll(jsonStr, fmt.Sprintf(`【global.%s】`, key), strings.ReplaceAll(field.ShowVals(), `"`, `\"`))
	}
	//再替换节点输出变量
	for nodeKey, output := range flow.outputs {
		for key, field := range output {
			jsonStr = strings.ReplaceAll(jsonStr, fmt.Sprintf(`【%s.%s】`, nodeKey, key), strings.ReplaceAll(field.ShowVals(), `"`, `\"`))
		}
	}
	//这个变成了旧数据兼容
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

func SysGlobalVariables() []string { //固定值,runtime时不可变更
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
		global:      common.SimpleFields{},                //没有系统全局变量了
		outputs:     make(map[string]common.SimpleFields), //记录每个节点输出的变量
		inputs:      make(map[string]common.SimpleFields), //输入参数
		curNodeKey:  startNodeKey,                         //开始节点
		runNodeKeys: make([]string, 0),
		runLogs:     make([]string, 0),
	}
	go func(flow *WorkFlow) {
		defer flow.ticker.Stop()
		select {
		case <-flow.context.Done():
			return //流程已结束
		case <-flow.ticker.C:
		}
		flow.Logs(`工作流执行超时...`)
		flow.isTimeout = true
		flow.cancel()
	}(flow)
	var err error
	//循环节点单独运行测试时变量注入
	if flow.params.IsTestLoopNodeRun {
		err = FlowRunningLoopTest(flow)
	} else if flow.params.IsTestBatchNodeRun {
		err = FlowRunningBatchTest(flow)
	} else {
		err = flow.Running() //运行流程
	}
	if err == nil { //额外的校验逻辑
		if flow.isTimeout {
			err = errors.New(`工作流执行超时`)
		} else if !flow.isFinish {
			err = errors.New(`工作流未到结束节点`)
		}
	}
	go flow.Ending() //记录runtime日志
	return flow, err //返回数据
}

func BaseCallWorkFlow(params *WorkFlowParams) (flow *WorkFlow, nodeLogs []common.NodeLog, err error) {
	if len(params.RealRobot) == 0 { //未传参时,使用params里面的
		params.RealRobot = params.Robot
	}
	var startNodeKey string
	if params.Draft.IsDraft {
		startNodeKey = params.Draft.StartNodeKey
		if len(startNodeKey) == 0 {
			err = errors.New(`没有开始节点数据`)
			return
		}
	} else {
		startNodeKey = params.RealRobot[`start_node_key`]
		if len(startNodeKey) == 0 {
			err = errors.New(`工作流未发布`)
			return
		}
	}
	if params.TriggerParams.TriggerType == 0 {
		params.TriggerParams.TriggerType = TriggerTypeChat //默认会话触发
	}
	flow, err = RunningWorkFlow(params, startNodeKey)
	if flow != nil {
		nodeLogs = flow.nodeLogs
	}
	return
}

func CallWorkFlow(params *WorkFlowParams, debugLog *[]any, monitor *common.Monitor, isSwitchManual *bool) (content string, requestTime int64, libUseTime common.LibUseTime, list []msql.Params, err error) {
	flow, _, err := BaseCallWorkFlow(params)
	if flow != nil && len(flow.nodeLogs) > 0 {
		monitor.NodeLogs = flow.nodeLogs //记录监控数据
	}
	if flow == nil || err != nil {
		return
	}

	content = TakeOutputReply(flow)
	if len(content) == 0 {
		err = errors.New(`工作流没有可以回复的内容返回`)
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

func BuildFunctionTools(robot msql.Params) ([]adaptor.FunctionTool, bool) {
	if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeChat || len(robot[`work_flow_ids`]) == 0 {
		return nil, false
	}
	//判断func call能力
	if err := common.CheckSupportFuncCall(cast.ToInt(robot[`admin_user_id`]), cast.ToInt(robot[`model_config_id`]), robot[`use_model`]); err != nil {
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

func TakeOutputReply(flow *WorkFlow) string {
	if len(flow.output) <= 0 {
		return ``
	}
	//take finish node messages
	outputSlice := make([]string, 1000)
	for outputKey, _ := range flow.output {
		if !strings.HasPrefix(outputKey, define.FinishReplyPrefixKey) {
			continue
		}
		params := strings.Split(strings.TrimPrefix(outputKey, define.FinishReplyPrefixKey), `_`)
		if len(params) != 2 {
			logs.Warning(`finish node output message format error：%s`, outputKey)
			continue
		}
		//定义的第几个输出
		messageIdx := cast.ToInt(params[1]) - 1
		takeContent := cast.ToString(flow.output[outputKey].GetVal(common.TypString))
		if len(takeContent) == 0 {
			logs.Warning(`finish node output message is empty：%s，%s`, outputKey, takeContent)
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
		return strings.Join(contents, "\n")
	}
	//take any one
	var replyContentNodes = []string{`special.llm_reply_content`, `special.question_optimize_reply_content`, `special.mcp_reply_content`}
	content := ``
	for _, fieldsKey := range replyContentNodes {
		content = cast.ToString(flow.output[fieldsKey].GetVal(common.TypString))
		if len(content) > 0 {
			return content
		}
	}
	return ``
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
			nodeName = `开始节点`
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
