// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ollama

import "time"

type Duration struct {
	time.Duration
}

// EmbeddingRequest is the request passed to [Client.Embeddings].
type EmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`

	// Prompt is the textual prompt to embed.
	Prompt string `json:"prompt"`

	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options,omitempty"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}
