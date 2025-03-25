// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"path/filepath"
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
	if !tool.InArrayString(category, []string{`library_file`, `app_avatar`, `robot_avatar`, `icon`, `library_image`, `library_doc_image`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `category`))))
		return
	}
	fileHeader, _ := c.FormFile(`file`)
	allowExts := define.ImageAllowExt
	filesize := define.ImageLimitSize
	if tool.InArray(category, []string{`library_image`}) {
		filesize = define.LibImageLimitSize
	}
	if tool.InArray(category, []string{`library_doc_image`}) {
		ext := ""
		filesize = define.LibImageLimitSize
		if fileHeader != nil {
			ext = strings.ToLower(strings.TrimLeft(filepath.Ext(fileHeader.Filename), `.`))
		}
		if tool.InArray(ext, define.VideoAllowExt) {
			filesize = define.LibFileLimitSize
		}
		if tool.InArray(ext, define.AudioAllowExt) {
			filesize = define.LibFileLimitSize
		}
		allowExts = append(allowExts, define.VideoAllowExt...)
		allowExts = append(allowExts, define.AudioAllowExt...)
	}

	uploadInfo, err := common.SaveUploadedFile(fileHeader, filesize, userId, category, allowExts)
	c.String(http.StatusOK, lib_web.FmtJson(uploadInfo, err))
}
