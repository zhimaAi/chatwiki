// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	RobotModule           = `robot`
	LibraryModule         = `library`
	OfficialAccountModule = `official_account`
)

type Ability struct {
	Name                 string      `json:"name" form:"name"`                                     // Function name
	ModuleType           string      `json:"module_type" form:"module_type"`                       // Module type, identifies which module it belongs to
	AbilityType          string      `json:"ability_type" form:"ability_type"`                     // Function type, globally unique, max length 100, usually module_key_function_key
	Introduction         string      `json:"introduction" form:"introduction"`                     // Function introduction
	Details              string      `json:"details" form:"details"`                               // Function details
	Icon                 string      `json:"icon" form:"icon"`                                     // Function icon
	SupportChannelsList  []string    `json:"support_channels_list" form:"support_channels_list"`   // List of supported channels
	ShowSelect           int         `json:"show_select" form:"show_select"`                       // Whether to show selection
	DefaultSelectedValue int         `json:"default_selected_value" form:"default_selected_value"` // Default value
	RobotOnlyShow        int         `json:"robot_only_show" form:"robot_only_show"`               // Only display in robot function center, does not show switch and fixed menu
	UserConfig           msql.Params `json:"user_config" form:"user_config"`                       // User enabled module configuration
	RobotConfig          msql.Params `json:"robot_config" form:"robot_config"`                     // Configuration when module is robot
	Menu                 define.Menu `json:"menu" form:"menu"`                                     // Menu configuration
}

const DefaultModuleType = `system` // Default system module

// GetAllAbilityList gets all function list
func GetAllAbilityList() []Ability {
	var abilityList []Ability

	// Robot module functions
	abilityList = append(abilityList, RobotAbilityList...)

	// Official account module functions
	abilityList = append(abilityList, LibraryAbilityList...)

	// Default module, all unassigned are default modules
	for _, ability := range abilityList {
		if ability.ModuleType == `` {
			ability.ModuleType = DefaultModuleType
		}
	}

	return abilityList
}

// GetUserAbilityByModuleType gets user module functions
func GetUserAbilityByModuleType(adminUserId int, moduleType string) ([]msql.Params, error) {
	data, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`module_type`, moduleType).Order(`id desc`).Select()
	return data, err
}

// GetUserAbility gets robot ability list
func GetUserAbility(adminUserId int) ([]msql.Params, error) {
	data, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if err == nil && len(data) > 0 {
		for _, item := range data {
			delete(item, "create_time")
			delete(item, "update_time")
		}
	}
	return data, err
}

// GetUseAbilityListByUserId gets the list of functions currently used by the account
func GetUseAbilityListByUserId(adminUserId int) ([]Ability, error) {
	var list []Ability
	configList, err := GetUserAbility(adminUserId)
	configMap := make(map[string]msql.Params)
	if err == nil && len(configList) > 0 {
		for _, robotConfig := range configList {
			configMap[robotConfig["ability_type"]] = robotConfig
		}
	}
	for _, item := range GetAllAbilityList() {
		userConfig, isOk := configMap[item.AbilityType]
		if isOk {
			if userConfig["switch_status"] == cast.ToString(define.SwitchOn) {
				// Show if configured and enabled
				item.ShowSelect = define.SwitchOn
				item.RobotConfig = userConfig
				list = append(list, item)
			}
		}
	}
	return list, err
}

// CheckUseAbilityByAbilityType checks if user is using a specific function
func CheckUseAbilityByAbilityType(adminUserId int, abilityType string) bool {
	data, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Find()
	if err == nil && len(data) > 0 {
		return true
	}
	return false
}

// GetAbilityList gets function center list
func GetAbilityList(adminUserId int, specifyAbilityType string) ([]Ability, error) {
	var list []Ability
	configList, err := GetUserAbility(adminUserId)
	configMap := make(map[string]msql.Params)
	if err == nil && len(configList) > 0 {
		for _, userConfig := range configList {
			configMap[userConfig["ability_type"]] = userConfig
		}
	}

	// Display list abilities
	abilityList := GetAllAbilityList()
	for _, item := range abilityList {
		if specifyAbilityType != `` && item.AbilityType != specifyAbilityType {
			continue
		}
		userConfig, isOk := configMap[item.AbilityType]
		if isOk {
			// Show selection if configured
			item.ShowSelect = define.SwitchOn
			item.UserConfig = userConfig
			list = append(list, item)
		} else {
			// Create default configuration if not configured
			info, err := SaveUserAbility(adminUserId, item.AbilityType, item.DefaultSelectedValue)
			if err != nil {
				logs.Error(err.Error())
				return nil, err
			}
			item.ShowSelect = define.SwitchOn
			item.UserConfig = info
			list = append(list, item)
		}
	}
	return list, err
}

// GetAbilityDetail
func GetAbilityDetail(adminUserId int, abilityType string) (Ability, error) {
	var result Ability

	hasFound := false
	for _, item := range GetAllAbilityList() {
		if item.AbilityType != abilityType {
			continue
		}
		hasFound = true
		result = item
	}
	if !hasFound {
		return Ability{}, errors.New("no data")
	}

	configList, err := GetUserAbility(adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return Ability{}, err
	}

	var config msql.Params
	for _, userConfig := range configList {
		if userConfig[`ability_type`] == abilityType {
			config = userConfig
		}
	}

	if len(config) == 0 {
		info, err := SaveUserAbility(adminUserId, abilityType, result.DefaultSelectedValue)
		if err != nil {
			logs.Error(err.Error())
			return Ability{}, err
		}
		result.ShowSelect = define.SwitchOn
		result.UserConfig = info
	} else {
		result.ShowSelect = define.SwitchOn
		result.UserConfig = config
	}

	return result, nil
}

// GetModuleTypeByAbilityType gets the module type
func GetModuleTypeByAbilityType(abilityType string) string {
	for _, item := range GetAllAbilityList() {
		if item.AbilityType == abilityType {
			return item.ModuleType
		}
	}
	return DefaultModuleType
}

// SaveUserAbility saves user ability
func SaveUserAbility(adminUserId int, abilityType string, switchStatus int) (msql.Params, error) {
	// Add and update
	oldData, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Find()
	if err != nil {
		return nil, err
	}
	// Module type
	moduleType := GetModuleTypeByAbilityType(abilityType)

	// Save configuration
	saveData := msql.Datas{
		"admin_user_id": adminUserId,
		"module_type":   moduleType,
		"ability_type":  abilityType,
		"switch_status": switchStatus,
		"update_time":   tool.Time2Int(),
	}
	if len(oldData) > 0 {
		_, err = msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Update(saveData)
	} else {
		saveData["create_time"] = tool.Time2Int()
		_, err = msql.Model(`ability`, define.Postgres).Insert(saveData)
		if err != nil {
			return nil, err
		}
	}

	oldData, err = msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Find()
	// Handle enable/disable function logic
	AbilitySwitchHandler(adminUserId, abilityType, cast.ToInt(oldData[`switch_status`]))

	return oldData, err
}

func AbilitySwitchHandler(adminUserId int, abilityType string, switchStatus int) {
	if switchStatus == define.SwitchOn {
		// Enable function
		switch abilityType {
		case OfficialAbilityCustomMenu:
			break
		}
	} else {
		// Disable function
		switch abilityType {
		case OfficialAbilityCustomMenu:
			// Delete all official account menus
			DeleteAllOfficialCustomMenuToWx(adminUserId, ``)
			break
		}

	}
}
