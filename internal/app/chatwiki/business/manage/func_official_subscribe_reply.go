// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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

// SubscribeReplyRequest generic request structure
type SubscribeReplyRequest struct {
	ID              int    `form:"id" json:"id"`
	Appid           string `form:"appid" json:"appid"`
	RuleType        string `form:"rule_type" json:"rule_type"`
	DurationType    string `form:"duration_type" json:"duration_type"`
	WeekDuration    string `form:"week_duration" json:"week_duration"`
	StartDay        string `form:"start_day" json:"start_day"`
	EndDay          string `form:"end_day" json:"end_day"`
	StartDuration   string `form:"start_duration" json:"start_duration"`
	EndDuration     string `form:"end_duration" json:"end_duration"`
	PriorityNum     int    `form:"priority_num" json:"priority_num"`
	ReplyInterval   int    `form:"reply_interval" json:"reply_interval"`
	SubscribeSource string `form:"subscribe_source" json:"subscribe_source"`
	ReplyContent    string `form:"reply_content" json:"reply_content"`
	ReplyNum        int    `form:"reply_num" json:"reply_num"`
	SwitchStatus    int    `form:"switch_status" json:"switch_status"`
}

// SubscribeReplySwitchStatusRequest switch status request structure
type SubscribeReplySwitchStatusRequest struct {
	ID           int    `form:"id" json:"id" binding:"required"`
	Appid        string `json:"appid" form:"appid"`
	SwitchStatus int    `form:"switch_status" json:"switch_status"`
}

// SubscribeReplyListFilterRequest list filter request structure
type SubscribeReplyListFilterRequest struct {
	Appid     string `json:"appid" form:"appid"`
	RuleType  string `json:"rule_type" form:"rule_type"`
	ReplyType string `json:"reply_type" form:"reply_type"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// SaveRobotSubscribeReply saves subscribe reply rule (create or update)
func SaveRobotSubscribeReply(c *gin.Context) {
	var req SubscribeReplyRequest

	// Get parameters
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Parse WeekDuration
	var weekDuration []int
	if req.WeekDuration != "" {
		weekDurationStr := strings.Split(req.WeekDuration, ",")
		for _, v := range weekDurationStr {
			weekDuration = append(weekDuration, cast.ToInt(strings.TrimSpace(v)))
		}
	}

	// Parse SubscribeSource
	var subscribeSource []string
	if req.SubscribeSource != "" {
		subscribeSource = strings.Split(req.SubscribeSource, ",")
		// Remove spaces
		for i, v := range subscribeSource {
			subscribeSource[i] = strings.TrimSpace(v)
		}
	}

	// Parse ReplyContent
	var replyContent []common.ReplyContent
	if req.ReplyContent != "" {
		if err := tool.JsonDecodeUseNumber(req.ReplyContent, &replyContent); err != nil {
			common.FmtError(c, `param_err`, "invalid reply_content format")
			return
		}
	}

	// Parse ReplyType
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

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// If updating, check if rule exists and belongs to current user and robot
	if req.ID > 0 {
		ruleInfo, err := common.GetRobotSubscribeReply(req.ID, adminUserId)
		if err != nil {
			logs.Error("Get robot subscribe reply error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}

		if len(ruleInfo) == 0 {
			common.FmtError(c, `rule_not_exist`)
			return
		}
	}

	// Check type
	if req.SwitchStatus == define.SwitchOn {
		result := common.CheckSubscribeSourceRepeatedlyEnable(req.RuleType, subscribeSource, req.Appid, req.ID)
		if cast.ToBool(result["is_repeat"]) {
			errorMsg := cast.ToString(result["error_msg"])
			subscribeSourceName := cast.ToString(result["subscribe_source_name"])
			logs.Error(errorMsg + ` ` + subscribeSourceName)
			common.FmtError(c, errorMsg, subscribeSourceName)
			return
		}
	}

	// Save rule (create or update)
	id, err := common.SaveRobotSubscribeReply(
		req.ID,
		adminUserId,
		req.Appid,
		req.RuleType,
		req.DurationType,
		weekDuration, // Use parsed values
		req.StartDay,
		req.EndDay,
		req.StartDuration,
		req.EndDuration,
		req.PriorityNum,
		req.ReplyInterval,
		subscribeSource, // Use parsed values
		replyContent,    // Use parsed values
		replyType,       // Use parsed values
		req.ReplyNum,
		req.SwitchStatus,
	)

	if err != nil {
		logs.Error("SaveRobotSubscribeReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteRobotSubscribeReply deletes subscribe reply rule
func DeleteRobotSubscribeReply(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))

	// Check parameters
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Delete rule
	err := common.DeleteRobotSubscribeReply(id, adminUserId)
	if err != nil {
		logs.Error("DeleteRobotSubscribeReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetRobotSubscribeReply gets single subscribe reply rule
func GetRobotSubscribeReply(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))

	// Check parameters
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get rule info
	ruleInfo, err := common.GetRobotSubscribeReply(id, adminUserId)
	if err != nil {
		logs.Error("Get robot subscribe reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	// Process return data
	var weekDuration []int
	var replyType, subscribeSource []string
	var replyContent []common.ReplyContent

	// Parse JSON data
	if ruleInfo["week_duration"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["week_duration"], &weekDuration)
	}

	if ruleInfo["reply_type"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_type"], &replyType)
	}

	if ruleInfo["subscribe_source"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["subscribe_source"], &subscribeSource)
	}

	if ruleInfo["reply_content"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_content"], &replyContent)
	}
	// Format smart menu message
	replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilitySubscribeReply)

	result := map[string]interface{}{
		"id":               ruleInfo["id"],
		"appid":            ruleInfo["appid"],
		"rule_type":        ruleInfo["rule_type"],
		"duration_type":    ruleInfo["duration_type"],
		"week_duration":    weekDuration,
		"start_day":        ruleInfo["start_day"],
		"end_day":          ruleInfo["end_day"],
		"start_duration":   ruleInfo["start_duration"],
		"end_duration":     ruleInfo["end_duration"],
		"priority_num":     ruleInfo["priority_num"],
		"reply_interval":   ruleInfo["reply_interval"],
		"subscribe_source": subscribeSource,
		"reply_content":    replyContent,
		"reply_type":       replyType,
		"switch_status":    ruleInfo["switch_status"],
		"reply_num":        ruleInfo["reply_num"],
		"create_time":      ruleInfo["create_time"],
		"update_time":      ruleInfo["update_time"],
	}

	common.FmtOk(c, result)
}

// GetRobotSubscribeReplyList gets subscribe reply rule list
func GetRobotSubscribeReplyList(c *gin.Context) {
	var req SubscribeReplyListFilterRequest

	// Get parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Set default pagination parameters
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get rule list
	list, total, err := common.GetRobotSubscribeReplyListWithFilter(adminUserId, req.Appid, req.RuleType, req.ReplyType, req.Page, req.Size)
	if err != nil {
		logs.Error("GetRobotSubscribeReplyListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	defaultPriorityNum := 1

	// Process return data
	var result []map[string]interface{}
	for _, item := range list {
		var weekDuration []int
		var replyType, subscribeSource []string
		var replyContent []common.ReplyContent

		priorityNum := cast.ToInt(item["priority_num"])
		if req.RuleType == define.RuleTypeSubscribeDefault {
			// Default type
			if priorityNum != defaultPriorityNum {
				priorityNum = defaultPriorityNum
				// Update sort order
				_ = common.UpdateRobotSubscribeReplyPriorityNum(cast.ToInt(item["id"]), adminUserId, priorityNum)
			}
			defaultPriorityNum++
		}
		// Parse JSON data
		if item["week_duration"] != "" {
			_ = tool.JsonDecodeUseNumber(item["week_duration"], &weekDuration)
		}

		if item["reply_type"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_type"], &replyType)
		}

		if item["subscribe_source"] != "" {
			_ = tool.JsonDecodeUseNumber(item["subscribe_source"], &subscribeSource)
		}

		if item["reply_content"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_content"], &replyContent)
		}
		// Format smart menu message
		replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilityReceivedMessageReply)

		result = append(result, map[string]interface{}{
			"id":               item["id"],
			"appid":            item["appid"],
			"rule_type":        item["rule_type"],
			"duration_type":    item["duration_type"],
			"week_duration":    weekDuration,
			"start_day":        item["start_day"],
			"end_day":          item["end_day"],
			"start_duration":   item["start_duration"],
			"end_duration":     item["end_duration"],
			"priority_num":     priorityNum,
			"reply_interval":   item["reply_interval"],
			"subscribe_source": subscribeSource,
			"reply_content":    replyContent,
			"reply_type":       replyType,
			"switch_status":    item["switch_status"],
			"reply_num":        item["reply_num"],
			"create_time":      item["create_time"],
			"update_time":      item["update_time"],
		})
	}

	// Return paginated data
	response := map[string]interface{}{
		"list":  result,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}

// UpdateRobotSubscribeReplySwitchStatus updates subscribe reply rule switch status
func UpdateRobotSubscribeReplySwitchStatus(c *gin.Context) {
	var req SubscribeReplySwitchStatusRequest

	// Get parameters
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Check if switch status parameter is valid (0: off, 1: on)
	if req.SwitchStatus != define.SwitchOff && req.SwitchStatus != define.SwitchOn {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check if rule exists and belongs to current user and robot
	ruleInfo, err := common.GetRobotSubscribeReply(req.ID, adminUserId)
	if err != nil {
		logs.Error("Get robot subscribe reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	// Update switch status
	var result map[string]interface{}
	if req.SwitchStatus == define.SwitchOn {
		// Parse subscribe_source
		var subscribeSource []string
		if ruleInfo["subscribe_source"] != "" {
			_ = tool.JsonDecodeUseNumber(ruleInfo["subscribe_source"], &subscribeSource)
		}

		result = common.CheckSubscribeSourceRepeatedlyEnable(ruleInfo["rule_type"], subscribeSource, ruleInfo["appid"], req.ID)
		if cast.ToBool(result["is_repeat"]) {
			errorMsg := cast.ToString(result["error_msg"])
			subscribeSourceName := cast.ToString(result["subscribe_source_name"])
			logs.Error(errorMsg + ` ` + subscribeSourceName)
			common.FmtError(c, errorMsg, subscribeSourceName)
			return
		}
	}

	err = common.UpdateRobotSubscribeReplySwitchStatus(req.ID, adminUserId, req.SwitchStatus)
	if err != nil {
		logs.Error("UpdateRobotSubscribeReplySwitchStatus error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"switch_status": req.SwitchStatus})
}

// UpdateRobotSubscribeReplyPriorityNum updates subscribe reply rule priority
func UpdateRobotSubscribeReplyPriorityNum(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	appid := cast.ToString(c.PostForm("appid"))
	priorityNum := cast.ToInt(c.PostForm("priority_num"))

	// Check parameters
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if appid == `` {
		common.FmtError(c, `param_lack`, `appid`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check if rule exists and belongs to current user and robot
	ruleInfo, err := common.GetRobotSubscribeReply(id, adminUserId)
	if err != nil {
		logs.Error("Get robot subscribe reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	// Update rule priority
	err = common.UpdateRobotSubscribeReplyPriorityNum(id, adminUserId, priorityNum)
	if err != nil {
		logs.Error("UpdateRobotSubscribeReplyPriorityNum error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}
