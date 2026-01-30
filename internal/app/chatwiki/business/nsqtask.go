// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/thumbnail"
	"chatwiki/internal/pkg/wechat"
	"context"
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/broadcasting/request"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"

	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
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
		define.LangEnUs,
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

	var embeddingKey = `embedding`
	if common.GetVectorDims(embedding) == 2000 {
		embeddingKey = `embedding2000` //固定2000维度向量的文档
	}

	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`status`:      define.VectorStatusConverted,
		embeddingKey:  embedding,
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
		define.LangEnUs,
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
	logs.Debug(fmt.Sprintf(`nsq:%s`, msg))
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

	returnDir := filepath.Join(fmt.Sprintf(`%s/upload/chat_ai/%d/%s/%s`, filepath.Dir(define.AppRoot), adminUserId, `library_file`, tool.Date(`Ym`)))
	thumbFolderPath := filepath.Join(fmt.Sprintf(`/upload/chat_ai/%d/%s/%s/`, adminUserId, `library_file`, tool.Date(`Ym`)))

	//生成缩略图
	thumbPath, _, _, genErr := thumbnail.GenerateThumbnail(common.GetFileByLink(uploadInfo.Link), returnDir, thumbFolderPath)

	if genErr != nil {
		logs.Error("generate thumbnail error：" + genErr.Error())
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
		`thumb_path`:          thumbPath,
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

func CrawlFeishuDoc(msg string, _ ...string) error {
	logs.Debug(fmt.Sprintf(`nsq:%s`, msg))
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

	// 获取飞书文档详情
	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.FeiShuRobot, `app_id`: file[`feishu_app_id`], `app_secret`: file[`feishu_app_secret`]})
	feishuApp, ok := app.(wechat.FeishuInterface)
	if !ok {
		logs.Error(`feishu app init failed`)
		return nil
	}
	title, content, err := feishuApp.GetDocFileDetail(file[`feishu_document_id`])
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	md5Hash := tool.MD5(content)
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, adminUserId, `library_file`, tool.Date(`Ym`), md5Hash, `md`)
	link, err := common.WriteFileByString(objectKey, content)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	returnDir := filepath.Join(fmt.Sprintf(`%s/upload/chat_ai/%d/%s/%s`, filepath.Dir(define.AppRoot), adminUserId, `library_file`, tool.Date(`Ym`)))
	thumbFolderPath := filepath.Join(fmt.Sprintf(`/upload/chat_ai/%d/%s/%s/`, adminUserId, `library_file`, tool.Date(`Ym`)))

	//生成缩略图
	thumbPath, _, _, genErr := thumbnail.GenerateThumbnail(common.GetFileByLink(link), returnDir, thumbFolderPath)

	if genErr != nil {
		logs.Error("generate thumbnail error：" + genErr.Error())
	}

	//if common.IsUrl(link) {
	//	err = os.Rename(outputPath, filepath.Dir(define.AppRoot)+thumbPath)
	//	if err != nil {
	//		logs.Error("move file error：" + err.Error())
	//	}
	//}

	// update file status
	_, err = m.Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`status`:              define.FileStatusWaitSplit,
		`file_name`:           title,
		`file_url`:            link,
		`thumb_path`:          thumbPath,
		`file_ext`:            `md`,
		`file_size`:           len(content),
		`update_time`:         tool.Time2Int(),
		`doc_last_renew_time`: tool.Time2Int(),
	})
	lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: fileId})
	if err != nil {
		logs.Error(err.Error())
		return nil
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
		return nil //task status incorrect
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
		errMsg = `parameter parsing error:` + err.Error()
		return nil
	}
	lang := define.LangEnUs //default is english
	switch cast.ToUint(task[`source`]) {
	case define.ExportSourceSession: //会话记录导出
		fileUrl, err = common.RunSessionExport(lang, params)
	case define.ExportSourceLibFileDoc: //知识库文档导出
		fileUrl, fileName, err = common.RunLibFileDocExport(lang, params)
	default:
		errMsg = `invalid source type:` + task[`source`]
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
	var list define.DocSplitItems
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
					data[`split_errmsg`] = fmt.Sprintf(`json parsing failed:%s`, err.Error())
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
		errMsg = `file extraction failed`
		status = define.FAQFileStatusExtractFailed
	}
	common.UpdateLibFileFaqStatus(fileId, adminUserId, status, errMsg)
	return nil
}

func OfficialAccountDraftSync(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)

	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return err
	}

	//验证admin_user_id 和 app_id 获取secret
	appInfo, err := common.GetWechatAppInfo(`access_key`, cast.ToString(data["access_key"]))
	if err != nil {
		logs.Error(`failed to get wechat app info:` + err.Error())
	}

	limit := cast.ToInt(data["limit"])
	admin_user_id := cast.ToInt(data["admin_user_id"])

	//初始化client
	client, err := common.GetOfficialAccountApp(appInfo["app_id"], appInfo["app_secret"])
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	count := 20 //每轮拉取条数
	round := math.Ceil(cast.ToFloat64(limit) / cast.ToFloat64(count))

	if limit == 0 { //获取全部的草稿，最多拉取100轮，2千个草稿
		round = 100
	}

	m := msql.Model(`wechat_official_account_draft`, define.Postgres)
	respData := define.OfficialAccountDraftListStruct{}

	for i := 0; i < cast.ToInt(round); i++ {
		offset := i * count
		params := &object.HashMap{
			"offset":     offset,
			"count":      count,
			"no_content": 1, //不要草稿内容明细
		}

		_, err = client.Base.BaseClient.HttpPostJson(context.Background(), "/cgi-bin/draft/batchget", params, nil, nil, &respData)

		//no more data, exit
		if len(respData.Item) == 0 {
			break
		}

		//遍历数据结果
		for _, s := range respData.Item {
			updateData := msql.Datas{
				`admin_user_id`:     admin_user_id,
				`app_id`:            appInfo["app_id"],
				`update_time`:       time.Now().Unix(),
				`article_type`:      s.Content.NewsItem[0].ArticleType,
				`thumb_url`:         s.Content.NewsItem[0].ThumbUrl,
				`title`:             s.Content.NewsItem[0].Title,
				`digest`:            s.Content.NewsItem[0].Digest,
				`draft_create_time`: s.Content.CreateTime,
				`draft_update_time`: s.Content.UpdateTime,
			}
			res, e := m.Where("media_id", s.MediaId).Update(updateData)

			if e != nil {
				logs.Error(`execution error:` + e.Error())
			}

			//if update count is 0, insert data
			if res == 0 {
				updateData[`create_time`] = time.Now().Unix()
				updateData[`media_id`] = s.MediaId
				m.Insert(updateData)
			}
		}
	}

	return nil
}

func OfficialAccountCommentSync(msg string, _ ...string) error {

	taskInfo := define.DelayTaskEvent{}
	if err := tool.JsonDecodeUseNumber(msg, &taskInfo); err != nil {
		logs.Error(`task:%s,err:%s`, msg, err.Error())
		return err
	}

	//验证任务类型
	if taskInfo.Type != define.OfficialAccountBatchSendSyncCommentTask {
		logs.Error(`invalid task type` + cast.ToString(taskInfo.Type))
		return nil
	}

	//查询任务详情信息
	taskData, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Alias("a").
		Join("wechat_official_account_draft b", "a.draft_id = b.id ", "inner").
		Join("chat_ai_wechat_app c", "a.access_key = c.access_key and a.admin_user_id = c.admin_user_id ", "inner").
		Where(`a.id`, cast.ToString(taskInfo.TaskId)).Where(`a.admin_user_id`, cast.ToString(taskInfo.AdminUserId)).
		Field(`a.*,b.media_id,c.app_id,c.app_secret`).Find()

	if err != nil {
		logs.Error(`failed to query task info:` + err.Error())
		return err
	}

	//	获取评论列表
	client, err := common.GetOfficialAccountApp(taskData["app_id"], taskData["app_secret"])
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	beginIndex := 0
	pageSize := 50

	maxCommentId := cast.ToInt(taskData["max_comment_id"])
	newMaxCommentId := cast.ToInt(taskData["max_comment_id"])
	returnSync := false

	for true {
		params := &object.HashMap{
			"msg_data_id": taskData["msg_data_id"],
			"index":       0,
			"begin":       beginIndex,
			"count":       pageSize,
			"type":        0,
		}
		respData := define.OfficialAccountCommentResp{}
		_, err = client.Base.BaseClient.HttpPostJson(context.Background(), "/cgi-bin/comment/list", params, nil, nil, &respData)
		if len(respData.Comment) == 0 {
			break
		}

		beginIndex += pageSize

		for _, comment := range respData.Comment {
			//if current comment id is less than recorded comment id, it has been synced, skip
			if comment.UserCommentId <= maxCommentId {
				returnSync = true
				break
			}

			//update max comment id
			if newMaxCommentId < comment.UserCommentId {
				newMaxCommentId = comment.UserCommentId
			}

			updateData := msql.Datas{
				"admin_user_id":       taskData["admin_user_id"],
				"update_time":         time.Now().Unix(),
				"msg_data_id":         taskData["msg_data_id"],
				"access_key":          taskData["access_key"],
				"task_id":             taskData["id"],
				"draft_id":            taskData["draft_id"],
				"user_comment_id":     comment.UserCommentId,
				"comment_create_time": comment.CreateTime,
				"content_text":        comment.Content,
				"comment_type":        comment.CommentType,
				"open_id":             comment.Openid,
				"reply_comment_text":  comment.Reply.Content,
				"reply_create_time":   comment.Reply.CreateTime,
			}

			//if no text content, default to AI selection
			if updateData["content_text"] == "" {
				updateData["ai_comment_rule_status"] = 1
			}

			res, err := msql.Model(`wechat_official_comment_list`, define.Postgres).
				Where(`admin_user_id`, taskData["admin_user_id"]).
				Where(`msg_data_id`, taskData["msg_data_id"]).
				Where(`user_comment_id`, cast.ToString(comment.UserCommentId)).Update(updateData)
			if err != nil {
				logs.Error(`update data failed:` + err.Error())
				continue
			}

			if res == 0 {
				updateData["comment_rule_id"] = taskData["comment_rule_id"]
				updateData["create_time"] = time.Now().Unix()

				_, err := msql.Model(`wechat_official_comment_list`, define.Postgres).Insert(updateData, "id")
				if err != nil {
					logs.Error(`data write failed:` + err.Error())
					continue
				}
			}
		}

		if returnSync {
			break
		}
	}

	//评论同步完了，更新一下数据
	msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(taskInfo.TaskId)).Where(`admin_user_id`, cast.ToString(taskInfo.AdminUserId)).Update(msql.Datas{
		`update_time`:            time.Now().Unix(),
		`last_comment_sync_time`: time.Now().Unix(),
		`max_comment_id`:         newMaxCommentId,
	})

	//if auto selection not enabled, exit
	if cast.ToInt(taskData["ai_comment_status"]) == 0 || !common.CheckUseAbilityByAbilityType(taskInfo.AdminUserId, common.OfficialAccountAbilityAIComment) {
		logs.Debug(`rule or account does not enable AI comment selection`)
		return nil
	}

	commentList, err := msql.Model(`wechat_official_comment_list`, define.Postgres).
		Where(`admin_user_id`, taskData["admin_user_id"]).
		Where(`msg_data_id`, taskData["msg_data_id"]).Where(`ai_comment_rule_status`, "0").Select()

	if err != nil {
		logs.Error(`failed to get comment list for AI selection:` + err.Error())
		return err
	}

	for _, params := range commentList {
		taskInfo := lib_define.OfficialAccountCommentAiCheckReq{
			AdminUserId:   params["admin_user_id"],
			CommentId:     cast.ToInt(params["id"]),
			TaskId:        cast.ToInt(params["task_id"]),
			MsgDataId:     params["msg_data_id"],
			AccessKey:     params["access_key"],
			CommentRuleId: cast.ToInt(params["comment_rule_id"]),
			UserCommentId: cast.ToInt(params["user_comment_id"]),
			CommentText:   params["content_text"],
		}
		if err := common.AddJobs(define.OfficialAccountCommentAiCheckTopic, tool.JsonEncodeNoError(taskInfo)); err != nil {
			logs.Error(`nsq production error, fallback to sync logic:%s`, err.Error())
			_ = OfficialAccountCommentAiCheck(tool.JsonEncodeNoError(taskInfo))
		}
	}
	return nil
}

func OfficialAccountCommentAiCheck(msg string, _ ...string) error {

	taskInfo := lib_define.OfficialAccountCommentAiCheckReq{}

	if err := tool.JsonDecodeUseNumber(msg, &taskInfo); err != nil {
		logs.Error(`task:%s,err:%s`, msg, err.Error())
		return err
	}

	query := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where(`admin_user_id`, cast.ToString(taskInfo.AdminUserId))

	if taskInfo.CommentRuleId == 0 { //默认规则
		query.Where("is_default", "1")
	} else {
		query.Where(`id`, cast.ToString(taskInfo.CommentRuleId))
	}

	//查询评论规则详情
	ruleInfo, err := query.Find()
	if err != nil {
		logs.Error(`failed to get comment rule:%s`, err.Error())
		return err
	}

	//自定义规则，需要开启
	if cast.ToInt(ruleInfo["is_default"]) == 0 && ruleInfo["switch"] != define.BaseOpen {
		logs.Error(`comment rule not enabled`)
		return nil
	}

	ai_comment_rule_text := []string{} //触发条件
	ai_comment_result := map[int]int{} //处理结果
	reply_context := ""

	//开始执行判断逻辑
	lang := define.LangZhCn
	if ruleInfo["delete_comment_switch"] == define.BaseOpen { //自动删除评论判断
		nodeRule := manage.BridgeDeleteCommentRule{}
		_ = tool.JsonDecode(ruleInfo["delete_comment_rule"], &nodeRule)
		hasKeywords := false
		aiDelete := false

		for _, checkType := range nodeRule.Type {
			//敏感词检测，且敏感词数大于1
			if checkType == define.CommentCheckTypeHintKeywords && len(nodeRule.Keywords) > 0 { //敏感词检测
				for _, keyword := range nodeRule.Keywords {
					if strings.Contains(taskInfo.CommentText, keyword) { //包含关键词
						hasKeywords = true
						ai_comment_rule_text = append(ai_comment_rule_text, i18n.Show(lang, `ai_comment_delete_sensitive_word`, keyword))
						break
					}
				}
			}

			//AI检测，且存在提示词
			if checkType == define.CommentCheckTypeAICheck && nodeRule.Prompt != "" {
				userPrompt := i18n.Show(lang, `ai_comment_delete_check_prompt`, taskInfo.CommentText, nodeRule.Prompt)
				checkRes := getAiCommentCheckRes(taskInfo.AdminUserId, userPrompt, ruleInfo["use_model"], ruleInfo["model_config_id"])

				if checkRes.NeedDelete { //如果需要删除
					ai_comment_rule_text = append(ai_comment_rule_text, i18n.Show(lang, `ai_comment_delete_ai_check`))
					aiDelete = checkRes.NeedDelete
				}
			}

		}

		//只选了一个判断条件，触发敏感词或者AI检测通过之后，就需要删除
		if len(nodeRule.Type) == 1 && (hasKeywords || aiDelete) {
			ai_comment_result[define.CommentExecTypeDelete] = define.CommentExecTypeDelete
		}

		//选了两个条件，需要同时满足，且两个条件都通过
		if len(nodeRule.Type) == 2 && nodeRule.Condition == 1 && hasKeywords && aiDelete {
			ai_comment_result[define.CommentExecTypeDelete] = define.CommentExecTypeDelete
		}

		//选了两个条件，满足任意一个，且有一个通过
		if len(nodeRule.Type) == 2 && nodeRule.Condition == 2 && (hasKeywords || aiDelete) {
			ai_comment_result[define.CommentExecTypeDelete] = define.CommentExecTypeDelete
		}
	}

	//如果没有删除评论操作，继续判断
	if _, ok := ai_comment_result[define.CommentExecTypeDelete]; !ok && ruleInfo["reply_comment_switch"] == define.BaseOpen { //自动回复判断

		nodeRule := manage.BridgeReplyCommentRule{}
		_ = tool.JsonDecode(ruleInfo["reply_comment_rule"], &nodeRule)

		userPrompt := i18n.Show(lang, `ai_comment_reply_check_prompt`, taskInfo.CommentText, nodeRule.CheckReplyPrompt)
		checkRes := getAiCommentCheckRes(taskInfo.AdminUserId, userPrompt, ruleInfo["use_model"], ruleInfo["model_config_id"])

		//如果需要回复，且使用固定回复内容
		if checkRes.NeedReply && nodeRule.ReplyType == 1 {
			reply_context = nodeRule.ReplyPrompt
			ai_comment_result[define.CommentExecTypeReply] = define.CommentExecTypeReply
			ai_comment_rule_text = append(ai_comment_rule_text, i18n.Show(lang, `ai_comment_default_reply`))
		}

		//如果需要回复，使用AI回复
		if checkRes.NeedReply && nodeRule.ReplyType == 2 {
			userPrompt := i18n.Show(lang, `ai_comment_reply_generation_prompt`, taskInfo.CommentText, nodeRule.ReplyPrompt)
			checkRes = getAiCommentCheckRes(taskInfo.AdminUserId, userPrompt, ruleInfo["use_model"], ruleInfo["model_config_id"])

			reply_context = checkRes.ReplyContent
			ai_comment_rule_text = append(ai_comment_rule_text, i18n.Show(lang, `ai_comment_auto_reply`))
			ai_comment_result[define.CommentExecTypeReply] = define.CommentExecTypeReply
		}

	}

	//如果没有删除评论操作，继续判断
	if _, ok := ai_comment_result[define.CommentExecTypeDelete]; !ok && ruleInfo["elect_comment_switch"] == define.BaseOpen { //自动精选判断
		nodeRule := manage.BridgeElectCommentRule{}
		_ = tool.JsonDecode(ruleInfo["elect_comment_rule"], &nodeRule)

		userPrompt := i18n.Show(lang, `ai_comment_top_check_prompt`, taskInfo.CommentText, nodeRule.Prompt)
		checkRes := getAiCommentCheckRes(taskInfo.AdminUserId, userPrompt, ruleInfo["use_model"], ruleInfo["model_config_id"])

		if checkRes.NeedTop {
			ai_comment_result[define.CommentExecTypeTop] = define.CommentExecTypeTop
			ai_comment_rule_text = append(ai_comment_rule_text, i18n.Show(lang, `ai_comment_auto_top`))
		}
	}

	commentResult := []int{}

	updateData := msql.Datas{
		"update_time":            time.Now().Unix(),
		"ai_exec_time":           time.Now().Unix(),
		"ai_comment_rule_status": 1,
		"comment_rule_id":        ruleInfo["id"],
		"ai_comment_rule_text":   strings.Join(ai_comment_rule_text, ","),
	}
	for key, _ := range ai_comment_result {
		commentResult = append(commentResult, key)
		if key == define.CommentExecTypeDelete { //删除评论
			_, err := manage.BridgeDeleteComment(taskInfo.AccessKey, taskInfo.MsgDataId, 0, taskInfo.UserCommentId)
			if err == nil {
				updateData["delete_status"] = 1
			}
		}

		if key == define.CommentExecTypeReply && reply_context != "" { //回复评论
			_, err := manage.BridgeReplyComment(taskInfo.AccessKey, taskInfo.MsgDataId, reply_context, 0, taskInfo.UserCommentId)
			if err == nil {
				updateData["reply_comment_text"] = reply_context
				updateData["reply_create_time"] = time.Now().Unix()
			}
		}

		if key == define.CommentExecTypeTop { //置顶评论
			_, err := manage.BridgeMarkElect(taskInfo.AccessKey, taskInfo.MsgDataId, 0, taskInfo.UserCommentId)
			if err == nil {
				updateData["comment_type"] = 1
			}
		}
	}

	updateData["ai_comment_result"] = tool.JsonEncodeNoError(commentResult)

	_, err = msql.Model(`wechat_official_comment_list`, define.Postgres).
		Where("admin_user_id", taskInfo.AdminUserId).
		Where("id", cast.ToString(taskInfo.CommentId)).
		Where("user_comment_id", cast.ToString(taskInfo.UserCommentId)).Update(updateData)

	if err != nil {
		logs.Error(`data update failed:` + err.Error())
	}

	return nil
}

func getAiCommentCheckRes(AdminUserId, userPrompt, use_model, model_config_id string) lib_define.OfficialAccountCommentAiCheckRes {

	checkRes := lib_define.OfficialAccountCommentAiCheckRes{}
	messages := []adaptor.ZhimaChatCompletionMessage{
		{Role: `system`, Content: define.OfficialAccountCommentCheckPrompt},
		{Role: `user`, Content: userPrompt},
	}

	chatResp, _, err := common.RequestChat(define.LangZhCn, cast.ToInt(AdminUserId), AdminUserId, nil, lib_define.AppYunPc,
		cast.ToInt(model_config_id), use_model, messages, nil, 0.5, 2000)
	if err != nil {
		logs.Error(`AI detection failed:` + err.Error())
		return checkRes
	}

	err = tool.JsonDecodeUseNumber(chatResp.Result, &checkRes)
	if err != nil {
		logs.Error(`AI detection result parsing failed:` + err.Error())
	}

	return checkRes
}

func OfficialAccountBatchSend(msg string, _ ...string) error {
	logs.Debug(`consumer received message:%s`, msg)

	taskInfo := define.DelayTaskEvent{}
	if err := tool.JsonDecodeUseNumber(msg, &taskInfo); err != nil {
		logs.Error(`task:%s,err:%s`, msg, err.Error())
		return err
	}

	//验证任务类型
	if taskInfo.Type != define.OfficialAccountBatchSendDelayTask {
		logs.Error(`invalid task type` + cast.ToString(taskInfo.Type))
		return nil
	}

	if !common.CheckUseAbilityByAbilityType(taskInfo.AdminUserId, common.OfficialAccountAbilityBatchSend) {
		logs.Debug(`account batch send not enabled:` + cast.ToString(taskInfo.AdminUserId))
		return nil
	}

	//此处执行发送任务动作，需要任务状态是开启的
	taskData, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Alias("a").
		Join("wechat_official_account_draft b", "a.draft_id = b.id ", "inner").
		Join("chat_ai_wechat_app c", "a.access_key = c.access_key and a.admin_user_id = c.admin_user_id ", "inner").
		Where(`a.id`, cast.ToString(taskInfo.TaskId)).Where(`a.admin_user_id`, cast.ToString(taskInfo.AdminUserId)).
		Where(`a.open_status`, define.BaseOpen).
		Field(`a.*,b.media_id,c.app_id,c.app_secret`).Find()

	if err != nil {
		logs.Error(`failed to query task info:` + err.Error())
		return err
	}

	if cast.ToInt(taskData["send_status"]) != 0 {
		logs.Error(`invalid task status:` + cast.ToString(taskData["send_status"]))
		return err
	}

	//初始化client
	client, err := common.GetOfficialAccountApp(taskData["app_id"], taskData["app_secret"])
	if err != nil {
		logs.Error(`wechat client init failed:` + err.Error())
		return err
	}

	//创建发送任务
	resp := power.HashMap{}
	res, err := client.Broadcasting.SendNews(context.Background(), taskData["media_id"], &request.Reception{
		//ToUser: []string{"o0DNE6i3WNEKnFRu7XXEzs9wLMtQ", "o0DNE6kjqc121GkNXSUng71jo5Kc"},
		//Filter: nil,
		ToUser: []string{},
		Filter: &request.Filter{
			IsToAll: true,
		},
	}, &resp)

	if err != nil {
		logs.Error(err.Error())
	}

	officialAccountBatchSendRes := lib_define.OfficialAccountBatchSendRes{}
	_ = tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(res), &officialAccountBatchSendRes)

	if officialAccountBatchSendRes.Errcode != 0 {
		logs.Error(`batch send task execution failed:` + officialAccountBatchSendRes.Errmsg)

		//变更任务状态
		_, _ = msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(taskInfo.TaskId)).Where(`admin_user_id`, cast.ToString(taskInfo.AdminUserId)).Update(msql.Datas{
			`update_time`: time.Now().Unix(),
			`send_res`:    tool.JsonEncodeNoError(res),
			`send_status`: define.BatchSendStatusErr, //send failed
		})
		return nil
	}

	//变更任务状态
	_, _ = msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(taskInfo.TaskId)).Where(`admin_user_id`, cast.ToString(taskInfo.AdminUserId)).Update(msql.Datas{
		`update_time`: time.Now().Unix(),
		`send_status`: define.BatchSendStatusExec,
		`msg_data_id`: officialAccountBatchSendRes.MsgDataId,
		`send_msg_id`: officialAccountBatchSendRes.MsgId,
		`send_res`:    tool.JsonEncodeNoError(res),
	})

	_, err = manage.BridgeChangeCommentStatus(cast.ToString(taskData["access_key"]), cast.ToString(officialAccountBatchSendRes.MsgDataId), 0, cast.ToInt(taskData["comment_status"]))
	if err != nil {
		return err
	}

	return nil
}

func ImportLibFileFaq(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	adminUserId := cast.ToInt(data[`admin_user_id`])
	libraryId := cast.ToInt(data[`library_id`])
	fileId := cast.ToInt(data[`file_id`])
	ids := cast.ToString(data[`ids`])
	token := cast.ToString(data[`token`])
	common.ImportFAQFile(adminUserId, libraryId, fileId, ids, token, false)
	return nil
}
