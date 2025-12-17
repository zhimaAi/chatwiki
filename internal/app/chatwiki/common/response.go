// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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

func FmtErrorWithCode(c *gin.Context, code int, msg string, params ...string) {
	data := struct{}{}
	err := errors.New(i18n.Show(GetLang(c), msg, params))
	c.String(code, lib_web.FmtJsonWithCode(code, data, err))
	c.Abort()
}

func FmtOk(c *gin.Context, data interface{}) {
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func FmtBridgeResponse(c *gin.Context, data any, httpStatus int, err error) {
	if httpStatus == -1 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	} else if httpStatus == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(data, nil))
	} else {
		c.String(httpStatus, lib_web.FmtJsonWithCode(httpStatus, struct{}{}, err))
	}
}

type response struct {
	Object    string `json:"object"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	RequestId string `json:"requestId"`
}

func FmtOpenAiErr(c *gin.Context, code int, msg string, params ...string) {
	if code == 0 {
		code = http.StatusBadRequest
	}
	err := errors.New(i18n.Show(GetLang(c), msg, params))
	c.JSON(http.StatusOK, response{
		Object:  "error",
		Message: err.Error(),
		Code:    code,
	})
	c.Abort()
}

func FmtOpenAiOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
	c.Abort()
}
