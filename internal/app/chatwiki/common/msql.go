// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/casbin"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func ToStringMap(data msql.Datas, adds ...any) msql.Params {
	params := msql.Params{}
	for key, val := range data {
		params[key] = cast.ToString(val)
	}
	for i := 0; i < len(adds); i = i + 2 {
		val := ``
		if len(adds) > i+1 {
			val = cast.ToString(adds[i+1])
		}
		params[cast.ToString(adds[i])] = val
	}
	return params
}

type RobotCacheBuildHandler struct{ RobotKey string }

func (h *RobotCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.robot_info.%s`, h.RobotKey)
}
func (h *RobotCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, h.RobotKey).Find()
}

func GetRobotInfo(robotKey string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotCacheBuildHandler{RobotKey: robotKey}, &result, time.Hour)
	return result, err
}

type RobotApiKeyCacheBuildHandler struct{ RobotKey string }

func (h *RobotApiKeyCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.robot_apikey.%s`, h.RobotKey)
}
func (h *RobotApiKeyCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`chat_ai_robot_apikey`, define.Postgres).Where(`robot_key`, h.RobotKey).Order("id desc").Select()
	if err == nil && len(data) > 0 {
		for _, item := range data {
			delete(item, "create_time")
			delete(item, "update_time")
			delete(item, "admin_user_id")
		}
	}
	return data, err
}

func GetRobotApikeyInfo(robotKey string) ([]msql.Params, error) {
	result := make([]msql.Params, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotApiKeyCacheBuildHandler{RobotKey: robotKey}, &result, time.Hour*24*7)
	return result, err
}

type CustomerCacheBuildHandler struct {
	Openid      string
	AdminUserId int
}

func (h *CustomerCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.customer_info.v2.%d.%s`, h.AdminUserId, h.Openid)
}
func (h *CustomerCacheBuildHandler) GetCacheData() (any, error) {
	m := msql.Model(`chat_ai_customer`, define.Postgres)
	customer, err := m.Where(`openid`, h.Openid).Where(`admin_user_id`, cast.ToString(h.AdminUserId)).Find()
	if err == nil && len(customer) > 0 {
		up := msql.Datas{}
		if len(customer[`name`]) == 0 {
			up[`name`] = `访客` + tool.Random(4)
		}
		if len(customer[`avatar`]) == 0 {
			up[`avatar`] = define.DefaultCustomerAvatar
		}
		if len(up) > 0 {
			_, _ = m.Where(`id`, customer[`id`]).Update(up)
			return h.GetCacheData()
		}
	}
	return customer, err
}

func GetCustomerInfo(openid string, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &CustomerCacheBuildHandler{Openid: openid, AdminUserId: adminUserId}, &result, time.Hour)
	return result, err
}

func InsertOrUpdateCustomer(openid string, adminUserId int, upData msql.Datas) {
	customer, err := GetCustomerInfo(openid, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if upData == nil {
		upData = make(msql.Datas)
	}
	m := msql.Model(`chat_ai_customer`, define.Postgres)
	if len(customer) == 0 { //new customer
		upData[`openid`] = openid
		upData[`admin_user_id`] = adminUserId
		upData[`create_time`] = tool.Time2Int()
		upData[`update_time`] = tool.Time2Int()
		_, err = m.Insert(upData)
	} else {
		delete(upData, `is_background`) //first effect
		if len(upData) == 0 {
			return
		}
		upData[`update_time`] = tool.Time2Int()
		_, err = m.Where(`id`, customer[`id`]).Update(upData)
	}
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &CustomerCacheBuildHandler{Openid: openid, AdminUserId: adminUserId})
}

type LibraryCacheBuildHandler struct{ LibraryId int }

func (h *LibraryCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.library_info.%d`, h.LibraryId)
}
func (h *LibraryCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(h.LibraryId)).Find()
}

func GetLibraryInfo(libraryId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibraryCacheBuildHandler{LibraryId: libraryId}, &result, time.Hour)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

type LibFileCacheBuildHandler struct{ FileId int }

func (h *LibFileCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.file_info.%d`, h.FileId)
}
func (h *LibFileCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(h.FileId)).Find()
}

func GetLibFileInfo(fileId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibFileCacheBuildHandler{FileId: fileId}, &result, time.Hour)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

func GetMatchLibraryParagraphByVectorSimilarity(robot msql.Params, openid, appType, question string, libraryIds string, size int, similarity float64, searchType int) ([]msql.Params, error) {
	result := make([]msql.Params, 0)
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeVector}) {
		return result, nil
	}
	group := make(map[int]map[string][]string)
	for _, libraryId := range strings.Split(libraryIds, `,`) {
		library, err := GetLibraryInfo(cast.ToInt(libraryId), 0)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(library) == 0 {
			continue
		}
		modelConfigId := cast.ToInt(library[`model_config_id`])
		if _, ok := group[modelConfigId]; !ok {
			group[modelConfigId] = make(map[string][]string)
		}
		useModel := library[`use_model`]
		if _, ok := group[modelConfigId][useModel]; !ok {
			group[modelConfigId][useModel] = make([]string, 0)
		}
		group[modelConfigId][useModel] = append(group[modelConfigId][useModel], libraryId)
	}
	list := make(define.SimilarityResult, 0)
	if len(group) == 0 {
		return result, nil
	}
	wg := &sync.WaitGroup{}
	for modelConfigId := range group {
		for useModel, libraryIds := range group[modelConfigId] {
			wg.Add(1)
			go func(wg *sync.WaitGroup, robot msql.Params, openid, appType string, modelConfigId int, useModel, question string, libraryIds string, size int, list *define.SimilarityResult) {
				defer wg.Done()
				embedding, err := GetVector2000(cast.ToInt(robot[`admin_user_id`]), openid, robot, msql.Params{}, msql.Params{}, modelConfigId, useModel, question)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				embeddingArr := strings.Split(embedding, ",")
				subList, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
					Alias("a").
					Join(`chat_ai_library_file_data_index b`, `a.id=b.data_id`, `left`).
					Where(`a.library_id`, `in`, libraryIds).
					Where(`b.status`, cast.ToString(define.VectorStatusConverted)).
					Where("vector_dims(b.embedding)", cast.ToString(len(embeddingArr))).
					Field(`a.*`).
					Field(fmt.Sprintf(`max(1-(b.embedding<=>'%s')) as similarity`, embedding)).
					Order(`similarity desc`).
					Group(`a.id`).
					Limit(size).
					Select()
				if err != nil {
					logs.Error(err.Error())
					return
				}
				*list = append(*list, subList...)
			}(wg, robot, openid, appType, modelConfigId, useModel, question, strings.Join(libraryIds, `,`), size, &list)
		}
	}
	wg.Wait()
	//sort by similarity
	sort.Sort(list)
	//similarity verify
	for i, one := range list {
		if i >= size {
			break
		}
		if cast.ToFloat64(one[`similarity`]) < similarity {
			break
		}
		result = append(result, one)
	}
	return result, nil
}

func GetMatchLibraryParagraphByFullTextSearch(question, libraryIds string, size int, similarity float64, searchType int) ([]msql.Params, error) {
	list := make([]msql.Params, 0)
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeFullText}) {
		return list, nil
	}
	question = strings.ReplaceAll(question, `'`, ` `)
	queryTokens, err := msql.Model(fmt.Sprintf(`ts_parse('zhparser', '%s')`, question), define.Postgres).ColumnArr(`token`)
	if err != nil {
		return nil, err
	}

	ids, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, `in`, libraryIds).
		Where(fmt.Sprintf(`to_tsvector('zhima_zh_parser',upper(content))@@to_tsquery('zhima_zh_parser',upper('%s'))`, strings.Join(queryTokens, " | "))).
		Limit(500).ColumnArr(`id`)
	if err != nil {
		return list, err
	}
	if len(ids) == 0 {
		return list, nil
	}

	list, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Alias("a").
		Join("chat_ai_library_file_data b", "a.data_id=b.id", "left").
		Where(`a.id`, `in`, strings.Join(ids, `,`)).
		Where(`b.id is not null`).
		Field(`b.*,a.id as index_id`).
		Field(fmt.Sprintf(`ts_rank(to_tsvector('zhima_zh_parser',upper(a.content)),to_tsquery('zhima_zh_parser',upper('%s'))) as rank`, strings.Join(queryTokens, " | "))).
		Order(`rank DESC`).Limit(size).Select()
	if err != nil {
		return nil, err
	}

	listIds := make([]string, 0)
	for _, one := range list {
		listIds = append(listIds, cast.ToString(one[`index_id`]))
	}

	answerTokensResult, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Alias(`a`).
		Join(`LATERAL ts_parse('zhparser', a.content) as b`, `true`, `LEFT`).
		Where(`id`, `in`, strings.Join(listIds, `,`)).
		Field(`a.id, string_agg(b.token, ',') AS tokens`).
		Group(`a.id`).Select()
	if err != nil {
		return nil, err
	}

	similarities := make(map[int]float64)
	for _, one := range answerTokensResult {
		answerTokens := strings.Split(one[`tokens`], ",")
		score := overlapCoefficient(queryTokens, answerTokens)
		similarities[cast.ToInt(one[`id`])] = score
	}

	// add similarity field
	var result []msql.Params
	bestScores := make(map[interface{}]msql.Params) // unique
	for _, one := range list {
		score := similarities[cast.ToInt(one[`index_id`])]
		if score < similarity {
			continue
		}
		id := one[`id`]
		if existing, exists := bestScores[id]; !exists || cast.ToFloat64(existing[`similarity`]) < score {
			one[`similarity`] = cast.ToString(score)
			bestScores[id] = one
		}
	}

	// convert map to slice
	for _, one := range bestScores {
		result = append(result, one)
	}

	// sort
	sort.Slice(result, func(i, j int) bool {
		similarityI := cast.ToFloat64(result[i]["similarity"])
		similarityJ := cast.ToFloat64(result[j]["similarity"])
		return similarityI > similarityJ
	})

	return result, err
}

func GetMatchLibraryParagraphByMergeRerank(openid, appType, question string, size int, vectorList, searchList []msql.Params, robot msql.Params) ([]msql.Params, error) {
	if len(robot) == 0 || cast.ToInt(robot[`rerank_status`]) == 0 {
		return nil, nil //not rerank config
	}
	//merge and remove duplication
	ms := map[string]struct{}{}
	for i := range vectorList {
		ms[vectorList[i][`id`]] = struct{}{}
	}
	list := vectorList
	for i := range searchList {
		if _, ok := ms[searchList[i][`id`]]; ok {
			continue //duplication skip
		}
		ms[searchList[i][`id`]] = struct{}{}
		list = append(list, searchList[i])
	}
	if len(list) == 0 {
		return nil, nil
	}
	// Rerank resorted
	chunks := make([]string, 0)
	for _, one := range list {
		chunks = append(chunks, one[`content`])
	}
	rerankReq := &adaptor.ZhimaRerankReq{
		Enable:   true,
		Query:    question,
		Passages: chunks,
		Data:     list,
		TopK:     size,
	}
	return RerankData(cast.ToInt(robot[`admin_user_id`]), openid, appType, robot, rerankReq)
}

func GetMatchLibraryParagraphList(openid, appType, question string, optimizedQuestions []string, libraryIds string, size int, similarity float64, searchType int, robot msql.Params) (_ []msql.Params, libUseTime LibUseTime, _ error) {
	result := make([]msql.Params, 0)
	if len(libraryIds) == 0 {
		return result, libUseTime, nil
	}
	if len(question) == 0 {
		return nil, libUseTime, errors.New(`question cannot be empty`)
	}

	fetchSize := 4 * size
	var vectorList, searchList []msql.Params

	temp := time.Now()
	for _, q := range append(optimizedQuestions, question) {
		list, err := GetMatchLibraryParagraphByVectorSimilarity(robot, openid, appType, q, libraryIds, fetchSize, similarity, searchType)
		if err != nil {
			logs.Error(err.Error())
		}
		vectorList = append(vectorList, list...)
		list, err = GetMatchLibraryParagraphByFullTextSearch(q, libraryIds, fetchSize, similarity, searchType)
		if err != nil {
			logs.Error(err.Error())
		}
		searchList = append(searchList, list...)
	}
	libUseTime.RecallTime = time.Now().Sub(temp).Milliseconds()

	//问题优化后,召回的内容重新按相似度排序
	sort.Slice(vectorList, func(i, j int) bool {
		return cast.ToFloat64(vectorList[i][`similarity`]) > cast.ToFloat64(vectorList[j][`similarity`])
	})
	sort.Slice(searchList, func(i, j int) bool {
		return cast.ToFloat64(searchList[i][`similarity`]) > cast.ToFloat64(searchList[j][`similarity`])
	})

	temp = time.Now()
	rerankList, err := GetMatchLibraryParagraphByMergeRerank(openid, appType, question, fetchSize, vectorList, searchList, robot)
	libUseTime.RerankTime = time.Now().Sub(temp).Milliseconds()
	if err != nil {
		logs.Error(err.Error())
	}
	//RRF sort
	list := (&RRF{}).
		Add(DataSource{List: vectorList, Key: `id`, Fixed: 60}).
		Add(DataSource{List: searchList, Key: `id`, Fixed: 60}).
		Add(DataSource{List: rerankList, Key: `id`, Fixed: 58}).Sort()

	//return
	for i, one := range list {
		if i >= size {
			break
		}
		//replenish file info
		fileInfo, _ := GetLibFileInfo(cast.ToInt(one[`file_id`]), 0)
		one[`file_name`] = fileInfo[`file_name`]
		result = append(result, one)
	}
	return result, libUseTime, nil
}

type DialogueCacheBuildHandler struct{ DialogueId int }

func (h *DialogueCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.dialogue_info.%d`, h.DialogueId)
}
func (h *DialogueCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_dialogue`, define.Postgres).Where(`id`, cast.ToString(h.DialogueId)).Find()
}

func GetDialogueInfo(dialogueId, adminUserId, robotId int, openid string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &DialogueCacheBuildHandler{DialogueId: dialogueId}, &result, time.Hour)
	if err == nil && ((adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId) ||
		(robotId != 0 && cast.ToInt(result[`robot_id`]) != robotId) || (len(openid) != 0 && result[`openid`] != openid)) {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

type ModelConfigCacheBuildHandler struct{ ModelConfigId int }

func (h *ModelConfigCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.model_config.%d`, h.ModelConfigId)
}
func (h *ModelConfigCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_model_config`, define.Postgres).Where(`id`, cast.ToString(h.ModelConfigId)).Find()
}

func GetModelConfigInfo(modelId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &ModelConfigCacheBuildHandler{ModelConfigId: modelId}, &result, time.Hour*12)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

func GetDefaultLlmConfig(adminUserId int) (int, string, bool) {
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		return 0, ``, false
	}
	for _, config := range configs {
		if !tool.InArrayString(Llm, strings.Split(config[`model_types`], `,`)) {
			continue
		}
		modelInfo, ok := GetModelInfoByDefine(config[`model_define`])
		if ok && len(modelInfo.LlmModelList) > 0 {
			return cast.ToInt(config[`id`]), modelInfo.LlmModelList[0], true
		}
	}
	return 0, ``, false
}

func SaveVector(adminUserID, libraryID, fileID, dataID int64, vectorType, content string) (int64, error) {
	m := msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	info, err := m.
		Where(`data_id`, cast.ToString(dataID)).
		Where(`type`, vectorType).
		Field(`id,content`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	if len(info) == 0 {
		id, err := m.Insert(msql.Datas{
			`admin_user_id`: adminUserID,
			`library_id`:    libraryID,
			`file_id`:       fileID,
			`data_id`:       dataID,
			`type`:          vectorType,
			`content`:       content,
			`status`:        define.VectorStatusInitial,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}, `id`)
		if err != nil {
			logs.Error(err.Error())
			return 0, err
		}
		return id, nil
	} else {
		if info[`content`] == content {
			return 0, nil
		} else {
			_, err = m.
				Where(`id`, info[`id`]).
				Update(msql.Datas{
					`status`:  define.VectorStatusInitial,
					`errmsg`:  ``,
					`content`: content,
				})
			if err != nil {
				logs.Error(err.Error())
				return 0, err
			}
			return cast.ToInt64(info[`id`]), nil
		}
	}
}

func GetOptimizedQuestions(param *define.ChatRequestParam, contextList []map[string]string) ([]string, error) {
	histories := ""
	for _, item := range contextList {
		histories += "Q: " + item[`question`] + "\n"
		histories += "A: " + item[`answer`] + "\n"
	}
	prompt := strings.ReplaceAll(define.PromptDefaultQuestionOptimize, `{{query}}`, param.Question)
	prompt = strings.ReplaceAll(prompt, `{{histories}}`, histories)

	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}

	var result []string
	chatResp, _, err := RequestChat(
		param.AdminUserId,
		param.Openid,
		param.Robot,
		param.AppType,
		cast.ToInt(param.Robot[`model_config_id`]),
		param.Robot[`use_model`],
		messages,
		nil,
		cast.ToFloat32(param.Robot[`temperature`]),
		200,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(chatResp.Result), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// overlap coefficient
func overlapCoefficient(setA, setB []string) float64 {
	setAMap := make(map[string]bool)
	setBMap := make(map[string]bool)

	for _, item := range setA {
		setAMap[item] = true
	}
	for _, item := range setB {
		setBMap[item] = true
	}

	intersectionSize := 0
	for item := range setAMap {
		if setBMap[item] {
			intersectionSize++
		}
	}

	minSize := min(len(setAMap), len(setBMap))

	if minSize == 0 {
		return 0
	}

	return float64(intersectionSize) / float64(minSize)
}

func ClientSideNeedLogin(adminUserId int) bool {
	info, err := msql.Model(define.TableUser, define.Postgres).Where(`id`, cast.ToString(adminUserId)).
		Where(`is_deleted`, define.Normal).Field(`client_side_login_switch`).Find()
	if err != nil {
		logs.Error(err.Error())
	}
	if len(info) == 0 {
		return true
	}
	return cast.ToInt(info[`client_side_login_switch`]) == define.SwitchOn
}

func CheckPermission(userId int, permission string) bool {
	userRoles, err := msql.Model(define.TableUser, define.Postgres).Where(`id`, cast.ToString(userId)).
		Where(`is_deleted`, define.Normal).Value(`user_roles`)
	if err != nil {
		logs.Error(err.Error())
	}
	if len(userRoles) == 0 {
		return false
	}
	rules, err := casbin.Handler.GetPolicyForUser(userRoles)
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	rolePermission := make([]string, 0)
	for _, rule := range rules {
		if len(rule) > 1 {
			if strings.ContainsAny(rule[1], `/`) {
				continue
			}
			rolePermission = append(rolePermission, rule[1])
		}
	}
	return tool.InArrayString(permission, rolePermission)
}

type UploadFormFileHandler struct{ TaskId string }

func (h *UploadFormFileHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.upload_form_file_proc.%s`, h.TaskId)
}

func (h *UploadFormFileHandler) GetCacheData() (any, error) {
	return `{}`, nil
}

func GetUploadFormFileProc(taskId string) (*define.UploadFormFile, error) {
	result := new(define.UploadFormFile)
	err := lib_redis.GetCacheWithBuild(define.Redis, &UploadFormFileHandler{TaskId: taskId}, &result, time.Hour)
	return result, err
}
func SetUploadFormFileProc(taskId string, uploadForm *define.UploadFormFile, ttl time.Duration) error {
	handler := UploadFormFileHandler{TaskId: taskId}
	str, _ := json.Marshal(uploadForm)
	_, err := define.Redis.Set(context.Background(), handler.GetCacheKey(), string(str), time.Second*ttl).Result()
	return err
}

type NodeCacheBuildHandler struct {
	RobotId  uint
	DataType uint
	NodeKey  string
}

func (h *NodeCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.work.flow.node..%d.%d.%s`, h.RobotId, h.DataType, h.NodeKey)
}
func (h *NodeCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`work_flow_node`, define.Postgres).Where(`robot_id`, cast.ToString(h.RobotId)).
		Where(`data_type`, cast.ToString(h.DataType)).Where(`node_key`, h.NodeKey).Find()
}
func GetRobotNode(robotId uint, nodeKey string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &NodeCacheBuildHandler{RobotId: robotId, DataType: define.DataTypeRelease, NodeKey: nodeKey}, &result, time.Hour)
	return result, err
}
