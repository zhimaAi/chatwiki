// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/casbin"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type AdminUserReq struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AdminLogin(c *gin.Context) {
	var (
		dPass = 0
		req   AdminUserReq
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	info, err := msql.Model(define.TableUser, define.Postgres).Where(`user_name`, req.UserName).Where("is_deleted", define.Normal).
		Where(fmt.Sprintf(`password=MD5(concat(%s,salt))`, msql.ToString(req.Password))).Field(`id,user_name,user_roles,avatar,nick_name,parent_id,login_switch,expire_time`).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `user_or_pwd_err`)
		return
	}
	if common.CheckUserLogin(cast.ToInt(info[`login_switch`]), cast.ToInt(info["expire_time"])) {
		common.FmtError(c, `client_side_cannot_login`)
		return
	}
	data, err := common.GetToken(info[`id`], info[`user_name`], info["parent_id"])
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	m := msql.Model(define.TableUser, define.Postgres)
	_, err = m.Where("id", info[`id`]).Update(msql.Datas{
		"login_time": time.Now().Unix(),
		"login_ip":   lib_web.GetClientIP(c),
	})
	if err != nil {
		logs.Error(err.Error())
	}

	// check is default password need to reset
	if cast.ToInt(data["user_roles"]) == casbin.Root && req.Password == define.DefaultPasswd {
		dPass = 1
	}
	data["d_pass"] = dPass
	data["user_roles"] = info["user_roles"]
	data["avatar"] = info["avatar"]
	data["nick_name"] = info["nick_name"]

	roleInfo, err := msql.Model(define.TableRole, define.Postgres).Where("id", info["user_roles"]).Find()
	if cast.ToInt(roleInfo["role_type"]) != 0 {
		data["role_type"] = roleInfo["role_type"]
	}

	common.FmtOk(c, data)
}

func GetAdminUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := common.ParseToken(token)
	if err != nil {
		c.String(http.StatusUnauthorized, lib_web.FmtJson(nil, err))
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		c.String(http.StatusUnauthorized, lib_web.FmtJson(nil, errors.New(`system error`)))
		return userId
	}
	if cast.ToInt(data["parent_id"]) <= 0 {
		return cast.ToInt(data["user_id"])
	}
	return cast.ToInt(data["parent_id"])
}

func getAdminUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := common.ParseToken(token)
	if err != nil {
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		return userId
	}
	if cast.ToInt(data["parent_id"]) <= 0 {
		return cast.ToInt(data["user_id"])
	}
	return cast.ToInt(data["parent_id"])
}

func getLoginUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := common.ParseToken(token)
	if err != nil {
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		return userId
	}
	return userId
}

func GetLoginUserInfo(c *gin.Context) jwt.MapClaims {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := common.ParseToken(token)
	if err != nil {
		return nil
	}
	return data
}

func GetUserInfo(userId string) (msql.Params, error) {
	return msql.Model(define.TableUser, define.Postgres).Where(`id`, userId).Field("id,user_name").Find()
}

func CheckLogin(c *gin.Context) {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	data, err := common.ParseToken(token)
	if err != nil {
		c.String(http.StatusUnauthorized, lib_web.FmtJson(nil, err))
		return
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		c.String(http.StatusUnauthorized, lib_web.FmtJson(nil, errors.New(`system error`)))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

type SaveProfileReq struct {
	Id int `form:"id" json:"id" binding:"required"`
	//Avatar        string `form:"avatar" json:"avatar" binding:"required_without=OldPassword"`
	OldPassword   string `form:"old_password" json:"old_password" binding:"required_with=Password,nefield=Password|len=0,max=32,omitempty" msg:"not equal password or length less 33"`
	Password      string `form:"password" json:"password" binding:"required_with=OldPassword,max=32,omitempty" `
	CheckPassword string `form:"check_password" json:"check_password" binding:"required_with=Password,eqfield=Password,omitempty"`
}

func SaveProfile(c *gin.Context) {
	var (
		err      error
		insertId int64
	)
	//get params
	var req SaveProfileReq
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	user := getLoginUserId(c)
	if req.Id != user {
		common.FmtError(c, `user_not_own`)
		return
	}
	data := msql.Datas{
		"update_time": time.Now().Unix(),
	}
	m := msql.Model(define.TableUser, define.Postgres)
	// save ..
	//headImg uploaded
	fileHeader, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageAvatarLimitSize, cast.ToInt(user), `user_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		data["avatar"] = uploadInfo.Link
	}
	if req.Password != "" {
		salt := tool.Random(20)
		password := tool.MD5(req.Password + salt)
		data["password"] = password
		data["salt"] = salt
		// check pass
		info, err := msql.Model(define.TableUser, define.Postgres).Where(`id`, cast.ToString(req.Id)).
			Where(fmt.Sprintf(`password=MD5(concat(%s,salt))`, msql.ToString(req.OldPassword))).Field(`id,user_name`).Find()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(info) == 0 {
			common.FmtError(c, `user_or_pwd_err`)
			return
		}
	}
	_, err = m.Where("id", cast.ToString(req.Id)).Update(data)
	if err != nil {
		common.FmtError(c, `user_save_err`)
		return
	}
	common.FmtOk(c, insertId)
}

func RefreshUserToken(c *gin.Context) {
	user := GetLoginUserInfo(c)
	if len(user) <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	data, err := common.GetToken(user[`user_id`], user[`user_name`], user["parent_id"])
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}

func LoginSwitch(c *gin.Context) {
	if getAdminUserId(c) == 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	userId := cast.ToInt(c.PostForm(`user_id`))
	if userId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	updata := fmt.Sprintf(`login_switch=1-login_switch,update_time=%v`, tool.Time2Int())
	m := msql.Model(define.TableUser, define.Postgres)
	// save ..
	_, err := m.Where("id", cast.ToString(userId)).Update2(updata)
	if err != nil {
		common.FmtError(c, `user_save_err`)
		return
	}
	common.FmtOk(c, userId)
}
