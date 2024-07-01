// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/casbin"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func GetCompany(c *gin.Context) {
	adminUserId := cast.ToInt(c.Query(`admin_user_id`))
	if adminUserId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	company, err := msql.Model(define.TableCompany, define.Postgres).Where(`parent_id`, cast.ToString(adminUserId)).Order(`id asc`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(company, nil))
}

func ClientSideGetRobotList(c *gin.Context) {
	adminUserId := cast.ToInt(c.Query(`admin_user_id`))
	if adminUserId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if common.ClientSideNeedLogin(adminUserId) {
		user := manage.GetLoginUserInfo(c)
		if len(user) == 0 {
			return
		}
		//check the ownership of the login user
		parentId := cast.ToInt(user[`parent_id`])
		if parentId == 0 {
			parentId = cast.ToInt(user[`user_id`])
		}
		if parentId != adminUserId {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `user_not_exist`))))
			return
		}
		//check ClientSideManage permission
		if !common.CheckPermission(cast.ToInt(user[`user_id`]), define.ClientSideManage) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `client_side_not_permission`))))
			return
		}
	}
	//robot list
	list, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_name,robot_intro,robot_avatar,robot_key`).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//company info
	company, err := msql.Model(define.TableCompany, define.Postgres).Where(`parent_id`, cast.ToString(adminUserId)).Order(`id asc`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`h5_domain`: define.Config.WebService[`h5_domain`], `list`: list, `company`: company}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func ClientSideLogin(c *gin.Context) {
	adminUserId := cast.ToInt(c.PostForm(`admin_user_id`))
	userName := strings.TrimSpace(c.PostForm(`user_name`))
	password := strings.TrimSpace(c.PostForm(`password`))
	if adminUserId <= 0 || len(userName) == 0 || len(password) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := msql.Model(define.TableUser, define.Postgres).Where(`user_name`, userName).Where("is_deleted", define.Normal).
		Where(fmt.Sprintf(`password=MD5(concat(%s,salt))`, msql.ToString(password))).Field(`id,user_name,user_roles,avatar,nick_name,parent_id`).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `user_or_pwd_err`)
		return
	}
	//check the ownership of the login user
	parentId := cast.ToInt(info[`parent_id`])
	if parentId == 0 {
		parentId = cast.ToInt(info[`id`])
	}
	if parentId != adminUserId {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `user_not_exist`))))
		return
	}
	//check ClientSideManage permission
	if !common.CheckPermission(cast.ToInt(info[`id`]), define.ClientSideManage) {
		extra := map[string]any{`type`: `client_side_cannot_login`}
		c.String(http.StatusOK, lib_web.FmtJson(extra, errors.New(i18n.Show(common.GetLang(c), `client_side_cannot_login`))))
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
	var dPass = 0
	if cast.ToInt(data["user_roles"]) == casbin.Root && password == define.DefaultPasswd {
		dPass = 1
	}
	data["d_pass"] = dPass
	data["user_roles"] = info["user_roles"]
	data["avatar"] = info["avatar"]
	data["nick_name"] = info["nick_name"]
	common.FmtOk(c, data)
}
