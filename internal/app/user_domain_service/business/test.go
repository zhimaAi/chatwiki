// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
