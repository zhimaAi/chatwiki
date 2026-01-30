// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package middlewares

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/i18n"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// I18nPlaceholderMiddleware 国际化占位符替换中间件
func I18nPlaceholderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get(`Accept`)
		if !isJSONContentType(accept) {
			return
		}
		// 创建响应包装器来捕获响应内容
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

		//logs.Debug(`处理返回请求`)
		// 继续处理请求
		c.Next()

		// 检查响应是否为JSON格式
		contentType := c.Writer.Header().Get("Content-Type")
		if isJSONContentType(contentType) {
			// 检查响应体是否包含占位符
			bodyStr := string(writer.body)
			if containsI18nPlaceholders(bodyStr) {
				// 获取请求语言
				lang := common.GetLang(c)

				// 替换JSON响应中的占位符
				modifiedBody := replaceI18nPlaceholders(writer.body, lang)

				//logs.Debug(`结束处理 - 有占位符需要替换`)
				// 清空Content-Length头部，因为内容长度可能已更改
				c.Header("Content-Length", "")
				// 设置新的响应内容（替换原始响应）
				c.DataFromReader(writer.statusCode, int64(len(modifiedBody)), contentType, strings.NewReader(string(modifiedBody)), nil)
			} else {
				//logs.Debug(`结束处理 - 无占位符，使用原始响应`)
				// 使用原始响应内容
				c.Header("Content-Length", "")
				c.DataFromReader(writer.statusCode, int64(len(writer.body)), contentType, strings.NewReader(string(writer.body)), nil)
			}
		} else {
			//logs.Debug(`结束处理 - 非JSON内容类型，使用原始响应`)
			// 使用原始响应内容
			c.Header("Content-Length", "")
			c.DataFromReader(writer.statusCode, int64(len(writer.body)), contentType, strings.NewReader(string(writer.body)), nil)
		}
	}
}

// isJSONContentType 检查内容类型是否为JSON
func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false // 如果没有指定Content-Type，假设它不是 JSON
	}
	if strings.Contains(strings.ToLower(contentType), "text/plain") {
		return true // 如果是text文本的，假设它是JSON
	}
	return strings.Contains(strings.ToLower(contentType), "application/json")
}

// containsI18nPlaceholders 检查字符串是否包含国际化占位符
func containsI18nPlaceholders(text string) bool {
	return i18n.PlaceholderRegex.MatchString(text)
}

// responseWriter 包装 gin.ResponseWriter 以捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body        []byte
	statusCode  int
	wroteHeader bool
}

// Write 捕获响应内容，但不立即写入原始响应
func (rw *responseWriter) Write(data []byte) (int, error) {
	// 将数据追加到内部缓冲区，但不写入原始响应
	rw.body = append(rw.body, data...)
	// 返回写入的字节数
	return len(data), nil
}

// WriteHeader 捕获状态码，但不立即写入原始响应
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.statusCode = statusCode
		rw.wroteHeader = true
		// 不写入原始响应
	}
}

// Header 返回响应头
func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// replaceI18nPlaceholders 替换JSON响应中的国际化占位符
func replaceI18nPlaceholders(jsonData []byte, lang string) []byte {
	// 直接在JSON字符串上进行占位符替换
	jsonStr := string(jsonData)
	result := i18n.ReplacePlaceholdersInString(jsonStr, lang)
	return []byte(result)
}
