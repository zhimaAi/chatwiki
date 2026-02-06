// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// CheckSkipCallLlm check if skip llm call
func CheckSkipCallLlm(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.AppInfo) > 0 && len(in.params.ReceivedMessageType) > 0 && in.params.ReceivedMessageType != lib_define.MsgTypeText {
		switch in.params.ReceivedMessageType {
		case lib_define.MsgTypeText:
			return pipeline.PipeContinue
		case lib_define.MsgTypeImage:
			if cast.ToBool(in.params.Robot[`question_multiple_switch`]) {
				return pipeline.PipeContinue
			}
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// BuildOpenApiContent build open api custom context
func BuildOpenApiContent(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	out.messages = common.BuildOpenApiContent(in.params, out.messages)
	return pipeline.PipeContinue
}

// BuildFunctionTools build function tools
func BuildFunctionTools(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.Robot[`form_ids`]) > 0 {
		formIdList := strings.Split(in.params.Robot[`form_ids`], `,`)
		out.functionTools, out.Error = common.BuildFunctionTools(formIdList, in.params.AdminUserId)
		if out.Error != nil {
			in.exitChat = true // exit flag
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
	}
	// chat bot supports related workflow
	var workFlowFuncCall []adaptor.FunctionTool
	workFlowFuncCall, in.needRunWorkFlow = work_flow.BuildFunctionTools(in.params.Lang, in.params.Robot)
	if in.needRunWorkFlow {
		out.functionTools = append(out.functionTools, workFlowFuncCall...)
	}
	return pipeline.PipeContinue
}

// SetLlmStartTime set llm request start time
func SetLlmStartTime(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.llmStartTime = time.Now()
	return pipeline.PipeContinue
}

// BuildImmediatelyReplyHandle build immediate reply handle
func BuildImmediatelyReplyHandle(in *ChatInParam, out *ChatOutParam) func(replyContent common.ReplyContent) {
	return func(replyContent common.ReplyContent) {
		out.replyContentList = append(out.replyContentList, replyContent)
		in.Stream(sse.Event{Event: `reply_content_list`, Data: out.replyContentList})
	}
}

// DoApplicationTypeFlow workflow bot logic
func DoApplicationTypeFlow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		workFlowParams := &work_flow.WorkFlowParams{ChatRequestParam: in.params, CurMsgId: int(out.cMsgId), DialogueId: in.dialogueId, SessionId: in.sessionId}
		workFlowParams.ImmediatelyReplyHandle = BuildImmediatelyReplyHandle(in, out)
		replyContentList := []common.ReplyContent{}
		out.content, out.requestTime, in.monitor.LibUseTime, out.list, replyContentList, out.Error = work_flow.CallWorkFlow(workFlowParams, &out.debugLog, in.monitor, &in.isSwitchManual)
		if len(out.replyContentList) == 0 {
			out.replyContentList = replyContentList
		} else {
			out.replyContentList = append(out.replyContentList, replyContentList...)
		}
		if out.Error != nil {
			in.exitChat = true // exit flag
			out.debugLog = append(out.debugLog, map[string]string{`type`: `cur_question`, `content`: in.params.Question})
			in.Stream(sse.Event{Event: `debug`, Data: out.debugLog}) // render prompt log
			common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
		} else {
			// show citation logic
			callLlm := pipeline.NewPipeline(in, out)
			callLlm.Pipe(StartQuoteFile)       // show citation start signal
			callLlm.Pipe(DisposeQuoteFilePush) // process citation files and push
			callLlm.Process()
			in.Stream(sse.Event{Event: `recall_time`, Data: in.monitor.LibUseTime.RecallTime})
			in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatByChatCache get content from chat cache
func DoChatByChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.hitCache {
		out.chatResp, out.requestTime, out.Error = common.ResponseMessagesFromCache(in.answerMessageId, in.useStream, in.chanStream)
		out.content = out.chatResp.Result
		in.Stream(sse.Event{Event: `notice`, Data: `content source: chat cache`})
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeDirect direct connection mode chat logic
func DoChatTypeDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeDirect {
		DoRequestChatUnify(in, out) // request llm
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeMixture mixture mode chat logic
func DoChatTypeMixture(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeMixture {
		if content, ok := common.CheckQaDirectReply(out.list, in.params.Robot); ok {
			out.content = content
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		} else {
			DoRequestChatUnify(in, out) // request llm
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeLibrary knowledge base only mode chat logic
func DoChatTypeLibrary(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeLibrary {
		if !in.needRunWorkFlow && len(out.list) == 0 {
			DisposeUnknownQuestionPrompt(in, out) // unknown question (no workflow scene)
		} else {
			if len(out.list) == 0 {
				in.waitChooseWorkFlow = true // wait for workflow selection
			}
			if content, ok := common.CheckQaDirectReply(out.list, in.params.Robot); ok {
				out.content = content
				in.Stream(sse.Event{Event: `sending`, Data: out.content})
			} else {
				DoRequestChatUnify(in, out) // request llm
			}
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckReplyByChatCache check reply from chat cache
func CheckReplyByChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.hitCache {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoRelationWorkFlow chat bot supports related workflow
func DoRelationWorkFlow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if out.Error == nil && in.needRunWorkFlow {
		workFlowRobot, workFlowGlobal := work_flow.ChooseWorkFlowRobot(cast.ToString(in.params.AdminUserId), out.chatResp.FunctionToolCalls)
		if len(workFlowRobot) == 0 { // no workflow returned by llm
			if in.waitChooseWorkFlow {
				DisposeUnknownQuestionPrompt(in, out) // unknown question (workflow scene)
				return pipeline.PipeStop
			}
			in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		} else { // build workflow params and execute workflow
			workFlowParams := work_flow.BuildWorkFlowParams(*in.params, workFlowRobot, workFlowGlobal, int(out.cMsgId), in.dialogueId, in.sessionId)
			workFlowParams.ImmediatelyReplyHandle = BuildImmediatelyReplyHandle(in, out)
			replyContentList := []common.ReplyContent{}
			out.content, out.requestTime, _, _, replyContentList, out.Error = work_flow.CallWorkFlow(workFlowParams, &out.debugLog, in.monitor, &in.isSwitchManual)
			if len(out.replyContentList) == 0 {
				out.replyContentList = replyContentList
			} else {
				out.replyContentList = append(out.replyContentList, replyContentList...)
			}
			if out.Error == nil {
				in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
				in.Stream(sse.Event{Event: `sending`, Data: out.content})
			} else {
				common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
			}
			in.saveRobotChatCache = false // related workflow not saved to chat cache
		}
	}
	return pipeline.PipeStop
}
