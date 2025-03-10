// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"archive/zip"
	"bytes"
	"chatwiki/internal/app/chatwiki/define"
	"fmt"

	"github.com/zhimaAi/go_tools/tool"
)

func GenerateChatwikiZip(userId int, domain, outputZip string) error {
	sourceZip := define.AppRoot + `data/client_side.zip`
	configJson := `resources/app.asar.unpacked/resources/config.json`
	reader, err := zip.OpenReader(sourceZip)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	defer func(writer *zip.Writer) {
		_ = writer.Close()
	}(writer)
	for _, file := range reader.File {
		if file.Name == configJson {
			continue
		}
		if err = writer.Copy(file); err != nil {
			return err
		}
	}
	fw, err := writer.Create(configJson)
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`{"ADMIN_USER_ID":%d,"BASE_API_URL":"%s"}`, userId, domain)
	if _, err = fw.Write([]byte(content)); err != nil {
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}
	return tool.WriteFile(outputZip, buffer.String())
}
