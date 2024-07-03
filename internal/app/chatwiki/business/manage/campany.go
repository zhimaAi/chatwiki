// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"time"

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
		common.FmtError(c, `user_no_login`)
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
	common.FmtOk(c, data)
}
