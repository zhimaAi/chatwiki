// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/pkg/lib_define"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func WecomRobot(c *gin.Context) {
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
	appInfo, err := common.GetWechatAppInfo(`app_id`, cast.ToString(message[`aibotid`]))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(appInfo) == 0 || appInfo[`app_type`] != lib_define.AppWecomRobot {
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
	//将消息降维处理
	msgtype := cast.ToString(message[`msgtype`])
	if info, ok := message[msgtype].(map[string]any); ok {
		for key, val := range info {
			message[key] = val
		}
		delete(message, msgtype)
	} else {
		logs.Error(`消息结构异常:%v`, message)
	}
	if from, ok := message[`from`].(map[string]any); ok {
		for key, val := range from {
			message[key] = val
		}
		delete(message, `from`)
	}
	//start processing the message
	adminUserId := cast.ToInt(appInfo[`admin_user_id`])
	nsqMsg := map[string]any{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: tool.MD5(cast.ToString(message[`userid`])),
		`CreateTime`:   timestamp,
		`MsgType`:      msgtype,
		`MsgId`:        cast.ToString(message[`msgid`]),
		`appid`:        appInfo[`app_id`],
		`nickname`:     cast.ToString(message[`userid`]),
		`chattype`:     message[`chattype`],     //会话类型，single\group，分别表示：单聊\群聊
		`response_url`: message[`response_url`], //支持主动回复消息的临时url
	}
	switch msgtype {
	case `event`:
		eventtype := cast.ToString(message[`eventtype`])
		switch eventtype {
		case `enter_chat`: //使用被动回复换欢迎语
			content := common.BuildWelcomesReply(robot[`welcomes`])
			if len(strings.TrimSpace(content)) == 0 {
				c.String(http.StatusOK, lib_define.SUCCESS)
				return
			}
			msg := map[string]any{`msgtype`: `text`, `text`: map[string]any{`content`: content}}
			reply, err := common.WecomMsgEncrypt(msg, ``)
			if err != nil {
				logs.Error(err.Error())
			}
			c.JSON(http.StatusOK, reply)
			return
		default:
			logs.Warning(`wecom_robot unsupported eventtype:%s, msg:%v`, eventtype, message)
			c.String(http.StatusOK, lib_define.SUCCESS)
			return
		}
	case `text`:
		nsqMsg[`Content`] = cast.ToString(message[`content`])
	case `image`:
		nsqMsg[`oss_url`], nsqMsg[`filename`] = common.GetWecomFileByUrl(cast.ToString(message[`url`]), adminUserId)
	case `mixed`:
		nsqMsg[`MsgType`] = `text`
		nsqMsg[`Content`] = common.ParseMsgItem(message[`msg_item`], adminUserId)
	case `voice`:
		nsqMsg[`MsgType`] = `text`
		nsqMsg[`Content`] = cast.ToString(message[`content`])
	case `file`:
		if fileurl, filename := common.GetWecomFileByUrl(cast.ToString(message[`url`]), adminUserId); len(fileurl) > 0 {
			nsqMsg[`MsgType`] = `text`
			nsqMsg[`Content`] = fmt.Sprintf(`%s(%s)`, filename, fileurl)
		} else {
			nsqMsg[`oss_url`], nsqMsg[`filename`] = ``, ``
		}
	case `stream`:
		//not supported for now
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	default:
		logs.Warning(`wecom_robot unsupported msgtype:%s, msg:%v`, msgtype, message)
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	go common.PushNSQ(nsqMsg)
	c.String(http.StatusOK, lib_define.SUCCESS)
}
