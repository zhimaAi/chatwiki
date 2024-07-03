// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"chatwiki/internal/app/chatwiki/llm/api/baai"
	"chatwiki/internal/app/chatwiki/llm/api/cohere"
	"chatwiki/internal/app/chatwiki/llm/api/jina"
	"chatwiki/internal/app/chatwiki/llm/api/xinference"
	"fmt"
	"sort"

	"github.com/zhimaAi/go_tools/msql"
)

type ZhimaRerankReq struct {
	Enable   bool
	Query    string        `json:"query" toml:"query"`
	Passages []string      `json:"passages" toml:"passages"`
	Data     []msql.Params `json:"data"`
	TopK     int           `json:"top_k" toml:"top_k"`
}
type RerankData struct {
	Index          int     `json:"index"`
	Text           string  `json:"text"`
	RelevanceScore float64 `json:"relevance_score"`
}
type ZhimaRerankResp struct {
	Data []*RerankData `json:"data"`
}

func (a *Adaptor) CreateRerank(params *ZhimaRerankReq) ([]msql.Params, error) {
	data := make([]RerankData, 0)
	switch a.meta.Corp {
	case "baai":
		client := baai.NewClient(a.meta.EndPoint, a.meta.Model, a.meta.APIKey)
		req := &baai.CreateRerankReq{
			Model:    a.meta.Model,
			Query:    params.Query,
			Passages: params.Passages,
			TopK:     params.TopK,
		}
		res, err := client.CreateRerank(req)
		if err != nil || len(res.Results) <= 0 {
			return nil, err
		}
		for _, item := range res.Results {
			data = append(data, RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "cohere":
		client := cohere.NewClient(a.meta.APIKey)
		req := cohere.ReRankRequest{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.ReRank(req)
		if err != nil || len(res.Results) <= 0 {
			return nil, err
		}
		for _, item := range res.Results {
			data = append(data, RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "jina":
		client := jina.NewClient(a.meta.APIKey)
		req := jina.ReRankRequest{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.ReRank(req)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Results {
			data = append(data, RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "xinference":
		client := xinference.NewClient(a.meta.EndPoint, a.meta.APIVersion, a.meta.Model)
		req := &xinference.CreateRerankReq{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.CreateRerank(req)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Results {
			data = append(data, RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	}
	return rerankData(params, data), nil
}

func rerankData(req *ZhimaRerankReq, rerankData []RerankData) []msql.Params {
	if len(rerankData) <= 0 {
		return req.Data
	}
	newData := make([]msql.Params, 0)
	sort.Slice(rerankData, func(i, j int) bool {
		return rerankData[i].RelevanceScore > rerankData[j].RelevanceScore
	})
	for _, item := range rerankData {
		// topN filter
		if req.TopK > 0 && len(newData) >= req.TopK {
			continue
		}
		if req.Data[item.Index] != nil {
			req.Data[item.Index]["relevance_score"] = fmt.Sprintf("%v", item.RelevanceScore)
			newData = append(newData, req.Data[item.Index])
		}
	}
	if len(newData) == 0 {
		return req.Data
	}
	return newData
}
