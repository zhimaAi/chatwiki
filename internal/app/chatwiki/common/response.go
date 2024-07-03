// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FmtError(c *gin.Context, msg string, params ...string) {
	data := struct{}{}
	err := errors.New(i18n.Show(GetLang(c), msg, params))
	c.String(http.StatusOK, lib_web.FmtJson(data, err))
	c.Abort()
}

func FmtOk(c *gin.Context, data interface{}) {
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}
