// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/llm/adaptor"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
		message[`content`] = define.LocalUploadPrefix + `default/robot_avatar.png`
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
	for _, modelInfo := range define.ModelList {
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
	nickname := strings.TrimSpace(c.PostForm(`nickname`))
	name := strings.TrimSpace(c.PostForm(`name`))
	avatar := strings.TrimSpace(c.PostForm(`avatar`))
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
	prompt := strings.ReplaceAll(define.PromptDefaultQuestionGuide, `{{histories}}`, histories)
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `user`, Content: prompt}}

	//var messages []adaptor.ZhimaChatCompletionMessage
	//messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: define.PromptDefaultQuestionGuide})
	//for _, msg := range questionGuideMessages {
	//	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: msg.Role, Content: msg.Content})
	//}
	//messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: `请按要求回答`})

	chatResp, _, err := common.RequestChat(
		chatBaseParam.AdminUserId,
		chatBaseParam.Openid,
		chatBaseParam.Robot,
		chatBaseParam.AppType,
		cast.ToInt(chatBaseParam.Robot[`model_config_id`]),
		chatBaseParam.Robot[`use_model`],
		messages,
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

	c.String(http.StatusOK, lib_web.FmtJson(guides, nil))
}

func getChatRequestParam(c *gin.Context) *define.ChatRequestParam {
	chatBaseParam, err := common.CheckChatRequest(c)
	isClose := false
	return &define.ChatRequestParam{
		ChatBaseParam: chatBaseParam,
		Error:         err,
		Lang:          common.GetLang(c),
		Question:      strings.TrimSpace(c.PostForm(`question`)),
		DialogueId:    cast.ToInt(c.PostForm(`dialogue_id`)),
		Prompt:        strings.TrimSpace(c.PostForm(`prompt`)),
		LibraryIds:    strings.TrimSpace(c.PostForm(`library_ids`)),
		IsClose:       &isClose,
	}
}

func DoChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event) (msql.Params, error) {
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
	sessionId, err := common.GetSessionId(params.ChatBaseParam, dialogueId)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	chanStream <- sse.Event{Event: `dialogue_id`, Data: dialogueId}
	chanStream <- sse.Event{Event: `session_id`, Data: sessionId}
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
	lastChat := msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: message[`content`],
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat)
	//message push
	customer, err := common.GetCustomerInfo(params.Openid, params.AdminUserId)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	chanStream <- sse.Event{Event: `customer`, Data: customer}
	chanStream <- sse.Event{Event: `c_message`, Data: common.ToStringMap(message, `id`, id)}
	//obtain the data required for gpt
	chanStream <- sse.Event{Event: `robot`, Data: params.Robot}
	debugLog := make([]any, 0) //debug log
	var messages []adaptor.ZhimaChatCompletionMessage
	var list []msql.Params
	var recallTime int64
	if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeDirect {
		messages, list, err = buildDirectChatRequestMessage(params, id, dialogueId, &debugLog)
	} else {
		recallStart := time.Now()
		messages, list, err = buildLibraryChatRequestMessage(params, id, dialogueId, &debugLog)
		recallTime = time.Now().Sub(recallStart).Milliseconds()
		chanStream <- sse.Event{Event: `recall_time`, Data: recallTime}
	}
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}

	var (
		content, menuJson string
		requestTime       int64
		chatResp          = adaptor.ZhimaChatCompletionResponse{}
	)
	msgType := define.MsgTypeText
	if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeDirect {
		if useStream {
			chatResp, requestTime, err = common.RequestChatStream(
				params.AdminUserId,
				params.Openid,
				params.Robot,
				params.AppType,
				cast.ToInt(params.Robot[`model_config_id`]),
				params.Robot[`use_model`],
				messages,
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
				cast.ToFloat32(params.Robot[`temperature`]),
				cast.ToInt(params.Robot[`max_token`]),
			)
		}
		content = chatResp.Result
		if err != nil {
			logs.Error(err.Error())
			sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
		}
	} else if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeMixture {
		if len(list) == 0 {
			if useStream {
				chatResp, requestTime, err = common.RequestChatStream(
					params.AdminUserId,
					params.Openid,
					params.Robot,
					params.AppType,
					cast.ToInt(params.Robot[`model_config_id`]),
					params.Robot[`use_model`],
					messages,
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
					cast.ToFloat32(params.Robot[`temperature`]),
					cast.ToInt(params.Robot[`max_token`]),
				)
			}
			content = chatResp.Result
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
				chanStream <- sse.Event{Event: `sending`, Data: content}
			} else {
				if useStream {
					chatResp, requestTime, err = common.RequestChatStream(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
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
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				}
				content = chatResp.Result
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
				chanStream <- sse.Event{Event: `sending`, Data: content}
			} else { // ask gpt
				if useStream {
					chatResp, requestTime, err = common.RequestChatStream(
						params.AdminUserId,
						params.Openid,
						params.Robot,
						params.AppType,
						cast.ToInt(params.Robot[`model_config_id`]),
						params.Robot[`use_model`],
						messages,
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
						cast.ToFloat32(params.Robot[`temperature`]),
						cast.ToInt(params.Robot[`max_token`]),
					)
				}
				content = chatResp.Result
				if err != nil {
					logs.Error(err.Error())
					sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
				}
			}
		}
	}

	if *params.IsClose { //client break
		return nil, errors.New(`client break`)
	}
	//push prompt log
	debugLog = append(debugLog, map[string]string{`type`: `cur_answer`, `content`: content})
	chanStream <- sse.Event{Event: `debug`, Data: debugLog}
	//dispose answer source
	quoteFile, ms := make([]msql.Params, 0), map[string]struct{}{}
	for _, one := range list {
		if _, ok := ms[one[`file_id`]]; ok {
			continue //remove duplication
		}
		ms[one[`file_id`]] = struct{}{}
		quoteFile = append(quoteFile, msql.Params{
			`id`:        one[`file_id`],
			`file_name`: one[`file_name`],
		})
	}
	quoteFileJson, _ := tool.JsonEncode(quoteFile)
	//database dispose
	message = msql.Datas{
		`admin_user_id`: params.AdminUserId,
		`robot_id`:      params.Robot[`id`],
		`openid`:        params.Openid,
		`dialogue_id`:   dialogueId,
		`session_id`:    sessionId,
		`is_customer`:   define.MsgFromRobot,
		`request_time`:  requestTime,
		`recall_time`:   recallTime,
		`msg_type`:      msgType,
		`content`:       content,
		`menu_json`:     menuJson,
		`quote_file`:    quoteFileJson,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	lastChat = msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: message[`content`],
	}
	id, err = msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(params.Lang, `sys_err`)}
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat)
	//message push
	chanStream <- sse.Event{Event: `ai_message`, Data: common.ToStringMap(message, `id`, id)}
	if len(quoteFile) > 0 && cast.ToBool(params.Robot[`answer_source_switch`]) {
		chanStream <- sse.Event{Event: `quote_file`, Data: quoteFile}
	}
	//save answer source
	if len(list) > 0 {
		asm := msql.Model(`chat_ai_answer_source`, define.Postgres)
		for _, one := range list {
			_, err := asm.Insert(msql.Datas{
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
	chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
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

func buildLibraryChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, []msql.Params, error) {
	if len(params.Prompt) == 0 { //no custom is used
		params.Prompt = params.Robot[`prompt`]
	}
	if len(params.LibraryIds) == 0 || !common.CheckIds(params.LibraryIds) { //no custom is used
		params.LibraryIds = params.Robot[`library_ids`]
	}

	contextList := buildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))

	//question optimize
	var optimizedQuestions []string
	if cast.ToBool(params.Robot[`enable_question_optimize`]) && len(params.LibraryIds) > 0 {
		var err error
		optimizedQuestions, err = common.GetOptimizedQuestions(params, contextList)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	//convert match
	list, err := common.GetMatchLibraryParagraphList(
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
	if err != nil {
		return nil, nil, err
	}

	//part1:prompt
	responseTypeMsg := buildChatResponseType(cast.ToInt(params.Robot["show_type"]), params.Lang)
	prompt := params.Prompt
	prompt = prompt + "\n\n" + define.PromptDefaultAnswerImage
	prompt = prompt + "\n\n" + responseTypeMsg
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: prompt})

	//part2:library
	for _, one := range list {
		var images []string
		err = tool.JsonDecode(one[`images`], &images)
		if err != nil {
			logs.Error(err.Error())
		}
		if cast.ToInt(one[`type`]) == define.ParagraphTypeNormal {
			messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: common.EmbTextImages(one[`content`], images)})
			*debugLog = append(*debugLog, map[string]string{`type`: `library`, `content`: common.EmbTextImages(one[`content`], images)})
		} else {
			messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: "question: " + one[`question`] + "\nanswer: " + common.EmbTextImages(one[`answer`], images)})
			*debugLog = append(*debugLog, map[string]string{`type`: `library`, `content`: "question: " + one[`question`] + "\nanswer: " + common.EmbTextImages(one[`answer`], images)})
		}
	}

	//part3:context_qa
	// Add a parameter if you need to clarify the distinction
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}

	//part4:cur_question
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: responseTypeMsg + params.Question})
	*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: responseTypeMsg + params.Question})

	return messages, list, nil
}

func buildDirectChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, []msql.Params, error) {
	var messages []adaptor.ZhimaChatCompletionMessage
	// Add a parameter if you need to clarify the distinction
	responseTypeMsg := buildChatResponseType(cast.ToInt(params.Robot["show_type"]), params.Lang)
	//messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: responseTypeMsg})
	//*debugLog = append(*debugLog, map[string]string{`type`: `system`, `content`: responseTypeMsg})
	contextList := buildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: responseTypeMsg + params.Question})
	*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: responseTypeMsg + params.Question})
	return messages, []msql.Params{}, nil
}

func buildChatContextPair(openid string, robotId, dialogueId, curMsgId, contextPair int) []map[string]string {
	contextList := make([]map[string]string, 0)
	if contextPair <= 0 {
		return contextList //no context required
	}
	list, err := msql.Model(`chat_ai_message`, define.Postgres).Where(`openid`, openid).
		Where(`robot_id`, cast.ToString(robotId)).Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`msg_type`, cast.ToString(define.MsgTypeText)).Where(`id`, `<`, cast.ToString(curMsgId)).
		Order(`id desc`).Field(`id,content,is_customer`).Limit(contextPair * 4).Select()
	if err != nil {
		logs.Error(err.Error())
	}
	if len(list) == 0 {
		return contextList
	}
	//reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	//foreach
	for i := 0; i < len(list)-1; i++ {
		if cast.ToInt(list[i][`is_customer`]) == define.MsgFromCustomer && cast.ToInt(list[i+1][`is_customer`]) == define.MsgFromRobot {
			contextList = append(contextList, map[string]string{`question`: list[i][`content`], `answer`: list[i+1][`content`]})
			i++ //skip answer
		}
	}
	//cut out
	if len(contextList) > contextPair {
		contextList = contextList[len(contextList)-contextPair:]
	}
	return contextList
}

func buildChatResponseType(showType int, lang string) string {
	result := ""
	if showType == define.RobotMarkdownResponse {
		result = fmt.Sprintf("(%s)", i18n.Show(lang, `chat_show_type`, `markdown格式`))
	}
	return result
}
