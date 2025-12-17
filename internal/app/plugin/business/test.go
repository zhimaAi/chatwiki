// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
