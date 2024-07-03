// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package baai

import (
	"chatwiki/internal/app/chatwiki/llm/common"
	"errors"
	"io"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type EmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
type EmbeddingData struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}
type EmbeddingResponse struct {
	Data  []EmbeddingData `json:"data"`
	Usage EmbeddingUsage  `json:"usage"`
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/v1/embeddings"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}

	jsonData, _ := tool.JsonEncode(req)
	logs.Info("baai_CreateEmbeddings,req:%+v", jsonData)
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
		return EmbeddingResponse{}, errors.New("baai response no embedding result")
	}

	return result, err
}

type SimilarityRequest struct {
	Query []float64   `json:"query"`
	Input [][]float64 `json:"input"`
	Model string      `json:"model"`
}

type SimilarityResponse struct {
	Data []float64 `json:"data"`
}

func (c *Client) ComputeSimilarity(req SimilarityRequest) (SimilarityResponse, error) {

	url := c.EndPoint + "/v1/similarity"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}

	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return SimilarityResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(responseRaw.Body)

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return SimilarityResponse{}, err
	}

	var result SimilarityResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return SimilarityResponse{}, err
	}
	if len(result.Data) <= 0 {
		return SimilarityResponse{}, errors.New("baai response no similarity result")
	}

	return result, err
}
