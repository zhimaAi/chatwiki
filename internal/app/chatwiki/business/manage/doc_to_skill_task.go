// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetDocToSkillTaskList(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	filter := define.DocToSkillTaskListFilter{Status: -1, Page: 1, Size: define.DocToSkillTaskDefaultPageSize}
	if err := c.ShouldBindQuery(&filter); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(filter, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.GetDocToSkillTaskList(common.GetLang(c), adminUserId, filter)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func CreateDocToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskCreateParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	form, err := c.MultipartForm()
	if err != nil || form == nil || len(form.File[`files`]) == 0 {
		common.FmtError(c, `param_err`, `files`)
		return
	}
	if len(form.File[`files`]) > define.DocToSkillTaskMaxFileCount {
		common.FmtError(c, `param_err`, `files`)
		return
	}
	sourceFiles, uploadErrors := common.SaveUploadedFileMulti(
		form, `files`, define.DocToSkillTaskFileLimitSize, adminUserId, `doc_to_skill_file`, define.DocToSkillTaskAllowExt,
	)
	if len(uploadErrors) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(uploadErrors[0])))
		return
	}
	id, err := common.CreateDocToSkillTask(common.GetLang(c), adminUserId, params, sourceFiles)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

func StopDocToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskIDParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.StopDocToSkillTask(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func RegenerateDocToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskIDParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	id, err := common.RegenerateDocToSkillTask(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

func GetDocToSkillTaskInfo(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskIDParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.GetDocToSkillTaskDetail(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func DownloadDocToSkillFile(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskIDParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	file, fileName, err := common.GetDocToSkillTaskDownloadFile(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.FileAttachment(file, fileName)
}

func InstallDocToSkill(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.DocToSkillTaskInstallParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	lockKey := define.LockPreKey + `ClawbotUserSkill.` + cast.ToString(adminUserId)
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		common.FmtError(c, `op_lock`)
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	data, err := common.InstallDocToSkillTask(common.GetLang(c), adminUserId, params.ID, params.Overwrite)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}
