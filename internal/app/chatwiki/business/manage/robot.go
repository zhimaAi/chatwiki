// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetRobotList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	m := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_name,robot_intro,robot_avatar,robot_key,application_type`).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`)

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
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot, define.RoleTypeAdmin}) {
		managedRobotIdList := GetUserManagedData(userId, `managed_robot_list`)
		m.Where(`id`, `in`, strings.Join(managedRobotIdList, `,`))
	}

	list, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
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
	robotIntro := strings.TrimSpace(c.PostForm(`robot_intro`))
	robotAvatar := ``
	promptType := cast.ToInt(c.DefaultPostForm(`prompt_type`, `1`))
	prompt := strings.TrimSpace(c.PostForm(`prompt`))
	promptStruct := strings.TrimSpace(c.DefaultPostForm(`prompt_struct`, common.GetDefaultPromptStruct()))
	libraryIds := strings.TrimSpace(c.PostForm(`library_ids`))
	formIds := strings.TrimSpace(c.PostForm(`form_ids`))
	welcomes := strings.TrimSpace(c.PostForm(`welcomes`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	temperature := cast.ToFloat32(c.DefaultPostForm(`temperature`, `0.5`))
	maxToken := cast.ToInt(c.DefaultPostForm(`max_token`, `2000`))
	contextPair := cast.ToInt(c.DefaultPostForm(`context_pair`, `6`))
	topK := cast.ToInt(c.DefaultPostForm(`top_k`, `5`))
	similarity := cast.ToFloat32(c.DefaultPostForm(`similarity`, `0.6`))
	searchType := cast.ToInt(c.DefaultPostForm(`search_type`, `1`))
	chatType := cast.ToInt(c.DefaultPostForm(`chat_type`, `1`))
	unknownQuestionPrompt := strings.TrimSpace(c.DefaultPostForm(`unknown_question_prompt`, `{"content":"哎呀，这个问题我暂时还不太清楚呢～（对手指）"}`))

	libraryQaDirectReplySwitch := cast.ToBool(c.PostForm(`library_qa_direct_reply_switch`))
	libraryQaDirectReplyScore := cast.ToFloat32(c.PostForm(`library_qa_direct_reply_score`))

	mixtureQaDirectReplySwitch := cast.ToBool(c.PostForm(`mixture_qa_direct_reply_switch`))
	mixtureQaDirectReplyScore := cast.ToFloat32(c.PostForm(`mixture_qa_direct_reply_score`))
	answerSourceSwitch := cast.ToBool(c.DefaultPostForm(`answer_source_switch`, `true`))

	enableQuestionOptimize := cast.ToBool(c.DefaultPostForm(`enable_question_optimize`, `false`))
	enableQuestionGuide := cast.ToBool(c.DefaultPostForm(`enable_question_guide`, `true`))
	questionGuideNum := cast.ToInt(c.DefaultPostForm(`question_guide_num`, `3`))
	enableCommonQuestion := cast.ToBool(c.DefaultPostForm(`enable_common_question`, `true`))
	commonQuestionList := strings.TrimSpace(c.DefaultPostForm(`common_question_list`, `[]`))
	thinkSwitch := strings.TrimSpace(c.DefaultPostForm(`think_switch`, `1`))
	feedbackSwitch := strings.TrimSpace(c.DefaultPostForm(`feedback_switch`, `1`))
	sensitiveWordsSwitch := strings.TrimSpace(c.DefaultPostForm(`sensitive_words_switch`, `0`))

	//set default value
	if id == 0 {
		robotAvatar = define.LocalUploadPrefix + `default/robot_avatar.svg`
		if modelConfigId == 0 && len(useModel) == 0 {
			var existLlm bool
			modelConfigId, useModel, existLlm = common.GetDefaultLlmConfig(userId)
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
	if promptStruct, err = common.CheckPromptConfig(promptType, promptStruct); err != nil {
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
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText, define.SearchTypeGraph}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
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
		libraryIds = strings.Join(validLibraryIds, `,`)
	}

	//headImg uploaded
	fileHeader, _ := c.FormFile(`robot_avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `robot_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		robotAvatar = uploadInfo.Link
	}
	//check model_config_id and use_model
	config, err := common.GetModelConfigInfo(modelConfigId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(config) == 0 || !tool.InArrayString(common.Llm, strings.Split(config[`model_types`], `,`)) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_config_id`))))
		return
	}
	modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
	if !tool.InArrayString(useModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(config["model_define"]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
		return
	}
	//check form
	if len(formIds) > 0 {
		if modelInfo.SupportedFunctionCallList == nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `does_not_support_form`))))
			return
		}
		if modelInfo.CheckFancCallRequest != nil {
			if err = modelInfo.CheckFancCallRequest(modelInfo, config, useModel); err != nil {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `form_ids_err`))))
				return
			}
		} else if !tool.InArrayString(useModel, modelInfo.SupportedFunctionCallList) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `form_ids_err`))))
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
		config, err := common.GetModelConfigInfo(rerankModelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(config) == 0 || !tool.InArrayString(common.Rerank, strings.Split(config[`model_types`], `,`)) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `rerank_model_config_id`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(rerankUseModel, modelInfo.RerankModelList) {
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

	//check common_questions
	commonQuestionList, err = common.CheckCommonQuestionJson(c, commonQuestionList)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//database dispose
	data := msql.Datas{
		`robot_name`:               robotName,
		`robot_intro`:              robotIntro,
		`prompt_type`:              promptType,
		`prompt`:                   prompt,
		`prompt_struct`:            promptStruct,
		`library_ids`:              libraryIds,
		`form_ids`:                 formIds,
		`welcomes`:                 welcomes,
		`model_config_id`:          modelConfigId,
		`use_model`:                useModel,
		`rerank_status`:            rerankStatus,
		`rerank_model_config_id`:   rerankModelConfigId,
		`rerank_use_model`:         rerankUseModel,
		`temperature`:              temperature,
		`max_token`:                maxToken,
		`context_pair`:             contextPair,
		`top_k`:                    topK,
		`similarity`:               similarity,
		`search_type`:              searchType,
		`chat_type`:                chatType,
		`answer_source_switch`:     answerSourceSwitch,
		`enable_question_optimize`: enableQuestionOptimize,
		`enable_question_guide`:    enableQuestionGuide,
		`question_guide_num`:       questionGuideNum,
		`enable_common_question`:   enableCommonQuestion,
		`common_question_list`:     commonQuestionList,
		`think_switch`:             thinkSwitch,
		`feedback_switch`:          feedbackSwitch,
		`sensitive_words_switch`:   sensitiveWordsSwitch,
		`update_time`:              tool.Time2Int(),
	}
	if len(robotAvatar) > 0 {
		data[`robot_avatar`] = robotAvatar
	}
	if len(unknownQuestionPrompt) > 0 {
		data[`unknown_question_prompt`] = unknownQuestionPrompt
	}
	if chatType == define.ChatTypeLibrary {
		data[`library_qa_direct_reply_switch`] = libraryQaDirectReplySwitch
		if libraryQaDirectReplySwitch {
			data[`library_qa_direct_reply_score`] = libraryQaDirectReplyScore
		}
	}
	if chatType == define.ChatTypeMixture {
		data[`mixture_qa_direct_reply_switch`] = mixtureQaDirectReplySwitch
		if mixtureQaDirectReplySwitch {
			data[`mixture_qa_direct_reply_score`] = mixtureQaDirectReplyScore
		}
	}

	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
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
		data[`admin_user_id`] = userId
		data[`robot_key`] = robotKey
		data[`create_time`] = data[`update_time`]
		id, err = m.Insert(data, `id`)
		// add robot api key
		if err == nil {
			_ = AddUserMangedData(getLoginUserId(c), `managed_robot_list`, id)
		}

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
	robotIntro := strings.TrimSpace(c.PostForm(`robot_intro`))
	robotAvatar := define.LocalUploadPrefix + `default/workflow_avatar.svg`
	//check required
	if len(robotName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
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
	//format check
	welcomes, _ := common.CheckMenuJson(i18n.Show(common.GetLang(c), `default_welcomes`))
	unknownQuestionPrompt, _ := common.CheckMenuJson(``)
	//database dispose
	data := msql.Datas{
		`admin_user_id`:           userId,
		`robot_name`:              robotName,
		`robot_intro`:             robotIntro,
		`robot_avatar`:            robotAvatar,
		`application_type`:        define.ApplicationTypeFlow,
		`create_time`:             tool.Time2Int(),
		`update_time`:             tool.Time2Int(),
		`welcomes`:                welcomes,
		`unknown_question_prompt`: unknownQuestionPrompt,
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
	data[`robot_key`] = robotKey
	id, err := m.Insert(data, `id`)
	// add robot api key
	if err == nil {
		_ = AddUserMangedData(getLoginUserId(c), `managed_robot_list`, id)
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
	if len(info[`prompt_struct`]) == 0 {
		info[`prompt_struct`] = tool.JsonEncodeNoError(common.GetEmptyPromptStruct()) //旧数据默认给空值
	}
	//configure external service parameters
	info[`h5_domain`] = define.Config.WebService[`h5_domain`]
	info[`pc_domain`] = define.Config.WebService[`pc_domain`]
	info[`prompt_struct_default`] = common.GetDefaultPromptStruct() //提供给前端的默认值

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
	}()
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func deleteRobotRelationData(robotId int, robotKey string) error {
	if robotId <= 0 || robotKey == "" {
		return nil
	}
	err := deleteFastCommandByRobotId(robotId)
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
	promptStruct, err := common.CreatePromptByAi(demand, adminUserId, cast.ToInt(info[`model_config_id`]), info[`use_model`])
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	data := map[string]any{`promptStruct`: promptStruct, `markdown`: common.BuildPromptStruct(define.PromptTypeStruct, ``, promptStruct)}
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
	data[`robot_name`] = createNewName(info[`robot_name`])
	data[`create_time`] = tool.Time2Int()
	data[`update_time`] = tool.Time2Int()
	newId, err := m.Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	//work_flow
	if cast.ToInt(info[`application_type`]) == define.ApplicationTypeFlow {
		workFlowNodeCopy(userId, fromId, newId)
	}
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
		`update_time`: tool.Time2Int(),
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
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
