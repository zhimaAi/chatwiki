// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
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

func OpenGetLibraryList(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//models
	reqM := manage.BridgeGetModelConfigOptionReq{
		ModelType: common.TextEmbedding,
	}
	models, httpStatus, err := manage.BridgeGetModelConfigOption(adminUserId, adminUserId, common.GetLang(c), &reqM)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, models, httpStatus, err)
		return
	}
	//library list
	reqL := manage.BridgeLibraryListReq{}
	err = common.RequestParamsBind(&reqL, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(reqL, err, common.GetLang(c)).Error())
		return
	}
	if reqL.Type != `` && !tool.InArrayInt(cast.ToInt(reqL.Type), define.LibraryTypes[:2]) {
		common.FmtError(c, `param_err`, `type`)
		return
	}
	libraryList, httpStatus, err := manage.BridgeGetLibraryList(adminUserId, adminUserId, common.GetLang(c), &reqL)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libraryList, httpStatus, err)
		return
	}
	//library group list
	reqG := manage.BridgeLibraryListGroupReq{}
	reqG.Type = reqL.Type
	libraryGroupList, httpStatus, err := manage.BridgeGetLibraryListGroup(adminUserId, adminUserId, common.GetLang(c), &reqG)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libraryList, httpStatus, err)
		return
	}
	//fill model info
	returnData := make([]map[string]any, 0)
	for _, library := range libraryList {
		library[`model_define`] = ``
		for _, model := range models {
			modelConfig, ok := model[`model_config`].(msql.Params)
			if !ok {
				break
			}
			if cast.ToInt(library[`model_config_id`]) == cast.ToInt(modelConfig[`id`]) {
				library[`model_define`] = modelConfig[`model_define`]
			}
		}
		//take params
		returnData = append(returnData, common.TakeParamsFromMap(library, `id`, `avatar`, `type`, `library_intro`, `library_name`,
			`file_size`, `file_total`, `model_config_id`, `model_define`, `use_model`, `group_id`))
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`library_list`:       returnData,
		`library_group_list`: libraryGroupList,
	}, nil))
}

func OpenGetModelConfigOption(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	req := manage.BridgeGetModelConfigOptionReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.ModelType = strings.TrimSpace(c.Query(`model_type`))
	list, httpStatus, err := manage.BridgeGetModelConfigOption(adminUserId, adminUserId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}
	returnData := make([]map[string]any, 0)
	for _, modelInfo := range list {
		configM := common.TakeParamsFromMap(modelInfo[`model_config`], `id`, `model_define`)
		infoM := common.TakeParamsFromMap(modelInfo[`model_info`], `vector_model_list`, `llm_model_list`)
		returnData = append(returnData, map[string]any{
			`model_config`: configM,
			`model_info`:   infoM,
		})
	}
	common.FmtBridgeResponse(c, returnData, httpStatus, err)
}

func OpenCreateLibraryGeneral(c *gin.Context) {
	openCreateLibrary(c, cast.ToString(define.GeneralLibraryType))
}

func OpenCreateLibraryQA(c *gin.Context) {
	openCreateLibrary(c, cast.ToString(define.QALibraryType))
}

func openCreateLibrary(c *gin.Context, typ string) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	req := manage.BridgeCreateLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.Type = typ
	req.FileAvatar, _ = c.FormFile(`avatar`)
	data, httpStatus, err := manage.BridgeCreateLibrary(adminUserId, adminUserId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, data, httpStatus, err)
		return
	}
	returnData := common.TakeParamsFromMap(data, `id`)
	common.FmtBridgeResponse(c, returnData, httpStatus, err)
}

func OpenEditLibrary(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	req := manage.BridgeEditLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.FileAvatar, _ = c.FormFile(`avatar`)
	req.Type = cast.ToString(define.GeneralLibraryType)
	list, httpStatus, err := manage.BridgeEditLibrary(c, adminUserId, adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func OpenEditLibraryQA(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	req := manage.BridgeEditLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.FileAvatar, _ = c.FormFile(`avatar`)
	req.Type = cast.ToString(define.QALibraryType)
	list, httpStatus, err := manage.BridgeEditLibrary(c, adminUserId, adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func OpenGetSeparatorsList(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	list, httpStatus, err := manage.BridgeGetSeparatorsList(adminUserId, adminUserId, common.GetLang(c))
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func OpenDeleteLibrary(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	req := manage.BridgeDeleteLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	data, httpStatus, err := manage.BridgeDeleteLibrary(adminUserId, adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, data, httpStatus, err)
}

func OpenGetLibFileList(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//lib file list
	req := manage.BridgeGetLibFileListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)
	libFileData, httpStatus, err := manage.BridgeGetLibFileList(adminUserId, adminUserId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libFileData, httpStatus, err)
		return
	}
	libraryInfo, ok := libFileData[`info`].(msql.Params)
	if !ok {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//lib group list
	reqG := manage.BridgeGetLibraryGroupReq{}
	reqG.LibraryId = req.LibraryId
	if cast.ToInt(libraryInfo[`type`]) == define.GeneralLibraryType {
		reqG.GroupType = cast.ToString(define.LibraryGroupTypeFile)
	} else if cast.ToInt(libraryInfo[`type`]) == define.QALibraryType {
		reqG.GroupType = cast.ToString(define.LibraryGroupTypeQA)
	}
	libGroupData, httpStatus, err := manage.BridgeGetLibraryGroup(adminUserId, adminUserId, common.GetLang(c), &reqG)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libGroupData, httpStatus, err)
		return
	}
	data := make(map[string]any)
	returnList := make([]map[string]any, 0)
	libFileList, ok := libFileData[`list`].([]msql.Params)
	if ok {
		for _, libFile := range libFileList {
			returnList = append(returnList, common.TakeParamsFromMap(libFile, `id`, `file_name`, `file_ext`, `file_size`, `total_hits`,
				`today_hits`, `yesterday_hits`, `status`, `update_time`))
		}
	}
	data[`lib_file_list`] = returnList
	data[`library_group_list`] = libGroupData
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func OpenGetParagraphList(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//paragraph list
	req := manage.BridgeGetParagraphListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.Status = c.DefaultQuery(`status`, `-1`)
	req.GraphStatus = c.DefaultQuery(`graph_status`, `-1`)
	req.CategoryId = c.DefaultQuery(`category_id`, `-1`)
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)
	libParagraphList, httpStatus, err := manage.BridgeGetParagraphList(adminUserId, adminUserId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libParagraphList, httpStatus, err)
		return
	}
	libraryInfo, ok := libParagraphList[`info`].(msql.Params)
	if !ok {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//lib group list
	reqG := manage.BridgeGetLibraryGroupReq{}
	reqG.LibraryId = req.LibraryId
	reqG.GroupType = cast.ToString(libraryInfo[`type`])
	libGroupData, httpStatus, err := manage.BridgeGetLibraryGroup(adminUserId, adminUserId, common.GetLang(c), &reqG)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libParagraphList, httpStatus, err)
		return
	}
	data := make(map[string]any)
	returnList := make([]map[string]any, 0)
	libFileList, ok := libParagraphList[`list`].([]map[string]any)
	if ok {
		for _, libFile := range libFileList {
			returnList = append(returnList, common.TakeParamsFromMap(libFile, `id`, `question`, `answer`, `total_hits`,
				`today_hits`, `yesterday_hits`))
		}
	}
	data[`lib_paragraph_list`] = returnList
	data[`library_group_list`] = libGroupData
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func OpenGetLibraryGeneralParagraphList(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//paragraph list
	req := manage.BridgeGetParagraphListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.Status = c.DefaultQuery(`status`, `-1`)
	req.GraphStatus = c.DefaultQuery(`graph_status`, `-1`)
	req.CategoryId = c.DefaultQuery(`category_id`, `-1`)
	req.GroupId = c.DefaultQuery(`group_id`, `-1`)
	libParagraphList, httpStatus, err := manage.BridgeGetParagraphList(adminUserId, adminUserId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, libParagraphList, httpStatus, err)
		return
	}
	returnList := make([]map[string]any, 0)
	libFileList, ok := libParagraphList[`list`].([]map[string]any)
	if ok {
		for _, libFile := range libFileList {
			returnList = append(returnList, common.TakeParamsFromMap(libFile, `id`, `title`, `content`, `status`))
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(returnList, nil))
}

func OpenAddLibraryFileQA(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//add file params
	req := manage.BridgeAddLibraryFileReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.IsQaDoc = cast.ToString(define.DocTypeQa)
	req.DocType = c.DefaultPostForm(`doc_type`, cast.ToString(define.DocTypeLocal))
	//chunk params
	chunkParam := define.ChunkParam{}
	err = common.RequestParamsBind(&chunkParam, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	openAddLibraryFile(c, &req, &chunkParam)
}

func OpenAddLibraryFileGeneralLocal(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//add file params
	req := manage.BridgeAddLibraryFileReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.DocType = c.DefaultPostForm(`doc_type`, cast.ToString(define.DocTypeLocal))
	//chunk params
	chunkParam := define.ChunkParam{}
	err = common.RequestParamsBind(&chunkParam, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	openAddLibraryFile(c, &req, &chunkParam)
}

func openAddLibraryFile(c *gin.Context, req *manage.BridgeAddLibraryFileReq, chunkParam *define.ChunkParam) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	data, httpStatus, err := manage.BridgeAddLibraryFile(adminUserId, adminUserId, common.GetLang(c), req, chunkParam, c)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, data, httpStatus, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

type BatchAddLibraryQaItem struct {
	Question         string `json:"question"`
	Answer           string `json:"answer"`
	SimilarQuestions string `json:"similar_questions"`
}

func OpenBatchAddLibraryQa(c *gin.Context) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//build token
	user := common.GetUserInfo(adminUserId)
	if user == nil {
		user = msql.Params{}
	}
	jwtParams, err := common.GetToken(adminUserId, user[`user_name`], 0)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//batch save
	batchAddLibraryQaItems := make([]BatchAddLibraryQaItem, 0)
	err = tool.JsonDecode(c.PostForm(`qa_lists`), &batchAddLibraryQaItems)
	if err != nil {
		common.FmtError(c, `param_err`, `qa_lists`)
		return
	}
	resultList := make([]string, 0)
	for _, item := range batchAddLibraryQaItems {
		req := manage.BridgeSaveParagraphReq{}
		req.LibraryId = c.PostForm(`library_id`)
		req.GroupId = c.PostForm(`group_id`)
		req.Question = item.Question
		req.Answer = item.Answer
		req.SimilarQuestions = item.SimilarQuestions
		req.Token = cast.ToString(jwtParams[`token`])
		_, httpStatus, err := manage.BridgeSaveParagraph(adminUserId, adminUserId, common.GetLang(c), &req)
		if httpStatus != 0 {
			resultList = append(resultList, err.Error())
		} else {
			resultList = append(resultList, `success`)
		}
	}
	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		`result_list`: resultList,
	}, nil))
}

func OpenEditGeneralParagraph(c *gin.Context) {
	//bind post
	req := manage.BridgeSaveParagraphReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if cast.ToInt(req.Id) == 0 {
		common.FmtError(c, `param_err`, `id`)
		return
	}
	openEditParagraph(c, req)
}

func OpenEditQAParagraph(c *gin.Context) {
	//bind post
	req := manage.BridgeSaveParagraphReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if cast.ToInt(req.Id) == 0 {
		common.FmtError(c, `param_err`, `id`)
		return
	}
	openEditParagraph(c, req)
}

func openEditParagraph(c *gin.Context, req manage.BridgeSaveParagraphReq) {
	adminUserId := parseToken(c)
	if adminUserId == 0 {
		return
	}
	//build token
	user := common.GetUserInfo(adminUserId)
	if len(user) == 0 {
		common.FmtError(c, `system_error`)
		return
	}
	jwtParams, err := common.GetToken(adminUserId, user[`user_name`], 0)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	req.Token = cast.ToString(jwtParams[`token`])
	_, httpStatus, err := manage.BridgeSaveParagraph(adminUserId, adminUserId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, nil, httpStatus, err)
}

func parseToken(c *gin.Context) int {
	headers, err := common.ParseAuthorizationToken(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return 0
	}
	if !common.CheckRobotKey(headers[`robot_key`]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return 0
	}
	robot, err := common.GetRobotInfo(headers[`robot_key`])
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return 0
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return 0
	}
	return cast.ToInt(robot[`admin_user_id`])
}
