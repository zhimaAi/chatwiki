// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/lib/pq"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func SaveUnknownIssueRecord(lang string, adminUserId int, robot msql.Params, question string) {
	question = GetFirstQuestionByInput(question) // Special handling for multimodal input
	m := msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres)
	id, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`robot_id`, robot[`id`]).
		Where(`stats_day`, cast.ToString(tool.GetYmd(0))).Where(`question`, question).Value(`id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	if cast.ToUint(id) == 0 { // No data
		ins := msql.Datas{
			`admin_user_id`: adminUserId,
			`robot_id`:      robot[`id`],
			`stats_day`:     tool.GetYmd(0),
			`question`:      question,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}
		newId, err := m.Insert(ins, `id`)
		if err != nil {
			var sqlerr *pq.Error
			if errors.As(err, &sqlerr) && sqlerr.Code == `23505` { // Unique index constraint
				SaveUnknownIssueRecord(lang, adminUserId, robot, question)
				return
			}
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return
		}
		id = cast.ToString(newId)
		// Asynchronous processing of unknown issue summary logic
		if cast.ToUint(robot[`unknown_summary_status`]) > 0 {
			go SaveUnknownIssueSummaryRecord(lang, adminUserId, robot, MbSubstr(question, 0, MaxContent))
		}
	}
	// Start updating data
	sqlraw := fmt.Sprintf(`trigger_total=trigger_total+1,update_time=%d`, tool.Time2Int())
	if _, err = m.Where(`id`, id).Update2(sqlraw); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
}

func DisposeStringList(sliceStr string, items ...string) []string {
	sliceList, temporary := make([]string, 0), make([]string, 0)
	_ = tool.JsonDecodeUseNumber(sliceStr, &temporary)
	for _, item := range append(temporary, items...) {
		if !tool.InArrayString(item, sliceList) {
			sliceList = append(sliceList, item)
		}
	}
	return sliceList
}

func SaveUnknownIssueSummaryRecord(lang string, adminUserId int, robot msql.Params, question string) {
	embedding, err := GetVector2000(lang, adminUserId, robot[`admin_user_id`], robot, msql.Params{}, msql.Params{},
		cast.ToInt(robot[`unknown_summary_model_config_id`]), robot[`unknown_summary_use_model`], question)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	triggerDay := cast.ToString(tool.GetYmd(0))
	m := msql.Model(`chat_ai_unknown_issue_summary`, define.Postgres)
	summary, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).Where(`trigger_day`, triggerDay).
		Where(`vector_dims(embedding)`, cast.ToString(len(strings.Split(embedding, `,`)))).
		Field(fmt.Sprintf(`max(1-(embedding<=>'%s')) as similarity`, embedding)).
		Group(`id`).Order(`similarity DESC`).Field(`id,question,unknown_list`).Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	if len(summary) == 0 || cast.ToFloat64(summary[`similarity`]) < cast.ToFloat64(robot[`unknown_summary_similarity`]) { // Insert data
		ins := msql.Datas{
			`admin_user_id`: adminUserId,
			`robot_id`:      robot[`id`],
			`trigger_day`:   triggerDay,
			`question`:      question,
			`embedding`:     embedding,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}
		_, err = m.Insert(ins)
	} else { // Update data
		if summary[`question`] == question {
			return // Same question, skip update
		}
		unknownList := DisposeStringList(summary[`unknown_list`], question)
		unknownListStr := tool.JsonEncodeNoError(unknownList)
		if summary[`unknown_list`] == unknownListStr {
			return // Same question, skip update
		}
		data := msql.Datas{
			`unknown_list`:  unknownListStr,
			`unknown_total`: len(unknownList),
			`update_time`:   tool.Time2Int(),
		}
		_, err = m.Where(`id`, summary[`id`]).Update(data)
	}
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
}

func ExportUnknownIssueSummary(lang string, list []msql.Params, ext string) (string, error) {
	if define.IsDocxFile(ext) {
		var lineBreak = "\r" // Special line break for docx
		doc := document.New()
		var imageConfig = &document.ImageConfig{Size: &document.ImageSize{Width: 145, KeepAspectRatio: true}}
		for idx, params := range list {
			if idx > 0 { // Add two line breaks
				doc.AddParagraph(lineBreak + lineBreak)
			}
			// Question
			para := doc.AddParagraph(``)
			para.AddFormattedText(i18n.Show(lang, `export_question_header`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`question`], &document.TextFormat{FontSize: 12})
			// Similar questions
			para.AddFormattedText(lineBreak+i18n.Show(lang, `export_similar_questions_header`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
			for _, item := range DisposeStringList(params[`unknown_list`]) {
				para.AddFormattedText(lineBreak+item, &document.TextFormat{FontSize: 12})
			}
			// Answer
			para.AddFormattedText(lineBreak+i18n.Show(lang, `export_answer_header`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`answer`], &document.TextFormat{FontSize: 12})
			// Images
			for _, imgUrl := range DisposeStringList(params[`images`]) {
				if !LinkExists(imgUrl) {
					continue
				}
				if _, err := doc.AddImageFromFile(GetFileByLink(imgUrl), imageConfig); err != nil {
					logs.Error(err.Error())
				}
			}
		}
		md := tool.MD5(tool.JsonEncodeNoError(list) + time.Now().String() + tool.Random(10))
		filepath := `static/public/download/` + md[:2] + `/` + md[2:] + `.docx`
		return filepath, doc.Save(filepath)
	} else {
		fields := tool.Fields{
			{Field: `question`, Header: i18n.Show(lang, `export_question_header`)},
			{Field: `unknown_list`, Header: i18n.Show(lang, `export_similar_questions_header`)},
			{Field: `answer`, Header: i18n.Show(lang, `export_answer_header`)},
		}
		data := make([]map[string]any, len(list))
		for idx, params := range list {
			for _, imgUrl := range DisposeStringList(params[`images`]) {
				params[`answer`] += fmt.Sprintf("\r\n{{!!%s!!}}", imgUrl)
			}
			data[idx] = map[string]any{
				`question`: params[`question`], `answer`: params[`answer`],
				`unknown_list`: strings.Join(DisposeStringList(params[`unknown_list`]), "\r\n"),
			}
		}
		filepath, _, err := tool.ExcelExportPro(data, fields, i18n.Show(lang, `export_unknown_issue_summary`), `static/public/download`)
		return filepath, err
	}
}
