// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/spark"

type SparkStreamResult struct {
	*spark.ChatCompletionStream
}

func (r *SparkStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseSpark, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseSpark.Payload.Choices.Text[0].Content,
	}, nil
}
