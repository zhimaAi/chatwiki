// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/wechat"
	wechatCommon "chatwiki/internal/app/chatwiki/wechat/common"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func AppPush(msg string, _ ...string) error {
	//parse message body
	message := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(msg, &message); err != nil {
		logs.Error(`msg:%s,err:%s`, msg, err.Error())
		return nil
	}
	//discard what is not needed
	msgType := strings.ToLower(cast.ToString(message[`MsgType`]))
	event := strings.ToLower(cast.ToString(message[`Event`]))
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
	push := &define.PushMessage{
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
		SendWelcome(push)
	case lib_define.MsgTypeText:
		SendReply(push)
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

func SendWelcome(push *define.PushMessage) {
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
	errcode, err := app.SendText(push.Openid, content)
	if err != nil {
		logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
		return
	}
}

func SendReply(push *define.PushMessage) {
	app, err := wechat.GetApplication(push.AppInfo)
	if err != nil {
		logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
		return
	}
	isClose := false
	params := &define.ChatRequestParam{
		ChatBaseParam: &define.ChatBaseParam{
			AppType:     push.AppInfo[`app_type`],
			AppInfo:     push.AppInfo,
			Openid:      push.Openid,
			AdminUserId: push.AdminUserId,
			Robot:       push.Robot,
			Customer:    push.Customer,
		},
		Question:   push.Content,
		MsgId:      cast.ToString(push.Message[`MsgId`]),
		DialogueId: common.GetLastDialogueId(push.AdminUserId, cast.ToInt(push.Robot[`id`]), push.Openid),
		IsClose:    &isClose,
	}
	//specify the language to use based on the content
	if common.IsContainChinese(push.Content) {
		params.Lang = define.LangZhCn
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

	var content string
	switch cast.ToInt(message[`msg_type`]) {
	case define.MsgTypeText:
		content = message[`content`]
	case define.MsgTypeMenu:
		content, err = BuildSendMenu(message[`menu_json`])
		if err != nil {
			logs.Error(`msg:%s,err:%s`, push.MsgRaw, err.Error())
			return
		}
	}
	if len(content) == 0 {
		return
	}

	text, images := common.GetImgInMessage(content, true)
	if len(text) > 0 {
		errcode, err := app.SendText(push.Openid, text)
		if err != nil {
			logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
		}
	}
	if len(images) > 0 {
		for _, image := range images {
			errcode, err := app.SendImage(push.Openid, image)
			if err != nil {
				logs.Error(`msg:%s,errcode:%d,err:%s`, push.MsgRaw, errcode, err.Error())
			}
		}
	}
	return
}
