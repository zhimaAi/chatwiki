// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package jina

type EmbeddingRequest struct {
	Input        []string `json:"input"`
	Model        string   `json:"model"`
	EncodingType string   `json:"encoding_type"`
}

type EmbeddingResponse struct {
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  Usage  `json:"usage"`
	Data   []Data `json:"data"`
}

type Usage struct {
	TotalTokens  int `json:"total_tokens"`
	PromptTokens int `json:"prompt_tokens"`
}

type Data struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}
