// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/pkg/lib_redis"
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

	// ====== 元数据及其值（列表也返回）======
	// 内置 is_show：来自 chat_ai_library
	libraryShow, _ := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(libraryId)).
		Field(`show_meta_source,show_meta_update_time,show_meta_create_time,show_meta_group`).
		Find()

	// 取该知识库的自定义 schema 列表（前端直接渲染表格）
	schemaList, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Order(`id asc`).
		Field(`id,library_id,name,key,type,is_show`).
		Select()
	schemaKeyList := make([]string, 0, len(schemaList))
	schemaMap := make(map[string]msql.Params, len(schemaList))
	if err == nil {
		for _, one := range schemaList {
			k := strings.TrimSpace(one[`key`])
			if k == `` {
				continue
			}
			// 按 is_show 控制是否返回
			if cast.ToInt(one[`is_show`]) != define.SwitchOn {
				continue
			}
			schemaKeyList = append(schemaKeyList, k)
			schemaMap[k] = one
		}
	}
	// 分组 id -> name 映射
	groupMap := map[int]string{0: `未分组`}
	groupList, err := msql.Model(`chat_ai_library_group`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Where(`group_type`, cast.ToString(define.LibraryGroupTypeFile)).
		Field(`id,group_name`).
		Select()
	if err == nil {
		for _, g := range groupList {
			groupMap[cast.ToInt(g[`id`])] = g[`group_name`]
		}
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

	// 给列表每条记录补齐 group_name/meta_list（内置+自定义）
	listAny := make([]map[string]any, 0, len(list))
	for _, item := range list {
		obj := make(map[string]any, len(item)+4)
		for k, v := range item {
			obj[k] = v
		}
		gid := cast.ToInt(item[`group_id`])
		gname, ok := groupMap[gid]
		if !ok {
			gname = `未分组`
		}
		obj[`group_name`] = gname

		metaMap := make(map[string]any)
		metaStr := strings.TrimSpace(cast.ToString(item[`metadata`]))
		if metaStr == `` {
			metaStr = `{}`
		}
		_ = tool.JsonDecode(metaStr, &metaMap)

		// meta_list：前端表格直接渲染（新功能不考虑兼容，不返回 meta_values）
		metaList := make([]map[string]any, 0, len(define.BuiltinMetaSchemaList)+len(schemaKeyList))
		builtinValueMap := map[string]any{
			define.BuiltinMetaKeyUpdateTime: cast.ToInt(item[`update_time`]),
			define.BuiltinMetaKeyCreateTime: cast.ToInt(item[`create_time`]),
			define.BuiltinMetaKeySource:     cast.ToInt(item[`doc_type`]),
			define.BuiltinMetaKeyGroup:      gname,
		}
		for _, b := range define.BuiltinMetaSchemaList {
			isShow := 0
			switch b.Key {
			case define.BuiltinMetaKeySource:
				isShow = cast.ToInt(libraryShow[`show_meta_source`])
			case define.BuiltinMetaKeyUpdateTime:
				isShow = cast.ToInt(libraryShow[`show_meta_update_time`])
			case define.BuiltinMetaKeyCreateTime:
				isShow = cast.ToInt(libraryShow[`show_meta_create_time`])
			case define.BuiltinMetaKeyGroup:
				isShow = cast.ToInt(libraryShow[`show_meta_group`])
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
		for _, k := range schemaKeyList {
			s := schemaMap[k]
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
		obj[`meta_list`] = metaList
		listAny = append(listAny, obj)
	}

	data := map[string]any{`info`: info, `list`: listAny, `count_data`: countData, `total`: total, `page`: page, `size`: size}
	return data, 0, nil
}

type BridgeAddLibraryFileReq struct {
	LibraryId                 string `form:"library_id"`
	LibraryKey                string `form:"library_key"`
	DocType                   string `form:"doc_type"`
	Urls                      string `form:"urls"`
	FileName                  string `form:"file_name"`
	Content                   string `form:"content"`
	Title                     string `form:"title"`
	IsQaDoc                   string `form:"is_qa_doc"`
	QaIndexType               string `form:"qa_index_type"`
	DocAutoRenewFrequency     string `form:"doc_auto_renew_frequency"`
	DocAutoRenewMinute        string `form:"doc_auto_renew_minute"`
	AnswerLable               string `form:"answer_lable"`
	AnswerColumn              string `form:"answer_column"`
	QuestionLable             string `form:"question_lable"`
	QuestionColumn            string `form:"question_column"`
	SimilarColumn             string `form:"similar_column"`
	SimilarLabel              string `form:"similar_label"`
	PdfParseType              string `form:"pdf_parse_type"`
	GroupId                   string `form:"group_id"`
	OfficialArticleId         string `form:"official_article_id"`
	OfficialArticleUpdateTime int64  `form:"official_article_update_time"`
	FeishuDocumentIdList      string `form:"feishu_document_id_list"`
	FeishuAppId               string `form:"feishu_app_id"`
	FeishuAppSecret           string `form:"feishu_app_secret"`
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
	fileIds, err := addLibFile(c.Request.MultipartForm, common.GetLang(c), adminUserId, libraryId, cast.ToInt(info[`type`]), chunkParam, req)
	if err != nil {
		return nil, -1, err
	}
	return map[string]any{`file_ids`: fileIds}, 0, nil
}

func BridgeDelLibraryFile(adminUserId int, ids, lang string) error {
	for _, id := range strings.Split(ids, `,`) {
		id := cast.ToInt(id)
		if id <= 0 {
			return errors.New(i18n.Show(lang, `param_lack`))
		}
		info, err := common.GetLibFileInfo(id, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		if len(info) == 0 {
			return errors.New(i18n.Show(lang, `no_data`))
		}
		_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(id)).Delete()
		if err != nil {
			logs.Error(err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
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

		if common.GetNeo4jStatus(adminUserId) {
			for _, dataId := range dataIdList {
				err = common.NewGraphDB(adminUserId).DeleteByData(cast.ToInt(dataId))
				if err != nil {
					logs.Error(err.Error())
				}
			}
		}
	}
	return nil
}
