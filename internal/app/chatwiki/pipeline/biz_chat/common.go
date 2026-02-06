// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

// DoRequestChatUnify unified logic for requesting large language model
func DoRequestChatUnify(in *ChatInParam, out *ChatOutParam) {
	if !in.needRunWorkFlow && in.useStream {
		out.chatResp, out.requestTime, out.Error = common.RequestChatStream(
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
		)
	} else {
		out.chatResp, out.requestTime, out.Error = common.RequestChat(
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
		)
	}
	out.content = out.chatResp.Result
	out.reasoningContent = out.chatResp.ReasoningContent
	if out.Error != nil {
		common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
	} else {
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
