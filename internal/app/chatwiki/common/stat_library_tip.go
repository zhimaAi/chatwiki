// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func StatLibraryTotal(adminUserId int) (map[string]any, error) {
	tip, err := msql.Model(`chat_ai_stat_library_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Sum(`tip`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return map[string]any{
		`tip_total`: cast.ToInt(tip),
	}, nil
}

func StatLibraryTipUp(result []msql.Params, robot msql.Params) {
	robotId := cast.ToInt(robot[`id`])
	adminUserId := cast.ToInt(robot[`admin_user_id`])
	libraryIds := make([]int, 0)
	dateYmd := time.Now().Format(`20060102`)
	cTime := tool.Time2Int()
	for _, item := range result {
		libraryId := cast.ToInt(item[`library_id`])
		libraryIds = append(libraryIds, libraryId)
		params := []any{adminUserId, robotId, libraryId, cast.ToInt(item[`file_id`]), cast.ToInt(item[`id`]), dateYmd, 1, cTime, cTime}
		sql := fmt.Sprintf(`INSERT INTO chat_ai_stat_library_data_robot_tip (admin_user_id,robot_id,library_id,file_id,data_id,date_ymd,tip,create_time,update_time)
				VALUES (%d,%d,%d,%d,%d,%s,%d,%d,%d)
				ON CONFLICT (admin_user_id, data_id, date_ymd,robot_id)
				DO UPDATE SET tip = chat_ai_stat_library_data_robot_tip.tip + EXCLUDED.tip , update_time = EXCLUDED.update_time;`, params...)
		_, err := msql.RawExec(define.Postgres, sql, nil)
		if err != nil {
			logs.Error(err.Error() + "\n" + sql)
		}
	}
	for _, libraryId := range libraryIds {
		params := []any{adminUserId, robotId, libraryId, dateYmd, 1, cTime, cTime}
		sql := fmt.Sprintf(`INSERT INTO chat_ai_stat_library_robot_tip (admin_user_id, robot_id,library_id, date_ymd, tip, create_time, update_time)
				VALUES (%d,%d,%d,%s,%d,%d,%d)
				ON CONFLICT (admin_user_id, library_id, date_ymd,robot_id)
				DO UPDATE SET tip = chat_ai_stat_library_robot_tip.tip + EXCLUDED.tip , update_time = EXCLUDED.update_time;`, params...)
		_, err := msql.RawExec(define.Postgres, sql, nil)
		if err != nil {
			logs.Error(err.Error() + "\n" + sql)
		}
	}
}

func StatLibraryDataSort(adminUserId, libraryId, page, size int, beginDateYmd, endDateYmd string) (map[string]any, error) {
	m := msql.Model(`chat_ai_stat_library_data_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	mtotal := msql.Model(`chat_ai_stat_library_data_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	if libraryId != 0 {
		m.Where(`library_id`, cast.ToString(libraryId))
		mtotal.Where(`library_id`, cast.ToString(libraryId))
	}
	sorts, err := m.Field(`data_id,sum(tip) as tip`).Group(`data_id`).Order(`tip desc`).Limit((page-1)*size, size).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	total, err := mtotal.Group(`data_id`).Count(`data_id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	mt := msql.Model(`chat_ai_stat_library_data_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	if libraryId != 0 {
		mt.Where(`library_id`, cast.ToString(libraryId))
	}
	totalTip, err := mt.Sum(`tip`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	list, err := statLibraryFormatByDataId(sorts, cast.ToInt(totalTip))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return map[string]any{
		`list`:  list,
		`total`: total,
	}, nil
}

func statLibraryFormatByDataId(sorts []msql.Params, totalTip int) ([]msql.Params, error) {
	if len(sorts) == 0 {
		return sorts, nil
	}
	dataIds := make([]string, 0)
	for _, item := range sorts {
		if !tool.InArray(item[`data_id`], dataIds) {
			if cast.ToInt(item[`data_id`]) > 0 {
				dataIds = append(dataIds, item[`data_id`])
			}
		}
	}
	datas := make([]msql.Params, 0)
	var err error
	if len(dataIds) > 0 {
		datas, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, `in`, strings.Join(dataIds, `,`)).
			Field(`id,library_id,file_id`).Select()
		if err != nil {
			logs.Error(err.Error())
			return nil, err
		}
	}

	libraryIds := make([]string, 0)
	libraryFileIds := make([]string, 0)
	for _, item := range sorts {
		for _, dataInfo := range datas {
			if item[`data_id`] == dataInfo[`id`] {
				item[`library_id`] = dataInfo[`library_id`]
				item[`library_file_id`] = dataInfo[`file_id`]
			}
		}
		dataIds = append(dataIds, item[`data_id`])
		if !tool.InArray(item[`library_id`], libraryIds) {
			if (cast.ToInt(item[`library_id`])) > 0 {
				libraryIds = append(libraryIds, item[`library_id`])
			}
		}
		if !tool.InArray(item[`library_file_id`], libraryFileIds) {
			if (cast.ToInt(item[`library_file_id`])) > 0 {
				libraryFileIds = append(libraryFileIds, item[`library_file_id`])
			}
		}
	}
	//data
	dataInfos := make([]msql.Params, 0)
	if len(dataIds) > 0 {
		dataInfos, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Where(`id`, `in`, strings.Join(dataIds, `,`)).Field(`id as data_id,type,content,question,answer,group_id,file_id`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			return nil, err
		}
	}

	//library
	libraryInfos := make([]msql.Params, 0)
	if len(libraryIds) > 0 {
		libraryInfos, err = msql.Model(`chat_ai_library`, define.Postgres).
			Where(`id`, `in`, strings.Join(libraryIds, `,`)).Field(`id as library_id,library_name,type`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			return nil, err
		}
	}

	//file
	libraryFileInfos := make([]msql.Params, 0)
	if len(libraryFileIds) > 0 {
		libraryFileInfos, err = msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, `in`, strings.Join(libraryFileIds, `,`)).Field(`id as library_file_id,file_name,group_id`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			return nil, err
		}
	}

	//data group
	dataGroupIds := make([]string, 0)
	for _, item := range dataInfos {
		if !tool.InArray(item[`group_id`], dataGroupIds) {
			dataGroupIds = append(dataGroupIds, item[`group_id`])
		}
	}
	//file group
	fileGroupIds := make([]string, 0)
	for _, item := range libraryFileInfos {
		if !tool.InArray(item[`group_id`], fileGroupIds) {
			fileGroupIds = append(fileGroupIds, item[`group_id`])
		}
	}
	dataGroupInfos := make([]msql.Params, 0)
	fileGroupInfos := make([]msql.Params, 0)
	if len(dataGroupIds) > 0 {
		dataGroupInfos, _ = msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`id`, `in`, strings.Join(dataGroupIds, `,`)).Field(`id as group_id,group_name`).
			Select()
	}
	if len(fileGroupIds) > 0 {
		fileGroupInfos, _ = msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`id`, `in`, strings.Join(fileGroupIds, `,`)).Field(`id as group_id,group_name`).
			Select()
	}

	for _, item := range sorts {
		for _, dataInfo := range dataInfos {
			if item[`data_id`] == dataInfo[`data_id`] {
				item[`type`] = dataInfo[`type`]
				item[`content`] = dataInfo[`content`]
				item[`question`] = dataInfo[`question`]
				item[`answer`] = dataInfo[`answer`]
				item[`group_id`] = dataInfo[`group_id`]
			}
		}
		//普通知识库 通过file找到group
		if cast.ToInt(item[`type`]) == define.ParagraphTypeNormal {
			for _, fileInfo := range libraryFileInfos {
				if item[`library_file_id`] == fileInfo[`library_file_id`] {
					item[`group_id`] = fileInfo[`group_id`]
				}
			}
		}
		for _, libraryInfo := range libraryInfos {
			if item[`library_id`] == libraryInfo[`library_id`] {
				item[`library_name`] = libraryInfo[`library_name`]
			}
		}
		if cast.ToInt(item[`type`]) == define.ParagraphTypeDocQA {
			for _, groupInfo := range dataGroupInfos {
				if item[`group_id`] == groupInfo[`group_id`] {
					item[`group_name`] = groupInfo[`group_name`]
				}
			}
		} else {
			for _, groupInfo := range fileGroupInfos {
				if item[`group_id`] == groupInfo[`group_id`] {
					item[`group_name`] = groupInfo[`group_name`]
				}
			}
		}

		for _, libraryFileInfo := range libraryFileInfos {
			if item[`library_file_id`] == libraryFileInfo[`library_file_id`] {
				item[`library_file_name`] = libraryFileInfo[`file_name`]
			}
		}
		item[`percentage`] = fmt.Sprintf("%.2f", cast.ToFloat64(item[`tip`])/cast.ToFloat64(totalTip)*100)
	}
	return sorts, nil
}

func StatLibraryDataRobotDetail(adminUserId, libraryId, dataId int, beginDateYmd, endDateYmd string) ([]msql.Params, error) {
	m := msql.Model(`chat_ai_stat_library_data_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd)).
		Where(`data_id`, cast.ToString(dataId))
	if libraryId != 0 {
		m.Where(`library_id`, cast.ToString(libraryId))
	}
	details, err := m.Field(`robot_id,sum(tip) as tip`).Group(`robot_id`).Order(`tip desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	if len(details) == 0 {
		return []msql.Params{}, nil
	}
	robotId := make([]string, 0)
	for _, item := range details {
		if !tool.InArray(item[`robot_id`], robotId) {
			robotId = append(robotId, item[`robot_id`])
		}
	}
	robotInfos, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, `in`, strings.Join(robotId, `,`)).Field(`id as robot_id,robot_name`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	for _, item := range details {
		for _, robotInfo := range robotInfos {
			if item[`robot_id`] == robotInfo[`robot_id`] {
				item[`robot_name`] = robotInfo[`robot_name`]
			}
		}
	}
	return details, nil
}

func StatLibrarySort(adminUserId, page, size int, beginDateYmd, endDateYmd string) (map[string]any, error) {
	m := msql.Model(`chat_ai_stat_library_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	mtotal := msql.Model(`chat_ai_stat_library_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	sorts, err := m.Field(`library_id,sum(tip) as tip`).Group(`library_id`).Order(`tip desc`).Limit((page-1)*size, size).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	total, err := mtotal.Group(`library_id`).Count(`library_id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	mt := msql.Model(`chat_ai_stat_library_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	totalTip, err := mt.Sum(`tip`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	list, err := statLibraryFormatById(sorts, cast.ToInt(totalTip))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return map[string]any{
		`list`:  list,
		`total`: total,
	}, nil
}

func statLibraryFormatById(sorts []msql.Params, totalTip int) ([]msql.Params, error) {
	if len(sorts) == 0 {
		return sorts, nil
	}
	libraryIds := make([]string, 0)
	for _, item := range sorts {
		if !tool.InArray(item[`library_id`], libraryIds) {
			libraryIds = append(libraryIds, item[`library_id`])
		}
	}
	//library
	libraryInfos, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`id`, `in`, strings.Join(libraryIds, `,`)).Field(`id as library_id,library_name`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	for _, item := range sorts {
		for _, libraryInfo := range libraryInfos {
			if item[`library_id`] == libraryInfo[`library_id`] {
				item[`library_name`] = libraryInfo[`library_name`]
			}
		}
		item[`percentage`] = fmt.Sprintf("%.2f", cast.ToFloat64(item[`tip`])/cast.ToFloat64(totalTip)*100)
	}
	return sorts, nil
}

func StatLibraryRobotDetail(adminUserId, libraryId int, beginDateYmd, endDateYmd string) ([]msql.Params, error) {
	m := msql.Model(`chat_ai_stat_library_robot_tip`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`date_ymd`, `between`, fmt.Sprintf(`%s,%s`, beginDateYmd, endDateYmd))
	m.Where(`library_id`, cast.ToString(libraryId))
	details, err := m.Field(`robot_id,sum(tip) as tip`).Group(`robot_id`).Order(`tip desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	if len(details) == 0 {
		return []msql.Params{}, nil
	}
	robotId := make([]string, 0)
	for _, item := range details {
		if !tool.InArray(item[`robot_id`], robotId) {
			robotId = append(robotId, item[`robot_id`])
		}
	}
	robotInfos, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, `in`, strings.Join(robotId, `,`)).Field(`id as robot_id,robot_name`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	for _, item := range details {
		for _, robotInfo := range robotInfos {
			if item[`robot_id`] == robotInfo[`robot_id`] {
				item[`robot_name`] = robotInfo[`robot_name`]
			}
		}
	}
	return details, nil
}
