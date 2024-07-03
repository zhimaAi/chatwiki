// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package hunyuan

import (
	"encoding/json"
	"errors"
	"io"

	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tencentHunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type ChatCompletionStream struct {
	events chan common.SSEvent
	closed bool
}

func (c *ChatCompletionStream) Recv() (tencentHunyuan.ChatCompletionsResponseParams, error) {
	for {
		event, ok := <-c.events
		if !ok {
			c.closed = true
			return tencentHunyuan.ChatCompletionsResponseParams{}, io.EOF
		}
		var response tencentHunyuan.ChatCompletionsResponseParams
		err := json.Unmarshal(event.Data, &response)
		if err != nil {
			return tencentHunyuan.ChatCompletionsResponseParams{}, err
		}
		if len(response.Choices) < 1 {
			return tencentHunyuan.ChatCompletionsResponseParams{}, errors.New("tencent response no response text")
		}
		return response, nil
	}
}
func (c *ChatCompletionStream) Close() error {
	if !c.closed {
		close(c.events)
	}
	return nil
}
