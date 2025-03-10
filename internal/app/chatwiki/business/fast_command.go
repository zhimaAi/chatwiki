// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

type GetFastCommandListReq struct {
	RobotKey string `form:"robot_key" json:"robot_key" binding:"required"`
	OpenId   string `form:"open_id" json:"open_id"`
	AppId    int    `form:"app_id,default=-1" json:"app_id,default=-1" binding:"oneof=-1 -2"`
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
	data, err := msql.Model(define.TableFastCommand, define.Postgres).Where("robot_id", cast.ToString(robot["id"])).Where("app_id", cast.ToString(req.AppId)).Order("sort asc,id desc").Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}
