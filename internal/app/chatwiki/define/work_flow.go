// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

type NodeLogs struct {
	StartTime any `json:"start_time"`
	Output    struct {
		LlmResult struct {
			CompletionToken any `json:"completion_token"`
			PromptToken     any `json:"prompt_token"`
		} `json:"llm_result"`
	} `json:"output"`
	UseTime any `json:"use_time"`
}

const (
	FinishNodeOutTypeVariable = "variable"
	FinishNodeOutTypeMessage  = "message"
)

const FinishReplyPrefixKey = `special.finish_reply_content_`

const (
	// WorkFlowStatusRunning means the workflow is still running.
	WorkFlowStatusRunning = 0
	// WorkFlowStatusCompleted means the workflow completed normally.
	WorkFlowStatusCompleted = 1
	// WorkFlowStatusStopped means the workflow was stopped manually or by cancellation.
	WorkFlowStatusStopped = 2
)

const (
	LibraryImportContent = `content` // Import by content
	LibraryImportUrl     = `url`     // Import by URL
)

const (
	LibraryImportRepeatNotImport = `not import` // Continue import when duplicate
	LibraryImportRepeatUpdate    = `update`     // Update when duplicate
	LibraryImportRepeatImport    = `import`     // Continue import when URL exists
)

const (
	QuestionAnswerTypeText = `text`
	QuestionAnswerTypeMenu = `menu`
)
