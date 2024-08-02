// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"time"
)

type AddRobotApiKeyReq struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

func AddRobotApikey(c *gin.Context) {
	var req AddRobotApiKeyReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if !common.CheckRobotKey(req.RobotKey) {
		common.FmtError(c, `param_invalid`, `robot_key`)
		return
	}
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
	token := common.GetAuthorizationToken(req.RobotKey)
	if token == "" {
		common.FmtError(c, `sys_err`)
		return
	}
	insertData := msql.Datas{
		"admin_user_id": GetAdminUserId(c),
		"key":           token,
		"robot_key":     req.RobotKey,
		"create_time":   time.Now().Unix(),
		"update_time":   time.Now().Unix(),
	}
	id, err := msql.Model(define.TableChatAiRobotApikey, define.Postgres).Insert(insertData, `id`)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotApiKeyCacheBuildHandler{RobotKey: req.RobotKey})
	common.FmtOk(c, id)
}

type DeleteRobotApikeyReq struct {
	Id       int    `form:"id" json:"id" binding:"required"`
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

func DeleteRobotApikey(c *gin.Context) {
	var req DeleteRobotApikeyReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	id, err := msql.Model(define.TableChatAiRobotApikey, define.Postgres).Where("id", cast.ToString(req.Id)).Where("admin_user_id", cast.ToString(GetAdminUserId(c))).Delete()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotApiKeyCacheBuildHandler{RobotKey: req.RobotKey})
	common.FmtOk(c, id)
}

type ListRobotApikeyReq struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

func ListRobotApikey(c *gin.Context) {
	var req ListRobotApikeyReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if !common.CheckRobotKey(req.RobotKey) {
		common.FmtError(c, `param_invalid`, `robot_key`)
		return
	}
	data, err := common.GetRobotApikeyInfo(req.RobotKey)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	endPoint := define.Config.WebService["api_domain"]
	common.FmtOk(c, map[string]any{
		"list":      data,
		"end_point": endPoint,
	})
}

type UpdateRobotApikeyReq struct {
	Id       int    `form:"id" json:"id" binding:"required"`
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
}

func UpdateRobotApikey(c *gin.Context) {
	var req UpdateRobotApikeyReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	sqlRaw := fmt.Sprintf("status = 1-status,update_time=%v", time.Now().Unix())
	id, err := msql.Model(define.TableChatAiRobotApikey, define.Postgres).Where("id", cast.ToString(req.Id)).Where("admin_user_id", cast.ToString(GetAdminUserId(c))).Update2(sqlRaw)
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotApiKeyCacheBuildHandler{RobotKey: req.RobotKey})
	common.FmtOk(c, id)
}

func addDefaultApiKey(c *gin.Context, robotKey string) {
	token := common.GetAuthorizationToken(robotKey)
	if token == "" {
		return
	}
	insertData := msql.Datas{
		"admin_user_id": GetAdminUserId(c),
		"key":           token,
		"robot_key":     robotKey,
		"create_time":   time.Now().Unix(),
		"update_time":   time.Now().Unix(),
	}
	_, err := msql.Model(define.TableChatAiRobotApikey, define.Postgres).Insert(insertData, `id`)
	if err != nil {
		logs.Error("add default api key err:", err)
		return
	}
}

func deleteRobotApiKey(robotKey string) error {
	_, err := msql.Model(define.TableChatAiRobotApikey, define.Postgres).Where("robot_key", robotKey).Delete()
	if err != nil {
		return err
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotApiKeyCacheBuildHandler{RobotKey: robotKey})
	return nil
}
