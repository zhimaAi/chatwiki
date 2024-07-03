// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import "strings"

const LocalUploadPrefix = `/upload/`

const ImageLimitSize = 100 * 1024         //100KB
const ImageAvatarLimitSize = 1024 * 1024  //1m
const LibFileLimitSize = 10 * 1024 * 1024 //10MB

var ImageAllowExt = []string{`heic`, `gif`, `jpg`, `jpeg`, `png`, `swf`, `bmp`, `webp`}
var LibFileAllowExt = []string{`pdf`, `docx`, `txt`, `md`, `xlsx`, `csv`, `html`}

func IsTableFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == `xlsx` || ext == `csv`
}
