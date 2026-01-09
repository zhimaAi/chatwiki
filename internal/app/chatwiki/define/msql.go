// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "chatwiki/internal/pkg/lib_define"

const Postgres = lib_define.Postgres

const (
	DefaultRoleRoot  = "所有者"
	DefaultRoleAdmin = "管理员"
	DefaultRoleUser  = "成员"
)

const (
	TableRole              = "role"
	TableUser              = "public.user"
	TableMenu              = "menu"
	TableCompany           = "company"
	TableFastCommand       = "fast_command"
	TableChatAiRobotApikey = "chat_ai_robot_apikey"
	TableRule              = "casbin_rule"
	TableChatAiWechatApp   = "chat_ai_wechat_app"
	TableUseGuideProcess   = "use_guide_process"
)

const (
	RoleTypeRoot  = 1
	RoleTypeAdmin = 2
	RoleTypeUser  = 3
)

const (
	Normal  = "0"
	Deleted = "1"
)

const (
	BaseOpen  = "1"
	BaseClose = "0"
)

const (
	RobotManage      = `RobotManage`
	LibraryManage    = `LibraryManage`
	SystemManage     = `SystemManage`
	ClientSideManage = `ClientSideManage`
)

const (
	GeneralLibraryType  = 0
	QALibraryType       = 2
	OpenLibraryType     = 1
	OfficialLibraryType = 3
)

var LibraryTypes = [...]int{
	GeneralLibraryType,  //0普通知识库
	QALibraryType,       //2问答知识库
	OpenLibraryType,     //1对外知识库
	OfficialLibraryType, //3公众号知识库
}

const (
	OpenLibraryAccessRights = 1
	AiSummary               = 1
)

const (
	ExportSourceSession    uint = 1 //会话记录导出
	ExportSourceLibFileDoc uint = 2 //知识库文档
)

var ExportSourceList = []map[string]any{
	{`source`: ExportSourceSession, `source_name`: `会话记录导出`},
	{`source`: ExportSourceLibFileDoc, `source_name`: `知识库文档导出`},
}

const (
	ExportStatusWaiting = 0
	ExportStatusRunning = 1
	ExportStatusSucceed = 2
	ExportStatusError   = 3
)

const (
	ApplicationTypeChat = 0
	ApplicationTypeFlow = 1
)

const (
	DataTypeDraft   uint = 1
	DataTypeRelease uint = 2
)

const (
	PromptTypeCustom = 0
	PromptTypeStruct = 1
)

const (
	PdfParseTypeText         = 1
	PdfParseTypeOcr          = 2
	PdfParseTypeOcrWithImage = 3
	PdfParseTypeOcrAli       = 4
)

const (
	OpenLibraryDefault = `default`
)

const (
	McpServerPublished = 1
	McpServerDraft     = 0
)

const (
	McpClientTypeSse  = 1
	McpClientTypeHttp = 2
)

// 1指定事件,2指定时间,3手动执行
const (
	RunTypeEvent = 1
	RunTypeCron  = 2
	RunTypeHand  = 3
)

// 1 删除评论,2 回复评论,3置顶评论
const (
	CommentExecTypeDelete = 1
	CommentExecTypeReply  = 2
	CommentExecTypeTop    = 3
)

var CommentExecTypeMap = map[int]string{
	CommentExecTypeDelete: "自动删除",
	CommentExecTypeReply:  "自动回复",
	CommentExecTypeTop:    "自动精选",
}

const (
	CommentCheckTypeHintKeywords = 1
	CommentCheckTypeAICheck      = 2
)

const (
	BatchSendStatusDelete = -1
	BatchSendStatusWait   = 0
	BatchSendStatusExec   = 1
	BatchSendStatusSucc   = 2
	BatchSendStatusErr    = 3
)
const DefaultLibDocBanner = `/upload/default/open_doc_home_default_banner.png`

const (
	SyncOfficialHistoryTypeHalfYear  = 1
	SyncOfficialHistoryTypeOneYear   = 2
	SyncOfficialHistoryTypeThreeYear = 3
	SyncOfficialHistoryTypeAll       = 10
)

var SyncOfficialHistoryTypeList = [...]int{
	SyncOfficialHistoryTypeHalfYear,
	SyncOfficialHistoryTypeOneYear,
	SyncOfficialHistoryTypeThreeYear,
	SyncOfficialHistoryTypeAll,
}

const (
	SyncOfficialContentStatusNotStart = 1
	SyncOfficialContentStatusWorking  = 2
	SyncOfficialContentStatusFailed   = 3
)

const (
	IsDefault  = 1 //library and robot is default
	NotDefault = 2 //library and robot is not default
)

var (
	RobotPaymentAuthCodeUsageStatusPending   = 1 // 未使用
	RobotPaymentAuthCodeUsageStatusExchanged = 2 // 已兑换
	RobotPaymentAuthCodeUsageStatusUsed      = 3 // 已使用
	RobotPaymentPackageTypeCount             = 1
	RobotPaymentPackageTypeDuration          = 2
)

const (
	RobotPaymentAuthCodePrefix             = "GLY##"
	RobotPaymentAuthCodeSuffix             = "##GLY"
	RobotPaymentAuthCodeManagerCachePrefix = "robot_payment_auth_code_manager_"
)
