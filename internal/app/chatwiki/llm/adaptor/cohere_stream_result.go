// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/cohere"

type CohereStreamResult struct {
	*cohere.ChatCompletionStream
}

func (r *CohereStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseCohere, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseCohere.Text,
	}, nil
}
