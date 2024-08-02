// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package openai

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"io"
)

type Client struct {
	APIKey   string
	EndPoint string
	ErrResp  common.ErrorResponseInterface
}

func NewClient(EndPoint, apiKey string, ErrResp common.ErrorResponseInterface) *Client {
	return &Client{
		APIKey:   apiKey,
		EndPoint: EndPoint,
		ErrResp:  ErrResp,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/embeddings"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	if req.Model == "text-embedding-3-large" {
		req.Dimensions = 1536 //compatible
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return EmbeddingResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(responseRaw.Body)

	err = common.HttpCheckError(responseRaw, c.ErrResp)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	var result EmbeddingResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return EmbeddingResponse{}, err
	}
	if len(result.Data) <= 0 {
		return EmbeddingResponse{}, errors.New("OpenAI response no embedding result")
	}

	return result, err
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/chat/completions"
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
	err = common.HttpCheckError(responseRaw, c.ErrResp)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	var result ChatCompletionResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	if len(result.Choices) <= 0 {
		return ChatCompletionResponse{}, errors.New("OpenAI response no choices result")
	}

	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/chat/completions"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}
	req.Stream = true
	req.StreamOptions = &StreamOptions{IncludeUsage: true}
	responseRaw, err := common.HttpStreamPost(url, headers, nil, req)
	if err != nil {
		return nil, err
	}

	err = common.HttpCheckError(responseRaw, c.ErrResp)
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
