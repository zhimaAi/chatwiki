// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type MetaSearchCondition struct {
	Key   string `json:"key"`
	Type  int    `json:"type"`  // 0 string, 1 time, 2 number (same as define.LibraryMetaType*)
	Op    int    `json:"op"`    // define.MetaOp*
	Value string `json:"value"` // Can be empty (empty/not empty operators)
}

type metaSearchTarget int

const (
	metaSearchTargetFile      metaSearchTarget = 1 // chat_ai_library_file (filter by file_id)
	metaSearchTargetParagraph metaSearchTarget = 2 // chat_ai_library_file_data (filter by data_id=id)
)

// IsCustomMetaKey Custom metadata key fixed format: key_number
func IsCustomMetaKey(key string) bool {
	ok, _ := regexp.MatchString(`^key_\d+$`, key)
	return ok
}

func sqlQuote(s string) string {
	// Minimize escaping: double single quotes
	return strings.ReplaceAll(s, `'`, `''`)
}

// normalizeBuiltinMetaValue normalizes built-in metadata condition values before SQL building.
// It is applied after all placeholder/chat-variable replacements, so it only focuses on:
// - Source: map user-facing labels (CN/EN) to underlying doc_type int (as string)
// - Time:   accept human-readable time strings and convert to unix timestamp (seconds)
// - Group:  accept numeric group_id and map to group_name (string) used in SQL expressions
func normalizeBuiltinMetaValue(cond MetaSearchCondition, target metaSearchTarget) MetaSearchCondition {
	key := strings.TrimSpace(cond.Key)
	if !define.IsBuiltinMetaKey(key) {
		return cond
	}
	val := strings.TrimSpace(cond.Value)
	if val == "" {
		return cond
	}

	switch key {
	case define.BuiltinMetaKeySource:
		cond.Value = normalizeSourceValue(val)
	case define.BuiltinMetaKeyCreateTime, define.BuiltinMetaKeyUpdateTime:
		if cond.Type == define.LibraryMetaTypeTime {
			cond.Value = normalizeTimeValue(val)
		}
	}
	return cond
}

// normalizeSourceValue maps human readable source labels to internal doc_type values.
// Returns a stringified int when mapping is successful, otherwise returns original value.
// Uses i18n to get both Chinese and English values for matching.
func normalizeSourceValue(val string) string {
	v := strings.TrimSpace(val)
	if v == "" {
		return v
	}
	// Already numeric, treat as doc_type value.
	if ok, _ := regexp.MatchString(`^\d+$`, v); ok {
		return v
	}

	lower := strings.ToLower(v)

	// Build mapping from i18n values (both zh-CN and en-US) to doc_type
	// Get i18n values for source labels
	sourceManualQAZh := strings.ToLower(i18n.Show("zh-CN", "source_manual_qa"))
	sourceManualQAEn := strings.ToLower(i18n.Show("en-US", "source_manual_qa"))
	sourceImportQAZh := strings.ToLower(i18n.Show("zh-CN", "source_import_qa"))
	sourceImportQAEn := strings.ToLower(i18n.Show("en-US", "source_import_qa"))
	sourceImportFsZh := strings.ToLower(i18n.Show("zh-CN", "source_import_fs"))
	sourceImportFsEn := strings.ToLower(i18n.Show("en-US", "source_import_fs"))

	nameMap := map[string]int{
		// Local documents
		"local":     define.DocTypeLocal,
		"local_doc": define.DocTypeLocal,
		// Online documents
		"online":     define.DocTypeOnline,
		"online_doc": define.DocTypeOnline,
		// Custom documents (generic)
		"custom": define.DocTypeCustom,
		// DIY/Q&A type (self-built QA) - use i18n values
		"diy":            define.DocTypeDiy,
		sourceManualQAZh: define.DocTypeDiy,
		sourceManualQAEn: define.DocTypeDiy,
		// Official / imported QA - use i18n values
		"official":       define.DocTypeOfficial,
		"importqa":       define.DocTypeOfficial,
		sourceImportQAZh: define.DocTypeOfficial,
		sourceImportQAEn: define.DocTypeOfficial,
		// Feishu - use i18n values
		"feishu":         define.DocTypeFeishu,
		sourceImportFsZh: define.DocTypeFeishu,
		sourceImportFsEn: define.DocTypeFeishu,
	}

	if code, ok := nameMap[lower]; ok {
		return cast.ToString(code)
	}
	return v
}

// normalizeTimeValue parses common time string formats and converts them to unix timestamp (seconds).
// If parsing fails or value already looks like a numeric timestamp, the original value is returned.
func normalizeTimeValue(val string) string {
	v := strings.TrimSpace(val)
	if v == "" {
		return v
	}
	// Already a numeric timestamp.
	if ok, _ := regexp.MatchString(`^\d{1,20}$`, v); ok {
		return v
	}

	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, v, time.Local); err == nil {
			return cast.ToString(t.Unix())
		}
	}
	// Parsing failed, return original for upstream validation/logging.
	logs.Error(fmt.Sprintf("[meta_search] time value parse failed, use raw value: %s", v))
	return v
}

// normalizeGroupValue maps group_id (numeric) to group_name for built-in "group" metadata.
// For non-numeric values we assume user already passed group_name.
func normalizeGroupValue(val string, target metaSearchTarget) string {
	v := strings.TrimSpace(val)
	if v == "" {
		return v
	}
	// Non-numeric -> treat as group_name directly.
	if ok, _ := regexp.MatchString(`^\d+$`, v); !ok {
		return v
	}

	// group_type depends on target: file vs paragraph(QA)
	var groupType int
	switch target {
	case metaSearchTargetParagraph:
		groupType = define.LibraryGroupTypeQA
	default:
		groupType = define.LibraryGroupTypeFile
	}

	info, err := msql.Model(`chat_ai_library_group`, define.Postgres).
		Where(`id`, v).
		Where(`group_type`, cast.ToString(groupType)).
		Field(`group_name`).
		Find()
	if err != nil {
		logs.Error(fmt.Sprintf("[meta_search] query group_name by id failed, id=%s, err=%v", v, err))
		return v
	}
	if len(info) == 0 {
		return v
	}
	name := strings.TrimSpace(cast.ToString(info[`group_name`]))
	if name == "" {
		return v
	}
	return name
}

func buildMetaFieldExpr(cond MetaSearchCondition, target metaSearchTarget) (expr string, joinGroup bool, joinFile bool, numeric bool, err error) {
	key := strings.TrimSpace(cond.Key)
	if define.IsBuiltinMetaKey(key) {
		switch key {
		case define.BuiltinMetaKeySource:
			// doc_type is int, treated as text here (for string operators)
			if target == metaSearchTargetParagraph {
				// Paragraph source: directly map d.type (consistent with GetParagraphList display)
				// type=1 => source=4, type=2 => source=5
				return "CASE WHEN d.type=1 THEN '4' WHEN d.type=2 THEN '5' ELSE '' END", false, false, false, nil
			}
			return "COALESCE(f.doc_type::text,'')", false, false, false, nil
		case define.BuiltinMetaKeyGroup:
			// Group filtering by group name (consistent with frontend display)
			if target == metaSearchTargetParagraph {
				return "CASE WHEN d.group_id=0 THEN '' ELSE COALESCE(g.group_name,'') END", true, false, false, nil
			}
			return "CASE WHEN f.group_id=0 THEN '' ELSE COALESCE(g.group_name,'') END", true, false, false, nil
		case define.BuiltinMetaKeyCreateTime:
			if target == metaSearchTargetParagraph {
				return "d.create_time", false, false, true, nil
			}
			return "f.create_time", false, false, true, nil
		case define.BuiltinMetaKeyUpdateTime:
			if target == metaSearchTargetParagraph {
				return "d.update_time", false, false, true, nil
			}
			return "f.update_time", false, false, true, nil
		default:
			return "", false, false, false, errors.New("invalid builtin key")
		}
	}

	// Custom key fixed format: key_number
	if !IsCustomMetaKey(key) {
		return "", false, false, false, errors.New("invalid meta key")
	}

	// Custom meta: from jsonb metadata
	switch cond.Type {
	case define.LibraryMetaTypeString:
		if target == metaSearchTargetParagraph {
			return fmt.Sprintf("COALESCE(d.metadata->>'%s','')", key), false, false, false, nil
		}
		return fmt.Sprintf("COALESCE(f.metadata->>'%s','')", key), false, false, false, nil
	case define.LibraryMetaTypeNumber:
		// Safe cast: values are trimmed first; if not numeric format then NULL
		if target == metaSearchTargetParagraph {
			// Note: Postgres regex (~) is POSIX, does not support \d, must use [0-9]
			return fmt.Sprintf("CASE WHEN btrim(d.metadata->>'%s') ~ '^-?[0-9]+(\\\\.[0-9]+)?$' THEN btrim(d.metadata->>'%s')::numeric ELSE NULL END", key, key), false, false, true, nil
		}
		return fmt.Sprintf("CASE WHEN btrim(f.metadata->>'%s') ~ '^-?[0-9]+(\\\\.[0-9]+)?$' THEN btrim(f.metadata->>'%s')::numeric ELSE NULL END", key, key), false, false, true, nil
	case define.LibraryMetaTypeTime:
		// time uses integer timestamp (max 20 digits)
		if target == metaSearchTargetParagraph {
			return fmt.Sprintf("CASE WHEN btrim(d.metadata->>'%s') ~ '^[0-9]{1,20}$' THEN btrim(d.metadata->>'%s')::bigint ELSE NULL END", key, key), false, false, true, nil
		}
		return fmt.Sprintf("CASE WHEN btrim(f.metadata->>'%s') ~ '^[0-9]{1,20}$' THEN btrim(f.metadata->>'%s')::bigint ELSE NULL END", key, key), false, false, true, nil
	default:
		return "", false, false, false, errors.New("invalid meta type")
	}
}

func buildMetaConditionSQL(cond MetaSearchCondition, target metaSearchTarget) (sql string, needGroupJoin bool, needFileJoin bool, err error) {
	expr, joinGroup, joinFile, numeric, err := buildMetaFieldExpr(cond, target)
	if err != nil {
		return "", false, false, err
	}
	needGroupJoin = joinGroup
	needFileJoin = joinFile

	op := cond.Op
	val := strings.TrimSpace(cond.Value)
	valQ := sqlQuote(val)

	// Empty/not empty (for numbers: NULL or 0 is considered empty; for strings: empty string is considered empty)
	if op == define.MetaOpEmpty {
		if numeric {
			// For built-in create/update_time: 0 is considered empty; for custom numbers: NULL is considered empty
			if define.IsBuiltinMetaKey(strings.TrimSpace(cond.Key)) && (cond.Key == define.BuiltinMetaKeyCreateTime || cond.Key == define.BuiltinMetaKeyUpdateTime) {
				return fmt.Sprintf("(%s = 0)", expr), needGroupJoin, needFileJoin, nil
			}
			return fmt.Sprintf("(%s IS NULL)", expr), needGroupJoin, needFileJoin, nil
		}
		return fmt.Sprintf("(%s = '')", expr), needGroupJoin, needFileJoin, nil
	}
	if op == define.MetaOpNotEmpty {
		if numeric {
			if define.IsBuiltinMetaKey(strings.TrimSpace(cond.Key)) && (cond.Key == define.BuiltinMetaKeyCreateTime || cond.Key == define.BuiltinMetaKeyUpdateTime) {
				return fmt.Sprintf("(%s <> 0)", expr), needGroupJoin, needFileJoin, nil
			}
			return fmt.Sprintf("(%s IS NOT NULL)", expr), needGroupJoin, needFileJoin, nil
		}
		return fmt.Sprintf("(%s <> '')", expr), needGroupJoin, needFileJoin, nil
	}

	// Operators that require value
	if val == "" {
		return "", false, false, errors.New("value required")
	}

	// String operators
	if !numeric {
		switch op {
		case define.MetaOpIs:
			return fmt.Sprintf("(%s = '%s')", expr, valQ), needGroupJoin, needFileJoin, nil
		case define.MetaOpIsNot:
			return fmt.Sprintf("(%s <> '%s')", expr, valQ), needGroupJoin, needFileJoin, nil
		case define.MetaOpContains:
			return fmt.Sprintf("(%s ILIKE '%%%s%%')", expr, valQ), needGroupJoin, needFileJoin, nil
		case define.MetaOpNotContains:
			return fmt.Sprintf("(%s NOT ILIKE '%%%s%%')", expr, valQ), needGroupJoin, needFileJoin, nil
		default:
			return "", false, false, errors.New("invalid op for string")
		}
	}

	// Numeric/time operators
	// Here val must be a number (validated upstream, double-checked here)
	if ok, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, val); !ok {
		return "", false, false, errors.New("invalid numeric value")
	}
	switch op {
	case define.MetaOpIs, define.MetaOpEq:
		return fmt.Sprintf("(%s IS NOT NULL AND %s = %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpIsNot:
		// "Is not/Not equal" does not include null/invalid values; use MetaOpEmpty / MetaOpNotEmpty for empty checks
		return fmt.Sprintf("(%s IS NOT NULL AND %s <> %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpGt:
		return fmt.Sprintf("(%s IS NOT NULL AND %s > %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpLt:
		return fmt.Sprintf("(%s IS NOT NULL AND %s < %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpGte:
		return fmt.Sprintf("(%s IS NOT NULL AND %s >= %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpLte:
		return fmt.Sprintf("(%s IS NOT NULL AND %s <= %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	default:
		return "", false, false, errors.New("invalid op for numeric/time")
	}
}

func getAllowedIdSetByRobotMetaSearch(adminUserId int, libraryIds string, robot msql.Params, target metaSearchTarget) (map[string]struct{}, error) {
	if cast.ToInt(robot[`meta_search_switch`]) != define.MetaSearchSwitchOn {
		return nil, nil
	}
	if len(libraryIds) == 0 || !CheckIds(libraryIds) {
		return nil, nil
	}
	metaType := cast.ToInt(robot[`meta_search_type`])
	if metaType == 0 {
		metaType = define.MetaSearchTypeAnd
	}
	raw := strings.TrimSpace(robot[`meta_search_condition_list`])
	if raw == "" || raw == "{}" || raw == "null" {
		return nil, nil
	}
	conds := make([]MetaSearchCondition, 0)
	if err := tool.JsonDecode(raw, &conds); err != nil {
		return nil, err
	}
	if len(conds) == 0 {
		return nil, nil
	}

	parts := make([]string, 0, len(conds))
	needGroupJoin := false
	needFileJoin := false
	for idx, c := range conds {
		// Normalize built-in meta condition values (source/time/group) after all placeholder replacements.
		conds[idx] = normalizeBuiltinMetaValue(c, target)

		sql, joinGroup, joinFile, err := buildMetaConditionSQL(conds[idx], target)
		if err != nil {
			return nil, err
		}
		if joinGroup {
			needGroupJoin = true
		}
		if joinFile {
			needFileJoin = true
		}
		parts = append(parts, sql)
	}
	joiner := " AND "
	if metaType == define.MetaSearchTypeOr {
		joiner = " OR "
	}
	where := "(" + strings.Join(parts, joiner) + ")"

	switch target {
	case metaSearchTargetParagraph:
		m := msql.Model(`chat_ai_library_file_data`, define.Postgres).
			Alias(`d`).
			Where(`d.admin_user_id`, cast.ToString(adminUserId)).
			Where(`d.library_id`, `in`, libraryIds).
			Where(`d.delete_time`, `0`)
		if needFileJoin {
			// source requires file.doc_type
			m.Join(`chat_ai_library_file f`, `d.file_id=f.id AND f.delete_time=0`, `left`)
		}
		if needGroupJoin {
			// Group name filtering requires join (group_type=LibraryGroupTypeQA)
			m.Join(`chat_ai_library_group g`, fmt.Sprintf(`d.group_id=g.id AND g.library_id=d.library_id AND g.group_type=%d`, define.LibraryGroupTypeQA), `left`)
		}
		m.Where(where)

		rows, err := m.Field(`d.id`).Select()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return nil, err
		}
		set := make(map[string]struct{}, len(rows))
		for _, r := range rows {
			set[cast.ToString(r[`id`])] = struct{}{}
		}
		return set, nil
	default:
		m := msql.Model(`chat_ai_library_file`, define.Postgres).
			Alias(`f`).
			Where(`f.admin_user_id`, cast.ToString(adminUserId)).
			Where(`f.library_id`, `in`, libraryIds).
			Where(`f.delete_time`, `0`)
		if needGroupJoin {
			// Group name filtering requires join (group_type=LibraryGroupTypeFile)
			m.Join(`chat_ai_library_group g`, fmt.Sprintf(`f.group_id=g.id AND g.library_id=f.library_id AND g.group_type=%d`, define.LibraryGroupTypeFile), `left`)
		}
		m.Where(where)

		rows, err := m.Field(`f.id`).Select()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return nil, err
		}
		set := make(map[string]struct{}, len(rows))
		for _, r := range rows {
			set[cast.ToString(r[`id`])] = struct{}{}
		}
		return set, nil
	}
}

func filterParamsByMetaSet(list []msql.Params, allowedFile map[string]struct{}, allowedData map[string]struct{}) []msql.Params {
	// nil means no filtering
	if allowedFile == nil && allowedData == nil {
		return list
	}
	out := make([]msql.Params, 0, len(list))
	for _, one := range list {
		typ := cast.ToInt(one[`type`])
		if typ == define.ParagraphTypeNormal {
			// Normal (type=1): filter by file table
			if allowedFile == nil {
				out = append(out, one)
				continue
			}
			if _, ok := allowedFile[cast.ToString(one[`file_id`])]; ok {
				out = append(out, one)
			}
			continue
		}
		// QA/Others: filter by paragraph table (id=data_id)
		if allowedData == nil {
			out = append(out, one)
			continue
		}
		if _, ok := allowedData[cast.ToString(one[`id`])]; ok {
			out = append(out, one)
		}
	}
	return out
}

func listHasNormalType(lists ...[]msql.Params) bool {
	for _, l := range lists {
		for _, one := range l {
			if cast.ToInt(one[`type`]) == define.ParagraphTypeNormal {
				return true
			}
		}
	}
	return false
}

func listHasNonNormalType(lists ...[]msql.Params) bool {
	for _, l := range lists {
		for _, one := range l {
			if cast.ToInt(one[`type`]) != define.ParagraphTypeNormal {
				return true
			}
		}
	}
	return false
}

// ApplyRobotMetaSearchFilter Simplest usage: input any number of lists:
// - Normal paragraphs (type=1) filtered by chat_ai_library_file metadata (file_id dimension)
// - QA/Other paragraphs (type!=1) filtered by chat_ai_library_file_data metadata (data_id=id dimension)
// Returns original when filtering is not enabled.
func ApplyRobotMetaSearchFilter(adminUserId int, libraryIds string, robot msql.Params, lists ...[]msql.Params) ([][]msql.Params, error) {
	// Filtering not enabled: return as-is
	if cast.ToInt(robot[`meta_search_switch`]) != define.MetaSearchSwitchOn {
		out := make([][]msql.Params, 0, len(lists))
		for _, l := range lists {
			out = append(out, l)
		}
		return out, nil
	}

	var allowedFile map[string]struct{}
	var allowedData map[string]struct{}
	var err error

	if listHasNormalType(lists...) {
		allowedFile, err = getAllowedIdSetByRobotMetaSearch(adminUserId, libraryIds, robot, metaSearchTargetFile)
		if err != nil {
			return nil, err
		}
	}
	if listHasNonNormalType(lists...) {
		allowedData, err = getAllowedIdSetByRobotMetaSearch(adminUserId, libraryIds, robot, metaSearchTargetParagraph)
		if err != nil {
			return nil, err
		}
	}

	// No filtering: return as-is
	if allowedFile == nil && allowedData == nil {
		out := make([][]msql.Params, 0, len(lists))
		for _, l := range lists {
			out = append(out, l)
		}
		return out, nil
	}

	out := make([][]msql.Params, 0, len(lists))
	for _, l := range lists {
		out = append(out, filterParamsByMetaSet(l, allowedFile, allowedData))
	}
	return out, nil
}
