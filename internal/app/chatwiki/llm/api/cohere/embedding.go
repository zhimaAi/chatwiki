// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package cohere

type EmbeddingRequest struct {
	Texts          []string `json:"texts"`
	Model          string   `json:"model"`
	InputType      string   `json:"input_type"`
	EmbeddingTypes []string `json:"embedding_types"`
	Truncate       string   `json:"truncate,omitempty"`
}

type EmbeddingResponse struct {
	ResponseType string       `json:"response_type"`
	ID           string       `json:"id"`
	Embeddings   [][]float64  `json:"embeddings"`
	Texts        []string     `json:"texts"`
	Meta         ResponseMeta `json:"meta"`
}
