// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type BridgeLibraryListReq struct {
	Type         string `form:"type"`
	LibraryName  string `form:"library_name"`
	Ids          string `form:"ids"`
	ShowOpenDocs string `form:"show_open_docs"`
}

func BridgeGetLibraryList(adminUserId, userId int, lang string, req *BridgeLibraryListReq) ([]msql.Params, int, error) {
	m := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	ids := strings.TrimSpace(req.Ids)
	if len(ids) > 0 {
		m.Where(`id`, `in`, ids)
	}
	libraryName := strings.TrimSpace(req.LibraryName)
	if len(libraryName) > 0 {
		m.Where(`library_name`, `like`, libraryName)
	}
	typ := cast.ToString(req.Type)
	showOpenDocs := cast.ToInt(req.ShowOpenDocs)
	if typ == "" {
		typ = fmt.Sprintf(`%v,%v`, define.GeneralLibraryType, define.QALibraryType)
	} else if !tool.InArrayInt(cast.ToInt(typ), define.LibraryTypes[:]) {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `type`))
	}
	if showOpenDocs == define.SwitchOn {
		m.Where(fmt.Sprintf(`(type in (%v) or (type=%v and use_model_switch = %v))`, typ, define.OpenLibraryType, define.SwitchOn))
	} else {
		m.Where(`type`, `in`, cast.ToString(typ))
	}
	if userId <= 0 {
		return nil, http.StatusUnauthorized, errors.New(i18n.Show(lang, `user_no_login`))
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(userInfo) == 0 {
		return nil, http.StatusUnauthorized, errors.New(i18n.Show(lang, `user_no_login`))
	}
	if cast.ToInt(typ) != define.OpenLibraryType {
		// managedLibraryIdList := GetUserManagedData(userId, `managed_library_list`)
		if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot}) {
			managedLibraryIdList := []string{`0`}
			permissionData, _ := common.GetAllPermissionManage(adminUserId, cast.ToString(userId), define.IdentityTypeUser, define.ObjectTypeLibrary)
			for _, permission := range permissionData {
				managedLibraryIdList = append(managedLibraryIdList, cast.ToString(permission[`object_id`]))
			}
			//m.Where(`id`, `in`, strings.Join(managedLibraryIdList, `,`))
		}
	}
	list, err := m.
		Field(`id,type,access_rights,avatar,library_name,library_intro,avatar,graph_switch,graph_model_config_id,model_config_id,use_model,graph_use_model,create_time,group_id`).
		Order(`id desc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(list) == 0 {
		return list, 0, nil
	}
	//stats data
	libraryIds := make([]string, 0)
	newList := make([]msql.Params, 0)
	for _, params := range list {
		if cast.ToInt(params[`type`]) == define.OpenLibraryType {
			if !checkIsPartner2(userId, cast.ToInt(params[`id`]), define.PartnerRightsEdit) {
				continue
			}
		}
		libraryIds = append(libraryIds, params[`id`])
		params[`file_total`], params[`file_size`] = `0`, `0`
		newList = append(newList, params)
	}
	if len(libraryIds) > 0 {
		stats, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`delete_time`, `0`).
			Where(`library_id`, `in`, strings.Join(libraryIds, `,`)).Group(`library_id`).
			ColumnMap(`COUNT(1) as file_total,SUM(file_size) as file_size`, `library_id`)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		for _, params := range newList {
			params[`library_key`] = common.BuildLibraryKey(cast.ToInt(params[`id`]), cast.ToInt(params[`create_time`]))
			if data, ok := stats[params[`id`]]; ok && len(data) > 0 {
				params[`file_total`] = stats[params[`id`]][`file_total`]
				params[`file_size`] = stats[params[`id`]][`file_size`]
			}
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
	return newList, 0, nil
}

type BridgeLibraryListGroupReq struct {
	Type string `form:"type"`
}

func BridgeGetLibraryListGroup(adminUserId, userId int, lang string, req *BridgeLibraryListGroupReq) ([]msql.Params, int, error) {
	m := msql.Model(`chat_ai_library_list_group`, define.Postgres)
	wheres := [][]string{{`admin_user_id`, cast.ToString(adminUserId)}}
	list, err := m.Where2(wheres).Field(`id,group_name`).Order(`sort desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	list = append([]msql.Params{{`id`: `0`, `group_name`: `未分组`}}, list...)
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(userInfo) == 0 {
		return nil, http.StatusUnauthorized, errors.New(i18n.Show(lang, `user_no_login`))
	}
	//check permission
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot}) {
		managedRobotIdList := []string{`0`}
		permissionData, _ := common.GetAllPermissionManage(adminUserId, cast.ToString(userId), define.IdentityTypeUser, define.ObjectTypeLibrary)
		for _, permission := range permissionData {
			managedRobotIdList = append(managedRobotIdList, cast.ToString(permission[`object_id`]))
		}
		//wheres = append(wheres, []string{`id`, `in`, strings.Join(managedRobotIdList, `,`)})
	}
	if req.Type == "" {
		req.Type = fmt.Sprintf(`%v,%v`, define.GeneralLibraryType, define.QALibraryType)
	} else {
		req.Type = fmt.Sprintf(`%v`, cast.ToInt(req.Type))
	}
	//统计数据
	stats, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where2(wheres).
		Where(`type`, `in`, req.Type).
		Group(`group_id`).ColumnObj(`COUNT(1) AS total`, `group_id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	for i, params := range list {
		list[i][`total`] = cast.ToString(cast.ToInt(stats[params[`id`]]))
	}
	return list, 0, nil
}

type BridgeCreateLibraryReq struct {
	LibraryName                      string `form:"library_name"`
	LibraryIntro                     string `form:"library_intro"`
	AiSummary                        string `form:"ai_summary"`
	SummaryModelConfigId             string `form:"summary_model_config_id"`
	AiSummaryModel                   string `form:"ai_summary_model"`
	Type                             string `form:"type"`
	AccessRights                     string `form:"access_rights"`
	AvatarFromTemplate               string `form:"avatar_from_template"`
	FileAvatar                       *multipart.FileHeader
	ModelConfigId                    string `form:"model_config_id"`
	UseModel                         string `form:"use_model"`
	ChunkType                        string `form:"chunk_type"`
	GraphSwitch                      string `form:"graph_switch"`
	GraphModelConfigId               string `form:"graph_model_config_id"`
	GraphUseModel                    string `form:"graph_use_model"`
	GroupId                          string `form:"group_id"`
	NormalChunkDefaultSeparatorsNo   string `form:"normal_chunk_default_separators_no"`
	NormalChunkDefaultChunkSize      string `form:"normal_chunk_default_chunk_size"`
	NormalChunkDefaultChunkOverlap   string `form:"normal_chunk_default_chunk_overlap"`
	NormalChunkDefaultNotMergedText  string `form:"normal_chunk_default_not_merged_text"`
	SemanticChunkDefaultChunkSize    string `form:"semantic_chunk_default_chunk_size"`
	SemanticChunkDefaultChunkOverlap string `form:"semantic_chunk_default_chunk_overlap"`
	SemanticChunkDefaultThreshold    string `form:"semantic_chunk_default_threshold"`
	AiChunkPrumpt                    string `form:"ai_chunk_prumpt"`
	AiChunkModel                     string `form:"ai_chunk_model"`
	AiChunkModelConfigId             string `form:"ai_chunk_model_config_id"`
	AiChunkSize                      string `form:"ai_chunk_size"`
	QaIndexType                      string `form:"qa_index_type"`
	FatherChunkParagraphType         string `form:"father_chunk_paragraph_type"`
	FatherChunkSeparatorsNo          string `form:"father_chunk_separators_no"`
	FatherChunkChunkSize             string `form:"father_chunk_chunk_size"`
	SonChunkSeparatorsNo             string `form:"son_chunk_separators_no"`
	SonChunkChunkSize                string `form:"son_chunk_chunk_size"`
}

func BridgeCreateLibrary(adminUserId, loginUserId int, lang string, req *BridgeCreateLibraryReq) (map[string]any, int, error) {
	//get params
	libraryName := strings.TrimSpace(req.LibraryName)
	libraryIntro := strings.TrimSpace(req.LibraryIntro)
	aiSummary := cast.ToInt(req.AiSummary)
	summaryModelConfigId := cast.ToInt(req.SummaryModelConfigId)
	aiSummaryModel := strings.TrimSpace(req.AiSummaryModel)
	typ := cast.ToInt(req.Type)
	accessRights := cast.ToInt(req.AccessRights)
	avatar := strings.TrimSpace(req.AvatarFromTemplate)
	modelConfigId := cast.ToInt(req.ModelConfigId)
	useModel := strings.TrimSpace(req.UseModel)
	chunkType := cast.ToInt(req.ChunkType)
	graphSwitch := cast.ToInt(req.GraphSwitch)
	graphModelConfigId := cast.ToInt(req.GraphModelConfigId)
	graphUseModel := strings.TrimSpace(req.GraphUseModel)
	normalChunkDefaultSeparatorsNo := cast.ToString(req.NormalChunkDefaultSeparatorsNo)
	normalChunkDefaultChunkSize := cast.ToInt(req.NormalChunkDefaultChunkSize)
	normalChunkDefaultChunkOverlap := cast.ToInt(req.NormalChunkDefaultChunkOverlap)
	normalChunkDefaultNotMergedText := cast.ToBool(req.NormalChunkDefaultNotMergedText)
	semanticChunkDefaultChunkSize := cast.ToInt(req.SemanticChunkDefaultChunkSize)
	semanticChunkDefaultChunkOverlap := cast.ToInt(req.SemanticChunkDefaultChunkOverlap)
	semanticChunkDefaultThreshold := cast.ToInt(req.SemanticChunkDefaultThreshold)
	AiChunkPrumpt := cast.ToString(req.AiChunkPrumpt)
	AiChunkModel := strings.TrimSpace(req.AiChunkModel)
	AiChunkModelConfigId := cast.ToInt(req.AiChunkModelConfigId)
	AiChunkSize := cast.ToInt(req.AiChunkSize)
	qaIndexType := cast.ToInt(req.QaIndexType)
	groupId := cast.ToInt(req.GroupId)
	fatherChunkParagraphType := cast.ToInt(req.FatherChunkParagraphType)
	fatherChunkSeparatorsNo := strings.TrimSpace(req.FatherChunkSeparatorsNo)
	fatherChunkChunkSize := cast.ToInt(req.FatherChunkChunkSize)
	sonChunkSeparatorsNo := strings.TrimSpace(req.SonChunkSeparatorsNo)
	sonChunkChunkSize := cast.ToInt(req.SonChunkChunkSize)
	if len(libraryName) == 0 || !tool.InArrayInt(typ, define.LibraryTypes[:]) {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	if typ != define.OpenLibraryType && cast.ToInt(modelConfigId) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `model_config_id`))
	}
	if typ != define.OpenLibraryType && useModel == `` {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `use_model`))
	}
	chunkParam := define.ChunkParam{
		ChunkType:                        req.ChunkType,
		NormalChunkDefaultSeparatorsNo:   req.NormalChunkDefaultSeparatorsNo,
		NormalChunkDefaultChunkSize:      req.NormalChunkDefaultChunkSize,
		NormalChunkDefaultChunkOverlap:   req.NormalChunkDefaultChunkOverlap,
		NormalChunkDefaultNotMergedText:  req.NormalChunkDefaultNotMergedText,
		SemanticChunkDefaultChunkSize:    req.SemanticChunkDefaultChunkSize,
		SemanticChunkDefaultChunkOverlap: req.SemanticChunkDefaultChunkOverlap,
		SemanticChunkDefaultThreshold:    req.SemanticChunkDefaultThreshold,
		AiChunkPrumpt:                    AiChunkPrumpt,
		AiChunkModel:                     AiChunkModel,
		AiChunkModelConfigId:             req.AiChunkModelConfigId,
		AiChunkSize:                      req.AiChunkSize,
		QaIndexType:                      req.QaIndexType,
		FatherChunkParagraphType:         req.FatherChunkParagraphType,
		FatherChunkSeparatorsNo:          req.FatherChunkSeparatorsNo,
		FatherChunkChunkSize:             req.FatherChunkChunkSize,
		SonChunkSeparatorsNo:             req.SonChunkSeparatorsNo,
		SonChunkChunkSize:                req.SonChunkChunkSize,
	}
	err := ValidateChunkParam(adminUserId, &chunkParam, cast.ToString(typ), lang)
	if err != nil {
		return nil, -1, err
	}
	if summaryModelConfigId > 0 {
		summaryConfig, err := common.GetModelConfigInfo(summaryModelConfigId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		modelInfo, _ := common.GetModelInfoByDefine(summaryConfig[`model_define`])
		if !tool.InArrayString(aiSummaryModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(summaryConfig["model_define"]) {
			return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `ai_summary_model`))
		}
	}
	//headImg uploaded
	uploadInfo, err := common.SaveUploadedFile(req.FileAvatar, define.ImageLimitSize, adminUserId, `library_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}
	useModelSwitch := define.SwitchOff
	if useModel != "" && modelConfigId > 0 {
		_, err := common.GetModelConfigInfo(modelConfigId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		useModelSwitch = define.SwitchOn
	}
	//database dispose
	data := msql.Datas{
		`admin_user_id`:                        adminUserId,
		`creator`:                              loginUserId,
		`library_name`:                         libraryName,
		`library_intro`:                        libraryIntro,
		`model_config_id`:                      modelConfigId,
		`use_model`:                            useModel,
		`ai_summary`:                           aiSummary,
		`ai_summary_model`:                     aiSummaryModel,
		`summary_model_config_id`:              summaryModelConfigId,
		`type`:                                 typ,
		`access_rights`:                        accessRights,
		`create_time`:                          tool.Time2Int(),
		`update_time`:                          tool.Time2Int(),
		`chunk_type`:                           chunkType,
		`use_model_switch`:                     cast.ToString(useModelSwitch),
		`graph_switch`:                         graphSwitch,
		`graph_model_config_id`:                graphModelConfigId,
		`graph_use_model`:                      graphUseModel,
		`normal_chunk_default_separators_no`:   normalChunkDefaultSeparatorsNo,
		`normal_chunk_default_chunk_size`:      normalChunkDefaultChunkSize,
		`normal_chunk_default_chunk_overlap`:   normalChunkDefaultChunkOverlap,
		`normal_chunk_default_not_merged_text`: normalChunkDefaultNotMergedText,
		`semantic_chunk_default_chunk_size`:    semanticChunkDefaultChunkSize,
		`semantic_chunk_default_chunk_overlap`: semanticChunkDefaultChunkOverlap,
		`semantic_chunk_default_threshold`:     semanticChunkDefaultThreshold,
		`ai_chunk_prumpt`:                      AiChunkPrumpt,
		`ai_chunk_model`:                       AiChunkModel,
		`ai_chunk_model_config_id`:             AiChunkModelConfigId,
		`ai_chunk_size`:                        AiChunkSize,
		`qa_index_type`:                        qaIndexType,
		`group_id`:                             groupId,
		`father_chunk_paragraph_type`:          fatherChunkParagraphType,
		`father_chunk_separators_no`:           fatherChunkSeparatorsNo,
		`father_chunk_chunk_size`:              fatherChunkChunkSize,
		`son_chunk_separators_no`:              sonChunkSeparatorsNo,
		`son_chunk_chunk_size`:                 sonChunkChunkSize,
	}
	if len(avatar) > 0 {
		data[`avatar`] = avatar
	}
	libraryId, err := msql.Model(`chat_ai_library`, define.Postgres).Insert(data, `id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibraryCacheBuildHandler{LibraryId: int(libraryId)})
	if cast.ToInt(typ) == define.OpenLibraryType {
		partner := []string{cast.ToString(adminUserId)}
		if adminUserId != loginUserId {
			partner = append(partner, cast.ToString(loginUserId))
		}
		_ = common.SaveLibDocPartner(loginUserId, int(libraryId), define.PartnerRightsManage, 1, partner)
	} else {
		go AddDefaultPermissionManage(adminUserId, loginUserId, int(libraryId), define.ObjectTypeLibrary)
	}
	//common save
	//fileIds, err = addLibFile(c, userId, int(libraryId), cast.ToInt(typ))
	//if err != nil {
	//	return nil,-1, err))
	//	return
	//}
	// _ = AddUserMangedData(loginUserId, `managed_library_list`, libraryId)
	return map[string]any{`id`: cast.ToString(libraryId),
		`library_key`: common.BuildLibraryKey(cast.ToInt(libraryId), cast.ToInt(data[`create_time`]))}, 0, nil
}

type BridgeEditLibraryReq struct {
	Id                               string `form:"id"`
	LibraryName                      string `form:"library_name"`
	LibraryIntro                     string `form:"library_intro"`
	AiSummary                        string `form:"ai_summary"`
	SummaryModelConfigId             string `form:"summary_model_config_id"`
	AiSummaryModel                   string `form:"ai_summary_model"`
	ShareUrl                         string `form:"share_url"`
	Type                             string `form:"type"`
	AccessRights                     string `form:"access_rights"`
	AvatarFromTemplate               string `form:"avatar_from_template"`
	FileAvatar                       *multipart.FileHeader
	ModelConfigId                    string `form:"model_config_id"`
	UseModel                         string `form:"use_model"`
	UseModelSwitch                   string `form:"use_model_switch"`
	ChunkType                        string `form:"chunk_type"`
	GraphSwitch                      string `form:"graph_switch"`
	GraphModelConfigId               string `form:"graph_model_config_id"`
	GraphUseModel                    string `form:"graph_use_model"`
	GroupId                          string `form:"group_id"`
	NormalChunkDefaultSeparatorsNo   string `form:"normal_chunk_default_separators_no"`
	NormalChunkDefaultChunkSize      string `form:"normal_chunk_default_chunk_size"`
	NormalChunkDefaultChunkOverlap   string `form:"normal_chunk_default_chunk_overlap"`
	NormalChunkDefaultNotMergedText  string `form:"normal_chunk_default_not_merged_text"`
	SemanticChunkDefaultChunkSize    string `form:"semantic_chunk_default_chunk_size"`
	SemanticChunkDefaultChunkOverlap string `form:"semantic_chunk_default_chunk_overlap"`
	SemanticChunkDefaultThreshold    string `form:"semantic_chunk_default_threshold"`
	AiChunkPrumpt                    string `form:"ai_chunk_prumpt"`
	AiChunkModel                     string `form:"ai_chunk_model"`
	AiChunkModelConfigId             string `form:"ai_chunk_model_config_id"`
	AiChunkSize                      string `form:"ai_chunk_size"`
	QaIndexType                      string `form:"qa_index_type"`
	FatherChunkParagraphType         string `form:"father_chunk_paragraph_type"`
	FatherChunkSeparatorsNo          string `form:"father_chunk_separators_no"`
	FatherChunkChunkSize             string `form:"father_chunk_chunk_size"`
	SonChunkSeparatorsNo             string `form:"son_chunk_separators_no"`
	SonChunkChunkSize                string `form:"son_chunk_chunk_size"`
	StatisticsSet                    string `form:"statistics_set"`
	IconTemplateConfigId             string `form:"icon_template_config_id"`
}

func BridgeEditLibrary(adminUserId, loginUserId int, lang string, req *BridgeEditLibraryReq) (map[string]any, int, error) {
	id := cast.ToInt(req.Id)
	libraryName := strings.TrimSpace(req.LibraryName)
	libraryIntro := strings.TrimSpace(req.LibraryIntro)
	modelConfigId := cast.ToInt(req.ModelConfigId)
	useModel := strings.TrimSpace(req.UseModel)
	graphSwitch := cast.ToInt(req.GraphSwitch)
	graphModelConfigId := cast.ToInt(req.GraphModelConfigId)
	graphUseModel := strings.TrimSpace(req.GraphUseModel)
	useModelSwitch := cast.ToInt(req.UseModelSwitch)
	aiSummary := cast.ToInt(req.AiSummary)
	aiSummaryModel := cast.ToString(req.AiSummaryModel)
	summaryModelConfigId := cast.ToInt(req.SummaryModelConfigId)
	shareUrl := strings.TrimSpace(req.ShareUrl)
	accessRights := cast.ToInt(req.AccessRights)
	statisticsSet := strings.TrimSpace(req.StatisticsSet)
	typ := strings.TrimSpace(req.Type)
	chunkType := cast.ToInt(req.ChunkType)
	normalChunkDefaultSeparatorsNo := cast.ToString(req.NormalChunkDefaultSeparatorsNo)
	normalChunkDefaultChunkSize := cast.ToInt(req.NormalChunkDefaultChunkSize)
	normalChunkDefaultChunkOverlap := cast.ToInt(req.NormalChunkDefaultChunkOverlap)
	normalChunkDefaultNotMergedText := cast.ToBool(req.NormalChunkDefaultNotMergedText)
	semanticChunkDefaultChunkSize := cast.ToInt(req.SemanticChunkDefaultChunkSize)
	semanticChunkDefaultChunkOverlap := cast.ToInt(req.SemanticChunkDefaultChunkOverlap)
	semanticChunkDefaultThreshold := cast.ToInt(req.SemanticChunkDefaultThreshold)
	AiChunkPrumpt := cast.ToString(req.AiChunkPrumpt)
	AiChunkModel := strings.TrimSpace(req.AiChunkModel)
	AiChunkModelConfigId := cast.ToInt(req.AiChunkModelConfigId)
	AiChunkSize := cast.ToInt(req.AiChunkSize)
	qaIndexType := cast.ToInt(req.QaIndexType)
	iconTemplateConfigId := cast.ToInt(req.IconTemplateConfigId)
	fatherChunkParagraphType := cast.ToInt(req.FatherChunkParagraphType)
	fatherChunkSeparatorsNo := strings.TrimSpace(req.FatherChunkSeparatorsNo)
	fatherChunkChunkSize := cast.ToInt(req.FatherChunkChunkSize)
	sonChunkSeparatorsNo := strings.TrimSpace(req.SonChunkSeparatorsNo)
	sonChunkChunkSize := cast.ToInt(req.SonChunkChunkSize)
	if id <= 0 || len(libraryName) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	if cast.ToInt(typ) != define.OpenLibraryType {
		//check model_config_id and use_model
		config, err := common.GetModelConfigInfo(modelConfigId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(useModel, modelInfo.VectorModelList) && !common.IsMultiConfModel(config["model_define"]) {
			return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `use_model`))
		}
		if len(config) == 0 || !tool.InArrayString(common.TextEmbedding, strings.Split(config[`model_types`], `,`)) {
			return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `model_config_id`))
		}
	}
	if useModelSwitch == define.SwitchOn && (modelConfigId == 0 || useModel == "") {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `use_model`))
	}
	if summaryModelConfigId > 0 {
		summaryConfig, err := common.GetModelConfigInfo(summaryModelConfigId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
		}
		modelInfo, _ := common.GetModelInfoByDefine(summaryConfig[`model_define`])
		if !tool.InArrayString(aiSummaryModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(summaryConfig["model_define"]) {
			return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `ai_summary_model`))
		}
	}
	if graphSwitch == 1 && graphModelConfigId == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `graph_model_config_id`))
	}
	if graphSwitch == 1 && len(graphUseModel) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `graph_use_model`))
	}
	chunkParam := define.ChunkParam{
		ChunkType:                        req.ChunkType,
		NormalChunkDefaultSeparatorsNo:   req.NormalChunkDefaultSeparatorsNo,
		NormalChunkDefaultChunkSize:      req.NormalChunkDefaultChunkSize,
		NormalChunkDefaultChunkOverlap:   req.NormalChunkDefaultChunkOverlap,
		NormalChunkDefaultNotMergedText:  req.NormalChunkDefaultNotMergedText,
		SemanticChunkDefaultChunkSize:    req.SemanticChunkDefaultChunkSize,
		SemanticChunkDefaultChunkOverlap: req.SemanticChunkDefaultChunkOverlap,
		SemanticChunkDefaultThreshold:    req.SemanticChunkDefaultThreshold,
		AiChunkPrumpt:                    AiChunkPrumpt,
		AiChunkModel:                     AiChunkModel,
		AiChunkModelConfigId:             req.AiChunkModelConfigId,
		AiChunkSize:                      req.AiChunkSize,
		QaIndexType:                      req.QaIndexType,
		FatherChunkParagraphType:         req.FatherChunkParagraphType,
		FatherChunkSeparatorsNo:          req.FatherChunkSeparatorsNo,
		FatherChunkChunkSize:             req.FatherChunkChunkSize,
		SonChunkSeparatorsNo:             req.SonChunkSeparatorsNo,
		SonChunkChunkSize:                req.SonChunkChunkSize,
	}
	err := ValidateChunkParam(adminUserId, &chunkParam, typ, lang)
	if err != nil {
		return nil, -1, err
	}
	info, err := common.GetLibraryInfo(id, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	if chunkType == define.ChunkTypeFatherSon && cast.ToInt(info[`type`]) != define.GeneralLibraryType { //父子分段仅支持普通知识库
		return nil, -1, errors.New(i18n.Show(lang, `param_invalid`, `chunk_type`))
	}

	//headImg uploaded
	avatar := ""
	uploadInfo, err := common.SaveUploadedFile(req.FileAvatar, define.ImageLimitSize, adminUserId, `library_avatar`, define.ImageAllowExt)
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
		`normal_chunk_default_not_merged_text`: normalChunkDefaultNotMergedText,
		`normal_chunk_default_chunk_overlap`:   normalChunkDefaultChunkOverlap,
		`semantic_chunk_default_chunk_size`:    semanticChunkDefaultChunkSize,
		`semantic_chunk_default_chunk_overlap`: semanticChunkDefaultChunkOverlap,
		`semantic_chunk_default_threshold`:     semanticChunkDefaultThreshold,
		`ai_chunk_prumpt`:                      AiChunkPrumpt,
		`ai_chunk_model`:                       AiChunkModel,
		`ai_chunk_model_config_id`:             AiChunkModelConfigId,
		`ai_chunk_size`:                        AiChunkSize,
		`qa_index_type`:                        qaIndexType,
		`icon_template_config_id`:              iconTemplateConfigId,
		`father_chunk_paragraph_type`:          fatherChunkParagraphType,
		`father_chunk_separators_no`:           fatherChunkSeparatorsNo,
		`father_chunk_chunk_size`:              fatherChunkChunkSize,
		`son_chunk_separators_no`:              sonChunkSeparatorsNo,
		`son_chunk_chunk_size`:                 sonChunkChunkSize,
	}
	if len(avatar) > 0 {
		data[`avatar`] = avatar
	}
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
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
	// QA 切换索引方式
	if cast.ToInt(info[`type`]) == define.QALibraryType && qaIndexType != cast.ToInt(info[`qa_index_type`]) {
		go common.EmbeddingNewQAVector(id, cast.ToInt(info[`admin_user_id`]), qaIndexType)
	}

	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.LibraryCacheBuildHandler{LibraryId: id})
	return nil, 0, nil
}

func ValidateChunkParam(adminUserId int, chunkParam *define.ChunkParam, typ string, lang string) error {
	if cast.ToInt(typ) == define.QALibraryType || cast.ToInt(typ) == define.OpenLibraryType {
		return nil
	}
	chunkType := cast.ToInt(chunkParam.ChunkType)
	allowChunkTypes := []int{define.ChunkTypeNormal, define.ChunkTypeSemantic, define.ChunkTypeAi, define.ChunkTypeFatherSon}
	if !tool.InArrayInt(chunkType, allowChunkTypes) && cast.ToInt(typ) != define.OpenLibraryType {
		return errors.New(i18n.Show(lang, `param_invalid`, `chunk_type`))
	}
	if chunkType == define.ChunkTypeFatherSon {
		if cast.ToInt(typ) != define.GeneralLibraryType { //父子分段仅支持普通知识库
			return errors.New(i18n.Show(lang, `param_invalid`, `chunk_type`))
		}
		if !tool.InArrayInt(cast.ToInt(chunkParam.FatherChunkParagraphType), []int{define.FatherChunkParagraphTypeFullText, define.FatherChunkParagraphTypeSection}) {
			return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_paragraph_type`))
		}
		if cast.ToInt(chunkParam.FatherChunkParagraphType) != define.FatherChunkParagraphTypeFullText {
			if len(chunkParam.FatherChunkSeparatorsNo) == 0 {
				return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_separators_no`))
			}
			if cast.ToInt(chunkParam.FatherChunkChunkSize) < 0 {
				return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_chunk_size`))
			}
		}
		if len(chunkParam.SonChunkSeparatorsNo) == 0 {
			return errors.New(i18n.Show(lang, `param_invalid`, `son_chunk_separators_no`))
		}
		if cast.ToInt(chunkParam.SonChunkChunkSize) < 0 {
			return errors.New(i18n.Show(lang, `param_invalid`, `son_chunk_chunk_size`))
		}
	}
	if chunkType == define.ChunkTypeNormal {
		if len(chunkParam.NormalChunkDefaultSeparatorsNo) == 0 {
			return errors.New(i18n.Show(lang, `param_invalid`, `separators_no`))
		}
		if cast.ToInt(chunkParam.NormalChunkDefaultChunkSize) < define.SplitChunkMinSize ||
			cast.ToInt(chunkParam.NormalChunkDefaultChunkSize) > define.SplitChunkMaxSize {
			return errors.New(i18n.Show(lang, `chunk_size_err`, define.SplitChunkMinSize, define.SplitChunkMaxSize))
		}
		maxChunkOverlap := cast.ToInt(chunkParam.NormalChunkDefaultChunkSize) / 2
		if cast.ToInt(chunkParam.NormalChunkDefaultChunkOverlap) < 0 || cast.ToInt(chunkParam.NormalChunkDefaultChunkOverlap) > maxChunkOverlap {
			return errors.New(i18n.Show(lang, `chunk_overlap_err`, 0, maxChunkOverlap))
		}
	}
	if chunkType == define.ChunkTypeSemantic {
		if cast.ToInt(chunkParam.SemanticChunkDefaultChunkSize) < define.SplitChunkMinSize ||
			cast.ToInt(chunkParam.SemanticChunkDefaultChunkSize) > define.SplitChunkMaxSize {
			return errors.New(i18n.Show(lang, `semantic_chunk_size_err`, define.SplitChunkMinSize, define.SplitChunkMaxSize))
		}
		maxSemanticChunkOverlap := cast.ToInt(chunkParam.SemanticChunkDefaultChunkSize) / 2
		if cast.ToInt(chunkParam.SemanticChunkDefaultChunkOverlap) > maxSemanticChunkOverlap {
			return errors.New(i18n.Show(lang, `semantic_chunk_overlap_err`, 0, maxSemanticChunkOverlap))
		}
		if cast.ToInt(chunkParam.SemanticChunkDefaultThreshold) < 1 || cast.ToInt(chunkParam.SemanticChunkDefaultThreshold) > 100 {
			return errors.New(i18n.Show(lang, `semantic_chunk_threshold_err`, 1, 100))
		}
	}
	if chunkType == define.ChunkTypeAi {
		if ok := common.CheckModelIsValid(adminUserId, cast.ToInt(chunkParam.AiChunkModelConfigId), chunkParam.AiChunkModel, common.Llm); !ok {
			return errors.New(i18n.Show(lang, `param_invalid`, `ai_chunk_model`))
		}
		if len(chunkParam.AiChunkPrumpt) == 0 || len(chunkParam.AiChunkPrumpt) > 500 {
			return errors.New(i18n.Show(lang, `param_invalid`, `ai_chunk_prumpt`))
		}
	}
	if chunkType == define.ChunkTypeFatherSon {
		if !tool.InArrayInt(cast.ToInt(chunkParam.FatherChunkParagraphType), []int{define.FatherChunkParagraphTypeFullText, define.FatherChunkParagraphTypeSection}) {
			return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_paragraph_type`))
		}
		if cast.ToInt(chunkParam.FatherChunkParagraphType) != define.FatherChunkParagraphTypeFullText {
			if len(chunkParam.FatherChunkSeparatorsNo) == 0 {
				return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_separators_no`))
			}
			if cast.ToInt(chunkParam.FatherChunkChunkSize) < 0 {
				return errors.New(i18n.Show(lang, `param_invalid`, `father_chunk_chunk_size`))
			}
		}
		if len(chunkParam.SonChunkSeparatorsNo) == 0 {
			return errors.New(i18n.Show(lang, `param_invalid`, `son_chunk_separators_no`))
		}
		if cast.ToInt(chunkParam.SonChunkChunkSize) < 0 {
			return errors.New(i18n.Show(lang, `param_invalid`, `son_chunk_chunk_size`))
		}
	}
	return nil
}

type BridgeDeleteLibraryReq struct {
	Id string `form:"id"`
}

func BridgeDeleteLibrary(adminUserId, loginUserId int, lang string, req *BridgeDeleteLibraryReq) (map[string]any, int, error) {
	id := cast.ToInt(req.Id)
	if id <= 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	info, err := common.GetLibraryInfo(id, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `no_data`))
	}
	// check robot relation
	robotdata, err := common.GetLibraryRobotInfo(adminUserId, id)
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(robotdata) > 0 {
		return nil, -1, errors.New(i18n.Show(lang, `relation_robot`))
	}
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
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
	if common.GetNeo4jStatus(adminUserId) {
		err = common.NewGraphDB(adminUserId).DeleteByLibrary(id)
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
	return nil, 0, nil
}
