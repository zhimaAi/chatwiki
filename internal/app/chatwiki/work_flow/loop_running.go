// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
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
	Outputs        common.SimpleFields // Output of the loop node itself
	LoopNodeParams *LoopNodeParams
	LoopNode       msql.Params
}

// FlowRunningLoopTest Loop node individual run test
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
		// Replace output result map has been updated
		//flow.outputs = workFlowLoop.FlowLs[0].outputs
		// Merge output
		for _, flowL := range workFlowLoop.FlowLs {
			flow.nodeLogs = append(flow.nodeLogs, flowL.nodeLogs...)
			flow.runLogs = append(flow.runLogs, flowL.runLogs...)
		}
		flow.isFinish = true
	} else {
		err = errors.New(i18n.Show(flow.params.Lang, "workflow_node_loop_not_executed"))
	}
	return err
}

// LoopNodeRunning Loop node running
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
		err = errors.New(i18n.Show(flowLoop.Flow.params.Lang, "workflow_node_loop_config_empty"))
		return flowLoop, err
	}
	if cast.ToInt(flowLoop.LoopNode[`node_type`]) != NodeTypeLoop {
		err = errors.New(i18n.Show(flowLoop.Flow.params.Lang, "workflow_node_loop_error"))
		return flowLoop, err
	}
	nodeParams := &NodeParams{}
	err = tool.JsonDecode(flowLoop.LoopNode[`node_params`], &nodeParams)
	if err != nil {
		logs.Error(`Node parameter parsing error: %s`, err)
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
	// Run test parameter injection
	flowLoop.testFillParams(flowLoop.Flow)
	// Loop basis
	loopInFields := flowLoop.getLoopInFields()
	if len(loopInFields) == 0 {
		flowLoop.Flow.Logs(`Loop count not found`)
		return nextNodeKey, errors.New(i18n.Show(flowLoop.Flow.params.Lang, "workflow_node_loop_count_not_found"))
	}
	// Individual run test only allows once
	if flowLoop.Flow.params.IsTestLoopNodeRun {
		loopInFields = []*common.SimpleField{loopInFields[0]}
	}
	flowLoop.Flow.Logs(`Start running loop node %d times..`, len(loopInFields))
	// Inject intermediate variables
	flowLoop.Flow.Logs(`Injecting intermediate variables...`)
	flowLoop.initLoopIntermediate()
loopFor:
	for loopIndex, inField := range loopInFields {
		flowL := &WorkFlow{}
		if inField == nil {
			flowLoop.Flow.Logs(`Executing loop node.. %dth time`, loopIndex)
		} else {
			flowLoop.Flow.Logs(`Executing loop node.. %dth time.. loop params%s%s`, loopIndex, "\n", tool.JsonEncodeNoError(*inField))
		}
		var isLoopEnd bool
		isLoopEnd, flowL, err = flowLoop.loopNodeRunning(loopIndex, len(loopInFields), inField)
		flowLoop.FlowLs = append(flowLoop.FlowLs, flowL)
		// Interrupt handling
		if flowLoop.Flow.isTimeout {
			break
		}
		if isLoopEnd {
			flowLoop.Flow.Logs(`Terminate loop node execution`)
			break
		}
		select {
		case <-flowLoop.Flow.context.Done():
			break loopFor
		default: // Execute next loop
		}
	}
	// Inject loop node output
	flowLoop.TakeChildOutputs()
	nextNodeKey = flowLoop.LoopNode[`next_node_key`]
	return
}

// Get loop parameter array
func (flowLoop *WorkFlowLoop) getLoopInFields() []*common.SimpleField {
	// Find loop array variable from external output
	if flowLoop.LoopNodeParams.LoopType == common.LoopTypeNumber {
		loopInFields := make([]*common.SimpleField, 0)
		for i := 0; i < flowLoop.LoopNodeParams.LoopNumber.Int(); i++ {
			loopInFields = append(loopInFields, nil)
		}
		return loopInFields
	}
	for _, loopArray := range flowLoop.LoopNodeParams.LoopArrays {
		if loopArray.NodeKey() == `global` { // global variable
			for globalKey, globalVal := range flowLoop.Flow.global {
				if globalKey == loopArray.ChooseKey() {
					return flowLoop.inFieldAppend(loopArray.Key, globalVal)
				}
			}
		} else {
			for outNodeKey, nodeOutputs := range flowLoop.Flow.outputs {
				if outNodeKey != loopArray.NodeKey() { // Non-specified loop array output, next
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
	//Extract and accumulate tokens
	completionToken, promptToken := 0, 0
	//Extract the last output of each node (not the last loop round)
	childOutputs := make(map[string]common.SimpleFields)
	for index, flowL := range flowLoop.FlowLs {
		//Consumed token handling
		for _, nodeLog := range flowL.nodeLogs {
			nodeOutputTake := NodeOutputTake{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(nodeLog.Output), &nodeOutputTake)
			if err != nil {
				flowLoop.Flow.Logs(`Error: failed to extract child node output params in loop node %s`, err.Error())
				continue
			}
			completionToken += cast.ToInt(nodeOutputTake.LlmResult.CompletionToken)
			promptToken += cast.ToInt(nodeOutputTake.LlmResult.PromptToken)
		}
		//Output handling
		childFlowLogs[`loop_logs.for_`+cast.ToString(index+1)] = flowL.nodeLogs
		for childNodeKey, childOutput := range flowL.outputs {
			childOutputs[childNodeKey] = childOutput
		}
	}
	//Inject tokens into result
	flowLoop.Outputs[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &completionToken}}}
	flowLoop.Outputs[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &promptToken}}}
	loopNumber := len(flowLoop.FlowLs)
	flowLoop.Outputs[`loop_result.loop_number`] = common.SimpleField{Key: `loop_result.loop_number`, Typ: common.TypNumber, Vals: []common.Val{{Number: &loopNumber}}}
	//Extract output from intermediate variables
	for _, loopOutput := range flowLoop.LoopNodeParams.Output {
		bFind := false
		for _, intermediate := range flowLoop.LoopNodeParams.IntermediateParams {
			if loopOutput.Value == flowLoop.LoopNode[`node_key`]+`.`+intermediate.Key {
				flowLoop.Outputs[loopOutput.Key] = common.SimpleField{
					Key:  loopOutput.Key,
					Typ:  loopOutput.Typ,
					Vals: intermediate.SimpleField.Vals,
				}
				bFind = true
			}
		}
		if !bFind {
			flowLoop.Outputs[loopOutput.Key] = common.SimpleField{
				Key:  loopOutput.Key,
				Typ:  loopOutput.Typ,
				Vals: []common.Val{},
			}
		}
	}
	//Test output logs
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
		params:      flowLoop.Flow.params,                 //inherit params
		nodeLogs:    make([]common.NodeLog, 0),            //node logs
		StartTime:   tool.Time2Int(),                      //start time
		context:     flowLoop.Flow.context,                //inherit context
		global:      flowLoop.Flow.global,                 //inherited from the workflow owning this loop node
		outputs:     flowLoop.Flow.outputs,                //inherit outputs of all nodes; modifications reflect to parent flow
		inputs:      make(map[string]common.SimpleFields), //input parameters
		runNodeKeys: make([]string, 0),                    //self run node keys
		runLogs:     make([]string, 0),                    //self  logs
		VersionId:   flowLoop.Flow.VersionId,              //inherit version id
		LoopIntermediate: LoopIntermediate{
			LoopNodeKey: flowLoop.LoopNode[`node_key`],
			Params:      &flowLoop.LoopNodeParams.IntermediateParams,
		}, //inject intermediate variables for variable assignment
	}
	//Inject loop parameters
	flowLoop.fillLoopInField(flowL, inField)
	//Inject intermediate variables into outputs
	flowLoop.fillLoopIntermediateToOutputs(flowL)
	//Find the entry child node
	flowL.curNodeKey, err = flowLoop.FindStartChildNodeKey()
	if err != nil {
		flowL.Logs(err.Error())
		return
	}
	if flowL.curNodeKey == `` {
		err = errors.New(i18n.Show(flowL.params.Lang, "workflow_node_loop_entry_node_not_found"))
		return
	}
	flowL.Logs(`Start running loop node child nodes %d/%d`, loopIndex, maxLoopNumber)
	for {
		var nodeInfo msql.Params
		flowL.curNode, nodeInfo, err = GetNodeByKey(flowL, cast.ToUint(flowL.params.RealRobot[`id`]), flowL.curNodeKey)
		if err != nil {
			flowL.Logs(err.Error())
			break
		}
		if flowL.curNode == nil {
			break //exit
		}
		flowL.Logs(`Loop node currently running child node:%s %s`, flowL.curNodeKey, nodeInfo[`node_name`])
		flowL.runNodeKeys = append(flowL.runNodeKeys, flowL.curNodeKey)
		var nextNodeKey string

		//node run start
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flowL.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		flowL.getNodeInputs()
		flowL.inputs[flowL.curNodeKey] = flowL.input
		if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeLoop {
			err = errors.New(i18n.Show(flowL.params.Lang, "workflow_node_loop_not_support_nested_loop", nodeLog.NodeName))
			break
		} else if cast.ToInt(nodeInfo[`node_type`]) == NodeTypeLoopEnd {
			isLoopEnd = true
			flowL.outputs[flowL.curNodeKey] = make(common.SimpleFields)
		} else {
			flowL.output, nextNodeKey, err = flowL.curNode.Running(flowL)
			flowL.outputs[flowL.curNodeKey] = flowL.output // Record variables output by each node
		}
		// Inject intermediate variables into outputs
		flowLoop.fillLoopIntermediateToOutputs(flowL)
		// Running parameter processing
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flowL.output))
		nodeLog.Input = common.GetFieldsObject(common.GetRecurveFields(flowL.input))
		nodeLog.NodeOutput = GetNodeOutput(nodeLog.Output)
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flowL.nodeLogs = append(flowL.nodeLogs, nodeLog)
		// Node running end
		flowLoop.Flow.Logs(`Result nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flowL.output) > 0 {
			flowLoop.Flow.Logs(`Output variables:%s`, tool.JsonEncodeNoError(flowL.output))
		}
		if isLoopEnd {
			break
		}
		if flowLoop.Flow.isTimeout || err != nil || len(nextNodeKey) == 0 {
			break // End
		}
		if nextNodeKey == flowLoop.LoopNode[`node_key`] {
			flowL.Logs(`Loop node round end`)
			break
		}
		// External interrupt listening processing
		select {
		case <-flowL.context.Done():
			goto flowExit
		default: // Execute next node
			flowL.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flowLoop.Flow.Logs(`Loop node %d/%d execution ended...`, loopIndex, maxLoopNumber)
	flowL.EndTime = tool.Time2Int()
	flowL.running = false // Running end
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

// Inject loop parameters
func (flowLoop *WorkFlowLoop) fillLoopInField(flowL *WorkFlow, inField *common.SimpleField) {
	if inField == nil {
		return
	}
	if _, ok := flowL.outputs[flowLoop.LoopNode[`node_key`]]; !ok {
		flowL.outputs[flowLoop.LoopNode[`node_key`]] = make(common.SimpleFields)
	}
	flowL.outputs[flowLoop.LoopNode[`node_key`]][inField.Key] = *inField
}

// Inject intermediate variables into outputs for child nodes
func (flowLoop *WorkFlowLoop) fillLoopIntermediateToOutputs(flowL *WorkFlow) {
	if _, ok := flowL.outputs[flowLoop.LoopNode[`node_key`]]; !ok {
		flowL.outputs[flowLoop.LoopNode[`node_key`]] = make(common.SimpleFields)
	}
	for _, loopIntermediate := range flowLoop.LoopNodeParams.IntermediateParams {
		flowL.outputs[flowLoop.LoopNode[`node_key`]][loopIntermediate.Key] = loopIntermediate.SimpleField
	}
}

// Initialize intermediate variables
func (flowLoop *WorkFlowLoop) initLoopIntermediate() {
	for key, intermediateParam := range flowLoop.LoopNodeParams.IntermediateParams {
		flowLoop.LoopNodeParams.IntermediateParams[key].Value = flowLoop.Flow.VariableReplace(intermediateParam.Value)

		// Variable replacement supports arrays etc.
		var data any = flowLoop.Flow.VariableReplace(intermediateParam.Value)
		if tool.InArrayString(intermediateParam.Typ, common.TypArrays[:]) {
			var temp []any // Special handling for array types
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
