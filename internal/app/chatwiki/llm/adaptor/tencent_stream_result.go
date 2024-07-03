// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/hunyuan"

type TencentStreamResult struct {
	*hunyuan.ChatCompletionStream
}

func (r *TencentStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseTencent, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: *responseTencent.Choices[0].Delta.Content,
	}, nil
}
