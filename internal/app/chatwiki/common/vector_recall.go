// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetVectorDims(embedding string) int {
	result := make([]float64, 0)
	err := tool.JsonDecodeUseNumber(embedding, &result)
	if err != nil { // Exception handling uses compatibility logic
		logs.Error(err.Error())
		return len(strings.Split(embedding, `,`))
	} else {
		return len(result)
	}
}

func VectorRecall(libraryIds, embedding string, size int) ([]msql.Params, error) {
	fetchSize := 10 * size // Since one segment corresponds to multiple vector indexes, ten times recall is used
	vectorDims := GetVectorDims(embedding)
	var embeddingKey = `embedding`
	var dimsSql = fmt.Sprintf(`vector_dims(embedding)=%d`, vectorDims)
	if vectorDims == 2000 {
		embeddingKey = `embedding2000` // Fixed 2000-dimensional vector document
		dimsSql = ``                   // Fixed dimensions do not require additional SQL statements
	}
	indexModel := msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	tempList, err := indexModel.Where(`delete_time`, `0`).
		Where(`status`, cast.ToString(define.VectorStatusConverted)).
		Where(`library_id`, `in`, libraryIds).Where(dimsSql).
		Field(`data_id`).Field(fmt.Sprintf(`1-(%s<=>'%s') similarity`, embeddingKey, embedding)).
		Order(fmt.Sprintf(`%s<=>'%s'`, embeddingKey, embedding)).Limit(fetchSize).Select()
	if err != nil {
		return nil, fmt.Errorf(`sql:%s,err:%s`, indexModel.GetLastSql(), err.Error())
	}
	dataIds, similarityMap := make([]string, 0), make(map[string]string)
	for _, one := range tempList {
		if _, ok := similarityMap[one[`data_id`]]; ok {
			continue // Skip segments that have already appeared
		}
		similarityMap[one[`data_id`]] = one[`similarity`]
		dataIds = append(dataIds, one[`data_id`])
		if len(dataIds) >= size {
			break // Segment count is sufficient
		}
	}
	if len(dataIds) == 0 {
		return nil, nil
	}
	// Get segment information
	dataModel := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	subList, err := dataModel.Where(`delete_time`, `0`).
		Where(`id`, `in`, strings.Join(dataIds, `,`)).Select()
	if err != nil {
		return nil, fmt.Errorf(`sql:%s,err:%s`, dataModel.GetLastSql(), err.Error())
	}
	if len(subList) == 0 {
		return nil, nil
	}
	// Add similarity field to segments
	for i, one := range subList {
		subList[i][`similarity`] = similarityMap[one[`id`]]
	}
	// Sort by similarity and truncate return
	list := make(define.SimilarityResult, 0)
	list = append(list, subList...)
	sort.Sort(list)
	if len(list) > size {
		list = list[:size]
	}
	return list, nil
}
