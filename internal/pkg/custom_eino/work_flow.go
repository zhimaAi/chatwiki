// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type WorkFlowTool struct {
	name   string
	desc   string
	params *schema.ParamsOneOf
	do     func(adaptor.FunctionToolCall) (string, error)
}

func BuildWorkFlowTool(name, desc string, params *schema.ParamsOneOf, do func(adaptor.FunctionToolCall) (string, error)) einotool.BaseTool {
	return &WorkFlowTool{name: name, desc: desc, params: params, do: do}
}

func (t *WorkFlowTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: t.name, Desc: t.desc, ParamsOneOf: t.params}, nil
}

func (t *WorkFlowTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	return t.do(adaptor.FunctionToolCall{Name: t.name, Arguments: argumentsInJSON})
}
