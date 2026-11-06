// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"maps"
)

func manualTransferOutput(state *define.ManualTransferState) common.SimpleFields {
	content, ok := state.Result()
	if !ok {
		return nil
	}
	return common.SimpleFields{
		`special.is_switch_manual`:  {Key: `special.is_switch_manual`, Typ: common.TypBoole, Vals: []common.Val{{Boole: &ok}}},
		`special.llm_reply_content`: {Key: `special.llm_reply_content`, Typ: common.TypString, Vals: []common.Val{{String: &content}}},
	}
}

func mergeManualTransferOutput(output common.SimpleFields, state *define.ManualTransferState) common.SimpleFields {
	manual := manualTransferOutput(state)
	if manual == nil {
		return output
	}
	merged := maps.Clone(output)
	if merged == nil {
		merged = common.SimpleFields{}
	}
	for key, field := range manual {
		merged[key] = field
	}
	return merged
}

// A transfer is a successful terminal result, not a pause or an exception.
func (flow *WorkFlow) applyManualTransfer() bool {
	if flow.params == nil || flow.params.ChatRequestParam == nil {
		return false
	}
	output := manualTransferOutput(flow.params.ManualTransfer)
	if output == nil {
		return false
	}
	flow.output = mergeManualTransferOutput(flow.output, flow.params.ManualTransfer)
	flow.isFinish = true
	flow.isStorage = false
	return true
}
