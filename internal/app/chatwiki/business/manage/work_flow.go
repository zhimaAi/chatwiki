// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_redis"
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
	"github.com/zhimaAi/go_tools/tool"
)

func GetNodeList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//format check
	robotKey := strings.TrimSpace(c.Query(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	//data check
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeFlow {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	dataType := cast.ToUint(c.DefaultQuery(`data_type`, cast.ToString(define.DataTypeDraft)))
	if !tool.InArray(dataType, []uint{define.DataTypeDraft, define.DataTypeRelease}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `data_type`))))
		return
	}
	m := msql.Model(`work_flow_node`, define.Postgres)
	list, err := m.Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, robot[`id`]).
		Where(`data_type`, cast.ToString(dataType)).Order(`id`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(list) == 0 && dataType == define.DataTypeDraft {
		_, err = m.Insert(msql.Datas{
			`admin_user_id`: userId,
			`data_type`:     dataType,
			`robot_id`:      robot[`id`],
			`node_type`:     work_flow.NodeTypeStart,
			`node_name`:     `流程开始`,
			`node_key`:      tool.MD5(time.Now().String() + tool.Random(6)),
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		GetNodeList(c)
		return
	}
	for i, node := range list {
		delete(list[i], `admin_user_id`)
		delete(list[i], `robot_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		list[i][`node_params`] = tool.JsonEncodeNoError(work_flow.DisposeNodeParams(cast.ToInt(node[`node_type`]), node[`node_params`]))
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveNodes(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//format check
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	//data check
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeFlow {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	dataType := cast.ToUint(c.DefaultPostForm(`data_type`, cast.ToString(define.DataTypeDraft)))
	if !tool.InArray(dataType, []uint{define.DataTypeDraft, define.DataTypeRelease}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `data_type`))))
		return
	}
	nodeList := make([]work_flow.WorkFlowNode, 0)
	if err := tool.JsonDecodeUseNumber(c.PostForm(`node_list`), &nodeList); err != nil || len(nodeList) == 0 {
		if err != nil && define.IsDev {
			logs.Error(`node_list err:%s`, err.Error())
		}
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `node_list`))))
		return
	}
	lockKey := define.LockPreKey + fmt.Sprintf(`%s.%s.%d`, `SaveNodes`, robotKey, dataType)
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	var startNodeKey, modelConfigIds, libraryIds string
	clearNodeKeys := make([]string, 0)
	m := msql.Model(`work_flow_node`, define.Postgres)
	if dataType == define.DataTypeRelease {
		if startNodeKey, modelConfigIds, libraryIds, err = work_flow.VerifyWorkFlowNodes(nodeList, userId); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			lib_redis.UnLock(define.Redis, lockKey) //unlock
			return
		}
		nodeKeys, err := m.Where(`robot_id`, robot[`id`]).Where(`data_type`, cast.ToString(dataType)).ColumnArr(`node_key`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			lib_redis.UnLock(define.Redis, lockKey) //unlock
			return
		}
		clearNodeKeys = append(clearNodeKeys, nodeKeys...)
	}
	//database dispose
	_ = m.Begin()
	_, err = m.Where(`robot_id`, robot[`id`]).Where(`data_type`, cast.ToString(dataType)).Delete()
	if err != nil {
		_ = m.Rollback()
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	for _, node := range nodeList {
		_, err = m.Insert(msql.Datas{
			`admin_user_id`:  userId,
			`data_type`:      dataType,
			`robot_id`:       robot[`id`],
			`node_type`:      node.NodeType,
			`node_name`:      node.NodeName,
			`node_key`:       node.NodeKey,
			`node_params`:    tool.JsonEncodeNoError(node.NodeParams),
			`node_info_json`: tool.JsonEncodeNoError(node.NodeInfoJson),
			`next_node_key`:  node.NextNodeKey,
			`create_time`:    tool.Time2Int(),
			`update_time`:    tool.Time2Int(),
		})
		if err != nil {
			_ = m.Rollback()
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			lib_redis.UnLock(define.Redis, lockKey) //unlock
			return
		}
		if dataType == define.DataTypeRelease {
			clearNodeKeys = append(clearNodeKeys, node.NodeKey)
		}
	}
	if dataType == define.DataTypeRelease {
		_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Update(msql.Datas{
			`start_node_key`:             startNodeKey,
			`library_ids`:                libraryIds,
			`work_flow_model_config_ids`: fmt.Sprintf(`{%s}`, modelConfigIds),
			`update_time`:                tool.Time2Int(),
		})
		if err != nil {
			_ = m.Rollback()
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			lib_redis.UnLock(define.Redis, lockKey) //unlock
			return
		}
	}
	_ = m.Commit()
	//clear cached data
	for _, nodeKey := range clearNodeKeys {
		lib_redis.DelCacheData(define.Redis, &common.NodeCacheBuildHandler{RobotId: cast.ToUint(robot[`id`]), DataType: define.DataTypeRelease, NodeKey: nodeKey})
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	lib_redis.UnLock(define.Redis, lockKey) //unlock
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func deleteWorkFlowByRobotId(robotId int) error {
	_, err := msql.Model(`work_flow_node`, define.Postgres).Where(`robot_id`, cast.ToString(robotId)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`work_flow_logs`, define.Postgres).Where(`robot_id`, cast.ToString(robotId)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}
