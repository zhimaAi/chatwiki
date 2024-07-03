// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package openai

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	Stream           bool                    `json:"stream,omitempty"`
	FrequencyPenalty int                     `json:"frequency_penalty,omitempty"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	N                int                     `json:"n,omitempty"`
	PresencePenalty  int                     `json:"presence_penalty,omitempty"`
	ResponseFormat   string                  `json:"response_format,omitempty"`
	Seed             int                     `json:"seed,omitempty"`
	Temperature      float64                 `json:"temperature,omitempty"`
	TopP             int                     `json:"top_p,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionChoice struct {
	Message ChatCompletionMessage `json:"message"`
}
type ChatCompletionStreamChoice struct {
	Index        int                   `json:"index"`
	Delta        ChatCompletionMessage `json:"delta"`
	FinishReason string                `json:"finish_reason"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage"`
}
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Choices []ChatCompletionStreamChoice `json:"choices,omitempty"`
	Usage   ChatCompletionUsage          `json:"usage"`
}

type ChatCompletionStream struct {
	*common.StreamReader[ChatCompletionStreamResponse]
}

var (
	headerData  = []byte("data:")
	errorPrefix = []byte(`data: {"error":`)
)

func (c *ChatCompletionStream) Recv() (ChatCompletionStreamResponse, error) {
	if c.StreamReader.IsFinished {
		return ChatCompletionStreamResponse{}, io.EOF
	}

	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := c.StreamReader.Reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			if readErr != io.EOF {
				c.StreamReader.UnmarshalError()
				if c.StreamReader.ErrorResponse != nil {
					return *new(ChatCompletionStreamResponse), fmt.Errorf("unmarshal error, %w", c.StreamReader.ErrorResponse.Error())
				}
				return *new(ChatCompletionStreamResponse), readErr
			} else {
				c.StreamReader.IsFinished = true
				return *new(ChatCompletionStreamResponse), io.EOF
			}
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
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
		if strings.TrimSpace(string(noPrefixLine)) == "[DONE]" {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionStreamResponse), io.EOF
		}

		var response ChatCompletionStreamResponse
		unmarshalErr := json.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionStreamResponse), unmarshalErr
		}
		if len(response.Choices) <= 0 {
			return *new(ChatCompletionStreamResponse), errors.New("response no choices result")
		}

		return response, nil
	}
}
