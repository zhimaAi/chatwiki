// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"errors"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

type BridgeGetLibraryGroupReq struct {
	LibraryId string `form:"library_id"`
	GroupType string `form:"group_type"`
}

func BridgeGetLibraryGroup(adminUserId, userId int, lang string, req *BridgeGetLibraryGroupReq) ([]msql.Params, int, error) {
	libraryId := cast.ToInt(req.LibraryId)
	groupType := cast.ToInt(req.GroupType)
	if libraryId <= 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	m := msql.Model(`chat_ai_library_group`, define.Postgres)
	wheres := [][]string{{`admin_user_id`, cast.ToString(adminUserId)}, {`library_id`, cast.ToString(libraryId)}}
	list, err := m.Where2(wheres).Where(`group_type`, cast.ToString(groupType)).Field(`id,group_name`).Order(`sort desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	list = append([]msql.Params{{`id`: `0`, `group_name`: lib_define.Ungrouped}}, list...)

	switch groupType {
	case define.LibraryGroupTypeQA:
		// Stats
		stats, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where2(wheres).Where(`isolated`, `false`).Where(`delete_time`, `0`).
			Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		for i, params := range list {
			list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
		}
	case define.LibraryGroupTypeFile:
		// Stats
		stats, err := msql.Model(`chat_ai_library_file`, define.Postgres).
			Where2(wheres).Where(`delete_time`, `0`).
			Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		for i, params := range list {
			list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
		}
	}
	return list, 0, nil
}
