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
	"strconv"
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

type messageFeedbackContext struct {
	feedback msql.Params
	question msql.Params
	answer   msql.Params
	session  msql.Params
}

// Resolve the conversation from the owned feedback, never from client-supplied session or customer IDs.
func loadMessageFeedbackContext(c *gin.Context, userId int) (*messageFeedbackContext, error) {
	id, err := strconv.Atoi(c.Query(`id`))
	if err != nil || id <= 0 {
		return nil, errors.New(`param_lack`)
	}
	result, err := msql.Model(`message_feedback`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(`sys_err`)
	}
	if len(result) == 0 {
		return nil, errors.New(`no_data`)
	}
	if !common.CheckOperateRights(userId, getLoginUserId(c), define.IdentityTypeUser,
		cast.ToInt(result[`robot_id`]), define.ObjectTypeRobot, define.PermissionQueryRights, common.GetLang(c)) {
		return nil, errors.New(`query_no_permission`)
	}
	// get question content
	question, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, result[`robot_id`]).
		Where(`id`, result[`customer_message_id`]).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(`sys_err`)
	}
	if len(question) == 0 {
		return nil, errors.New(`no_data`)
	}

	// get answer content and quote file
	answer, err := msql.Model(`chat_ai_message`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, result[`robot_id`]).
		Where(`id`, result[`ai_message_id`]).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(`sys_err`)
	}
	if len(answer) == 0 {
		return nil, errors.New(`no_data`)
	}
	if question[`is_customer`] != `1` || answer[`is_customer`] != `0` ||
		question[`openid`] != answer[`openid`] || question[`dialogue_id`] != answer[`dialogue_id`] ||
		question[`session_id`] != answer[`session_id`] || cast.ToInt(answer[`session_id`]) <= 0 {
		return nil, errors.New(`no_data`)
	}
	session, err := msql.Model(`chat_ai_session`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, result[`robot_id`]).
		Where(`id`, answer[`session_id`]).Where(`openid`, answer[`openid`]).
		Where(`dialogue_id`, answer[`dialogue_id`]).
		Field(`id,admin_user_id,robot_id,openid,dialogue_id,app_type,app_id,rel_user_id`).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(`sys_err`)
	}
	if len(session) == 0 {
		return nil, errors.New(`no_data`)
	}
	return &messageFeedbackContext{feedback: result, question: question, answer: answer, session: session}, nil
}

func GetMessageFeedbackDetail(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	context, err := loadMessageFeedbackContext(c, userId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), err.Error()))))
		return
	}
	result, answer := context.feedback, context.answer

	// get quote file content
	type quoteFile struct {
		FileName string `json:"file_name"`
		Id       string `json:"id"`
	}
	var quoteFiles []quoteFile
	err = json.Unmarshal([]byte(answer[`quote_file`]), &quoteFiles)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	var fileIds []string
	for _, quoteFile := range quoteFiles {
		fileIds = append(fileIds, quoteFile.Id)
	}
	var answerSources []msql.Params
	if len(fileIds) > 0 {
		answerSources, err = msql.Model(`chat_ai_answer_source`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`message_id`, result[`ai_message_id`]).
			Where(`file_id`, `in`, strings.Join(fileIds, `,`)).
			Field(`type,content,question,answer,file_id`).Select()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

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

	result[`question`] = context.question[`content`]
	result[`answer`] = answer[`content`]
	result[`quotes`] = quotes

	// The chat customer cache builder may backfill names and avatars; viewing feedback must not write them.
	customer, err := msql.Model(`chat_ai_customer`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`openid`, answer[`openid`]).
		Field(`name,nickname,avatar`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if cast.ToInt(context.session[`rel_user_id`]) > 0 {
		user, err := msql.Model(define.TableUser, define.Postgres).
			Where(`id`, context.session[`rel_user_id`]).
			Field(`id,parent_id,nick_name,avatar`).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if user[`id`] == cast.ToString(userId) || user[`parent_id`] == cast.ToString(userId) {
			common.FillRelUserInfo2(customer, cast.ToInt(context.session[`rel_user_id`]), user)
		}
	}
	channelName := context.session[`app_type`]
	for _, channel := range common.GetChannelList(userId, cast.ToUint(result[`robot_id`])) {
		if channel.AppType == context.session[`app_type`] && (channel.AppId == `` || channel.AppId == context.session[`app_id`]) {
			channelName = channel.AppName
			break
		}
	}
	data := map[string]any{`customer`: customer, `session`: context.session, `channel_name`: channelName}
	for key, value := range result {
		data[key] = value
	}
	common.FmtOk(c, data)
}

func GetMessageFeedbackMessages(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	minId, maxId := 0, 0
	for key, target := range map[string]*int{`min_id`: &minId, `max_id`: &maxId} {
		if raw, exists := c.GetQuery(key); exists {
			value, err := strconv.Atoi(raw)
			if err != nil || value <= 0 {
				common.FmtError(c, `param_invalid`, key)
				return
			}
			*target = value
		}
	}
	if minId > 0 && maxId > 0 {
		common.FmtError(c, `param_invalid`, `min_id,max_id`)
		return
	}
	context, err := loadMessageFeedbackContext(c, adminUserId)
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	query := func() *msql.Builder {
		return msql.Model(`chat_ai_message`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`robot_id`, context.feedback[`robot_id`]).
			Where(`session_id`, context.session[`id`]).Where(`dialogue_id`, context.session[`dialogue_id`]).
			Where(`openid`, context.session[`openid`])
	}
	const size = 20
	var before, after []msql.Params
	anchor := context.answer[`id`]
	if minId > 0 {
		anchor = cast.ToString(minId)
	} else if maxId > 0 {
		anchor = cast.ToString(maxId)
	}
	if maxId == 0 {
		before, err = query().Where(`id`, `<`, anchor).Order(`id desc`).Limit(size + 1).Select()
	}
	if err == nil && minId == 0 {
		after, err = query().Where(`id`, `>`, anchor).Order(`id asc`).Limit(size + 1).Select()
	}
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	hasBefore, hasAfter := len(before) > size, len(after) > size
	before, after = before[:min(len(before), size)], after[:min(len(after), size)]
	list := make([]msql.Params, 0, len(before)+len(after)+1)
	for i := len(before) - 1; i >= 0; i-- {
		list = append(list, before[i])
	}
	if minId == 0 && maxId == 0 {
		list = append(list, context.answer)
	}
	list = append(list, after...)
	common.FmtOk(c, map[string]any{`list`: list, `has_more_before`: hasBefore, `has_more_after`: hasAfter})
}
