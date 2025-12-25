// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func SaveTriggerOfficialConfig(adminUserId string, trigger TriggerConfig, robot msql.Params, lang string) error {
	if !tool.InArray(trigger.TriggerOfficialConfig.MsgType, []string{define.TriggerOfficialMessage, define.TriggerOfficialSubscribeUnScribe,
		define.TriggerOfficialMenuClick, define.TriggerOfficialQrCodeScan}) {
		return errors.New(i18n.Show(lang, `param_err`, `trigger_json`))
	}
	appidInfos, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, adminUserId).Where(`app_type`, `official_account`).Field(`app_id`).Select()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	var appIds = make([]string, 0)
	for _, appInfo := range appidInfos {
		appIds = append(appIds, cast.ToString(appInfo[`app_id`]))
	}
	chooseAppIds := strings.Split(trigger.TriggerOfficialConfig.AppIds, `,`)
	if len(chooseAppIds) == 0 {
		return errors.New(i18n.Show(lang, `param_err`, `appid`))
	}
	for _, appid := range chooseAppIds {
		if !tool.InArray(appid, appIds) {
			return errors.New(i18n.Show(lang, `param_err`, `appid`))
		}
	}
	for _, appid := range chooseAppIds {
		_, err = msql.Model(`work_flow_trigger`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: cast.ToInt(robot[`admin_user_id`]),
			`robot_id`:      cast.ToString(robot[`id`]),
			`trigger_type`:  TriggerTypeOfficial,
			`trigger_json`:  tool.JsonEncodeNoError(trigger),
			`find_key`:      fmt.Sprintf(`%s.%s`, appid, trigger.TriggerOfficialConfig.MsgType),
			`create_time`:   time.Now().Unix(),
			`update_time`:   time.Now().Unix(),
		}, `id`)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		lib_redis.DelCacheData(define.Redis, common.TriggerOfficialCacheBuildHandler{
			AppId:   appid,
			MsgType: trigger.TriggerOfficialConfig.MsgType,
		})
	}
	return nil
}

func TriggerOfficialVerifyStartNode(robotId, msgType, appid string, trigger TriggerConfig) (isOk bool, robot msql.Params) {
	//base verify
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, robotId).Field(`admin_user_id,robot_key,robot_name`).Find()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	if len(robotInfo) == 0 {
		logs.Debug(LogTriggerPrefix + `robot not exist`)
		return
	}
	robot, err = common.GetRobotInfo(robotInfo[`robot_key`])
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	if len(robot) == 0 {
		logs.Debug(LogTriggerPrefix + `robot not exist`)
		return
	}
	//trigger plugin switch
	nodeParam, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`robot_id`, robotId).Where(`node_type`, cast.ToString(NodeTypeStart)).Value(`node_params`)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	nodeParams := NodeParams{}
	err = tool.JsonDecode(nodeParam, &nodeParams)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	if len(nodeParams.Start.TriggerList) == 0 {
		logs.Debug(LogTriggerPrefix + `triggers is empty`)
		return
	}

	for _, triggerVal := range nodeParams.Start.TriggerList {
		if triggerVal.TriggerSwitch == false {
			continue
		}
		if triggerVal.TriggerType != trigger.TriggerType {
			continue
		}
		//official trigger
		if triggerVal.TriggerOfficialConfig.MsgType != msgType {
			continue
		}
		if !tool.InArray(appid, strings.Split(triggerVal.TriggerOfficialConfig.AppIds, `,`)) {
			continue
		}
		isOk = true
	}
	return
}

func StartOfficial(message map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error(LogTriggerPrefix+`panic recover %v`, r)
		}
	}()
	appid := strings.TrimSpace(cast.ToString(message[`appid`]))
	//check app
	appInfo, err := common.GetWechatAppInfo(`app_id`, appid)
	if err != nil {
		logs.Error(`get app info error :%s`, err.Error())
		return
	}
	if len(appInfo) == 0 {
		return
	}
	if appInfo[`app_type`] != lib_define.AppOfficeAccount {
		return
	}
	openid := strings.TrimSpace(cast.ToString(message[`FromUserName`]))
	msgType := strings.ToLower(strings.TrimSpace(cast.ToString(message[`MsgType`])))
	event := strings.ToLower(strings.TrimSpace(cast.ToString(message[`Event`])))
	if tool.InArray(msgType, []string{lib_define.MsgTypeText, lib_define.MsgTypeImage}) {
		msgType = define.TriggerOfficialMessage
	} else if msgType == lib_define.MsgTypeEvent {
		if event == lib_define.EventSubscribe {
			if cast.ToString(message[`EventKey`]) != `` {
				msgType = define.TriggerOfficialQrCodeScan
			} else {
				msgType = define.TriggerOfficialSubscribeUnScribe
			}
		} else if event == lib_define.EventUnSubscribe {
			msgType = define.TriggerOfficialSubscribeUnScribe
		} else if tool.InArray(event, []string{lib_define.EventView, lib_define.EventClick, lib_define.EventViewMiniprogram}) {
			msgType = define.TriggerOfficialMenuClick
		} else if event == lib_define.EventScan {
			msgType = define.TriggerOfficialQrCodeScan
		}
	}
	if msgType == `` {
		return
	}

	workflowTriggers := make([]msql.Params, 0)
	err = lib_redis.GetCacheWithBuild(define.Redis, common.TriggerOfficialCacheBuildHandler{
		AppId:   appid,
		MsgType: msgType,
	}, &workflowTriggers, time.Hour*24)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	for _, trigger := range workflowTriggers {
		triggerConfig := TriggerConfig{}
		err := tool.JsonDecode(trigger[`trigger_json`], &triggerConfig)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			continue
		}
		isOk, robot := TriggerOfficialVerifyStartNode(trigger[`robot_id`], msgType, appid, triggerConfig)
		if !isOk {
			logs.Error(LogTriggerPrefix + `触发器验证失败`)
			continue
		}
		go func(outputs []TriggerOutputParam) {
			testParams := make(map[string]any)
			for _, output := range outputs {
				if _, ok := message[output.Key]; ok {
					testParams[output.Key] = message[output.Key]
				} else {
					testParams[output.Key] = ``
				}
			}
			startTime := time.Now().Unix()
			workFlowParams := &WorkFlowParams{
				ChatRequestParam: &define.ChatRequestParam{
					ChatBaseParam: &define.ChatBaseParam{
						Openid:      openid,
						AdminUserId: cast.ToInt(robot[`admin_user_id`]),
						Robot:       robot,
					},
				},
				TriggerParams: TriggerParams{
					TriggerType: TriggerTypeOfficial,
					TestParams:  testParams,
				},
			}
			_, _, err := BaseCallWorkFlow(workFlowParams)
			setRunResult(trigger[`id`], startTime, err)
		}(triggerConfig.Outputs)
	}

	return
}
