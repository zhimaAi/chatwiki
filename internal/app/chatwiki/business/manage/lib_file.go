// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
)

func GetLibFileList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := cast.ToInt(c.Query(`library_id`))
	if libraryId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	wheres := [][]string{
		{`admin_user_id`, cast.ToString(userId)}, {`library_id`, cast.ToString(libraryId)},
	}
	status := cast.ToString(c.Query(`status`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	m := msql.Model(`chat_ai_library_file`, define.Postgres).
		Alias(`f`).
		Join(`chat_ai_library_file_data d`, `f.id=d.file_id`, `left`).
		Where(`f.admin_user_id`, cast.ToString(userId)).
		Where(`f.library_id`, cast.ToString(libraryId)).
		Group(`f.id`).
		Field(`f.*, count(d.id) as paragraph_count`).
		Field(`count(case when d.graph_status = 3 then 1 else null end) as graph_err_count`).
		Field(`
			COALESCE(
    			(SELECT graph_err_msg FROM chat_ai_library_file_data WHERE file_id = f.id AND graph_err_msg <> '' LIMIT 1),
    			'no error'
  			) AS graph_err_msg
		`)
	fileName := strings.TrimSpace(c.Query(`file_name`))
	if len(fileName) > 0 {
		m.Where(`file_name`, `like`, fileName)
		wheres = append(wheres, []string{`file_name`, `like`, fileName})
	}
	if status != "" {
		m.Where(`f.status`, `in`, status)
		// wheres = append(wheres, []string{`status`,`in`, status})
	}
	list, total, err := m.Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	countData, err := getLibFileCount(wheres)
	var graphEntityCountRes *neo4j.EagerResult
	var idList []string
	for _, item := range list {
		idList = append(idList, cast.ToString(item[`id`]))
	}
	if len(idList) > 0 && common.GetNeo4jStatus(userId) {
		graphEntityCountRes, err = common.NewGraphDB(userId).GetEntityCount(idList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	for _, item := range list {
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
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func getLibFileCount(wheres [][]string) (data map[string]int, err error) {
	m := msql.Model(`chat_ai_library_file`, define.Postgres).Alias(`f`).
		Group(`status`).
		Field(`status,count(id) as count`)
	for _, v := range wheres {
		m.Where(v...)
	}
	count, err := m.Select()
	data = map[string]int{
		`learned_count`:      0,
		`learned_wait_count`: 0,
		`learned_err_count`:  0,
	}
	for _, item := range count {
		if tool.InArrayInt(cast.ToInt(item[`status`]), []int{define.FileStatusException, define.FileStatusPartException}) {
			data[`learned_err_count`] += cast.ToInt(item[`count`])
		}
		if tool.InArrayInt(cast.ToInt(item[`status`]), []int{define.FileStatusLearned}) {
			data[`learned_count`] += cast.ToInt(item[`count`])
		}
		if tool.InArrayInt(cast.ToInt(item[`status`]), []int{define.FileStatusWaitSplit, define.FileStatusChunking}) {
			data[`learned_wait_count`] += cast.ToInt(item[`count`])
		}
	}
	return data, err
}

func addLibFile(c *gin.Context, userId, libraryId, libraryType int) ([]int64, error) {
	m := msql.Model(`chat_ai_library_file`, define.Postgres)

	//get params
	docType := cast.ToInt(c.DefaultPostForm(`doc_type`, cast.ToString(define.DocTypeLocal)))
	docUrls := strings.TrimSpace(c.PostForm(`urls`))
	fileName := strings.TrimSpace(c.PostForm(`file_name`))
	content := strings.TrimSpace(c.PostForm(`content`))
	title := strings.TrimSpace(c.PostForm(`title`))
	isQaDoc := cast.ToInt(c.PostForm(`is_qa_doc`))
	qaIndexType := cast.ToInt(c.PostForm(`qa_index_type`))
	docAutoRenewFrequency := cast.ToInt(c.PostForm(`doc_auto_renew_frequency`))
	docAutoRenewMinute := cast.ToInt(c.PostForm(`doc_auto_renew_minute`))
	answerLable := strings.TrimSpace(c.PostForm(`answer_lable`))
	answerColumn := strings.TrimSpace(c.PostForm(`answer_column`))
	questionLable := strings.TrimSpace(c.PostForm(`question_lable`))
	questionColumn := strings.TrimSpace(c.PostForm(`question_column`))
	similarColumn := strings.TrimSpace(c.PostForm(`similar_column`))
	similarLabel := strings.TrimSpace(c.PostForm(`similar_label`))
	pdfParseType := cast.ToInt(c.PostForm(`pdf_parse_type`))
	//document uploaded
	var libraryFiles []*define.UploadInfo
	switch docType {
	case define.DocTypeDiy: // diy library
		if libraryType != define.OpenLibraryType {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_type`))
		}
		md5Hash := tool.MD5(content)
		ext := define.LibDocFileAllowExt[0]
		objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_file`, tool.Date(`Ym`), md5Hash, ext)
		link, err := common.WriteFileByString(objectKey, content)
		if err != nil {
			return nil, err
		}
		libraryFiles = append(libraryFiles, &define.UploadInfo{Name: title,
			Size: int64(len(content)), Ext: ext, Link: link})
	case define.DocTypeCustom: // custom library
		if len(fileName) == 0 || (isQaDoc == define.DocTypeQa && !tool.InArrayInt(qaIndexType, []int{define.QAIndexTypeQuestionAndAnswer, define.QAIndexTypeQuestion})) {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))
		}
		if (libraryType == define.GeneralLibraryType && isQaDoc == define.DocTypeQa) ||
			(libraryType == define.QALibraryType && isQaDoc != define.DocTypeQa) {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `is_qa_doc`))
		}
		libraryFiles = append(libraryFiles, &define.UploadInfo{
			Name: fileName, Size: 0, Ext: `-`, Custom: true,
			Link: define.LocalUploadPrefix + `default/empty_document.pdf`,
		})
	case define.DocTypeOnline: // online library
		if len(docUrls) == 0 {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))
		}
		if libraryType == define.QALibraryType {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_type`))
		}

		type UrlItem struct {
			URL    string `json:"url"`
			Remark string `json:"remark"`
		}
		var urlItems []UrlItem
		err := json.Unmarshal([]byte(docUrls), &urlItems)
		if err != nil {
			return nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_urls`))
		}

		// check url is valid
		for _, urlItem := range urlItems {
			if len(urlItem.URL) == 0 {
				return nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))
			}
			if _, err := url.Parse(urlItem.URL); err != nil {
				return nil, errors.New(i18n.Show(common.GetLang(c), `invalid_url`, `url`))
			}
		}

		// save document and push to nsq
		var fileIds []int64
		for _, urlItem := range urlItems {
			insData := msql.Datas{
				`admin_user_id`:            userId,
				`library_id`:               libraryId,
				`status`:                   define.FileStatusWaitCrawl,
				`file_ext`:                 "html",
				`create_time`:              tool.Time2Int(),
				`update_time`:              tool.Time2Int(),
				`is_table_file`:            cast.ToInt(false),
				`doc_type`:                 define.DocTypeOnline,
				`doc_url`:                  urlItem.URL,
				`doc_auto_renew_frequency`: docAutoRenewFrequency,
				`doc_auto_renew_minute`:    docAutoRenewMinute,
				`remark`:                   urlItem.Remark,
			}
			fileId, err := m.Insert(insData, `id`)
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			fileIds = append(fileIds, fileId)
			if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `admin_user_id`: userId}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(define.CrawlArticleTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
		return fileIds, nil
	case define.DocTypeLocal: // document uploaded
		libFileAloowExts := define.LibFileAllowExt
		switch libraryType {
		case define.OpenLibraryType:
			libFileAloowExts = define.LibDocFileAllowExt
		case define.QALibraryType:
			libFileAloowExts = define.QALibFileAllowExt
		}
		libraryFiles, _ = common.SaveUploadedFileMulti(c, `library_files`, define.LibFileLimitSize, userId, `library_file`, libFileAloowExts)
		if len(libraryFiles) == 0 {
			return nil, errors.New(i18n.Show(common.GetLang(c), `upload_empty`))
		}
	default:
		return nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_type`))
	}
	libraryInfo, _ := common.GetLibraryInfo(libraryId, userId)
	//database dispose
	fileIds := make([]int64, 0)
	var autoSplit bool
	splitParams := common.DefaultSplitParams()
	// question and answer to auto learning
	if len(libraryFiles) == 1 {
		autoSplit = true
	}
	if answerLable != "" || answerColumn != "" {
		autoSplit = true
		splitParams.QuestionColumn = questionColumn
		splitParams.SimilarColumn = similarColumn
		splitParams.AnswerColumn = answerColumn
		splitParams.QuestionLable = questionLable
		splitParams.SimilarLabel = similarLabel
		splitParams.AnswerLable = answerLable
		splitParams.IsQaDoc = isQaDoc
		splitParams.ChunkSize = 512
		splitParams.IsTableFile = cast.ToInt(c.PostForm(`is_table_file`))
	}
	for _, uploadInfo := range libraryFiles {
		status := define.FileStatusWaitSplit
		isTableFile := define.IsTableFile(uploadInfo.Ext)

		if define.IsPdfFile(uploadInfo.Ext) && (pdfParseType == define.PdfParseTypeOcr ||
			pdfParseType == define.PdfParseTypeOcrWithImage ||
			pdfParseType == define.PdfParseTypeOcrAli) {
			status = define.FileStatusInitial
		}
		insData := msql.Datas{
			`admin_user_id`:        userId,
			`library_id`:           libraryId,
			`file_url`:             uploadInfo.Link,
			`file_name`:            uploadInfo.Name,
			`status`:               status,
			`chunk_size`:           512,
			`chunk_overlap`:        0,
			`separators_no`:        `11,12`,
			`enable_extract_image`: true,
			`pdf_parse_type`:       pdfParseType,
			`is_qa_doc`:            isQaDoc,
			`answer_lable`:         answerLable,
			`answer_column`:        answerColumn,
			`question_lable`:       questionLable,
			`question_column`:      questionColumn,
			`similar_label`:        similarLabel,
			`similar_column`:       similarColumn,
			`file_ext`:             uploadInfo.Ext,
			`file_size`:            uploadInfo.Size,
			`create_time`:          tool.Time2Int(),
			`update_time`:          tool.Time2Int(),
			`is_table_file`:        cast.ToInt(isTableFile),
			`doc_type`:             uploadInfo.GetDocType(),
			`doc_url`:              uploadInfo.DocUrl,
		}
		if define.IsPdfFile(uploadInfo.Ext) {
			page, err := api.PageCountFile(common.GetFileByLink(uploadInfo.Link))
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			insData[`ocr_pdf_total`] = page
			if pdfParseType == define.PdfParseTypeOcrAli && page > 1000 {
				return nil, errors.New(i18n.Show(common.GetLang(c), `exceed_max_pdf_page_count`))
			}
		}
		if uploadInfo.Custom {
			insData[`status`] = define.FileStatusLearned
			insData[`html_url`] = insData[`file_url`]
			insData[`is_qa_doc`] = isQaDoc
			if qaIndexType != 0 {
				insData[`qa_index_type`] = qaIndexType
			}
		}
		fileId, err := m.Insert(insData, `id`)
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: int(fileId)})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		fileIds = append(fileIds, fileId)
		if uploadInfo.Custom {
			continue
		}
		if define.IsPdfFile(uploadInfo.Ext) && (pdfParseType == define.PdfParseTypeOcr || pdfParseType == define.PdfParseTypeOcrWithImage) {
			if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `file_url`: uploadInfo.Link}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(define.ConvertHtmlTopic, message); err != nil {
				logs.Error(err.Error())
			}
			continue
		}
		if define.IsPdfFile(uploadInfo.Ext) && pdfParseType == define.PdfParseTypeOcrAli {
			jobId, err := common.SubmitOdcParserJob(c, userId, uploadInfo.Link)
			if err != nil {
				logs.Error(err.Error())
				_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
					`status`:      define.FileStatusException,
					`errmsg`:      err.Error(),
					`update_time`: tool.Time2Int(),
				})
				if err != nil {
					logs.Error(err.Error())
				}
				lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: cast.ToInt(fileId)})
			} else {
				_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
					`ali_ocr_job_id`: jobId,
					`update_time`:    tool.Time2Int(),
				})
				if err != nil {
					logs.Error(err.Error())
				}
				lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: cast.ToInt(fileId)})
			}
			continue
		}

		switch cast.ToInt(libraryInfo[`type`]) {
		case define.GeneralLibraryType:
			if isTableFile {
				autoSplit = true
			}
			splitParams.FileExt = uploadInfo.Ext
			splitParams.ChunkType = cast.ToInt(libraryInfo[`chunk_type`])
			if splitParams.ChunkType == define.ChunkTypeNormal {
				splitParams.SeparatorsNo = cast.ToString(libraryInfo[`normal_chunk_default_separators_no`])
				splitParams.ChunkSize = cast.ToInt(libraryInfo[`normal_chunk_default_chunk_size`])
				splitParams.ChunkOverlap = cast.ToInt(libraryInfo[`normal_chunk_default_chunk_overlap`])
			} else if splitParams.ChunkType == define.ChunkTypeSemantic {
				splitParams.SemanticChunkSize = cast.ToInt(libraryInfo[`semantic_chunk_default_chunk_size`])
				splitParams.SemanticChunkOverlap = cast.ToInt(libraryInfo[`semantic_chunk_default_chunk_overlap`])
				splitParams.SemanticChunkThreshold = cast.ToInt(libraryInfo[`semantic_chunk_default_threshold`])
			} else if splitParams.ChunkType == define.ChunkTypeAi {
				splitParams.AiChunkModel = cast.ToString(libraryInfo[`ai_chunk_model`])
				splitParams.AiChunkModelConfigId = cast.ToInt(libraryInfo[`ai_chunk_model_config_id`])
				splitParams.AiChunkPrumpt = cast.ToString(libraryInfo[`ai_chunk_prumpt`])
				splitParams.AiChunkSize = cast.ToInt(libraryInfo[`ai_chunk_size`])
				splitParams.AiChunkNew = !isTableFile
				splitParams.AiChunkTaskId = uuid.New().String()
			}
			splitParams.SemanticChunkModelConfigId = cast.ToInt(libraryInfo[`model_config_id`])
			splitParams.SemanticChunkUseModel = libraryInfo[`use_model`]
		case define.QALibraryType:
			if !define.IsTableFile(uploadInfo.Ext) {
				autoSplit = true
			}
		case define.OpenLibraryType:
			if define.IsMdFile(uploadInfo.Ext) {
				autoSplit = true
			}
		}
		if autoSplit {
			go common.AutoSplitLibFile(userId, int(fileId), splitParams)
		}
	}
	return fileIds, nil
}

func AddLibraryFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	if libraryId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//common save
	fileIds, err := addLibFile(c, userId, libraryId, cast.ToInt(info[`type`]))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`file_ids`: fileIds}, nil))
}

func DelLibraryFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	ids := cast.ToString(c.PostForm(`id`))
	for _, id := range strings.Split(ids, `,`) {
		id := cast.ToInt(id)
		if id <= 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
			return
		}
		info, err := common.GetLibFileInfo(id, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(info) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})
		//dispose relation data
		dataIdList, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`file_id`, cast.ToString(id)).
			Where(`category_id`, `0`).
			ColumnArr(`id`)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(dataIdList) == 0 {
			continue
		}

		_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`id`, `in`, strings.Join(dataIdList, `,`)).
			Delete()
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
			Where(`data_id`, `in`, strings.Join(dataIdList, `,`)).
			Delete()
		if err != nil {
			logs.Error(err.Error())
		}
		_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`file_id`, cast.ToString(id)).
			Update(msql.Datas{`isolated`: true})
		if err != nil {
			logs.Error(err.Error())
		}

		if common.GetNeo4jStatus(userId) {
			for _, dataId := range dataIdList {
				err = common.NewGraphDB(userId).DeleteByData(cast.ToInt(dataId))
				if err != nil {
					logs.Error(err.Error())
				}
			}
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetLibFileInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `file_deleted`))))
		return
	}
	library, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, info[`library_id`]).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(library) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	info[`library_name`] = library[`library_name`]
	info[`library_type`] = library[`type`]
	info[`default_chunk_type`] = library[`chunk_type`]
	info[`normal_chunk_default_separators_no`] = library[`normal_chunk_default_separators_no`]
	info[`normal_chunk_default_chunk_size`] = library[`normal_chunk_default_chunk_size`]
	info[`normal_chunk_default_chunk_overlap`] = library[`normal_chunk_default_chunk_overlap`]
	info[`semantic_chunk_default_chunk_size`] = library[`semantic_chunk_default_chunk_size`]
	info[`semantic_chunk_default_chunk_overlap`] = library[`semantic_chunk_default_chunk_overlap`]
	info[`semantic_chunk_default_threshold`] = library[`semantic_chunk_default_threshold`]
	info[`default_model_config_id`] = library[`model_config_id`]
	info[`default_use_model`] = library[`use_model`]
	if info[`ai_chunk_model`] == "" {
		info[`ai_chunk_model`] = library[`ai_chunk_model`]
	}
	if cast.ToInt(info[`ai_chunk_model_config_id`]) == 0 {
		info[`ai_chunk_model_config_id`] = library[`ai_chunk_model_config_id`]
	}
	if info[`ai_chunk_prumpt`] == "" {
		info[`ai_chunk_prumpt`] = library[`ai_chunk_prumpt`]
	}
	if cast.ToInt(info[`ai_chunk_size`]) == 0 {
		info[`ai_chunk_size`] = library[`ai_chunk_size`]
	}
	info[`graph_switch`] = library[`graph_switch`]
	if !common.GetNeo4jStatus(userId) {
		info[`graph_switch`] = "0"
	}

	var separators []string
	for _, noStr := range strings.Split(info[`separators_no`], `,`) {
		if len(noStr) == 0 {
			continue
		}
		no := cast.ToInt(noStr)
		separators = append(separators, cast.ToString(define.SeparatorsList[no-1][`name`]))
	}
	info[`separators`] = strings.Join(separators, ", ")

	c.String(http.StatusOK, lib_web.FmtJson(info, nil))
}

func GetLibRawFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `file_deleted`))))
		return
	}
	if !common.LinkExists(info[`file_url`]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.File(common.GetFileByLink(info[`file_url`]))
	return
}

func GetLibRawFileOnePage(c *gin.Context) {
	adminUserId := cast.ToInt(c.Query(`admin_user_id`))
	id := cast.ToInt(c.Query(`id`))
	page := cast.ToInt(c.Query(`page`))
	if adminUserId == 0 || id <= 0 || page <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `file_deleted`))))
		return
	}
	if !common.LinkExists(info[`file_url`]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	file := common.GetFileByLink(info[`file_url`])
	outDir := define.UploadDir + fmt.Sprintf(`pdf_split/%s`, tool.Random(8))
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(outDir)
	_ = tool.MkDirAll(outDir)
	if err := api.SplitFile(file, outDir, 1, nil); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	filename := strings.TrimSuffix(filepath.Base(file), `.pdf`)
	item := fmt.Sprintf(`%s/%s_%d.pdf`, outDir, filename, page)
	if !tool.IsFile(item) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	fmt.Printf(item)
	c.File(item)
	return
}

func GetSeparatorsList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	list := make([]map[string]any, 0)
	for _, item := range define.SeparatorsList {
		name := i18n.Show(common.GetLang(c), cast.ToString(item[`name`]))
		list = append(list, map[string]any{`no`: item[`no`], `name`: name})
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetLibFileSplit(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.Query(`id`))
	if fileId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	pdfPageNum := cast.ToInt(c.Query(`pdf_page_num`))
	if pdfPageNum < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	splitParams := define.SplitParams{
		IsDiySplit:                 cast.ToInt(c.Query(`is_diy_split`)),
		SeparatorsNo:               strings.TrimSpace(c.Query(`separators_no`)),
		Separators:                 make([]string, 0),
		ChunkSize:                  cast.ToInt(c.Query(`chunk_size`)),
		ChunkOverlap:               cast.ToInt(c.Query(`chunk_overlap`)),
		IsQaDoc:                    cast.ToInt(c.Query(`is_qa_doc`)),
		QuestionLable:              strings.TrimSpace(c.Query(`question_lable`)),
		SimilarLabel:               strings.TrimSpace(c.Query(`similar_label`)),
		AnswerLable:                strings.TrimSpace(c.Query(`answer_lable`)),
		QuestionColumn:             strings.TrimSpace(c.Query(`question_column`)),
		SimilarColumn:              strings.TrimSpace(c.Query(`similar_column`)),
		AnswerColumn:               strings.TrimSpace(c.Query(`answer_column`)),
		EnableExtractImage:         cast.ToBool(c.Query(`enable_extract_image`)),
		ChunkType:                  cast.ToInt(c.Query(`chunk_type`)),
		SemanticChunkSize:          cast.ToInt(c.Query(`semantic_chunk_size`)),
		SemanticChunkOverlap:       cast.ToInt(c.Query(`semantic_chunk_overlap`)),
		SemanticChunkThreshold:     cast.ToInt(c.Query(`semantic_chunk_threshold`)),
		SemanticChunkModelConfigId: cast.ToInt(c.Query(`semantic_chunk_model_config_id`)),
		SemanticChunkUseModel:      strings.TrimSpace(c.Query(`semantic_chunk_use_model`)),
		AiChunkPrumpt:              cast.ToString(c.Query(`ai_chunk_prumpt`)),
		AiChunkModel:               strings.TrimSpace(c.Query(`ai_chunk_model`)),
		AiChunkModelConfigId:       cast.ToInt(c.Query(`ai_chunk_model_config_id`)),
		AiChunkSize:                cast.ToInt(c.Query(`ai_chunk_size`)),
		AiChunkTaskId:              strings.TrimSpace(c.Query(`ai_chunk_task_id`)),
		AiChunkPreview:             cast.ToBool(c.Query(`ai_chunk_preview`)),
	}
	if splitParams.ChunkType == define.ChunkTypeSemantic {
		if splitParams.SemanticChunkModelConfigId <= 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `semantic_chunk_model_config_id`))))
			return
		}
		if len(splitParams.SemanticChunkUseModel) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `semantic_chunk_use_model`))))
			return
		}
	} else if splitParams.ChunkType == define.ChunkTypeAi && splitParams.AiChunkTaskId == "" {
		if ok := common.CheckModelIsValid(userId, splitParams.AiChunkModelConfigId, splitParams.AiChunkModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_model`))))
			return
		}
		if len(splitParams.AiChunkPrumpt) == 0 || utf8.RuneCountInString(splitParams.AiChunkPrumpt) > 500 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_prumpt`))))
			return
		}
		splitParams.AiChunkTaskId = uuid.New().String()
		splitParams.AiChunkNew = true
	}
	list, wordTotal, splitParams, err := common.GetLibFileSplit(userId, fileId, pdfPageNum, splitParams, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), err.Error()))))
		return
	}
	data := map[string]any{`split_params`: splitParams, `list`: list, `word_total`: wordTotal}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func SaveLibFileSplit(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.PostForm(`id`))
	wordTotal := cast.ToInt(c.PostForm(`word_total`))
	splitParams, list := define.SplitParams{}, make([]define.DocSplitItem, 0)
	qaIndexType := cast.ToInt(c.PostForm(`qa_index_type`))
	pdfPageNum := cast.ToInt(c.PostForm(`pdf_page_num`))
	if pdfPageNum < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	err := tool.JsonDecodeUseNumber(c.PostForm(`split_params`), &splitParams)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `split_params`))))
		return
	}
	if splitParams.ChunkType == define.ChunkTypeAi && splitParams.AiChunkTaskId == "" {
		if ok := common.CheckModelIsValid(userId, splitParams.AiChunkModelConfigId, splitParams.AiChunkModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_model`))))
			return
		}
		if len(splitParams.AiChunkPrumpt) == 0 || len(splitParams.AiChunkPrumpt) > 500 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_prumpt`))))
			return
		}
		splitParams.AiChunkTaskId = uuid.New().String()
		splitParams.AiChunkNew = true
	}

	if splitParams.ChunkType == define.ChunkTypeAi && !splitParams.AiChunkNew {
		// do nothging
		// 同步orc解析pdf的情况很耗时，直接使用前端传入的值保存
	} else if !(pdfPageNum > 0 && (splitParams.PdfParseType == define.PdfParseTypeOcr || splitParams.PdfParseType == define.PdfParseTypeOcrWithImage)) {
		list, wordTotal, splitParams, err = common.GetLibFileSplit(userId, fileId, pdfPageNum, splitParams, common.GetLang(c))
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
	}

	if err = tool.JsonDecodeUseNumber(c.PostForm(`list`), &list); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `list`))))
		return
	}
	if fileId <= 0 || wordTotal <= 0 || len(list) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if splitParams.ChunkType == define.ChunkTypeAi && splitParams.AiChunkNew {
		go func ()  {
			if err = common.SaveAISplitDocs(userId, fileId, wordTotal, qaIndexType, splitParams, list, pdfPageNum, common.GetLang(c)); err != nil {
				c.String(http.StatusOK, lib_web.FmtJson(nil, err))
				return
			}
		}()
		c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
	} else {
		dataIds, err := common.SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType, splitParams, list, pdfPageNum, common.GetLang(c))
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		c.String(http.StatusOK, lib_web.FmtJson(dataIds, nil))
	}
}

func GetLibFileExcelTitle(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if info[`is_table_file`] != cast.ToString(define.FileIsTable) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `is_not_excel`))))
		return
	}
	rows, err := common.ParseTabFile(info[`file_url`], info[`file_ext`])
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(rows) < 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `excel_less_row`))))
		return
	}

	var data = make(map[string]string)
	for i, v := range rows[0] {
		column, err := common.IdentifierFromColumnIndex(i)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		data[column] = v
	}

	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
	return
}

func EditLibFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	docAutoRenewFrequency := cast.ToInt(c.PostForm(`doc_auto_renew_frequency`))
	docAutoRenewMinute := cast.ToInt(c.PostForm(`doc_auto_renew_minute`))
	fileName := cast.ToString(c.PostForm(`file_name`))
	if id <= 0 || len(fileName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if docAutoRenewFrequency > 5 || docAutoRenewFrequency < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_auto_renew_frequency`))))
		return
	}
	if docAutoRenewMinute < 0 || docAutoRenewMinute > 3600 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_auto_renew_minute`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	updateData := msql.Datas{
		`file_name`:   fileName,
		`update_time`: tool.Time2Int(),
	}
	if docAutoRenewFrequency != cast.ToInt(info[`doc_auto_renew_minute`]) && docAutoRenewFrequency > 0 {
		updateData[`doc_auto_renew_frequency`] = docAutoRenewFrequency
		//updateData[`doc_last_renew_time`] = 0
	}
	if docAutoRenewMinute != cast.ToInt(info[`doc_auto_renew_minute`]) {
		updateData[`doc_auto_renew_minute`] = docAutoRenewMinute
		//updateData[`doc_last_renew_time`] = 0
	}

	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Update(updateData)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}

	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ManualCrawl(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToInt(info[`doc_type`]) != define.DocTypeOnline {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if message, err := tool.JsonEncode(map[string]any{`file_id`: info[`id`], `admin_user_id`: userId}); err != nil {
		logs.Error(err.Error())
	} else if err := common.AddJobs(define.CrawlArticleTopic, message); err != nil {
		logs.Error(err.Error())
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ConstructGraph(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	if !common.GetNeo4jStatus(userId) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `graph_not_opened`))))
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if cast.ToInt(info[`graph_status`]) != define.GraphStatusNotStart {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if err = common.ConstructGraph(id); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetFileGraphInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	if !common.GetNeo4jStatus(userId) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `graph_not_opened`))))
		return
	}
	fileId := cast.ToInt(c.Query(`file_id`))
	dataId := cast.ToInt(c.Query(`data_id`))
	searchTerm := strings.TrimSpace(c.Query(`search_term`))

	if fileId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(fileId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	paragraphList, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Alias(`d`).
		Join(`chat_ai_library_file f`, `d.file_id = f.id`, `left`).
		Where(`d.admin_user_id`, cast.ToString(userId)).
		Where(`d.file_id`, cast.ToString(fileId)).
		Field(`d.id,d.file_id,d.content,f.file_name`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	paragraphMap := make(map[string]map[string]string)
	for _, paragraph := range paragraphList {
		paragraphId := cast.ToString(paragraph["id"])
		paragraphMap[paragraphId] = paragraph
	}

	// 获取节点数据，添加搜索条件筛选
	nodeRes, err := common.NewGraphDB(userId).GetFileNodes(fileId, dataId, searchTerm)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// 准备节点数据
	nodes := make([]map[string]any, 0)
	nodeIds := make(map[int64]bool)

	for _, record := range nodeRes.Records {
		nodeId, _ := record.Get("id")
		name, _ := record.Get("name")
		dataId, _ := record.Get("data_id")

		// 将节点ID添加到nodeIds集合中
		if idInt64, ok := nodeId.(int64); ok {
			nodeIds[idInt64] = true
		}

		node := map[string]any{
			"id":      nodeId,
			"name":    name,
			"data_id": dataId,
		}

		// 将段落数据添加到节点中
		if dataIdStr := cast.ToString(dataId); dataIdStr != "" && paragraphMap[dataIdStr] != nil {
			paragraph := paragraphMap[dataIdStr]
			node["file_id"] = paragraph["file_id"]
			node["file_name"] = paragraph["file_name"]
			node["content"] = paragraph["content"]
		}

		nodes = append(nodes, node)
	}

	// 如果没有找到节点，直接返回空结果
	if len(nodes) == 0 {
		result := map[string]any{
			"nodes": nodes,
			"edges": []map[string]any{},
		}
		c.String(http.StatusOK, lib_web.FmtJson(result, nil))
		return
	}

	// 获取边数据，添加搜索条件筛选
	edgeRes, err := common.NewGraphDB(userId).GetFileRelationships(fileId, dataId, searchTerm)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// 准备边数据，只保留与筛选后节点相关的边
	edges := make([]map[string]any, 0)
	// 用于去重的映射，记录已经添加的边的方向
	edgeDirections := make(map[string]bool)

	for _, record := range edgeRes.Records {
		edgeId, _ := record.Get("id")
		fromId, _ := record.Get("from_id")
		toId, _ := record.Get("to_id")
		label, _ := record.Get("label")

		// 检查边的两端是否都在筛选后的节点集合中
		fromIdInt64, fromOk := fromId.(int64)
		toIdInt64, toOk := toId.(int64)

		if !fromOk || !toOk {
			continue
		}

		// 只有当边的两端都在筛选后的节点列表中，才考虑添加这条边
		if nodeIds[fromIdInt64] && nodeIds[toIdInt64] {
			// 创建方向键，格式为：较小ID-较大ID
			var directionKey string
			if fromIdInt64 < toIdInt64 {
				directionKey = fmt.Sprintf("%d-%d", fromIdInt64, toIdInt64)
			} else {
				directionKey = fmt.Sprintf("%d-%d", toIdInt64, fromIdInt64)
			}

			// 如果这个方向已经被添加过，则跳过
			if edgeDirections[directionKey] {
				continue
			}

			// 标记这个方向已经被添加
			edgeDirections[directionKey] = true

			edges = append(edges, map[string]any{
				"id":      edgeId,
				"from_id": fromId,
				"to_id":   toId,
				"label":   label,
			})
		}
	}

	// 构建结果
	result := map[string]any{
		"nodes": nodes,
		"edges": edges,
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func RestudyLibraryFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	pdfParseType := cast.ToInt(c.PostForm(`pdf_parse_type`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !tool.InArray(pdfParseType, []int{define.PdfParseTypeText, define.PdfParseTypeOcr, define.PdfParseTypeOcrWithImage, define.PdfParseTypeOcrAli}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `pdf_parse_type`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if !(define.IsPdfFile(info[`file_ext`])) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `only_support_pdf`))))
		return
	}
	page, err := api.PageCountFile(common.GetFileByLink(info[`file_url`]))
	if err != nil {
		logs.Error(err.Error())
	}
	if page > 1000 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exceed_max_pdf_page_count`))))
		return
	}

	// update status
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{
			`status`:         cast.ToString(define.FileStatusInitial),
			`pdf_parse_type`: pdfParseType,
			`ocr_pdf_index`:  0,
			`ocr_pdf_total`:  page,
		})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	if pdfParseType == define.PdfParseTypeOcrAli {
		jobId, err := common.SubmitOdcParserJob(c, userId, info["file_url"])
		if err != nil {
			logs.Error(err.Error())
			_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
				`status`:      define.FileStatusException,
				`errmsg`:      err.Error(),
				`update_time`: tool.Time2Int(),
			})
			if err != nil {
				logs.Error(err.Error())
			}
			lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})
		} else {
			_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
				`ali_ocr_job_id`: jobId,
				`update_time`:    tool.Time2Int(),
			})
			if err != nil {
				logs.Error(err.Error())
			}
			lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})
		}
	} else {
		// delete cache
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})

		// push to nsq task
		if message, err := tool.JsonEncode(map[string]any{`file_id`: id, `admin_user_id`: info[`admin_user_id`], `file_url`: info[`file_url`]}); err != nil {
			logs.Error(err.Error())
		} else if err := common.AddJobs(define.ConvertHtmlTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func RenewLibraryFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if info[`doc_type`] != cast.ToString(define.DocTypeOnline) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `is_not_online`))))
		return
	}
	if php2go.InArray(info[`status`], []string{cast.ToString(define.FileStatusWaitCrawl), cast.ToString(define.FileStatusCrawling)}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `doc_is_crawling`))))
		return
	}

	// update status
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{`status`: cast.ToString(define.FileStatusWaitCrawl)})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}

	// delete cache
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})

	// push to nsq task
	if message, err := tool.JsonEncode(map[string]any{`file_id`: id, `admin_user_id`: info[`admin_user_id`]}); err != nil {
		logs.Error(err.Error())
	} else if err := common.AddJobs(define.CrawlArticleTopic, message); err != nil {
		logs.Error(err.Error())
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ReconstructVector(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	dataList, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Where(`file_id`, cast.ToString(id)).
		Where(`status`, cast.ToString(define.VectorStatusException)).
		Field(`id,file_id`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
		return
	}
	for _, data := range dataList {
		message, err := tool.JsonEncode(map[string]any{`id`: data[`id`], `file_id`: data[`file_id`]})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if err = common.AddJobs(define.ConvertVectorTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ReconstructCategoryVector(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	dataList, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Alias(`idx`).
		Where(`library_id`, cast.ToString(id)).
		Where(`status`, cast.ToString(define.VectorStatusException)).
		Where(`exists (select 1 from chat_ai_library_file_data where id = idx.data_id and category_id > 0)`).
		Field(`id`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
		return
	}
	for _, data := range dataList {
		message, err := tool.JsonEncode(map[string]any{`id`: data[`id`], `file_id`: 0})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if err = common.AddJobs(define.ConvertVectorTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ReconstructGraph(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	if !common.GetNeo4jStatus(userId) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `graph_not_opened`))))
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if !tool.InArrayInt(cast.ToInt(info[`graph_status`]), []int{define.GraphStatusPartlyConverted, define.GraphStatusException}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
		return
	}

	dataList, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`file_id`, cast.ToString(id)).
		Where(`graph_status`, cast.ToString(define.GraphStatusException)).
		Field(`id,file_id`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for _, data := range dataList {
		message, err := tool.JsonEncode(map[string]any{`id`: data[`id`], `file_id`: data[`file_id`]})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if err = common.AddJobs(define.ConvertGraphTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func CancelOcrPdf(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if cast.ToInt(info[`status`]) != define.FileStatusInitial {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `status_exception`))))
		return
	}
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{"status": define.FileStatusCancelled})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ReadLibFileExcelTitle(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB 最大内存
		logs.Error(err.Error())
		common.FmtError(c, `file_err`)
		return
	}

	// save upload file
	if c.Request.MultipartForm == nil || len(c.Request.MultipartForm.File) == 0 {
		common.FmtError(c, `upload_empty`)
		return
	}
	for _, fileHeader := range c.Request.MultipartForm.File[`library_files`] {
		if fileHeader == nil {
			common.FmtError(c, `upload_empty`)
			return
		}
		uploadInfo, err := common.ReadUploadedFile(fileHeader, define.LibFileLimitSize, define.FormFileAllowExt)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `file_err`)
			return
		}
		if uploadInfo == nil || uploadInfo.Columns == "" {
			common.FmtError(c, `upload_empty`)
			return
		}

		splitData := strings.Split(uploadInfo.Columns, "\r\n")
		title := make([]string, 0)
		var data = make(map[string]string)
		for key, item := range splitData {
			upData := strings.Split(item, ",")
			if len(upData) < len(title) {
				continue
			}
			if key == 0 {
				title = upData
				for i, v := range title {
					column, err := common.IdentifierFromColumnIndex(i)
					if err != nil {
						c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
						return
					}
					data[column] = v
				}
				break
			}
		}
		common.FmtOk(c, data)
		return
	}
}

func GetLibFileSplitAiChunks(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	info, err := common.GetLibFileInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	taskId := strings.TrimSpace(c.Query(`task_id`))
	var list = make(map[string]any)
	if err := lib_redis.GetCacheWithBuild(define.Redis, &common.LibFileSplitAiChunksBacheHandle{TaskId: taskId}, &list, 1*time.Hour); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, list)
}
