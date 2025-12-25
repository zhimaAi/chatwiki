// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
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

func SaveUnknownIssueRecord(adminUserId int, robot msql.Params, question string) {
	question = GetFirstQuestionByInput(question) //多模态输入特殊处理
	m := msql.Model(`chat_ai_unknown_issue_stats`, define.Postgres)
	id, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`robot_id`, robot[`id`]).
		Where(`stats_day`, cast.ToString(tool.GetYmd(0))).Where(`question`, question).Value(`id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	if cast.ToUint(id) == 0 { //没有数据
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
			if errors.As(err, &sqlerr) && sqlerr.Code == `23505` { //唯一索引约束
				SaveUnknownIssueRecord(adminUserId, robot, question)
				return
			}
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return
		}
		id = cast.ToString(newId)
		//未知问题总结逻辑异步处理
		if cast.ToUint(robot[`unknown_summary_status`]) > 0 {
			go SaveUnknownIssueSummaryRecord(adminUserId, robot, MbSubstr(question, 0, MaxContent))
		}
	}
	//开始更新数据
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

func SaveUnknownIssueSummaryRecord(adminUserId int, robot msql.Params, question string) {
	embedding, err := GetVector2000(adminUserId, robot[`admin_user_id`], robot, msql.Params{}, msql.Params{},
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
	if len(summary) == 0 || cast.ToFloat64(summary[`similarity`]) < cast.ToFloat64(robot[`unknown_summary_similarity`]) { //插入数据
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
	} else { //更新数据
		if summary[`question`] == question {
			return //相同的问题,跳过更新
		}
		unknownList := DisposeStringList(summary[`unknown_list`], question)
		unknownListStr := tool.JsonEncodeNoError(unknownList)
		if summary[`unknown_list`] == unknownListStr {
			return //相同的问题,跳过更新
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

func ExportUnknownIssueSummary(list []msql.Params, ext string) (string, error) {
	if define.IsDocxFile(ext) {
		var lineBreak = "\r" //docx的换行符,比较特殊
		doc := document.New()
		var imageConfig = &document.ImageConfig{Size: &document.ImageSize{Width: 145, KeepAspectRatio: true}}
		for idx, params := range list {
			if idx > 0 { //添加两个换行
				doc.AddParagraph(lineBreak + lineBreak)
			}
			//问题
			para := doc.AddParagraph(``)
			para.AddFormattedText(`问题：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`question`], &document.TextFormat{FontSize: 12})
			//相似问法
			para.AddFormattedText(lineBreak+`相似问法：`, &document.TextFormat{Bold: true, FontSize: 14})
			for _, item := range DisposeStringList(params[`unknown_list`]) {
				para.AddFormattedText(lineBreak+item, &document.TextFormat{FontSize: 12})
			}
			//答案
			para.AddFormattedText(lineBreak+`答案：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`answer`], &document.TextFormat{FontSize: 12})
			//图片
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
		fields := tool.Fields{{Field: "question", Header: "问题"}, {Field: "unknown_list", Header: "相似问题"}, {Field: "answer", Header: "答案"}}
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
		filepath, _, err := tool.ExcelExportPro(data, fields, `未知问题总结`, `static/public/download`)
		return filepath, err
	}
}
