// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package gemini

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"io"
)

type Client struct {
	EndPoint string
	APIKey   string
	Model    string
}

func NewClient(APIKey, Model string) *Client {
	return &Client{
		EndPoint: "https://generativelanguage.googleapis.com/v1",
		APIKey:   APIKey,
		Model:    Model,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {
	url := c.EndPoint + "/models/" + c.Model + ":embedContent"
	params := []common.Param{
		{Key: "key", Value: c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, nil, params, req)
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
	url := c.EndPoint + "/models/" + c.Model + ":generateContent"
	params := []common.Param{
		{Key: "key", Value: c.APIKey},
	}

	responseRaw, err := common.HttpPost(url, nil, params, req)
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
	if len(result.Candidates) <= 0 || len(result.Candidates[0].Content.Parts) <= 0 {
		return ChatCompletionResponse{}, errors.New("gemini response candidates is empty")
	}
	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/models/" + c.Model + ":streamGenerateContent"
	params := []common.Param{
		{Key: "key", Value: c.APIKey},
	}

	responseRaw, err := common.HttpStreamPost(url, nil, params, req)
	if err != nil {
		return nil, err
	}

	err = common.HttpCheckErrors[*ErrorResponse](responseRaw)
	if err != nil {
		return nil, err
	}

	var errResp ErrorResponse
	streamResp := &common.StreamReader[ChatCompletionResponse]{
		EmptyMessagesLimit: 3000,
		Reader:             bufio.NewReader(responseRaw.Body),
		Response:           responseRaw,
		ErrAccumulator:     common.NewErrorAccumulator(),
		ErrorResponse:      &errResp,
		HttpHeader:         responseRaw.Header,
	}

	return &ChatCompletionStream{StreamReader: streamResp}, nil
}
