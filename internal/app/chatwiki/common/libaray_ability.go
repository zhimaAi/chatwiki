// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

const (
	OfficialAccountAbilityBatchSend = `official_account_batch_send` //公众号文章群发
	OfficialAccountAbilityAIComment = `official_account_ai_comment` //公众号文章AI评论
	LibraryAbilityOfficialAccount   = `library_ability_official_account`
)

var LibraryAbilityList = []Ability{
	{
		Name:          "公众号文章群发",
		ModuleType:    RobotModule,                     //模块类型
		AbilityType:   OfficialAccountAbilityBatchSend, //全局唯一值 类型
		Introduction:  "群发文章给公众号用户，该群发支持使用AI评论精选功能，自动删评，回复，精选评论",
		Details:       "群发文章给公众号用户，该群发支持使用AI评论精选功能，自动删评，回复，精选评论",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"公众号",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "文章群发",
			UniKey: "exploreArticleGroupSend",
			Path:   "/explore/index/article-group-send",
			Children: []*define.Menu{
				{
					Name:   "群发管理",
					UniKey: "exploreArticleGroupSendGroupSend",
					Path:   "/explore/index/article-group-send/group-send",
				},
				{
					Name:   "草稿箱",
					UniKey: "exploreArticleGroupSendDraftBox",
					Path:   "/explore/index/article-group-send/draft-box",
				},
			},
		},
	},
	{
		Name:          "公众号AI评论精选  ",
		ModuleType:    RobotModule,                     //模块类型
		AbilityType:   OfficialAccountAbilityAIComment, //全局唯一值 类型
		Introduction:  "对系统内发布的文章群发支持自动删评，回复，精选评论 ",
		Details:       "对系统内发布的文章群发支持自动删评，回复，精选评论 ",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"公众号",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "AI评论管理",
			UniKey: "exploreAiCommentManagement",
			Path:   "/explore/index/ai-comment-management",
			Children: []*define.Menu{
				{
					Name:   "默认规则'",
					UniKey: "exploreAiCommentManagementDefaultRule",
					Path:   "/explore/index/ai-comment-management/default-rule",
				},
				{
					Name:   "自定义规则",
					UniKey: "exploreAiCommentManagementCustomRule",
					Path:   "/explore/index/ai-comment-management/custom-rule",
				},
				{
					Name:   "新建自定义规则",
					UniKey: "exploreAiCommentManagementCreateCustomRule",
					Path:   "/explore/index/ai-comment-management/create-custom-rule",
				},
				{
					Name:   "评论处理记录",
					UniKey: "exploreAiCommentManagementCommentProcessingRecord",
					Path:   "/explore/index/ai-comment-management/comment-processing-record",
				},
			},
		},
	}, {
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
