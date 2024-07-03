// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/tool"
)

func Upload(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	category := strings.TrimSpace(c.PostForm(`category`))
	if !tool.InArrayString(category, []string{`library_file`, `app_avatar`, `robot_avatar`, `icon`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `category`))))
		return
	}
	fileHeader, _ := c.FormFile(`file`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, category, define.ImageAllowExt)
	c.String(http.StatusOK, lib_web.FmtJson(uploadInfo, err))
}
