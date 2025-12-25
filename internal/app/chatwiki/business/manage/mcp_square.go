// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

// GetMcpSquareTypeList 获取模板分类列表
func GetMcpSquareTypeList(c *gin.Context) {
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetMcpTypeList`, nil)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	respData, ok := resp.Data.([]interface{})
	if !ok {
		err = errors.New(`invalid data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(respData, nil))
}

// GetMcpSquareList 获取模板列表
func GetMcpSquareList(c *gin.Context) {
	body := make(map[string]any)

	title := strings.TrimSpace(c.Query(`title`))
	source := strings.TrimSpace(c.Query(`source`))
	filterType := cast.ToInt(c.Query(`filter_type`))
	id := cast.ToInt(c.Query(`id`))
	page := cast.ToInt(c.Query(`page`))
	size := cast.ToInt(c.Query(`size`))
	if id > 0 {
		body[`id`] = id
	}
	if len(title) > 0 {
		body[`title`] = title
	}
	if len(source) > 0 {
		body[`source`] = source
	}
	if filterType > 0 {
		body[`filter_type`] = filterType
	}
	if page > 0 {
		body[`page`] = page
	}
	if size > 0 {
		body[`size`] = size
	}
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetMcpList`, body)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(resp.Data, nil))
}
