// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
)

var CheckFileLearnedMutex sync.Map
var CheckFileGraphLearnedMutex sync.Map

func ConvertHtml(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	fileId, link := cast.ToInt(data[`file_id`]), cast.ToString(data[`file_url`])
	if fileId <= 0 || len(link) == 0 {
		logs.Error(`data exception:%s`, msg)
		return nil
	}
	if !common.LinkExists(link) {
		logs.Error(`file does not exist:%s`, msg)
		return nil
	}
	info, err := common.GetLibFileInfo(fileId, 0)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if len(info) == 0 {
		logs.Error(`no data:%s`, msg)
		return nil
	}
	if cast.ToInt(info[`status`]) != define.FileStatusInitial {
		logs.Error(`abnormal state:%s/%v`, msg, info[`status`])
		return nil
	}
	//convert html
	htmlUrl, err := common.ConvertHtml(link, cast.ToInt(info[`admin_user_id`]))
	if err != nil && err.Error() == `Service Unavailable` {
		logs.Error(`service unavailable. try again in one minute:%s`, msg)
		_ = common.AddJobs(define.ConvertHtmlTopic, msg, time.Minute)
		return nil
	}
	m := msql.Model(`chat_ai_library_file`, define.Postgres)
	if err != nil {
		_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
			`status`:      define.FileStatusException,
			`errmsg`:      err.Error(),
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
		}
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
		return nil
	}
	_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`html_url`:    htmlUrl,
		`status`:      define.FileStatusWaitSplit,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})

	//create default lib file split
	splitParams := common.DefaultSplitParams()
	common.AutoSplitLibFile(cast.ToInt(info[`admin_user_id`]), fileId, splitParams)

	return nil
}

func ConvertVector(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	id, fileId := cast.ToInt(data[`id`]), cast.ToInt(data[`file_id`])
	if id <= 0 || fileId <= 0 {
		logs.Error(`data exception:%s`, msg)
		return nil
	}
	file, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	info, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if len(info) == 0 {
		logs.Error(`no data:%s`, msg)
		return nil
	}
	//start convert
	library, _ := common.GetLibraryInfo(cast.ToInt(info[`library_id`]), cast.ToInt(info[`admin_user_id`]))
	skipUseModel := cast.ToInt(library[`type`]) == define.OpenLibraryType && cast.ToInt(library[`use_model_switch`]) != define.SwitchOn
	if skipUseModel {
		return nil
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`status`:      define.VectorStatusConverting,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	embedding, err := common.GetVector2000(
		cast.ToInt(info[`admin_user_id`]),
		info[`admin_user_id`],
		msql.Params{},
		info,
		file,
		cast.ToInt(library[`model_config_id`]),
		library[`use_model`],
		info[`content`],
	)
	if err != nil {
		_, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
			`status`:      define.VectorStatusException,
			`errmsg`:      err.Error(),
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			return nil
		}

		//check finish
		CheckFileLearned(fileId, define.FileStatusPartException)
		return nil
	}

	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`status`:      define.VectorStatusConverted,
		`embedding`:   embedding,
		`errmsg`:      `success`,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	//check finish
	CheckFileLearned(fileId, define.FileStatusLearned)
	return nil
}

// ConvertGraph Convert Knowledge Graph
func ConvertGraph(msg string, _ ...string) error {
	logs.Debug(`nsq:graph:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	id, fileId := cast.ToInt(data[`id`]), cast.ToInt(data[`file_id`])
	if id <= 0 || fileId <= 0 {
		logs.Error(`data exception:%s`, msg)
		return nil
	}
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	// get file data
	info, err := m.Where(`id`, cast.ToString(id)).Where(`file_id`, cast.ToString(fileId)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if len(info) == 0 {
		logs.Error(`no data:%s`, msg)
		return nil
	}

	library, err := common.GetLibraryInfo(cast.ToInt(info[`library_id`]), cast.ToInt(info[`admin_user_id`]))
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if !cast.ToBool(library[`graph_switch`]) {
		logs.Error("graph not switch")
		return nil
	}

	constructGraphInit(fileId, id)

	// construct graph
	content := info[`content`]
	if cast.ToInt(info[`type`]) != define.ParagraphTypeNormal {
		content = "question: " + info[`question`] + "\n\nanswer: " + info[`answer`]
	}
	prompt := strings.ReplaceAll(define.PromptDefaultGraphConstruct, `{{content}}`, content)
	messages := []adaptor.ZhimaChatCompletionMessage{{Role: `user`, Content: prompt}}
	chatResp, _, err := common.RequestChat(
		cast.ToInt(info[`admin_user_id`]),
		info[`admin_user_id`],
		msql.Params{},
		"",
		cast.ToInt(library[`graph_model_config_id`]),
		library[`graph_use_model`],
		messages,
		nil,
		0.1,
		3000,
	)
	if err != nil {
		logs.Error(err.Error())
		constructGraphFailed(fileId, id, err.Error())
		CheckFileGraphLearned(fileId)
		return nil
	}
	chatResp.Result = strings.TrimPrefix(chatResp.Result, "```json")
	chatResp.Result = strings.TrimSuffix(chatResp.Result, "```")

	var graphData []map[string]interface{}
	if err := tool.JsonDecode(chatResp.Result, &graphData); err != nil {
		logs.Error(`graph data parsing failure:%s/%s`, chatResp.Result, err.Error())
		constructGraphFailed(fileId, id, err.Error())
		CheckFileGraphLearned(fileId)
		return nil
	}

	hasError := false
	for _, triple := range graphData {
		subject := cast.ToString(triple["subject"])
		predicate := cast.ToString(triple["predicate"])
		object := cast.ToString(triple["object"])
		confidence := cast.ToFloat64(triple["confidence"])

		if len(subject) > 0 && len(predicate) > 0 && len(object) > 0 {
			graphDB := common.NewGraphDB("graphrag")
			// 替换关系名称中的特殊字符
			sanitizedPredicate := strings.ReplaceAll(predicate, ".", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, "...", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, " ", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, "-", "_")

			createGraphSQL := fmt.Sprintf(`
				cypher('graphrag', $$ 
					MERGE (s:Entity {name: '%s', library_id: %d, file_id: %d, data_id: %d})
					MERGE (o:Entity {name: '%s', library_id: %d, file_id: %d, data_id: %d})
					CREATE (s)-[r:%s {confidence: %f, library_id: %d, file_id: %d, data_id: %d}]->(o)
					RETURN r
				$$) as (r agtype)`,
				subject, cast.ToInt(info[`library_id`]), fileId, id,
				object, cast.ToInt(info[`library_id`]), fileId, id,
				sanitizedPredicate, confidence, cast.ToInt(info[`library_id`]), fileId, id)
			_, err = graphDB.ExecuteCypher(createGraphSQL)
			if err != nil {
				logs.Error(`create graph error: %s, sql is: `, err.Error(), createGraphSQL)
				hasError = true
				constructGraphFailed(fileId, id, err.Error())
				break
			}
		}
	}
	if !hasError {
		constructGraphSucceed(id)
	}

	CheckFileGraphLearned(fileId)
	return nil
}

func constructGraphInit(fileId, dataId int) {
	_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, cast.ToString(dataId)).
		Update(msql.Datas{
			`graph_status`: define.GraphStatusWorking,
			`update_time`:  tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`id`, cast.ToString(fileId)).
		Where(`graph_status`, cast.ToString(define.GraphStatusInitial)).
		Update(msql.Datas{
			`graph_status`: define.GraphStatusWorking,
			`update_time`:  tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
}

func constructGraphFailed(fileId, dataId int, errMsg string) {
	_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, cast.ToString(dataId)).
		Update(msql.Datas{
			`graph_status`:  define.GraphStatusException,
			`graph_err_msg`: errMsg,
			`update_time`:   tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
}

func constructGraphSucceed(dataId int) {
	_, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, cast.ToString(dataId)).
		Update(msql.Datas{
			`graph_status`: define.GraphStatusConverted,
			`update_time`:  tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
}

func CheckFileLearned(fileId, status int) {
	mtx, _ := CheckFileLearnedMutex.LoadOrStore(fileId, &sync.Mutex{})
	mutex := mtx.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	m := msql.Model(`chat_ai_library_file_data_index`, define.Postgres)

	total, err := m.Where(`file_id`, cast.ToString(fileId)).Where(`status`, cast.ToString(define.VectorStatusInitial)).Count(`1`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if total > 0 {
		return //not finish
	}

	// finished
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`status`:      status,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
}

func CheckFileGraphLearned(fileId int) {
	mtx, _ := CheckFileGraphLearnedMutex.LoadOrStore(fileId, &sync.Mutex{})
	mutex := mtx.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)

	res, err := m.Where(`file_id`, cast.ToString(fileId)).
		Field(fmt.Sprintf(`count(case when graph_status = %d then 1 else null end) as notstart_count`, define.GraphStatusNotStart)).
		Field(fmt.Sprintf(`count(case when graph_status = %d then 1 else null end) as working_count`, define.GraphStatusWorking)).
		Field(fmt.Sprintf(`count(case when graph_status = %d then 1 else null end) as converted_count`, define.GraphStatusConverted)).
		Field(fmt.Sprintf(`count(case when graph_status = %d then 1 else null end) as exception_count`, define.GraphStatusException)).
		Field(fmt.Sprintf(`count(1) as total_count`)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	if cast.ToInt(res[`notstart_count`]) > 0 || cast.ToInt(res[`working_count`]) > 0 {
		return // not finish
	}

	finalStatus := define.GraphStatusConverted
	if cast.ToInt(res[`exception_count`]) > 0 {
		finalStatus = define.GraphStatusPartlyConverted
		if cast.ToInt(res[`exception_count`]) == cast.ToInt(res[`total_count`]) {
			finalStatus = define.GraphStatusException
		}
	} else {
		finalStatus = define.GraphStatusConverted
	}
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`graph_status`: finalStatus,
		`update_time`:  tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return
	}

	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
}

func CrawlArticle(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	fileId := cast.ToInt(cast.ToInt(data[`file_id`]))
	adminUserId := cast.ToInt(cast.ToInt(data[`admin_user_id`]))
	if fileId <= 0 || adminUserId <= 0 {
		logs.Error(`data exception:%s`, msg)
		return nil
	}

	// check file id
	m := msql.Model(`chat_ai_library_file`, define.Postgres)
	file, err := common.GetLibFileInfo(fileId, cast.ToInt(data[`admin_user_id`]))
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if len(file) == 0 {
		logs.Error(`library not found:%s`, msg)
		return nil
	}

	// update file status
	_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`status`: define.FileStatusCrawling,
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	//start crawl
	uploadInfo, err := common.SaveUrlPage(cast.ToInt(file[`admin_user_id`]), file[`doc_url`], "library_file")
	if err != nil {
		_, err := m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
			`status`: define.FileStatusCrawlException,
			`errmsg`: err.Error(),
		})
		if err != nil {
			logs.Error(err.Error())
		}
		return nil
	}

	// update file status
	_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`status`:              define.FileStatusInitial,
		`file_name`:           uploadInfo.Name,
		`file_size`:           uploadInfo.Size,
		`file_url`:            uploadInfo.Link,
		`update_time`:         tool.Time2Int(),
		`doc_last_renew_time`: tool.Time2Int(),
	})
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	// convert html
	if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId, `file_url`: uploadInfo.Link}); err != nil {
		logs.Error(err.Error())
	} else if err := common.AddJobs(define.ConvertHtmlTopic, message); err != nil {
		logs.Error(err.Error())
	}
	return nil
}

func ExportTask(id string, _ ...string) error {
	logs.Debug(`nsq:%s`, id)
	m := msql.Model(`chat_ai_export_task`, define.Postgres)
	task, err := m.Where(`id`, id).Where(`status`, cast.ToString(define.ExportStatusWaiting)).Field(`source,params`).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if len(task) == 0 {
		return nil //任务状态不对
	}
	_, err = m.Where(`id`, id).Where(`status`, cast.ToString(define.ExportStatusWaiting)).Update(msql.Datas{
		`status`:      define.ExportStatusRunning,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	var status = define.ExportStatusError
	var fileUrl string
	var errMsg = `unknown`
	defer func() {
		_, err := m.Where(`id`, id).Where(`status`, cast.ToString(define.ExportStatusRunning)).Update(msql.Datas{
			`status`:      status,
			`file_url`:    fileUrl,
			`err_msg`:     errMsg,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	params := make(map[string]any)
	if err = tool.JsonDecodeUseNumber(task[`params`], &params); err != nil {
		errMsg = `参数解析错误:` + err.Error()
		return nil
	}
	switch cast.ToUint(task[`source`]) {
	case define.ExportSourceSession: //会话记录导出
		fileUrl, err = common.RunSessionExport(params)
	default:
		errMsg = `来源类型错误:` + task[`source`]
	}
	if err != nil { //导出失败
		errMsg = err.Error()
	} else { //导出成功
		status = define.ExportStatusSucceed
		errMsg = `SUCCEED`
	}
	return nil
}
