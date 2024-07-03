// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"chatwiki/internal/app/chatwiki/llm/api/baai"
)

type ZhimaSimilarityRequest struct {
	Model string      `json:"model"`
	Query []float64   `json:"query"`
	Input [][]float64 `json:"input"`
}

type ZhimaSimilarityResponse struct {
	Result []float64 `json:"result"`
}

func (a *Adaptor) CreateSimilarity(req ZhimaSimilarityRequest) (ZhimaSimilarityResponse, error) {
	switch a.meta.Corp {
	case "baai":
		client := baai.NewClient(a.meta.EndPoint, a.meta.Model, a.meta.APIKey)
		req := baai.SimilarityRequest{
			Model: req.Model,
			Query: req.Query,
			Input: req.Input,
		}
		res, err := client.ComputeSimilarity(req)
		if err != nil {
			return ZhimaSimilarityResponse{}, err
		}
		return ZhimaSimilarityResponse{
			Result: res.Data,
		}, nil
	}
	return ZhimaSimilarityResponse{}, nil
}
