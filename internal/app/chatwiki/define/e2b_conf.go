// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const TableChatAiRobotE2bConf = "chat_ai_robot_e2b_conf"

const (
	E2bConfApiKeyMaxLen        = 500
	E2bConfApiBaseUrlMaxLen    = 1000
	E2bConfSandboxDomainMaxLen = 1000
	E2bConfTemplateMaxLen      = 500
	E2bConfCommandUserMaxLen   = 100
)

type E2bConfGetParams struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

type E2bConfParams struct {
	RobotKey       string `form:"robot_key" json:"robot_key" binding:"required"`
	SwitchStatus   int    `form:"switch_status" json:"switch_status" binding:"oneof=0 1"`
	ApiKey         string `form:"api_key" json:"api_key"`
	ApiBaseUrl     string `form:"api_base_url" json:"api_base_url"`
	SandboxDomain  string `form:"sandbox_domain" json:"sandbox_domain"`
	Template       string `form:"template" json:"template"`
	Timeout        int    `form:"timeout" json:"timeout"`
	CommandTimeout int    `form:"command_timeout" json:"command_timeout"`
	CommandUser    string `form:"command_user" json:"command_user"`
}
