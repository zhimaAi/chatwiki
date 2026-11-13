// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPromptOptimizeContext(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	lang := common.GetLang(c)
	params := define.PromptOptimizeRobotParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, "param_err", middlewares.GetValidateErr(params, err, lang).Error())
		return
	}
	data, err := common.GetPromptOptimizeContext(lang, adminUserID, getLoginUserId(c), params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func GetPromptOptimizeRecords(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	lang := common.GetLang(c)
	params := define.PromptOptimizeRecordsParams{Size: define.PromptOptimizeDefaultPageSize}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, "param_err", middlewares.GetValidateErr(params, err, lang).Error())
		return
	}
	data, err := common.GetPromptOptimizeRecords(lang, adminUserID, getLoginUserId(c), params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func GetPromptOptimizeVersions(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	lang := common.GetLang(c)
	params := define.PromptOptimizeRobotParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, "param_err", middlewares.GetValidateErr(params, err, lang).Error())
		return
	}
	data, err := common.GetPromptOptimizeVersions(lang, adminUserID, getLoginUserId(c), params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func GetPromptOptimizeVersion(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	lang := common.GetLang(c)
	params := define.PromptOptimizeVersionParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.FmtError(c, "param_err", middlewares.GetValidateErr(params, err, lang).Error())
		return
	}
	data, err := common.GetPromptOptimizeVersion(lang, adminUserID, getLoginUserId(c), params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

func OptimizePrompt(c *gin.Context) {
	adminUserID := GetAdminUserId(c)
	if adminUserID == 0 {
		return
	}
	lang := common.GetLang(c)
	params := define.PromptOptimizeGenerateParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, "param_err", middlewares.GetValidateErr(params, err, lang).Error())
		return
	}
	prepared, err := common.PreparePromptOptimize(lang, adminUserID, getLoginUserId(c), params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Accel-Buffering", "no")
	events := make(chan sse.Event, 16)
	go common.RunPromptOptimize(c.Request.Context(), lang, prepared, events)
	for event := range events {
		if c.Request.Context().Err() == nil {
			c.SSEvent(event.Event, event.Data)
			c.Writer.Flush()
		}
	}
}
