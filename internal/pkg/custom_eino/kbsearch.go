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

const KbsearchToolName = "search_knowledge"

type KbsearchTool struct {
	do func(query string) (string, error)
}

func BuildKbsearchTool(do func(query string) (string, error)) einotool.BaseTool {
	return &KbsearchTool{do: do}
}

func (t *KbsearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: KbsearchToolName,
		Desc: "Search the knowledge base for passages relevant to a user question.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     schema.String,
				Desc:     "Search query or user question",
				Required: true,
			},
		}),
	}, nil
}

func (t *KbsearchTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	var input struct {
		Query string `json:"query"`
	}
	if err := tool.JsonDecodeUseNumber(argumentsInJSON, &input); err != nil {
		return "", fmt.Errorf("parse %s arguments: %w", KbsearchToolName, err)
	}
	if strings.TrimSpace(input.Query) == `` {
		return `query can not be empty`, nil
	}
	return t.do(input.Query)
}
