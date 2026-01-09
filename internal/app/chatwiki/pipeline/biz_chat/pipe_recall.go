// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/pipeline"
	"fmt"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
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

// CheckPaymentSkipAiAndWorkflow 检查是否跳过AI和工作流
func CheckPaymentSkipAiAndWorkflow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.paymentSkipAIAndWorkflow {
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

// CheckRobotPaymentAuthCode 检查授权码次数
func CheckRobotPaymentAuthCode(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	//获取收费配置
	setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
		Where(`robot_id`, in.params.Robot[`id`]).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return pipeline.PipeContinue
	}
	if len(setting) == 0 {
		logs.Error(fmt.Sprintf("机器人%s没有配置应用收费", in.params.Robot[`id`]))
		return pipeline.PipeContinue
	}

	//管理员无限制使用
	if in.isPaymentManager {
		return pipeline.PipeContinue
	}

	//先检查试用
	if cast.ToInt(setting[`try_count`]) > 0 {
		info, err := msql.Model(`robot_payment_user_try_count`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
			Where(`robot_id`, in.params.Robot[`id`]).
			Where(`openid`, in.params.Openid).
			Find()
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}
		if len(info) == 0 {
			_, err = msql.Model(`robot_payment_user_try_count`, define.Postgres).
				Insert(msql.Datas{
					`admin_user_id`: in.params.AdminUserId,
					`robot_id`:      in.params.Robot[`id`],
					`openid`:        in.params.Openid,
					`create_time`:   tool.Time2Int(),
					`update_time`:   tool.Time2Int(),
					`try_count`:     1,
				})
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}
			return pipeline.PipeContinue
		} else if cast.ToInt(info[`try_count`]) < cast.ToInt(setting[`try_count`]) {
			_, err = msql.Model(`robot_payment_user_try_count`, define.Postgres).Where(`id`, info[`id`]).
				Update(msql.Datas{`try_count`: cast.ToInt(info[`try_count`]) + 1, `update_time`: tool.Time2Int()})
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}
			return pipeline.PipeContinue
		}
	}

	//优先使用时长套餐
	packageInfo, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
		Where(`robot_id`, in.params.Robot[`id`]).
		Where(`exchanger_openid`, in.params.Openid).
		Where(`package_type`, cast.ToString(define.RobotPaymentPackageTypeDuration)).
		Where(`usage_status`, cast.ToString(define.RobotPaymentAuthCodeUsageStatusExchanged)).
		Order(`id asc`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return pipeline.PipeContinue
	}
	if len(packageInfo) == 0 {
		//再考虑次数套餐
		packageInfo, err = msql.Model(`robot_payment_auth_code`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
			Where(`robot_id`, in.params.Robot[`id`]).
			Where(`exchanger_openid`, in.params.Openid).
			Where(`package_type`, cast.ToString(define.RobotPaymentPackageTypeCount)).
			Where(`usage_status`, cast.ToString(define.RobotPaymentAuthCodeUsageStatusExchanged)).
			Order(`id asc`).
			Find()
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}
	}

	if len(packageInfo) > 0 { // 增加使用次数
		usedCount := cast.ToInt(packageInfo[`used_count`]) + 1
		data := msql.Datas{
			`used_count`:  usedCount,
			`update_time`: tool.Time2Int(),
		}
		if usedCount >= cast.ToInt(packageInfo[`package_count`]) {
			data[`usage_status`] = define.RobotPaymentAuthCodeUsageStatusUsed
			data[`use_time`] = tool.Time2Int()
			data[`use_date`] = time.Now().Format(`2006-01-02`)
		}
		_, err = msql.Model(`robot_payment_auth_code`, define.Postgres).
			Where(`id`, packageInfo[`id`]).
			Update(data)
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}
		return pipeline.PipeContinue
	} else { // 没有使用次数
		var replyContent common.ReplyContent
		replyContent.ReplyType = common.ReplyTypeImg
		replyContent.ThumbURL = setting[`package_poster`]
		out.replyContentList = append(out.replyContentList, replyContent)
		if in.params.PassiveId > 0 {
			out.content = fmt.Sprintf("![package_poster](%s)", setting[`package_poster`])
		}
		return pipeline.PipeStop
	}
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
