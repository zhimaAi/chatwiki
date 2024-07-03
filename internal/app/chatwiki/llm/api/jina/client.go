// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package jina

import (
	"chatwiki/internal/app/chatwiki/llm/common"
	"io"
)

type Client struct {
	EndPoint string
	APIKey   string
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://api.jina.ai/v1",
		APIKey:   APIKey,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {
	url := c.EndPoint + "/embeddings"
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

func (c *Client) ReRank(req ReRankRequest) (ReRankResponse, error) {
	url := c.EndPoint + "/rerank"
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
