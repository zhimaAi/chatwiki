// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var CheckFileLearnedMutex sync.Map

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
	if !tool.IsFile(common.GetFileByLink(link)) {
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
	lang := define.LangEnUs

	splitParams := define.SplitParams{}
	splitParams.ChunkSize = 512
	splitParams.ChunkOverlap = 0
	splitParams.SeparatorsNo = `11,12`
	splitParams.EnableExtractImage = true
	list, wordTotal, err := common.GetLibFileSplit(cast.ToInt(info[`admin_user_id`]), fileId, splitParams, lang)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = common.SaveLibFileSplit(cast.ToInt(info[`admin_user_id`]), fileId, wordTotal, define.QAIndexTypeQuestionAndAnswer, splitParams, list, lang)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

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
	if cast.ToInt(info[`status`]) != define.VectorStatusInitial {
		logs.Error(`abnormal state:%s/%v`, msg, info[`status`])
		return nil
	}
	//start convert
	library, _ := common.GetLibraryInfo(cast.ToInt(info[`library_id`]), cast.ToInt(info[`admin_user_id`]))
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
		CheckFileLearned(fileId)
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
	CheckFileLearned(fileId)
	return nil
}

func CheckFileLearned(fileId int) {
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
		`status`:      define.FileStatusLearned,
		`update_time`: tool.Time2Int(),
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
