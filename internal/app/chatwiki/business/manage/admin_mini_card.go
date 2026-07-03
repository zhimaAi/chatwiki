// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/middlewares"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

type AdminMiniCardRequest struct {
	ID       int    `form:"id" json:"id"`
	Title    string `form:"title" json:"title" binding:"required"`
	Appid    string `form:"appid" json:"appid" binding:"required"`
	PagePath string `form:"page_path" json:"page_path" binding:"required"`
	ThumbURL string `form:"thumb_url" json:"thumb_url" binding:"required"`
}

type AdminMiniCardListRequest struct {
	Keyword string `form:"keyword" json:"keyword"`
	Appid   string `form:"appid" json:"appid"`
	Page    int    `form:"page" json:"page"`
	Size    int    `form:"size" json:"size"`
}

func GetAdminMiniCardList(c *gin.Context) {
	var req AdminMiniCardListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	list, total, err := common.GetAdminMiniCardList(adminUserId, strings.TrimSpace(req.Keyword), strings.TrimSpace(req.Appid), req.Page, req.Size)
	if err != nil {
		logs.Error("GetAdminMiniCardList error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	result := make([]map[string]any, 0, len(list))
	for _, item := range list {
		result = append(result, common.FormatAdminMiniCard(item))
	}
	common.FmtOk(c, map[string]any{
		`list`:  result,
		`total`: total,
		`page`:  req.Page,
		`size`:  req.Size,
	})
}

func AddAdminMiniCard(c *gin.Context) {
	var req AdminMiniCardRequest
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	id, err := common.SaveAdminMiniCard(0, adminUserId, strings.TrimSpace(req.Title), strings.TrimSpace(req.Appid), strings.TrimSpace(req.PagePath), strings.TrimSpace(req.ThumbURL))
	if err != nil {
		logs.Error("AddAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{`id`: id})
}

func DeleteAdminMiniCard(c *gin.Context) {
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	card, err := common.GetAdminMiniCard(id, adminUserId)
	if err != nil {
		logs.Error("GetAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(card) == 0 {
		common.FmtError(c, `admin_mini_card_not_exist`)
		return
	}

	relationCount, err := common.CountAdminMiniCardRelations(adminUserId, id)
	if err != nil {
		logs.Error("CountAdminMiniCardRelations error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if relationCount > 0 {
		common.FmtError(c, `admin_mini_card_in_use`)
		return
	}

	if err = common.DeleteAdminMiniCard(id, adminUserId); err != nil {
		logs.Error("DeleteAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, nil)
}

func UpdateAdminMiniCard(c *gin.Context) {
	var req AdminMiniCardRequest
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.ID <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	card, err := common.GetAdminMiniCard(req.ID, adminUserId)
	if err != nil {
		logs.Error("GetAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(card) == 0 {
		common.FmtError(c, `admin_mini_card_not_exist`)
		return
	}

	id, err := common.SaveAdminMiniCard(req.ID, adminUserId, strings.TrimSpace(req.Title), strings.TrimSpace(req.Appid), strings.TrimSpace(req.PagePath), strings.TrimSpace(req.ThumbURL))
	if err != nil {
		logs.Error("UpdateAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{`id`: id})
}

func GetAdminMiniCard(c *gin.Context) {
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	card, err := common.GetAdminMiniCard(id, adminUserId)
	if err != nil {
		logs.Error("GetAdminMiniCard error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(card) == 0 {
		common.FmtError(c, `admin_mini_card_not_exist`)
		return
	}
	common.FmtOk(c, common.FormatAdminMiniCard(card))
}
