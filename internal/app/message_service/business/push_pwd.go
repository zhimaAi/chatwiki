// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/pkg/lib_define"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
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

	var robot msql.Params
	//获取机器人信息
	if cast.ToString(appInfo[`robot_key`]) != `` {
		robot, err = common.GetRobotInfo(appInfo[`robot_key`])
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_define.SUCCESS)
			return
		}
	}

	//未认证公众号的文本消息特殊处理
	msgTypes := []string{lib_define.MsgTypeText}
	if len(robot) > 0 && cast.ToBool(robot[`question_multiple_switch`]) {
		msgTypes = append(msgTypes, lib_define.MsgTypeImage)
	}
	if len(robot) > 0 && appType == lib_define.AppOfficeAccount && tool.InArrayString(msgType, msgTypes) && !lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`]) {
		var echo string
		if logid, serial, ok := common.CheckQueryAiReply(message); ok && logid > 0 && serial >= 0 { //查询AI回复
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
