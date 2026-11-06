// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	llm "github.com/zhimaAi/llm_adaptor/v2"
)

const (
	llmHTTPLogURLMaxLength     = 8192
	llmHTTPLogMethodMaxLength  = 64
	llmHTTPLogHeadersMaxLength = 256 * 1024
	llmHTTPLogErrorMaxLength   = 16384
	llmHTTPLogBodyMaxBytes     = 16 * 1024 * 1024
)

// SaveLLMHTTPLog persists one SDK HTTP round trip without affecting model calls.
func SaveLLMHTTPLog(_ context.Context, entry llm.HTTPLog) {
	data, err := llmHTTPLogData(entry)
	if err == nil {
		_, err = msql.Model("llm_http_log", define.Postgres).Insert(data)
	}
	if err != nil {
		logs.Error("save llm http log failed: status_code=%d, err=%s", entry.StatusCode, err.Error())
	}
}

func llmHTTPLogData(entry llm.HTTPLog) (msql.Datas, error) {
	requestHeaders, err := tool.JsonEncode(entry.RequestHeaders)
	if err != nil {
		return nil, fmt.Errorf("encode request headers: %w", err)
	}
	responseHeaders, err := tool.JsonEncode(entry.ResponseHeaders)
	if err != nil {
		return nil, fmt.Errorf("encode response headers: %w", err)
	}
	errorMessage := ""
	if entry.Error != nil {
		errorMessage = entry.Error.Error()
	}
	dayTime := entry.RequestTime
	if dayTime.IsZero() {
		dayTime = entry.EndTime
	}
	if dayTime.IsZero() {
		dayTime = time.Now()
	}
	return msql.Datas{
		"request_time":      llmHTTPLogTime(entry.RequestTime),
		"url":               truncateLLMHTTPLogText(entry.URL, llmHTTPLogURLMaxLength),
		"method":            truncateLLMHTTPLogText(entry.Method, llmHTTPLogMethodMaxLength),
		"request_headers":   truncateLLMHTTPLogText(requestHeaders, llmHTTPLogHeadersMaxLength),
		"request_body":      truncateLLMHTTPLogBody(entry.RequestBody),
		"response_time":     llmHTTPLogTime(entry.ResponseTime),
		"end_time":          llmHTTPLogTime(entry.EndTime),
		"duration":          entry.Duration.Nanoseconds(),
		"status_code":       entry.StatusCode,
		"response_headers":  truncateLLMHTTPLogText(responseHeaders, llmHTTPLogHeadersMaxLength),
		"response_body":     truncateLLMHTTPLogBody(entry.ResponseBody),
		"response_body_eof": entry.ResponseBodyEOF,
		"error":             truncateLLMHTTPLogText(errorMessage, llmHTTPLogErrorMaxLength),
		"is_error":          entry.Error != nil || entry.StatusCode >= 400,
		"belong_day":        cast.ToInt(tool.Date("Ymd", dayTime.Unix())),
	}, nil
}

func llmHTTPLogTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func truncateLLMHTTPLogText(value string, maxLength int) string {
	// PostgreSQL text columns reject NUL bytes and invalid UTF-8.
	value = strings.ReplaceAll(strings.ToValidUTF8(value, "\uFFFD"), "\x00", `\u0000`)
	length := 0
	for end := range value {
		if length == maxLength {
			return value[:end]
		}
		length++
	}
	return value
}

func truncateLLMHTTPLogBody(body []byte) []byte {
	if len(body) == 0 {
		return []byte{}
	}
	if len(body) > llmHTTPLogBodyMaxBytes {
		return body[:llmHTTPLogBodyMaxBytes]
	}
	return body
}
