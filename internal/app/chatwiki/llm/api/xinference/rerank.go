// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package xinference

type (
	CreateRerankReq struct {
		Model           string   `json:"model"`
		Query           string   `json:"query"`
		Documents       []string `json:"documents"`
		TopN            int      `json:"top_n,omitempty"`
		ReturnDocuments bool     `json:"return_documents,omitempty"`
		MaxChunksPerDoc int      `json:"max_chunks_per_doc,omitempty"`
	}
	CreateRerankRes struct {
		Results []Data `json:"results"`
		Id      string `json:"id"`
	}
	Data struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
		Document       string  `json:"document"`
	}
)
