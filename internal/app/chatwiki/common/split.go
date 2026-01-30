// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/textsplitter"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/sse"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
	"golang.org/x/image/webp"
)

func GetLibFileSplit(userId, fileId, pdfPageNum int, splitParams define.SplitParams, lang string) (list define.DocSplitItems, wordTotal int, _splitParams define.SplitParams, err error) {
	info, err := GetLibFileInfo(fileId, userId)
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	if len(info) == 0 {
		err = errors.New(i18n.Show(lang, `no_data`))
		return
	}
	if !tool.InArrayInt(cast.ToInt(info[`status`]), []int{define.FileStatusWaitSplit, define.FileStatusLearned}) {
		err = errors.New(i18n.Show(lang, `status_exception`))
		return
	}
	splitParams.FileExt = info[`file_ext`]
	library, err := GetLibraryInfo(cast.ToInt(info[`library_id`]), userId)
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	if len(library) == 0 {
		err = errors.New(i18n.Show(lang, `no_data`))
		return
	}
	splitParams.IsTableFile = cast.ToInt(info[`is_table_file`])
	splitParams, err = CheckSplitParams(library, splitParams, lang)
	if err != nil {
		return
	}
	if splitParams.PdfParseType == 0 {
		splitParams.PdfParseType = cast.ToInt(info[`pdf_parse_type`])
		if splitParams.PdfParseType == 0 {
			splitParams.PdfParseType = define.PdfParseTypeText
		}
	}
	if (cast.ToInt(library[`type`]) == define.GeneralLibraryType && splitParams.IsQaDoc == define.DocTypeQa) ||
		(cast.ToInt(library[`type`]) == define.QALibraryType && splitParams.IsQaDoc != define.DocTypeQa) {
		err = errors.New(i18n.Show(lang, `param_invalid`, `is_qa_doc`))
		return
	}
	if cast.ToInt(info[`is_table_file`]) == define.FileIsTable && splitParams.IsQaDoc == define.DocTypeQa {
		list, wordTotal, err = ReadQaTab(info[`file_url`], info[`file_ext`], splitParams)
	} else if cast.ToInt(info[`is_table_file`]) == define.FileIsTable && splitParams.IsQaDoc != define.DocTypeQa {
		list, wordTotal, err = ReadTab(info[`file_url`], info[`file_ext`])
	} else if define.IsDocxFile(info[`file_ext`]) {
		list, wordTotal, err = ReadDocx(info[`file_url`], userId)
	} else if define.IsOfdFile(info[`file_ext`]) {
		list, wordTotal, err = ReadOfd(info[`file_url`], userId)
	} else if define.IsTxtFile(info[`file_ext`]) || define.IsMdFile(info[`file_ext`]) {
		list, wordTotal, err = ReadTxt(info[`file_url`])
	} else if define.IsHtmlFile(info[`file_ext`]) {
		list, wordTotal, err = ReadHtmlContent(info[`file_url`], userId)
	} else if define.IsPdfFile(info[`file_ext`]) && splitParams.PdfParseType == define.PdfParseTypeText {
		list, wordTotal, err = ReadPdf(info[`file_url`], pdfPageNum, lang)
	} else if define.IsPdfFile(info[`file_ext`]) && (splitParams.PdfParseType == define.PdfParseTypeOcr || splitParams.PdfParseType == define.PdfParseTypeOcrWithImage) && pdfPageNum > 0 {
		list, wordTotal, err = OcrReadOnePagePdf(userId, splitParams.PdfParseType, info[`file_url`], pdfPageNum, lang)
	} else {
		if len(info[`html_url`]) == 0 || !LinkExists(info[`html_url`]) { //compatible with old data
			list, wordTotal, err = ConvertAndReadHtmlContent(cast.ToInt(info[`id`]), info[`file_url`], userId, splitParams.PdfParseType)
		} else {
			list, wordTotal, err = ReadHtmlContent(info[`html_url`], userId)
		}
	}

	if err != nil {
		return
	}
	if len(list) == 0 || wordTotal == 0 {
		err = errors.New(i18n.Show(lang, `doc_empty`))
		return
	}

	//fix: pq: invalid byte sequence for encoding "UTF8": 0x00
	regReplace := regexp.MustCompile(`[\x00-\x08\x0E-\x1F\x7F]`)
	for i, item := range list {
		list[i].Title = regReplace.ReplaceAllString(item.Title, ``)
		list[i].Content = regReplace.ReplaceAllString(item.Content, ``)
		list[i].Question = regReplace.ReplaceAllString(item.Question, ``)
		for j, similar := range list[i].SimilarQuestionList {
			list[i].SimilarQuestionList[j] = regReplace.ReplaceAllString(similar, ``)
		}
		list[i].Answer = regReplace.ReplaceAllString(item.Answer, ``)
	}

	//fix:sql err
	quotesReplacer := strings.NewReplacer(`'`, `"`)
	for i, item := range list {
		list[i].Title = quotesReplacer.Replace(item.Title)
		list[i].Content = quotesReplacer.Replace(item.Content)
		list[i].Question = quotesReplacer.Replace(item.Question)
		for j, similar := range list[i].SimilarQuestionList {
			list[i].SimilarQuestionList[j] = quotesReplacer.Replace(similar)
		}
		list[i].Answer = quotesReplacer.Replace(item.Answer)
	}

	// split by document type
	if splitParams.IsQaDoc == define.DocTypeQa {
		if cast.ToInt(info[`is_table_file`]) != define.FileIsTable {
			list = QaDocSplit(splitParams, list)
		}
	} else {
		list, err = MultDocSplit(cast.ToInt(info[`admin_user_id`]), fileId, pdfPageNum, splitParams, list)
	}
	if err != nil {
		return
	}
	list.UnifySetNumber() //对number统一编号
	for i := range list {
		if splitParams.IsQaDoc == define.DocTypeQa {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question) + utf8.RuneCountInString(list[i].Answer)
		} else {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
		}
	}
	_splitParams = splitParams

	return
}

func SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType int, splitParams define.SplitParams, list define.DocSplitItems, pdfPageNum int, lang string) ([]int64, error) {
	info, err := GetLibFileInfo(fileId, userId)
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return []int64{}, errors.New(i18n.Show(lang, `no_data`))
	}
	if !tool.InArrayInt(cast.ToInt(info[`status`]), []int{define.FileStatusWaitSplit, define.FileStatusChunking, define.FileStatusLearned}) {
		return []int64{}, errors.New(i18n.Show(lang, `status_exception`))
	}

	//check params
	if splitParams.IsQaDoc == define.DocTypeQa { // qa
		for i := range list {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question + list[i].Answer)
			if utf8.RuneCountInString(list[i].Question) < 1 || utf8.RuneCountInString(list[i].Question) > MaxContent {
				return []int64{}, errors.New(i18n.Show(lang, `length_err`, i+1))
			}
			if utf8.RuneCountInString(list[i].Answer) < 1 || utf8.RuneCountInString(list[i].Answer) > MaxContent {
				return []int64{}, errors.New(i18n.Show(lang, `length_err`, i+1))
			}
		}
	} else {
		for i := range list {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
			if list[i].WordTotal < 1 || list[i].WordTotal > MaxContent {
				return []int64{}, errors.New(i18n.Show(lang, `length_err`, i+1))
			}
		}
	}
	list.UnifySetNumber() //对number统一编号

	if splitParams.IsQaDoc == define.DocTypeQa {
		if qaIndexType != define.QAIndexTypeQuestionAndAnswer && qaIndexType != define.QAIndexTypeQuestion {
			return []int64{}, errors.New(i18n.Show(lang, `param_invalid`, `qa_index_type`))
		}
	}
	//add lock dispose
	if !lib_redis.AddLock(define.Redis, define.LockPreKey+`SaveLibFileSplit`+cast.ToString(fileId), time.Minute*5) {
		err = errors.New(i18n.Show(lang, `op_lock`))
	}
	defer func(fileId int) {
		lib_redis.UnLock(define.Redis, define.LockPreKey+`SaveLibFileSplit`+cast.ToString(fileId))
	}(fileId)

	//database dispose
	m := msql.Model(`chat_ai_library_file`, define.Postgres)
	err = m.Begin()
	defer func() {
		_ = m.Rollback()
	}()
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}

	//database dispose
	shouldDeleteQuery := msql.Model("chat_ai_library_file_data", define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`file_id`, cast.ToString(fileId))
	if pdfPageNum > 0 {
		shouldDeleteQuery.Where(`page_num`, cast.ToString(pdfPageNum))
	}
	shouldDeleteIds, err := shouldDeleteQuery.Where("category_id", "0").ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}

	shouldIsolatedQuery := msql.Model("chat_ai_library_file_data", define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`file_id`, cast.ToString(fileId))
	if pdfPageNum > 0 {
		shouldIsolatedQuery.Where(`page_num`, cast.ToString(pdfPageNum))
	}
	shouldIsolatedIds, err := shouldIsolatedQuery.Where("category_id", ">", "0").ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(shouldDeleteIds) > 0 {
		_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`id`, `in`, strings.Join(shouldDeleteIds, `,`)).
			Delete()
		if err != nil {
			logs.Error(err.Error())
			return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
		}
		_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
			Where(`data_id`, `in`, strings.Join(shouldDeleteIds, `,`)).
			Delete()
		if err != nil {
			logs.Error(err.Error())
			return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	if len(shouldIsolatedIds) > 0 {
		_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`id`, `in`, strings.Join(shouldIsolatedIds, `,`)).
			Update(msql.Datas{"isolated": true})
		if err != nil {
			logs.Error(err.Error())
			return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
		}
	}

	var (
		aiChunkErrMsg = ""
		aiChunkErrNum = 0
		indexIds      []int64
		dataIds       []int64
		library, _    = GetLibraryData(cast.ToInt(info[`library_id`]))
		skipUseModel  = cast.ToInt(library[`type`]) == define.OpenLibraryType && cast.ToInt(library[`use_model_switch`]) != define.SwitchOn
	)
	for i, item := range list {
		if utf8.RuneCountInString(item.Content) > MaxContent || utf8.RuneCountInString(item.Question) > MaxContent || utf8.RuneCountInString(item.Answer) > MaxContent {
			return []int64{}, errors.New(i18n.Show(lang, `length_err`, i+1))
		}

		data := msql.Datas{
			`admin_user_id`: info[`admin_user_id`],
			`library_id`:    info[`library_id`],
			`file_id`:       fileId,
			`number`:        item.Number,
			`page_num`:      item.PageNum,
			`title`:         item.Title,
			`word_total`:    item.WordTotal,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
			//父子分段
			`father_chunk_paragraph_number`: item.FatherChunkParagraphNumber,
		}

		// 如果知识库的的类型是问答知识库
		// 则把知识库文件中的group_id塞到分段表的group_id
		if cast.ToInt(library[`type`]) == define.QALibraryType {
			data[`group_id`] = info[`group_id`]
		}

		if splitParams.IsQaDoc == define.DocTypeQa {
			if splitParams.IsTableFile == define.FileIsTable {
				data[`type`] = define.ParagraphTypeExcelQA
			} else {
				data[`type`] = define.ParagraphTypeDocQA
			}
			data[`question`] = strings.TrimSpace(item.Question)
			data[`answer`] = strings.TrimSpace(item.Answer)
			similarQuestions, err := tool.JsonEncode(item.SimilarQuestionList)
			if err != nil {
				logs.Error(err.Error())
			} else {
				data[`similar_questions`] = similarQuestions
			}
			if len(item.Images) > 0 {
				jsonImages, err := CheckLibraryImage(item.Images)
				if err != nil {
					_ = m.Rollback()
					return []int64{}, errors.New(i18n.Show(lang, `param_invalid`, `images`))
				}
				data[`images`] = jsonImages
			} else {
				data[`images`] = `[]`
			}

			// 仅对问答知识库处理
			if cast.ToInt(library[`type`]) == define.QALibraryType {
				// 在这里检测，如果问题相同，就删除原有问答再新增新问答
				// 不过有一个问题，因为这个方法是通过goroutine方式运行，有一个极限的情况
				// 当用户快速上传多个文件，而多个文件中存在相同问题，则可能还是会出现问题相同，答案不同的情况
				qaQuestionExistQuery := msql.Model("chat_ai_library_file_data", define.Postgres).
					Where(`admin_user_id`, cast.ToString(userId)).
					Where(`library_id`, cast.ToString(info[`library_id`])).
					Where(`question`, cast.ToString(data[`question`]))
				qaQuestionExistIds, err := qaQuestionExistQuery.ColumnArr(`id`)
				if err != nil {
					logs.Error(err.Error())
					return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
				}

				// 将重复的问答直接删除问答和向量索引
				if len(qaQuestionExistIds) > 0 {
					_, err = msql.Model("chat_ai_library_file_data", define.Postgres).
						Where(`admin_user_id`, cast.ToString(userId)).
						Where(`library_id`, cast.ToString(info[`library_id`])).
						Where(`id`, `in`, strings.Join(qaQuestionExistIds, `,`)).
						Delete()
					if err != nil {
						logs.Error(err.Error())
						return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
					}
					_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
						Where(`admin_user_id`, cast.ToString(userId)).
						Where(`library_id`, cast.ToString(info[`library_id`])).
						Where(`data_id`, `in`, strings.Join(qaQuestionExistIds, `,`)).
						Delete()
					if err != nil {
						logs.Error(err.Error())
						return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
					}
				}
			}

			// 实际添加分段文件
			id, err := msql.Model("chat_ai_library_file_data", define.Postgres).Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
			dataIds = append(dataIds, id)
			vectorID, err := SaveVector(
				cast.ToInt64(info[`admin_user_id`]),
				cast.ToInt64(info[`library_id`]),
				cast.ToInt64(fileId),
				id,
				cast.ToString(define.VectorTypeQuestion),
				strings.TrimSpace(item.Question),
			)
			if err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
			indexIds = append(indexIds, vectorID)
			if err = DeleteLibraryFileDataIndex(cast.ToString(id), cast.ToString(define.VectorTypeSimilarQuestion)); err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
			for _, similarQuestion := range item.SimilarQuestionList {
				if len(similarQuestion) <= 0 {
					continue
				}
				vectorID, err := SaveVector(
					cast.ToInt64(info[`admin_user_id`]),
					cast.ToInt64(info[`library_id`]),
					cast.ToInt64(fileId),
					id,
					cast.ToString(define.VectorTypeSimilarQuestion),
					strings.TrimSpace(similarQuestion),
				)
				if err != nil {
					logs.Error(err.Error())
					return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
				}
				indexIds = append(indexIds, vectorID)
			}

			if qaIndexType == define.QAIndexTypeQuestionAndAnswer {
				vectorID, err = SaveVector(
					cast.ToInt64(info[`admin_user_id`]),
					cast.ToInt64(info[`library_id`]),
					cast.ToInt64(fileId),
					id,
					cast.ToString(define.VectorTypeAnswer),
					strings.TrimSpace(item.Answer),
				)
				if err != nil {
					logs.Error(err.Error())
					return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
				}
				indexIds = append(indexIds, vectorID)
			}
		} else {
			data[`type`] = define.ParagraphTypeNormal
			data[`content`] = strings.TrimSpace(item.Content)
			if item.AiChunkErrMsg != "" {
				data[`split_err_msg`] = item.AiChunkErrMsg
				data[`split_status`] = define.SplitStatusException
				aiChunkErrMsg = item.AiChunkErrMsg
				aiChunkErrNum++
			}
			if len(item.Images) > 0 {
				jsonImages, err := CheckLibraryImage(item.Images)
				if err != nil {
					return []int64{}, errors.New(i18n.Show(lang, `param_invalid`, `images`))
				}
				data[`images`] = jsonImages
			} else {
				data[`images`] = `[]`
			}

			id, err := msql.Model("chat_ai_library_file_data", define.Postgres).Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
			dataIds = append(dataIds, id)
			vectorID, err := SaveVector(
				cast.ToInt64(info[`admin_user_id`]),
				cast.ToInt64(info[`library_id`]),
				cast.ToInt64(fileId),
				id,
				cast.ToString(define.VectorTypeParagraph),
				strings.TrimSpace(item.Content),
			)
			if err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
			if item.AiChunkErrMsg != "" {
				continue
			}
			indexIds = append(indexIds, vectorID)
		}
	}

	status := define.FileStatusLearning
	errmsg := `success`
	if len(list) <= 0 {
		status = define.FileStatusException
		errmsg = i18n.Show(lang, `doc_empty`)
	}
	if aiChunkErrNum == len(indexIds) {
		status = define.FileStatusException
		errmsg = aiChunkErrMsg
	} else if aiChunkErrNum > 0 {
		status = define.FileStatusLearned
	}
	data := msql.Datas{
		`status`:                         status,
		`errmsg`:                         errmsg,
		`word_total`:                     wordTotal,
		`split_total`:                    len(list),
		`is_qa_doc`:                      splitParams.IsQaDoc,
		`is_diy_split`:                   splitParams.IsDiySplit,
		`separators_no`:                  splitParams.SeparatorsNo,
		`chunk_size`:                     splitParams.ChunkSize,
		`chunk_overlap`:                  splitParams.ChunkOverlap,
		`question_lable`:                 splitParams.QuestionLable,
		`similar_label`:                  splitParams.SimilarLabel,
		`answer_lable`:                   splitParams.AnswerLable,
		`question_column`:                splitParams.QuestionColumn,
		`similar_column`:                 splitParams.SimilarColumn,
		`answer_column`:                  splitParams.AnswerColumn,
		`enable_extract_image`:           splitParams.EnableExtractImage,
		`chunk_type`:                     splitParams.ChunkType,
		`semantic_chunk_size`:            splitParams.SemanticChunkSize,
		`semantic_chunk_overlap`:         splitParams.SemanticChunkOverlap,
		`semantic_chunk_threshold`:       splitParams.SemanticChunkThreshold,
		`semantic_chunk_use_model`:       splitParams.SemanticChunkUseModel,
		`semantic_chunk_model_config_id`: splitParams.SemanticChunkModelConfigId,
		`ai_chunk_model`:                 splitParams.AiChunkModel,
		`ai_chunk_model_config_id`:       splitParams.AiChunkModelConfigId,
		`ai_chunk_size`:                  splitParams.AiChunkSize,
		`ai_chunk_prumpt`:                splitParams.AiChunkPrumpt,
		`ai_chunk_task_id`:               splitParams.AiChunkTaskId,
		`pdf_parse_type`:                 splitParams.PdfParseType,
		`update_time`:                    tool.Time2Int(),
		`not_merged_text`:                splitParams.NotMergedText,
		`father_chunk_paragraph_type`:    splitParams.FatherChunkParagraphType,
		`father_chunk_separators_no`:     splitParams.FatherChunkSeparatorsNo,
		`father_chunk_chunk_size`:        splitParams.FatherChunkChunkSize,
		`son_chunk_separators_no`:        splitParams.SonChunkSeparatorsNo,
		`son_chunk_chunk_size`:           splitParams.SonChunkChunkSize,
	}
	if qaIndexType != 0 {
		data[`qa_index_type`] = qaIndexType
	}
	_, err = m.Where(`id`, cast.ToString(fileId)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(list) <= 0 {
		err = m.Commit()
		lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})
		return []int64{}, errors.New(i18n.Show(lang, `doc_empty`))
	}
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})

	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
	}

	if skipUseModel {
		return []int64{}, err
	}
	//async task:convert vector
	for _, id := range indexIds {
		message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if err = AddJobs(define.ConvertVectorTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}

	if GetNeo4jStatus(userId) {
		err = NewGraphDB(userId).DeleteByFile(fileId)
		if err != nil {
			logs.Error(err.Error())
			return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
		}
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, cast.ToString(fileId)).
			Update(msql.Datas{`graph_status`: define.GraphStatusNotStart})
		if err != nil {
			logs.Error(err.Error())
			return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
		}
		lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})
		if cast.ToInt(info[`graph_status`]) == define.GraphStatusConverted {
			if err = ConstructGraph(cast.ToInt(info[`id`])); err != nil {
				logs.Error(err.Error())
				return []int64{}, errors.New(i18n.Show(lang, `sys_err`))
			}
		}
	}

	return dataIds, nil
}

func RecursiveCharacter(items define.DocSplitItems, separators []string, chunkSize, chunkOverlap int, notMergedText bool) (define.DocSplitItems, error) {
	split := textsplitter.NewRecursiveCharacter()
	split.Separators = append(separators, ``)
	split.ChunkSize = chunkSize
	split.ChunkOverlap = chunkOverlap
	split.NotMergedText = notMergedText
	split.LenFunc = func(s string) int {
		//将内容包含图片的部分统计为字符-(长度1),这里不能为长度0,最后一段会被丢弃
		return utf8.RuneCountInString(regexp.MustCompile(`\{\{!!(.+?)!!}}`).ReplaceAllString(s, `-`))
	}
	list := make(define.DocSplitItems, 0)
	for _, item := range items {
		//内容包含图片:前面追加一个第一顺位的切分字符,避免图片出现在两段里面
		contents, err := split.SplitText(strings.ReplaceAll(item.Content, `{{!!`, split.Separators[0]+`{{!!`))
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		for _, content := range contents {
			if len(content) == 0 {
				continue
			}
			//将追加的第一顺位的切分字符剔除掉
			content = strings.ReplaceAll(content, split.Separators[0]+`{{!!`, `{{!!`)
			content, images := ExtractTextImages(content)
			if len(content) == 0 {
				continue
			}
			list = append(list, define.DocSplitItem{
				PageNum: item.PageNum, Content: content, Images: images,
				FatherChunkParagraphNumber: item.FatherChunkParagraphNumber,
			})
		}
	}
	return list, nil
}

func GetSeparatorsByNo(separatorsNo, lang string) ([]string, error) {
	separators := make([]string, 0)
	if len(separatorsNo) == 0 {
		return separators, nil
	}
	list := make([]any, 0)
	if err := tool.JsonDecodeUseNumber(separatorsNo, &list); err != nil {
		for _, no := range strings.Split(separatorsNo, `,`) {
			list = append(list, cast.ToInt(no))
		}
	}
	replacer := strings.NewReplacer(`\\n`, `\n`, `\\r`, `\r`, `\\t`, `\t`, `\n`, "\n", `\r`, "\r", `\t`, "\t")
	for i, no := range list {
		//自定义分段标识符
		if noStr, ok := no.(string); ok {
			separators = append(separators, replacer.Replace(noStr))
			continue
		}
		//系统预设分段标识符
		no := cast.ToInt(no)
		if no < 1 || no > len(define.SeparatorsList) {
			return separators, errors.New(i18n.Show(lang, `param_invalid`, `separators_no.`+cast.ToString(i)))
		}
		code := define.SeparatorsList[no-1][`code`]
		if realCode, ok := code.([]string); ok {
			separators = append(separators, realCode...)
		} else {
			separators = append(separators, cast.ToString(code))
		}
	}
	return separators, nil
}

func MergeDocSplitItems(items define.DocSplitItems) define.DocSplitItems {
	if len(items) <= 1 {
		return items
	}
	contents := make([]string, len(items))
	for i, item := range items {
		contents[i] = item.Content
	}
	items[0].Content = strings.Join(contents, "\r\n")
	return define.DocSplitItems{items[0]}
}

func MultDocSplit(adminUserId, fileId, pdfPageNum int, splitParams define.SplitParams, items define.DocSplitItems) (define.DocSplitItems, error) {
	var (
		err         error
		previewText string
		separators  []string
	)
	if splitParams.ChunkType == define.ChunkTypeNormal {
		separators = append(splitParams.Separators, textsplitter.DefaultOptions().Separators...) //补充默认分隔符
		return RecursiveCharacter(items, separators, splitParams.ChunkSize, splitParams.ChunkOverlap, splitParams.NotMergedText)
	} else if splitParams.ChunkType == define.ChunkTypeFatherSon {
		if cast.ToBool(splitParams.IsTableFile) { //表格按分段特殊处理
			for i := range items { //标记父编号
				items[i].FatherChunkParagraphNumber = i + 1
			}
			return RecursiveCharacter(items, textsplitter.DefaultOptions().Separators,
				textsplitter.DefaultOptions().ChunkSize, textsplitter.DefaultOptions().ChunkOverlap, false)
		}
		items = MergeDocSplitItems(items) //pdf等文件先合并成一个
		if splitParams.FatherChunkParagraphType != define.FatherChunkParagraphTypeFullText {
			separators, err = GetSeparatorsByNo(splitParams.FatherChunkSeparatorsNo, ``)
			items, err = RecursiveCharacter(items, separators, splitParams.FatherChunkChunkSize, 0, true)
			if err != nil {
				return items, err
			}
		} else { //出于性能原因,超过10000个标记的文本将被自动截断
			items, err = RecursiveCharacter(items, nil, 10000, 0, true)
			if err != nil {
				return items, err
			}
		}
		for i := range items { //标记父编号
			items[i].FatherChunkParagraphNumber = i + 1
		}
		separators, err = GetSeparatorsByNo(splitParams.SonChunkSeparatorsNo, ``)
		if err != nil {
			return items, err
		}
		return RecursiveCharacter(items, separators, splitParams.SonChunkChunkSize, 0, true)
	} else if splitParams.ChunkType == define.ChunkTypeAi {
		if define.IsTableFile(splitParams.FileExt) {
			return items, nil
		}
		list := make(define.DocSplitItems, 0)
		if splitParams.AiChunkNew {
			var aiSplitItems = make(define.DocSplitItems, len(items))
			copy(aiSplitItems, items)
			go func() {
				list, err = AISplitDocs(cast.ToInt(adminUserId), fileId, splitParams, aiSplitItems)
				if err != nil {
					logs.Error(err.Error())
					return
				}
			}()
			return items, err
		} else {
			// get db data
			data, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
				Where(`file_id`, cast.ToString(fileId)).
				Order(`page_num,father_chunk_paragraph_number,number`).Select()
			if err != nil {
				logs.Error(err.Error())
			}
			for _, item := range data {
				var images []string
				if len(item[`images`]) > 0 {
					_ = tool.JsonDecode(item[`images`], &images)
				}
				if splitParams.PdfParseType >= define.PdfParseTypeText && cast.ToInt(item[`page_num`]) != pdfPageNum && pdfPageNum > 0 {
					continue
				}
				list = append(list, define.DocSplitItem{
					Content:   item[`content`],
					Number:    cast.ToInt(item[`number`]),
					Images:    images,
					Title:     item[`title`],
					PageNum:   cast.ToInt(item[`page_num`]),
					WordTotal: cast.ToInt(item[`word_total`]),
				})
			}
			return list, err
		}
	} else {
		split := NewSemanticSplitterClient()
		split.GoRoutineNum = 5
		split.SemanticChunkSize = splitParams.SemanticChunkSize
		split.SemanticChunkOverlap = splitParams.SemanticChunkOverlap
		split.SemanticChunkThreshold = splitParams.SemanticChunkThreshold
		split.AdminUserId = adminUserId
		split.ModelConfigId = splitParams.SemanticChunkModelConfigId
		split.UseModel = splitParams.SemanticChunkUseModel
		list := make(define.DocSplitItems, 0)
		for _, item := range items {
			contents, err := split.SplitText(item.Content)
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			for _, content := range contents {
				if len(content) == 0 {
					continue
				}
				content, images := ExtractTextImages(content)
				if len(content) == 0 {
					continue
				}
				list = append(list, define.DocSplitItem{PageNum: item.PageNum, Content: content, Images: images})
				previewText += content
				if splitParams.ChunkPreview && utf8.RuneCountInString(previewText) >= splitParams.ChunkPreviewSize {
					return list, err
				}
			}
		}
		return list, err
	}
}

func ConvertAndReadHtmlContent(fileId int, fileUrl string, userId int, pdfParseType int) (define.DocSplitItems, int, error) {
	htmlUrl, err := ConvertHtml(fileId, fileUrl, userId, pdfParseType)
	if err != nil {
		return nil, 0, err
	}

	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`html_url`:    htmlUrl,
		`status`:      define.FileStatusWaitSplit,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		return nil, 0, err
	}
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})

	return ReadHtmlContent(htmlUrl, userId)
}

func RequestConvertService(file, fromFormat string, pdfParseType int) (content string, err error) {
	useOcr := false
	extractImage := false
	if fromFormat == "pdf" {
		useOcr = true
	}
	if pdfParseType == define.PdfParseTypeOcrWithImage {
		extractImage = true
	}
	request := curl.Post(define.Config.WebService[`converter`]+`/convert`).
		PostFile(`file`, file).
		Param(`from_format`, fromFormat).
		Param(`to_format`, `html`).
		Param(`use_ocr`, cast.ToString(useOcr)).
		Param(`extract_images`, cast.ToString(extractImage))
	if resp, err := request.Response(); err != nil {
		return ``, err
	} else if resp.StatusCode != http.StatusOK {
		return ``, errors.New(content)
	}
	return request.String()
}

func PdfConvertHtml(fileId int, fileLink string, pdfParseType int) (content string, err error) {
	page, err := api.PageCountFile(fileLink)
	if err != nil { //获取页码出错
		return
	}
	outDir := define.UploadDir + fmt.Sprintf(`pdf_split/%s`, tool.Random(8)) //随机生成切分后的目录
	defer func(path string) {
		_ = os.RemoveAll(path) //结束后删除目录
	}(outDir)
	_ = tool.MkDirAll(outDir) //确保输出目录存在
	if err = api.SplitFile(fileLink, outDir, 1, nil); err != nil {
		return
	}
	//来自spanFileName,千万不要改,大写会出问题!!!
	filename := strings.TrimSuffix(filepath.Base(fileLink), `.pdf`)
	for idx := 1; idx <= page; idx++ {
		item := fmt.Sprintf(`%s/%s_%d.pdf`, outDir, filename, idx)
		if !tool.IsFile(item) {
			continue //预防文件不存在的情况
		}
		onePage, err := RequestConvertService(item, `pdf`, pdfParseType)
		if err != nil {
			return ``, fmt.Errorf(`[%d/%d]:%v`, idx, page, err)
		}
		info, err := msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, cast.ToString(fileId)).
			Find()
		if err != nil {
			return ``, fmt.Errorf(`get library file error, file id = %d, err = %v`, fileId, err)
		}
		if len(info) == 0 {
			return ``, fmt.Errorf("cannot find library file id = %d", fileId)
		}
		if cast.ToInt(info[`status`]) == define.FileStatusCancelled {
			return ``, fmt.Errorf("pdf parse cancelled, file_id = %d", fileId)
		}
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, cast.ToString(fileId)).
			Update(msql.Datas{"ocr_pdf_index": idx})
		if err != nil {
			return ``, err
		}

		content += onePage
	}
	return
}

func ConvertHtml(fileId int, link string, userId int, pdfParseType int) (content string, err error) {
	if !LinkExists(link) {
		return ``, errors.New(`file not exist:` + link)
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(link), `.`))
	if ext == `pdf` { //切分成每一页再转换合并
		content, err = PdfConvertHtml(fileId, GetFileByLink(link), pdfParseType)
	} else { //直接请求转换服务
		content, err = RequestConvertService(GetFileByLink(link), ext, pdfParseType)
	}
	if err != nil {
		return ``, err
	}
	//替换base64编码的图片
	content, err = ReplaceBase64Img(content, userId)
	if err != nil {
		return ``, err
	}
	//保存html文件
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.html`, userId,
		`convert`, tool.Date(`Ym`), tool.MD5(content))
	url, err := WriteFileByString(objectKey, content)
	if err != nil {
		return ``, err
	}
	return url, nil
}

func ReplaceBase64Img(content string, userId int) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	doc.Find("img").Each(func(index int, item *goquery.Selection) {
		src, exists := item.Attr("src")
		if exists && IsUrl(src) {
			newTag := fmt.Sprintf("<b>{{!!%s!!}}</b>", src)
			item.ReplaceWithHtml(newTag)
		}
		if exists && strings.HasPrefix(src, "data:image") {
			parts := strings.Split(src, ";")
			if len(parts) < 2 {
				logs.Debug(fmt.Sprintf("could not find base64 data"))
				return
			}
			format := strings.TrimPrefix(parts[0], "data:image/")
			if format == "svg+xml" {
				return
			}
			base64Data := strings.TrimPrefix(parts[1], "base64,")
			imgData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			if format == `webp` {
				imgData, err = ConvertWebPToPNG(imgData)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				format = `png`
			}
			// save to file
			objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_image`, tool.Date(`Ym`), tool.MD5(string(imgData)), format)
			imgUrl, err := WriteFileByString(objectKey, string(imgData))
			if err != nil {
				logs.Error(err.Error())
				return
			}
			// Replace img tag with a span tag
			newTag := fmt.Sprintf("<b>{{!!%s!!}}</b>", imgUrl)
			item.ReplaceWithHtml(newTag)
		}
	})
	content, err = doc.Html()
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	if !utf8.ValidString(content) {
		content = tool.Convert(content, `gbk`, `utf-8`)
	}
	return content, nil
}

func ReplaceRemoteImg(content string, userId int) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	doc.Find("img").Each(func(index int, item *goquery.Selection) {
		src, exists := item.Attr("src")
		if exists && (strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
			// 下载远程图片
			resp, err := http.Get(src)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					logs.Error(err.Error())
				}
			}(resp.Body)

			if resp.StatusCode != http.StatusOK {
				logs.Error(fmt.Sprintf("failed to download image: %s, status code: %d", src, resp.StatusCode))
				return
			}

			// 读取图片数据
			imgData, err := io.ReadAll(resp.Body)
			if err != nil {
				logs.Error(err.Error())
				return
			}

			// 确定图片格式
			contentType := resp.Header.Get("Content-Type")
			format := "jpg" // 默认格式
			if strings.Contains(contentType, "png") {
				format = "png"
			} else if strings.Contains(contentType, "gif") {
				format = "gif"
			} else if strings.Contains(contentType, "webp") {
				// 如果是webp格式，转换为PNG
				imgData, err = ConvertWebPToPNG(imgData)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				format = "png"
			}

			// 保存图片到本地
			objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_image`, tool.Date(`Ym`), tool.MD5(string(imgData)), format)
			imgUrl, err := WriteFileByString(objectKey, string(imgData))
			if err != nil {
				logs.Error(err.Error())
				return
			}

			// 替换图片标签
			newTag := fmt.Sprintf("<b>{{!!%s!!}}</b>", imgUrl)
			item.ReplaceWithHtml(newTag)
		}
	})

	content, err = doc.Html()
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}

	if !utf8.ValidString(content) {
		content = tool.Convert(content, `gbk`, `utf-8`)
	}

	return content, nil
}

func ReadHtmlContent(htmlUrl string, userId int) (define.DocSplitItems, int, error) {
	if !LinkExists(htmlUrl) {
		return nil, 0, errors.New(`file not exist:` + htmlUrl)
	}
	content, err := tool.ReadFile(GetFileByLink(htmlUrl))
	if err != nil {
		return nil, 0, err
	}
	return ParseHtmlContent(userId, content)
}

func ParseHtmlContent(userId int, content string) (define.DocSplitItems, int, error) {
	content = strings.ReplaceAll(content, `<!DOCTYPE html>`, ``)
	pages := strings.Split(content, `<meta charset="UTF-8"/>`)
	list := make(define.DocSplitItems, 0)
	wordTotal := 0
	pageNum := 0
	for _, pageContent := range pages {
		pageContent, err := ReplaceBase64Img(pageContent, userId)
		if err != nil {
			logs.Error(err.Error())
			continue
		}

		pageContent = strip.StripTags(pageContent)
		pageContent = strings.TrimSpace(pageContent)
		if len(pageContent) == 0 {
			if pageNum > 0 {
				pageNum += 1
			}
			continue
		}
		pageNum += 1
		list = append(list, define.DocSplitItem{Content: pageContent, PageNum: pageNum, WordTotal: utf8.RuneCountInString(pageContent)})
		wordTotal += utf8.RuneCountInString(pageContent)
	}

	return list, wordTotal, nil
}

func ParseOnePageHtmlContent(userId int, content string, pageNum int) (define.DocSplitItems, int, error) {
	content = strings.ReplaceAll(content, `<!DOCTYPE html>`, ``)
	list := make(define.DocSplitItems, 0)
	wordTotal := 0
	pageContent, err := ReplaceBase64Img(content, userId)
	if err != nil {
		logs.Error(err.Error())
		return nil, wordTotal, err
	}
	pageContent = strip.StripTags(pageContent)
	pageContent = strings.TrimSpace(pageContent)
	if len(pageContent) == 0 {
		return nil, wordTotal, err
	}
	list = append(list, define.DocSplitItem{Content: pageContent, PageNum: pageNum, WordTotal: utf8.RuneCountInString(pageContent)})
	wordTotal += utf8.RuneCountInString(pageContent)

	return list, wordTotal, nil
}

func ParseTabFile(fileUrl, fileExt string) ([][]string, error) {
	if len(fileUrl) == 0 {
		return nil, errors.New(`file link cannot be empty`)
	}
	if !LinkExists(fileUrl) {
		return nil, errors.New(`file not exist:` + fileUrl)
	}
	//read excel
	var rows = make([][]string, 0)
	if fileExt == `csv` {
		content, err := tool.ReadFile(GetFileByLink(fileUrl))
		if err != nil {
			return nil, err
		}
		if !utf8.ValidString(content) {
			content = tool.Convert(content, `gbk`, `utf-8`)
		}
		rows, err = csv.NewReader(strings.NewReader(content)).ReadAll()
		if err != nil {
			return nil, err
		}
	} else {
		f, err := excelize.OpenFile(GetFileByLink(fileUrl))
		if err != nil {
			return nil, err
		}
		rows, err = f.GetRows(f.GetSheetName(f.GetActiveSheetIndex()))
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func ReadTab(fileUrl, fileExt string) (define.DocSplitItems, int, error) {
	rows, err := ParseTabFile(fileUrl, fileExt)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) < 2 {
		return nil, 0, errors.New(`excel_less_row`)
	}
	//line collection
	list := make(define.DocSplitItems, 0)
	wordTotal := 0
	for i := range rows {
		pairs := make([]string, 0)
		for j := range rows[i] {
			wordTotal += utf8.RuneCountInString(rows[i][j])
			if i == 0 { //excel head
				continue
			}
			if len(rows[i][j]) == 0 || len(rows[0]) <= j {
				continue
			}
			pairs = append(pairs, fmt.Sprintf(`%s:%s`, rows[0][j], rows[i][j]))
		}
		if len(pairs) == 0 {
			continue
		}
		list = append(list, define.DocSplitItem{PageNum: i, Content: strings.Join(pairs, `;`)})
	}
	return list, wordTotal, nil
}

func ReadTxt(fileUrl string) (define.DocSplitItems, int, error) {
	if !LinkExists(fileUrl) {
		return nil, 0, errors.New(`file not exist:` + fileUrl)
	}
	content, err := tool.ReadFile(GetFileByLink(fileUrl))
	if err != nil {
		return nil, 0, err
	}
	if !utf8.ValidString(content) {
		content = tool.Convert(content, `gbk`, `utf-8`)
	}
	list := define.DocSplitItems{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}

func ColumnIndexFromIdentifier(identifier string) (int, error) {
	if identifier == "" {
		return -1, errors.New("identifier cannot be empty")
	}
	colIndex := 0
	for _, char := range identifier {
		if !unicode.IsUpper(char) || char < 'A' || char > 'Z' {
			return -1, errors.New("invalid identifier: identifiers must be uppercase letters (A-Z)")
		}
		colIndex = colIndex*26 + int(char-'A'+1)
	}
	return colIndex - 1, nil
}

func IdentifierFromColumnIndex(index int) (string, error) {
	if index < 0 {
		return "", errors.New("index cannot be negative")
	}
	var identifier strings.Builder
	for index >= 0 {
		charIndex := (index % 26) + 'A'
		identifier.WriteByte(byte(charIndex))
		index = (index / 26) - 1
	}
	result := identifier.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func ReadQaTab(fileUrl, fileExt string, splitParams define.SplitParams) (define.DocSplitItems, int, error) {
	rows, err := ParseTabFile(fileUrl, fileExt)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) < 2 {
		return nil, 0, errors.New(`excel_less_row`)
	}
	questionIndex, err := ColumnIndexFromIdentifier(splitParams.QuestionColumn)
	if err != nil {
		return nil, 0, err
	}
	answerIndex, err := ColumnIndexFromIdentifier(splitParams.AnswerColumn)
	if err != nil {
		return nil, 0, err
	}
	similarIndex := -1
	if len(splitParams.SimilarColumn) > 0 {
		similarIndex, err = ColumnIndexFromIdentifier(splitParams.SimilarColumn)
		if err != nil {
			return nil, 0, err
		}
	}
	if questionIndex == answerIndex || questionIndex == similarIndex || answerIndex == similarIndex {
		return nil, 0, errors.New(`excel question index cannot be equal to answer`)
	}

	//line collection
	list := make(define.DocSplitItems, 0)
	wordTotal := 0
	for i, row := range rows[1:] {
		var question, answer string
		var similarQuestionList []string
		if len(row) > answerIndex {
			answer = row[answerIndex]
		}
		if len(row) > questionIndex {
			question = row[questionIndex]
		}
		if similarIndex > 0 && len(row) > similarIndex && len(row[similarIndex]) > 0 {
			// 通过换行符分割相似问题
			similarQuestionList = append(similarQuestionList, strings.FieldsFunc(row[similarIndex], func(r rune) bool {
				return r == '\n' || r == '\r'
			})...)
		}
		answer, images := ExtractTextImages(answer)
		if len(answer) == 0 || len(question) == 0 {
			continue
		}

		wordTotal += utf8.RuneCountInString(question + answer)
		list = append(list, define.DocSplitItem{PageNum: i + 1, Question: question, SimilarQuestionList: similarQuestionList, Answer: answer, Images: images})
	}

	return list, wordTotal, nil
}

func QaDocSplit(splitParams define.SplitParams, items define.DocSplitItems) define.DocSplitItems {
	list := make(define.DocSplitItems, 0)
	for i, item := range items {
		for _, section := range strings.Split(item.Content, splitParams.QuestionLable) {
			if len(strings.TrimSpace(section)) == 0 {
				continue
			}

			var question string
			var similarList []string
			var answer string
			var images []string

			if splitParams.SimilarLabel == "" {
				qa := strings.SplitN(section, splitParams.AnswerLable, 2)
				if len(qa) == 0 {
					continue
				} else if len(qa) == 1 {
					question = qa[0]
					continue
				} else if len(qa) == 2 {
					question = qa[0]
					answer = qa[1]
				} else {
					continue
				}
			} else {
				splitParams.SimilarLabel = strings.ReplaceAll(splitParams.SimilarLabel, ":", "：")
				splitParams.AnswerLable = strings.ReplaceAll(splitParams.AnswerLable, ":", "：")
				similarIndex := strings.Index(section, splitParams.SimilarLabel)
				answerIndex := strings.Index(section, splitParams.AnswerLable)

				if similarIndex != -1 && answerIndex != -1 && similarIndex < answerIndex {
					// 情况 1: 找到了 Similar 标签并在 Answer 标签之前（预期的 Q/S/A 结构）
					question = section[:similarIndex]

					// 相似问题文本块在 Similar 标签和 Answer 标签之间
					similarBlock := section[similarIndex+len(splitParams.SimilarLabel) : answerIndex]
					// 通过换行符分割相似问题
					similarList = append(similarList, strings.FieldsFunc(similarBlock, func(r rune) bool {
						return r == '\n' || r == '\r'
					})...)
					// similarList = append(similarList, similarBlock)

					// 答案部分从 Answer 标签之后开始
					answer = section[answerIndex+len(splitParams.AnswerLable):]
				} else if answerIndex != -1 {
					// 情况 2: 找到了 Answer 标签，但 Similar 标签缺失或在 Answer 之后
					// 将从块的开始到 Answer 标签之前的所有内容视为问题。
					// 相似问题列表将为空。
					question = section[:answerIndex]
					similarList = []string{} // 没有找到相似问题

					// 答案部分从 Answer 标签之后开始
					answer = section[answerIndex+len(splitParams.AnswerLable):]
				} else {
					// 情况 3: Answer 和 Similar 都没有按预期方式找到（或者只找到了 Similar 但没有 Answer）。
					// 这个块不符合预期的 Q&A 格式，跳过。
					continue
				}
			}

			question = strings.TrimSpace(question)
			answer = strings.TrimSpace(answer)
			answer, images = ExtractTextImages(answer) // 对去除空白后的答案文本进行处理

			// 检查处理后是否得到了有效的问题和答案
			if len(question) == 0 || len(answer) == 0 {
				// fmt.Printf("警告: 跳过处理后问题或答案为空的块 (页码 %d, 块索引 %d)。\n", i+1, sectionIndex) // 可选的日志记录
				continue // 如果问题或答案为空则跳过
			}

			// 将成功提取的项添加到列表中
			list = append(list, define.DocSplitItem{
				PageNum:             i + 1, // 假设 i 是原始项/页的索引，从 0 开始
				Question:            question,
				SimilarQuestionList: similarList, // 存储相似问题列表
				Answer:              answer,
				Images:              images,
			})
		}
	}
	return list
}

func ExtractTextImages(content string) (string, []string) {
	re := regexp.MustCompile(`\{\{\!\!(.+?)\!\!\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	images := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	content = re.ReplaceAllString(content, "")
	return content, images
}

func ExtractTextImagesPlaceholders(content string) (string, []string) {
	re := regexp.MustCompile(`\{\{\!\!(.+?)\!\!\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	images := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	content = re.ReplaceAllString(content, lib_define.ImagePlaceholder)
	return content, images
}

func InTextImagesPlaceholders(content string, images []string) (string, []string, []string) {

	count := strings.Count(content, lib_define.ImagePlaceholder)
	imgs := make([]string, 0)
	for i := 0; i < count; i++ {
		if len(images) > 0 {
			imgs = append(imgs, images[0])
			images = images[1:]
		}
	}
	content = strings.ReplaceAll(content, lib_define.ImagePlaceholder, "")
	return content, imgs, images
}

func EmbTextImages(content string, images []string) string {
	var imgTags []string
	for _, image := range images {
		imgTags = append(imgTags, fmt.Sprintf(`<img src="%s">`, image))
	}

	return content + "\n" + strings.Join(imgTags, " ")
}

func MbSubstr(s string, start, length int) string {
	runes := []rune(s)
	if start >= len(runes) {
		return ""
	}
	if start+length > len(runes) {
		length = len(runes) - start
	}
	return string(runes[start : start+length])
}

func ConvertWebPToPNG(webpData []byte) ([]byte, error) {
	img, err := webp.Decode(bytes.NewReader(webpData))
	if err != nil {
		return nil, fmt.Errorf("error decoding WebP image: %w", err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("error encoding PNG image: %w", err)
	}

	return buf.Bytes(), nil
}

func DefaultSplitParams() define.SplitParams {
	return define.SplitParams{
		ChunkSize:          512,
		ChunkOverlap:       0,
		SeparatorsNo:       `12,11`,
		EnableExtractImage: true,
	}
}

func AutoSplitLibFile(adminUserId, fileId int, splitParams define.SplitParams) {
	lang := define.LangEnUs
	list, wordTotal, splitParams, err := GetLibFileSplit(adminUserId, fileId, 0, splitParams, lang)
	if err != nil {
		UpdateLibFileData(adminUserId, fileId, msql.Datas{`status`: define.FileStatusException, `errmsg`: err.Error()})
		logs.Error(err.Error())
		return
	}
	// default type
	if splitParams.QaIndexType == 0 {
		splitParams.QaIndexType = define.QAIndexTypeQuestionAndAnswer
	}
	if splitParams.ChunkType == define.ChunkTypeAi && splitParams.AiChunkNew {
		if err = SaveAISplitDocs(adminUserId, fileId, wordTotal, splitParams.QaIndexType, splitParams, list, 0, lang); err != nil {
			logs.Error(err.Error())
			UpdateLibFileData(adminUserId, fileId, msql.Datas{`status`: define.FileStatusException, `errmsg`: err.Error()})
			return
		}
	} else {
		// 知识库问答文件导入走这里@sizz 20250903
		_, err = SaveLibFileSplit(adminUserId, fileId, wordTotal, splitParams.QaIndexType, splitParams, list, 0, lang)
		if err != nil {
			logs.Error(err.Error())
			return
		}
	}
}

func UpdateLibFileData(adminUserId, fileId int, data msql.Datas) error {
	_, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, cast.ToString(fileId)).Update(data)
	return err
}

func AISplitDocs(adminUserId, fileId int, splitParams define.SplitParams, list define.DocSplitItems) (define.DocSplitItems, error) {
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	contents := ""
	var errMsg = ""
	contentMap := make(map[int]define.DocSplitItem)
	index := 1
	maxToken := 0
	currChan := make(chan struct{}, 8)
	if !CheckModelIsValid(adminUserId, cast.ToInt(splitParams.AiChunkModelConfigId), splitParams.AiChunkModel, Llm) {
		errMsg = `model not valid`
		return nil, errors.New(errMsg)
	}
	if !(splitParams.ParagraphChunk || splitParams.ChunkPreview) {
		UpdateLibFileData(adminUserId, fileId, msql.Datas{`status`: define.FileStatusChunking})
	}

	var submitContent define.DocSplitItems
	spliter := NewAiSpliterClient(splitParams.AiChunkSize)
	for _, item := range list {
		// pdf 按页提交
		if define.IsPdfFile(splitParams.FileExt) {
			spliteItem := define.DocSplitItem{
				PageNum: item.PageNum,
				Content: item.Content,
			}
			submitContent = append(submitContent, spliteItem)
			continue
		} else {
			chunks, err := spliter.SplitText(item.Content)
			if err != nil {
				errMsg = err.Error()
				return nil, errors.New(errMsg)
			}
			for _, chunk := range chunks {
				spliteItem := define.DocSplitItem{
					PageNum: item.PageNum,
					Content: chunk,
					Images:  item.Images,
				}
				submitContent = append(submitContent, spliteItem)
			}
		}
	}
	type chunkResult struct {
		Chunk string `json:"chunk"`
	}
	for _, item := range submitContent {
		contents = item.Content
		// ai split
		if len(contents) == 0 {
			continue
		}
		if len(contents) <= 2 {
			lock.Lock()
			contentMap[index] = define.DocSplitItem{
				PageNum: item.PageNum,
				Content: contents,
				Images:  item.Images,
			}
			lock.Unlock()
			continue
		}
		currChan <- struct{}{}
		wg.Add(1)
		go func(wg *sync.WaitGroup, contentMap map[int]define.DocSplitItem, index int, contents string, errMsg *string) {
			defer wg.Done()
			defer func() { <-currChan }()
			messages := []adaptor.ZhimaChatCompletionMessage{
				{
					Role:    `system`,
					Content: splitParams.AiChunkPrumpt + define.AiChunkPrumptSuffix,
				},
				{
					Role:    `user`,
					Content: contents,
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			retryTimes := 0
			chanStream := make(chan sse.Event, 100)
			go func() {
				for {
					select {
					case <-chanStream:
						continue
					case <-ctx.Done():
						return
					}
				}
			}()
		Retry:
			chatResp, _, err := RequestChatStream(
				define.LangEnUs,
				adminUserId,
				"",
				msql.Params{},
				"",
				cast.ToInt(splitParams.AiChunkModelConfigId),
				splitParams.AiChunkModel,
				messages,
				nil,
				chanStream,
				0.1,
				maxToken,
			)
			docSplitItem := define.DocSplitItem{
				PageNum: item.PageNum,
				Content: contents,
				Images:  item.Images,
			}
			if err != nil && len(chatResp.Result) <= 0 {
				logs.Error(err.Error())
				if retryTimes <= 3 {
					retryTimes++
					goto Retry
				}
				*errMsg = err.Error()
				lock.Lock()
				docSplitItem.AiChunkErrMsg = err.Error()
				contentMap[index] = docSplitItem
				lock.Unlock()
			} else {
				lock.Lock()
				resp := []chunkResult{}
				if err := tool.JsonDecode(chatResp.Result, &resp); err != nil {
					docSplitItem.AiChunkErrMsg = fmt.Sprintf(`ai return format error:%s`, err.Error())
				} else {
					docSplitItem.Content = chatResp.Result
				}
				contentMap[index] = docSplitItem
				lock.Unlock()
			}
			cancel()
		}(wg, contentMap, index, contents, &errMsg)
		index++
		if splitParams.ChunkPreview {
			break
		}
	}

	wg.Wait()
	var newList define.DocSplitItems
	var imageInsertIndex = map[int][]string{}
	var number = 1

	for i := 1; i <= len(contentMap); i++ {
		if contentMap[i].AiChunkErrMsg != "" {
			newItem := define.DocSplitItem{
				PageNum:       contentMap[i].PageNum,
				Images:        contentMap[i].Images,
				Number:        number,
				Content:       contentMap[i].Content,
				WordTotal:     utf8.RuneCountInString(contentMap[i].Content),
				AiChunkErrMsg: contentMap[i].AiChunkErrMsg,
			}
			newList = append(newList, newItem)
			number++
		} else {
			allContents := []chunkResult{}
			if err := tool.JsonDecode(contentMap[i].Content, &allContents); err != nil {
				errMsg = fmt.Sprintf(`ai return format error:%s`, err.Error())
			}
			for _, item := range allContents {
				if len(item.Chunk) == 0 {
					continue
				}
				content, images := ExtractTextImages(item.Chunk)
				if len(images) > 0 && len(content) == 0 {
					index := min(len(newList)-1, 0)
					imageInsertIndex[index] = images
				} else {
					newItem := define.DocSplitItem{
						Images:    images,
						PageNum:   contentMap[i].PageNum,
						Number:    number,
						Content:   content,
						WordTotal: utf8.RuneCountInString(content),
					}
					if len(images) <= 0 && len(contentMap[i].Images) > 0 {
						newItem.Images = contentMap[i].Images
					}
					newList = append(newList, newItem)
					number++
				}
			}
		}
	}
	if len(imageInsertIndex) > 0 {
		for index, images := range imageInsertIndex {
			newList[index].Images = append(newList[index].Images, images...)
		}
	}
	// save data
	redisClient := LibFileSplitAiChunksBacheHandle{TaskId: splitParams.AiChunkTaskId}
	redisClient.SaveCacheData(LibFileSplitAiChunksCache{
		List:   newList,
		ErrMsg: errMsg,
	})
	if errMsg != "" {
		logs.Error(errMsg)
		return newList, errors.New(errMsg)
	}
	if !(splitParams.ParagraphChunk || splitParams.ChunkPreview) {
		UpdateLibFileData(adminUserId, fileId, msql.Datas{`status`: define.FileStatusWaitSplit})
	}
	return newList, nil
}

func SaveAISplitDocs(userId, fileId, wordTotal, qaIndexType int, splitParams define.SplitParams, list define.DocSplitItems, pdfPageNum int, lang string) (err error) {
	// save data
	taskId := splitParams.AiChunkTaskId
	ticker := time.NewTicker(1 * time.Second)
	timeout := 0
	errChan := make(chan error, 1)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				timeout++
				var list = LibFileSplitAiChunksCache{}
				if err = lib_redis.GetCacheWithBuild(define.Redis, &LibFileSplitAiChunksBacheHandle{TaskId: taskId}, &list, 1*time.Hour); err != nil {
					logs.Error(err.Error())
					errChan <- err
					return
				}
				if list.ErrMsg != "" {
					errChan <- errors.New(list.ErrMsg)
				}
				if len(list.List) == 0 {
					continue
				}
				if timeout == 3600 {
					errChan <- errors.New(`ai request timeout`)
					return
				}
				SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType, splitParams, list.List, pdfPageNum, lang)
				errChan <- nil
				return
			}
		}
	}()
	lib_redis.DelCacheData(define.Redis, &LibFileSplitAiChunksBacheHandle{TaskId: taskId})
	if err = <-errChan; err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}
