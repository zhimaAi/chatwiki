// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetParagraphSplit(userId, fileId, pdfPageNum int, fileDataIds string, splitParams define.SplitParams, lang string) (list []define.DocSplitItem, wordTotal int, _splitParams define.SplitParams, err error) {
	info, err := GetLibFileInfo(fileId, userId)
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	if len(info) == 0 {
		err = errors.New(i18n.Show(lang, `no_data`))
		return
	}
	library, err := GetLibraryInfo(cast.ToInt(info[`library_id`]), userId)
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	if len(library) == 0 {
		err = errors.New(i18n.Show(lang, `no_data`))
		return
	}
	splitParams.IsTableFile = cast.ToInt(info[`is_table_file`])
	splitParams, err = CheckSplitParams(library, splitParams, lang)
	if err != nil {
		return
	}
	// get file data
	list, err = GetFileDataInfo(fileDataIds, fileId, pdfPageNum, lang)
	// split by document type
	if splitParams.IsQaDoc == define.DocTypeQa {
		if cast.ToInt(info[`is_table_file`]) != define.FileIsTable {
			list = QaDocSplit(splitParams, list)
		}
	} else {
		list, err = MultDocSplit(cast.ToInt(info[`admin_user_id`]), fileId, 0, splitParams, list)
	}
	if err != nil {
		return
	}
	for i := range list {
		list[i].Number = i + 1 //serial number
		if splitParams.IsQaDoc == define.DocTypeQa {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question) + utf8.RuneCountInString(list[i].Answer)
		} else {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
		}
	}
	_splitParams = splitParams

	return
}

func GetFileDataInfo(id string, FileDataId, pdfPageNum int, lang string) (list []define.DocSplitItem, err error) {
	data, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`id`, `in`, id).Select()
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	for _, item := range data {
		var images []string
		if item[`images`] != "" {
			tool.JsonDecode(cast.ToString(item[`images`]), &images)
		}
		list = append(list, define.DocSplitItem{
			FileDataId: cast.ToInt(item[`id`]),
			PageNum:    cast.ToInt(item[`page_num`]),
			Title:      cast.ToString(item[`title`]),
			Content:    cast.ToString(item[`content`]),
			Answer:     item[`answer`],
			Images:     images,
			WordTotal:  cast.ToInt(item[`word_total`]),
			Number:     cast.ToInt(item[`number`]),
		})
	}
	return list, nil
}

func SaveSplitParagraph(adminUserId, fileId, dataId int, list []define.DocSplitItem) ([]define.DocSplitItem, error) {
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	data, err := m.Where(`file_id`, cast.ToString(fileId)).Order(`number asc`).Field(`id,number`).Select()
	if err != nil {
		logs.Error(err.Error())
		return list, err
	}
	var number int
	for _, item := range data {
		if cast.ToInt(item[`id`]) == dataId {
			number = cast.ToInt(item[`number`])
			for i := range list {
				list[i].Number = number + i
			}
			break
		}
	}
	for _, item := range data {
		if cast.ToInt(item[`id`]) == dataId {
			m.Where(`id`, item[`id`]).Delete()
			msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`data_id`, item[`id`]).Delete()
			continue
		}
		if cast.ToInt(item[`number`]) < number {
			continue
		}
		m.Where(`id`, item[`id`]).Update(msql.Datas{`number`: cast.ToInt(item[`number`]) + len(list) - 1, `update_time`: cast.ToString(tool.Time2Int())})
	}
	return list, nil
}
