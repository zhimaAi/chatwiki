// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/openai"

type OpenAIStreamResult struct {
	*openai.ChatCompletionStream
}

func (r *OpenAIStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseOpenAI, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}

	return ZhimaChatCompletionResponse{
		Result: responseOpenAI.Choices[0].Delta.Content,
	}, nil
}
