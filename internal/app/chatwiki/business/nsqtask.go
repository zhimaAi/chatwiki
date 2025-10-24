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
	m := msql.Model(`chat_ai_library_file`, define.Postgres)

	//convert html
	htmlUrl, err := common.ConvertHtml(fileId, link, cast.ToInt(info[`admin_user_id`]), cast.ToInt(info[`pdf_parse_type`]))
	if err != nil && err.Error() == `Service Unavailable` {
		logs.Error(`service unavailable. try again in one minute:%s`, msg)
		_ = common.AddJobs(define.ConvertHtmlTopic, msg, time.Minute)
		return nil
	}
	if err != nil && strings.Contains(err.Error(), `pdf parse cancelled`) {
		logs.Error(err.Error())
		_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
			`status`:      define.FileStatusCancelled,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
		}
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
		return nil
	}
	if err != nil {
		logs.Error(err.Error())
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
	if len(info[`async_split_params`]) > 0 {
		if err = tool.JsonDecodeUseNumber(info[`async_split_params`], &splitParams); err != nil {
			logs.Error(err.Error())
		}
	}
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
	if id <= 0 || fileId < 0 {
		logs.Error(`data exception:%s`, msg)
		return nil
	}
	var file msql.Params
	var err error
	if fileId > 0 {
		if file, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Find(); err != nil {
			logs.Error(err.Error())
		}
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
		logs.Debug("graph not switch")
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
		4096,
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
			// 替换关系名称中的特殊字符
			sanitizedPredicate := strings.ReplaceAll(predicate, ".", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, "...", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, " ", "_")
			sanitizedPredicate = strings.ReplaceAll(sanitizedPredicate, "-", "_")

			_, err = common.NewGraphDB(cast.ToInt(info[`admin_user_id`])).ConstructEntity(
				subject,
				object,
				sanitizedPredicate,
				cast.ToInt(info[`library_id`]),
				fileId,
				id,
				confidence,
			)
			if err != nil {
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
		logs.Error(err.Error())
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
		`async_split_params`:  ``, //清空之前的参数
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
	var (
		fileUrl  string
		fileName string
	)
	var errMsg = `unknown`
	defer func() {
		datas := msql.Datas{
			`status`:      status,
			`file_url`:    fileUrl,
			`err_msg`:     errMsg,
			`update_time`: tool.Time2Int(),
		}
		if len(fileName) > 0 {
			datas[`file_name`] = fileName
		}
		_, err := m.Where(`id`, id).Where(`status`, cast.ToString(define.ExportStatusRunning)).Update(datas)
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
	case define.ExportSourceLibFileDoc: //知识库文档导出
		fileUrl, fileName, err = common.RunLibFileDocExport(params)
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

func ExtractFaqFiles(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	dataIds, fileId := cast.ToString(data[`ids`]), cast.ToInt(data[`file_id`])
	var file msql.Params
	var err error
	if file, err = common.GetFaqFilesInfo(fileId, 0); err != nil {
		logs.Error(err.Error())
	}
	if len(file) == 0 {
		logs.Error(`file not found:%s`, msg)
		return nil
	}
	lib_redis.DelCacheData(define.Redis, &common.FaqFilesCacheBuildHandler{FileId: fileId})

	adminUserId := cast.ToInt(file[`admin_user_id`])
	splitFaqParams := define.SplitFaqParams{
		ChunkType:          cast.ToInt(file[`chunk_type`]),
		ChunkSize:          cast.ToInt(file[`chunk_size`]),
		SeparatorsNo:       file[`separators_no`],
		ChunkPrompt:        file[`chunk_prompt`],
		ChunkModel:         file[`chunk_model`],
		ChunkModelConfigId: cast.ToInt(file[`chunk_model_config_id`]),
		ExtractType:        define.FAQExtractTypeAI,
	}
	// FAQFileStatusAnalyzing
	status := define.FAQFileStatusAnalyzing
	errMsg := ""
	if err = common.UpdateLibFileFaqStatus(fileId, adminUserId, status, errMsg); err != nil {
		logs.Error(err.Error())
		return nil
	}
	var list []define.DocSplitItem
	m := msql.Model(`chat_ai_faq_files_data`, define.Postgres)
	if len(dataIds) > 0 {
		result, err := m.Where(`id`, `in`, dataIds).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Select()
		if err != nil {
			logs.Error(err.Error())
		}
		for _, item := range result {
			var images []string
			if err = tool.JsonDecode(cast.ToString(item[`images`]), &images); err != nil {
				logs.Error(err.Error())
			}
			list = append(list, define.DocSplitItem{
				Number:  cast.ToInt(item[`number`]),
				PageNum: cast.ToInt(item[`page_num`]),
				Title:   cast.ToString(item[`title`]),
				Content: cast.ToString(item[`content`]),
				Images:  images,
			})
		}
	} else {
		list, err = common.GetLibFileFaqSplit(fileId, adminUserId, splitFaqParams)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	// FAQ chunks
	splitList, err := common.MultSplitFaqFiles(list, splitFaqParams)
	status = define.FAQFileStatusExtracting
	if err = common.UpdateLibFileFaqStatus(fileId, adminUserId, status, errMsg); err != nil {
		logs.Error(err.Error())
		return nil
	}
	//FAQ Extracting
	newList := make(chan define.DocSplitItem, 10)
	go common.ExtractLibFaqFiles(adminUserId, splitFaqParams, splitList, newList)
	if file, err = common.GetFaqFilesInfo(fileId, 0); err != nil {
		logs.Error(err.Error())
	}
	if len(file) == 0 {
		logs.Error(`file not found:%s`, msg)
		return nil
	}
	isNew := true
	qaModel := msql.Model(`chat_ai_faq_files_data_qa`, define.Postgres)
	if len(dataIds) > 0 {
		isNew = false
		_, err = m.Where(`id`, `in`, dataIds).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Delete()
		_, err = qaModel.Where(`data_id`, `in`, dataIds).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Delete()
	} else {
		_, err = m.Where(`file_id`, `in`, cast.ToString(fileId)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Delete()
		_, err = qaModel.Where(`file_id`, `in`, cast.ToString(fileId)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Delete()
	}
	if err != nil {
		logs.Error(err.Error())
		errMsg = err.Error()
		common.UpdateLibFileFaqStatus(fileId, adminUserId, define.FAQFileStatusExtractFailed, errMsg)
		return nil
	}
	// save faq
	status = define.FAQFileStatusExtracted
	type chunkResult struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	failCount := 0
	total := 0
	for {
		select {
		case item, ok := <-newList:
			if !ok {
				goto EndExtract
			}
			total++
			data := msql.Datas{
				`admin_user_id`: file[`admin_user_id`],
				`file_id`:       fileId,
				`number`:        item.Number,
				`split_status`:  define.FAQFileSplitStatusSuccess,
				`page_num`:      item.PageNum,
				`title`:         item.Title,
				`content`:       item.Content,
				`images`:        tool.JsonEncodeNoError(item.Images),
				`word_total`:    item.WordTotal,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			}
			allContents := []chunkResult{}
			if item.AiChunkErrMsg == "" {
				if err := tool.JsonDecode(item.Answer, &allContents); err != nil {
					data[`split_errmsg`] = fmt.Sprintf(`json解析失败:%s`, err.Error())
					data[`split_status`] = define.FAQFileSplitStatusFailed
					data[`err_content`] = item.Answer
				}
			} else {
				data[`split_errmsg`] = item.AiChunkErrMsg
				data[`split_status`] = define.FAQFileSplitStatusFailed
			}
			if data[`split_status`] == define.FAQFileSplitStatusFailed {
				failCount++
			}
			dataId, err := m.Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				errMsg = err.Error()
				continue
			}
			var (
				images = item.Images
				imgs   = make([]string, 0)
			)
			for k, chunk := range allContents {
				if len(chunk.Question+chunk.Answer) == 0 {
					continue
				}
				var answer string
				answer, imgs, images = common.InTextImagesPlaceholders(chunk.Answer, images)
				_, err := qaModel.Insert(msql.Datas{
					`admin_user_id`: adminUserId,
					`file_id`:       fileId,
					`data_id`:       dataId,
					`number`:        k + 1,
					`page_num`:      item.PageNum,
					`question`:      chunk.Question,
					`answer`:        answer,
					`images`:        tool.JsonEncodeNoError(imgs),
					`create_time`:   tool.Time2Int(),
					`update_time`:   tool.Time2Int(),
				})
				if err != nil {
					logs.Error(err.Error())
					errMsg = err.Error()
					continue
				}
			}
		}
	}

EndExtract:
	if failCount == total && isNew {
		errMsg = `提取文件失败`
		status = define.FAQFileStatusExtractFailed
	}
	common.UpdateLibFileFaqStatus(fileId, adminUserId, status, errMsg)
	return nil
}
