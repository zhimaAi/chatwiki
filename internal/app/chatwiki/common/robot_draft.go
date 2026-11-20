// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
)

// RobotNonDraftFields take effect immediately and are never applied through the draft.
var RobotNonDraftFields = []string{
	`robot_name`, `robot_intro`, `robot_avatar`, `group_id`, `en_name`,
	`external_config_h5`, `external_config_pc`,
	`is_top`, `sort_num`, `update_time`,
}

type RobotDraftData struct {
	Robot       msql.Params               `json:"robot"`
	LangConfigs []RobotMultilingualConfig `json:"lang_configs"`
}

func newRobotDraftData() *RobotDraftData {
	return &RobotDraftData{Robot: msql.Params{}, LangConfigs: []RobotMultilingualConfig{}}
}

func getRobotDraftRow(robotId int) (msql.Params, error) {
	if robotId <= 0 {
		return msql.Params{}, nil
	}
	return msql.Model(define.TableChatAiRobotDraft, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Find()
}

// GetRobotDraft returns the draft payload and its update time, nil means no draft.
func GetRobotDraft(robotId int) (*RobotDraftData, int, error) {
	row, err := getRobotDraftRow(robotId)
	if err != nil || len(row) == 0 {
		return nil, 0, err
	}
	raw := strings.TrimSpace(row[`draft_data`])
	if len(raw) == 0 {
		return newRobotDraftData(), cast.ToInt(row[`update_time`]), nil
	}
	var payload struct {
		Robot       map[string]any            `json:"robot"`
		LangConfigs []RobotMultilingualConfig `json:"lang_configs"`
	}
	if err = tool.JsonDecodeUseNumber(raw, &payload); err != nil {
		return nil, 0, err
	}
	data := newRobotDraftData()
	for key, val := range payload.Robot {
		data.Robot[key] = cast.ToString(val)
	}
	data.LangConfigs = NormalizeRobotMultilingualConfigs(payload.LangConfigs)
	return data, cast.ToInt(row[`update_time`]), nil
}

// SaveRobotDraftData writes the complete draft payload, keeping a single row per robot.
func SaveRobotDraftData(adminUserId, robotId int, data *RobotDraftData) error {
	if adminUserId <= 0 || robotId <= 0 || data == nil {
		return nil
	}
	if data.Robot == nil {
		data.Robot = msql.Params{}
	}
	if data.LangConfigs == nil {
		data.LangConfigs = []RobotMultilingualConfig{}
	}
	now := tool.Time2Int()
	row, err := getRobotDraftRow(robotId)
	if err != nil {
		return err
	}
	m := msql.Model(define.TableChatAiRobotDraft, define.Postgres)
	if len(row) > 0 {
		_, err = m.Where(`robot_id`, cast.ToString(robotId)).Update(msql.Datas{
			`admin_user_id`: adminUserId,
			`draft_data`:    tool.JsonEncodeNoError(data),
			`update_time`:   now,
		})
		return err
	}
	_, err = m.Insert(msql.Datas{
		`admin_user_id`: adminUserId,
		`robot_id`:      robotId,
		`draft_data`:    tool.JsonEncodeNoError(data),
		`create_time`:   now,
		`update_time`:   now,
	})
	return err
}

// SaveRobotDraftRobotConfig updates only the robot configuration part of the draft.
func SaveRobotDraftRobotConfig(adminUserId, robotId int, robot msql.Params) error {
	draft, _, err := GetRobotDraft(robotId)
	if err != nil {
		return err
	}
	if draft == nil {
		draft = newRobotDraftData()
	}
	draft.Robot = robot
	return SaveRobotDraftData(adminUserId, robotId, draft)
}

// SaveRobotDraftLangConfigs merges multilingual configs into the draft by lang_key.
func SaveRobotDraftLangConfigs(adminUserId, robotId int, configs []RobotMultilingualConfig) error {
	draft, _, err := GetRobotDraft(robotId)
	if err != nil {
		return err
	}
	if draft == nil {
		draft = newRobotDraftData()
	}
	if len(draft.LangConfigs) == 0 {
		published, err := GetRobotMultilingualConfigList(robotId)
		if err != nil {
			return err
		}
		for _, one := range published {
			draft.LangConfigs = append(draft.LangConfigs, robotMultilingualConfigFromParams(one))
		}
	}
	langMap := make(map[string]RobotMultilingualConfig, len(draft.LangConfigs)+len(configs))
	order := make([]string, 0, len(draft.LangConfigs)+len(configs))
	for _, one := range draft.LangConfigs {
		langKey := NormalizeRobotLangKey(one.LangKey)
		if len(langKey) == 0 {
			continue
		}
		if _, ok := langMap[langKey]; !ok {
			order = append(order, langKey)
		}
		langMap[langKey] = one
	}
	for _, one := range NormalizeRobotMultilingualConfigs(configs) {
		if _, ok := langMap[one.LangKey]; !ok {
			order = append(order, one.LangKey)
		}
		langMap[one.LangKey] = one
	}
	draft.LangConfigs = make([]RobotMultilingualConfig, 0, len(order))
	for _, langKey := range order {
		draft.LangConfigs = append(draft.LangConfigs, langMap[langKey])
	}
	return SaveRobotDraftData(adminUserId, robotId, draft)
}

func robotMultilingualConfigFromParams(one msql.Params) RobotMultilingualConfig {
	return RobotMultilingualConfig{
		LangKey:                 one[`lang_key`],
		Welcomes:                one[`welcomes`],
		UnknownQuestionPrompt:   one[`unknown_question_prompt`],
		TipsBeforeAnswerSwitch:  cast.ToBool(one[`tips_before_answer_switch`]),
		TipsBeforeAnswerContent: one[`tips_before_answer_content`],
		EnableCommonQuestion:    cast.ToBool(one[`enable_common_question`]),
		CommonQuestionList:      one[`common_question_list`],
	}
}

func DeleteRobotDraft(robotId int) error {
	if robotId <= 0 {
		return nil
	}
	_, err := msql.Model(define.TableChatAiRobotDraft, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Delete()
	return err
}

// RobotInfoWithDraft returns robot info with draft fields for admin responses.
func RobotInfoWithDraft(info msql.Params) map[string]any {
	result := make(map[string]any, len(info)+4)
	for key, val := range info {
		result[key] = val
	}
	if _, ok := result[`published_time`]; !ok {
		result[`published_time`] = `0`
	}
	result[`has_draft`] = false
	result[`draft_update_time`] = 0
	result[`draft_data`] = nil
	draft, updateTime, err := GetRobotDraft(cast.ToInt(info[`id`]))
	if err != nil || draft == nil {
		return result
	}
	if len(draft.Robot) == 0 && len(draft.LangConfigs) == 0 {
		return result
	}
	result[`has_draft`] = true
	result[`draft_update_time`] = updateTime
	result[`draft_data`] = draft
	return result
}

// ApplyRobotDraftForChat overlays the draft config for the robot owner's chat preview.
func ApplyRobotDraftForChat(c *gin.Context, robot msql.Params) (msql.Params, error) {
	if cast.ToInt(robot[`admin_user_id`]) <= 0 || cast.ToInt(robot[`admin_user_id`]) != getDraftAdminUserId(c) {
		return nil, errors.New(i18n.Show(GetLang(c), `auth_no_permission`))
	}
	draft, _, err := GetRobotDraft(cast.ToInt(robot[`id`]))
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return robot, nil
	}
	result := make(msql.Params, len(robot)+len(draft.Robot)+1)
	for key, val := range robot {
		result[key] = val
	}
	immutableKeys := append([]string{`id`, `robot_key`, `admin_user_id`, `create_time`, `creator`, `published_time`,
		`work_flow_ids`, `library_ids`, `form_ids`, `default_library_id`}, RobotNonDraftFields...)
	for key, val := range draft.Robot {
		if tool.InArrayString(key, immutableKeys) {
			continue
		}
		result[key] = val
	}
	if len(draft.LangConfigs) > 0 {
		result[`multi_lang_configs`] = tool.JsonEncodeNoError(draft.LangConfigs)
	}
	return result, nil
}

func getDraftAdminUserId(c *gin.Context) int {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	if len(token) == 0 {
		token = c.PostForm(`token`)
	}
	data, err := ParseToken(token)
	if err != nil {
		return 0
	}
	userId := cast.ToInt(data[`user_id`])
	if userId <= 0 {
		return 0
	}
	if cast.ToInt(data[`parent_id`]) <= 0 {
		return userId
	}
	return cast.ToInt(data[`parent_id`])
}
