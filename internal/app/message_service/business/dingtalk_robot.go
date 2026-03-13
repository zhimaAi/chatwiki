// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/pkg/lib_define"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func DingTalkPush(c *gin.Context) {
	accessKey := strings.TrimSpace(c.Param(`access_key`))
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)

	body, err := io.ReadAll(c.Request.Body)

	respData := lib_define.DingtalkMsgEvent{}
	err = tool.JsonDecodeUseNumber(string(body), &respData)
	if err != nil {
		logs.Error("错误了：" + err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	//组装消息推送
	nsqMsg := map[string]interface{}{}
	nsqMsg[`appid`] = appInfo[`app_id`]
	nsqMsg[`ToUserName`] = respData.ConversationId
	nsqMsg[`FromUserName`] = respData.SenderStaffId
	nsqMsg[`CreateTime`] = respData.CreateAt
	nsqMsg[`MsgType`] = respData.Msgtype
	nsqMsg[`Content`] = respData.Text.Content
	nsqMsg[`MsgId`] = respData.MsgId
	nsqMsg[`RobotCode`] = respData.RobotCode
	nsqMsg[`SenderNick`] = respData.SenderNick
	nsqMsg[`SessionType`] = respData.ConversationType //1：单聊，2：群聊

	if respData.Msgtype == lib_define.DingTalkMsgTypeImage { //图片内容
		nsqMsg[`Content`], _ = tool.JsonEncode(respData.Content)
	}

	go common.PushNSQ(nsqMsg)

	c.JSON(http.StatusOK, nil)
	return
}
