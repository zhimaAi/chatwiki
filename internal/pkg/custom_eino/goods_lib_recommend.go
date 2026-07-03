// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"
	"fmt"
	"strings"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/zhimaAi/go_tools/tool"
)

const GoodsLibRecommendToolName = "recommend_goods"

type GoodsLibRecommendTool struct {
	do func(query, searchType string, maxCount int) (string, error)
}

func BuildGoodsLibRecommendTool(do func(query, searchType string, maxCount int) (string, error)) einotool.BaseTool {
	return &GoodsLibRecommendTool{do: do}
}

func (t *GoodsLibRecommendTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: GoodsLibRecommendToolName,
		Desc: "Search enabled products in the goods library and recommend matching items for product-related requests such as buying, finding, comparing, or learning about products. Results are limited to the configured goods group scope. Use an empty query for broad recommendations within that scope.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     schema.String,
				Desc:     "Optional product keyword. It is matched fuzzily against goods_id and goods_name only. Leave empty when the request is a broad recommendation or the user does not provide a product ID/name keyword.",
				Required: false,
			},
			"search_type": {
				Type:     schema.String,
				Desc:     "Optional result field set. Use detail for full product fields, including images, description, QA, and custom_info. Use basic for core fields only. Omit for detail behavior.",
				Required: false,
			},
			"max_count": {
				Type:     schema.Integer,
				Desc:     "Optional maximum number of products to query. Use a positive integer. Omitted or non-positive values default to 9999.",
				Required: false,
			},
		}),
	}, nil
}

func (t *GoodsLibRecommendTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	var input struct {
		Query      string `json:"query"`
		SearchType string `json:"search_type"`
		MaxCount   int    `json:"max_count"`
	}
	if err := tool.JsonDecodeUseNumber(argumentsInJSON, &input); err != nil {
		return "", fmt.Errorf("parse %s arguments: %w", GoodsLibRecommendToolName, err)
	}
	if input.MaxCount <= 0 {
		input.MaxCount = 9999
	}
	return t.do(strings.TrimSpace(input.Query), strings.ToLower(strings.TrimSpace(input.SearchType)), input.MaxCount)
}
