// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

const ManualTransferToolName = "zmxkf_internal_transfer_to_human"

type ManualTransferTool struct {
	do func() (string, error)
}

func BuildManualTransferTool(do func() (string, error)) einotool.BaseTool {
	return &ManualTransferTool{do: do}
}

func (t *ManualTransferTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ManualTransferToolName,
		Desc:        "Transfer the current conversation to human customer service when the user needs human assistance. Routing and the transfer message are configured by the administrator. If the tool reports failure, explain the reason to the user and continue helping; do not claim a successful transfer.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}, nil
}

func (t *ManualTransferTool) InvokableRun(_ context.Context, _ string, _ ...einotool.Option) (string, error) {
	return t.do()
}
