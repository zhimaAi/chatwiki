// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
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

func EmbeddingNewQAVector(libraryId, adminUserId, qaIndexType int) {
	//async task:convert vector
	var (
		vectorData = make(map[string][]int)
		m          = msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	)
	answerIds, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Where(`type`, `in`, cast.ToString(define.VectorTypeAnswer)).
		ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(answerIds) > 0 {
		m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, cast.ToString(libraryId)).
			Where(`id`, `in`, strings.Join(answerIds, `,`)).
			Delete()
	}
	if qaIndexType == define.QAIndexTypeQuestionAndAnswer {
		if len(answerIds) <= 0 {
			// 新增
			page := 1
			size := 200
			for {
				m := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
					Where(`library_id`, cast.ToString(libraryId)).
					Field(`id,file_id,answer`)
				list, _, err := m.Order(`id desc`).Paginate(page, size)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				if len(list) == 0 {
					break
				}
				for _, item := range list {
					vectorID, err := SaveVector(int64(adminUserId), cast.ToInt64(libraryId), 0, cast.ToInt64(item[`id`]), cast.ToString(define.VectorTypeAnswer), item[`answer`])
					if err != nil {
						logs.Error(err.Error())
						return
					}
					vectorData[item[`file_id`]] = append(vectorData[item[`file_id`]], cast.ToInt(vectorID))
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
		ColumnArr(`data_id`)
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

func AddDefaultLibrary(token, robotName, libraryIds, robotKey string, adminUserId int) (int, error) {
	libraryData := map[string]string{
		`library_name`:                       robotName + `默认知识库`,
		`chunk_type`:                         cast.ToString(define.ChunkTypeNormal),
		`normal_chunk_default_separators_no`: `12,11`,
		`normal_chunk_default_chunk_size`:    `512`,
		`type`:                               cast.ToString(define.QALibraryType),
	}

	modelConfig := GetDefaultEmbeddingConfig(adminUserId)
	if len(modelConfig) > 0 {
		libraryData[`use_model`] = modelConfig[`vector_model`]
		libraryData[`model_config_id`] = modelConfig[`id`]
	}

	var (
		res lib_web.Response
		err error
	)
	req := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/createLibrary`, define.Config.WebService[`port`])).Header(`token`, token)
	for key, item := range libraryData {
		req.Param(key, item)
	}
	if err = req.ToJSON(&res); err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		return 0, fmt.Errorf(`%s`, cast.ToString(res.Msg))
	}
	id := cast.ToInt(cast.ToStringMap(res.Data)[`id`])
	if len(libraryIds) > 0 {
		libraryIds = fmt.Sprintf(`%d,%s`, id, libraryIds)
	} else {
		libraryIds = cast.ToString(id)
	}
	// 更新机器人的知识库关联信息
	if _, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Update(msql.Datas{
		`library_ids`:        libraryIds,
		`default_library_id`: id,
		`update_time`:        tool.Time2Int(),
	}); err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	lib_redis.DelCacheData(define.Redis, &RobotCacheBuildHandler{RobotKey: robotKey})
	return id, err
}
