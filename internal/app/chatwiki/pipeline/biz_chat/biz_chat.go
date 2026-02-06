// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/pipeline"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type ChatInParam struct {
	// request parameters
	params     *define.ChatRequestParam
	useStream  bool
	chanStream chan sse.Event
	// runtime parameters
	monitor         *common.Monitor
	dialogueId      int
	sessionId       int
	answerMessageId string
	llmStartTime    time.Time // llm request start time
	// flags
	needRunWorkFlow          bool
	waitChooseWorkFlow       bool
	showQuoteFile            bool
	startQuoteFile           bool
	hitCache                 bool
	keywordSkipAI            bool
	exitChat                 bool // direct exit flag
	saveRobotChatCache       bool
	isSwitchManual           bool
	chanStreamClosed         bool
	robotAbilityPayment      bool
	isPaymentManager         bool // is payment manager
	paymentSkipAIAndWorkflow bool // payment skip ai and workflow
}

// Stream push stream events to frontend
func (in *ChatInParam) Stream(event sse.Event) {
	if in.useStream && in.chanStream != nil {
		in.chanStream <- event
	}
}

type ChatOutParam struct {
	// return parameters
	AiMessage msql.Params
	Error     error
	// runtime parameters
	cMsgId           int64                                // client message id
	cMessage         msql.Params                          // client message
	debugLog         []any                                // prompt log
	messages         []adaptor.ZhimaChatCompletionMessage // llm context
	functionTools    []adaptor.FunctionTool               // function tool
	list             []msql.Params                        // knowledge segment recall
	replyContentList []common.ReplyContent                // non ai extra reply content list
	quoteFileJson    string                               // quote knowledge base file storage json
	chatResp         adaptor.ZhimaChatCompletionResponse  // llm response structure
	reasoningContent string                               // ai reasoning process
	content          string                               // ai reply content
	msgType          int                                  // reply message type
	menuJson         string                               // unknown question reply menu
	requestTime      int64                                // llm start return answer time
	aiMsgId          int64                                // AiMessage.id
}

func DoChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event) (out *ChatOutParam) {
	// builder
	in := &ChatInParam{
		params: params, useStream: useStream, chanStream: chanStream,
		monitor:       common.NewMonitor(params), // chat monitor initialization
		showQuoteFile: cast.ToBool(params.Robot[`answer_source_switch`]),
	}
	out = &ChatOutParam{
		debugLog:         make([]any, 0), // debug log
		messages:         make([]adaptor.ZhimaChatCompletionMessage, 0),
		list:             make([]msql.Params, 0),
		replyContentList: make([]common.ReplyContent, 0),
		quoteFileJson:    `[]`,               // default value
		msgType:          define.MsgTypeText, // default value
	}
	// defer
	defer func() {
		if in.monitor != nil {
			in.monitor.DebugLog = out.debugLog
			in.monitor.Save(out.Error) // save monitor data
		}
		in.chanStreamClosed = true // mark chan closed first
		if in.chanStream != nil {
			defer close(in.chanStream) // close sse pipe
		}
	}()

	// request
	request := pipeline.NewPipeline(in, out)
	request.Pipe(CheckChanStream)          // check stream output pipe
	request.Pipe(SseKeepAlive)             // stream output keep alive
	request.Pipe(StreamPing)               // push ping to frontend
	request.Pipe(CheckParams)              // check request parameters
	request.Pipe(FilterLibrary)            // filter knowledge base
	request.Pipe(CloseOpenApiReceiver)     // close open_api receiver
	request.Pipe(GetDialogueId)            // verify dialogue or create dialogue
	request.Pipe(GetSessionId)             // get session id
	request.Pipe(CustomerPush)             // push customer info
	request.Pipe(UpChatPromptVariablesByC) // save session prompt variables
	request.Pipe(SaveCustomerMsg)          // save customer message
	request.Pipe(UpLastChatByC)            // update last_chat
	request.Pipe(WebsocketNotifyByC)       // reception change notification
	request.Pipe(SetRobotAbilityPayment)   // set robot payment switch flag
	request.Pipe(CheckPaymentManager)      // set current session is auth code manager
	request.Process()
	if out.Error != nil {
		return // terminate logic after error
	}

	// switch to manual
	switchManual := pipeline.NewPipeline(in, out)
	switchManual.Pipe(CheckManualReplyPauseRobotReply) // manual intervention and pause robot reply
	switchManual.Pipe(CheckKeywordSwitchManual)        // keyword switch to manual
	switchManual.Pipe(CheckIntentionSwitchManual)      // switch to manual by user intention
	if switchManual.Process() == pipeline.PipeStop {
		return // terminate logic after switch to manual
	}

	// initialize llm
	initLlm := pipeline.NewPipeline(in, out)
	initLlm.Pipe(RobotInfoPush) // push robot info
	initLlm.Process()           // no error, continue execution

	// function center
	funcCenter := pipeline.NewPipeline(in, out)
	funcCenter.Pipe(CheckKeywordReply)              // keyword detection processing
	funcCenter.Pipe(SetRobotPaymentAuthCodeManager) // set auth code manager
	funcCenter.Pipe(GetRobotPaymentAuthCodePackage) // manager get auth code package
	funcCenter.Pipe(GetRobotPaymentAuthCodeContent) // manager get auth code
	funcCenter.Pipe(ExchangeRobotPaymentAuthCode)   // exchange auth code
	funcCenter.Pipe(QueryRobotPaymentAuthCodeRight) // view auth code benefits
	funcCenter.Pipe(CheckReceivedMessageReply)      // received message reply processing
	funcCenter.Pipe(PushReplyContentList)           // push reply content list
	funcCenter.Process()                            // ignore error, continue execution

	// recall+context
	recall := pipeline.NewPipeline(in, out)
	recall.Pipe(CheckKeywordSkipAi)            // check skip ai reply
	recall.Pipe(CheckPaymentSkipAiAndWorkflow) // check skip ai and workflow reply
	recall.Pipe(CheckReplyByChatCache)         // check reply from chat cache
	recall.Pipe(CheckWorkFlowRobot)            // check workflow robot
	recall.Pipe(CheckChatTypeDirect)           // direct connect mode logic
	recall.Pipe(CheckChatTypeNotDirect)        // mixture and knowledge base only mode
	recall.Process()
	if in.exitChat || out.Error != nil {
		return // terminate logic on error or switch to manual
	}

	// call llm
	callLlm := pipeline.NewPipeline(in, out)
	callLlm.Pipe(CheckSkipCallLlm)              // check skip llm call
	callLlm.Pipe(CheckKeywordSkipAi)            // check skip ai reply
	callLlm.Pipe(CheckPaymentSkipAiAndWorkflow) // check skip ai and workflow reply
	callLlm.Pipe(BuildOpenApiContent)           // open api custom context processing
	callLlm.Pipe(BuildFunctionTools)            // build function tool
	callLlm.Pipe(SetLlmStartTime)               // set llm request start time
	callLlm.Pipe(CheckRobotPaymentAuthCode)     // consume auth code count
	callLlm.Pipe(DoApplicationTypeFlow)         // workflow robot logic
	callLlm.Pipe(DoChatByChatCache)             // get response by chat cache
	callLlm.Pipe(DoChatTypeDirect)              // direct connect mode chat logic
	callLlm.Pipe(DoChatTypeMixture)             // mixture mode chat logic
	callLlm.Pipe(DoChatTypeLibrary)             // knowledge base only mode chat logic
	callLlm.Process()
	if in.exitChat {
		return // terminate logic on non ai error
	}

	// workflow
	workFlow := pipeline.NewPipeline(in, out)
	workFlow.Pipe(CheckSkipCallLlm)              // check skip llm call
	workFlow.Pipe(CheckPaymentSkipAiAndWorkflow) // check skip ai and workflow reply
	workFlow.Pipe(DoRelationWorkFlow)            // chat robot support related workflow
	workFlow.Process()                           // ignore error, continue execution

	// ending
	ending := pipeline.NewPipeline(in, out)
	ending.Pipe(CheckSaveRobotMsg)         // check need save ai message
	ending.Pipe(SetMonitorFromLlm)         // record llm monitor data
	ending.Pipe(OfficeAccountPassiveReply) // unauthenticated office account message special processing
	ending.Pipe(DisposeClientBreak)        // handle client disconnect logic
	ending.Pipe(PushDebugLog)              // push render prompt log
	ending.Pipe(SaveRobotMsg)              // save robot message
	ending.Pipe(SaveRobotChatCache)        // save robot chat cache
	ending.Pipe(UpLastChatByAi)            // update last_chat
	ending.Pipe(WebsocketNotifyByAi)       // reception change notification
	ending.Pipe(SaveAnswerSource)          // save answer source
	ending.Pipe(AdditionAiMessage)         // append robot message return fields
	ending.Pipe(PushAiMessageFinish)       // push robot message and finish flag
	ending.Process()                       // cleanup work, ignore result
	return
}
