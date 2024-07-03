// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ollama

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

func NewClient(EndPoint, DeploymentName string) *Client {
	return &Client{
		EndPoint:       EndPoint,
		DeploymentName: DeploymentName,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/api/embeddings"
	responseRaw, err := common.HttpPost(url, nil, nil, req)
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
	if len(result.Embedding) <= 0 {
		return EmbeddingResponse{}, errors.New("ollama response No embedding result")
	}
	return result, err
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/api/chat"
	bools := false
	req.Stream = &bools
	responseRaw, err := common.HttpPost(url, nil, nil, req)
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
	if len(result.Message.Content) <= 0 {
		return ChatCompletionResponse{}, errors.New("ollama response No choices result")
	}
	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/api/chat"
	responseRaw, err := common.HttpStreamPost(url, nil, nil, req)
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
