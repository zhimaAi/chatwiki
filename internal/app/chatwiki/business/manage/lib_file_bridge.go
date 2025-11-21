// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
)

type BridgeGetLibFileListReq struct {
	LibraryId string `form:"library_id"`
	Status    string `form:"status"`
	Page      string `form:"page"`
	Size      string `form:"size"`
	SortField string `form:"sort_field"`
	SortType  string `form:"sort_type"`
	GroupId   string `form:"group_id"`
	FileName  string `form:"file_name"`
}

func BridgeGetLibFileList(adminUserId, loginUserId int, lang string, req *BridgeGetLibFileListReq) (map[string]any, int, error) {
	libraryId := cast.ToInt(req.LibraryId)
	if libraryId <= 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	wheres := [][]string{
		{`admin_user_id`, cast.ToString(adminUserId)}, {`library_id`, cast.ToString(libraryId)}, {`delete_time`, `0`},
	}
	status := cast.ToString(req.Status)
	page := max(1, cast.ToInt(req.Page))
	size := max(1, cast.ToInt(req.Size))
	sortField := cast.ToString(req.SortField)
	sortType := cast.ToString(req.SortType)
	// groupId := cast.ToInt(c.Query(`group_id`))
	// 全部时给一个默认值
	groupId := cast.ToInt(req.GroupId)
	m := msql.Model(`chat_ai_library_file`, define.Postgres).
		Alias(`f`).
		Join(`chat_ai_library_file_data d`, `f.id=d.file_id`, `left`).
		Where(`f.admin_user_id`, cast.ToString(adminUserId)).
		Where(`f.library_id`, cast.ToString(libraryId)).
		Where(`f.delete_time`, `0`).
		Group(`f.id`).
		Field(`f.*, count(d.id) as paragraph_count`).
		Field(`count(case when d.graph_status = 3 then 1 else null end) as graph_err_count`).
		Field(`
			COALESCE(
    			(SELECT graph_err_msg FROM chat_ai_library_file_data WHERE file_id = f.id AND graph_err_msg <> '' LIMIT 1),
    			'no error'
  			) AS graph_err_msg
		`).
		Field(`COALESCE((SELECT SUM(yesterday_hits) FROM chat_ai_library_file_data WHERE file_id = f.id), 0) as yesterday_hits`).
		Field(`COALESCE((SELECT SUM(today_hits) FROM chat_ai_library_file_data WHERE file_id = f.id), 0) as today_hits`).
		Field(`COALESCE((SELECT SUM(total_hits) FROM chat_ai_library_file_data WHERE file_id = f.id), 0) as total_hits`)
	fileName := strings.TrimSpace(req.FileName)
	if len(fileName) > 0 {
		m.Where(`file_name`, `like`, fileName)
		wheres = append(wheres, []string{`file_name`, `like`, fileName})
	}
	if status != "" {
		m.Where(`f.status`, `in`, status)
		// wheres = append(wheres, []string{`status`,`in`, status})
	}
	if groupId >= 0 {
		m.Where(`f.group_id`, cast.ToString(groupId))
		wheres = append(wheres, []string{`f.group_id`, cast.ToString(groupId)})
	}
	sortFields := []string{`yesterday_hits`, `today_hits`, `total_hits`}
	if tool.InArray(sortField, sortFields) {
		m.Order(sortField + ` ` + sortType)
	} else {
		m.Order(`id desc`)
	}
	list, total, err := m.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	countData, err := GetLibFileCount(wheres)
	var graphEntityCountRes *neo4j.EagerResult
	var idList []string
	for _, item := range list {
		idList = append(idList, cast.ToString(item[`id`]))
	}
	if len(idList) > 0 && common.GetNeo4jStatus(adminUserId) {
		graphEntityCountRes, err = common.NewGraphDB(adminUserId).GetEntityCount(idList)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	libFileIds := make([]string, 0)
	for _, item := range list {
		libFileIds = append(libFileIds, item[`id`])
		item[`graph_entity_count`] = `0`
		if graphEntityCountRes == nil {
			continue
		}
		for _, record := range graphEntityCountRes.Records {
			fileId, exists1 := record.Get("file_id")
			count, exists2 := record.Get("count")
			if exists1 && exists2 && fileId == cast.ToInt64(item[`id`]) {
				item[`graph_entity_count`] = cast.ToString(count)
			}
		}
	}
	data := map[string]any{`info`: info, `list`: list, `count_data`: countData, `total`: total, `page`: page, `size`: size}
	return data, 0, nil
}

type BridgeAddLibraryFileReq struct {
	LibraryId             string `form:"library_id"`
	LibraryKey            string `form:"library_key"`
	DocType               string `form:"doc_type"`
	Urls                  string `form:"urls"`
	FileName              string `form:"file_name"`
	Content               string `form:"content"`
	Title                 string `form:"title"`
	IsQaDoc               string `form:"is_qa_doc"`
	QaIndexType           string `form:"qa_index_type"`
	DocAutoRenewFrequency string `form:"doc_auto_renew_frequency"`
	DocAutoRenewMinute    string `form:"doc_auto_renew_minute"`
	AnswerLable           string `form:"answer_lable"`
	AnswerColumn          string `form:"answer_column"`
	QuestionLable         string `form:"question_lable"`
	QuestionColumn        string `form:"question_column"`
	SimilarColumn         string `form:"similar_column"`
	SimilarLabel          string `form:"similar_label"`
	PdfParseType          string `form:"pdf_parse_type"`
	GroupId               string `form:"group_id"`
}

func BridgeAddLibraryFile(adminUserId, loginUserId int, lang string, req *BridgeAddLibraryFileReq, chunkParam *define.ChunkParam, c *gin.Context) (map[string]any, int, error) {
	libraryId := cast.ToInt(req.LibraryId)
	if libraryId <= 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	if chunkParam != nil && cast.ToInt(chunkParam.SetChunk) != 0 {
		err = ValidateChunkParam(adminUserId, chunkParam, info[`type`], lang)
		if err != nil {
			return nil, -1, err
		}
	}
	//common save
	fileIds, err := addLibFile(c, adminUserId, libraryId, cast.ToInt(info[`type`]), chunkParam, req)
	if err != nil {
		return nil, -1, err
	}
	return map[string]any{`file_ids`: fileIds}, 0, nil
}
