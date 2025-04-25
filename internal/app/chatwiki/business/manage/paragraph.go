// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

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
		Order(`a.page_num asc, a.number asc, a.id desc`)
	if status >= 0 {
		query.Where(`b.status`, cast.ToString(status))
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
	similarQuestion := strings.TrimSpace(c.PostForm(`similar_questions`))
	images := c.PostFormArray(`images`)
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
	fileInfo, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Find()
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
		if len(question) < 1 || len(question) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
		if len(answer) < 1 || len(answer) > common.MaxContent {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `length_error`))))
			return
		}
	} else {
		if len(content) < 1 || len(content) > common.MaxContent {
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
		`library_id`:    fileInfo[`library_id`],
		`file_id`:       fileId,
		`title`:         title,
		`images`:        jsonImages,
		`update_time`:   tool.Time2Int(),
	}
	var vectorIds []int64
	if cast.ToInt(fileInfo[`is_qa_doc`]) == define.DocTypeQa {
		data[`word_total`] = utf8.RuneCountInString(question + answer)
		data[`content`] = ``
		data[`question`] = question
		data[`answer`] = answer
		data[`similar_questions`] = similarQuestion
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
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeQuestion), question)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorIds = append(vectorIds, vectorID)

		if fileInfo[`type`] == cast.ToString(define.QAIndexTypeQuestionAndAnswer) {
			vectorID, err = common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
				fileId, id, cast.ToString(define.VectorTypeAnswer), question)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
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
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		vectorID, err := common.SaveVector(int64(userId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeParagraph), content)
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

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
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

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
