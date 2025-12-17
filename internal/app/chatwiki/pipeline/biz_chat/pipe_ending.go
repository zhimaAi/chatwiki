// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
	"errors"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// SetMonitorFromLlm 记录llm的监控数据
func SetMonitorFromLlm(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if !in.llmStartTime.IsZero() { //未设置llmStartTime的跳过
		in.monitor.LlmCallTime = time.Now().Sub(in.llmStartTime).Milliseconds()
	}
	in.monitor.RequestTime, in.monitor.Error = out.requestTime, out.Error
	out.Error = nil //清空llm的调用错误
	return pipeline.PipeContinue
}

// OfficeAccountPassiveReply 未认证公众号的消息特殊处理
func OfficeAccountPassiveReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.params.AppType == lib_define.AppOfficeAccount && in.params.PassiveId > 0 {
		common.PassiveReplyLogNotify(in.params.PassiveId, in.params.Question, out.content)
	}
	return pipeline.PipeContinue
}

// DisposeClientBreak 处理客户端断开逻辑
func DisposeClientBreak(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if *in.params.IsClose {
		out.Error = errors.New(`client break`)
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// PushDebugLog 推送渲染Prompt日志
func PushDebugLog(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	out.debugLog = append(out.debugLog, map[string]string{`type`: `cur_answer`, `content`: out.content})
	in.Stream(sse.Event{Event: `debug`, Data: out.debugLog})
	return pipeline.PipeContinue
}

// SaveRobotMsg 保存机器人消息
func SaveRobotMsg(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	aiMessage := msql.Datas{
		`admin_user_id`:          in.params.AdminUserId,
		`robot_id`:               in.params.Robot[`id`],
		`openid`:                 in.params.Openid,
		`dialogue_id`:            in.dialogueId,
		`session_id`:             in.sessionId,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           out.requestTime,
		`recall_time`:            in.monitor.LibUseTime.RecallTime,
		`msg_type`:               out.msgType,
		`content`:                out.content,
		`reasoning_content`:      out.reasoningContent,
		`is_valid_function_call`: out.chatResp.IsValidFunctionCall,
		`menu_json`:              out.menuJson,
		`quote_file`:             out.quoteFileJson,
		`create_time`:            tool.Time2Int(),
		`update_time`:            tool.Time2Int(),
	}
	if len(in.params.Robot) > 0 {
		aiMessage[`nickname`] = `` //none
		aiMessage[`name`] = in.params.Robot[`robot_name`]
		aiMessage[`avatar`] = in.params.Robot[`robot_avatar`]
	}
	if len(out.replyContentList) > 0 { //关键词回复 触发内容
		aiMessage[`reply_content_list`] = tool.JsonEncodeNoError(out.replyContentList)
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(aiMessage, `id`)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	out.aiMsgId = id
	out.AiMessage = common.ToStringMap(aiMessage, `id`, id)
	in.Stream(sse.Event{Event: `ai_message`, Data: out.AiMessage})
	return pipeline.PipeContinue
}

// SaveRobotChatCache 保存机器人聊天缓存
func SaveRobotChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.saveRobotChatCache {
		_ = common.SetRobotMessageCache(in.params.Robot[`robot_key`], in.params.Question, cast.ToString(out.aiMsgId), in.params.Robot[`cache_config`])
	}
	return pipeline.PipeContinue
}

// UpLastChatByAi 更新last_chat
func UpLastChatByAi(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	lastChat := msql.Datas{
		`last_chat_time`:    out.AiMessage[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(out.AiMessage[`content`]), 0, 1000),
	}
	common.UpLastChat(in.dialogueId, in.sessionId, lastChat, define.MsgFromRobot)
	return pipeline.PipeContinue
}

// WebsocketNotifyByAi 接待变更通知
func WebsocketNotifyByAi(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	common.ReceiverChangeNotify(in.params.AdminUserId, `ai_message`, out.AiMessage)
	return pipeline.PipeContinue
}

// SaveAnswerSource 保存答案来源
func SaveAnswerSource(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(out.list) > 0 {
		asm := msql.Model(`chat_ai_answer_source`, define.Postgres)
		for _, one := range out.list {
			_, err := asm.Insert(msql.Datas{
				`admin_user_id`: in.params.AdminUserId,
				`message_id`:    out.aiMsgId,
				`file_id`:       one[`file_id`],
				`paragraph_id`:  one[`id`],
				`word_total`:    one[`word_total`],
				`similarity`:    one[`similarity`],
				`title`:         one[`title`],
				`type`:          one[`type`],
				`content`:       one[`content`],
				`question`:      one[`question`],
				`answer`:        one[`answer`],
				`images`:        one[`images`],
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
			if err != nil {
				logs.Error(`sql:%s,err:%s`, asm.GetLastSql(), err.Error())
			}
		}
	}
	return pipeline.PipeContinue
}

// AdditionAiMessage 追加机器人消息返回字段
func AdditionAiMessage(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	aiMessage := msql.Datas{
		`prompt_tokens`:     out.chatResp.PromptToken,
		`completion_tokens`: out.chatResp.CompletionToken,
		`use_model`:         in.params.Robot[`use_model`],
		`is_switch_manual`:  in.isSwitchManual,
	}
	for key, val := range out.AiMessage {
		aiMessage[key] = val
	}
	AdditionQuoteLib(in.params, out.list, &aiMessage) //quote_lib
	out.AiMessage = common.ToStringMap(aiMessage)     //重新赋值
	return pipeline.PipeContinue
}

// PushAiMessageFinish 推送机器人消息+完成标志
func PushAiMessageFinish(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.Stream(sse.Event{Event: `data`, Data: out.AiMessage})
	in.Stream(sse.Event{Event: `finish`, Data: tool.Time2Int()})
	return pipeline.PipeContinue
}
