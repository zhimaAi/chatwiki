// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import "chatwiki/internal/app/chatwiki/define"

const (
	OfficialAccountAbilityBatchSend = `official_account_batch_send` //公众号文章群发
	OfficialAccountAbilityAIComment = `official_account_ai_comment` //公众号文章AI评论
)

var OfficialAccountAbility = []Ability{
	{
		Name:          "文章群发",
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
		Name:          "AI评论管理精选",
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
	},
}
