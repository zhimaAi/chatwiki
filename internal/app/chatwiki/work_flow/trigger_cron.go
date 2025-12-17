// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var LogTriggerPrefix = `work_flow trigger：`
var cro *cron.Cron

func init() {
	cro = cron.New()
	cro.Start()
}

func SaveTriggerCronConfig(trigger TriggerConfig, robot msql.Params, lang string) error {

	id, err := msql.Model(`work_flow_trigger`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: cast.ToInt(robot[`admin_user_id`]),
		`robot_id`:      cast.ToString(robot[`id`]),
		`trigger_type`:  TriggerTypeCron,
		`trigger_json`:  tool.JsonEncodeNoError(trigger),
		`create_time`:   time.Now().Unix(),
		`update_time`:   time.Now().Unix(),
	}, `id`)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	cronEntryId := 0
	if trigger.TriggerCronConfig.Type == CronTypeCrontab {
		cronEntryId, err = StartCron(cast.ToString(id), trigger, robot, lang)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
		logs.Debug(LogTriggerPrefix+`add entry id %d`, cronEntryId)
		_, err = msql.Model(`work_flow_trigger`, define.Postgres).
			Where(`id`, cast.ToString(id)).Update(msql.Datas{
			`cron_entry_id`: int(cronEntryId),
		})
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	return nil
}

func RemoveEntry(entryId int) {
	cro.Remove(cron.EntryID(entryId))
}

func TriggerCronSelectTimeRun() {
	list, err := msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`trigger_type`, cast.ToString(TriggerTypeCron)).
		Where(`cron_entry_id`, `0`).
		Where(`is_finish`, `0`).Select()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	logs.Debug(LogTriggerPrefix+`trigger select time list %d`, len(list))
	t := time.Now()
	for _, trigger := range list {
		logs.Debug(LogTriggerPrefix+`trigger id %s`, trigger[`id`])
		triggerConfig := TriggerConfig{}
		err = tool.JsonDecode(trigger[`trigger_json`], &triggerConfig)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			continue
		}
		if !triggerConfig.TriggerSwitch {
			logs.Debug(LogTriggerPrefix + ` not switch`)
			continue
		}
		if triggerConfig.TriggerType != TriggerTypeCron {
			logs.Debug(LogTriggerPrefix + ` is not cron`)
			continue
		}
		if triggerConfig.TriggerCronConfig.Type != CronTypeSelectTime {
			logs.Debug(LogTriggerPrefix + ` is not select time`)
			continue
		}
		if triggerConfig.TriggerCronConfig.HourMinute != fmt.Sprintf(`%02d:%02d`, t.Hour(), t.Minute()) {
			logs.Debug(LogTriggerPrefix + ` not time ` + triggerConfig.TriggerCronConfig.HourMinute + ` ` + fmt.Sprintf(`%02d:%02d`, t.Hour(), t.Minute()))
			continue
		}
		if triggerConfig.TriggerCronConfig.EveryType == EveryTypeWeek {
			if cast.ToInt(triggerConfig.TriggerCronConfig.WeekNumber) != cast.ToInt(t.Weekday()) {
				logs.Debug(LogTriggerPrefix + ` not week ` + cast.ToString(cast.ToInt(triggerConfig.TriggerCronConfig.WeekNumber)) + ` ` + cast.ToString(cast.ToInt(t.Weekday())))
				continue
			}
		} else if triggerConfig.TriggerCronConfig.EveryType == EveryTypeMonth {
			if cast.ToInt(triggerConfig.TriggerCronConfig.MonthDay) != cast.ToInt(t.Day()) {
				logs.Debug(LogTriggerPrefix + ` not month day ` + cast.ToString(cast.ToInt(triggerConfig.TriggerCronConfig.MonthDay)) + ` ` + cast.ToString(cast.ToInt(t.Day())))
				continue
			}
		}
		startTime := time.Now().Unix()
		go TriggerCronRun(trigger[`admin_user_id`], trigger[`robot_id`], triggerConfig, func(err error) {
			setRunResult(trigger[`id`], startTime, err)
		})
	}
}

func setRunResult(id string, startTime int64, err error) {
	errmsg := `success`
	if err != nil {
		if utf8.RuneCountInString(err.Error()) > 500 {
			runes := []rune(err.Error())
			errmsg = string(runes[0:500])
		} else {
			errmsg = err.Error()
		}
	}
	_, err = msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`id`, id).Update(msql.Datas{
		`is_finish`: 1,
		`last_msg`: tool.JsonEncodeNoError(map[string]any{
			`start_time`: startTime,
			`end_time`:   time.Now().Unix(),
			`errmsg`:     errmsg,
		}),
		`update_time`: time.Now().Unix(),
	})
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
	}
}

func CheckTriggerSwitchStatus(adminUserId, triggerType string) bool {
	info := msql.Params{}
	err := lib_redis.GetCacheWithBuild(define.Redis, common.TriggerConfigCacheBuildHandler{
		AdminUserId: cast.ToInt(adminUserId),
		TriggerType: triggerType,
	}, &info, time.Hour*12)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return false
	}
	if cast.ToInt(info[`switch_status`]) > 0 {
		return true
	}
	return false
}

func TriggerCronRun(adminUserId, robotId string, trigger TriggerConfig, finishCall func(error)) {
	logs.Debug(LogTriggerPrefix+`trigger run robotId:%s config:%s`, robotId, tool.JsonEncodeNoError(trigger))
	//if !CheckTriggerSwitchStatus(adminUserId, cast.ToString(trigger.TriggerType)) {
	//	logs.Debug(LogTriggerPrefix + `trigger switch status is 0`)
	//	return
	//}
	isOk, robot := TriggerCronVerifyStartNode(robotId, trigger)
	if !isOk {
		logs.Debug(LogTriggerPrefix + `trigger not exist`)
		return
	}
	workFlowParams := &WorkFlowParams{
		ChatRequestParam: &define.ChatRequestParam{
			ChatBaseParam: &define.ChatBaseParam{
				Openid:      "",
				AdminUserId: cast.ToInt(robot[`admin_user_id`]),
				Robot:       robot,
			},
		},
		TriggerParams: TriggerParams{
			TriggerType: TriggerTypeCron,
			TestParams:  map[string]any{},
		},
	}
	_, _, err := BaseCallWorkFlow(workFlowParams)
	finishCall(err)
}

func TriggerCronVerifyStartNode(robotId string, trigger TriggerConfig) (isOk bool, robot msql.Params) {
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
		if triggerVal.TriggerCronConfig.Type != trigger.TriggerCronConfig.Type {
			continue
		}
		if triggerVal.TriggerCronConfig.WeekNumber != trigger.TriggerCronConfig.WeekNumber {
			continue
		}
		if triggerVal.TriggerCronConfig.EveryType != trigger.TriggerCronConfig.EveryType {
			continue
		}
		if triggerVal.TriggerCronConfig.HourMinute != trigger.TriggerCronConfig.HourMinute {
			continue
		}
		if triggerVal.TriggerCronConfig.MonthDay != trigger.TriggerCronConfig.MonthDay {
			continue
		}
		if triggerVal.TriggerCronConfig.LinuxCrontab != trigger.TriggerCronConfig.LinuxCrontab {
			continue
		}
		isOk = true
	}
	return
}

func StartCron(triggerId string, trigger TriggerConfig, robot msql.Params, lang string) (int, error) {
	if !common.CheckLinuxCrontab(trigger.TriggerCronConfig.LinuxCrontab) {
		return 0, errors.New(i18n.Show(lang, `trigger_cron_linux_crontab_error`))
	}
	fmt.Println(LogTriggerPrefix + fmt.Sprintf(`create cron task %q`, trigger.TriggerCronConfig.LinuxCrontab))
	entryId, err := cro.AddFunc(trigger.TriggerCronConfig.LinuxCrontab, func() {
		startTime := time.Now().Unix()
		go TriggerCronRun(robot[`admin_user_id`], robot[`id`], trigger, func(err error) {
			setRunResult(triggerId, startTime, err)
		})
	})
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return int(entryId), nil
}

func StartLoadCronTriggers() {
	logs.Debug(LogTriggerPrefix + `start load trigger cron`)
	list, err := msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`trigger_type`, cast.ToString(TriggerTypeCron)).Select()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	for _, trigger := range list {
		triggerConfig := TriggerConfig{}
		err = tool.JsonDecode(trigger[`trigger_json`], &triggerConfig)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			continue
		}
		isOk, robot := TriggerCronVerifyStartNode(trigger[`robot_id`], triggerConfig)
		if !isOk {
			logs.Debug(LogTriggerPrefix + `trigger not exist`)
			return
		}
		_, err = StartCron(trigger[`id`], triggerConfig, robot, ``)
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			continue
		}
		logs.Debug(LogTriggerPrefix+`load trigger cron ok (%s %s)`, trigger[`id`], triggerConfig.TriggerCronConfig.LinuxCrontab)
	}
}

func DeleteRobotFollow(adminUserId, robotId int) {
	list, err := msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).Select()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	for _, trigger := range list {
		if trigger[`cron_entry_id`] != "" {
			RemoveEntry(cast.ToInt(trigger[`cron_entry_id`]))
		}
		_, err := msql.Model(`work_flow_trigger`, define.Postgres).
			Where(`id`, trigger[`id`]).Delete()
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			continue
		}
		logs.Debug(LogTriggerPrefix+`delete trigger cron ok (%s)`, trigger[`id`])
	}
}
