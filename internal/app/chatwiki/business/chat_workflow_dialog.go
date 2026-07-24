// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var workflowDialogTempSeq int64

// workflowDialogStreamState records what has already been pushed so the SSE
// stream can avoid repeating the same reply again at workflow finish.
type workflowDialogStreamState struct {
	replyContentList             []common.ReplyContent
	hasPushedReply               bool
	hasCurrentRunTextReply       bool
	hasCurrentRunStructuredReply bool
	hasExplicitFinishText        bool
	lastPushedText               string
	lastReplyJSON                string
}

// CallWorkFlowDialog runs the draft workflow in a dialog-style SSE mode.
// Unlike /chat/request, it stays inside the current workflow and does not write
// formal dialogue/session/message records.
func CallWorkFlowDialog(c *gin.Context) {
	_ = c.Request.ParseMultipartForm(define.DefaultMultipartMemory)
	c.Header(`Content-Type`, `text/event-stream`)
	c.Header(`Cache-Control`, `no-cache`)
	c.Header(`Connection`, `keep-alive`)
	if define.IsDev {
		c.Header(`Access-Control-Allow-Origin`, `*`)
	}

	c.Set(`from_work_flow`, true)
	params := getChatRequestParam(c)
	if params != nil {
		params.HeaderToken = c.GetHeader(`token`)
		if len(params.HeaderToken) == 0 {
			params.HeaderToken = c.Query(`token`)
		}
	}

	workFlowParams, err := buildCallWorkFlowDialogParams(c, params)
	chanStream := make(chan sse.Event)
	go func() {
		defer close(chanStream)
		if params == nil {
			streamWorkflowDialogError(chanStream, errors.New(`chat request params is nil`))
			return
		}
		if params.Error != nil {
			streamWorkflowDialogError(chanStream, params.Error)
			return
		}
		if err != nil {
			streamWorkflowDialogError(chanStream, err)
			return
		}
		runCallWorkFlowDialog(workFlowParams, chanStream)
	}()

	c.Stream(func(_ io.Writer) bool {
		if event, ok := <-chanStream; ok {
			if data, ok := event.Data.(string); ok {
				event.Data = strings.ReplaceAll(data, "\r", ``)
			}
			c.SSEvent(event.Event, event.Data)
			return true
		}
		return false
	})
	if params != nil && params.IsClose != nil {
		*params.IsClose = true
	}
	for range chanStream {
	}
}

// buildCallWorkFlowDialogParams reuses the existing draft workflow assembly
// path, then adds the stricter entry checks required by this test-only API.
func buildCallWorkFlowDialogParams(c *gin.Context, chatRequestParam *define.ChatRequestParam) (*work_flow.WorkFlowParams, error) {
	if chatRequestParam == nil {
		return nil, errors.New(`chat request params is nil`)
	}
	if chatRequestParam.Error != nil {
		return nil, chatRequestParam.Error
	}
	if cast.ToInt(chatRequestParam.Robot[`application_type`]) != define.ApplicationTypeFlow {
		return nil, errors.New(i18n.Show(chatRequestParam.Lang, `param_invalid`, `robot_key`))
	}
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam: chatRequestParam,
		DialogueId:       cast.ToInt(c.PostForm(`dialogue_id`)),
		SessionId:        cast.ToInt(c.PostForm(`session_id`)),
		IsDialogMode:     true, //workflow test dialog mode
		Draft:            work_flow.Draft{IsDraft: true, QuestionMultipleSwitch: cast.ToBool(c.PostForm(`question_multiple_switch`))},
	}
	var err error
	workFlowParams.Draft.NodeMaps, err = msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).
		Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).
		ColumnMap(`*`, `node_key`)
	if err != nil {
		return nil, err
	}
	nodeList := make([]work_flow.WorkFlowNode, 0, len(workFlowParams.Draft.NodeMaps))
	for _, params := range workFlowParams.Draft.NodeMaps {
		node := work_flow.WorkFlowNode{
			NodeType:      common.MixedInt(cast.ToInt(params[`node_type`])),
			NodeName:      params[`node_name`],
			NodeKey:       params[`node_key`],
			NodeParams:    work_flow.NodeParams{},
			NodeInfoJson:  make(map[string]any),
			NextNodeKey:   params[`next_node_key`],
			LoopParentKey: params[`loop_parent_key`],
		}
		_ = tool.JsonDecodeUseNumber(params[`node_params`], &node.NodeParams)
		_ = tool.JsonDecodeUseNumber(params[`node_info_json`], &node.NodeInfoJson)
		nodeList = append(nodeList, node)
	}
	var errNodeKey string
	workFlowParams.Draft.StartNodeKey, _, _, _, err, errNodeKey = work_flow.VerifyWorkFlowNodes(nodeList, chatRequestParam.AdminUserId, chatRequestParam.Lang)
	if err != nil {
		return nil, err
	}
	hasChatTrigger, err := hasEnabledChatTrigger(workFlowParams.Draft.NodeMaps, workFlowParams.Draft.StartNodeKey, chatRequestParam.Lang)
	if err != nil {
		return nil, err
	}
	if !hasChatTrigger {
		return nil, errors.New(i18n.Show(chatRequestParam.Lang, `no_enabled_trigger_for_scenario`))
	}
	workFlowParams.Draft.QuestionMultipleSwitch, err = getChatTriggerQuestionMultipleSwitch(
		workFlowParams.Draft.NodeMaps,
		workFlowParams.Draft.StartNodeKey,
		chatRequestParam.Lang,
	)
	if err != nil {
		return nil, err
	}
	_ = errNodeKey
	return workFlowParams, nil
}

// runCallWorkFlowDialog maps workflow callbacks to SSE events and keeps the
// temporary dialogue/session ids stable for question-node resume.
func runCallWorkFlowDialog(params *work_flow.WorkFlowParams, chanStream chan sse.Event) {
	if params.DialogueId <= 0 || params.SessionId <= 0 {
		params.DialogueId, params.SessionId = generateWorkflowDialogTempIDs()
	}
	state := &workflowDialogStreamState{
		replyContentList: make([]common.ReplyContent, 0),
	}
	params.ImmediatelyReplyHandle = func(replyContent common.ReplyContent) {
		state.replyContentList = append(state.replyContentList, replyContent)
		state.hasPushedReply = true
		state.hasCurrentRunStructuredReply = true
		if text := takeReplyText(replyContent); text != `` {
			state.lastPushedText = text
		}
		state.lastReplyJSON = tool.JsonEncodeNoError([]common.ReplyContent{replyContent})
		chanStream <- sse.Event{Event: `reply_content_list`, Data: state.replyContentList}
	}
	params.NodeLogHandle = func(nodeLog common.NodeLog) {
		chanStream <- sse.Event{Event: `workflow_node_log`, Data: nodeLog}
		if nodeLog.NodeType == work_flow.NodeTypeFinish {
			state.hasExplicitFinishText = hasExplicitFinishText(nodeLog)
			return
		}
		replyList := takeStructuredReplyContentListFromNodeLog(nodeLog)
		if len(replyList) > 0 {
			state.hasCurrentRunStructuredReply = true
			replyJSON := tool.JsonEncodeNoError(replyList)
			if replyJSON == state.lastReplyJSON {
				return
			}
			state.replyContentList = append(state.replyContentList, replyList...)
			state.hasPushedReply = true
			state.lastReplyJSON = replyJSON
			if text := takeReplyText(replyList[0]); text != `` {
				state.lastPushedText = text
			}
			chanStream <- sse.Event{Event: `reply_content_list`, Data: state.replyContentList}
			return
		}
		if takeNodeLogTextReply(nodeLog) != `` {
			state.hasCurrentRunTextReply = true
		}
		if nodeLog.NodeType == work_flow.NodeTypeQuestion {
			if text := takeNodeLogTextReply(nodeLog); text != `` {
				state.hasPushedReply = true
				state.lastPushedText = text
				chanStream <- sse.Event{Event: `sending`, Data: text}
			}
		}
	}

	chanStream <- sse.Event{Event: `ping`, Data: `ping`}
	chanStream <- sse.Event{Event: `dialogue_id`, Data: params.DialogueId}
	chanStream <- sse.Event{Event: `session_id`, Data: params.SessionId}

	flow, nodeLogs, err := work_flow.BaseCallWorkFlow(params)
	if err != nil {
		streamWorkflowDialogError(chanStream, err)
		return
	}
	content, replyContentList := work_flow.TakeOutputReply(flow)
	if len(replyContentList) > 0 {
		currentReplyJSON := tool.JsonEncodeNoError(replyContentList)
		if shouldSendWorkflowDialogFinalReplyList(currentReplyJSON, state.lastReplyJSON, state.hasCurrentRunStructuredReply) {
			state.replyContentList = append(state.replyContentList, replyContentList...)
			chanStream <- sse.Event{Event: `reply_content_list`, Data: state.replyContentList}
		}
	} else if (state.hasCurrentRunTextReply || state.hasExplicitFinishText) &&
		shouldSendWorkflowDialogFinalText(content, state.lastPushedText, state.hasPushedReply, replyContentList) {
		chanStream <- sse.Event{Event: `sending`, Data: content}
	}
	// Always push a finish event so the client knows the current round ended,
	// whether the workflow completed normally or paused at a question node
	// waiting for user input. The finish payload carries an is_paused flag so
	// the client can tell "wait for input" apart from "done".
	paused := work_flow.IsStoragePaused(flow)
	useToken, useMills := common.TakeWorkFlowTestUseToken(nodeLogs)
	chanStream <- sse.Event{Event: `finish`, Data: map[string]any{
		`status`:     map[bool]string{true: `paused`, false: `done`}[paused],
		`is_paused`:  paused,
		`dialog_id`:  params.DialogueId,
		`session_id`: params.SessionId,
		`use_token`:  useToken,
		`use_mills`:  useMills,
	}}
}

// streamWorkflowDialogError keeps the SSE contract consistent on failures.
func streamWorkflowDialogError(chanStream chan sse.Event, err error) {
	if err == nil {
		err = errors.New(`unknown error`)
	}
	chanStream <- sse.Event{Event: `error`, Data: err.Error()}
	chanStream <- sse.Event{Event: `finish`, Data: `[DONE]`}
}

// generateWorkflowDialogTempIDs creates int-compatible temporary ids without
// touching the formal dialogue/session tables.
func generateWorkflowDialogTempIDs() (int, int) {
	seq := atomic.AddInt64(&workflowDialogTempSeq, 1) % 10000
	base := time.Now().UnixMilli()*10000 + seq*10
	return cast.ToInt(base), cast.ToInt(base + 1)
}

// hasEnabledChatTrigger enforces that only workflows with an enabled chat
// trigger can enter this dialog-style debug interface.
func hasEnabledChatTrigger(nodeMaps map[string]msql.Params, startNodeKey, lang string) (bool, error) {
	trigger, found, err := getEnabledChatTrigger(nodeMaps, startNodeKey, lang)
	if err != nil {
		return false, err
	}
	return found && trigger.TriggerType == work_flow.TriggerTypeChat, nil
}

// getChatTriggerQuestionMultipleSwitch reads the multimodal-input switch from
// the enabled chat trigger on the draft workflow start node.
func getChatTriggerQuestionMultipleSwitch(nodeMaps map[string]msql.Params, startNodeKey, lang string) (bool, error) {
	trigger, found, err := getEnabledChatTrigger(nodeMaps, startNodeKey, lang)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	return trigger.TriggerChatConfig.QuestionMultipleSwitch, nil
}

// getEnabledChatTrigger returns the first enabled chat trigger configured on
// the start node so dialog debugging can follow workflow trigger settings.
func getEnabledChatTrigger(nodeMaps map[string]msql.Params, startNodeKey, lang string) (work_flow.TriggerConfig, bool, error) {
	startNode, ok := nodeMaps[startNodeKey]
	if !ok {
		return work_flow.TriggerConfig{}, false, errors.New(i18n.Show(lang, "workflow_no_start_node_data"))
	}
	rawNodeParams := startNode[`node_params`]
	nodeParams := work_flow.NodeParams{}
	if err := tool.JsonDecodeUseNumber(rawNodeParams, &nodeParams); err != nil {
		return work_flow.TriggerConfig{}, false, err
	}
	if len(nodeParams.Start.TriggerList) == 0 {
		nodeParams = work_flow.DisposeNodeParams(work_flow.NodeTypeStart, rawNodeParams, lang)
	}
	for _, trigger := range nodeParams.Start.TriggerList {
		if trigger.TriggerType == work_flow.TriggerTypeChat && trigger.TriggerSwitch {
			return trigger, true, nil
		}
	}
	return work_flow.TriggerConfig{}, false, nil
}

// takeStructuredReplyContentListFromNodeLog extracts menu/specified replies
// from node output so they can be streamed as soon as the node completes.
func takeStructuredReplyContentListFromNodeLog(nodeLog common.NodeLog) []common.ReplyContent {
	special, ok := nodeLog.Output[`special`].(map[string]any)
	if !ok {
		return nil
	}
	rawReplyList, ok := special[`reply_content_list`]
	if !ok {
		return nil
	}
	replyList := make([]common.ReplyContent, 0)
	if err := tool.JsonDecode(cast.ToString(rawReplyList), &replyList); err != nil {
		return nil
	}
	return replyList
}

// takeNodeLogTextReply extracts default text replies from serialized output.
func takeNodeLogTextReply(nodeLog common.NodeLog) string {
	special, ok := nodeLog.Output[`special`].(map[string]any)
	if !ok {
		return ``
	}
	for _, outputKey := range []string{`llm_reply_content`, `question_optimize_reply_content`, `mcp_reply_content`} {
		if content := cast.ToString(special[outputKey]); content != `` {
			return content
		}
	}
	return ``
}

// hasExplicitFinishText distinguishes configured finish messages from output
// inherited from the node that ran before the finish node.
func hasExplicitFinishText(nodeLog common.NodeLog) bool {
	for outputKey := range nodeLog.Output {
		if strings.HasPrefix(outputKey, define.FinishReplyPrefixKey) {
			return true
		}
	}
	return false
}

// takeReplyText picks the most user-visible text from a structured reply item
// so finish-stage deduplication can compare by content.
func takeReplyText(replyContent common.ReplyContent) string {
	if replyContent.Description != `` {
		return replyContent.Description
	}
	if replyContent.Title != `` {
		return replyContent.Title
	}
	return ``
}

// shouldSendWorkflowDialogFinalText prevents duplicate terminal text when the
// workflow has already pushed the same reply during execution.
func shouldSendWorkflowDialogFinalText(content, lastPushedText string, hasPushedReply bool, finalReplyList []common.ReplyContent) bool {
	if len(finalReplyList) > 0 {
		return false
	}
	content = strings.TrimSpace(content)
	if content == `` {
		return false
	}
	if hasPushedReply && strings.TrimSpace(lastPushedText) == content {
		return false
	}
	return true
}

func shouldSendWorkflowDialogFinalReplyList(currentReplyJSON, lastReplyJSON string, hasCurrentRunStructuredReply bool) bool {
	return hasCurrentRunStructuredReply && currentReplyJSON != lastReplyJSON
}

// is paused
func shouldEmitWorkflowDialogFinish(flow *work_flow.WorkFlow, err error) bool {
	if err != nil || flow == nil {
		return false
	}
	return !work_flow.IsStoragePaused(flow)
}
