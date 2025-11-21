// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/feishu_robot"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func PushPwd(c *gin.Context) {
	appType := strings.TrimSpace(c.Param(`app_type`))
	if !tool.InArrayString(appType, []string{lib_define.AppOfficeAccount, lib_define.AppMini}) {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	accessKey := strings.TrimSpace(c.Param(`access_key`))
	//check app
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(appInfo) == 0 || appInfo[`app_type`] != appType {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	//inspection robot
	robot, err := common.GetRobotInfo(appInfo[`robot_key`])
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	signature := c.Query(`signature`)
	timestamp := c.Query(`timestamp`)
	nonce := c.Query(`nonce`)
	if echostr := c.Query(`echostr`); len(echostr) > 0 {
		if common.VerifySignature(signature, timestamp, nonce) {
			c.String(http.StatusOK, echostr)
		} else {
			c.String(http.StatusOK, `signature verification failure`)
		}
		return
	}
	encryptType := c.Query(`encrypt_type`)
	msgSignature := c.Query(`msg_signature`)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logs.Error(err.Error())
	}
	//body is empty
	if len(body) == 0 || len(strings.TrimSpace(string(body))) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		logs.Notice(strings.Join([]string{c.ClientIP(), c.Request.Method, c.Request.Host, c.Request.RequestURI}, ` | `))
		return
	}
	message := common.GetMessage(string(body), signature, encryptType, msgSignature, nonce, timestamp)
	if tool.InArrayString(strings.ToUpper(cast.ToString(message[`Event`])), lib_define.DiscardEvents) {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	//inject the required parameters
	message[`appid`] = appInfo[`app_id`]
	//special recovery structure
	msgType := strings.ToLower(cast.ToString(message[`MsgType`]))
	event := strings.ToLower(cast.ToString(message[`Event`]))
	if msgType == lib_define.MsgTypeEvent && tool.InArrayString(event, []string{`minigame_deliver_goods`, `minigame_coin_deliver_completed`}) {
		echo := lib_define.SuccessJson
		if len(body) > 0 && body[0] == '<' {
			echo = lib_define.SuccessXml
		}
		c.String(http.StatusOK, echo)
		go common.PushNSQ(message)
		return
	}
	//未认证公众号的文本消息特殊处理
	if appType == lib_define.AppOfficeAccount && msgType == lib_define.MsgTypeText && !lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`]) {
		var echo string
		if logid, serial, ok := common.CheckQueryAiReply(message); ok { //查询AI回复
			echo = common.GetAiReply(robot, message, logid, serial)
		} else { //普通消息处理
			echo = common.WaitAiReply(robot, message)
		}
		if len(echo) > 0 { //被动回复
			c.String(http.StatusOK, common.BuildXmlStr(appInfo, message, echo))
		} else { //失败了的回复success
			c.String(http.StatusOK, lib_define.SUCCESS)
		}
		return
	}
	c.String(http.StatusOK, lib_define.SUCCESS)
	go common.PushNSQ(message)
}

func WechatKefu(c *gin.Context) {
	msgSignature := c.Query(`msg_signature`)
	timestamp := c.Query(`timestamp`)
	nonce := c.Query(`nonce`)
	if echostr := c.Query(`echostr`); len(echostr) > 0 {
		if replyEchoStr, err := common.MsgDecrypt(echostr, msgSignature, nonce, timestamp); err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, replyEchoStr)
		}
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logs.Error(err.Error())
	}
	//body is empty
	if len(body) == 0 || len(strings.TrimSpace(string(body))) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		logs.Notice(strings.Join([]string{c.ClientIP(), c.Request.Method, c.Request.Host, c.Request.RequestURI}, ` | `))
		return
	}
	message := common.GetMessage(string(body), `signature`, `aes`, msgSignature, nonce, timestamp)
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`ToUserName`]))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(appInfo) == 0 || appInfo[`app_type`] != lib_define.AppWechatKefu {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	switch cast.ToString(message[`MsgType`]) {
	case `event`:
		switch cast.ToString(message[`Event`]) {
		case `kf_msg_or_event`:
			common.SetMessageLasttimeCache(cast.ToString(message[`ToUserName`]), cast.ToInt(message[`CreateTime`]))
			common.SetSyncMsgTokenCache(cast.ToString(message[`ToUserName`]), cast.ToString(message[`Token`]))
			go common.GetKfMsgOrEvent(appInfo)
		default:
			logs.Info(`%+v`, message)
		}
	default:
		logs.Info(`%+v`, message)
	}
	c.String(http.StatusOK, lib_define.SUCCESS)
}

type reqPostData struct {
	Encrypt string `json:"encrypt"`
}

func FeishuPush(c *gin.Context) {

	accessKey := strings.TrimSpace(c.Param(`access_key`))
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	//消息体解析失败
	reqData := reqPostData{}
	err = tool.JsonDecode(string(body), &reqData)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	//消息解析
	message, err := common.GetFeiShuMessage(reqData.Encrypt, appInfo["encrypt_key"])

	respData := define.FeishuMsgEvent{}
	err = tool.JsonDecodeUseNumber(message, &respData)
	if err != nil {
		logs.Error("错误了：" + err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	if respData.Type == "url_verification" { //回调链接验证，直接退出
		c.JSON(http.StatusOK, respData)
		return
	}

	toUserName := ``
	fromUserName := ``
	if respData.Event.Message.ChatId != `` {
		//消息类型
		toUserName = respData.Event.Message.ChatId
		fromUserName = respData.Event.Sender.SenderId.OpenId

		if fromUserName == `` {
			logs.Error("没有用户信息错误了：" + message)
			c.String(http.StatusOK, lib_define.SUCCESS)
			return
		}

		if respData.Event.Message.MessageType == lib_define.FeShuMsgTypeText { //非文本消息，将原文json推送
			feishuMsgContent := feishu_robot.FeiShuTextMsgContent{}
			err = tool.JsonDecode(respData.Event.Message.Content, &feishuMsgContent)
			if err != nil {
				logs.Error("错误了：" + err.Error())
				c.JSON(http.StatusOK, respData)
				return
			}
			respData.Event.Message.Content = feishuMsgContent.Text
		}

		//组装消息推送
		nsqMsg := map[string]interface{}{}
		nsqMsg[`appid`] = appInfo[`app_id`]
		nsqMsg[`ToUserName`] = toUserName
		nsqMsg[`FromUserName`] = fromUserName
		nsqMsg[`CreateTime`] = respData.Event.Message.CreateTime
		nsqMsg[`MsgType`] = respData.Event.Message.MessageType
		nsqMsg[`Content`] = respData.Event.Message.Content
		nsqMsg[`MsgId`] = respData.Event.Message.MessageId
		nsqMsg[`SessionType`] = respData.Event.Message.ChatType

		go common.PushNSQ(nsqMsg)

	} else if respData.Event.ChatID != `` {
		//事件类型
		toUserName = respData.Event.ChatID
		fromUserName = respData.Event.OperatorID.OpenID

		if fromUserName == `` {
			logs.Error("没有用户信息错误了：" + message)
			c.String(http.StatusOK, lib_define.SUCCESS)
			return
		}
		nsqMsg := map[string]interface{}{}
		nsqMsg[`appid`] = appInfo[`app_id`]
		nsqMsg[`ToUserName`] = toUserName
		nsqMsg[`FromUserName`] = fromUserName
		nsqMsg[`CreateTime`] = respData.Header.CreateTime
		nsqMsg[`MsgType`] = respData.Header.EventType
		nsqMsg[`Content`] = ``
		nsqMsg[`MsgId`] = respData.Header.EventId
		nsqMsg[`SessionType`] = `event`

		// go common.PushNSQ(nsqMsg)
		logs.Info("飞书事件消息暂时不发送给nsq：" + message)
	}

	c.JSON(http.StatusOK, respData)
	return
}
