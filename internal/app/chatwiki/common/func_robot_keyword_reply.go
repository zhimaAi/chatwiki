// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const MaxKeywordNum = 500
const MaxKeywordReplyRuleNum = 1000

// RobotKeywordReplyCacheBuildHandler Robot keyword reply rule cache
type RobotKeywordReplyCacheBuildHandler struct {
	RobotId int
}

type RobotKeywordReply struct {
	ID           int            `json:"id"`
	AdminUserID  int            `json:"admin_user_id"`
	RobotID      int            `json:"robot_id"`
	Name         string         `json:"name"`
	FullKeyword  []string       `json:"full_keyword"`
	HalfKeyword  []string       `json:"half_keyword"`
	ReplyContent []ReplyContent `json:"reply_content"`
	ReplyType    []string       `json:"reply_type"`
	SwitchStatus int            `json:"switch_status"`
	ReplyNum     int            `json:"reply_num"`
}

func (h *RobotKeywordReplyCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.func_chat_robot_keyword_reply.%d`, h.RobotId)
}

func (h *RobotKeywordReplyCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where(`robot_id`, cast.ToString(h.RobotId)).Where(`switch_status`, cast.ToString(define.SwitchOn)).Order(`id desc`).Select()
	// Convert
	list := make([]RobotKeywordReply, 0)
	if err == nil && len(data) > 0 {
		for _, item := range data {
			//
			var replyContent = make([]ReplyContent, 0)
			_ = tool.JsonDecodeUseNumber(item[`reply_content`], &replyContent)
			if len(replyContent) == 0 {
				continue
			}
			list = append(list, RobotKeywordReply{
				ID:           cast.ToInt(item[`id`]),
				AdminUserID:  cast.ToInt(item[`admin_user_id`]),
				RobotID:      cast.ToInt(item[`robot_id`]),
				Name:         item[`name`],
				FullKeyword:  DisposeStringList(item[`full_keyword`]),
				HalfKeyword:  DisposeStringList(item[`half_keyword`]),
				ReplyContent: replyContent, // Process reply content
				ReplyType:    DisposeStringList(item[`reply_type`]),
				SwitchStatus: cast.ToInt(item[`switch_status`]),
				ReplyNum:     cast.ToInt(item[`reply_num`]),
			})
		}
	}
	return list, err
}

// GetRobotKeywordReplyListByRobotId Robot keyword reply rule list
func GetRobotKeywordReplyListByRobotId(robotId int) ([]RobotKeywordReply, error) {
	result := make([]RobotKeywordReply, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotKeywordReplyCacheBuildHandler{RobotId: robotId}, &result, time.Hour)
	if err != nil {
		return nil, err
	}
	return result, err
}

// SaveRobotKeywordReply Save keyword reply rule (create or update)
func SaveRobotKeywordReply(id, adminUserID, robotID int, name string, fullKeyword, halfKeyword []string, replyContent []ReplyContent, replyType []string, replyNum int) (int64, error) {
	fullKeywordJson, _ := tool.JsonEncode(fullKeyword)
	halfKeywordJson, _ := tool.JsonEncode(halfKeyword)
	replyContentJson, _ := tool.JsonEncode(replyContent)

	// Compatible judgment
	if len(replyType) == 0 && len(replyContent) > 0 {
		// Loop replyContent to assign reply_type to replyType
		for _, reply := range replyContent {
			if len(reply.ReplyType) == 0 {
				replyType = append(replyType, reply.ReplyType)
			}
		}
	}

	replyTypeJson, _ := tool.JsonEncode(replyType)

	data := msql.Datas{
		"admin_user_id": adminUserID,
		"robot_id":      robotID,
		"name":          name,
		"full_keyword":  fullKeywordJson,
		"half_keyword":  halfKeywordJson,
		"reply_content": replyContentJson,
		"reply_type":    replyTypeJson,
		"reply_num":     replyNum,
		"update_time":   tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// Create new record
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Insert(data, "id")
	} else {
		// Update existing record
		_, err = msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotKeywordReplyCacheBuildHandler{RobotId: robotID})
	}
	return newId, err
}

// DeleteRobotKeywordReply Delete keyword reply rule
func DeleteRobotKeywordReply(id, robotID int) error {
	_, err := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotKeywordReplyCacheBuildHandler{RobotId: robotID})
	}
	return err
}

// GetRobotKeywordReply Get a single keyword reply rule
func GetRobotKeywordReply(id int) (msql.Params, error) {
	return msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("id", cast.ToString(id)).Find()
}

// GetRobotKeywordReplyUseRuleNum Get the number of keyword reply rules in use by the robot
func GetRobotKeywordReplyUseRuleNum(robotID int) (int, error) {
	total, err := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where(`robot_id`, cast.ToString(robotID)).Where(`switch_status`, cast.ToString(define.SwitchOn)).Count()
	return total, err
}

// GetRobotKeywordReplyListWithFilter Get keyword reply rule list (with filters and pagination)
func GetRobotKeywordReplyListWithFilter(robotID int, keyword, replyType string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID))

	// Fuzzy query by keyword (in full_keyword or half_keyword)
	if keyword != "" {
		// Query using PostgreSQL's jsonb operator
		model.Where("( (name ILIKE '%" + keyword + "%') OR  (full_keyword ?| array['" + keyword + "'] OR half_keyword ?| array['" + keyword + "']) )")
	}

	// Query by reply type
	if replyType != "" {
		// Query using PostgreSQL's jsonb operator
		model.Where("reply_type ?| array['" + replyType + "']")
	}

	// Add pagination
	list, total, err := model.Order("id DESC").Paginate(page, size)
	return list, total, err
}

// UpdateRobotKeywordReplySwitchStatus Update keyword reply rule switch status
func UpdateRobotKeywordReplySwitchStatus(id, robotID, switchStatus int) error {
	data := msql.Datas{
		"switch_status": switchStatus,
		"update_time":   tool.Time2Int(),
	}

	_, err := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotKeywordReplyCacheBuildHandler{RobotId: robotID})
	}
	return err
}

// CheckKeyWordRepeat Check if keyword is repeated
func CheckKeyWordRepeat(robotID int, keyword string, id int) (bool, string) {
	// Check if repeated, return repeat status and rule name
	model := msql.Model(`func_chat_robot_keyword_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID))
	// Fuzzy query by keyword (in full_keyword or half_keyword)
	if keyword != "" {
		// Query using PostgreSQL's jsonb operator
		model.Where(" (full_keyword ?| array['" + keyword + "'] OR half_keyword ?| array['" + keyword + "']) ")
	} else {
		return false, ""
	}
	checkRule, err := model.Order("id DESC").Find()
	if err != nil {
		return false, ""
	}
	if len(checkRule) == 0 {
		return false, ""
	}
	return checkRule["id"] != cast.ToString(id), checkRule["name"]
}
