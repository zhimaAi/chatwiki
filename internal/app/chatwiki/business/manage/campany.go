// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zhimaAi/go_tools/logs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

type SaveCompanyReq struct {
	Id   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name" binding:"max=15"`
}

func SaveCompany(c *gin.Context) {
	var (
		err      error
		insertId int64
	)
	//get params
	var req SaveCompanyReq
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	// login user
	user := GetLoginUserInfo(c)
	if user == nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	data := msql.Datas{
		"name":        req.Name,
		"update_time": time.Now().Unix(),
	}
	//headImg uploaded
	fileHeader, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageAvatarLimitSize, cast.ToInt(user["user_id"]), `company_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		data["avatar"] = uploadInfo.Link
	}
	m := msql.Model(define.TableCompany, define.Postgres)
	// save ..
	if req.Id > 0 {
		_, err = m.Where("id", cast.ToString(req.Id)).Update(data)
	} else {
		data["parent_id"] = user["user_id"]
		data["create_time"] = time.Now().Unix()
		insertId, err = m.Insert(data, "id")
	}
	if err != nil {
		common.FmtError(c, `company_save_err`)
		return
	}
	common.FmtOk(c, insertId)
}

type SaveTopNavigateReq struct {
	Id          int    `form:"id" json:"id" binding:"required"`
	TopNavigate string `form:"top_navigate" json:"top_navigate"`
}

func SaveTopNavigate(c *gin.Context) {
	var (
		err      error
		insertId int64
	)
	//get params
	var req SaveTopNavigateReq
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	var adminUserId = GetAdminUserId(c)
	// login user
	user := GetLoginUserInfo(c)
	if user == nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	data := msql.Datas{
		"top_navigate": req.TopNavigate,
		"update_time":  time.Now().Unix(),
	}

	m := msql.Model(define.TableCompany, define.Postgres)
	// save ..
	_, err = m.Where("id", cast.ToString(req.Id)).Where(`parent_id`, cast.ToString(adminUserId)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `company_save_err`)
		return
	}
	common.FmtOk(c, insertId)
}

func GetCompany(c *gin.Context) {
	var (
		data = msql.Params{}
		err  error
	)
	userId := getLoginUserId(c)
	if userId == 0 {
		data, err = msql.Model(define.TableCompany, define.Postgres).Order("id asc").Find()
	} else {
		parentId := GetAdminUserId(c)
		data, err = msql.Model(define.TableCompany, define.Postgres).Where("parent_id", cast.ToString(parentId)).Order("id asc").Find()
	}
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if data == nil {
		common.FmtError(c, `company_not_exist`)
		return
	}
	data[`neo4j_status`] = cast.ToString(common.GetNeo4jStatus(cast.ToInt(data[`parent_id`])))
	data[`wechat_ip`] = define.Config.WebService[`wechat_ip`]
	common.FmtOk(c, data)
}

func SaveAliOcr(c *gin.Context) {
	aliOcrKey := strings.TrimSpace(c.PostForm(`ali_ocr_key`))
	aliOcrSecret := strings.TrimSpace(c.PostForm(`ali_ocr_secret`))
	aliOcrSwitch := cast.ToInt(c.PostForm(`ali_ocr_switch`))

	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	if aliOcrSwitch != 1 && aliOcrSwitch != 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ali_ocr_switch`))))
		return
	}

	_, err := msql.Model(`company`, define.Postgres).
		Where(`parent_id`, cast.ToString(userId)).
		Update(msql.Datas{
			"ali_ocr_key":    aliOcrKey,
			"ali_ocr_secret": aliOcrSecret,
			"ali_ocr_switch": aliOcrSwitch,
		})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func CheckAliOcr(c *gin.Context) {
	aliOcrKey := strings.TrimSpace(c.PostForm(`ali_ocr_key`))
	aliOcrSecret := strings.TrimSpace(c.PostForm(`ali_ocr_secret`))
	if len(aliOcrKey) != 0 && len(aliOcrSecret) != 0 {
		if err := common.CheckAliOcr(aliOcrKey, aliOcrSecret); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `ali_ocr_check_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
