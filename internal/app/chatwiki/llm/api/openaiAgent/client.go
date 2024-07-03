// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package openai_agent

import (
	"chatwiki/internal/app/chatwiki/llm/api/openai"
	"chatwiki/internal/app/chatwiki/llm/common"
)

type Client struct {
	APIKey       string
	EndPoint     string
	ErrResp      common.ErrorResponseInterface
	OpenAIClient *openai.Client
}

func NewClient(apiEndpoint, APIKey, apiVersion string) *Client {
	return &Client{
		EndPoint: apiEndpoint,
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: apiEndpoint + "/" + apiVersion,
			APIKey:   APIKey,
			ErrResp:  &openai.ErrorResponse{},
		},
	}
}
