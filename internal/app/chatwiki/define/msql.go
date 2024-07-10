// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const Postgres = `postgres`

const (
	DefaultRoleRoot  = "所有者"
	DefaultRoleAdmin = "管理员"
	DefaultRoleUser  = "成员"
)

const (
	TableRole        = "role"
	TableUser        = "public.user"
	TableMenu        = "menu"
	TableCompany     = "company"
	TableFastCommand = "fast_command"
	TableRule        = "casbin_rule"
)

const (
	DefaultRoleIdRoot  = 1
	DefaultRoleIdAdmin = 2
	DefaultRoleIdUser  = 3
)

const (
	Normal  = "0"
	Deleted = "1"
)

const (
	RobotTextResponse     = 0
	RobotMarkdownResponse = 1
	RobotManage           = `RobotManage`
	LibraryManage         = `LibraryManage`
	SystemManage          = `SystemManage`
	ClientSideManage      = `ClientSideManage`
)
