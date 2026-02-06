// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/pkg/thumbnail"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Thumbnail(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		filename = c.GetHeader("X-File-Name")
	}

	if filename == "" {
		c.String(http.StatusInternalServerError, "empty file name")
		return
	}

	// 2. Get Request Body
	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "read body error:"+err.Error())
		return
	}
	defer c.Request.Body.Close()

	if len(content) == 0 {
		c.String(http.StatusInternalServerError, "body empty ")
		return
	}

	// 3. generate thumbnail
	thumbBytes, thumbName, err := thumbnail.GenerateThumbnail(content, filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "system error: "+err.Error())
		return
	}

	// 4. response
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", thumbName))
	c.Data(http.StatusOK, "image/png", thumbBytes)

}
