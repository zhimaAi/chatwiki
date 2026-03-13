// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_define"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// useNewDialogueFlex accepts use_new_dialogue as number 0/1 or bool true/false in JSON and normalizes to 0/1
type useNewDialogueFlex int

func (v *useNewDialogueFlex) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*v = 1
		} else {
			*v = 0
		}
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*v = useNewDialogueFlex(i)
	return nil
}

// ChatClawChatReq chat request (streaming only)
type ChatClawChatReq struct {
	RobotKey            string                               `json:"robot_key" binding:"required"`
	Content             string                               `json:"content" binding:"required"`
	Messages            []adaptor.ZhimaChatCompletionMessage `json:"messages"`
	DialogueId          int                                  `json:"dialogue_id"`
	UseNewDialogue      useNewDialogueFlex                   `json:"use_new_dialogue"` // 1 or true = new session, 0 or false = continue existing dialogue
	QuoteLib            bool                                 `json:"quote_lib"`
	Global              map[string]any                       `json:"global"`
	ChatPromptVariables string                               `json:"chat_prompt_variables"`
}

// ChatClawChat chat API: requires token auth, robot_key from request body, AppType=chat_claw, streaming only
func ChatClawChat(c *gin.Context) {
	_ = c.Request.ParseMultipartForm(define.DefaultMultipartMemory)
	claims, _, err := common.GetChatClawAuthClaims(c)
	if err != nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userId := cast.ToInt(claims["user_id"])
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	var req ChatClawChatReq
	if err = c.ShouldBindJSON(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	robotKey := strings.TrimSpace(req.RobotKey)
	if !common.CheckRobotKey(robotKey) {
		common.FmtError(c, `param_invalid`, `robot_key`)
		return
	}
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(robot) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	adminUserId := cast.ToInt(robot[`admin_user_id`])

	openid := fmt.Sprintf("chat_claw_client_%d", userId)
	customer, err := common.GetCustomerInfo(openid, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if req.ChatPromptVariables != "" {
		if c.Request.PostForm == nil {
			c.Request.PostForm = make(url.Values)
		}
		c.Request.PostForm.Set("chat_prompt_variables", req.ChatPromptVariables)
	}
	chatPromptVariables := GetPromptVariables(cast.ToString(robot[`id`]), c)
	isClose := false
	if req.Global == nil {
		req.Global = make(map[string]any)
	}
	useNewDialogue := 0
	if int(req.UseNewDialogue) != 0 {
		useNewDialogue = 1
	}
	chatBaseParam := &define.ChatBaseParam{
		AppType:        lib_define.ChatClawClient,
		Openid:         openid,
		AdminUserId:    adminUserId,
		Robot:          robot,
		Customer:       customer,
		UseNewDialogue: useNewDialogue,
	}
	dialogueId := req.DialogueId
	if useNewDialogue == 1 {
		// Create new dialogue at entry to avoid being affected by old dialogue_id in the pipeline
		dialogueId, err = common.GetDialogueId(chatBaseParam, strings.TrimSpace(req.Content))
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}
	params := &define.ChatRequestParam{
		ChatBaseParam:       chatBaseParam,
		Lang:                common.GetLang(c),
		Question:            strings.TrimSpace(req.Content),
		DialogueId:          dialogueId,
		IsClose:             &isClose,
		WorkFlowGlobal:      req.Global,
		QuoteLib:            req.QuoteLib,
		OpenApiContent:      tool.JsonEncodeNoError(req.Messages),
		ChatPromptVariables: chatPromptVariables,
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	if define.IsDev {
		c.Header("Access-Control-Allow-Origin", "*")
	}
	chanStream := make(chan sse.Event)
	go func() {
		_, _ = DoChatRequest(params, true, chanStream)
	}()
	c.Stream(func(w io.Writer) bool {
		event, ok := <-chanStream
		if !ok {
			return false
		}
		data := event.Data
		if s, ok := data.(string); ok {
			data = strings.ReplaceAll(s, "\r", "")
		}
		c.SSEvent(event.Event, data)
		return true
	})
	*params.IsClose = true
	for range chanStream {
	}
}
