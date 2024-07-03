// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package azure

type EmbeddingRequest struct {
	Input []string `json:"input"`
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
