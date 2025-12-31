// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/casbin"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/go-redis/redis/v8"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	return fmt.Sprintf(lib_define.RedisPrefixRobotInfo, h.RobotKey)
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
			if len(customer[`nickname`]) > 0 {
				up[`name`] = customer[`nickname`]
			} else {
				up[`name`] = `访客` + tool.Random(4)
			}
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
	data, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(h.LibraryId)).Find()
	if len(data) > 0 {
		data[`library_key`] = BuildLibraryKey(cast.ToInt(data[`id`]), cast.ToInt(data[`create_time`]))
	}
	return data, err
}

func GetLibraryInfo(libraryId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibraryCacheBuildHandler{LibraryId: libraryId}, &result, time.Hour)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

func GetLibraryData(libraryId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibraryCacheBuildHandler{LibraryId: libraryId}, &result, time.Hour)
	return result, err
}

type LibFileCacheBuildHandler struct{ FileId int }

func (h *LibFileCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.file_info.%d`, h.FileId)
}
func (h *LibFileCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(h.FileId)).Where(`delete_time`, `0`).Find()
}

func GetLibFileInfo(fileId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &LibFileCacheBuildHandler{FileId: fileId}, &result, time.Hour)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

func GetMatchLibraryParagraphByVectorSimilarity(adminUserId int, robot msql.Params, openid, appType, question string, libraryIds string, size int, similarity float64, searchType int) ([]msql.Params, error) {
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
			go func(wg *sync.WaitGroup, adminUserId int, robot msql.Params, openid, appType string, modelConfigId int, useModel, question string, libraryIds string, size int, list *define.SimilarityResult) {
				defer wg.Done()
				embedding, err := GetVector2000(adminUserId, openid, robot, msql.Params{}, msql.Params{}, modelConfigId, useModel, question)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				subList, err := VectorRecall(libraryIds, embedding, size)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				*list = append(*list, subList...)
			}(wg, adminUserId, robot, openid, appType, modelConfigId, useModel, question, strings.Join(libraryIds, `,`), size, &list)
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

func GetMatchLibraryParagraphByGraphSimilarity(robot msql.Params, openid, appType, question string, libraryIds string, size int, searchType int) ([]msql.Params, error) {
	result := make([]msql.Params, 0)
	if GetNeo4jStatus(cast.ToInt(robot[`admin_user_id`])) || !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeGraph}) {
		return result, nil
	}

	// Input validation
	if len(question) == 0 {
		logs.Error("Question is empty")
		return result, errors.New("question cannot be empty")
	}

	libraryIdList, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, robot[`admin_user_id`]).
		Where(`id`, `in`, libraryIds).
		Where(`graph_switch`, cast.ToString(define.SwitchOn)).
		ColumnArr("id")
	libraryIds = strings.Join(libraryIdList, ",")

	if len(libraryIds) == 0 {
		logs.Error("no enabled graph library")
		return result, errors.New("no enabled graph library")
	}

	if size <= 0 {
		size = 10 // Set default value
	}

	// 1. 从问题中提取实体
	extractEntitiesPrompt := strings.ReplaceAll(define.PromptDefaultEntityExtract, `{{question}}`, question)
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: extractEntitiesPrompt}}
	chatResp, _, err := RequestChat(
		cast.ToInt(robot[`admin_user_id`]),
		openid,
		robot,
		appType,
		cast.ToInt(robot[`model_config_id`]),
		robot[`use_model`],
		messages,
		nil,
		0.1,
		500,
	)
	if err != nil {
		logs.Error("Failed to extract entities: %s", err.Error())
		return result, err
	}

	// Clean and parse LLM response
	chatResp.Result = strings.TrimSpace(chatResp.Result)
	chatResp.Result = strings.TrimPrefix(chatResp.Result, "```json")
	chatResp.Result = strings.TrimPrefix(chatResp.Result, "```")
	chatResp.Result = strings.TrimSuffix(chatResp.Result, "```")
	chatResp.Result = strings.TrimSpace(chatResp.Result)

	// 2. 解析LLM提取的实体列表
	var entities []string
	err = json.Unmarshal([]byte(chatResp.Result), &entities)
	if err != nil {
		logs.Error("Failed to parse entities: %s, raw data: %s", err.Error(), chatResp.Result)
		return result, err
	}

	if len(entities) == 0 {
		logs.Info("No valid entities extracted")
		return result, nil
	}

	// 记录提取的实体数量
	logs.Info("Extracted %d entities from question[%s]", len(entities), question)

	// 3. 使用图数据库查询相关实体
	graphDB := NewGraphDB(cast.ToInt(robot[`admin_user_id`]))
	libraryIdsArr := strings.Split(libraryIds, ",")

	// 存储所有查询结果
	allResults := make([]*neo4j.Record, 0)
	dataIds := make(map[int]float64) // 用于去重和存储最高置信度

	// 为每个实体分配合理的查询限制
	perEntityLimit := size * 4
	if len(entities) > 0 {
		perEntityLimit = perEntityLimit / len(entities)
		if perEntityLimit < 10 {
			perEntityLimit = 10 // 确保最小查询限制
		}
	}

	// 4. 对每个实体执行优化后的深度查询（使用协程并行执行）
	queryStartTime := time.Now() // 查询开始时间
	var wg sync.WaitGroup
	var mu sync.Mutex // 用于保护 allResults 的并发访问

	// 为每个实体启动一个协程
	for i, entityName := range entities {
		if len(entityName) == 0 {
			logs.Error("Empty entity at index %d", i)
			continue
		}

		wg.Add(1)
		go func(index int, entity string) {
			defer wg.Done()
			logs.Info("Querying entity[%d]: %s", index+1, entity)
			res, err := graphDB.FindRelatedEntities(entity, libraryIdsArr, perEntityLimit, 3)
			if err != nil {
				logs.Error("Failed to query entity %s: %s", entity, err.Error())
				return
			}

			logs.Info("Found %d related triples for entity[%d]", len(res.Records), index+1)

			// 安全地追加结果
			mu.Lock()
			allResults = append(allResults, res.Records...)
			mu.Unlock()
		}(i, entityName)
	}

	// 等待所有查询完成
	wg.Wait()

	// 计算并记录查询时间
	queryDuration := time.Since(queryStartTime)
	logs.Info("Parallel entity query completed in %v for %d entities", queryDuration, len(entities))

	// 如果没有找到相关实体，直接返回
	if len(allResults) == 0 {
		logs.Info("No related entities found")
		return result, nil
	}

	// 5. 收集相关的数据ID和置信度，根据深度调整置信度
	for _, record := range allResults {
		dataIdValue, dataIdExists := record.Get("data_id")
		confidenceValue, confidenceExists := record.Get("confidence")
		depthValue, depthExists := record.Get("depth")
		if !dataIdExists || !confidenceExists || !depthExists {
			continue
		}
		dataId := cast.ToInt(dataIdValue)
		confidence := cast.ToFloat64(confidenceValue)
		depth := cast.ToInt(depthValue)
		if dataId <= 0 {
			continue
		}
		if confidence <= 0 {
			confidence = 0.5
		}
		if depth > 0 {
			// 每增加一级深度，降低20%的置信度
			depthFactor := 1.0 - float64(depth-1)*0.2
			if depthFactor < 0.4 {
				depthFactor = 0.4 // 最低保留40%的原始置信度
			}
			confidence = confidence * depthFactor
		}
		// 保存最高置信度
		if existingConf, exists := dataIds[dataId]; exists {
			if confidence > existingConf {
				dataIds[dataId] = confidence
			}
		} else {
			dataIds[dataId] = confidence
		}
	}

	// 6. 查询对应的段落数据
	if len(dataIds) > 0 {
		logs.Info("Found %d related data IDs", len(dataIds))

		dataIdList := make([]string, 0)
		for id := range dataIds {
			dataIdList = append(dataIdList, cast.ToString(id))
		}

		paragraphs, err := msql.Model("chat_ai_library_file_data", define.Postgres).
			Where(`delete_time`, `0`).
			Where("id", "in", strings.Join(dataIdList, ",")).
			Select()
		if err != nil {
			logs.Error("Failed to query paragraph data: %s", err.Error())
			return result, err
		}

		// 7. 添加置信度作为相似度得分
		for _, paragraph := range paragraphs {
			id := cast.ToInt(paragraph["id"])
			if conf, exists := dataIds[id]; exists {
				paragraph["similarity"] = cast.ToString(conf)
				result = append(result, paragraph)
			}
		}

		// 按相似度降序排序
		sort.Slice(result, func(i, j int) bool {
			return cast.ToFloat64(result[i]["similarity"]) > cast.ToFloat64(result[j]["similarity"])
		})

		// 限制返回大小
		if len(result) > size {
			result = result[:size]
		}

		logs.Info("Final query results: %d items", len(result))
	} else {
		logs.Info("No related data IDs found")
	}

	return result, nil
}

func GetMatchLibraryParagraphByFullTextSearch(question, libraryIds string, size, searchType int) ([]msql.Params, error) {
	list := make([]msql.Params, 0)
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeFullText}) {
		return list, nil
	}
	question = strings.ReplaceAll(question, `'`, ` `)
	question = strings.ReplaceAll(strings.ReplaceAll(question, "\r\n", ""), "\n", "")
	queryTokens, err := msql.Model(fmt.Sprintf(`ts_parse('zhparser', '%s')`, question), define.Postgres).ColumnArr(`token`)
	if err != nil {
		return nil, err
	}

	ids, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, `in`, libraryIds).
		Where(fmt.Sprintf(`to_tsvector('zhima_zh_parser',upper(content))@@to_tsquery('zhima_zh_parser',upper('%s'))`, strings.Join(queryTokens, " | "))).
		Limit(5000).ColumnArr(`id`)
	if err != nil {
		return list, err
	}
	if len(ids) == 0 {
		return list, nil
	}

	return msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Alias("a").
		Join("chat_ai_library_file_data b", "a.data_id=b.id", "left").
		Where(`a.delete_time`, `0`).
		Where(`a.id`, `in`, strings.Join(ids, `,`)).
		Where(`b.id is not null`).
		Field(`b.*,a.id as index_id`).
		Field(fmt.Sprintf(`ts_rank(to_tsvector('zhima_zh_parser',upper(a.content)),to_tsquery('zhima_zh_parser',upper('%s'))) as similarity`, question)).
		Order(`similarity DESC`).Limit(size).Select()
}

func GetMatchLibraryDataIdsByFullText(content, libraryIds string, size int) ([]string, error) {
	queryTokens := []string{content}
	ids, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, `in`, libraryIds).
		Where(`delete_time`, `0`).
		Where(fmt.Sprintf(`to_tsvector('zhima_zh_parser',upper(content))@@to_tsquery('zhima_zh_parser',upper('%s'))`, strings.Join(queryTokens, " | "))).
		Limit(size).ColumnArr(`id`)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	return ids, nil
}

func GetMatchLibraryDataIdsByLike(content, libraryIds string, size int) ([]string, error) {
	ids := make([]string, 0)
	// 精准搜索
	ids, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, `in`, libraryIds).
		Where(`delete_time`, `0`).
		Where(fmt.Sprintf(`content similar to '%%%s%%'`, content)).
		Limit(size).ColumnArr(`id`)
	if err != nil {
		return ids, err
	}
	if len(ids) <= 0 {
		return []string{`0`}, nil
	}
	return ids, nil
}

func GetMatchFileParagraphIdsByFullTextSearch(question, fileIds string) ([]string, error) {
	ids := make([]string, 0)
	// 精准搜索
	ids, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`file_id`, `in`, fileIds).
		Where(`delete_time`, `0`).
		Where(fmt.Sprintf(`content similar to '%%%s%%'`, question)).
		Limit(1000).ColumnArr(`id`)
	if err != nil {
		return ids, err
	}
	if len(ids) <= 0 {
		return []string{`0`}, nil
	}
	return ids, nil
}

func GetMatchLibraryParagraphByMergeRerank(openid, appType, question string, list []msql.Params, robot msql.Params) ([]msql.Params, error) {
	if len(robot) == 0 || cast.ToInt(robot[`rerank_status`]) == 0 {
		return nil, nil //not rerank config
	}
	if len(list) == 0 {
		return nil, nil
	}
	// Rerank resorted
	chunks := make([]string, 0)
	for _, one := range list {
		if cast.ToInt(one[`type`]) != define.ParagraphTypeNormal {
			var similarQuestions []string
			if err := tool.JsonDecode(one[`similar_questions`], &similarQuestions); err != nil {
				logs.Error(err.Error())
			}
			var similar string
			if len(similarQuestions) > 0 {
				similar = fmt.Sprintf("\n相似问法：%s", strings.Join(similarQuestions, `/`))
			}
			one[`content`] = fmt.Sprintf("问题:%s%s\n答案:%s", one[`question`], similar, one[`answer`])
		}
		chunks = append(chunks, one[`content`])
	}
	rerankReq := &adaptor.ZhimaRerankReq{
		Enable:   true,
		Query:    question,
		Passages: chunks,
		Data:     list,
		TopK:     min(500, len(list)),
	}
	return RerankData(cast.ToInt(robot[`admin_user_id`]), openid, appType, robot, rerankReq)
}

func GetMatchLibraryParagraphList(openid, appType, question string, optimizedQuestions []string, libraryIds string, size int, similarity float64, searchType int, robot msql.Params) (_ []msql.Params, libUseTime LibUseTime, _ error) {
	result := make([]msql.Params, 0)
	if len(libraryIds) == 0 {
		return result, libUseTime, nil
	}
	question = GetFirstQuestionByInput(question) //多模态输入特殊处理
	if len(question) == 0 {
		return result, libUseTime, nil
	}

	fetchSize := 4 * size
	var vectorList, searchList, graphList []msql.Params
	adminUserId := cast.ToInt(robot[`admin_user_id`])

	temp := time.Now()
	for _, q := range append(optimizedQuestions, question) {
		list, err := GetMatchLibraryParagraphByVectorSimilarity(adminUserId, robot, openid, appType, q, libraryIds, fetchSize, similarity, searchType)
		if err != nil {
			logs.Error(err.Error())
		}
		vectorList = append(vectorList, changeListContent(list)...)
		list, err = GetMatchLibraryParagraphByGraphSimilarity(robot, openid, appType, q, libraryIds, fetchSize, searchType)
		if err != nil {
			logs.Error(err.Error())
		}
		graphList = append(graphList, changeListContent(list)...)
		list, err = GetMatchLibraryParagraphByFullTextSearch(q, libraryIds, fetchSize, searchType)
		if err != nil {
			logs.Error(err.Error())
		}
		searchList = append(searchList, changeListContent(list)...)
	}
	libUseTime.RecallTime = time.Now().Sub(temp).Milliseconds()

	//由于存在问题优化,需要根据相似度再次排序
	sort.Slice(vectorList, func(i, j int) bool {
		return cast.ToFloat64(vectorList[i][`similarity`]) > cast.ToFloat64(vectorList[j][`similarity`])
	})
	sort.Slice(searchList, func(i, j int) bool {
		return cast.ToFloat64(searchList[i][`similarity`]) > cast.ToFloat64(searchList[j][`similarity`])
	})
	sort.Slice(graphList, func(i, j int) bool {
		return cast.ToFloat64(graphList[i][`similarity`]) > cast.ToFloat64(graphList[j][`similarity`])
	})

	//由于存在问题优化,导致召回的内容会出现重复的
	vectorList = SliceMsqlParamsUnique(vectorList, `id`)
	searchList = SliceMsqlParamsUnique(searchList, `id`)
	graphList = SliceMsqlParamsUnique(graphList, `id`)

	//RRF sort
	weights := ParseRrfWeight(adminUserId, robot[`rrf_weight`])
	list := (&RRF{}).
		Add(DataSource{List: vectorList, Key: `id`, Fixed: 50, Weight: weights.Vector}).
		Add(DataSource{List: searchList, Key: `id`, Fixed: 50, Weight: weights.Search}).
		Add(DataSource{List: graphList, Key: `id`, Fixed: 50, Weight: weights.Graph}).Sort()

	//由于存在问题优化和4倍召回,需要截取一下
	if len(list) > size {
		list = list[:size]
	}

	//rerank 重排逻辑
	temp = time.Now()
	rerankList, err := GetMatchLibraryParagraphByMergeRerank(openid, appType, question, list, robot)
	libUseTime.RerankTime = time.Now().Sub(temp).Milliseconds()
	if err != nil {
		logs.Error(err.Error())
	}
	if len(rerankList) > 0 {
		list = rerankList
	}

	//父子分段-替换成对应的父分段内容+去重
	list = FatherSonChunkReplace(list)

	//return
	var ids []string
	for i, one := range list {
		if i >= size {
			break
		}
		// Supplement file info
		fileInfo, _ := GetLibFileInfo(cast.ToInt(one[`file_id`]), 0)
		one[`file_name`] = fileInfo[`file_name`]
		ids = append(ids, one[`id`])
		result = append(result, one)
	}
	go UpdateParagraphHits(strings.Join(ids, `,`), 1)
	if len(result) > 0 {
		go statDailyRequestLibraryTip(adminUserId, robot, appType, cast.ToString(StatsTypeDailyLibraryTipCount))
		go StatLibraryTipUp(result, robot)
	}
	go statDailyRequestLibraryTip(adminUserId, robot, appType, cast.ToString(StatsTypeDailyAiMsgCount))
	return result, libUseTime, nil
}

func FatherSonChunkReplace(list []msql.Params) []msql.Params {
	existMap := make(map[string]struct{})
	result := make([]msql.Params, 0)
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	for _, one := range list {
		if cast.ToUint(one[`father_chunk_paragraph_number`]) > 0 { //父子分段逻辑,子段替换成父段内容
			if _, ok := existMap[one[`father_chunk_paragraph_number`]]; ok {
				continue //已存在对应的父分段内容,去重
			}
			contents, err := m.Where(`content`, `<>`, ``).Where(`file_id`, one[`file_id`]).
				Where(`father_chunk_paragraph_number`, one[`father_chunk_paragraph_number`]).
				Order(`page_num,father_chunk_paragraph_number,number`).ColumnArr(`content`)
			if err != nil {
				logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			}
			if len(contents) > 0 { //返回对应的父分段内容
				one[`content`] = strings.Join(contents, "\r\n")
				existMap[one[`father_chunk_paragraph_number`]] = struct{}{}
			}
		}
		result = append(result, one)
	}
	return result
}

func changeListContent(list []msql.Params) []msql.Params {
	for idx, one := range list {
		if one[`content`] == "" && one[`answer`] != "" {
			list[idx][`content`] = one[`answer`]
		}
	}
	return list
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
	return fmt.Sprintf(`chatwiki.model_config.v20251210.%d`, h.ModelConfigId)
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

type ModelListCacheBuildHandler struct{ ModelConfigId int }

func (h *ModelListCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.model_list.%d`, h.ModelConfigId)
}
func (h *ModelListCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_model_list`, define.Postgres).
		Where(`model_config_id`, cast.ToString(h.ModelConfigId)).Order(`id`).Select()
}

func GetModelListInfo(modelConfigId int) ([]msql.Params, error) {
	result := make([]msql.Params, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &ModelListCacheBuildHandler{ModelConfigId: modelConfigId}, &result, time.Hour*12)
	return result, err
}

func GetDefaultLlmConfig(adminUserId int) (int, string, bool) {
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		return 0, ``, false
	}
	sort.Slice(configs, func(i, j int) bool {
		return tool.InArrayString(configs[i][`model_define`], []string{ModelChatWiki}) //优先选取的模型
	})
	for _, config := range configs {
		if !tool.InArrayString(Llm, strings.Split(config[`model_types`], `,`)) {
			continue
		}
		if modelInfo, ok := GetModelInfoByConfig(adminUserId, cast.ToInt(config[`id`])); ok {
			if models := modelInfo.GetLlmModelList(); len(models) > 0 {
				return cast.ToInt(config[`id`]), models[0], true
			}
		}
	}
	return 0, ``, false
}

func GetDefaultEmbeddingConfig(adminUserId int) (int, string, bool) {
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		return 0, ``, false
	}
	sort.Slice(configs, func(i, j int) bool {
		return tool.InArrayString(configs[i][`model_define`], []string{ModelChatWiki}) //优先选取的模型
	})
	for _, config := range configs {
		if !tool.InArrayString(TextEmbedding, strings.Split(config[`model_types`], `,`)) {
			continue
		}
		if modelInfo, ok := GetModelInfoByConfig(adminUserId, cast.ToInt(config[`id`])); ok {
			if models := modelInfo.GetVectorModelList(); len(models) > 0 {
				return cast.ToInt(config[`id`]), models[0], true
			}
		}
	}
	return 0, ``, false
}

type WechatAppCacheBuildHandler struct {
	Field string
	Value string
}

func (h *WechatAppCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(lib_define.RedisPrefixAppInfo, h.Field, h.Value)
}
func (h *WechatAppCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(h.Field, h.Value).Find()
}

func GetWechatAppInfo(field, value string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &WechatAppCacheBuildHandler{Field: field, Value: value}, &result, time.Hour)
	return result, err
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
	if vectorType == cast.ToString(define.VectorTypeSimilarQuestion) {
		info = nil
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

type WeChatDialogueCacheBuildHandler struct {
	AdminUserId, RobotId int
	Openid               string
}

func (h *WeChatDialogueCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.wechat_dialogue.%d.%d.%s`, h.AdminUserId, h.RobotId, h.Openid)
}
func (h *WeChatDialogueCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_dialogue`, define.Postgres).Where(`admin_user_id`, cast.ToString(h.AdminUserId)).
		Where(`robot_id`, cast.ToString(h.RobotId)).Where(`openid`, h.Openid).Order(`id DESC`).Value(`id`)
}

func GetLastDialogueId(adminUserId, robotId int, openid string) int {
	var dialogueId any
	err := lib_redis.GetCacheWithBuild(define.Redis,
		&WeChatDialogueCacheBuildHandler{AdminUserId: adminUserId, RobotId: robotId, Openid: openid}, &dialogueId, time.Hour)
	if err != nil {
		logs.Error(err.Error())
	}
	return cast.ToInt(dialogueId)
}

func GetOptimizedQuestions(param *define.ChatRequestParam, contextList []map[string]string) ([]string, error) {
	question := GetFirstQuestionByInput(param.Question) //多模态输入特殊处理
	if len(question) == 0 {
		return []string{}, nil //输入没有文本内容时,跳过问题优化
	}
	histories := ""
	for _, item := range contextList {
		histories += "Q: " + item[`question`] + "\n"
		histories += "A: " + item[`answer`] + "\n"
	}

	var prompt string
	if len(param.Robot[`optimize_question_dialogue_background`]) > 0 {
		prompt = strings.ReplaceAll(define.PromptDefaultQuestionOptimize, `{{dialogue_background}}`, `# 对话背景：\n`+param.Robot[`optimize_question_dialogue_background`])
	} else {
		prompt = strings.ReplaceAll(define.PromptDefaultQuestionOptimize, `{{dialogue_background}}`, ``)
	}

	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `system`, Content: prompt}}
	for _, item := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: item[`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: item[`answer`]})
	}
	messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: question})

	modelConfigId := cast.ToInt(param.Robot[`model_config_id`])
	useModel := param.Robot[`use_model`]
	if cast.ToInt(param.Robot[`optimize_question_model_config_id`]) > 0 {
		modelConfigId = cast.ToInt(param.Robot[`optimize_question_model_config_id`])
		useModel = param.Robot[`optimize_question_use_model`]
	}

	chatResp, _, err := RequestChat(
		param.AdminUserId,
		param.Openid,
		param.Robot,
		param.AppType,
		modelConfigId,
		useModel,
		messages,
		nil,
		cast.ToFloat32(param.Robot[`temperature`]),
		200,
	)
	if err != nil {
		return nil, err
	}

	return []string{chatResp.Result}, nil
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

func LibraryAiSummary(lang, question, prompt string, optimizedQuestions []string, libraryIds string, size, maxToken int, similarity, temperature float64, searchType int, robot msql.Params, chanStream chan sse.Event, summarySwitch int) error {
	defer close(chanStream)
	chanStream <- sse.Event{Event: `ping`, Data: tool.Time2Int()}
	if summarySwitch == define.SwitchOff {
		chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
		return nil
	}
	list, _, err := GetMatchLibraryParagraphList("", "", question, []string{}, libraryIds, size, similarity, searchType, robot)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(lang, `sys_err`)}
		return err
	}
	if len(list) == 0 {
		chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}
		return nil
	}
	quoteFile := []msql.Params{}
	quoteFileMap := make(map[int]bool)
	summary := []adaptor.ZhimaChatCompletionMessage{
		{
			Role:    `user`,
			Content: question,
		},
		{
			Role:    `system`,
			Content: prompt,
		},
	}

	prompt, _ = FormatSystemPrompt("", list)
	for _, item := range list {
		file, err := GetLibFileInfo(cast.ToInt(item[`file_id`]), cast.ToInt(item[`admin_user_id`]))
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if _, ok := quoteFileMap[cast.ToInt(file[`id`])]; ok {
			continue
		}
		quoteFile = append(quoteFile, msql.Params{
			`id`:        file[`id`],
			`file_name`: file[`file_name`],
		})
	}

	summary = append(summary, adaptor.ZhimaChatCompletionMessage{
		Role:    `assistant`,
		Content: prompt,
	})
	// to ai summary
	chatResp, requestTime, err := RequestSearchStream(
		cast.ToInt(robot[`admin_user_id`]),
		cast.ToInt(robot[`model_config_id`]),
		strings.TrimSpace(robot[`use_model`]),
		robot,
		summary,
		[]adaptor.FunctionTool{},
		chanStream,
		float32(temperature),
		maxToken,
	)
	logs.Info(`LibraryAiSummary:%v`, chatResp)
	logs.Info(`LibraryAiSummary:%v`, requestTime)
	if err != nil {
		logs.Error(err.Error())
		chanStream <- sse.Event{Event: `error`, Data: i18n.Show(lang, `sys_err`)}
	}
	//message push
	if len(quoteFile) > 0 {
		chanStream <- sse.Event{Event: `quote_file`, Data: quoteFile}
	}
	chanStream <- sse.Event{Event: `finish`, Data: tool.Time2Int()}

	return nil
}

type LibFileSplitAiChunksBacheHandle struct{ TaskId string }

func (h *LibFileSplitAiChunksBacheHandle) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.lib.file.split.ai.chunks.%s`, h.TaskId)
}
func (h *LibFileSplitAiChunksBacheHandle) GetCacheData() (any, error) {
	return nil, nil
}

type LibFileSplitAiChunksCache struct {
	ErrMsg string               `json:"err_msg"`
	List   define.DocSplitItems `json:"list"`
}

func (h *LibFileSplitAiChunksBacheHandle) SaveCacheData(list LibFileSplitAiChunksCache) error {
	_, err := define.Redis.Set(context.Background(), h.GetCacheKey(), tool.JsonEncodeNoError(list), 1*time.Hour).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return err
	}
	return nil
}

func GetNeo4jStatus(adminUserId int) bool {
	if status, exists := define.Neo4jStatus[adminUserId]; exists {
		return status
	} else {
		return false
	}
}

func SetNeo4jStatus(adminUserId int, status bool) {
	if define.Neo4jStatus == nil {
		define.Neo4jStatus = make(map[int]bool)
	}
	define.Neo4jStatus[adminUserId] = status
}

func ConstructGraph(id int) error {
	dataList, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`file_id`, cast.ToString(id)).
		Where(`delete_time`, `0`).
		Field(`id,file_id`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(dataList) == 0 {
		return err
	}
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{`graph_status`: define.GraphStatusInitial})
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: id})
	for _, data := range dataList {
		message, err := tool.JsonEncode(map[string]any{`id`: data[`id`], `file_id`: data[`file_id`]})
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if err = AddJobs(define.ConvertGraphTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}
	return nil
}

type FaqFilesCacheBuildHandler struct{ FileId int }

func (h *FaqFilesCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki_faq_files_info.%d`, h.FileId)
}
func (h *FaqFilesCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`chat_ai_faq_files`, define.Postgres).Where(`id`, cast.ToString(h.FileId)).Find()
}

func GetFaqFilesInfo(fileId, adminUserId int) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &FaqFilesCacheBuildHandler{FileId: fileId}, &result, time.Hour)
	if err == nil && adminUserId != 0 && cast.ToInt(result[`admin_user_id`]) != adminUserId {
		result = make(msql.Params) //attribution error. null data returned
	}
	return result, err
}

func DeleteLibraryFileDataIndex(dataIdList, vectorType string) error {
	_, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).
		Where(`data_id`, `in`, dataIdList).
		Where(`type`, `in`, vectorType).
		Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}

type TokenAppLimitConfigCacheBuildHandler struct {
	AdminUserId  int
	TokenAppType string
	RobotId      int
}

func (h *TokenAppLimitConfigCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`token.app.limit.config.%d.%s.%d`, h.AdminUserId, h.TokenAppType, h.RobotId)
}
func (h *TokenAppLimitConfigCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(h.AdminUserId)).
		Where(`token_app_type`, h.TokenAppType).Where(`robot_id`, cast.ToString(h.RobotId)).Find()
}

type TokenAppUseCacheBuildHandler struct {
	AdminUserId  int
	TokenAppType string
	RobotId      int
	DateYmd      string
}

func (h *TokenAppUseCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`token.app.use.incr.%s.%d.%s.%d`, h.DateYmd, h.AdminUserId, h.TokenAppType, h.RobotId)
}
func (h *TokenAppUseCacheBuildHandler) GetCacheData() (any, error) {
	if h.DateYmd == `` {
		return GetTokenAppLimitUse(h.AdminUserId, h.RobotId, h.TokenAppType)
	}
	return GetTokenAppUseByDate(h.AdminUserId, h.RobotId, h.TokenAppType, h.DateYmd)
}

// 获取配置
func GetAdminConfig(UserId int) map[string]string {
	m := msql.Model(``, define.Postgres)
	list, _ := m.Table(`admin_user_config`).Where(`admin_user_id`, cast.ToString(UserId)).Select()

	if len(list) == 0 { //如果没有配置，给一下默认配置
		defaultList := make(map[string]string)
		defaultList["draft_exptime"] = "10"
		defaultList["comment_pull_days"] = "7"
		defaultList["comment_pull_limit"] = "10"
		return defaultList
	}

	return list[0]
}
