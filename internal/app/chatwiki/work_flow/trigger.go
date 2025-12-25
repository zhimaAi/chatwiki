// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	TriggerTypeChat     = 1 //会话触发器
	TriggerTypeTest     = 2 //测试触发器
	TriggerTypeCron     = 3 //定时触发器
	TriggerTypeOfficial = 4 //公众号触发器
)

const (
	CronTypeSelectTime = `select_time`   //选择触发时间
	CronTypeCrontab    = `linux_crontab` //linux crontab
)

const (
	EveryTypeDay   = `day`   //每天
	EveryTypeWeek  = `week`  //每周
	EveryTypeMonth = `month` //每月
)

type TriggerOutputParam struct {
	StartNodeParam
	Variable string `json:"variable"` //对应的全部变量
}

type TriggerChatConfig struct {
	QuestionMultipleSwitch bool `json:"question_multiple_switch"`
}

type TriggerCronConfig struct {
	Type         string `json:"type"`          //1 选择触发时间,2 linux crontab代码
	LinuxCrontab string `json:"linux_crontab"` //Type为2时 linux crontab
	EveryType    string `json:"every_type"`    //Type为1时 day每天 week每周 month每月
	HourMinute   string `json:"hour_minute"`   //Type为1时 触发的小时分 例如 14:01
	WeekNumber   string `json:"week_number"`   //Type为1时 everyType为week时 存储每周几 0为周日，1-6周一到周六
	MonthDay     string `json:"month_day"`     //Type为1时 everyType为month时 存储每月几号 示例：14表示每月14号
}

type TriggerOfficialConfig struct {
	Event   string `json:"event"`
	MsgType string `json:"msg_type"`
	AppIds  string `json:"app_ids"` //wx1,wx2
}

type TriggerConfig struct {
	TriggerType           uint                  `json:"trigger_type"`
	TriggerName           string                `json:"trigger_name"`
	TriggerIcon           string                `json:"trigger_icon"`
	TriggerSwitch         bool                  `json:"trigger_switch"`
	Outputs               []TriggerOutputParam  `json:"outputs"`
	TriggerChatConfig     TriggerChatConfig     `json:"chat_config"`
	TriggerCronConfig     TriggerCronConfig     `json:"cron_config"`
	TriggerOfficialConfig TriggerOfficialConfig `json:"trigger_official_config"`
}

func (trigger *TriggerConfig) SetGlobalValue(flow *WorkFlow) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(`触发器执行错误 %s`, err)
			logs.Debug(`触发器执行错误 %s`, debug.Stack())
		}
	}()
	assignParams := make(map[string]any)
	switch trigger.TriggerType {
	case TriggerTypeChat: //会话触发器
		assignParams[`openid`] = flow.params.Openid
		assignParams[`question`] = flow.params.Question
		assignParams[`conversationid`] = flow.params.SessionId
		//会话触发器开启多模态输入后的处理逻辑
		if trigger.TriggerChatConfig.QuestionMultipleSwitch {
			if questionMultiple, ok := common.ParseInputQuestion(flow.params.Question); ok {
				assignParams[`question_multiple`] = common.QuestionMultipleAppendImageDomain(questionMultiple)
				assignParams[`question`] = common.GetQuestionByQuestionMultiple(questionMultiple)
			}
		}
	case TriggerTypeTest: //会话触发器
		assignParams = flow.params.TriggerParams.TestParams
	case TriggerTypeOfficial: //公众号触发器
		assignParams = flow.params.TriggerParams.TestParams
	case TriggerTypeCron: //定时触发器
	default:
		logs.Warning(`触发器:%s[%d]赋值逻辑未处理...`, trigger.TriggerName, trigger.TriggerType)
		return
	}
	if len(assignParams) > 0 { //存在传入参数时,给自定义变量赋值
		_ = tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(assignParams), &assignParams)
		for _, output := range trigger.Outputs {
			key, _ := strings.CutPrefix(output.Variable, `global.`)
			if field, ok := flow.global[key]; ok {
				flow.global[key] = field.SetVals(assignParams[output.Key])
			}
		}
	}
}

// GetTriggerConfigByType 根据类型获取触发器配置信息
func GetTriggerConfigByType(triggerType uint) (TriggerConfig, bool) {
	switch triggerType {
	case TriggerTypeChat:
		return GetTriggerChatConfig(), true
	case TriggerTypeTest:
		return GetTriggerTestConfig(), true
	case TriggerTypeCron:
		return GetTriggerCronConfig(msql.Params{`name`: `定时触发器`, `icon`: `/public/trigger_cron_icon.svg`}), true
	case TriggerTypeOfficial:
		return GetTriggerOfficialConfig(msql.Params{`name`: `公众号触发器`, `icon`: `/public/trigger_official_icon.svg`}), true
	}
	return TriggerConfig{}, false
}

// GetTriggerOutputsByType 根据类型获取触发器输出配置
func GetTriggerOutputsByType(triggerType uint) ([]TriggerOutputParam, bool) {
	triggerConfig, exist := GetTriggerConfigByType(triggerType)
	if exist {
		return triggerConfig.Outputs, true
	}
	return nil, false
}

func GetTriggerChatConfig() TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeChat,
		TriggerName:   `会话触发器`,
		TriggerIcon:   `/public/trigger_chat_icon.svg`,
		TriggerSwitch: true,
		Outputs: []TriggerOutputParam{
			{StartNodeParam: StartNodeParam{Key: `openid`, Typ: common.TypString, Required: false, Desc: `用户id`}, Variable: `global.openid`},
			{StartNodeParam: StartNodeParam{Key: `question`, Typ: common.TypString, Required: false, Desc: `用户咨询的问题`}, Variable: `global.question`},
			{StartNodeParam: StartNodeParam{Key: `question_multiple`, Typ: common.TypArrObject, Required: false, Desc: `多模态输入`}, Variable: `global.question_multiple`},
			{StartNodeParam: StartNodeParam{Key: `conversationid`, Typ: common.TypNumber, Required: false, Desc: `会话ID`}, Variable: `global.conversationid`},
		},
	}
}

func GetTriggerTestConfig() TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeTest,
		TriggerName:   `测试触发器`,
		TriggerIcon:   `/public/trigger_test_icon.jpeg`,
		TriggerSwitch: true,
		Outputs: []TriggerOutputParam{
			{StartNodeParam: StartNodeParam{Key: `test_str`, Typ: common.TypString, Required: true, Desc: `测试字符串`}, Variable: `global.test_str`},
			{StartNodeParam: StartNodeParam{Key: `test_num`, Typ: common.TypNumber, Required: true, Desc: `测试数字`}, Variable: `global.test_num`},
			{StartNodeParam: StartNodeParam{Key: `test_arr`, Typ: common.TypArrString, Required: true, Desc: `测试数组`}, Variable: `global.test_arr`},
		},
	}
}

func GetTriggerCronConfig(trigger msql.Params) TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeCron,
		TriggerName:   trigger[`name`],
		TriggerIcon:   trigger[`icon`],
		TriggerSwitch: true,
		TriggerCronConfig: TriggerCronConfig{
			Type:         CronTypeSelectTime,
			LinuxCrontab: "",
			EveryType:    EveryTypeDay,
			HourMinute:   "09:00",
			WeekNumber:   "",
			MonthDay:     "",
		},
	}
}

func GetTriggerOfficialConfig(trigger msql.Params) TriggerConfig {
	return TriggerConfig{
		TriggerType:           TriggerTypeOfficial,
		TriggerName:           trigger[`name`],
		TriggerIcon:           trigger[`icon`],
		TriggerSwitch:         true,
		TriggerOfficialConfig: TriggerOfficialConfig{},
	}
}

func GetTriggerConfigList(adminUserId int, lang string) ([]TriggerConfig, error) {
	triggerList := make([]TriggerConfig, 0)
	triggerList = append(triggerList, GetTriggerChatConfig()) //会话触发器
	if define.IsDev {
		triggerList = append(triggerList, GetTriggerTestConfig()) //测试触发器
	}
	list, err := TriggerList(adminUserId, lang)
	if err != nil {
		return nil, err
	}
	for _, trigger := range list {
		if cast.ToInt(trigger[`switch_status`]) == 0 {
			continue
		}
		if cast.ToInt(trigger[`trigger_type`]) == TriggerTypeCron {
			triggerList = append(triggerList, GetTriggerCronConfig(trigger))
		} else if cast.ToInt(trigger[`trigger_type`]) == TriggerTypeOfficial {
			triggerList = append(triggerList, GetTriggerOfficialConfig(trigger))
		}
	}
	return triggerList, nil
}

func TriggerList(adminUserId int, lang string) ([]msql.Params, error) {
	list, err := msql.Model(`trigger_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id asc`).Select()
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(list) == 0 {
		TriggerInitDefault(adminUserId)
		list, err = msql.Model(`trigger_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id asc`).Select()
		if err != nil {
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	//create official trigger
	boolCreateOfficial := true
	for _, val := range list {
		if cast.ToInt(val[`trigger_type`]) == TriggerTypeOfficial {
			boolCreateOfficial = false
			break
		}
	}
	if boolCreateOfficial {
		TriggerInitOfficial(adminUserId)
		list, err = msql.Model(`trigger_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id asc`).Select()
		if err != nil {
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	for key, val := range list {
		if cast.ToInt(val[`trigger_type`]) == TriggerTypeCron {
			list[key][`icon`] = `/public/trigger_cron_icon.svg`
		} else if cast.ToInt(val[`trigger_type`]) == TriggerTypeOfficial {
			list[key][`icon`] = `/public/trigger_official_icon.svg`
		}
	}
	return list, nil
}

func TriggerInitOfficial(adminUserId int) {
	lockKey := `trigger_init_official`
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	_, err := msql.Model(`trigger_config`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`switch_status`: 1,
		`name`:          `公众号触发器`,
		`trigger_type`:  TriggerTypeOfficial,
		`intro`:         `安装该插件后，当公众号推送外部事件时，自动触发工作流`,
		`author`:        `chatwiki`,
		`from_type`:     define.FromInherited,
		`create_time`:   time.Now().Unix(),
		`update_time`:   time.Now().Unix(),
	})
	if err != nil {
		logs.Error(err.Error())
	} else {
		lib_redis.DelCacheData(define.Redis, common.TriggerConfigCacheBuildHandler{
			AdminUserId: adminUserId,
			TriggerType: cast.ToString(TriggerTypeOfficial),
		})
	}
}

func TriggerInitDefault(adminUserId int) {
	lockKey := `trigger_init_cron`
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	_, err := msql.Model(`trigger_config`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`switch_status`: 1,
		`name`:          `定时触发器`,
		`trigger_type`:  TriggerTypeCron,
		`intro`:         `安装该插件后，支持在设定的时间自动触发工作流`,
		`author`:        `chatwiki`,
		`from_type`:     define.FromInherited,
		`create_time`:   time.Now().Unix(),
		`update_time`:   time.Now().Unix(),
	})
	if err != nil {
		logs.Error(err.Error())
	} else {
		lib_redis.DelCacheData(define.Redis, common.TriggerConfigCacheBuildHandler{
			AdminUserId: adminUserId,
			TriggerType: cast.ToString(TriggerTypeCron),
		})
	}
}

// SaveTriggerConfig 触发器保存
func SaveTriggerConfig(robot msql.Params, node *WorkFlowNode, lang string) error {
	if len(node.NodeParams.Start.TriggerList) == 0 {
		return nil
	}
	//clean olds
	existAllTriggers, err := msql.Model(`work_flow_trigger`, define.Postgres).
		Where(`robot_id`, robot[`id`]).Select()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	for _, existTrigger := range existAllTriggers {
		if cast.ToInt(existTrigger[`cron_entry_id`]) > 0 {
			logs.Debug(LogTriggerPrefix+`remove entry id %s %s`, existTrigger[`id`], existTrigger[`cron_entry_id`])
			RemoveEntry(cast.ToInt(existTrigger[`cron_entry_id`]))
		}
	}
	_, err = msql.Model(`work_flow_trigger`, define.Postgres).Where(`robot_id`, robot[`id`]).Delete()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	for _, trigger := range node.NodeParams.Start.TriggerList {
		switch trigger.TriggerType {
		case TriggerTypeCron:
			err = SaveTriggerCronConfig(trigger, robot, lang)
		case TriggerTypeOfficial:
			err = SaveTriggerOfficialConfig(robot[`admin_user_id`], trigger, robot, lang)
		}
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	return nil
}
