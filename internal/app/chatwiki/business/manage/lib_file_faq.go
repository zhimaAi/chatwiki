// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetFAQConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 从数据库获取FAQ配置信息
	config, err := msql.Model("chat_ai_faq_files", define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Order("id DESC").
		Field(`id,chunk_type,chunk_prompt,chunk_model,chunk_model_config_id,chunk_size,chunk_overlap,separators_no`).
		Find()
	if err != nil {
		common.FmtError(c, "sys_err")
		return
	}
	common.FmtOk(c, config)
}

func AddFAQFile(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//common save
	fileIds, err := addFAQFile(c, adminUserId)
	if err != nil {
		if err.Error() == `vip_limits` {
			common.FmtError(c, `vip_limits`)
			return
		} else if err.Error() == `chunk_prompt_too_long` {
			common.FmtError(c, `chunk_prompt_too_long`)
			return
		}
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{`file_ids`: fileIds})
}

func RenewFAQFileData(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	fileId := cast.ToInt(c.PostForm(`id`))
	if fileId == 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}

	// 查询文件信息
	file, err := common.GetFaqFilesInfo(fileId, adminUserId)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(file) == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	// 查询失败的数据ID列表
	var faildIds []string
	failedData, err := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Field("id,split_status").Select()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(failedData) > 0 {
		for _, item := range failedData {
			if cast.ToInt(item["split_status"]) == define.FAQFileSplitStatusFailed {
				faildIds = append(faildIds, item[`id`])
			}
		}
	}
	// 重新提取失败的数据
	if message, err := tool.JsonEncode(map[string]any{
		"file_id": fileId,
		"ids":     strings.Join(faildIds, ","),
	}); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	} else if err := common.AddJobs(define.ExtractFaqFilesTopic, message); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.UpdateLibFileFaqStatus(fileId, adminUserId, define.FAQFileStatusQueuing, "")
	common.FmtOk(c, nil)
}

// GetFAQFileList 获取FAQ文件列表
func GetFAQFileList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 获取分页参数
	page := cast.ToInt(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize := cast.ToInt(c.Query("size"))
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询FAQ文件列表
	m := msql.Model("chat_ai_faq_files", define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId))
	// 分页查询
	list, total, err := m.Order("id DESC").
		Paginate(page, pageSize)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}

	// 获取每个文件的统计数据
	// 获取所有文件ID
	fileIds := make([]string, 0, len(list))
	for i := range list {
		fileIds = append(fileIds, cast.ToString(list[i]["id"]))
	}

	// 批量查询成功数
	successCountMap := make(map[string]int)
	successList, _ := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", "in", strings.Join(fileIds, ",")).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("split_status", cast.ToString(define.FAQFileSplitStatusSuccess)).
		Group("file_id").
		Field("file_id, COUNT(*) as count").Select()
	for _, item := range successList {
		successCountMap[cast.ToString(item["file_id"])] = cast.ToInt(item["count"])
	}

	// 批量查询失败数
	failCountMap := make(map[string]int)
	failList, _ := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", "in", strings.Join(fileIds, ",")).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("split_status", cast.ToString(define.FAQFileSplitStatusFailed)).
		Group("file_id").
		Field("file_id, COUNT(*) as count").Select()
	for _, item := range failList {
		failCountMap[cast.ToString(item["file_id"])] = cast.ToInt(item["count"])
	}

	// 批量查询QA总数
	qaCountMap := make(map[string]int)
	qaList, _ := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("file_id", "in", strings.Join(fileIds, ",")).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Group("file_id").
		Field("file_id, COUNT(*) as count").Select()
	for _, item := range qaList {
		qaCountMap[cast.ToString(item["file_id"])] = cast.ToInt(item["count"])
	}

	// 赋值统计数据
	for i := range list {
		fileId := cast.ToString(list[i]["id"])
		list[i]["success_count"] = cast.ToString(successCountMap[fileId])
		list[i]["fail_count"] = cast.ToString(failCountMap[fileId])
		list[i]["qa_count"] = cast.ToString(qaCountMap[fileId])
	}
	common.FmtOk(c, map[string]any{
		"list":  list,
		"total": total,
		"page":  page,
	})
}

func GetFAQFileChunks(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 获取文件ID
	fileId := cast.ToInt(c.Query("id"))
	splitStatus := c.DefaultQuery("split_status", cast.ToString(define.FAQFileSplitStatusFailed))
	if fileId == 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	// 获取分页参数
	data, err := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("split_status", cast.ToString(splitStatus)).Select()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	for i := range data {
		data[i]["content"] = strings.ReplaceAll(cast.ToString(data[i]["content"]), "[图片占位符]", "")
	}
	common.FmtOk(c, data)
}

// DeleteFAQFile 删除FAQ文件
func DeleteFAQFile(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	fileId := cast.ToInt(c.PostForm("id"))
	if fileId == 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}

	m := msql.Model("chat_ai_faq_files", define.Postgres)
	m.Begin()
	value, err := m.Where("id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Value(`id`)
	if len(value) == 0 {
		m.Rollback()
		common.FmtError(c, `no_data`)
		return
	}
	// 删除主表数据
	_, err = m.Where("id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Delete()
	if err != nil {
		m.Rollback()
		common.FmtError(c, `sys_err`)
		return
	}

	// 删除分析数据
	_, err = msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Delete()
	if err != nil {
		m.Rollback()
		common.FmtError(c, `sys_err`)
		return
	}

	// 删除QA数据
	_, err = msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Delete()
	if err != nil {
		m.Rollback()
		common.FmtError(c, `sys_err`)
		return
	}

	lib_redis.DelCacheData(define.Redis, &common.FaqFilesCacheBuildHandler{FileId: fileId})
	m.Commit()
	common.FmtOk(c, nil)
}

// GetFAQFileInfo 获取FAQ文件状态
func GetFAQFileInfo(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	fileId := cast.ToInt(c.Query("id"))
	if fileId == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	// 查询文件状态
	file, err := msql.Model("chat_ai_faq_files", define.Postgres).
		Where("id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Find()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(file) == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	// 获取统计数据
	successCount, _ := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("split_status", cast.ToString(define.FAQFileSplitStatusSuccess)).
		Count()

	failCount, _ := msql.Model("chat_ai_faq_files_data", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("split_status", cast.ToString(define.FAQFileSplitStatusFailed)).
		Count()

	qaCount, _ := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Count()

	file["success_count"] = cast.ToString(successCount)
	file["fail_count"] = cast.ToString(failCount)
	file["qa_count"] = cast.ToString(qaCount)
	common.FmtOk(c, file)
}

func GetFAQFileQAList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 获取文件ID
	fileId := cast.ToInt(c.Query("id"))
	if fileId == 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	isImport := c.Query("is_import")
	// 获取分页参数
	page := cast.ToInt(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize := cast.ToInt(c.Query("size"))
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询QA列表
	m := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId))
	if isImport != "" {
		m.Where("is_import", isImport)
	}
	// 分页查询
	list, _, err := m.Order("id ASC").
		Paginate(page, pageSize)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	for i := range list {
		libraryInfo, _ := common.GetLibraryInfo(cast.ToInt(list[i][`library_id`]), adminUserId)
		list[i][`library_name`] = libraryInfo[`library_name`]
	}
	// 查询总数
	totalData, _ := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("file_id", cast.ToString(fileId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Group(`is_import`).
		Field(`is_import,count(id) as count`).Select()
	total := 0
	importTotal := 0
	for _, item := range totalData {
		total += cast.ToInt(item[`count`])
		if cast.ToInt(item[`is_import`]) == define.SwitchOn {
			importTotal = cast.ToInt(item[`count`])
		}
	}
	common.FmtOk(c, map[string]any{
		"list":         list,
		"total":        total,
		`import_total`: importTotal,
		"page":         page,
	})
}

func ExportFAQFileAllQA(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 获取文件ID
	fileId := cast.ToInt(c.Query("id"))
	ext := c.DefaultQuery("ext", `xlsx`)
	if fileId == 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	if !tool.InArrayString(ext, []string{`xlsx`, `docx`}) {
		common.FmtError(c, `param_invalid`, `ext`)
		return
	}
	fileInfo, err := common.GetFaqFilesInfo(fileId, adminUserId)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(fileInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	page := 1
	pageSize := 500
	data := make([]msql.Params, 0)
	// 查询QA列表
	for {
		m := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
			Where("file_id", cast.ToString(fileId)).
			Where("admin_user_id", cast.ToString(adminUserId)).Field("question,answer,images")
		// 分页查询
		list, _, err := m.Order("id ASC").
			Paginate(page, pageSize)
		if err != nil {
			common.FmtError(c, `sys_err`)
			return
		}
		if len(list) == 0 {
			break
		}
		page++
		data = append(data, list...)
	}
	filePath, err := common.ExportFAQFileAllQA(data, ext, `lib_faq`)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	c.FileAttachment(filePath, fmt.Sprintf(`FAQ分段结果导出%s.%s`, tool.Date(`Y-m-d-H-i-s`), ext))
	go func(filePath string) {
		time.Sleep(1 * time.Minute)
		_ = os.Remove(filePath)
	}(filePath)
}

func SaveFAQFileQA(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// 获取参数
	id := cast.ToInt(c.PostForm("id"))
	if id == 0 {
		common.FmtError(c, "param_invalid", "id")
		return
	}

	question := strings.TrimSpace(c.PostForm("question"))
	answer := strings.TrimSpace(c.PostForm("answer"))
	images := c.PostFormArray(`images`)
	if question == "" || answer == "" {
		common.FmtError(c, "param_invalid", "question or answer")
		return
	}

	// 查询数据是否存在
	m := msql.Model("chat_ai_faq_files_data_qa", define.Postgres)
	data, err := m.Where("id", cast.ToString(id)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Value(`id`)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(data) == 0 {
		common.FmtError(c, "no_data")
		return
	}

	updateData := msql.Datas{
		"question":    question,
		"answer":      answer,
		"update_time": tool.Time2Int(),
	}
	// 提取答案中的图片
	jsonImages, err := common.CheckLibraryImage(images)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if len(jsonImages) > 0 && len(images) > 0 {
		updateData[`images`] = jsonImages
	}
	// 更新数据
	_, err = m.Where("id", cast.ToString(id)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Update(updateData)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

func DeleteFAQFileQA(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	ids := c.PostForm("ids")
	if len(ids) == 0 {
		common.FmtError(c, `param_invalid`, `ids`)
		return
	}

	// 批量删除数据
	_, err := msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("id", "in", ids).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Delete()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, nil)
}

func ImportParagraph(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	// 获取文件ID
	libraryId := cast.ToInt(c.PostForm("library_id"))
	fileId := cast.ToInt(c.PostForm("file_id"))
	ids := cast.ToString(c.PostForm("ids"))
	// 根据id或file_id查询QA数据并保存到段落
	if libraryId == 0 || fileId == 0 {
		common.FmtError(c, "param_invalid", "library_id or file_id")
		return
	}
	common.ImportFAQFile(adminUserId, libraryId, fileId, ids, c.GetHeader(`token`), true)
	common.FmtOk(c, nil)
}

func addFAQFile(c *gin.Context, adminUserId int) ([]int64, error) {
	m := msql.Model(`chat_ai_faq_files`, define.Postgres)
	//get params
	chunkType := cast.ToInt(c.PostForm(`chunk_type`))
	chunkSize := cast.ToInt(c.PostForm(`chunk_size`))
	chunkPrompt := strings.TrimSpace(c.PostForm(`chunk_prompt`))
	chunkModel := strings.TrimSpace(c.PostForm(`chunk_model`))
	chunkModelConfigId := cast.ToInt(c.PostForm(`chunk_model_config_id`))
	separatorsNo := strings.TrimSpace(c.PostForm(`separators_no`))
	//document uploaded
	var (
		libraryFiles []*define.UploadInfo
	)
	if utf8.RuneCountInString(chunkPrompt) > 5000 {
		return nil, errors.New(`chunk_prompt_too_long`)
	}
	// uploaded file
	libFileAloowExts := define.FAQLibFileAllowExt
	libraryFiles, _ = common.SaveUploadedFileMulti(c, `faq_files`, define.LibFileLimitSize, adminUserId, `faq_files`, libFileAloowExts)
	if len(libraryFiles) == 0 {
		return nil, errors.New(i18n.Show(common.GetLang(c), `upload_empty`))
	}
	fileSize := 0
	for _, uploadInfo := range libraryFiles {
		fileSize += int(uploadInfo.Size)
	}
	// dispatch queue
	fileIds := make([]int64, 0)
	for _, uploadInfo := range libraryFiles {
		status := define.FAQFileStatusQueuing
		insData := msql.Datas{
			`admin_user_id`:         adminUserId,
			`creator`:               getLoginUserId(c),
			`file_url`:              uploadInfo.Link,
			`file_name`:             uploadInfo.Name,
			`file_ext`:              uploadInfo.Ext,
			`file_size`:             uploadInfo.Size,
			`status`:                status,
			`chunk_size`:            chunkSize,
			`separators_no`:         separatorsNo,
			`chunk_type`:            chunkType,
			`chunk_prompt`:          chunkPrompt,
			`chunk_model`:           chunkModel,
			`chunk_model_config_id`: chunkModelConfigId,
			`create_time`:           tool.Time2Int(),
			`update_time`:           tool.Time2Int(),
		}
		fileId, err := m.Insert(insData, `id`)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		//dispatch queue
		if message, err := tool.JsonEncode(map[string]any{`file_id`: fileId}); err != nil {
			logs.Error(err.Error())
		} else if err := common.AddJobs(define.ExtractFaqFilesTopic, message); err != nil {
			logs.Error(err.Error())
		}
		fileIds = append(fileIds, fileId)
	}
	return fileIds, nil
}
