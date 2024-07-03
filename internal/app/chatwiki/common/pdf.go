// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/tool"
)

func ConvertToPdf(link string, userId int) (string, error) {
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(link), `.`))
	if ext == `pdf` {
		return link, nil
	}
	host, path, name := define.Config.WebService[`pdf`], ``, ``
	if len(host) == 0 {
		return ``, errors.New(`pdf convert not deploy`)
	}
	if ext == `md` {
		path, name = `/forms/chromium/convert/markdown`, `markdown.md`
	} else {
		path, name = `/forms/libreoffice/convert`, `file`+filepath.Ext(link)
	}
	request := curl.Post(host+path).PostFile(`files`, GetFileByLink(link), name)
	if ext == `md` {
		assist := define.AppRoot + `data/md_assist.html`
		request.PostFile(`htmlfile`, assist, `index.html`)
	}
	content, err := request.String()
	if err != nil {
		return ``, err
	}
	resp, err := request.Response()
	if err != nil {
		return ``, err
	}
	if resp.StatusCode != http.StatusOK {
		return ``, errors.New(content)
	}
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.pdf`, userId,
		`convert`, tool.Date(`Ym`), tool.MD5(content))
	pdfUrl, err := WriteFileByString(objectKey, content)
	if err != nil {
		return ``, err
	}
	return pdfUrl, nil
}
