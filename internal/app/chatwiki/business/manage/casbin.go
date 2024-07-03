// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/casbin"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func CheckPermission(c *gin.Context) {
	user := GetLoginUserId(c)
	if user <= 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	var data = make(map[string]interface{})

	info, err := msql.Model(define.TableUser, define.Postgres).Where("id", cast.ToString(user)).Field("user_roles").Find()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if info == nil {
		common.FmtError(c, `user_not_exist`)
		return
	}
	rules, err := casbin.Handler.GetPolicyForUser(info["user_roles"])
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	rolePermission := make([]string, 0)
	for _, rule := range rules {
		if len(rule) >= 1 {
			if strings.ContainsAny(rule[1], "/") {
				continue
			}
			rolePermission = append(rolePermission, rule[1])
		}
	}
	data["role_permission"] = rolePermission
	data["menu"] = getMenu()
	data["user_roles"] = info["user_roles"]
	common.FmtOk(c, data)
}

type GetMnuRes struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	UniKey   string       `json:"uni_key"`
	Path     string       `json:"path"`
	ParentId int          `json:"parent_id"`
	Child    []*GetMnuRes `json:"child"`
}

func getMenu() map[string]*GetMnuRes {

	var res = make(map[string]*GetMnuRes)
	m := msql.Model(define.TableMenu, define.Postgres)
	data, err := m.Where("is_deleted", define.Normal).Where("parent_id", "0").Select()
	if err != nil {
		logs.Error("getMenu err:%s", err)
		return res
	}
	for _, v := range data {
		if v["parent_id"] == "0" {
			res[v["id"]] = &GetMnuRes{
				Id:     v["id"],
				Name:   v["name"],
				UniKey: v["uni_key"],
				Path:   v["path"],
				Child:  make([]*GetMnuRes, 0),
			}
		}
	}
	return res
}

func GetMenu(c *gin.Context) {
	common.FmtOk(c, getMenu())
}

type AddMenuReq struct {
	Id       int    `form:"id" json:"id"`
	Name     string `form:"name" json:"name"`
	UniKey   string `form:"uni_key" json:"uni_key"`
	Path     string `form:"path" json:"path" binding:"required"`
	ParentId int    `form:"parent_id" json:"parent_id" binding:"oneof=0 1 2 3,omitempty"`
}

func SaveMenu(c *gin.Context) {
	var (
		req AddMenuReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	userInfo := GetLoginUserInfo(c)
	m := msql.Model(define.TableMenu, define.Postgres)
	var insertId = int64(req.Id)
	if req.Id > 0 {
		menu, err := m.Where("id", cast.ToString(req.Id)).Find()
		if err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		updates := msql.Datas{
			"update_time":  time.Now().Unix(),
			"operate_id":   userInfo["user_id"],
			"operate_name": userInfo["user_name"],
			"name":         req.Name,
			"parent_id":    cast.ToInt(req.ParentId),
		}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.UniKey != "" {
			updates["uni_key"] = req.UniKey
		}
		if req.Path != "" {
			updates["path"] = fmt.Sprintf("%s,%s", menu["path"], req.Path)
		}
		if req.ParentId != 0 {
			updates["parent_id"] = req.ParentId
		}
		_, err = m.Where("id", cast.ToString(req.Id)).Update(updates)
	} else {
		insertId, err = m.Insert(msql.Datas{
			"name":         req.Name,
			"path":         req.Path,
			"operate_id":   userInfo["user_id"],
			"operate_name": userInfo["user_name"],
			"uni_key":      req.UniKey,
			"parent_id":    req.ParentId,
		})
	}

	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, insertId)
}

type DelMenuReq struct {
	Id string `json:"id" binding:"required"`
}

func DelMenu(c *gin.Context) {
	var req DelMenuReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	_, err := msql.Model(define.TableMenu, define.Postgres).Where("id", req.Id).Update(msql.Datas{
		"is_deleted": define.Deleted,
	})
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, "ok")
}

type GetUserListReq struct {
	Page   int    `form:"page" json:"page"`
	Size   int    `form:"size" json:"size" binding:"max=200"`
	Search string `form:"search" json:"search"`
}

func GetUserList(c *gin.Context) {
	var req GetUserListReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 200
	}
	parentId := GetAdminUserId(c)
	//parentId := getParentId(userId)
	page := req.Page
	m := msql.Model(define.TableUser, define.Postgres).Where("parent_id", cast.ToString(parentId)).Where("is_deleted", define.Normal)
	if req.Search != "" {
		str := "%" + req.Search + "%"
		m.Where(fmt.Sprintf("( user_name like '%v' or nick_name like '%v')", str, str))
	} else {
		m.WhereOr("id", cast.ToString(parentId))
	}
	data, total, err := m.Order("id asc").Paginate(page, req.Size)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	roles, err := msql.Model(define.TableRole, define.Postgres).Select()
	roleMap := make(map[string]msql.Params)
	for _, role := range roles {
		roleMap[role["id"]] = role
	}
	for _, user := range data {
		user["role"] = ""
		if name, ok := roleMap[user["user_roles"]]["name"]; ok {
			user["role"] = name
		}
		user["salt"] = ""
	}
	result := map[string]interface{}{
		"total":    total,
		"list":     data,
		"has_more": page*req.Size < total,
	}
	common.FmtOk(c, result)
}

type SaveUserReq struct {
	Id            int    `form:"id" json:"id" binding:"omitempty"`
	UserName      string `form:"user_name" json:"user_name" binding:"required,max=100,alphanum|containsany= _-."`
	NickName      string `form:"nick_name" json:"nick_name" binding:"required,max=100"`
	UserRoles     int    `form:"user_roles" json:"user_roles" binding:"required,min=2"`
	Password      string `form:"password" json:"password" binding:"max=32,omitempty"`
	CheckPassword string `form:"check_password" json:"check_password" binding:"required_with=Password,eqfield=Password,omitempty"`
}

func SaveUser(c *gin.Context) {
	var (
		err      error
		insertId int64
	)
	//get params
	var req SaveUserReq
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
	roleId := cast.ToString(req.UserRoles)
	// check user is deleted
	roleInfo, _ := msql.Model(define.TableRole, define.Postgres).Where("id", roleId).Field("id,is_deleted").Find()
	if roleInfo == nil {
		common.FmtError(c, `role_not_exist`)
		return
	}
	if roleInfo["is_deleted"] == define.Deleted {
		common.FmtError(c, `role_is_deleted`)
		return
	}
	data := msql.Datas{
		"nick_name":    req.NickName,
		"operate_id":   user["user_id"],
		"operate_name": user["user_name"],
		"user_roles":   req.UserRoles,
		"update_time":  time.Now().Unix(),
	}
	//headImg uploaded
	fileHeader, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageAvatarLimitSize, cast.ToInt(user["user_id"]), `user_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		data["avatar"] = uploadInfo.Link
	}

	if req.Password != "" {
		salt := tool.Random(20)
		password := tool.MD5(req.Password + salt)
		data["password"] = password
		data["salt"] = salt
	}
	m := msql.Model(define.TableUser, define.Postgres)
	// save ..
	if req.Id > 0 {
		_, err = m.Where("id", cast.ToString(req.Id)).Update(data)
	} else {
		if req.Password == "" {
			common.FmtError(c, `param_err`, "password")
			return
		}
		data["parent_id"] = GetAdminUserId(c)
		data["user_name"] = req.UserName
		data["create_time"] = time.Now().Unix()
		insertId, err = m.Insert(data, "id")
	}
	if err != nil {
		common.FmtError(c, `user_save_err`, catchErr(err))
		return
	}
	// save user role
	if success, err := casbin.Handler.UpdateUserRole(req.UserName, roleId); err != nil || !success {
		common.FmtError(c, `user_save_err`, "role add error")
		return
	}
	common.FmtOk(c, insertId)
}
func catchErr(err error) string {
	msg := err.Error()
	if strings.ContainsAny(msg, "duplicate key") {
		return "account exists"
	}
	return err.Error()
}

type ResetPassReq struct {
	Id            int    `form:"id" json:"id" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required,min=6,max=32"`
	CheckPassword string `form:"check_password" json:"check_password" binding:"required,eqfield=Password"`
}

func ResetPass(c *gin.Context) {
	var (
		err      error
		insertId int64
	)
	//get params
	var req ResetPassReq
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	salt := tool.Random(20)
	password := tool.MD5(req.Password + salt)
	data := msql.Datas{
		"password":    password,
		"salt":        salt,
		"update_time": time.Now().Unix(),
	}
	m := msql.Model(define.TableUser, define.Postgres)
	// reset ..
	_, err = m.Where("id", cast.ToString(req.Id)).Update(data)
	if err != nil {
		common.FmtError(c, `user_save_err`)
		return
	}
	common.FmtOk(c, insertId)
}

type DelUserReq struct {
	Id int `form:"id" json:"id" binding:"required,min=2"`
}

func DeleteUser(c *gin.Context) {
	var req = DelUserReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	// login user
	user := GetLoginUserInfo(c)
	if user == nil {
		common.FmtError(c, `user_no_login`)
		return
	}
	// del user
	m := msql.Model(define.TableUser, define.Postgres)
	userInfo, err := m.Where("id", cast.ToString(req.Id)).Field("id,user_name").Find()
	if err != nil {
		common.FmtError(c, `param_lock`)
		return
	}
	if len(userInfo) <= 0 {
		common.FmtOk(c, "ok")
		return
	}
	result, _ := casbin.MatchRootFunc(userInfo["user_name"])
	isRoot := result.(bool)
	if isRoot {
		common.FmtError(c, `root_user`)
		return
	}
	_, err = m.Where("id", cast.ToString(req.Id)).Update(msql.Datas{
		"update_time":  time.Now().Unix(),
		"operate_id":   user["user_id"],
		"operate_name": user["user_name"],
		"is_deleted":   1,
	})
	if err != nil {
		common.FmtError(c, `user_del_err`)
		return
	}
	// del user from role
	if _, err := casbin.Handler.DeleteUserRole(0, userInfo["user_name"]); err != nil {
		common.FmtError(c, `user_del_err`)
		return
	}
	common.FmtOk(c, "ok")
}

type (
	GetUserReq struct {
		Id int `form:"id" json:"id" binding:"required"`
	}
	GetUserRes struct {
		Id             int      `form:"id" json:"id"`
		UserName       string   `form:"user_name" json:"user_name"`
		CreateTime     int64    `form:"create_time" json:"create_time"`
		UpdateTime     int64    `form:"update_time" json:"update_time"`
		LoginIp        string   `form:"login_ip" json:"login_ip"`
		NickName       string   `form:"nick_name" json:"nick_name"`
		Avatar         string   `form:"avatar" json:"avatar"`
		IsDeleted      int      `form:"is_deleted" json:"is_deleted"`
		OperateId      int      `form:"operate_id" json:"operate_id"`
		OperateName    string   `form:"operate_name" json:"operate_name"`
		UserRole       string   `form:"user_role" json:"user_role"`
		RolePermission []string `form:"role_permission" json:"role_permission"`
	}
)

func GetUser(c *gin.Context) {
	var req = GetUserReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	data, err := msql.Model(define.TableUser, define.Postgres).Where("id", cast.ToString(req.Id)).Find()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	if data == nil {
		common.FmtError(c, `user_not_exist`)
		return
	}
	// get role and permission
	rules, err := casbin.Handler.GetPolicyForUser(data["user_roles"])
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	rolePermission := make([]string, 0)
	for _, rule := range rules {
		if len(rule) >= 1 {
			if strings.ContainsAny(rule[1], "/") {
				continue
			}
			rolePermission = append(rolePermission, rule[1])
		}
	}
	userRoleInfo := GetUserRes{
		Id:             cast.ToInt(data["id"]),
		UserName:       data["user_name"],
		CreateTime:     cast.ToInt64(data["create_time"]),
		UpdateTime:     cast.ToInt64(data["update_time"]),
		LoginIp:        data["login_ip"],
		NickName:       data["nick_name"],
		Avatar:         data["avatar"],
		IsDeleted:      cast.ToInt(data["is_deleted"]),
		OperateId:      cast.ToInt(data["operate_id"]),
		OperateName:    data["operate_name"],
		UserRole:       data["user_roles"],
		RolePermission: rolePermission,
	}
	common.FmtOk(c, userRoleInfo)
}

type (
	GetRoleListReq struct {
		Page   int    `form:"page" json:"page"`
		Size   int    `form:"size" json:"size" binding:"max=200"`
		Search string `form:"search" json:"search"`
	}
)

func GetRoleList(c *gin.Context) {
	var req = GetRoleListReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 200
	}
	page := req.Page
	parentId := GetAdminUserId(c)
	//parentId := getParentId(userId)
	m := msql.Model(define.TableRole, define.Postgres).Where("is_deleted", define.Normal).Where("parent_id", "in", fmt.Sprintf("%s,0", cast.ToString(parentId))).Order("id asc")
	if req.Search != "" {
		m.Where([]string{"name", "like", req.Search}...)
	}
	data, total, err := m.Order("id asc").Paginate(page, req.Size)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	// get role info
	result := map[string]interface{}{
		"total":    total,
		"list":     data,
		"has_more": page*req.Size < total,
	}
	common.FmtOk(c, result)
}

type (
	SaveRoleReq struct {
		Id   int    `form:"id" json:"id" binding:"omitempty"`
		Name string `form:"name" json:"name" binding:"required,max=100"`
		Mark string `form:"mark" json:"mark" binding:"max=500"`
		//RoleType       int    `form:"role_type" json:"role_type" binding:"required,oneof=2 3"`
		RoleType int    `form:"role_type" json:"role_type"`
		UniKeys  string `form:"uni_keys" json:"uni_keys"`
	}
)

func SaveRole(c *gin.Context) {
	var (
		req      = SaveRoleReq{}
		err      error
		insertId int64
	)
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
		"operate_id":   user["user_id"],
		"operate_name": user["user_name"],
		"role_type":    req.RoleType,
		"name":         req.Name,
		"mark":         req.Mark,
		"update_time":  time.Now().Unix(),
	}
	m := msql.Model(define.TableRole, define.Postgres)
	// save ..
	roleId := cast.ToString(req.Id)
	if req.Id > 0 {
		if req.Id <= casbin.Root {
			common.FmtError(c, `root_user`)
			return
		}
		if _, err = m.Where("id", roleId).Update(data); err != nil {
			logs.Error("SaveRole:%+v,err:%s", data, err)
			common.FmtError(c, `role_save_err`, err.Error())
			return
		}
		_, err = casbin.Handler.DelRoleRules(roleId)
	} else {
		data["create_time"] = time.Now().Unix()
		data["create_name"] = user["user_name"]
		data["parent_id"] = GetAdminUserId(c)
		insertId, err = m.Insert(data, "id")
		roleId = cast.ToString(insertId)
	}
	if err != nil {
		common.FmtError(c, `role_save_err`)
		return
	}
	rolePermissions := make([][]string, 0)
	rolePermissionsReq := make(map[string]bool)
	if len(req.UniKeys) > 0 {
		uniKeys := strings.Split(req.UniKeys, ",")
		for _, item := range uniKeys {
			rolePermissionsReq[item] = true
			rolePermissions = append(rolePermissions, []string{roleId, item, "GET"})
		}
		_, err = casbin.Handler.AddPolicies(rolePermissions)
		if err != nil {
			common.FmtError(c, `role_save_err`)
			return
		}
	}
	common.FmtOk(c, insertId)
}

type DelRoleReq struct {
	Id int `form:"id" json:"id" binding:"required,min=4"`
}

func DelRole(c *gin.Context) {
	var req = DelRoleReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	// login user
	user := GetLoginUserInfo(c)
	if user == nil {
		common.FmtError(c, `user_no_login`)
		return
	}
	// del role
	m := msql.Model(define.TableRole, define.Postgres)
	roleInfo, err := m.Where("id", cast.ToString(req.Id)).Field("id").Find()
	if err != nil {
		common.FmtError(c, `param_lock`)
		return
	}
	if len(roleInfo) <= 0 {
		common.FmtOk(c, "ok")
		return
	}
	users, err := casbin.Handler.GetUsersForRole(roleInfo["id"])
	if len(users) > 0 || err != nil {
		common.FmtError(c, `role_has_users`)
		return
	}
	_, err = m.Where("id", cast.ToString(req.Id)).Update(msql.Datas{
		"operate_id":   user["user_id"],
		"operate_name": user["user_name"],
		"update_time":  time.Now().Unix(),
		"is_deleted":   define.Deleted,
	})
	if err != nil {
		common.FmtError(c, `role_del_err`)
		return
	}
	// del user from role
	if _, err := casbin.Handler.DeleteRole(roleInfo["id"]); err != nil {
		common.FmtError(c, `role_del_err`)
		return
	}
	common.FmtOk(c, "ok")
}

type (
	GetRoleReq struct {
		Id int `form:"id" json:"id" binding:"required"`
	}
	GetRoleRes struct {
		Id             int      `json:"id"`
		Name           string   `json:"name"`
		Mark           string   `json:"mark"`
		IsDeleted      int      `json:"is_deleted"`
		CreateTime     int64    `json:"create_time"`
		UpdateTime     int64    `json:"update_time"`
		OperateName    string   `json:"operate_name"`
		OperateId      int      `json:"operate_id"`
		RolePermission []string `json:"role_permission"`
		RoleType       int      `json:"role_type"`
	}
)

func GetRole(c *gin.Context) {
	var req = GetRoleReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	m := msql.Model(define.TableRole, define.Postgres)
	roleInfo, err := m.Where("id", cast.ToString(req.Id)).Field("id,name,mark,role_type").Find()
	if err != nil {
		common.FmtError(c, `param_lock`)
		return
	}
	if len(roleInfo) <= 0 {
		common.FmtOk(c, "ok")
		return
	}
	// get role and permission
	rules, err := casbin.Handler.GetPolicyForUser(roleInfo["id"])
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	rolePermission := make([]string, 0)
	for _, rule := range rules {
		if len(rule) >= 1 {
			if strings.ContainsAny(rule[1], "/") {
				continue
			}
			rolePermission = append(rolePermission, rule[1])
		}
	}
	roleData := &GetRoleRes{
		Id:             cast.ToInt(roleInfo["id"]),
		Name:           roleInfo["name"],
		Mark:           roleInfo["mark"],
		RoleType:       cast.ToInt(roleInfo["role_type"]),
		IsDeleted:      cast.ToInt(roleInfo["is_deleted"]),
		CreateTime:     cast.ToInt64(roleInfo["create_time"]),
		UpdateTime:     cast.ToInt64(roleInfo["update_time"]),
		OperateName:    roleInfo["operate_name"],
		OperateId:      cast.ToInt(roleInfo["operate_id"]),
		RolePermission: rolePermission,
	}
	common.FmtOk(c, roleData)
}
