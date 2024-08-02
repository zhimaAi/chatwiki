// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"chatwiki/internal/app/chatwiki/llm/api/ollama"
)

type OllamaStreamResult struct {
	*ollama.ChatCompletionStream
}

func (r *OllamaStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseOllama, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result:          responseOllama.Message.Content,
		PromptToken:     responseOllama.Metrics.PromptEvalCount,
		CompletionToken: responseOllama.Metrics.EvalCount,
	}, nil
}
