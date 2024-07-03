// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package azure

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"io"
)

type Client struct {
	EndPoint       string
	APIVersion     string
	APIKey         string
	DeploymentName string
}

func NewClient(EndPoint, APIVersion, APIKey, DeploymentName string) *Client {
	return &Client{
		EndPoint:       EndPoint,
		APIKey:         APIKey,
		APIVersion:     APIVersion,
		DeploymentName: DeploymentName,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/openai/deployments/" + c.DeploymentName + "/embeddings"
	headers := []common.Header{
		{Key: "api-key", Value: c.APIKey},
	}
	params := []common.Param{
		{Key: "api-version", Value: c.APIVersion},
	}
	responseRaw, err := common.HttpPost(url, headers, params, req)
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
	if err != nil {
		return EmbeddingResponse{}, err
	}
	if len(result.Data) <= 0 {
		return EmbeddingResponse{}, errors.New("azure response No embedding result")
	}

	return result, err
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/openai/deployments/" + c.DeploymentName + "/chat/completions"
	params := []common.Param{
		{Key: "api-version", Value: c.APIVersion},
	}
	headers := []common.Header{
		{Key: "api-key", Value: c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, headers, params, req)
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
	if len(result.Choices) <= 0 {
		return ChatCompletionResponse{}, errors.New("azure response No choices result")
	}

	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/openai/deployments/" + c.DeploymentName + "/chat/completions"
	params := []common.Param{
		{Key: "api-version", Value: c.APIVersion},
	}
	headers := []common.Header{
		{Key: "api-key", Value: c.APIKey},
	}
	req.Stream = true
	responseRaw, err := common.HttpStreamPost(url, headers, params, req)
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
