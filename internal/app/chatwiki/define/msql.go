// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const Postgres = `postgres`

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
	ExportSourceSession uint = 1 //会话记录导出
)

var ExportSourceList = []map[string]any{
	{`source`: ExportSourceSession, `source_name`: `会话记录导出`},
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
