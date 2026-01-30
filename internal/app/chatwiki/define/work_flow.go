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
	LibraryImportContent = `content` //按内容导入
	LibraryImportUrl     = `url`     //按url导入
)

const (
	LibraryImportRepeatNotImport = `not import` //重复时继续导入
	LibraryImportRepeatUpdate    = `update`     //重复时更新
	LibraryImportRepeatImport    = `import`     //url存在时继续导入
)

const (
	QuestionAnswerTypeText = `text`
	QuestionAnswerTypeMenu = `menu`
)
