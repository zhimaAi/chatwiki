// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package spark

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gorilla/websocket"
)

type ChatCompletionRequest struct {
	Header    RequestHeader  `json:"header"`
	Parameter Parameter      `json:"parameter"`
	Payload   RequestPayload `json:"payload"`
}
type RequestHeader struct {
	APPID string `json:"app_id"`
}
type Parameter struct {
	Chat Chat `json:"chat"`
}
type RequestPayload struct {
	Message RequestMessage `json:"message"`
}
type RequestMessage struct {
	Text []ChatCompletionMessage `json:"text"`
}
type Chat struct {
	Domain      string  `json:"domain"`
	Temperature float64 `json:"temperature"`
	TopK        int     `json:"top_k"`
	MaxTokens   int     `json:"max_tokens"`
	Auditing    string  `json:"auditing"`
	Stream      bool    `json:"stream"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Header  ResponseHeader  `json:"header"`
	Payload ResponsePayload `json:"payload"`
}
type ResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"`
}
type ResponsePayload struct {
	Choices Choice `json:"choices"`
	Usage   Usage  `json:"usage"`
}
type Choice struct {
	Status int                     `json:"status"`
	Seq    int                     `json:"seq"`
	Text   []ChatCompletionMessage `json:"text"`
}
type Usage struct {
	Text Text `json:"text"`
}
type Text struct {
	QuestionTokens   int `json:"question_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type ChatCompletionStream struct {
	conn *websocket.Conn
}

func (c *ChatCompletionStream) Recv() (ChatCompletionResponse, error) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return ChatCompletionResponse{}, err
		}

		var response ChatCompletionResponse
		err = json.Unmarshal(msg, &response)
		if err != nil {
			return ChatCompletionResponse{}, err
		}
		if response.Header.Code != 0 {
			return ChatCompletionResponse{}, errors.New(response.Header.Message)
		}
		if response.Payload.Choices.Status == 2 {
			return response, io.EOF
		}
		if len(response.Payload.Choices.Text) <= 0 {
			return ChatCompletionResponse{}, errors.New("no text in response")
		}
		return response, nil
	}
}

func (c *ChatCompletionStream) Close() error {
	return c.conn.Close()
}
