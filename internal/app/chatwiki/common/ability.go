// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	RobotModule = `robot`
)

type Ability struct {
	Name                string      `json:"name" form:"name"`                                   //功能名称
	ModuleType          string      `json:"module_type" form:"module_type"`                     //模块类型 区分属于哪个模块
	AbilityType         string      `json:"ability_type" form:"ability_type"`                   //功能类型 全局唯一 长度最大100 一般是 模块key_功能key
	Introduction        string      `json:"introduction" form:"introduction"`                   //功能简介
	Details             string      `json:"details" form:"details"`                             //功能详情
	Icon                string      `json:"icon" form:"icon"`                                   //功能图标
	SupportChannelsList []string    `json:"support_channels_list" form:"support_channels_list"` //支持的渠道列表
	ShowSelect          int         `json:"show_select" form:"show_select"`                     //是否显示选择
	RobotOnlyShow       int         `json:"robot_only_show" form:"robot_only_show"`             //机器人功能中心中仅显示 这个不展示开关和固定菜单
	UserConfig          msql.Params `json:"user_config" form:"user_config"`                     //用户开启模块配置
	RobotConfig         msql.Params `json:"robot_config" form:"robot_config"`                   //模块是robot时候的配置
	Menu                define.Menu `json:"menu" form:"menu"`                                   //菜单配置
}

const DefaultModuleType = `system` //默认系统模块

var DefaultUserConfig = msql.Params{
	"switch_status": cast.ToString(define.SwitchOff), //开启状态
}

// GetAllAbilityList 获取所有功能列表
func GetAllAbilityList() []Ability {
	var abilityList []Ability

	//机器人模块的功能
	abilityList = append(abilityList, RobotAbilityList...)

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
	for _, item := range GetAllAbilityList() {
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
			//没有配置的 查看是否显示选择
			item.UserConfig = DefaultUserConfig
			if item.ShowSelect == define.SwitchOn {
				list = append(list, item)
			}
		}
	}
	return list, err
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
func SaveUserAbility(adminUserId int, abilityType string, switchStatus int) error {
	//添加和更新
	oldData, err := msql.Model(`ability`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`ability_type`, abilityType).Find()
	if err != nil {
		return err
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
	}
	return err
}
