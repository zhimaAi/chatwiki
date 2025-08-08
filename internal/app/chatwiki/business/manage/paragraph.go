// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
)

func GetParagraphList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.Query(`file_id`))
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
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	status := cast.ToInt(c.Query(`status`))
	graphStatus := cast.ToInt(c.Query(`graph_status`))
	categoryId := cast.ToInt(c.Query(`category_id`))
	query := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Alias("a").
		Join("chat_ai_library_file_data_index b", "a.id=b.data_id", "inner").
		Where(`a.admin_user_id`, cast.ToString(userId)).Where(`a.file_id`, cast.ToString(fileId)).
		Where(`a.isolated`, "false").
		Where(`a.delete_time`, `0`).
		Field(`a.*`).
		Field(`
			CASE 
    			WHEN bool_and(b.status = 0) THEN 0
    			WHEN bool_and(b.status = 1) THEN 1
    			WHEN bool_and(b.status = 3) THEN 3
    			ELSE 2
			END AS status		
		`).
		Field(`
			COALESCE(
    			(SELECT errmsg FROM chat_ai_library_file_data_index WHERE data_id = a.id AND errmsg IS NOT NULL LIMIT 1),
    			'no error'
  			) AS errmsg
		`).
		Group(`a.id`).
		Order(`a.page_num asc, a.number asc`)
	if status >= 0 {
		if status == define.SplitStatusException {
			query.Where(`a.split_status`, cast.ToString(status))
		} else {
			query.Where(`b.status`, cast.ToString(status))
		}
	}
	if graphStatus >= 0 {
		query.Where(`a.graph_status`, cast.ToString(graphStatus))
	}
	if categoryId >= 0 {
		query.Where(`a.category_id`, cast.ToString(categoryId))
	}
	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	var formatedList []map[string]any
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		var images []string
		err = json.Unmarshal([]byte(item[`images`]), &images)
		if err != nil {
			continue
		}
		tempItem[`images`] = images
		formatedList = append(formatedList, tempItem)
	}

	data := map[string]any{`info`: info, `list`: formatedList, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetParagraphCount(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fileId := cast.ToInt(c.Query(`file_id`))
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
	query := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Alias("a").
		Join("chat_ai_library_file_data_index b", "a.id=b.data_id", "inner").
		Where(`a.admin_user_id`, cast.ToString(userId)).Where(`a.file_id`, cast.ToString(fileId)).
		Where(`a.isolated`, "false").
		Where(`a.delete_time`, `0`).
		Field(`a.split_status,b.status`)
	list, err := query.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	var formatedList = map[string]int{
		`vector_status_initial`:    0,
		`vector_status_converted`:  0,
		`vector_status_exception`:  0,
		`vector_status_converting`: 0,
		`split_status_exception`:   0,
	}
	for _, item := range list {
		formatedList[`total`] += 1
		switch cast.ToInt(item[`split_status`]) {
		case define.SplitStatusException:
			formatedList[`split_status_exception`] += 1
			continue
		}
		switch cast.ToInt(item[`status`]) {
		case define.VectorStatusInitial:
			formatedList[`vector_status_initial`] += 1
		case define.VectorStatusConverted:
			formatedList[`vector_status_converted`] += 1
		case define.VectorStatusConverting:
			formatedList[`vector_status_converting`] += 1
		default:
			formatedList[`vector_status_exception`] += 1
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(formatedList, nil))
}

func GetCategoryParagraphList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	libraryId := cast.ToInt(c.Query(`library_id`))
	status := cast.ToInt(c.Query(`status`))
	categoryId := cast.ToInt(c.Query(`category_id`))
	library, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(library) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	query := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Alias("a").
		Join(`chat_ai_library_file f`, `a.file_id=f.id`, `left`).
		Join("chat_ai_library_file_data_index b", "a.id=b.data_id", "inner").
		Where(`a.admin_user_id`, cast.ToString(userId)).Where(`a.library_id`, cast.ToString(libraryId)).
		Where(`a.category_id`, `>`, `0`).
		Where(`a.delete_time`, `0`).
		Field(`a.*,f.file_name`).
		Field(`
			CASE 
    			WHEN bool_and(b.status = 0) THEN 0
    			WHEN bool_and(b.status = 1) THEN 1
    			WHEN bool_and(b.status = 3) THEN 3
    			ELSE 2
			END AS status		
		`).
		Field(`
			COALESCE(
    			(SELECT errmsg FROM chat_ai_library_file_data_index WHERE data_id = a.id AND errmsg IS NOT NULL LIMIT 1),
    			'no error'
  			) AS errmsg
		`).
		Group(`a.id,f.file_name`).
		Order(`a.create_time asc`)
	if status >= 0 {
		query.Where(`b.status`, cast.ToString(status))
	}
	if categoryId >= 0 {
		query.Where(`a.category_id`, cast.ToString(categoryId))
	}
	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	var formatedList []map[string]any
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		var images []string
		err = json.Unmarshal([]byte(item[`images`]), &images)
		if err != nil {
			continue
		}
		tempItem[`images`] = images
		formatedList = append(formatedList, tempItem)
	}

	data := map[string]any{`info`: library, `list`: formatedList, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func SaveCategoryParagraph(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt64(c.PostForm(`id`))
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	title := strings.TrimSpace(c.PostForm(`title`))
	content := strings.TrimSpace(c.PostForm(`content`))
	question := strings.TrimSpace(c.PostForm(`question`))
	similarQuestions := strings.TrimSpace(c.PostForm(`similar_questions`))
	answer := strings.TrimSpace(c.PostForm(`answer`))
	images := c.PostFormArray(`images`)
	categoryId := cast.ToInt(c.PostForm(`category_id`))
	if id < 0 || libraryId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	library, err := common.GetLibraryInfo(libraryId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	if id > 0 {
		data, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(data) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}

	if cast.ToInt(library[`type`]) == define.QALibraryType {
		if len(question) < 1 || utf8.RuneCountInString(question) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
		if len(answer) < 1 || utf8.RuneCountInString(answer) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
	} else {
		if len(content) < 1 || utf8.RuneCountInString(content) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
	}
	jsonImages, err := common.CheckLibraryImage(images)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `images`))))
		return
	}

	_ = m.Begin()
	data := msql.Datas{
		`admin_user_id`: userId,
		`library_id`:    libraryId,
		`title`:         title,
		`images`:        jsonImages,
		`category_id`:   categoryId,
		`update_time`:   tool.Time2Int(),
	}
	var vectorIds []int64
	if cast.ToInt(library[`type`]) == define.QALibraryType {
		data[`word_total`] = utf8.RuneCountInString(question + answer)
		data[`question`] = question
		data[`similar_questions`] = similarQuestions
		data[`answer`] = answer
		if id > 0 {
			_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		} else {
			data[`type`] = define.ParagraphTypeDocQA
			data[`create_time`] = data[`update_time`]
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(libraryId), 0, id, cast.ToString(define.VectorTypeQuestion), question)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorIds = append(vectorIds, vectorID)

		vectorID, err = common.SaveVector(int64(userId), cast.ToInt64(libraryId), 0, id, cast.ToString(define.VectorTypeAnswer), question)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorIds = append(vectorIds, vectorID)
	} else {
		data[`word_total`] = utf8.RuneCountInString(content)
		data[`content`] = content
		if id > 0 {
			_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		} else {
			data[`type`] = define.ParagraphTypeNormal
			data[`create_time`] = data[`update_time`]
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(libraryId), 0, id, cast.ToString(define.VectorTypeParagraph), content)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorIds = append(vectorIds, vectorID)
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	//async task:convert vector
	for _, id := range vectorIds {
		if message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: 0}); err != nil {
			logs.Error(err.Error())
			continue
		} else {
			if err = common.AddJobs(define.ConvertVectorTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func getParagraphAddNumber(c *gin.Context, fileId int64) int {
	if number := cast.ToInt(c.PostForm(`number`)); number > 0 {
		return number
	}
	maxNumber, _ := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`file_id`, cast.ToString(fileId)).Max(`number`)
	return cast.ToInt(maxNumber) + 1
}

func SaveParagraph(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt64(c.PostForm(`id`))
	fileId := cast.ToInt64(c.PostForm(`file_id`))
	title := strings.TrimSpace(c.PostForm(`title`))
	content := strings.TrimSpace(c.PostForm(`content`))
	question := strings.TrimSpace(c.PostForm(`question`))
	answer := strings.TrimSpace(c.PostForm(`answer`))
	similarQuestions := strings.TrimSpace(c.PostForm(`similar_questions`))
	images := c.PostFormArray(`images`)
	categoryId := cast.ToInt(c.PostForm(`category_id`))
	if id < 0 || fileId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	if id > 0 {
		fileIdStr, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`file_id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		fileId = cast.ToInt64(fileIdStr)
	}
	if fileId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	fileInfo, err := common.GetLibFileInfo(int(fileId), userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(fileInfo) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if cast.ToInt(fileInfo[`is_qa_doc`]) == define.DocTypeQa {
		if len(question) < 1 || utf8.RuneCountInString(question) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
		if len(answer) < 1 || utf8.RuneCountInString(answer) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
	} else {
		if len(content) < 1 || utf8.RuneCountInString(content) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
	}
	jsonImages, err := common.CheckLibraryImage(images)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `images`))))
		return
	}
	if imagesJson := strings.TrimSpace(c.PostForm(`images_json`)); len(imagesJson) > 0 {
		jsonImages = imagesJson //适用于特殊场景,直接传递参数
	}

	_ = m.Begin()
	data := msql.Datas{
		`admin_user_id`: userId,
		`library_id`:    fileInfo[`library_id`],
		`file_id`:       fileId,
		`title`:         title,
		`images`:        jsonImages,
		`category_id`:   categoryId,
		`update_time`:   tool.Time2Int(),
	}
	var vectorIds []int64
	if cast.ToInt(fileInfo[`is_qa_doc`]) == define.DocTypeQa {
		data[`word_total`] = utf8.RuneCountInString(question + answer)
		data[`content`] = ``
		data[`question`] = question
		data[`answer`] = answer
		data[`similar_questions`] = similarQuestions
		if id > 0 {
			_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		} else {
			data[`type`] = define.ParagraphTypeDocQA
			data[`create_time`] = data[`update_time`]
			data[`number`] = getParagraphAddNumber(c, fileId)
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeQuestion), question)
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		vectorIds = append(vectorIds, vectorID)
		similarQuestionArr := make([]string, 0)
		tool.JsonDecode(similarQuestions, &similarQuestionArr)
		if err = common.DeleteLibraryFileDataIndex(cast.ToString(id), cast.ToString(define.VectorTypeSimilarQuestion)); err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for _, similarQuestion := range similarQuestionArr {
			vectorID, err := common.SaveVector(
				cast.ToInt64(userId),
				cast.ToInt64(fileInfo[`library_id`]),
				fileId,
				id,
				cast.ToString(define.VectorTypeSimilarQuestion),
				strings.TrimSpace(similarQuestion),
			)
			if err != nil {
				logs.Error(err.Error())
				_ = m.Rollback()
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			vectorIds = append(vectorIds, vectorID)
		}

		if fileInfo[`type`] == cast.ToString(define.QAIndexTypeQuestionAndAnswer) {
			vectorID, err = common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
				fileId, id, cast.ToString(define.VectorTypeAnswer), question)
			if err != nil {
				logs.Error(err.Error())
				_ = m.Rollback()
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			vectorIds = append(vectorIds, vectorID)
		}
	} else {
		data[`word_total`] = utf8.RuneCountInString(content)
		data[`content`] = content
		data[`question`] = ``
		data[`answer`] = ``
		if id > 0 {
			_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		} else {
			data[`type`] = define.ParagraphTypeNormal
			data[`create_time`] = data[`update_time`]
			data[`number`] = getParagraphAddNumber(c, fileId)
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeParagraph), content)
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		vectorIds = append(vectorIds, vectorID)
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	//async task:convert vector
	for _, id := range vectorIds {
		if message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId}); err != nil {
			logs.Error(err.Error())
			continue
		} else {
			if err = common.AddJobs(define.ConvertVectorTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}

	if common.GetNeo4jStatus(userId) {
		message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId})
		if err != nil {
			logs.Error(err.Error())
		} else {
			if err = common.AddJobs(define.ConvertGraphTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(m.Where(`id`, cast.ToString(id)).Find()))
}

func SaveSplitParagraph(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt64(c.PostForm(`data_id`))
	fileId := cast.ToInt64(c.PostForm(`file_id`))
	categoryId := cast.ToInt(c.PostForm(`category_id`))
	similarQuestions := cast.ToInt(c.PostForm(`similar_questions`))
	if id < 0 || fileId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	fileInfo, err := common.GetLibFileInfo(int(fileId), adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(fileInfo) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	list := make([]define.DocSplitItem, 0)
	if err := tool.JsonDecodeUseNumber(c.PostForm(`list`), &list); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `list`))))
		return
	}
	list, err = common.SaveSplitParagraph(adminUserId, cast.ToInt(fileId), cast.ToInt(id), list)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	_ = m.Begin()
	var vectorIds []int64
	for _, item := range list {
		jsonImages, err := common.CheckLibraryImage(item.Images)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `images`))))
			return
		}
		data := msql.Datas{
			`admin_user_id`: adminUserId,
			`library_id`:    fileInfo[`library_id`],
			`file_id`:       fileId,
			`title`:         item.Title,
			`images`:        jsonImages,
			`category_id`:   categoryId,
			`update_time`:   tool.Time2Int(),
		}
		if cast.ToInt(fileInfo[`is_qa_doc`]) == define.DocTypeQa {
			data[`word_total`] = utf8.RuneCountInString(item.Question + item.Answer)
			data[`content`] = ``
			data[`question`] = item.Question
			data[`answer`] = item.Answer
			data[`similar_questions`] = similarQuestions
			data[`type`] = define.ParagraphTypeDocQA
			data[`create_time`] = data[`update_time`]
			data[`number`] = item.Number
			id, err = m.Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
			vectorID, err := common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
				fileId, id, cast.ToString(define.VectorTypeQuestion), item.Question)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
			vectorIds = append(vectorIds, vectorID)
			if fileInfo[`type`] == cast.ToString(define.QAIndexTypeQuestionAndAnswer) {
				vectorID, err = common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
					fileId, id, cast.ToString(define.VectorTypeAnswer), item.Question)
				if err != nil {
					logs.Error(err.Error())
					c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
					_ = m.Rollback()
					return
				}
				vectorIds = append(vectorIds, vectorID)
			}
		} else {
			data[`word_total`] = utf8.RuneCountInString(item.Content)
			data[`content`] = item.Content
			data[`question`] = ``
			data[`answer`] = ``
			data[`type`] = define.ParagraphTypeNormal
			data[`create_time`] = data[`update_time`]
			data[`number`] = item.Number
			data[`page_num`] = item.PageNum
			id, err = m.Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
			vectorID, err := common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
				fileId, id, cast.ToString(define.VectorTypeParagraph), item.Content)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
			vectorIds = append(vectorIds, vectorID)
		}

	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		m.Rollback()
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//async task:convert vector
	for _, id := range vectorIds {
		if message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId}); err != nil {
			logs.Error(err.Error())
			continue
		} else {
			if err = common.AddJobs(define.ConvertVectorTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}

	if common.GetNeo4jStatus(adminUserId) {
		message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId})
		if err != nil {
			logs.Error(err.Error())
		} else {
			if err = common.AddJobs(define.ConvertGraphTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetSplitParagraph(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	dataIds := cast.ToString(c.Query(`data_ids`))
	fileId := cast.ToInt(c.Query(`file_id`))
	if len(dataIds) <= 0 || fileId <= 0 {
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
		ParagraphChunk:             true,
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
	} else if splitParams.ChunkType == define.ChunkTypeAi {
		if ok := common.CheckModelIsValid(adminUserId, splitParams.AiChunkModelConfigId, splitParams.AiChunkModel, common.Llm); !ok {
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
	list, wordTotal, splitParams, err := common.GetParagraphSplit(adminUserId, fileId, pdfPageNum, dataIds, splitParams, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), err.Error()))))
		return
	}
	data := map[string]any{`split_params`: splitParams, `list`: list, `word_total`: wordTotal}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func DeleteParagraph(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	data, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(data) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	if cast.ToInt(data[`category_id`]) > 0 {
		_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{"isolated": true})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	} else {
		_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`data_id`, cast.ToString(id)).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if common.GetNeo4jStatus(userId) {
			err = common.NewGraphDB(userId).DeleteByData(id)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UpdateParagraphCategory(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	categoryId := cast.ToInt(c.PostForm(`category_id`))
	if id <= 0 || categoryId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	data, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`delete_time`, `0`).
		Where(`id`, cast.ToString(id)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(data) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{`category_id`: categoryId})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	if cast.ToBool(data[`isolated`]) {
		DeleteParagraph(c)
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GenerateSimilarQuestions(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := cast.ToString(c.PostForm(`use_model`))
	question := strings.TrimSpace(c.PostForm(`question`))
	answer := strings.TrimSpace(c.PostForm(`answer`))
	num := cast.ToInt(c.PostForm(`num`))
	if libraryId == 0 || modelConfigId == 0 || len(useModel) == 0 || len(question) == 0 || len(answer) == 0 || num <= 0 || num > 20 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	library, err := common.GetLibraryInfo(libraryId, userId)
	if len(library) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	prompt := strings.ReplaceAll(define.PromptGenerateSimilarQuestions, `{{num}}`, cast.ToString(num))
	prompt = strings.ReplaceAll(prompt, `{{question}}`, question)
	prompt = strings.ReplaceAll(prompt, `{{answer}}`, answer)

	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `user`, Content: prompt}}

	chatResp, _, err := common.RequestChat(
		userId,
		``,
		msql.Params{},
		``,
		modelConfigId,
		useModel,
		messages,
		nil,
		0.1,
		1024,
	)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	chatResp.Result = strings.TrimPrefix(chatResp.Result, "```json")
	chatResp.Result = strings.TrimSuffix(chatResp.Result, "```")
	var result []string
	err = json.Unmarshal([]byte(chatResp.Result), &result)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func GenerateAiPrompt(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	aiChunkModel := strings.TrimSpace(c.PostForm(`ai_prompt_model`))
	aiChunkModelConfigId := cast.ToInt(c.PostForm(`ai_prompt_model_config_id`))
	aiPromptQuestion := strings.TrimSpace(c.PostForm(`ai_prompt_question`))
	// check conf
	if !common.CheckModelIsValid(adminUserId, aiChunkModelConfigId, aiChunkModel, common.Llm) {
		common.FmtError(c, `param_invalid`, `ai_prompt_model`)
		return
	}
	if len(aiPromptQuestion) == 0 {
		common.FmtError(c, `question_empty`)
		return
	}
	maxTokens := 0
	// 生成AI提示词
	prompt := fmt.Sprintf(define.PromptAiGenerate, 500)
	messages := []adaptor.ZhimaChatCompletionMessage{
		{
			Role:    "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: aiPromptQuestion,
		},
	}
	chatResp, _, err := common.RequestChat(
		adminUserId,
		"",
		msql.Params{},
		"",
		cast.ToInt(aiChunkModelConfigId),
		aiChunkModel,
		messages,
		nil,
		0.1,
		maxTokens,
	)
	if err != nil && len(chatResp.Result) <= 0 {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	common.FmtOk(c, chatResp.Result)
}
