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

// UnifiedMessageType 统一MsgType
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
	//统一MsgType
	msgType = UnifiedMessageType(msgType)
	message[`MsgType`] = msgType
	//统一消息类型结束
	event := strings.ToLower(cast.ToString(message[`Event`]))

	//点击菜单 关注取关 私信消息 扫描带参数二维码事件
	go work_flow.StartOfficial(message)
	//菜单点击事件处理
	if msgType == lib_define.MsgTypeEvent && tool.InArrayString(event, []string{lib_define.EventMenuClick}) {
		err := MenuClickHandler(message, msg)
		if err != nil {
			logs.Error(`菜单点击处理失败msg:%s,err:%s`, msg, err.Error())
			return err
		}

		return nil
	}

	//关注事件
	if msgType == lib_define.MsgTypeEvent && tool.InArrayString(event, []string{lib_define.EventSubscribe}) {
		err := SubscribeEventHandler(message, msg)
		if err != nil {
			logs.Error(`关注后回复msg:%s,err:%s`, msg, err.Error())
			return err
		}

		return nil
	}

	// 事件过滤 ？？ 只限制某些事件
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
		//其他类型
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
		logs.Error(`welcomes:%s,err:%s`, push.Robot[`welcomes`], err.Error())
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

// MenuClickHandler 菜单点击事件处理
func MenuClickHandler(message map[string]any, msg string) error {
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`appid`]))
	if err != nil {
		logs.Error(`无对应的公众号 msg:%s,err:%s`, msg, err.Error())
	}
	if len(appInfo) == 0 {
		return nil
	}

	eventKey := cast.ToString(message[`EventKey`])
	if tool.IsNumeric(eventKey) {
		menuInfo, err := common.GetOfficialCustomMenuInfo(cast.ToInt(eventKey))
		if err != nil {
			logs.Error(`无发查询到菜单 msg:%s,err:%s`, msg, err.Error())
			return nil
		}

		if menuInfo.ID == 0 {
			logs.Error(`无对应的菜单 msg:%s`, msg)
			return nil
		}
		appAdminUserId := cast.ToInt(appInfo[`admin_user_id`])
		if menuInfo.AdminUserID != appAdminUserId {
			logs.Error(`菜单不是对应公众号的菜单 msg:%s`, msg)
			return nil
		}

		switch menuInfo.ChooseActItem {
		case common.OfficialCustomMenuActTypeSendMessage: //发送消息
			//发送回复消息
			if len(menuInfo.ActParams.ReplyContent) == 0 {
				logs.Error(`菜单无回复消息 msg:%s,err:%s`, msg)
				return nil
			}
			var replyList []common.ReplyContent
			//获取openid
			openid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
			if appInfo[`app_type`] == lib_define.AppWechatKefu { //external_userid integrates with open_kfid
				openid = wechatCommon.GetWechatKefuOpenid(cast.ToInt(appInfo[`admin_user_id`]), message)
			}
			//构建push
			push := &lib_define.PushMessage{
				MsgRaw:      msg,
				Message:     message,
				AdminUserId: cast.ToInt(appInfo[`admin_user_id`]),
				CreateTime:  cast.ToInt(message[`CreateTime`]),
				Openid:      openid,
				Content:     ``,
				AppInfo:     appInfo,
			}
			//构建app进行推送
			app, err := wechat.GetApplication(appInfo)
			if err != nil {
				logs.Error(`初始化app失败 msg:%s,err:%s`, push.MsgRaw, err.Error())
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
				logs.Error(`菜单无回复消息 `)
				return nil
			}
			//发送回复消息
			SendReplyContentList(push, replyList, app)
			break
		case common.OfficialCustomMenuActTypeJumpURL:
			//跳转链接 不需要处理
			break
		}
	}
	return nil
}

// SubscribeEventHandler 关注事件
func SubscribeEventHandler(message map[string]any, msg string) error {
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`appid`]))
	if err != nil {
		logs.Error(`无对应的公众号 msg:%s,err:%s`, msg, err.Error())
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
	//订阅回复
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

// SendSubscribeReply 订阅回复
func SendSubscribeReply(push *lib_define.PushMessage) {
	if len(push.AppInfo) == 0 || push.AppInfo[`app_type`] != lib_define.AppOfficeAccount {
		//不是公众号的关注回复
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

	//关键词回复
	useAbility := common.CheckUseAbilityByAbilityType(push.AdminUserId, common.RobotAbilitySubscribeReply)
	if !useAbility {
		//关键词回复没开启
		logs.Error(`关键词回复功能没开启 msg:%s,err:%s`, push.MsgRaw)
		return
	}
	//获取关注场景
	subscribeScene, err := app.GetSubscribeScene(push.Openid)
	if err != nil {
		logs.Error(`无法获取关注场景 msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	//关注回复的消息
	message, err := common.SubscribeReplyHandle(params, subscribeScene)
	if err != nil {
		logs.Error(`无关注回复消息 msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	//没消息不回复
	if len(message) == 0 {
		return
	}

	//发送回复的消息
	ReplyContentListHandle(push, message, app)
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
	//下载图片
	ImageMediaIdToOssUrl(push, receivedMessageType, app, params)
	//下载缩略图
	ThumbMediaIdToOssUrl(push, app, params)

	//构建成多模态输入数据格式
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
	params.Question = push.Content //将question替换成多模态输入数据格式

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

	//未认证公众号的消息特殊处理
	if params.AppType == lib_define.AppOfficeAccount && params.PassiveId > 0 {
		return //走被动回复,发送消息接口没有权限
	}

	//发送回复的消息
	SendReplyMessageHandle(push, message, app, err, params)
	return
}

func SendReplyMessageHandle(push *lib_define.PushMessage, message msql.Params, app wechat.ApplicationInterface, err error, params *define.ChatRequestParam) bool {
	//判断是否有关键词回复
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

// ReplyContentListHandle 回复处理
func ReplyContentListHandle(push *lib_define.PushMessage, message msql.Params, app wechat.ApplicationInterface) {
	replyContentListJson, isKeyword := message[`reply_content_list`]
	if isKeyword {
		var replyContent []common.ReplyContent
		_ = tool.JsonDecodeUseNumber(replyContentListJson, &replyContent)
		SendReplyContentList(push, replyContent, app)
	}
}

// SendReplyContentList 发送回复内容列表
func SendReplyContentList(push *lib_define.PushMessage, replyContent []common.ReplyContent, app wechat.ApplicationInterface) {
	if len(replyContent) > 0 {
		for _, keywordReply := range replyContent {
			//有关键词回复的处理
			//兼容类型
			checkType := keywordReply.ReplyType
			if checkType == `` && keywordReply.Type != `` {
				checkType = keywordReply.Type
			}
			switch checkType {
			case common.ReplyTypeImageText: //图文
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`图片不存在，url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendImageTextLink(push.Openid, keywordReply.URL, keywordReply.Title, keywordReply.Description, localThumbURL, keywordReply.ThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeText: //文本
				errCode, err := app.SendText(push.Openid, keywordReply.Description, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeUrl: //链接
				errCode, err := app.SendUrl(push.Openid, keywordReply.URL, keywordReply.Title, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeImg: //图片
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`图片不存在，url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendImage(push.Openid, localThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeCard: //小程序卡片
				localThumbURL := common.GetFileByLink(keywordReply.ThumbURL)
				if localThumbURL == `` {
					logs.Error(`图片不存在，url:%s`, keywordReply.ThumbURL)
					break
				}
				errCode, err := app.SendMiniProgramPage(push.Openid, keywordReply.Appid, keywordReply.Title, keywordReply.PagePath, localThumbURL, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break
			case common.ReplyTypeSmartMenu: //智能菜单
				errCode, err := app.SendSmartMenu(push.Openid, keywordReply.SmartMenu, push)
				if err != nil {
					logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errCode, err.Error())
				}
				break

			default:
				//其他类型消息 没指定类型的发送兼容
				break
			}
		}
	}
}

// SendReceivedMessageReply 发送收到消息的回复
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
		logs.Error(`非消息类型，不处理消息：%s`, push.MsgRaw)
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
	//下载图片
	ImageMediaIdToOssUrl(push, receivedMessageType, app, params)
	//下载缩略图
	ThumbMediaIdToOssUrl(push, app, params)

	//记录收到的消息
	message, err := common.OnlyReceivedMessageReply(params)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	//没消息不回复
	if len(message) == 0 {
		return
	}

	//发送回复的消息
	SendReplyMessageHandle(push, message, app, err, params)
	return
}

// ThumbMediaIdToOssUrl 缩略图处理
func ThumbMediaIdToOssUrl(push *lib_define.PushMessage, app wechat.ApplicationInterface, params *define.ChatRequestParam) {
	thumbMediaId := cast.ToString(push.Message[`ThumbMediaId`])
	if thumbMediaId != `` {
		thumbMedia, h, _, err := app.GetFileByMedia(thumbMediaId, push)
		if err != nil {
			logs.Error(`下载缩略图错误 thumbMedia：%s, msg:%s,err:%s`, thumbMedia, push.MsgRaw, err.Error())
			return
		}
		uploadInfo, err := common.SaveImageByMedia(thumbMedia, h, push.AdminUserId, `received_message_images`, define.ImageAllowExt)
		if err != nil {
			logs.Error(`上传缩略图文件获取链接失败：%s, msg:%s,err:%s`, thumbMediaId, push.MsgRaw, err.Error())
			return
		}
		//上传到oss获取链接
		params.ThumbMediaIdToOssUrl = uploadInfo.Link
	}
}

// OfficeAccountImageDownloadPriority 公众号图片下载优先逻辑
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
	//上传到oss获取链接
	params.MediaIdToOssUrl = uploadInfo.Link
	return true
}

// ImageMediaIdToOssUrl 图片消息处理
func ImageMediaIdToOssUrl(push *lib_define.PushMessage, receivedMessageType string, app wechat.ApplicationInterface, params *define.ChatRequestParam) {
	if OfficeAccountImageDownloadPriority(push, receivedMessageType, params) {
		return //公众号图片下载优先逻辑
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
			logs.Error(`下载图片错误 mediaId：%s, msg:%s,err:%s`, mediaId, push.MsgRaw, err.Error())
			return
		}
		uploadInfo, err := common.SaveImageByMedia(media, h, push.AdminUserId, `received_message_images`, define.ImageAllowExt)
		if err != nil {
			logs.Error(`上传文件获取链接失败：%s, msg:%s,err:%s`, mediaId, push.MsgRaw, err.Error())
			return
		}
		//上传到oss获取链接
		params.MediaIdToOssUrl = uploadInfo.Link
	}
}
