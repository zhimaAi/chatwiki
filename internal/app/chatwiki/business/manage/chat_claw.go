// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type ChatClawLoginReq struct {
	UserName  string `form:"user_name"  json:"user_name"  binding:"required"`
	Password  string `form:"password"   json:"password"   binding:"required"`
	OsType    string `form:"os_type"    json:"os_type"`
	OsVersion string `form:"os_version" json:"os_version"`
}

type ChatClawRobotListReq struct {
	ApplicationType int `form:"application_type"`
	OnlyOpen        int `form:"only_open"`
}

type ChatClawLibraryListReq struct {
	Type        string `form:"type"`
	LibraryName string `form:"library_name"`
	Ids         string `form:"ids"`
	OnlyOpen    int    `form:"only_open"`
}

type ChatClawSwitchReq struct {
	Id           int `form:"id" json:"id"`
	SwitchStatus int `form:"switch_status" json:"switch_status"`
}

type ChatClawForceOfflineReq struct {
	Id     int    `form:"id" json:"id"`
	Reason string `form:"reason" json:"reason"`
}

type ChatClawRefreshTokenReq struct {
	OsType    string `form:"os_type" json:"os_type"`
	OsVersion string `form:"os_version" json:"os_version"`
}

func ChatClawLogin(c *gin.Context) {
	var req ChatClawLoginReq
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	info, err := msql.Model(define.TableUser, define.Postgres).Where(`user_name`, req.UserName).Where("is_deleted", define.Normal).
		Where(fmt.Sprintf(`password=MD5(concat(%s,salt))`, msql.ToString(req.Password))).Field(`id,user_name,user_roles,avatar,nick_name,parent_id`).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `user_or_pwd_err`)
		return
	}
	data, err := common.GetChatClawToken(info[`id`], info[`user_name`], info["parent_id"])
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	serverIP := lib_web.GetClientIP(c)
	m := msql.Model(define.TableUser, define.Postgres)
	_, err = m.Where("id", info[`id`]).Update(msql.Datas{
		"login_time": time.Now().Unix(),
		"login_ip":   serverIP,
	})
	if err != nil {
		logs.Error(err.Error())
	}

	// Record token issuance log. The client IP is derived from request metadata.
	_, logErr := msql.Model(define.TableChatClawTokenLog, define.Postgres).Insert(msql.Datas{
		"user_id":     cast.ToInt(info[`id`]),
		"token":       cast.ToString(data[`token`]),
		"token_hash":  common.GetTokenSha256(cast.ToString(data[`token`])),
		"os_type":     req.OsType,
		"os_version":  req.OsVersion,
		"client_ip":   serverIP,
		"expired_at":  cast.ToInt64(data[`exp`]),
		"status":      common.ChatClawTokenStatusActive,
		"create_time": tool.Time2Int(),
	})
	if logErr != nil {
		logs.Error("chatclaw token log insert error: %s", logErr.Error())
	}

	data["user_roles"] = info["user_roles"]
	data["avatar"] = info["avatar"]
	data["nick_name"] = info["nick_name"]
	common.FmtOk(c, data)
}

func getChatClawAuthUser(c *gin.Context) (int, int, bool) {
	claims, _, err := common.GetChatClawAuthClaims(c)
	if err != nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return 0, 0, false
	}
	userId := cast.ToInt(claims[`user_id`])
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return 0, 0, false
	}
	adminUserId := cast.ToInt(claims[`parent_id`])
	if adminUserId <= 0 {
		adminUserId = userId
	}
	return adminUserId, userId, true
}

func ChatClawGetRobotList(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawRobotListReq{
		ApplicationType: cast.ToInt(c.DefaultQuery(`application_type`, `-1`)),
		OnlyOpen:        cast.ToInt(c.DefaultQuery(`only_open`, `0`)),
	}
	if !tool.InArrayInt(req.OnlyOpen, []int{define.SwitchOff, define.SwitchOn}) {
		common.FmtError(c, `param_invalid`, `only_open`)
		return
	}
	list, httpStatus, err := BridgeGetRobotListWithOption(adminUserId, userId, common.GetLang(c), req.ApplicationType, &BridgeRobotListOption{
		OnlyOpen:            req.OnlyOpen,
		IncludeChatClawInfo: true,
	})
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func ChatClawGetEnabledRobotList(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawRobotListReq{
		ApplicationType: cast.ToInt(c.DefaultQuery(`application_type`, `-1`)),
		OnlyOpen:        define.SwitchOn,
	}
	list, httpStatus, err := BridgeGetRobotListWithOption(adminUserId, userId, common.GetLang(c), req.ApplicationType, &BridgeRobotListOption{
		OnlyOpen:            req.OnlyOpen,
		IncludeChatClawInfo: true,
	})
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func ChatClawGetLibraryList(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawLibraryListReq{
		Type:        c.DefaultQuery(`type`, ``),
		LibraryName: c.DefaultQuery(`library_name`, ``),
		Ids:         c.DefaultQuery(`ids`, ``),
		OnlyOpen:    cast.ToInt(c.DefaultQuery(`only_open`, `0`)),
	}
	list, httpStatus, err := BridgeGetChatClawLibraryList(adminUserId, userId, common.GetLang(c), &BridgeChatClawLibraryListReq{
		Type:        req.Type,
		LibraryName: req.LibraryName,
		Ids:         req.Ids,
		OnlyOpen:    req.OnlyOpen,
	})
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}
	common.FmtOk(c, list)
}

func ChatClawGetEnabledLibraryList(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawLibraryListReq{
		Type:        c.DefaultQuery(`type`, ``),
		LibraryName: c.DefaultQuery(`library_name`, ``),
		Ids:         c.DefaultQuery(`ids`, ``),
		OnlyOpen:    define.SwitchOn,
	}
	list, httpStatus, err := BridgeGetChatClawLibraryList(adminUserId, userId, common.GetLang(c), &BridgeChatClawLibraryListReq{
		Type:        req.Type,
		LibraryName: req.LibraryName,
		Ids:         req.Ids,
		OnlyOpen:    req.OnlyOpen,
	})
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}
	common.FmtOk(c, list)
}

func ChatClawGetLibFileList(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := BridgeGetLibFileListReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)
	data, httpStatus, err := BridgeGetLibFileList(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func ChatClawGetParagraphList(c *gin.Context) {
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}

	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		logs.Error("[chatclaw.getParagraphList] auth failed raw_query=%s token_len=%d", c.Request.URL.RawQuery, len(token))
		return
	}

	req := BridgeGetParagraphListReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		logs.Error("[chatclaw.getParagraphList] bind req failed raw_query=%s err=%v", c.Request.URL.RawQuery, err)
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.Status = c.DefaultQuery(`status`, `-1`)
	req.GraphStatus = c.DefaultQuery(`graph_status`, `-1`)
	req.CategoryId = c.DefaultQuery(`category_id`, `-1`)
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)

	data, httpStatus, err := BridgeGetParagraphList(adminUserId, userId, common.GetLang(c), &req)
	if httpStatus != 0 || err != nil {
		logs.Error("[chatclaw.getParagraphList] bridge failed http_status=%d err=%v req=%+v", httpStatus, err, req)
		common.FmtBridgeResponse(c, data, httpStatus, err)
		return
	}

	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func ChatClawGetLibraryGroup(c *gin.Context) {
	adminUserId, userId, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := BridgeGetLibraryGroupReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	list, httpStatus, err := BridgeGetLibraryGroup(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func ChatClawUpdateRobotSwitchStatus(c *gin.Context) {
	adminUserId, _, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawSwitchReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Id <= 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	if !tool.InArrayInt(req.SwitchStatus, []int{define.SwitchOff, define.SwitchOn}) {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}
	info, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(req.Id)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	_, err = msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(req.Id)).
		Update(msql.Datas{
			`chat_claw_switch_status`: req.SwitchStatus,
			`update_time`:             tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{`id`: req.Id, `switch_status`: req.SwitchStatus})
}

func ChatClawUpdateLibrarySwitchStatus(c *gin.Context) {
	adminUserId, _, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}
	req := ChatClawSwitchReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Id <= 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	if !tool.InArrayInt(req.SwitchStatus, []int{define.SwitchOff, define.SwitchOn}) {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}
	info, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`id`, cast.ToString(req.Id)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	_, err = msql.Model(`chat_ai_library`, define.Postgres).
		Where(`id`, cast.ToString(req.Id)).
		Update(msql.Datas{
			`chat_claw_switch_status`: req.SwitchStatus,
			`update_time`:             tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{`id`: req.Id, `switch_status`: req.SwitchStatus})
}

// ChatClawTokenLogList returns paginated token issuance logs for current user.
func ChatClawTokenLogList(c *gin.Context) {
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	size := cast.ToInt(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	m := msql.Model(define.TableChatClawTokenLog, define.Postgres).
		Where("user_id", cast.ToString(userId)).
		Order("create_time desc")
	list, total, err := m.Paginate(page, size)
	if err != nil {
		logs.Error("chatclaw token log list error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{"list": list, "total": total})
}

// ChatClawTokenForceOffline revokes a token manually. Without id, revokes the current request's token.
func ChatClawTokenForceOffline(c *gin.Context) {
	req := ChatClawForceOfflineReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if req.Reason == "" {
		req.Reason = "manual_offline"
	}

	var logId int
	var userId int

	if req.Id > 0 {
		userId = getLoginUserId(c)
		if userId <= 0 {
			common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
			return
		}
		info, err := msql.Model(define.TableChatClawTokenLog, define.Postgres).
			Where("id", cast.ToString(req.Id)).
			Where("user_id", cast.ToString(userId)).
			Field("id,status").
			Find()
		if err != nil {
			logs.Error("chatclaw token force offline find error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(info) == 0 {
			common.FmtError(c, `no_data`)
			return
		}
		if cast.ToInt(info["status"]) == common.ChatClawTokenStatusRevoked {
			common.FmtOk(c, map[string]any{"id": req.Id, "status": common.ChatClawTokenStatusRevoked})
			return
		}
		logId = cast.ToInt(info["id"])
	} else {
		claims, oldToken, err := common.GetChatClawAuthClaims(c)
		if err != nil {
			common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
			return
		}
		userId = cast.ToInt(claims["user_id"])
		if userId <= 0 {
			common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
			return
		}
		oldTokenHash := common.GetTokenSha256(oldToken)
		oldLog, err := msql.Model(define.TableChatClawTokenLog, define.Postgres).
			Where("user_id", cast.ToString(userId)).
			Where("token_hash", oldTokenHash).
			Order("id desc").
			Field("id,status").
			Find()
		if err != nil {
			oldLog, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
				Where("user_id", cast.ToString(userId)).
				Where("token", oldToken).
				Order("id desc").
				Field("id,status").
				Find()
			if err != nil {
				logs.Error("chatclaw token force offline find current token error: %s", err.Error())
				common.FmtError(c, `sys_err`)
				return
			}
		}
		if len(oldLog) == 0 {
			common.FmtError(c, `no_data`)
			return
		}
		if cast.ToInt(oldLog["status"]) == common.ChatClawTokenStatusRevoked {
			common.FmtOk(c, map[string]any{"id": oldLog["id"], "status": common.ChatClawTokenStatusRevoked})
			return
		}
		logId = cast.ToInt(oldLog["id"])
	}

	_, err := msql.Model(define.TableChatClawTokenLog, define.Postgres).
		Where("id", cast.ToString(logId)).
		Where("user_id", cast.ToString(userId)).
		Update(msql.Datas{
			"status":        common.ChatClawTokenStatusRevoked,
			"revoke_time":   tool.Time2Int(),
			"revoke_by":     userId,
			"revoke_reason": req.Reason,
		})
	if err != nil {
		logs.Error("chatclaw token force offline update error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{"id": logId, "status": common.ChatClawTokenStatusRevoked})
}

// ChatClawRefreshToken rotates token and immediately revokes the current one.
func ChatClawRefreshToken(c *gin.Context) {
	req := ChatClawRefreshTokenReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	claims, oldToken, err := common.GetChatClawAuthClaims(c)
	if err != nil {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userId := cast.ToInt(claims["user_id"])
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	oldTokenHash := common.GetTokenSha256(oldToken)
	oldLog, err := msql.Model(define.TableChatClawTokenLog, define.Postgres).
		Where("user_id", cast.ToString(userId)).
		Where("token_hash", oldTokenHash).
		Order("id desc").
		Field("id").
		Find()
	if err != nil {
		logs.Error("chatclaw refresh token find old token by hash error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(oldLog) == 0 {
		oldLog, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
			Where("user_id", cast.ToString(userId)).
			Where("token", oldToken).
			Order("id desc").
			Field("id").
			Find()
		if err != nil {
			logs.Error("chatclaw refresh token find old token by token error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}
	if len(oldLog) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	_, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).
		Where("id", cast.ToString(oldLog["id"])).
		Where("user_id", cast.ToString(userId)).
		Update(msql.Datas{
			"status":        common.ChatClawTokenStatusRevoked,
			"revoke_time":   tool.Time2Int(),
			"revoke_by":     userId,
			"revoke_reason": "refresh_renew",
		})
	if err != nil {
		logs.Error("chatclaw refresh token revoke old token error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	data, err := common.GetChatClawToken(claims["user_id"], claims["user_name"], claims["parent_id"])
	if err != nil {
		logs.Error("chatclaw refresh token create new token error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	serverIP := lib_web.GetClientIP(c)
	_, err = msql.Model(define.TableChatClawTokenLog, define.Postgres).Insert(msql.Datas{
		"user_id":     userId,
		"token":       cast.ToString(data["token"]),
		"token_hash":  common.GetTokenSha256(cast.ToString(data["token"])),
		"os_type":     req.OsType,
		"os_version":  req.OsVersion,
		"client_ip":   serverIP,
		"expired_at":  cast.ToInt64(data["exp"]),
		"status":      common.ChatClawTokenStatusActive,
		"create_time": tool.Time2Int(),
	})
	if err != nil {
		logs.Error("chatclaw refresh token save new token error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{
		"token": data["token"],
		"exp":   data["exp"],
	})
}

// ChatClawLibraryRecall executes knowledge-base recall via ChatClaw authentication.
func ChatClawLibraryRecall(c *gin.Context) {
	adminUserId, _, ok := getChatClawAuthUser(c)
	if !ok {
		return
	}

	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	libraryIds := strings.TrimSpace(c.PostForm(`id`))
	question := strings.TrimSpace(c.PostForm(`question`))
	size := cast.ToInt(c.PostForm(`size`))
	similarity := cast.ToFloat64(c.PostForm(`similarity`))
	searchType := cast.ToInt(c.DefaultPostForm(`search_type`, `1`))
	rrfWeight := strings.TrimSpace(c.PostForm(`rrf_weight`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	rerankStatus := strings.TrimSpace(c.DefaultPostForm(`rerank_status`, `1`))
	recallType := cast.ToString(c.PostForm(`recall_type`))

	if libraryIds == "" || question == "" || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText, define.SearchTypeGraph}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	if err := common.CheckRrfWeight(rrfWeight, common.GetLang(c)); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if modelConfigId > 0 || useModel != "" {
		if ok := common.CheckModelIsValid(adminUserId, modelConfigId, useModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
			return
		}
	}

	robot := msql.Params{
		`recall_type`:   recallType,
		`rrf_weight`:    rrfWeight,
		`admin_user_id`: cast.ToString(adminUserId),
	}
	for _, libraryId := range strings.Split(libraryIds, `,`) {
		libId := cast.ToInt(strings.TrimSpace(libraryId))
		if libId <= 0 {
			continue
		}
		info, err := common.GetLibraryInfo(libId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(info) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		robotName, _ := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`rerank_status`, `1`).
			Where(`rerank_model_config_id`, cast.ToString(rerankModelConfigID)).
			Value(`robot_name`)
		if rerankModelConfigID > 0 && cast.ToInt(rerankStatus) == define.SwitchOn {
			robot[`rerank_status`] = rerankStatus
			robot[`rerank_model_config_id`] = cast.ToString(rerankModelConfigID)
			robot[`rerank_use_model`] = cast.ToString(rerankUseModel)
			robot[`robot_name`] = robotName
		}
		if searchType == define.SearchTypeGraph {
			if !cast.ToBool(info[`graph_switch`]) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `graph is not enabled`))))
				return
			}
			robot[`admin_user_id`] = info[`admin_user_id`]
			robot[`model_config_id`] = info[`graph_model_config_id`]
			robot[`use_model`] = info[`graph_use_model`]
			robot[`id`] = strconv.Itoa(0)
		}
		if modelConfigId > 0 && useModel != "" {
			robot[`model_config_id`] = cast.ToString(modelConfigId)
			robot[`use_model`] = useModel
		}
	}

	list, _, err := common.GetMatchLibraryParagraphList(
		common.GetLang(c),
		cast.ToString(adminUserId),
		lib_define.ChatClawClient,
		"",
		question,
		[]string{},
		libraryIds,
		size,
		similarity,
		searchType,
		robot,
	)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	for _, item := range list {
		library, err := common.GetLibraryInfo(cast.ToInt(item[`library_id`]), adminUserId)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		item[`library_name`] = library[`library_name`]
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}
