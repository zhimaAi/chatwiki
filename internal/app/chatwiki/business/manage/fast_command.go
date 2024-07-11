// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

type SaveFastCommandReq struct {
	Id      int    `form:"id" json:"id"`
	RobotID int    `form:"robot_id" json:"robot_id" binding:"required"`
	Title   string `form:"title" json:"title" binding:"required,max=20"`
	Typ     int    `form:"typ" json:"typ" binding:"required,oneof=1 2"`
	Content string `form:"content" json:"content" binding:"required,max=500"`
}

func SaveFastCommand(c *gin.Context) {
	var (
		req SaveFastCommandReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	m := msql.Model(define.TableFastCommand, define.Postgres)
	var insertId = int64(req.Id)
	if req.Id > 0 {
		updates := msql.Datas{
			"robot_id":      req.RobotID,
			"title":         req.Title,
			"typ":           req.Typ,
			"content":       req.Content,
			"admin_user_id": cast.ToString(adminUserId),
			"update_time":   time.Now().Unix(),
		}
		_, err = m.Where("id", cast.ToString(req.Id)).Update(updates)
	} else {
		insertId, err = m.Insert(msql.Datas{
			"robot_id":      req.RobotID,
			"title":         req.Title,
			"typ":           req.Typ,
			"content":       req.Content,
			"admin_user_id": cast.ToString(adminUserId),
			"create_time":   time.Now().Unix(),
			"update_time":   time.Now().Unix(),
		})
	}
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, insertId)
}

type GetFastCommandListReq struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
	OpenId   string `form:"open_id" json:"open_id"`
}

func GetFastCommandList(c *gin.Context) {
	var (
		req GetFastCommandListReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if !common.CheckRobotKey(req.RobotKey) {
		common.FmtError(c, `param_invalid`, `robot_key`)
		return
	}
	//data check
	robot, err := common.GetRobotInfo(req.RobotKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(robot) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	data, err := msql.Model(define.TableFastCommand, define.Postgres).Where("robot_id", cast.ToString(robot["id"])).Order("sort asc,id desc").Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}

type GetFastCommandInfoReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func GetFastCommandInfo(c *gin.Context) {
	var (
		req GetFastCommandInfoReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	m := msql.Model(define.TableFastCommand, define.Postgres)
	data, err := m.Where("id", cast.ToString(req.Id)).Where("admin_user_id", cast.ToString(adminUserId)).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}

type DeleteFastCommandReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func DeleteFastCommand(c *gin.Context) {
	var (
		req DeleteFastCommandReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	m := msql.Model(define.TableFastCommand, define.Postgres)
	data, err := m.Where("id", cast.ToString(req.Id)).Where("admin_user_id", cast.ToString(adminUserId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}

type UpdateFastCommandReq struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

func UpdateFastCommandSwitch(c *gin.Context) {
	var (
		req UpdateFastCommandReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	//format check
	robotKey := req.RobotKey
	if !common.CheckRobotKey(robotKey) {
		common.FmtError(c, `param_invalid`, `robot_key`)
		return
	}
	//data check
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
	}
	if len(robot) == 0 {
		common.FmtError(c, `no_data`)
	}
	m := msql.Model("chat_ai_robot", define.Postgres)
	sqlRaw := fmt.Sprintf("fast_command_switch = 1-fast_command_switch")
	data, err := m.Where("id", cast.ToString(robot["id"])).Update2(sqlRaw)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	common.FmtOk(c, data)
}

type (
	SortFastCommandReq struct {
		Data []SortFastCommandItem `json:"data"`
	}
	SortFastCommandItem struct {
		Id   int `form:"id" json:"id" binding:"required"`
		Sort int `form:"sort" json:"sort" binding:"required"`
	}
)

func SortFastCommand(c *gin.Context) {
	var (
		req SortFastCommandReq
		err error
	)
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	m := msql.Model(define.TableFastCommand, define.Postgres)
	m.Begin()
	count := 0
	for _, v := range req.Data {
		_, err = m.Where("id", cast.ToString(v.Id)).Where("admin_user_id", cast.ToString(adminUserId)).Update(msql.Datas{
			"sort":        v.Sort,
			"update_time": time.Now().Unix(),
		})
		if err != nil {
			logs.Error(err.Error())
			m.Rollback()
			common.FmtError(c, `sys_err`)
		}
		count++
	}
	m.Commit()
	common.FmtOk(c, count)
}
