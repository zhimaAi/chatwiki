// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
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

// FlowRunningBatchTest 批处理节点单独运行测试
func FlowRunningBatchTest(flow *WorkFlow) error {
	defer func() {
		if r := recover(); r != nil {
			logs.Debug(`报错了 %s`, debug.Stack())
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
		err = errors.New(`未执行批处理节点`)
	}
	return err
}

// BatchNodeRunning 批处理节点运行
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
		err = errors.New(`批处理节点配置不能为空`)
		return flowBatch, err
	}
	if cast.ToInt(flowBatch.BatchNode[`node_type`]) != NodeTypeBatch {
		err = errors.New(`错误的批处理节点`)
		return flowBatch, err
	}
	nodeParams := &NodeParams{}
	err = tool.JsonDecode(flowBatch.BatchNode[`node_params`], &nodeParams)
	if err != nil {
		logs.Error(`节点参数解析错误: %s`, err)
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
		flowBatch.MainFlow.Logs(`未找到批处理数据`)
		return nextNodeKey, errors.New(`未找到批处理数据`)
	}
	//单独运行测试只允许一次
	if flowBatch.MainFlow.params.IsTestBatchNodeRun {
		batchDatas = []*common.SimpleField{batchDatas[0]}
	}
	if len(batchDatas) > flowBatch.BatchNodeParams.MaxRunNumber.Int() {
		flowBatch.MainFlow.Logs(`批处理传入数据%d，最大执行%d`, len(batchDatas), flowBatch.BatchNodeParams.MaxRunNumber.Int())
		batchDatas = batchDatas[:flowBatch.BatchNodeParams.MaxRunNumber.Int()]
	}
	flowBatch.MainFlow.Logs(`开始批处理数据%d次..`, len(batchDatas))
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
				flowBatch.MainFlow.Logs(`批处理节点运行错误: %s`, err)
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
		if batchArray.NodeKey() == `global` { //开始节点
			for globalKey, globalVal := range flowBatch.MainFlow.global {
				if globalKey == batchArray.ChooseKey() {
					return flowBatch.inFieldAppend(batchArray.Key, globalVal)
				}
			}
		} else { //非开始节点
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
		logs.Debug(`到这里了 11`)
		//take token
		for _, nodeLog := range flowL.nodeLogs {
			nodeOutputTake := NodeOutputTake{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(nodeLog.Output), &nodeOutputTake)
			if err != nil {
				flowBatch.MainFlow.Logs(`异常:批处理节点提取子节点输出参数错误 %s`, err.Error())
				continue
			}
			completionToken += cast.ToInt(nodeOutputTake.LlmResult.CompletionToken)
			promptToken += cast.ToInt(nodeOutputTake.LlmResult.PromptToken)
		}
		logs.Debug(`上一个工作流的输出 %s`, tool.JsonEncodeNoError(flowL.outputs))
		logs.Debug(`准备提取的数据 %s`, tool.JsonEncodeNoError(flowBatch.BatchNodeParams.Output))
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
	logs.Debug(`到这里了4`)
	//batch output
	for _, outField := range flowBatch.BatchNodeParams.Output {
		flowBatch.Outputs[outField.Key] = outField.SimpleField
	}
	logs.Debug(`到这里了5`)
	//注入token到结果中
	batchNumber := len(flowBatch.FlowLs)
	flowBatch.Outputs[`llm_result.completion_token`] = common.SimpleField{Key: `llm_result.completion_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &completionToken}}}
	flowBatch.Outputs[`llm_result.prompt_token`] = common.SimpleField{Key: `llm_result.prompt_token`, Typ: common.TypNumber, Vals: []common.Val{{Number: &promptToken}}}
	flowBatch.Outputs[`batch_result.batch_number`] = common.SimpleField{Key: `batch_result.batch_number`, Typ: common.TypNumber, Vals: []common.Val{{Number: &batchNumber}}}

	//测试输出日志
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
		global:      flowBatch.MainFlow.global,            //继承自所属工作流
		outputs:     make(map[string]common.SimpleFields), //深度拷贝主流程的输出
		runNodeKeys: make([]string, 0),                    //self run node keys
		runLogs:     make([]string, 0),                    //self  logs
		VersionId:   flowBatch.MainFlow.VersionId,         //inherit version id
	}
	flowBatch.copyOutputs(flowL.outputs)
	//批处理本次参数注入
	flowBatch.fillBatchInField(flowL, dataField)
	//找到入口子节点
	flowL.curNodeKey, err = flowBatch.FindStartChildNodeKey()
	if err != nil {
		flowL.Logs(err.Error())
		return
	}
	if flowL.curNodeKey == `` {
		err = errors.New(`未找到批处理的开始节点`)
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
			break //退出
		}
		flowL.Logs(`批处理节点当前运行子节点:%s %s`, flowL.curNodeKey, nodeInfo[`node_name`])
		flowL.runNodeKeys = append(flowL.runNodeKeys, flowL.curNodeKey)
		var nextNodeKey string

		//节点运行开始
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flowL.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		if tool.InArray(cast.ToInt(nodeInfo[`node_type`]), []int{NodeTypeLoop, NodeTypeBatch}) {
			err = errors.New(fmt.Sprintf(`批处理节点中不支持循环节点或者批处理节点 %s`, nodeLog.NodeName))
			break
		} else {
			flowL.output, nextNodeKey, err = flowL.curNode.Running(flowL)
			flowL.outputs[flowL.curNodeKey] = flowL.output //记录每个节点输出的变量
		}
		//运行参数处理
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flowL.output))
		nodeLog.ErrorMsg = fmt.Sprintf(`%v`, err)
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flowL.nodeLogs = append(flowL.nodeLogs, nodeLog)
		//节点运行结束
		flowBatch.MainFlow.Logs(`结果nextNodeKey:%s,err:%v`, nextNodeKey, err)
		if len(flowL.output) > 0 {
			flowBatch.MainFlow.Logs(`输出变量:%s`, tool.JsonEncodeNoError(flowL.output))
		}
		if flowBatch.MainFlow.isTimeout || err != nil || len(nextNodeKey) == 0 {
			break //结束
		}
		//外部中断监听处理
		logs.Debug(`监听处理 %v`, flowL.context)
		select {
		case <-flowL.context.Done():
			goto flowExit
		default: //执行下一个节点
			flowL.curNodeKey = nextNodeKey
		}
	}
flowExit:
	flowL.EndTime = tool.Time2Int()
	flowL.running = false //运行结束
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
