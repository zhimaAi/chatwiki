// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"encoding/json"
	"errors"
	"math"
	"mime/multipart"
	neturl "net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/lib/pq"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetGoodsLibGroupList(lang string, adminUserId int) (map[string]any, error) {
	list, err := getGoodsLibGroups(adminUserId)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
	countList, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`group_id,count(*) AS goods_count`).
		Group(`group_id`).
		Select()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	countMap := make(map[int64]int, len(countList))
	total := 0
	for _, item := range countList {
		groupId := cast.ToInt64(item[`group_id`])
		count := cast.ToInt(item[`goods_count`])
		countMap[groupId] = count
		total += count
	}
	nodes := make([]*define.GoodsLibGroup, 0, len(list))
	for _, item := range list {
		groupId := cast.ToInt64(item[`id`])
		nodes = append(nodes, &define.GoodsLibGroup{
			ID:         groupId,
			ParentID:   cast.ToInt64(item[`parent_id`]),
			GroupName:  item[`group_name`],
			Level:      cast.ToInt(item[`level`]),
			Sort:       cast.ToInt(item[`sort`]),
			GoodsCount: countMap[groupId],
			Children:   make([]*define.GoodsLibGroup, 0),
		})
	}
	tree := buildGoodsLibGroupTree(nodes, 0)
	for _, node := range tree {
		calculateGoodsLibGroupTotal(node)
	}
	return map[string]any{
		`list`:            tree,
		`total`:           total,
		`ungrouped_count`: countMap[0],
	}, nil
}

func getGoodsLibGroups(adminUserId int) ([]msql.Params, error) {
	m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
	list, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,parent_id,group_name,level,sort`).
		Order(`level asc,parent_id asc,sort asc,id asc`).
		Select()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	return list, err
}

func buildGoodsLibGroupTree(list []*define.GoodsLibGroup, parentId int64) []*define.GoodsLibGroup {
	tree := make([]*define.GoodsLibGroup, 0)
	for _, node := range list {
		if node.ParentID != parentId {
			continue
		}
		node.Children = buildGoodsLibGroupTree(list, node.ID)
		tree = append(tree, node)
	}
	return tree
}

func calculateGoodsLibGroupTotal(node *define.GoodsLibGroup) int {
	total := node.GoodsCount
	for _, child := range node.Children {
		total += calculateGoodsLibGroupTotal(child)
	}
	node.TotalGoodsCount = total
	return total
}

func SaveGoodsLibGroup(lang string, adminUserId int, params define.GoodsLibSaveGroupParams) (int64, error) {
	params.GroupName = strings.TrimSpace(params.GroupName)
	if params.ID < 0 || params.ParentID < 0 || len(params.GroupName) == 0 ||
		utf8.RuneCountInString(params.GroupName) > define.GoodsLibGroupNameMaxLength {
		return 0, errors.New(i18n.Show(lang, `param_invalid`, `group`))
	}
	list, err := getGoodsLibGroups(adminUserId)
	if err != nil {
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	groupMap := make(map[int64]msql.Params, len(list))
	for _, item := range list {
		groupMap[cast.ToInt64(item[`id`])] = item
	}
	currentLevel := 0
	currentParentId := int64(0)
	if params.ID > 0 {
		current, ok := groupMap[params.ID]
		if !ok {
			return 0, errors.New(i18n.Show(lang, `no_data`))
		}
		currentLevel = cast.ToInt(current[`level`])
		currentParentId = cast.ToInt64(current[`parent_id`])
	}
	parentLevel := 0
	if params.ParentID > 0 {
		parent, ok := groupMap[params.ParentID]
		if !ok {
			return 0, errors.New(i18n.Show(lang, `param_invalid`, `parent_id`))
		}
		if params.ParentID == params.ID || goodsLibGroupHasAncestor(groupMap, params.ParentID, params.ID) {
			return 0, errors.New(i18n.Show(lang, `param_invalid`, `parent_id`))
		}
		parentLevel = cast.ToInt(parent[`level`])
	}
	newLevel := parentLevel + 1
	if newLevel > define.GoodsLibMaxGroupLevel {
		return 0, errors.New(i18n.Show(lang, `goods_group_level_max`, define.GoodsLibMaxGroupLevel))
	}
	maxChildDepth := 0
	if params.ID > 0 {
		maxChildDepth = goodsLibGroupMaxChildDepth(list, params.ID, 0, map[int64]bool{})
	}
	if newLevel+maxChildDepth > define.GoodsLibMaxGroupLevel {
		return 0, errors.New(i18n.Show(lang, `goods_group_level_max`, define.GoodsLibMaxGroupLevel))
	}
	now := tool.Time2Int()
	if params.ID == 0 {
		m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
		maxSort, err := m.
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`parent_id`, cast.ToString(params.ParentID)).
			Max(`sort`)
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return 0, errors.New(i18n.Show(lang, `sys_err`))
		}
		id, err := m.Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`parent_id`:     params.ParentID,
			`group_name`:    params.GroupName,
			`level`:         newLevel,
			`sort`:          cast.ToInt(maxSort) + 1,
			`create_time`:   now,
			`update_time`:   now,
		}, `id`)
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return 0, errors.New(i18n.Show(lang, `sys_err`))
		}
		return id, nil
	}
	m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
	if err = m.Begin(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	updateData := msql.Datas{
		`parent_id`:   params.ParentID,
		`group_name`:  params.GroupName,
		`level`:       newLevel,
		`update_time`: now,
	}
	if currentParentId != params.ParentID {
		maxSort, queryErr := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`parent_id`, cast.ToString(params.ParentID)).
			Max(`sort`)
		if queryErr != nil {
			_ = m.Rollback()
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), queryErr.Error())
			return 0, errors.New(i18n.Show(lang, `sys_err`))
		}
		updateData[`sort`] = cast.ToInt(maxSort) + 1
	}
	if _, err = m.Where(`id`, cast.ToString(params.ID)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Update(updateData); err != nil {
		_ = m.Rollback()
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	levelOffset := newLevel - currentLevel
	if levelOffset != 0 {
		descendantIds := goodsLibGroupDescendantIds(list, params.ID)
		for _, descendantId := range descendantIds {
			info := groupMap[descendantId]
			if _, err = m.Where(`id`, cast.ToString(descendantId)).
				Where(`admin_user_id`, cast.ToString(adminUserId)).
				Update(msql.Datas{
					`level`:       cast.ToInt(info[`level`]) + levelOffset,
					`update_time`: now,
				}); err != nil {
				_ = m.Rollback()
				logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
				return 0, errors.New(i18n.Show(lang, `sys_err`))
			}
		}
	}
	if err = m.Commit(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return params.ID, nil
}

func goodsLibGroupHasAncestor(groupMap map[int64]msql.Params, groupId, ancestorId int64) bool {
	if ancestorId <= 0 {
		return false
	}
	visited := make(map[int64]bool)
	for groupId > 0 && !visited[groupId] {
		visited[groupId] = true
		info, ok := groupMap[groupId]
		if !ok {
			return false
		}
		parentId := cast.ToInt64(info[`parent_id`])
		if parentId == ancestorId {
			return true
		}
		groupId = parentId
	}
	return false
}

func goodsLibGroupMaxChildDepth(list []msql.Params, parentId int64, depth int, visited map[int64]bool) int {
	if visited[parentId] {
		return depth
	}
	visited[parentId] = true
	maxDepth := depth
	for _, item := range list {
		if cast.ToInt64(item[`parent_id`]) != parentId {
			continue
		}
		childDepth := goodsLibGroupMaxChildDepth(list, cast.ToInt64(item[`id`]), depth+1, visited)
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	delete(visited, parentId)
	return maxDepth
}

func goodsLibGroupDescendantIds(list []msql.Params, parentId int64) []int64 {
	result := make([]int64, 0)
	visited := map[int64]bool{parentId: true}
	var findChildren func(int64)
	findChildren = func(id int64) {
		for _, item := range list {
			if cast.ToInt64(item[`parent_id`]) != id {
				continue
			}
			childId := cast.ToInt64(item[`id`])
			if visited[childId] {
				continue
			}
			visited[childId] = true
			result = append(result, childId)
			findChildren(childId)
		}
	}
	findChildren(parentId)
	return result
}

func DeleteGoodsLibGroup(lang string, adminUserId int, id int64) error {
	if id <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `id`))
	}
	list, err := getGoodsLibGroups(adminUserId)
	if err != nil {
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	found := false
	for _, item := range list {
		if cast.ToInt64(item[`id`]) == id {
			found = true
			break
		}
	}
	if !found {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	ids := append([]int64{id}, goodsLibGroupDescendantIds(list, id)...)
	values := make([]any, 0, len(ids))
	for _, groupId := range ids {
		values = append(values, groupId)
	}
	m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
	if err = m.Begin(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	now := tool.Time2Int()
	if _, err = m.Table(define.TableGoodsLibLibrary).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		WhereIn(`group_id`, values...).
		Update(msql.Datas{`group_id`: 0, `update_time`: now}); err != nil {
		_ = m.Rollback()
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if _, err = m.Table(define.TableGoodsLibGroup).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		WhereIn(`id`, values...).
		Delete(); err != nil {
		_ = m.Rollback()
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if err = m.Commit(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	return nil
}

func SortGoodsLibGroup(lang string, adminUserId int, items []define.GoodsLibGroupSortItem) error {
	if len(items) == 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `data`))
	}
	ids := make([]any, 0, len(items))
	seen := make(map[int64]bool, len(items))
	for _, item := range items {
		if item.ID <= 0 || item.Sort < 0 || seen[item.ID] {
			return errors.New(i18n.Show(lang, `param_invalid`, `data`))
		}
		seen[item.ID] = true
		ids = append(ids, item.ID)
	}
	m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
	count, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		WhereIn(`id`, ids...).
		Count()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if count != len(items) {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	if err = m.Begin(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	now := tool.Time2Int()
	for _, item := range items {
		if _, err = m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(item.ID)).
			Update(msql.Datas{`sort`: item.Sort, `update_time`: now}); err != nil {
			_ = m.Rollback()
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	if err = m.Commit(); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	return nil
}

func GetGoodsLibLibraryList(lang string, adminUserId int, filter define.GoodsLibListFilter) ([]map[string]any, int, error) {
	m := buildGoodsLibLibraryQuery(adminUserId, filter)
	list, total, err := m.Field(goodsLibLibrarySelectFields()).Order(`id desc`).Paginate(filter.Page, filter.Size)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	groups, err := getGoodsLibGroups(adminUserId)
	if err != nil {
		return nil, 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return formatGoodsLibLibraryList(lang, list, groups), total, nil
}

func ValidateGoodsLibraryListReq(filter define.GoodsLibListFilter) bool {
	return ValidateGoodsLibraryFilter(filter) && filter.Page >= 1 && filter.Size >= 1
}

func ValidateGoodsLibraryFilter(filter define.GoodsLibListFilter) bool {
	return filter.GroupID >= -1 &&
		(filter.SwitchStatus == -1 || filter.SwitchStatus == define.GoodsLibSwitchOff || filter.SwitchStatus == define.GoodsLibSwitchOn) &&
		utf8.RuneCountInString(strings.TrimSpace(filter.Keyword)) <= define.GoodsLibBaseInfoMaxLength
}

func GetGoodsLibLibraryInfo(lang string, adminUserId int, id int64) (map[string]any, error) {
	if id <= 0 {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `id`))
	}
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
	info, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Field(goodsLibLibrarySelectFields()).
		Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	groups, err := getGoodsLibGroups(adminUserId)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	groupNames := buildGoodsLibGroupNames(lang, groups)
	return formatGoodsLibLibraryItem(info, groupNames), nil
}

func buildGoodsLibLibraryQuery(adminUserId int, filter define.GoodsLibListFilter) *msql.Builder {
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	if len(filter.GroupIDs) > 0 && CheckIds(filter.GroupIDs) {
		m.Where(`group_id`, `in`, filter.GroupIDs)
	}
	if filter.GroupID >= 0 {
		m.Where(`group_id`, cast.ToString(filter.GroupID))
	}
	if filter.SwitchStatus == define.GoodsLibSwitchOff || filter.SwitchStatus == define.GoodsLibSwitchOn {
		m.Where(`switch_status`, cast.ToString(filter.SwitchStatus))
	}
	if filter.Keyword = strings.TrimSpace(filter.Keyword); len(filter.Keyword) > 0 {
		pattern := `%` + escapeGoodsLibLike(filter.Keyword) + `%`
		m.WhereRaw(`(goods_id ILIKE $1 ESCAPE '\' OR goods_name ILIKE $2 ESCAPE '\')`, pattern, pattern)
	}
	return m
}

func escapeGoodsLibLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	return strings.ReplaceAll(value, `_`, `\_`)
}

func goodsLibLibrarySelectFields() string {
	return `id,group_id,goods_id,goods_name,category,brand,price,stock,link,images,description,qa,custom_info,switch_status,create_time,update_time`
}

func formatGoodsLibLibraryList(lang string, list []msql.Params, groups []msql.Params) []map[string]any {
	groupNames := buildGoodsLibGroupNames(lang, groups)
	result := make([]map[string]any, 0, len(list))
	for _, item := range list {
		result = append(result, formatGoodsLibLibraryItem(item, groupNames))
	}
	return result
}

func formatGoodsLibLibraryItem(item msql.Params, groupNames map[int64]string) map[string]any {
	images := make([]string, 0)
	_ = json.Unmarshal([]byte(item[`images`]), &images)
	groupId := cast.ToInt64(item[`group_id`])
	return map[string]any{
		`id`:            cast.ToInt64(item[`id`]),
		`group_id`:      groupId,
		`group_names`:   goodsLibGroupNames(groupId, groupNames),
		`goods_id`:      item[`goods_id`],
		`goods_name`:    item[`goods_name`],
		`category`:      item[`category`],
		`brand`:         item[`brand`],
		`price`:         cast.ToFloat64(item[`price`]),
		`stock`:         cast.ToInt64(item[`stock`]),
		`link`:          item[`link`],
		`images`:        images,
		`description`:   item[`description`],
		`qa`:            item[`qa`],
		`custom_info`:   item[`custom_info`],
		`switch_status`: cast.ToInt(item[`switch_status`]),
		`create_time`:   cast.ToInt(item[`create_time`]),
		`update_time`:   cast.ToInt(item[`update_time`]),
	}
}

func buildGoodsLibGroupNames(lang string, groups []msql.Params) map[int64]string {
	ungrouped := i18n.Show(lang, `ungrouped_label`)
	groupNames := map[int64]string{0: ungrouped}
	groupMap := make(map[int64]msql.Params, len(groups))
	for _, group := range groups {
		groupMap[cast.ToInt64(group[`id`])] = group
	}
	var render func(groupId int64, visiting map[int64]bool) string
	render = func(groupId int64, visiting map[int64]bool) string {
		if groupId == 0 {
			return ungrouped
		}
		if name, ok := groupNames[groupId]; ok {
			return name
		}
		group, ok := groupMap[groupId]
		if !ok {
			return ungrouped
		}
		groupName := strings.TrimSpace(group[`group_name`])
		if len(groupName) == 0 {
			return ungrouped
		}
		if visiting[groupId] {
			return groupName
		}
		visiting[groupId] = true
		parentId := cast.ToInt64(group[`parent_id`])
		if parentId > 0 {
			parentName := render(parentId, visiting)
			if len(parentName) > 0 && parentName != ungrouped {
				groupName = parentName + `>` + groupName
			}
		}
		delete(visiting, groupId)
		groupNames[groupId] = groupName
		return groupName
	}
	for _, group := range groups {
		render(cast.ToInt64(group[`id`]), map[int64]bool{})
	}
	return groupNames
}

func goodsLibGroupNames(groupId int64, groupNames map[int64]string) string {
	if name, ok := groupNames[groupId]; ok && len(name) > 0 {
		return name
	}
	return groupNames[0]
}

func SaveGoodsLibLibrary(lang string, adminUserId int, params define.GoodsLibSaveParams) (int64, error) {
	if err := normalizeGoodsLibSaveParams(lang, &params); err != nil {
		return 0, err
	}
	if err := checkGoodsLibGroup(lang, adminUserId, params.GroupID); err != nil {
		return 0, err
	}
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
	if params.ID > 0 {
		info, err := m.Where(`id`, cast.ToString(params.ID)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Field(`id`).
			Find()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return 0, errors.New(i18n.Show(lang, `sys_err`))
		}
		if len(info) == 0 {
			return 0, errors.New(i18n.Show(lang, `no_data`))
		}
	}
	duplicateQuery := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`goods_id`, params.GoodsID)
	if params.ID > 0 {
		duplicateQuery.Where(`id`, `!=`, cast.ToString(params.ID))
	}
	duplicate, err := duplicateQuery.Value(`id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, duplicateQuery.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if cast.ToInt64(duplicate) > 0 {
		return 0, errors.New(i18n.Show(lang, `goods_id_exists`))
	}
	now := tool.Time2Int()
	data := goodsLibLibraryData(params)
	data[`update_time`] = now
	if params.SwitchStatus != nil {
		data[`switch_status`] = *params.SwitchStatus
	}
	if params.ID > 0 {
		_, err = m.Where(`id`, cast.ToString(params.ID)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Update(data)
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return 0, errors.New(i18n.Show(lang, `sys_err`))
		}
		return params.ID, nil
	}
	data[`admin_user_id`] = adminUserId
	data[`create_time`] = now
	if params.SwitchStatus == nil {
		data[`switch_status`] = define.GoodsLibSwitchOn
	}
	id, err := m.Insert(data, `id`)
	if isGoodsLibDuplicateError(err) {
		return 0, errors.New(i18n.Show(lang, `goods_id_exists`))
	}
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return id, nil
}

func goodsLibLibraryData(params define.GoodsLibSaveParams) msql.Datas {
	images, _ := json.Marshal(params.Images)
	return msql.Datas{
		`group_id`:    params.GroupID,
		`goods_id`:    params.GoodsID,
		`goods_name`:  params.GoodsName,
		`category`:    params.Category,
		`brand`:       params.Brand,
		`price`:       params.Price,
		`stock`:       params.Stock,
		`link`:        params.Link,
		`images`:      string(images),
		`description`: params.Description,
		`qa`:          params.QA,
		`custom_info`: params.CustomInfo,
	}
}

func normalizeGoodsLibSaveParams(lang string, params *define.GoodsLibSaveParams) error {
	params.GoodsID = strings.TrimSpace(params.GoodsID)
	params.GoodsName = strings.TrimSpace(params.GoodsName)
	params.Category = strings.TrimSpace(params.Category)
	params.Brand = strings.TrimSpace(params.Brand)
	params.Link = strings.TrimSpace(params.Link)
	params.Description = strings.TrimSpace(params.Description)
	params.QA = strings.TrimSpace(params.QA)
	params.CustomInfo = strings.TrimSpace(params.CustomInfo)
	if params.ID < 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `id`))
	}
	if params.GroupID < 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `group_id`))
	}
	if len(params.GoodsID) == 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_id`)))
	}
	if len(params.GoodsName) == 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_name`)))
	}
	if math.IsNaN(params.Price) || math.IsInf(params.Price, 0) || params.Price < 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_price`)))
	}
	if params.Stock < 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_stock`)))
	}
	if utf8.RuneCountInString(params.GoodsID) > define.GoodsLibBaseInfoMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_id`)))
	}
	if utf8.RuneCountInString(params.GoodsName) > define.GoodsLibBaseInfoMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_name`)))
	}
	if utf8.RuneCountInString(params.Category) > define.GoodsLibBaseInfoMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_category`)))
	}
	if utf8.RuneCountInString(params.Brand) > define.GoodsLibBaseInfoMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_brand`)))
	}
	if utf8.RuneCountInString(params.Link) > define.GoodsLibLinkMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_link`)))
	}
	if utf8.RuneCountInString(params.Description) > define.GoodsLibDetailMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_description`)))
	}
	if utf8.RuneCountInString(params.QA) > define.GoodsLibDetailMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_qa`)))
	}
	if utf8.RuneCountInString(params.CustomInfo) > define.GoodsLibDetailMaxLength {
		return errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_custom_info`)))
	}
	if params.SwitchStatus != nil && *params.SwitchStatus != define.GoodsLibSwitchOff && *params.SwitchStatus != define.GoodsLibSwitchOn {
		return errors.New(i18n.Show(lang, `param_invalid`, `switch_status`))
	}
	images, err := NormalizeGoodsLibImages(lang, params.Images)
	if err != nil {
		return err
	}
	params.Images = images
	return nil
}

func NormalizeGoodsLibImages(lang string, images []string) ([]string, error) {
	result := make([]string, 0, len(images))
	seen := make(map[string]bool, len(images))
	localPattern := regexp.MustCompile(`(?i)^/upload/chat_ai/\d+/goods_lib_image/\d+/[a-f0-9]{32}\.([a-z0-9]+)$`)
	for _, image := range images {
		image = strings.TrimSpace(image)
		if len(image) == 0 || seen[image] {
			continue
		}
		var ext string
		if strings.HasPrefix(image, `/`) {
			matches := localPattern.FindStringSubmatch(image)
			if len(matches) != 2 {
				return nil, errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_images`)))
			}
			ext = strings.ToLower(matches[1])
		} else {
			parsed, err := neturl.ParseRequestURI(image)
			if err != nil || (parsed.Scheme != `http` && parsed.Scheme != `https`) || len(parsed.Host) == 0 {
				return nil, errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_images`)))
			}
			ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(parsed.Path), `.`))
		}
		if !tool.InArrayString(ext, define.ImageAllowExt) {
			return nil, errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_images`)))
		}
		seen[image] = true
		result = append(result, image)
	}
	return result, nil
}

func checkGoodsLibGroup(lang string, adminUserId int, groupId int64) error {
	if groupId == 0 {
		return nil
	}
	m := msql.Model(define.TableGoodsLibGroup, define.Postgres)
	id, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(groupId)).
		Value(`id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if cast.ToInt64(id) == 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `group_id`))
	}
	return nil
}

func DeleteGoodsLibLibrary(lang string, adminUserId int, id int64) error {
	if id <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `id`))
	}
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
	rows, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if rows == 0 {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	return nil
}

func UpdateGoodsLibLibrarySwitch(lang string, adminUserId int, id int64, switchStatus int) error {
	if id <= 0 || (switchStatus != define.GoodsLibSwitchOff && switchStatus != define.GoodsLibSwitchOn) {
		return errors.New(i18n.Show(lang, `param_invalid`, `switch_status`))
	}
	m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
	rows, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{`switch_status`: switchStatus, `update_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if rows == 0 {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	return nil
}

func SaveGoodsLibImage(lang string, adminUserId int, fileHeader *multipart.FileHeader) (*define.UploadInfo, error) {
	uploadInfo, err := SaveUploadedFile(
		fileHeader,
		define.GoodsLibImageLimitSize,
		adminUserId,
		`goods_lib_image`,
		define.ImageAllowExt,
	)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	return uploadInfo, nil
}

func FormatGoodsLibRecommendResult(searchType string, list []map[string]any, total int) string {
	fields := []string{`group_names`, `goods_id`, `goods_name`, `category`, `brand`, `price`, `stock`, `link`, `images`, `description`, `qa`, `custom_info`}
	if searchType == define.GoodsLibRecommendSearchTypeBasic {
		fields = []string{`group_names`, `goods_id`, `goods_name`, `category`, `brand`, `price`, `stock`, `link`}
	}
	items := make([]map[string]any, 0, len(list))
	for index, item := range list {
		data := map[string]any{`sort`: index + 1}
		for _, field := range fields {
			data[field] = item[field]
		}
		items = append(items, data)
	}
	return tool.JsonEncodeNoError(map[string]any{
		`search_type`:    searchType,
		`total`:          total,
		`returned_count`: len(items),
		`items`:          items,
	})
}

func ImportGoodsLibLibrary(lang string, adminUserId int, groupId int64, fileHeader *multipart.FileHeader) (*define.GoodsLibImportResult, error) {
	if fileHeader == nil || fileHeader.Size <= 0 || fileHeader.Size > define.GoodsLibImportFileLimitSize {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `file`))
	}
	if err := checkGoodsLibGroup(lang, adminUserId, groupId); err != nil {
		return nil, err
	}
	uploadInfo, err := SaveUploadedFile(
		fileHeader,
		define.GoodsLibImportFileLimitSize,
		adminUserId,
		`goods_lib_import`,
		define.GoodsLibImportAllowExt,
	)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if uploadInfo == nil || len(uploadInfo.Link) == 0 {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `file`))
	}
	rows, err := ParseTabFile(uploadInfo.Link, uploadInfo.Ext)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `file`))
	}
	if len(rows) < 2 {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `file`))
	}
	headerMap := goodsLibImportHeaderMap(lang, rows[0])
	if _, ok := headerMap[`goods_id`]; !ok {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_id`)))
	}
	if _, ok := headerMap[`goods_name`]; !ok {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, i18n.Show(lang, `goods_import_header_goods_name`)))
	}
	result := &define.GoodsLibImportResult{
		Errors:  make([]define.GoodsLibImportError, 0),
		Headers: goodsLibImportHeaderDefinitions(lang),
	}
	hasDataRow := false
	for index, row := range rows[1:] {
		if goodsLibImportRowEmpty(row) {
			continue
		}
		hasDataRow = true
		result.TotalCount++
		params := define.GoodsLibSaveParams{
			GroupID:     groupId,
			GoodsID:     goodsLibImportCell(row, headerMap, `goods_id`),
			GoodsName:   goodsLibImportCell(row, headerMap, `goods_name`),
			Category:    goodsLibImportCell(row, headerMap, `category`),
			Brand:       goodsLibImportCell(row, headerMap, `brand`),
			Price:       cast.ToFloat64(goodsLibImportCell(row, headerMap, `price`)),
			Stock:       cast.ToInt64(goodsLibImportCell(row, headerMap, `stock`)),
			Link:        goodsLibImportCell(row, headerMap, `link`),
			Images:      splitGoodsLibImportImages(goodsLibImportCell(row, headerMap, `images`)),
			Description: goodsLibImportCell(row, headerMap, `description`),
			QA:          goodsLibImportCell(row, headerMap, `qa`),
			CustomInfo:  goodsLibImportCell(row, headerMap, `custom_info`),
		}
		if err = normalizeGoodsLibSaveParams(lang, &params); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, define.GoodsLibImportError{
				Row: index + 2, GoodsLibSaveParams: params, Message: err.Error(),
			})
			continue
		}
		m := msql.Model(define.TableGoodsLibLibrary, define.Postgres)
		id, queryErr := m.
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`goods_id`, params.GoodsID).
			Value(`id`)
		if queryErr != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), queryErr.Error())
			result.FailedCount++
			result.Errors = append(result.Errors, define.GoodsLibImportError{
				Row: index + 2, GoodsLibSaveParams: params, Message: i18n.Show(lang, `sys_err`),
			})
			continue
		}
		params.ID = cast.ToInt64(id)
		if _, saveErr := SaveGoodsLibLibrary(lang, adminUserId, params); saveErr != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, define.GoodsLibImportError{
				Row: index + 2, GoodsLibSaveParams: params, Message: saveErr.Error(),
			})
			continue
		}
		if params.ID > 0 {
			result.UpdatedCount++
		} else {
			result.CreatedCount++
		}
	}
	if !hasDataRow {
		return nil, errors.New(i18n.Show(lang, `param_invalid`, `file`))
	}
	return result, nil
}

func goodsLibImportHeaderDefinitions(lang string) []define.GoodsLibImportHeader {
	return []define.GoodsLibImportHeader{
		{Field: `goods_id`, Name: i18n.Show(lang, `goods_import_header_goods_id`)},
		{Field: `goods_name`, Name: i18n.Show(lang, `goods_import_header_goods_name`)},
		{Field: `category`, Name: i18n.Show(lang, `goods_import_header_category`)},
		{Field: `brand`, Name: i18n.Show(lang, `goods_import_header_brand`)},
		{Field: `price`, Name: i18n.Show(lang, `goods_import_header_price`)},
		{Field: `stock`, Name: i18n.Show(lang, `goods_import_header_stock`)},
		{Field: `link`, Name: i18n.Show(lang, `goods_import_header_link`)},
		{Field: `images`, Name: i18n.Show(lang, `goods_import_header_images`)},
		{Field: `description`, Name: i18n.Show(lang, `goods_import_header_description`)},
		{Field: `qa`, Name: i18n.Show(lang, `goods_import_header_qa`)},
		{Field: `custom_info`, Name: i18n.Show(lang, `goods_import_header_custom_info`)},
	}
}

func goodsLibImportHeaderMap(lang string, headers []string) map[string]int {
	aliases := make(map[string]string)
	for _, definition := range goodsLibImportHeaderDefinitions(lang) {
		aliases[normalizeGoodsLibImportHeader(definition.Name)] = definition.Field
	}
	result := make(map[string]int)
	for index, header := range headers {
		normalized := normalizeGoodsLibImportHeader(header)
		if field, ok := aliases[normalized]; ok {
			result[field] = index
		}
	}
	return result
}

func normalizeGoodsLibImportHeader(header string) string {
	header = strings.TrimPrefix(strings.TrimSpace(header), "\uFEFF")
	header = strings.ToLower(header)
	replacer := strings.NewReplacer(` `, ``, `_`, ``, `-`, ``, `*`, ``, `：`, ``, `:`, ``)
	return replacer.Replace(header)
}

func goodsLibImportCell(row []string, headerMap map[string]int, field string) string {
	index, ok := headerMap[field]
	if !ok || index < 0 || index >= len(row) {
		return ``
	}
	return row[index]
}

func goodsLibImportRowEmpty(row []string) bool {
	for _, value := range row {
		if len(strings.TrimSpace(value)) > 0 {
			return false
		}
	}
	return true
}

func splitGoodsLibImportImages(value string) []string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.ReplaceAll(value, `；`, "\n")
	value = strings.ReplaceAll(value, `;`, "\n")
	if value == `` {
		return []string{}
	}
	return strings.Split(value, "\n")
}

func ExportGoodsLibLibrary(lang string, adminUserId int, filter define.GoodsLibListFilter, template bool) (string, string, error) {
	headerDefinitions := goodsLibImportHeaderDefinitions(lang)
	fields := make(tool.Fields, 0, len(headerDefinitions))
	for _, definition := range headerDefinitions {
		fields = append(fields, struct {
			Field  string
			Header string
		}{Field: definition.Field, Header: definition.Name})
	}
	data := make([]map[string]any, 0)
	title := i18n.Show(lang, `goods_export_title`)
	if template {
		title = i18n.Show(lang, `goods_import_template_title`)
		data = append(data, map[string]any{
			`goods_id`:    i18n.Show(lang, `goods_import_example_goods_id`),
			`goods_name`:  i18n.Show(lang, `goods_import_example_goods_name`),
			`category`:    i18n.Show(lang, `goods_import_example_category`),
			`brand`:       i18n.Show(lang, `goods_import_example_brand`),
			`price`:       99.00,
			`stock`:       100,
			`link`:        `https://example.com/goods.html`,
			`images`:      `https://example.com/goods.jpg`,
			`description`: i18n.Show(lang, `goods_import_example_description`),
			`qa`:          i18n.Show(lang, `goods_import_example_qa`),
			`custom_info`: i18n.Show(lang, `goods_import_example_custom_info`),
		})
	} else {
		query := buildGoodsLibLibraryQuery(adminUserId, filter)
		list, err := query.
			Field(goodsLibLibrarySelectFields()).
			Order(`id desc`).
			Select()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, query.GetLastSql(), err.Error())
			return ``, ``, errors.New(i18n.Show(lang, `sys_err`))
		}
		data = make([]map[string]any, 0, len(list))
		for _, item := range list {
			images := make([]string, 0)
			_ = json.Unmarshal([]byte(item[`images`]), &images)
			data = append(data, map[string]any{
				`goods_id`:    item[`goods_id`],
				`goods_name`:  item[`goods_name`],
				`category`:    item[`category`],
				`brand`:       item[`brand`],
				`price`:       cast.ToFloat64(item[`price`]),
				`stock`:       cast.ToInt64(item[`stock`]),
				`link`:        item[`link`],
				`images`:      strings.Join(images, "\r\n"),
				`description`: item[`description`],
				`qa`:          item[`qa`],
				`custom_info`: item[`custom_info`],
			})
		}
	}
	filePath, fileName, err := tool.ExcelExportPro(data, fields, title, `static/public/download`)
	if err != nil {
		logs.Error(err.Error())
		return ``, ``, errors.New(i18n.Show(lang, `sys_err`))
	}
	return filePath, fileName, nil
}

func isGoodsLibDuplicateError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == `23505`
}
