// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type WorkFlowBatch struct {
	MainFlow        *WorkFlow
	FlowLs          []*WorkFlow
	Outputs         common.SimpleFields
	BatchNodeParams *BatchNodeParams
	BatchNode       msql.Params
}

// FlowRunningBatchTest Batch node separate run test
func FlowRunningBatchTest(flow *WorkFlow) error {
	defer func() {
		if r := recover(); r != nil {
			logs.Debug(`Error occurred %s`, debug.Stack())
		}
	}()
	_, childNode, err := GetNodeByKey(flow, cast.ToUint(flow.params.RealRobot[`id`]), flow.curNodeKey)
	if err != nil {
		return err
	}
	workFlowBatch, err := NewWorkFlowBatch(childNode[`loop_parent_key`], nil, flow)
	if err != nil {
		return err
	}
	_, err = workFlowBatch.ForRunning()
	if err != nil {
		return err
	}
	if len(workFlowBatch.FlowLs) > 0 {
		for _, flowL := range workFlowBatch.FlowLs {
			flow.nodeLogs = append(flow.nodeLogs, flowL.nodeLogs...)
			flow.runLogs = append(flow.runLogs, flowL.runLogs...)
		}
		flow.isFinish = true
	} else {
		err = errors.New(i18n.Show(flow.params.Lang, `batch_node_not_executed`))
	}
	return err
}

// BatchNodeRunning Batch node running
func BatchNodeRunning(nodeInfo msql.Params, flow *WorkFlow) (common.SimpleFields, string, error) {
	workFlowBatch, err := NewWorkFlowBatch(``, nodeInfo, flow)
	if err != nil {
		return make(common.SimpleFields), ``, err
	}
	nextNodeKey, err := workFlowBatch.ForRunning()
	if err != nil {
		return common.SimpleFields{}, ``, err
	}
	return workFlowBatch.Outputs, nextNodeKey, nil
}

func NewWorkFlowBatch(nodeKey string, node msql.Params, flow *WorkFlow) (*WorkFlowBatch, error) {
	var err error
	flowBatch := &WorkFlowBatch{
		MainFlow:        flow,
		FlowLs:          make([]*WorkFlow, 0),
		Outputs:         map[string]common.SimpleField{},
		BatchNodeParams: &BatchNodeParams{},
		BatchNode:       make(msql.Params),
	}
	if len(node) > 0 {
		flowBatch.BatchNode = node
	} else if nodeKey != `` {
		_, flowBatch.BatchNode, err = GetNodeByKey(flowBatch.MainFlow, cast.ToUint(flowBatch.MainFlow.params.RealRobot[`id`]), nodeKey)
		if err != nil {
			return flowBatch, err
		}
	} else {
		err = errors.New(i18n.Show(flowBatch.MainFlow.params.Lang, `batch_node_config_empty`))
		return flowBatch, err
	}
	if cast.ToInt(flowBatch.BatchNode[`node_type`]) != NodeTypeBatch {
		err = errors.New(i18n.Show(flowBatch.MainFlow.params.Lang, `batch_node_invalid`))
		return flowBatch, err
	}
	nodeParams := &NodeParams{}
	err = tool.JsonDecode(flowBatch.BatchNode[`node_params`], &nodeParams)
	if err != nil {
		logs.Error(`Node parameter parsing error: %s`, err)
		return flowBatch, err
	}
	flowBatch.BatchNodeParams = &nodeParams.Batch
	return flowBatch, err
}

func (flowBatch *WorkFlowBatch) ForRunning() (nextNodeKey string, err error) {
	nextNodeKey = ``
	flowBatch.testFillParams(flowBatch.MainFlow)
	batchDatas := flowBatch.getBatchDatas()
	if len(batchDatas) == 0 {
		flowBatch.MainFlow.Logs(`Batch data not found`)
		return nextNodeKey, errors.New(i18n.Show(flowBatch.MainFlow.params.Lang, `batch_data_not_found`))
	}
	// Only one run allowed for separate testing
	if flowBatch.MainFlow.params.IsTestBatchNodeRun {
		batchDatas = []*common.SimpleField{batchDatas[0]}
	}
	if len(batchDatas) > flowBatch.BatchNodeParams.MaxRunNumber.Int() {
		flowBatch.MainFlow.Logs(`Batch processing incoming data %d, max execution %d`, len(batchDatas), flowBatch.BatchNodeParams.MaxRunNumber.Int())
		batchDatas = batchDatas[:flowBatch.BatchNodeParams.MaxRunNumber.Int()]
	}
	flowBatch.MainFlow.Logs(`Start batch processing data %d times..`, len(batchDatas))
	var wg sync.WaitGroup
	cTask := make(chan any, flowBatch.BatchNodeParams.ChanNumber.Int())
	for _, inField := range batchDatas {
		wg.Add(1)
		cTask <- nil
		go func(data *common.SimpleField) {
			defer wg.Done()
			defer func() { <-cTask }()
			if flowBatch.MainFlow.isTimeout {
				return
			}
			flowL := &WorkFlow{}
			flowL, err := flowBatch.batchNodeRunning(data)
			if err != nil {
				flowBatch.MainFlow.Logs(`Batch node running error: %s`, err)
			} else {
				flowBatch.FlowLs = append(flowBatch.FlowLs, flowL)
			}
		}(inField)
	}
	wg.Wait()
	close(cTask)
	//take result
	flowBatch.TakeChildOutputs()
	nextNodeKey = flowBatch.BatchNode[`next_node_key`]
	return
}

func (flowBatch *WorkFlowBatch) getBatchDatas() []*common.SimpleField {
	for _, batchArray := range flowBatch.BatchNodeParams.BatchArrays {
		if batchArray.NodeKey() == `global` { // Start node
			for globalKey, globalVal := range flowBatch.MainFlow.global {
				if globalKey == batchArray.ChooseKey() {
					return flowBatch.inFieldAppend(batchArray.Key, globalVal)
				}
			}
		} else { // Non-start node
			for outNodeKey, nodeOutputs := range flowBatch.MainFlow.outputs {
				if outNodeKey != batchArray.NodeKey() {
					continue
				}
				for _, outField := range nodeOutputs {
					if outField.Typ == batchArray.Typ && outField.Key == batchArray.ChooseKey() {
						return flowBatch.inFieldAppend(batchArray.Key, outField)
					}
				}
			}
		}
	}
	return make([]*common.SimpleField, 0)
}

func (flowBatch *WorkFlowBatch) inFieldAppend(key string, outField common.SimpleField) []*common.SimpleField {
	inFields := make([]*common.SimpleField, 0)
	for _, val := range outField.Vals {
		switch outField.Typ {
		case common.TypArrFloat:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypFloat,
				Vals: []common.Val{val},
			})
		case common.TypArrObject:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypObject,
				Vals: []common.Val{val},
			})
		case common.TypArrBoole:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypBoole,
				Vals: []common.Val{val},
			})
		case common.TypArrNumber:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypNumber,
				Vals: []common.Val{val},
			})
		case common.TypArrParams:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypParams,
				Vals: []common.Val{val},
			})
		case common.TypArrString:
			inFields = append(inFields, &common.SimpleField{
				Key:  key,
				Typ:  common.TypString,
				Vals: []common.Val{val},
			})
		}
	}
	return inFields
}

func (flowBatch *WorkFlowBatch) TakeChildOutputs() {
	completionToken, promptToken := 0, 0
	for _, flowL := range flowBatch.FlowLs {
		//take token
		for _, nodeLog := range flowL.nodeLogs {
			nodeOutputTake := NodeOutputTake{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(nodeLog.Output), &nodeOutputTake)
			if err != nil {
				flowBatch.MainFlow.Logs(`Exception: Batch node extract child node output parameter error %s`, err.Error())
				continue
			}
			completionToken += cast.ToInt(nodeOutputTake.LlmResult.CompletionToken)
			promptToken += cast.ToInt(nodeOutputTake.LlmResult.PromptToken)
		}
		logs.Debug(`Output of previous workflow %s`, tool.JsonEncodeNoError(flowL.outputs))
		logs.Debug(`Data to be extracted %s`, tool.JsonEncodeNoError(flowBatch.BatchNodeParams.Output))
		//take result
		for childNodeKey, childOutput := range flowL.outputs {
			for _, childOutField := range childOutput {
				for outKey, outField := range flowBatch.BatchNodeParams.Output {
					if outField.Value == childNodeKey+`.`+childOutField.Key { //find the key
						flowBatch.BatchNodeParams.Output[outKey].Vals = append(flowBatch.BatchNodeParams.Output[outKey].Vals, childOutField.Vals...)
					}
				}
			}
		}
	}
	//batch output
	for _, outField := range flowBatch.BatchNodeParams.Output {
		flowBatch.Outputs[outField.Key] = outField.SimpleField
	}
	// Inject token into result
	batchNumber := len(flowBatch.FlowLs)
	flowBatch.Outputs[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &completionToken}}}
	flowBatch.Outputs[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &promptToken}}}
	flowBatch.Outputs[`batch_result.batch_number`] = common.SimpleField{Key: `batch_result.batch_number`, Typ: common.TypNumber, Vals: []common.Val{{Number: &batchNumber}}}

	// Test output log
	if define.IsDev {
		//for key, log := range childFlowLogs {
		//	flowBatch.Outputs[key] = common.SimpleField{Key: key, Typ: common.TypArrObject, Vals: []common.Val{{Object: log}}}
		//}
	}
}

func (flowBatch *WorkFlowBatch) batchNodeRunning(dataField *common.SimpleField) (flowL *WorkFlow, err error) {
	//new flow
	flowL = &WorkFlow{
		params:      flowBatch.MainFlow.params,            //inherit params
		nodeLogs:    make([]common.NodeLog, 0),            //node logs
		StartTime:   tool.Time2Int(),                      //start time
		context:     flowBatch.MainFlow.context,           //inherit context
		global:      flowBatch.MainFlow.global,            // Inherited from the owning workflow
		outputs:     make(map[string]common.SimpleFields), // Deep copy of main flow outputs
		inputs:      make(map[string]common.SimpleFields), // Input parameters do not inherit from main flow
		runNodeKeys: make([]string, 0),                    //self run node keys
		runLogs:     make([]string, 0),                    //self  logs
		VersionId:   flowBatch.MainFlow.VersionId,         //inherit version id
	}
	flowBatch.copyOutputs(flowL.outputs)
	// Batch parameter injection for this run
	flowBatch.fillBatchInField(flowL, dataField)
	// Find entry child node
	flowL.curNodeKey, err = flowBatch.FindStartChildNodeKey()
	if err != nil {
		flowL.Logs(err.Error())
		return
	}
	if flowL.curNodeKey == `` {
		err = errors.New(i18n.Show(flowL.params.Lang, `batch_start_node_not_found`))
		return
	}
	for {
		var nodeInfo msql.Params
		flowL.curNode, nodeInfo, err = GetNodeByKey(flowL, cast.ToUint(flowL.params.RealRobot[`id`]), flowL.curNodeKey)
		if err != nil {
			flowL.Logs(err.Error())
			break
		}
		if flowL.curNode == nil {
			break // Exit
		}
		flowL.Logs(`Batch node currently running child node:%s %s`, flowL.curNodeKey, nodeInfo[`node_name`])
		flowL.runNodeKeys = append(flowL.runNodeKeys, flowL.curNodeKey)
		var nextNodeKey string

		// Node execution start
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flowL.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		flowL.getNodeInputs()
		flowL.inputs[flowL.curNodeKey] = flowL.input
		if tool.InArray(cast.ToInt(nodeInfo[`node_type`]), []int{NodeTypeLoop, NodeTypeBatch}) {
			err = errors.New(i18n.Show(flowL.params.Lang, `batch_node_not_support_loop_or_batch_node`, nodeLog.NodeName))
			break
		} else {
			flowL.output, nextNodeKey, err = flowL.curNode.Running(flowL)
			flowL.outputs[flowL.curNodeKey] = flowL.output // Record output variables of each node
		}
		// Running parameter processing
		nodeLog.Input = common.GetFieldsObject(common.GetRecurveFields(flowL.input))
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flowL.output))
		nodeLog.NodeOutput = GetNodeOutput(nodeLog.Output)
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flowL.nodeLogs = append(flowL.nodeLogs, nodeLog)
		// Node execution end
		flowBatch.MainFlow.Logs(`Result nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flowL.output) > 0 {
			flowBatch.MainFlow.Logs(`Output variables:%s`, tool.JsonEncodeNoError(flowL.output))
		}
		if flowBatch.MainFlow.isTimeout || err != nil || len(nextNodeKey) == 0 {
			break // End
		}
		// External interrupt listener processing
		logs.Debug(`Listener processing %v`, flowL.context)
		select {
		case <-flowL.context.Done():
			goto flowExit
		default: // Execute next node
			flowL.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flowL.EndTime = tool.Time2Int()
	flowL.running = false // Running end
	return
}

func (flowBatch *WorkFlowBatch) copyOutputs(outputs map[string]common.SimpleFields) {
	if len(flowBatch.MainFlow.outputs) == 0 {
		return
	}
	for nodeKey, simpleFields := range flowBatch.MainFlow.outputs {
		outputs[nodeKey] = common.SimpleFields{}
		for key, field := range simpleFields {
			outputs[nodeKey][key] = field
		}
	}
}

func (flowBatch *WorkFlowBatch) testFillParams(flowL *WorkFlow) {
	if !flowBatch.MainFlow.params.IsTestBatchNodeRun {
		return
	}
	fillTestParamsToRunningParams(flowL, flowL.params.BatchTestParams)
}

// fill batch data
func (flowBatch *WorkFlowBatch) fillBatchInField(flowL *WorkFlow, inField *common.SimpleField) {
	if inField == nil {
		return
	}
	if _, ok := flowL.outputs[flowBatch.BatchNode[`node_key`]]; !ok {
		flowL.outputs[flowBatch.BatchNode[`node_key`]] = make(common.SimpleFields)
	}
	flowL.outputs[flowBatch.BatchNode[`node_key`]][inField.Key] = *inField
}

func (flowBatch *WorkFlowBatch) FindStartChildNodeKey() (string, error) {
	return msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, flowBatch.BatchNode[`admin_user_id`]).
		Where(`robot_id`, flowBatch.BatchNode[`robot_id`]).
		Where(`loop_parent_key`, flowBatch.BatchNode[`node_key`]).
		Where(`node_type`, cast.ToString(NodeTypeBatchStart)).
		Value(`node_key`)
}
