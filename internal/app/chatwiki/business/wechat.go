// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/wechat"
	wechatCommon "chatwiki/internal/pkg/wechat/common"
	"chatwiki/internal/pkg/wechat/feishu_robot"
	"chatwiki/internal/pkg/wechat/official_account"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// UnifiedMessageType unified MsgType
func UnifiedMessageType(msgType string) string {
	switch msgType {
	case lib_define.DingTalkMsgTypeImage:
		msgType = lib_define.MsgTypeImage
		break
	case lib_define.FeShuMsgTypeAudio:
		msgType = lib_define.MsgTypeVoice
		break
	case lib_define.FeShuMsgTypeMedia:
	case lib_define.DingTalkMsgTypeVideo:
		msgType = lib_define.MsgTypeVideo
		break
	}
	return msgType
}

func AppPush(msg string, _ ...string) error {
	//parse message body
	message := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(msg, &message); err != nil {
		logs.Error(`msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	//discard what is not needed
	msgType := strings.ToLower(cast.ToString(message[`MsgType`]))
	// unified MsgType
	msgType = UnifiedMessageType(msgType)
	message[`MsgType`] = msgType
	// end of unified message type
	event := strings.ToLower(cast.ToString(message[`Event`]))

	// click menu, subscribe/unsubscribe, private message, scan QR code with parameters event
	go work_flow.StartOfficial(message)
	// menu click event handling
	if msgType == lib_define.MsgTypeEvent && tool.InArrayString(event, []string{lib_define.EventMenuClick}) {
		err := MenuClickHandler(message, msg)
		if err != nil {
			logs.Error(`menu click handler failed msg:%s,err:%s`, msg, err.Error())
			return err
		}

		return nil
	}

	// subscribe event
	if msgType == lib_define.MsgTypeEvent && tool.InArrayString(event, []string{lib_define.EventSubscribe}) {
		err := SubscribeEventHandler(message, msg)
		if err != nil {
			logs.Error(`subscribe event handler failed msg:%s,err:%s`, msg, err.Error())
			return err
		}

		return nil
	}

	// event filtering - only restrict certain events
	if msgType == lib_define.MsgTypeEvent && !tool.InArrayString(event, []string{lib_define.EventEnterSession, lib_define.EventUserEnterTempsession}) {
		return nil
	}
	//check app
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`appid`]))
	if err != nil {
		logs.Error(`msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	if len(appInfo) == 0 {
		return nil
	}
	//inspection robot
	robot, err := common.GetRobotInfo(appInfo[`robot_key`])
	if err != nil {
		logs.Error(`msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	if len(robot) == 0 {
		return nil
	}
	//get parameter
	openid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
	if appInfo[`app_type`] == lib_define.AppWechatKefu { //external_userid integrates with open_kfid
		openid = wechatCommon.GetWechatKefuOpenid(cast.ToInt(robot[`admin_user_id`]), message)
	}
	push := &lib_define.PushMessage{
		MsgRaw:      msg,
		Message:     message,
		AdminUserId: cast.ToInt(robot[`admin_user_id`]),
		CreateTime:  cast.ToInt(message[`CreateTime`]),
		Openid:      openid,
		Content:     strings.TrimSpace(cast.ToString(message[`Content`])),
		AppInfo:     appInfo,
		Robot:       robot,
	}
	//save customer
	upData := msql.Datas{}
	if appInfo[`app_type`] == lib_define.AppWechatKefu &&
		lib_redis.AddLock(define.Redis, define.LockPreKey+`update_wxkf_customer.`+openid, time.Hour) {
		if app, err := wechat.GetApplication(push.AppInfo); err == nil {
			if info, _, err := app.GetCustomerInfo(openid); err == nil && len(info) > 0 {
				upData[`nickname`] = info[`nickname`]
				upData[`name`] = info[`nickname`]
				upData[`avatar`] = info[`avatar`]
			}
		}
	}
	common.InsertOrUpdateCustomer(push.Openid, push.AdminUserId, upData)
	customer, err := common.GetCustomerInfo(push.Openid, push.AdminUserId)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	if len(customer) == 0 {
		return nil
	}
	push.Customer = customer
	//processing response
	if define.IsDev {
		logs.Debug(`message:%s`, msg)
	}
	switch msgType {
	case lib_define.MsgTypeEvent:
		switch event {
		case lib_define.EventEnterSession:
			SendWelcome(push)
			break
		case lib_define.EventUserEnterTempsession:
			SendWelcome(push)
			break
		}
		break
	case lib_define.MsgTypeText, lib_define.MsgTypeImage, lib_define.MsgTypeVoice, lib_define.MsgTypeVideo:
		SendReply(push)
		break
	default:
		// other type
		SendReceivedMessageReply(push)
		break
	}
	return nil
}

func BuildSendMenu(menuJsonStr string) (string, error) {
	menuJson := define.MenuJsonStruct{}
	err := tool.JsonDecode(menuJsonStr, &menuJson)
	if err != nil {
		return ``, err
	}
	content := menuJson.Content
	if len(menuJson.Question) > 0 {
		if len(content) > 0 {
			content += "\r\n\r\n"
		}
		for _, question := range menuJson.Question {
			content += fmt.Sprintf(`<a href="weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=0">%s</a>`+"\r\n", question, question)
		}
	}
	return content, nil
}

func SendWelcome(push *lib_define.PushMessage) {
	if push.Robot == nil && len(push.Robot[`welcomes`]) == 0 {
		return
	}
	content, err := BuildSendMenu(push.Robot[`welcomes`])
	if err != nil {
		logs.Error(`build send menu failed welcomes:%s,err:%s`, push.Robot[`welcomes`], err.Error())
		return
	}
	app, err := wechat.GetApplication(push.AppInfo)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	//wechat_kefu special treatment
	welcomeCode := cast.ToString(push.Message[`welcome_code`])
	if len(push.AppInfo) > 0 && push.AppInfo[`app_type`] == lib_define.AppWechatKefu && len(welcomeCode) > 0 {
		if _, err := app.SendMsgOnEvent(welcomeCode, content); err == nil {
			return
		}
	}
	errcode, err := app.SendText(push.Openid, content, push)
	if err != nil {
		logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
		return
	}
}

// MenuClickHandler menu click event handling
func MenuClickHandler(message map[string]any, msg string) error {
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`appid`]))
	if err != nil {
		logs.Error(`no corresponding official account msg:%s,err:%s`, msg, err.Error())
	}
	if len(appInfo) == 0 {
		return nil
	}

	eventKey := cast.ToString(message[`EventKey`])
	if tool.IsNumeric(eventKey) {
		menuInfo, err := common.GetOfficialCustomMenuInfo(cast.ToInt(eventKey))
		if err != nil {
			logs.Error(`failed to query menu msg:%s,err:%s`, msg, err.Error())
			return nil
		}

		if menuInfo.ID == 0 {
			logs.Error(`no corresponding menu msg:%s`, msg)
			return nil
		}
		appAdminUserId := cast.ToInt(appInfo[`admin_user_id`])
		if menuInfo.AdminUserID != appAdminUserId {
			logs.Error(`menu does not belong to this official account msg:%s`, msg)
			return nil
		}

		switch menuInfo.ChooseActItem {
		case common.OfficialCustomMenuActTypeSendMessage: // send message
			// send reply message
			if len(menuInfo.ActParams.ReplyContent) == 0 {
				logs.Error(`menu has no reply content msg:%s,err:%s`, msg)
				return nil
			}
			var replyList []common.ReplyContent
			// get openid
			openid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
			if appInfo[`app_type`] == lib_define.AppWechatKefu { //external_userid integrates with open_kfid
				openid = wechatCommon.GetWechatKefuOpenid(cast.ToInt(appInfo[`admin_user_id`]), message)
			}
			// build push
			push := &lib_define.PushMessage{
				MsgRaw:      msg,
				Message:     message,
				AdminUserId: cast.ToInt(appInfo[`admin_user_id`]),
				CreateTime:  cast.ToInt(message[`CreateTime`]),
				Openid:      openid,
				Content:     ``,
				AppInfo:     appInfo,
			}
			// build app for pushing
			app, err := wechat.GetApplication(appInfo)
			if err != nil {
				logs.Error(`init app failed msg:%s,err:%s`, push.MsgRaw, err.Error())
				return nil
			}
			if menuInfo.ActParams.ReplyNum == 0 {
				replyList = menuInfo.ActParams.ReplyContent
			} else {
				selectReplyList := common.GetRandomSliceReply(menuInfo.ActParams.ReplyContent, menuInfo.ActParams.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
			replyList = common.FormatReplyListToDb(replyList, common.OfficialAbilityCustomMenu)
			if len(replyList) == 0 {
				logs.Error(`menu has no reply message`)
				return nil
			}
			// send reply message
			SendReplyContentList(push, replyList, app)
			break
		case common.OfficialCustomMenuActTypeJumpURL:
			// jump link no need to handle
			break
		}
	}
	return nil
}

// SubscribeEventHandler subscribe event
func SubscribeEventHandler(message map[string]any, msg string) error {
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`appid`]))
	if err != nil {
		logs.Error(`no corresponding official account msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	if len(appInfo) == 0 {
		return nil
	}

	//get parameter
	openid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
	if appInfo[`app_type`] == lib_define.AppWechatKefu { //external_userid integrates with open_kfid
		openid = wechatCommon.GetWechatKefuOpenid(cast.ToInt(appInfo[`admin_user_id`]), message)
	}
	// subscription reply
	push := &lib_define.PushMessage{
		MsgRaw:      msg,
		Message:     message,
		AdminUserId: cast.ToInt(appInfo[`admin_user_id`]),
		CreateTime:  cast.ToInt(message[`CreateTime`]),
		Openid:      openid,
		Content:     strings.TrimSpace(cast.ToString(message[`Content`])),
		AppInfo:     appInfo,
	}

	SendSubscribeReply(push)
	return nil
}

// SendSubscribeReply subscription reply
func SendSubscribeReply(push *lib_define.PushMessage) {
	if len(push.AppInfo) == 0 || push.AppInfo[`app_type`] != lib_define.AppOfficeAccount {
		// not a public account subscription reply
		return
	}
	app := &official_account.Application{AppID: push.AppInfo[`app_id`], Secret: push.AppInfo[`app_secret`]}
	params := &define.ChatRequestParam{
		ChatBaseParam: &define.ChatBaseParam{
			AppType:     push.AppInfo[`app_type`],
			AppInfo:     push.AppInfo,
			Openid:      push.Openid,
			AdminUserId: push.AdminUserId,
			Robot:       push.Robot,
			Customer:    push.Customer,
		},
		ReceivedMessageType: lib_define.EventSubscribe,
		Question:            push.Content,
	}

	// keyword reply
	useAbility := common.CheckUseAbilityByAbilityType(push.AdminUserId, common.RobotAbilitySubscribeReply)
	if !useAbility {
		// keyword reply not enabled
		logs.Error(`keyword reply feature not enabled msg:%s,err:%s`, push.MsgRaw)
		return
	}
	// get subscription scene
	subscribeScene, err := app.GetSubscribeScene(push.Openid)
	if err != nil {
		logs.Error(`failed to get subscribe scene msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	// subscription reply message
	message, err := common.SubscribeReplyHandle(params, subscribeScene)
	if err != nil {
		logs.Error(`no subscribe reply message msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	// no message no reply
	if len(message) == 0 {
		return
	}

	// send reply message
	ReplyContentListHandle(push, message, app)
	return
}

// ShowTypingStatusToUser authenticated WeChat public account + mini program, show typing status
func ShowTypingStatusToUser(appType string, appInfo, robot msql.Params) (showTyping bool) {
	switch appType {
	case lib_define.AppOfficeAccount:
		if lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`]) {
			showTyping = cast.ToBool(robot[`show_typing_gzh`])
		}
	case lib_define.AppMini:
		showTyping = cast.ToBool(robot[`show_typing_mini`])
	}
	return
}

func SendReply(push *lib_define.PushMessage) {
	app, err := wechat.GetApplication(push.AppInfo)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	isClose := false
	receivedMessageType := strings.ToLower(cast.ToString(push.Message[`MsgType`]))
	receivedMessageType = UnifiedMessageType(receivedMessageType)

	params := &define.ChatRequestParam{
		ChatBaseParam: &define.ChatBaseParam{
			AppType:     push.AppInfo[`app_type`],
			AppInfo:     push.AppInfo,
			Openid:      push.Openid,
			AdminUserId: push.AdminUserId,
			Robot:       push.Robot,
			Customer:    push.Customer,
		},
		ReceivedMessage:     push.Message,
		ReceivedMessageType: receivedMessageType,
		Question:            push.Content,
		MsgId:               cast.ToString(push.Message[`MsgId`]),
		PassiveId:           cast.ToInt64(push.Message[`passive_id`]),
		DialogueId:          common.GetLastDialogueId(push.AdminUserId, cast.ToInt(push.Robot[`id`]), push.Openid),
		IsClose:             &isClose,
	}
	//specify the language to use based on the content
	if common.IsContainChinese(push.Content) {
		params.Lang = define.LangZhCn
	}
	// download image
	ImageMediaIdToOssUrl(push, receivedMessageType, app, params)
	// download thumbnail
	ThumbMediaIdToOssUrl(push, app, params)

	// construct to multimodal input data format
	switch receivedMessageType {
	case lib_define.MsgTypeImage:
		push.Content = tool.JsonEncodeNoError(adaptor.QuestionMultiple{
			{Type: adaptor.TypeImage, ImageUrl: adaptor.ImageUrl{Url: params.MediaIdToOssUrl}},
		})
	case lib_define.MsgTypeVoice:
		push.Content = tool.JsonEncodeNoError(adaptor.QuestionMultiple{
			{Type: adaptor.TypeAudio, InputAudio: adaptor.InputAudio{Data: params.MediaIdToOssUrl}},
		})
	case lib_define.MsgTypeVideo:
		push.Content = tool.JsonEncodeNoError(adaptor.QuestionMultiple{
			{Type: adaptor.TypeVideo, VedioUrl: adaptor.VedioUrl{Url: params.MediaIdToOssUrl}},
		})
	}
	params.Question = push.Content // replace question with multimodal input data format

	// authenticated WeChat public account + mini program, show typing status
	if ShowTypingStatusToUser(params.AppType, params.AppInfo, params.Robot) {
		if errCode, e := app.SetTyping(push.Openid, lib_define.CommandTyping); e != nil {
			logs.Error(`customer:%s,errcode:%d,err:%s`, push.Openid, errCode, e.Error())
		}
	}

	chanStream := make(chan sse.Event)
	go func(chanStream chan sse.Event) {
		for event := range chanStream {
			if define.IsDev {
				event.Data, _ = tool.JsonEncode(event.Data)
				logs.Debug(`event:%v`, event)
			}
		}
	}(chanStream)
	message, err := DoChatRequest(params, false, chanStream)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	if len(message) == 0 {
		return
	}

	// special handling for unauthenticated public account messages
	if params.AppType == lib_define.AppOfficeAccount && params.PassiveId > 0 {
		return // passive reply, sending message interface has no permission
	}

	// send reply message
	SendReplyMessageHandle(push, message, app, err, params)
	return
}

func SendReplyMessageHandle(push *lib_define.PushMessage, message msql.Params, app wechat.ApplicationInterface, err error, params *define.ChatRequestParam) bool {
	// check if there are keyword replies
	ReplyContentListHandle(push, message, app)

	var content string
	switch cast.ToInt(message[`msg_type`]) {
	case define.MsgTypeText:
		content = message[`content`]
	case define.MsgTypeMenu:
		content, err = BuildSendMenu(message[`menu_json`])
		if err != nil {
			logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
			return true
		}
	}
	if len(content) == 0 {
		return true
	}
	text, images, voices := common.GetMessageInMessage(content, true)
	if len(text) > 0 {
		errcode, err := app.SendText(push.Openid, text, push)
		if err != nil {
			logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
		}
	}
	if len(images) > 0 {
		for _, image := range images {
			errcode, err := app.SendImage(push.Openid, image, push)
			if err != nil {
				logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
			}
		}
	}
	if len(voices) > 0 && tool.InArray(params.AppType, []string{lib_define.AppWechatKefu, lib_define.AppOfficeAccount}) {
		for _, voice := range voices {
			ext := strings.ToLower(filepath.Ext(voice))
			if !tool.InArray(ext, []string{`mp3`, `amr`}) {
				logs.Warning(`voice is not mp3 or amr ,%s`, voice)
			}
			if params.AppType == lib_define.AppWechatKefu && ext == `.mp3` {
				voice = common.Mp3ToAmr(params.AdminUserId, voice)
			}
			errcode, err := app.SendVoice(push.Openid, voice, push)
			if err != nil {
				logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
			}
		}
	}

	return false
}

// ReplyContentListHandle reply handling
func ReplyContentListHandle(push *lib_define.PushMessage, message msql.Params, app wechat.ApplicationInterface) {
	replyContentListJson, isKeyword := message[`reply_content_list`]
	if isKeyword {
		var replyContent []common.ReplyContent
		_ = tool.JsonDecodeUseNumber(replyContentListJson, &replyContent)
		SendReplyContentList(push, replyContent, app)
	}
}

// SendReplyContentList send reply content list
func SendReplyContentList(push *lib_define.PushMessage, replyContent []common.ReplyContent, app wechat.ApplicationInterface) {
	if len(replyContent) > 0 {
		for _, keywordReply := range replyContent {
			// handle keyword replies
			// compatible type
			checkType := keywordReply.ReplyType
			if checkType == `` && keywordReply.Type != `` {
				checkType = keywordReply.Type
			}
			switch checkType {
			case common.ReplyTypeImageText: // image-text
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`image does not exist url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendImageTextLink(push.Openid, keywordReply.URL, keywordReply.Title, keywordReply.Description, localThumbURL, keywordReply.ThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeText: // text
				errCode, err := app.SendText(push.Openid, keywordReply.Description, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeUrl: // link
				errCode, err := app.SendUrl(push.Openid, keywordReply.URL, keywordReply.Title, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeImg: // image
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`image does not exist url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendImage(push.Openid, localThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeCard: // mini program card
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`image does not exist url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendMiniProgramPage(push.Openid, keywordReply.Appid, keywordReply.Title, keywordReply.PagePath, localThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeSmartMenu: // smart menu
				errCode, err := app.SendSmartMenu(push.Openid, keywordReply.SmartMenu, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break

			default:
				// other message types, send compatible if no specific type
				break
			}
		}
	}
}

// SendReceivedMessageReply send received message reply
func SendReceivedMessageReply(push *lib_define.PushMessage) {
	app, err := wechat.GetApplication(push.AppInfo)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	isClose := false
	receivedMessageType := strings.ToLower(cast.ToString(push.Message[`MsgType`]))
	receivedMessageType = UnifiedMessageType(receivedMessageType)
	if receivedMessageType == `` {
		logs.Error(`not a message type, skip processing:%s`, push.MsgRaw)
		return
	}

	params := &define.ChatRequestParam{
		ChatBaseParam: &define.ChatBaseParam{
			AppType:     push.AppInfo[`app_type`],
			AppInfo:     push.AppInfo,
			Openid:      push.Openid,
			AdminUserId: push.AdminUserId,
			Robot:       push.Robot,
			Customer:    push.Customer,
		},
		ReceivedMessage:     push.Message,
		ReceivedMessageType: receivedMessageType,
		Question:            push.Content,
		MsgId:               cast.ToString(push.Message[`MsgId`]),
		DialogueId:          common.GetLastDialogueId(push.AdminUserId, cast.ToInt(push.Robot[`id`]), push.Openid),
		IsClose:             &isClose,
	}
	//specify the language to use based on the content
	if common.IsContainChinese(push.Content) {
		params.Lang = define.LangZhCn
	}
	// download image
	ImageMediaIdToOssUrl(push, receivedMessageType, app, params)
	// download thumbnail
	ThumbMediaIdToOssUrl(push, app, params)

	// record received message
	message, err := common.OnlyReceivedMessageReply(params)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	// no message no reply
	if len(message) == 0 {
		return
	}

	// send reply message
	SendReplyMessageHandle(push, message, app, err, params)
	return
}

// ThumbMediaIdToOssUrl thumbnail processing
func ThumbMediaIdToOssUrl(push *lib_define.PushMessage, app wechat.ApplicationInterface, params *define.ChatRequestParam) {
	thumbMediaId := cast.ToString(push.Message[`ThumbMediaId`])
	if thumbMediaId != `` {
		thumbMedia, h, _, err := app.GetFileByMedia(thumbMediaId, push)
		if err != nil {
			logs.Error(`download thumbnail error thumbmedia:%s, msg:%s,err:%s`, thumbMedia, push.MsgRaw, err.Error())
			return
		}
		uploadInfo, err := common.SaveImageByMedia(thumbMedia, h, push.AdminUserId, `received_message_images`, define.ImageAllowExt)
		if err != nil {
			logs.Error(`upload thumbnail file failed:%s, msg:%s,err:%s`, thumbMediaId, push.MsgRaw, err.Error())
			return
		}
		// upload to oss to get link
		params.ThumbMediaIdToOssUrl = uploadInfo.Link
	}
}

// OfficeAccountImageDownloadPriority public account image download priority logic
func OfficeAccountImageDownloadPriority(push *lib_define.PushMessage, receivedMessageType string, params *define.ChatRequestParam) (ok bool) {
	if push.AppInfo[`app_type`] != lib_define.AppOfficeAccount || receivedMessageType != lib_define.MsgTypeImage {
		return
	}
	picUrl := cast.ToString(push.Message[`PicUrl`])
	if len(picUrl) == 0 || !common.IsUrl(picUrl) {
		return
	}
	request := curl.Get(picUrl)
	resp, err := request.Response()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	media, err := request.Bytes()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	uploadInfo, err := common.SaveImageByMedia(media, resp.Header, push.AdminUserId, `received_message_images`, define.ImageAllowExt)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	// upload to oss to get link
	params.MediaIdToOssUrl = uploadInfo.Link
	return true
}

// ImageMediaIdToOssUrl image message processing
func ImageMediaIdToOssUrl(push *lib_define.PushMessage, receivedMessageType string, app wechat.ApplicationInterface, params *define.ChatRequestParam) {
	if OfficeAccountImageDownloadPriority(push, receivedMessageType, params) {
		return // public account image download priority logic
	}
	mediaId := cast.ToString(push.Message[`MediaId`])
	if push.AppInfo[`app_type`] == lib_define.FeiShuRobot && receivedMessageType == lib_define.MsgTypeImage {
		content := feishu_robot.FeiShuImgMsgContent{}
		_ = tool.JsonDecodeUseNumber(cast.ToString(push.Message[`Content`]), &content)
		mediaId = content.ImageKey
	}
	if push.AppInfo[`app_type`] == lib_define.DingTalkRobot && receivedMessageType == lib_define.MsgTypeImage {
		msgContent := lib_define.DingtalkImgContent{}
		_ = tool.JsonDecode(cast.ToString(push.Message["Content"]), &msgContent)
		mediaId = msgContent.DownloadCode
	}
	if receivedMessageType == lib_define.MsgTypeImage && mediaId != `` {
		media, h, _, err := app.GetFileByMedia(mediaId, push)
		if err != nil {
			logs.Error(`download image error mediaid:%s, msg:%s,err:%s`, mediaId, push.MsgRaw, err.Error())
			return
		}
		uploadInfo, err := common.SaveImageByMedia(media, h, push.AdminUserId, `received_message_images`, define.ImageAllowExt)
		if err != nil {
			logs.Error(`upload file failed to get link:%s, msg:%s,err:%s`, mediaId, push.MsgRaw, err.Error())
			return
		}
		// upload to oss to get link
		params.MediaIdToOssUrl = uploadInfo.Link
	}
}
