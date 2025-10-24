// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
	"github.com/zhimaAi/go_tools/tool"
)

func PushPwd(c *gin.Context) {
	appType := strings.TrimSpace(c.Param(`app_type`))
	if !tool.InArrayString(appType, []string{lib_define.AppOfficeAccount, lib_define.AppMini}) {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	accessKey := strings.TrimSpace(c.Param(`access_key`))
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
	if msgType == `event` && tool.InArrayString(event, []string{`minigame_deliver_goods`, `minigame_coin_deliver_completed`}) {
		echo := lib_define.SuccessJson
		if len(body) > 0 && body[0] == '<' {
			echo = lib_define.SuccessXml
		}
		c.String(http.StatusOK, echo)
		go common.PushNSQ(message)
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
