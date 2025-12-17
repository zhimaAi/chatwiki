// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import "github.com/zhimaAi/go_tools/msql"

func SliceUnifyUnique[Item any, Key comparable](list []Item, unique func(Item) Key) []Item {
	newList := make([]Item, 0)
	maps := map[Key]struct{}{}
	for _, item := range list {
		uniqueKey := unique(item)
		if _, ok := maps[uniqueKey]; ok {
			continue
		}
		maps[uniqueKey] = struct{}{}
		newList = append(newList, item)
	}
	return newList
}

func SliceMsqlParamsUnique(list []msql.Params, key string) []msql.Params {
	return SliceUnifyUnique(list, func(params msql.Params) string { return params[key] })
}
