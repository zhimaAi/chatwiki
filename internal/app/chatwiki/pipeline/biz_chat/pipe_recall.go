package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/pipeline"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// CheckKeywordSkipAi 检查是否跳过AI回复
func CheckKeywordSkipAi(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.keywordSkipAI {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckWorkFlowRobot 检查工作流机器人
func CheckWorkFlowRobot(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckChatTypeDirect 直连模式逻辑
func CheckChatTypeDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeDirect {
		out.messages, out.Error = common.BuildDirectChatRequestMessage(in.params, out.cMsgId, in.dialogueId, &out.debugLog)
		if out.Error != nil {
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckChatTypeNotDirect 混合模式和仅知识库模式
func CheckChatTypeNotDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	defer func() { //显示引文逻辑
		recall := pipeline.NewPipeline(in, out)
		recall.Pipe(DisposeQuoteFilePush) //处理引用知识库文件+推送
		recall.Process()
	}()
	//真正的召回pipeline
	recall := pipeline.NewPipeline(in, out)
	recall.Pipe(CheckRepetitionSwitchManual) //重复问题转人工
	recall.Pipe(StartQuoteFile)              //显示引文开始信号
	recall.Pipe(RecallByChatCache)           //通过聊天缓存召回知识库
	recall.Pipe(RecallByLibrarys)            //通过知识库召回分段
	recall.Process()
	if in.exitChat || out.Error != nil {
		return pipeline.PipeStop //出错或转人工终止逻辑
	}
	//处理未知问题统计
	isBackground := len(in.params.Customer) > 0 && cast.ToInt(in.params.Customer[`is_background`]) > 0
	if !isBackground && len(out.list) == 0 {
		common.SaveUnknownIssueRecord(in.params.AdminUserId, in.params.Robot, in.params.Question)
	}
	//未知问题转人工
	if msg, ok := IsUnknownSwitchManual(in.params, in.sessionId, in.dialogueId, out.list, in.chanStream); ok {
		out.AiMessage = msg
		in.exitChat = true //直接退出标志位
		return pipeline.PipeStop
	}
	return pipeline.PipeStop
}

// CheckRepetitionSwitchManual 重复问题转人工
func CheckRepetitionSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if msg, ok := IsRepetitionSwitchManual(in.params, in.sessionId, in.dialogueId, out.cMsgId, in.chanStream); ok {
		out.AiMessage = msg
		in.exitChat = true //直接退出标志位
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// StartQuoteFile 显示引文开始信号
func StartQuoteFile(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.showQuoteFile && !in.startQuoteFile {
		in.Stream(sse.Event{Event: `start_quote_file`, Data: tool.Time2Int()})
		in.startQuoteFile = true
	}
	return pipeline.PipeContinue
}

// RecallByChatCache 通过聊天缓存召回知识库
func RecallByChatCache(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.hitCache, in.answerMessageId = common.HitRobotMessageCache(in.params.Robot[`robot_key`], in.params.Question, in.params.Robot[`cache_config`])
	if in.hitCache {
		out.list, in.hitCache, out.Error = common.BuildLibraryMessagesFromCache(in.params.Robot[`robot_key`], in.answerMessageId)
		if out.Error != nil {
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
	}
	if in.hitCache {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// RecallByLibrarys 通过知识库召回分段
func RecallByLibrarys(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	out.messages, out.list, in.monitor.LibUseTime, out.Error = common.BuildLibraryChatRequestMessage(in.params, out.cMsgId, in.dialogueId, &out.debugLog)
	in.Stream(sse.Event{Event: `recall_time`, Data: in.monitor.LibUseTime.RecallTime})
	if out.Error != nil {
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DisposeQuoteFilePush 处理引用知识库文件+推送
func DisposeQuoteFilePush(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	var quoteFile []msql.Params
	quoteFile, out.quoteFileJson = common.DisposeQuoteFile(in.params.AdminUserId, out.list)
	if in.showQuoteFile && in.startQuoteFile {
		in.Stream(sse.Event{Event: `quote_file`, Data: quoteFile})
	}
	return pipeline.PipeContinue
}
