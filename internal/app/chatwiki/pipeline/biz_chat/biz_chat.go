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
	//请求参数
	params     *define.ChatRequestParam
	useStream  bool
	chanStream chan sse.Event
	//运行参数
	monitor         *common.Monitor
	dialogueId      int
	sessionId       int
	answerMessageId string
	llmStartTime    time.Time //llm请求开始时间
	//标志位
	needRunWorkFlow          bool
	waitChooseWorkFlow       bool
	showQuoteFile            bool
	startQuoteFile           bool
	hitCache                 bool
	keywordSkipAI            bool
	exitChat                 bool //直接退出标志位
	saveRobotChatCache       bool
	isSwitchManual           bool
	chanStreamClosed         bool
	robotAbilityPayment      bool
	isPaymentManager         bool //是否为应用收费管理员
	paymentSkipAIAndWorkflow bool //应用收费忽略AI和工作流
}

// Stream 给前端推送数据流事件
func (in *ChatInParam) Stream(event sse.Event) {
	if in.useStream && in.chanStream != nil {
		in.chanStream <- event
	}
}

type ChatOutParam struct {
	//返回参数
	AiMessage msql.Params
	Error     error
	//运行参数
	cMsgId           int64                                // C端消息id
	cMessage         msql.Params                          // C端消息
	debugLog         []any                                //Prompt日志
	messages         []adaptor.ZhimaChatCompletionMessage //大模型上下文
	functionTools    []adaptor.FunctionTool               //function tool
	list             []msql.Params                        //知识分段召回
	replyContentList []common.ReplyContent                //非ai额外回复内容列表
	quoteFileJson    string                               //引用知识库文件入库json
	chatResp         adaptor.ZhimaChatCompletionResponse  //大模型返回结构体
	reasoningContent string                               //AI思考过程(chatResp.ReasoningContent)
	content          string                               //AI回复内容(chatResp.Result)
	msgType          int                                  //回复消息的类型
	menuJson         string                               //未知问题回复的菜单
	requestTime      int64                                //大模型开始返回答案的时间
	aiMsgId          int64                                //AiMessage.id
}

func DoChatRequest(params *define.ChatRequestParam, useStream bool, chanStream chan sse.Event) (out *ChatOutParam) {
	//builder
	in := &ChatInParam{
		params: params, useStream: useStream, chanStream: chanStream,
		monitor:       common.NewMonitor(params), //聊天监控初始化
		showQuoteFile: cast.ToBool(params.Robot[`answer_source_switch`]),
	}
	out = &ChatOutParam{
		debugLog:         make([]any, 0), //debug log
		messages:         make([]adaptor.ZhimaChatCompletionMessage, 0),
		list:             make([]msql.Params, 0),
		replyContentList: make([]common.ReplyContent, 0),
		quoteFileJson:    `[]`,               //默认值
		msgType:          define.MsgTypeText, //默认值
	}
	//defer
	defer func() {
		if in.monitor != nil {
			in.monitor.DebugLog = out.debugLog
			in.monitor.Save(out.Error) //记录监控数据
		}
		in.chanStreamClosed = true //先标记chan关闭
		if in.chanStream != nil {
			defer close(in.chanStream) //关闭sse管道
		}
	}()

	//request
	request := pipeline.NewPipeline(in, out)
	request.Pipe(CheckChanStream)        //检查流式输出的管道
	request.Pipe(SseKeepAlive)           //流式输出保活
	request.Pipe(StreamPing)             //给前端推送ping
	request.Pipe(CheckParams)            //请求参数检查
	request.Pipe(FilterLibrary)          //过滤知识库
	request.Pipe(CloseOpenApiReceiver)   //close open_api receiver
	request.Pipe(GetDialogueId)          //校验对话或创建对话
	request.Pipe(GetSessionId)           //获取会话ID
	request.Pipe(CustomerPush)           //推送customer信息
	request.Pipe(SaveCustomerMsg)        //保存customer消息
	request.Pipe(UpLastChatByC)          //更新last_chat
	request.Pipe(WebsocketNotifyByC)     //接待变更通知
	request.Pipe(SetRobotAbilityPayment) //设置机器人应用收费开关标志
	request.Pipe(CheckPaymentManager)    //设置当前会话是否是授权码管理员
	request.Process()
	if out.Error != nil {
		return //出错后终止逻辑
	}

	//switch_manual
	switchManual := pipeline.NewPipeline(in, out)
	switchManual.Pipe(CheckManualReplyPauseRobotReply) //人工介入+人工回复后暂停机器人回复
	switchManual.Pipe(CheckKeywordSwitchManual)        //关键词转人工
	switchManual.Pipe(CheckIntentionSwitchManual)      //根据用户意图转人工
	if switchManual.Process() == pipeline.PipeStop {
		return //转人工后终止逻辑
	}

	//init_llm
	initLlm := pipeline.NewPipeline(in, out)
	initLlm.Pipe(RobotInfoPush) //推送机器人信息
	initLlm.Process()           //不存在错误,继续向下执行

	//func_center
	funcCenter := pipeline.NewPipeline(in, out)
	funcCenter.Pipe(CheckKeywordReply)              //关键词检测处理
	funcCenter.Pipe(SetRobotPaymentAuthCodeManager) //设置授权码管理员
	funcCenter.Pipe(GetRobotPaymentAuthCodePackage) //管理员获取授权码套餐
	funcCenter.Pipe(GetRobotPaymentAuthCodeContent) //管理员获取授权码
	funcCenter.Pipe(ExchangeRobotPaymentAuthCode)   //兑换授权码
	funcCenter.Pipe(QueryRobotPaymentAuthCodeRight) //查看授权码权益
	funcCenter.Pipe(CheckReceivedMessageReply)      //收到消息回复处理
	funcCenter.Pipe(PushReplyContentList)           //推送回复内容列表
	funcCenter.Process()                            //忽略错误,继续向下执行

	//recall+context
	recall := pipeline.NewPipeline(in, out)
	recall.Pipe(CheckKeywordSkipAi)            //检查是否跳过AI回复
	recall.Pipe(CheckPaymentSkipAiAndWorkflow) //检查是否跳过AI和工作流回复
	recall.Pipe(CheckReplyByChatCache)         //检查回复来自聊天缓存
	recall.Pipe(CheckWorkFlowRobot)            //检查工作流机器人
	recall.Pipe(CheckChatTypeDirect)           //直连模式逻辑
	recall.Pipe(CheckChatTypeNotDirect)        //混合模式和仅知识库模式
	recall.Process()
	if in.exitChat || out.Error != nil {
		return //出错或转人工终止逻辑
	}

	//call_llm
	callLlm := pipeline.NewPipeline(in, out)
	callLlm.Pipe(CheckSkipCallLlm)              //检查是否跳过llm调用
	callLlm.Pipe(CheckKeywordSkipAi)            //检查是否跳过AI回复
	callLlm.Pipe(CheckPaymentSkipAiAndWorkflow) //检查是否跳过AI和工作流回复
	callLlm.Pipe(BuildOpenApiContent)           //开放接口自定义上下文处理
	callLlm.Pipe(BuildFunctionTools)            //构建function tool
	callLlm.Pipe(SetLlmStartTime)               //设置llm请求开始时间
	callLlm.Pipe(CheckRobotPaymentAuthCode)     //消耗授权码次数
	callLlm.Pipe(DoApplicationTypeFlow)         //工作流机器人逻辑
	callLlm.Pipe(DoChatByChatCache)             //通过聊天缓存获取相应内容
	callLlm.Pipe(DoChatTypeDirect)              //直连模式聊天逻辑
	callLlm.Pipe(DoChatTypeMixture)             //混合模式聊天逻辑
	callLlm.Pipe(DoChatTypeLibrary)             //仅知识库模式聊天逻辑
	callLlm.Process()
	if in.exitChat {
		return //非AI出错后终止逻辑
	}

	//work_flow
	workFlow := pipeline.NewPipeline(in, out)
	workFlow.Pipe(CheckSkipCallLlm)              //检查是否跳过llm调用
	workFlow.Pipe(CheckPaymentSkipAiAndWorkflow) //检查是否跳过AI和工作流回复
	workFlow.Pipe(DoRelationWorkFlow)            //聊天机器人支持关联工作流
	workFlow.Process()                           //忽略错误,继续向下执行

	//ending
	ending := pipeline.NewPipeline(in, out)
	ending.Pipe(CheckSaveRobotMsg)         //检查是否需要保存AI消息
	ending.Pipe(SetMonitorFromLlm)         //记录llm的监控数据
	ending.Pipe(OfficeAccountPassiveReply) //未认证公众号的消息特殊处理
	ending.Pipe(DisposeClientBreak)        //处理客户端断开逻辑
	ending.Pipe(PushDebugLog)              //推送渲染Prompt日志
	ending.Pipe(SaveRobotMsg)              //保存机器人消息
	ending.Pipe(SaveRobotChatCache)        //保存机器人聊天缓存
	ending.Pipe(UpLastChatByAi)            //更新last_chat
	ending.Pipe(WebsocketNotifyByAi)       //接待变更通知
	ending.Pipe(SaveAnswerSource)          //保存答案来源
	ending.Pipe(AdditionAiMessage)         //追加机器人消息返回字段
	ending.Pipe(PushAiMessageFinish)       //推送机器人消息+完成标志
	ending.Process()                       //收尾工作,不用管结果了
	return
}
