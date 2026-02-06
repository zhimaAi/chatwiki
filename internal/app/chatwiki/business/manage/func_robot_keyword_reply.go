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

// KeywordReplyRequest generic request structure
type KeywordReplyRequest struct {
	ID           int    `form:"id" json:"id"`
	RobotID      int    `form:"robot_id" json:"robot_id" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	FullKeyword  string `form:"full_keyword" json:"full_keyword"`
	HalfKeyword  string `form:"half_keyword" json:"half_keyword"`
	ReplyContent string `form:"reply_content" json:"reply_content"`
	ReplyNum     int    `form:"reply_num" json:"reply_num"`
	ForcedEnable int    `form:"forced_enable" json:"forced_enable"`
}

// SwitchStatusRequest switch status request structure
type SwitchStatusRequest struct {
	ID           int `form:"id" json:"id" binding:"required"`
	RobotID      int `form:"robot_id" json:"robot_id" binding:"required"`
	SwitchStatus int `form:"switch_status" json:"switch_status"`
}

// ListFilterRequest list filter request structure
type ListFilterRequest struct {
	RobotID   int    `json:"robot_id" form:"robot_id" binding:"required"`
	Keyword   string `json:"keyword" form:"keyword"`
	RuleName  string `json:"rule_name" form:"rule_name"`
	ReplyType string `json:"reply_type" form:"reply_type"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// SaveRobotKeywordReply saves keyword reply rule (create or update)
func SaveRobotKeywordReply(c *gin.Context) {
	var req KeywordReplyRequest

	// Get parameters
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Parse FullKeyword
	var fullKeyword []string
	if req.FullKeyword != "" {
		fullKeyword = strings.Split(req.FullKeyword, ",")
		// Remove spaces
		for i, v := range fullKeyword {
			fullKeyword[i] = strings.TrimSpace(v)
		}
	}
	if len(fullKeyword) > common.MaxKeywordNum {
		common.FmtError(c, `robot_keyword_reply_max_num`, cast.ToString(common.MaxKeywordNum))
		return
	}

	// Parse HalfKeyword
	var halfKeyword []string
	if req.HalfKeyword != "" {
		halfKeyword = strings.Split(req.HalfKeyword, ",")
		// Remove spaces
		for i, v := range halfKeyword {
			halfKeyword[i] = strings.TrimSpace(v)
		}
	}

	if len(halfKeyword) > common.MaxKeywordNum {
		common.FmtError(c, `robot_keyword_reply_max_num`, cast.ToString(common.MaxKeywordNum))
		return
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
		// Backend merges automatically
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

	// Check if robot exists and belongs to current user
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// If updating, check if rule exists and belongs to current user and robot
	if req.ID > 0 {
		ruleInfo, err := common.GetRobotKeywordReply(req.ID)
		if err != nil {
			logs.Error("Get robot keyword reply error: %s", err.Error())
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

	// Save rule (create or update)
	id, err := common.SaveRobotKeywordReply(
		req.ID,
		adminUserId,
		req.RobotID,
		req.Name,
		fullKeyword,  // Use parsed values
		halfKeyword,  // Use parsed values
		replyContent, // Use parsed values
		replyType,    // Use parsed values
		req.ReplyNum,
	)

	if err != nil {
		logs.Error("SaveRobotKeywordReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if req.ForcedEnable != 0 {
		err = common.UpdateRobotKeywordReplySwitchStatus(cast.ToInt(id), req.RobotID, req.ForcedEnable)
		if err != nil {
			logs.Error("UpdateRobotKeywordReplySwitchStatus error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}
	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteRobotKeywordReply deletes keyword reply rule
func DeleteRobotKeywordReply(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))
	robotID := cast.ToInt(c.PostForm("robot_id"))

	// Check parameters
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check if robot exists and belongs to current user
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// Check if rule exists and belongs to current user and robot
	ruleInfo, err := common.GetRobotKeywordReply(id)
	if err != nil {
		logs.Error("Get robot keyword reply error: %s", err.Error())
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

	// Delete rule
	err = common.DeleteRobotKeywordReply(id, robotID)
	if err != nil {
		logs.Error("DeleteRobotKeywordReply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetRobotKeywordReply gets single keyword reply rule
func GetRobotKeywordReply(c *gin.Context) {
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
	ruleInfo, err := common.GetRobotKeywordReply(id)
	if err != nil {
		logs.Error("Get robot keyword reply error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(ruleInfo) == 0 {
		common.FmtError(c, `rule_not_exist`)
		return
	}

	// Check permission
	if cast.ToInt(ruleInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// Process return data
	var fullKeyword, halfKeyword, replyType []string
	var replyContent []common.ReplyContent

	// Parse JSON data
	if ruleInfo["full_keyword"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["full_keyword"], &fullKeyword)
	}

	if ruleInfo["half_keyword"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["half_keyword"], &halfKeyword)
	}

	if ruleInfo["reply_type"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_type"], &replyType)
	}

	if ruleInfo["reply_content"] != "" {
		_ = tool.JsonDecodeUseNumber(ruleInfo["reply_content"], &replyContent)
	}
	// Format smart menu message
	replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilityReceivedMessageReply)

	result := map[string]interface{}{
		"id":            ruleInfo["id"],
		"robot_id":      ruleInfo["robot_id"],
		"name":          ruleInfo["name"],
		"full_keyword":  fullKeyword,
		"half_keyword":  halfKeyword,
		"reply_content": replyContent,
		"reply_type":    replyType,
		"switch_status": ruleInfo["switch_status"],
		"reply_num":     ruleInfo["reply_num"],
		"create_time":   ruleInfo["create_time"],
		"update_time":   ruleInfo["update_time"],
	}

	common.FmtOk(c, result)
}

// GetRobotKeywordReplyList gets keyword reply rule list
func GetRobotKeywordReplyList(c *gin.Context) {
	var req ListFilterRequest

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

	// Check if robot exists and belongs to current user
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// Get rule list
	list, total, err := common.GetRobotKeywordReplyListWithFilter(req.RobotID, req.Keyword, req.ReplyType, req.Page, req.Size)
	if err != nil {
		logs.Error("GetRobotKeywordReplyListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// Process return data
	var result []map[string]interface{}
	for _, item := range list {
		var fullKeyword, halfKeyword, replyType []string
		var replyContent []common.ReplyContent

		// Parse JSON data
		if item["full_keyword"] != "" {
			_ = tool.JsonDecodeUseNumber(item["full_keyword"], &fullKeyword)
		}

		if item["half_keyword"] != "" {
			_ = tool.JsonDecodeUseNumber(item["half_keyword"], &halfKeyword)
		}

		if item["reply_type"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_type"], &replyType)
		}

		if item["reply_content"] != "" {
			_ = tool.JsonDecodeUseNumber(item["reply_content"], &replyContent)
		}
		// Format smart menu message
		replyContent = common.FormatReplyListToDb(replyContent, common.RobotAbilityKeywordReply)

		result = append(result, map[string]interface{}{
			"id":            item["id"],
			"robot_id":      item["robot_id"],
			"name":          item["name"],
			"full_keyword":  fullKeyword,
			"half_keyword":  halfKeyword,
			"reply_content": replyContent,
			"reply_type":    replyType,
			"switch_status": item["switch_status"],
			"reply_num":     item["reply_num"],
			"create_time":   item["create_time"],
			"update_time":   item["update_time"],
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

// UpdateRobotKeywordReplySwitchStatus updates keyword reply rule switch status
func UpdateRobotKeywordReplySwitchStatus(c *gin.Context) {
	var req SwitchStatusRequest

	// Get parameters
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Check enabled rule count
	useNum, err := common.GetRobotKeywordReplyUseRuleNum(req.RobotID)
	if err != nil {
		logs.Error("GetRobotKeywordReplyUseRuleNum error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if useNum >= common.MaxKeywordReplyRuleNum {
		common.FmtError(c, `robot_keyword_reply_rule_num_limit`, cast.ToString(common.MaxKeywordReplyRuleNum))
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

	// Check if robot exists and belongs to current user
	if checkRobotByAdminUserId(c, req.RobotID, adminUserId) {
		return
	}

	// Check if rule exists and belongs to current user and robot
	ruleInfo, err := common.GetRobotKeywordReply(req.ID)
	if err != nil {
		logs.Error("Get robot keyword reply error: %s", err.Error())
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

	// Update switch status
	err = common.UpdateRobotKeywordReplySwitchStatus(req.ID, req.RobotID, req.SwitchStatus)
	if err != nil {
		logs.Error("UpdateRobotKeywordReplySwitchStatus error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

func CheckKeyWordRepeat(c *gin.Context) {
	id := cast.ToInt(c.DefaultPostForm("id", "0"))
	robotID := cast.ToInt(c.DefaultPostForm("robot_id", "0"))
	keyword := c.DefaultPostForm("keyword", "")

	if robotID <= 0 {
		common.FmtError(c, `param_lack`, `robot_id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check if robot exists and belongs to current user
	if checkRobotByAdminUserId(c, robotID, adminUserId) {
		return
	}

	// Check if keyword is duplicate
	isRepeat, ruleName := common.CheckKeyWordRepeat(robotID, keyword, id)
	// Return response data
	response := map[string]interface{}{
		"is_repeat": isRepeat,
		"rule_name": ruleName,
	}
	common.FmtOk(c, response)
}
