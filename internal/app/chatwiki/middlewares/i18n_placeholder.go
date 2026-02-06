// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package middlewares

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/i18n"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// I18nPlaceholderMiddleware i18n placeholder replacement middleware
func I18nPlaceholderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get(`Accept`)
		if !isJSONContentType(accept) {
			return
		}
		// create response wrapper to capture response content
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

		// continue processing request
		c.Next()

		// check if response is JSON format
		contentType := c.Writer.Header().Get("Content-Type")
		if isJSONContentType(contentType) {
			// check if response body contains placeholders
			bodyStr := string(writer.body)
			if containsI18nPlaceholders(bodyStr) {
				// get request language
				lang := common.GetLang(c)

				// replace placeholders in JSON response
				modifiedBody := replaceI18nPlaceholders(writer.body, lang)

				// clear Content-Length header because content length may have changed
				c.Header("Content-Length", "")
				// set new response content (replace original response)
				c.DataFromReader(writer.statusCode, int64(len(modifiedBody)), contentType, strings.NewReader(string(modifiedBody)), nil)
			} else {
				// use original response content
				c.Header("Content-Length", "")
				c.DataFromReader(writer.statusCode, int64(len(writer.body)), contentType, strings.NewReader(string(writer.body)), nil)
			}
		} else {
			// use original response content
			c.Header("Content-Length", "")
			c.DataFromReader(writer.statusCode, int64(len(writer.body)), contentType, strings.NewReader(string(writer.body)), nil)
		}
	}
}

// isJSONContentType check if content type is JSON
func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false // if no Content-Type specified, assume it is not JSON
	}
	if strings.Contains(strings.ToLower(contentType), "text/plain") {
		return true // if it is text, assume it is JSON
	}
	return strings.Contains(strings.ToLower(contentType), "application/json")
}

// containsI18nPlaceholders check if string contains i18n placeholders
func containsI18nPlaceholders(text string) bool {
	return i18n.PlaceholderRegex.MatchString(text)
}

// responseWriter wrap gin.ResponseWriter to capture response content
type responseWriter struct {
	gin.ResponseWriter
	body        []byte
	statusCode  int
	wroteHeader bool
}

// Write capture response content, but do not write to original response immediately
func (rw *responseWriter) Write(data []byte) (int, error) {
	// append data to internal buffer, but do not write to original response
	rw.body = append(rw.body, data...)
	// return number of bytes written
	return len(data), nil
}

// WriteHeader capture status code, but do not write to original response immediately
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.statusCode = statusCode
		rw.wroteHeader = true
		// do not write to original response
	}
}

// Header return response header
func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// replaceI18nPlaceholders replace i18n placeholders in JSON response
func replaceI18nPlaceholders(jsonData []byte, lang string) []byte {
	// directly replace placeholders in JSON string
	jsonStr := string(jsonData)
	result := i18n.ReplacePlaceholdersInString(jsonStr, lang)
	return []byte(result)
}
