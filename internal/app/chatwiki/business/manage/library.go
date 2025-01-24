// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
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
		managedRobotIdList := GetUserManagedData(userId, `managed_library_list`)
		m.Where(`id`, `in`, strings.Join(managedRobotIdList, `,`))
	}

	list, err := m.Field(`id,library_name,library_intro`).Order(`id desc`).Select()
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
	for _, params := range list {
		libraryIds = append(libraryIds, params[`id`])
		params[`file_total`], params[`file_size`] = `0`, `0`
	}
	stats, err := msql.Model(`chat_ai_library_file`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, `in`, strings.Join(libraryIds, `,`)).Group(`library_id`).
		ColumnMap(`COUNT(1) as file_total,SUM(file_size) as file_size`, `library_id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for _, params := range list {
		if len(stats[params[`id`]]) == 0 {
			continue
		}
		params[`file_total`] = stats[params[`id`]][`file_total`]
		params[`file_size`] = stats[params[`id`]][`file_size`]
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
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
	for k, v := range info {
		data[k] = v
	}
	data[`is_offline`] = false
	for _, config := range common.GetModelList() {
		if info[`model_define`] == config.ModelDefine && config.IsOffline {
			data[`is_offline`] = true
		}
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
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	avatar := ""
	if len(libraryName) == 0 || modelConfigId <= 0 || len(useModel) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//check model_config_id and use_model
	config, err := common.GetModelConfigInfo(modelConfigId, userId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(config) == 0 || !tool.InArrayString(common.TextEmbedding, strings.Split(config[`model_types`], `,`)) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_config_id`))))
		return
	}
	modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
	if !tool.InArrayString(useModel, modelInfo.VectorModelList) && !common.IsMultiConfModel(config["model_define"]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
		return
	}
	//headImg uploaded
	fileAvatar, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileAvatar, define.ImageLimitSize, userId, `library_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}
	//database dispose
	data := msql.Datas{
		`admin_user_id`:   userId,
		`library_name`:    libraryName,
		`library_intro`:   libraryIntro,
		`model_config_id`: modelConfigId,
		`use_model`:       useModel,
		`create_time`:     tool.Time2Int(),
		`update_time`:     tool.Time2Int(),
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
	//common save
	fileIds, err := addLibFile(c, userId, int(libraryId))
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	_ = AddUserMangedData(getLoginUserId(c), `managed_library_list`, libraryId)

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`id`: libraryId, `file_ids`: fileIds}, nil))
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
	_, err = msql.Model(`chat_ai_library_file_data`, define.Postgres).Where(`library_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`library_id`, cast.ToString(id)).Delete()
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
	if id <= 0 || len(libraryName) == 0 {
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
	_, err = msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`library_name`:  libraryName,
		`library_intro`: libraryIntro,
		`update_time`:   tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
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
	libraryId := cast.ToInt(c.PostForm(`id`))
	question := strings.TrimSpace(c.PostForm(`question`))
	size := cast.ToInt(c.PostForm(`size`))
	similarity := cast.ToFloat64(c.PostForm(`similarity`))
	searchType := cast.ToInt(c.PostForm(`search_type`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	if libraryId <= 0 || len(question) == 0 || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if searchType != define.SearchTypeMixed && searchType != define.SearchTypeVector && searchType != define.SearchTypeFullText {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	info, err := common.GetLibraryInfo(libraryId, userId)
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
	robot := msql.Params{}
	if rerankModelConfigID > 0 {
		robot[`rerank_status`] = cast.ToString(1)
		robot[`rerank_model_config_id`] = cast.ToString(rerankModelConfigID)
		robot[`robot_name`] = robotName
	}

	list, err := common.GetMatchLibraryParagraphList("", "", question, []string{}, cast.ToString(libraryId), size, similarity, searchType, robot)
	c.String(http.StatusOK, lib_web.FmtJson(list, err))
}
