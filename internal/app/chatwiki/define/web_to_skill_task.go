// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "time"

const TableWebToSkillTask = `chat_ai_web_to_skill_task`

const (
	WebToSkillTaskStatusQueued   = 0
	WebToSkillTaskStatusRunning  = 1
	WebToSkillTaskStatusSucceed  = 2
	WebToSkillTaskStatusFailed   = 3
	WebToSkillTaskStatusStopping = 4
	WebToSkillTaskStatusStopped  = 5
)

const (
	WebToSkillTaskDefaultPageSize  = 20
	WebToSkillTaskDefaultTemp      = 1
	WebToSkillTaskDefaultMaxToken  = 32 * 1024 // 32k
	WebToSkillTaskStopKeyTTL       = 24 * time.Hour
	WebToSkillTaskDNSLookupTimeout = time.Second
)

const WebToSkillTaskStopKeyPrefix = `chatwiki.web_to_skill.stop.`

var WebToSkillTaskZipAllowExt = []string{`zip`}

type WebToSkillTaskListFilter struct {
	Status int `form:"status" json:"status"`
	Page   int `form:"page" json:"page"`
	Size   int `form:"size" json:"size"`
}

type WebToSkillTaskCreateParams struct {
	Urls          []string `form:"urls" json:"urls"`
	CustomPrompt  string   `form:"custom_prompt" json:"custom_prompt"`
	ModelConfigId int      `form:"model_config_id" json:"model_config_id" binding:"required,gt=0"`
	UseModel      string   `form:"use_model" json:"use_model" binding:"required"`
	Temperature   *float32 `form:"temperature" json:"temperature" binding:"omitempty,gte=0,lte=2"`
	MaxToken      *int     `form:"max_token" json:"max_token" binding:"omitempty,gt=0"`
}

type WebToSkillTaskIDParams struct {
	ID int64 `form:"id" json:"id" binding:"required,gt=0"`
}

type WebToSkillTaskInstallParams struct {
	ID        int64 `form:"id" json:"id" binding:"required,gt=0"`
	Overwrite bool  `form:"overwrite" json:"overwrite"`
}

type WebToSkillTaskItem struct {
	ID               int64    `json:"id"`
	AdminUserId      int      `json:"admin_user_id,omitempty"`
	TaskBatch        string   `json:"task_batch"`
	Urls             []string `json:"urls"`
	CustomPrompt     string   `json:"custom_prompt"`
	ModelConfigId    int      `json:"model_config_id"`
	UseModel         string   `json:"use_model"`
	Temperature      float32  `json:"temperature"`
	MaxToken         int      `json:"max_token"`
	Status           int      `json:"status"`
	SkillName        string   `json:"skill_name"`
	SkillDescription string   `json:"skill_description"`
	FileName         string   `json:"file_name"`
	FileUrl          string   `json:"file_url"`
	FileSize         int      `json:"file_size"`
	DebugLog         []any    `json:"debug_log,omitempty"`
	ErrMsg           string   `json:"err_msg"`
	StartTime        int      `json:"start_time"`
	EndTime          int      `json:"end_time"`
	CreateTime       int      `json:"create_time"`
	UpdateTime       int      `json:"update_time"`
}
