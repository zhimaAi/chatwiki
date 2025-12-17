// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func OpenDoc(c *gin.Context) {
	dodKey := c.Param(`key`)
	info, err := common.GetLibDocInfoByDocKey(dodKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) <= 0 {
		c.HTML(http.StatusOK, `404.html`, gin.H{})
		return
	}
	previewKey := c.Query(`preview`)
	if cast.ToInt(info[`is_pub`]) != define.IsPub && previewKey == "" {
		c.HTML(http.StatusOK, `404.html`, gin.H{})
		return
	}
	library, err := common.GetLibraryData(cast.ToInt(info[`library_id`]))
	if len(library) <= 0 || err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	// get preview data
	if previewKey != "" {
		if !common.CheckPreviewOpenDoc(library[`library_key`], previewKey) {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
		info = common.GetLibDocInfo(cast.ToInt(info[`id`]))
	} else {
		if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
			adminUserId := common.GetAdminUserId(c)
			if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
				c.HTML(http.StatusOK, `404.html`, gin.H{})
				return
			}
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(cast.ToInt(info[`library_id`]), previewKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	var (
		docKey = cast.ToString(info[`doc_key`])
	)
	prev, next := common.FindPrevAndNext(catalog, docKey)
	if previewKey != "" {
		draftInfo := common.GetLibDocFields(cast.ToInt(info[`id`]), `is_draft,content,draft_content`)
		if len(draftInfo) > 0 {
			if cast.ToInt(draftInfo[`is_draft`]) == define.IsDraft {
				info[`content`] = draftInfo[`draft_content`]
			}
		}
	}
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
		`preview_key`:    previewKey,
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
		c.HTML(http.StatusOK, `404.html`, gin.H{})
		return
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
	previewKey := c.Query(`preview`)
	// get preview data
	if previewKey != "" {
		if !common.CheckPreviewOpenDoc(library[`library_key`], previewKey) {
			c.HTML(http.StatusOK, `404.html`, gin.H{})
			return
		}
	} else {
		if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
			adminUserId := common.GetAdminUserId(c)
			if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
				c.HTML(http.StatusOK, `404.html`, gin.H{})
				return
			}
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(libraryId, previewKey)
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
		`preview_key`:    previewKey,
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
		catalog, err := common.GetLibDocCateLogByCache(libraryId, "")
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
		//  get catalog
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
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err = common.LibDocAiSummary(common.GetLang(c), libraryId, search, library, chanStream, IsClose)
		if err != nil {
			common.FmtError(c, `sys_err`)
			return
		}
	}()
	go func() {
		defer wg.Done()
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

func OpenBindLibList(c *gin.Context) {
	libraryKey := c.Query(`library_key`)
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	domain := lib_web.GetRequestDomain(c)
	if libraryId > 0 {
		library, _ := common.GetLibraryData(libraryId)
		if cast.ToString(library[`share_url`]) == "" {
			common.FmtOk(c, []msql.Params{library})
			return
		}
		domain = cast.ToString(library[`share_url`])
	}
	result, err := common.GetOpenBindLibList(domain, cast.ToString(define.OpenLibraryType))
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, result)
}

func OpenDocApi(c *gin.Context) {
	dodKey := c.Param(`key`)
	info, err := common.GetLibDocInfoByDocKey(dodKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) <= 0 {
		common.FmtOk(c, gin.H{
			`is_404`: 1,
		})
		return
	}
	previewKey := c.Query(`preview`)
	if cast.ToInt(info[`is_pub`]) != define.IsPub && previewKey == "" {
		common.FmtOk(c, gin.H{
			`is_404`: 1,
		})
		return
	}
	library, err := common.GetLibraryData(cast.ToInt(info[`library_id`]))
	if len(library) <= 0 || err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	// get preview data
	if previewKey != "" {
		if !common.CheckPreviewOpenDoc(library[`library_key`], previewKey) {
			common.FmtOk(c, gin.H{
				`is_404`: 1,
			})
			return
		}
		info = common.GetLibDocInfo(cast.ToInt(info[`id`]))
	} else {
		if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
			adminUserId := common.GetAdminUserId(c)
			if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
				common.FmtOk(c, gin.H{
					`is_404`: 1,
				})
				return
			}
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(cast.ToInt(info[`library_id`]), previewKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	var (
		docKey = cast.ToString(info[`doc_key`])
	)
	prev, next := common.FindPrevAndNext(catalog, docKey)

	if prev != nil && prev.IsDir == 1 {
		// 循环10次，找到第一个非目录的节点
		for i := 0; i < 10; i++ {
			// 如果第一个节点是个空文件夹
			if prev == nil || prev.DocKey == "" {
				break
			}
			prev, _ = common.FindPrevAndNext(catalog, prev.DocKey)
			if prev != nil && prev.IsDir == 0 {
				break
			}
		}
	}

	if next != nil && next.IsDir == 1 {
		// 循环10次，找到第一个非目录的节点
		for i := 0; i < 10; i++ {
			//如果最后一个节点是空文件夹
			if next == nil || next.DocKey == "" {
				break
			}
			_, next = common.FindPrevAndNext(catalog, next.DocKey)
			if next != nil && next.IsDir == 0 {
				break
			}
		}
	}

	if previewKey != "" {
		draftInfo := common.GetLibDocFields(cast.ToInt(info[`id`]), `is_draft,content,draft_content`)
		if len(draftInfo) > 0 {
			if cast.ToInt(draftInfo[`is_draft`]) == define.IsDraft {
				info[`content`] = draftInfo[`draft_content`]
			}
		}
	}
	common.FmtOk(c, gin.H{
		`is_404`:                  0,
		`is_index`:                cast.ToInt(info[`is_index`]),
		`doc_key`:                 docKey,
		`title`:                   info[`title`],
		`seo_title`:               info[`seo_title`],
		`seo_desc`:                info[`seo_desc`],
		`seo_keywords`:            info[`seo_keywords`],
		`create_time`:             info[`create_time`],
		`update_time`:             info[`update_time`],
		`library_title`:           library[`library_name`],
		`library_key`:             library[`library_key`],
		`library_avatar`:          library[`avatar`],
		`icon_template_config_id`: library[`icon_template_config_id`],
		`statistics_set`:          template.HTML(library[`statistics_set`]),
		`body`:                    template.HTML(info[`content`]),
		`catalog`:                 catalog,
		`prev_doc`:                prev,
		`next_doc`:                next,
		`preview_key`:             previewKey,
	})
}

func OpenHomeApi(c *gin.Context) {
	key := c.Param(`key`)
	var (
		libraryId int
	)
	libraryId = cast.ToInt(common.ParseLibraryKey(key))
	library, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(library) <= 0 {
		common.FmtOk(c, gin.H{
			`is_404`: 1,
		})
		return
	}

	title := library[`library_name`]
	content := library[`library_intro`]
	info := common.GetLibDocIndex(libraryId)
	if len(info) > 0 {
		title = info[`title`]
		content = info[`content`]
		if content == "" || cast.ToInt(info[`is_draft`]) == define.IsDraft {
			content = info[`draft_content`]
		}
	}
	previewKey := c.Query(`preview`)
	// get preview data
	if previewKey != "" {
		if !common.CheckPreviewOpenDoc(library[`library_key`], previewKey) {
			common.FmtOk(c, gin.H{
				`is_404`: 1,
			})
			return
		}
	} else {
		if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
			adminUserId := common.GetAdminUserId(c)
			if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
				common.FmtOk(c, gin.H{
					`is_404`: 1,
				})
				return
			}
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(libraryId, previewKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	// get question guide
	questionGuide := common.GetQuestionGuideList(libraryId)

	// 给默认的banner
	if len(cast.ToString(info[`banner_img_url`])) == 0 {
		info[`banner_img_url`] = define.DefaultLibDocBanner
	}

	// 渲染快捷文档
	quickDocContentValue := []msql.Datas{}
	if len(cast.ToString(info[`quick_doc_content`])) > 0 {
		var quickDocArr []msql.Datas // 快速文档内容数组
		// 将quickDocIdArr中的doc_id提取出来，存到quickDocIds中
		var quickDocIdArr []msql.Datas
		var quickDocIds []int
		if err := tool.JsonDecode(cast.ToString(info[`quick_doc_content`]), &quickDocIdArr); err != nil {
			logs.Error(err.Error())
			return
		}
		for _, quickDocId := range quickDocIdArr {
			if quickDocId[`doc_id`] != nil {
				quickDocIds = append(quickDocIds, int(cast.ToInt(quickDocId[`doc_id`])))
			}
		}
		for _, quickDocId := range quickDocIds {
			quickDocItem := msql.Datas{}
			quickDocInfo := common.GetLibDocInfo(quickDocId)
			if len(quickDocInfo) <= 0 {
				continue
			}
			// 这里要把删除的文档干掉
			if cast.ToInt(quickDocInfo[`delete_time`]) != 0 {
				continue
			}

			quickDocItem[`id`] = quickDocInfo[`id`]
			quickDocItem[`pid`] = quickDocInfo[`pid`]
			quickDocItem[`doc_key`] = quickDocInfo[`doc_key`]
			quickDocItem[`is_dir`] = quickDocInfo[`is_dir`]
			quickDocItem[`doc_icon`] = quickDocInfo[`doc_icon`]
			quickDocItem[`is_draft`] = quickDocInfo[`is_draft`]
			quickDocItem[`title`] = quickDocInfo[`title`]
			quickDocItem[`content`] = quickDocInfo[`content`]
			quickDocItem[`create_time`] = quickDocInfo[`create_time`]
			quickDocItem[`update_time`] = quickDocInfo[`update_time`]
			quickDocItem[`children`] = common.GetLibDocCateLog(quickDocId, libraryId, false)
			// 最后添加成数组
			quickDocArr = append(quickDocArr, quickDocItem)
		}
		quickDocContentValue = quickDocArr
	}

	// 再处理quick_doc_content
	quickDocContent := []msql.Datas{}
	if len(cast.ToString(info[`quick_doc_content`])) > 0 {
		err = tool.JsonDecode(cast.ToString(info[`quick_doc_content`]), &quickDocContent)
		if err != nil {
			logs.Error("json encode quick_doc_content error:%v", quickDocContent)
		}
	}

	common.FmtOk(c, gin.H{
		`is_404`:                  0,
		`doc_key`:                 info[`doc_key`],
		`doc_id`:                  info[`id`],
		`title`:                   title,
		`banner_img_url`:          info[`banner_img_url`],
		`quick_doc_content`:       quickDocContent,
		`quick_doc_content_value`: quickDocContentValue,
		`seo_title`:               info[`seo_title`],
		`seo_desc`:                info[`seo_desc`],
		`seo_keywords`:            info[`seo_keywords`],
		`library_title`:           library[`library_name`],
		`library_key`:             library[`library_key`],
		`library_avatar`:          library[`avatar`],
		`icon_template_config_id`: library[`icon_template_config_id`],
		`statistics_set`:          template.HTML(library[`statistics_set`]),
		`content`:                 content,
		`catalog`:                 catalog,
		`question_guide`:          questionGuide,
		`preview_key`:             previewKey,
	})
}

func OpenSearchApi(c *gin.Context) {
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
		common.FmtOk(c, gin.H{
			`is_404`: 1,
		})
		return
	}
	// private access rights
	if cast.ToInt(library[`access_rights`]) != define.OpenLibraryAccessRights {
		adminUserId := common.GetAdminUserId(c)
		if adminUserId <= 0 || cast.ToInt(library[`admin_user_id`]) != adminUserId {
			common.FmtOk(c, gin.H{
				`is_404`: 1,
			})
			return
		}
	}
	//  get catalog
	catalog, err := common.GetLibDocCateLogByCache(libraryId, "")
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, gin.H{
		`is_404`:                  0,
		`library_title`:           library[`library_name`],
		`library_avatar`:          library[`avatar`],
		`statistics_set`:          template.HTML(library[`statistics_set`]),
		`library_key`:             library[`library_key`],
		`icon_template_config_id`: library[`icon_template_config_id`],
		`search`:                  search,
		`catalog`:                 catalog,
	})
}
