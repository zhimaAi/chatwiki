// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/88250/lute/ast"

	"github.com/88250/lute"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

func OpenDoc(c *gin.Context) {
	dodKey := c.Param(`key`)
	info, err := common.GetLibDocInfoByDocKey(dodKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) <= 0 || cast.ToInt(info[`is_pub`]) != define.IsPub {
		common.FmtError(c, `no_data`)
		return
	}
	library, err := common.GetLibraryData(cast.ToInt(info[`library_id`]))
	if len(library) <= 0 || err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
		adminUserId := common.GetAdminUserId(c)
		if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(cast.ToInt(info[`library_id`]))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	var (
		docKey = cast.ToString(info[`doc_key`])
	)
	prev, next := common.FindPrevAndNext(catalog, docKey)
	htmlStr := luteMdToHTML(info[`content`])
	c.HTML(http.StatusOK, `document.html`, gin.H{
		`is_index`:       cast.ToInt(info[`is_index`]),
		`doc_key`:        docKey,
		`title`:          info[`title`],
		`seo_title`:      info[`seo_title`],
		`seo_desc`:       info[`seo_desc`],
		`seo_keywords`:   info[`seo_keywords`],
		`library_title`:  library[`library_name`],
		`library_key`:    library[`library_key`],
		`library_avatar`: library[`avatar`],
		`statistics_set`: template.HTML(library[`statistics_set`]),
		`body`:           template.HTML(htmlStr),
		`catalog`:        catalog,
		`prev_doc`:       prev,
		`next_doc`:       next,
	})
}

func luteMdToHTML(markdownText string) string {
	l := lute.New()
	l.SetFootnotes(true)
	l.SetAutoSpace(false)
	l.SetIndentCodeBlock(false)
	l.SetLinkBase("")
	l.SetGFMTaskListItem(true)

	// handle font color styles
	l.Md2HTMLRendererFuncs[ast.NodeText] = func(node *ast.Node, entering bool) (string, ast.WalkStatus) {
		if !entering {
			return "", ast.WalkContinue
		}
		text := node.TokensStr()

		// 匹配颜色和字体大小的正则表达式
		colorReg := `!!#([0-9a-fA-F]{6})`
		sizeReg := `!(\d+)`

		var style []string

		// 提取颜色
		if strings.Contains(text, "!!#") {
			if matches := regexp.MustCompile(colorReg).FindStringSubmatch(text); len(matches) > 1 {
				style = append(style, fmt.Sprintf("color: #%s", matches[1]))
				text = regexp.MustCompile(colorReg+` `).ReplaceAllString(text, "")
			}
		}

		// 提取字体大小
		if strings.Contains(text, "!") {
			if matches := regexp.MustCompile(sizeReg).FindStringSubmatch(text); len(matches) > 1 {
				size := cast.ToInt(matches[1])
				if size > 0 {
					style = append(style, fmt.Sprintf("font-size: %dpx", size))
					style = append(style, "line-height: 1em")
					text = regexp.MustCompile(sizeReg+` `).ReplaceAllString(text, "")
				}
			}
		}

		// 清理结尾的标记
		content := text
		if strings.HasSuffix(content, "!!") {
			content = content[:len(content)-2]
		}
		if strings.HasSuffix(content, "!") {
			content = content[:len(content)-1]
		}
		content = strings.TrimSpace(content)

		if len(style) > 0 && content != "" {
			return fmt.Sprintf(`<span style="%s">%s</span>`, strings.Join(style, "; "), content), ast.WalkContinue
		}

		return text, ast.WalkContinue
	}

	markdownText = processVideoSyntax(markdownText)
	return l.MarkdownStr("", markdownText)
}

func processVideoSyntax(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "!video[") {
			content := trimmedLine[7:]
			nameEnd := strings.Index(content, "]")
			if nameEnd != -1 && len(content) > nameEnd+1 {
				name := content[:nameEnd]
				url := strings.Trim(content[nameEnd+1:], "()")
				if url != "" {
					videoTag := fmt.Sprintf(`<video controls=“controls" src="%s">%s</video>`, url, name)
					result = append(result, videoTag)
					continue
				}
			}
		}
		if strings.HasPrefix(trimmedLine, "!audio[") {
			content := trimmedLine[7:]
			nameEnd := strings.Index(content, "]")
			if nameEnd != -1 && len(content) > nameEnd+1 {
				name := content[:nameEnd]
				url := strings.Trim(content[nameEnd+1:], "()")
				if url != "" {
					audioTag := fmt.Sprintf(`<audio controls=“controls" src="%s">%s</audio>`, url, name)
					result = append(result, audioTag)
					continue
				}
			}
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func OpenHome(c *gin.Context) {
	key := c.Param(`key`)
	libraryId := cast.ToInt(common.ParseLibraryKey(key))
	library, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(library) <= 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
		adminUserId := common.GetAdminUserId(c)
		if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
	}
	title := library[`library_name`]
	content := library[`library_intro`]
	info := common.GetLibDocIndex(libraryId)
	if len(info) > 0 {
		title = info[`title`]
		content = info[`content`]
		if content == "" {
			content = info[`draft_content`]
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(cast.ToInt(libraryId))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	// get question guide
	questionGuide := common.GetQuestionGuideList(libraryId)
	c.HTML(http.StatusOK, `home.html`, gin.H{
		`doc_key`:        info[`doc_key`],
		`doc_id`:         info[`id`],
		`title`:          title,
		`seo_title`:      info[`seo_title`],
		`seo_desc`:       info[`seo_desc`],
		`seo_keywords`:   info[`seo_keywords`],
		`library_title`:  library[`library_name`],
		`library_key`:    library[`library_key`],
		`library_avatar`: library[`avatar`],
		`statistics_set`: template.HTML(library[`statistics_set`]),
		`content`:        content,
		`catalog`:        catalog,
		`question_guide`: questionGuide,
	})
}

func OpenSearch(c *gin.Context) {
	typ := c.Param(`type`)
	libraryKey := c.Param(`lib_key`)
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	search, err := url.QueryUnescape(c.Query(`v`))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `param_invalid`, `search`)
		return
	}
	library, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(library) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// private access rights
	if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
		adminUserId := common.GetAdminUserId(c)
		if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
	}
	if typ == `html` {
		//  get catalog
		catalog, err := common.GetLibDocCateLogByCache(libraryId)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		c.HTML(http.StatusOK, `search.html`, gin.H{
			`library_title`:  library[`library_name`],
			`library_avatar`: library[`avatar`],
			`statistics_set`: template.HTML(library[`statistics_set`]),
			`library_key`:    library[`library_key`],
			`search`:         search,
			`catalog`:        catalog,
		})
	} else if typ == `query` {
		docInfo, _, _, err := common.LibDocSearch(common.GetLang(c), libraryId, search, library)
		if err != nil {
			common.FmtError(c, `sys_err`)
			return
		}
		common.FmtOk(c, docInfo)
	} else {
		common.FmtOk(c, nil)
	}
}

func OpenAiSummary(c *gin.Context) {
	libraryKey := c.Param(`lib_key`)
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	search, err := url.QueryUnescape(c.Query(`v`))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `param_invalid`, `search`)
		return
	}
	library, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(library) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// private access rights
	if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
		adminUserId := common.GetAdminUserId(c)
		if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
	}
	if cast.ToInt(library[`ai_summary`]) != define.AiSummary {
		common.FmtOk(c, nil)
		return
	}
	var (
		IsClose *bool
		wg      = sync.WaitGroup{}
	)
	chanStream := make(chan sse.Event)
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = common.LibDocAiSummary(common.GetLang(c), libraryId, search, library, chanStream, IsClose)
		if err != nil {
			common.FmtError(c, `sys_err`)
			return
		}
	}()
	go func() {
		c.Stream(func(_ io.Writer) bool {
			if event, ok := <-chanStream; ok {
				if data, ok := event.Data.(string); ok {
					event.Data = strings.ReplaceAll(data, "\r", ``)
				}
				c.SSEvent(event.Event, event.Data)
				return true
			}
			return false
		})
	}()
	wg.Wait()
}
