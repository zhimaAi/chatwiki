// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package claude

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	System      string    `json:"system,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float32   `json:"top_p,omitempty"`
}

type Metadata struct {
	UserId string `json:"user_id,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	Model        string    `json:"model"`
	StopReason   string    `json:"stop_reason"`
	StopSequence string    `json:"stop_sequence"`
	Usage        Usage     `json:"usage"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"input_token"`
	OutputTokens int `json:"output_token"`
}

type ChatCompletionStreamResponse struct {
	Type         string                 `json:"type"`
	Message      ChatCompletionResponse `json:"message"`
	Index        int                    `json:"index"`
	ContentBlock Content                `json:"content_block"`
	Delta        Delta                  `json:"delta"`
	Usage        Usage                  `json:"usage"`
}
type Delta struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	StopReason   string `json:"stop_reason"`
	EndTurn      string `json:"end_turn"`
	StopSequence string `json:"stop_sequence"`
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
		if response.Type == "message_stop" {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionStreamResponse), nil
		}

		return response, nil
	}
}
