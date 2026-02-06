// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/pipeline"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// CheckKeywordSkipAi check if skip ai response
func CheckKeywordSkipAi(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.keywordSkipAI {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckPaymentSkipAiAndWorkflow check if skip ai and workflow
func CheckPaymentSkipAiAndWorkflow(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.paymentSkipAIAndWorkflow {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckWorkFlowRobot check workflow robot
func CheckWorkFlowRobot(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`application_type`]) == define.ApplicationTypeFlow {
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckRobotPaymentAuthCode checks auth code usage count
func CheckRobotPaymentAuthCode(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	// get payment settings
	setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
		Where(`robot_id`, in.params.Robot[`id`]).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return pipeline.PipeContinue
	}
	if len(setting) == 0 {
		logs.Error(i18n.Show(in.params.Lang, "payment_robot_no_config", in.params.Robot[`id`]))
		return pipeline.PipeContinue
	}

	// admin has unlimited usage
	if in.isPaymentManager {
		return pipeline.PipeContinue
	}

	// check trial first
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

	// prioritize duration package
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
		// then consider count package
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

	if len(packageInfo) > 0 { // increment usage count
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
	} else { // no usage count left
		var replyContent common.ReplyContent
		replyContent.ReplyType = common.ReplyTypeImg
		replyContent.ThumbURL = setting[`package_poster`]
		out.replyContentList = append(out.replyContentList, replyContent)
		if in.params.PassiveId > 0 {
			out.content = i18n.Show(in.params.Lang, "payment_package_poster", setting[`package_poster`])
		}
		return pipeline.PipeStop
	}
}

// CheckChatTypeDirect handles direct mode logic
func CheckChatTypeDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if cast.ToInt(in.params.Robot[`chat_type`]) == define.ChatTypeDirect {
		out.messages, out.Error = common.BuildDirectChatRequestMessage(in.params, out.cMsgId, in.dialogueId, in.sessionId, &out.debugLog)
		if out.Error != nil {
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		}
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// CheckChatTypeNotDirect handles hybrid and knowledge base only mode
func CheckChatTypeNotDirect(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	defer func() { // citation display logic
		recall := pipeline.NewPipeline(in, out)
		recall.Pipe(DisposeQuoteFilePush) // process quote file and push
		recall.Process()
	}()
	// actual recall pipeline
	recall := pipeline.NewPipeline(in, out)
	recall.Pipe(CheckRepetitionSwitchManual) // transfer to human for repeated questions
	recall.Pipe(StartQuoteFile)              // show citation start signal
	recall.Pipe(RecallByChatCache)           // recall from knowledge base via chat cache
	recall.Pipe(RecallByLibrarys)            // recall segments from knowledge base
	recall.Process()
	if in.exitChat || out.Error != nil {
		return pipeline.PipeStop // stop logic for error or transfer to human
	}
	// process unknown issue statistics
	isBackground := len(in.params.Customer) > 0 && cast.ToInt(in.params.Customer[`is_background`]) > 0
	if !isBackground && len(out.list) == 0 {
		common.SaveUnknownIssueRecord(in.params.Lang, in.params.AdminUserId, in.params.Robot, in.params.Question)
	}
	// transfer to human for unknown issues
	if msg, ok := IsUnknownSwitchManual(in.params, in.sessionId, in.dialogueId, out.list, in.chanStream); ok {
		out.AiMessage = msg
		in.exitChat = true // direct exit flag
		return pipeline.PipeStop
	}
	return pipeline.PipeStop
}

// CheckRepetitionSwitchManual handles repeated question transfer to human
func CheckRepetitionSwitchManual(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if msg, ok := IsRepetitionSwitchManual(in.params, in.sessionId, in.dialogueId, out.cMsgId, in.chanStream); ok {
		out.AiMessage = msg
		in.exitChat = true // direct exit flag
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// StartQuoteFile shows citation start signal
func StartQuoteFile(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.showQuoteFile && !in.startQuoteFile {
		in.Stream(sse.Event{Event: `start_quote_file`, Data: tool.Time2Int()})
		in.startQuoteFile = true
	}
	return pipeline.PipeContinue
}

// RecallByChatCache recalls from knowledge base via chat cache
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

// RecallByLibrarys recalls segments from knowledge base
func RecallByLibrarys(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	out.messages, out.list, in.monitor.LibUseTime, out.Error = common.BuildLibraryChatRequestMessage(in.params, out.cMsgId, in.dialogueId, in.sessionId, &out.debugLog)
	in.Stream(sse.Event{Event: `recall_time`, Data: in.monitor.LibUseTime.RecallTime})
	if out.Error != nil {
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// DisposeQuoteFilePush processes quote file and push
func DisposeQuoteFilePush(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	var quoteFile []msql.Params
	quoteFile, out.quoteFileJson = common.DisposeQuoteFile(in.params.AdminUserId, out.list)
	if in.showQuoteFile && in.startQuoteFile {
		in.Stream(sse.Event{Event: `quote_file`, Data: quoteFile})
	}
	return pipeline.PipeContinue
}
