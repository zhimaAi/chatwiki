// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// RobotReceivedMessageReplyCacheBuildHandler Robot received message reply rule cache
type RobotReceivedMessageReplyCacheBuildHandler struct {
	RobotId  int
	RuleType string
}

const (
	MessageTypeAll     = define.MessageTypeAll
	MessageTypeSpecify = define.MessageTypeSpecify
)
const (
	RuleTypeMessageType = define.RuleTypeMessageType // Message type rule
	RuleTypeDuration    = define.RuleTypeDuration    // Time range rule
)

const (
	DurationTypeWeek      = define.DurationTypeWeek      // Duration type: week
	DurationTypeDay       = define.DurationTypeDay       // Duration type: day
	DurationTypeTimeRange = define.DurationTypeTimeRange // Time range rule

)

type RobotReceivedMessageReply struct {
	ID                 int            `json:"id"`
	AdminUserID        int            `json:"admin_user_id"`
	RobotID            int            `json:"robot_id"`
	RuleType           string         `json:"rule_type"`
	DurationType       string         `json:"duration_type"`
	WeekDuration       []int          `json:"week_duration"`
	StartDay           string         `json:"start_day"`
	EndDay             string         `json:"end_day"`
	StartDuration      string         `json:"start_duration"`
	EndDuration        string         `json:"end_duration"`
	PriorityNum        int            `json:"priority_num"`
	ReplyInterval      int            `json:"reply_interval"`
	MessageType        int            `json:"message_type"`
	SpecifyMessageType []string       `json:"specify_message_type"`
	ReplyContent       []ReplyContent `json:"reply_content"`
	ReplyType          []string       `json:"reply_type"`
	SwitchStatus       int            `json:"switch_status"`
	ReplyNum           int            `json:"reply_num"`
	CreateTime         int            `json:"create_time"`
	UpdateTime         int            `json:"update_time"`
}

func (h *RobotReceivedMessageReplyCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.func_chat_robot_received_message_reply.%d.%s`, h.RobotId, h.RuleType)
}

func (h *RobotReceivedMessageReplyCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where(`robot_id`, cast.ToString(h.RobotId)).Where(`rule_type`, h.RuleType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Order(`priority_num asc,id desc`).Select()
	// Convert
	list := make([]RobotReceivedMessageReply, 0)
	if err == nil && len(data) > 0 {
		for _, item := range data {
			//
			var replyContent = make([]ReplyContent, 0)
			_ = tool.JsonDecodeUseNumber(item[`reply_content`], &replyContent)
			if len(replyContent) == 0 {
				continue
			}

			// Parse week_duration field
			var weekDuration = make([]int, 0)
			_ = tool.JsonDecodeUseNumber(item[`week_duration`], &weekDuration)

			// Parse specify_message_type field
			var specifyMessageType = make([]string, 0)
			_ = tool.JsonDecodeUseNumber(item[`specify_message_type`], &specifyMessageType)

			// Parse reply_type field
			var replyType = make([]string, 0)
			_ = tool.JsonDecodeUseNumber(item[`reply_type`], &replyType)

			list = append(list, RobotReceivedMessageReply{
				ID:                 cast.ToInt(item[`id`]),
				AdminUserID:        cast.ToInt(item[`admin_user_id`]),
				RobotID:            cast.ToInt(item[`robot_id`]),
				RuleType:           item[`rule_type`],
				DurationType:       item[`duration_type`],
				WeekDuration:       weekDuration,
				StartDay:           item[`start_day`],
				EndDay:             item[`end_day`],
				StartDuration:      item[`start_duration`],
				EndDuration:        item[`end_duration`],
				PriorityNum:        cast.ToInt(item[`priority_num`]),
				ReplyInterval:      cast.ToInt(item[`reply_interval`]),
				MessageType:        cast.ToInt(item[`message_type`]),
				SpecifyMessageType: specifyMessageType,
				ReplyContent:       replyContent,
				ReplyType:          replyType,
				SwitchStatus:       cast.ToInt(item[`switch_status`]),
				ReplyNum:           cast.ToInt(item[`reply_num`]),
				CreateTime:         cast.ToInt(item[`create_time`]),
				UpdateTime:         cast.ToInt(item[`update_time`]),
			})
		}
	}
	return list, err
}

type ReceivedMessageReplyLastTimeCacheBuildHandler struct {
	RobotId  int
	RuleId   int
	Openid   string
	LastTime int
}

func (h *ReceivedMessageReplyLastTimeCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.func_chat_robot_received_message_reply_interval.%d.%d.%s`, h.RobotId, h.RuleId, h.Openid)
}

func (h *ReceivedMessageReplyLastTimeCacheBuildHandler) GetCacheData() (any, error) {
	if h.LastTime == 0 {
		h.LastTime = tool.Time2Int()
	}
	return map[string]interface{}{
		"last_time": h.LastTime,
	}, nil
}

// GetReceivedMessageReplyLastTime Get the last reply time for received message reply rules
func GetReceivedMessageReplyLastTime(robotId, ruleId int, openid string) (int, error) {
	result := map[string]interface{}{
		"last_time": 0,
	}
	err := lib_redis.GetOne(define.Redis, &ReceivedMessageReplyLastTimeCacheBuildHandler{RobotId: robotId, RuleId: ruleId, Openid: openid}, &result)
	if err != nil {
		return 0, err
	}
	return cast.ToInt(result["last_time"]), err
}

func SetReceivedMessageReplyLastTime(robotId, ruleId int, lastTime int, openid string) error {
	return lib_redis.SetOne(define.Redis, &ReceivedMessageReplyLastTimeCacheBuildHandler{RobotId: robotId, RuleId: ruleId, LastTime: lastTime, Openid: openid}, time.Hour)
}

// GetRobotReceivedMessageReplyListByRobotId Robot received message reply rule list
func GetRobotReceivedMessageReplyListByRobotId(robotId int, ruleType string) ([]RobotReceivedMessageReply, error) {
	result := make([]RobotReceivedMessageReply, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotId, RuleType: ruleType}, &result, time.Hour)
	if err != nil {
		return nil, err
	}
	return result, err
}

// SaveRobotReceivedMessageReply Save received message reply rule (create or update)
func SaveRobotReceivedMessageReply(id, adminUserID, robotID int, ruleType, durationType string, weekDuration []int, startDay, endDay, startDuration, endDuration string, priorityNum, replyInterval, messageType int, specifyMessageType []string, replyContent []ReplyContent, replyType []string, replyNum int, switchStatus int) (int64, error) {
	weekDurationJson, _ := tool.JsonEncode(weekDuration)
	specifyMessageTypeJson, _ := tool.JsonEncode(specifyMessageType)
	replyContentJson, _ := tool.JsonEncode(replyContent)
	replyTypeJson, _ := tool.JsonEncode(replyType)

	data := msql.Datas{
		"admin_user_id":        adminUserID,
		"robot_id":             robotID,
		"rule_type":            ruleType,
		"duration_type":        durationType,
		"week_duration":        weekDurationJson,
		"start_day":            startDay,
		"end_day":              endDay,
		"start_duration":       startDuration,
		"end_duration":         endDuration,
		"priority_num":         priorityNum,
		"reply_interval":       replyInterval,
		"message_type":         messageType,
		"specify_message_type": specifyMessageTypeJson,
		"reply_content":        replyContentJson,
		"reply_type":           replyTypeJson,
		"reply_num":            replyNum,
		"switch_status":        switchStatus,
		"update_time":          tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// Create new record
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Insert(data, "id")
	} else {
		// Update existing record
		_, err = msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotID, RuleType: ruleType})
	}
	return newId, err
}

// DeleteRobotReceivedMessageReply Delete received message reply rule
func DeleteRobotReceivedMessageReply(id, robotID int) error {
	oldOne, err := GetRobotReceivedMessageReply(id, robotID)
	if err != nil {
		return err
	}
	ruleType := oldOne[`rule_type`]

	_, err = msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotID, RuleType: ruleType})
	}
	return err
}

// UpdateRobotReceivedMessageReplyPriorityNum Update received message reply rule priority
func UpdateRobotReceivedMessageReplyPriorityNum(id, robotID, priorityNum int) error {
	oldOne, err := GetRobotReceivedMessageReply(id, robotID)
	if err != nil {
		return err
	}
	ruleType := oldOne[`rule_type`]

	data := msql.Datas{
		"priority_num": priorityNum,
		"update_time":  tool.Time2Int(),
	}

	_, err = msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotID, RuleType: ruleType})
	}
	return err

}

// GetRobotReceivedMessageReply Get a single received message reply rule
func GetRobotReceivedMessageReply(id int, robotID int) (msql.Params, error) {
	return msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Find()
}

// GetRobotReceivedMessageReplyListWithFilter Get received message reply rule list (with filters and pagination)
func GetRobotReceivedMessageReplyListWithFilter(robotID int, ruleType, replyType string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID))

	// Query by rule type
	if ruleType != "" {
		model.Where("rule_type", ruleType)
	}

	// Query by reply type
	if replyType != "" {
		// Query using PostgreSQL's jsonb operator
		model.Where("reply_type ?| array['" + replyType + "']")
	}

	// Add pagination
	list, total, err := model.Order("priority_num ASC,id DESC").Paginate(page, size)
	return list, total, err
}

// UpdateRobotReceivedMessageReplySwitchStatus Update received message reply rule switch status
func UpdateRobotReceivedMessageReplySwitchStatus(id, robotID, switchStatus int) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"is_repeat":     false,
		"switch_status": switchStatus,
		"error_msg":     ``,
		"message_type":  MessageTypeAll,
	}

	data := msql.Datas{
		"switch_status": switchStatus,
		"update_time":   tool.Time2Int(),
	}

	if switchStatus == define.SwitchOn {
		// Get message
		reply, err := GetRobotReceivedMessageReply(id, robotID)
		if err != nil {
			result["is_repeat"] = true
			result["error_msg"] = err.Error()
			return result, err
		}
		if len(reply) == 0 {
			result["is_repeat"] = true
			result["error_msg"] = `received_message_reply_rule_not_exist`
			return result, err
		}
		var specifyMessageType []string
		specifyMessageType = DisposeStringList(reply[`specify_message_type`])

		result = CheckSpecifyMessageTypeRepeatedlyEnable(reply[`rule_type`], cast.ToInt(reply[`message_type`]), specifyMessageType, robotID, id)
		if cast.ToBool(result[`is_repeat`]) {
			result["switch_status"] = switchStatus
			return result, nil
		}
	}

	_, err := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotID, RuleType: RuleTypeMessageType})
		lib_redis.DelCacheData(define.Redis, &RobotReceivedMessageReplyCacheBuildHandler{RobotId: robotID, RuleType: RuleTypeDuration})
	}
	return result, err
}

// CheckSpecifyMessageTypeRepeatedlyEnable Check if specified message type is repeatedly enabled
func CheckSpecifyMessageTypeRepeatedlyEnable(ruleType string, messageType int, specifyMessageType []string, robotID, id int) map[string]interface{} {

	result := map[string]interface{}{
		"is_repeat":            false,
		"error_msg":            ``,
		"rule_type":            ruleType,
		"message_type":         MessageTypeAll,
		"specify_message_type": specifyMessageType,
	}
	if ruleType != RuleTypeMessageType {
		return result
	}
	if messageType == MessageTypeAll {
		// Query if there are enabled message types, can prompt that it cannot be turned on because the same message type exists in default replies
		checkUseRule, err := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID)).Where("rule_type", RuleTypeMessageType).Where("switch_status", cast.ToString(define.SwitchOn)).Where("id", "!=", cast.ToString(id)).Find()
		if err != nil {
			result["is_repeat"] = true
			result["error_msg"] = err.Error()
			return result
		}
		if len(checkUseRule) > 0 {
			result["is_repeat"] = true
			result["error_msg"] = "received_message_reply_rule_repeat_enable"
			result["message_type"] = checkUseRule["message_type"]
			result["specify_message_type"] = DisposeStringList(checkUseRule["specify_message_type"])
			return result
		}
	} else {
		// Specify message type
		// Check if there is a fully enabled type
		checkUseRule, err := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID)).Where("rule_type", RuleTypeMessageType).Where("switch_status", cast.ToString(define.SwitchOn)).Where("message_type", cast.ToString(MessageTypeAll)).Where("id", "!=", cast.ToString(id)).Find()
		if err != nil {
			result["is_repeat"] = true
			result["error_msg"] = err.Error()
			return result
		}
		if len(checkUseRule) > 0 {
			result["is_repeat"] = true
			result["error_msg"] = "received_message_reply_rule_repeat_enable"
			result["message_type"] = checkUseRule["message_type"]
			result["specify_message_type"] = DisposeStringList(checkUseRule["specify_message_type"])
			return result
		}

		// Check if SpecifyMessageType is duplicated
		model := msql.Model(`func_chat_robot_received_message_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID)).Where("rule_type", RuleTypeMessageType).Where("switch_status", cast.ToString(define.SwitchOn)).Where("id", "!=", cast.ToString(id))
		if len(specifyMessageType) > 0 {
			model.Where("specify_message_type ?| array['" + strings.Join(specifyMessageType, `','`) + "']")
		}
		checkUseRule, err = model.Find()
		if err != nil {
			result["is_repeat"] = true
			result["error_msg"] = err.Error()
			return result
		}
		if len(checkUseRule) > 0 {
			result["is_repeat"] = true
			result["error_msg"] = "received_message_reply_rule_repeat_enable"
			result["message_type"] = checkUseRule["message_type"]
			result["specify_message_type"] = DisposeStringList(checkUseRule["specify_message_type"])
			return result
		}
	}
	return result
}
