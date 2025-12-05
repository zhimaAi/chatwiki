// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/pipeline"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// BuildOpenApiContent 开放接口自定义上下文处理
func BuildOpenApiContent(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	out.messages = common.BuildOpenApiContent(in.params, out.messages)
	return pipeline.PipeContinue
}

// BuildFunctionTools 构建function tool
func BuildFunctionTools(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.Robot[`form_ids`]) > 0 {
		formIdList := strings.Split(in.params.Robot[`form_ids`], `,`)
		out.functionTools, out.Error = common.BuildFunctionTools(formIdList, in.params.AdminUserId)
		if out.Error != nil {
			in.exitChat = true //直接退出标志位
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
	}
	//聊天机器人支持关联工作流
	var workFlowFuncCall []adaptor.FunctionTool
	workFlowFuncCall, in.needRunWorkFlow = work_flow.BuildFunctionTools(in.params.Robot)
	if in.needRunWorkFlow {
		out.functionTools = append(out.functionTools, workFlowFuncCall...)
	}
	return pipeline.PipeContinue
}

// SetLlmStartTime 设置llm请求开始时间
func SetLlmStartTime(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.llmStartTime = time.Now()
	return pipeline.PipeContinue
}

// DoApplicationTypeFlow 工作流机器人逻辑
func DoApplicationTypeFlow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		workFlowParams := &work_flow.WorkFlowParams{ChatRequestParam: in.params, CurMsgId: int(out.cMsgId), DialogueId: in.dialogueId, SessionId: in.sessionId}
		out.content, out.requestTime, in.monitor.LibUseTime, out.list, out.Error = work_flow.CallWorkFlow(workFlowParams, &out.debugLog, in.monitor, &in.isSwitchManual)
		if out.Error != nil {
			in.exitChat = true //直接退出标志位
			out.debugLog = append(out.debugLog, map[string]string{`type`: `cur_question`, `content`: in.params.Question})
			in.Stream(sse.Event{Event: `debug`, Data: out.debugLog}) //渲染Prompt日志
			common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
		} else {
			//显示引文逻辑
			callLlm := pipeline.NewPipeline(in, out)
			callLlm.Pipe(StartQuoteFile)       //显示引文开始信号
			callLlm.Pipe(DisposeQuoteFilePush) //处理引用知识库文件+推送
			callLlm.Process()
			in.Stream(sse.Event{Event: `recall_time`, Data: in.monitor.LibUseTime.RecallTime})
			in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatByChatCache 通过聊天缓存获取相应内容
func DoChatByChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.hitCache {
		out.chatResp, out.requestTime, out.Error = common.ResponseMessagesFromCache(in.answerMessageId, in.useStream, in.chanStream)
		out.content = out.chatResp.Result
		in.Stream(sse.Event{Event: `notice`, Data: `内容来源:聊天缓存`})
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeDirect 直连模式聊天逻辑
func DoChatTypeDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeDirect {
		DoRequestChatUnify(in, out) //请求大语言模型
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeMixture 混合模式聊天逻辑
func DoChatTypeMixture(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeMixture {
		if content, ok := common.CheckQaDirectReply(out.list, in.params.Robot); ok {
			out.content = content
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		} else {
			DoRequestChatUnify(in, out) //请求大语言模型
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoChatTypeLibrary 仅知识库模式聊天逻辑
func DoChatTypeLibrary(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeLibrary {
		if !in.needRunWorkFlow && len(out.list) == 0 {
			DisposeUnknownQuestionPrompt(in, out) //未知问题(未关联工作流场景)
		} else {
			if len(out.list) == 0 {
				in.waitChooseWorkFlow = true //等待选择的关联工作流
			}
			if content, ok := common.CheckQaDirectReply(out.list, in.params.Robot); ok {
				out.content = content
				in.Stream(sse.Event{Event: `sending`, Data: out.content})
			} else {
				DoRequestChatUnify(in, out) //请求大语言模型
			}
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckReplyByChatCache 检查回复来自聊天缓存
func CheckReplyByChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.hitCache {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DoRelationWorkFlow 聊天机器人支持关联工作流
func DoRelationWorkFlow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if out.Error == nil && in.needRunWorkFlow {
		workFlowRobot, workFlowGlobal := work_flow.ChooseWorkFlowRobot(out.chatResp.FunctionToolCalls)
		if len(workFlowRobot) == 0 { //大模型没有返回需要调用的工作流
			if in.waitChooseWorkFlow {
				DisposeUnknownQuestionPrompt(in, out) //未知问题(已关联工作流场景)
				return pipeline.PipeStop
			}
			in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
			in.Stream(sse.Event{Event: `sending`, Data: out.content})
		} else { //组装工作流请求参数,并执行工作流
			workFlowParams := work_flow.BuildWorkFlowParams(*in.params, workFlowRobot, workFlowGlobal, int(out.cMsgId), in.dialogueId, in.sessionId)
			out.content, out.requestTime, _, _, out.Error = work_flow.CallWorkFlow(workFlowParams, &out.debugLog, in.monitor, &in.isSwitchManual)
			if out.Error == nil {
				in.Stream(sse.Event{Event: `request_time`, Data: out.requestTime})
				in.Stream(sse.Event{Event: `sending`, Data: out.content})
			} else {
				common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
			}
			in.saveRobotChatCache = false //关联工作流不进聊天缓存
		}
	}
	return pipeline.PipeStop
}
