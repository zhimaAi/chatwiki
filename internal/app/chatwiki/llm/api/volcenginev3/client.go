// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package volcenginev3

import (
	"chatwiki/internal/app/chatwiki/llm/api/openai"
	"chatwiki/internal/app/chatwiki/llm/common"
)

type Client struct {
	Host   string
	Model  string
	AK     string
	SK     string
	Region string
}

func NewClient(Host, Model, AK, SK, Region string) *Client {
	return &Client{
		Host:   Host,
		Model:  Model,
		AK:     AK,
		SK:     SK,
		Region: Region,
	}
}

func getAccessToken(Region, Model, AK, SK string) (string, error) {
	tokenManager := common.GetTokenManagerInstance()
	return tokenManager.GetVolcengineAccessToken("https://open.volcengineapi.com", Region, Model, AK, SK)
}

func (c *Client) CreateChatCompletion(req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	accessToken, err := getAccessToken(c.Region, c.Model, c.AK, c.SK)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	OpenAIClient := openai.NewClient(c.Host, accessToken, &openai.ErrorResponse{})
	return OpenAIClient.CreateChatCompletion(req)
}

func (c *Client) CreateChatCompletionStream(req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	accessToken, err := getAccessToken(c.Region, c.Model, c.AK, c.SK)
	if err != nil {
		return nil, err
	}
	OpenAIClient := openai.NewClient(c.Host, accessToken, &openai.ErrorResponse{})
	return OpenAIClient.CreateChatCompletionStream(req)
}
