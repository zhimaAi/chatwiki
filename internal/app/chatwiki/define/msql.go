// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "chatwiki/internal/pkg/lib_define"

const Postgres = lib_define.Postgres

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
	GeneralLibraryType,  //0 general library
	QALibraryType,       //2 Q&A library
	OpenLibraryType,     //1 open library
	OfficialLibraryType, //3 official account library
}

const (
	OpenLibraryAccessRights = 1
	AiSummary               = 1
)

const (
	ExportSourceSession    uint = 1 // session record export
	ExportSourceLibFileDoc uint = 2 // library document
)

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

// 1 specific event, 2 specific time, 3 manual execution
const (
	RunTypeEvent = 1
	RunTypeCron  = 2
	RunTypeHand  = 3
)

// 1 delete comment, 2 reply to comment, 3 pin comment
const (
	CommentExecTypeDelete = 1
	CommentExecTypeReply  = 2
	CommentExecTypeTop    = 3
)

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
	RobotPaymentAuthCodeUsageStatusPending   = 1 // not used
	RobotPaymentAuthCodeUsageStatusExchanged = 2 // exchanged
	RobotPaymentAuthCodeUsageStatusUsed      = 3 // used
	RobotPaymentPackageTypeCount             = 1
	RobotPaymentPackageTypeDuration          = 2
)

const (
	RobotPaymentAuthCodePrefix             = "GLY##"
	RobotPaymentAuthCodeSuffix             = "##GLY"
	RobotPaymentAuthCodeManagerCachePrefix = "robot_payment_auth_code_manager_"
)
