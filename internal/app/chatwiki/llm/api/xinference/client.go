// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package xinference

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"github.com/zhimaAi/go_tools/logs"
	"io"
)

type Client struct {
	EndPoint       string
	APIVersion     string
	APIKey         string
	DeploymentName string
}

func NewClient(EndPoint, apiVersion, DeploymentName string) *Client {
	return &Client{
		EndPoint:       EndPoint,
		APIVersion:     apiVersion,
		DeploymentName: DeploymentName,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {
	url := c.EndPoint + "/" + c.APIVersion + "/embeddings"
	headers := []common.Header{
		{Key: "Context-Type", Value: "application/json"},
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
	if err != nil {
		return EmbeddingResponse{}, err
	}
	if len(result.Data) <= 0 {
		return EmbeddingResponse{}, errors.New("xinference response No embedding result")
	}
	return result, err
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {

	url := c.EndPoint + "/" + c.APIVersion + "/chat/completions"
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
	if len(result.Choices) <= 0 {
		return ChatCompletionResponse{}, errors.New("xinference response No embedding result")
	}
	logs.Info("CreateChatCompletion:%+v,res:%+v", req, result)
	return result, err
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {

	url := c.EndPoint + "/" + c.APIVersion + "/chat/completions"
	headers := []common.Header{
		{Key: "Context-Type", Value: "application/json"},
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

func (r *CreateRerankReq) validate() error {
	return nil
}
func (c *Client) CreateRerank(req *CreateRerankReq) (*CreateRerankRes, error) {
	result := &CreateRerankRes{}
	if err := req.validate(); err != nil {
		return result, err
	}
	url := c.EndPoint + "/" + c.APIVersion + "/rerank"
	headers := []common.Header{
		{Key: "Context-Type", Value: "application/json"},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return result, err
	}

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return result, err
	}

	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return result, err
	}
	logs.Info("CreateRerank:%+v,res:%+v", req, result)
	return result, err
}
