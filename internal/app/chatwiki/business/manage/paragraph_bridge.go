// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type BridgeGetParagraphListReq struct {
	LibraryId   string `form:"library_id"`
	FileId      string `form:"file_id"`
	Page        string `form:"page"`
	Size        string `form:"size"`
	Status      string `form:"status"`
	GraphStatus string `form:"graph_status"`
	CategoryId  string `form:"category_id"`
	GroupId     string `form:"group_id"`
	SortField   string `form:"sort_field"`
	SortType    string `form:"sort_type"`
	Search      string `form:"search"`
}

func BridgeGetParagraphList(adminUserId, loginUserId int, lang string, req *BridgeGetParagraphListReq) (map[string]any, int, error) {
	libraryId := cast.ToUint(req.LibraryId)
	fileId := cast.ToUint(req.FileId)
	if libraryId == 0 && fileId == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	var (
		info    = make(msql.Params)
		err     error
		paraIds []string
	)
	if fileId > 0 {
		info, _ = common.GetLibFileInfo(int(fileId), adminUserId)
	}
	page := max(1, cast.ToInt(req.Page))
	size := max(1, cast.ToInt(req.Size))
	status := cast.ToInt(req.Status)
	graphStatus := cast.ToInt(req.GraphStatus)
	categoryId := cast.ToInt(req.CategoryId)
	groupId := cast.ToInt(req.GroupId)
	sortField := req.SortField
	sortType := req.SortType
	search := cast.ToString(req.Search)
	if len(search) > 0 && fileId > 0 {
		paraIds, err = common.GetMatchFileParagraphIdsByFullTextSearch(search, cast.ToString(fileId))
		if err != nil {
			logs.Error(`GetMatchFileParagraphIdsByFullTextSearch err:%v`, err)
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	if len(search) > 0 && libraryId > 0 {
		paraIds, err = common.GetMatchLibraryDataIdsByLike(search, cast.ToString(libraryId), 1000)
		if err != nil {
			logs.Error(`GetMatchLibraryDataIdsByLike err:%v`, err)
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		if len(paraIds) == 0 {
			data := map[string]any{`info`: info, `list`: make([]map[string]any, 0), `total`: 0, `exception_total`: 0, `page`: page, `size`: size}
			return data, 0, nil
		}
	}
	query := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Alias("a").
		Join("chat_ai_library_file_data_index b", "a.id=b.data_id", "inner").
		Where(`a.admin_user_id`, cast.ToString(adminUserId)).
		Where(`a.isolated`, "false").
		Where(`a.delete_time`, `0`).
		Field(`a.*,STRING_AGG(b.status::text, ',') all_status`).
		Field(`a.*`).
		Field(`
			CASE 
    			WHEN bool_or(b.status = 2) THEN 2
    			WHEN bool_or(b.status = 3) THEN 3
    			WHEN bool_and(b.status = 0) THEN 0
    			WHEN bool_and(b.status = 1) THEN 1
    			ELSE 3
			END AS status		
		`).
		Field(`
			COALESCE(
    			(SELECT errmsg FROM chat_ai_library_file_data_index WHERE data_id = a.id AND errmsg IS NOT NULL LIMIT 1),
    			'no error'
  			) AS errmsg
		`).
		Group(`a.id`)
	exception := msql.Model(`chat_ai_library_file_data`, define.Postgres).Alias(`a`).
		Join(`chat_ai_library_file_data_index b`, `a.id=b.data_id`, `inner`).
		Where(`a.admin_user_id`, cast.ToString(adminUserId)).Where(`a.isolated`, `false`).
		Where(`b.status`, cast.ToString(define.VectorStatusException)).Group(`a.id`)
	if libraryId > 0 {
		query.Where(`a.library_id`, cast.ToString(libraryId))
		exception.Where(`a.library_id`, cast.ToString(libraryId))
	}
	if fileId > 0 {
		query.Where(`a.file_id`, cast.ToString(fileId))
		exception.Where(`a.file_id`, cast.ToString(fileId))
	}
	if status >= 0 {
		if status == define.SplitStatusException {
			query.Where(`a.split_status`, cast.ToString(status))
			exception.Where(`a.split_status`, cast.ToString(status))
		} else {
			query.Where(`b.status`, cast.ToString(status))
			exception.Where(`b.status`, cast.ToString(status))
		}
	}
	if graphStatus >= 0 {
		query.Where(`a.graph_status`, cast.ToString(graphStatus))
		exception.Where(`a.graph_status`, cast.ToString(graphStatus))
	}
	if categoryId >= 0 {
		query.Where(`a.category_id`, cast.ToString(categoryId))
		exception.Where(`a.category_id`, cast.ToString(categoryId))
	}
	if groupId >= 0 {
		query.Where(`a.group_id`, cast.ToString(groupId))
		exception.Where(`a.group_id`, cast.ToString(groupId))
	}
	if len(paraIds) > 0 {
		query.Where(`b.id`, `in`, strings.Join(paraIds, `,`))
		exception.Where(`b.id`, `in`, strings.Join(paraIds, `,`))
	}
	orderRaw := `a.page_num,a.father_chunk_paragraph_number,a.number`
	if len(sortField) > 0 && len(sortType) > 0 {
		orderRaw = fmt.Sprintf(`a.%s %s`, sortField, sortType)
	}
	query.Order(orderRaw)

	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	//===== meta_list (also returned for paragraphs) =====
	//Batch prepare: group_name (QA group), library meta schema, built-in is_show switch
	libraryIdSet := make(map[int]struct{})
	for _, it := range list {
		libraryIdSet[cast.ToInt(it[`library_id`])] = struct{}{}
	}
	libIds := make([]string, 0, len(libraryIdSet))
	for id := range libraryIdSet {
		if id > 0 {
			libIds = append(libIds, cast.ToString(id))
		}
	}

	// library_id -> show_meta_*
	libShowMap := make(map[int]msql.Params)
	if len(libIds) > 0 {
		libs, e := msql.Model(`chat_ai_library`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, strings.Join(libIds, `,`)).
			Field(`id,show_meta_source,show_meta_update_time,show_meta_create_time,show_meta_group`).
			Select()
		if e == nil {
			for _, l := range libs {
				libShowMap[cast.ToInt(l[`id`])] = l
			}
		}
	}

	// library_id -> schema list
	schemaByLib := make(map[int][]msql.Params)
	if len(libIds) > 0 {
		schemaList, e := msql.Model(`library_meta_schema`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, `in`, strings.Join(libIds, `,`)).
			Order(`id asc`).
			Field(`id,library_id,name,key,type,is_show`).
			Select()
		if e == nil {
			for _, s := range schemaList {
				schemaByLib[cast.ToInt(s[`library_id`])] = append(schemaByLib[cast.ToInt(s[`library_id`])], s)
			}
		}
	}

	//(library_id, group_id) -> group_name (QA groups only)
	groupNameMap := make(map[string]string)
	groupNameMap[`0:0`] = lib_define.Ungrouped
	if len(libIds) > 0 {
		groups, e := msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, `in`, strings.Join(libIds, `,`)).
			Where(`group_type`, cast.ToString(define.LibraryGroupTypeQA)).
			Field(`id,library_id,group_name`).
			Select()
		if e == nil {
			for _, g := range groups {
				k := fmt.Sprintf(`%d:%d`, cast.ToInt(g[`library_id`]), cast.ToInt(g[`id`]))
				groupNameMap[k] = g[`group_name`]
			}
		}
	}

	var formatedList = make([]map[string]any, 0)
	builtinMetaSchemaList := common.GetBuiltinMetaSchemaList(lang)
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

		//meta_list: built-in + custom (do not mutate DB metadata; parse from a.metadata)
		libId := cast.ToInt(item[`library_id`])
		//Source: mapped from chat_ai_library_file_data.type
		// type=1 => source=4，type=2 => source=5
		paraType := cast.ToInt(item[`type`])
		sourceVal := 0
		if paraType == define.ParagraphTypeNormal {
			sourceVal = define.DocTypeDiy
		} else if paraType == define.ParagraphTypeDocQA {
			sourceVal = define.DocTypeOfficial
		}
		gid := cast.ToInt(item[`group_id`])
		showCfg := libShowMap[libId]
		gkey := fmt.Sprintf(`%d:%d`, libId, gid)
		groupName := groupNameMap[gkey]
		if groupName == `` {
			groupName = lib_define.Ungrouped
		}
		metaStr := strings.TrimSpace(cast.ToString(item[`metadata`]))
		if metaStr == `` {
			metaStr = `{}`
		}
		metaMap := make(map[string]any)
		_ = tool.JsonDecode(metaStr, &metaMap)

		builtinValueMap := map[string]any{
			define.BuiltinMetaKeyUpdateTime: cast.ToInt(item[`update_time`]),
			define.BuiltinMetaKeyCreateTime: cast.ToInt(item[`create_time`]),
			define.BuiltinMetaKeySource:     sourceVal,
			define.BuiltinMetaKeyGroup:      groupName,
		}

		metaList := make([]map[string]any, 0, len(builtinMetaSchemaList)+len(schemaByLib[libId]))
		for _, b := range builtinMetaSchemaList {
			isShow := 0
			switch b.Key {
			case define.BuiltinMetaKeySource:
				isShow = cast.ToInt(showCfg[`show_meta_source`])
			case define.BuiltinMetaKeyUpdateTime:
				isShow = cast.ToInt(showCfg[`show_meta_update_time`])
			case define.BuiltinMetaKeyCreateTime:
				isShow = cast.ToInt(showCfg[`show_meta_create_time`])
			case define.BuiltinMetaKeyGroup:
				isShow = cast.ToInt(showCfg[`show_meta_group`])
			}
			if isShow != define.SwitchOn {
				continue
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
		for _, s := range schemaByLib[libId] {
			if cast.ToInt(s[`is_show`]) != define.SwitchOn {
				continue
			}
			k := strings.TrimSpace(s[`key`])
			if k == `` {
				continue
			}
			val, ok := metaMap[k]
			if !ok {
				val = ``
			}
			metaList = append(metaList, map[string]any{
				`id`:         cast.ToInt(s[`id`]),
				`library_id`: cast.ToInt(s[`library_id`]),
				`name`:       s[`name`],
				`key`:        k,
				`type`:       cast.ToInt(s[`type`]),
				`value`:      val,
				`is_show`:    cast.ToInt(s[`is_show`]),
				`is_builtin`: 0,
			})
		}
		tempItem[`meta_list`] = metaList

		formatedList = append(formatedList, tempItem)
	}

	exceptionTotal, err := exception.Count(`1`)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	data := map[string]any{`info`: info, `list`: formatedList, `total`: total, `exception_total`: exceptionTotal, `page`: page, `size`: size}
	return data, 0, nil
}

type BridgeSaveParagraphReq struct {
	Id               string   `form:"id"`
	FileId           string   `form:"file_id"`
	Title            string   `form:"title"`
	Content          string   `form:"content"`
	Question         string   `form:"question"`
	Answer           string   `form:"answer"`
	SimilarQuestions string   `form:"similar_questions"`
	Images           []string `form:"images"`
	CategoryId       string   `form:"category_id"`
	GroupId          string   `form:"group_id"`
	LibraryId        string   `form:"library_id"`
	ImagesJson       string   `form:"images_json"`
	Token            string   `url:"token"`
	HeaderToken      string
}

func BridgeSaveParagraph(adminUserId, loginUserId int, lang string, req *BridgeSaveParagraphReq) (msql.Params, int, error) {
	id := cast.ToInt64(req.Id)
	fileId := cast.ToInt64(req.FileId)
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	question := strings.TrimSpace(req.Question)
	answer := strings.TrimSpace(req.Answer)
	similarQuestions := strings.TrimSpace(req.SimilarQuestions)
	images := req.Images
	categoryId := cast.ToInt(req.CategoryId)
	groupId := cast.ToUint(req.GroupId)
	if id < 0 || fileId < 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	if id > 0 {
		fileIdStr, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`file_id`)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		if cast.ToUint(fileIdStr) == 0 { //paragraph segment does not exist
			return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `id`))
		}
		fileId = cast.ToInt64(fileIdStr)
	}
	if fileId == 0 { //no file specified; create a default custom document
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
		token := req.HeaderToken
		if len(token) == 0 {
			token = req.Token
		}
		fileId, err = getLibraryDefaultFile(lang, libraryId, adminUserId, token)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	fileInfo, err := common.GetLibFileInfo(int(fileId), adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(fileInfo) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	if cast.ToInt(fileInfo[`is_qa_doc`]) == define.DocTypeQa {
		if len(question) < 1 || utf8.RuneCountInString(question) > common.MaxContent {
			return nil, -1, errors.New(i18n.Show(lang, `length_error`))
		}
		if len(answer) < 1 || utf8.RuneCountInString(answer) > common.MaxContent {
			return nil, -1, errors.New(i18n.Show(lang, `length_error`))
		}
	} else {
		if len(content) < 1 || utf8.RuneCountInString(content) > common.MaxContent {
			return nil, -1, errors.New(i18n.Show(lang, `length_error`))
		}
	}
	jsonImages, err := common.CheckLibraryImage(images)
	if err != nil {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `images`))
	}
	if imagesJson := strings.TrimSpace(req.ImagesJson); len(imagesJson) > 0 {
		jsonImages = imagesJson //for special scenarios, pass through parameters directly
	}

	_ = m.Begin()
	data := msql.Datas{
		`admin_user_id`: adminUserId,
		`library_id`:    fileInfo[`library_id`],
		`file_id`:       fileId,
		`title`:         title,
		`images`:        jsonImages,
		`category_id`:   categoryId,
		`group_id`:      groupId,
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
			if fatherChunkParagraphNumber, number, sqlErr := common.GetAddParagraphNumbers(fileId); sqlErr == nil {
				data[`father_chunk_paragraph_number`] = fatherChunkParagraphNumber
				data[`number`] = number
			} else {
				_ = m.Rollback()
				return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
			}
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		vectorID, err := common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeQuestion), question)
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		vectorIds = append(vectorIds, vectorID)
		similarQuestionArr := make([]string, 0)
		tool.JsonDecode(similarQuestions, &similarQuestionArr)
		if err = common.DeleteLibraryFileDataIndex(cast.ToString(id), cast.ToString(define.VectorTypeSimilarQuestion)); err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		for _, similarQuestion := range similarQuestionArr {
			vectorID, err := common.SaveVector(
				cast.ToInt64(adminUserId),
				cast.ToInt64(fileInfo[`library_id`]),
				fileId,
				id,
				cast.ToString(define.VectorTypeSimilarQuestion),
				strings.TrimSpace(similarQuestion),
			)
			if err != nil {
				logs.Error(err.Error())
				_ = m.Rollback()
				return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
			}
			vectorIds = append(vectorIds, vectorID)
		}

		if fileInfo[`type`] == cast.ToString(define.QAIndexTypeQuestionAndAnswer) {
			vectorID, err = common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
				fileId, id, cast.ToString(define.VectorTypeAnswer), question)
			if err != nil {
				logs.Error(err.Error())
				_ = m.Rollback()
				return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
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
			if fatherChunkParagraphNumber, number, sqlErr := common.GetAddParagraphNumbers(fileId); sqlErr == nil {
				data[`father_chunk_paragraph_number`] = fatherChunkParagraphNumber
				data[`number`] = number
			} else {
				_ = m.Rollback()
				return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
			}
			id, err = m.Insert(data, `id`)
		}
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		vectorID, err := common.SaveVector(int64(adminUserId), cast.ToInt64(fileInfo[`library_id`]),
			fileId, id, cast.ToString(define.VectorTypeParagraph), content)
		if err != nil {
			logs.Error(err.Error())
			_ = m.Rollback()
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		vectorIds = append(vectorIds, vectorID)
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
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
	info, err := m.Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	return info, 0, nil
}

func BridgeDeleteParagraph(adminUserId int, ids, lang string) error {
	if len(ids) <= 0 {
		return errors.New(i18n.Show(lang, `param_lack`))
	}

	data, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, `in`, cast.ToString(ids)).Find()
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(data) == 0 {
		return errors.New(i18n.Show(lang, `no_data`))
	}

	if cast.ToInt(data[`category_id`]) > 0 {
		_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, `in`, cast.ToString(ids)).Update(msql.Datas{"isolated": true})
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `no_data`))
		}
	} else {
		_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, `in`, cast.ToString(ids)).Delete()
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		// clear stat library data robot tip
		common.CleanStatLibraryDataTip(adminUserId, strings.Split(ids, `,`))
		_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`data_id`, `in`, cast.ToString(ids)).Delete()
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		if common.GetNeo4jStatus(adminUserId) {
			for _, id := range strings.Split(ids, `,`) {
				err = common.NewGraphDB(adminUserId).DeleteByData(cast.ToInt(id))
				if err != nil {
					logs.Error(err.Error())
					return errors.New(i18n.Show(lang, `sys_err`))
				}
			}
		}
	}
	return nil
}
