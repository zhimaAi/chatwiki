// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
	category := strings.TrimSpace(c.PostForm(`category`))
	if !tool.InArrayString(category, []string{`library_file`, `app_avatar`, `received_message_images`, `robot_avatar`, `http_tool_avatar`, `icon`, `library_image`, `library_doc_image`, `chat_image`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `category`))))
		return
	}
	var identity any
	if tool.InArray(category, []string{`chat_image`}) && len(c.GetHeader(`token`)) == 0 && len(c.Query(`token`)) == 0 {
		robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
		if !common.CheckRobotKey(robotKey) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
			return
		}
		if robot, err := common.GetRobotInfo(robotKey); err != nil || len(robot) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
			return
		}
		identity = robotKey
	} else { //管理有台操作,校验登录态
		var userId int
		if userId = GetAdminUserId(c); userId == 0 {
			return
		}
		identity = userId
	}
	fileHeader, _ := c.FormFile(`file`)
	allowExts := define.ImageAllowExt
	filesize := define.ImageLimitSize
	if tool.InArray(category, []string{`http_tool_avatar`, `library_image`, `received_message_images`}) {
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
	if tool.InArray(category, []string{`chat_image`}) {
		filesize = define.ChatImageLimitSize //聊天测试、webapp、网页插件的多模态输入-上传图片
	}

	uploadInfo, err := common.SaveUploadedFile(fileHeader, filesize, identity, category, allowExts)
	c.String(http.StatusOK, lib_web.FmtJson(uploadInfo, err))
}
