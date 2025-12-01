// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
	loopTestParams := make([]any, 0)
	_ = tool.JsonDecodeUseNumber(c.DefaultPostForm(`loop_test_params`, `[]`), &loopTestParams)
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
		LoopTestParams: loopTestParams,
	}
}

func OnlyReceivedMessageReply(params *define.ChatRequestParam) (msql.Params, error) {
	monitor := common.NewMonitor(params)
	message, err := OnlyReceivedMessageReplyHandle(params, monitor)
	if len(message) > 0 {
		monitor.Save(err)
	}
	return message, err
}

// getMsgTypeByReceivedMessageType 不清楚为啥数据库中存的是int类型，这里做兼容
func getMsgTypeByReceivedMessageType(ReceivedMessageType string) int {
	switch ReceivedMessageType {
	case lib_define.MsgTypeText:
		return define.MsgTypeText
	case lib_define.MsgTypeImage:
		return define.MsgTypeImage
	case lib_define.MsgTypeVoice:
		return define.MsgTypeVoice
	case lib_define.MsgTypeVideo:
		return define.MsgTypeVideo
	case lib_define.MsgTypeShortVideo:
		return define.MsgTypeShortVideo
	case lib_define.MsgTypeMinirogrampage:
		return define.MsgTypeMinirogrampage
	case lib_define.MsgTypeLocation:
		return define.MsgTypeLocation
	case lib_define.MsgTypeLink:
		return define.MsgTypeLink
	case lib_define.MsgTypeEvent:
		return define.MsgTypeEvent
	}
	return define.MsgTypeOther
}

func OnlyReceivedMessageReplyHandle(params *define.ChatRequestParam, monitor *common.Monitor) (msql.Params, error) {
	var err error
	dialogueId := params.DialogueId
	sessionId, err := common.GetSessionId(params, dialogueId)
	customer, err := common.GetCustomerInfo(params.Openid, params.AdminUserId)
	//msgType := getMsgTypeByReceivedMessageType(params.ReceivedMessageType)
	msgType := define.MsgTypeText
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	var receivedMessageJson string
	if len(params.ReceivedMessage) > 0 {
		//回复内容
		receivedMessageJson = tool.JsonEncodeNoError(params.ReceivedMessage)
	}

	//显示内容
	showContent, isName := lib_define.MsgTypeNameMap[params.ReceivedMessageType]
	if !isName {
		showContent = `未知`
	}
	showContent = `收到【` + showContent + `】类型的消息`
	//展示图片消息
	if params.ReceivedMessageType == lib_define.MsgTypeImage && params.MediaIdToOssUrl != `` {
		msgType = define.MsgTypeImage
		showContent = params.MediaIdToOssUrl
	}

	message := msql.Datas{
		`admin_user_id`:             params.AdminUserId,
		`robot_id`:                  params.Robot[`id`],
		`openid`:                    params.Openid,
		`dialogue_id`:               dialogueId,
		`session_id`:                sessionId,
		`is_customer`:               define.MsgFromCustomer,
		`msg_type`:                  msgType,
		`content`:                   showContent,
		`received_message_type`:     params.ReceivedMessageType,
		`received_message`:          receivedMessageJson,
		`media_id_to_oss_url`:       params.MediaIdToOssUrl,
		`thumb_media_id_to_oss_url`: params.ThumbMediaIdToOssUrl,
		`menu_json`:                 ``,
		`quote_file`:                `[]`,
		`create_time`:               tool.Time2Int(),
		`update_time`:               tool.Time2Int(),
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
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromCustomer)
	//websocket notify
	common.ReceiverChangeNotify(params.AdminUserId, `c_message`, common.ToStringMap(message, `id`, id))

	debugLog := make([]any, 0) //debug log
	defer func() {
		monitor.DebugLog = debugLog //记录监控数据
	}()

	var receivedMessageReplyList []common.ReplyContent
	//收到消息回复处理
	receivedMessageReplyList, _ = buildReceivedMessageReply(params, params.ReceivedMessageType, &debugLog)

	if len(receivedMessageReplyList) == 0 {
		logs.Error(`received message reply list is empty`)
		return msql.Params{}, nil
	}
	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
		llmStartTime                        = time.Now()
	)

	//记录监控数据
	monitor.LlmCallTime = time.Now().Sub(llmStartTime).Milliseconds()
	monitor.RequestTime, monitor.Error = requestTime, err

	if *params.IsClose { //client break
		return nil, errors.New(`client break`)
	}

	quoteFile, _ := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	quoteFileJson, _ := tool.JsonEncode(quoteFileForSave)

	message = msql.Datas{
		`admin_user_id`:          params.AdminUserId,
		`robot_id`:               params.Robot[`id`],
		`openid`:                 params.Openid,
		`dialogue_id`:            dialogueId,
		`session_id`:             sessionId,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           requestTime,
		`recall_time`:            monitor.LibUseTime.RecallTime,
		`msg_type`:               define.MsgTypeText,
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
	if len(receivedMessageReplyList) > 0 {
		//回复内容
		receivedMessageReplyListJson := tool.JsonEncodeNoError(receivedMessageReplyList)
		message[`reply_content_list`] = receivedMessageReplyListJson
	}

	lastChat = msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err = msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromRobot)
	//websocket notify
	common.ReceiverChangeNotify(params.AdminUserId, `ai_message`, common.ToStringMap(message, `id`, id))

	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
	message["use_model"] = params.Robot["use_model"]
	return common.ToStringMap(message, `id`, id), nil
}

// SubscribeReplyHandle 关注后回复处理
func SubscribeReplyHandle(params *define.ChatRequestParam, subscribeScene string) (msql.Params, error) {
	var err error
	dialogueId := params.DialogueId
	sessionId, err := common.GetSessionId(params, dialogueId)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	//显示内容
	var subscribeReplyList []common.ReplyContent
	//收到消息回复处理
	subscribeReplyList, _ = buildSubscribeReply(params, subscribeScene)

	if len(subscribeReplyList) == 0 {
		logs.Error(`subscribe reply list is empty`)
		return msql.Params{}, nil
	}
	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
	)

	if *params.IsClose { //client break
		return nil, errors.New(`client break`)
	}

	quoteFile, _ := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	quoteFileJson, _ := tool.JsonEncode(quoteFileForSave)

	message := msql.Datas{
		`admin_user_id`:          params.AdminUserId,
		`robot_id`:               params.Robot[`id`],
		`openid`:                 params.Openid,
		`dialogue_id`:            dialogueId,
		`session_id`:             sessionId,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           requestTime,
		`recall_time`:            tool.Time2Int(),
		`msg_type`:               define.MsgTypeText,
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
	if len(subscribeReplyList) > 0 {
		//回复内容
		subscribeReplyListJson := tool.JsonEncodeNoError(subscribeReplyList)
		message[`reply_content_list`] = subscribeReplyListJson
	}

	lastChat := msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	common.UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromRobot)
	//websocket notify
	common.ReceiverChangeNotify(params.AdminUserId, `ai_message`, common.ToStringMap(message, `id`, id))

	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
	message["use_model"] = params.Robot["use_model"]
	return common.ToStringMap(message, `id`, id), nil
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
	//close open_api receiver
	CloseReceiverFromAppOpenApi(params)
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
	//manual_reply_pause_robot_reply
	if IsManualReplyPauseRobotReply(params, sessionId, chanStream) {
		return nil, nil
	}
	//keyword_switch_manual
	if msg, ok := IsKeywordSwitchManual(params, sessionId, dialogueId, chanStream); ok {
		return msg, nil
	}
	//intention_switch_manual
	if msg, ok := IsIntentionSwitchManual(params, sessionId, dialogueId, monitor, chanStream); ok {
		return msg, nil
	}
	isBackground := len(params.Customer) > 0 && cast.ToInt(params.Customer[`is_background`]) > 0
	//obtain the data required for gpt
	chanStream <- sse.Event{Event: `robot`, Data: params.Robot}
	debugLog := make([]any, 0) //debug log
	defer func() {
		monitor.DebugLog = debugLog //记录监控数据
	}()
	var messages []adaptor.ZhimaChatCompletionMessage
	var list []msql.Params
	var (
		showQuoteFile   = cast.ToBool(params.Robot[`answer_source_switch`])
		startQuoteFile  bool
		hitCache        bool
		answerMessageId = ""
	)

	var keywordSkipAI bool
	var replyContentList []common.ReplyContent
	var keywordReplyList []common.ReplyContent
	//默认不跳过
	keywordSkipAI = false
	//关键词检测处理
	keywordReplyList, keywordSkipAI, err = buildKeywordReplyMessage(params, &debugLog)
	if len(keywordReplyList) > 0 {
		replyContentList = append(replyContentList, keywordReplyList...)
	}

	if !keywordSkipAI {
		//收到消息回复处理
		receivedMessageReplyList, _ := buildReceivedMessageReply(params, lib_define.MsgTypeText, &debugLog)
		if len(receivedMessageReplyList) > 0 {
			replyContentList = append(replyContentList, receivedMessageReplyList...)
		}
	}

	//构造发送给ai的请求消息参数 messages
	if keywordSkipAI {
		//跳过ai回复 不需要构造发送给ai的请求消息参数

	} else if cast.ToInt(params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		//nothing to do
	} else if cast.ToInt(params.Robot[`chat_type`]) == define.ChatTypeDirect {
		//直连模式 不考虑知识库
		messages, list, err = buildDirectChatRequestMessage(params, id, dialogueId, &debugLog)
	} else {
		//混合模式 和 仅知识库模式
		//repetition_switch_manual
		if msg, ok := IsRepetitionSwitchManual(params, sessionId, dialogueId, id, chanStream); ok {
			return msg, nil
		}
		if showQuoteFile {
			chanStream <- sse.Event{Event: `start_quote_file`, Data: tool.Time2Int()}
			startQuoteFile = true
		}
		hitCache, answerMessageId = common.HitRobotMessageCache(params.Robot[`robot_key`], params.Question, params.Robot[`cache_config`])
		if hitCache {
			list, hitCache, err = common.BuildLibraryMessagesFromCache(params.Robot[`robot_key`], answerMessageId)
		} else {
			messages, list, monitor.LibUseTime, err = buildLibraryChatRequestMessage(params, id, dialogueId, &debugLog)

			chanStream <- sse.Event{Event: `recall_time`, Data: monitor.LibUseTime.RecallTime}
			if !isBackground && len(list) == 0 { //未知问题统计
				common.SaveUnknownIssueRecord(params.AdminUserId, params.Robot, params.Question)
			}
			//unknown_switch_manual
			if msg, ok := IsUnknownSwitchManual(params, sessionId, dialogueId, list, chanStream); ok {
				return msg, nil
			}
		}
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
	var workFlowFuncCall []adaptor.FunctionTool
	var needRunWorkFlow = false
	if !keywordSkipAI { //非关键词跳过ai
		//构建 工作流
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
		workFlowFuncCall, needRunWorkFlow = work_flow.BuildFunctionTools(params.Robot)
		if needRunWorkFlow {
			functionTools = append(functionTools, workFlowFuncCall...)
		}
	}

	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
		llmStartTime                        = time.Now()
		saveRobotChatCache                  bool
		isSwitchManual                      bool
	)

	//关键词处理

	msgType := define.MsgTypeText
	if len(replyContentList) > 0 {
		chanStream <- sse.Event{Event: `reply_content_list`, Data: replyContentList}
	}

	if len(keywordReplyList) > 0 && keywordSkipAI {
		//关键词直接回复 跳过ai处理
	} else if cast.ToInt(params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		workFlowParams := &work_flow.WorkFlowParams{ChatRequestParam: params, CurMsgId: int(id), DialogueId: dialogueId, SessionId: sessionId}
		content, requestTime, monitor.LibUseTime, list, err = work_flow.CallWorkFlow(workFlowParams, &debugLog, monitor, &isSwitchManual)
		if err != nil {
			sendDefaultUnknownQuestionPrompt(params, err.Error(), chanStream, &content)
			debugLog = append(debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
			chanStream <- sse.Event{Event: `debug`, Data: debugLog} //渲染Prompt日志
			return nil, err
		}
		chanStream <- sse.Event{Event: `recall_time`, Data: monitor.LibUseTime.RecallTime}
		chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
		chanStream <- sse.Event{Event: `sending`, Data: content}
	} else if hitCache {
		// response from cache
		chatResp, requestTime, err = common.ResponseMessagesFromCache(params.Robot[`robot_key`], answerMessageId, useStream, chanStream)
		content = chatResp.Result
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
			} else {
				saveRobotChatCache = true
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
							content += fmt.Sprintf("\n![img](%s)", imageUrl)
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
				} else {
					saveRobotChatCache = true
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
							content += fmt.Sprintf("\n![img](%s)", imageUrl)
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
				} else {
					saveRobotChatCache = true
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
			content, requestTime, _, _, err = work_flow.CallWorkFlow(workFlowParams, &debugLog, monitor, &isSwitchManual)
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

	//未认证公众号的消息特殊处理
	if params.AppType == lib_define.AppOfficeAccount && params.PassiveId > 0 {
		PassiveReplyLogNotify(params.PassiveId, params.Question, content)
	}

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
	if len(replyContentList) > 0 {
		//关键词回复 触发内容
		replyContentListJson := tool.JsonEncodeNoError(replyContentList)
		message[`reply_content_list`] = replyContentListJson
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
	if saveRobotChatCache {
		go common.SetRobotMessageCache(params.Robot[`robot_key`], params.Question, cast.ToString(id), params.Robot[`cache_config`])
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
	AdditionQuoteLib(params, list, &message) //quote_lib
	chanStream <- sse.Event{Event: `data`, Data: message}
	chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
	return common.ToStringMap(message, `id`, id, `is_switch_manual`, isSwitchManual), nil
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

// GetRandomSliceReply 从回复内容列表中随机选择指定数量的条目
// 如果请求数量大于列表总数，则返回全部内容
// 如果请求数量小于等于0，则返回空列表
func GetRandomSliceReply(replyList []common.ReplyContent, num int) []common.ReplyContent {
	// 边界条件检查
	if len(replyList) == 0 || num <= 0 {
		return []common.ReplyContent{}
	}

	// 如果请求数量大于总数，则返回全部
	if num >= len(replyList) {
		return replyList
	}

	// 创建结果切片
	result := make([]common.ReplyContent, 0, num)

	// 创建索引切片并随机打乱
	indexes := make([]int, len(replyList))
	for i := range indexes {
		indexes[i] = i
	}

	// Fisher-Yates shuffle 算法随机打乱索引
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(indexes) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		indexes[i], indexes[j] = indexes[j], indexes[i]
	}

	// 取前num个元素
	for i := 0; i < num; i++ {
		result = append(result, replyList[indexes[i]])
	}

	return result
}

// buildKeywordReplyMessage 构建关键词回复消息
func buildKeywordReplyMessage(params *define.ChatRequestParam, debugLog *[]any) ([]common.ReplyContent, bool, error) {
	//part0:init messages
	var replyList []common.ReplyContent
	//判断是否关键词跳过ai回复
	var keywordSkipAI = false

	//判断开关
	robotId := cast.ToInt(params.Robot[`id`])
	adminUserId := cast.ToInt(params.Robot[`admin_user_id`])
	//关键词回复
	robotAbilityConfig := common.GetRobotAbilityConfigByAbilityType(adminUserId, robotId, common.RobotAbilityAutoReply)
	if len(robotAbilityConfig) == 0 {
		//关键词回复没开启
		return replyList, false, nil
	}

	//获取所有关键词缓存
	robotKeywordReplyList, err := common.GetRobotKeywordReplyListByRobotId(robotId)
	if err != nil {
		return replyList, false, err
	}

	//关键词回复是否跳过ai
	keywordSkipAI = cast.ToInt(robotAbilityConfig[`ai_reply_status`]) != define.SwitchOn

	//问题判断
	question := strings.TrimSpace(params.Question)

	//循环判断  构造消息
	for _, robotKeywordReply := range robotKeywordReplyList {
		//判断关键词
		if robotKeywordReply.SwitchStatus != define.SwitchOn {
			continue
		}
		keywordFlag := false
		// 精确匹配 FullKeyword
		for _, keyword := range robotKeywordReply.FullKeyword {
			if question == keyword {
				//匹配成功 构造消息
				keywordFlag = true
				break
			}
		}

		// 包含匹配 HalfKeyword
		for _, keyword := range robotKeywordReply.HalfKeyword {
			if strings.Contains(question, keyword) {
				//匹配成功 构造消息
				keywordFlag = true
				break
			}
		}

		if keywordFlag {
			//匹配成功 判断回复类型
			if robotKeywordReply.ReplyNum == 0 {
				replyList = append(replyList, robotKeywordReply.ReplyContent...)
			} else {
				//随机选择 ReplyNum 条
				//数组中随机切取 ReplyNum 条
				selectReplyList := GetRandomSliceReply(robotKeywordReply.ReplyContent, robotKeywordReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
		}
	}
	//判断是否继续ai
	if len(replyList) == 0 {
		//没有匹配到关键词 继续ai
		return replyList, false, nil
	}
	//循环replyList标记来源
	for i := range replyList {
		//标记来源
		replyList[i].SendSource = common.RobotAbilityKeywordReply
	}
	//返回消息
	return replyList, keywordSkipAI, nil
}

// buildSubscribeReply 构建关注后回复消息
func buildSubscribeReply(params *define.ChatRequestParam, subscribeScene string) ([]common.ReplyContent, error) {
	//part0:init messages
	var replyList []common.ReplyContent
	//判断开关
	robotId := cast.ToInt(params.Robot[`id`])
	adminUserId := cast.ToInt(params.Robot[`admin_user_id`])
	appid := cast.ToString(params.AppInfo[`appid`]) // 从 params.AppInfo 获取 appid
	//关键词回复
	robotAbilityConfig := common.GetRobotAbilityConfigByAbilityType(adminUserId, robotId, common.RobotAbilitySubscribeReply)
	if len(robotAbilityConfig) == 0 {
		//关键词回复没开启
		return replyList, nil
	}

	// 1. 获取今天是星期几的int值
	// time.Weekday：Sunday=0, Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6
	weekday := cast.ToInt(time.Now().Weekday())
	//获取所有关键词缓存
	subscribeReplyList, err := common.GetRobotSubscribeReplyListByAppid(robotId, appid, define.RuleTypeSubscribeSource)
	if err != nil {

		return replyList, err
	}

	needCheck := true

	//来源检测
	for _, subscribeReply := range subscribeReplyList {
		//判断关键词
		if subscribeReply.SwitchStatus != define.SwitchOn {
			continue
		}
		if subscribeReply.RuleType != define.RuleTypeSubscribeSource {
			//不是来源的
			continue
		}
		//检测来源
		if !tool.InArrayString(subscribeScene, subscribeReply.SubscribeSource) {
			//不在指定来源中
			continue
		}
		//匹配到消息 则跳过 消息检测
		needCheck = false
		//匹配成功 构造消息
		if subscribeReply.ReplyNum == 0 {
			replyList = append(replyList, subscribeReply.ReplyContent...)
		} else {
			//随机选择 ReplyNum 条
			//数组中随机切取 ReplyNum 条
			selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
			if len(selectReplyList) > 0 {
				replyList = append(replyList, selectReplyList...)
			}
		}
		break
	}

	if needCheck && len(replyList) == 0 {
		//时间检测
		subscribeReplyList, err = common.GetRobotSubscribeReplyListByAppid(robotId, appid, define.RuleTypeSubscribeDuration)
		if err != nil {

			return replyList, err
		}

		//循环判断  构造消息
		for _, subscribeReply := range subscribeReplyList {
			//判断关键词
			if subscribeReply.SwitchStatus != define.SwitchOn {
				continue
			}
			if subscribeReply.RuleType != define.RuleTypeSubscribeDuration {
				//不是时间规则
				continue
			}
			//时间规则 且 开启
			checkFlag := false
			switch subscribeReply.DurationType {
			case common.DurationTypeWeek:
				// 周
				for _, day := range subscribeReply.WeekDuration {
					if day == weekday {
						checkFlag = true
						break
					}
				}
				break
			case common.DurationTypeDay:
				// 每天
				checkFlag = true
				break
			case common.DurationTypeTimeRange:
				// 时间范围
				checkFlag = isTodayInDateRange(subscribeReply.StartDay, subscribeReply.EndDay)
				break
			default:
				// 默认
				break
			}
			//判断是否继续
			if !checkFlag {
				continue
			}
			//对比时间 不在范围的跳过
			if !nowInHHmmRangeSimple(subscribeReply.StartDuration, subscribeReply.EndDuration) {
				continue
			}
			//匹配到消息 则跳过 消息检测
			needCheck = false
			//判断时间间隔
			if subscribeReply.ReplyInterval > 0 {
				//判断时间间隔
				var lastTime int
				lastTime, err = common.GetReceivedMessageReplyLastTime(robotId, subscribeReply.ID, params.Openid)
				if err != nil {
					continue
				}
				nextTime := lastTime + subscribeReply.ReplyInterval
				if nextTime > tool.Time2Int() {
					//时间间隔内中 不满足
					break
				}
				// 设置时间间隔
				err = common.SetReceivedMessageReplyLastTime(robotId, subscribeReply.ID, tool.Time2Int(), params.Openid)
				if err != nil {
					return nil, err
				}
			}

			//匹配成功 构造消息
			if subscribeReply.ReplyNum == 0 {
				replyList = append(replyList, subscribeReply.ReplyContent...)
			} else {
				//随机选择 ReplyNum 条
				//数组中随机切取 ReplyNum 条
				selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
			break
		}
	}

	if needCheck && len(replyList) == 0 {
		//按默认关注开启判断
		subscribeReplyList, err = common.GetRobotSubscribeReplyListByAppid(robotId, appid, define.RuleTypeSubscribeDefault)
		if err != nil {

			return replyList, err
		}
		for _, subscribeReply := range subscribeReplyList {
			if subscribeReply.SwitchStatus != define.SwitchOn {
				continue
			}
			if subscribeReply.RuleType != define.RuleTypeSubscribeDefault {
				//不是默认类型
				continue
			}

			//匹配成功 构造消息
			if subscribeReply.ReplyNum == 0 {
				replyList = append(replyList, subscribeReply.ReplyContent...)
			} else {
				//随机选择 ReplyNum 条
				//数组中随机切取 ReplyNum 条
				selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
		}
	}

	//判断是否继续ai
	if len(replyList) == 0 {
		//没有匹配到关键词 继续ai
		return replyList, nil
	}
	//循环replyList标记来源
	for i := range replyList {
		//标记来源
		replyList[i].SendSource = common.RobotAbilitySubscribeReply
	}
	//返回消息
	return replyList, nil
}

func buildReceivedMessageReply(params *define.ChatRequestParam, messageType string, debugLog *[]any) ([]common.ReplyContent, error) {
	//part0:init messages
	var replyList []common.ReplyContent
	//判断开关
	robotId := cast.ToInt(params.Robot[`id`])
	adminUserId := cast.ToInt(params.Robot[`admin_user_id`])
	//关键词回复
	robotAbilityConfig := common.GetRobotAbilityConfigByAbilityType(adminUserId, robotId, common.RobotAbilityAutoReply)
	if len(robotAbilityConfig) == 0 {
		//关键词回复没开启
		return replyList, nil
	}

	// 1. 获取今天是星期几的int值
	// time.Weekday：Sunday=0, Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6
	weekday := cast.ToInt(time.Now().Weekday())
	//获取所有关键词缓存
	receivedMessageRuleList, err := common.GetRobotReceivedMessageReplyListByRobotId(robotId, common.RuleTypeDuration)
	if err != nil {

		return replyList, err
	}

	messageTypeCheck := true
	//循环判断  构造消息
	for _, receivedMessageRule := range receivedMessageRuleList {
		//判断关键词
		if receivedMessageRule.SwitchStatus != define.SwitchOn {
			continue
		}
		if receivedMessageRule.RuleType != common.RuleTypeDuration {
			//不是时间规则
			continue
		}
		//时间规则 且 开启
		checkFlag := false
		switch receivedMessageRule.DurationType {
		case common.DurationTypeWeek:
			// 周
			for _, day := range receivedMessageRule.WeekDuration {
				if day == weekday {
					checkFlag = true
					break
				}
			}
			break
		case common.DurationTypeDay:
			// 每天
			checkFlag = true
			break
		case common.DurationTypeTimeRange:
			// 时间范围
			checkFlag = isTodayInDateRange(receivedMessageRule.StartDay, receivedMessageRule.EndDay)
			break
		default:
			// 默认
			break
		}
		//判断是否继续
		if !checkFlag {
			continue
		}
		//对比时间 不在范围的跳过
		if !nowInHHmmRangeSimple(receivedMessageRule.StartDuration, receivedMessageRule.EndDuration) {
			continue
		}
		//匹配到消息 则跳过 消息检测
		messageTypeCheck = false
		//判断时间间隔
		if receivedMessageRule.ReplyInterval > 0 {
			//判断时间间隔
			var lastTime int
			lastTime, err = common.GetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, params.Openid)
			if err != nil {
				continue
			}
			nextTime := lastTime + receivedMessageRule.ReplyInterval
			if nextTime > tool.Time2Int() {
				//时间间隔内中 不满足
				break
			}
			// 设置时间间隔
			err = common.SetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, tool.Time2Int(), params.Openid)
			if err != nil {
				return nil, err
			}
		}

		//匹配成功 构造消息
		if receivedMessageRule.ReplyNum == 0 {
			replyList = append(replyList, receivedMessageRule.ReplyContent...)
		} else {
			//随机选择 ReplyNum 条
			//数组中随机切取 ReplyNum 条
			selectReplyList := GetRandomSliceReply(receivedMessageRule.ReplyContent, receivedMessageRule.ReplyNum)
			if len(selectReplyList) > 0 {
				replyList = append(replyList, selectReplyList...)
			}
		}
		break
	}

	if messageTypeCheck && len(replyList) == 0 {
		//按消息类型判断
		receivedMessageRuleList, err = common.GetRobotReceivedMessageReplyListByRobotId(robotId, common.RuleTypeMessageType)
		if err != nil {

			return replyList, err
		}
		for _, receivedMessageRule := range receivedMessageRuleList {
			if receivedMessageRule.SwitchStatus != define.SwitchOn {
				continue
			}
			if receivedMessageRule.RuleType != common.RuleTypeMessageType {
				//不是指定消息类型
				continue
			}
			checkFlag := false
			switch receivedMessageRule.MessageType {
			case common.MessageTypeAll:
				checkFlag = true
				break
			case common.MessageTypeSpecify:
				//指定消息类型
				for _, msgType := range receivedMessageRule.SpecifyMessageType {
					if messageType == msgType {
						checkFlag = true
						break
					}
				}
				break
			default:
				// 默认
				break
			}
			//判断是否继续
			if !checkFlag {
				continue
			}
			//判断时间间隔
			if receivedMessageRule.ReplyInterval > 0 {
				//判断时间间隔
				var lastTime int
				lastTime, err = common.GetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, params.Openid)
				if err != nil {
					continue
				}
				nextTime := lastTime + receivedMessageRule.ReplyInterval
				if nextTime > tool.Time2Int() {
					//时间间隔内中 不满足
					break
				}
				// 设置时间间隔
				err = common.SetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, tool.Time2Int(), params.Openid)
				if err != nil {
					return nil, err
				}
			}
			//匹配成功 构造消息
			if receivedMessageRule.ReplyNum == 0 {
				replyList = append(replyList, receivedMessageRule.ReplyContent...)
			} else {
				//随机选择 ReplyNum 条
				//数组中随机切取 ReplyNum 条
				selectReplyList := GetRandomSliceReply(receivedMessageRule.ReplyContent, receivedMessageRule.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}

			break
		}
	}

	//判断是否继续ai
	if len(replyList) == 0 {
		//没有匹配到关键词 继续ai
		return replyList, nil
	}
	//循环replyList标记来源
	for i := range replyList {
		//标记来源
		replyList[i].SendSource = common.RobotAbilityReceivedMessageReply
	}
	//返回消息
	return replyList, nil
}

// 简洁但健壮版（推荐用于通用场景）
func isTodayInDateRange(start, end string) bool {
	today := time.Now()
	sd, _ := time.Parse("2006-01-02", start)
	ed, _ := time.Parse("2006-01-02", end)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, sd.Location())
	return !today.Before(sd) && !today.After(ed)
}

func nowInHHmmRangeSimple(start, end string) bool {
	now := time.Now()
	loc := now.Location()

	// 解析HH:mm到临时时间（日期会被忽略）
	startTime, _ := time.Parse("15:04", start)
	endTime, _ := time.Parse("15:04", end)

	// 构造今天的起止时间（只取时分）
	startT := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, loc)
	endT := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, loc)

	return !now.Before(startT) && !now.After(endT)
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

func GetDialogueSession(params *define.ChatRequestParam) (int, int, error) {
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
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
	}
	if len(chatRequestParam.Question) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(chatRequestParam.Lang, `question_empty`))))
		return
	}
	dialogueId, sessionId, err := GetDialogueSession(chatRequestParam)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	isDraft := cast.ToBool(c.PostForm(`is_draft`))
	workFlowParams := &work_flow.WorkFlowParams{
		ChatRequestParam: chatRequestParam,
		DialogueId:       dialogueId,
		SessionId:        sessionId,
		Draft:            work_flow.Draft{IsDraft: isDraft},
	}
	if isDraft { //运行测试(草稿)
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
		if workFlowParams.Draft.StartNodeKey, _, _, err = work_flow.VerifyWorkFlowNodes(nodeList, chatRequestParam.AdminUserId); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
	}
	_, nodeLogs, err := work_flow.BaseCallWorkFlow(workFlowParams)
	c.String(http.StatusOK, lib_web.FmtJson(nodeLogs, err))
}

func CallLoopWorkFlow(c *gin.Context) {
	chatRequestParam := getChatRequestParam(c)
	if chatRequestParam.Error != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, chatRequestParam.Error))
		return
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
	if workFlowParams.Draft.StartNodeKey, _, _, err = work_flow.VerityLoopWorkflowNodes(chatRequestParam.AdminUserId, loopWorkFlowNode, nodeList); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	_, nodeLogs, err := work_flow.BaseCallWorkFlow(workFlowParams)
	c.String(http.StatusOK, lib_web.FmtJson(nodeLogs, err))
}

func CallLoopWorkFlowParams(c *gin.Context) {
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
	//循环数组
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
	//中间变量
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
	//是否需要问题
	isNeedQuestion := work_flow.FindKeyIsUse(childNodes, `global.question`)
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`loop_test_params`: loopTestParams,
		`is_need_question`: isNeedQuestion,
	}, err))
}
