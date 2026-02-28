// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"net/http"
	"strings"
	"time"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func GetWorkbenchConfig(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	config, err := msql.Model(`workbench_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(config) == 0 {
		_, err := msql.Model(`workbench_config`, define.Postgres).Insert(msql.Datas{
			"admin_user_id":         adminUserId,
			"default_robot_id":      0,
			"enable_last_app_entry": 1,
			"create_time":           time.Now().Unix(),
			"update_time":           time.Now().Unix(),
		})
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		config, err = msql.Model(`workbench_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}
	config[`default_robot_key`] = ``
	defaultRobotId := cast.ToInt(config[`default_robot_id`])
	if defaultRobotId > 0 {
		robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(defaultRobotId)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(robot) > 0 {
			config[`default_robot_key`] = robot[`robot_key`]
		}
	}
	home := getHomeConfig(adminUserId, userId, config)
	common.FmtOk(c, map[string]any{
		`config`: config,
		`home`:   home,
	})
}

type WorkbenchConfigReq struct {
	DefaultRobotID     int `form:"default_robot_id" json:"default_robot_id"`
	EnableLastAppEntry int `form:"enable_last_app_entry" json:"enable_last_app_entry"` // 1open, 0close
}

func SaveWorkbenchConfig(c *gin.Context) {
	var req WorkbenchConfigReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	if req.DefaultRobotID > 0 {
		robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(req.DefaultRobotID)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(robot) == 0 {
			common.FmtError(c, `robot_info_not_exist`)
			return
		}
	}

	now := int(time.Now().Unix())
	data := msql.Datas{
		"default_robot_id":      req.DefaultRobotID,
		"enable_last_app_entry": req.EnableLastAppEntry,
		"update_time":           now,
	}

	existingConfig, err := msql.Model(`workbench_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	var result interface{}
	if len(existingConfig) > 0 {
		_, err = msql.Model(`workbench_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Update(data)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		result = nil
	} else {
		data["admin_user_id"] = adminUserId
		data["create_time"] = now
		_, err = msql.Model(`workbench_config`, define.Postgres).
			Insert(data)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		result = nil
	}
	common.FmtOk(c, result)
}

func GetRobotHistoryVisit(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	limit := 5
	m := msql.Model(`robot_history_visit`, define.Postgres).
		Alias(`h`).
		Join(`chat_ai_robot r`, `h.robot_id = r.id and r.admin_user_id = h.admin_user_id`, `left`).
		Where(`h.admin_user_id`, cast.ToString(adminUserId)).
		Where(`h.user_id`, cast.ToString(userId)).
		Order(`h.update_time desc`).
		Field(`h.id, h.robot_id, h.create_time, h.update_time, r.robot_name, r.robot_avatar, r.robot_key, r.application_type`).
		Limit(limit)

	list, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	out := make([]msql.Params, 0, len(list))
	for _, item := range list {
		if len(item[`robot_key`]) > 0 {
			out = append(out, item)
		}
	}
	common.FmtOk(c, out)
}

type RecordRobotVisitReq struct {
	RobotID int `form:"robot_id" json:"robot_id" binding:"required"`
}

func RecordRobotVisit(c *gin.Context) {
	var req RecordRobotVisitReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(req.RobotID)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(robot) == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	now := int(time.Now().Unix())
	exist, err := msql.Model(`robot_history_visit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(req.RobotID)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(exist) > 0 {
		_, err = msql.Model(`robot_history_visit`, define.Postgres).
			Where(`id`, cast.ToString(exist[`id`])).
			Update(msql.Datas{`update_time`: now})
	} else {
		_, err = msql.Model(`robot_history_visit`, define.Postgres).
			Insert(msql.Datas{
				`admin_user_id`: adminUserId,
				`user_id`:       userId,
				`robot_id`:      req.RobotID,
				`create_time`:   now,
				`update_time`:   now,
			})
	}
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

type TopRobotReq struct {
	RobotID int `form:"robot_id" json:"robot_id" binding:"required"`
}

func TopRobot(c *gin.Context) {
	var req TopRobotReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(req.RobotID)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(robot) == 0 {
		common.FmtError(c, `robot_not_exist_or_no_permission`)
		return
	}

	now := int(time.Now().Unix())

	existingTop, err := msql.Model(`workbench_top_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(req.RobotID)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	var result interface{}
	if len(existingTop) > 0 {
		_, err = msql.Model(`workbench_top_robot`, define.Postgres).
			Where(`id`, cast.ToString(existingTop["id"])).
			Delete()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		result = map[string]interface{}{"action": "untopped"}
	} else {
		_, err = msql.Model(`workbench_top_robot`, define.Postgres).
			Insert(msql.Datas{
				"admin_user_id": adminUserId,
				"user_id":       userId,
				"robot_id":      req.RobotID,
				"create_time":   now,
				"update_time":   now,
			})
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		result = map[string]interface{}{"action": "topped"}
	}

	common.FmtOk(c, result)
}

func WorkbenchTeamRobotList(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Query chat robots (application_type = 0), all are displayed
	chatRobots, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Alias(`r`).
		Field(`r.id, r.robot_name, r.application_type, r.robot_intro, r.robot_avatar, r.robot_key, r.group_id`).
		Where(`r.admin_user_id`, cast.ToString(adminUserId)).
		Where(`r.application_type`, `0`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// Query published workflow robots (application_type = 1 and data_type = 2 in work_flow_node)
	workflowRobots, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Alias(`r`).
		Field(`r.id, r.robot_name, r.application_type, r.robot_intro, r.robot_avatar, r.robot_key,  r.group_id`).
		Where(`r.admin_user_id`, cast.ToString(adminUserId)).
		Where(`r.application_type`, `1`).
		Where(`EXISTS (SELECT 1 FROM work_flow_node n WHERE n.robot_id = r.id AND n.admin_user_id = r.admin_user_id AND n.data_type = ` + cast.ToString(define.DataTypeRelease) + `)`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// Merge chat robots and published workflow robots
	robotList := make([]msql.Params, 0, len(chatRobots)+len(workflowRobots))
	robotList = append(robotList, chatRobots...)
	robotList = append(robotList, workflowRobots...)

	var topRobotIds []string
	for _, robot := range robotList {

		robotId := cast.ToString(robot["id"])
		topRobotIds = append(topRobotIds, robotId)
	}

	var userTopRobots []msql.Params
	if len(topRobotIds) > 0 {
		userTopRobots, err = msql.Model(`workbench_top_robot`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`user_id`, cast.ToString(userId)).
			Where(`robot_id`, `in`, strings.Join(topRobotIds, `,`)).
			Select()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}

	topRobotMap := make(map[string]string)
	for _, topRobot := range userTopRobots {
		robotId := cast.ToString(topRobot["robot_id"])
		topRobotMap[robotId] = topRobot[`id`]
	}

	for i := range robotList {
		robotId := cast.ToString(robotList[i]["id"])
		robotList[i]["top_id"] = cast.ToString(topRobotMap[robotId])
	}

	common.FmtOk(c, robotList)
}

func getHomeConfig(adminUserId, userId int, config msql.Params) msql.Params {
	homeConfig := msql.Params{
		"robot_id":  "0",
		"robot_key": "",
	}
	defaultRobotID := cast.ToInt(config["default_robot_id"])
	enableLastAppEntry := cast.ToInt(config["enable_last_app_entry"])
	//last visit robot
	if enableLastAppEntry == 1 {
		history, err := msql.Model(`robot_history_visit`, define.Postgres).
			Alias(`h`).
			Join(`chat_ai_robot r`, `h.robot_id = r.id and r.admin_user_id = h.admin_user_id`, `left`).
			Where(`h.admin_user_id`, cast.ToString(adminUserId)).
			Where(`h.user_id`, cast.ToString(userId)).
			Where(`r.robot_key`, `<>`, ``).
			Order(`h.update_time desc`).
			Field(`r.id, r.robot_key,r.robot_name,r.robot_intro,r.robot_avatar`).
			Find()
		if err == nil && len(history) > 0 {
			homeConfig = getHomeComfigParams(history)
			return homeConfig
		}
	}
	//default robot
	if defaultRobotID > 0 {
		robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(defaultRobotID)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Field(`id, robot_key,robot_name,robot_intro,robot_avatar`).
			Find()
		if err == nil && len(robot) > 0 {
			homeConfig = getHomeComfigParams(robot)
			return homeConfig
		}
	}
	//last created robot
	chatRobot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`application_type`, `0`).
		Order(`create_time desc`).
		Field(`id, robot_key, create_time,robot_name,robot_intro,robot_avatar`).
		Find()
	if err != nil {
		logs.Error("getHomeConfig: get chat robot error: %s", err.Error())
	}
	//last publishd workflow robot
	workflowRobot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Alias(`r`).
		Join(`work_flow_node n`, `r.id = n.robot_id and r.admin_user_id = n.admin_user_id`, `inner`).
		Where(`r.admin_user_id`, cast.ToString(adminUserId)).
		Where(`r.application_type`, `1`).
		Where(`n.data_type`, cast.ToString(define.DataTypeRelease)).
		Order(`n.update_time desc`).
		Field(`r.id, r.robot_key, n.update_time,r.robot_name,r.robot_intro,r.robot_avatar`).
		Find()
	if err != nil {
		logs.Error("getHomeConfig: get workflow robot error: %s", err.Error())
	}
	var chatRobotTime, workflowRobotTime int
	if len(chatRobot) > 0 {
		chatRobotTime = cast.ToInt(chatRobot["create_time"])
	}
	if len(workflowRobot) > 0 {
		workflowRobotTime = cast.ToInt(workflowRobot["update_time"])
	}
	if chatRobotTime > 0 && workflowRobotTime > 0 {
		if chatRobotTime >= workflowRobotTime {
			homeConfig = getHomeComfigParams(chatRobot)
		} else {
			homeConfig = getHomeComfigParams(workflowRobot)
		}
	} else if chatRobotTime > 0 {
		homeConfig = getHomeComfigParams(chatRobot)
	} else if workflowRobotTime > 0 {
		homeConfig = getHomeComfigParams(workflowRobot)
	}
	return homeConfig
}

func getHomeComfigParams(data msql.Params) msql.Params {
	return msql.Params{
		"robot_id":     data["id"],
		"robot_key":    data["robot_key"],
		"robot_name":   data["robot_name"],
		"robot_intro":  data["robot_intro"],
		"robot_avatar": data["robot_avatar"],
	}
}
