// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"os"
)

var IsDev bool
var Env string

const AppRoot = `internal/app/chatwiki/`
const UploadDir = AppRoot + `upload/`
const TemplateDir = `/html-template/open/doc/`
const TemplateStaticDir = `/static/html-template/open/doc/`
const AppName = "chatwiki"

// Version 会在编译时通过 -ldflags 注入
var Version = "V2025-12-09"

func GetTemplatesPath() string {
	filePath, _ := os.Getwd()
	return filePath + TemplateDir
}

func GetTemplatesStaticPath() string {
	filePath, _ := os.Getwd()
	return filePath + TemplateStaticDir
}
