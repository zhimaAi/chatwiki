// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetE2bConf(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	params := define.E2bConfGetParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.GetE2bConf(common.GetLang(c), adminUserId, params.RobotKey)
	c.String(http.StatusOK, lib_web.FmtJson(data, err))
}

func SaveE2bConf(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	params := define.E2bConfParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	data, err := common.SaveE2bConf(common.GetLang(c), adminUserId, params)
	c.String(http.StatusOK, lib_web.FmtJson(data, err))
}
