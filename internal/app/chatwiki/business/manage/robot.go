// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type metaSearchCondition struct {
	Key   string `json:"key"`
	Type  int    `json:"type"`  //0 string, 1 time, 2 number (same as define.LibraryMetaType*)
	Op    int    `json:"op"`    // define.MetaOp*
	Value string `json:"value"` //can be empty (empty/not-empty operators)
}

func validateMetaSearchConfig(lang string, metaSwitch, metaType int, raw string) (int, int, string, error) {
	//Off by default
	if metaSwitch != define.MetaSearchSwitchOn {
		return define.MetaSearchSwitchOff, define.MetaSearchTypeAnd, `[]`, nil
	}
	// type
	if metaType == 0 {
		metaType = define.MetaSearchTypeAnd
	}
	if metaType != define.MetaSearchTypeAnd && metaType != define.MetaSearchTypeOr {
		return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_type`))
	}
	// list
	raw = strings.TrimSpace(raw)
	if raw == `` {
		raw = `[]`
	}
	conds := make([]metaSearchCondition, 0)
	if err := tool.JsonDecode(raw, &conds); err != nil {
		logs.Error(err.Error())
		return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
	}
	if len(conds) > define.MetaSearchMaxConditions {
		return 0, 0, ``, errors.New(i18n.Show(lang, `max_conditions_exceed`, define.MetaSearchMaxConditions))
	}

	for _, c := range conds {
		c.Key = strings.TrimSpace(c.Key)
		if c.Key == `` {
			return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
		}
		//Fixed key format: built-in key or key_number
		if !define.IsBuiltinMetaKey(c.Key) && !common.IsCustomMetaKey(c.Key) {
			return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
		}
		if !define.IsLibraryMetaTypeValid(c.Type) {
			return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
		}

		//Op rules
		switch c.Type {
		case define.LibraryMetaTypeString:
			if !tool.InArrayInt(c.Op, []int{define.MetaOpIs, define.MetaOpIsNot, define.MetaOpContains, define.MetaOpNotContains, define.MetaOpEmpty, define.MetaOpNotEmpty}) {
				return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
			}
		case define.LibraryMetaTypeNumber, define.LibraryMetaTypeTime:
			if !tool.InArrayInt(c.Op, []int{define.MetaOpIs, define.MetaOpIsNot, define.MetaOpEmpty, define.MetaOpNotEmpty, define.MetaOpGt, define.MetaOpEq, define.MetaOpLt, define.MetaOpGte, define.MetaOpLte}) {
				return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
			}
		}

		//Value rules
		v := strings.TrimSpace(c.Value)
		needValue := !(c.Op == define.MetaOpEmpty || c.Op == define.MetaOpNotEmpty)
		if needValue {
			if v == `` {
				return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
			}
			if utf8.RuneCountInString(v) > 20 {
				return 0, 0, ``, errors.New(i18n.Show(lang, `meta_condition_value_too_long`, 20))
			}
			//For number/time, the value must be numeric
			if c.Type == define.LibraryMetaTypeNumber {
				if ok, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, v); !ok {
					return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
				}
			}
			if c.Type == define.LibraryMetaTypeTime {
				if ok, _ := regexp.MatchString(`^\d{1,20}$`, v); !ok {
					return 0, 0, ``, errors.New(i18n.Show(lang, `param_invalid`, `meta_search_condition_list`))
				}
			}
		}
	}

	out := tool.JsonEncodeNoError(conds)
	return define.MetaSearchSwitchOn, metaType, out, nil
}

func GetRobotList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	m := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_name,robot_intro,robot_avatar,robot_key,application_type,creator,start_node_key,group_id,sort_num,is_top`).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`is_top desc, sort_num desc ,id desc`)

	applicationType := cast.ToInt(c.DefaultQuery(`application_type`, `-1`))
	if applicationType >= 0 { //filter by application type
		m.Where(`application_type`, cast.ToString(applicationType))
	}

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	//check permission
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot}) {
		// managedRobotIdList := GetUserManagedData(userId, `managed_robot_list`)
		managedRobotIdList := []string{`0`}
		permissionData, _ := common.GetAllPermissionManage(adminUserId, cast.ToString(userId), define.IdentityTypeUser, define.ObjectTypeRobot)
		for _, permission := range permissionData {
			managedRobotIdList = append(managedRobotIdList, cast.ToString(permission[`object_id`]))
		}
		//m.Where(`id`, `in`, strings.Join(managedRobotIdList, `,`))
	}

	list, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	//Add has_published field
	//Collect robot IDs with application_type=1
	var workflowRobotIds []string
	for i := range list {
		applicationType := cast.ToInt(list[i][`application_type`])
		if applicationType == 1 {
			workflowRobotIds = append(workflowRobotIds, cast.ToString(list[i][`id`]))
		}
	}

	//Batch query robots that have published versions
	publishedRobotMap := make(map[string]bool)
	if len(workflowRobotIds) > 0 {
		publishedRobots, err := msql.Model(`work_flow_version`, define.Postgres).
			Where(`robot_id`, `in`, strings.Join(workflowRobotIds, `,`)).
			Group(`robot_id`).
			ColumnArr(`robot_id`)
		if err != nil {
			logs.Error(err.Error())
		} else {
			for _, robotId := range publishedRobots {
				publishedRobotMap[robotId] = true
			}
		}
	}

	//Set has_published field
	for i := range list {
		applicationType := cast.ToInt(list[i][`application_type`])
		if applicationType == 0 {
			//application_type=0 (chat): has_published=1
			list[i][`has_published`] = `1`
		} else if applicationType == 1 {
			//application_type=1 (workflow): check for published versions
			robotId := cast.ToString(list[i][`id`])
			if publishedRobotMap[robotId] {
				list[i][`has_published`] = `1`
			} else {
				list[i][`has_published`] = `0`
			}
		} else {
			//Other types default to 0
			list[i][`has_published`] = `0`
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveRobot(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	robotName := strings.TrimSpace(c.PostForm(`robot_name`))
	checkName := cast.ToInt(c.PostForm(`check_name`))
	robotIntro := strings.TrimSpace(c.PostForm(`robot_intro`))
	robotAvatar := strings.TrimSpace(c.DefaultPostForm(`avatar_from_template`, ``))
	promptType := cast.ToInt(c.DefaultPostForm(`prompt_type`, `1`))
	prompt := strings.TrimSpace(c.PostForm(`prompt`))
	promptStruct := strings.TrimSpace(c.DefaultPostForm(`prompt_struct`, common.GetDefaultPromptStruct(common.GetLang(c))))
	libraryIds := strings.TrimSpace(c.PostForm(`library_ids`))
	formIds := strings.TrimSpace(c.PostForm(`form_ids`))
	welcomes := strings.TrimSpace(c.PostForm(`welcomes`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	temperature := cast.ToFloat32(c.DefaultPostForm(`temperature`, `0.5`))
	maxToken := cast.ToInt(c.DefaultPostForm(`max_token`, `2000`))
	contextPair := cast.ToInt(c.DefaultPostForm(`context_pair`, `6`))
	enableThinking := cast.ToUint(c.PostForm(`enable_thinking`))
	questionMultipleSwitch := cast.ToUint(c.PostForm(`question_multiple_switch`))
	topK := cast.ToInt(c.DefaultPostForm(`top_k`, `5`))
	similarity := cast.ToFloat32(c.DefaultPostForm(`similarity`, `0.6`))
	recallNeighborSwitch := cast.ToBool(c.DefaultPostForm(`recall_neighbor_switch`, `false`))
	recallNeighborBeforeNum := cast.ToInt(c.DefaultPostForm(`recall_neighbor_before_num`, `1`))
	recallNeighborAfterNum := cast.ToInt(c.DefaultPostForm(`recall_neighbor_after_num`, `1`))

	searchType := cast.ToInt(c.DefaultPostForm(`search_type`, `1`))
	rrfWeight := strings.TrimSpace(c.PostForm(`rrf_weight`))
	chatType := cast.ToInt(c.DefaultPostForm(`chat_type`, `1`))
	unknownQuestionPromptJson := tool.JsonEncodeNoError(define.MenuJsonStruct{Content: lib_define.DefaultUnknownQuestionPromptContent})
	unknownQuestionPrompt := strings.TrimSpace(c.DefaultPostForm(`unknown_question_prompt`, unknownQuestionPromptJson))

	libraryQaDirectReplySwitch := cast.ToBool(c.DefaultPostForm(`library_qa_direct_reply_switch`, `true`))
	libraryQaDirectReplyScore := cast.ToFloat32(c.DefaultPostForm(`library_qa_direct_reply_score`, `0.9`))

	mixtureQaDirectReplySwitch := cast.ToBool(c.DefaultPostForm(`mixture_qa_direct_reply_switch`, `true`))
	mixtureQaDirectReplyScore := cast.ToFloat32(c.DefaultPostForm(`mixture_qa_direct_reply_score`, `0.9`))
	answerSourceSwitch := cast.ToBool(c.DefaultPostForm(`answer_source_switch`, `false`))

	enableQuestionOptimize := cast.ToBool(c.DefaultPostForm(`enable_question_optimize`, `false`))
	optimizeQuestionModelConfigId := cast.ToInt(c.DefaultPostForm(`optimize_question_model_config_id`, "0"))
	optimizeQuestionUseModel := strings.TrimSpace(c.DefaultPostForm(`optimize_question_use_model`, ``))
	optimizeQuestionDialogueBackground := strings.TrimSpace(c.DefaultPostForm(`optimize_question_dialogue_background`, ``))
	enableQuestionGuide := cast.ToBool(c.DefaultPostForm(`enable_question_guide`, `true`))
	questionGuideNum := cast.ToInt(c.DefaultPostForm(`question_guide_num`, `3`))
	enableCommonQuestion := cast.ToBool(c.DefaultPostForm(`enable_common_question`, `true`))
	commonQuestionList := strings.TrimSpace(c.DefaultPostForm(`common_question_list`, `[]`))
	thinkSwitch := strings.TrimSpace(c.DefaultPostForm(`think_switch`, `1`))
	feedbackSwitch := strings.TrimSpace(c.DefaultPostForm(`feedback_switch`, `1`))
	sensitiveWordsSwitch := strings.TrimSpace(c.DefaultPostForm(`sensitive_words_switch`, `0`))
	cacaheConfig := strings.TrimSpace(c.DefaultPostForm(`cache_config`, `{"cache_switch":1,"valid_time":86400}`))
	groupId := cast.ToInt(c.PostForm(`group_id`))
	promptRoleType := cast.ToInt(c.PostForm(`prompt_role_type`))
	opTypeRelationLibrary := cast.ToInt(c.PostForm(`op_type_relation_library`))
	//metadata
	metaSearchSwitch := cast.ToInt(c.DefaultPostForm(`meta_search_switch`, `0`))
	metaSearchType := cast.ToInt(c.DefaultPostForm(`meta_search_type`, cast.ToString(define.MetaSearchTypeAnd)))
	metaSearchConditionList := strings.TrimSpace(c.DefaultPostForm(`meta_search_condition_list`, `[]`))
	isDefault := cast.ToInt(c.DefaultPostForm(`is_default`, `1`))
	if !tool.InArrayInt(isDefault, []int{define.IsDefault, define.NotDefault}) {
		isDefault = define.IsDefault
	}
	//set default value
	if id == 0 {
		robotAvatar = define.LocalUploadPrefix + `default/robot_avatar.svg`
		if modelConfigId == 0 && len(useModel) == 0 {
			var existLlm bool
			modelConfigId, useModel, existLlm = common.GetDefaultLlmConfig(common.GetLang(c), userId)
			if !existLlm {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `not_llm_config`))))
				return
			}
		}
		if len(welcomes) == 0 {
			welcomes = i18n.Show(common.GetLang(c), `default_welcomes`)
		}
	}
	//check required
	if id < 0 || len(robotName) == 0 || len(welcomes) == 0 || modelConfigId <= 0 || len(useModel) == 0 || maxToken < 0 || topK <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if checkName != 0 {
		//Auto-increment name suffix based on max existing value
		RobotCount, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_name`, `like`, robotName+`%`).Count(`id`)
		if err == nil {
			if RobotCount > 0 {
				robotName = robotName + `_` + cast.ToString(RobotCount+1)
			}
		}
	}
	if promptStruct, err = common.CheckPromptConfig(common.GetLang(c), promptType, promptStruct); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if temperature < 0 || temperature > 2 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `temperature`))))
		return
	}
	if similarity < 0 || similarity > 1 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `similarity`))))
		return
	}

	if recallNeighborSwitch != true && recallNeighborSwitch != false {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `recall_neighbor_switch`))))
		return
	}
	if recallNeighborBeforeNum < 0 || recallNeighborBeforeNum > 5 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `recall_neighbor_before_num`))))
		return
	}
	if recallNeighborAfterNum < 0 || recallNeighborAfterNum > 5 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `recall_neighbor_after_num`))))
		return
	}

	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText, define.SearchTypeGraph}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	if id == 0 && len(rrfWeight) == 0 { //fill default values on create
		rrfWeight = tool.JsonEncodeNoError(common.GetDefaultRrfWeight(userId))
	}
	if err = common.CheckRrfWeight(rrfWeight, common.GetLang(c)); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if questionGuideNum < 1 || questionGuideNum > 10 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `question_optimize_num`))))
		return
	}
	//data check
	var count int
	var robotKey string
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	if id > 0 {
		robotKey, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`robot_key`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(robotKey) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	} else {
		count, err = m.Where(`admin_user_id`, cast.ToString(userId)).Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if count >= define.MaxRobotNum {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `max_robot_num`, define.MaxRobotNum))))
			return
		}
	}
	//format check
	welcomes, err = common.CheckMenuJson(welcomes)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	unknownQuestionPrompt, err = common.CheckMenuJson(unknownQuestionPrompt)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	if len(libraryIds) > 0 {
		//format check
		if !common.CheckIds(libraryIds) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `library_ids_err`))))
			return
		}
		//data check
		var validLibraryIds []string
		for _, libraryId := range strings.Split(libraryIds, `,`) {
			info, err := common.GetLibraryInfo(cast.ToInt(libraryId), userId)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			if len(info) != 0 {
				validLibraryIds = append(validLibraryIds, libraryId)
			}
		}
		// default library check
		defaultLibraryId, _ := msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
			Where(`id`, cast.ToString(id)).Value(`default_library_id`)
		if len(defaultLibraryId) > 0 && !tool.InArrayString(defaultLibraryId, strings.Split(libraryIds, `,`)) {
			common.FmtError(c, `default_library_remove`)
			return
		}
		libraryIds = strings.Join(validLibraryIds, `,`)
	}

	//headImg uploaded
	fileHeader, _ := c.FormFile(`robot_avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `robot_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		robotAvatar = uploadInfo.Link
	}
	//check model_config_id and use_model
	if ok := common.CheckModelIsValid(userId, modelConfigId, useModel, common.Llm); !ok {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
		return
	}

	//check optimize_question_model_config_id and optimize_question_use_model
	if optimizeQuestionModelConfigId > 0 {
		if ok := common.CheckModelIsValid(userId, optimizeQuestionModelConfigId, optimizeQuestionUseModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `optimize_question_use_model`))))
			return
		}
	}

	//check form
	if len(formIds) > 0 {
		//Check func-call capability
		if err = common.CheckSupportFuncCall(common.GetLang(c), userId, modelConfigId, useModel); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		if !common.CheckIds(formIds) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `form_ids_err`))))
			return
		}
		var validFormIds []string
		for _, formId := range strings.Split(formIds, `,`) {
			info, err := common.GetFormInfo(cast.ToInt(formId), userId)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			if len(info) != 0 {
				validFormIds = append(validFormIds, formId)
			}
		}
		formIds = strings.Join(validFormIds, `,`)
	}
	//check rerank config
	rerankStatus := cast.ToInt(c.PostForm(`rerank_status`))
	rerankModelConfigId := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	if rerankStatus != 0 || rerankModelConfigId != 0 || len(rerankUseModel) != 0 {
		if rerankModelConfigId <= 0 || len(rerankUseModel) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
			return
		}
		if ok := common.CheckModelIsValid(userId, rerankModelConfigId, rerankUseModel, common.Rerank); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `rerank_use_model`))))
			return
		}
	}
	if chatType != define.ChatTypeLibrary && chatType != define.ChatTypeDirect && chatType != define.ChatTypeMixture {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `chat_type`))))
		return
	}
	//check qa_direct_reply
	if libraryQaDirectReplyScore < 0.0 || libraryQaDirectReplyScore > 1.0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, "library_qa_direct_reply_score"))))
		return
	}
	if mixtureQaDirectReplyScore < 0.0 || mixtureQaDirectReplyScore > 1.0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, "mixture_qa_direct_reply_score"))))
		return
	}
	maxSortNum, _ := GetMaxRobotNum(userId)

	//check common_questions
	commonQuestionList, err = common.CheckCommonQuestionJson(c, commonQuestionList)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	tipsBeforeAnswerSwitch := cast.ToBool(c.DefaultPostForm(`tips_before_answer_switch`, `true`))
	tipsBeforeAnswerContent := strings.TrimSpace(c.DefaultPostForm(`tips_before_answer_content`, i18n.Show(common.GetLang(c), `thinking_please_wait`))) //"Thinking... please wait"

	// logs.Info(`tipsBeforeAnswerContent: %v`, utf8.RuneCountInString(tipsBeforeAnswerContent))
	if len(tipsBeforeAnswerContent) > 30 && utf8.RuneCountInString(tipsBeforeAnswerContent) > 10 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `tips_before_answer_content`))))
		return
	}

	// check meta search config
	metaSearchSwitch, metaSearchType, metaSearchConditionList, err = validateMetaSearchConfig(common.GetLang(c), metaSearchSwitch, metaSearchType, metaSearchConditionList)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//database dispose
	data := msql.Datas{
		`robot_name`:                            robotName,
		`robot_intro`:                           robotIntro,
		`prompt_type`:                           promptType,
		`prompt`:                                prompt,
		`prompt_struct`:                         promptStruct,
		`library_ids`:                           libraryIds,
		`form_ids`:                              formIds,
		`welcomes`:                              welcomes,
		`model_config_id`:                       modelConfigId,
		`use_model`:                             useModel,
		`rerank_status`:                         rerankStatus,
		`rerank_model_config_id`:                rerankModelConfigId,
		`rerank_use_model`:                      rerankUseModel,
		`temperature`:                           temperature,
		`max_token`:                             maxToken,
		`context_pair`:                          contextPair,
		`enable_thinking`:                       enableThinking,
		`question_multiple_switch`:              questionMultipleSwitch,
		`top_k`:                                 topK,
		`similarity`:                            similarity,
		`recall_neighbor_switch`:                recallNeighborSwitch,
		`recall_neighbor_before_num`:            recallNeighborBeforeNum,
		`recall_neighbor_after_num`:             recallNeighborAfterNum,
		`search_type`:                           searchType,
		`rrf_weight`:                            rrfWeight,
		`chat_type`:                             chatType,
		`library_qa_direct_reply_switch`:        libraryQaDirectReplySwitch,
		`library_qa_direct_reply_score`:         libraryQaDirectReplyScore,
		`mixture_qa_direct_reply_switch`:        mixtureQaDirectReplySwitch,
		`mixture_qa_direct_reply_score`:         mixtureQaDirectReplyScore,
		`answer_source_switch`:                  answerSourceSwitch,
		`enable_question_optimize`:              enableQuestionOptimize,
		`optimize_question_model_config_id`:     optimizeQuestionModelConfigId,
		`optimize_question_dialogue_background`: optimizeQuestionDialogueBackground,
		`optimize_question_use_model`:           optimizeQuestionUseModel,
		`enable_question_guide`:                 enableQuestionGuide,
		`question_guide_num`:                    questionGuideNum,
		`enable_common_question`:                enableCommonQuestion,
		`common_question_list`:                  commonQuestionList,
		`think_switch`:                          thinkSwitch,
		`feedback_switch`:                       feedbackSwitch,
		`sensitive_words_switch`:                sensitiveWordsSwitch,
		`cache_config`:                          cacaheConfig,
		`prompt_role_type`:                      promptRoleType,
		`tips_before_answer_switch`:             tipsBeforeAnswerSwitch,
		`tips_before_answer_content`:            tipsBeforeAnswerContent,
		`is_top`:                                0,
		`sort_num`:                              maxSortNum + 1,
		`meta_search_switch`:                    metaSearchSwitch,
		`meta_search_type`:                      metaSearchType,
		`meta_search_condition_list`:            metaSearchConditionList,
		`update_time`:                           tool.Time2Int(),
	}
	if len(robotAvatar) > 0 {
		data[`robot_avatar`] = robotAvatar
	}
	if len(unknownQuestionPrompt) > 0 {
		data[`unknown_question_prompt`] = unknownQuestionPrompt
	}

	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		if err == nil && opTypeRelationLibrary > 0 && isDefault == define.NotDefault {
			_ = common.SetStepFinish(userId, define.StepRelationLibrary)
		}
	} else {
		for i := 0; i < 5; i++ {
			tempKey := tool.Random(10)
			if robot, e := common.GetRobotInfo(tempKey); e == nil && len(robot) == 0 {
				robotKey = tempKey
				break
			}
			time.Sleep(time.Nanosecond) //sleep 1 ns
		}
		if len(robotKey) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		loginUserId := getLoginUserId(c)
		data[`admin_user_id`] = userId
		data[`robot_key`] = robotKey
		data[`create_time`] = data[`update_time`]
		data[`creator`] = loginUserId
		data[`group_id`] = groupId
		data[`is_default`] = isDefault
		id, err = m.Insert(data, `id`)
		// add robot api key

		if err == nil {
			if isDefault == define.NotDefault {
				_ = common.SetStepFinish(userId, define.StepCreateRobot)
			}
			addDefaultApiKey(c, robotKey)
			// _ = AddUserMangedData(loginUserId, `managed_robot_list`, id)
			go AddDefaultPermissionManage(userId, loginUserId, int(id), define.ObjectTypeRobot)
			// add default library
			_, _ = common.AddDefaultLibrary(common.GetLang(c), c.GetHeader(`token`), robotName, libraryIds, robotKey, userId)
		}
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	common.ClearMCPServerCache(userId)
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func AddFlowRobot(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).Field(`u.*,r.role_type`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	//get params
	robotName := strings.TrimSpace(c.PostForm(`robot_name`))
	checkName := cast.ToInt(c.PostForm(`check_name`))
	enName := strings.TrimSpace(c.PostForm(`en_name`))
	if !common.CheckEnName(cast.ToString(userId), enName, `0`) {
		common.FmtError(c, `param_err`, "en_name")
		return
	}
	groupId := cast.ToInt(c.PostForm(`group_id`))
	robotIntro := strings.TrimSpace(c.PostForm(`robot_intro`))
	robotAvatar := strings.TrimSpace(c.DefaultPostForm(`avatar_from_template`, define.LocalUploadPrefix+`default/workflow_avatar.svg`))
	isDefault := cast.ToInt(c.DefaultPostForm(`is_default`, cast.ToString(define.IsDefault)))
	//check required
	if len(robotName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if checkName != 0 {
		//Auto-increment name suffix based on max existing value
		RobotCount, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_name`, `like`, robotName+`%`).Count(`id`)
		if err == nil {
			if RobotCount > 0 {
				robotName = robotName + `_` + cast.ToString(RobotCount+1)
			}
		}
	}
	//data check
	var robotKey string
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	count, err := m.Where(`admin_user_id`, cast.ToString(userId)).Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if count >= define.MaxRobotNum {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `max_robot_num`, define.MaxRobotNum))))
		return
	}
	//headImg uploaded
	fileHeader, _ := c.FormFile(`robot_avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `robot_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		robotAvatar = uploadInfo.Link
	}
	maxSortNum, _ := GetMaxRobotNum(userId)
	//format check
	welcomes, _ := common.CheckMenuJson(i18n.Show(common.GetLang(c), `default_welcomes`))
	unknownQuestionPrompt, _ := common.CheckMenuJson(``)
	//database dispose
	data := msql.Datas{
		`admin_user_id`:           userId,
		`robot_name`:              robotName,
		`robot_intro`:             robotIntro,
		`robot_avatar`:            robotAvatar,
		`group_id`:                groupId,
		`application_type`:        define.ApplicationTypeFlow,
		`create_time`:             tool.Time2Int(),
		`update_time`:             tool.Time2Int(),
		`welcomes`:                welcomes,
		`unknown_question_prompt`: unknownQuestionPrompt,
		`is_top`:                  0,
		`sort_num`:                maxSortNum + 1,
		`en_name`:                 enName,
	}
	for i := 0; i < 5; i++ {
		tempKey := tool.Random(10)
		if robot, e := common.GetRobotInfo(tempKey); e == nil && len(robot) == 0 {
			robotKey = tempKey
			break
		}
		time.Sleep(time.Nanosecond) //sleep 1 ns
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	loginUserId := getLoginUserId(c)
	data[`robot_key`] = robotKey
	data[`creator`] = loginUserId
	data[`is_default`] = isDefault
	id, err := m.Insert(data, `id`)
	// add robot api key
	if err == nil {
		addDefaultApiKey(c, robotKey)
		// _ = AddUserMangedData(loginUserId, `managed_robot_list`, id)
		go AddDefaultPermissionManage(userId, loginUserId, int(id), define.ObjectTypeRobot)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func GetRobotMetaSchemaList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(userId)).
		Field(`id,library_ids`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	libraryIdsStr := strings.TrimSpace(robot[`library_ids`])
	if libraryIdsStr == `` {
		c.String(http.StatusOK, lib_web.FmtJson([]map[string]any{}, nil))
		return
	}

	//Intersect by name: return only names present in all libraries
	idArr := strings.Split(libraryIdsStr, ",")
	libIdSet := make(map[string]struct{}, len(idArr))
	libIds := make([]string, 0, len(idArr))
	for _, s := range idArr {
		s = strings.TrimSpace(s)
		if s == `` {
			continue
		}
		if _, ok := libIdSet[s]; ok {
			continue
		}
		libIdSet[s] = struct{}{}
		libIds = append(libIds, s)
	}
	if len(libIds) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson([]map[string]any{}, nil))
		return
	}

	//Deduplicate (built-in by key; custom by name)
	seenBuiltinKey := make(map[string]bool)
	result := make([]map[string]any, 0, 32)

	//Built-in meta: keep one copy
	builtinMetaSchemaList := common.GetBuiltinMetaSchemaList(common.GetLang(c))
	for _, b := range builtinMetaSchemaList {
		k := b.Key
		if !seenBuiltinKey[k] {
			seenBuiltinKey[k] = true
			result = append(result, map[string]any{
				`id`:         0,
				`name`:       b.Name,
				`key`:        b.Key,
				`type`:       b.Type,
				`is_show`:    1,
				`is_builtin`: 1,
			})
		}
	}

	//Custom meta: intersect by name (type must match)
	customList, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`library_id`, `in`, strings.Join(libIds, `,`)).
		Order(`id asc`).
		Field(`id,library_id,name,key,type,is_show`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	type agg struct {
		count    int
		typ      int
		name     string
		key      string
		conflict bool
		seenLib  map[string]struct{}
	}
	aggMap := make(map[string]*agg, 64) // name => agg
	for _, item := range customList {
		name := strings.TrimSpace(item[`name`])
		if name == `` {
			continue
		}
		k := strings.TrimSpace(item[`key`])
		if k == `` || define.IsBuiltinMetaKey(k) {
			continue
		}
		libId := strings.TrimSpace(item[`library_id`])
		if libId == `` {
			libId = cast.ToString(item[`library_id`])
		}
		typ := cast.ToInt(item[`type`])
		a, ok := aggMap[name]
		if !ok {
			a = &agg{
				typ:     typ,
				name:    name,
				key:     k,
				seenLib: make(map[string]struct{}, len(libIds)),
			}
			aggMap[name] = a
		}
		//Do not double-count duplicate names within the same library
		if _, ok := a.seenLib[libId]; ok {
			continue
		}
		a.seenLib[libId] = struct{}{}
		a.count++
		if a.typ != typ {
			a.conflict = true
		}
	}

	for _, a := range aggMap {
		if a == nil || a.conflict {
			continue
		}
		if a.count != len(libIds) {
			continue
		}
		result = append(result, map[string]any{
			`id`:         0,
			`library_id`: 0,
			`name`:       a.name,
			`key`:        a.key,
			`type`:       a.typ,
			`is_show`:    1,
			`is_builtin`: 0,
		})
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func EditExternalConfig(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt(c.PostForm(`id`))
	externalConfigH5 := strings.TrimSpace(c.PostForm(`external_config_h5`))
	externalConfigPc := strings.TrimSpace(c.PostForm(`external_config_pc`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{}
	if len(externalConfigH5) > 0 {
		data[`external_config_h5`] = externalConfigH5
	}
	if len(externalConfigPc) > 0 {
		data[`external_config_pc`] = externalConfigPc
	}
	if len(data) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	data[`update_time`] = tool.Time2Int()
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// GetDefaultRrfWeight returns default RRF weight configuration
func GetDefaultRrfWeight(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(common.GetDefaultRrfWeight(adminUserId), nil))
}

func GetRobotInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	// add default library
	if cast.ToInt(info[`default_library_id`]) <= 0 && cast.ToInt(info[`application_type`]) == define.ApplicationTypeChat {
		_, _ = common.AddDefaultLibrary(common.GetLang(c), c.GetHeader(`token`), info[`robot_name`], info[`library_ids`], info[`robot_key`], userId)
		info, _ = common.GetRobotInfo(info[`robot_key`])
	}
	if len(info[`rrf_weight`]) == 0 { //fill default values
		info[`rrf_weight`] = tool.JsonEncodeNoError(common.GetDefaultRrfWeight(userId))
	}
	if len(info[`prompt_struct`]) == 0 {
		info[`prompt_struct`] = tool.JsonEncodeNoError(common.GetEmptyPromptStruct(common.GetLang(c))) //for legacy data, default to empty
	}
	if cast.ToInt(info[`prompt_type`]) == define.PromptTypeStruct { //for legacy data normalization
		info[`prompt_struct`], _ = common.CheckPromptConfig(common.GetLang(c), define.PromptTypeStruct, info[`prompt_struct`])
	}
	//configure external service parameters
	info[`image_domain`] = define.Config.WebService[`image_domain`]
	info[`h5_domain`] = define.Config.WebService[`h5_domain`]
	info[`pc_domain`] = define.Config.WebService[`pc_domain`]
	info[`prompt_struct_default`] = common.GetDefaultPromptStruct(common.GetLang(c)) //default value for the frontend
	info[`wechat_ip`] = define.Config.WebService[`wechat_ip`]
	info[`push_wechat_kefu`] = fmt.Sprintf(`%s/push_pwd/wechat_kefu`, define.Config.WebService[`push_domain`])
	info[`push_token`] = lib_define.SignToken
	info[`push_aeskey`] = lib_define.AesKey
	c.String(http.StatusOK, lib_web.FmtJson(info, nil))
}

func DeleteRobot(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	info, err := m.Field(`robot_key`).Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	_, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: info[`robot_key`]})
	//dispose relation data
	_, err = msql.Model(`chat_ai_message`, define.Postgres).Where(`robot_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	go func() {
		err := deleteRobotRelationData(id, info[`robot_key`])
		if err != nil {
			logs.Error(err.Error())
		}
		work_flow.DeleteRobotFollow(userId, id)
	}()
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func deleteRobotRelationData(robotId int, robotKey string) error {
	if robotId <= 0 || robotKey == "" {
		return nil
	}
	err := deleteRobotApiKey(robotKey)
	err = deleteFastCommandByRobotId(robotId)
	err = deleteWorkFlowByRobotId(robotId)
	return err
}

func CreatePromptByAi(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	demand := strings.TrimSpace(c.Query(`demand`))
	if id <= 0 || len(demand) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	promptStruct, err := common.CreatePromptByAi(common.GetLang(c), demand, adminUserId, cast.ToInt(info[`model_config_id`]), info[`use_model`])
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`promptStruct`: promptStruct, `markdown`: common.BuildPromptStruct(common.GetLang(c), define.PromptTypeStruct, ``, promptStruct)}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func RobotCopy(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fromId := cast.ToInt64(c.PostForm(`from_id`))
	if fromId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	//data check
	count, err := m.Where(`admin_user_id`, cast.ToString(userId)).Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if count >= define.MaxRobotNum {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `max_robot_num`, define.MaxRobotNum))))
		return
	}
	//robot check
	info, err := m.Where(`id`, cast.ToString(fromId)).Where(`admin_user_id`, cast.ToString(userId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := make(msql.Datas)
	for key, val := range info {
		if !tool.InArrayString(key, []string{`id`, `robot_key`, `robot_name`, `start_node_key`, `work_flow_model_config_ids`}) {
			data[key] = val
		}
	}
	var robotKey string
	for i := 0; i < 5; i++ {
		tempKey := tool.Random(10)
		if robot, e := common.GetRobotInfo(tempKey); e == nil && len(robot) == 0 {
			robotKey = tempKey
			break
		}
		time.Sleep(time.Nanosecond) //sleep 1 ns
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data[`robot_key`] = robotKey
	data[`creator`] = getLoginUserId(c)
	data[`robot_name`] = createNewName(info[`robot_name`])
	data[`create_time`] = tool.Time2Int()
	data[`update_time`] = tool.Time2Int()
	if len(info[`en_name`]) > 0 {
		data[`en_name`] = tool.Random(50)
	}
	newId, err := m.Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	// add robot api key
	addDefaultApiKey(c, robotKey)
	//_ = AddUserMangedData(getLoginUserId(c), `managed_robot_list`, newId)
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	//work_flow
	if cast.ToInt(info[`application_type`]) == define.ApplicationTypeFlow {
		workFlowNodeCopy(userId, fromId, newId)
	}
	go AddDefaultPermissionManage(userId, getLoginUserId(c), int(newId), define.ObjectTypeRobot)
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func createNewName(name string) string {
	match := regexp.MustCompile(`^(.+?)(_(\d+))?$`).FindStringSubmatch(name)
	if len(match) != 4 {
		return fmt.Sprintf(`%s_%d`, name, 1)
	}
	return fmt.Sprintf(`%s_%d`, match[1], cast.ToInt(match[3])+1)
}

func workFlowNodeCopy(userId int, fromId, newId int64) {
	m := msql.Model(`work_flow_node`, define.Postgres)
	list, err := m.Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(fromId)).
		Where(`data_type`, cast.ToString(define.DataTypeDraft)).Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for _, node := range list {
		data := make(msql.Datas)
		for key, val := range node {
			if !tool.InArrayString(key, []string{`id`, `robot_id`}) {
				data[key] = val
			}
		}
		data[`robot_id`] = newId
		_, err = m.Insert(data, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
	}
}

func EditBaseInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	robotName := strings.TrimSpace(c.PostForm(`robot_name`))
	robotIntro := strings.TrimSpace(c.PostForm(`robot_intro`))
	groupId := cast.ToInt(c.PostForm(`group_id`))
	robotAvatar := ``
	//check required
	if id <= 0 || len(robotName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	robotInfo, err := common.GetRobotInfo(robotKey)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//headImg uploaded
	fileHeader, _ := c.FormFile(`robot_avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `robot_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		robotAvatar = uploadInfo.Link
	}
	//database dispose
	data := msql.Datas{
		`robot_name`:  robotName,
		`robot_intro`: robotIntro,
		`group_id`:    groupId,
		`update_time`: tool.Time2Int(),
	}
	if cast.ToInt(robotInfo[`application_type`]) == define.ApplicationTypeFlow {
		enName := strings.TrimSpace(c.PostForm(`en_name`))
		if !common.CheckEnName(cast.ToString(userId), enName, cast.ToString(id)) {
			common.FmtError(c, `param_err`, "en_name")
			return
		}
		data[`en_name`] = enName
	}
	if len(robotAvatar) > 0 {
		data[`robot_avatar`] = robotAvatar
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	common.ClearMCPServerCache(userId)
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func RelationWorkFlow(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	workFlowIds := strings.TrimSpace(c.PostForm(`work_flow_ids`))
	//check required
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robot, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`application_type`, cast.ToString(define.ApplicationTypeChat)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if len(workFlowIds) > 0 {
		if !common.CheckIds(workFlowIds) { //format check
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `work_flow_ids`))))
			return
		}
		robotIds, err := m.Where(`application_type`, cast.ToString(define.ApplicationTypeFlow)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, `in`, workFlowIds).ColumnArr(`id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		workFlowIds = strings.Join(robotIds, `,`)
	}
	if len(workFlowIds) > 0 { //check func-call capability
		if err = common.CheckSupportFuncCall(common.GetLang(c), adminUserId, cast.ToInt(robot[`model_config_id`]), robot[`use_model`]); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
	}
	//database dispose
	data := msql.Datas{
		`work_flow_ids`: workFlowIds,
		`update_time`:   tool.Time2Int(),
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robot[`robot_key`]})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robot[`robot_key`])))
}

func SetUnknownIssueSummary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	unknownSummaryStatus := cast.ToInt(c.PostForm(`unknown_summary_status`))
	unknownSummaryModelConfigId := cast.ToInt(c.PostForm(`unknown_summary_model_config_id`))
	unknownSummaryUseModel := strings.TrimSpace(c.PostForm(`unknown_summary_use_model`))
	unknownSummarySimilarity := cast.ToFloat32(c.DefaultPostForm(`unknown_summary_similarity`, `0.8`))
	//check required
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if unknownSummaryStatus != 0 || unknownSummaryModelConfigId != 0 || len(unknownSummaryUseModel) != 0 {
		if unknownSummaryModelConfigId <= 0 || len(unknownSummaryUseModel) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
			return
		}
		if ok := common.CheckModelIsValid(adminUserId, unknownSummaryModelConfigId, unknownSummaryUseModel, common.TextEmbedding); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `unknown_summary_use_model`))))
			return
		}
	}
	if unknownSummarySimilarity < 0 || unknownSummarySimilarity > 1 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `unknown_summary_similarity`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{
		`unknown_summary_status`:          unknownSummaryStatus,
		`unknown_summary_model_config_id`: unknownSummaryModelConfigId,
		`unknown_summary_use_model`:       unknownSummaryUseModel,
		`unknown_summary_similarity`:      unknownSummarySimilarity,
		`update_time`:                     tool.Time2Int(),
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func RobotAutoAdd(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	var (
		robotInfo msql.Params
		err       error
	)
	userId := getLoginUserId(c)
	if userId != adminUserId {
		common.FmtOk(c, robotInfo)
		return
	}
	robot, _ := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`application_type`, cast.ToString(define.ApplicationTypeChat)).
		Order(`id desc`).
		ColumnArr(`id`)
	if len(robot) == 0 {
		if robotInfo, err = common.RobotAutoAdd(common.GetLang(c), c.GetHeader(`token`), adminUserId); err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
	}
	common.FmtOk(c, robotInfo)
}

func RelationLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := strings.TrimSpace(c.PostForm(`id`))
	if len(id) == 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}
	defaultLibraryId := cast.ToInt(c.DefaultPostForm(`default_library_id`, `0`))
	libraryIds := strings.TrimSpace(c.PostForm(`library_ids`))
	robotInfo, _ := msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).Field(`robot_key,library_ids,default_library_id`).Find()
	if defaultLibraryId > 0 {
		robotLibraryIds := []string{cast.ToString(defaultLibraryId)}
		for _, libraryId := range strings.Split(robotInfo[`library_ids`], `,`) {
			if !tool.InArrayString(libraryId, robotLibraryIds) {
				robotLibraryIds = append(robotLibraryIds, libraryId)
			}
		}
		libraryIds = strings.Join(robotLibraryIds, `,`)
	} else {
		//data check
		defaultLibraryId := cast.ToString(robotInfo[`default_library_id`])
		if !tool.InArrayString(defaultLibraryId, strings.Split(libraryIds, `,`)) {
			common.FmtError(c, `default_library_remove`)
			return
		}
	}
	updateData := msql.Datas{
		`library_ids`: libraryIds,
		`update_time`: tool.Time2Int(),
	}
	if defaultLibraryId > 0 {
		updateData[`default_library_id`] = defaultLibraryId
	}
	if _, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Update(updateData); err != nil {
		common.FmtError(c, `sys_err`, err.Error())
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotInfo[`robot_key`]})
	common.FmtOk(c, nil)
}

func CleanRobotChatCache(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := strings.TrimSpace(c.PostForm(`id`))
	if len(id) == 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if len(robotKey) == 0 {
		common.FmtError(c, `param_lack`, `robot_key`)
		return
	}

	robotInfo, _ := msql.Model(`chat_ai_robot`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).Field(`robot_key`).Find()
	if len(robotInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if robotInfo[`robot_key`] != robotKey {
		common.FmtError(c, `robot_key_not_match`)
		return
	}

	//clear cached data
	_ = common.CleanRobotMessageCache(id, robotKey)
	common.FmtOk(c, nil)
}
