// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package volcenginev2

type EmbeddingRequest struct {
	Input []string
}

type EmbeddingResponse struct {
	Data  []Embedding `json:"data"`
	Usage Usage       `json:"usage"`
}

type Embedding struct {
	Embedding []float32 `json:"embedding"`
	Object    string    `json:"object"`
	Index     int       `json:"index"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
