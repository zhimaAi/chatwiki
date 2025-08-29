// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"

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
	chatBaseParam, err := common.CheckChatRequest(c)
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

	// add corp name field to robot info
	var corpName string
	for _, modelInfo := range common.GetModelList() {
		if len(modelConfig[`model_define`]) == 0 || modelInfo.ModelDefine == modelConfig[`model_define`] {
			corpName = modelInfo.ModelName
		}
	}
	robotInfo := chatBaseParam.Robot
	robotInfo[`corp_name`] = corpName
	if len(modelConfig[`deployment_name`]) > 0 {
		robotInfo[`use_model`] = modelConfig[`deployment_name`]
	}
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
	data := map[string]any{
		`message`:  common.ToStringMap(message),
		`robot`:    chatBaseParam.Robot,
		`customer`: chatBaseParam.Customer,
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
	*params.IsClose = true //set flag
	for range chanStream {
		//discard unpushed data flows
	}
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
	return &define.ChatRequestParam{
		ChatBaseParam:  chatBaseParam,
		Error:          err,
		Lang:           common.GetLang(c),
		Question:       strings.TrimSpace(c.PostForm(`question`)),
		DialogueId:     cast.ToInt(c.PostForm(`dialogue_id`)),
		Prompt:         strings.TrimSpace(c.PostForm(`prompt`)),
		LibraryIds:     strings.TrimSpace(c.PostForm(`library_ids`)),
		IsClose:        &isClose,
		WorkFlowGlobal: workFlowGlobal,
	}
}

func DoChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event) (msql.Params, error) {
	monitor := common.NewMonitor(params)
	message, err := doChatRequest(params, useStream, chanStream, monitor)
	monitor.Save(err)
	return message, err
}

func doChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event, monitor *common.Monitor) (msql.Params, error) {
	defer close(chanStream)
	chanStream <- sse.Event{Event: `ping`, Data: tool.Time2Int()}
	//check params
	if params.Error != nil {
		chanStream <- sse.Event{Event: `error`, Data: params.Error.Error()}
		return nil, params.Error
	}
	if len(params.Question) == 0 {
		err := errors.New(i18n.Show(params.Lang, `question_empty`))
		chanStream <- sse.Event{Event: `error`, Data: err.Error()}
		return nil, err
	}
	//get dialogue_id and session_id
	var err error
	dialogueId := params.DialogueId
	if dialogueId > 0 {
		dialogue, err := common.GetDialogueInfo(dialogueId, params.AdminUserId, cast.ToInt(params.Robot[`id`]), params.Openid)
		if err != nil {
			logs.Error(err.Error())
			chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
			return nil, err
		}
		if len(dialogue) == 0 {
			err := errors.New(i18n.Show(params.Lang, `param_invalid`, `dialogue_id`))
			chanStream <- sse.Event{Event: `error`, Data: err}
			return nil, err
		}
	} else {
		dialogueId, err = common.GetDialogueId(params.ChatBaseParam, params.Question)
		if err != nil {
			logs.Error(err.Error())
			chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
			return nil, err
		}
	}
	sessionId, err := common.GetSessionId(params, dialogueId)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	chanStream <- sse.Event{Event: `dialogue_id`, Data: dialogueId}
	chanStream <- sse.Event{Event: `session_id`, Data: sessionId}
	//customer push
	customer, err := common.GetCustomerInfo(params.Openid, params.AdminUserId)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	chanStream <- sse.Event{Event: `customer`, Data: customer}
	//database dispose
	message := msql.Datas{
		`admin_user_id`: params.AdminUserId,
		`robot_id`:      params.Robot[`id`],
		`openid`:        params.Openid,
		`dialogue_id`:   dialogueId,
		`session_id`:    sessionId,
		`is_customer`:   define.MsgFromCustomer,
		`msg_type`:      define.MsgTypeText,
		`content`:       params.Question,
		`menu_json`:     ``,
		`quote_file`:    `[]`,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	if len(customer) > 0 {
		message[`nickname`] = customer[`nickname`]
		message[`name`] = customer[`name`]
		message[`avatar`] = customer[`avatar`]
	}
	lastChat := msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromCustomer)
	//message push
	chanStream <- sse.Event{Event: `c_message`, Data: common.ToStringMap(message, `id`, id)}
	//websocket notify
	common.ReceiverChangeNotify(params.AdminUserId, `c_message`, common.ToStringMap(message, `id`, id))
	//obtain the data required for gpt
	chanStream <- sse.Event{Event: `robot`, Data: params.Robot}
	debugLog := make([]any, 0) //debug log
	defer func() {
		monitor.DebugLog = debugLog //记录监控数据
	}()
	var messages []adaptor.ZhimaChatCompletionMessage
	var list []msql.Params
	var (
		showQuoteFile  = cast.ToBool(params.Robot[`answer_source_switch`])
		startQuoteFile bool
	)
	if cast.ToInt(params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		//nothing to do
	} else if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeDirect {
		messages, list, err = buildDirectChatRequestMessage(params, id, dialogueId, &debugLog)
	} else {
		if showQuoteFile {
			chanStream <- sse.Event{Event: `start_quote_file`, Data: tool.Time2Int()}
			startQuoteFile = true
		}
		messages, list, monitor.LibUseTime, err = buildLibraryChatRequestMessage(params, id, dialogueId, &debugLog)

		chanStream <- sse.Event{Event: `recall_time`, Data: monitor.LibUseTime.RecallTime}
	}

	if err != nil {
		logs.Error(err.Error())
		if startQuoteFile {
			chanStream <- sse.Event{Event: `quote_file`, Data: `[]`}
		}
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	//save answer source
	var fileSourceMap = make(map[string][]msql.Datas)
	if len(list) > 0 {
		for _, one := range list {
			var images []string
			if err := tool.JsonDecode(one[`images`], &images); err != nil {
				logs.Error(err.Error())
			}
			fileSourceMap[one[`file_id`]] = append(fileSourceMap[one[`file_id`]], msql.Datas{
				`admin_user_id`: params.AdminUserId,
				`file_id`:       one[`file_id`],
				`paragraph_id`:  one[`id`],
				`word_total`:    one[`word_total`],
				`similarity`:    one[`similarity`],
				`title`:         one[`title`],
				`type`:          one[`type`],
				`content`:       one[`content`],
				`question`:      one[`question`],
				`answer`:        one[`answer`],
				`images`:        images,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
		}
	}
	messages = common.BuildOpenApiContent(params, messages)
	//dispose answer source and quote_file
	quoteFile, ms := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	for _, one := range list {
		if _, ok := ms[one[`file_id`]]; ok {
			continue //remove duplication
		}
		library, err := common.GetLibraryInfo(cast.ToInt(one[`library_id`]), cast.ToInt(one[`admin_user_id`]))
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		ms[one[`file_id`]] = struct{}{}
		quoteFileForSave = append(quoteFileForSave, msql.Params{
			`id`:           one[`file_id`],
			`library_id`:   library[`id`],
			`library_name`: library[`library_name`],
			`file_name`:    one[`file_name`],
		})
		quoteFile = append(quoteFile, msql.Params{
			`id`:                 one[`file_id`],
			`library_id`:         library[`id`],
			`library_name`:       library[`library_name`],
			`file_name`:          one[`file_name`],
			`answer_source_data`: tool.JsonEncodeNoError(fileSourceMap[one[`file_id`]]),
		})
	}
	if showQuoteFile && startQuoteFile {
		chanStream <- sse.Event{Event: `quote_file`, Data: quoteFile}
	}
	quoteFileJson, _ := tool.JsonEncode(quoteFileForSave)
	var functionTools []adaptor.FunctionTool
	if len(params.Robot[`form_ids`]) > 0 {
		formIdList := strings.Split(params.Robot[`form_ids`], `,`)
		functionTools, err = common.BuildFunctionTools(formIdList, params.AdminUserId)
		if err != nil {
			logs.Error(err.Error())
			chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
			return nil, err
		}
	}
	//聊天机器人支持关联工作流
	workFlowFuncCall, needRunWorkFlow := work_flow.BuildFunctionTools(params.Robot)
	if needRunWorkFlow {
		functionTools = append(functionTools, workFlowFuncCall...)
	}

	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
		llmStartTime                        = time.Now()
	)
	msgType := define.MsgTypeText
	if cast.ToInt(params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		workFlowParams := &work_flow.WorkFlowParams{ChatRequestParam: params, CurMsgId: int(id), DialogueId: dialogueId, SessionId: sessionId}
		content, requestTime, monitor.LibUseTime, list, err = work_flow.CallWorkFlow(workFlowParams, &debugLog, monitor)
		if err != nil {
			sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
			debugLog = append(debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
			chanStream <- sse.Event{Event: `debug`, Data: debugLog} //渲染Prompt日志
			return nil, err
		}
		chanStream <- sse.Event{Event: `recall_time`, Data: monitor.LibUseTime.RecallTime}
		chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
		chanStream <- sse.Event{Event: `sending`, Data: content}
	} else if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeDirect {
		if !needRunWorkFlow && useStream {
			chatResp, requestTime, err = common.RequestChatStream(
				params.AdminUserId,
				params.Openid,
				params.Robot,
				params.AppType,
				cast.ToInt(params.Robot[`model_config_id`]),
				params.Robot[`use_model`],
				messages,
				functionTools,
				chanStream,
				cast.ToFloat32(params.Robot[`temperature`]),
				cast.ToInt(params.Robot[`max_token`]),
			)
		} else {
			chatResp, requestTime, err = common.RequestChat(
				params.AdminUserId,
				params.Openid,
				params.Robot,
				params.AppType,
				cast.ToInt(params.Robot[`model_config_id`]),
				params.Robot[`use_model`],
				messages,
				functionTools,
				cast.ToFloat32(params.Robot[`temperature`]),
				cast.ToInt(params.Robot[`max_token`]),
			)
		}
		content = chatResp.Result
		reasoningContent = chatResp.ReasoningContent
		if err != nil {
			logs.Error(err.Error())
			sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
		}
	} else if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeMixture {
		if len(list) == 0 {
			if !needRunWorkFlow && useStream {
				chatResp, requestTime, err = common.RequestChatStream(
					params.AdminUserId,
					params.Openid,
					params.Robot,
					params.AppType,
					cast.ToInt(params.Robot[`model_config_id`]),
					params.Robot[`use_model`],
					messages,
					functionTools,
					chanStream,
					cast.ToFloat32(params.Robot[`temperature`]),
					cast.ToInt(params.Robot[`max_token`]),
				)
			} else {
				chatResp, requestTime, err = common.RequestChat(
					params.AdminUserId,
					params.Openid,
					params.Robot,
					params.AppType,
					cast.ToInt(params.Robot[`model_config_id`]),
					params.Robot[`use_model`],
					messages,
					functionTools,
					cast.ToFloat32(params.Robot[`temperature`]),
					cast.ToInt(params.Robot[`max_token`]),
				)
			}
			content = chatResp.Result
			reasoningContent = chatResp.ReasoningContent
			if err != nil {
				logs.Error(err.Error())
				sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
			}
		} else {
			if cast.ToBool(params.Robot[`mixture_qa_direct_reply_switch`]) &&
				cast.ToInt(list[0][`type`]) != define.ParagraphTypeNormal &&
				len(list[0][`similarity`]) > 0 &&
				cast.ToFloat32(list[0][`similarity`]) >= cast.ToFloat32(params.Robot[`mixture_qa_direct_reply_score`]) {
				content = list[0][`answer`]

				if len(list[0][`images`]) > 0 {
					var imageSlice []string
					image_decode_err := json.Unmarshal([]byte(cast.ToString(list[0][`images`])), &imageSlice)
					if image_decode_err == nil {
						// 成功解码JSON到切片
						for _, imageUrl := range imageSlice {
							content += `<br/>![img](` + imageUrl + `)`
						}
					}
				}

				chanStream <- sse.Event{Event: `sending`, Data: content}
			} else {
				if !needRunWorkFlow && useStream {
					chatResp, requestTime, err = common.RequestChatStream(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
						functionTools,
						chanStream,
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				} else {
					chatResp, requestTime, err = common.RequestChat(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
						functionTools,
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				}
				content = chatResp.Result
				reasoningContent = chatResp.ReasoningContent
				if err != nil {
					logs.Error(err.Error())
					sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
				}
			}
		}
	} else {
		if len(list) == 0 {
			unknownQuestionPrompt := define.MenuJsonStruct{}
			_ = tool.JsonDecodeUseNumber(params.Robot[`unknown_question_prompt`], &unknownQuestionPrompt)
			if len(unknownQuestionPrompt.Content) == 0 && len(unknownQuestionPrompt.Question) == 0 {
				sendDefaultUnknownQuestionPrompt(params, `unknown_question_prompt not config`, chanStream, &content)
			} else {
				msgType = define.MsgTypeMenu
				content = unknownQuestionPrompt.Content
				menuJson, _ = tool.JsonEncode(unknownQuestionPrompt)
			}
		} else {
			// direct answer
			if cast.ToBool(params.Robot[`library_qa_direct_reply_switch`]) &&
				cast.ToInt(list[0][`type`]) != define.ParagraphTypeNormal &&
				len(list[0][`similarity`]) > 0 &&
				cast.ToFloat32(list[0][`similarity`]) >= cast.ToFloat32(params.Robot[`library_qa_direct_reply_score`]) {
				content = list[0][`answer`]

				if len(list[0][`images`]) > 0 {
					var imageSlice []string
					image_decode_err := json.Unmarshal([]byte(cast.ToString(list[0][`images`])), &imageSlice)
					if image_decode_err == nil {
						// 成功解码JSON到切片
						for _, imageUrl := range imageSlice {
							content += `<br/>![img](` + imageUrl + `)`
						}
					}
				}

				chanStream <- sse.Event{Event: `sending`, Data: content}
			} else { // ask gpt
				if !needRunWorkFlow && useStream {
					chatResp, requestTime, err = common.RequestChatStream(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
						functionTools,
						chanStream,
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				} else {
					chatResp, requestTime, err = common.RequestChat(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
						functionTools,
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				}
				content = chatResp.Result
				reasoningContent = chatResp.ReasoningContent
				if err != nil {
					logs.Error(err.Error())
					sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
				}
			}
		}
	}
	//聊天机器人支持关联工作流
	if err == nil && needRunWorkFlow {
		workFlowRobot, workFlowGlobal := work_flow.ChooseWorkFlowRobot(chatResp.FunctionToolCalls)
		if len(workFlowRobot) == 0 { //大模型没有返回需要调用的工作流
			chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
			chanStream <- sse.Event{Event: `sending`, Data: content}
		} else { //组装工作流请求参数,并执行工作流
			workFlowParams := work_flow.BuildWorkFlowParams(*params, workFlowRobot, workFlowGlobal, int(id), dialogueId, sessionId)
			content, requestTime, _, _, err = work_flow.CallWorkFlow(workFlowParams, &debugLog, monitor)
			if err == nil {
				chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
				chanStream <- sse.Event{Event: `sending`, Data: content}
			} else {
				logs.Error(err.Error())
				sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
			}
		}
	}

	//记录监控数据
	monitor.LlmCallTime = time.Now().Sub(llmStartTime).Milliseconds()
	monitor.RequestTime, monitor.Error = requestTime, err

	if *params.IsClose { //client break
		return nil, errors.New(`client break`)
	}
	debugLog = append(debugLog, map[string]string{`type`: `cur_answer`, `content`: content})
	//push prompt log
	chanStream <- sse.Event{Event: `debug`, Data: debugLog}
	//dispose answer source
	if cast.ToInt(params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		fileSourceMap = map[string][]msql.Datas{}
		for _, one := range list {
			var images []string
			if err := tool.JsonDecode(one[`images`], &images); err != nil {
				logs.Error(err.Error())
			}
			fileSourceMap[one[`file_id`]] = append(fileSourceMap[one[`file_id`]], msql.Datas{
				`admin_user_id`: params.AdminUserId,
				`file_id`:       one[`file_id`],
				`paragraph_id`:  one[`id`],
				`word_total`:    one[`word_total`],
				`similarity`:    one[`similarity`],
				`title`:         one[`title`],
				`type`:          one[`type`],
				`content`:       one[`content`],
				`question`:      one[`question`],
				`answer`:        one[`answer`],
				`images`:        images,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
		}
		for _, one := range list {
			if _, ok := ms[one[`file_id`]]; ok {
				continue //remove duplication
			}
			library, err := common.GetLibraryInfo(cast.ToInt(one[`library_id`]), cast.ToInt(one[`admin_user_id`]))
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			ms[one[`file_id`]] = struct{}{}
			quoteFile = append(quoteFile, msql.Params{
				`id`:                 one[`file_id`],
				`file_name`:          one[`file_name`],
				`library_id`:         one[`library_id`],
				`library_name`:       library[`library_name`],
				`answer_source_data`: tool.JsonEncodeNoError(fileSourceMap[one[`file_id`]]),
			})
			quoteFileForSave = append(quoteFileForSave, msql.Params{
				`id`:           one[`file_id`],
				`file_name`:    one[`file_name`],
				`library_id`:   one[`library_id`],
				`library_name`: library[`library_name`],
			})
		}
		if len(quoteFile) > 0 && showQuoteFile {
			chanStream <- sse.Event{Event: `quote_file`, Data: quoteFile}
		}
		quoteFileJson, _ = tool.JsonEncode(quoteFileForSave)
	}
	//database dispose
	message = msql.Datas{
		`admin_user_id`:          params.AdminUserId,
		`robot_id`:               params.Robot[`id`],
		`openid`:                 params.Openid,
		`dialogue_id`:            dialogueId,
		`session_id`:             sessionId,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           requestTime,
		`recall_time`:            monitor.LibUseTime.RecallTime,
		`msg_type`:               msgType,
		`content`:                content,
		`reasoning_content`:      reasoningContent,
		`is_valid_function_call`: chatResp.IsValidFunctionCall,
		`menu_json`:              menuJson,
		`quote_file`:             quoteFileJson,
		`create_time`:            tool.Time2Int(),
		`update_time`:            tool.Time2Int(),
	}
	if len(params.Robot) > 0 {
		message[`nickname`] = `` //none
		message[`name`] = params.Robot[`robot_name`]
		message[`avatar`] = params.Robot[`robot_avatar`]
	}
	lastChat = msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err = msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromRobot)
	//message push
	chanStream <- sse.Event{Event: `ai_message`, Data: common.ToStringMap(message, `id`, id)}
	//websocket notify
	common.ReceiverChangeNotify(params.AdminUserId, `ai_message`, common.ToStringMap(message, `id`, id))
	//save answer source only workflow
	if len(list) > 0 {
		asm := msql.Model(`chat_ai_answer_source`, define.Postgres)
		for _, one := range list {
			_, err = asm.Insert(msql.Datas{
				`admin_user_id`: params.AdminUserId,
				`message_id`:    id,
				`file_id`:       one[`file_id`],
				`paragraph_id`:  one[`id`],
				`word_total`:    one[`word_total`],
				`similarity`:    one[`similarity`],
				`title`:         one[`title`],
				`type`:          one[`type`],
				`content`:       one[`content`],
				`question`:      one[`question`],
				`answer`:        one[`answer`],
				`images`:        one[`images`],
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
			if err != nil {
				logs.Error(`sql:%s,err:%s`, asm.GetLastSql(), err.Error())
			}
		}
	}
	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
	message["use_model"] = params.Robot["use_model"]
	chanStream <- sse.Event{Event: `data`, Data: message}
	chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
	return common.ToStringMap(message, `id`, id), nil
}

func sendDefaultUnknownQuestionPrompt(params *define.ChatRequestParam, errmsg string, chanStream chan sse.Event, content *string) {
	chanStream <- sse.Event{Event: `error`, Data: `SYSERR:` + errmsg}
	code := `unknown`
	if ms := regexp.MustCompile(`ERROR\s+CODE:\s?(.*)`).FindStringSubmatch(errmsg); len(ms) > 1 {
		code = ms[1]
	}
	*content = i18n.Show(params.Lang, `gpt_error`, code)
	chanStream <- sse.Event{Event: `sending`, Data: *content}
}

func buildLibraryChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, []msql.Params, common.LibUseTime, error) {
	if len(params.Prompt) == 0 { //no custom is used
		params.Prompt = common.BuildPromptStruct(cast.ToInt(params.Robot[`prompt_type`]), params.Robot[`prompt`], params.Robot[`prompt_struct`])
	}
	if len(params.LibraryIds) == 0 || !common.CheckIds(params.LibraryIds) { //no custom is used
		params.LibraryIds = params.Robot[`library_ids`]
	}

	contextList := common.BuildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))

	//question optimize
	var questionopTime int64
	var optimizedQuestions []string
	if cast.ToBool(params.Robot[`enable_question_optimize`]) && len(params.LibraryIds) > 0 {
		var err error
		temp := time.Now()
		optimizedQuestions, err = common.GetOptimizedQuestions(params, contextList)
		questionopTime = time.Now().Sub(temp).Milliseconds()
		if err != nil {
			logs.Error(err.Error())
		}
	}

	//convert match
	list, libUseTime, err := common.GetMatchLibraryParagraphList(
		params.Openid,
		params.AppType,
		params.Question,
		optimizedQuestions,
		params.LibraryIds,
		cast.ToInt(params.Robot[`top_k`]),
		cast.ToFloat64(params.Robot[`similarity`]),
		cast.ToInt(params.Robot[`search_type`]),
		params.Robot,
	)
	libUseTime.QuestionOp = questionopTime
	if err != nil {
		return nil, nil, libUseTime, err
	}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	//part1:prompt
	roleType := define.PromptRoleTypeMap[cast.ToInt(params.Robot[`prompt_role_type`])]
	prompt, libraryContent := common.FormatSystemPrompt(params.Prompt, list)
	if roleType == define.PromptRoleUser {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: libraryContent})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: libraryContent})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: roleType, Content: prompt})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: prompt})
	}
	//part2:context_qa
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	//part3:question,prompt+question
	if roleType == define.PromptRoleUser {
		content := strings.Join([]string{params.Prompt, params.Question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: content})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: params.Question})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
	}
	return messages, list, libUseTime, nil
}

func buildDirectChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, []msql.Params, error) {
	if len(params.Prompt) == 0 { //no custom is used
		params.Prompt = common.BuildPromptStruct(cast.ToInt(params.Robot[`prompt_type`]), params.Robot[`prompt`], params.Robot[`prompt_struct`])
	}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	//part1:prompt
	prompt, _ := common.FormatSystemPrompt(params.Prompt, nil)
	roleType := define.PromptRoleTypeMap[cast.ToInt(params.Robot[`prompt_role_type`])]
	if roleType != define.PromptRoleUser {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: roleType, Content: prompt})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: prompt})
	}
	//part2:context_qa
	contextList := common.BuildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	//part3:cur_question
	if roleType == define.PromptRoleUser {
		content := strings.Join([]string{prompt, params.Question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: content})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: params.Question})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
	}
	return messages, []msql.Params{}, nil
}
