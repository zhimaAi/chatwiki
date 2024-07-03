// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"chatwiki/internal/app/chatwiki/llm/api/xinference"
)

type XinferenceStreamResult struct {
	*xinference.ChatCompletionStream
}

func (r *XinferenceStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseAzure, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseAzure.Choices[0].Delta.Content,
	}, nil
}
