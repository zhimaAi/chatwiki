// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_define"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type (
	ChatMessagesReq struct {
		Content  string `form:"content" json:"content" binding:"required"`
		OpenID   string `form:"open_id" json:"open_id" binding:"required"`
		Stream   bool   `form:"stream" json:"stream,omitempty"`
		RobotKey string
	}
	ChatMessagesRes struct {
		MessageId      string               `json:"message_id"`
		ConversationId string               `json:"conversation_id"`
		CreateAt       int64                `json:"create_at"`
		Answer         string               `json:"answer"`
		Image          []string             `json:"image,omitempty"`
		MetaData       ChatMessagesMetaData `json:"metadata,omitempty"`
	}
	ChatMessagesMetaData struct {
		Usage Usage `json:"usage,omitempty"`
	}
	Usage struct {
		PromptTokens     int `json:"prompt_tokens,omitempty"`
		CompletionTokens int `json:"completion_tokens,omitempty"`
	}
)

func ChatMessages(c *gin.Context) {
	var req = ChatMessagesReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	// token check
	headers, err := common.ParseAuthorizationToken(c)
	if err != nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, err.Error())
		return
	}
	req.RobotKey = cast.ToString(headers["robot_key"])
	params, err := req.buildChatRequestParam(c)
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	chanStream := make(chan sse.Event)
	if req.Stream {
		c.Header(`Content-Type`, `text/event-stream`)
		c.Header(`Cache-Control`, `no-cache`)
		c.Header(`Connection`, `keep-alive`)
		if define.IsDev {
			c.Header(`Access-Control-Allow-Origin`, `*`)
		}
		go func() {
			_, _ = DoChatRequest(params, req.Stream, chanStream)
		}()
		c.Stream(func(_ io.Writer) bool {
			if event, ok := <-chanStream; ok {
				c.SSEvent(event.Event, event.Data)
				return true
			}
			return false
		})
	} else {
		go func(chanStream chan sse.Event) {
			for event := range chanStream {
				if define.IsDev {
					event.Data, _ = tool.JsonEncode(event.Data)
					logs.Debug(`event:%v`, event)
				}
			}
		}(chanStream)
		params.Robot["show_type"] = cast.ToString(define.RobotTextResponse)
		message, err := DoChatRequest(params, req.Stream, chanStream)
		if err != nil {
			logs.Error("%s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		res := &ChatMessagesRes{
			MessageId:      common.BuildMessageId("messageId", cast.ToString(message["id"]), cast.ToInt(message["create_time"])),
			ConversationId: common.BuildMessageId("dialogueId", cast.ToString(message["dialogue_id"]), cast.ToInt(message["create_time"])),
			CreateAt:       cast.ToInt64(message["create_time"]),
			Answer:         message["content"],
			MetaData: ChatMessagesMetaData{Usage: Usage{
				PromptTokens:     cast.ToInt(message["prompt_tokens"]),
				CompletionTokens: cast.ToInt(message["completion_tokens"]),
			}},
		}
		msg, imgs := common.GetImgInMessage(res.Answer)
		if len(imgs) > 0 {
			res.Image = imgs
			res.Answer = msg
		}
		common.FmtOk(c, res)
	}
}

func (r *ChatMessagesReq) buildChatRequestParam(c *gin.Context) (*define.ChatRequestParam, error) {
	//  openId parse
	if !common.CheckRobotKey(r.RobotKey) {
		return nil, fmt.Errorf(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))
	}
	robot, err := common.GetRobotInfo(r.RobotKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	if len(robot) == 0 {
		common.FmtError(c, `no_data`)
		return nil, fmt.Errorf(`no_data`)
	}
	adminUserId := cast.ToInt(robot[`admin_user_id`])
	if !common.IsChatOpenid(r.OpenID) {
		return nil, fmt.Errorf(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))
	}
	customer, err := common.GetCustomerInfo(r.OpenID, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	chatBaseParam := &define.ChatBaseParam{
		AppType:     lib_define.AppYunH5,
		Openid:      r.OpenID,
		AdminUserId: adminUserId,
		Robot:       robot,
		Customer:    customer,
	}
	if len(customer) == 0 {
		go saveCustomerInfo(c, chatBaseParam)
	}
	dialogueId, err := msql.Model("chat_ai_dialogue", define.Postgres).Where("robot_id", robot["id"]).
		Where("openid", r.OpenID).Order("id desc").Limit(1).Value("id")
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	isClose := false
	return &define.ChatRequestParam{
		ChatBaseParam: chatBaseParam,
		Error:         nil,
		Lang:          common.GetLang(c),
		Question:      strings.TrimSpace(r.Content),
		DialogueId:    cast.ToInt(dialogueId),
		Prompt:        strings.TrimSpace(robot["prompt"]),
		LibraryIds:    strings.TrimSpace(robot["library_ids"]),
		IsClose:       &isClose,
	}, nil
}
