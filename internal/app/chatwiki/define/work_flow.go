// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
