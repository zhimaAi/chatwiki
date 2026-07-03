// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const (
	TableAdminMiniCard         = `admin_mini_card`
	TableAdminMiniCardRelation = `admin_mini_card_relation`

	AdminMiniCardTargetLibraryQA        = `library_qa`
	AdminMiniCardTargetLibraryParagraph = `library_paragraph`
	AdminMiniCardTargetRobotPrompt      = `robot_prompt`
)

var (
	adminMiniCardPromptRegexp       = regexp.MustCompile(`【chat_mini_card_id:(\d+)】`)
	miniCardPromptShortRegexp       = regexp.MustCompile(`【mini_card:(\d+)】`)
	rawMiniCardPromptRegexp         = regexp.MustCompile(`【(?:chat_mini_card_id|mini_card):(\d+)】`)
	rawMiniCardPromptMarkerPrefixes = []string{`【chat_mini_card_id:`, `【mini_card:`}
)

func SaveAdminMiniCard(id, adminUserID int, title, appid, pagePath, thumbURL string) (int64, error) {
	data := msql.Datas{
		`admin_user_id`: adminUserID,
		`title`:         title,
		`appid`:         appid,
		`page_path`:     pagePath,
		`thumb_url`:     thumbURL,
		`update_time`:   tool.Time2Int(),
	}

	if id <= 0 {
		data[`create_time`] = tool.Time2Int()
		return msql.Model(TableAdminMiniCard, define.Postgres).Insert(data, `id`)
	}

	_, err := msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`delete_time`, `0`).
		Update(data)
	return int64(id), err
}

func GetAdminMiniCard(id, adminUserID int) (msql.Params, error) {
	return msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`delete_time`, `0`).
		Find()
}

// GetAdminMiniCardsByIDs returns active mini cards keyed by card ID.
func GetAdminMiniCardsByIDs(adminUserID int, ids []int) (map[int]map[string]any, error) {
	result := make(map[int]map[string]any)
	if adminUserID <= 0 || len(ids) == 0 {
		return result, nil
	}
	list, err := msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`id`, `in`, joinAdminMiniCardIDs(ids)).
		Where(`delete_time`, `0`).
		Field(`id,admin_user_id,title,appid,page_path,thumb_url,create_time,update_time`).
		Select()
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		result[cast.ToInt(item[`id`])] = FormatAdminMiniCard(item)
	}
	return result, nil
}

func GetAdminMiniCardList(adminUserID int, keyword, appid string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`delete_time`, `0`)
	if keyword != `` {
		model.Where(`title LIKE ?`, `%`+keyword+`%`)
	}
	if appid != `` {
		model.Where(`appid`, appid)
	}
	return model.Order(`id DESC`).Paginate(page, size)
}

func DeleteAdminMiniCard(id, adminUserID int) error {
	_, err := msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`delete_time`, `0`).
		Update(msql.Datas{
			`delete_time`: tool.Time2Int(),
			`update_time`: tool.Time2Int(),
		})
	return err
}

func CountAdminMiniCardRelations(adminUserID, miniCardID int) (int, error) {
	return msql.Model(TableAdminMiniCardRelation, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`mini_card_id`, cast.ToString(miniCardID)).
		Where(`delete_time`, `0`).
		Count()
}

func FormatAdminMiniCard(item msql.Params) map[string]any {
	return map[string]any{
		`id`:            item[`id`],
		`admin_user_id`: item[`admin_user_id`],
		`title`:         item[`title`],
		`appid`:         item[`appid`],
		`page_path`:     item[`page_path`],
		`thumb_url`:     item[`thumb_url`],
		`reply_type`:    ReplyTypeCard,
		`create_time`:   item[`create_time`],
		`update_time`:   item[`update_time`],
	}
}

// ParseAdminMiniCardIDs parses comma-separated mini card IDs, trims spaces, and removes duplicates.
func ParseAdminMiniCardIDs(raw string) ([]int, error) {
	raw = strings.TrimSpace(raw)
	if raw == `` {
		return make([]int, 0), nil
	}
	parts := strings.Split(raw, `,`)
	ids := make([]int, 0, len(parts))
	idSet := make(map[int]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == `` {
			return nil, errors.New(`mini_card invalid`)
		}
		id := cast.ToInt(part)
		if id <= 0 || cast.ToString(id) != part {
			return nil, errors.New(`mini_card invalid`)
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

// ExtractAdminMiniCardIDsFromPrompt extracts mini card IDs from full-width prompt tags.
func ExtractAdminMiniCardIDsFromPrompt(prompt string) []int {
	return extractMiniCardIDsFromPromptByRegexp(adminMiniCardPromptRegexp, prompt)
}

// ExtractMiniCardIDsFromShortPrompt extracts mini card IDs from short prompt tags.
func ExtractMiniCardIDsFromShortPrompt(prompt string) []int {
	return extractMiniCardIDsFromPromptByRegexp(miniCardPromptShortRegexp, prompt)
}

func extractMiniCardIDsFromPromptByRegexp(re *regexp.Regexp, prompt string) []int {
	matches := re.FindAllStringSubmatch(prompt, -1)
	ids := make([]int, 0, len(matches))
	idSet := make(map[int]struct{}, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		id := cast.ToInt(match[1])
		if id <= 0 {
			continue
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

// RemoveMiniCardPromptMarkers removes raw prompt mini card markers from text.
func RemoveMiniCardPromptMarkers(content string) string {
	if len(content) == 0 || !rawMiniCardPromptRegexp.MatchString(content) {
		return content
	}
	return rawMiniCardPromptRegexp.ReplaceAllString(content, ``)
}

// IsMiniCardPromptMarkerPrefix reports whether content may be an unfinished raw mini card marker.
func IsMiniCardPromptMarkerPrefix(content string) bool {
	if len(content) == 0 {
		return false
	}
	for _, prefix := range rawMiniCardPromptMarkerPrefixes {
		if len(content) <= len(prefix) {
			if strings.HasPrefix(prefix, content) {
				return true
			}
			continue
		}
		if !strings.HasPrefix(content, prefix) {
			continue
		}
		for _, r := range content[len(prefix):] {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// ValidateAdminMiniCards checks whether all card IDs belong to the admin user and are not deleted.
func ValidateAdminMiniCards(adminUserID int, ids []int) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	count, err := msql.Model(TableAdminMiniCard, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`id`, `in`, joinAdminMiniCardIDs(ids)).
		Where(`delete_time`, `0`).
		Count()
	if err != nil {
		return false, err
	}
	return count == len(ids), nil
}

// SaveAdminMiniCardRelations replaces active mini card relations for one target object.
func SaveAdminMiniCardRelations(adminUserID, libraryID int, targetType string, targetID int, miniCardIDs []int) error {
	if adminUserID <= 0 || targetID <= 0 || targetType == `` {
		return nil
	}
	now := tool.Time2Int()
	model := msql.Model(TableAdminMiniCardRelation, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`target_type`, targetType).
		Where(`target_id`, cast.ToString(targetID)).
		Where(`delete_time`, `0`)
	activeRelations, err := model.Field(`mini_card_id`).Select()
	if err != nil {
		return err
	}
	activeCardIDs := make(map[int]struct{}, len(activeRelations))
	for _, relation := range activeRelations {
		activeCardIDs[cast.ToInt(relation[`mini_card_id`])] = struct{}{}
	}
	newCardIDs := make(map[int]struct{}, len(miniCardIDs))
	for _, miniCardID := range miniCardIDs {
		newCardIDs[miniCardID] = struct{}{}
	}
	deleteIDs := make([]int, 0)
	for miniCardID := range activeCardIDs {
		if _, ok := newCardIDs[miniCardID]; !ok {
			deleteIDs = append(deleteIDs, miniCardID)
		}
	}
	if len(deleteIDs) > 0 {
		if _, err = msql.Model(TableAdminMiniCardRelation, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserID)).
			Where(`target_type`, targetType).
			Where(`target_id`, cast.ToString(targetID)).
			Where(`mini_card_id`, `in`, joinAdminMiniCardIDs(deleteIDs)).
			Where(`delete_time`, `0`).
			Update(msql.Datas{
				`delete_time`: now,
				`update_time`: now,
			}); err != nil {
			return err
		}
	}
	for _, miniCardID := range miniCardIDs {
		if _, ok := activeCardIDs[miniCardID]; ok {
			continue
		}
		if _, err = msql.Model(TableAdminMiniCardRelation, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: adminUserID,
			`mini_card_id`:  miniCardID,
			`library_id`:    libraryID,
			`target_type`:   targetType,
			`target_id`:     targetID,
			`create_time`:   now,
			`update_time`:   now,
		}); err != nil {
			return err
		}
	}
	return nil
}

// ClearAdminMiniCardRelationsByTargetIds soft-deletes mini card relations for the given target IDs.
func ClearAdminMiniCardRelationsByTargetIds(adminUserID int, targetType string, targetIDs []int) error {
	if adminUserID <= 0 || targetType == `` || len(targetIDs) == 0 {
		return nil
	}
	now := tool.Time2Int()
	_, err := msql.Model(TableAdminMiniCardRelation, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserID)).
		Where(`target_type`, targetType).
		Where(`target_id`, `in`, joinAdminMiniCardIDs(targetIDs)).
		Where(`delete_time`, `0`).
		Update(msql.Datas{
			`delete_time`: now,
			`update_time`: now,
		})
	return err
}

// GetAdminMiniCardsByTargets returns active mini cards grouped by target ID and ordered by card ID descending.
func GetAdminMiniCardsByTargets(adminUserID int, targetType string, targetIDs []int) (map[int][]map[string]any, error) {
	result := make(map[int][]map[string]any, len(targetIDs))
	for _, targetID := range targetIDs {
		result[targetID] = make([]map[string]any, 0)
	}
	if adminUserID <= 0 || targetType == `` || len(targetIDs) == 0 {
		return result, nil
	}
	list, err := msql.Model(TableAdminMiniCardRelation, define.Postgres).
		Alias(`r`).
		Join(TableAdminMiniCard+` c`, `r.mini_card_id=c.id`, `inner`).
		Where(`r.admin_user_id`, cast.ToString(adminUserID)).
		Where(`r.target_type`, targetType).
		Where(`r.target_id`, `in`, joinAdminMiniCardIDs(targetIDs)).
		Where(`r.delete_time`, `0`).
		Where(`c.delete_time`, `0`).
		Field(`r.target_id,c.id,c.admin_user_id,c.title,c.appid,c.page_path,c.thumb_url,c.create_time,c.update_time`).
		Order(`c.id DESC`).
		Select()
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		targetID := cast.ToInt(item[`target_id`])
		result[targetID] = append(result[targetID], FormatAdminMiniCard(item))
	}
	return result, nil
}

func joinAdminMiniCardIDs(ids []int) string {
	items := make([]string, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			items = append(items, cast.ToString(id))
		}
	}
	return strings.Join(items, `,`)
}
