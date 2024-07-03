// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package claude

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"io"
)

type Client struct {
	EndPoint   string
	APIKey     string
	APIVersion string
}

func NewClient(APIKey, APIVersion string) *Client {
	return &Client{
		EndPoint:   "https://api.anthropic.com",
		APIKey:     APIKey,
		APIVersion: APIVersion,
	}
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/v1/messages"
	headers := []common.Header{
		{Key: "x-api-key", Value: c.APIKey},
		{Key: "anthropic-version", Value: c.APIVersion},
	}

	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(responseRaw.Body)

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	var result ChatCompletionResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	if len(result.Content) <= 0 {
		return ChatCompletionResponse{}, errors.New("claude response no content data")
	}

	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/v1/messages"
	headers := []common.Header{
		{Key: "x-api-key", Value: c.APIKey},
		{Key: "anthropic-version", Value: c.APIVersion},
	}

	req.Stream = true
	responseRaw, err := common.HttpStreamPost(url, headers, nil, req)
	if err != nil {
		return nil, err
	}

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return nil, err
	}

	var errResp ErrorResponse
	streamResp := &common.StreamReader[ChatCompletionStreamResponse]{
		EmptyMessagesLimit: 300,
		Reader:             bufio.NewReader(responseRaw.Body),
		Response:           responseRaw,
		ErrAccumulator:     common.NewErrorAccumulator(),
		ErrorResponse:      &errResp,
		HttpHeader:         responseRaw.Header,
	}

	return &ChatCompletionStream{StreamReader: streamResp}, nil
}
