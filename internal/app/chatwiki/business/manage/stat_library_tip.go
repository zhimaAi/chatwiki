// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func StatLibraryTotal(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	stat, err := common.StatLibraryTotal(userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}

func StatLibraryDataSort(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := strings.TrimSpace(c.PostForm(`library_id`))
	page := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`page`, `1`)))
	size := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`size`, `100`)))
	if page < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `page`))))
	}
	if size < 0 || size > 1000 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `size`))))
	}

	beginDateYmd, endDateYmd := checkBeginEndDateYmd(c)
	if beginDateYmd == "" || endDateYmd == "" {
		return
	}
	stat, err := common.StatLibraryDataSort(userId, cast.ToInt(libraryId), page, size, beginDateYmd, endDateYmd)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}

func StatLibrarySort(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	beginDateYmd, endDateYmd := checkBeginEndDateYmd(c)
	if beginDateYmd == "" || endDateYmd == "" {
		return
	}
	page := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`page`, `1`)))
	size := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`size`, `100`)))
	if page < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `page`))))
	}
	if size < 0 || size > 1000 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `size`))))
	}
	stat, err := common.StatLibrarySort(userId, page, size, beginDateYmd, endDateYmd)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}

func StatLibraryDataRobotDetail(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := strings.TrimSpace(c.PostForm(`library_id`))
	if libraryId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `library_id`))))
		return
	}
	dataId := strings.TrimSpace(c.PostForm(`data_id`))
	if dataId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `data_id`))))
		return
	}
	beginDateYmd, endDateYmd := checkBeginEndDateYmd(c)
	if beginDateYmd == "" || endDateYmd == "" {
		return
	}
	stat, err := common.StatLibraryDataRobotDetail(userId, cast.ToInt(libraryId), cast.ToInt(dataId), beginDateYmd, endDateYmd)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}

func StatLibraryRobotDetail(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	libraryId := strings.TrimSpace(c.PostForm(`library_id`))
	if libraryId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `library_id`))))
		return
	}
	groupId := strings.TrimSpace(c.DefaultPostForm(`group_id`, `-1`))
	beginDateYmd, endDateYmd := checkBeginEndDateYmd(c)
	if beginDateYmd == "" || endDateYmd == "" {
		return
	}
	var stat []msql.Params
	var err error
	if cast.ToInt(groupId) != -1 {
		stat, err = common.StatLibraryRobotGroupDetail(userId, cast.ToInt(libraryId), groupId, beginDateYmd, endDateYmd)
	} else {
		stat, err = common.StatLibraryRobotDetail(userId, cast.ToInt(libraryId), beginDateYmd, endDateYmd)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}

func checkBeginEndDateYmd(c *gin.Context) (string, string) {
	beginDateYmd := strings.TrimSpace(c.PostForm(`begin_date_ymd`))
	if beginDateYmd == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `begin_date_ymd`))))
		return ``, ``
	}
	endDateYmd := strings.TrimSpace(c.PostForm(`end_date_ymd`))
	if endDateYmd == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date_ymd`))))
		return ``, ``
	}
	if !common.IsBasicDate(beginDateYmd) || !common.IsValidDate(beginDateYmd) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `begin_date_ymd`))))
		return ``, ``
	}
	if !common.IsBasicDate(endDateYmd) || !common.IsValidDate(endDateYmd) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date_ymd`))))
		return ``, ``
	}
	if beginDateYmd > endDateYmd {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `begin_date_ymd`))))
		return ``, ``
	}
	return beginDateYmd, endDateYmd
}

func StatLibraryGroupSort(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	beginDateYmd, endDateYmd := checkBeginEndDateYmd(c)
	if beginDateYmd == "" || endDateYmd == "" {
		return
	}
	page := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`page`, `1`)))
	size := cast.ToInt(strings.TrimSpace(c.DefaultPostForm(`size`, `100`)))
	libraryId := cast.ToInt(strings.TrimSpace(c.PostForm(`library_id`)))
	if page < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `page`))))
	}
	if size < 0 || size > 1000 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `size`))))
	}
	stat, err := common.StatLibraryGroupDataSort(userId, libraryId, page, size, beginDateYmd, endDateYmd)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(stat, nil))
	return
}
