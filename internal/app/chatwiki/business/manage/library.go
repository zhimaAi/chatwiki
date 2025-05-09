// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
)

func GetLibraryList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	m := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	libraryName := strings.TrimSpace(c.Query(`library_name`))
	if len(libraryName) > 0 {
		m.Where(`library_name`, `like`, libraryName)
	}
	typ := cast.ToString(c.Query(`type`))
	if typ == "" {
		typ = fmt.Sprintf(`%v,%v`, define.GeneralLibraryType, define.QALibraryType)
	} else if !tool.InArrayInt(cast.ToInt(typ), define.LibraryTypes[:]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `type`))))
		return
	}
	m.Where(`type`, `in`, cast.ToString(typ))
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
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot, define.RoleTypeAdmin}) && cast.ToInt(typ) != define.OpenLibraryType {
		managedRobotIdList := GetUserManagedData(userId, `managed_library_list`)
		m.Where(`id`, `in`, strings.Join(managedRobotIdList, `,`))
	}

	list, err := m.
		Field(`id,type,access_rights,avatar,library_name,library_intro,avatar,graph_switch,graph_model_config_id,graph_use_model,create_time`).
		Order(`id desc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(list) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(list, nil))
		return
	}
	//stats data
	libraryIds := make([]string, 0)
	newList := make([]msql.Params, 0)
	for _, params := range list {
		if cast.ToInt(params[`type`]) == define.OpenLibraryType {
			if !checkIsPartner(c, cast.ToInt(params[`id`]), define.PartnerRightsEdit) {
				continue
			}
		}
		libraryIds = append(libraryIds, params[`id`])
		params[`file_total`], params[`file_size`] = `0`, `0`
		newList = append(newList, params)
	}
	if len(libraryIds) > 0 {
		stats, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, `in`, strings.Join(libraryIds, `,`)).Group(`library_id`).
			ColumnMap(`COUNT(1) as file_total,SUM(file_size) as file_size`, `library_id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for _, params := range newList {
			params[`library_key`] = common.BuildLibraryKey(cast.ToInt(params[`id`]), cast.ToInt(params[`create_time`]))
			if len(stats[params[`id`]]) == 0 {
				continue
			}
			params[`file_total`] = stats[params[`id`]][`file_total`]
			params[`file_size`] = stats[params[`id`]][`file_size`]
			if len(params[`avatar`]) == 0 {
				params[`avatar`] = define.LocalUploadPrefix + `default/library_avatar.png`
			}
			robotInfo, err := common.GetLibraryRobotInfo(userId, cast.ToInt(params[`id`]))
			if err != nil {
				logs.Error(err.Error())
			}
			params[`robot_nums`] = cast.ToString(len(robotInfo))
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(newList, nil))
}

func GetLibraryInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := msql.Model(`chat_ai_library`, define.Postgres).
		Alias("a").
		Join(`chat_ai_model_config b`, "a.model_config_id=b.id", "left").
		Where(`a.id`, cast.ToString(id)).
		Where(`a.admin_user_id`, cast.ToString(userId)).
		Field(`a.*`).
		Field(`b.model_define`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	data := msql.Datas{}
	if cast.ToInt(info[`ai_chunk_size`]) == 0 {
		info[`ai_chunk_size`] = cast.ToString(define.SplitAiChunkMaxSize)
	}
	for k, v := range info {
		data[k] = v
	}
	data[`is_offline`] = false
	data[`library_key`] = common.BuildLibraryKey(cast.ToInt(data[`id`]), cast.ToInt(data[`create_time`]))
	robotInfo, err := common.GetLibraryRobotInfo(userId, id)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data[`robot_nums`] = len(robotInfo)
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetLibraryRobotInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	data, err := common.GetLibraryRobotInfo(userId, id)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func CreateLibrary(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	libraryName := strings.TrimSpace(c.PostForm(`library_name`))
	libraryIntro := strings.TrimSpace(c.PostForm(`library_intro`))
	aiSummary := cast.ToInt(c.PostForm(`ai_summary`))
	summaryModelConfigId := cast.ToInt(c.PostForm(`summary_model_config_id`))
	aiSummaryModel := strings.TrimSpace(c.PostForm(`ai_summary_model`))
	typ := cast.ToInt(c.PostForm(`type`))
	accessRights := cast.ToInt(c.PostForm(`access_rights`))
	chunkType := cast.ToInt(c.PostForm(`chunk_type`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	graphSwitch := cast.ToInt(c.PostForm(`graph_switch`))
	graphModelConfigId := cast.ToInt(c.PostForm(`graph_model_config_id`))
	graphUseModel := strings.TrimSpace(c.PostForm(`graph_use_model`))
	normalChunkDefaultSeparatorsNo := cast.ToString(c.PostForm(`normal_chunk_default_separators_no`))
	normalChunkDefaultChunkSize := cast.ToInt(c.PostForm(`normal_chunk_default_chunk_size`))
	normalChunkDefaultChunkOverlap := cast.ToInt(c.PostForm(`normal_chunk_default_chunk_overlap`))
	semanticChunkDefaultChunkSize := cast.ToInt(c.PostForm(`semantic_chunk_default_chunk_size`))
	semanticChunkDefaultChunkOverlap := cast.ToInt(c.PostForm(`semantic_chunk_default_chunk_overlap`))
	semanticChunkDefaultThreshold := cast.ToInt(c.PostForm(`semantic_chunk_default_threshold`))
	AiChunkPrumpt := cast.ToString(c.PostForm(`ai_chunk_prumpt`))
	AiChunkModel := strings.TrimSpace(c.PostForm(`ai_chunk_model`))
	AiChunkModelConfigId := cast.ToInt(c.PostForm(`ai_chunk_model_config_id`))
	AiChunkSize := cast.ToInt(c.PostForm(`ai_chunk_size`))
	avatar := ""
	if len(libraryName) == 0 || !tool.InArrayInt(typ, define.LibraryTypes[:]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if summaryModelConfigId > 0 {
		summaryConfig, err := common.GetModelConfigInfo(summaryModelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(summaryConfig[`model_define`])
		if !tool.InArrayString(aiSummaryModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(summaryConfig["model_define"]) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_summary_model`))))
			return
		}
	}
	//headImg uploaded
	fileAvatar, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileAvatar, define.ImageLimitSize, userId, `library_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}
	useModelSwitch := define.SwitchOff
	if useModel != "" && modelConfigId > 0 {
		_, err := common.GetModelConfigInfo(modelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		useModelSwitch = define.SwitchOn
	}
	loginUserId := getLoginUserId(c)
	//database dispose
	data := msql.Datas{
		`admin_user_id`:                        userId,
		`creator`:                              loginUserId,
		`library_name`:                         libraryName,
		`library_intro`:                        libraryIntro,
		`ai_summary`:                           aiSummary,
		`ai_summary_model`:                     aiSummaryModel,
		`summary_model_config_id`:              summaryModelConfigId,
		`type`:                                 typ,
		`access_rights`:                        accessRights,
		`create_time`:                          tool.Time2Int(),
		`update_time`:                          tool.Time2Int(),
		`chunk_type`:                           chunkType,
		`model_config_id`:                      modelConfigId,
		`use_model`:                            useModel,
		`use_model_switch`:                     cast.ToString(useModelSwitch),
		`graph_switch`:                         graphSwitch,
		`graph_model_config_id`:                graphModelConfigId,
		`graph_use_model`:                      graphUseModel,
		`normal_chunk_default_separators_no`:   normalChunkDefaultSeparatorsNo,
		`normal_chunk_default_chunk_size`:      normalChunkDefaultChunkSize,
		`normal_chunk_default_chunk_overlap`:   normalChunkDefaultChunkOverlap,
		`semantic_chunk_default_chunk_size`:    semanticChunkDefaultChunkSize,
		`semantic_chunk_default_chunk_overlap`: semanticChunkDefaultChunkOverlap,
		`semantic_chunk_default_threshold`:     semanticChunkDefaultThreshold,
		`ai_chunk_prumpt`:                      AiChunkPrumpt,
		`ai_chunk_model`:                       AiChunkModel,
		`ai_chunk_model_config_id`:             AiChunkModelConfigId,
		`ai_chunk_size`:                        AiChunkSize,
	}
	if len(avatar) > 0 {
		data[`avatar`] = avatar
	}
	libraryId, err := msql.Model(`chat_ai_library`, define.Postgres).Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibraryCacheBuildHandler{LibraryId: int(libraryId)})
	if cast.ToInt(typ) == define.OpenLibraryType {
		partner := []string{cast.ToString(userId)}
		if userId != loginUserId {
			partner = append(partner, cast.ToString(loginUserId))
		}
		_ = common.SaveLibDocPartner(loginUserId, int(libraryId), define.PartnerRightsManage, 1, partner)
	}
	//common save
	//fileIds, err = addLibFile(c, userId, int(libraryId), cast.ToInt(typ))
	//if err != nil {
	//	c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	//	return
	//}
	_ = AddUserMangedData(loginUserId, `managed_library_list`, libraryId)

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`id`: libraryId,
		`library_key`: common.BuildLibraryKey(cast.ToInt(libraryId), cast.ToInt(data[`create_time`]))}, nil))
}

func DeleteLibrary(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := common.GetLibraryInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibraryCacheBuildHandler{LibraryId: id})
	//dispose relation data
	fileModel := msql.Model(`chat_ai_library_file`, define.Postgres)
	fileIds, err := fileModel.Where(`library_id`, cast.ToString(id)).ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = fileModel.Where(`library_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	for _, fileId := range fileIds {
		//clear cached data
		lib_redis.DelCacheData(define.Redis, &common.LibFileCacheBuildHandler{FileId: cast.ToInt(fileId)})
	}
	if common.GetNeo4jStatus(userId) {
		err = common.NewGraphDB(userId).DeleteByLibrary(id)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`library_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`chat_ai_library_file_doc`, define.Postgres).Where(`library_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func EditLibrary(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	libraryName := strings.TrimSpace(c.PostForm(`library_name`))
	libraryIntro := strings.TrimSpace(c.PostForm(`library_intro`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	graphSwitch := cast.ToInt(c.PostForm(`graph_switch`))
	graphModelConfigId := cast.ToInt(c.PostForm(`graph_model_config_id`))
	graphUseModel := strings.TrimSpace(c.PostForm(`graph_use_model`))
	useModelSwitch := cast.ToInt(c.PostForm(`use_model_switch`))
	aiSummary := cast.ToInt(c.PostForm(`ai_summary`))
	aiSummaryModel := cast.ToString(c.PostForm(`ai_summary_model`))
	summaryModelConfigId := cast.ToInt(c.PostForm(`summary_model_config_id`))
	shareUrl := strings.TrimSpace(c.PostForm(`share_url`))
	accessRights := cast.ToInt(c.PostForm(`access_rights`))
	statisticsSet := strings.TrimSpace(c.PostForm(`statistics_set`))
	typ := strings.TrimSpace(c.PostForm(`type`))
	chunkType := cast.ToInt(c.PostForm(`chunk_type`))
	normalChunkDefaultSeparatorsNo := cast.ToString(c.PostForm(`normal_chunk_default_separators_no`))
	normalChunkDefaultChunkSize := cast.ToInt(c.PostForm(`normal_chunk_default_chunk_size`))
	normalChunkDefaultChunkOverlap := cast.ToInt(c.PostForm(`normal_chunk_default_chunk_overlap`))
	semanticChunkDefaultChunkSize := cast.ToInt(c.PostForm(`semantic_chunk_default_chunk_size`))
	semanticChunkDefaultChunkOverlap := cast.ToInt(c.PostForm(`semantic_chunk_default_chunk_overlap`))
	semanticChunkDefaultThreshold := cast.ToInt(c.PostForm(`semantic_chunk_default_threshold`))
	AiChunkPrumpt := cast.ToString(c.PostForm(`ai_chunk_prumpt`))
	AiChunkModel := strings.TrimSpace(c.PostForm(`ai_chunk_model`))
	AiChunkModelConfigId := cast.ToInt(c.PostForm(`ai_chunk_model_config_id`))
	AiChunkSize := cast.ToInt(c.PostForm(`ai_chunk_size`))
	if id <= 0 || len(libraryName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if cast.ToInt(typ) != define.OpenLibraryType {
		//check model_config_id and use_model
		config, err := common.GetModelConfigInfo(modelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(useModel, modelInfo.VectorModelList) && !common.IsMultiConfModel(config["model_define"]) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
			return
		}
		if len(config) == 0 || !tool.InArrayString(common.TextEmbedding, strings.Split(config[`model_types`], `,`)) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_config_id`))))
			return
		}
	}
	if useModelSwitch == define.SwitchOn && (modelConfigId == 0 || useModel == "") {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
		return
	}
	if summaryModelConfigId > 0 {
		summaryConfig, err := common.GetModelConfigInfo(summaryModelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(summaryConfig[`model_define`])
		if !tool.InArrayString(aiSummaryModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(summaryConfig["model_define"]) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_summary_model`))))
			return
		}
	}
	if graphSwitch == 1 && graphModelConfigId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `graph_model_config_id`))))
		return
	}
	if graphSwitch == 1 && len(graphUseModel) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `graph_use_model`))))
		return
	}

	if !tool.InArrayInt(chunkType, []int{define.ChunkTypeNormal, define.ChunkTypeSemantic, define.ChunkTypeAi}) && cast.ToInt(typ) != define.OpenLibraryType {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `chunk_type`))))
		return
	}
	if chunkType == define.ChunkTypeNormal {
		if len(normalChunkDefaultSeparatorsNo) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `separators_no`))))
			return
		}
		if normalChunkDefaultChunkSize < define.SplitChunkMinSize || normalChunkDefaultChunkSize > define.SplitChunkMaxSize {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `chunk_size_err`, define.SplitChunkMinSize, define.SplitChunkMaxSize))))
			return
		}
		maxChunkOverlap := normalChunkDefaultChunkSize / 2
		if normalChunkDefaultChunkOverlap < 0 || normalChunkDefaultChunkOverlap > maxChunkOverlap {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `chunk_overlap_err`, 0, maxChunkOverlap))))
			return
		}
	}
	if chunkType == define.ChunkTypeSemantic {
		if semanticChunkDefaultChunkSize < define.SplitChunkMinSize || semanticChunkDefaultChunkSize > define.SplitChunkMaxSize {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `semantic_chunk_size_err`, define.SplitChunkMinSize, define.SplitChunkMaxSize))))
			return
		}
		maxSemanticChunkOverlap := semanticChunkDefaultChunkSize / 2
		if semanticChunkDefaultChunkOverlap > maxSemanticChunkOverlap {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `semantic_chunk_overlap_err`, 0, maxSemanticChunkOverlap))))
			return
		}
		if semanticChunkDefaultThreshold < 1 || semanticChunkDefaultThreshold > 100 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `semantic_chunk_threshold_err`, 1, 100))))
			return
		}
	}

	if chunkType == define.ChunkTypeAi {
		if ok := common.CheckModelIsValid(userId, AiChunkModelConfigId, AiChunkModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_model`))))
			return
		}
		if len(AiChunkPrumpt) == 0 || len(AiChunkPrumpt) > 500 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `ai_chunk_prumpt`))))
			return
		}
	}

	info, err := common.GetLibraryInfo(id, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//headImg uploaded
	avatar := ""
	fileAvatar, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileAvatar, define.ImageLimitSize, userId, `library_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}
	data := msql.Datas{
		`library_name`:                         libraryName,
		`library_intro`:                        libraryIntro,
		`model_config_id`:                      modelConfigId,
		`use_model`:                            useModel,
		`use_model_switch`:                     useModelSwitch,
		`ai_summary`:                           aiSummary,
		`share_url`:                            shareUrl,
		`ai_summary_model`:                     aiSummaryModel,
		`summary_model_config_id`:              summaryModelConfigId,
		`graph_switch`:                         graphSwitch,
		`graph_model_config_id`:                graphModelConfigId,
		`graph_use_model`:                      graphUseModel,
		`access_rights`:                        accessRights,
		`statistics_set`:                       statisticsSet,
		`update_time`:                          tool.Time2Int(),
		`chunk_type`:                           chunkType,
		`normal_chunk_default_separators_no`:   normalChunkDefaultSeparatorsNo,
		`normal_chunk_default_chunk_size`:      normalChunkDefaultChunkSize,
		`normal_chunk_default_chunk_overlap`:   normalChunkDefaultChunkOverlap,
		`semantic_chunk_default_chunk_size`:    semanticChunkDefaultChunkSize,
		`semantic_chunk_default_chunk_overlap`: semanticChunkDefaultChunkOverlap,
		`semantic_chunk_default_threshold`:     semanticChunkDefaultThreshold,
		`ai_chunk_prumpt`:                      AiChunkPrumpt,
		`ai_chunk_model`:                       AiChunkModel,
		`ai_chunk_model_config_id`:             AiChunkModelConfigId,
		`ai_chunk_size`:                        AiChunkSize,
	}
	if len(avatar) > 0 {
		data[`avatar`] = avatar
	}
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	// embedding new vector
	if cast.ToInt(info[`type`]) == define.OpenLibraryType && useModelSwitch == define.SwitchOn {
		if cast.ToInt(info[`use_model_switch`]) != useModelSwitch || (info[`use_model`] != useModel || cast.ToInt(info[`model_config_id`]) != modelConfigId) {
			if err = common.AddFileDataIndex(id, cast.ToInt(info[`admin_user_id`])); err == nil {
				go common.EmbeddingNewVector(id, cast.ToInt(info[`admin_user_id`]))
			}
		}
	} else if info[`use_model`] != useModel || cast.ToInt(info[`model_config_id`]) != modelConfigId {
		go common.EmbeddingNewVector(id, cast.ToInt(info[`admin_user_id`]))
	}

	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibraryCacheBuildHandler{LibraryId: id})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func LibraryRecallTest(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	libraryIds := cast.ToString(c.PostForm(`id`))
	question := strings.TrimSpace(c.PostForm(`question`))
	size := cast.ToInt(c.PostForm(`size`))
	similarity := cast.ToFloat64(c.PostForm(`similarity`))
	searchType := cast.ToInt(c.PostForm(`search_type`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	rerankStatus := strings.TrimSpace(c.DefaultPostForm(`rerank_status`, `1`))
	recallType := cast.ToString(c.PostForm(`recall_type`))
	if len(libraryIds) <= 0 || len(question) == 0 || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if searchType != define.SearchTypeMixed && searchType != define.SearchTypeVector && searchType != define.SearchTypeFullText && searchType != define.SearchTypeGraph {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	if modelConfigId > 0 || useModel != "" {
		//check model_config_id and use_model
		config, err := common.GetModelConfigInfo(modelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(useModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(config["model_define"]) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
			return
		}
		if len(config) == 0 || !tool.InArrayString(common.Llm, strings.Split(config[`model_types`], `,`)) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_config_id`))))
			return
		}
	}
	robot := msql.Params{
		`recall_type`: recallType,
	}
	for _, libraryId := range strings.Split(libraryIds, `,`) {
		info, err := common.GetLibraryInfo(cast.ToInt(libraryId), userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(info) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		robotName, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`rerank_status`, `1`).Where(`rerank_model_config_id`, cast.ToString(rerankModelConfigID)).Value(`robot_name`)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		}
		if rerankModelConfigID > 0 && cast.ToInt(rerankStatus) == define.SwitchOn {
			robot[`rerank_status`] = cast.ToString(rerankStatus)
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

	list, _, err := common.GetMatchLibraryParagraphList("", "", question, []string{}, libraryIds, size, similarity, searchType, robot)
	for _, item := range list {
		library, err := common.GetLibraryInfo(cast.ToInt(item[`library_id`]), userId)
		if err != nil {
			logs.Error(err.Error())
		}
		item[`library_name`] = library[`library_name`]
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, err))
}

func RelationRobot(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotIds := strings.TrimSpace(c.PostForm(`robot_ids`))
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	if len(robotIds) == 0 || libraryId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	data, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_key,library_ids`).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, `in`, robotIds).Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	for _, item := range data {
		ids := strings.Split(item[`library_ids`], ",")
		if !tool.InArrayString(cast.ToString(libraryId), ids) {
			ids = append(ids, cast.ToString(libraryId))
		}
		_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, cast.ToString(item[`id`])).Update(msql.Datas{`library_ids`: strings.Join(ids, ",")})
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: item[`robot_key`]})
	}
	common.FmtOk(c, nil)
}
