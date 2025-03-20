// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"github.com/88250/lute"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	// 创建 Lute 实例
	l := lute.New()
	l.SetFootnotes(true)
	// 将 Markdown 转换为 HTML
	str := l.MarkdownStr("", markdownText)
	return str
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
