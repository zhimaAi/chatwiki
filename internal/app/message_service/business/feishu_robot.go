// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/feishu_robot"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

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

	respData := lib_define.FeishuMsgEvent{}
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
