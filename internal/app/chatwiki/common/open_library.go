// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func SaveLibDoc(adminUserId, userId, libraryId, docId, fileId, pid, isIndex, draft int, title string, info msql.Params, content string, isDir int, docIcon string) (int, error) {
	var (
		err     error
		summary string
	)
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	// save file doc info
	updateData := msql.Datas{
		`pid`:         pid,
		`file_id`:     fileId,
		`title`:       title,
		`is_draft`:    draft,
		`content`:     content,
		`summary`:     summary,
		`update_time`: tool.Time2Int(),
		`is_dir`:      isDir,
		`doc_icon`:    docIcon,
	}
	updateData[`edit_user`] = userId
	if isIndex == define.LibDocIndex {
		updateData[`is_pub`] = define.IsPub
		if draft == define.IsDraft {
			delete(updateData, `content`)
			updateData[`draft_content`] = content
		} else {
			updateData[`draft_content`] = ""
			updateData[`edit_user`] = 0
		}
	} else {
		if draft == define.IsDraft {
			delete(updateData, `content`)
			updateData[`draft_content`] = content
		} else {
			updateData[`is_pub`] = define.IsPub
			updateData[`draft_content`] = ""
			updateData[`edit_user`] = 0
		}
	}
	if fileId == 0 {
		delete(updateData, `file_id`)
	}
	if docId <= 0 {
		updateData[`admin_user_id`] = adminUserId
		updateData[`library_id`] = libraryId
		updateData[`doc_key`] = BuildDocKey(adminUserId)
		updateData[`sort`] = time.Now().UnixMilli()
		updateData[`is_index`] = isIndex
		updateData[`create_time`] = tool.Time2Int()
		id, err := m.Insert(updateData, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
		docId = int(id)
	} else {
		_, err = m.Where(`id`, cast.ToString(docId)).Update(updateData)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return docId, err
}

func ChangeLibDoc(docId int, updateData msql.Datas) (int, error) {
	var (
		err error
	)
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	updateData[`update_time`] = tool.Time2Int()
	_, err = m.Where(`id`, cast.ToString(docId)).Update(updateData)
	if err != nil {
		logs.Error(err.Error())
	}
	return docId, err
}

func DeleteLibDocInfo(docId int, info msql.Params) error {
	if docId <= 0 {
		return nil
	}
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	_, err := m.Where(`id`, cast.ToString(docId)).Update(msql.Datas{
		`delete_time`: tool.Time2Int(),
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		return err
	}
	return err
}

func DeleteLibDocRelationData(id int) error {
	// update status
	_, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	logs.Info(`DeleteLibDocRelationData:%v,%s`, id, err)
	if err != nil {
		return err
	}
	_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`file_id`, cast.ToString(id)).Delete()
	logs.Info(`DeleteLibDocRelationData:%v,%s`, id, err)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`file_id`, cast.ToString(id)).Delete()
	logs.Info(`DeleteLibDocRelationData:%v,%s`, id, err)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	// delete cache
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: id})
	return err
}

type LibDocCacheBuildHandler struct{ DocKey string }

func (h *LibDocCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.doc_info.v1.%s`, h.DocKey)
}
func (h *LibDocCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_library_file_doc`, define.Postgres).Where(`doc_key`, h.DocKey).Find()
}

func GetLibDocInfoByDocKey(docKey string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibDocCacheBuildHandler{DocKey: docKey}, &result, time.Hour*3)
	return result, err
}

func GetLibDocInfo(docId int) msql.Params {
	info := msql.Params{}
	if docId <= 0 {
		return info
	}
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	info, _ = m.Where(`id`, cast.ToString(docId)).Find()
	return info
}

func GetLibDocIndex(libraryId int) msql.Params {
	info := msql.Params{}
	if libraryId <= 0 {
		return info
	}
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	info, _ = m.Where(`library_id`, cast.ToString(libraryId)).Where(`is_index`, `1`).Find()
	return info
}

func GetLibDocFields(docId int, fields string) msql.Params {
	info := msql.Params{}
	if docId <= 0 {
		return info
	}
	if fields == "" {
		fields = `*`
	}
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	info, _ = m.Where(`id`, cast.ToString(docId)).Field(fields).Find()
	return info
}

func GetLibDocInfoByFileId(docId int) msql.Params {
	info := msql.Params{}
	if docId <= 0 {
		return info
	}
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	info, _ = m.Where(`file_id`, cast.ToString(docId)).Field(`id,doc_key,title,pid,library_id,is_draft,is_pub,create_time,update_time,is_dir,doc_icon`).Find()
	return info
}

func GetLibDocChildren(pid []string, fetchAll bool) []msql.Params {
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	if !fetchAll {
		m.Where(`is_pub`, `1`)
	}
	children, _ := m.Where(`pid`, `in`, strings.Join(pid, `,`)).Where(`delete_time`, `0`).Field(`id,doc_key,pid,file_id,is_draft,sort,title,create_time,update_time,is_dir,doc_icon`).Order(`sort desc`).Select()
	return children
}

func GetLibDocAllChildren(pid []string, fetchAll bool) []msql.Params {
	children := GetLibDocChildren(pid, fetchAll)
	pids := []string{}
	for _, item := range children {
		pids = append(pids, item[`id`])
	}
	data := GetLibDocChildren(pids, fetchAll)
	children = append(children, data...)
	return children
}

type LibraryCatalog struct {
	Id         int               `json:"id"`
	DocKey     string            `json:"doc_key"`
	Pid        int               `json:"pid"`
	FileId     int               `json:"file_id"`
	Title      string            `json:"title"`
	Sort       int               `json:"sort"`
	DocIcon    string            `json:"doc_icon"`
	IsDir      int               `json:"is_dir"`
	CreateTime int               `json:"create_time"`
	UpdateTime int               `json:"update_time"`
	Children   []*LibraryCatalog `json:"children"`
	HasChild   int               `json:"has_child"`
	NextKey    string            `json:"nest_key"`
	PrevNode   *LibraryCatalog   `json:"prev_node"`
	NextNode   *LibraryCatalog   `json:"next_node"`
	Level      int               `json:"level"`
	HasBother  bool
}

type LibraryCatalogCacheBuildHandle struct {
	LibraryId  int
	PreviewKey string
}

func (h *LibraryCatalogCacheBuildHandle) GetCacheKey() string {
	if h.PreviewKey != "" {
		return fmt.Sprintf(`chatwiki.library_catalog.v1.preview.%d`, h.LibraryId)
	}
	return fmt.Sprintf(`chatwiki.library_catalog.v1.%d`, h.LibraryId)
}

func (h *LibraryCatalogCacheBuildHandle) GetCacheData() (any, error) {
	var (
		page = 1
		size = 500
		tree = make([]*LibraryCatalog, 0)
	)
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres)
	for {
		m.Where(`library_id`, cast.ToString(h.LibraryId)).Where(`delete_time`, `0`)
		if h.PreviewKey == "" {
			m.Where(`is_pub`, `1`)
		}
		info, _, err := m.Field(`id,doc_key,pid,is_index,file_id,sort,title,create_time,update_time,doc_icon,is_dir`).Order(`sort desc`).Paginate(page, size)
		if len(info) == 0 || err != nil {
			break
		}
		for _, item := range info {
			if cast.ToInt(item[`is_index`]) == define.LibDocIndex {
				continue
			}
			data := &LibraryCatalog{
				Id:         cast.ToInt(item[`id`]),
				DocKey:     item[`doc_key`],
				Pid:        cast.ToInt(item[`pid`]),
				FileId:     cast.ToInt(item[`file_id`]),
				Title:      item[`title`],
				Sort:       cast.ToInt(item[`sort`]),
				CreateTime: cast.ToInt(item[`create_time`]),
				UpdateTime: cast.ToInt(item[`update_time`]),
				DocIcon:    item[`doc_icon`],
				IsDir:      cast.ToInt(item[`is_dir`]),
				Children:   make([]*LibraryCatalog, 0),
			}
			tree = append(tree, data)
		}
		page++
	}
	tree = ConvertSliceToTree(tree, 0, 0)
	return tree, nil
}

// ConvertSliceToTree
func ConvertSliceToTree(nodes []*LibraryCatalog, parentId, level int) []*LibraryCatalog {
	var (
		tree []*LibraryCatalog
	)
	for _, node := range nodes {
		if node.Pid == parentId {
			node.Level = level
			node.Children = ConvertSliceToTree(nodes, node.Id, node.Level+1)
			node.HasChild = 1
			if len(node.Children) > 0 {
				node.NextNode = node.Children[0]
			}
			tree = append(tree, node)
		}
	}
	return tree
}

func FindPrevAndNext(nodes []*LibraryCatalog, docKey string) (prev *LibraryCatalog, next *LibraryCatalog) {
	root := &LibraryCatalog{Children: nodes}
	target := findNode(root, docKey)
	var list []*LibraryCatalog
	dfs(root, &list)
	prev = findPrevious(list, target)
	next = findNext(list, target)
	return prev, next
}

// 深度优先遍历，将节点按顺序存储到切片中
func dfs(root *LibraryCatalog, nodes *[]*LibraryCatalog) {
	if root == nil {
		return
	}
	*nodes = append(*nodes, root) // 将当前节点加入切片
	for _, child := range root.Children {
		dfs(child, nodes) // 递归遍历子节点
	}
}

// 查找上一个节点
func findPrevious(nodes []*LibraryCatalog, target *LibraryCatalog) *LibraryCatalog {
	for i, node := range nodes {
		if node == target && i > 0 {
			if nodes[i-1].Id <= 0 {
				break
			}
			return nodes[i-1]
		}
	}
	return nil
}

// 查找下一个节点
func findNext(nodes []*LibraryCatalog, target *LibraryCatalog) *LibraryCatalog {
	for i, node := range nodes {
		if node == target && i < len(nodes)-1 {
			if nodes[i+1].Id <= 0 {
				break
			}
			return nodes[i+1]
		}
	}
	return nil
}

// 查找节点
func findNode(root *LibraryCatalog, targetValue string) *LibraryCatalog {
	if root == nil {
		return nil
	}
	if root.DocKey == targetValue {
		return root
	}
	for _, child := range root.Children {
		found := findNode(child, targetValue)
		if found != nil {
			return found
		}
	}
	return nil
}

func GetLibDocCateLogByCache(pid int, preview string) ([]*LibraryCatalog, error) {
	result := make([]*LibraryCatalog, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibraryCatalogCacheBuildHandle{LibraryId: pid, PreviewKey: preview}, &result, time.Hour*3)
	return result, err
}

func GetLibDocCateLog(pid, libraryId int, fetchAll bool) []msql.Params {
	m := msql.Model(`chat_ai_library_file_doc`, define.Postgres).Where(`library_id`, cast.ToString(libraryId)).Where(`pid`, cast.ToString(pid)).Where(`delete_time`, `0`).Where(`is_index`, `0`)
	if !fetchAll {
		m.Where(`is_pub`, `1`)
	}
	info, _ := m.Field(`id,doc_key,pid,is_draft,file_id,sort,title,create_time,update_time,is_dir,doc_icon`).Order(`sort desc`).Select()
	docIds := []string{}
	for _, item := range info {
		docIds = append(docIds, item[`id`])
	}
	children := GetLibDocChildren(docIds, fetchAll)
	childrenMap := make(map[string]int)
	for _, item := range children {
		childrenMap[item[`pid`]]++
	}
	for _, item := range info {
		item[`children_num`] = `0`
		if num, ok := childrenMap[item[`id`]]; ok {
			item[`children_num`] = cast.ToString(num)
		}
	}
	return info
}

func LibDocSearch(lang string, libraryId int, search string, library msql.Params) ([]msql.Params, []adaptor.ZhimaChatCompletionMessage, []msql.Params, error) {
	// mixed search
	var (
		result                 = make([]msql.Params, 0)
		summary                = make([]adaptor.ZhimaChatCompletionMessage, 0)
		fetchSize              = 1000
		summarySie             = 10
		similarity             = 0.1
		searchType             = define.SearchTypeMixed
		vectorList, searchList []msql.Params
		quoteFile              = make([]msql.Params, 0)
	)
	if cast.ToInt(library[`use_model_switch`]) == define.SwitchOn {
		list, err := GetMatchLibraryParagraphByVectorSimilarity(lang, cast.ToInt(library[`admin_user_id`]), msql.Params{}, "", "", search, cast.ToString(libraryId), fetchSize, similarity, searchType)
		if err != nil {
			logs.Error(err.Error())
			return result, summary, quoteFile, err
		}
		vectorList = append(vectorList, list...)
	}
	list, err := GetMatchLibraryParagraphByFullTextSearch(search, cast.ToString(libraryId), fetchSize, searchType)
	if err != nil {
		logs.Error(err.Error())
		return result, summary, quoteFile, nil
	}
	searchList = append(searchList, list...)

	//RRF sort
	list = (&RRF{}).
		Add(DataSource{List: vectorList, Key: `id`, Fixed: 50, Weight: 70}).
		Add(DataSource{List: searchList, Key: `id`, Fixed: 50, Weight: 30}).Sort()

	// quoteFile
	summary = append(summary, adaptor.ZhimaChatCompletionMessage{
		Role:    `user`,
		Content: search,
	})
	for key, one := range list {
		if cast.ToInt(one[`number`]) != 1 {
			continue
		}
		//replenish file library
		docInfo := GetLibDocInfoByFileId(cast.ToInt(one[`file_id`]))
		if len(docInfo) <= 0 {
			continue
		}
		one[`title`] = docInfo[`title`]
		one[`pid`] = docInfo[`pid`]
		one[`doc_id`] = docInfo[`id`]
		one[`doc_key`] = docInfo[`doc_key`]
		one[`is_dir`] = docInfo[`is_dir`]
		one[`doc_icon`] = docInfo[`doc_icon`]
		one[`create_time`] = docInfo[`create_time`]
		one[`update_time`] = docInfo[`update_time`]
		if key < summarySie {
			summary = append(summary, adaptor.ZhimaChatCompletionMessage{
				Role:    `user`,
				Content: one[`content`],
			})
			quoteFile = append(quoteFile, msql.Params{
				`doc_id`:    one[`doc_id`],
				`doc_key`:   one[`doc_key`],
				`file_name`: docInfo[`title`],
			})
		}
		result = append(result, one)
	}
	return result, summary, quoteFile, nil
}

func LibDocAiSummary(lang string, libraryId int, search string, library msql.Params, chanStream chan sse.Event, isClose *bool) ([]msql.Params, error) {
	defer close(chanStream)
	chanStream <- sse.Event{Event: `ping`, Data: tool.Time2Int()}
	// mixed search
	var (
		temperature float32 = 0.5
		maxToken            = 2000
	)
	list, summary, quoteFile, err := LibDocSearch(lang, libraryId, search, library)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(lang, `sys_err`)}
		return list, err
	}

	// to ai summary
	if cast.ToInt(library[`ai_summary`]) == define.SwitchOn {
		chatResp, requestTime, err := RequestSearchStream(
			lang,
			cast.ToInt(library[`admin_user_id`]),
			cast.ToInt(library[`summary_model_config_id`]),
			strings.TrimSpace(library[`ai_summary_model`]),
			library,
			summary,
			[]adaptor.FunctionTool{},
			chanStream,
			temperature,
			maxToken,
		)
		logs.Info(`%v`, chatResp)
		logs.Info(`%v`, requestTime)
		if err != nil {
			logs.Error(err.Error())
			chanStream <- sse.Event{Event: `error`, Data: i18n.Show(lang, `sys_err`)}
		}
		//message push
		if len(quoteFile) > 0 {
			chanStream <- sse.Event{Event: `quote_file`, Data: quoteFile}
		}
		chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
	}
	return nil, nil
}

func BuildLibraryKey(id, createTime int) string {
	str := fmt.Sprintf("%s_%v_%d", `lib`, id, createTime)
	encodeStr := tool.Base64Encode(str)
	return encodeStr
}

func ParseLibraryKey(key string) string {
	data, err := tool.Base64Decode(key)
	if err != nil {
		return ""
	}
	message := strings.SplitN(data, "_", 3)
	if len(message) >= 3 {
		return cast.ToString(message[1])
	}
	return ""
}

func BuildDocKey(adminUserId int) string {
	str := fmt.Sprintf("%s_%v_%v", `doc`, int64(adminUserId)+time.Now().UnixMicro(), tool.Random(4))
	encodeStr := tool.MD5(str)
	return encodeStr
}

func SaveQuestionGuide(adminUserId, libraryId, docId int64, question string) (int64, error) {
	var (
		err error
	)
	m := msql.Model(`chat_ai_library_question_guide`, define.Postgres)
	if docId <= 0 {
		// save file doc info
		insData := msql.Datas{
			`admin_user_id`: adminUserId,
			`library_id`:    libraryId,
			`question`:      question,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}
		docId, err = m.Insert(insData, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
	} else {
		updateData := msql.Datas{
			`admin_user_id`: adminUserId,
			`library_id`:    libraryId,
			`question`:      question,
			`update_time`:   tool.Time2Int(),
		}
		_, err = m.Where(`id`, cast.ToString(docId)).Update(updateData)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return docId, err
}

func GetQuestionGuideList(id int) (list []msql.Params) {
	list, err := msql.Model(`chat_ai_library_question_guide`, define.Postgres).Where(`library_id`, cast.ToString(id)).Select()
	if err != nil {
		logs.Error(err.Error())
	}
	return
}

func GetPartnerInfo(userId, libraryId int) msql.Params {
	info := msql.Params{}
	m := msql.Model(`chat_ai_library_partner`, define.Postgres)
	info, _ = m.Where(`user_id`, cast.ToString(userId)).Where(`library_id`, cast.ToString(libraryId)).Find()
	return info
}

func SaveLibDocPartner(userId, libraryId, rights, typ int, userIds []string) error {
	m := msql.Model(`chat_ai_library_partner`, define.Postgres)
	if typ == 1 {
		for _, user := range userIds {
			datas := msql.Datas{
				`user_id`:        user,
				`library_id`:     libraryId,
				`operate_rights`: rights,
				`create_time`:    tool.Time2Int(),
				`update_time`:    tool.Time2Int(),
				`creator`:        userId,
			}
			_, err := m.Insert(datas)
			if err != nil {
				logs.Error(err.Error())
				return err
			}
		}
	} else {
		for _, user := range userIds {
			datas := msql.Datas{
				`operate_rights`: rights,
				`update_time`:    tool.Time2Int(),
			}
			_, err := m.Where(`library_id`, cast.ToString(libraryId)).Where(`user_id`, cast.ToString(user)).Update(datas)
			if err != nil {
				logs.Error(err.Error())
				return err
			}
		}
	}
	return nil
}

func DeleteLibDocPartner(id int) error {
	_, err := msql.Model(`chat_ai_library_partner`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	return err
}

func LibDocPartnerList(libraryIds []string, page, size int) ([]msql.Params, int, error) {
	list, total, err := msql.Model(`chat_ai_library_partner`, define.Postgres).Where(`library_id`, `in`, strings.Join(libraryIds, `,`)).Order(`id asc`).Paginate(page, size)
	var (
		userIds = []string{}
		userMap = map[string]msql.Params{}
	)
	if len(list) > 0 {
		for _, item := range list {
			userIds = append(userIds, item[`user_id`])
		}
		userData, err := msql.Model(define.TableUser, define.Postgres).Where("id", `in`, strings.Join(userIds, `,`)).Where("is_deleted", define.Normal).Field(`id,user_name,avatar,nick_name,user_roles`).Select()
		if err != nil {
			logs.Error(err.Error())
			return nil, 0, err
		}
		roles, err := msql.Model(define.TableRole, define.Postgres).Select()
		roleMap := make(map[string]msql.Params)
		for _, role := range roles {
			roleMap[role["id"]] = role
		}
		for _, user := range userData {
			user["role_name"] = ""
			user["role_type"] = ""
			if role, ok := roleMap[user["user_roles"]]; ok {
				user["role_name"] = role[`name`]
				user["role_type"] = role[`role_type`]
			}
			userMap[user[`id`]] = user
		}
		for _, item := range list {
			if user, ok := userMap[item[`user_id`]]; ok {
				item[`user_name`] = user[`user_name`]
				item[`nick_name`] = user[`nick_name`]
				item[`avatar`] = user[`avatar`]
				item[`role_type`] = user[`role_type`]
			}
		}
	}
	return list, total, err
}

type LibDocPreviewCacheBuildHandler struct {
	LibraryKey string
	PreviewKey string
}

func (h *LibDocPreviewCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.doc_preview_info.v1.%s.%s`, h.LibraryKey, h.PreviewKey)
}
func (h *LibDocPreviewCacheBuildHandler) GetCacheData() (any, error) {
	return nil, nil
}

func (h *LibDocPreviewCacheBuildHandler) SaveCacheData(value interface{}, expire time.Duration) error {
	_, err := define.Redis.Set(context.Background(), h.GetCacheKey(), value, expire).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return err
	}
	return nil
}

func CheckPreviewOpenDoc(libraryKey, previewKey string) bool {
	handler := &LibDocPreviewCacheBuildHandler{
		LibraryKey: libraryKey,
		PreviewKey: previewKey,
	}
	data, err := define.Redis.Get(context.Background(), handler.GetCacheKey()).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return false
	}
	if cast.ToInt(data) <= 0 {
		return false
	}
	return true
}

func GetOpenBindLibList(domain, typ string) ([]msql.Params, error) {
	var (
		err      error
		bindList = []msql.Params{}
	)
	list, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`type`, cast.ToString(typ)).
		Field(`id,library_name,avatar,access_rights,share_url,create_time`).Order(`id asc`).Select()
	for _, item := range list {
		if cast.ToInt(item[`access_rights`]) != define.OpenLibraryAccessRights {
			continue
		}
		// 去除http(s)前缀后进行域名匹配
		shareUrl := strings.TrimPrefix(strings.TrimPrefix(item[`share_url`], "https://"), "http://")
		domainTrim := strings.TrimPrefix(strings.TrimPrefix(domain, "https://"), "http://")
		if strings.Contains(shareUrl, domainTrim) {
			item[`library_key`] = BuildLibraryKey(cast.ToInt(item[`id`]), cast.ToInt(item[`create_time`]))
			bindList = append(bindList, item)
		}
	}
	return bindList, err
}
