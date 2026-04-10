// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// Language keys used by robot multilingual configuration.
const (
	RobotLangCh = "zh-CN"
	RobotLangEn = "en-US"
	RobotLangJa = "ja"
	RobotLangTw = "zh-TW"
	RobotLangVi = "vi"
	RobotLangHk = "zh-HK"
	RobotLangTh = "th"
	RobotLangId = "id"
	RobotLangMs = "ms"
	RobotLangPt = "pt"
	RobotLangEs = "es"
	RobotLangKo = "ko"
	RobotLangFr = "fr"
	RobotLangRu = "ru"
	RobotLangTr = "tr"
)

// TableRobotMultilingualConfig is the storage table for robot multilingual settings.
const TableRobotMultilingualConfig = "chat_ai_robot_multilingual_config"

// RobotSupportedLangKeys defines all supported robot language keys in fixed order.
var RobotSupportedLangKeys = []string{
	RobotLangCh, RobotLangEn, RobotLangJa, RobotLangTw, RobotLangVi,
	RobotLangHk, RobotLangTh, RobotLangId, RobotLangMs, RobotLangPt,
	RobotLangEs, RobotLangKo, RobotLangFr, RobotLangRu, RobotLangTr,
}

type RobotMultilingualConfig struct {
	LangKey                 string `json:"lang_key"`
	Welcomes                string `json:"welcomes"`
	UnknownQuestionPrompt   string `json:"unknown_question_prompt"`
	TipsBeforeAnswerSwitch  any    `json:"tips_before_answer_switch"`
	TipsBeforeAnswerContent string `json:"tips_before_answer_content"`
	EnableCommonQuestion    any    `json:"enable_common_question"`
	CommonQuestionList      string `json:"common_question_list"`
}

type RobotThinkingTipsConfig struct {
	TipsBeforeAnswerContent string `json:"tips_before_answer_content"`
	TipsBeforeAnswerSwitch  bool   `json:"tips_before_answer_switch"`
}

// NormalizeRobotMultilingualConfigs normalizes imported multilingual configs before storage.
func NormalizeRobotMultilingualConfigs(list []RobotMultilingualConfig) []RobotMultilingualConfig {
	result := make([]RobotMultilingualConfig, 0, len(list))
	for _, one := range list {
		langKey := NormalizeRobotLangKey(one.LangKey)
		if langKey == `` {
			continue
		}
		result = append(result, RobotMultilingualConfig{
			LangKey:                 langKey,
			Welcomes:                one.Welcomes,
			UnknownQuestionPrompt:   one.UnknownQuestionPrompt,
			TipsBeforeAnswerSwitch:  NormalizeFlexibleBool(one.TipsBeforeAnswerSwitch),
			TipsBeforeAnswerContent: one.TipsBeforeAnswerContent,
			EnableCommonQuestion:    NormalizeFlexibleBool(one.EnableCommonQuestion),
			CommonQuestionList:      one.CommonQuestionList,
		})
	}
	return result
}

// NormalizeFlexibleBool accepts bool-like values from mixed JSON payloads such as true/false or "true"/"false".
func NormalizeFlexibleBool(val any) bool {
	if val == nil {
		return false
	}
	return cast.ToBool(val)
}

// NormalizeRobotLangKey converts various locale formats into normalized robot lang keys.
func NormalizeRobotLangKey(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "", "zh", "zh-cn", "zh_cn", "ch":
		return RobotLangCh
	case "en", "en-us", "en_us":
		return RobotLangEn
	case "ja", "ja-jp", "ja_jp", "jp":
		return RobotLangJa
	case "zh-tw", "zh_tw", "tw":
		return RobotLangTw
	case "vi", "vi-vn", "vi_vn":
		return RobotLangVi
	case "zh-hk", "zh_hk", "hk":
		return RobotLangHk
	case "th", "th-th", "th_th":
		return RobotLangTh
	case "id", "id-id", "id_id", "in":
		return RobotLangId
	case "ms", "ms-my", "ms_my", "my":
		return RobotLangMs
	case "pt", "pt-pt", "pt_pt", "pt-br", "pt_br":
		return RobotLangPt
	case "es", "es-es", "es_es":
		return RobotLangEs
	case "ko", "ko-kr", "ko_kr":
		return RobotLangKo
	case "fr", "fr-fr", "fr_fr":
		return RobotLangFr
	case "ru", "ru-ru", "ru_ru":
		return RobotLangRu
	case "tr", "tr-tr", "tr_tr":
		return RobotLangTr
	default:
		return RobotLangEn
	}
}

func findRobotMultilingualConfigByLang(robotId int, langKey string) (msql.Params, error) {
	if robotId <= 0 {
		return msql.Params{}, nil
	}
	langKey = NormalizeRobotLangKey(langKey)
	list, err := GetRobotMultilingualConfigList(robotId)
	if err != nil {
		return nil, err
	}
	for _, one := range list {
		if NormalizeRobotLangKey(one[`lang_key`]) == langKey {
			return one, nil
		}
	}
	return msql.Params{}, nil
}

// defaultRobotLangValues returns per-language default values for multilingual robot fields.
func defaultRobotLangValues(langKey string) (string, string, string, bool, bool, string) {
	switch langKey {
	case RobotLangCh:
		return i18n.Show(define.LangZhCn, `default_welcomes`),
			tool.JsonEncodeNoError(define.MenuJsonStruct{Content: lib_define.DefaultUnknownQuestionPromptContent, Question: []string{}}),
			i18n.Show(define.LangZhCn, `thinking_please_wait`), true,
			false,
			`[]`
	case RobotLangEn:
		return `{"content":"Hello, welcome to consult!","question":[]}`,
			`{"content":"Oops, I'm not sure about this question yet.","question":[]}`,
			`Thinking, please wait...`, true,
			false,
			`[]`
	default:
		return ``, `{"content":"","question":[]}`, ``, false, false, `[]`
	}
}

// BuildRobotMultilingualConfigsByLegacy builds full multilingual rows using legacy robot fields as source.
func BuildRobotMultilingualConfigsByLegacy(welcomes, unknownQuestionPrompt, tipsBeforeAnswerContent string, tipsBeforeAnswerSwitch bool, enableCommonQuestion bool, commonQuestionList string) []RobotMultilingualConfig {
	result := make([]RobotMultilingualConfig, 0, len(RobotSupportedLangKeys))
	for _, langKey := range RobotSupportedLangKeys {
		defWelcomes, defUnknownPrompt, defTipsBeforeAnswerContent, defTipsBeforeAnswerSwitch, defEnableCommonQuestion, defCommonQuestionList := defaultRobotLangValues(langKey)
		one := RobotMultilingualConfig{
			LangKey:                 langKey,
			Welcomes:                defWelcomes,
			UnknownQuestionPrompt:   defUnknownPrompt,
			TipsBeforeAnswerSwitch:  defTipsBeforeAnswerSwitch,
			TipsBeforeAnswerContent: defTipsBeforeAnswerContent,
			EnableCommonQuestion:    defEnableCommonQuestion,
			CommonQuestionList:      defCommonQuestionList,
		}
		// build from old data
		if langKey == RobotLangCh {
			one.Welcomes = welcomes
			one.UnknownQuestionPrompt = unknownQuestionPrompt
			one.TipsBeforeAnswerSwitch = tipsBeforeAnswerSwitch
			one.TipsBeforeAnswerContent = tipsBeforeAnswerContent
			one.EnableCommonQuestion = enableCommonQuestion
			one.CommonQuestionList = commonQuestionList
		}
		result = append(result, one)
	}
	return NormalizeRobotMultilingualConfigs(result)
}

// BuildRobotMultilingualConfigsFromRobot builds multilingual rows from robot legacy fields.
func BuildRobotMultilingualConfigsFromRobot(robot msql.Params) []RobotMultilingualConfig {
	defWelcomes, defUnknownPrompt, defTipsBeforeAnswerContent, defTipsBeforeAnswerSwitch, defEnableCommonQuestion, defCommonQuestionList :=
		defaultRobotLangValues(RobotLangCh)

	welcomes := strings.TrimSpace(robot[`welcomes`])
	if welcomes == `` {
		welcomes = defWelcomes
	}
	unknownQuestionPrompt := strings.TrimSpace(robot[`unknown_question_prompt`])
	if unknownQuestionPrompt == `` {
		unknownQuestionPrompt = defUnknownPrompt
	}
	tipsBeforeAnswerSwitch := cast.ToBool(robot[`tips_before_answer_switch`])
	if strings.TrimSpace(robot[`tips_before_answer_switch`]) == `` {
		tipsBeforeAnswerSwitch = defTipsBeforeAnswerSwitch
	}
	tipsBeforeAnswerContent := strings.TrimSpace(robot[`tips_before_answer_content`])
	if tipsBeforeAnswerContent == `` && tipsBeforeAnswerSwitch {
		tipsBeforeAnswerContent = defTipsBeforeAnswerContent
	}
	enableCommonQuestion := cast.ToBool(robot[`enable_common_question`])
	if strings.TrimSpace(robot[`enable_common_question`]) == `` {
		enableCommonQuestion = defEnableCommonQuestion
	}
	commonQuestionList := strings.TrimSpace(robot[`common_question_list`])
	if commonQuestionList == `` {
		commonQuestionList = defCommonQuestionList
	}

	return BuildRobotMultilingualConfigsByLegacy(
		welcomes,
		unknownQuestionPrompt,
		tipsBeforeAnswerContent,
		tipsBeforeAnswerSwitch,
		enableCommonQuestion,
		commonQuestionList,
	)
}

// SaveRobotMultilingualConfigs replaces all multilingual rows of a robot and inserts the new set.
func SaveRobotMultilingualConfigs(adminUserId, robotId int, list []RobotMultilingualConfig) error {
	if adminUserId <= 0 || robotId <= 0 {
		return nil
	}
	list = NormalizeRobotMultilingualConfigs(list)
	if len(list) == 0 {
		return nil
	}
	m := msql.Model(TableRobotMultilingualConfig, define.Postgres)
	if _, err := m.Where(`robot_id`, cast.ToString(robotId)).Delete(); err != nil {
		return err
	}
	now := tool.Time2Int()
	for _, one := range list {
		if _, err := m.Insert(msql.Datas{
			`admin_user_id`:              adminUserId,
			`robot_id`:                   robotId,
			`lang_key`:                   NormalizeRobotLangKey(one.LangKey),
			`welcomes`:                   one.Welcomes,
			`unknown_question_prompt`:    one.UnknownQuestionPrompt,
			`tips_before_answer_switch`:  NormalizeFlexibleBool(one.TipsBeforeAnswerSwitch),
			`tips_before_answer_content`: one.TipsBeforeAnswerContent,
			`enable_common_question`:     NormalizeFlexibleBool(one.EnableCommonQuestion),
			`common_question_list`:       one.CommonQuestionList,
			`create_time`:                now,
			`update_time`:                now,
		}); err != nil {
			logs.Error(`save robot multi lang error：%s`, err.Error())
			return err
		}
	}
	return nil
}

// SaveRobotMultilingualConfig upserts one multilingual row by robot_id + lang_key.
func SaveRobotMultilingualConfig(adminUserId, robotId int, one RobotMultilingualConfig) error {
	if adminUserId <= 0 || robotId <= 0 {
		return nil
	}
	langKey := NormalizeRobotLangKey(one.LangKey)
	if len(langKey) == 0 {
		return nil
	}
	now := tool.Time2Int()
	m := msql.Model(TableRobotMultilingualConfig, define.Postgres)
	exist, err := findRobotMultilingualConfigByLang(robotId, langKey)
	if err != nil {
		return err
	}
	data := msql.Datas{
		`admin_user_id`:              adminUserId,
		`robot_id`:                   robotId,
		`lang_key`:                   langKey,
		`welcomes`:                   one.Welcomes,
		`unknown_question_prompt`:    one.UnknownQuestionPrompt,
		`tips_before_answer_switch`:  NormalizeFlexibleBool(one.TipsBeforeAnswerSwitch),
		`tips_before_answer_content`: one.TipsBeforeAnswerContent,
		`enable_common_question`:     NormalizeFlexibleBool(one.EnableCommonQuestion),
		`common_question_list`:       one.CommonQuestionList,
		`update_time`:                now,
	}
	if len(exist) > 0 {
		_, err = m.
			Where(`id`, exist[`id`]).
			Update(data)
		return err
	}
	data[`create_time`] = now
	_, err = m.Insert(data)
	return err
}

// BuildLegacyRobotFieldsFromMultilingualConfig converts zh multilingual config back to legacy robot fields.
func BuildLegacyRobotFieldsFromMultilingualConfig(one RobotMultilingualConfig) (msql.Datas, bool) {
	if NormalizeRobotLangKey(one.LangKey) != RobotLangCh {
		return nil, false
	}
	return msql.Datas{
		`welcomes`:                   one.Welcomes,
		`unknown_question_prompt`:    one.UnknownQuestionPrompt,
		`tips_before_answer_switch`:  NormalizeFlexibleBool(one.TipsBeforeAnswerSwitch),
		`tips_before_answer_content`: one.TipsBeforeAnswerContent,
		`enable_common_question`:     NormalizeFlexibleBool(one.EnableCommonQuestion),
		`common_question_list`:       one.CommonQuestionList,
	}, true
}

// GetRobotMultilingualConfigList returns multilingual config rows of the specified robot.
func GetRobotMultilingualConfigList(robotId int) ([]msql.Params, error) {
	if robotId <= 0 {
		return []msql.Params{}, nil
	}
	list, err := msql.Model(TableRobotMultilingualConfig, define.Postgres).
		Where(`robot_id`, cast.ToString(robotId)).
		Order(`id asc`).
		Select()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetRobotMultilingualConfig returns one multilingual config row by robot and language key.
func GetRobotMultilingualConfig(robotId int, langKey string) (msql.Params, error) {
	return findRobotMultilingualConfigByLang(robotId, langKey)
}

func applyEmptyRobotMultilingualConfig(robot msql.Params, lang string) msql.Params {
	robot[`welcomes`] = ``
	robot[`unknown_question_prompt`] = tool.JsonEncodeNoError(define.MenuJsonStruct{Content: ``, Question: []string{}})
	robot[`common_question_list`] = `[]`
	robot[`tips_before_answer_content`] = ``
	robot[`tips_before_answer_switch`] = `false`
	robot[`enable_common_question`] = `false`
	robot[`lang_key`] = lang
	return robot
}

// ApplyRobotMultilingualConfig overlays robot legacy fields by language, with English fallback.
func ApplyRobotMultilingualConfig(robot msql.Params, lang string) msql.Params {
	if len(robot) == 0 {
		return robot
	}

	lang = NormalizeRobotLangKey(lang)
	if len(lang) == 0 {
		lang = RobotLangCh
	}
	if len(robot[`multi_lang_configs`]) == 0 {
		robotId := cast.ToInt(robot[`id`])
		if robotId > 0 {
			if list, err := GetRobotMultilingualConfigList(robotId); err == nil && len(list) > 0 {
				robot[`multi_lang_configs`] = tool.JsonEncodeNoError(list)
			}
		}
	}
	var list []msql.Params
	_ = tool.JsonDecodeUseNumber(robot[`multi_lang_configs`], &list)
	if len(list) == 0 {
		return applyEmptyRobotMultilingualConfig(robot, lang)
	}
	langMap := make(map[string]msql.Params)
	for _, one := range list {
		langMap[NormalizeRobotLangKey(one[`lang_key`])] = one
	}
	targetCfg, ok := langMap[lang]
	if !ok || len(targetCfg) == 0 {
		targetCfg = langMap[RobotLangEn]
		if len(targetCfg) == 0 {
			return applyEmptyRobotMultilingualConfig(robot, lang)
		}
	}
	pick := func(key string) string {
		if len(targetCfg[key]) > 0 {
			return targetCfg[key]
		}
		return ``
	}
	robot[`welcomes`] = pick(`welcomes`)
	robot[`unknown_question_prompt`] = pick(`unknown_question_prompt`)
	if len(robot[`unknown_question_prompt`]) == 0 {
		robot[`unknown_question_prompt`] = tool.JsonEncodeNoError(define.MenuJsonStruct{Content: ``, Question: []string{}})
	}
	robot[`common_question_list`] = pick(`common_question_list`)
	if len(robot[`common_question_list`]) == 0 {
		robot[`common_question_list`] = `[]`
	}
	robot[`tips_before_answer_content`] = pick(`tips_before_answer_content`)
	if len(targetCfg[`tips_before_answer_switch`]) > 0 {
		robot[`tips_before_answer_switch`] = targetCfg[`tips_before_answer_switch`]
	} else {
		robot[`tips_before_answer_switch`] = `false`
	}
	if len(targetCfg[`enable_common_question`]) > 0 {
		robot[`enable_common_question`] = targetCfg[`enable_common_question`]
	} else {
		robot[`enable_common_question`] = `false`
	}
	robot[`lang_key`] = lang
	return robot
}
