// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/pipeline"
	"errors"
	"strings"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// CheckChanStream 检查流式输出的管道
func CheckChanStream(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.chanStream == nil {
		out.Error = errors.New(`channel stream is nil`)
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// StreamPing 给前端推送ping
func StreamPing(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.Stream(sse.Event{Event: `ping`, Data: tool.Time2Int()})
	return pipeline.PipeContinue
}

// CheckParams 请求参数检查
func CheckParams(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.params.Error != nil {
		out.Error = in.params.Error
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	if len(in.params.Question) == 0 {
		out.Error = errors.New(i18n.Show(in.params.Lang, `question_empty`))
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

func FilterLibrary(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.LibraryIds) > 0 {
		libraryIdList := strings.Split(in.params.LibraryIds, ",")
		hasEnabledOfficialAccount, err := common.CheckHasEnabledOfficialAccount(cast.ToInt(in.params.AdminUserId))
		if err != nil {
			out.Error = in.params.Error
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
		filteredLibIdList := libraryIdList

		// 如果关闭了公众号知识库，需要移除公众号类库
		if !hasEnabledOfficialAccount {
			libIds, err := msql.Model(`chat_ai_library`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
				Where(`type`, cast.ToString(define.OfficialLibraryType)).
				ColumnArr(`id`)
			if err != nil {
				out.Error = in.params.Error
				in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
				return pipeline.PipeStop
			}
			filteredLibIdList = []string{}
			for _, libId := range libraryIdList {
				if !tool.InArrayString(libId, libIds) {
					filteredLibIdList = append(filteredLibIdList, libId)
				}
			}
		}
		in.params.LibraryIds = strings.Join(filteredLibIdList, ",")
	}
	return pipeline.PipeContinue
}

// CloseOpenApiReceiver close open_api receiver
func CloseOpenApiReceiver(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	CloseReceiverFromAppOpenApi(in.params)
	return pipeline.PipeContinue
}

// GetDialogueId 校验对话或创建对话
func GetDialogueId(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.params.DialogueId > 0 {
		dialogue, err := common.GetDialogueInfo(in.params.DialogueId, in.params.AdminUserId, cast.ToInt(in.params.Robot[`id`]), in.params.Openid)
		if err != nil {
			out.Error = err
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
		if len(dialogue) == 0 {
			out.Error = errors.New(i18n.Show(in.params.Lang, `param_invalid`, `dialogue_id`))
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
	} else {
		dialogueId, err := common.GetDialogueId(in.params.ChatBaseParam, in.params.Question)
		if err != nil {
			out.Error = err
			in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
			return pipeline.PipeStop
		}
		in.params.DialogueId = dialogueId
	}
	//记录对话ID
	in.dialogueId = in.params.DialogueId
	in.Stream(sse.Event{Event: `dialogue_id`, Data: in.dialogueId})
	return pipeline.PipeContinue
}

// GetSessionId 获取会话ID
func GetSessionId(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	sessionId, err := common.GetSessionId(in.params, in.dialogueId)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	//记录会话ID
	in.sessionId = sessionId
	in.Stream(sse.Event{Event: `session_id`, Data: in.sessionId})
	return pipeline.PipeContinue
}

// CustomerPush 推送customer信息
func CustomerPush(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	customer, err := common.GetCustomerInfo(in.params.Openid, in.params.AdminUserId)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	in.params.Customer = customer //更新一下
	in.Stream(sse.Event{Event: `customer`, Data: customer})
	return pipeline.PipeContinue
}

// SaveCustomerMsg 保存customer消息
func SaveCustomerMsg(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	cMessage := msql.Datas{
		`admin_user_id`: in.params.AdminUserId,
		`robot_id`:      in.params.Robot[`id`],
		`openid`:        in.params.Openid,
		`dialogue_id`:   in.dialogueId,
		`session_id`:    in.sessionId,
		`is_customer`:   define.MsgFromCustomer,
		`msg_type`:      define.MsgTypeText,
		`content`:       in.params.Question,
		`menu_json`:     ``,
		`quote_file`:    `[]`,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	if len(in.params.Customer) > 0 {
		cMessage[`nickname`] = in.params.Customer[`nickname`]
		cMessage[`name`] = in.params.Customer[`name`]
		cMessage[`avatar`] = in.params.Customer[`avatar`]
	}
	if in.params.RelUserId > 0 {
		common.FillRelUserInfo(cMessage, in.params.RelUserId)
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(cMessage, `id`)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	out.cMsgId = id
	out.cMessage = common.ToStringMap(cMessage, `id`, id)
	in.Stream(sse.Event{Event: `c_message`, Data: out.cMessage})
	return pipeline.PipeContinue
}

// UpLastChatByC 更新last_chat
func UpLastChatByC(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	lastChat := msql.Datas{
		`last_chat_time`:    out.cMessage[`create_time`],
		`last_chat_message`: common.MbSubstr(cast.ToString(out.cMessage[`content`]), 0, 1000),
		`rel_user_id`:       in.params.RelUserId,
	}
	common.UpLastChat(in.dialogueId, in.sessionId, lastChat, define.MsgFromCustomer)
	return pipeline.PipeContinue
}

// WebsocketNotifyByC 接待变更通知
func WebsocketNotifyByC(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	common.ReceiverChangeNotify(in.params.AdminUserId, `c_message`, out.cMessage)
	return pipeline.PipeContinue
}
