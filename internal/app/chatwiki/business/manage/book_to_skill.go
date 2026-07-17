// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

// CreateBookToSkillTask creates a new book-to-skill generation task
// Each uploaded file creates a separate task, executed asynchronously via NSQ
func CreateBookToSkillTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	// Get params
	robotId := cast.ToInt(c.PostForm("robot_id"))
	skillName := c.PostForm("skill_name")
	modelConfigId := cast.ToInt(c.PostForm("model_config_id"))
	useModel := c.PostForm("use_model")
	templateContent := c.PostForm("template_content")
	lang := c.PostForm("lang")

	// Validate required
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "param_lack"))))
		return
	}
	if skillName = filterValidUtf8(skillName, define.BookToSkillMaxFileName); skillName == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "param_invalid", "skill_name"))))
		return
	}

	// Validate robot
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, int64(robotId))
	if !ok {
		return
	}

	// Validate model
	if !common.CheckModelIsValid(adminUserId, modelConfigId, useModel, common.Llm) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "param_invalid", "model_config_id"))))
		return
	}

	// Get uploaded files
	form := c.Request.MultipartForm
	if form == nil || len(form.File["files"]) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "param_lack"))))
		return
	}
	if len(form.File["files"]) > define.BookToSkillMaxFileCount {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "too_many_files"))))
		return
	}

	// Save uploaded files
	uploadInfos, uploadErrors := common.SaveUploadedFileMulti(
		form, "files", define.BookToSkillFileLimitSize,
		adminUserId, "book_to_skill_file", define.BookToSkillAllowExt,
	)
	if len(uploadErrors) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(uploadErrors[0])))
		return
	}
	if len(uploadInfos) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "param_lack"))))
		return
	}

	// Create one task per uploaded file, execute via NSQ
	multiFile := len(uploadInfos) > 1
	taskIds := make([]int, 0, len(uploadInfos))

	for i, uploadInfo := range uploadInfos {
		// Suffix skill_name with index when multiple files
		taskSkillName := skillName
		if multiFile {
			taskSkillName = fmt.Sprintf("%s_%d", skillName, i+1)
		}

		// Create task record
		task := &common.BookToSkillTask{
			AdminUserId:     adminUserId,
			RobotId:         robotId,
			RobotKey:        robotKey,
			SkillName:       taskSkillName,
			ModelConfigId:   modelConfigId,
			UseModel:        useModel,
			TemplateContent: templateContent,
			Lang:            lang,
		}
		taskId, err := common.CreateBookToSkillTask(task, templateContent)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
			return
		}
		task.Id = taskId

		// Init task working directory + copy scripts/templates (requires taskId)
		if err := common.InitBookToSkillTaskDir(robotKey, taskId, taskSkillName); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}

		// Prepare single file to working_dir
		if err := common.PrepareBookToSkillFiles(task, []*define.UploadInfo{uploadInfo}); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}

		// Push NSQ message for async execution
		msg := tool.JsonEncodeNoError(map[string]any{"task_id": taskId})
		if err := common.AddJobs(define.BookToSkillTopic, msg); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
			return
		}

		taskIds = append(taskIds, taskId)
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{"task_ids": taskIds}, nil))
}

// GetBookToSkillTaskList returns paginated task list
func GetBookToSkillTaskList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	robotId := cast.ToInt(c.Query("robot_id"))
	status := cast.ToInt(c.Query("status"))
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	list, total, err := common.GetBookToSkillTaskList(adminUserId, robotId, status, page, pageSize)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil))
}

// GetBookToSkillTaskProgress returns task status and log
func GetBookToSkillTaskProgress(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.Query("task_id"))
	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		"status":      task.Status,
		"status_text": statusToText(task.Status, task.Lang),
		"error_msg":   task.ErrorMsg,
		"skill_dir":   task.SkillDir,
	}, nil))
}

// StopBookToSkillTask stops a running task
func StopBookToSkillTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.PostForm("task_id"))
	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}
	if task.Status != define.BookToSkillStatusPending && task.Status != define.BookToSkillStatusRunning {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "status_not_allowed"))))
		return
	}

	// Create .stop flag file so the running executor can detect and gracefully stop
	if err := common.CreateStopFlag(task.RobotKey, task.Id); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
		return
	}

	if err := common.StopBookToSkillTask(taskId, nil); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// RetryBookToSkillTask retries a failed or stopped task from scratch.
// It clears the previous working directory, resets all progress/result fields,
// re-initializes the task directory, and re-starts the generation pipeline.
func RetryBookToSkillTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.PostForm("task_id"))
	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}
	if task.Status != define.BookToSkillStatusFailed && task.Status != define.BookToSkillStatusStopped && task.Status != define.BookToSkillStatusPending {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "status_not_allowed"))))
		return
	}

	// Clear old directory, reset DB, re-init dir, re-prepare files
	if err := common.RetryBookToSkillTask(task); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	// Push NSQ message to start execution
	msg := tool.JsonEncodeNoError(map[string]any{"task_id": taskId})
	if err := common.AddJobs(define.BookToSkillTopic, msg); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// GetBookToSkillTaskLog returns task execution log
func GetBookToSkillTaskLog(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.Query("task_id"))
	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		"error_msg": task.ErrorMsg,
	}, nil))
}

// InstallBookToSkill installs a generated skill to a robot
func InstallBookToSkill(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.PostForm("task_id"))
	targetRobotId := cast.ToInt(c.PostForm("target_robot_id"))

	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}
	if task.Status != define.BookToSkillStatusSuccess {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "status_not_allowed"))))
		return
	}

	// Get target robot key
	if targetRobotId <= 0 {
		targetRobotId = task.RobotId
	}
	targetRobotKey, ok := common.GetClawbotRobotKey(c, adminUserId, int64(targetRobotId))
	if !ok {
		return
	}

	if err := common.InstallBookToSkill(task, targetRobotKey); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// DownloadBookToSkill downloads generated skill files as a zip
func DownloadBookToSkill(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	taskId := cast.ToInt(c.Query("task_id"))
	task, err := common.GetBookToSkillTask(taskId, adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}
	if task.SkillDir == "" || !tool.IsDir(task.SkillDir) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "no_data"))))
		return
	}

	// Create zip and stream as download
	zipPath := task.SkillDir + ".zip"
	if !tool.IsFile(zipPath) {
		if err := common.ZipDir(task.SkillDir, zipPath); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), "sys_err"))))
			return
		}
	}
	c.FileAttachment(zipPath, task.SkillName+".zip")
}

func filterValidUtf8(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if !define.SkillNameRegexp.MatchString(s) {
		return ""
	}
	if utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		s = string(runes[:maxLen])
	}
	return s
}

func statusToText(status int, lang string) string {
	switch status {
	case define.BookToSkillStatusPending:
		return i18n.Show(lang, "bts_status_pending")
	case define.BookToSkillStatusRunning:
		return i18n.Show(lang, "bts_status_running")
	case define.BookToSkillStatusSuccess:
		return i18n.Show(lang, "bts_status_success")
	case define.BookToSkillStatusFailed:
		return i18n.Show(lang, "bts_status_failed")
	case define.BookToSkillStatusStopped:
		return i18n.Show(lang, "bts_status_stopped")
	default:
		return i18n.Show(lang, "bts_status_unknown")
	}
}
