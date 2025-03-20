// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
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
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	m := msql.Model(`chat_ai_library_file`, define.Postgres).
		Alias(`f`).
		Join(`chat_ai_library_file_data d`, `f.id=d.file_id`, `left`).
		Where(`f.admin_user_id`, cast.ToString(userId)).
		Where(`f.library_id`, cast.ToString(libraryId)).
		Group(`f.id`).
		Field(`f.*, count(d.id) as paragraph_count`)
	fileName := strings.TrimSpace(c.Query(`file_name`))
	if len(fileName) > 0 {
		m.Where(`file_name`, `like`, fileName)
	}
	list, total, err := m.Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`info`: info, `list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
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
	//database dispose
	fileIds := make([]int64, 0)

	for _, uploadInfo := range libraryFiles {
		status := define.FileStatusInitial
		isTableFile := define.IsTableFile(uploadInfo.Ext)
		if !define.IsPdfFile(uploadInfo.Ext) {
			status = define.FileStatusWaitSplit
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
			`file_ext`:             uploadInfo.Ext,
			`file_size`:            uploadInfo.Size,
			`create_time`:          tool.Time2Int(),
			`update_time`:          tool.Time2Int(),
			`is_table_file`:        cast.ToInt(isTableFile),
			`doc_type`:             uploadInfo.GetDocType(),
			`doc_url`:              uploadInfo.DocUrl,
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
		} else {
			fileIds = append(fileIds, fileId)
			if status == define.FileStatusInitial && !uploadInfo.Custom { //async task:convert html
				if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `file_url`: uploadInfo.Link}); err != nil {
					logs.Error(err.Error())
				} else if err := common.AddJobs(define.ConvertHtmlTopic, message); err != nil {
					logs.Error(err.Error())
				}
			}
			if define.IsMdFile(uploadInfo.Ext) { //markdown文档自动切分
				go common.AutoSplitLibFile(userId, int(fileId))
			}
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
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})
	//dispose relation data
	_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`file_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`file_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
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
	splitParams := define.SplitParams{
		IsDiySplit:         cast.ToInt(c.Query(`is_diy_split`)),
		SeparatorsNo:       strings.TrimSpace(c.Query(`separators_no`)),
		Separators:         make([]string, 0),
		ChunkSize:          cast.ToInt(c.Query(`chunk_size`)),
		ChunkOverlap:       cast.ToInt(c.Query(`chunk_overlap`)),
		IsQaDoc:            cast.ToInt(c.Query(`is_qa_doc`)),
		QuestionLable:      strings.TrimSpace(c.Query(`question_lable`)),
		AnswerLable:        strings.TrimSpace(c.Query(`answer_lable`)),
		QuestionColumn:     strings.TrimSpace(c.Query(`question_column`)),
		AnswerColumn:       strings.TrimSpace(c.Query(`answer_column`)),
		EnableExtractImage: cast.ToBool(c.Query(`enable_extract_image`)),
	}
	list, wordTotal, err := common.GetLibFileSplit(userId, fileId, splitParams, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
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
	if err := tool.JsonDecodeUseNumber(c.PostForm(`split_params`), &splitParams); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `split_params`))))
		return
	}
	list, wordTotal, err := common.GetLibFileSplit(userId, fileId, splitParams, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	if err := tool.JsonDecodeUseNumber(c.PostForm(`list`), &list); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `list`))))
		return
	}
	if fileId <= 0 || wordTotal <= 0 || len(list) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	err = common.SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType, splitParams, list, common.GetLang(c))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
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
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if docAutoRenewFrequency < 0 || docAutoRenewFrequency > 5 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `doc_auto_renew_frequency`))))
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

	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`doc_auto_renew_frequency`: docAutoRenewFrequency,
	})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}

	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: id})

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
