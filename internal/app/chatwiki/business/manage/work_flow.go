// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

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
	filtered := make([]msql.Params, 0)
	for i, node := range list {
		nodeType := cast.ToInt(node[`node_type`])
		if !tool.InArrayInt(nodeType, work_flow.NodeTypes[:]) {
			logs.Error(`节点类型非法:%s`, node[`node_type`])
			continue
		}
		delete(list[i], `admin_user_id`)
		delete(list[i], `robot_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		list[i][`node_params`] = tool.JsonEncodeNoError(work_flow.DisposeNodeParams(cast.ToInt(node[`node_type`]), node[`node_params`]))
		filtered = append(filtered, list[i])
	}
	c.String(http.StatusOK, lib_web.FmtJson(filtered, nil))
}

func SaveNodes(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//format check
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	draftSaveType := strings.TrimSpace(c.PostForm(`draft_save_type`))
	lastUpdateTime := strings.TrimSpace(c.PostForm(`draft_save_time`))
	reCoverSave := strings.TrimSpace(c.PostForm(`re_cover_save`))
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

	//获取编辑锁内容
	lockEditKey := fmt.Sprintf("draft_lock:user_id:%s,robot:%s,", userId, robotKey)
	lockKeyMd5 := tool.MD5(lockEditKey)

	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["user_agent"] = c.Request.UserAgent()
	lockValue, _ := tool.JsonEncode(filtered)

	lockRes, err := define.Redis.Get(context.Background(), lockKeyMd5).Result()
	if lockRes != "" && lockRes != lockValue { //有编辑锁，且锁不是自己的
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_draft_edit_permission`))))
		return
	}

	//重新设置一下锁的时间
	corpConfig := common.GetAdminConfig(userId)
	draft_exptime := cast.ToInt(corpConfig["draft_exptime"])
	_, _ = define.Redis.SetEX(context.Background(), lockKeyMd5, lockValue, time.Duration(draft_exptime)*time.Minute).Result()

	if cast.ToInt(reCoverSave) == 0 { //需要验证
		//获取一下库内的草稿信息
		flowInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Find()
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		dbDraftSaveTime := cast.ToInt(flowInfo["draft_save_time"])
		//如果库里面记录了草稿保存时间，且 保存时间大于提交上来的最后保存时间，判定为此次保存是落后版本
		if dbDraftSaveTime != 0 && dbDraftSaveTime > cast.ToInt(lastUpdateTime) {
			c.String(http.StatusOK, lib_web.FmtJson(map[string]int{
				"behind_draft": 1,
			}, errors.New(i18n.Show(common.GetLang(c), `behind_draft`))))
			return
		}
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
			`admin_user_id`:   userId,
			`data_type`:       dataType,
			`robot_id`:        robot[`id`],
			`node_type`:       node.NodeType,
			`node_name`:       node.NodeName,
			`node_key`:        node.NodeKey,
			`node_params`:     tool.JsonEncodeNoError(node.NodeParams),
			`node_info_json`:  tool.JsonEncodeNoError(node.NodeInfoJson),
			`next_node_key`:   node.NextNodeKey,
			`loop_parent_key`: node.LoopParentKey,
			`create_time`:     tool.Time2Int(),
			`update_time`:     tool.Time2Int(),
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

	//保存上次编辑人信息
	_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Update(msql.Datas{
		`last_edit_ip`:         lib_web.GetClientIP(c),
		`last_edit_user_agent`: c.Request.UserAgent(),
		`draft_save_type`:      draftSaveType,
		`draft_save_time`:      tool.Time2Int(),
		`update_time`:          tool.Time2Int(),
	})

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

func TestCodeRun(c *gin.Context) {
	language := strings.TrimSpace(c.DefaultPostForm(`language`, `javaScript`))
	MainFunc := strings.TrimSpace(c.DefaultPostForm(`main_func`, `function main(){return {}}`))
	params := strings.TrimSpace(c.DefaultPostForm(`params`, `{}`))
	Params := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(params, &Params); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	code := http.StatusOK
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-c.Request.Context().Done():
			code = 499 //客户端主动断开
			cancel()
		}
	}()
	data := lib_define.CodeRunBody{MainFunc: MainFunc, Params: Params}
	result, err := common.RequestCodeRun(ctx, language, data)
	c.String(code, lib_web.FmtJson(result, err))
}

func WorkFlowNextVersion(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	robot, err := common.CheckRobotKey2(robotKey, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotId := cast.ToInt(robot[`id`])
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	id, err := m.Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).Value(`id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToInt(id) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	version, err := msql.Model(`work_flow_version`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).Order(`id desc`).Limit(1).Value(`version`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	version = common.NextWorkFlowVersion(version)
	if version == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`version`: version, `version_params`: strings.Split(version, `.`)}, nil))
	return
}

func WorkFlowVersions(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	robot, err := common.CheckRobotKey2(robotKey, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotId := cast.ToInt(robot[`id`])
	versions, err := msql.Model(`work_flow_version`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).Field(`id as version_id,version,version_desc,create_time,update_time,last_edit_ip,last_edit_user_agent`).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`versions`: versions}, nil))
	return
}

func WorkFlowVersionDetail(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	robot, err := common.CheckRobotKey2(robotKey, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotId := cast.ToInt(robot[`id`])
	versionId := cast.ToInt(c.PostForm(`version_id`))
	_, err = common.CheckWorkFlowVersionId(versionId, userId, common.GetLang(c))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	m := msql.Model(`work_flow_node_version`, define.Postgres)
	list, err := m.Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`work_flow_version_id`, cast.ToString(versionId)).Order(`id`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	filtered := make([]msql.Params, 0)
	for i, node := range list {
		nodeType := cast.ToInt(node[`node_type`])
		if !tool.InArrayInt(nodeType, work_flow.NodeTypes[:]) {
			logs.Error(`节点类型非法:%s`, node[`node_type`])
			continue
		}
		list[i][`data_type`] = cast.ToString(define.DataTypeDraft)
		delete(list[i], `admin_user_id`)
		delete(list[i], `robot_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		list[i][`node_params`] = tool.JsonEncodeNoError(work_flow.DisposeNodeParams(cast.ToInt(node[`node_type`]), node[`node_params`]))
		filtered = append(filtered, list[i])
	}
	c.String(http.StatusOK, lib_web.FmtJson(filtered, nil))
	return
}

func WorkFlowPublishVersion(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	version := strings.TrimSpace(c.PostForm(`version`))
	versionDesc := strings.TrimSpace(c.PostForm(`version_desc`))
	if utf8.RuneCountInString(versionDesc) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `desc`))))
		return
	}

	//获取编辑锁内容
	lockEditKey := fmt.Sprintf("draft_lock:user_id:%s,robot:%s,", userId, robotKey)
	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["user_agent"] = c.Request.UserAgent()
	lockValue, _ := tool.JsonEncode(filtered)

	lockKeyMd5 := tool.MD5(lockEditKey)

	lockRes, err := define.Redis.Get(context.Background(), lockKeyMd5).Result()
	if lockRes != lockValue { //没有编辑锁
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_draft_edit_permission`))))
	}

	robot, err := common.CheckRobotKey2(robotKey, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	robotId := cast.ToInt(robot[`id`])
	err = common.CheckWorkFlowVersion(robotId, userId, version, common.GetLang(c))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	dataType := define.DataTypeRelease
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
	if startNodeKey, modelConfigIds, libraryIds, err = work_flow.VerifyWorkFlowNodes(nodeList, userId); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	nodeKeys, err := msql.Model(`work_flow_node`, define.Postgres).Where(`robot_id`, robot[`id`]).Where(`data_type`, cast.ToString(dataType)).ColumnArr(`node_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	clearNodeKeys = append(clearNodeKeys, nodeKeys...)
	//start
	m := msql.Model(``, define.Postgres)
	_ = m.Begin()
	m.Reset()
	versionId, versionErr := m.Table(`work_flow_version`).Insert(map[string]any{
		"robot_id":             robotId,
		"admin_user_id":        userId,
		"version":              version,
		"version_desc":         versionDesc,
		"create_time":          time.Now().Unix(),
		"update_time":          time.Now().Unix(),
		`last_edit_ip`:         lib_web.GetClientIP(c),
		`last_edit_user_agent`: c.Request.UserAgent(),
	}, `id`)
	if versionErr != nil {
		_ = m.Rollback()
		logs.Error(versionErr.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	m.Reset()
	_, err = m.Table(`work_flow_node`).Where(`robot_id`, robot[`id`]).Where(`data_type`, cast.ToString(dataType)).Delete()
	if err != nil {
		_ = m.Rollback()
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	for _, node := range nodeList {
		m.Reset()
		_, err = m.Table(`work_flow_node`).Insert(msql.Datas{
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
		m.Reset()
		_, err = m.Table(`work_flow_node_version`).Insert(msql.Datas{
			`admin_user_id`:        userId,
			`work_flow_version_id`: versionId,
			`robot_id`:             robot[`id`],
			`node_type`:            node.NodeType,
			`node_name`:            node.NodeName,
			`node_key`:             node.NodeKey,
			`node_params`:          tool.JsonEncodeNoError(node.NodeParams),
			`node_info_json`:       tool.JsonEncodeNoError(node.NodeInfoJson),
			`next_node_key`:        node.NextNodeKey,
			`create_time`:          tool.Time2Int(),
			`update_time`:          tool.Time2Int(),
		})
		if err != nil {
			_ = m.Rollback()
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			lib_redis.UnLock(define.Redis, lockKey)
			return
		}
		clearNodeKeys = append(clearNodeKeys, node.NodeKey)
	}
	m.Reset()
	_, err = m.Table(`chat_ai_robot`).Where(`robot_key`, robotKey).Update(msql.Datas{
		`start_node_key`:             startNodeKey,
		`library_ids`:                libraryIds,
		`work_flow_model_config_ids`: fmt.Sprintf(`{%s}`, modelConfigIds),
		`update_time`:                tool.Time2Int(),
		`last_edit_ip`:               lib_web.GetClientIP(c),
		`last_edit_user_agent`:       c.Request.UserAgent(),
	})
	if err != nil {
		_ = m.Rollback()
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
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

// 获取企业配置
func GetAdminConfig(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	config := common.GetAdminConfig(userId)

	c.String(http.StatusOK, lib_web.FmtJson(config, nil))
}

// 设置草稿箱编辑锁过期时间
func SaveDraftExTime(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	draft_exptime := cast.ToInt(strings.TrimSpace(c.PostForm(`draft_exptime`)))

	if draft_exptime < 10 || draft_exptime > 60 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `draft_exptime_exceed_limit`))))
		return
	}

	m := msql.Model(``, define.Postgres)

	list, _ := m.Table(`admin_user_config`).Where(`admin_user_id`, cast.ToString(userId)).Select()

	if len(list) == 0 { //没有配置
		_, err := m.Table(`admin_user_config`).Insert(map[string]any{
			"admin_user_id": userId,
			"create_time":   time.Now().Unix(),
			"update_time":   time.Now().Unix(),
			`draft_exptime`: draft_exptime,
		}, `id`)

		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		_, err := m.Table(`admin_user_config`).Where(`admin_user_id`, cast.ToString(userId)).Update(msql.Datas{
			`draft_exptime`: draft_exptime,
			`update_time`:   tool.Time2Int(),
		})

		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

	list, _ = m.Table(`admin_user_config`).Where(`admin_user_id`, cast.ToString(userId)).Select()
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
	return

}

// 获取草稿箱编辑锁
func GetDraftKey(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	robotKey := strings.TrimSpace(c.Query(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}

	lockKey := fmt.Sprintf("draft_lock:user_id:%s,robot:%s,", userId, robotKey)
	lockKeyMd5 := tool.MD5(lockKey)

	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["user_agent"] = c.Request.UserAgent()
	lockValue, _ := tool.JsonEncode(filtered)

	logs.Error("锁内容：" + lockKey + "，锁结果：" + lockValue)
	corpConfig := common.GetAdminConfig(userId)
	draft_exptime := cast.ToInt(corpConfig["draft_exptime"])
	lockRes, lockVal, lockTtl := lib_redis.AddValueLock(define.Redis, lockKeyMd5, lockValue, time.Duration(draft_exptime)*time.Minute)

	logs.Error("写锁内容：" + lockVal)

	filtered["lock_res"] = lockRes
	filtered["lock_ttl"] = lockTtl
	filtered["lockKey"] = lockKey
	filtered["lockValue"] = lockValue
	filtered["lockVal"] = lockVal
	filtered["Header"] = c.Request.Header
	filtered["is_self"] = lockVal == lockValue

	//获取锁失败，看下是谁锁的
	if !lockRes && lockVal != lockValue {
		lockValueMap := make(map[string]any)
		err := tool.JsonDecode(lockVal, &lockValueMap)
		if err == nil {
			filtered["remote_addr"] = lockValueMap["remote_addr"]
			filtered["user_agent"] = lockValueMap["user_agent"]
			filtered["login_user_id"] = lockValueMap["login_user_id"]
		}
	}

	loginUser, err := msql.Model(define.TableUser, define.Postgres).Where(`id`, cast.ToString(filtered["login_user_id"])).Field("id,user_name").Find()
	if err == nil {
		filtered["login_user_name"] = loginUser["user_name"]
	}

	c.String(http.StatusOK, lib_web.FmtJson(filtered, nil))

}
