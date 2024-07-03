// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package cohere

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"io"
)

type Client struct {
	EndPoint string
	APIKey   string
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://api.cohere.com",
		APIKey:   APIKey,
	}
}

func (c *Client) ReRank(req ReRankRequest) (ReRankResponse, error) {
	url := c.EndPoint + "/v1/rerank"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return ReRankResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(responseRaw.Body)

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return ReRankResponse{}, err
	}

	var result ReRankResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	return result, err
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {
	url := c.EndPoint + "/v1/embed"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return EmbeddingResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(responseRaw.Body)

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return EmbeddingResponse{}, err
	}

	var result EmbeddingResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	return result, err
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/v1/chat"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	req.Stream = false
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
	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/v1/chat"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	req.Stream = true
	responseRaw, err := common.HttpPost(url, headers, nil, req)
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
