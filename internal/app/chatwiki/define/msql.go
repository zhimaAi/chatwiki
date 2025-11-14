// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
	RobotManage      = `RobotManage`
	LibraryManage    = `LibraryManage`
	SystemManage     = `SystemManage`
	ClientSideManage = `ClientSideManage`
)

const (
	GeneralLibraryType = 0
	QALibraryType      = 2
	OpenLibraryType    = 1
)

var LibraryTypes = [...]int{
	GeneralLibraryType, //0普通知识库
	QALibraryType,      //2问答知识库
	OpenLibraryType,    //1对外知识库
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

const DefaultLibDocBanner = `/upload/default/open_doc_home_default_banner.png`
