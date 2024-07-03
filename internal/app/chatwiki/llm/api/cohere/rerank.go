// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package cohere

type ReRankRequest struct {
	Model           string   `json:"model"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	TopN            int      `json:"top_n,omitempty"`
	RankField       []string `json:"rank_field,omitempty"`
	ReturnDocuments bool     `json:"return_documents,omitempty"`
	MaxChunksPerDoc int      `json:"max_chunks_per_doc,omitempty"`
}

type ReRankResponse struct {
	ID      string   `json:"id"`
	Results []Result `json:"results"`
	Meta    Meta     `json:"meta"`
}

type Result struct {
	Index          int      `json:"index"`
	Document       Document `json:"document,omitempty"`
	RelevanceScore float64  `json:"relevance_score"`
}

type Document struct {
	Text string `json:"text"`
}

type Meta struct {
	APIVersion  APIVersion  `json:"apiVersion,omitempty"`
	BilledUnits BilledUnits `json:"billed_units,omitempty"`
	Tokens      Tokens      `json:"tokens,omitempty"`
}

type APIVersion struct {
	Version        string `json:"version,omitempty"`
	IsDeprecated   bool   `json:"is_deprecated,omitempty"`
	IsExperimental bool   `json:"is_experimental,omitempty"`
}

type BilledUnits struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	SearchUnits  int `json:"search_units,omitempty"`
}

type Tokens struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}
