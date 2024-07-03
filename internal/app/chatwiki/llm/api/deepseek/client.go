// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package deepseek

import "chatwiki/internal/app/chatwiki/llm/api/openai"

type Client struct {
	EndPoint     string
	APIKey       string
	OpenAIClient *openai.Client
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://api.deepseek.com",
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: "https://api.deepseek.com",
			APIKey:   APIKey,
			ErrResp:  &ErrorResponse{},
		},
	}
}
