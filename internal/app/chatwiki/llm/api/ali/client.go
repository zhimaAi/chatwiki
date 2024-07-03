// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ali

import (
	"chatwiki/internal/app/chatwiki/llm/api/openai"
	"chatwiki/internal/app/chatwiki/llm/common"
)

type Client struct {
	EndPoint     string
	APIKey       string
	OpenAIClient *openai.Client // proxy openai
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://dashscope.aliyuncs.com",
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:   APIKey,
			ErrResp:  &openai.ErrorResponse{},
		},
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/api/v1/services/embeddings/text-embedding/text-embedding"
	headers := []common.Header{
		{Key: "Authorization", Value: c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return EmbeddingResponse{}, err
	}

	var result EmbeddingResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	return result, err
}
