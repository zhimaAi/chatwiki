// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"regexp"

	"github.com/cloudwego/eino/schema"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func SendDefaultUnknownQuestionPrompt(in *ChatInParam, out *ChatOutParam, errmsg string) {
	in.Stream(sse.Event{Event: `error`, Data: `SYSERR:` + errmsg})
	code := `unknown`
	if ms := regexp.MustCompile(`ERROR\s+CODE:\s?(.*)`).FindStringSubmatch(errmsg); len(ms) > 1 {
		code = ms[1]
	}
	logs.Error(`robot_id:%s,openid:%s,gpt_error:%s`, in.params.Robot[`id`], in.params.Openid, errmsg)
	out.content = i18n.Show(in.params.Lang, `gpt_error`, code)
	StreamContent(in, out.content)
}

func StreamContent(in *ChatInParam, content string) {
	if cast.ToInt(in.params.Robot[`application_type`]) == define.ApplicationTypeClaw {
		in.Stream(sse.Event{Event: `stream_message`, Data: schema.AssistantMessage(content, nil)})
		return
	}
	in.Stream(sse.Event{Event: `sending`, Data: content})
}

// DoRequestChatUnify unified logic for requesting large language model
func DoRequestChatUnify(in *ChatInParam, out *ChatOutParam) {
	if !in.needRunWorkFlow && in.useStream {
		out.chatResp, out.requestTime, out.Error = common.RequestChatStream(
			in.params.StopCtx,
			in.params.Lang,
			in.params.AdminUserId,
			in.params.Openid,
			in.params.Robot,
			in.params.AppType,
			cast.ToInt(in.params.Robot[`model_config_id`]),
			in.params.Robot[`use_model`],
			out.messages,
			out.functionTools,
			in.chanStream,
			cast.ToFloat32(in.params.Robot[`temperature`]),
			cast.ToInt(in.params.Robot[`max_token`]),
			common.ToReasoningEffort(in.params.Robot[`enable_thinking`], in.params.Robot[`reasoning_effort`]),
		)
	} else {
		out.chatResp, out.requestTime, out.Error = common.RequestChat(
			in.params.StopCtx,
			in.params.Lang,
			in.params.AdminUserId,
			in.params.Openid,
			in.params.Robot,
			in.params.AppType,
			cast.ToInt(in.params.Robot[`model_config_id`]),
			in.params.Robot[`use_model`],
			out.messages,
			out.functionTools,
			cast.ToFloat32(in.params.Robot[`temperature`]),
			cast.ToInt(in.params.Robot[`max_token`]),
			common.ToReasoningEffort(in.params.Robot[`enable_thinking`], in.params.Robot[`reasoning_effort`]),
		)
	}
	out.content = out.chatResp.Result()
	out.reasoningContent = out.chatResp.ReasoningContent()
	if out.Error != nil {
		SendDefaultUnknownQuestionPrompt(in, out, out.Error.Error())
	} else {
		if content, ok := common.ReplaceMiniCardMarkersForRobotPromptReply(in.params.AdminUserId, out.content); ok {
			out.content = content
			out.chatResp.SetResult(out.content)
		}
		if cast.ToInt(in.params.Robot[`chat_type`]) != define.ChatTypeDirect {
			in.saveRobotChatCache = true
		}
	}
}

// DisposeUnknownQuestionPrompt handle unknown question prompt and questions
func DisposeUnknownQuestionPrompt(in *ChatInParam, out *ChatOutParam) {
	unknownQuestionPrompt := define.MenuJsonStruct{}
	_ = tool.JsonDecodeUseNumber(in.params.Robot[`unknown_question_prompt`], &unknownQuestionPrompt)
	if len(unknownQuestionPrompt.Content) == 0 && len(unknownQuestionPrompt.Question) == 0 {
		unknownQuestionPrompt.Content = lib_define.DefaultUnknownQuestionPromptContent // default value
	}
	out.msgType = define.MsgTypeMenu
	out.content = unknownQuestionPrompt.Content
	out.menuJson = tool.JsonEncodeNoError(unknownQuestionPrompt)
	in.saveRobotChatCache = false // unknown questions not saved to chat cache
}
