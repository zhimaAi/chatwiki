// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"net/http"
	"time"
)

func SyncDraftList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	access_key := cast.ToString(c.Query(`access_key`))
	limit := cast.ToInt(c.Query(`limit`))

	message, err := tool.JsonEncode(map[string]any{
		"admin_user_id": adminUserId,
		"access_key":    access_key,
		"limit":         limit,
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	err = common.AddJobs(define.OfficialAccountDraftSyncTopic, message)

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))

}

func DraftList(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	req := BridgeGetOfficialDraftListReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	req.Title = cast.ToString(c.Query(`title`))
	req.Page = cast.ToString(c.DefaultQuery(`page`, `1`))
	req.Size = cast.ToString(c.DefaultQuery(`size`, `10`))
	req.AppId = cast.ToString(c.DefaultQuery(`app_id`, ``))
	req.GroupId = cast.ToInt(c.DefaultQuery(`group_id`, `-1`))

	data, httpStatus, err := BridgeGetOfficialDraftList(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)

}

func DraftGroupList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	data, httpStatus, err := BridgeGetOfficialDraftGroupList(adminUserId)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func SaveDraftGroup(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	id := cast.ToInt(c.PostForm(`id`))
	group_name := cast.ToString(c.PostForm(`group_name`))

	if group_name == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_name`))))
		return
	}

	model := msql.Model(`wechat_official_account_draft_group`, define.Postgres)
	data, _ := model.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`group_name`, group_name).Find()

	if id != 0 && cast.ToInt(data[`id`]) != 0 && id != cast.ToInt(data[`id`]) { //编辑场景
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `duplicated_draft_group_name`))))
		return
	}

	if id == 0 {
		model.Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`create_time`:   time.Now().Unix(),
			`update_time`:   time.Now().Unix(),
			`group_name`:    group_name,
		})
	} else {
		model.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
			`group_name`:  group_name,
			`update_time`: time.Now().Unix(),
		})
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DeleteDraftGroup(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	id := cast.ToInt(c.DefaultPostForm(`id`, `0`))
	if id == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `id`))))
		return
	}

	_, err := msql.Model(`wechat_official_account_draft_group`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, cast.ToString(id)).Delete()

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	_, err = msql.Model(`wechat_official_account_draft`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`group_id`, cast.ToString(id)).Update(msql.Datas{
		`update_time`: time.Now().Unix(),
		`group_id`:    0,
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))

}

func MoveDraftGroup(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	draft_ids := cast.ToString(c.PostForm(`draft_id`))
	if len(draft_ids) <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	app_id := cast.ToString(c.DefaultPostForm(`app_id`, `0`))
	group_id := cast.ToInt(c.DefaultPostForm(`group_id`, `0`))
	if group_id == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `group_id`))))
		return
	}
	if app_id == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `app_id`))))
		return
	}

	_, err := msql.Model(`wechat_official_account_draft`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`app_id`, app_id).
		Where(`id`, `in`, draft_ids).
		Update(msql.Datas{
			`update_time`: time.Now().Unix(),
			`group_id`:    group_id,
		})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func CreateBathSendTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	req := BridgeCreateSendTaskReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	data, httpStatus, err := BridgeCreateSendTask(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)

}

func BathSendTaskList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	req := BridgeBatchSendTaskListReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	data, httpStatus, err := BridgeSendTaskList(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)

}

func SetBatchSendTaskTopStatus(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToInt(c.PostForm(`task_id`))
	is_top := cast.ToInt(c.PostForm(`is_top`))

	res, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
		`is_top`:      is_top,
		`update_time`: time.Now().Unix(),
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func SetBatchSendTaskOpenStatus(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToInt(c.PostForm(`task_id`))
	open_status := cast.ToInt(c.PostForm(`open_status`))

	taskInfo, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	res, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
		`open_status`: open_status,
		`update_time`: time.Now().Unix(),
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//需要开启状态
	if open_status == cast.ToInt(define.BaseOpen) {
		BridgeAddDelayTask(adminUserId, task_id, cast.ToInt64(taskInfo["send_time"]))
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func DeleteBatchSendTask(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToInt(c.PostForm(`task_id`))
	taskInfo, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//发送中或者发送完成的，不支持删除
	if cast.ToInt(taskInfo["send_status"]) == 1 || cast.ToInt(taskInfo["send_status"]) == 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `task_running`))))
		return
	}

	updateRes, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
		`send_status`: -1,
		`update_time`: time.Now().Unix(),
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(updateRes, nil))
}

func ChangeBatchTaskCommentRuleStatus(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToInt(c.PostForm(`task_id`))
	ai_comment_status := cast.ToInt(c.PostForm(`ai_comment_status`))

	m := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres)

	taskInfo, err := m.Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	update := msql.Datas{
		`update_time`:       time.Now().Unix(),
		`ai_comment_status`: ai_comment_status,
	}

	if cast.ToInt(taskInfo["comment_rule_id"]) == 0 && ai_comment_status == cast.ToInt(define.BaseOpen) { //如果开启规则，且当前群发任务没有关联规则ID，查一下默认规则
		ruleInfo, _ := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`is_default`, "1").Find()
		if cast.ToInt(ruleInfo["id"]) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `default_comment_rule_not_found`))))
			return
		}
		update[`comment_rule_id`] = ruleInfo["id"]
	}

	res, err := m.Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(update)

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func ChangeBatchTaskCommentRule(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToInt(c.PostForm(`task_id`))
	rule_id := cast.ToInt(c.PostForm(`rule_id`))

	res, err := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`id`, cast.ToString(task_id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
		`update_time`:     time.Now().Unix(),
		`comment_rule_id`: rule_id,
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))
}

func ChangeCommentStatus(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	task_id := cast.ToString(c.PostForm(`task_id`))
	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_status := cast.ToInt(c.PostForm(`comment_status`))

	model := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres)
	//先查询一下任务状态，
	taskInfo, err := model.Where(`id`, task_id).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`access_key`, access_key).Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	if cast.ToInt(taskInfo["send_status"]) != 1 && cast.ToInt(taskInfo["send_status"]) != 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `task_not_sending`))))
		return
	}

	res, err := BridgeChangeCommentStatus(access_key, msg_id, index, comment_status)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	model.Where(`id`, task_id).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`access_key`, access_key).Update(msql.Datas{
		`update_time`:    time.Now().Unix(),
		`comment_status`: comment_status,
	})

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))
}

func MarkElect(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_id := cast.ToInt(c.PostForm(`comment_id`))

	res, err := BridgeMarkElect(access_key, msg_id, index, comment_id)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func UnMarkElect(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_id := cast.ToInt(c.PostForm(`comment_id`))

	res, err := BridgeUnMarkElect(access_key, msg_id, index, comment_id)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func ReplyComment(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_id := cast.ToInt(c.PostForm(`comment_id`))
	repl_content := cast.ToString(c.PostForm(`repl_content`))

	res, err := BridgeReplyComment(access_key, msg_id, repl_content, index, comment_id)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func DeleteComment(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_id := cast.ToInt(c.PostForm(`comment_id`))

	res, err := BridgeDeleteComment(access_key, msg_id, index, comment_id)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func DeleteCommentReply(c *gin.Context) {

	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	access_key := cast.ToString(c.PostForm(`access_key`))
	msg_id := cast.ToString(c.PostForm(`msg_id`))
	index := cast.ToInt(c.PostForm(`index`))
	comment_id := cast.ToInt(c.PostForm(`comment_id`))

	res, err := BridgeDeleteCommentReply(access_key, msg_id, index, comment_id)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))

}

func GetCommentList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	req := BridgeGetCommentListReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	req.Page = cast.ToString(c.DefaultQuery(`page`, `1`))
	req.Size = cast.ToString(c.DefaultQuery(`size`, `10`))
	req.TaskId = cast.ToString(c.DefaultQuery(`task_id`, ``))
	req.CommentText = cast.ToString(c.DefaultQuery(`comment_text`, ``))
	req.CheckResult = cast.ToString(c.DefaultQuery(`check_result`, ``))

	data, httpStatus, err := BridgeCommentList(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)

}

func SaveCommentRule(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	DeleteCommentRule := BridgeDeleteCommentRule{}
	ReplyCommentRule := BridgeReplyCommentRule{}
	ElectCommentRule := BridgeElectCommentRule{}

	_ = tool.JsonDecodeUseNumber(cast.ToString(c.PostForm(`delete_comment_rule`)), &DeleteCommentRule)
	_ = tool.JsonDecodeUseNumber(cast.ToString(c.PostForm(`reply_comment_rule`)), &ReplyCommentRule)
	_ = tool.JsonDecodeUseNumber(cast.ToString(c.PostForm(`elect_comment_rule`)), &ElectCommentRule)

	req := BridgeSaveCommentRuleReq{}
	req.DeleteCommentRule = DeleteCommentRule
	req.ReplyCommentRule = ReplyCommentRule
	req.ElectCommentRule = ElectCommentRule
	req.RuleName = cast.ToString(c.PostForm(`rule_name`))
	req.UseModel = cast.ToString(c.PostForm(`use_model`))
	req.DeleteCommentSwitch = cast.ToInt(c.PostForm(`delete_comment_switch`))
	req.ReplyCommentSwitch = cast.ToInt(c.PostForm(`reply_comment_switch`))
	req.ElectCommentSwitch = cast.ToInt(c.PostForm(`elect_comment_switch`))
	req.IsDefault = cast.ToInt(c.PostForm(`is_default`))
	req.ModelConfigId = cast.ToInt(c.PostForm(`model_config_id`))
	req.Id = cast.ToInt(c.PostForm(`id`))

	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	data, httpStatus, err := BridgeSaveCommentRule(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func DeleteCommentRule(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	rule_id := cast.ToString(c.PostForm(`id`))

	updateCommentRuleId := msql.Datas{
		"update_time":     time.Now().Unix(),
		"comment_rule_id": 0,
	}
	//查询默认规则信息，
	ruleInfo, err := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("is_default", "1").Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//如果有默认规则，改为默认规则生效
	if cast.ToInt(ruleInfo["id"]) != 0 {
		updateCommentRuleId["comment_rule_id"] = ruleInfo["id"]
	}

	query := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId)).Where("comment_rule_id", rule_id)

	_, err = query.Update(updateCommentRuleId)

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	_, err = msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("id", cast.ToString(rule_id)).Delete()

	c.String(http.StatusOK, lib_web.FmtJson(nil, err))
}

func ChangeCommentRuleStatus(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	rule_id := cast.ToString(c.PostForm(`id`))
	change_fields := cast.ToString(c.PostForm(`change_fields`))
	switch_status := cast.ToString(c.PostForm(`switch_status`))

	ruleInfo, err := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("id", cast.ToString(rule_id)).Find()

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if ruleInfo["is_default"] == "1" && change_fields == "switch" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `default_comment_rule_not_edit`))))
	}

	res, err := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("id", cast.ToString(rule_id)).Update(msql.Datas{
		"update_time": time.Now().Unix(),
		change_fields: switch_status,
	})

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, err))

}

func GetCommentRuleList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	req := BridgeGetCommentRuleListReq{}
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	req.Page = cast.ToString(c.DefaultQuery(`page`, `1`))
	req.Size = cast.ToString(c.DefaultQuery(`size`, `10`))
	req.RuleName = cast.ToString(c.DefaultQuery(`rule_name`, ``))
	req.IsDefault = cast.ToInt(c.DefaultQuery(`is_default`, `0`))

	data, httpStatus, err := BridgeGetCommentRuleList(adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func GetCommentRuleInfo(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	rule_id := cast.ToString(c.Query(`id`))

	res, err := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("id", cast.ToString(rule_id)).Find()

	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	tempItem := make(map[string]any)
	for key, item := range res {
		tempItem[key] = item
	}

	tempItem["task_info"], _ = msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`comment_rule_id`, cast.ToString(tempItem["id"])).Field("id,task_name").Select()

	c.String(http.StatusOK, lib_web.FmtJson(tempItem, err))

}
