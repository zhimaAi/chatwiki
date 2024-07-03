// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"sort"

	"github.com/zhimaAi/go_tools/msql"
)

type DataSource struct {
	List  []msql.Params
	Key   string
	Fixed int
}

type RRF struct {
	dataSources []DataSource
}

func (r *RRF) Add(ds DataSource) *RRF {
	if len(ds.List) == 0 {
		return r
	}
	if r.dataSources == nil {
		r.dataSources = make([]DataSource, 0)
	}
	r.dataSources = append(r.dataSources, ds)
	return r
}

type RrfScoreItem struct {
	Key   string
	Score float64
}

type RrfScoreList []RrfScoreItem

func (r RrfScoreList) Len() int {
	return len(r)
}
func (r RrfScoreList) Less(i, j int) bool {
	return r[i].Score > r[j].Score
}
func (r RrfScoreList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *RRF) Sort() []msql.Params {
	//sum score
	rrfMap := make(map[string]float64)
	for i := range r.dataSources {
		for j := range r.dataSources[i].List {
			score := 1 / float64(j+1+r.dataSources[i].Fixed)
			rrfMap[r.dataSources[i].List[j][r.dataSources[i].Key]] += score
		}
	}
	//sort by score
	itemList := make(RrfScoreList, 0)
	for key := range rrfMap {
		itemList = append(itemList, RrfScoreItem{Key: key, Score: rrfMap[key]})
	}
	sort.Sort(itemList)
	//build all data
	alllistMap := make(map[string]msql.Params)
	for i := range r.dataSources {
		for j := range r.dataSources[i].List {
			if _, ok := alllistMap[r.dataSources[i].List[j][r.dataSources[i].Key]]; ok {
				continue //duplication
			}
			alllistMap[r.dataSources[i].List[j][r.dataSources[i].Key]] = r.dataSources[i].List[j]
		}
	}
	//return
	result := make([]msql.Params, 0)
	for i := range itemList {
		if one, ok := alllistMap[itemList[i].Key]; ok {
			result = append(result, one)
		}
	}
	return result
}
