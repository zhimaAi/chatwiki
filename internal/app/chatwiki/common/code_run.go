// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/tool"
)

func RequestCodeRun(ctx context.Context, api string, data lib_define.CodeRunBody) (string, error) {
	domain := define.Config.WebService[`code_run`]
	body, err := tool.JsonEncode(data)
	if err != nil {
		return ``, err
	}
	timestamp := tool.Time2Int()
	nonce := tool.Random(20)
	sign := tool.MD5(body + cast.ToString(timestamp) + nonce + lib_define.CodeRunKey)
	link := fmt.Sprintf(`%s/%s?sign=%s&timestamp=%d&nonce=%s`, domain, api, sign, timestamp, nonce)
	request := curl.Post(link).Body(body).Header(`Content-Type`, `application/json`)
	resp, err := request.WithContext(ctx).Response()
	if err != nil {
		return ``, err
	}
	if resp.StatusCode != http.StatusOK {
		return ``, errors.New(fmt.Sprintf(`SYSTEM ERROR:%d`, resp.StatusCode))
	}
	code := lib_web.Response{}
	if err = request.ToJSON(&code); err != nil {
		return ``, err
	}
	if code.Res != lib_web.CommonSuccess {
		return ``, errors.New(code.Msg)
	}
	return cast.ToString(code.Data), nil
}
