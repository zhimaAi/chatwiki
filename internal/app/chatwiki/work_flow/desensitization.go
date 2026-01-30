// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func NodesConfDesensitize(adminUserId int, nodes []msql.Params, lang string) []msql.Params {
	for i, node := range nodes {
		nodes[i] = nodeConfDesensitize(adminUserId, node, lang)
	}
	return nodes
}

func nodeConfDesensitize(adminUserId int, node msql.Params, lang string) msql.Params {
	nodeType := cast.ToInt(node[`node_type`])
	//nodeParams := NodeParams{}
	//_ = tool.JsonDecodeUseNumber(node[`node_params`], &nodeParams)
	nodeParams := DisposeNodeParams(nodeType, node[`node_params`], lang)
	nodeInfoStr := node[`node_info_json`]
	switch nodeType {
	case NodeTypeStart:
		nodeParams.Start, nodeInfoStr = startNodeDesensitize(nodeParams.Start, nodeInfoStr)
	case NodeTypePlugin:
		nodeParams.Plugin, nodeInfoStr = PluginNodeDesensitize(adminUserId, nodeParams.Plugin, nodeInfoStr)
	}
	node[`node_params`] = tool.JsonEncodeNoError(nodeParams)
	node[`node_info_json`] = nodeInfoStr
	return node
}

func startNodeDesensitize(start StartNodeParams, nodeInfoStr string) (StartNodeParams, string) {
	replace := make(map[string]any) // Correct node_info_json.dataRaw data
	for i, trigger := range start.TriggerList {
		switch trigger.TriggerType {
		case TriggerTypeOfficial:
			replace[fmt.Sprintf(`app_ids#%s`, trigger.TriggerOfficialConfig.AppIds)] = ``
			start.TriggerList[i].TriggerOfficialConfig.AppIds = ``
		}
	}
	return start, AmendNodeinfojson(nodeInfoStr, replace)
}

func PluginNodeDesensitize(adminUserId int, plugin PluginNodeParams, nodeInfoStr string) (PluginNodeParams, string) {
	// First check if there are parameters
	arguments, exists := plugin.Params[`arguments`]
	if !exists || arguments == nil {
		//logs.Debug("Plugin has no parameters, no replacement needed")
		return plugin, nodeInfoStr
	}
	pluginArguments := arguments.(map[string]any)
	// Check if the plugin has sensitive information and perform replacement
	u := define.Config.Plugin[`endpoint`] + "/manage/plugin/local-plugins/run"
	resp := &lib_web.Response{}
	request := curl.Post(u).Header(`admin_user_id`, cast.ToString(adminUserId))
	request.Param("name", plugin.Name)
	request.Param("action", "default/export-filtering")
	params, err := json.Marshal(plugin.Params)
	if err != nil {
		//logs.Debug("Failed to convert node params to json for filtering: %v", err.Error())
		plugin.NoAuthFilter = true
		return plugin, nodeInfoStr
	}

	err = request.Param("params", string(params)).ToJSON(resp)
	if err != nil {
		return plugin, nodeInfoStr
	}
	if resp.Res != 0 {
		//logs.Debug("Replacement filtering failed resp %v", resp.Res)
		plugin.NoAuthFilter = true
		return plugin, nodeInfoStr
	}
	// verify schema
	filterFieldData, ok := resp.Data.(map[string]any)
	if !ok {
		//logs.Debug("No data filtering data")
		plugin.NoAuthFilter = false
		return plugin, nodeInfoStr
	}
	replace := make(map[string]any) // Correct node_info_json.dataRaw data

	for key, toValue := range filterFieldData {
		value, isOk := pluginArguments[key]
		if !isOk {
			// No configuration, do not replace
			continue
		}
		if value == toValue {
			// No replacement needed
			continue
		}
		replace[fmt.Sprintf(`%s#%v`, key, value)] = toValue
		pluginArguments[key] = toValue
	}
	rendering, exists := plugin.Params[`rendering`]
	if exists && rendering != nil {
		// If there are rendering fields, replace them
		pluginRendering := rendering.([]any)
		for i, renderObject := range pluginRendering {
			renderItem := renderObject.(map[string]any)
			key := cast.ToString(renderItem[`key`])
			toValue, isOk := filterFieldData[key]
			if isOk {
				renderItem[`value`] = toValue
				pluginRendering[i] = renderItem
			}
		}
		plugin.Params[`rendering`] = pluginRendering
	}
	plugin.Params[`config_name`] = ``
	plugin.Params[`arguments`] = pluginArguments
	plugin.NoAuthFilter = false
	return plugin, AmendNodeinfojson(nodeInfoStr, replace)
}
