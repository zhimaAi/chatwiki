// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type MetaSearchCondition struct {
	Key   string `json:"key"`
	Type  int    `json:"type"`  // 0 string,1 time,2 number（同 define.LibraryMetaType*）
	Op    int    `json:"op"`    // define.MetaOp*
	Value string `json:"value"` // 可为空（为空/不为空操作符）
}

type metaSearchTarget int

const (
	metaSearchTargetFile      metaSearchTarget = 1 // chat_ai_library_file（按 file_id 过滤）
	metaSearchTargetParagraph metaSearchTarget = 2 // chat_ai_library_file_data（按 data_id=id 过滤）
)

// IsCustomMetaKey 自定义元数据 key 固定格式：key_数字
func IsCustomMetaKey(key string) bool {
	ok, _ := regexp.MatchString(`^key_\d+$`, key)
	return ok
}

func sqlQuote(s string) string {
	// 最小化转义：单引号翻倍
	return strings.ReplaceAll(s, `'`, `''`)
}

func buildMetaFieldExpr(cond MetaSearchCondition, target metaSearchTarget) (expr string, joinGroup bool, joinFile bool, numeric bool, err error) {
	key := strings.TrimSpace(cond.Key)
	if define.IsBuiltinMetaKey(key) {
		switch key {
		case define.BuiltinMetaKeySource:
			// doc_type 是 int，这里按 text 处理（便于 string 操作符）
			if target == metaSearchTargetParagraph {
				// 段落来源：直接用 d.type 映射（与 GetParagraphList 展示一致）
				// type=1 => source=4，type=2 => source=5
				return "CASE WHEN d.type=1 THEN '4' WHEN d.type=2 THEN '5' ELSE '' END", false, false, false, nil
			}
			return "COALESCE(f.doc_type::text,'')", false, false, false, nil
		case define.BuiltinMetaKeyGroup:
			// group 过滤按分组名（与前端展示一致）
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

	// 自定义 key 格式固定：key_数字
	if !IsCustomMetaKey(key) {
		return "", false, false, false, errors.New("invalid meta key")
	}

	// 自定义 meta：来自 jsonb metadata
	switch cond.Type {
	case define.LibraryMetaTypeString:
		if target == metaSearchTargetParagraph {
			return fmt.Sprintf("COALESCE(d.metadata->>'%s','')", key), false, false, false, nil
		}
		return fmt.Sprintf("COALESCE(f.metadata->>'%s','')", key), false, false, false, nil
	case define.LibraryMetaTypeNumber:
		// 安全 cast：值统一先 btrim 去掉空格；不满足数字格式则 NULL
		if target == metaSearchTargetParagraph {
			// 注意：Postgres 正则（~）是 POSIX，不支持 \\d，必须用 [0-9]
			return fmt.Sprintf("CASE WHEN btrim(d.metadata->>'%s') ~ '^-?[0-9]+(\\\\.[0-9]+)?$' THEN btrim(d.metadata->>'%s')::numeric ELSE NULL END", key, key), false, false, true, nil
		}
		return fmt.Sprintf("CASE WHEN btrim(f.metadata->>'%s') ~ '^-?[0-9]+(\\\\.[0-9]+)?$' THEN btrim(f.metadata->>'%s')::numeric ELSE NULL END", key, key), false, false, true, nil
	case define.LibraryMetaTypeTime:
		// time 用整数时间戳（最多 20 位）
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

	// 空/不空（对数值：NULL 或 0 视为空；对字符串：空串视为空）
	if op == define.MetaOpEmpty {
		if numeric {
			// 对内置 create/update_time：0 视为空；对自定义数值：NULL 视为空
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

	// 需要 value 的操作符
	if val == "" {
		return "", false, false, errors.New("value required")
	}

	// string ops
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

	// numeric/time ops
	// 这里 val 必须是数字（上游保存已校验，这里再兜底）
	if ok, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, val); !ok {
		return "", false, false, errors.New("invalid numeric value")
	}
	switch op {
	case define.MetaOpIs, define.MetaOpEq:
		return fmt.Sprintf("(%s IS NOT NULL AND %s = %s)", expr, expr, val), needGroupJoin, needFileJoin, nil
	case define.MetaOpIsNot:
		// “不是/不等于”不包含空值/非法值；空/不空请用 MetaOpEmpty / MetaOpNotEmpty
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
	for _, c := range conds {
		sql, joinGroup, joinFile, err := buildMetaConditionSQL(c, target)
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
			// source 需要 file.doc_type
			m.Join(`chat_ai_library_file f`, `d.file_id=f.id AND f.delete_time=0`, `left`)
		}
		if needGroupJoin {
			// 分组名过滤需要 join（group_type=LibraryGroupTypeQA）
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
			// 分组名过滤需要 join（group_type=LibraryGroupTypeFile）
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
	// nil 表示不做过滤
	if allowedFile == nil && allowedData == nil {
		return list
	}
	out := make([]msql.Params, 0, len(list))
	for _, one := range list {
		typ := cast.ToInt(one[`type`])
		if typ == define.ParagraphTypeNormal {
			// 普通（type=1）：按文件表过滤
			if allowedFile == nil {
				out = append(out, one)
				continue
			}
			if _, ok := allowedFile[cast.ToString(one[`file_id`])]; ok {
				out = append(out, one)
			}
			continue
		}
		// QA/其它：按段落表过滤（id=data_id）
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

// ApplyRobotMetaSearchFilter 最简用法：输入任意数量列表：
// - 普通段(type=1)按 chat_ai_library_file 元数据过滤（file_id 维度）
// - QA/其它段(type!=1)按 chat_ai_library_file_data 元数据过滤（data_id=id 维度）
// 未开启过滤时原样返回。
func ApplyRobotMetaSearchFilter(adminUserId int, libraryIds string, robot msql.Params, lists ...[]msql.Params) ([][]msql.Params, error) {
	// 未开启过滤：原样返回
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

	// 无过滤：原样返回
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
