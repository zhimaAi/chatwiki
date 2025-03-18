// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"strings"
)

func EmbeddingNewVector(libraryId, adminUserId int) {
	//async task:convert vector
	var (
		vectorData = make(map[string][]int)
		m          = msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	)
	// get file data
	page := 1
	size := 200
	for {
		m := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, cast.ToString(libraryId)).
			Field(`id,file_id`)
		list, _, err := m.Order(`id desc`).Paginate(page, size)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		if len(list) == 0 {
			break
		}
		var ids, fileIds = []string{}, []string{}
		for _, item := range list {
			vectorData[item[`file_id`]] = append(vectorData[item[`file_id`]], cast.ToInt(item[`id`]))
			ids = append(ids, item[`id`])
			fileIds = append(fileIds, item[`file_id`])
		}
		// finished
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, `in`, strings.Join(fileIds, `,`)).Update(msql.Datas{
			`status`:      define.FileStatusLearning,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			return
		}
		if _, err = m.Where(`id`, `in`, strings.Join(ids, `,`)).Update(msql.Datas{
			`status`:      define.VectorStatusInitial,
			`update_time`: tool.Time2Int(),
		}); err != nil {
			logs.Error(err.Error())
		}
		page++
	}
	for fileId, ids := range vectorData {
		for _, id := range ids {
			if message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId}); err != nil {
				logs.Error(err.Error())
			} else if err := AddJobs(define.ConvertVectorTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}
}

func AddFileDataIndex(libraryId, adminUserId int) error {
	dataIds, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Field(`id,file_id,content`).Select()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	indexIds, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Field(`data_id`).GetFields()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	insertIds := []msql.Params{}
	for _, item := range dataIds {
		if !tool.InArrayString(item[`id`], indexIds) {
			insertIds = append(insertIds, item)
		}
	}
	for _, item := range insertIds {
		_, err := SaveVector(int64(adminUserId), int64(libraryId), cast.ToInt64(item[`file_id`]), cast.ToInt64(item[`id`]), cast.ToString(define.VectorTypeParagraph), item[`content`])
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return err
}
