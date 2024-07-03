// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/claude"

type ClaudeStreamResult struct {
	*claude.ChatCompletionStream
}

func (r *ClaudeStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseClaude, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseClaude.Delta.Text,
	}, nil
}
