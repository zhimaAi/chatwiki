// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"strings"
)

const LocalUploadPrefix = `/upload/`

const ImageLimitSize = 100 * 1024           //100KB
const ImageAvatarLimitSize = 1024 * 1024    //1m
const LibFileLimitSize = 100 * 1024 * 1024  //100MB
const LibImageLimitSize = 2 * 1024 * 1024   //2M
const ChatImageLimitSize = 10 * 1024 * 1024 //10M

var ImageAllowExt = []string{`heic`, `gif`, `jfif`, `jpg`, `jpeg`, `jpe`, `png`, `swf`, `bmp`, `webp`}
var LibFileAllowExt = []string{`pdf`, `docx`, `ofd`, `txt`, `md`, `xlsx`, `csv`, `html`}
var FormFileAllowExt = []string{`json`, `xlsx`, `csv`}
var LibDocFileAllowExt = []string{`md`}
var QALibFileAllowExt = []string{`docx`, `xlsx`, `csv`}
var VideoAllowExt = []string{`mp4`}
var AudioAllowExt = []string{`mp3`}
var FAQLibFileAllowExt = []string{`md`, `docx`, `txt`}
var AllExt = make([]string, 0)

func init() {
	AllExt = append(AllExt, ImageAllowExt...)
	AllExt = append(AllExt, LibFileAllowExt...)
	AllExt = append(AllExt, FormFileAllowExt...)
	AllExt = append(AllExt, LibDocFileAllowExt...)
	AllExt = append(AllExt, QALibFileAllowExt...)
	AllExt = append(AllExt, VideoAllowExt...)
	AllExt = append(AllExt, AudioAllowExt...)
	AllExt = append(AllExt, FAQLibFileAllowExt...)
}

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
