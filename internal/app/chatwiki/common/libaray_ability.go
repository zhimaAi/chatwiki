// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

const LibraryAbilityOfficialAccount = `library_ability_official_account`

var LibraryAbilityList = []Ability{
	{
		Name:                 "微信公众号知识库",
		ModuleType:           LibraryModule,                 //模块类型
		AbilityType:          LibraryAbilityOfficialAccount, //全局唯一值 类型
		Introduction:         "功能开启后，可以在知识库>公众号知识库模块中绑定微信公众号，自动获取公众号已发布的文章",
		Details:              "查看详情中的信息 前端自己定义",
		Icon:                 "iconfont icon-official-account",
		ShowSelect:           define.SwitchOn,
		DefaultSelectedValue: define.SwitchOn,
		SupportChannelsList: []string{
			"已认证公众号，服务号",
		},
		Menu: define.Menu{
			Name:     "微信公众号知识库",
			UniKey:   "LibraryAbilityOfficialAccount",
			Path:     "/library/list",
			Children: []*define.Menu{},
		},
	},
}

// CheckHasEnabledOfficialAccount 检查是否已启动公众号知识库开关
func CheckHasEnabledOfficialAccount(adminUserId int) (bool, error) {
	data, err := msql.Model(`ability`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`module_type`, LibraryModule).
		Where(`ability_type`, LibraryAbilityOfficialAccount).
		Find()
	if err != nil {
		return false, err
	}
	if len(data) == 0 {
		return false, nil
	}
	if cast.ToInt(data[`switch_status`]) == 0 {
		return false, nil
	}
	return true, nil
}
