// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type (
	ChatMessagesReq struct {
		Content  any                                  `form:"content" json:"content" binding:"required"`
		Messages []adaptor.ZhimaChatCompletionMessage `form:"messages" json:"messages"`
		OpenID   string                               `form:"open_id" json:"open_id" binding:"required"`
		Stream   bool                                 `form:"stream" json:"stream,omitempty"`
		QuoteLib bool                                 `form:"quote_lib" json:"quote_lib,omitempty"`
		Global   map[string]any                       `form:"global" json:"global"`
		RobotKey string
	}
	ChatMessagesRes struct {
		MessageId      string               `json:"message_id"`
		ConversationId string               `json:"conversation_id"`
		CreateAt       int64                `json:"create_at"`
		RawAnswer      string               `json:"raw_answer,omitempty"`
		Answer         string               `json:"answer"`
		Image          []string             `json:"image,omitempty"`
		Voice          []string             `json:"voice,omitempty"`
		MetaData       ChatMessagesMetaData `json:"metadata,omitempty"`
		QuoteLib       any                  `json:"quote_lib,omitempty"`
		QuoteFile      any                  `json:"quote_file,omitempty"`
		IsSwitchManual *bool                `json:"is_switch_manual,omitempty"`
		//功能中心-自动回复:(关键词回复+收到消息回复)
		ReplyContentList []common.ReplyContent `json:"reply_content_list,omitempty"`
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
				if data, ok := event.Data.(string); ok {
					event.Data = strings.ReplaceAll(data, "\r", ``)
				}
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
			RawAnswer:      message["content"],
			MetaData: ChatMessagesMetaData{Usage: Usage{
				PromptTokens:     cast.ToInt(message["prompt_tokens"]),
				CompletionTokens: cast.ToInt(message["completion_tokens"]),
			}},
		}
		msg, imgs, voices := common.GetMessageInMessage(res.Answer, false)
		if len(imgs) > 0 || len(voices) > 0 {
			res.Image = imgs
			res.Voice = voices
			res.Answer = msg
		}
		if params.QuoteLib {
			_ = tool.JsonDecodeUseNumber(message[`quote_lib`], &res.QuoteLib)
			_ = tool.JsonDecodeUseNumber(message[`quote_file`], &res.QuoteFile)
		}
		if isSwitchManual := cast.ToBool(message[`is_switch_manual`]); isSwitchManual {
			res.IsSwitchManual = &isSwitchManual
		}
		if len(message[`reply_content_list`]) > 0 { //返回功能中心自动回复内容
			_ = tool.JsonDecodeUseNumber(message[`reply_content_list`], &res.ReplyContentList)
			for i, item := range res.ReplyContentList {
				if len(item.ThumbURL) > 0 && !common.IsUrl(item.ThumbURL) {
					res.ReplyContentList[i].ThumbURL = define.Config.WebService[`image_domain`] + item.ThumbURL
				}
				if len(item.Pic) > 0 && !common.IsUrl(item.Pic) {
					res.ReplyContentList[i].Pic = define.Config.WebService[`image_domain`] + item.Pic
				}
			}
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
	if len(c.GetString(`wechatapp_appid`)) > 0 {
		if !common.IsXkfOpenid(r.OpenID) {
			return nil, fmt.Errorf(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))
		}
	} else {
		if !common.IsChatOpenid(r.OpenID) {
			return nil, fmt.Errorf(i18n.Show(common.GetLang(c), `param_invalid`, `openid`))
		}
	}
	customer, err := common.GetCustomerInfo(r.OpenID, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	chatBaseParam := &define.ChatBaseParam{
		AppType:     lib_define.AppOpenApi,
		Openid:      r.OpenID,
		AdminUserId: adminUserId,
		Robot:       robot,
		Customer:    customer,
	}
	//更新CustomerInfo
	saveCustomerInfo(c, chatBaseParam)
	chatBaseParam.Customer, err = common.GetCustomerInfo(r.OpenID, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	isClose := false
	return &define.ChatRequestParam{
		ChatBaseParam:  chatBaseParam,
		Lang:           common.GetLang(c),
		Question:       common.GetQuestionByContent(r.Content),
		OpenApiContent: tool.JsonEncodeNoError(r.Messages),
		WechatappAppid: c.GetString(`wechatapp_appid`),
		IsClose:        &isClose,
		WorkFlowGlobal: r.Global,
		QuoteLib:       r.QuoteLib,
	}, nil
}

// Completions compatible openai standard api
func Completions(c *gin.Context) {
	c.String(http.StatusNotFound, `开源版本不支持!`)
}

func GetRobotInfo(c *gin.Context) {
	headers, err := common.ParseAuthorizationToken(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if !common.CheckRobotKey(headers[`robot_key`]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	robot, err := common.GetRobotInfo(headers[`robot_key`])
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(robot, nil))
}

func OpenWorkFlowWebHook(c *gin.Context) {
	res, err := work_flow.StartWebhook(c)
	if err != nil {
		logs.Error(`workflow webhook error:%v`, err)
		c.String(http.StatusOK, ``)
	}
	c.String(http.StatusOK, res)
}
