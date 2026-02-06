// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.
package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func GenThumbnail(realFilePath, userId string) string {

	returnDir := filepath.Join(fmt.Sprintf(`%s/upload/chat_ai/%s/%s/%s`, filepath.Dir(define.AppRoot), userId, `library_file`, tool.Date(`Ym`)))
	thumbFolderPath := filepath.Join(fmt.Sprintf(`/upload/chat_ai/%s/%s/%s/`, userId, `library_file`, tool.Date(`Ym`)))

	baseName := filepath.Base(realFilePath)

	file, err := os.Open(realFilePath)
	if err != nil {
		logs.Error("open file error:" + err.Error())
		return ""
	}
	defer file.Close()

	domain := define.Config.WebService[`gen_thumbnail`]

	requestURL := fmt.Sprintf("%s/thumbnail?filename=%s", domain, baseName)

	req, err := http.NewRequest("POST", requestURL, file)
	if err != nil {
		logs.Error("request client error:" + err.Error())
		return ""
	}
	// set Content-Type
	req.Header.Set("Content-Type", "application/octet-stream")

	// --- D. send request ---
	client := &http.Client{
		Timeout: 30 * time.Second, // set timeout 30 seconds
	}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("request  error:" + err.Error())
		return ""
	}
	defer resp.Body.Close()

	// --- E. check status ---
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logs.Error("system error (Code %d): %s", resp.StatusCode, string(bodyBytes))
		return ""
	}

	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	fileName := nameWithoutExt + "_thumb.png"

	thumbPath := filepath.Join(thumbFolderPath, fileName)

	dir := filepath.Dir(returnDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	fileByte, _ := io.ReadAll(resp.Body)

	_err := os.WriteFile(filepath.Join(returnDir, fileName), fileByte, 0644)
	if _err != nil {
		logs.Error("save thumbnail error：" + _err.Error())
		return ""
	}

	return thumbPath

}
