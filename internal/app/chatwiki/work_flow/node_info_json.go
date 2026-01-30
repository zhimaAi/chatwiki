// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func AmendNodeinfojson(nodeInfoStr string, replace map[string]any) string {
	if len(replace) == 0 {
		return nodeInfoStr
	}
	nodeInfoJson := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(nodeInfoStr, &nodeInfoJson); err != nil {
		logs.Error(err.Error())
	}
	if _, ok := nodeInfoJson[`dataRaw`]; !ok {
		return nodeInfoStr
	}
	var dataRaw any
	if err := tool.JsonDecodeUseNumber(cast.ToString(nodeInfoJson[`dataRaw`]), &dataRaw); err != nil {
		logs.Error(err.Error())
	}
	dataRaw = AmendDataRaw(dataRaw, replace)
	nodeInfoJson[`dataRaw`] = tool.JsonEncodeNoError(dataRaw)
	return tool.JsonEncodeNoError(nodeInfoJson)
}

func AmendDataRaw(dataRaw any, replace map[string]any, keys ...string) any {
	switch realData := dataRaw.(type) {
	case []any:
		newData := make([]interface{}, len(realData))
		for idx, val := range realData {
			newData[idx] = AmendDataRaw(val, replace)
		}
		return newData
	case map[string]any:
		newData := make(map[string]interface{})
		for key, val := range realData {
			newData[key] = AmendDataRaw(val, replace, key)
		}
		return newData
	default:
		if len(keys) > 0 {
			if newVal, ok := replace[fmt.Sprintf(`%s#%v`, keys[0], realData)]; ok {
				return newVal //replace with new value
			}
			if newVal, ok := replace[keys[0]]; ok {
				return newVal //replace with new value
			}
		}
		return realData
	}
}
