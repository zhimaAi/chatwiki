// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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

// Version will be injected via -ldflags during compilation
var Version = "V2025-12-09"
var IsPublicNetWork = 1

func GetTemplatesPath() string {
	filePath, _ := os.Getwd()
	return filePath + TemplateDir
}

func GetTemplatesStaticPath() string {
	filePath, _ := os.Getwd()
	return filePath + TemplateStaticDir
}
