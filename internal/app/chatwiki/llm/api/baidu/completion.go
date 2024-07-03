// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package baidu

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model           string                  `json:"model"`
	Messages        []ChatCompletionMessage `json:"messages"`
	Stream          bool                    `json:"stream,omitempty"`
	Temperature     float64                 `json:"temperature,omitempty"`
	TopP            float32                 `json:"top_p,omitempty"`
	PenaltyScore    float32                 `json:"penalty_score,omitempty"`
	System          string                  `json:"system,omitempty"`
	Stop            []string                `json:"stop,omitempty"`
	DisableSearch   bool                    `json:"disable_search,omitempty"`
	EnableCitation  bool                    `json:"enable_citation,omitempty"`
	EnableTrace     bool                    `json:"enable_trace,omitempty"`
	MaxOutputTokens int                     `json:"max_output_tokens,omitempty"`
	ResponseFormat  string                  `json:"response_format,omitempty"`
	UserId          string                  `json:"user_id,omitempty"`
}

type ChatCompletionResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	FinishReason     string `json:"finish_reason"`
	Usage            Usage  `json:"usage"`
	Result           string `json:"result"`
}
type ChatCompletionStreamResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	IsTruncated      bool   `json:"is_truncated"`
	IsEnd            bool   `json:"is_end"`
	NeedClearHistory bool   `json:"need_clear_history"`
	FinishReason     string `json:"finish_reason"`
	Usage            Usage  `json:"usage"`
	Result           string `json:"result"`
}
type ChatCompletionStream struct {
	*common.StreamReader[ChatCompletionStreamResponse]
}

func (c *ChatCompletionStream) Recv() (ChatCompletionStreamResponse, error) {
	if c.StreamReader.IsFinished {
		return ChatCompletionStreamResponse{}, io.EOF
	}

	var emptyMessagesCount uint
	var headerData = []byte("data: ")
	var errorPrefix = []byte(`{"error`)

	for {
		rawLine, readErr := c.StreamReader.Reader.ReadBytes('\n')
		if readErr != nil {
			if readErr != io.EOF {
				c.StreamReader.UnmarshalError()
				if c.StreamReader.ErrorResponse != nil {
					return *new(ChatCompletionStreamResponse), fmt.Errorf("unmarshal error, %w", c.StreamReader.ErrorResponse.Error())
				}
				return *new(ChatCompletionStreamResponse), readErr
			}
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if !bytes.HasPrefix(noSpaceLine, headerData) {
			if bytes.HasPrefix(noSpaceLine, errorPrefix) {
				var errResp ErrorResponse
				err := json.Unmarshal(noSpaceLine, &errResp)
				if err != nil {
					return *new(ChatCompletionStreamResponse), fmt.Errorf("unmarshal error, %w", c.StreamReader.ErrorResponse.Error())
				} else {
					errResp.SetHTTPStatusCode(c.Response.StatusCode)
					return *new(ChatCompletionStreamResponse), errResp.Error()
				}
			}

			writeErr := c.StreamReader.ErrAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return *new(ChatCompletionStreamResponse), writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > c.StreamReader.EmptyMessagesLimit {
				return *new(ChatCompletionStreamResponse), errors.New("stream has sent too many empty messages")
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)

		var response ChatCompletionStreamResponse
		unmarshalErr := json.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionStreamResponse), unmarshalErr
		}
		if response.IsEnd {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionStreamResponse), nil
		}

		return response, nil
	}
}
