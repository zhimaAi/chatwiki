// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"net/http"
	"strings"
	"time"
)

func StatMessageFeedback(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	todayStats, err := messageFeedbackStats(userId, robotId, now.BeginningOfDay().Unix(), 0)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	yesterdayStats, err := messageFeedbackStats(userId, robotId, now.BeginningOfDay().Add(-time.Hour*24).Unix(), now.BeginningOfDay().Unix())
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	weekStats, err := messageFeedbackStats(userId, robotId, now.BeginningOfDay().Add(-time.Hour*24*6).Unix(), 0)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	totalStats, err := messageFeedbackStats(userId, robotId, 0, now.EndOfMinute().Unix())

	data := map[string]any{`today_stats`: todayStats, `yesterday_stats`: yesterdayStats, `week_stats`: weekStats, `total_stats`: totalStats}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func messageFeedbackStats(adminUserId, robotId int, start, end int64) (msql.Params, error) {
	m := msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId))
	if start > 0 {
		m.Where(`create_time`, `>=`, cast.ToString(start))
	}
	if end > 0 {
		m.Where(`create_time`, `<`, cast.ToString(end))
	}
	stats, err := m.Field(`count(case when type=1 then 1 end) as like_count,count(case when type=2 then 1 end) as dislike`).Find()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetMessageFeedbackList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	_type := cast.ToInt(c.Query(`type`))
	startDate := strings.TrimSpace(c.Query(`start_date`))
	endDate := strings.TrimSpace(c.Query(`end_date`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))

	m := msql.Model(`message_feedback`, define.Postgres).
		Alias(`f`).
		Join(`chat_ai_message m1`, `f.ai_message_id=m1.id`, `left`).
		Join(`chat_ai_message m2`, `f.customer_message_id=m2.id`, `left`).
		Where(`f.admin_user_id`, cast.ToString(userId)).
		Where(`f.robot_id`, cast.ToString(robotId)).
		Order(`f.id desc`).Field(`f.*,m1.content as answer,m2.content as question,m1.quote_file`)
	if len(startDate) > 0 {
		startTimeStamp, err := now.Parse(startDate)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `start_date`))))
			return
		}
		m.Where(`f.create_time`, `>=`, cast.ToString(startTimeStamp.Unix()))
	}
	if len(endDate) > 0 {
		endTimeStamp, err := now.Parse(endDate + ` 23:59:59`)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `end_date`))))
			return
		}
		m.Where(`f.create_time`, `<`, cast.ToString(endTimeStamp.Unix()))
	}
	if _type > 0 {
		m.Where(`type`, cast.ToString(_type))
	}
	list, total, err := m.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetMessageFeedbackDetail(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	result, err := msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// get question content
	question, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, result[`customer_message_id`]).
		Value(`content`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(question) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	// get answer content and quote file
	answer, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, result[`ai_message_id`]).
		Field(`content,quote_file`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(answer) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	// get quote file content
	type quoteFile struct {
		FileName string `json:"file_name"`
		Id       string `json:"id"`
	}
	var quoteFiles []quoteFile
	err = json.Unmarshal([]byte(answer[`quote_file`]), &quoteFiles)

	var fileIds []string
	for _, quoteFile := range quoteFiles {
		fileIds = append(fileIds, quoteFile.Id)
	}
	answerSources, err := msql.Model(`chat_ai_answer_source`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`message_id`, result[`ai_message_id`]).
		Where(`file_id`, `in`, strings.Join(fileIds, `,`)).
		Field(`type,content,question,answer,file_id`).
		Select()

	// join quote answer
	type QuoteAnswer struct {
		Type     int
		Content  string
		Question string
		Answer   string
		FileId   int
		FileName string
	}
	var quoteAnswers []QuoteAnswer
	for _, quoteFile := range quoteFiles {
		for _, answerSource := range answerSources {
			if quoteFile.Id == answerSource[`file_id`] {
				quoteAnswers = append(quoteAnswers, QuoteAnswer{
					Type:     cast.ToInt(answerSource[`type`]),
					Content:  answerSource[`content`],
					Question: answerSource[`question`],
					Answer:   answerSource[`answer`],
					FileId:   cast.ToInt(quoteFile.Id),
					FileName: quoteFile.FileName,
				})
			}
		}
	}
	quotes, err := tool.JsonEncode(quoteAnswers)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	result[`question`] = question
	result[`answer`] = answer[`content`]
	result[`quotes`] = quotes

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}
