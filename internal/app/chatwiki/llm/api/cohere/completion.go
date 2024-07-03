// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package cohere

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"fmt"
	"io"
)

type ChatCompletionRequest struct {
	Message           string        `json:"message"`
	Stream            bool          `json:"stream,omitempty"`
	Preamble          string        `json:"preamble,omitempty"`
	ChatHistory       []ChatHistory `json:"chat_history,omitempty"`
	ConversationID    string        `json:"conversation_id,omitempty"`
	PromptTruncation  string        `json:"prompt_truncation,omitempty"`
	Connectors        []Connector   `json:"connectors,omitempty"`
	SearchQueriesOnly bool          `json:"search_queries_only,omitempty"`
	Documents         []Document    `json:"documents,omitempty"`
	CitationQuality   string        `json:"citation_quality,omitempty"`
	Temperature       float64       `json:"temperature,omitempty"`
	MaxTokens         int           `json:"max_tokens,omitempty"`
	MaxInputTokens    int           `json:"max_input_tokens,omitempty"`
	K                 int           `json:"k,omitempty"`
	P                 int           `json:"p,omitempty"`
	Seed              int           `json:"seed,omitempty"`
	StopSequences     []string      `json:"stop_sequences,omitempty"`
	FrequencyPenalty  int           `json:"frequency_penalty,omitempty"`
	PresencePenalty   int           `json:"presence_penalty,omitempty"`
}

type ChatHistory struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type Connector struct {
	ID              string `json:"id"`
	UserAccessToken string `json:"access_token,omitempty"`
}

type ChatCompletionResponse struct {
	Text         string       `json:"text"`
	GenerationId string       `json:"generation_id"`
	Documents    []Document   `json:"documents"`
	Meta         ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	APIVersion  APIVersion  `json:"api_version"`
	BilledUnits BilledUnits `json:"billed_units"`
	Tokens      Tokens      `json:"tokens"`
	Warnings    []string    `json:"warnings"`
}

type ChatCompletionStreamResponse struct {
	IsFinished   bool     `json:"is_finished"`
	EventType    string   `json:"event_type"`
	Text         string   `json:"text"`
	Response     Response `json:"response"`
	FinishReason string   `json:"finish_reason"`
}

type Response struct {
	ResponseID   string        `json:"response_id"`
	Text         string        `json:"text"`
	GenerationID string        `json:"generation_id"`
	ChatHistory  []ChatHistory `json:"chat_history"`
	FinishReason string        `json:"finish_reason"`
	Meta         ResponseMeta  `json:"meta"`
}

type ChatCompletionStream struct {
	*common.StreamReader[ChatCompletionStreamResponse]
}

func (c *ChatCompletionStream) Recv() (ChatCompletionStreamResponse, error) {
	if c.StreamReader.IsFinished {
		return ChatCompletionStreamResponse{}, io.EOF
	}

	for {
		rawLine, readErr := c.StreamReader.Reader.ReadBytes('\n')
		if readErr != nil {
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
		var response ChatCompletionStreamResponse
		unmarshalErr := json.Unmarshal(noSpaceLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionStreamResponse), unmarshalErr
		}
		if response.EventType == "stream-end" || response.IsFinished {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionStreamResponse), nil
		}

		return response, nil
	}
}
