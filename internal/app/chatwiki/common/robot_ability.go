// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	RobotAbilityAutoReply            = `robot_auto_reply`             //自动回复功能 包含关键词回复 和 收到消息回复
	RobotAbilityKeywordReply         = `robot_keyword_reply`          //关键词回复功能
	RobotAbilityReceivedMessageReply = `robot_received_message_reply` //收到消息回复功能
	RobotAbilitySubscribeReply       = `robot_subscribe_reply`        //关注后回复功能
	RobotAbilitySmartMenu            = `robot_smart_menu`             //智能菜单功能
	OfficialAbilityCustomMenu        = `official_custom_menu`         //自定义菜单功能
)

var DefaultRobotConfig = msql.Params{
	"switch_status": cast.ToString(define.SwitchOff), //开启状态
	"fixed_menu":    cast.ToString(define.SwitchOff), //固定菜单
}

// RobotAbilityList 机器人能力列表
var RobotAbilityList = []Ability{
	{
		Name:          "自动回复",
		ModuleType:    RobotModule,           //模块类型
		AbilityType:   RobotAbilityAutoReply, //全局唯一值 类型
		Introduction:  "用户触发了关键词后，通过设置的关键词回复规则，回复指定的内容",
		Details:       "查看详情中的信息 前端自己定义",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOff,
		SupportChannelsList: []string{
			"公众号",
			"微信客服",
			"WebApp",
			"飞书",
			"钉钉",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "自动回复",
			UniKey: "RobotAbilityAutoReply",
			Path:   "/robot/ability/auto-reply",
			Children: []*define.Menu{
				{
					Name:   "关键词回复",
					UniKey: "RobotAbilityKeyWordReply",
				},
				{
					Name:   "收到消息回复",
					UniKey: "RobotAbilityReceivedMessageReply",
				},
			},
		},
	},
	{
		Name:          "关注后回复",
		ModuleType:    RobotModule,                //模块类型 //本来应该在公众号模块 但wendy要求在机器人的功能中心中展示
		AbilityType:   RobotAbilitySubscribeReply, //全局唯一值 类型
		Introduction:  "关注公众号后，自动给用户回复消息。支持按照时间段以及关注来源设置",
		Details:       "查看详情中的信息 前端自己定义",
		Icon:          "iconfont icon-subscribe-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"公众号",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "关注后回复",
			UniKey:   "RobotAbilitySubscribeReply",
			Path:     "/explore/index/subscribe-reply",
			Children: []*define.Menu{},
		},
	},
	{
		Name:          "自定义菜单",
		ModuleType:    RobotModule,               //模块类型 //本来应该在公众号模块 但wendy要求在机器人的功能中心中展示
		AbilityType:   OfficialAbilityCustomMenu, //全局唯一值 类型
		Introduction:  "支特对公众号菜单自定义，多公众号可单独设置",
		Details:       "查看详情中的信息 前端自己定义",
		Icon:          "iconfont icon-official-custom-menu",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"公众号",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "自定义菜单",
			UniKey:   "OfficialAbilityCustomMenu",
			Path:     "/explore/index/custom-menu",
			Children: []*define.Menu{},
		},
	},
	{
		Name:          "智能菜单",
		ModuleType:    RobotModule,           //模块类型
		AbilityType:   RobotAbilitySmartMenu, //全局唯一值 类型
		Introduction:  "用户点击菜单后，通过设置的菜单回复规则，回复指定的内容",
		Details:       "查看详情中的信息 前端自己定义",
		Icon:          "iconfont icon-smart-menu",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOff,
		SupportChannelsList: []string{
			"公众号",
			"微信客服",
			"WebApp",
			"飞书",
			"钉钉",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "智能菜单",
			UniKey:   "RobotAbilitySmartMenu",
			Path:     "/robot/ability/smart-menu",
			Children: []*define.Menu{},
		},
	},
}

const (
	RobotAbilitySwitchOn  = 1
	RobotAbilitySwitchOff = 0
)

// GetRobotAbility 机器人能力列表
func GetRobotAbility(robotId int) ([]msql.Params, error) {
	data, err := msql.Model(`chat_robot_ability`, define.Postgres).Where(`robot_id`, cast.ToString(robotId)).Order(`id desc`).Select()
	if err == nil && len(data) > 0 {
		for _, item := range data {
			delete(item, "create_time")
			delete(item, "update_time")
			delete(item, "admin_user_id")
		}
	}
	return data, err
}

// GetRobotUseAbilityListByRobotId 获取机器人功能列表
func GetRobotUseAbilityListByRobotId(robotId int, adminUserId int) ([]Ability, error) {
	var list []Ability
	userConfigList, err := GetUserAbilityByModuleType(adminUserId, RobotModule)
	userConfigMap := make(map[string]msql.Params)
	if err == nil && len(userConfigList) > 0 {
		for _, userConfig := range userConfigList {
			userConfigMap[userConfig["ability_type"]] = userConfig
		}
	}
	robotConfigList, err := GetRobotAbility(robotId)
	robotConfigMap := make(map[string]msql.Params)
	if err == nil && len(robotConfigList) > 0 {
		for _, robotConfig := range robotConfigList {
			robotConfigMap[robotConfig["ability_type"]] = robotConfig
		}
	}
	for _, item := range RobotAbilityList {
		userConfig, isUserOk := userConfigMap[item.AbilityType]
		if !isUserOk {
			//用户没有配置的
			continue
		}
		item.UserConfig = userConfig
		if userConfig["switch_status"] != cast.ToString(define.SwitchOn) {
			//用户关闭配置的
			continue
		}
		robotConfig, isOk := robotConfigMap[item.AbilityType]
		if isOk {
			//机器人开启的
			if robotConfig["switch_status"] == cast.ToString(define.SwitchOn) {
				//有配置的 且开启的显示
				item.ShowSelect = define.SwitchOn
				item.RobotConfig = robotConfig
				list = append(list, item)
			}
		}
	}
	return list, err
}

// GetRobotAbilityConfigByAbilityType 检查机器人是否使用某个功能
func GetRobotAbilityConfigByAbilityType(adminUserId int, robotId int, abilityType string) msql.Params {
	//用户是否开启功能
	if CheckUseAbilityByAbilityType(adminUserId, abilityType) {
		//机器人是否开启功能
		data, err := msql.Model(`chat_robot_ability`, define.Postgres).Where(`robot_id`, cast.ToString(robotId)).Where(`ability_type`, abilityType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Find()
		if err == nil && len(data) > 0 {
			return data
		}
	}
	return msql.Params{}
}

// GetRobotAbilityList 机器人功能中心列表
func GetRobotAbilityList(robotId int, adminUserId int, specifyAbilityType string) ([]Ability, error) {
	var list []Ability
	userConfigList, err := GetUserAbilityByModuleType(adminUserId, RobotModule)
	userConfigMap := make(map[string]msql.Params)
	if err == nil && len(userConfigList) > 0 {
		for _, userConfig := range userConfigList {
			userConfigMap[userConfig["ability_type"]] = userConfig
		}
	}
	robotConfigList, err := GetRobotAbility(robotId)
	robotConfigMap := make(map[string]msql.Params)
	if err == nil && len(robotConfigList) > 0 {
		for _, robotConfig := range robotConfigList {
			robotConfigMap[robotConfig["ability_type"]] = robotConfig
		}
	}

	//显示列表能力
	for _, item := range GetAllAbilityList() {
		if specifyAbilityType != `` && item.AbilityType != specifyAbilityType {
			continue
		}
		userConfig, isUserOk := userConfigMap[item.AbilityType]
		if !isUserOk {
			//用户没有配置的
			continue
		}
		item.UserConfig = userConfig
		if userConfig["switch_status"] != cast.ToString(define.SwitchOn) {
			//用户关闭配置的
			continue
		}
		robotConfig, isOk := robotConfigMap[item.AbilityType]
		if isOk {
			//有配置的 直接显示选择
			item.ShowSelect = define.SwitchOn
			item.RobotConfig = robotConfig
			list = append(list, item)
		} else {
			//没有配置的 查看是否显示选择
			item.RobotConfig = DefaultRobotConfig
			if item.ShowSelect == define.SwitchOn {
				list = append(list, item)
			}
		}
	}
	return list, err
}

// SaveRobotAbilitySwitchStatus 保存机器人能力开关状态
func SaveRobotAbilitySwitchStatus(adminUserId int, robotId int, abilityType string, switchStatus int) error {
	// 添加和更新
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// 保存配置
	saveData := msql.Datas{
		"admin_user_id": adminUserId,
		"robot_id":      robotId,
		"ability_type":  abilityType,
		"switch_status": switchStatus,
		"update_time":   tool.Time2Int(),
	}

	model := msql.Model(`chat_robot_ability`, define.Postgres)
	if len(oldData) > 0 {
		_, err = model.Where(`id`, oldData["id"]).Update(saveData)
	} else {
		saveData["create_time"] = tool.Time2Int()
		_, err = model.Insert(saveData)
	}
	return err
}

// SaveRobotAbilityFixedMenu 保存机器人能力固定菜单状态
func SaveRobotAbilityFixedMenu(adminUserId int, robotId int, abilityType string, fixedMenu int) error {
	// 添加和更新
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// 保存配置
	saveData := msql.Datas{
		"admin_user_id": adminUserId,
		"robot_id":      robotId,
		"ability_type":  abilityType,
		"fixed_menu":    fixedMenu,
		"update_time":   tool.Time2Int(),
	}

	model := msql.Model(`chat_robot_ability`, define.Postgres)
	if len(oldData) > 0 {
		_, err = model.Where(`id`, oldData["id"]).Update(saveData)
	} else {
		saveData["create_time"] = tool.Time2Int()
		_, err = model.Insert(saveData)
	}
	return err
}

// SaveRobotAbilityAiReplyStatus 保存机器人能力AI回复状态
func SaveRobotAbilityAiReplyStatus(adminUserId int, robotId int, abilityType string, aiReplyStatus int) error {
	// 添加和更新
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// 保存配置
	saveData := msql.Datas{
		"admin_user_id":   adminUserId,
		"robot_id":        robotId,
		"ability_type":    abilityType,
		"ai_reply_status": aiReplyStatus,
		"update_time":     tool.Time2Int(),
	}

	model := msql.Model(`chat_robot_ability`, define.Postgres)
	if len(oldData) > 0 {
		_, err = model.Where(`id`, oldData["id"]).Update(saveData)
	} else {
		saveData["create_time"] = tool.Time2Int()
		_, err = model.Insert(saveData)
	}
	return err
}

// SaveRobotAbility 保存机器人能力配置
func SaveRobotAbility(adminUserId int, robotId int, abilityType string, switchStatus int, fixedMenu int, aiReplyStatus int) error {
	// 添加和更新
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// 保存配置
	saveData := msql.Datas{
		"admin_user_id":   adminUserId,
		"robot_id":        robotId,
		"ability_type":    abilityType,
		"switch_status":   switchStatus,
		"fixed_menu":      fixedMenu,
		"ai_reply_status": aiReplyStatus,
		"update_time":     tool.Time2Int(),
	}

	model := msql.Model(`chat_robot_ability`, define.Postgres)
	if len(oldData) > 0 {
		_, err = model.Where(`id`, oldData["id"]).Update(saveData)
	} else {
		saveData["create_time"] = tool.Time2Int()
		_, err = model.Insert(saveData)
	}
	return err
}

// GetRobotAbilityByRobotId 根据robot_id获取机器人能力列表
func GetRobotAbilityByRobotId(robotId int) ([]msql.Params, error) {
	data, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Order(`id desc`).
		Select()
	if err == nil && len(data) > 0 {
		for _, item := range data {
			delete(item, "create_time")
			delete(item, "update_time")
			delete(item, "admin_user_id")
		}
	}
	return data, err
}
