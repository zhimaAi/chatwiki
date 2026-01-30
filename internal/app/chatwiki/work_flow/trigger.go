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
	TriggerTypeChat     = 1 // Chat Trigger
	TriggerTypeTest     = 2 // Test Trigger
	TriggerTypeCron     = 3 // Cron Trigger
	TriggerTypeOfficial = 4 // Official Account Trigger
	TriggerTypeWebHook  = 5 // Webhook Trigger
)

const (
	CronTypeSelectTime = `select_time`   // Select trigger time
	CronTypeCrontab    = `linux_crontab` // linux crontab
)

const (
	EveryTypeDay   = `day`   // Every day
	EveryTypeWeek  = `week`  // Every week
	EveryTypeMonth = `month` // Every month
)

type TriggerOutputParam struct {
	StartNodeParam
	Variable string `json:"variable"` // Corresponding all variables
}

type TriggerChatConfig struct {
	QuestionMultipleSwitch bool `json:"question_multiple_switch"`
}

type TriggerCronConfig struct {
	Type         string `json:"type"`          // 1 Select trigger time, 2 linux crontab code
	LinuxCrontab string `json:"linux_crontab"` // When Type is 2, linux crontab
	EveryType    string `json:"every_type"`    // When Type is 1, day/week/month
	HourMinute   string `json:"hour_minute"`   // When Type is 1, trigger hour/minute e.g. 14:01
	WeekNumber   string `json:"week_number"`   // When Type is 1 and everyType is week, store day of week (0 for Sunday, 1-6 for Mon-Sat)
	MonthDay     string `json:"month_day"`     // When Type is 1 and everyType is month, store day of month (e.g., 14 means 14th)
}

type TriggerOfficialConfig struct {
	Event   string `json:"event"`
	MsgType string `json:"msg_type"`
	AppIds  string `json:"app_ids"` // wx1, wx2
}

const (
	WebHookResponseTypeNow             = `now`
	WebHookResponseTypeMessageVariable = `message_variable`
)

type TriggerWebHookConfig struct {
	Url                string               `json:"url"`             // Request URL
	Method             string               `json:"method"`          // Method GET POST
	SwitchVerify       string               `json:"switch_verify"`   // Whether to authenticate, 1 yes, 0 no
	SwitchAllowIp      string               `json:"switch_allow_ip"` // Whether to verify IP, 1 yes, 0 no
	AllowIps           string               `json:"allow_ips"`       // IP whitelist, multiple separated by comma, max 1000
	Params             common.RecurveFields `json:"params"`          // URL parameters
	RequestContentType string               `json:"request_content_type"`
	Form               common.RecurveFields `json:"form"`
	XForm              common.RecurveFields `json:"x_form"`
	Json               common.RecurveFields `json:"json"`
	ResponseType       string               `json:"response_type"` // Return type: now returns json immediately, message_variable returns message and variables
	ResponseNow        string               `json:"response_now"`  // String returned immediately
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
	TriggerWebHookConfig  TriggerWebHookConfig  `json:"trigger_web_hook_config"`
}

func (trigger *TriggerConfig) SetGlobalValue(flow *WorkFlow) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(`trigger run faild %s`, err)
			logs.Debug(`trigger run faild %s`, debug.Stack())
		}
	}()
	assignParams := make(map[string]any)
	switch trigger.TriggerType {
	case TriggerTypeChat: // Chat Trigger
		assignParams[`openid`] = flow.params.Openid
		assignParams[`question`] = flow.params.Question
		assignParams[`conversationid`] = flow.params.SessionId
		// Logic for handling chat trigger when multimodal input is enabled
		if trigger.TriggerChatConfig.QuestionMultipleSwitch {
			if questionMultiple, ok := common.ParseInputQuestion(flow.params.Question); ok {
				assignParams[`question_multiple`] = common.QuestionMultipleAppendImageDomain(questionMultiple)
				assignParams[`question`] = common.GetQuestionByQuestionMultiple(questionMultiple)
			}
		}
	case TriggerTypeTest: // Test Trigger
		assignParams = flow.params.TriggerParams.TestParams
	case TriggerTypeOfficial: // Official Account Trigger
		assignParams = flow.params.TriggerParams.TestParams
	case TriggerTypeWebHook: // Webhook Trigger
		assignParams = flow.params.TriggerParams.TestParams
		trigger.Outputs = flow.params.TriggerParams.TriggerOutputs
	case TriggerTypeCron: // Cron Trigger
	default:
		logs.Warning(`Trigger:%s[%d] assignment logic not handled...`, trigger.TriggerName, trigger.TriggerType)
		return
	}
	if len(assignParams) > 0 { // When input parameters exist, assign values to custom variables
		_ = tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(assignParams), &assignParams)
		for _, output := range trigger.Outputs {
			key, _ := strings.CutPrefix(output.Variable, `global.`)
			if field, ok := flow.global[key]; ok {
				flow.global[key] = field.SetVals(assignParams[output.Key])
			}
		}
	}
}

// GetTriggerConfigByType gets trigger config by type
func GetTriggerConfigByType(triggerType uint, lang string) (TriggerConfig, bool) {
	switch triggerType {
	case TriggerTypeChat:
		return GetTriggerChatConfig(lang), true
	case TriggerTypeTest:
		return GetTriggerTestConfig(lang), true
	case TriggerTypeCron:
		return GetTriggerCronConfig(msql.Params{`name`: i18n.Show(lang, `timed_trigger`), `icon`: `/public/trigger_cron_icon.svg`}), true
	case TriggerTypeOfficial:
		return GetTriggerOfficialConfig(msql.Params{`name`: i18n.Show(lang, `official_account_trigger`), `icon`: `/public/trigger_official_icon.svg`}), true
	case TriggerTypeWebHook:
		return GetTriggerOfficialConfig(msql.Params{`name`: i18n.Show(lang, `webhook_trigger`), `icon`: `/public/trigger_webhook_icon.svg`}), true
	}
	return TriggerConfig{}, false
}

// GetTriggerOutputsByType gets trigger output config by type
func GetTriggerOutputsByType(triggerType uint, lang string) ([]TriggerOutputParam, bool) {
	triggerConfig, exist := GetTriggerConfigByType(triggerType, lang)
	if exist {
		return triggerConfig.Outputs, true
	}
	return nil, false
}

func GetTriggerChatConfig(lang string) TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeChat,
		TriggerName:   i18n.Show(lang, `chat_trigger`),
		TriggerIcon:   `/public/trigger_chat_icon.svg`,
		TriggerSwitch: true,
		Outputs: []TriggerOutputParam{
			{StartNodeParam: StartNodeParam{Key: `openid`, Typ: common.TypString, Required: false, Desc: i18n.Show(lang, `user_id`)}, Variable: `global.openid`},
			{StartNodeParam: StartNodeParam{Key: `question`, Typ: common.TypString, Required: false, Desc: i18n.Show(lang, `user_question`)}, Variable: `global.question`},
			{StartNodeParam: StartNodeParam{Key: `question_multiple`, Typ: common.TypArrObject, Required: false, Desc: i18n.Show(lang, `multi_modal_input`)}, Variable: `global.question_multiple`},
			{StartNodeParam: StartNodeParam{Key: `conversationid`, Typ: common.TypNumber, Required: false, Desc: i18n.Show(lang, `excel_header_session_id`)}, Variable: `global.conversationid`},
		},
	}
}

func GetTriggerTestConfig(lang string) TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeTest,
		TriggerName:   i18n.Show(lang, `test_trigger`),
		TriggerIcon:   `/public/trigger_test_icon.jpeg`,
		TriggerSwitch: true,
		Outputs: []TriggerOutputParam{
			{StartNodeParam: StartNodeParam{Key: `test_str`, Typ: common.TypString, Required: true, Desc: i18n.Show(lang, `test_string`)}, Variable: `global.test_str`},
			{StartNodeParam: StartNodeParam{Key: `test_num`, Typ: common.TypNumber, Required: true, Desc: i18n.Show(lang, `test_number`)}, Variable: `global.test_num`},
			{StartNodeParam: StartNodeParam{Key: `test_arr`, Typ: common.TypArrString, Required: true, Desc: i18n.Show(lang, `test_array`)}, Variable: `global.test_arr`},
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

func GetTriggerWebHookConfig(trigger msql.Params, robotKey string) TriggerConfig {
	return TriggerConfig{
		TriggerType:   TriggerTypeWebHook,
		TriggerName:   trigger[`name`],
		TriggerIcon:   trigger[`icon`],
		TriggerSwitch: true,
		TriggerWebHookConfig: TriggerWebHookConfig{
			Url: GetWebHookUrl(robotKey),
		},
	}
}

func GetWebHookUrl(robotKey string) string {
	return define.Config.WebService["api_domain"] + `/open/workflow/webhook/` + robotKey + `/` + tool.Random(10)
}

func GetTriggerConfigList(adminUserId int, robotKey, lang string) ([]TriggerConfig, error) {
	triggerList := make([]TriggerConfig, 0)
	triggerList = append(triggerList, GetTriggerChatConfig(lang)) // Chat Trigger
	if define.IsDev {
		triggerList = append(triggerList, GetTriggerTestConfig(lang)) // Test Trigger
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
		} else if cast.ToInt(trigger[`trigger_type`]) == TriggerTypeWebHook {
			triggerList = append(triggerList, GetTriggerWebHookConfig(trigger, robotKey))
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
		TriggerInitDefault(adminUserId, lang)
		list, err = msql.Model(`trigger_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id asc`).Select()
		if err != nil {
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	//create official trigger
	boolCreateOfficial := true
	boolCreateWebhook := true
	for _, val := range list {
		if cast.ToInt(val[`trigger_type`]) == TriggerTypeOfficial {
			boolCreateOfficial = false
			continue
		}
		if cast.ToInt(val[`trigger_type`]) == TriggerTypeWebHook {
			boolCreateWebhook = false
			continue
		}
	}
	if boolCreateOfficial {
		TriggerInitOfficial(adminUserId, lang)
		list, err = msql.Model(`trigger_config`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id asc`).Select()
		if err != nil {
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	if boolCreateWebhook {
		TriggerInitWebhook(adminUserId, lang)
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
		} else if cast.ToInt(val[`trigger_type`]) == TriggerTypeWebHook {
			list[key][`icon`] = `/public/trigger_webhook_icon.svg`
		}
	}
	return list, nil
}

func TriggerInitOfficial(adminUserId int, lang string) {
	lockKey := `trigger_init_official`
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	_, err := msql.Model(`trigger_config`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`switch_status`: 1,
		`name`:          i18n.Show(lang, `official_account_trigger`),
		`trigger_type`:  TriggerTypeOfficial,
		`intro`:         i18n.Show(lang, `official_account_trigger_intro`),
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

func TriggerInitWebhook(adminUserId int, lang string) {
	lockKey := `trigger_init_webhook`
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	_, err := msql.Model(`trigger_config`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`switch_status`: 1,
		`name`:          i18n.Show(lang, `webhook_trigger`),
		`trigger_type`:  TriggerTypeWebHook,
		`intro`:         i18n.Show(lang, `webhook_trigger_intro`),
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
			TriggerType: cast.ToString(TriggerTypeWebHook),
		})
	}
}

func TriggerInitDefault(adminUserId int, lang string) {
	lockKey := `trigger_init_cron`
	if !lib_redis.AddLock(define.Redis, lockKey, time.Second*5) {
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	_, err := msql.Model(`trigger_config`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`switch_status`: 1,
		`name`:          i18n.Show(lang, `timed_trigger`),
		`trigger_type`:  TriggerTypeCron,
		`intro`:         i18n.Show(lang, `timed_trigger_intro`),
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

// SaveTriggerConfig Save trigger configuration
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
		case TriggerTypeWebHook:
			err = SaveTriggerWebhookConfig(robot[`admin_user_id`], trigger, robot, lang)
		}
		if err != nil {
			logs.Error(LogTriggerPrefix + err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	return nil
}
