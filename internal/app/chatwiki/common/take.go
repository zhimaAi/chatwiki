// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func TakeParamsFromMap(data any, keys ...string) map[string]any {
	eData, err := tool.JsonEncode(data)
	if err != nil {
		logs.Error(err.Error())
		return map[string]any{}
	}
	mData := make(map[string]any)
	tData := make(map[string]any)
	err = tool.JsonDecode(eData, &mData)
	if err != nil {
		logs.Error(err.Error())
		return map[string]any{}
	}
	for key, val := range mData {
		if tool.InArray(key, keys) {
			tData[key] = val
		}
	}
	return tData
}

func TakeResultCombine(data ...map[string]any) map[string]any {
	if len(data) == 0 {
		return map[string]any{}
	}
	if len(data) == 1 {
		return data[0]
	}
	cData := make(map[string]any)
	for _, dataVal := range data {
		for k, v := range dataVal {
			cData[k] = v
		}
	}
	return cData
}
