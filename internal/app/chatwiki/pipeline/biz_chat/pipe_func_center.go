// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// CheckKeywordReply 关键词检测处理
func CheckKeywordReply(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.AppInfo) > 0 && len(in.params.ReceivedMessageType) > 0 && in.params.ReceivedMessageType != lib_define.MsgTypeText {
		return pipeline.PipeContinue //微信应用等,非文本消息跳过关键词回复逻辑
	}
	var keywordReplyList []common.ReplyContent
	keywordReplyList, in.keywordSkipAI, _ = common.BuildKeywordReplyMessage(in.params)
	if len(keywordReplyList) > 0 {
		out.replyContentList = append(out.replyContentList, keywordReplyList...)
	}
	return pipeline.PipeContinue
}

// SetRobotPaymentAuthCodeManager 设置授权码管理员
func SetRobotPaymentAuthCodeManager(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText
	ctx := context.Background()
	// 管理员发送设置码
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

			// 从访客表里查询详细信息
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
			replyContent.Description = `已将您设置为【授权码】管理员，后续发送【授权码】将返回可用的授权码`
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

// GetRobotPaymentAuthCodePackage 管理员获取授权码套餐
func GetRobotPaymentAuthCodePackage(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText
	if in.params.Question == `授权码` && in.isPaymentManager {
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
			replyContent.Description = "点击以下智能菜单将为您生成对应的授权码："
			for _, countPackageInfo := range countPackageInfoList {
				content := fmt.Sprintf("%s【%d次】--%s元", countPackageInfo.Name, countPackageInfo.Count, cast.ToString(countPackageInfo.Price))
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
			replyContent.Description = "点击以下智能菜单将为您生成对应的授权码："
			for _, durationPackageInfo := range durationPackageInfoList {
				content := fmt.Sprintf("%s【%d天】--%s元", durationPackageInfo.Name, durationPackageInfo.Duration, cast.ToString(durationPackageInfo.Price))
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

// GetRobotPaymentAuthCodeContent 管理员获取授权码
func GetRobotPaymentAuthCodeContent(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent

	if in.isPaymentManager {
		reg1 := regexp.MustCompile(`^(.+?)【(\d+)次】--(.+?)元$`)
		reg2 := regexp.MustCompile(`^(.+?)【(\d+)天】--(.+?)元$`)
		matches1 := reg1.FindStringSubmatch(in.params.Question)
		matches2 := reg2.FindStringSubmatch(in.params.Question)

		if len(matches1) > 0 || len(matches2) > 0 {
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

			// 次数套餐
			if len(matches1) > 0 {
				var countPackageInfoList []*define.RobotPaymentCountPackage
				err = json.Unmarshal([]byte(setting[`count_package`]), &countPackageInfoList)
				if err != nil {
					logs.Error(err.Error())
					return pipeline.PipeContinue
				}
				for _, countPackageInfo := range countPackageInfoList {
					format := fmt.Sprintf("%s【%d次】--%s元", countPackageInfo.Name, countPackageInfo.Count, cast.ToString(countPackageInfo.Price))
					if format == in.params.Question {
						packageId = countPackageInfo.Id
						packageName = countPackageInfo.Name
						packageCount = countPackageInfo.Count
						packagePrice = countPackageInfo.Price
						content = `###` + fmt.Sprintf(`%dC`, packageCount) + tool.Random(15) + "###"
						break
					}
				}
			} else if len(matches2) > 0 {
				var durationPackageInfoList []define.RobotPaymentDurationPackage
				err = json.Unmarshal([]byte(setting[`duration_package`]), &durationPackageInfoList)
				if err != nil {
					logs.Error(err.Error())
					return pipeline.PipeContinue
				}
				for _, durationPackageInfo := range durationPackageInfoList {
					format := fmt.Sprintf("%s【%d天】--%s元", durationPackageInfo.Name, durationPackageInfo.Duration, cast.ToString(durationPackageInfo.Price))
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

			// 从访客表里查询详细信息
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
			replyContent.Description = fmt.Sprintf("您的授权码为：\n%s\n将本内容发送到公众号兑换", content)
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

// ExchangeRobotPaymentAuthCode 兑换授权码
func ExchangeRobotPaymentAuthCode(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText

	re := regexp.MustCompile(`###.+?###`)
	authCodeContent := re.FindString(in.params.Question)
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
			replyContent.Description = fmt.Sprintf("兑换失败。\n请检查授权码是否正确")
			out.replyContentList = append(out.replyContentList, replyContent)
			if in.params.PassiveId > 0 {
				out.content = replyContent.Description
			}
			in.paymentSkipAIAndWorkflow = true
			return pipeline.PipeContinue
		}

		if in.isPaymentManager {
			// 组装公众号消息
			replyContent.Description += fmt.Sprintf("管理员无需兑换，可直接使用工作流，以下为授权码日志")
			replyContent.Description += fmt.Sprintf("\n授权码：\n%s", authCodeContent)
			if cast.ToInt(info[`package_type`]) == define.RobotPaymentPackageTypeCount {
				replyContent.Description += fmt.Sprintf("\n所属套餐：\n%s【%s次】--%s元", info[`package_name`], info[`package_count`], info[`package_price`])
			} else {
				replyContent.Description += fmt.Sprintf("\n所属套餐：\n%s【%s天】--%s元", info[`package_name`], info[`package_duration`], info[`package_price`])
			}
			if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusPending {
				replyContent.Description += fmt.Sprintf("\n兑换状态：待使用")
			} else if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusExchanged {
				replyContent.Description += fmt.Sprintf("\n兑换状态：已兑换")
				replyContent.Description += fmt.Sprintf("\n兑换人：%s", info[`exchanger_name`])
				replyContent.Description += fmt.Sprintf("\n兑换时间：%s", time.Unix(cast.ToInt64(info[`exchange_time`]), 0).Format("06-01-02 15:04:05"))
			} else if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusUsed {
				replyContent.Description += fmt.Sprintf("\n兑换状态：已使用")
				replyContent.Description += fmt.Sprintf("\n使用时间：%s", time.Unix(cast.ToInt64(info[`use_time`]), 0).Format("06-01-02 15:04:05"))
			}
			replyContent.Description += fmt.Sprintf("\n创建时间：%s", time.Unix(cast.ToInt64(info[`create_time`]), 0).Format("06-01-02 15:04:05"))
			replyContent.Description += fmt.Sprintf("\n创建人：%s", info[`creator_name`])
		} else {
			if cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusExchanged || cast.ToInt(info[`usage_status`]) == define.RobotPaymentAuthCodeUsageStatusUsed {
				replyContent.Description = fmt.Sprintf("兑换失败。\n%s\n%s已兑换", authCodeContent, time.Unix(cast.ToInt64(info[`use_time`]), 0).Format("06-01-02 15:04:05"))
			} else {
				// 从访客表里查询详细信息
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

				// 兑换逻辑
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

				// 组装公众号消息
				replyContent.Description = fmt.Sprintf("您已兑换%s", info[`package_name`])
				if cast.ToInt(info[`package_type`]) == define.RobotPaymentPackageTypeCount {
					replyContent.Description += fmt.Sprintf("\n剩余可用：%d次\n", cast.ToInt(info[`package_count`])-cast.ToInt(info[`used_count`]))
				} else {
					replyContent.Description += fmt.Sprintf("\n剩余可用：%d天", cast.ToInt(info[`package_duration`])-cast.ToInt(info[`used_duration`]))
					replyContent.Description += fmt.Sprintf("\n剩余可用：%d次\n", cast.ToInt(info[`package_count`])-cast.ToInt(info[`used_count`]))
				}
				replyContent.Description += `回复 <a href="weixin://bizmsgmenu?msgmenucontent=我的权益&msgmenuid=1">我的权益</a> 可查看使用情况`
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

// QueryRobotPaymentAuthCodeRight 查看授权码权益
func QueryRobotPaymentAuthCodeRight(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	//应用收费没开启
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}

	var replyContent common.ReplyContent
	replyContent.ReplyType = common.ReplyTypeText

	if in.params.Question == "我的权益" {
		// 先计算时长套餐
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

		replyContent.Description = `当前`
		if remainDuration > 0 {
			replyContent.Description += fmt.Sprintf("\n剩余可用：%d天", remainDuration)
		}
		replyContent.Description += fmt.Sprintf("\n剩余可用：%d次", remainCount)

		if len(authCodeList) > 0 {
			replyContent.Description += fmt.Sprintf("\n兑换记录：")
			for _, authCode := range authCodeList {
				if cast.ToInt(authCode[`package_type`]) == define.RobotPaymentPackageTypeCount {
					replyContent.Description += fmt.Sprintf("\n%s 兑换%d次",
						time.Unix(cast.ToInt64(authCode[`exchange_time`]), 0).Format("06-01-02 15:04:05"),
						cast.ToInt(authCode[`package_count`]),
					)
				} else {
					replyContent.Description += fmt.Sprintf("\n%s 兑换%d天",
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

// CheckReceivedMessageReply 收到消息回复处理
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

// PushReplyContentList 推送回复内容列表
func PushReplyContentList(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(out.replyContentList) > 0 {
		in.Stream(sse.Event{Event: `reply_content_list`, Data: out.replyContentList})
	}
	return pipeline.PipeContinue
}
