// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import "chatwiki/internal/app/chatwiki/llm/api/baidu"

type BaiduStreamResult struct {
	*baidu.ChatCompletionStream
}

func (r *BaiduStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseBaidu, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result: responseBaidu.Result,
	}, nil
}
