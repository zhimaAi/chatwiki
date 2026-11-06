// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"net/http"
	"strings"
	"time"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

// answerSourceAccess is resolved from a verified JWT or an existing visitor's
// message. A public robot key alone never authorizes a message or a document.
type answerSourceAccess struct {
	AdminUserID int
	RobotID     int
	OpenID      string
	MessageID   int
	AdminAccess bool
}

func getAnswerSourceAccess(c *gin.Context, messageID int) (answerSourceAccess, bool) {
	var access answerSourceAccess
	robotKey := strings.TrimSpace(c.Query("robot_key"))
	if !common.CheckRobotKey(robotKey) {
		common.FmtError(c, "param_invalid", "robot_key")
		return access, false
	}
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	var loginAdminID int
	if token != "" {
		loginAdminID = GetAdminUserId(c)
		if loginAdminID <= 0 {
			return access, false
		}
		access.AdminAccess = true
	}
	robot, err := msql.Model("chat_ai_robot", define.Postgres).Where("robot_key", robotKey).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err")
		return access, false
	}
	adminID := cast.ToInt(robot["admin_user_id"])
	if adminID <= 0 || (access.AdminAccess && loginAdminID != adminID) {
		common.FmtError(c, "no_data")
		return access, false
	}
	m := msql.Model("chat_ai_message", define.Postgres).Alias("m").
		Where("m.id", cast.ToString(messageID)).Where("m.admin_user_id", cast.ToString(adminID)).
		Where("m.robot_id", robot["id"]).Where("m.is_customer", "0")
	if !access.AdminAccess {
		openid := strings.TrimSpace(c.Query("openid"))
		if !common.IsXkfOpenid(openid) || !cast.ToBool(robot["answer_source_switch"]) {
			common.FmtError(c, "no_data")
			return access, false
		}
		m.Join("chat_ai_dialogue g", "m.dialogue_id = g.id and m.admin_user_id = g.admin_user_id and m.robot_id = g.robot_id and m.openid = g.openid", "inner").
			Where("m.openid", openid).Where("g.is_background", "0")
	}
	message, err := m.Field("m.id,m.openid").Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err")
		return access, false
	}
	if len(message) == 0 {
		common.FmtError(c, "no_data")
		return access, false
	}
	access.AdminUserID, access.RobotID = adminID, cast.ToInt(robot["id"])
	access.OpenID, access.MessageID = message["openid"], messageID
	return access, true
}

// answerSourceQuery binds both ends of the stored reference to the same tenant
// and file. Checking only c.file_id would allow inconsistent joins to grant access.
func answerSourceQuery(adminID, messageID, fileID int) *msql.Builder {
	return msql.Model("chat_ai_answer_source", define.Postgres).Alias("c").
		Join("chat_ai_library_file_data d", "c.paragraph_id = d.id and c.file_id = d.file_id and c.admin_user_id = d.admin_user_id", "inner").
		Join("chat_ai_library_file f", "d.file_id = f.id and d.admin_user_id = f.admin_user_id", "inner").
		Where("c.admin_user_id", cast.ToString(adminID)).Where("c.message_id", cast.ToString(messageID)).
		Where("c.file_id", cast.ToString(fileID)).Where("d.delete_time", "0").Where("f.delete_time", "0")
}

func getLibFilePreviewInfo(preview common.LibFilePagePreview) (msql.Params, error) {
	m := answerSourceQuery(preview.AdminUserID, preview.MessageID, preview.FileID).
		Join("chat_ai_message m", "c.message_id = m.id and c.admin_user_id = m.admin_user_id", "inner").
		Join("chat_ai_robot r", "m.robot_id = r.id and m.admin_user_id = r.admin_user_id", "inner").
		Where("m.robot_id", cast.ToString(preview.RobotID)).Where("m.openid", preview.OpenID).
		Where("m.is_customer", "0").Where("c.paragraph_id", cast.ToString(preview.ParagraphID)).
		Where("d.page_num", cast.ToString(preview.Page)).Where("f.file_ext", "pdf")
	if !preview.AdminAccess {
		m.Join("chat_ai_dialogue g", "m.dialogue_id = g.id and m.admin_user_id = g.admin_user_id and m.robot_id = g.robot_id and m.openid = g.openid", "inner").
			Where("g.is_background", "0").Where("r.answer_source_switch", "1")
	}
	return m.Field("f.*").Find()
}

func GetChatLibRawFileOnePage(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	preview, err := common.ParseLibFilePagePreview(c.Query("preview_token"), time.Now())
	if err != nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, "auth_no_permission")
		return
	}
	info, err := getLibFilePreviewInfo(preview)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, "sys_err")
		return
	}
	if len(info) == 0 {
		common.FmtError(c, "no_data")
		return
	}
	// File and page come only from the verified ticket, never other query fields.
	serveLibRawFileOnePage(c, info, preview.Page)
}
