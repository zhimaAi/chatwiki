// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type WorkFlowLoop struct {
	Flow           *WorkFlow
	FlowLs         []*WorkFlow
	Outputs        common.SimpleFields //循环节点本身的输出
	LoopNodeParams *LoopNodeParams
	LoopNode       msql.Params
}

// FlowRunningLoopTest 循环节点单独运行测试
func FlowRunningLoopTest(flow *WorkFlow) error {
	_, childNode, err := GetNodeByKey(flow, cast.ToUint(flow.params.RealRobot[`id`]), flow.curNodeKey)
	if err != nil {
		return err
	}
	workFlowLoop, err := NewWorkFlowLoop(childNode[`loop_parent_key`], nil, flow)
	if err != nil {
		return err
	}
	_, err = workFlowLoop.ForRunning()
	if err != nil {
		return err
	}
	if len(workFlowLoop.FlowLs) > 0 {
		//替换输出结果 map已经更新了
		//flow.outputs = workFlowLoop.FlowLs[0].outputs
		//合并输出
		for _, flowL := range workFlowLoop.FlowLs {
			flow.nodeLogs = append(flow.nodeLogs, flowL.nodeLogs...)
			flow.runLogs = append(flow.runLogs, flowL.runLogs...)
		}
		flow.isFinish = true
	} else {
		err = errors.New(`未执行循环节点`)
	}
	return err
}

// LoopNodeRunning 循环节点运行
func LoopNodeRunning(nodeInfo msql.Params, flow *WorkFlow) (common.SimpleFields, string, error) {
	workFlowLoop, err := NewWorkFlowLoop(``, nodeInfo, flow)
	if err != nil {
		return make(common.SimpleFields), ``, err
	}
	nextNodeKey, err := workFlowLoop.ForRunning()
	if err != nil {
		return common.SimpleFields{}, ``, err
	}
	return workFlowLoop.Outputs, nextNodeKey, nil
}

func NewWorkFlowLoop(nodeKey string, node msql.Params, flow *WorkFlow) (*WorkFlowLoop, error) {
	var err error
	flowLoop := &WorkFlowLoop{
		Flow:           flow,
		FlowLs:         make([]*WorkFlow, 0),
		Outputs:        map[string]common.SimpleField{},
		LoopNodeParams: &LoopNodeParams{},
		LoopNode:       make(msql.Params),
	}
	if len(node) > 0 {
		flowLoop.LoopNode = node
	} else if nodeKey != `` {
		_, flowLoop.LoopNode, err = GetNodeByKey(flowLoop.Flow, cast.ToUint(flowLoop.Flow.params.RealRobot[`id`]), nodeKey)
		if err != nil {
			return flowLoop, err
		}
	} else {
		err = errors.New(`循环节点配置不能为空`)
		return flowLoop, err
	}
	if cast.ToInt(flowLoop.LoopNode[`node_type`]) != NodeTypeLoop {
		err = errors.New(`错误的循环节点`)
		return flowLoop, err
	}
	nodeParams := &NodeParams{}
	err = tool.JsonDecode(flowLoop.LoopNode[`node_params`], &nodeParams)
	if err != nil {
		logs.Error(`节点参数解析错误: %s`, err)
		return flowLoop, err
	}
	flowLoop.LoopNodeParams = &nodeParams.Loop
	return flowLoop, err
}

type NodeOutputTake struct {
	LlmResult struct {
		CompletionToken int `json:"completion_token"`
		PromptToken     int `json:"prompt_token"`
	} `json:"llm_result"`
}

func (flowLoop *WorkFlowLoop) ForRunning() (nextNodeKey string, err error) {
	nextNodeKey = ``
	//运行测试参数注入
	flowLoop.testFillParams(flowLoop.Flow)
	//循环基准
	loopInFields := flowLoop.getLoopInFields()
	if len(loopInFields) == 0 {
		flowLoop.Flow.Logs(`未找到循环次数`)
		return nextNodeKey, errors.New(`未找到循环次数`)
	}
	//单独运行测试只允许一次
	if flowLoop.Flow.params.IsTestLoopNodeRun {
		loopInFields = []*common.SimpleField{loopInFields[0]}
	}
	flowLoop.Flow.Logs(`开始运行循环节点%d次..`, len(loopInFields))
	//注入中间变量
	flowLoop.Flow.Logs(`注入中间变量...`)
	flowLoop.initLoopIntermediate()
loopFor:
	for loopIndex, inField := range loopInFields {
		flowL := &WorkFlow{}
		if inField == nil {
			flowLoop.Flow.Logs(`进行循环节点执行..第%d次`, loopIndex)
		} else {
			flowLoop.Flow.Logs(`进行循环节点执行..第%d次..循环参数%s%s`, loopIndex, "\n", tool.JsonEncodeNoError(*inField))
		}
		var isLoopEnd bool
		isLoopEnd, flowL, err = flowLoop.loopNodeRunning(loopIndex, len(loopInFields), inField)
		flowLoop.FlowLs = append(flowLoop.FlowLs, flowL)
		//中断处理
		if flowLoop.Flow.isTimeout {
			break
		}
		if isLoopEnd {
			flowLoop.Flow.Logs(`终止循环节点执行`)
			break
		}
		select {
		case <-flowLoop.Flow.context.Done():
			break loopFor
		default: //执行下一次循环
		}
	}
	//注入循环节点输出
	flowLoop.TakeChildOutputs()
	nextNodeKey = flowLoop.LoopNode[`next_node_key`]
	return
}

// 获取循环参数数组
func (flowLoop *WorkFlowLoop) getLoopInFields() []*common.SimpleField {
	//从外部输出中找到循环数组变量
	if flowLoop.LoopNodeParams.LoopType == common.LoopTypeNumber {
		loopInFields := make([]*common.SimpleField, 0)
		for i := 0; i < flowLoop.LoopNodeParams.LoopNumber.Int(); i++ {
			loopInFields = append(loopInFields, nil)
		}
		return loopInFields
	}
	for _, loopArray := range flowLoop.LoopNodeParams.LoopArrays {
		if loopArray.NodeKey() == `global` { //全局变量
			for globalKey, globalVal := range flowLoop.Flow.global {
				if globalKey == loopArray.ChooseKey() {
					return flowLoop.inFieldAppend(loopArray.Key, globalVal)
				}
			}
		} else {
			for outNodeKey, nodeOutputs := range flowLoop.Flow.outputs {
				if outNodeKey != loopArray.NodeKey() { //非指定的循环数组输出 下一个
					continue
				}
				for _, outField := range nodeOutputs {
					if outField.Typ == loopArray.Typ && outField.Key == loopArray.ChooseKey() {
						return flowLoop.inFieldAppend(loopArray.Key, outField)
					}
				}
			}
		}
	}
	return make([]*common.SimpleField, 0)
}

func (flowLoop *WorkFlowLoop) inFieldAppend(key string, outField common.SimpleField) []*common.SimpleField {
	loopInFields := make([]*common.SimpleField, 0)
	for _, val := range outField.Vals {
		switch outField.Typ {
		case common.TypArrFloat:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypFloat,
				Vals: []common.Val{val},
			})
		case common.TypArrObject:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypObject,
				Vals: []common.Val{val},
			})
		case common.TypArrBoole:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypBoole,
				Vals: []common.Val{val},
			})
		case common.TypArrNumber:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypNumber,
				Vals: []common.Val{val},
			})
		case common.TypArrParams:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypParams,
				Vals: []common.Val{val},
			})
		case common.TypArrString:
			loopInFields = append(loopInFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypString,
				Vals: []common.Val{val},
			})
		}
	}
	return loopInFields
}

func (flowLoop *WorkFlowLoop) TakeChildOutputs() {
	childFlowLogs := make(map[string][]common.NodeLog, 0)
	//提取token 累加
	completionToken, promptToken := 0, 0
	//提取每个节点最后一次的输出 不是最后一轮
	childOutputs := make(map[string]common.SimpleFields)
	for index, flowL := range flowLoop.FlowLs {
		//消耗的token处理
		for _, nodeLog := range flowL.nodeLogs {
			nodeOutputTake := NodeOutputTake{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(nodeLog.Output), &nodeOutputTake)
			if err != nil {
				flowLoop.Flow.Logs(`异常:循环节点提取子节点输出参数错误 %s`, err.Error())
				continue
			}
			completionToken += cast.ToInt(nodeOutputTake.LlmResult.CompletionToken)
			promptToken += cast.ToInt(nodeOutputTake.LlmResult.PromptToken)
		}
		//输出处理
		childFlowLogs[`loop_logs.for_`+cast.ToString(index+1)] = flowL.nodeLogs
		for childNodeKey, childOutput := range flowL.outputs {
			childOutputs[childNodeKey] = childOutput
		}
	}
	//注入token到结果中
	flowLoop.Outputs[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &completionToken}}}
	flowLoop.Outputs[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &promptToken}}}
	loopNumber := len(flowLoop.FlowLs)
	flowLoop.Outputs[`loop_result.loop_number`] = common.SimpleField{Key: `loop_result.loop_number`, Typ: common.TypNumber, Vals: []common.Val{{Number: &loopNumber}}}
	//从中间变量提取输出
	for _, loopOutput := range flowLoop.LoopNodeParams.Output {
		bFind := false
		for _, intermediate := range flowLoop.LoopNodeParams.IntermediateParams {
			if loopOutput.Value == flowLoop.LoopNode[`node_key`]+`.`+intermediate.Key {
				flowLoop.Outputs[loopOutput.Key] = intermediate.SimpleField
				bFind = true
			}
		}
		if !bFind {
			flowLoop.Outputs[loopOutput.Key] = common.SimpleField{
				Key:  loopOutput.ChooseKey(),
				Typ:  loopOutput.Typ,
				Vals: []common.Val{},
			}
		}
	}
	//测试输出日志
	if define.IsDev {
		//for key, log := range childFlowLogs {
		//	flowLoop.Outputs[key] = common.SimpleField{Key: key, Typ: common.TypArrObject, Vals: []common.Val{{Object: log}}}
		//}
	}
}

func (flowLoop *WorkFlowLoop) loopNodeRunning(loopIndex, maxLoopNumber int, inField *common.SimpleField) (isLoopEnd bool, flowL *WorkFlow, err error) {
	//new flow
	isLoopEnd = false
	flowL = &WorkFlow{
		params:      flowLoop.Flow.params,      //inherit params
		nodeLogs:    make([]common.NodeLog, 0), //node logs
		StartTime:   tool.Time2Int(),           //start time
		context:     flowLoop.Flow.context,     //inherit context
		global:      flowLoop.Flow.global,      //继承自循环节点所属工作流
		outputs:     flowLoop.Flow.outputs,     //继承所有节点的输出 所有对outputs的修改都会反应到主工作流
		runNodeKeys: make([]string, 0),         //self run node keys
		runLogs:     make([]string, 0),         //self  logs
		VersionId:   flowLoop.Flow.VersionId,   //inherit version id
		LoopIntermediate: LoopIntermediate{
			LoopNodeKey: flowLoop.LoopNode[`node_key`],
			Params:      &flowLoop.LoopNodeParams.IntermediateParams,
		}, //注入中间变量 供变量赋值使用
	}
	//循环参数注入
	flowLoop.fillLoopInField(flowL, inField)
	//将中间变量注入到outputs
	flowLoop.fillLoopIntermediateToOutputs(flowL)
	//找到入口子节点
	flowL.curNodeKey, err = flowLoop.FindStartChildNodeKey()
	if err != nil {
		flowL.Logs(err.Error())
		return
	}
	if flowL.curNodeKey == `` {
		err = errors.New(`未找到循环节点的入口子节点`)
		return
	}
	flowL.Logs(`开始运行循环节点子节点 第%d/%d次`, loopIndex, maxLoopNumber)
	for {
		var nodeInfo msql.Params
		flowL.curNode, nodeInfo, err = GetNodeByKey(flowL, cast.ToUint(flowL.params.RealRobot[`id`]), flowL.curNodeKey)
		if err != nil {
			flowL.Logs(err.Error())
			break
		}
		if flowL.curNode == nil {
			break //退出
		}
		flowL.Logs(`循环节点当前运行子节点:%s %s`, flowL.curNodeKey, nodeInfo[`node_name`])
		flowL.runNodeKeys = append(flowL.runNodeKeys, flowL.curNodeKey)
		var nextNodeKey string

		//节点运行开始
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flowL.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeLoop {
			err = errors.New(fmt.Sprintf(`循环节点中不支持循环节点 %s`, nodeLog.NodeName))
			break
		} else if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeLoopEnd {
			isLoopEnd = true
			flowL.outputs[flowL.curNodeKey] = make(common.SimpleFields)
		} else {
			flowL.output, nextNodeKey, err = flowL.curNode.Running(flowL)
			flowL.outputs[flowL.curNodeKey] = flowL.output //记录每个节点输出的变量
		}
		//将中间变量注入到outputs
		flowLoop.fillLoopIntermediateToOutputs(flowL)
		//运行参数处理
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flowL.output))
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flowL.nodeLogs = append(flowL.nodeLogs, nodeLog)
		//节点运行结束
		flowLoop.Flow.Logs(`结果nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flowL.output) > 0 {
			flowLoop.Flow.Logs(`输出变量:%s`, tool.JsonEncodeNoError(flowL.output))
		}
		if isLoopEnd {
			break
		}
		if flowLoop.Flow.isTimeout || err != nil || len(nextNodeKey) == 0 {
			break //结束
		}
		if nextNodeKey == flowLoop.LoopNode[`node_key`] {
			flowL.Logs(`循环节点本轮结束`)
			break
		}
		//外部中断监听处理
		select {
		case <-flowL.context.Done():
			goto flowExit
		default: //执行下一个节点
			flowL.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flowLoop.Flow.Logs(`循环节点第%d/%d次本次执行结束...`, loopIndex, maxLoopNumber)
	flowL.EndTime = tool.Time2Int()
	flowL.running = false //运行结束
	return
}

type TestFillVal struct {
	NodeKey  string `json:"node_key"`
	NodeName string `json:"node_name"`
	Field    struct {
		Sys      bool   `json:"sys"`
		Key      string `json:"key"`
		Typ      string `json:"typ"`
		Required bool   `json:"required"`
		Vals     any    `json:"Vals"`
	} `json:"field"`
}

func (flowLoop *WorkFlowLoop) testFillParams(flowL *WorkFlow) {
	if !flowLoop.Flow.params.IsTestLoopNodeRun {
		return
	}
	fillTestParamsToRunningParams(flowL, flowL.params.LoopTestParams)
}

// 注入循环参数
func (flowLoop *WorkFlowLoop) fillLoopInField(flowL *WorkFlow, inField *common.SimpleField) {
	if inField == nil {
		return
	}
	if _, ok := flowL.outputs[flowLoop.LoopNode[`node_key`]]; !ok {
		flowL.outputs[flowLoop.LoopNode[`node_key`]] = make(common.SimpleFields)
	}
	flowL.outputs[flowLoop.LoopNode[`node_key`]][inField.Key] = *inField
}

// 注入中间变量到outputs 供子节点使用
func (flowLoop *WorkFlowLoop) fillLoopIntermediateToOutputs(flowL *WorkFlow) {
	if _, ok := flowL.outputs[flowLoop.LoopNode[`node_key`]]; !ok {
		flowL.outputs[flowLoop.LoopNode[`node_key`]] = make(common.SimpleFields)
	}
	for _, loopIntermediate := range flowLoop.LoopNodeParams.IntermediateParams {
		flowL.outputs[flowLoop.LoopNode[`node_key`]][loopIntermediate.Key] = loopIntermediate.SimpleField
	}
}

// 初始化中间变量
func (flowLoop *WorkFlowLoop) initLoopIntermediate() {
	for key, intermediateParam := range flowLoop.LoopNodeParams.IntermediateParams {
		flowLoop.LoopNodeParams.IntermediateParams[key].Value = flowLoop.Flow.VariableReplace(intermediateParam.Value)

		//变量替换 支持数组等
		var data any = flowLoop.Flow.VariableReplace(intermediateParam.Value)
		if tool.InArrayString(intermediateParam.Typ, common.TypArrays[:]) {
			var temp []any //数组类型特殊处理
			for _, item := range strings.Split(cast.ToString(data), `、`) {
				temp = append(temp, item)
			}
			data = temp
		}
		flowLoop.LoopNodeParams.IntermediateParams[key].SimpleField = intermediateParam.SimpleField.SetVals(data)
	}
}

func (flowLoop *WorkFlowLoop) FindStartChildNodeKey() (string, error) {
	return msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, flowLoop.LoopNode[`admin_user_id`]).
		Where(`robot_id`, flowLoop.LoopNode[`robot_id`]).
		Where(`loop_parent_key`, flowLoop.LoopNode[`node_key`]).
		Where(`node_type`, cast.ToString(NodeTypeLoopStart)).
		Value(`node_key`)
}
