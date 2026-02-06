// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/pipeline"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// CheckChanStream check stream output pipe
func CheckChanStream(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if in.chanStream == nil {
		out.Error = errors.New(`channel stream is nil`)
		return pipeline.PipeStop
	}
	return pipeline.PipeContinue
}

// SseKeepAlive stream output keep alive
func SseKeepAlive(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	go func() {
		defer func() { // prevent: send on closed channel
			if r := recover(); r != nil {
				logs.Error(`Recovered in SseKeepAlive:%v`, r)
			}
		}()
		for {
			time.Sleep(3 * time.Second)
			if in.chanStreamClosed {
				return // closed, no more push
			}
			in.Stream(sse.Event{Event: `keep-alive`, Data: tool.Date()})
		}
	}()
	return pipeline.PipeContinue
}

// StreamPing push ping to frontend
func StreamPing(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	in.Stream(sse.Event{Event: `ping`, Data: tool.Time2Int()})
	return pipeline.PipeContinue
}

// CheckParams check request parameters
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
	// when input is multi-modal, perform basic parameter validation
	if questionMultiple, ok := common.ParseInputQuestion(in.params.Question); ok {
		for idx := range questionMultiple {
			var imageTotal int
			switch questionMultiple[idx].Type {
			case adaptor.TypeText:
				if len(questionMultiple[idx].Text) == 0 {
					out.Error = errors.New(i18n.Show(in.params.Lang, `param_invalid`, fmt.Sprintf(`question.%d.text`, idx)))
				}
			case adaptor.TypeImage:
				imageTotal++
				if len(questionMultiple[idx].ImageUrl.Url) == 0 {
					out.Error = errors.New(i18n.Show(in.params.Lang, `param_invalid`, fmt.Sprintf(`question.%d.image_url.url`, idx)))
				}
			case adaptor.TypeAudio, adaptor.TypeVideo:
				//current not check
			default:
				out.Error = errors.New(i18n.Show(in.params.Lang, `param_invalid`, fmt.Sprintf(`question.%d.type`, idx)))
			}
			if out.Error == nil && imageTotal > 10 {
				out.Error = errors.New(i18n.Show(in.params.Lang, `max_upload_images_limit`))
			}
			if out.Error != nil {
				in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
				return pipeline.PipeStop
			}
		}
		in.params.Question = tool.JsonEncodeNoError(questionMultiple)
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

		// if official account knowledge base is closed, remove official account libraries
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

// GetDialogueId validate or create dialogue
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
	// record dialogue id
	in.dialogueId = in.params.DialogueId
	in.Stream(sse.Event{Event: `dialogue_id`, Data: in.dialogueId})
	return pipeline.PipeContinue
}

// GetSessionId get session id
func GetSessionId(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	sessionId, err := common.GetSessionId(in.params, in.dialogueId)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	// record session id
	in.sessionId = sessionId
	in.Stream(sse.Event{Event: `session_id`, Data: in.sessionId})
	return pipeline.PipeContinue
}

// CustomerPush push customer info
func CustomerPush(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	customer, err := common.GetCustomerInfo(in.params.Openid, in.params.AdminUserId)
	if err != nil {
		out.Error = err
		in.Stream(sse.Event{Event: `error`, Data: out.Error.Error()})
		return pipeline.PipeStop
	}
	in.params.Customer = customer // update
	in.Stream(sse.Event{Event: `customer`, Data: customer})
	return pipeline.PipeContinue
}

// SaveCustomerMsg save customer message
func SaveCustomerMsg(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	msgType, content := define.MsgTypeText, in.params.Question
	// wechat and other apps: special handling for message types
	switch in.params.ReceivedMessageType {
	case lib_define.MsgTypeImage:
		msgType, content = define.MsgTypeImage, in.params.MediaIdToOssUrl
	case lib_define.MsgTypeVoice, lib_define.MsgTypeVideo:
		showContent := lib_define.MsgTypeNameMap[in.params.ReceivedMessageType]
		msgType, content = define.MsgTypeText, i18n.Show(in.params.Lang, `received_message_type`, showContent)
	}
	// when input is multi-modal: change message type for database
	if questionMultiple, ok := common.ParseInputQuestion(content); ok {
		msgType, content = define.MsgTypeMixed, tool.JsonEncodeNoError(questionMultiple)
	}
	cMessage := msql.Datas{
		`admin_user_id`:             in.params.AdminUserId,
		`robot_id`:                  in.params.Robot[`id`],
		`openid`:                    in.params.Openid,
		`dialogue_id`:               in.dialogueId,
		`session_id`:                in.sessionId,
		`is_customer`:               define.MsgFromCustomer,
		`msg_type`:                  msgType,
		`content`:                   content,
		`received_message_type`:     in.params.ReceivedMessageType,
		`received_message`:          tool.JsonEncodeNoError(in.params.ReceivedMessage),
		`media_id_to_oss_url`:       in.params.MediaIdToOssUrl,
		`thumb_media_id_to_oss_url`: in.params.ThumbMediaIdToOssUrl,
		`menu_json`:                 ``,
		`quote_file`:                `[]`,
		`create_time`:               tool.Time2Int(),
		`update_time`:               tool.Time2Int(),
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

// UpLastChatByC update last_chat
func UpLastChatByC(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	content := common.GetFirstQuestionByInput(out.cMessage[`content`]) // multi-modal input special handling
	lastChat := msql.Datas{
		`last_chat_time`:    out.cMessage[`create_time`],
		`last_chat_message`: common.MbSubstr(content, 0, 1000),
		`rel_user_id`:       in.params.RelUserId,
	}
	common.UpLastChat(in.dialogueId, in.sessionId, lastChat, define.MsgFromCustomer)
	return pipeline.PipeContinue
}

func UpChatPromptVariablesByC(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if len(in.params.ChatPromptVariables) == 0 {
		return pipeline.PipeContinue
	}
	upData := msql.Datas{
		`chat_prompt_variables`: in.params.ChatPromptVariables,
		`rel_user_id`:           in.params.RelUserId,
	}
	var fillVariables []map[string]any
	err := tool.JsonDecode(in.params.ChatPromptVariables, &fillVariables)
	if err != nil {
		logs.Error(err.Error())
		return pipeline.PipeContinue
	}
	common.UpChatPromptVariables(in.dialogueId, in.sessionId, upData)
	in.Stream(sse.Event{Event: `chat_prompt_variables`, Data: map[string]any{
		`dialogue_id`:        in.dialogueId,
		`session_id`:         in.sessionId,
		`need_fill_variable`: false,
		`fill_variables`:     fillVariables,
		`wait_variables`:     nil,
	}})
	return pipeline.PipeContinue
}

// WebsocketNotifyByC receiver change notification
func WebsocketNotifyByC(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	common.ReceiverChangeNotify(in.params.AdminUserId, `c_message`, out.cMessage)
	return pipeline.PipeContinue
}

// SetRobotAbilityPayment SetRobotAbilityPayment
func SetRobotAbilityPayment(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	robotAbilityConfig := common.GetRobotAbilityConfigByAbilityType(in.params.AdminUserId, cast.ToInt(in.params.Robot[`id`]), common.RobotAbilityPayment)
	if len(robotAbilityConfig) != 0 {
		if len(in.params.AppInfo) > 0 &&
			len(in.params.ReceivedMessageType) > 0 &&
			in.params.ReceivedMessageType == lib_define.MsgTypeText &&
			tool.InArrayString(in.params.AppInfo[`app_type`], []string{lib_define.AppOfficeAccount, lib_define.AppMini, lib_define.AppWechatKefu}) {
			in.robotAbilityPayment = true
		}
	}
	return pipeline.PipeContinue
}

// CheckPaymentManager check if current session is auth code manager
func CheckPaymentManager(in *ChatInParam, out *ChatOutParam) pipeline.PipeResult {
	if !in.robotAbilityPayment {
		return pipeline.PipeContinue
	}
	info, err := msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(in.params.AdminUserId)).
		Where(`robot_id`, in.params.Robot[`id`]).
		Where(`manager_openid`, in.params.Openid).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return pipeline.PipeContinue
	}
	if len(info) == 0 {
		return pipeline.PipeContinue
	}
	in.isPaymentManager = true
	return pipeline.PipeContinue
}
