// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

const (
	OfficialAccountAbilityBatchSend = `official_account_batch_send` // Official account article mass sending
	OfficialAccountAbilityAIComment = `official_account_ai_comment` // Official account article AI comments
	LibraryAbilityOfficialAccount   = `library_ability_official_account`
)

var LibraryAbilityList = []Ability{
	{
		Name:          "[[ZM--OfficialAccountBatchSendName--ZM]]",
		ModuleType:    RobotModule,                     // Module type
		AbilityType:   OfficialAccountAbilityBatchSend, // Globally unique value type
		Introduction:  "[[ZM--OfficialAccountBatchSendIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "[[ZM--ArticleGroupSendName--ZM]]",
			UniKey: "exploreArticleGroupSend",
			Path:   "/explore/index/article-group-send",
			Children: []*define.Menu{
				{
					Name:   "[[ZM--GroupSendName--ZM]]",
					UniKey: "exploreArticleGroupSendGroupSend",
					Path:   "/explore/index/article-group-send/group-send",
				},
				{
					Name:   "[[ZM--DraftBoxName--ZM]]",
					UniKey: "exploreArticleGroupSendDraftBox",
					Path:   "/explore/index/article-group-send/draft-box",
				},
			},
		},
	},
	{
		Name:          "[[ZM--OfficialAccountAICommentName--ZM]]",
		ModuleType:    RobotModule,                     // Module type
		AbilityType:   OfficialAccountAbilityAIComment, // Globally unique value type
		Introduction:  "[[ZM--OfficialAccountAICommentIntro--ZM]]",
		Details:       "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:          "iconfont icon-keyword-reply",
		ShowSelect:    define.SwitchOn,
		RobotOnlyShow: define.SwitchOn,
		SupportChannelsList: []string{
			"[[ZM--OfficialAccountChannel--ZM]]",
		},
		RobotConfig: DefaultRobotConfig,
		Menu: define.Menu{
			Name:   "[[ZM--AICommentManagementName--ZM]]",
			UniKey: "exploreAiCommentManagement",
			Path:   "/explore/index/ai-comment-management",
			Children: []*define.Menu{
				{
					Name:   "[[ZM--DefaultRuleName--ZM]]",
					UniKey: "exploreAiCommentManagementDefaultRule",
					Path:   "/explore/index/ai-comment-management/default-rule",
				},
				{
					Name:   "[[ZM--CustomRuleName--ZM]]",
					UniKey: "exploreAiCommentManagementCustomRule",
					Path:   "/explore/index/ai-comment-management/custom-rule",
				},
				{
					Name:   "[[ZM--CreateCustomRuleName--ZM]]",
					UniKey: "exploreAiCommentManagementCreateCustomRule",
					Path:   "/explore/index/ai-comment-management/create-custom-rule",
				},
				{
					Name:   "[[ZM--CommentProcessingRecordName--ZM]]",
					UniKey: "exploreAiCommentManagementCommentProcessingRecord",
					Path:   "/explore/index/ai-comment-management/comment-processing-record",
				},
			},
		},
	}, {
		Name:                 "[[ZM--LibraryAbilityOfficialAccountName--ZM]]",
		ModuleType:           LibraryModule,                 // Module type
		AbilityType:          LibraryAbilityOfficialAccount, // Globally unique value type
		Introduction:         "[[ZM--LibraryAbilityOfficialAccountIntro--ZM]]",
		Details:              "[[ZM--DetailsPlaceholder--ZM]]",
		Icon:                 "iconfont icon-official-account",
		ShowSelect:           define.SwitchOn,
		DefaultSelectedValue: define.SwitchOn,
		SupportChannelsList: []string{
			"[[ZM--AuthenticatedOfficialAccountChannel--ZM]]",
		},
		Menu: define.Menu{
			Name:     "[[ZM--LibraryAbilityOfficialAccountName--ZM]]",
			UniKey:   "LibraryAbilityOfficialAccount",
			Path:     "/library/list",
			Children: []*define.Menu{},
		},
	},
}

// CheckHasEnabledOfficialAccount checks if the official account knowledge base switch is enabled
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
