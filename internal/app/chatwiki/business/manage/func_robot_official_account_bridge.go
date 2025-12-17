// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"context"
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"strings"
	"time"
)

type BridgeGetOfficialDraftListReq struct {
	Page    string `form:"page"`
	Size    string `form:"size"`
	Title   string `form:"title"`
	AppId   string `form:"app_id"`
	GroupId int    `form:"group_id"`
}

func BridgeGetOfficialDraftList(adminUserId int, lang string, req *BridgeGetOfficialDraftListReq) (map[string]any, int, error) {
	if req.AppId == "" {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	var (
		err error
	)
	page := max(1, cast.ToInt(req.Page))
	size := max(1, cast.ToInt(req.Size))

	query := msql.Model(`wechat_official_account_draft`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`app_id`, req.AppId).
		Field(`*`)

	if req.GroupId >= 0 {
		query.Where(`group_id`, cast.ToString(req.GroupId))
	}

	if req.Title != "" {
		query.Where(`title`, `like`, req.Title)
	}

	orderRaw := `draft_create_time desc`

	query.Order(orderRaw)

	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	var formatedList = make([]map[string]any, 0)
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		formatedList = append(formatedList, tempItem)
	}

	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	data := map[string]any{`list`: formatedList, `total`: total, `page`: page, `size`: size}
	return data, 0, nil
}

type DraftGroupStruct struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
}

func BridgeGetOfficialDraftGroupList(adminUserId int) ([]DraftGroupStruct, int, error) {

	data, _ := msql.Model(`wechat_official_account_draft_group`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`create_time ASC`).Field(`id,group_name`).Select()

	respData := []DraftGroupStruct{
		{Id: -1, GroupName: "全部"},
		{Id: 0, GroupName: "未分组"},
	}

	for _, datum := range data {
		respData = append(respData, DraftGroupStruct{
			Id:        cast.ToInt(datum[`id`]),
			GroupName: cast.ToString(datum[`group_name`]),
		})
	}

	return respData, 0, nil
}

type BridgeCreateSendTaskReq struct {
	TaskId        int    `form:"task_id"`
	DraftId       int    `form:"draft_id"`
	AppId         string `form:"app_id"`
	AccessKey     string `form:"access_key"`
	SendTime      int64  `form:"send_time"`
	CommentStatus int    `form:"comment_status"`
	CommentRuleId int    `form:"comment_rule_id"`
	OpenStatus    int    `form:"open_status"`
	ToUserType    int    `form:"to_user_type"`
	ToUser        string `form:"to_user"`
	TaskName      string `form:"task_name"`
}

func BridgeCreateSendTask(adminUserId int, lang string, req *BridgeCreateSendTaskReq) (map[string]any, int, error) {

	model := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres)
	insertData := msql.Datas{
		`admin_user_id`:   adminUserId,
		`update_time`:     time.Now().Unix(),
		`app_id`:          req.AppId,
		`access_key`:      req.AccessKey,
		`draft_id`:        req.DraftId,
		`comment_status`:  req.CommentStatus,
		`open_status`:     req.OpenStatus,
		`send_time`:       req.SendTime,
		`to_user_type`:    req.ToUserType,
		`task_name`:       req.TaskName,
		`comment_rule_id`: req.CommentRuleId,
	}

	id := req.TaskId

	if req.TaskId != 0 {
		//查询任务信息
		taskData, _ := model.Where(`id`, cast.ToString(req.TaskId)).Find()

		if cast.ToInt(taskData["id"]) == 0 {
			return nil, -1, errors.New(i18n.Show(lang, `no_data`))
		}

		//不是待发送的状态，不允许编辑
		if cast.ToInt(taskData["send_status"]) != define.BatchSendStatusWait {
			return nil, -1, errors.New(i18n.Show(lang, `task_not_wait`))
		}

		model.Where(`id`, cast.ToString(req.TaskId)).Update(insertData)
	} else {
		//查询默认规则
		if req.CommentRuleId == 0 { //没有指定生效规则，查询默认规则
			ruleInfo, _ := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserId)).Where("is_default", "1").Find()
			insertData["comment_rule_id"] = cast.ToInt(ruleInfo["id"])
		}

		insertData["create_time"] = time.Now().Unix()
		insert_id, err := model.Insert(insertData, `id`)

		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		id = int(insert_id)
	}

	//任务开启状态 创建群发定时任务
	if req.OpenStatus == 1 {
		BridgeAddDelayTask(adminUserId, int(id), req.SendTime)
	}

	return map[string]any{`id`: id}, 0, nil
}

func BridgeAddDelayTask(adminUserId, tasId int, sendTime int64) {
	//默认延迟10秒执行，如果发送时间，超过当前时间，计算一下延迟执行时间
	delayTime := int64(5)
	if sendTime > time.Now().Unix() {
		delayTime = sendTime - time.Now().Unix()
	}

	common.AddDelayTask(define.DelayTaskEvent{
		BaseDelayTask: define.BaseDelayTask{Type: define.OfficialAccountBatchSendDelayTask},
		AdminUserId:   adminUserId,
		TaskId:        tasId,
	}, delayTime)
}

type BridgeBatchSendTaskListReq struct {
	Page  string `form:"page"`
	Size  string `form:"size"`
	AppId string `form:"app_id"`
}

func BridgeSendTaskList(adminUserId int, lang string, req *BridgeBatchSendTaskListReq) (map[string]any, int, error) {

	page := max(1, cast.ToInt(req.Page))
	size := max(10, cast.ToInt(req.Size))

	query := msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Alias("a").
		Join("wechat_official_account_draft b", "a.draft_id = b.id ", "inner").
		Join("wechat_official_account_comment_rule c", "a.comment_rule_id = c.id", "left").
		Join("wechat_official_account_draft_group d", "b.group_id = d.id", "left").
		Where(`a.admin_user_id`, cast.ToString(adminUserId)).
		Where(`a.send_status`, `>=`, "0").
		Field(`a.*,b.thumb_url,b.title,b.digest,c.rule_name comment_rule_name,d.group_name,c.is_default`)

	if req.AppId != "" {
		query.Where(`a.app_id`, req.AppId)
	}

	orderRaw := `a.is_top DESC,a.create_time desc`

	query.Order(orderRaw)

	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	var formatedList = make([]map[string]any, 0)
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		if item["comment_rule_id"] == "0" && item["comment_rule_name"] == "" {
			tempItem["comment_rule_name"] = "暂无规则"
		}

		formatedList = append(formatedList, tempItem)
	}

	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	data := map[string]any{`list`: formatedList, `total`: total, `page`: page, `size`: size}
	return data, 0, nil
}

func BridgeChangeCommentStatus(access_key, msg_id string, index, common_status int) (*response.ResponseOfficialAccount, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}

	if common_status == 1 {
		res, err := client.Comment.Open(context.Background(), msg_id, index)
		if res.ErrCode != 0 {
			return res, errors.New(res.ErrMsg)
		}

		return res, err
	} else {
		res, err := client.Comment.Close(context.Background(), msg_id, index)
		if res.ErrCode != 0 {
			return res, errors.New(res.ErrMsg)
		}

		return res, err
	}

}

func BridgeMarkElect(access_key, msg_id string, index, comment_id int) (*response.ResponseOfficialAccount, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}

	res, err := client.Comment.MarkElect(context.Background(), msg_id, index, comment_id)
	if res.ErrCode != 0 {
		return res, errors.New(res.ErrMsg)
	}

	return res, err
}

func BridgeUnMarkElect(access_key, msg_id string, index, comment_id int) (*response.ResponseOfficialAccount, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}
	res, err := client.Comment.UnmarkElect(context.Background(), msg_id, index, comment_id)
	if res.ErrCode != 0 {
		return res, errors.New(res.ErrMsg)
	}

	return res, err
}

func BridgeReplyComment(access_key, msg_id, repl_content string, index, comment_id int) (*response.ResponseOfficialAccount, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}

	res, err := client.Comment.Reply(context.Background(), msg_id, index, comment_id, repl_content)
	if res.ErrCode != 0 {
		return res, errors.New(res.ErrMsg)
	}

	return res, err
}

func BridgeDeleteComment(access_key, msg_id string, index, comment_id int) (*response.ResponseOfficialAccount, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}

	res, err := client.Comment.Delete(context.Background(), msg_id, index, comment_id)

	if res.ErrCode != 0 {
		return res, errors.New(res.ErrMsg)
	}

	_, err = msql.Model(`wechat_official_comment_list`, define.Postgres).
		Where("msg_data_id", msg_id).
		Where("access_key", access_key).
		Where("user_comment_id", cast.ToString(comment_id)).Update(msql.Datas{
		`update_time`:   time.Now().Unix(),
		`delete_status`: 1,
	})

	if err != nil {
		logs.Error("修改数据执行错误：" + err.Error())
	}

	return res, err
}

func BridgeDeleteCommentReply(access_key, msg_id string, index, comment_id int) (map[string]any, error) {

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		return nil, err
	}

	params := &object.HashMap{
		"msg_data_id":     msg_id,
		"user_comment_id": comment_id,
		"index":           0,
	}

	respData := map[string]any{}

	_, err = client.Base.BaseClient.HttpPostJson(context.Background(), "/cgi-bin/comment/reply/delete", params, nil, nil, &respData)
	if err != nil {
		logs.Error("错误内容：" + err.Error())
	}

	if cast.ToInt(respData["errcode"]) != 0 {
		return nil, errors.New(cast.ToString(respData["errmsg"]))
	}

	_, err = msql.Model(`wechat_official_comment_list`, define.Postgres).
		Where("msg_data_id", msg_id).
		Where("access_key", access_key).
		Where("user_comment_id", cast.ToString(comment_id)).Update(msql.Datas{
		`update_time`:        time.Now().Unix(),
		`reply_comment_text`: "",
		`reply_create_time`:  0,
	})

	if err != nil {
		logs.Error("修改数据执行错误：" + err.Error())
	}

	return respData, err
}

type BatchSendRes struct {
	MsgId     int64  `json:"msg_id"`
	MsgStatus string `json:"msg_status"`
}

func BridgeGetSendStatus(access_key, msg_id string) (*BatchSendRes, error) {
	respData := BatchSendRes{}

	client, _, err := getOfficialAppClient(access_key)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	res, err := client.Broadcasting.Status(context.Background(), msg_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	err = tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(res), &respData)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return &respData, err
}

func getOfficialAppClient(access_key string) (*officialAccount.OfficialAccount, msql.Params, error) {
	appInfo, err := common.GetWechatAppInfo("access_key", access_key)
	if err != nil {
		return nil, nil, err
	}

	client, err := common.GetOfficialAccountApp(appInfo["app_id"], appInfo["app_secret"])
	if err != nil {
		return nil, appInfo, err
	}

	return client, appInfo, nil
}

type BridgeDeleteCommentRule struct {
	Type      []int    `json:"type"`      //1:触发敏感词，2：AI检测
	Keywords  []string `json:"keywords"`  //敏感词列表
	Condition int      `json:"condition"` //条件，1：全部满足，2：满足其中一个
	Priority  int      `json:"priority"`  //优先级，1：触发敏感词，2：AI检测
	Prompt    string   `json:"prompt"`    //AI 检测提示词
}
type BridgeReplyCommentRule struct {
	CheckReplyPrompt string `json:"check_reply_prompt"` //判断是否回复提示词
	ReplyType        int    `json:"reply_type"`         //回复类型。1：指定内容，2：AI生成回复
	ReplyPrompt      string `json:"reply_prompt"`       //回复内容或者回复提示词
}
type BridgeElectCommentRule struct {
	Prompt string `json:"prompt"` //提示词
}

type BridgeSaveCommentRuleReq struct {
	Id                  int                     `from:"id"`
	RuleName            string                  `from:"rule_name"`
	UseModel            string                  `from:"use_model"`
	DeleteCommentSwitch int                     `json:"delete_comment_switch"`
	DeleteCommentRule   BridgeDeleteCommentRule `json:"delete_comment_rule"`
	ReplyCommentSwitch  int                     `json:"reply_comment_switch"`
	ReplyCommentRule    BridgeReplyCommentRule  `json:"reply_comment_rule"`
	ElectCommentSwitch  int                     `json:"elect_comment_switch"`
	ElectCommentRule    BridgeElectCommentRule  `json:"elect_comment_rule"`
	IsDefault           int                     `json:"is_default"`
	ModelConfigId       int                     `json:"model_config_id"`
}

func BridgeSaveCommentRule(adminUserId int, lang string, req *BridgeSaveCommentRuleReq) (map[string]any, int, error) {
	insertData := msql.Datas{
		`admin_user_id`:         adminUserId,
		`update_time`:           time.Now().Unix(),
		`rule_name`:             req.RuleName,
		`use_model`:             req.UseModel,
		`delete_comment_switch`: req.DeleteCommentSwitch,
		`delete_comment_rule`:   tool.JsonEncodeNoError(req.DeleteCommentRule),
		`reply_comment_switch`:  req.ReplyCommentSwitch,
		`reply_comment_rule`:    tool.JsonEncodeNoError(req.ReplyCommentRule),
		`elect_comment_switch`:  req.ElectCommentSwitch,
		`elect_comment_rule`:    tool.JsonEncodeNoError(req.ElectCommentRule),
		`is_default`:            req.IsDefault,
		`model_config_id`:       req.ModelConfigId,
	}
	m := msql.Model(`wechat_official_account_comment_rule`, define.Postgres)

	//如果创建默认规则，验证一下是否存在重复的默认规则
	if req.IsDefault == 1 {
		m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`is_default`, cast.ToString(req.IsDefault))

		if req.Id != 0 {
			m.Where(`id`, "!=", cast.ToString(req.Id))
		}
		ruleList, _ := m.Select()
		//新建场景，当前账户下有默认规则，退出
		if len(ruleList) != 0 {
			return nil, -1, errors.New(i18n.Show(lang, `default_comment_rule_already_exists`))
		}

		//默认规则，默认开启
		insertData["switch"] = 1
	}

	//编辑场景
	if req.Id != 0 {
		_, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, cast.ToString(req.Id)).Update(insertData)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}

	} else {
		insertData[`create_time`] = time.Now().Unix()
		id, err := m.Insert(insertData, `id`)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		req.Id = int(id)
	}

	return map[string]any{`id`: req.Id}, 0, nil
}

type BridgeGetCommentRuleListReq struct {
	Page      string `form:"page"`
	Size      string `form:"size"`
	RuleName  string `form:"rule_name"`
	IsDefault int    `form:"is_default"`
}

func BridgeGetCommentRuleList(adminUserId int, lang string, req *BridgeGetCommentRuleListReq) (map[string]any, int, error) {

	var (
		err error
	)

	page := max(1, cast.ToInt(req.Page))
	size := max(1, cast.ToInt(req.Size))

	query := msql.Model(`wechat_official_account_comment_rule`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`*`)

	if req.IsDefault != -1 {
		query.Where(`is_default`, cast.ToString(req.IsDefault))
	}

	if req.RuleName != "" {
		query.Where(`rule_name`, `like`, req.RuleName)
	}

	orderRaw := `is_default DESC,create_time desc`

	query.Order(orderRaw)

	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	var formatedList = make([]map[string]any, 0)
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		//查询规则关联群发任务
		tempItem["task_info"], _ = msql.Model(`wechat_official_account_batch_send_task`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`comment_rule_id`, cast.ToString(tempItem["id"])).Field("id,task_name").Select()
		formatedList = append(formatedList, tempItem)
	}

	data := map[string]any{`list`: formatedList, `total`: total, `page`: page, `size`: size}
	return data, 0, nil
}

type BridgeGetCommentListReq struct {
	Page        string `form:"page"`
	Size        string `form:"size"`
	TaskId      string `form:"task_id"`
	CommentText string `form:"comment_text"`
	CheckResult string `form:"check_result"`
}

func BridgeCommentList(adminUserId int, lang string, req *BridgeGetCommentListReq) (map[string]any, int, error) {

	var (
		err error
	)
	page := max(1, cast.ToInt(req.Page))
	size := max(1, cast.ToInt(req.Size))

	query := msql.Model(`wechat_official_comment_list`, define.Postgres).Alias(`a`).
		Join("wechat_official_account_draft b", "a.draft_id = b.id", "left").
		Join("wechat_official_account_comment_rule c", "a.comment_rule_id = c.id", "left").
		Where(`a.admin_user_id`, cast.ToString(adminUserId)).
		Field(`a.*,b.title as  draft_title,c.rule_name`)

	if req.CommentText != "" {
		query.Where(`a.content_text`, `like`, req.CommentText)
	}

	if req.TaskId != "" {
		query.Where(`a.task_id`, req.TaskId)
	}

	if req.CheckResult != "0" {
		query.Where("( a.ai_comment_result @> '[" + req.CheckResult + "]'::jsonb )")
	}

	orderRaw := `a.comment_create_time DESC,a.user_comment_id desc`

	query.Order(orderRaw)

	list, total, err := query.Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	var formatedList = make([]map[string]any, 0)
	for _, item := range list {
		tempItem := make(map[string]any)
		for k, v := range item {
			tempItem[k] = v
		}

		comment_result := []int{}
		commentResultText := []string{}
		_ = tool.JsonDecodeUseNumber(cast.ToString(tempItem["ai_comment_result"]), &comment_result)
		for _, result := range comment_result {
			commentResultText = append(commentResultText, define.CommentExecTypeMap[result])
		}
		tempItem["ai_comment_result_text"] = strings.Join(commentResultText, ",")
		tempItem["ai_comment_rule_text"] = strings.Split(cast.ToString(tempItem["ai_comment_rule_text"]), ",")

		formatedList = append(formatedList, tempItem)
	}

	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}

	data := map[string]any{`list`: formatedList, `total`: total, `page`: page, `size`: size, `type`: define.CommentExecTypeMap}
	return data, 0, nil
}
