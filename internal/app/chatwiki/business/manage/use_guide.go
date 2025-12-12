// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/pkg/lib_web"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UseGuideProcess(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	guideData, err := common.GetUseGuide(adminUserId, common.GetLang(c))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(guideData, nil))
}
