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
	Name                 string      `json:"name" form:"name"`                                     //功能名称
	ModuleType           string      `json:"module_type" form:"module_type"`                       //模块类型 区分属于哪个模块
	AbilityType          string      `json:"ability_type" form:"ability_type"`                     //功能类型 全局唯一 长度最大100 一般是 模块key_功能key
	Introduction         string      `json:"introduction" form:"introduction"`                     //功能简介
	Details              string      `json:"details" form:"details"`                               //功能详情
	Icon                 string      `json:"icon" form:"icon"`                                     //功能图标
	SupportChannelsList  []string    `json:"support_channels_list" form:"support_channels_list"`   //支持的渠道列表
	ShowSelect           int         `json:"show_select" form:"show_select"`                       //是否显示选择
	DefaultSelectedValue int         `json:"default_selected_value" form:"default_selected_value"` // 默认值
	RobotOnlyShow        int         `json:"robot_only_show" form:"robot_only_show"`               //机器人功能中心中仅显示 这个不展示开关和固定菜单
	UserConfig           msql.Params `json:"user_config" form:"user_config"`                       //用户开启模块配置
	RobotConfig          msql.Params `json:"robot_config" form:"robot_config"`                     //模块是robot时候的配置
	Menu                 define.Menu `json:"menu" form:"menu"`                                     //菜单配置
}

const DefaultModuleType = `system` //默认系统模块

// GetAllAbilityList 获取所有功能列表
func GetAllAbilityList() []Ability {
	var abilityList []Ability

	//机器人模块的功能
	abilityList = append(abilityList, RobotAbilityList...)

	//公众号模块的功能
	abilityList = append(abilityList, LibraryAbilityList...)

	//默认模块 没赋值的全是默认模块
	for _, ability := range abilityList {
		if ability.ModuleType == `` {
			ability.ModuleType = DefaultModuleType
		}
	}

	return abilityList
}

// GetUserAbilityByModuleType 获取用户模块功能
func GetUserAbilityByModuleType(adminUserId int, moduleType string) ([]msql.Params, error) {
	data, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`module_type`, moduleType).Order(`id desc`).Select()
	return data, err
}

// GetUserAbility 机器人能力列表
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

// GetUseAbilityListByUserId 获取账号正在使用的功能列表
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
				//有配置的 且开启的显示
				item.ShowSelect = define.SwitchOn
				item.RobotConfig = userConfig
				list = append(list, item)
			}
		}
	}
	return list, err
}

// CheckUseAbilityByAbilityType 检查用户是否使用某个功能
func CheckUseAbilityByAbilityType(adminUserId int, abilityType string) bool {
	data, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Find()
	if err == nil && len(data) > 0 {
		return true
	}
	return false
}

// GetAbilityList 功能中心列表
func GetAbilityList(adminUserId int, specifyAbilityType string) ([]Ability, error) {
	var list []Ability
	configList, err := GetUserAbility(adminUserId)
	configMap := make(map[string]msql.Params)
	if err == nil && len(configList) > 0 {
		for _, userConfig := range configList {
			configMap[userConfig["ability_type"]] = userConfig
		}
	}

	//显示列表能力
	abilityList := GetAllAbilityList()
	for _, item := range abilityList {
		if specifyAbilityType != `` && item.AbilityType != specifyAbilityType {
			continue
		}
		userConfig, isOk := configMap[item.AbilityType]
		if isOk {
			//有配置的 直接显示选择
			item.ShowSelect = define.SwitchOn
			item.UserConfig = userConfig
			list = append(list, item)
		} else {
			//没有配置的 创建默认配置
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

// GetModuleTypeByAbilityType 获取模块类型
func GetModuleTypeByAbilityType(abilityType string) string {
	for _, item := range GetAllAbilityList() {
		if item.AbilityType == abilityType {
			return item.ModuleType
		}
	}
	return DefaultModuleType
}

// SaveUserAbility 保存用户功能
func SaveUserAbility(adminUserId int, abilityType string, switchStatus int) (msql.Params, error) {
	//添加和更新
	oldData, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Find()
	if err != nil {
		return nil, err
	}
	//模块类型
	moduleType := GetModuleTypeByAbilityType(abilityType)

	//保存配置
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
	//处理启用关闭功能逻辑
	AbilitySwitchHandler(adminUserId, abilityType, cast.ToInt(oldData[`switch_status`]))

	return oldData, err
}

func AbilitySwitchHandler(adminUserId int, abilityType string, switchStatus int) {
	if switchStatus == define.SwitchOn {
		//开启功能
		switch abilityType {
		case OfficialAbilityCustomMenu:
			break
		}
	} else {
		//关闭功能
		switch abilityType {
		case OfficialAbilityCustomMenu:
			//删除所有公众号的菜单
			DeleteAllOfficialCustomMenuToWx(adminUserId, ``)
			break
		}

	}
}
