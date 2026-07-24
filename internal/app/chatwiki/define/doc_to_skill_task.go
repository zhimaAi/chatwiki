// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "time"

const TableDocToSkillTask = `chat_ai_doc_to_skill_task`

const (
	DocToSkillTaskStatusQueued   = 0
	DocToSkillTaskStatusRunning  = 1
	DocToSkillTaskStatusSucceed  = 2
	DocToSkillTaskStatusFailed   = 3
	DocToSkillTaskStatusStopping = 4
	DocToSkillTaskStatusStopped  = 5
)

const (
	DocToSkillTaskDefaultPageSize = 20
	DocToSkillTaskDefaultTemp     = 1
	DocToSkillTaskDefaultMaxToken = 32 * 1024
	DocToSkillTaskMaxFileCount    = 20
	DocToSkillTaskFileLimitSize   = 100 * 1024 * 1024
	DocToSkillTaskMaxIteration    = 1000
	DocToSkillReductionMaxTokens  = 64 * 1024
	DocToSkillReductionMinTokens  = 8 * 1024
	DocToSkillReductionKeepRounds = 4
	DocToSkillTaskStopKeyTTL      = 24 * time.Hour
)

const (
	DocToSkillTaskStopKeyPrefix = `chatwiki.doc_to_skill.stop.`
	DocToSkillWorkDir           = `clawbot/working_dir/doc-to-skill/<task_batch>`
)

var (
	DocToSkillTaskAllowExt    = []string{`txt`, `md`, `docx`, `pdf`}
	DocToSkillTaskZipAllowExt = []string{`zip`}
)

type DocToSkillTaskListFilter struct {
	Status int `form:"status" json:"status"`
	Page   int `form:"page" json:"page"`
	Size   int `form:"size" json:"size"`
}

type DocToSkillTaskCreateParams struct {
	CustomPrompt  string   `form:"custom_prompt" json:"custom_prompt"`
	ModelConfigId int      `form:"model_config_id" json:"model_config_id" binding:"required,gt=0"`
	UseModel      string   `form:"use_model" json:"use_model" binding:"required"`
	Temperature   *float32 `form:"temperature" json:"temperature" binding:"omitempty,gte=0,lte=2"`
	MaxToken      *int     `form:"max_token" json:"max_token" binding:"omitempty,gt=0"`
}

type DocToSkillTaskIDParams struct {
	ID int64 `form:"id" json:"id" binding:"required,gt=0"`
}

type DocToSkillTaskInstallParams struct {
	ID        int64 `form:"id" json:"id" binding:"required,gt=0"`
	Overwrite bool  `form:"overwrite" json:"overwrite"`
}

type DocToSkillTaskItem struct {
	ID               int64         `json:"id"`
	AdminUserId      int           `json:"admin_user_id,omitempty"`
	TaskBatch        string        `json:"task_batch"`
	SourceFiles      []*UploadInfo `json:"source_files"`
	CustomPrompt     string        `json:"custom_prompt"`
	ModelConfigId    int           `json:"model_config_id"`
	UseModel         string        `json:"use_model"`
	Temperature      float32       `json:"temperature"`
	MaxToken         int           `json:"max_token"`
	Status           int           `json:"status"`
	SkillName        string        `json:"skill_name"`
	SkillDescription string        `json:"skill_description"`
	FileName         string        `json:"file_name"`
	FileUrl          string        `json:"file_url"`
	FileSize         int           `json:"file_size"`
	DebugLog         []any         `json:"debug_log,omitempty"`
	ErrMsg           string        `json:"err_msg"`
	StartTime        int           `json:"start_time"`
	EndTime          int           `json:"end_time"`
	CreateTime       int           `json:"create_time"`
	UpdateTime       int           `json:"update_time"`
}
