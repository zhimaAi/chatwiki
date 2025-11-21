// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/middlewares"
	"errors"
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
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeLibraryListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	list, httpStatus, err := BridgeGetLibraryList(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
	return
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
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeCreateLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		logs.Error(`bind error %s`, err.Error())
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	logs.Debug(`CreateLibrary req %#v`, req)
	req.FileAvatar, _ = c.FormFile(`avatar`)
	list, httpStatus, err := BridgeCreateLibrary(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func DeleteLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeDeleteLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	list, httpStatus, err := BridgeDeleteLibrary(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func EditLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeEditLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.FileAvatar, _ = c.FormFile(`avatar`)
	list, httpStatus, err := BridgeEditLibrary(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
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
		`recall_type`:   recallType,
		`admin_user_id`: cast.ToString(userId),
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
