// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, `pong`)
}

func NoMethod(c *gin.Context) {
	c.String(http.StatusOK, `NoMethod`)
}

func NoRoute(c *gin.Context) {
	c.String(http.StatusOK, `NoRoute`)
}

func Test(c *gin.Context) {
	c.String(http.StatusOK, `test`)
}

func Test1(c *gin.Context) {
	c.String(http.StatusOK, `test111`)
}

func TestDomain(c *gin.Context) {
	// save
	fileName := c.Query(`file`)
	m := msql.Model(`chat_ai_file_info`, define.Postgres)
	info, err := m.Where(`file_name`, fileName).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	c.String(http.StatusOK, `%v`, info[`file_content`])
}
