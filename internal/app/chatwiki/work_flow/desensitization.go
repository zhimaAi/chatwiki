// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func NodesConfDesensitize(nodes []msql.Params) []msql.Params {
	for i, node := range nodes {
		nodes[i] = nodeConfDesensitize(node)
	}
	return nodes
}

func nodeConfDesensitize(node msql.Params) msql.Params {
	nodeType := cast.ToInt(node[`node_type`])
	//nodeParams := NodeParams{}
	//_ = tool.JsonDecodeUseNumber(node[`node_params`], &nodeParams)
	nodeParams := DisposeNodeParams(nodeType, node[`node_params`])
	nodeInfoStr := node[`node_info_json`]
	switch nodeType {
	case NodeTypeStart:
		nodeParams.Start, nodeInfoStr = startNodeDesensitize(nodeParams.Start, nodeInfoStr)
	case NodeTypePlugin:
		nodeParams.Plugin, nodeInfoStr = pluginNodeDesensitize(nodeParams.Plugin, nodeInfoStr)
	}
	node[`node_params`] = tool.JsonEncodeNoError(nodeParams)
	node[`node_info_json`] = nodeInfoStr
	return node
}

func startNodeDesensitize(start StartNodeParams, nodeInfoStr string) (StartNodeParams, string) {
	replace := make(map[string]any) //修正node_info_json.dataRaw数据
	for i, trigger := range start.TriggerList {
		switch trigger.TriggerType {
		case TriggerTypeOfficial:
			replace[fmt.Sprintf(`app_ids#%s`, trigger.TriggerOfficialConfig.AppIds)] = ``
			start.TriggerList[i].TriggerOfficialConfig.AppIds = ``
		}
	}
	return start, AmendNodeinfojson(nodeInfoStr, replace)
}

func pluginNodeDesensitize(plugin PluginNodeParams, nodeInfoStr string) (PluginNodeParams, string) {
	//todo ...
	return plugin, nodeInfoStr
}
