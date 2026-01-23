// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetLibFileList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeGetLibFileListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)
	data, httpStatus, err := BridgeGetLibFileList(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func GetLibFileCount(wheres [][]string) (data map[string]int, err error) {
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

func addLibFile(multipartForm *multipart.Form, lang string, userId, libraryId, libraryType int, chunkParam *define.ChunkParam, addFileParam *BridgeAddLibraryFileReq) ([]int64, error) {
	m := msql.Model(`chat_ai_library_file`, define.Postgres)

	//get params
	docType := cast.ToInt(addFileParam.DocType)
	docUrls := strings.TrimSpace(addFileParam.Urls)
	fileName := strings.TrimSpace(addFileParam.FileName)
	content := strings.TrimSpace(addFileParam.Content)
	title := strings.TrimSpace(addFileParam.Title)
	isQaDoc := cast.ToInt(addFileParam.IsQaDoc)
	qaIndexType := cast.ToInt(addFileParam.QaIndexType)
	docAutoRenewFrequency := cast.ToInt(addFileParam.DocAutoRenewFrequency)
	docAutoRenewMinute := cast.ToInt(addFileParam.DocAutoRenewMinute)
	answerLable := strings.TrimSpace(addFileParam.AnswerLable)
	answerColumn := strings.TrimSpace(addFileParam.AnswerColumn)
	questionLable := strings.TrimSpace(addFileParam.QuestionLable)
	questionColumn := strings.TrimSpace(addFileParam.QuestionColumn)
	similarColumn := strings.TrimSpace(addFileParam.SimilarColumn)
	similarLabel := strings.TrimSpace(addFileParam.SimilarLabel)
	pdfParseType := cast.ToInt(addFileParam.PdfParseType)
	// 问答知识库（libraryType == define.QALibraryType groupId是chat_ai_library_file 和 chat_ai_library_file_data 的group_id字段；
	groupId := max(0, cast.ToInt(addFileParam.GroupId))
	officialArticleId := strings.TrimSpace(addFileParam.OfficialArticleId)
	officialArticleUpdateTime := cast.ToInt64(addFileParam.OfficialArticleUpdateTime)
	feishuDocumentIdList := strings.TrimSpace(addFileParam.FeishuDocumentIdList)
	feishuAppId := strings.TrimSpace(addFileParam.FeishuAppId)
	feishuAppSecret := strings.TrimSpace(addFileParam.FeishuAppSecret)
	//document uploaded
	var libraryFiles []*define.UploadInfo
	switch docType {
	case define.DocTypeDiy: // diy library
		if libraryType != define.OpenLibraryType {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_type`))
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
			return nil, errors.New(i18n.Show(lang, `param_lack`))
		}
		if (libraryType == define.GeneralLibraryType && isQaDoc == define.DocTypeQa) ||
			(libraryType == define.QALibraryType && isQaDoc != define.DocTypeQa) {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `is_qa_doc`))
		}
		libraryFiles = append(libraryFiles, &define.UploadInfo{
			Name: fileName, Size: 0, Ext: `-`, Custom: true,
			Link: define.LocalUploadPrefix + `default/empty_document.pdf`,
		})
	case define.DocTypeOnline: // online library
		if len(docUrls) == 0 {
			return nil, errors.New(i18n.Show(lang, `param_lack`))
		}
		if libraryType == define.QALibraryType {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_type`))
		}

		type UrlItem struct {
			URL    string `json:"url"`
			Remark string `json:"remark"`
		}
		var urlItems []UrlItem
		err := json.Unmarshal([]byte(docUrls), &urlItems)
		if err != nil {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_urls`))
		}

		// check url is valid
		for _, urlItem := range urlItems {
			if len(urlItem.URL) == 0 {
				return nil, errors.New(i18n.Show(lang, `param_lack`))
			}
			if _, err := url.Parse(urlItem.URL); err != nil {
				return nil, errors.New(i18n.Show(lang, `invalid_url`, `url`))
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
				`group_id`:                 groupId,
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
		libraryFiles, _ = common.SaveUploadedFileMulti(multipartForm, `library_files`, define.LibFileLimitSize, userId, `library_file`, libFileAloowExts)
		if len(libraryFiles) == 0 {
			return nil, errors.New(i18n.Show(lang, `upload_empty`))
		}
	case define.DocTypeOfficial: // 公众号文章
		if libraryType != define.OfficialLibraryType {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_type`))
		}
		md5Hash := tool.MD5(content)
		objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_file`, tool.Date(`Ym`), md5Hash, `html`)
		link, err := common.WriteFileByString(objectKey, content)
		if err != nil {
			return nil, err
		}
		libraryFiles = append(libraryFiles, &define.UploadInfo{Name: title, Size: int64(len(content)), Ext: `html`, Link: link})
	case define.DocTypeFeishu: // 飞书知识库
		if libraryType != define.GeneralLibraryType {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_type`))
		}
		if strings.TrimSpace(feishuDocumentIdList) == "" || len(feishuAppId) == 0 || len(feishuAppSecret) == 0 {
			return nil, errors.New(i18n.Show(lang, `param_lack`))
		}
		// feishu_document_id 去重
		feishuDocumentIdSet := make(map[string]struct{})
		for _, rawId := range strings.Split(feishuDocumentIdList, ",") {
			docId := strings.TrimSpace(rawId)
			if docId == "" {
				continue
			}
			feishuDocumentIdSet[docId] = struct{}{}
		}
		if len(feishuDocumentIdSet) == 0 {
			return nil, errors.New(i18n.Show(lang, `param_lack`))
		}

		existRows, existErr := msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`library_id`, cast.ToString(libraryId)).
			Where(`doc_type`, cast.ToString(define.DocTypeFeishu)).
			Where(`delete_time`, `0`).
			Field(`id,feishu_document_id`).
			Select()
		if existErr != nil {
			logs.Error(existErr.Error())
			return nil, existErr
		}
		existIdMap := make(map[string]int64, len(existRows))
		for _, row := range existRows {
			docId := strings.TrimSpace(row[`feishu_document_id`])
			if docId == "" {
				continue
			}
			existIdMap[docId] = cast.ToInt64(row[`id`])
		}

		// save document and push to nsq
		var fileIds []int64
		for feishuDocumentId := range feishuDocumentIdSet {
			// 已存在则更新并复用（同一个知识库同一个feishu_document_id只保留一条）
			if existId, ok := existIdMap[feishuDocumentId]; ok && existId > 0 {
				_, updateErr := msql.Model(`chat_ai_library_file`, define.Postgres).
					Where(`id`, cast.ToString(existId)).
					Update(msql.Datas{
						`status`:             define.FileStatusWaitCrawl,
						`update_time`:        tool.Time2Int(),
						`is_table_file`:      cast.ToInt(false),
						`group_id`:           groupId,
						`feishu_app_id`:      feishuAppId,
						`feishu_app_secret`:  feishuAppSecret,
						`feishu_document_id`: feishuDocumentId,
					})
				if updateErr != nil {
					logs.Error(updateErr.Error())
					continue
				}
				lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: int(existId)})
				fileIds = append(fileIds, existId)
				if message, err := tool.JsonEncode(map[string]any{`file_id`: existId, `admin_user_id`: userId}); err != nil {
					logs.Error(err.Error())
				} else if err := common.AddJobs(define.CrawlFeishuDocTopic, message); err != nil {
					logs.Error(err.Error())
				}
				continue
			}
			insData := msql.Datas{
				`admin_user_id`:      userId,
				`library_id`:         libraryId,
				`status`:             define.FileStatusWaitCrawl,
				`create_time`:        tool.Time2Int(),
				`update_time`:        tool.Time2Int(),
				`is_table_file`:      cast.ToInt(false),
				`doc_type`:           define.DocTypeFeishu,
				`group_id`:           groupId,
				`feishu_app_id`:      feishuAppId,
				`feishu_app_secret`:  feishuAppSecret,
				`feishu_document_id`: feishuDocumentId,
			}

			fileId, insertErr := m.Insert(insData, `id`)
			if insertErr != nil {
				// 可能并发/历史数据导致唯一约束冲突：补查存在则更新复用，否则记录错误
				old, findErr := msql.Model(`chat_ai_library_file`, define.Postgres).
					Where(`admin_user_id`, cast.ToString(userId)).
					Where(`library_id`, cast.ToString(libraryId)).
					Where(`doc_type`, cast.ToString(define.DocTypeFeishu)).
					Where(`feishu_document_id`, feishuDocumentId).
					Where(`delete_time`, `0`).
					Find()
				if findErr != nil {
					logs.Error(findErr.Error())
					continue
				}
				if len(old) > 0 {
					existId := cast.ToInt64(old[`id`])
					existIdMap[feishuDocumentId] = existId
					_, updateErr := msql.Model(`chat_ai_library_file`, define.Postgres).
						Where(`id`, cast.ToString(existId)).
						Update(msql.Datas{
							`status`:             define.FileStatusWaitCrawl,
							`update_time`:        tool.Time2Int(),
							`is_table_file`:      cast.ToInt(false),
							`group_id`:           groupId,
							`feishu_app_id`:      feishuAppId,
							`feishu_app_secret`:  feishuAppSecret,
							`feishu_document_id`: feishuDocumentId,
						})
					if updateErr != nil {
						logs.Error(updateErr.Error())
						continue
					}
					lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: int(existId)})
					fileIds = append(fileIds, existId)
					if message, err := tool.JsonEncode(map[string]any{`file_id`: existId, `admin_user_id`: userId}); err != nil {
						logs.Error(err.Error())
					} else if err := common.AddJobs(define.CrawlFeishuDocTopic, message); err != nil {
						logs.Error(err.Error())
					}
					continue
				}
				logs.Error(insertErr.Error())
				continue
			}
			existIdMap[feishuDocumentId] = fileId
			fileIds = append(fileIds, fileId)
			if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `admin_user_id`: userId}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(define.CrawlFeishuDocTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
		return fileIds, nil
	default:
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `doc_type`))
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
		splitParams.QaIndexType = cast.ToInt(libraryInfo[`qa_index_type`])
		splitParams.ChunkSize = 512
	}
	for _, uploadInfo := range libraryFiles {
		status := define.FileStatusWaitSplit
		isTableFile := define.IsTableFile(uploadInfo.Ext)
		splitParams.IsTableFile = cast.ToInt(isTableFile)

		if define.IsPdfFile(uploadInfo.Ext) && (pdfParseType == define.PdfParseTypeOcr ||
			pdfParseType == define.PdfParseTypeOcrWithImage ||
			pdfParseType == define.PdfParseTypeOcrAli) {
			status = define.FileStatusInitial
		}
		insData := msql.Datas{
			`admin_user_id`:                userId,
			`library_id`:                   libraryId,
			`file_url`:                     uploadInfo.Link,
			`file_name`:                    uploadInfo.Name,
			`status`:                       status,
			`chunk_size`:                   512,
			`chunk_overlap`:                0,
			`separators_no`:                `12,11`,
			`enable_extract_image`:         true,
			`pdf_parse_type`:               pdfParseType,
			`is_qa_doc`:                    isQaDoc,
			`answer_lable`:                 answerLable,
			`answer_column`:                answerColumn,
			`question_lable`:               questionLable,
			`question_column`:              questionColumn,
			`similar_label`:                similarLabel,
			`similar_column`:               similarColumn,
			`file_ext`:                     uploadInfo.Ext,
			`file_size`:                    uploadInfo.Size,
			`create_time`:                  tool.Time2Int(),
			`update_time`:                  tool.Time2Int(),
			`is_table_file`:                cast.ToInt(isTableFile),
			`doc_type`:                     uploadInfo.GetDocType(),
			`doc_url`:                      uploadInfo.DocUrl,
			`group_id`:                     groupId,
			`official_article_id`:          officialArticleId,
			`official_article_update_time`: officialArticleUpdateTime,
		}
		if define.IsPdfFile(uploadInfo.Ext) {
			page, err := api.PageCountFile(common.GetFileByLink(uploadInfo.Link))
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			insData[`ocr_pdf_total`] = page
			if pdfParseType == define.PdfParseTypeOcrAli && page > 1000 {
				return nil, errors.New(i18n.Show(lang, `exceed_max_pdf_page_count`))
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
		var fileId int64
		var err error
		if len(officialArticleId) > 0 { // 公众号知识库需要根据article_id来做更新操作
			old, err := m.Where(`library_id`, cast.ToString(libraryId)).Where(`official_article_id`, officialArticleId).Find()
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			insData[`html_url`] = insData[`doc_url`]
			if len(old) == 0 {
				fileId, err = m.Insert(insData, `id`)
			} else {
				_, err = m.Where(`id`, old[`id`]).Update(insData)
				fileId = cast.ToInt64(old[`id`])
			}
		} else { // 普通知识库直接插入
			fileId, err = m.Insert(insData, `id`)
			//set use guide finish
			if docType == define.DocTypeLocal && cast.ToInt(libraryInfo[`is_default`]) == define.NotDefault && define.IsPdfFile(uploadInfo.Ext) {
				_ = common.SetStepFinish(userId, define.StepImportPdf)
			}
			if docType == define.DocTypeLocal && cast.ToInt(libraryInfo[`is_default`]) == define.NotDefault && define.IsDocxFile(uploadInfo.Ext) {
				_ = common.SetStepFinish(userId, define.StepImportWord)
			}
		}

		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: int(fileId)})
		if err != nil {
			logs.Error(err.Error())
			continue
		}

		fileIds = append(fileIds, fileId)
		if uploadInfo.Custom {
			continue
		}
		switch cast.ToInt(libraryInfo[`type`]) {
		case define.GeneralLibraryType:
			if isTableFile {
				autoSplit = true
			}
			splitParams.FileExt = uploadInfo.Ext
			//set chunk param
			setChunkParam(chunkParam, &splitParams, libraryInfo, isTableFile)
			splitParams.SemanticChunkModelConfigId = cast.ToInt(libraryInfo[`model_config_id`])
			splitParams.SemanticChunkUseModel = libraryInfo[`use_model`]
		case define.QALibraryType:
			if define.IsTableFile(uploadInfo.Ext) || define.IsDocxFile(uploadInfo.Ext) {
				autoSplit = true
			}
		case define.OpenLibraryType:
			splitParams.ChunkType = define.ChunkTypeNormal
			if define.IsMdFile(uploadInfo.Ext) {
				autoSplit = true
			}
		case define.OfficialLibraryType:
			splitParams.ChunkType = define.ChunkTypeNormal
			if define.IsMdFile(uploadInfo.Ext) {
				autoSplit = true
			}
		}
		if define.IsPdfFile(uploadInfo.Ext) && (pdfParseType == define.PdfParseTypeOcr || pdfParseType == define.PdfParseTypeOcrWithImage) {
			_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
				`async_split_params`: tool.JsonEncodeNoError(splitParams),
				`update_time`:        tool.Time2Int(),
			})
			if err != nil {
				logs.Error(err.Error())
			}
			lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: cast.ToInt(fileId)})
			if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `file_url`: uploadInfo.Link}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(define.ConvertHtmlTopic, message); err != nil {
				logs.Error(err.Error())
			}
			continue
		}
		if define.IsPdfFile(uploadInfo.Ext) && pdfParseType == define.PdfParseTypeOcrAli {
			jobId, err := common.SubmitOdcParserJob(lang, userId, uploadInfo.Link)
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
					`async_split_params`: tool.JsonEncodeNoError(splitParams),
					`ali_ocr_job_id`:     jobId,
					`update_time`:        tool.Time2Int(),
				})
				if err != nil {
					logs.Error(err.Error())
				}
				lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: cast.ToInt(fileId)})
			}
			continue
		}
		if autoSplit {
			go common.AutoSplitLibFile(userId, int(fileId), splitParams)
		}
	}
	return fileIds, nil
}

func setChunkParam(chunkParam *define.ChunkParam, splitParams *define.SplitParams, libraryInfo msql.Params,
	isTableFile bool) {
	if chunkParam == nil || cast.ToInt(chunkParam.SetChunk) == 0 {
		splitParams.ChunkType = cast.ToInt(libraryInfo[`chunk_type`])
		if splitParams.ChunkType == define.ChunkTypeNormal {
			splitParams.SeparatorsNo = cast.ToString(libraryInfo[`normal_chunk_default_separators_no`])
			splitParams.ChunkSize = cast.ToInt(libraryInfo[`normal_chunk_default_chunk_size`])
			splitParams.ChunkOverlap = cast.ToInt(libraryInfo[`normal_chunk_default_chunk_overlap`])
			splitParams.NotMergedText = cast.ToBool(libraryInfo[`normal_chunk_default_not_merged_text`])
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
		} else if splitParams.ChunkType == define.ChunkTypeFatherSon {
			splitParams.FatherChunkParagraphType = cast.ToInt(libraryInfo[`father_chunk_paragraph_type`])
			splitParams.FatherChunkSeparatorsNo = libraryInfo[`father_chunk_separators_no`]
			splitParams.FatherChunkChunkSize = cast.ToInt(libraryInfo[`father_chunk_chunk_size`])
			splitParams.SonChunkSeparatorsNo = libraryInfo[`son_chunk_separators_no`]
			splitParams.SonChunkChunkSize = cast.ToInt(libraryInfo[`son_chunk_chunk_size`])
			splitParams.NotMergedText = true //父子分段不合并较小分段
		}
		return
	}
	//lib file chunk custom
	splitParams.ChunkType = cast.ToInt(chunkParam.ChunkType)
	if splitParams.ChunkType == define.ChunkTypeNormal {
		splitParams.SeparatorsNo = chunkParam.NormalChunkDefaultSeparatorsNo
		splitParams.ChunkSize = cast.ToInt(chunkParam.NormalChunkDefaultChunkSize)
		splitParams.ChunkOverlap = cast.ToInt(chunkParam.NormalChunkDefaultChunkOverlap)
		splitParams.NotMergedText = cast.ToBool(chunkParam.NormalChunkDefaultNotMergedText)
	} else if splitParams.ChunkType == define.ChunkTypeSemantic {
		splitParams.SemanticChunkSize = cast.ToInt(chunkParam.SemanticChunkDefaultChunkSize)
		splitParams.SemanticChunkOverlap = cast.ToInt(chunkParam.SemanticChunkDefaultChunkOverlap)
		splitParams.SemanticChunkThreshold = cast.ToInt(chunkParam.SemanticChunkDefaultThreshold)
	} else if splitParams.ChunkType == define.ChunkTypeAi {
		splitParams.AiChunkModel = cast.ToString(chunkParam.AiChunkModel)
		splitParams.AiChunkModelConfigId = cast.ToInt(chunkParam.AiChunkModelConfigId)
		splitParams.AiChunkPrumpt = cast.ToString(chunkParam.AiChunkPrumpt)
		splitParams.AiChunkSize = cast.ToInt(chunkParam.AiChunkSize)
		splitParams.AiChunkNew = !isTableFile
		splitParams.AiChunkTaskId = uuid.New().String()
	} else if splitParams.ChunkType == define.ChunkTypeFatherSon {
		splitParams.FatherChunkParagraphType = cast.ToInt(chunkParam.FatherChunkParagraphType)
		splitParams.FatherChunkSeparatorsNo = chunkParam.FatherChunkSeparatorsNo
		splitParams.FatherChunkChunkSize = cast.ToInt(chunkParam.FatherChunkChunkSize)
		splitParams.SonChunkSeparatorsNo = chunkParam.SonChunkSeparatorsNo
		splitParams.SonChunkChunkSize = cast.ToInt(chunkParam.SonChunkChunkSize)
		splitParams.NotMergedText = true //父子分段不合并较小分段
	}
}

func AddLibraryFile(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeAddLibraryFileReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.DocType = c.DefaultPostForm(`doc_type`, cast.ToString(define.DocTypeLocal))
	data, httpStatus, err := BridgeAddLibraryFile(adminUserId, userId, common.GetLang(c), &req, nil, c)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func DelLibraryFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	ids := cast.ToString(c.PostForm(`id`))

	err := BridgeDelLibraryFile(userId, ids, common.GetLang(c))
	c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	return
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
	info[`normal_chunk_default_not_merged_text`] = library[`normal_chunk_default_not_merged_text`]
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

	separators, _ := common.GetSeparatorsByNo(info[`separators_no`], common.GetLang(c))
	info[`separators`] = strings.Join(separators, ", ")

	// ====== 元数据及其值 ======
	// 1) 内置元数据值：固定从 chat_ai_library_file 字段映射
	groupId := cast.ToInt(info[`group_id`])
	groupName := `未分组`
	if groupId > 0 {
		groupInfo, err := msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`library_id`, cast.ToString(cast.ToInt(info[`library_id`]))).
			Where(`group_type`, cast.ToString(define.LibraryGroupTypeFile)).
			Where(`id`, cast.ToString(groupId)).
			Field(`group_name`).
			Find()
		if err == nil && len(groupInfo[`group_name`]) > 0 {
			groupName = groupInfo[`group_name`]
		}
	}
	info[`group_name`] = groupName

	// 2) 自定义元数据值：来自 chat_ai_library_file.metadata(jsonb)
	metaMap := make(map[string]any)
	metaStr := strings.TrimSpace(cast.ToString(info[`metadata`]))
	if metaStr == `` {
		metaStr = `{}` // 容错
	}
	_ = tool.JsonDecode(metaStr, &metaMap)

	// 3) 组装元数据表格行（前端直接渲染）
	builtinValueMap := map[string]any{
		define.BuiltinMetaKeyUpdateTime: cast.ToInt(info[`update_time`]),
		define.BuiltinMetaKeyCreateTime: cast.ToInt(info[`create_time`]),
		define.BuiltinMetaKeySource:     cast.ToInt(info[`doc_type`]),
		define.BuiltinMetaKeyGroup:      groupName,
	}

	// 自定义：只返回已配置的 schema key（避免返回无效/脏key）
	schemaList, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`library_id`, cast.ToString(cast.ToInt(info[`library_id`]))).
		Order(`id asc`).
		Select()

	// 4) 生成表格行：数据名/数据类型/数据值
	metaList := make([]map[string]any, 0, len(define.BuiltinMetaSchemaList)+len(schemaList))
	// 内置：数据值来自文件字段；is_show 来自 chat_ai_library 的 show_meta_*
	for _, b := range define.BuiltinMetaSchemaList {
		isShow := 0
		switch b.Key {
		case define.BuiltinMetaKeySource:
			isShow = cast.ToInt(library[`show_meta_source`])
		case define.BuiltinMetaKeyUpdateTime:
			isShow = cast.ToInt(library[`show_meta_update_time`])
		case define.BuiltinMetaKeyCreateTime:
			isShow = cast.ToInt(library[`show_meta_create_time`])
		case define.BuiltinMetaKeyGroup:
			isShow = cast.ToInt(library[`show_meta_group`])
		}
		metaList = append(metaList, map[string]any{
			`name`:       b.Name,
			`key`:        b.Key,
			`type`:       b.Type,
			`value`:      builtinValueMap[b.Key],
			`is_show`:    isShow,
			`is_builtin`: 1,
		})
	}

	if err == nil {
		for _, schema := range schemaList {
			k := strings.TrimSpace(schema[`key`])
			if k == `` {
				continue
			}
			val, ok := metaMap[k]
			if !ok {
				val = ``
			}
			metaList = append(metaList, map[string]any{
				`id`:         cast.ToInt(schema[`id`]),
				`library_id`: cast.ToInt(schema[`library_id`]),
				`name`:       schema[`name`],
				`key`:        k,
				`type`:       cast.ToInt(schema[`type`]),
				`value`:      val,
				`is_show`:    cast.ToInt(schema[`is_show`]),
				`is_builtin`: 0,
			})
		}
	}
	// info 是 msql.Params(map[string]string)，这里组装一个 any-map 返回，避免把 map 赋给 string
	resp := make(map[string]any, len(info)+4)
	for k, v := range info {
		resp[k] = v
	}
	resp[`group_name`] = groupName
	resp[`meta_list`] = metaList

	c.String(http.StatusOK, lib_web.FmtJson(resp, nil))
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

func unifyGetLibFileSplit(c *gin.Context, chunkPreview bool) {
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
		NotMergedText:              cast.ToBool(c.Query(`not_merged_text`)),
		FatherChunkParagraphType:   cast.ToInt(c.Query(`father_chunk_paragraph_type`)),
		FatherChunkSeparatorsNo:    strings.TrimSpace(c.Query(`father_chunk_separators_no`)),
		FatherChunkChunkSize:       cast.ToInt(c.Query(`father_chunk_chunk_size`)),
		SonChunkSeparatorsNo:       strings.TrimSpace(c.Query(`son_chunk_separators_no`)),
		SonChunkChunkSize:          cast.ToInt(c.Query(`son_chunk_chunk_size`)),
	}
	if chunkPreview { //预览逻辑
		splitParams.ChunkPreview = true
		splitParams.ChunkPreviewSize = define.SplitPreviewChunkMaxSize
		if define.IsDev {
			splitParams.ChunkPreviewSize = 500
		}
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
	} else if splitParams.ChunkType == define.ChunkTypeAi && (splitParams.AiChunkTaskId == "" || chunkPreview) {
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
	} else if splitParams.ChunkType == define.ChunkTypeFatherSon {
		if !tool.InArrayInt(splitParams.FatherChunkParagraphType, []int{define.FatherChunkParagraphTypeFullText, define.FatherChunkParagraphTypeSection}) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `father_chunk_paragraph_type`))))
			return
		}
		if splitParams.FatherChunkParagraphType != define.FatherChunkParagraphTypeFullText {
			if len(splitParams.FatherChunkSeparatorsNo) == 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `father_chunk_separators_no`))))
				return
			}
			if splitParams.FatherChunkChunkSize < 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `father_chunk_chunk_size`))))
				return
			}
		}
		if len(splitParams.SonChunkSeparatorsNo) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `son_chunk_separators_no`))))
			return
		}
		if splitParams.SonChunkChunkSize < 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `son_chunk_chunk_size`))))
			return
		}
	}

	list, wordTotal, splitParams, err := common.GetLibFileSplit(userId, fileId, pdfPageNum, splitParams, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`split_params`: splitParams, `list`: list, `word_total`: wordTotal}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetLibFileSplit(c *gin.Context) {
	unifyGetLibFileSplit(c, false)
}

func GetLibFileSplitPreview(c *gin.Context) {
	unifyGetLibFileSplit(c, true)
}

func SaveLibFileSplit(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.PostForm(`id`))
	wordTotal := cast.ToInt(c.PostForm(`word_total`))
	splitParams, list := define.SplitParams{}, make(define.DocSplitItems, 0)
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
	if splitParams.ChunkAsync {
		saveLibFileSplitAsync(c)
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
		go func() {
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

func saveLibFileSplitAsync(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.PostForm(`id`))
	wordTotal := cast.ToInt(c.PostForm(`word_total`))
	splitParams, list := define.SplitParams{}, make(define.DocSplitItems, 0)
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
	if splitParams.ChunkType == define.ChunkTypeAi {
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
	if splitParams.ChunkAsync {
		// 异步保存
		go func() {
			if err = common.UpdateLibFileData(userId, fileId, msql.Datas{`status`: define.FileStatusChunking}); err != nil {
				logs.Error(err.Error())
				return
			}
			list, wordTotal, splitParams, err = common.GetLibFileSplit(userId, fileId, pdfPageNum, splitParams, common.GetLang(c))
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, err))
				return
			}
			if splitParams.ChunkType == define.ChunkTypeAi {
				if err = common.SaveAISplitDocs(userId, fileId, wordTotal, qaIndexType, splitParams, list, pdfPageNum, common.GetLang(c)); err != nil {
					logs.Error(err.Error())
					return
				}
			} else {
				_, err = common.SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType, splitParams, list, pdfPageNum, common.GetLang(c))
				if err != nil {
					logs.Error(err.Error())
					return
				}
			}
		}()
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
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
		Where(`f.delete_time`, `0`).
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
			`status`:             cast.ToString(define.FileStatusInitial),
			`pdf_parse_type`:     pdfParseType,
			`ocr_pdf_index`:      0,
			`ocr_pdf_total`:      page,
			`async_split_params`: ``, //清空之前的参数
		})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	if pdfParseType == define.PdfParseTypeOcrAli {
		jobId, err := common.SubmitOdcParserJob(common.GetLang(c), userId, info["file_url"])
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
		Where(`delete_time`, `0`).
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
		Where(`delete_time`, `0`).
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
		Where(`delete_time`, `0`).
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

func CreateExportLibFileTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 类型判断
	exportType := cast.ToInt(c.Query("export_type")) // 1: upload file 2:qa docs
	libraryId := cast.ToUint(c.Query("library_id"))
	if !tool.InArrayInt(exportType, []int{define.ExportLibFileUpload, define.ExportQALibDocs}) {
		common.FmtError(c, `param_invalid`, `export_type`)
		return
	}
	params := map[string]any{
		`admin_user_id`: adminUserId,
		`file_id`:       strings.TrimSpace(c.Query("file_id")),
		`library_id`:    libraryId,
		`data_ids`:      strings.TrimSpace(c.Query("data_ids")),
		`group_id`:      strings.TrimSpace(c.Query("group_id")),
		`export_type`:   exportType,
	}
	id, err := common.CreateExportTask(uint(adminUserId), 0, define.ExportSourceLibFileDoc, ``, params)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, id)
}

func DownloadLibraryFile(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	fileId := cast.ToInt(c.Query(`id`))
	if fileId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	fileInfo, _ := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(fileId)).
		Find()
	if len(fileInfo) == 0 || len(fileInfo[`file_url`]) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// 导出上传文件
	filePath := fileInfo[`file_url`]
	fileName := fileInfo[`file_name`]
	if !strings.HasSuffix(fileName, fileInfo[`file_ext`]) {
		fileName = fileName + `.` + fileInfo[`file_ext`]
	}
	if !common.LinkExists(filePath) {
		common.FmtError(c, `no_data`)
		return
	}
	c.FileAttachment(common.GetFileByLink(filePath), fileName)
}

func SaveMetadata(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	fileId := cast.ToInt(c.PostForm(`file_id`))
	if fileId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}

	// 批量保存：list=[{"key":"key_1","value":"xxx"},...]
	listRaw := strings.TrimSpace(c.PostForm(`list`))
	if listRaw == `` || listRaw == `[]` || listRaw == `null` || listRaw == `{}` {
		common.FmtError(c, `param_lack`)
		return
	}
	type MetaKV struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}
	reqList := make([]MetaKV, 0)
	if err := tool.JsonDecodeUseNumber(listRaw, &reqList); err != nil || len(reqList) == 0 {
		common.FmtError(c, `param_err`)
		return
	}
	updateMap := make(map[string]any, len(reqList))
	keySet := make(map[string]struct{}, len(reqList))
	for _, it := range reqList {
		k := strings.TrimSpace(it.Key)
		if k == `` {
			continue
		}
		// 内置元数据固定值：不允许修改（直接忽略）
		if define.IsBuiltinMetaKey(k) {
			continue
		}
		// 统一转字符串存储（与过滤逻辑一致：metadata->>key 再 cast）
		updateMap[k] = cast.ToString(it.Value)
		keySet[k] = struct{}{}
	}
	if len(updateMap) == 0 {
		common.FmtOk(c, nil)
		return
	}

	fileInfo, err := msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(fileId)).
		Where(`delete_time`, `0`).
		Field(`id,library_id,metadata`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(fileInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	libraryId := cast.ToInt(fileInfo[`library_id`])

	// 校验所有 key 必须存在于该知识库的自定义 meta schema 中
	keyList := make([]string, 0, len(keySet))
	for k := range keySet {
		keyList = append(keyList, k)
	}
	okKeys, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Where(`key`, `in`, strings.Join(keyList, `,`)).
		ColumnArr(`key`)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	okSet := make(map[string]struct{}, len(okKeys))
	for _, k := range okKeys {
		okSet[cast.ToString(k)] = struct{}{}
	}
	for k := range keySet {
		if _, ok := okSet[k]; !ok {
			common.FmtError(c, `param_err`)
			return
		}
	}

	if len(okSet) == 0 {
		common.FmtError(c, `param_err`)
		return
	}

	metaMap := make(map[string]any)
	metaStr := strings.TrimSpace(cast.ToString(fileInfo[`metadata`]))
	if metaStr == `` {
		metaStr = `{}`
	}
	_ = tool.JsonDecode(metaStr, &metaMap)
	for k, v := range updateMap {
		metaMap[k] = v
	}
	metaJson, err := tool.JsonEncode(metaMap)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(fileId)).
		Update(msql.Datas{`metadata`: metaJson, `update_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	// clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})

	common.FmtOk(c, nil)
}
