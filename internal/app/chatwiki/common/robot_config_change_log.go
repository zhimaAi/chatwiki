// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const TableRobotConfigChangeLog = `chat_ai_robot_config_change_log`

// change modules, one per handler
const (
	RobotChangeModuleRelationFlow    = `relation_workflow` // relation workflow
	RobotChangeModuleRelationLibrary = `relation_library`  // relation knowledge library
	RobotChangeModuleRobotSetting    = `robot_setting`     // model/prompt/search/chat setting
	RobotChangeModuleLangConfig      = `lang_config`       // multilingual config
	RobotChangeModuleChatVariable    = `chat_variable`     // chat variable add/edit/delete
)

// translateConfigChangeText returns the localized text for key under the given language.
// When the key is missing, beego/i18n returns the key itself, so we fall back to the given value.
func translateConfigChangeText(lang, key, fallback string) string {
	text := i18n.Show(lang, key)
	if text == key || text == `` {
		return fallback
	}
	return text
}

// getModuleSection returns the localized UI section name for a change module.
func getModuleSection(lang, module string) string {
	return translateConfigChangeText(lang, `cfglog_section_`+module, ``)
}

// GetConfigChangeModuleSection exposes the localized module section for other packages
// so section names are maintained only in the locale files.
func GetConfigChangeModuleSection(lang, module string) string {
	return getModuleSection(lang, module)
}

// noise fields ignored when diffing
var robotChangeIgnoreFields = map[string]struct{}{
	`update_time`:          {},
	`last_edit_ip`:         {},
	`last_edit_user_agent`: {},
	`draft_save_time`:      {},
	`draft_save_type`:      {},
	`sort_num`:             {},
	`is_top`:               {},
}

// RobotChangeDetailItem is one field-level change entry stored in change_detail.
type RobotChangeDetailItem struct {
	Section string `json:"section"`
	Field   string `json:"field"`
	Label   string `json:"label"`
	Before  string `json:"before"`
	After   string `json:"after"`
}

// getFieldLabel returns the localized label for a field, or field name if unknown.
func getFieldLabel(lang, field string) string {
	return translateConfigChangeText(lang, `cfglog_field_`+field, field)
}

// GetConfigChangeFieldLabel exposes the localized field label for other packages
// so field labels are maintained only in the locale files.
func GetConfigChangeFieldLabel(lang, field string) string {
	return getFieldLabel(lang, field)
}

// NormalizeValue normalizes a value for comparison
func NormalizeValue(v any) string {
	switch val := v.(type) {
	case string:
		// try to normalize JSON
		var js any
		if err := json.Unmarshal([]byte(val), &js); err == nil {
			normalized, _ := json.Marshal(js)
			return string(normalized)
		}
		return strings.TrimSpace(val)
	case float64, float32:
		// format float without unnecessary decimals
		f := cast.ToFloat64(val)
		if math.Mod(f, 1.0) == 0 {
			return strconv.FormatInt(int64(f), 10)
		}
		return strconv.FormatFloat(f, 'f', -1, 64)
	case int, int32, int64, uint, uint32, uint64:
		return cast.ToString(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return cast.ToString(val)
	}
}

// ValuesEqual compares two values with tolerance for different representations
func ValuesEqual(a, b any) bool {
	normA := NormalizeValue(a)
	normB := NormalizeValue(b)

	// check bool-like values
	boolValues := map[string]bool{
		"t":     true,
		"true":  true,
		"1":     true,
		"f":     false,
		"false": false,
		"0":     false,
	}

	aBool, aIsBool := boolValues[strings.ToLower(normA)]
	bBool, bIsBool := boolValues[strings.ToLower(normB)]

	if aIsBool && bIsBool {
		return aBool == bBool
	}

	// check float-like with tolerance
	aFloat, aErr := strconv.ParseFloat(normA, 64)
	bFloat, bErr := strconv.ParseFloat(normB, 64)

	if aErr == nil && bErr == nil {
		// tolerance of 1e-9
		return math.Abs(aFloat-bFloat) < 1e-9
	}

	return normA == normB
}

// DiffConfigMaps generates diff details between before and after maps
func DiffConfigMaps(lang string, before msql.Params, after msql.Datas, defaultSection string) []RobotChangeDetailItem {
	details := make([]RobotChangeDetailItem, 0)

	// check all fields in after
	for field, afterVal := range after {
		if _, ignore := robotChangeIgnoreFields[field]; ignore {
			continue
		}

		beforeVal, exists := before[field]
		if !exists {
			beforeVal = ""
		}

		if !ValuesEqual(beforeVal, afterVal) {
			details = append(details, RobotChangeDetailItem{
				Section: defaultSection,
				Field:   field,
				Label:   getFieldLabel(lang, field),
				Before:  cast.ToString(beforeVal),
				After:   cast.ToString(afterVal),
			})
		}
	}

	return details
}

// SaveRobotConfigChangeLog diffs before and after and saves the log
func SaveRobotConfigChangeLog(lang string, adminUserId, operUserId int, robotId int64, robotKey, module string, before msql.Params, after msql.Datas) {
	beforeCopy := make(msql.Params, len(before))
	for k, v := range before {
		beforeCopy[k] = v
	}
	afterCopy := make(msql.Datas, len(after))
	for k, v := range after {
		afterCopy[k] = v
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("SaveRobotConfigChangeLog panic: ", r)
			}
		}()

		section := getModuleSection(lang, module)
		details := DiffConfigMaps(lang, beforeCopy, afterCopy, section)

		if len(details) == 0 {
			return
		}

		// build before/after content maps with only changed fields
		beforeContent := make(map[string]any)
		afterContent := make(map[string]any)
		for _, d := range details {
			beforeContent[d.Field] = d.Before
			afterContent[d.Field] = d.After
		}

		insertRobotConfigChangeLog(lang, adminUserId, operUserId, robotId, robotKey, module,
			cast.ToInt(beforeCopy["application_type"]),
			tool.JsonEncodeNoError(beforeContent),
			tool.JsonEncodeNoError(afterContent),
			details,
		)
	}()
}

// SaveRobotConfigChangeLogDetail saves a log with pre-built details
func SaveRobotConfigChangeLogDetail(lang string, adminUserId, operUserId int, robotId int64, robotKey, module string, applicationType int, details []RobotChangeDetailItem) {
	if len(details) == 0 {
		return
	}

	// copy details
	detailsCopy := make([]RobotChangeDetailItem, len(details))
	copy(detailsCopy, details)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("SaveRobotConfigChangeLogDetail panic: ", r)
			}
		}()

		// build before/after content
		beforeContent := make(map[string]any)
		afterContent := make(map[string]any)
		for _, d := range detailsCopy {
			beforeContent[d.Field] = d.Before
			afterContent[d.Field] = d.After
		}

		insertRobotConfigChangeLog(lang, adminUserId, operUserId, robotId, robotKey, module,
			applicationType,
			tool.JsonEncodeNoError(beforeContent),
			tool.JsonEncodeNoError(afterContent),
			detailsCopy,
		)
	}()
}

// chatVariableLogFields lists the chat variable fields tracked in the change log, in display order.
var chatVariableLogFields = []string{
	`variable_type`,
	`variable_key`,
	`variable_name`,
	`max_input_length`,
	`default_value`,
	`must_input`,
	`options`,
}

// BuildChatVariableChangeDetails builds field-level change items for a chat variable.
// before is the snapshot before the change (nil/empty means a newly created variable);
// after is the snapshot after the change (nil means the variable was deleted).
// Fields whose before/after values are equal are skipped, so an unchanged save produces no detail.
func BuildChatVariableChangeDetails(lang string, before msql.Params, after map[string]any) []RobotChangeDetailItem {
	section := getModuleSection(lang, RobotChangeModuleChatVariable)
	details := make([]RobotChangeDetailItem, 0, len(chatVariableLogFields))
	for _, field := range chatVariableLogFields {
		beforeVal := before[field]
		afterVal := after[field]
		if ValuesEqual(beforeVal, afterVal) {
			continue
		}
		details = append(details, RobotChangeDetailItem{
			Section: section,
			Field:   field,
			Label:   getFieldLabel(lang, field),
			Before:  cast.ToString(beforeVal),
			After:   cast.ToString(afterVal),
		})
	}
	return details
}

// insertRobotConfigChangeLog is the shared low-level writer
func insertRobotConfigChangeLog(lang string, adminUserId, operUserId int, robotId int64, robotKey, module string, applicationType int, beforeContent, afterContent string, details []RobotChangeDetailItem) {
	operUserName := lib_define.UnknownUser
	if userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Where("id", cast.ToString(operUserId)).Field("user_name").Find(); err == nil && userInfo["user_name"] != "" {
		operUserName = userInfo["user_name"]
	}

	section := getModuleSection(lang, module)
	// use first detail section if available
	if len(details) > 0 && details[0].Section != "" {
		section = details[0].Section
	}

	_, err := msql.Model(TableRobotConfigChangeLog, define.Postgres).Insert(msql.Datas{
		"admin_user_id":    adminUserId,
		"robot_id":         robotId,
		"robot_key":        robotKey,
		"application_type": applicationType,
		"module":           module,
		"section":          section,
		"oper_user_id":     operUserId,
		"oper_user_name":   operUserName,
		"before_content":   beforeContent,
		"after_content":    afterContent,
		"change_detail":    tool.JsonEncodeNoError(details),
		"create_time":      tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
	}
}
