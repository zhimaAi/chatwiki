// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetWebToSkillTaskList(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	filter := define.WebToSkillTaskListFilter{
		Status: -1,
		Page:   1,
		Size:   define.WebToSkillTaskDefaultPageSize,
	}
	if err := c.ShouldBindQuery(&filter); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(filter, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.GetWebToSkillTaskList(common.GetLang(c), adminUserId, filter)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func CreateWebToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskCreateParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	params.Urls = append(params.Urls, c.PostFormArray(`urls[]`)...)
	if len(params.Urls) == 0 {
		params.Urls = common.ParseWebToSkillURLText(c.PostForm(`urls`))
	}
	id, err := common.CreateWebToSkillTask(common.GetLang(c), adminUserId, params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

func StopWebToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskIDParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.StopWebToSkillTask(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func RegenerateWebToSkillTask(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskIDParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	id, err := common.RegenerateWebToSkillTask(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

func GetWebToSkillTaskInfo(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskIDParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.GetWebToSkillTaskDetail(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func DownloadWebToSkillFile(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskIDParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	file, fileName, err := common.GetWebToSkillTaskDownloadFile(common.GetLang(c), adminUserId, params.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.FileAttachment(file, fileName)
}

func InstallWebToSkill(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.WebToSkillTaskInstallParams{}
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

	data, err := common.InstallWebToSkillTask(common.GetLang(c), adminUserId, params.ID, params.Overwrite)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}
