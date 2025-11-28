// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"strings"
)

const LocalUploadPrefix = `/upload/`

const ImageLimitSize = 100 * 1024          //100KB
const ImageAvatarLimitSize = 1024 * 1024   //1m
const LibFileLimitSize = 100 * 1024 * 1024 //100MB
const LibImageLimitSize = 2 * 1024 * 1024  // 2M

var ImageAllowExt = []string{`heic`, `gif`, `jpg`, `jpeg`, `png`, `swf`, `bmp`, `webp`}
var LibFileAllowExt = []string{`pdf`, `docx`, `ofd`, `txt`, `md`, `xlsx`, `csv`, `html`}
var FormFileAllowExt = []string{`json`, `xlsx`, `csv`}
var LibDocFileAllowExt = []string{`md`}
var QALibFileAllowExt = []string{`docx`, `xlsx`, `csv`}
var VideoAllowExt = []string{`mp4`}
var AudioAllowExt = []string{`mp3`}
var FAQLibFileAllowExt = []string{`md`, `docx`, `txt`}

func IsTableFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `xlsx` || ext == `csv`
}

func IsDocxFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `docx`
}

func IsOfdFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `ofd`
}

func IsTxtFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `txt`
}

func IsMdFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `md`
}

func IsPdfFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `pdf`
}

func IsHtmlFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `html`
}
