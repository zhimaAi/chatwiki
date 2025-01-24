// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"encoding/json"
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
		AppType:     lib_define.AppOpenApi,
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

type ChatCompletionMessage struct {
	Role    string `json:"role,omitempty" binding:"required"`
	Content string `json:"content,omitempty"`
}

type ChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages" binding:"required,dive"`
	Stream      bool                    `json:"stream,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Temperature float64                 `json:"temperature,omitempty"`
	RobotKey    string
}

type ChatCompletionResponse struct {
	ID      string              `json:"id,omitempty"`
	Created int                 `json:"created,omitempty" `
	Usage   ChatCompletionUsage `json:"usage,omitempty" `
	Model   string              `json:"model,omitempty"`
	Choices []interface{}       `json:"choices,omitempty"`
	Object  string              `json:"object,omitempty"`
}

type ChatCompletionUsage struct {
	CompletionTokens int `json:"completion_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty" `
	TotalTokens      int `json:"total_tokens,omitempty"`
}

// Completions compatible openai standard api
func Completions(c *gin.Context) {
	var req = ChatCompletionRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtOpenAiErr(c, http.StatusBadRequest, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	// token check
	headers, err := common.ParseAuthorizationToken(c)
	if err != nil {
		common.FmtOpenAiErr(c, http.StatusUnauthorized, err.Error())
		return
	}
	req.RobotKey = cast.ToString(headers["robot_key"])
	params, err := req.buildChatRequestParam(c)
	if err != nil {
		common.FmtOpenAiErr(c, http.StatusBadRequest, `sys_err`)
		return
	}
	msg, _ := tool.JsonEncode(req.Messages)
	if define.IsDev {
		logs.Debug("请求数据原始:%+v", msg)
		logs.Debug("请求数据解析后问题:%+v", params.Question)
		logs.Debug("请求数据解析后提示词:%+v", params.Prompt)
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
		responseId := common.BuildOpenAiMsgId()
		streamResponse(c, responseId, chanStream)
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
			common.FmtOpenAiErr(c, http.StatusBadRequest, err.Error())
			return
		}
		msg, _ := common.GetImgInMessage(message["content"])
		streamResp := openAiRes{
			model:            params.Robot["use_model"],
			content:          msg,
			isFinish:         true,
			promptTokens:     cast.ToInt(message["prompt_tokens"]),
			completionTokens: cast.ToInt(message["completion_tokens"]),
		}
		res := formatStandardOpenAiRes(streamResp)
		// only data response
		common.FmtOpenAiOk(c, res)
	}
}

func (r *ChatCompletionRequest) buildChatRequestParam(c *gin.Context) (*define.ChatRequestParam, error) {
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
	openId := common.BuildOpenId(r.RobotKey)
	customer, err := common.GetCustomerInfo(openId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	if r.MaxTokens > 0 {
		robot["max_token"] = cast.ToString(r.MaxTokens)
	}
	if r.Temperature > 0 {
		robot["temperature"] = cast.ToString(r.Temperature)
	}
	chatBaseParam := &define.ChatBaseParam{
		AppType:     lib_define.AppOpenApi,
		Openid:      openId,
		AdminUserId: adminUserId,
		Robot:       robot,
		Customer:    customer,
	}
	if len(customer) == 0 {
		go saveCustomerInfo(c, chatBaseParam)
	}
	dialogueId, err := msql.Model("chat_ai_dialogue", define.Postgres).Where("robot_id", robot["id"]).
		Where("openid", openId).Order("id desc").Limit(1).Value("id")
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf(`sys_err`)
	}
	isClose := false
	question := ""
	openApiContent := ""
	prompt := ""
	if len(r.Messages) > 0 {
		msgArr := make([]ChatCompletionMessage, 0)
		for key, item := range r.Messages {
			if item.Role == "user" {
				question = item.Content
			}
			if key+1 == len(r.Messages) {
				continue
			}
			msgArr = append(msgArr, item)
		}
		openApiContent, _ = tool.JsonEncode(msgArr)
	}
	return &define.ChatRequestParam{
		ChatBaseParam:  chatBaseParam,
		Lang:           common.GetLang(c),
		Question:       strings.TrimSpace(question),
		OpenApiContent: strings.TrimSpace(openApiContent),
		DialogueId:     cast.ToInt(dialogueId),
		Prompt:         strings.TrimSpace(prompt),
		LibraryIds:     strings.TrimSpace(robot["library_ids"]),
		IsClose:        &isClose,
	}, nil
}

func streamResponse(c *gin.Context, responseId string, chanStream chan sse.Event) {
	c.Stream(func(w io.Writer) bool {
		if event, ok := <-chanStream; ok {
			var resp interface{}
			flusher, _ := w.(http.Flusher)
			switch event.Event {
			case "sending":
				content, ers := event.Data.(string)
				if !ers {
					return false
				}
				streamResp := openAiRes{
					id:       responseId,
					content:  content,
					isStream: true,
				}
				resp = formatStandardOpenAiRes(streamResp)
				bts, _ := json.Marshal(resp)
				_, err := fmt.Fprintf(w, "data: %s\n\n", string(bts))
				if err != nil {
					logs.Error(err.Error())
					return false
				}
			case "data":
				content, ers := event.Data.(msql.Datas)
				if !ers {
					return false
				}
				streamResp := openAiRes{
					id:               responseId,
					model:            cast.ToString(content["use_model"]),
					completionTokens: cast.ToInt(content["completion_tokens"]),
					promptTokens:     cast.ToInt(content["prompt_tokens"]),
					isStream:         true,
					isFinish:         true,
				}
				resp = formatStandardOpenAiRes(streamResp)
				bts, _ := json.Marshal(resp)
				_, err := fmt.Fprintf(w, "data: %s\n\n", string(bts))
				if err != nil {
					logs.Error(err.Error())
					return false
				}
			case "finish":
				resp = "[DONE]"
				_, err := fmt.Fprintf(w, "data: %s\n\n", resp)
				if err != nil {
					logs.Error(err.Error())
					return false
				}
			}
			flusher.Flush()
			return true
		}
		return false
	})
}

type openAiRes struct {
	id               string
	content          string
	model            string
	promptTokens     int
	completionTokens int
	isFinish         bool
	isStream         bool
}

func formatStandardOpenAiRes(response openAiRes) interface{} {
	choices := map[string]interface{}{
		"finish_reason": nil,
		"index":         0,
	}
	message := ChatCompletionMessage{
		Content: response.content,
	}
	object := "chat.completion"
	if response.id == "" {
		response.id = common.BuildOpenAiMsgId()
	}
	if response.isStream {
		choices["delta"] = message
		object = "chat.completion.chunk"
	} else {
		choices["message"] = message
	}
	if response.isFinish {
		choices["finish_reason"] = "stop"
		if response.isStream {
			choices["delta"] = ChatCompletionMessage{}
		}
	}

	resp := ChatCompletionResponse{
		ID:      response.id,
		Created: tool.Time2Int(),
		Usage: ChatCompletionUsage{
			CompletionTokens: response.completionTokens,
			PromptTokens:     response.promptTokens,
			TotalTokens:      response.completionTokens + response.promptTokens,
		},
		Model:   response.model,
		Choices: []interface{}{choices},
		Object:  object,
	}
	return resp
}
