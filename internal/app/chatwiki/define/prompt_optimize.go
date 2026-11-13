// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"time"

	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
)

const (
	TablePromptOptimizeRecord       = "chat_ai_prompt_optimize_record"
	TablePromptOptimizeVersion      = "chat_ai_prompt_optimize_version"
	PromptOptimizeStatusPending     = "pending"
	PromptOptimizeStatusSuccess     = "success"
	PromptOptimizeStatusFailed      = "failed"
	PromptOptimizeStatusCancelled   = "cancelled"
	PromptOptimizeSceneConfig       = "config"
	PromptOptimizeSceneConversation = "conversation"
	PromptOptimizeModeOptimize      = "optimize"
	PromptOptimizeModeRewrite       = "rewrite"
	PromptOptimizePromptOpenTag     = "<optimized_prompt>"
	PromptOptimizePromptCloseTag    = "</optimized_prompt>"
	PromptOptimizeReplyOpenTag      = "<explanation>"
	PromptOptimizeReplyCloseTag     = "</explanation>"
)

const (
	PromptOptimizeHistoryWindow   = 7 * 24 * time.Hour
	PromptOptimizeHistoryRounds   = 100
	PromptOptimizeDefaultPageSize = 20
	PromptOptimizePageSizeLimit   = 100
	PromptOptimizeRequestIDLength = 64
	PromptOptimizeRequestTimeout  = 5 * time.Minute
	PromptOptimizeTemperature     = 0.5
	PromptOptimizeMaxTokens       = 4096
)

type PromptOptimizeRound struct {
	ID        int64
	CreatedAt int64
	Status    string
	Input     string
	Prompt    string
	Reply     string
}

type PromptOptimizeRobotParams struct {
	RobotID int `form:"robot_id" json:"robot_id" binding:"required,gt=0"`
}
type PromptOptimizeRecordsParams struct {
	PromptOptimizeRobotParams
	BeforeID int64 `form:"before_id" json:"before_id" binding:"gte=0"`
	Size     int   `form:"size" json:"size" binding:"gt=0"`
}
type PromptOptimizeVersionParams struct {
	PromptOptimizeRobotParams
	VersionID int64 `form:"version_id" json:"version_id" binding:"required,gt=0"`
}
type PromptOptimizeGenerateParams struct {
	PromptOptimizeRobotParams
	Scene           string `form:"scene" json:"scene" binding:"required,oneof=config conversation"`
	Mode            string `form:"mode" json:"mode" binding:"required,oneof=optimize rewrite"`
	Demand          string `form:"demand" json:"demand" binding:"required"`
	RequestID       string `form:"request_id" json:"request_id" binding:"required"`
	ModelConfigID   int    `form:"model_config_id" json:"model_config_id" binding:"required,gt=0"`
	UseModel        string `form:"use_model" json:"use_model" binding:"required"`
	ReasoningEffort string `form:"reasoning_effort" json:"reasoning_effort"`
	BaseVersionID   int64  `form:"base_version_id" json:"base_version_id" binding:"gte=0"`
	MessageID       int64  `form:"message_id" json:"message_id" binding:"gte=0"`
}
type PromptOptimizeName struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}
type PromptOptimizeContext struct {
	Prompt          string               `json:"prompt"`
	ApplicationType int                  `json:"application_type"`
	ModelConfigID   string               `json:"model_config_id"`
	UseModel        string               `json:"use_model"`
	Names           []PromptOptimizeName `json:"names"`
}
type PromptOptimizeRecordList struct {
	List    []msql.Params `json:"list"`
	HasMore bool          `json:"has_more"`
}
type PromptOptimizeVersion struct {
	ID         int64  `json:"id"`
	Version    int    `json:"version"`
	RecordID   int64  `json:"record_id"`
	Prompt     string `json:"prompt"`
	CreateTime int64  `json:"create_time"`
}
type PromptOptimizeOutput struct {
	Prompt string
	Reply  string
}
type PromptOptimizeDelta struct {
	Text string `json:"text"`
}
type PromptOptimizeDone struct {
	RecordID         int64                  `json:"record_id"`
	Reply            string                 `json:"reply"`
	ReasoningContent string                 `json:"reasoning_content"`
	Version          *PromptOptimizeVersion `json:"version,omitempty"`
}
type PromptOptimizePrepared struct {
	Robot            msql.Params
	RecordID         int64
	Params           PromptOptimizeGenerateParams
	Messages         []chat.Message
	Names            []PromptOptimizeName
	HistoryRounds    int
	BasePromptLength int
}
type PromptOptimizeInput struct {
	Request      string               `json:"request"`
	BasePrompt   string               `json:"base_prompt"`
	Capabilities []PromptOptimizeName `json:"capabilities"`
	Question     string               `json:"question,omitempty"`
	Answer       string               `json:"answer,omitempty"`
}
