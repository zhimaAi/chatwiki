// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type WorkFlowParams struct {
	*define.ChatRequestParam
	CurMsgId   int
	DialogueId int
	SessionId  int
}

type WorkFlow struct {
	params      *WorkFlowParams
	nodeLogs    []common.NodeLog
	StartTime   int
	EndTime     int
	context     context.Context
	cancel      context.CancelFunc
	ticker      *time.Ticker //流程超时
	isTimeout   bool
	global      common.SimpleFields
	output      common.SimpleFields
	outputs     map[string]common.SimpleFields
	curNodeKey  string
	runNodeKeys []string
	curNode     NodeAdapter
	runLogs     []string
	running     bool
	isFinish    bool
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
	for {
		var nodeInfo msql.Params
		flow.Logs(`当前运行节点:%s`, flow.curNodeKey)
		flow.curNode, nodeInfo, err = GetNodeByKey(cast.ToUint(flow.params.Robot[`id`]), flow.curNodeKey)
		if err != nil {
			flow.Logs(err.Error())
		}
		if flow.curNode == nil {
			break //退出
		}
		flow.runNodeKeys = append(flow.runNodeKeys, flow.curNodeKey)
		var nextNodeKey string
		//节点运行开始
		nodeLog := common.NodeLog{
			StartTime: time.Now().UnixMilli(),
			NodeKey:   flow.curNodeKey,
			NodeName:  nodeInfo[`node_name`],
			NodeType:  cast.ToInt(nodeInfo[`node_type`]),
		}
		flow.output, nextNodeKey, err = flow.curNode.Running(flow)
		flow.outputs[flow.curNodeKey] = flow.output //记录每个节点输出的变量
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = flow.output
		nodeLog.Error = err
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flow.nodeLogs = append(flow.nodeLogs, nodeLog)
		//节点运行结束
		if !flow.isFinish {
			flow.Logs(`结果nextNodeKey:%s,err:%v`, nextNodeKey, err)
			if len(flow.output) > 0 {
				flow.Logs(`输出变量:%s`, tool.JsonEncodeNoError(flow.output))
			}
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

func (flow *WorkFlow) Ending() {
	flow.Logs(`保存工作流运行日志...`)
	_, err := msql.Model(`work_flow_logs`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: flow.params.AdminUserId,
		`robot_id`:      flow.params.Robot[`id`],
		`openid`:        flow.global[`openid`].GetVal(common.TypString),
		`run_node_keys`: strings.Join(flow.runNodeKeys, `,`),
		`run_logs`:      tool.JsonEncodeNoError(flow.runLogs),
		`create_time`:   flow.StartTime, //这里放开始时间
		`update_time`:   flow.EndTime,   //这里放结束时间
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

func (flow *WorkFlow) GetVariable(key string) (field common.SimpleField, exist bool) {
	if field, exist = flow.output[key]; exist {
		return
	}
	if strings.HasPrefix(key, `global.`) {
		realKey := strings.TrimPrefix(key, `global.`)
		if field, exist = flow.global[realKey]; exist {
			return
		}
	}
	return
}

func SysGlobalVariables() []string { //固定值,runtime时不可变更
	return []string{`global.question`, `global.openid`}
}

func RunningWorkFlow(params *WorkFlowParams) (*WorkFlow, error) {
	ctx, cancel := context.WithCancel(context.Background())
	flow := &WorkFlow{
		params:    params,
		nodeLogs:  make([]common.NodeLog, 0),
		StartTime: tool.Time2Int(),
		context:   ctx,
		cancel:    cancel,
		ticker:    time.NewTicker(time.Minute * 5), //DIY
		global: common.SimpleFields{
			`question`: common.SimpleField{Sys: true, Key: `question`, Desc: tea.String(`用户消息`), Typ: common.TypString, Vals: []common.Val{{String: &params.Question}}},
			`openid`:   common.SimpleField{Sys: true, Key: `openid`, Desc: tea.String(`用户openid`), Typ: common.TypString, Vals: []common.Val{{String: &params.Openid}}},
		},
		outputs:     make(map[string]common.SimpleFields), //记录每个节点输出的变量
		curNodeKey:  params.Robot[`start_node_key`],       //开始节点
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
	err := flow.Running() //运行流程
	if flow.isTimeout {
		err = errors.New(`工作流执行超时`)
	} else if !flow.isFinish {
		err = errors.New(`工作流未到结束节点`)
	}
	go flow.Ending() //记录runtime日志
	return flow, err //返回数据
}

func CallWorkFlow(params *WorkFlowParams, debugLog *[]any, monitor *common.Monitor) (content string, requestTime int64, libUseTime common.LibUseTime, list []msql.Params, err error) {
	if len(params.Robot[`start_node_key`]) == 0 {
		err = errors.New(`工作流机器人未发布`)
		return
	}
	flow, err := RunningWorkFlow(params)
	if flow != nil && len(flow.nodeLogs) > 0 {
		monitor.NodeLogs = flow.nodeLogs //记录监控数据
	}
	if flow == nil || err != nil {
		return
	}
	content = cast.ToString(flow.output[`special.llm_reply_content`].GetVal(common.TypString))
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
	return
}
