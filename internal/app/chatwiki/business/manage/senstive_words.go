// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetSensitiveWordsListReq struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
}

func GetSensitiveWordsList(c *gin.Context) {
	var (
		req GetSensitiveWordsListReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	data, total, err := common.GetSensitiveWordsList(getAdminUserId(c), req.Page, req.Size)
	if err != nil {
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	common.FmtOk(c, map[string]any{
		"list":     data,
		"total":    total,
		`has_more`: req.Page*req.Size < int(total),
	})
}

type SaveSensitiveWordsReq struct {
	Id          int64  `form:"id" json:"id"`
	Words       string `form:"words" json:"words" binding:"required"`
	TriggerType int    `form:"trigger_type" json:"trigger_type"`
	RobotIds    string `form:"robot_ids"  json:"robot_ids"`
}

func SaveSensitiveWords(c *gin.Context) {
	var (
		adminUserId = getAdminUserId(c)
		req         SaveSensitiveWordsReq
		err         error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if len(strings.Split(strings.ReplaceAll(strings.ReplaceAll(req.Words, "\r\n", "\n"), "\r", "\n"), "\n")) > 1000 {
		common.FmtError(c, `file_data_limits`)
		return
	}
	id, err := common.SaveSensitiveWords(int64(adminUserId), req.Id, req.Words, req.RobotIds, req.TriggerType)
	if err != nil {
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.SenitiveWordsCacheHandle{AdminUserId: adminUserId})
	common.FmtOk(c, map[string]any{
		"id": id,
	})
}

type SwitchSensitiveWordsReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func SwitchSensitiveWords(c *gin.Context) {
	var (
		adminUserId = getAdminUserId(c)
		req         SwitchSensitiveWordsReq
		err         error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if _, err = common.SwitchSensitiveWords(adminUserId, req.Id); err != nil {
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.SenitiveWordsCacheHandle{AdminUserId: adminUserId})
	common.FmtOk(c, map[string]any{
		"id": req.Id,
	})
}

type DeleteSensitiveWordsReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func DeleteSensitiveWords(c *gin.Context) {
	var (
		adminUserId = getAdminUserId(c)
		req         DeleteSensitiveWordsReq
		err         error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if _, err = common.DeleteSensitiveWords(adminUserId, req.Id); err != nil {
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.SenitiveWordsCacheHandle{AdminUserId: adminUserId})
	common.FmtOk(c, nil)
}
