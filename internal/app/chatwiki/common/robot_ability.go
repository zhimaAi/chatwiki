// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	RobotAbilityAutoReply            = `robot_auto_reply`             // Auto-reply function including keyword reply and received message reply
	RobotAbilityPayment              = `robot_payment`                // Application payment
	RobotAbilityKeywordReply         = `robot_keyword_reply`          // Keyword reply function
	RobotAbilityReceivedMessageReply = `robot_received_message_reply` // Received message reply function
	RobotAbilitySubscribeReply       = `robot_subscribe_reply`        // Subscribe reply function
	RobotAbilitySmartMenu            = `robot_smart_menu`             // Smart menu function
	OfficialAbilityCustomMenu        = `official_custom_menu`         // Custom menu function
)

var DefaultRobotConfig = msql.Params{
	"switch_status":   cast.ToString(define.SwitchOff), // Enable status
	"fixed_menu":      cast.ToString(define.SwitchOff), // Fixed menu
	"show_select":     cast.ToString(define.SwitchOn),  // Whether to show selection
	"robot_only_show": cast.ToString(define.SwitchOff), // Only display in robot function center, this does not show switches and fixed menus

}

// RobotAbilityList Robot capability list
var RobotAbilityList = []Ability{
	{
		Name:          "[[ZM--AutoReplyName--ZM]]",
		ModuleType:    RobotModule,           // Module type
		AbilityType:   RobotAbilityAutoReply, // Globally unique value type
		Introduction:  "[[ZM--AutoReplyIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOff,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
			"[[ZM--WeChatCustomerServiceChannel--ZM]]",
			"[[ZM--WebAppChannel--ZM]]",
			"[[ZM--FeiShuChannel--ZM]]",
			"[[ZM--DingTalkChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "[[ZM--AutoReplyName--ZM]]",
			UniKey: "RobotAbilityAutoReply",
			Path:   "/robot/ability/auto-reply",
			Children: []*define.Menu{
				{
					Name:   "[[ZM--KeywordReplyName--ZM]]",
					UniKey: "RobotAbilityKeyWordReply",
				},
				{
					Name:   "[[ZM--ReceivedMessageReplyName--ZM]]",
					UniKey: "RobotAbilityReceivedMessageReply",
				},
			},
		},
	},
	{
		Name:          "[[ZM--PaymentName--ZM]]",
		ModuleType:    RobotModule,         // Module type
		AbilityType:   RobotAbilityPayment, // Globally unique value type
		Introduction:  "[[ZM--PaymentIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-keyword-payment",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOff,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "[[ZM--PaymentName--ZM]]",
			UniKey:   "RobotAbilityPayment",
			Path:     "/robot/ability/payment",
			Children: []*define.Menu{},
		},
	},
	{
		Name:          "[[ZM--SmartMenuName--ZM]]",
		ModuleType:    RobotModule,           // Module type
		AbilityType:   RobotAbilitySmartMenu, // Globally unique value type
		Introduction:  "[[ZM--SmartMenuIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-smart-menu",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOff,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
			"[[ZM--WeChatCustomerServiceChannel--ZM]]",
			"[[ZM--WebAppChannel--ZM]]",
			"[[ZM--FeiShuChannel--ZM]]",
			"[[ZM--DingTalkChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "[[ZM--SmartMenuName--ZM]]",
			UniKey:   "RobotAbilitySmartMenu",
			Path:     "/robot/ability/smart-menu",
			Children: []*define.Menu{},
		},
	},
	{
		Name:          "[[ZM--SubscribeReplyName--ZM]]",
		ModuleType:    RobotModule,                // Module type // Originally should be in official account module but wendy requested to display in robot function center
		AbilityType:   RobotAbilitySubscribeReply, // Globally unique value type
		Introduction:  "[[ZM--SubscribeReplyIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-subscribe-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "[[ZM--SubscribeReplyMenuName--ZM]]",
			UniKey:   "RobotAbilitySubscribeReply",
			Path:     "/explore/index/subscribe-reply",
			Children: []*define.Menu{},
		},
	},
	{
		Name:          "[[ZM--CustomMenuName--ZM]]",
		ModuleType:    RobotModule,               // Module type // Originally should be in official account module but wendy requested to display in robot function center
		AbilityType:   OfficialAbilityCustomMenu, // Globally unique value type
		Introduction:  "[[ZM--CustomMenuIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-official-custom-menu",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:     "[[ZM--CustomMenuMenuName--ZM]]",
			UniKey:   "OfficialAbilityCustomMenu",
			Path:     "/explore/index/custom-menu",
			Children: []*define.Menu{},
		},
	},
}

const (
	RobotAbilitySwitchOn  = 1
	RobotAbilitySwitchOff = 0
)

// GetRobotAbility Robot capability list
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

// GetRobotUseAbilityListByRobotId Get robot function list
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
			// User has not configured
			continue
		}
		item.UserConfig = userConfig
		if userConfig["switch_status"] != cast.ToString(define.SwitchOn) {
			// User has disabled the configuration
			continue
		}
		robotConfig, isOk := robotConfigMap[item.AbilityType]
		if isOk {
			// Robot enabled
			if robotConfig["switch_status"] == cast.ToString(define.SwitchOn) {
				// Configured and enabled ones are displayed
				item.ShowSelect = define.SwitchOn
				item.RobotConfig = robotConfig
				list = append(list, item)
			}
		}
	}
	return list, err
}

// GetRobotAbilityConfigByAbilityType Check if robot uses a certain function
func GetRobotAbilityConfigByAbilityType(adminUserId int, robotId int, abilityType string) msql.Params {
	// Whether user has enabled the function
	if CheckUseAbilityByAbilityType(adminUserId, abilityType) {
		// Whether robot has enabled the function
		data, err := msql.Model(`chat_robot_ability`, define.Postgres).Where(`robot_id`, cast.ToString(robotId)).Where(`ability_type`, abilityType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Find()
		if err == nil && len(data) > 0 {
			return data
		}
	}
	return msql.Params{}
}

// GetRobotAbilityList Robot function center list
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

	// Display list capabilities
	for _, item := range GetAllAbilityList() {
		if specifyAbilityType != `` && item.AbilityType != specifyAbilityType {
			continue
		}
		if specifyAbilityType == `` { // Display non-specified ones
			userConfig, isUserOk := userConfigMap[item.AbilityType]
			if !isUserOk {
				// User has not configured
				continue
			}
			item.UserConfig = userConfig
			if userConfig["switch_status"] != cast.ToString(define.SwitchOn) {
				// User has disabled the configuration
				continue
			}
		} else { // Display specified configuration
			userConfig, isUserOk := userConfigMap[item.AbilityType]
			if isUserOk {
				item.UserConfig = userConfig
			} else {
				item.UserConfig = DefaultRobotConfig // Default closed
			}
		}
		robotConfig, isOk := robotConfigMap[item.AbilityType]
		if isOk {
			// Configured ones, directly display selection
			item.ShowSelect = define.SwitchOn
			item.RobotConfig = robotConfig
			list = append(list, item)
		} else {
			// Unconfigured ones, check whether to display selection
			item.RobotConfig = DefaultRobotConfig
			if item.ShowSelect == define.SwitchOn {
				list = append(list, item)
			}
		}
	}
	return list, err
}

// SaveRobotAbilitySwitchStatus Save robot capability switch status
func SaveRobotAbilitySwitchStatus(adminUserId int, robotId int, abilityType string, switchStatus int) error {
	// Add and update
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// Save configuration
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

// SaveRobotAbilityFixedMenu Save robot capability fixed menu status
func SaveRobotAbilityFixedMenu(adminUserId int, robotId int, abilityType string, fixedMenu int) error {
	// Add and update
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// Save configuration
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

// SaveRobotAbilityAiReplyStatus Save robot capability AI reply status
func SaveRobotAbilityAiReplyStatus(adminUserId int, robotId int, abilityType string, aiReplyStatus int) error {
	// Add and update
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// Save configuration
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

// SaveRobotAbility Save robot capability configuration
func SaveRobotAbility(adminUserId int, robotId int, abilityType string, switchStatus int, fixedMenu int, aiReplyStatus int) error {
	// Add and update
	oldData, err := msql.Model(`chat_robot_ability`, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`ability_type`, abilityType).
		Find()
	if err != nil {
		return err
	}

	// Save configuration
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

// GetRobotAbilityByRobotId Get robot capability list by robot_id
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
