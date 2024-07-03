// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package hunyuan

import (
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentHunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type Client struct {
	SecretID  string
	SecretKey string
	Region    string
}

func NewClient(secretID, secretKey, region string) *Client {
	return &Client{
		SecretID:  secretID,
		SecretKey: secretKey,
		Region:    region,
	}
}

func (c *Client) CreateEmbeddings(req tencentHunyuan.GetEmbeddingRequest) (tencentHunyuan.GetEmbeddingResponseParams, error) {
	credential := common.NewCredential(c.SecretID, c.SecretKey)
	client, err := tencentHunyuan.NewClient(credential, c.Region, profile.NewClientProfile())
	if err != nil {
		return tencentHunyuan.GetEmbeddingResponseParams{}, err
	}
	response, err := client.GetEmbedding(&req)
	var tencentCloudSDKError *tencentErrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		return tencentHunyuan.GetEmbeddingResponseParams{}, err
	}
	if err != nil {
		return tencentHunyuan.GetEmbeddingResponseParams{}, err
	}
	if len(response.Response.Data) < 1 {
		return tencentHunyuan.GetEmbeddingResponseParams{}, errors.New("response no data")
	}
	return *response.Response, nil
}

func (c *Client) CreateChatCompletion(req tencentHunyuan.ChatCompletionsRequest) (tencentHunyuan.ChatCompletionsResponseParams, error) {

	credential := common.NewCredential(c.SecretID, c.SecretKey)
	client, err := tencentHunyuan.NewClient(credential, c.Region, profile.NewClientProfile())
	if err != nil {
		return tencentHunyuan.ChatCompletionsResponseParams{}, err
	}
	req.Stream = common.BoolPtr(false)
	response, err := client.ChatCompletions(&req)
	var tencentCloudSDKError *tencentErrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		return tencentHunyuan.ChatCompletionsResponseParams{}, err
	}
	if err != nil {
		return tencentHunyuan.ChatCompletionsResponseParams{}, err
	}
	return *response.Response, nil
}

func (c *Client) CreateChatCompletionStream(req tencentHunyuan.ChatCompletionsRequest) (*ChatCompletionStream, error) {
	credential := common.NewCredential(c.SecretID, c.SecretKey)
	client, err := tencentHunyuan.NewClient(credential, c.Region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	req.Stream = common.BoolPtr(true)
	response, err := client.ChatCompletions(&req)
	var tencentCloudSDKError *tencentErrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &ChatCompletionStream{events: response.Events, closed: false}, nil
}
