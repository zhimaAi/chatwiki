// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/gemini"

type GeminiStreamResult struct {
	*gemini.ChatCompletionStream
}

func (c *GeminiStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseGemini, err := c.ChatCompletionStream.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseGemini.Candidates[0].Content.Parts[0].Text,
	}, nil
}
