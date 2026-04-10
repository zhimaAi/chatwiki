// Copyright 漏 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		return MigrateRobotMultilingualConfig()
	}, nil)
}

func MigrateRobotMultilingualConfig() error {
	list, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,admin_user_id,welcomes,unknown_question_prompt,tips_before_answer_content,tips_before_answer_switch,enable_common_question,common_question_list`).
		Select()
	if err != nil {
		logs.Error(`MigrateRobotMultilingualConfig query robots err:%s`, err.Error())
		return nil
	}
	for _, one := range list {
		robotId := cast.ToInt(one[`id`])
		adminUserId := cast.ToInt(one[`admin_user_id`])
		if robotId <= 0 || adminUserId <= 0 {
			continue
		}
		robotLangConfigs := common.BuildRobotMultilingualConfigsByLegacy(
			one[`welcomes`],
			one[`unknown_question_prompt`],
			one[`tips_before_answer_content`],
			cast.ToBool(one[`tips_before_answer_switch`]),
			cast.ToBool(one[`enable_common_question`]),
			one[`common_question_list`],
		)
		for i := range robotLangConfigs {
			if robotLangConfigs[i].LangKey == common.RobotLangCh {
				robotLangConfigs[i].TipsBeforeAnswerSwitch = cast.ToBool(one[`tips_before_answer_switch`])
			}
		}
		if err = common.SaveRobotMultilingualConfigs(adminUserId, robotId, robotLangConfigs); err != nil {
			logs.Error(`MigrateRobotMultilingualConfig robot_id:%d err:%s`, robotId, err.Error())
			return nil
		}
	}
	return nil
}
