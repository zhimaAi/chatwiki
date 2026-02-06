// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/pipeline/biz_chat"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func GetWsUrl(c *gin.Context) {
	openid := strings.TrimSpace(c.Query(`openid`))
	if !common.IsChatOpenid(openid) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))))
		return
	}
	var wsUrl string
	if cast.ToBool(define.Config.WebService[`ws_use_ssl`]) {
		wsUrl = fmt.Sprintf(`wss://%s/ws?openid=%s`, define.Config.WebService[`ws_domain`], openid)
	} else {
		wsUrl = fmt.Sprintf(`ws://%s/ws?openid=%s`, define.Config.WebService[`ws_domain`], openid)
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`ws_url`: wsUrl}, nil))
	//sending debug message
	if cast.ToInt(c.Query(`debug`)) > 0 {
		//debug text push
		message := msql.Datas{
			`admin_user_id`: 1,
			`robot_id`:      1,
			`openid`:        openid,
			`dialogue_id`:   1,
			`session_id`:    1,
			`is_customer`:   define.MsgFromRobot,
			`msg_type`:      define.MsgTypeText,
			`content`:       `websocket text push ...`,
			`menu_json`:     ``,
			`quote_file`:    `[]`,
			`create_time`:   tool.Time2Int() + 5,
			`update_time`:   tool.Time2Int() + 5,
		}
		messageStr, _ := tool.JsonEncode(map[string]any{"openid": openid, "message": message})
		_ = common.AddJobs(lib_define.WsMessagePushTopic, messageStr, time.Second*5)
		//debug menu push
		message[`msg_type`] = define.MsgTypeMenu
		message[`content`] = `[menu]`
		message[`menu_json`] = `{"content":"menu_content","question":["question_1", "question_2"]}`
		message[`create_time`] = tool.Time2Int() + 10
		message[`update_time`] = tool.Time2Int() + 10
		messageStr, _ = tool.JsonEncode(map[string]any{"openid": openid, "message": message})
		_ = common.AddJobs(lib_define.WsMessagePushTopic, messageStr, time.Second*10)
		//debug image push
		message[`msg_type`] = define.MsgTypeImage
		message[`content`] = define.LocalUploadPrefix + `default/robot_avatar.svg`
		message[`menu_json`] = ``
		message[`create_time`] = tool.Time2Int() + 15
		message[`update_time`] = tool.Time2Int() + 15
		messageStr, _ = tool.JsonEncode(map[string]any{"openid": openid, "message": message})
		_ = common.AddJobs(lib_define.WsMessagePushTopic, messageStr, time.Second*15)
	}
}

func IsOnLine(c *gin.Context) {
	openid := strings.TrimSpace(c.Query(`openid`))
	if !common.IsChatOpenid(openid) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))))
		return
	}
	var link string
	if cast.ToBool(define.Config.WebService[`ws_use_ssl`]) {
		link = fmt.Sprintf(`https://%s/isOnLine`, define.Config.WebService[`ws_domain`])
	} else {
		link = fmt.Sprintf(`http://%s/isOnLine`, define.Config.WebService[`ws_domain`])
	}
	response, err := curl.Get(link).Param(`openid`, openid).String()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, response)
}

func ChatMessage(c *gin.Context) {
	chatBaseParam, err := common.CheckChatRequest(c, true)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//get params
	dialogueId := cast.ToUint(c.PostForm(`dialogue_id`))
	minId := cast.ToUint(c.PostForm(`min_id`))
	size := max(1, cast.ToInt(c.PostForm(`size`)))
	m := msql.Model(`chat_ai_message`, define.Postgres).
		Alias(`m`).
		Join(`message_feedback f`, `m.id=f.ai_message_id`, `left`).
		Where(`m.openid`, chatBaseParam.Openid).Where(`m.robot_id`, chatBaseParam.Robot[`id`]).
		Field(`m.*,case when f.type=1 then 1 when f.type=2 then 2 else 0 end as feedback_type`)
	if dialogueId > 0 {
		m.Where(`m.dialogue_id`, cast.ToString(dialogueId))
	}
	if minId > 0 {
		m.Where(`m.id`, `<`, cast.ToString(minId))
	}
	list, err := m.Limit(size).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`robot`: chatBaseParam.Robot, `customer`: chatBaseParam.Customer, `list`: list}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func AddChatMessageFeedback(c *gin.Context) {

	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	aiMessageId := cast.ToInt(c.PostForm(`ai_message_id`))
	customerMessageId := cast.ToInt(c.PostForm(`customer_message_id`))
	_type := cast.ToInt(c.PostForm(`type`))
	content := strings.TrimSpace(c.DefaultPostForm(`content`, ``))
	if aiMessageId <= 0 || customerMessageId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if _type != 1 && _type != 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `type`))))
		return
	}

	lockKey := define.LockPreKey + `SaveMessageFeedback` + cast.ToString(aiMessageId)
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer func(lockKey string) {
		lib_redis.UnLock(define.Redis, lockKey)
	}(lockKey)

	aiMessage, err := msql.Model(`chat_ai_message`, define.Postgres).Where(`is_customer`, `0`).Where(`id`, cast.ToString(aiMessageId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(aiMessage) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	customerMessage, err := msql.Model(`chat_ai_message`, define.Postgres).Where(`is_customer`, `1`).Where(`id`, cast.ToString(customerMessageId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(customerMessage) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	row, err := msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`ai_message_id`, cast.ToString(aiMessageId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(row) > 0 {
		if row[`type`] == cast.ToString(_type) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
			return
		} else {
			_, err = msql.Model(`message_feedback`, define.Postgres).Where(`id`, row[`id`]).Delete()
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
		}
	}

	modelConfig, err := common.GetModelConfigInfo(cast.ToInt(chatBaseParam.Robot[`model_config_id`]), chatBaseParam.AdminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	robotInfo := chatBaseParam.Robot
	robotInfo[`corp_name`] = common.GetModelNameByDefine(common.GetLang(c), modelConfig[`model_define`])
	robotJson, err := tool.JsonEncode(robotInfo)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	_, err = msql.Model(`message_feedback`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`:       chatBaseParam.AdminUserId,
		`robot_id`:            chatBaseParam.Robot[`id`],
		`ai_message_id`:       aiMessageId,
		`customer_message_id`: customerMessageId,
		`type`:                _type,
		`robot`:               robotJson,
		`content`:             content,
		`create_time`:         tool.Time2Int(),
		`update_time`:         tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DelChatMessageFeedback(c *gin.Context) {
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	aiMessageId := cast.ToInt(c.PostForm(`ai_message_id`))
	if aiMessageId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	message, err := msql.Model(`chat_ai_message`, define.Postgres).Where(`is_customer`, `0`).Where(`id`, cast.ToString(aiMessageId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(message) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	exists, err := msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`ai_message_id`, cast.ToString(aiMessageId)).
		Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if exists == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
		return
	}
	_, err = msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`robot_id`, chatBaseParam.Robot[`id`]).
		Where(`ai_message_id`, cast.ToString(aiMessageId)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func saveCustomerInfo(c *gin.Context, chatBaseParam *define.ChatBaseParam) {
	nickname := strings.TrimSpace(c.DefaultPostForm(`nickname`, c.Query(`nickname`)))
	name := strings.TrimSpace(c.DefaultPostForm(`name`, c.Query(`name`)))
	avatar := strings.TrimSpace(c.DefaultPostForm(`avatar`, c.Query(`avatar`)))
	upData := msql.Datas{}
	if len(chatBaseParam.Customer) == 0 || chatBaseParam.Customer[`nickname`] != nickname {
		upData[`nickname`] = nickname
	}
	if len(chatBaseParam.Customer) == 0 || chatBaseParam.Customer[`name`] != name {
		upData[`name`] = name
	}
	if len(chatBaseParam.Customer) == 0 || chatBaseParam.Customer[`avatar`] != avatar {
		upData[`avatar`] = avatar
	}
	if len(chatBaseParam.Customer) == 0 && cast.ToInt(c.PostForm(`is_background`)) > 0 {
		upData[`is_background`] = 1 //background create
	}
	common.InsertOrUpdateCustomer(chatBaseParam.Openid, chatBaseParam.AdminUserId, upData)
}

func ChatWelcome(c *gin.Context) {
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	dialogueId := cast.ToInt(c.PostForm(`dialogue_id`))
	//database dispose
	saveCustomerInfo(c, chatBaseParam)
	chatBaseParam.Customer, _ = common.GetCustomerInfo(chatBaseParam.Openid, chatBaseParam.AdminUserId)
	//build message
	message := msql.Datas{
		`admin_user_id`: chatBaseParam.AdminUserId,
		`robot_id`:      chatBaseParam.Robot[`id`],
		`openid`:        chatBaseParam.Openid,
		`is_customer`:   define.MsgFromRobot,
		`msg_type`:      define.MsgTypeMenu,
		`content`:       i18n.Show(common.GetLang(c), `welcomes`),
		`menu_json`:     chatBaseParam.Robot[`welcomes`],
		`quote_file`:    `[]`,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	dialogId, err := common.GetDialogueIdNoCreate(chatBaseParam)
	if err != nil {
		logs.Error(err.Error())
	}
	sessionId, err := common.GetSessionIdNoCreate(dialogId)
	if err != nil {
		logs.Error(err.Error())
	}
	data := map[string]any{
		`message`:       common.ToStringMap(message),
		`robot`:         chatBaseParam.Robot,
		`customer`:      chatBaseParam.Customer,
		`dialog_id`:     dialogId,
		`session_id`:    sessionId,
		`chat_variable`: manage.GetChatRobotVariables(dialogueId, chatBaseParam),
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func ChatRequest(c *gin.Context) {
	//preinitialize:c.Stream can close body,future get c.PostForm exception
	_ = c.Request.ParseMultipartForm(define.DefaultMultipartMemory)
	c.Header(`Content-Type`, `text/event-stream`)
	c.Header(`Cache-Control`, `no-cache`)
	c.Header(`Connection`, `keep-alive`)
	if define.IsDev {
		c.Header(`Access-Control-Allow-Origin`, `*`)
	}
	params := getChatRequestParam(c)
	chanStream := make(chan sse.Event)
	go func() {
		_, _ = DoChatRequest(params, true, chanStream)
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
	*params.IsClose = true // set flag
	for range chanStream {
		// discard unpushed data flows
	}
}

func ChatRequestNotStream(c *gin.Context) {
	//preinitialize:c.Stream can close body,future get c.PostForm exception
	_ = c.Request.ParseMultipartForm(define.DefaultMultipartMemory)
	c.Header(`Connection`, `keep-alive`)
	if define.IsDev {
		c.Header(`Access-Control-Allow-Origin`, `*`)
	}
	params := getChatRequestParam(c)
	// Execute non-streaming chat request
	chanStream := make(chan sse.Event)
	go func() {
		for range chanStream {
			// Consume stream events to prevent blocking
		}
	}()
	out, err := DoChatRequest(params, false, chanStream)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotInfo, err := common.GetRobotInfo(c.PostForm(`robot_key`))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToInt(robotInfo[`is_default`]) == define.NotDefault {
		_ = common.SetStepFinish(cast.ToInt(robotInfo[`admin_user_id`]), define.StepTestRobot)
	}
	c.String(http.StatusOK, lib_web.FmtJson(out, nil))
}

type QuestionGuideMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ChatQuestionGuide(c *gin.Context) {
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	dialogId := cast.ToInt(c.PostForm(`dialogue_id`))
	if dialogId == 0 {
		logs.Error(`dialogue_id is empty`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_empty`, `dialogue_id`))))
		return
	}

	recentMessages, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`openid`, chatBaseParam.Openid).
		Where(`robot_id`, chatBaseParam.Robot[`id`]).
		Where(`dialogue_id`, cast.ToString(dialogId)).
		Where(`msg_type`, cast.ToString(define.MsgTypeText)).
		Limit(8).
		Order(`id desc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(recentMessages) == 0 {
		logs.Error(`recentMessages is empty`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	sort.Slice(recentMessages, func(i, j int) bool {
		return cast.ToInt(recentMessages[i][`id`]) < cast.ToInt(recentMessages[j][`id`])
	})

	var questionGuideMessages []QuestionGuideMessage
	for _, one := range recentMessages {
		if len(one[`content`]) == 0 {
			continue
		}
		if cast.ToInt(one[`is_customer`]) == define.MsgFromCustomer {
			questionGuideMessages = append(questionGuideMessages, QuestionGuideMessage{Role: `user`, Content: one[`content`]})
		} else {
			questionGuideMessages = append(questionGuideMessages, QuestionGuideMessage{Role: `assistant`, Content: one[`content`]})
		}
	}

	if cast.ToBool(chatBaseParam.Robot[`enable_question_guide`]) == false {
		logs.Error(`enable_question_guide is closed`)
		c.String(http.StatusOK, lib_web.FmtJson([]string{}, nil))
		return
	}

	histories := ""
	for _, msg := range questionGuideMessages {
		if msg.Role == `user` {
			histories += "Q: " + msg.Content
		} else {
			histories += "A: " + msg.Content
		}
	}
	prompt := strings.ReplaceAll(define.PromptDefaultQuestionGuide, `{{num}}`, chatBaseParam.Robot[`question_guide_num`])
	prompt = strings.ReplaceAll(prompt, `{{histories}}`, histories)
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `user`, Content: prompt}}

	chatResp, _, err := common.RequestChat(
		common.GetLang(c),
		chatBaseParam.AdminUserId,
		chatBaseParam.Openid,
		chatBaseParam.Robot,
		chatBaseParam.AppType,
		cast.ToInt(chatBaseParam.Robot[`model_config_id`]),
		chatBaseParam.Robot[`use_model`],
		messages,
		nil,
		cast.ToFloat32(chatBaseParam.Robot[`temperature`]),
		200,
	)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson([]string{}, nil))
		return
	}
	content := chatResp.Result
	var guides []string
	err = json.Unmarshal([]byte(content), &guides)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson([]string{}, nil))
		return
	}
	if len(guides) >= cast.ToInt(chatBaseParam.Robot[`question_guide_num`]) {
		guides = guides[:cast.ToInt(chatBaseParam.Robot[`question_guide_num`])]
	}
	c.String(http.StatusOK, lib_web.FmtJson(guides, nil))
}

func getChatRequestParam(c *gin.Context) *define.ChatRequestParam {
	chatBaseParam, err := common.CheckChatRequest(c)
	isClose := false
	workFlowGlobal := make(map[string]any)
	_ = tool.JsonDecodeUseNumber(c.DefaultPostForm(`global`, `{}`), &workFlowGlobal)
	question := strings.TrimSpace(c.PostForm(`question`))
	openid := strings.TrimSpace(c.PostForm(`openid`))
	loopTestParams := work_flow.TakeTestParams(question, openid, c.DefaultPostForm(`loop_test_params`, `[]`), &workFlowGlobal)
	batchTestParams := work_flow.TakeTestParams(question, openid, c.DefaultPostForm(`batch_test_params`, `[]`), &workFlowGlobal)
	testParams := work_flow.TakeTestParams(question, openid, c.DefaultPostForm(`test_params`, `[]`), &workFlowGlobal)
	return &define.ChatRequestParam{
		ChatBaseParam:       chatBaseParam,
		Error:               err,
		Lang:                common.GetLang(c),
		Question:            question,
		DialogueId:          cast.ToInt(c.PostForm(`dialogue_id`)),
		Prompt:              strings.TrimSpace(c.PostForm(`prompt`)),
		LibraryIds:          strings.TrimSpace(c.PostForm(`library_ids`)),
		IsClose:             &isClose,
		WorkFlowGlobal:      workFlowGlobal,
		LoopTestParams:      loopTestParams,
		BatchTestParams:     batchTestParams,
		ChatPromptVariables: c.PostForm(`chat_prompt_variables`),
		TestParams:          testParams,
	}
}

func DoChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event) (msql.Params, error) {
	out := biz_chat.DoChatRequest(params, useStream, chanStream)
	return out.AiMessage, out.Error
}

func GetDialogueSession(params *define.ChatRequestParam) (int, int, error) {
	if len(params.Openid) == 0 {
		return 0, 0, nil // from workflow, openid is empty, no dialog id or session id generated
	}
	var err error
	dialogueId := params.DialogueId
	if dialogueId > 0 {
		dialogue, err := common.GetDialogueInfo(dialogueId, params.AdminUserId, cast.ToInt(params.Robot[`id`]), params.Openid)
		if err != nil {
			logs.Error(err.Error())
			return 0, 0, err
		}
		if len(dialogue) == 0 {
			return 0, 0, errors.New(i18n.Show(params.Lang, `param_invalid`, `dialogue_id`))
		}
	} else {
		dialogueId, err = common.GetDialogueId(params.ChatBaseParam, params.Question)
		if err != nil {
			logs.Error(err.Error())
			return 0, 0, err
		}
	}
	sessionId, err := common.GetSessionId(params, dialogueId)
	if err != nil {
		logs.Error(err.Error())
		return 0, 0, err
	}
	return dialogueId, sessionId, nil
}

func CallWorkFlow(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
	}
	chatRequestParam.HeaderToken = c.GetHeader(`token`)
	if len(chatRequestParam.HeaderToken) == 0 {
		chatRequestParam.HeaderToken = c.Query(`token`)
	}
	dialogueId, sessionId, err := GetDialogueSession(chatRequestParam)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if dialogueId == 0 {
		dialogueId = cast.ToInt(c.PostForm(`dialogue_id`))
	}
	if sessionId == 0 {
		sessionId = cast.ToInt(c.PostForm(`session_id`))
	}
	isDraft := cast.ToBool(c.PostForm(`is_draft`))
	questionMultipleSwitch := cast.ToBool(c.PostForm(`question_multiple_switch`))
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam: chatRequestParam,
		DialogueId:       dialogueId,
		SessionId:        sessionId,
		Draft:            work_flow.Draft{IsDraft: isDraft, QuestionMultipleSwitch: questionMultipleSwitch},
	}
	if isDraft { // run test (draft)
		workFlowParams.Draft.NodeMaps, err = msql.Model(`work_flow_node`, define.Postgres).
			Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).Where(`robot_id`, chatRequestParam.Robot[`id`]).
			Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
		nodeList := make([]work_flow.WorkFlowNode, 0)
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
		if workFlowParams.Draft.StartNodeKey, _, _, _, err, errNodeKey = work_flow.VerifyWorkFlowNodes(nodeList, chatRequestParam.AdminUserId, chatRequestParam.Lang); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
				`err_node_key`: errNodeKey,
			}, err))
			return
		}
	}
	_, nodeLogs, err := work_flow.BaseCallWorkFlow(workFlowParams)
	useToken, useMills := common.TakeWorkFlowTestUseToken(nodeLogs)
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`node_logs`:  nodeLogs,
		`use_token`:  useToken,
		`use_mills`:  useMills,
		`dialog_id`:  workFlowParams.DialogueId,
		`session_id`: workFlowParams.SessionId,
	}, err))
}

func CallLoopWorkFlow(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
	}
	chatRequestParam.HeaderToken = c.GetHeader(`token`)
	if len(chatRequestParam.HeaderToken) == 0 {
		chatRequestParam.HeaderToken = c.Query(`token`)
	}
	dialogueId, sessionId, err := GetDialogueSession(chatRequestParam)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam:  chatRequestParam,
		DialogueId:        dialogueId,
		SessionId:         sessionId,
		Draft:             work_flow.Draft{IsDraft: true},
		IsTestLoopNodeRun: true,
	}
	loopNodeKey := c.PostForm(`loop_node_key`)
	if loopNodeKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(chatRequestParam.Lang, `param_invalid`, `loop_node_key`))))
		return
	}
	loopNodeInfo, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).
		Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`node_key`, loopNodeKey).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if cast.ToInt(loopNodeInfo[`node_type`]) != work_flow.NodeTypeLoop {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(chatRequestParam.Lang, `no_data`))))
		return
	}
	//build loop node
	loopWorkFlowNode := work_flow.WorkFlowNode{
		NodeType:     common.MixedInt(cast.ToInt(loopNodeInfo[`node_type`])),
		NodeName:     loopNodeInfo[`node_name`],
		NodeKey:      loopNodeInfo[`node_key`],
		NodeParams:   work_flow.NodeParams{},
		NodeInfoJson: make(map[string]any),
	}
	_ = tool.JsonDecodeUseNumber(loopNodeInfo[`node_params`], &loopWorkFlowNode.NodeParams)
	_ = tool.JsonDecodeUseNumber(loopNodeInfo[`node_info_json`], &loopWorkFlowNode.NodeInfoJson)
	//query child node
	workFlowParams.Draft.NodeMaps, err = msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).
		Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`loop_parent_key`, loopNodeKey).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
	nodeList := make([]work_flow.WorkFlowNode, 0)
	//build child node
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
	workFlowParams.Draft.NodeMaps[loopNodeKey] = loopNodeInfo
	if workFlowParams.Draft.StartNodeKey, _, _, err = work_flow.VerityLoopWorkflowNodes(chatRequestParam.AdminUserId, loopWorkFlowNode, nodeList, work_flow.LoopAllowNodeTypes, lib_define.CircularNode, chatRequestParam.Lang); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	_, nodeLogs, err := work_flow.BaseCallWorkFlow(workFlowParams)
	useToken, useMills := common.TakeWorkFlowTestUseToken(nodeLogs)
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`node_logs`: nodeLogs,
		`use_token`: useToken,
		`use_mills`: useMills,
	}, err))
}

func CallBatchWorkFlow(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
	}
	chatRequestParam.HeaderToken = c.GetHeader(`token`)
	if len(chatRequestParam.HeaderToken) == 0 {
		chatRequestParam.HeaderToken = c.Query(`token`)
	}
	dialogueId, sessionId, err := GetDialogueSession(chatRequestParam)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam:   chatRequestParam,
		DialogueId:         dialogueId,
		SessionId:          sessionId,
		Draft:              work_flow.Draft{IsDraft: true},
		IsTestBatchNodeRun: true,
	}
	batchNodeKey := c.PostForm(`batch_node_key`)
	if batchNodeKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(chatRequestParam.Lang, `param_invalid`, `loop_node_key`))))
		return
	}
	batchNodeInfo, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).
		Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`node_key`, batchNodeKey).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if cast.ToInt(batchNodeInfo[`node_type`]) != work_flow.NodeTypeBatch {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(chatRequestParam.Lang, `no_data`))))
		return
	}
	//build batch node
	batchWorkFlowNode := work_flow.WorkFlowNode{
		NodeType:     common.MixedInt(cast.ToInt(batchNodeInfo[`node_type`])),
		NodeName:     batchNodeInfo[`node_name`],
		NodeKey:      batchNodeInfo[`node_key`],
		NodeParams:   work_flow.NodeParams{},
		NodeInfoJson: make(map[string]any),
	}
	_ = tool.JsonDecodeUseNumber(batchNodeInfo[`node_params`], &batchWorkFlowNode.NodeParams)
	_ = tool.JsonDecodeUseNumber(batchNodeInfo[`node_info_json`], &batchWorkFlowNode.NodeInfoJson)
	//query child node
	workFlowParams.Draft.NodeMaps, err = msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).
		Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`loop_parent_key`, batchNodeKey).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
	nodeList := make([]work_flow.WorkFlowNode, 0)
	//build child node
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
	workFlowParams.Draft.NodeMaps[batchNodeKey] = batchNodeInfo
	if workFlowParams.Draft.StartNodeKey, _, _, err = work_flow.VerityLoopWorkflowNodes(chatRequestParam.AdminUserId, batchWorkFlowNode, nodeList, work_flow.BatchAllowNodeTypes, lib_define.BatchProcessing, chatRequestParam.Lang); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	_, nodeLogs, err := work_flow.BaseCallWorkFlow(workFlowParams)
	useToken, useMills := common.TakeWorkFlowTestUseToken(nodeLogs)
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`node_logs`: nodeLogs,
		`use_token`: useToken,
		`use_mills`: useMills,
	}, err))
}

func CallLoopWorkFlowParams(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(chatBaseParam.Robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	loopNodeKey := c.PostForm(`loop_node_key`)
	if loopNodeKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `loop_node_key`))))
		return
	}
	var loopNode *work_flow.WorkFlowNode
	nodeMaps, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`robot_id`, chatBaseParam.Robot[`id`]).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
	nodeList := make([]work_flow.WorkFlowNode, 0)
	for _, params := range nodeMaps {
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
		if node.NodeKey == loopNodeKey {
			loopNode = &node
		}
		nodeList = append(nodeList, node)
	}
	if loopNode == nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	// loop array
	loopTestParams := make([]common.LoopTestParams, 0)
	if loopNode.NodeParams.Loop.LoopType == common.LoopTypeArray {
		for _, loopField := range loopNode.NodeParams.Loop.LoopArrays {
			findNode := work_flow.FindNodeByUseKey(nodeList, loopField.Value)
			if findNode == nil {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
				return
			}
			loopTestParams = append(loopTestParams, common.LoopTestParams{
				NodeKey:  findNode.NodeKey,
				NodeName: findNode.NodeName,
				Field: common.SimpleField{
					Sys:      false,
					Key:      loopField.Value,
					Typ:      loopField.Typ,
					Vals:     make([]common.Val, 0),
					Required: true,
				},
			})
		}
	}
	childNodes := make([]work_flow.WorkFlowNode, 0)
	for _, node := range nodeList {
		if node.LoopParentKey == loopNode.NodeKey {
			childNodes = append(childNodes, node)
		}
	}
	// intermediate variable
	for _, intermediateField := range loopNode.NodeParams.Loop.IntermediateParams {
		if intermediateField.Key == `` {
			continue
		}
		isUse := work_flow.FindKeyIsUse(childNodes, loopNode.NodeKey+`.`+intermediateField.Key)
		if !isUse {
			continue
		}
		loopTestParams = append(loopTestParams, common.LoopTestParams{
			NodeKey:  loopNode.NodeKey,
			NodeName: loopNode.NodeName,
			Field: common.SimpleField{
				Sys:      false,
				Key:      intermediateField.Key,
				Typ:      intermediateField.Typ,
				Vals:     make([]common.Val, 0),
				Required: false,
			},
		})
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`loop_test_params`: loopTestParams,
		`is_need_question`: work_flow.FindKeyIsUse(childNodes, `global.question`), // whether question is needed
		`is_need_openid`:   work_flow.FindKeyIsUse(childNodes, `global.openid`),   // whether openid is needed
	}, err))
}

func CallBatchWorkFlowParams(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatBaseParam, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(chatBaseParam.Robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	batchNodeKey := c.PostForm(`batch_node_key`)
	if batchNodeKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `batch_node_key`))))
		return
	}
	var batchNode *work_flow.WorkFlowNode
	nodeMaps, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(chatBaseParam.AdminUserId)).
		Where(`robot_id`, chatBaseParam.Robot[`id`]).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
	nodeList := make([]work_flow.WorkFlowNode, 0)
	for _, params := range nodeMaps {
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
		if node.NodeKey == batchNodeKey {
			batchNode = &node
		}
		nodeList = append(nodeList, node)
	}
	if batchNode == nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	// loop array
	batchTestParams := make([]common.BatchTestParams, 0)
	for _, batchField := range batchNode.NodeParams.Batch.BatchArrays {
		findNode := work_flow.FindNodeByUseKey(nodeList, batchField.Value)
		if findNode == nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		batchTestParams = append(batchTestParams, common.BatchTestParams{
			NodeKey:  findNode.NodeKey,
			NodeName: findNode.NodeName,
			Field: common.SimpleField{
				Sys:      false,
				Key:      batchField.Value,
				Typ:      batchField.Typ,
				Vals:     make([]common.Val, 0),
				Required: true,
			},
		})
	}
	childNodes := make([]work_flow.WorkFlowNode, 0)
	for _, node := range nodeList {
		if node.LoopParentKey == batchNode.NodeKey {
			childNodes = append(childNodes, node)
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`loop_test_params`: batchTestParams,
		`is_need_question`: work_flow.FindKeyIsUse(childNodes, `global.question`), // whether question is needed
		`is_need_openid`:   work_flow.FindKeyIsUse(childNodes, `global.openid`),   // whether openid is needed
	}, err))
}

func CallWorkFlowHttpTest(c *gin.Context) {
	c.Set(`from_work_flow`, true) // set flag, do not validate openid as empty
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
	}
	chatRequestParam.HeaderToken = c.GetHeader(`token`)
	if len(chatRequestParam.HeaderToken) == 0 {
		chatRequestParam.HeaderToken = c.Query(`token`)
	}
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam: chatRequestParam,
	}
	var curlNodeKey = c.PostForm(`curl_node_key`)
	if curlNodeKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `curl_node_key`))))
		return
	}
	var err error
	workFlowParams.Draft.NodeMaps, err = msql.Model(`work_flow_node`, define.Postgres).
		Where(`admin_user_id`, chatRequestParam.Robot[`admin_user_id`]).Where(`robot_id`, chatRequestParam.Robot[`id`]).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).ColumnMap(`*`, `node_key`)
	nodeList := make([]work_flow.WorkFlowNode, 0)
	curlNode := &work_flow.WorkFlowNode{}
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
		if node.NodeKey == curlNodeKey {
			curlNode = &node
		}
	}
	//Verify the workflow nodes
	var errNodeKey string
	workFlowParams.Draft.StartNodeKey, _, _, _, err, errNodeKey = work_flow.VerifyWorkFlowNodes(nodeList, chatRequestParam.AdminUserId, chatRequestParam.Lang)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
			`err_node_key`: errNodeKey,
		}, err))
		return
	}
	ret, err := work_flow.CallHttpTest(workFlowParams, curlNode)
	c.String(http.StatusOK, lib_web.FmtJson(ret, err))
}

func ChatEditVariables(c *gin.Context) {
	_, err := common.CheckChatRequest(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	dialogId := cast.ToInt(c.PostForm(`dialogue_id`))
	if dialogId == 0 {
		logs.Error(`dialogue_id is empty`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_empty`, `dialogue_id`))))
		return
	}
	sessionId := cast.ToInt(c.PostForm(`session_id`))
	if sessionId == 0 {
		logs.Error(`session_id is empty`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_empty`, `session_id`))))
		return
	}
	chatPromptVariables := cast.ToString(c.PostForm(`chat_prompt_variables`))
	if len(chatPromptVariables) == 0 {
		logs.Error(`chat_prompt_variables is empty`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_empty`, `chat_prompt_variables`))))
		return
	}
	upData := msql.Datas{
		`chat_prompt_variables`: chatPromptVariables,
	}
	common.UpChatPromptVariables(dialogId, sessionId, upData)
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
