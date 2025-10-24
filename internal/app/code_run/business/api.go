// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

func GetCodeRunBody(c *gin.Context) (param lib_define.CodeRunBody, err error) {
	sign := strings.TrimSpace(c.Query(`sign`))
	if len(sign) == 0 {
		return param, errors.New(`sign empty`)
	}
	timestamp := cast.ToInt(c.Query(`timestamp`))
	if tool.AbsInt(tool.Time2Int()-timestamp) > 300 {
		return param, errors.New(`timestamp error`)
	}
	nonce := strings.TrimSpace(c.Query(`nonce`))
	if len(nonce) == 0 {
		return param, errors.New(`nonce empty`)
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return param, err
	}
	if len(body) == 0 {
		return param, errors.New(`body empty`)
	}
	if sign != tool.MD5(string(body)+cast.ToString(timestamp)+nonce+lib_define.CodeRunKey) {
		return param, errors.New(`sign error`)
	}
	if err = tool.JsonDecode(string(body), &param); err != nil {
		return param, err
	}
	return param, nil
}
