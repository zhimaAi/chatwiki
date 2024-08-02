// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"chatwiki/internal/app/chatwiki/llm/api/spark"
	"errors"
	"io"
)

type SparkStreamResult struct {
	*spark.ChatCompletionStream
}

func (r *SparkStreamResult) Read() (resp ZhimaChatCompletionResponse, err error) {
	responseSpark, err := r.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			resp = ZhimaChatCompletionResponse{
				Result:          responseSpark.Payload.Choices.Text[0].Content,
				PromptToken:     responseSpark.Payload.Usage.Text.PromptTokens,
				CompletionToken: responseSpark.Payload.Usage.Text.CompletionTokens,
			}
		}
	} else {
		resp = ZhimaChatCompletionResponse{
			Result:          responseSpark.Payload.Choices.Text[0].Content,
			PromptToken:     responseSpark.Payload.Usage.Text.PromptTokens,
			CompletionToken: responseSpark.Payload.Usage.Text.CompletionTokens,
		}
	}

	return
}
