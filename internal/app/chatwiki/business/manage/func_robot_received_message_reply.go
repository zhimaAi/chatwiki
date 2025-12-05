// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// ReceivedMessageReplyRequest 通用请求结构
type ReceivedMessageReplyRequest struct {
	ID                 int    `form:"id" json:"id"`
	RobotID            int    `form:"robot_id" json:"robot_id" binding:"required"`
	RuleType           string `form:"rule_type" json:"rule_type"`
	DurationType       string `form:"duration_type" json:"duration_type"`
	WeekDuration       string `form:"week_duration" json:"week_duration"`
	StartDay           string `form:"start_day" json:"start_day"`
	EndDay             string `form:"end_day" json:"end_day"`
	StartDuration      string `form:"start_duration" json:"start_duration"`
	EndDuration        string `form:"end_duration" json:"end_duration"`
	PriorityNum        int    `form:"priority_num" json:"priority_num"`
	ReplyInterval      int    `form:"reply_interval" json:"reply_interval"`
	MessageType        int    `form:"message_type" json:"message_type"`
	SpecifyMessageType string `form:"specify_message_type" json:"specify_message_type"`
	ReplyContent       string `form:"reply_content" json:"reply_content"`
	ReplyNum           int    `form:"reply_num" json:"reply_num"`
	SwitchStatus       int    `form:"switch_status" json:"switch_status"`
}

// ReceivedMessageReplySwitchStatusRequest 开关状态请求结构
type ReceivedMessageReplySwitchStatusRequest struct {
	ID           int `form:"id" json:"id" binding:"required"`
	RobotID      int `form:"robot_id" json:"robot_id" binding:"required"`
	SwitchStatus int `form:"switch_status" json:"switch_status"`
}

// ReceivedMessageReplyListFilterRequest 列表过滤请求结构
type ReceivedMessageReplyListFilterRequest struct {
	RobotID   int    `json:"robot_id" form:"robot_id" binding:"required"`
	RuleType  string `json:"rule_type" form:"rule_type"`
	ReplyType string `json:"reply_type" form:"reply_type"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// SaveRobotReceivedMessageReply 保存收到消息回复规则（创建或更新）
func SaveRobotReceivedMessageReply(c *gin.Context) {
	var req ReceivedMessageReplyRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 解析 WeekDuration
	var weekDuration []int
	if req.WeekDuration != "" {
		weekDurationStr := strings.Split(req.WeekDuration, ",")
		for _, v := range weekDurationStr {
			weekDuration = append(weekDuration, cast.ToInt(strings.TrimSpace(v)))
		}
	}

	// 解析 SpecifyMessageType
	var specifyMessageType []string
	if req.SpecifyMessageType != "" {
		specifyMessageType = strings.Split(req.SpecifyMessageType, ",")
		// 去除空格
		for i, v := range specifyMessageType {
			specifyMessageType[i] = strings.TrimSpace(v)
		}
	}

	// 解析 ReplyContent
	var replyContent []common.ReplyContent
	if req.ReplyContent != "" {
		if err := tool.JsonDecodeUseNumber(req.ReplyContent, &replyContent); err != nil {
			common.FmtError(c, `param_err`, "invalid reply_content format")
			return
		}
	}

	// 解析 ReplyType
	var replyType []string
	if len(replyContent) > 0 {
		for _, reply := range replyContent {
			if reply.ReplyType == `` && reply.Type != `` {
				replyType = append(replyType, reply.Type)
			} else {
				replyType = append(replyType, reply.ReplyType)
			}
		}
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// 如果是更新操作，检查规则是否存在且属于当前用户和机器人
	if req.ID > 0 {
		ruleInfo, err := common.GetRobotReceivedMessageReply(req.ID, req.RobotID)
		if err != nil {
			logs.Error("Get robot received message reply error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}

		if len(ruleInfo) == 0 {
			common.FmtError(c, `rule_not_exist`)
			return
		}

		if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId || cast.ToInt(ruleInfo["robot_id"]) != req.RobotID {
			common.FmtError(c, `auth_no_permission`)
			return
		}
	}
	//检测类型
	if req.SwitchStatus == define.SwitchOn {
		result := common.CheckSpecifyMessageTypeRepeatedlyEnable(req.RuleType, req.MessageType, specifyMessageType, req.RobotID, req.ID)
		if cast.ToBool(result["is_repeat"]) {
			common.FmtError(c, cast.ToString(result["error_msg"]))
			return
		}
	}

	// 保存规则（创建或更新）
	id, err := common.SaveRobotReceivedMessageReply(
		req.ID,
		adminUserId,
		req.RobotID,
		req.RuleType,
		req.DurationType,
		weekDuration, // 使用解析后的值
		req.StartDay,
		req.EndDay,
		req.StartDuration,
		req.EndDuration,
		req.PriorityNum,
		req.ReplyInterval,
		req.MessageType,
		specifyMessageType, // 使用解析后的值
		replyContent,       // 使用解析后的值
		replyType,          // 使用解析后的值
		req.ReplyNum,
		req.SwitchStatus,
	)

	if err != nil {
		logs.Error("SaveRobotReceivedMessageReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteRobotReceivedMessageReply 删除收到消息回复规则
func DeleteRobotReceivedMessageReply(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	robotID := cast.ToInt(c.PostForm("robot_id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// 检查规则是否存在且属于当前用户和机器人
	ruleInfo, err := common.GetRobotReceivedMessageReply(id, robotID)
	if err != nil {
		logs.Error("Get robot received message reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId || cast.ToInt(ruleInfo["robot_id"]) != robotID {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 删除规则
	err = common.DeleteRobotReceivedMessageReply(id, robotID)
	if err != nil {
		logs.Error("DeleteRobotReceivedMessageReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetRobotReceivedMessageReply 获取单个收到消息回复规则
func GetRobotReceivedMessageReply(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	robotID := cast.ToInt(c.Query("robot_id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取规则信息
	ruleInfo, err := common.GetRobotReceivedMessageReply(id, robotID)
	if err != nil {
		logs.Error("Get robot received message reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	// 检查权限
	if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 处理返回数据
	var weekDuration []int
	var specifyMessageType, replyType []string
	var replyContent []common.ReplyContent

	// 解析JSON数据
	if ruleInfo["week_duration"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["week_duration"], &weekDuration)
	}

	if ruleInfo["specify_message_type"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["specify_message_type"], &specifyMessageType)
	}

	if ruleInfo["reply_type"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_type"], &replyType)
	}

	if ruleInfo["reply_content"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_content"], &replyContent)
	}
	// 格式化智能菜单消息
	replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilityReceivedMessageReply)

	result := map[string]interface{}{
		"id":                   ruleInfo["id"],
		"robot_id":             ruleInfo["robot_id"],
		"rule_type":            ruleInfo["rule_type"],
		"duration_type":        ruleInfo["duration_type"],
		"week_duration":        weekDuration,
		"start_day":            ruleInfo["start_day"],
		"end_day":              ruleInfo["end_day"],
		"start_duration":       ruleInfo["start_duration"],
		"end_duration":         ruleInfo["end_duration"],
		"priority_num":         ruleInfo["priority_num"],
		"reply_interval":       ruleInfo["reply_interval"],
		"message_type":         ruleInfo["message_type"],
		"specify_message_type": specifyMessageType,
		"reply_content":        replyContent,
		"reply_type":           replyType,
		"switch_status":        ruleInfo["switch_status"],
		"reply_num":            ruleInfo["reply_num"],
		"create_time":          ruleInfo["create_time"],
		"update_time":          ruleInfo["update_time"],
	}

	common.FmtOk(c, result)
}

// GetRobotReceivedMessageReplyList 获取收到消息回复规则列表
func GetRobotReceivedMessageReplyList(c *gin.Context) {
	var req ReceivedMessageReplyListFilterRequest

	// 获取参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// 获取规则列表
	list, total, err := common.GetRobotReceivedMessageReplyListWithFilter(req.RobotID, req.RuleType, req.ReplyType, req.Page, req.Size)
	if err != nil {
		logs.Error("GetRobotReceivedMessageReplyListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// 处理返回数据
	var result []map[string]interface{}
	for _, item := range list {
		var weekDuration []int
		var specifyMessageType, replyType []string
		var replyContent []common.ReplyContent

		// 解析JSON数据
		if item["week_duration"] != "" {
			_ = tool.JsonDecodeUseNumber(item["week_duration"], &weekDuration)
		}

		if item["specify_message_type"] != "" {
			_ = tool.JsonDecodeUseNumber(item["specify_message_type"], &specifyMessageType)
		}

		if item["reply_type"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_type"], &replyType)
		}

		if item["reply_content"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_content"], &replyContent)
		}
		// 格式化智能菜单消息
		replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilityReceivedMessageReply)

		result = append(result, map[string]interface{}{
			"id":                   item["id"],
			"robot_id":             item["robot_id"],
			"rule_type":            item["rule_type"],
			"duration_type":        item["duration_type"],
			"week_duration":        weekDuration,
			"start_day":            item["start_day"],
			"end_day":              item["end_day"],
			"start_duration":       item["start_duration"],
			"end_duration":         item["end_duration"],
			"priority_num":         item["priority_num"],
			"reply_interval":       item["reply_interval"],
			"message_type":         item["message_type"],
			"specify_message_type": specifyMessageType,
			"reply_content":        replyContent,
			"reply_type":           replyType,
			"switch_status":        item["switch_status"],
			"reply_num":            item["reply_num"],
			"create_time":          item["create_time"],
			"update_time":          item["update_time"],
		})
	}

	// 返回分页数据
	response := map[string]interface{}{
		"list":  result,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}

// UpdateRobotReceivedMessageReplySwitchStatus 更新收到消息回复规则开关状态
func UpdateRobotReceivedMessageReplySwitchStatus(c *gin.Context) {
	var req ReceivedMessageReplySwitchStatusRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 检查开关状态参数是否合法 (0:关闭, 1:开启)
	if req.SwitchStatus != define.SwitchOff && req.SwitchStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// 检查规则是否存在且属于当前用户和机器人
	ruleInfo, err := common.GetRobotReceivedMessageReply(req.ID, req.RobotID)
	if err != nil {
		logs.Error("Get robot received message reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId || cast.ToInt(ruleInfo["robot_id"]) != req.RobotID {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 更新开关状态
	var result map[string]interface{}
	result, err = common.UpdateRobotReceivedMessageReplySwitchStatus(req.ID, req.RobotID, req.SwitchStatus)
	if err != nil {
		logs.Error("UpdateRobotReceivedMessageReplySwitchStatus error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if cast.ToBool(result[`is_repeat`]) {
		logs.Error(cast.ToString(result[`error_msg`]))
		common.FmtError(c, cast.ToString(result[`error_msg`]))
		return
	}

	common.FmtOk(c, result)
}

// UpdateRobotReceivedMessageReplyPriorityNum 更新收到消息回复规则优先级
func UpdateRobotReceivedMessageReplyPriorityNum(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	robotID := cast.ToInt(c.PostForm("robot_id"))
	priorityNum := cast.ToInt(c.PostForm("priority_num"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查机器人是否存在且属于当前用户
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// 检查规则是否存在且属于当前用户和机器人
	ruleInfo, err := common.GetRobotReceivedMessageReply(id, robotID)
	if err != nil {
		logs.Error("Get robot received message reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId || cast.ToInt(ruleInfo["robot_id"]) != robotID {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 删除规则
	err = common.UpdateRobotReceivedMessageReplyPriorityNum(id, robotID, priorityNum)
	if err != nil {
		logs.Error("UpdateRobotReceivedMessageReplyPriorityNum error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}
