// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/app/plugin/php"
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

func GetTriggerConfigList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	list, err := work_flow.GetTriggerConfigList(adminUserId, c.Query(`robot_key`), common.GetLang(c))
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

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
			`node_name`:     lib_define.TheProcessBegins,
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

	pluginName := []string{}
	for _, node := range list {
		nodeType := cast.ToInt(node[`node_type`])
		if nodeType == work_flow.NodeTypePlugin { //获取引用插件列表
			nodeParams := work_flow.DisposeNodeParams(nodeType, node[`node_params`], common.GetLang(c))
			pluginName = append(pluginName, nodeParams.Plugin.Name)
		}
	}

	remotePlugin := []define.RemotePluginStruct{}
	remotePluginMap := make(map[string]define.RemotePluginStruct)
	localPluginMap := make(map[string]string)

	//如果存在插件引用，查询一下插件状态
	if len(pluginName) > 0 {
		reqData := map[string]any{
			"name": strings.Join(pluginName, ","),
		}
		// 获取远程插件列表
		resp, _ := requestXiaokefu(`kf/ChatWiki/CommonGetPluginList`, reqData)
		err := tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(resp.Data), &remotePlugin)
		if err == nil {
			for _, pluginStruct := range remotePlugin {
				remotePluginMap[pluginStruct.Name] = pluginStruct
			}
		}

		//查询本地插件信息
		pluginList, _ := msql.Model(`plugin_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).Where(`name`, `in`, strings.Join(pluginName, `,`)).Select()
		for _, params := range pluginList {
			localPluginMap[params[`name`]] = params[`has_loaded`]
		}
	}

	filtered := make([]msql.Params, 0)
	for i, node := range list {
		nodeType := cast.ToInt(node[`node_type`])
		if !tool.InArrayInt(nodeType, work_flow.NodeTypes[:]) {
			logs.Error(`invalid node type:%s`, node[`node_type`])
			continue
		}
		nodeParams := work_flow.DisposeNodeParams(cast.ToInt(node[`node_type`]), node[`node_params`], common.GetLang(c))
		if nodeType == work_flow.NodeTypePlugin { //获取引用插件列表
			//logs.Debug(`开始处理插件脱敏:%s`, nodeParams.Plugin.Name)
			manifest, _ := php.GetPluginManifest(nodeParams.Plugin.Name)
			//logs.Debug(`判断前:%s`, nodeParams.Plugin.Name)
			if nodeParams.Plugin.NoAuthFilter {
				//敏感信息脱敏 兼容没安装插件插件的后续脱敏
				//logs.Debug(`判断后:%s`, nodeParams.Plugin.Name)
				nodeInfoStr := node[`node_info_json`]
				nodeParams.Plugin, nodeInfoStr = work_flow.PluginNodeDesensitize(userId, nodeParams.Plugin, nodeInfoStr)
				list[i][`node_info_json`] = nodeInfoStr
			}
			//logs.Debug(`处理完成:%s`, nodeParams.Plugin.Name)

			list[i][`has_loaded`] = localPluginMap[nodeParams.Plugin.Name]
			list[i][`plugin_name`] = nodeParams.Plugin.Name
			if manifest != nil {
				list[i][`plugin_version`] = manifest.Version
			} else {
				list[i][`plugin_version`] = `0.0.0`
			}
			list[i][`remote_plugin_version`] = remotePluginMap[nodeParams.Plugin.Name].LatestVersion
			list[i][`latest_version_detail_url`] = remotePluginMap[nodeParams.Plugin.Name].LatestVersionDetail.DownloadUrl
			list[i][`latest_version_detail_id`] = remotePluginMap[nodeParams.Plugin.Name].LatestVersionDetail.Id
		}
		if nodeType == work_flow.NodeTypeHttpTool {
			nodeParams := work_flow.DisposeNodeParams(nodeType, node[`node_params`], common.GetLang(c))
			list[i][`node_icon`] = nodeParams.Curl.ToolInfo.HttpToolAvatar
		}

		delete(list[i], `admin_user_id`)
		delete(list[i], `robot_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		list[i][`node_params`] = tool.JsonEncodeNoError(nodeParams)
		filtered = append(filtered, list[i])
	}
	c.String(http.StatusOK, lib_web.FmtJson(filtered, nil))
}

func GetStartNode(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(robotId)).
		Find()
	if len(robot) == 0 || cast.ToInt(robot[`application_type`]) != define.ApplicationTypeFlow {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	m := msql.Model(`work_flow_node`, define.Postgres)
	info, err := m.Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, robot[`id`]).
		Where(`node_key`, robot[`start_node_key`]).
		Where(`data_type`, cast.ToString(define.DataTypeRelease)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(info, nil))
}

// GetHttpAuthConfig get http node auth config
func GetHttpAuthConfig(c *gin.Context) {
	var userId = getLoginUserId(c)
	if userId == 0 {
		return
	}
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	m := msql.Model(`http_node_auth_config`, define.Postgres)
	list, err := m.Where(`staff_user_id`, cast.ToString(userId)).Where(`admin_user_id`, cast.ToString(adminUserId)).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))

}

// SaveHttpAuthConfig save http node auth config
func SaveHttpAuthConfig(c *gin.Context) {
	var userId = getLoginUserId(c)
	if userId == 0 {
		return
	}
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	// http_auth_config_list 是前端通过post form传过来的一个json数组，每个元素包含以下字段：auth_key、auth_value、auth_value_addto、auth_remark
	httpAuthConfigList := c.PostForm(`http_auth_config_list`)
	var authConfigList []map[string]interface{}
	tool.JsonDecodeUseNumber(httpAuthConfigList, &authConfigList)
	if len(authConfigList) > 100 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `the auth config should less then 100`))))
		return
	}

	// 首先将http_node_auth_config 表中staff_user_id、admin_user_id 的数据清空，然后循环插入httpAuthConfigList 中的数据
	m := msql.Model(`http_node_auth_config`, define.Postgres)
	_, err := m.Where(`staff_user_id`, cast.ToString(userId)).Where(`admin_user_id`, cast.ToString(adminUserId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	if len(authConfigList) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
		return
	}

	// 循环插入httpAuthConfigList 中的数据
	currentTime := tool.Time2Int()
	for _, item := range authConfigList {
		authKey := strings.TrimSpace(cast.ToString(item[`auth_key`]))
		if authKey == "" {
			continue
		}
		authRemark := strings.TrimSpace(cast.ToString(item[`auth_remark`]))
		// 如果auth_remark 大于10个字则跳过
		if len(authRemark) > 10 {
			continue
		}

		_, err = m.Insert(msql.Datas{
			`staff_user_id`:    userId,
			`admin_user_id`:    adminUserId,
			`auth_key`:         authKey,
			`auth_value`:       strings.TrimSpace(cast.ToString(item[`auth_value`])),
			`auth_value_addto`: strings.TrimSpace(cast.ToString(item[`auth_value_addto`])),
			`auth_remark`:      authRemark,
			`create_time`:      currentTime,
			`update_time`:      currentTime,
		})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
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
	uniIdentifier := strings.TrimSpace(c.PostForm(`uni_identifier`))
	realUserAgent := strings.TrimSpace(c.PostForm(`user_agent`))
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
	lockEditKey := fmt.Sprintf("user_id:%s,robot:%s,", userId, robotKey)
	lockKeyMd5 := define.LockPreKey + ".draft_lock." + tool.MD5(lockEditKey)

	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["uni_identifier"] = uniIdentifier
	lockValue, _ := tool.JsonEncode(filtered)

	lockRes, err := define.Redis.Get(context.Background(), lockKeyMd5).Result()
	if lockRes != "" && lockRes != lockValue { //有编辑锁，且锁不是自己的
		redisValueMap := make(map[string]any)
		_ = tool.JsonDecodeUseNumber(lockRes, &redisValueMap)

		userInfo, err := GetUserInfo(cast.ToString(redisValueMap["login_user_id"]))
		staffUserName := lib_define.UnknownUser
		if err == nil && userInfo["user_name"] != "" {
			staffUserName = userInfo["user_name"]
		}

		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_draft_edit_permission`, staffUserName, redisValueMap["remote_addr"]))))
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
	var questionMultipleSwitch bool
	clearNodeKeys := make([]string, 0)
	m := msql.Model(`work_flow_node`, define.Postgres)
	if dataType == define.DataTypeRelease {
		if startNodeKey, modelConfigIds, libraryIds, questionMultipleSwitch, err = work_flow.VerifyWorkFlowNodes(nodeList, userId, common.GetLang(c)); err != nil {
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
			`question_multiple_switch`:   cast.ToUint(questionMultipleSwitch),
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
	userInfo, err := GetUserInfo(cast.ToString(getLoginUserId(c)))
	editUserName := lib_define.UnknownUser
	if err == nil && userInfo["user_name"] != "" {
		editUserName = userInfo["user_name"]
	}
	_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Update(msql.Datas{
		`last_edit_ip`:         lib_web.GetClientIP(c),
		`last_edit_user_agent`: fmt.Sprintf(`editor: %s, %s`, editUserName, realUserAgent),
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

func VerifyTriggerCronConfig(adminUserId string, node *work_flow.WorkFlowNode, lang string) error {
	for _, trigger := range node.NodeParams.Start.TriggerList {
		var err error
		switch trigger.TriggerType {
		case work_flow.TriggerTypeCron:
			err = verifyTriggerCronConfig(trigger.TriggerCronConfig, lang)
		case work_flow.TriggerTypeOfficial:
			err = verifyTriggerOfficialConfig(adminUserId, trigger.TriggerOfficialConfig, lang)
		case work_flow.TriggerTypeWebHook:
			err = verifyTriggerWebHookConfig(trigger.TriggerWebHookConfig, lang)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyTriggerWebHookConfig(triggerWebHookConfig work_flow.TriggerWebHookConfig, lang string) error {
	if triggerWebHookConfig.Url == `` {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.url`))
	}
	if !tool.InArray(triggerWebHookConfig.Method, []string{http.MethodGet, http.MethodPost}) {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.method`))
	}
	if !tool.InArray(triggerWebHookConfig.RequestContentType, []string{define.ContentTypeNode, define.ContentTypeMultipart, define.ContentTypeFormUrl, define.ContentTypeJson}) {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.request_content_type`))
	}
	if triggerWebHookConfig.ResponseType == `` || !tool.InArray(triggerWebHookConfig.ResponseType, []string{work_flow.WebHookResponseTypeNow, work_flow.WebHookResponseTypeMessageVariable}) {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.response_type`))
	}
	if cast.ToInt(triggerWebHookConfig.SwitchAllowIp) > 0 && triggerWebHookConfig.AllowIps == `` {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.allow_ips`))
	}
	if len(strings.Split(triggerWebHookConfig.AllowIps, "\n")) > 1000 {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_web_hook_config.allow_ips`))
	}
	return nil
}

func verifyTriggerCronConfig(cronConfig work_flow.TriggerCronConfig, lang string) error {
	if cronConfig.Type == work_flow.CronTypeSelectTime {
		if len(strings.Split(cronConfig.HourMinute, `:`)) != 2 {
			return errors.New(i18n.Show(lang, `param_err`, `trigger_cron_config`))
		}
		if cronConfig.EveryType == work_flow.EveryTypeWeek && tool.InArrayString(cronConfig.WeekNumber, []string{`0`, `1`, `2`, `3`, `4`, `5`, `6`}) {
			return errors.New(i18n.Show(lang, `param_err`, `trigger_cron_config`))
		}
		if cronConfig.EveryType == work_flow.EveryTypeMonth && tool.InArrayString(cronConfig.MonthDay, []string{`1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`, `22`, `23`, `24`, `25`, `26`, `27`, `28`, `29`, `30`, `31`}) {
			return errors.New(i18n.Show(lang, `param_err`, `trigger_cron_config`))
		}
	} else if cronConfig.Type == work_flow.CronTypeCrontab {
		if cronConfig.LinuxCrontab == `` || !common.CheckLinuxCrontab(cronConfig.LinuxCrontab) {
			return errors.New(i18n.Show(lang, `linux_crontab_err`))
		}
	}
	return nil
}

func verifyTriggerOfficialConfig(adminUserId string, officialConfig work_flow.TriggerOfficialConfig, lang string) error {
	appidList, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, adminUserId).Where(`app_type`, `official_account`).ColumnArr(`app_id`)
	if err != nil {
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	chooseAppidList := strings.Split(officialConfig.AppIds, `,`)
	for _, appid := range chooseAppidList {
		if !tool.InArray(appid, appidList) {
			return errors.New(i18n.Show(lang, `param_err`, `trigger_official_config`))
		}
	}
	if !tool.InArray(officialConfig.MsgType, []string{
		define.TriggerOfficialMessage,
		define.TriggerOfficialSubscribeUnScribe,
		define.TriggerOfficialMenuClick,
		define.TriggerOfficialQrCodeScan,
	}) {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_official_config`))
	}
	return nil
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
			logs.Error(`invalid node type:%s`, node[`node_type`])
			continue
		}
		list[i][`data_type`] = cast.ToString(define.DataTypeDraft)
		delete(list[i], `admin_user_id`)
		delete(list[i], `robot_id`)
		delete(list[i], `create_time`)
		delete(list[i], `update_time`)
		list[i][`node_params`] = tool.JsonEncodeNoError(work_flow.DisposeNodeParams(cast.ToInt(node[`node_type`]), node[`node_params`], common.GetLang(c)))
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
	uniIdentifier := strings.TrimSpace(c.PostForm(`uni_identifier`))
	realUserAgent := strings.TrimSpace(c.PostForm(`user_agent`))
	if utf8.RuneCountInString(versionDesc) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `desc`))))
		return
	}

	//获取编辑锁内容
	lockEditKey := fmt.Sprintf("user_id:%s,robot:%s,", userId, robotKey)
	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["uni_identifier"] = uniIdentifier
	lockValue, _ := tool.JsonEncode(filtered)

	lockKeyMd5 := define.LockPreKey + ".draft_lock." + tool.MD5(lockEditKey)

	lockRes, err := define.Redis.Get(context.Background(), lockKeyMd5).Result()
	if lockRes != "" && lockRes != lockValue { //没有编辑锁
		redisValueMap := make(map[string]any)
		_ = tool.JsonDecodeUseNumber(lockRes, &redisValueMap)

		userInfo, err := GetUserInfo(cast.ToString(redisValueMap["login_user_id"]))
		staffUserName := lib_define.UnknownUser
		if err == nil && userInfo["user_name"] != "" {
			staffUserName = userInfo["user_name"]
		}

		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_draft_edit_permission`, staffUserName, redisValueMap["remote_addr"]))))
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
	var questionMultipleSwitch bool
	clearNodeKeys := make([]string, 0)
	if startNodeKey, modelConfigIds, libraryIds, questionMultipleSwitch, err = work_flow.VerifyWorkFlowNodes(nodeList, userId, common.GetLang(c)); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	//trigger verify
	for _, node := range nodeList {
		if node.NodeType != work_flow.NodeTypeStart {
			continue
		}
		err := VerifyTriggerCronConfig(cast.ToString(userId), &node, common.GetLang(c))
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			lib_redis.UnLock(define.Redis, lockKey) //unlock
			return
		}
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
	startNode := work_flow.WorkFlowNode{}
	m := msql.Model(``, define.Postgres)
	_ = m.Begin()
	userInfo, err := GetUserInfo(cast.ToString(getLoginUserId(c)))
	editUserName := lib_define.UnknownUser
	if err == nil && userInfo["user_name"] != "" {
		editUserName = userInfo["user_name"]
	}

	m.Reset()
	versionId, versionErr := m.Table(`work_flow_version`).Insert(map[string]any{
		"robot_id":             robotId,
		"admin_user_id":        userId,
		"version":              version,
		"version_desc":         versionDesc,
		"create_time":          time.Now().Unix(),
		"update_time":          time.Now().Unix(),
		`last_edit_ip`:         lib_web.GetClientIP(c),
		`last_edit_user_agent`: fmt.Sprintf(`editor: %s, %s`, editUserName, realUserAgent),
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
		if node.NodeType == work_flow.NodeTypeStart {
			startNode = node
		}
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
		`question_multiple_switch`:   cast.ToUint(questionMultipleSwitch),
		`last_edit_ip`:               lib_web.GetClientIP(c),
		`last_edit_user_agent`:       fmt.Sprintf(`editor: %s, %s`, editUserName, realUserAgent),
	})
	if err != nil {
		_ = m.Rollback()
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
	_ = m.Commit()
	//trigger save
	err = work_flow.SaveTriggerConfig(robot, &startNode, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		lib_redis.UnLock(define.Redis, lockKey) //unlock
		return
	}
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
	uniIdentifier := strings.TrimSpace(c.Query(`uni_identifier`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}

	lockKey := fmt.Sprintf("user_id:%s,robot:%s,", userId, robotKey)
	lockKeyMd5 := define.LockPreKey + ".draft_lock." + tool.MD5(lockKey)

	filtered := make(map[string]any)
	filtered["login_user_id"] = getLoginUserId(c)
	filtered["robot_key"] = robotKey
	filtered["remote_addr"] = lib_web.GetClientIP(c)
	filtered["uni_identifier"] = uniIdentifier
	lockValue, _ := tool.JsonEncode(filtered)

	corpConfig := common.GetAdminConfig(userId)
	draft_exptime := cast.ToInt(corpConfig["draft_exptime"])
	lockRes, lockVal, lockTtl := lib_redis.AddValueLock(define.Redis, lockKeyMd5, lockValue, time.Duration(draft_exptime)*time.Minute)

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
			filtered["login_user_id"] = lockValueMap["login_user_id"]
		}
	}

	//查一下当前草稿最后编辑人是谁
	flowInfo, _ := msql.Model(`chat_ai_robot`, define.Postgres).Where(`robot_key`, robotKey).Find()
	filtered["user_agent"] = flowInfo["last_edit_user_agent"]

	loginUser, err := GetUserInfo(cast.ToString(filtered["login_user_id"]))
	if err == nil {
		filtered["login_user_name"] = loginUser["user_name"]
	}

	c.String(http.StatusOK, lib_web.FmtJson(filtered, nil))

}

func TriggerConfigList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	title := strings.TrimSpace(c.Query(`title`))
	triggerList, err := work_flow.TriggerList(adminUserId, common.GetLang(c))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	result := make([]map[string]any, 0)
	for _, trigger := range triggerList {
		manifest := map[string]any{
			`name`:       `trigger_` + trigger[`trigger_type`],
			`version`:    `1.0.0`,
			`type`:       `trigger`,
			`compatible`: ``,
			`resource`:   ``,
		}
		if cast.ToInt(trigger[`switch_status`]) > 0 {
			manifest[`has_loaded`] = true
		} else {
			manifest[`has_loaded`] = false
		}
		boolShow := true
		if title != `` {
			boolShow = false
			if strings.Contains(trigger[`name`], title) {
				boolShow = true
			}
			if strings.Contains(lib_define.Trigger, title) {
				boolShow = true
			}
			if strings.Contains(trigger[`author`], title) {
				boolShow = true
			}
			if strings.Contains(trigger[`intro`], title) {
				boolShow = true
			}
		}
		if boolShow {
			result = append(result, map[string]any{
				`local`: manifest,
				`remote`: map[string]any{
					`title`:             trigger[`name`],
					`author`:            trigger[`author`],
					`icon`:              trigger[`icon`],
					`latest_version`:    `1.0.0`,
					`filter_type_title`: lib_define.Trigger,
					`description`:       trigger[`intro`],
					`trigger_config_id`: trigger[`id`],
					`trigger_type`:      trigger[`trigger_type`],
				},
			})
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
	return
}
