// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ali

type EmbeddingRequest struct {
	Input      Texts    `json:"input"`
	Model      string   `json:"model"`
	Parameters TextType `json:"parameters"`
}
type Texts struct {
	Texts []string `json:"texts"`
}
type TextType struct {
	TextType string `json:"text_type"`
}

type EmbeddingResponse struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestId string `json:"request_id"`
}

type Output struct {
	Embeddings []Embedding `json:"embeddings"`
}

type Embedding struct {
	TextIndex int       `json:"text_index"`
	Embedding []float64 `json:"embedding"`
}

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}
