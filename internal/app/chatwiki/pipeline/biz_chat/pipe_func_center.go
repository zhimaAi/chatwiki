// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// CheckKeywordReply handles keyword detection and reply
func CheckKeywordReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.AppInfo) > 0 && len(in.params.ReceivedMessageType) > 0 && in.params.ReceivedMessageType != lib_define.MsgTypeText {
		return pipeline.PipeContinue // skip keyword reply logic for non-text messages in WeChat app etc.
	}
	var keywordReplyList []common.ReplyContent
	keywordReplyList, in.keywordSkipAI, _ = common.BuildKeywordReplyMessage(in.params)
	if len(keywordReplyList) > 0 {
		out.replyContentList = append(out.replyContentList, keywordReplyList...)
	}
	return pipeline.PipeContinue
}

// SetRobotPaymentAuthCodeManager sets auth code manager
func SetRobotPaymentAuthCodeManager(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText
	ctx := context.Background()
	// admin sends setup code
	if strings.HasPrefix(in.params.Question, define.RobotPaymentAuthCodePrefix) && strings.HasSuffix(in.params.Question, define.RobotPaymentAuthCodeSuffix) {
		key := define.RobotPaymentAuthCodeManagerCachePrefix + in.params.Robot[`id`]
		val, err := define.Redis.Get(ctx, key).Result()
		if err == nil && len(val) > 0 {
			info, err := msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
				Where(`robot_id`, in.params.Robot[`id`]).
				Where(`manager_openid`, in.params.Openid).
				Find()
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}

			// query details from customer table
			avatar := ``
			name := ``
			customer, err := msql.Model(`chat_ai_customer`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
				Where(`openid`, in.params.Openid).
				Find()
			if err != nil {
				logs.Error(err.Error())
			}
			if len(customer) > 0 {
				avatar = customer[`avatar`]
				name = customer[`name`]
			}

			if len(info) > 0 {
				_, err = msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
					Where(`id`, info[`id`]).
					Update(msql.Datas{
						`manager_avatar`:   avatar,
						`manager_nickname`: name,
						`update_time`:      tool.Time2Int(),
					})
				if err != nil {
					logs.Error(err.Error())
				}
			} else {
				_, err = msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
					Where(`id`, info[`id`]).
					Insert(msql.Datas{
						`admin_user_id`:    in.params.AdminUserId,
						`robot_id`:         in.params.Robot[`id`],
						`manager_openid`:   in.params.Openid,
						`manager_avatar`:   avatar,
						`manager_nickname`: name,
						`create_time`:      tool.Time2Int(),
						`update_time`:      tool.Time2Int(),
					})
				if err != nil {
					logs.Error(err.Error())
				}
			}
			// default to zh-CN for auth code content (no bilingual match)
			replyContent.Description = i18n.Show(define.LangZhCn, "payment_set_manager_success")
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}
	}
	return pipeline.PipeContinue
}

// GetRobotPaymentAuthCodePackage admin gets auth code package
func GetRobotPaymentAuthCodePackage(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText
	// bilingual keyword match: auth code
	matched, matchedLang := i18n.KeywordAuthCode.Match(in.params.Question)
	if matched && in.isPaymentManager {
		// default to zh-CN, use en-US only when en-US match succeeds
		lang := define.LangZhCn
		if matchedLang == define.LangEnUs {
			lang = define.LangEnUs
		}
		setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
			Where(`robot_id`, in.params.Robot[`id`]).
			Find()
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}
		if len(setting) == 0 {
			return pipeline.PipeContinue
		}
		if len(setting[`count_package`]) > 0 && cast.ToInt(setting[`package_type`]) == define.RobotPaymentPackageTypeCount {
			var countPackageInfoList []define.RobotPaymentCountPackage
			err = json.Unmarshal([]byte(setting[`count_package`]), &countPackageInfoList)
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}
			replyContent.Description = i18n.Show(lang, "payment_select_package_prompt")
			for _, countPackageInfo := range countPackageInfoList {
				content := fmt.Sprintf(i18n.GetPaymentCountPackageFormat(lang), countPackageInfo.Name, countPackageInfo.Count, cast.ToString(countPackageInfo.Price))
				replyContent.Description += fmt.Sprintf("\n<a href=\"weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=%d\">%s</a>", content, countPackageInfo.Id, content)
			}
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}
		if len(setting[`duration_package`]) > 0 && cast.ToInt(setting[`package_type`]) == define.RobotPaymentPackageTypeDuration {
			var durationPackageInfoList []define.RobotPaymentDurationPackage
			err = json.Unmarshal([]byte(setting[`duration_package`]), &durationPackageInfoList)
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}
			replyContent.Description = i18n.Show(lang, "payment_select_package_prompt")
			for _, durationPackageInfo := range durationPackageInfoList {
				content := fmt.Sprintf(i18n.GetPaymentDurationPackageFormat(lang), durationPackageInfo.Name, durationPackageInfo.Duration, cast.ToString(durationPackageInfo.Price))
				replyContent.Description += fmt.Sprintf("\n<a href=\"weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=%d\">%s</a>", content, durationPackageInfo.Id, content)
			}
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}
	}
	return pipeline.PipeContinue
}

// GetRobotPaymentAuthCodeContent admin gets auth code content
func GetRobotPaymentAuthCodeContent(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent

	if in.isPaymentManager {
		// use bilingual regex to match package format
		countResult := i18n.PaymentCountPackagePattern.Match(in.params.Question)
		durationResult := i18n.PaymentDurationPackagePattern.Match(in.params.Question)

		if countResult.Matched || durationResult.Matched {
			// default to zh-CN, use en-US only when en-US match succeeds
			lang := define.LangZhCn
			if countResult.Matched && countResult.Lang == define.LangEnUs {
				lang = define.LangEnUs
			} else if durationResult.Matched && durationResult.Lang == define.LangEnUs {
				lang = define.LangEnUs
			}
			setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
				Where(`robot_id`, in.params.Robot[`id`]).
				Find()
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}
			if len(setting) == 0 {
				return pipeline.PipeContinue
			}

			var packageId int
			var packageName string
			var packageDuration int
			var packageCount int
			var packagePrice float32
			var content string

			// count package
			if countResult.Matched {
				var countPackageInfoList []*define.RobotPaymentCountPackage
				err = json.Unmarshal([]byte(setting[`count_package`]), &countPackageInfoList)
				if err != nil {
					logs.Error(err.Error())
					return pipeline.PipeContinue
				}
				for _, countPackageInfo := range countPackageInfoList {
					format := fmt.Sprintf(i18n.GetPaymentCountPackageFormat(lang), countPackageInfo.Name, countPackageInfo.Count, cast.ToString(countPackageInfo.Price))
					if format == in.params.Question {
						packageId = countPackageInfo.Id
						packageName = countPackageInfo.Name
						packageCount = countPackageInfo.Count
						packagePrice = countPackageInfo.Price
						content = `###` + fmt.Sprintf(`%dC`, packageCount) + tool.Random(15) + "###"
						break
					}
				}
			} else if durationResult.Matched {
				var durationPackageInfoList []define.RobotPaymentDurationPackage
				err = json.Unmarshal([]byte(setting[`duration_package`]), &durationPackageInfoList)
				if err != nil {
					logs.Error(err.Error())
					return pipeline.PipeContinue
				}
				for _, durationPackageInfo := range durationPackageInfoList {
					format := fmt.Sprintf(i18n.GetPaymentDurationPackageFormat(lang), durationPackageInfo.Name, durationPackageInfo.Duration, cast.ToString(durationPackageInfo.Price))
					if format == in.params.Question {
						packageId = durationPackageInfo.Id
						packageName = durationPackageInfo.Name
						packageDuration = durationPackageInfo.Duration
						packageCount = durationPackageInfo.Count
						packagePrice = durationPackageInfo.Price
						content = `###` + fmt.Sprintf(`%dD`, packageDuration) + tool.Random(15) + "###"
						break
					}
				}
			}

			if len(packageName) == 0 {
				return pipeline.PipeContinue
			}

			// query details from customer table
			name := ``
			customer, err := msql.Model(`chat_ai_customer`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
				Where(`openid`, in.params.Openid).
				Find()
			if err != nil {
				logs.Error(err.Error())
			}
			if len(customer) > 0 {
				name = customer[`name`]
			}

			item := msql.Datas{
				`admin_user_id`:    in.params.AdminUserId,
				`creator_id`:       0,
				`creator_name`:     name,
				`robot_id`:         in.params.Robot[`id`],
				`content`:          content,
				`package_type`:     setting[`package_type`],
				`package_id`:       packageId,
				`package_name`:     packageName,
				`package_duration`: packageDuration,
				`package_count`:    packageCount,
				`package_price`:    packagePrice,
				`usage_status`:     define.RobotPaymentAuthCodeUsageStatusPending,
				`remark`:           ``,
				`create_date`:      time.Now().Format(`2006-01-02`),
				`exchange_date`:    ``,
				`use_time`:         0,
				`use_date`:         ``,
				`create_time`:      tool.Time2Int(),
				`update_time`:      tool.Time2Int(),
			}
			_, err = msql.Model(`robot_payment_auth_code`, define.Postgres).Insert(item)
			if err != nil {
				logs.Error(err.Error())
				return pipeline.PipeContinue
			}

			replyContent.ReplyType = common.ReplyTypeText
			replyContent.Description = i18n.Show(lang, "payment_auth_code_generated", content)
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}
	}

	return pipeline.PipeContinue
}

// ExchangeRobotPaymentAuthCode exchanges auth code
func ExchangeRobotPaymentAuthCode(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText

	authCodeContent := i18n.PaymentAuthCodeContentPattern.FindString(in.params.Question)
	// if question only contains auth code, use zh-CN; otherwise use user's lang
	lang := in.params.Lang
	if authCodeContent == in.params.Question {
		lang = define.LangZhCn
	}
	if authCodeContent != "" {
		info, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
			Where(`robot_id`, in.params.Robot[`id`]).
			Where(`content`, authCodeContent).
			Find()
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}
		if len(info) == 0 {
			replyContent.Description = i18n.Show(lang, "payment_exchange_failed")
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}

		if in.isPaymentManager {
			// build official account message
			replyContent.Description += i18n.Show(lang, "payment_manager_no_exchange")
			replyContent.Description += fmt.Sprintf("\n%s\n%s", i18n.Show(lang, "payment_auth_code_label"), authCodeContent)
			if cast.ToInt(info[`package_type`]) == define.RobotPaymentPackageTypeCount {
				packageFormat := fmt.Sprintf(i18n.GetPaymentCountPackageFormat(lang), info[`package_name`], cast.ToInt(info[`package_count`]), info[`package_price`])
				replyContent.Description += fmt.Sprintf("\n%s\n%s", i18n.Show(lang, "payment_package_label"), packageFormat)
			} else {
				packageFormat := fmt.Sprintf(i18n.GetPaymentDurationPackageFormat(lang), info[`package_name`], cast.ToInt(info[`package_duration`]), info[`package_price`])
				replyContent.Description += fmt.Sprintf("\n%s\n%s", i18n.Show(lang, "payment_package_label"), packageFormat)
			}
			if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusPending {
				replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_status_pending"))
			} else if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusExchanged {
				replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_status_exchanged"))
				replyContent.Description += fmt.Sprintf("\n%s%s", i18n.Show(lang, "payment_exchanger_label"), info[`exchanger_name`])
				replyContent.Description += fmt.Sprintf("\n%s%s", i18n.Show(lang, "payment_exchange_time_label"), time.Unix(cast.ToInt64(info[`exchange_time`]), 0).Format("06-01-02 15:04:05"))
			} else if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusUsed {
				replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_status_used"))
				replyContent.Description += fmt.Sprintf("\n%s%s", i18n.Show(lang, "payment_use_time_label"), time.Unix(cast.ToInt64(info[`use_time`]), 0).Format("06-01-02 15:04:05"))
			}
			replyContent.Description += fmt.Sprintf("\n%s%s", i18n.Show(lang, "payment_create_time_label"), time.Unix(cast.ToInt64(info[`create_time`]), 0).Format("06-01-02 15:04:05"))
			replyContent.Description += fmt.Sprintf("\n%s%s", i18n.Show(lang, "payment_creator_label"), info[`creator_name`])
		} else {
			if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusExchanged || cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusUsed {
				replyContent.Description = i18n.Show(lang, "payment_exchange_failed_already", authCodeContent, time.Unix(cast.ToInt64(info[`exchange_time`]), 0).Format("06-01-02 15:04:05"))
			} else {
				// query details from customer table
				avatar := ``
				name := ``
				customer, err := msql.Model(`chat_ai_customer`, define.Postgres).
					Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
					Where(`openid`, in.params.Openid).
					Find()
				if err != nil {
					logs.Error(err.Error())
				}
				if len(customer) > 0 {
					avatar = customer[`avatar`]
					name = customer[`name`]
				}

				// exchange logic
				_, err = msql.Model(`robot_payment_auth_code`, define.Postgres).
					Where(`id`, info[`id`]).
					Update(msql.Datas{
						`update_time`:      tool.Time2Int(),
						`exchanger_openid`: in.params.Openid,
						`exchange_time`:    tool.Time2Int(),
						`exchange_date`:    time.Now().Format(`2006-01-02`),
						`exchanger_name`:   name,
						`exchanger_avatar`: avatar,
						`usage_status`:     define.RobotPaymentAuthCodeUsageStatusExchanged,
					})
				if err != nil {
					logs.Error(err.Error())
				}

				// build official account message
				replyContent.Description = i18n.Show(lang, "payment_exchange_success", info[`package_name`])
				if cast.ToInt(info[`package_type`]) == define.RobotPaymentPackageTypeCount {
					replyContent.Description += i18n.Show(lang, "payment_remain_count_with_newline", cast.ToInt(info[`package_count`])-cast.ToInt(info[`used_count`]))
				} else {
					replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_remain_duration", cast.ToInt(info[`package_duration`])-cast.ToInt(info[`used_duration`])))
					replyContent.Description += i18n.Show(lang, "payment_remain_count_with_newline", cast.ToInt(info[`package_count`])-cast.ToInt(info[`used_count`]))
				}
				// select link text based on language
				if lang == define.LangEnUs {
					replyContent.Description += i18n.Show(lang, "payment_check_benefits_link_en")
				} else {
					replyContent.Description += i18n.Show(lang, "payment_check_benefits_link")
				}
			}
		}
		out.replyContentList = append(out.replyContentList, replyContent)
		if in.params.PassiveId > 0 {
			out.content = replyContent.Description
		}
		in.paymentSkipAIAndWorkflow = true
		return pipeline.PipeContinue
	}
	return pipeline.PipeContinue
}

// QueryRobotPaymentAuthCodeRight views auth code benefits
func QueryRobotPaymentAuthCodeRight(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	// payment feature not enabled
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText

	// bilingual keyword match: my benefits
	matched, matchedLang := i18n.KeywordMyBenefits.Match(in.params.Question)
	if matched {
		// default to zh-CN, use en-US only when en-US match succeeds
		lang := define.LangZhCn
		if matchedLang == define.LangEnUs {
			lang = define.LangEnUs
		}
		// calculate duration package first
		authCodeList, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
			Where(`robot_id`, in.params.Robot[`id`]).
			Where(`exchanger_openid`, in.params.Openid).
			Order(`id desc`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			return pipeline.PipeContinue
		}

		remainCount := 0
		remainDuration := 0
		for _, authCode := range authCodeList {
			if cast.ToInt(authCode[`usage_status`]) != define.RobotPaymentAuthCodeUsageStatusExchanged {
				continue
			}
			if cast.ToInt(authCode[`package_type`]) == define.RobotPaymentPackageTypeCount {
				remainCount += cast.ToInt(authCode[`package_count`]) - cast.ToInt(authCode[`used_count`])
			} else {
				remainDuration += cast.ToInt(authCode[`package_duration`]) - cast.ToInt(authCode[`used_duration`])
				remainCount += cast.ToInt(authCode[`package_count`]) - cast.ToInt(authCode[`used_count`])
			}
		}

		replyContent.Description = i18n.Show(lang, "payment_current_status")
		if remainDuration > 0 {
			replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_remain_duration", remainDuration))
		}
		replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_remain_count", remainCount))

		if len(authCodeList) > 0 {
			replyContent.Description += fmt.Sprintf("\n%s", i18n.Show(lang, "payment_exchange_history"))
			for _, authCode := range authCodeList {
				if cast.ToInt(authCode[`package_type`]) == define.RobotPaymentPackageTypeCount {
					replyContent.Description += i18n.Show(lang, "payment_exchange_count_record",
						time.Unix(cast.ToInt64(authCode[`exchange_time`]), 0).Format("06-01-02 15:04:05"),
						cast.ToInt(authCode[`package_count`]),
					)
				} else {
					replyContent.Description += i18n.Show(lang, "payment_exchange_duration_record",
						time.Unix(cast.ToInt64(authCode[`exchange_time`]), 0).Format("06-01-02 15:04:05"),
						cast.ToInt(authCode[`package_duration`]),
					)
				}
			}
		}
		out.replyContentList = append(out.replyContentList, replyContent)
		if in.params.PassiveId > 0 {
			out.content = replyContent.Description
		}
		in.paymentSkipAIAndWorkflow = true
		return pipeline.PipeContinue
	}

	return pipeline.PipeContinue
}

// CheckReceivedMessageReply handles received message reply
func CheckReceivedMessageReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.keywordSkipAI {
		return pipeline.PipeContinue
	}
	receivedMessageReplyList, _ := common.BuildReceivedMessageReply(in.params, lib_define.MsgTypeText)
	if len(receivedMessageReplyList) > 0 {
		out.replyContentList = append(out.replyContentList, receivedMessageReplyList...)
	}
	return pipeline.PipeContinue
}

// PushReplyContentList pushes reply content list
func PushReplyContentList(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(out.replyContentList) > 0 {
		in.Stream(sse.Event{Event: `reply_content_list`, Data: out.replyContentList})
	}
	return pipeline.PipeContinue
}
