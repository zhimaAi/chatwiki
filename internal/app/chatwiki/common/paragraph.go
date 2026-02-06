// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetParagraphSplit(userId, fileId, pdfPageNum int, fileDataIds string, splitParams define.SplitParams, lang string) (list define.DocSplitItems, wordTotal int, _splitParams define.SplitParams, err error) {
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
	list.UnifySetNumber() // Unify number assignment
	for i := range list {
		if splitParams.IsQaDoc == define.DocTypeQa {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question) + utf8.RuneCountInString(list[i].Answer)
		} else {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
		}
	}
	_splitParams = splitParams

	return
}

func GetFileDataInfo(id string, FileDataId, pdfPageNum int, lang string) (list define.DocSplitItems, err error) {
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
			// Father-son chunking
			FatherChunkParagraphNumber: cast.ToInt(item[`father_chunk_paragraph_number`]),
		})
	}
	list.UnifySetNumber() // Unify number assignment
	return list, nil
}

func SaveSplitParagraph(adminUserId, fileId int, paragraph msql.Params, list define.DocSplitItems) (define.DocSplitItems, error) {
	if len(paragraph) == 0 { // PDF single page re-segmentation, compatibility handling
		fatherChunkParagraphNumber, number, sqlErr := GetAddParagraphNumbers(int64(fileId))
		if sqlErr != nil {
			logs.Error(sqlErr.Error())
			return list, sqlErr
		}
		// Set numbers for new data
		for i := range list {
			list[i].FatherChunkParagraphNumber, list[i].Number = fatherChunkParagraphNumber, number+i
		}
		return list, nil
	}
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	datas, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`file_id`, cast.ToString(fileId)).
		Where(`father_chunk_paragraph_number`, paragraph[`father_chunk_paragraph_number`]).
		Order(`page_num,father_chunk_paragraph_number,number`).Field(`id,number`).Select()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return list, err
	}
	indexModel := msql.Model(`chat_ai_library_file_data_index`, define.Postgres)
	var number int // Number variable
	for _, item := range datas {
		if item[`id`] == paragraph[`id`] { // Re-segmented paragraph
			// Delete paragraph and paragraph index
			_, err = m.Where(`id`, item[`id`]).Delete()
			if err != nil {
				logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			}
			_, err = indexModel.Where(`data_id`, item[`id`]).Delete()
			if err != nil {
				logs.Error(`sql:%s,err:%s`, indexModel.GetLastSql(), err.Error())
			}
			// Set numbers for new data
			for i := range list {
				number++
				list[i].FatherChunkParagraphNumber = cast.ToInt(paragraph[`father_chunk_paragraph_number`]) // Keep unchanged
				list[i].Number = number
				list[i].PageNum = cast.ToInt(paragraph[`page_num`]) // Keep unchanged
			}
		} else { // Other data directly modify number
			number++
			if cast.ToInt(item[`id`]) != number { // Correct number
				_, err = m.Where(`id`, item[`id`]).Update(msql.Datas{`number`: number, `update_time`: cast.ToString(tool.Time2Int())})
				if err != nil {
					logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
				}
			}
		}
	}
	return list, nil
}

func UpdateParagraphHits(dataIds string, today int) {
	sqlRaw := fmt.Sprintf(`update_time = %v`, tool.Time2Int())
	if today > 0 {
		sqlRaw = fmt.Sprintf(`%v,today_hits = today_hits + %v,total_hits = total_hits+%v`, sqlRaw, today, today)
	} else if today == 0 {
		sqlRaw = fmt.Sprintf(`%v,today_hits = 0`, sqlRaw)
	}
	msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`id`, `in`, cast.ToString(dataIds)).
		Update2(sqlRaw)
}

func GetAddParagraphNumbers(fileId int64) (int, int, error) {
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	maxFn, err := m.Where(`file_id`, cast.ToString(fileId)).Max(`father_chunk_paragraph_number`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, 0, err
	}
	if cast.ToInt(maxFn) > 0 { // Father-son chunking, start a new father chunk
		return cast.ToInt(maxFn) + 1, 1, nil
	}
	maxNumber, err := m.Where(`file_id`, cast.ToString(fileId)).Max(`number`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, 0, err
	}
	return 0, cast.ToInt(maxNumber) + 1, nil
}

func RefreshParagraphNumbers(fileId int64) {
	m := msql.Model(`chat_ai_library_file_data`, define.Postgres)
	datas, err := m.Where(`file_id`, cast.ToString(fileId)).
		Field(`id,father_chunk_paragraph_number,number`).
		Order(`page_num,father_chunk_paragraph_number,number`).
		Select()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	list := make(define.DocSplitItems, len(datas))
	for i, item := range datas {
		list[i] = define.DocSplitItem{
			Number:                     cast.ToInt(item[`number`]),
			FatherChunkParagraphNumber: cast.ToInt(item[`father_chunk_paragraph_number`]),
		}
	}
	list.UnifySetNumber() // Unify number assignment
	// Correct numbers for all data
	for i, item := range datas {
		if cast.ToInt(item[`id`]) != list[i].Number {
			_, err = m.Where(`id`, item[`id`]).Update(msql.Datas{`number`: list[i].Number, `update_time`: cast.ToString(tool.Time2Int())})
			if err != nil {
				logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			}
		}
	}
}
