// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/custom_eino"
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
)

func GetPromptOptimizeRobot(lang string, adminUserID, userID, robotID int) (msql.Params, error) {
	m := msql.Model("chat_ai_robot", define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("id", cast.ToString(robotID))
	robot, err := m.Find()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, "sys_err"))
	}
	if len(robot) == 0 || !tool.InArrayInt(cast.ToInt(robot["application_type"]), []int{define.ApplicationTypeChat, define.ApplicationTypeClaw}) {
		return nil, errors.New(i18n.Show(lang, "prompt_optimize_robot_unavailable"))
	}
	if !CheckObjectAccessRights(adminUserID, userID, define.IdentityTypeUser, robotID, define.ObjectTypeRobot, define.PermissionEditRights, lang) {
		return nil, errors.New(i18n.Show(lang, "auth_no_permission"))
	}
	return robot, nil
}

func promptOptimizeNames(lang string, robot msql.Params) ([]define.PromptOptimizeName, error) {
	names := make([]define.PromptOptimizeName, 0)
	if robot["work_flow_ids"] != "" {
		m := msql.Model("chat_ai_robot", define.Postgres).Where("admin_user_id", robot["admin_user_id"]).Where("id", "in", robot["work_flow_ids"]).Where("application_type", cast.ToString(define.ApplicationTypeFlow)).Field("robot_key,en_name,robot_name,robot_intro")
		rows, err := m.Select()
		if err != nil {
			logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
			return nil, errors.New(i18n.Show(lang, "sys_err"))
		}
		for _, r := range rows {
			name := r["en_name"]
			if name == "" {
				name = "work_flow_#_" + r["robot_key"]
			}
			names = append(names, define.PromptOptimizeName{Name: name, Label: r["robot_name"], Type: "Tool", Description: r["robot_intro"]})
		}
	}
	if cast.ToInt(robot["application_type"]) == define.ApplicationTypeClaw {
		if !cast.ToBool(robot["search_knowledge_close"]) {
			names = append(names, define.PromptOptimizeName{Name: custom_eino.KbsearchToolName, Label: i18n.Show(lang, "prompt_optimize_knowledge_tool"), Type: "Tool"})
		}
		if !cast.ToBool(robot["query_local_docs_close"]) {
			names = append(names, define.PromptOptimizeName{Name: define.SkillReservedName, Label: i18n.Show(lang, "prompt_optimize_local_docs_skill"), Type: "Skill"})
		}
		// Only the cloud edition stores this flag; an absent value stays disabled.
		if cast.ToBool(robot["open_agent_manual_tool"]) {
			names = append(names, define.PromptOptimizeName{Name: custom_eino.ManualTransferToolName, Label: i18n.Show(lang, "prompt_optimize_manual_tool"), Type: "Tool"})
		}
		if cast.ToBool(robot["open_agent_write_file_tool"]) {
			names = append(names, define.PromptOptimizeName{Name: "write_file", Label: "write_file", Type: "Tool"})
		}
		if cast.ToBool(robot["open_agent_edit_file_tool"]) {
			names = append(names, define.PromptOptimizeName{Name: "edit_file", Label: "edit_file", Type: "Tool"})
		}
		if cast.ToBool(robot["open_agent_execute_tool"]) {
			names = append(names, define.PromptOptimizeName{Name: "execute", Label: "execute", Type: "Tool"})
		}
		if cast.ToBool(robot["goods_lib_recommend_switch"]) {

			names = append(names, define.PromptOptimizeName{Name: custom_eino.GoodsLibRecommendToolName, Label: i18n.Show(lang, "prompt_optimize_goods_tool"), Type: "Tool"})
		}
		m := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).Where("admin_user_id", robot["admin_user_id"]).Where("robot_key", robot["robot_key"]).Field("skill_name,remark_name,description")
		skills, err := m.Select()
		if err != nil {
			logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
			return nil, errors.New(i18n.Show(lang, "sys_err"))
		}
		for _, s := range skills {
			names = append(names, define.PromptOptimizeName{Name: s["skill_name"], Label: s["remark_name"], Type: "Skill", Description: s["description"]})
		}
	}
	seen := make(map[string]bool)
	uniqueNames := make([]define.PromptOptimizeName, 0, len(names))
	for _, name := range names {
		key := name.Type + ":" + name.Name
		if name.Name == "" || seen[key] {
			continue
		}
		seen[key] = true
		uniqueNames = append(uniqueNames, name)
	}
	return uniqueNames, nil
}

func GetPromptOptimizeContext(lang string, adminUserID, userID int, params define.PromptOptimizeRobotParams) (define.PromptOptimizeContext, error) {
	robot, err := GetPromptOptimizeRobot(lang, adminUserID, userID, params.RobotID)
	if err != nil {
		return define.PromptOptimizeContext{}, err
	}
	names, err := promptOptimizeNames(lang, robot)
	if err != nil {
		return define.PromptOptimizeContext{}, err
	}
	return define.PromptOptimizeContext{
		Prompt:          BuildPromptStruct(lang, cast.ToInt(robot["prompt_type"]), robot["prompt"], robot["prompt_struct"]),
		ApplicationType: cast.ToInt(robot["application_type"]), ModelConfigID: robot["model_config_id"], UseModel: robot["use_model"], Names: names,
	}, nil
}

func GetPromptOptimizeRecords(lang string, adminUserID, userID int, params define.PromptOptimizeRecordsParams) (define.PromptOptimizeRecordList, error) {
	_, err := GetPromptOptimizeRobot(lang, adminUserID, userID, params.RobotID)
	if err != nil {
		return define.PromptOptimizeRecordList{}, err
	}
	size := min(define.PromptOptimizePageSizeLimit, max(1, params.Size))
	m := msql.Model(define.TablePromptOptimizeRecord, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID))
	if params.BeforeID > 0 {
		m.Where("id", "<", cast.ToString(params.BeforeID))
	}
	rows, err := m.Field("id,demand,scene,mode,summary,reply,reasoning_content,version_id,status,error,create_time,(NULLIF(input_context,'')::jsonb->'capabilities')::text AS capabilities,char_length(NULLIF(input_context,'')::jsonb->>'base_prompt') AS base_prompt_length").Order("id desc").Limit(size).Select()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return define.PromptOptimizeRecordList{}, errors.New(i18n.Show(lang, "sys_err"))
	}
	return define.PromptOptimizeRecordList{List: rows, HasMore: len(rows) == size}, nil
}

func GetPromptOptimizeVersions(lang string, adminUserID, userID int, params define.PromptOptimizeRobotParams) ([]msql.Params, error) {
	_, err := GetPromptOptimizeRobot(lang, adminUserID, userID, params.RobotID)
	if err != nil {
		return nil, err
	}
	m := msql.Model(define.TablePromptOptimizeVersion, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID))
	rows, err := m.Field("id,version,record_id,create_time").Order("version desc").Select()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, "sys_err"))
	}
	return rows, nil
}

func GetPromptOptimizeVersion(lang string, adminUserID, userID int, params define.PromptOptimizeVersionParams) (msql.Params, error) {
	_, err := GetPromptOptimizeRobot(lang, adminUserID, userID, params.RobotID)
	if err != nil {
		return nil, err
	}
	m := msql.Model(define.TablePromptOptimizeVersion, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID)).Where("id", cast.ToString(params.VersionID))
	row, err := m.Field("id,version,record_id,prompt,create_time").Find()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, "sys_err"))
	}
	if len(row) == 0 {
		return nil, errors.New(i18n.Show(lang, "prompt_optimize_version_unavailable"))
	}
	return row, nil
}

func PreparePromptOptimize(lang string, adminUserID, userID int, params define.PromptOptimizeGenerateParams) (define.PromptOptimizePrepared, error) {
	prepared := define.PromptOptimizePrepared{}
	params.Demand = strings.TrimSpace(params.Demand)
	params.RequestID = strings.TrimSpace(params.RequestID)
	params.ReasoningEffort = strings.TrimSpace(params.ReasoningEffort)
	if err := ValidateReasoningEffort(lang, params.ReasoningEffort, "reasoning_effort"); err != nil {
		return prepared, err
	}
	if params.ReasoningEffort == "" {
		params.ReasoningEffort = string(chat.ReasoningEffortNone)
	}
	if params.Demand == "" {
		return prepared, errors.New(i18n.Show(lang, "param_invalid", "demand"))
	}
	if params.RequestID == "" || utf8.RuneCountInString(params.RequestID) > define.PromptOptimizeRequestIDLength {
		return prepared, errors.New(i18n.Show(lang, "param_invalid", "request_id"))
	}
	if params.Scene == define.PromptOptimizeSceneConversation && params.MessageID <= 0 {
		return prepared, errors.New(i18n.Show(lang, "param_invalid", "message_id"))
	}
	system, err := PromptOptimizeSystemPrompt(lang, params.Scene, params.Mode)
	if err != nil {
		return prepared, err
	}
	robot, err := GetPromptOptimizeRobot(lang, adminUserID, userID, params.RobotID)
	if err != nil {
		return prepared, err
	}
	m := msql.Model(define.TablePromptOptimizeRecord, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID)).Where("request_id", params.RequestID)
	old, err := m.Find()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return prepared, errors.New(i18n.Show(lang, "sys_err"))
	}
	if len(old) > 0 {
		return prepared, errors.New(i18n.Show(lang, "prompt_optimize_request_exists"))
	}
	base := BuildPromptStruct(lang, cast.ToInt(robot["prompt_type"]), robot["prompt"], robot["prompt_struct"])
	if params.BaseVersionID > 0 {
		m = msql.Model(define.TablePromptOptimizeVersion, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID)).Where("id", cast.ToString(params.BaseVersionID))
		version, e := m.Find()
		if e != nil {
			logs.Error("sql:%s,err:%s", m.GetLastSql(), e.Error())
			return prepared, errors.New(i18n.Show(lang, "sys_err"))
		}
		if len(version) == 0 {
			return prepared, errors.New(i18n.Show(lang, "prompt_optimize_version_unavailable"))
		}
		base = version["prompt"]
	}
	names, err := promptOptimizeNames(lang, robot)
	if err != nil {
		return prepared, err
	}
	input := define.PromptOptimizeInput{Request: params.Demand, BasePrompt: base, Capabilities: names}
	if params.MessageID > 0 {
		m = msql.Model("chat_ai_message", define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID)).Where("id", cast.ToString(params.MessageID)).Where("is_customer", "0")
		answer, e := m.Find()
		if e != nil {
			logs.Error("sql:%s,err:%s", m.GetLastSql(), e.Error())
			return prepared, errors.New(i18n.Show(lang, "sys_err"))
		}
		if len(answer) == 0 || strings.TrimSpace(answer["content"]) == "" {
			return prepared, errors.New(i18n.Show(lang, "prompt_optimize_message_unavailable"))
		}
		input.Answer = answer["content"]
		customerMessageID := cast.ToInt64(answer["customer_message_id"])
		// Keep the selected reply even when legacy data has no linked question.
		if customerMessageID > 0 {
			m = msql.Model("chat_ai_message", define.Postgres).Where("admin_user_id", cast.ToString(adminUserID)).Where("robot_id", cast.ToString(params.RobotID)).Where("dialogue_id", answer["dialogue_id"]).Where("session_id", answer["session_id"]).Where("openid", answer["openid"]).Where("is_customer", "1").Where("id", cast.ToString(customerMessageID))
			question, e := m.Find()
			if e != nil {
				logs.Error("sql:%s,err:%s", m.GetLastSql(), e.Error())
				return prepared, errors.New(i18n.Show(lang, "sys_err"))
			}
			if len(question) > 0 && strings.TrimSpace(question["content"]) != "" {
				input.Question = question["content"]
			}
		}
	}
	inputText, err := tool.JsonEncode(input)
	if err != nil {
		logs.Error("encode prompt optimization input: %s", err.Error())
		return prepared, errors.New(i18n.Show(lang, "sys_err"))
	}
	now := time.Now()
	m = msql.Model(define.TablePromptOptimizeRecord, define.Postgres).Alias("r").
		Join(define.TablePromptOptimizeVersion+" v", "v.id=r.version_id AND v.admin_user_id=r.admin_user_id AND v.robot_id=r.robot_id", "left").
		Where("r.admin_user_id", cast.ToString(adminUserID)).Where("r.robot_id", cast.ToString(params.RobotID)).Where("r.status", define.PromptOptimizeStatusSuccess).
		Where("r.create_time", ">=", cast.ToString(now.Add(-define.PromptOptimizeHistoryWindow).Unix())).
		Field("r.id,r.create_time,r.status,r.input_context,r.reply,v.prompt").Order("r.create_time desc,r.id desc").Limit(define.PromptOptimizeHistoryRounds)
	history, err := m.Select()
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return prepared, errors.New(i18n.Show(lang, "sys_err"))
	}
	messages := []chat.Message{{Role: chat.RoleSystem, Content: chat.TextContent(system)}}
	rounds := make([]define.PromptOptimizeRound, 0, len(history))
	for _, h := range history {
		rounds = append(rounds, define.PromptOptimizeRound{ID: cast.ToInt64(h["id"]), CreatedAt: cast.ToInt64(h["create_time"]), Status: h["status"], Input: h["input_context"], Prompt: h["prompt"], Reply: h["reply"]})
	}
	historyRounds := 0
	for _, h := range SelectPromptOptimizeHistory(rounds, now) {
		reply := promptOptimizeHistoryReply(h)
		if strings.TrimSpace(reply) == "" {
			continue
		}
		messages = append(messages, chat.Message{Role: chat.RoleUser, Content: chat.TextContent(h.Input)}, chat.Message{Role: chat.RoleAssistant, Content: chat.TextContent(reply)})
		historyRounds++
	}
	messages = append(messages, chat.Message{Role: chat.RoleUser, Content: chat.TextContent(inputText)})
	createTime := time.Now().Unix()
	m = msql.Model(define.TablePromptOptimizeRecord, define.Postgres)
	id, err := m.Insert(msql.Datas{"admin_user_id": adminUserID, "robot_id": params.RobotID, "request_id": params.RequestID, "demand": params.Demand, "scene": params.Scene, "mode": params.Mode, "model_config_id": params.ModelConfigID, "use_model": params.UseModel, "base_version_id": params.BaseVersionID, "message_id": params.MessageID, "input_context": inputText, "status": define.PromptOptimizeStatusPending, "create_time": createTime, "update_time": createTime}, "id")
	if err != nil {
		logs.Error("sql:%s,err:%s", m.GetLastSql(), err.Error())
		return prepared, errors.New(i18n.Show(lang, "sys_err"))
	}
	return define.PromptOptimizePrepared{Robot: robot, Params: params, RecordID: id, Messages: messages, Names: names, HistoryRounds: historyRounds, BasePromptLength: utf8.RuneCountInString(base)}, nil
}

// Lock only during version allocation and persistence, never during model inference.
func savePromptOptimizeResult(ctx context.Context, lang string, prepared define.PromptOptimizePrepared, modelOutput define.PromptOptimizeOutput, reasoningContent string) (define.PromptOptimizeDone, error) {
	result := define.PromptOptimizeDone{RecordID: prepared.RecordID, Reply: modelOutput.Reply, ReasoningContent: reasoningContent}
	db, err := msql.GetDB(define.Postgres)
	if err != nil {
		logs.Error("prompt optimization database: %s", err.Error())
		return result, errors.New(i18n.Show(lang, "sys_err"))
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logs.Error("sql:BEGIN prompt optimization,err:%s", err.Error())
		return result, errors.New(i18n.Show(lang, "sys_err"))
	}
	defer tx.Rollback()
	robotID, adminID := prepared.Robot["id"], prepared.Robot["admin_user_id"]
	now := time.Now().Unix()
	var versionID int64
	var summary string
	if modelOutput.Prompt != "" {
		sqlText := "SELECT id FROM chat_ai_robot WHERE id=$1 AND admin_user_id=$2 FOR UPDATE"
		var locked int
		if err = tx.QueryRowContext(ctx, sqlText, robotID, adminID).Scan(&locked); err != nil {
			logs.Error("sql:%s,err:%s", sqlText, err.Error())
			return result, errors.New(i18n.Show(lang, "sys_err"))
		}
		sqlText = "SELECT COALESCE(MAX(version),0)+1 FROM " + define.TablePromptOptimizeVersion + " WHERE robot_id=$1 AND admin_user_id=$2"
		var version int
		if err = tx.QueryRowContext(ctx, sqlText, robotID, adminID).Scan(&version); err != nil {
			logs.Error("sql:%s,err:%s", sqlText, err.Error())
			return result, errors.New(i18n.Show(lang, "sys_err"))
		}
		sqlText = "INSERT INTO " + define.TablePromptOptimizeVersion + " (admin_user_id,robot_id,record_id,version,prompt,create_time) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		if err = tx.QueryRowContext(ctx, sqlText, adminID, robotID, prepared.RecordID, version, modelOutput.Prompt, now).Scan(&versionID); err != nil {
			logs.Error("sql:%s,err:%s", sqlText, err.Error())
			return result, errors.New(i18n.Show(lang, "sys_err"))
		}
		result.Version = &define.PromptOptimizeVersion{ID: versionID, Version: version, RecordID: prepared.RecordID, Prompt: modelOutput.Prompt, CreateTime: now}
		summary = i18n.Show(lang, "prompt_optimize_generated_version", version)
	}
	sqlText := "UPDATE " + define.TablePromptOptimizeRecord + " SET status=$1,version_id=$2,update_time=$3,summary=$4,reply=$5,reasoning_content=$6 WHERE id=$7 AND admin_user_id=$8 AND robot_id=$9 AND status=$10"
	updated, err := tx.ExecContext(ctx, sqlText, define.PromptOptimizeStatusSuccess, versionID, now, summary, result.Reply, result.ReasoningContent, prepared.RecordID, adminID, robotID, define.PromptOptimizeStatusPending)
	if err != nil {
		logs.Error("sql:%s,err:%s", sqlText, err.Error())
		return result, errors.New(i18n.Show(lang, "sys_err"))
	}
	if affected, affectedErr := updated.RowsAffected(); affectedErr != nil || affected != 1 {
		logs.Error("sql:%s,rows:%d,err:%v", sqlText, affected, affectedErr)
		return result, errors.New(i18n.Show(lang, "sys_err"))
	}
	if err = tx.Commit(); err != nil {
		logs.Error("sql:COMMIT prompt optimization,err:%s", err.Error())
		return result, errors.New(i18n.Show(lang, "sys_err"))
	}
	return result, nil
}

// RunPromptOptimize owns model inference, record state and stream lifecycle.
// The controller only forwards the returned business events to the HTTP client.
func RunPromptOptimize(parent context.Context, lang string, prepared define.PromptOptimizePrepared, output chan sse.Event) {
	defer close(output)
	ctx, cancel := context.WithTimeout(parent, define.PromptOptimizeRequestTimeout)
	defer cancel()
	sendOrAbort(ctx, output, sse.Event{Event: "prepared", Data: map[string]any{"record_id": prepared.RecordID, "history_rounds": prepared.HistoryRounds, "names": prepared.Names, "base_prompt_length": prepared.BasePromptLength}})
	events := make(chan sse.Event, 16)
	modelDone := make(chan error, 1)
	go func() {
		defer close(events)
		// Enable the existing reasoning event only for this request. Keep it outside business robot accounting and fallback.
		modelRobot := msql.Params{"think_switch": cast.ToString(define.SwitchOn)}
		_, _, err := RequestChatStream(ctx, lang, cast.ToInt(prepared.Robot["admin_user_id"]), prepared.Robot["admin_user_id"], modelRobot, "", prepared.Params.ModelConfigID, prepared.Params.UseModel, prepared.Messages, nil, events, define.PromptOptimizeTemperature, define.PromptOptimizeMaxTokens, chat.ReasoningEffort(prepared.Params.ReasoningEffort))
		modelDone <- err
	}()
	var parser promptOptimizeOutputParser
	var reasoning strings.Builder
	valid := true
	for event := range events {
		text, ok := event.Data.(string)
		if !ok || !valid || ctx.Err() != nil {
			continue
		}
		switch event.Event {
		case "reasoning_content":
			reasoning.WriteString(text)
			sendOrAbort(ctx, output, sse.Event{Event: "reasoning_content", Data: define.PromptOptimizeDelta{Text: text}})
		case "sending":
			var promptDelta, replyDelta string
			promptDelta, replyDelta, valid = parser.push(text)
			if !valid {
				cancel()
				continue
			}
			if promptDelta != "" {
				sendOrAbort(ctx, output, sse.Event{Event: "prompt_delta", Data: define.PromptOptimizeDelta{Text: promptDelta}})
			}
			if replyDelta != "" {
				sendOrAbort(ctx, output, sse.Event{Event: "reply_delta", Data: define.PromptOptimizeDelta{Text: replyDelta}})
			}
		}
	}
	requestErr := <-modelDone
	var result define.PromptOptimizeDone
	var err error
	switch {
	case parent.Err() != nil:
		err = PromptOptimizeError(lang, parent.Err())
	case !valid:
		err = errors.New(i18n.Show(lang, "prompt_optimize_invalid_output"))
	case ctx.Err() != nil:
		err = PromptOptimizeError(lang, ctx.Err())
	case requestErr != nil:
		logs.Error("prompt optimization model request,robot:%s,record:%d,err:%s", prepared.Robot["id"], prepared.RecordID, requestErr.Error())
		err = PromptOptimizeError(lang, requestErr)
	default:
		modelOutput, complete := parser.finish()
		if !complete {
			err = errors.New(i18n.Show(lang, "prompt_optimize_invalid_output"))
		} else {
			result, err = savePromptOptimizeResult(ctx, lang, prepared, modelOutput, reasoning.String())
			if err != nil && ctx.Err() != nil {
				err = PromptOptimizeError(lang, ctx.Err())
			}
		}
	}
	if err != nil {
		status := define.PromptOptimizeStatusFailed
		if errors.Is(parent.Err(), context.Canceled) {
			status = define.PromptOptimizeStatusCancelled
		}
		m := msql.Model(define.TablePromptOptimizeRecord, define.Postgres).Where("admin_user_id", prepared.Robot["admin_user_id"]).Where("robot_id", prepared.Robot["id"]).Where("id", cast.ToString(prepared.RecordID))
		m.Where("status", define.PromptOptimizeStatusPending)
		if _, updateErr := m.Update(msql.Datas{"status": status, "error": err.Error(), "update_time": time.Now().Unix()}); updateErr != nil {
			logs.Error("sql:%s,err:%s", m.GetLastSql(), updateErr.Error())
		}
		sendOrAbort(parent, output, sse.Event{Event: "error", Data: map[string]string{"message": err.Error()}})
		return
	}
	sendOrAbort(parent, output, sse.Event{Event: "done", Data: result})
}
