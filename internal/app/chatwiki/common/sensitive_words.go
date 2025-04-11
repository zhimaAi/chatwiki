// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
)

func GetSensitiveWordsList(adminUserId, page, size int) ([]msql.Params, int, error) {
	data, total, err := msql.Model(`sensitive_words`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Order("id desc").Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}
	if len(data) == 0 {
		return data, total, nil
	}
	var wordIds, robotIds []string
	for _, v := range data {
		wordIds = append(wordIds, v["id"])
	}
	relationData, err := msql.Model(`sensitive_words_relation`, define.Postgres).Where("words_id", `in`, strings.Join(wordIds, `,`)).Field(`words_id,robot_id`).Select()
	if err != nil {
		return nil, 0, err
	}
	if len(relationData) == 0 {
		return data, total, nil
	}
	relationMap := map[string][]string{}
	for _, v := range relationData {
		robotIds = append(robotIds, v["robot_id"])
		relationMap[v["words_id"]] = append(relationMap[v["words_id"]], v["robot_id"])
	}
	robotData, err := msql.Model(`chat_ai_robot`, define.Postgres).Where("id", `in`, strings.Join(robotIds, `,`)).Field(`id,robot_name`).Select()
	if err != nil {
		return nil, 0, err
	}
	robotMap := map[string]msql.Params{}
	for _, v := range robotData {
		robotMap[v["id"]] = v
	}
	for k, v := range data {
		robotData := []msql.Params{}
		data[k]["robot_data"] = `[]`
		if len(relationMap[v["id"]]) > 0 {
			for _, vv := range relationMap[v["id"]] {
				if data, ok := robotMap[vv]; ok {
					robotData = append(robotData, data)
				}
			}
			data[k]["robot_data"] = tool.JsonEncodeNoError(robotData)
		}
	}
	return data, total, nil
}

func DeleteSensitiveWords(adminUserId, id int) (int64, error) {
	if _, err := msql.Model(`sensitive_words`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where(`id`, cast.ToString(id)).Delete(); err != nil {
		return 0, err
	}
	return msql.Model(`sensitive_words_relation`, define.Postgres).Where("words_id", cast.ToString(id)).Delete()
}

func SaveSensitiveWords(adminUserId, id int64, words, robotIds string, triggerType int) (int64, error) {
	m := msql.Model(`sensitive_words`, define.Postgres)
	err := m.Begin()
	if err != nil {
		return 0, err
	}
	if id > 0 {
		if _, err := msql.Model(`sensitive_words_relation`, define.Postgres).Where("words_id", cast.ToString(id)).Delete(); err != nil {
			m.Rollback()
			return 0, err
		}
		_, err := msql.Model(`sensitive_words`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where(`id`, cast.ToString(id)).Update(msql.Datas{
			"words":        words,
			"trigger_type": triggerType,
			`update_time`:  tool.Time2Int(),
		})
		if err != nil {
			m.Rollback()
			return 0, err
		}
	} else {
		id, err = msql.Model(`sensitive_words`, define.Postgres).Insert(msql.Datas{
			`words`:         words,
			`trigger_type`:  triggerType,
			`admin_user_id`: adminUserId,
			`status`:        define.SwitchOn,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}, `id`)
		if err != nil {
			m.Rollback()
			return 0, err
		}
	}
	robotIdsArr := []string{"0"}
	if triggerType > 0 {
		robotIdsArr = strings.Split(robotIds, ",")
	}
	for _, v := range robotIdsArr {
		if _, err := msql.Model(`sensitive_words_relation`, define.Postgres).Insert(msql.Datas{
			"words_id":    id,
			"robot_id":    v,
			`create_time`: tool.Time2Int(),
			`update_time`: tool.Time2Int(),
		}); err != nil {
			m.Rollback()
			return 0, err
		}
	}
	m.Commit()
	return id, err
}

func SwitchSensitiveWords(adminUserId, id int) (int64, error) {
	sqlRaw := fmt.Sprintf(`update_time = %d,status = 1-status`, tool.Time2Int())
	if _, err := msql.Model(`sensitive_words`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where(`id`, cast.ToString(id)).Update2(sqlRaw); err != nil {
		return 0, err
	}
	if _, err := msql.Model(`sensitive_words_relation`, define.Postgres).Where("words_id", cast.ToString(id)).Update2(sqlRaw); err != nil {
		return 0, err
	}
	return 0, nil
}

func CheckSensitiveWords(question string, adminUserId, robotId int) (bool, []string) {
	var wordsArr []string
	data, err := GetSenitiveWordsCache(adminUserId)
	if err != nil {
		return false, wordsArr
	}
	for _, v := range data {
		words := strings.Split(strings.ReplaceAll(strings.ReplaceAll(cast.ToString(v["words"]), "\r\n", "\n"), "\r", "\n"), "\n")
		for _, vv := range words {
			word := strings.TrimSpace(vv)
			if cast.ToInt(v[`trigger_type`]) == 0 || tool.InArrayString(cast.ToString(robotId), cast.ToStringSlice(v[`robot_ids`])) {
				if strings.Contains(question, word) && !tool.InArrayString(word, wordsArr) {
					wordsArr = append(wordsArr, vv)
				}
			}
		}
	}
	return len(wordsArr) > 0, wordsArr
}

type SenitiveWordsCacheHandle struct {
	AdminUserId int
}

func (h *SenitiveWordsCacheHandle) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.sensitive_words.%d`, h.AdminUserId)
}

func (h *SenitiveWordsCacheHandle) GetCacheData() (any, error) {
	var wordsData = []map[string]any{}
	data, err := msql.Model(`sensitive_words`, define.Postgres).Where("admin_user_id", cast.ToString(h.AdminUserId)).Where(`status`, cast.ToString(define.SwitchOn)).Field(`id,words,trigger_type`).Select()
	if err != nil {
		return wordsData, err
	}
	if len(data) == 0 {
		return wordsData, nil
	}
	var ids []string
	for _, v := range data {
		ids = append(ids, v["id"])
	}
	relationData, err := msql.Model(`sensitive_words_relation`, define.Postgres).Where("words_id", `in`, strings.Join(ids, `,`)).Field(`words_id,robot_id`).Select()
	if err != nil {
		return wordsData, err
	}
	if len(relationData) == 0 {
		return wordsData, nil
	}
	relationMap := map[string][]string{}
	for _, v := range relationData {
		relationMap[v["words_id"]] = append(relationMap[v["words_id"]], v["robot_id"])
	}
	for _, v := range data {
		wordsData = append(wordsData, map[string]any{
			"words":        v["words"],
			"trigger_type": cast.ToInt(v["trigger_type"]),
			"robot_ids":    relationMap[v["id"]],
		})
	}
	return wordsData, nil
}

func GetSenitiveWordsCache(adminUserId int) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &SenitiveWordsCacheHandle{AdminUserId: adminUserId}, &result, time.Hour)
	return result, err
}
