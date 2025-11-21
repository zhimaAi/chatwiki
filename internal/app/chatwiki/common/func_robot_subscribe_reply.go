// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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

// RobotSubscribeReplyCacheBuildHandler 关注后回复规则缓存
type RobotSubscribeReplyCacheBuildHandler struct {
	RobotId  int
	Appid    string
	RuleType string
}

type SubscribeReplyLastTimeCacheBuildHandler struct {
	RobotId  int
	RuleId   int
	Openid   string
	LastTime int
}

func (h *SubscribeReplyLastTimeCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.func_chat_robot_subscribe_reply_interval.%d.%d.%s`, h.RobotId, h.RuleId, h.Openid)
}

func (h *SubscribeReplyLastTimeCacheBuildHandler) GetCacheData() (any, error) {
	if h.LastTime == 0 {
		h.LastTime = tool.Time2Int()
	}
	return map[string]interface{}{
		"last_time": h.LastTime,
	}, nil
}

// GetSubscribeReplyLastTime 获取关注后回复规则最后回复时间
func GetSubscribeReplyLastTime(robotId, ruleId int, openid string) (int, error) {
	result := map[string]interface{}{
		"last_time": 0,
	}
	err := lib_redis.GetOne(define.Redis, &SubscribeReplyLastTimeCacheBuildHandler{RobotId: robotId, RuleId: ruleId, Openid: openid}, &result)
	if err != nil {
		return 0, err
	}
	return cast.ToInt(result["last_time"]), err
}

func SetSubscribeReplyLastTime(robotId, ruleId int, lastTime int, openid string) error {
	return lib_redis.SetOne(define.Redis, &SubscribeReplyLastTimeCacheBuildHandler{RobotId: robotId, RuleId: ruleId, LastTime: lastTime, Openid: openid}, time.Hour)
}

const (
	SubscribeRuleTypeDefault  = `subscribe_reply_default`  // 默认规则
	SubscribeRuleTypeSource   = `subscribe_reply_source`   // 来源规则
	SubscribeRuleTypeDuration = `subscribe_reply_duration` // 时间规则
)

type RobotSubscribeReply struct {
	ID                 int            `json:"id"`
	AdminUserID        int            `json:"admin_user_id"`
	RobotID            int            `json:"robot_id"`
	Appid              string         `json:"appid"`
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
	SubscribeSource    []string       `json:"subscribe_source"`
	SwitchStatus       int            `json:"switch_status"`
	ReplyNum           int            `json:"reply_num"`
	CreateTime         int            `json:"create_time"`
	UpdateTime         int            `json:"update_time"`
}

func (h *RobotSubscribeReplyCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.func_chat_robot_subscribe_reply.%d.%s.%s`, h.RobotId, h.Appid, h.RuleType)
}

func (h *RobotSubscribeReplyCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where(`robot_id`, cast.ToString(h.RobotId)).Where(`appid`, h.Appid).Where(`rule_type`, h.RuleType).Where(`switch_status`, cast.ToString(define.SwitchOn)).Order(`priority_num asc,id desc`).Select()
	//转换
	list := make([]RobotSubscribeReply, 0)
	if err == nil && len(data) > 0 {
		for _, item := range data {
			//
			var replyContent = make([]ReplyContent, 0)
			_ = tool.JsonDecodeUseNumber(item[`reply_content`], &replyContent)
			if len(replyContent) == 0 {
				continue
			}

			// 解析 week_duration 字段
			var weekDuration = make([]int, 0)
			_ = tool.JsonDecodeUseNumber(item[`week_duration`], &weekDuration)

			// 解析 specify_message_type 字段
			var specifyMessageType = make([]string, 0)
			_ = tool.JsonDecodeUseNumber(item[`specify_message_type`], &specifyMessageType)

			// 解析 reply_type 字段
			var replyType = make([]string, 0)
			_ = tool.JsonDecodeUseNumber(item[`reply_type`], &replyType)

			// 解析 subscribe_source 字段
			var subscribeSource = make([]string, 0)
			_ = tool.JsonDecodeUseNumber(item[`subscribe_source`], &subscribeSource)

			list = append(list, RobotSubscribeReply{
				ID:                 cast.ToInt(item[`id`]),
				AdminUserID:        cast.ToInt(item[`admin_user_id`]),
				RobotID:            cast.ToInt(item[`robot_id`]),
				Appid:              item[`appid`],
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
				SubscribeSource:    subscribeSource,
				SwitchStatus:       cast.ToInt(item[`switch_status`]),
				ReplyNum:           cast.ToInt(item[`reply_num`]),
				CreateTime:         cast.ToInt(item[`create_time`]),
				UpdateTime:         cast.ToInt(item[`update_time`]),
			})
		}
	}
	return list, err
}

// GetRobotSubscribeReplyListByAppid 公众号关注后回复规则列表
func GetRobotSubscribeReplyListByAppid(robotId int, appid, ruleType string) ([]RobotSubscribeReply, error) {
	result := make([]RobotSubscribeReply, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &RobotSubscribeReplyCacheBuildHandler{RobotId: robotId, Appid: appid, RuleType: ruleType}, &result, time.Hour)
	if err != nil {
		return nil, err
	}
	return result, err
}

// SaveRobotSubscribeReply 保存关注后回复规则（创建或更新）
func SaveRobotSubscribeReply(id, adminUserID, robotID int, appid, ruleType, durationType string, weekDuration []int, startDay, endDay, startDuration, endDuration string, priorityNum, replyInterval, messageType int, specifyMessageType, subscribeSource []string, replyContent []ReplyContent, replyType []string, replyNum int, switchStatus int) (int64, error) {
	weekDurationJson, _ := tool.JsonEncode(weekDuration)
	specifyMessageTypeJson, _ := tool.JsonEncode(specifyMessageType)
	subscribeSourceJson, _ := tool.JsonEncode(subscribeSource)
	replyContentJson, _ := tool.JsonEncode(replyContent)
	replyTypeJson, _ := tool.JsonEncode(replyType)

	data := msql.Datas{
		"admin_user_id":        adminUserID,
		"robot_id":             robotID,
		"appid":                appid,
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
		"subscribe_source":     subscribeSourceJson,
		"reply_content":        replyContentJson,
		"reply_type":           replyTypeJson,
		"reply_num":            replyNum,
		"switch_status":        switchStatus,
		"update_time":          tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// 创建新记录
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Insert(data, "id")
	} else {
		// 更新现有记录
		_, err = msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSubscribeReplyCacheBuildHandler{RobotId: robotID, Appid: appid, RuleType: ruleType})
	}
	return newId, err
}

// DeleteRobotSubscribeReply 删除关注后回复规则
func DeleteRobotSubscribeReply(id, robotID int) error {
	//先获取删除的记录
	oldOne, err := GetRobotSubscribeReply(id, robotID)
	if err != nil {
		return err
	}
	appid := cast.ToString(oldOne[`appid`])
	ruleType := cast.ToString(oldOne[`rule_type`])
	_, err = msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSubscribeReplyCacheBuildHandler{RobotId: robotID, Appid: appid, RuleType: ruleType})
	}
	return err
}

// GetRobotSubscribeReply 获取单个关注后回复规则
func GetRobotSubscribeReply(id int, robotID int) (msql.Params, error) {
	return msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("id", cast.ToString(id)).Where("robot_id", cast.ToString(robotID)).Find()
}

// GetRobotSubscribeReplyListWithFilter 获取关注后回复规则列表（带过滤条件和分页）
func GetRobotSubscribeReplyListWithFilter(robotID int, appid, ruleType, replyType string, page, size int) ([]msql.Params, int, error) {
	model := msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("robot_id", cast.ToString(robotID))

	// 根据appid查询
	if appid != "" {
		model.Where("appid", appid)
	}
	// 根据规则类型查询
	if ruleType != "" {
		model.Where("rule_type", ruleType)
	}

	// 根据回复类型查询
	if replyType != "" {
		// 使用 PostgreSQL 的 jsonb 操作符进行查询
		model.Where("reply_type ?| array['" + replyType + "']")
	}

	// 添加分页
	list, total, err := model.Order("id DESC").Paginate(page, size)
	return list, total, err
}

// UpdateRobotSubscribeReplySwitchStatus 更新关注后回复规则开关状态
func UpdateRobotSubscribeReplySwitchStatus(id, robotID, switchStatus int) error {
	//先获取更新的记录
	oldOne, err := GetRobotSubscribeReply(id, robotID)
	if err != nil {
		return err
	}
	appid := cast.ToString(oldOne[`appid`])
	ruleType := cast.ToString(oldOne[`rule_type`])

	data := msql.Datas{
		"switch_status": switchStatus,
		"update_time":   tool.Time2Int(),
	}

	_, err = msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSubscribeReplyCacheBuildHandler{RobotId: robotID, Appid: appid, RuleType: ruleType})
	}
	return err
}

// UpdateRobotSubscribeReplyPriorityNum 更新关注后回复规则优先级
func UpdateRobotSubscribeReplyPriorityNum(id, robotID, priorityNum int) error {
	//先获取更新的记录
	oldOne, err := GetRobotSubscribeReply(id, robotID)
	if err != nil {
		return err
	}
	appid := cast.ToString(oldOne[`appid`])
	ruleType := cast.ToString(oldOne[`rule_type`])

	data := msql.Datas{
		"priority_num": priorityNum,
		"update_time":  tool.Time2Int(),
	}

	_, err = msql.Model(`func_chat_robot_subscribe_reply`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &RobotSubscribeReplyCacheBuildHandler{RobotId: robotID, Appid: appid, RuleType: ruleType})
	}
	return err
}
