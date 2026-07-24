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

type WorkFlowRestoreData struct {
	NodeLogs    []common.NodeLog
	Global      common.SimpleFields
	Outputs     map[string]common.SimpleFields
	CurNodeKey  string
	RunNodeKeys []string
	RunLogs     []string
}

func SetWorkFlowStorage(flow *WorkFlow) {
	data := &WorkFlowRestoreData{
		NodeLogs:    flow.nodeLogs,
		Global:      flow.global,
		Outputs:     flow.outputs,
		CurNodeKey:  flow.curNodeKey,
		RunNodeKeys: flow.runNodeKeys,
		RunLogs:     flow.runLogs,
	}
	dialogId, sessionId := getWorkFlowStorageIDs(flow)
	boolCreate := true
	if sessionId > 0 {
		_, storeId, _ := GetWorkFlowRestore(dialogId, sessionId)
		if cast.ToInt(storeId) > 0 {
			_, err := msql.Model(`work_flow_storage_cache`, define.Postgres).Where(`id`, storeId).Update(msql.Datas{
				"storage":     tool.JsonEncodeNoError(data),
				"log_id":      flow.LogID, // preserve the stable identity of this run
				"update_time": time.Now().Unix(),
			})
			if err != nil {
				logs.Error(err.Error())
			}
			boolCreate = false
		}
	}
	if boolCreate {
		newId, err := msql.Model(`work_flow_storage_cache`, define.Postgres).Insert(msql.Datas{
			"admin_user_id": flow.params.AdminUserId,
			"robot_id":      flow.params.Robot[`id`],
			"dialog_id":     dialogId,
			"session_id":    sessionId,
			"openid":        flow.params.Openid,
			"storage":       tool.JsonEncodeNoError(data),
			"log_id":        flow.LogID, // preserve the stable identity of this run
			"create_time":   time.Now().Unix(),
			"update_time":   time.Now().Unix(),
		}, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
		if shouldRewriteDraftStorageSessionID(flow) {
			flow.params.SessionId = cast.ToInt(newId)
			_, err = msql.Model(`work_flow_storage_cache`, define.Postgres).Where(`id`, cast.ToString(newId)).Update(msql.Datas{
				`session_id`:  newId,
				"update_time": time.Now().Unix(),
			})
			if err != nil {
				logs.Error(err.Error())
			}
		}
	}
}

func shouldRewriteDraftStorageSessionID(flow *WorkFlow) bool {
	if flow == nil || flow.params == nil {
		return false
	}
	return flow.params.Draft.IsDraft && !flow.params.IsDialogMode
}

// get storage ids
func getWorkFlowStorageIDs(flow *WorkFlow) (dialogId, sessionId int) {
	if flow == nil || flow.params == nil {
		return 0, 0
	}
	if flow.params.Draft.IsDraft && !flow.params.IsDialogMode {
		return flow.params.AdminUserId, flow.params.SessionId
	}
	return flow.params.DialogueId, flow.params.SessionId
}

func GetWorkFlowRestore(dialogId, sessionId int) (data *WorkFlowRestoreData, s string, logID int64) {
	data = &WorkFlowRestoreData{}
	if dialogId == 0 || sessionId == 0 {
		return
	}
	//2 hours cache
	storage, err := msql.Model(`work_flow_storage_cache`, define.Postgres).
		Where(`dialog_id`, cast.ToString(dialogId)).
		Where(`session_id`, cast.ToString(sessionId)).
		Where(`create_time`, `>`, cast.ToString(time.Now().Unix()-7200)).Order(`id desc`).Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(storage) == 0 {
		return
	}
	err = tool.JsonDecode(storage[`storage`], data)
	if err != nil {
		logs.Error(err.Error())
	}
	return data, storage[`id`], cast.ToInt64(storage[`log_id`])
}

func IsStoragePaused(flow *WorkFlow) bool {
	if flow == nil {
		return false
	}
	return flow.isStorage
}

func DelWorkFlowStorage(dialogId, sessionId int) {
	_, err := msql.Model(`work_flow_storage_cache`, define.Postgres).
		Where(`dialog_id`, cast.ToString(dialogId)).
		Where(`session_id`, cast.ToString(sessionId)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
}

// DelWorkFlowStorageByLogID removes paused snapshots by log_id.
func DelWorkFlowStorageByLogID(logID int64) {
	if logID <= 0 {
		return
	}
	_, err := msql.Model(`work_flow_storage_cache`, define.Postgres).
		Where(`log_id`, cast.ToString(logID)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
}

func WorkFlowRestore(flow *WorkFlow) (err error) {
	restoreDialogId, restoreSessionId := getWorkFlowStorageIDs(flow)
	defer func() {
		DelWorkFlowStorage(restoreDialogId, restoreSessionId)
	}()
	logs.Debug(`to restore，dialog_id %v session_id %v`, restoreDialogId, restoreSessionId)
	data, _, logID := GetWorkFlowRestore(restoreDialogId, restoreSessionId)
	if data == nil || data.CurNodeKey == `` {
		return
	}
	// Restore reuses the original log identity to avoid inserting another log on completion.
	if logID > 0 {
		flow.LogID = logID
		if GetWorkFlowLogStatus(logID) != define.WorkFlowStatusRunning {
			flow.isStopped = true
			return errors.New(i18n.Show(flow.params.Lang, `status_exception`))
		}
	}
	//start params
	startNode, _, err := GetNodeByKey(flow, cast.ToUint(flow.params.RealRobot[`id`]), flow.curNodeKey)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	_, _, err = startNode.Running(flow)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//fill
	flow.nodeLogs = data.NodeLogs
	flow.outputs = data.Outputs
	flow.curNodeKey = data.CurNodeKey
	flow.runNodeKeys = data.RunNodeKeys
	flow.runLogs = data.RunLogs
	if _, exist := flow.outputs[flow.curNodeKey]; !exist {
		return
	}
	// Sync the current user input into global.question before matching a
	// smart-menu option. The value assembled from the restored context (or the
	// stale global carried by the request) must not be used to pick a menu
	// branch; otherwise the reply never hits the option the user clicked.
	if len(flow.params.Question) > 0 {
		q := flow.params.Question
		flow.global[`question`] = common.SimpleField{
			Key:  `question`,
			Typ:  common.TypString,
			Vals: []common.Val{{String: &q}},
		}
	}
	menuQuestion := ``
	//last node output
	if question, exist := flow.global[`question`]; exist {
		flow.outputs[flow.curNodeKey][`question`] = common.SimpleField{
			Sys:  question.Sys,
			Key:  strings.TrimPrefix(question.Key, `global.`),
			Typ:  question.Typ,
			Vals: question.Vals,
		}
		if len(question.Vals) > 0 {
			menuQuestion = *question.Vals[0].String
		}
	}
	if question, exist := flow.global[`question_multiple`]; exist {
		flow.outputs[flow.curNodeKey][`question_multiple`] = common.SimpleField{
			Sys:  question.Sys,
			Key:  strings.TrimPrefix(question.Key, `global.`),
			Typ:  question.Typ,
			Vals: question.Vals,
		}
	}

	flow.output = flow.outputs[flow.curNodeKey]
	for globalKey, globalVal := range data.Global {
		flow.global[globalKey] = globalVal
	}
	if len(flow.nodeLogs) > 0 {
		nodeLog := flow.nodeLogs[len(flow.nodeLogs)-1]
		nodeLog.EndTime = time.Now().UnixMilli()
		nodeLog.Output = common.GetFieldsObject(common.GetRecurveFields(flow.output))
		nodeLog.UseTime = nodeLog.EndTime - nodeLog.StartTime
		flow.nodeLogs[len(flow.nodeLogs)-1] = nodeLog
	}
	var nodeInfo msql.Params
	if flow.params.Draft.IsDraft {
		nodeInfo = flow.params.Draft.NodeMaps[flow.curNodeKey]
	} else {
		nodeInfo, err = common.GetRobotNode(cast.ToUint(flow.params.Robot[`id`]), flow.curNodeKey)
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(flow.params.Lang, "workflow_node_get_fail"))
		}
	}
	nodeParams := NodeParams{}
	err = tool.JsonDecode(nodeInfo[`node_params`], &nodeParams)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(flow.params.Lang, "workflow_node_params_get_fail"))
	}
	defaultNextNodeKey := ``
	if nodeParams.Question.AnswerType == define.QuestionAnswerTypeMenu {
		flow.curNodeKey = ``
		for _, replyContent := range nodeParams.Question.ReplyContentList {
			if replyContent.ReplyType != common.ReplyTypeSmartMenu {
				continue
			}
			if len(replyContent.SmartMenu.MenuContent) == 0 {
				continue
			}
			for menuIndex, menuContent := range replyContent.SmartMenu.MenuContent {
				if menuContent.MenuType == `-1` {
					defaultNextNodeKey = menuContent.NextNodeKey
					continue
				}
				replaceMenuContent := flow.VariableReplace(menuContent.Content)
				if tool.InArray(menuQuestion, []string{replaceMenuContent, fmt.Sprintf(`%d %s`, menuIndex+1, replaceMenuContent)}) {
					flow.curNodeKey = menuContent.NextNodeKey
					break
				}
			}
		}
		if flow.curNodeKey == `` {
			flow.curNodeKey = defaultNextNodeKey
		}
	} else {
		flow.curNodeKey = nodeInfo[`next_node_key`]
	}
	return
}
