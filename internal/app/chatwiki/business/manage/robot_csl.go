// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func RobotExport(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	robotCsl, err := CreateRobotCsl(id, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	content, err := robotCsl.Output()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	temFile := `static/public/download/` + tool.MD5(content) + `.csl`
	if err = tool.WriteFile(temFile, content); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.FileAttachment(temFile, robotCsl.FileName)
	_ = os.Remove(temFile) //delete temp file
}

func RobotImport(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	fileHeader, err := c.FormFile(`file`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(fileHeader.Filename), `.`))
	if ext != `csl` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf(`不支持的文件格式:%s`, ext)))
		return
	}
	reader, err := fileHeader.Open()
	defer func(reader multipart.File) {
		_ = reader.Close()
	}(reader)
	bs, err := io.ReadAll(reader)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(bs) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`文件内容不能为空`)))
		return
	}
	robotCsl, err := common.ParseRobotCsl(string(bs))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	c.String(http.StatusOK, lib_web.FmtJson(ApplyRobotCsl(adminUserId, getLoginUserId(c), token, robotCsl)))
}

func CreateRobotCsl(id, adminUserId int, simple ...bool) (robotCsl *common.RobotCsl, err error) {
	robotCsl = common.NewRobotCsl()
	if id <= 0 {
		err = errors.New(`机器人ID参数错误`)
		return
	}
	if adminUserId <= 0 {
		err = errors.New(`管理员ID参数错误`)
		return
	}
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		return
	}
	if len(robot) == 0 {
		err = errors.New(`机器人信息不存在`)
		return
	}
	robotCsl.Robot = robot
	robotName := regexp.MustCompile(`[^a-zA-Z0-9\p{Han}]`).ReplaceAllString(robot[`robot_name`], ``)
	if len(robotName) > 0 {
		robotCsl.FileName = fmt.Sprintf(`%s_%s.csl`, robotName, tool.Date(`YmdHis`))
	}
	//文件分段精选配置
	category, err := msql.Model(`category`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).ColumnObj(`type`, `id`)
	if err != nil {
		return
	}
	robotCsl.Category = category
	//工作流,知识库,数据表引用采集
	switch cast.ToInt(robot[`application_type`]) {
	case define.ApplicationTypeChat:
		if len(robot[`work_flow_ids`]) > 0 {
			for _, workFlowId := range strings.Split(robot[`work_flow_ids`], `,`) {
				children, subErr := CreateRobotCsl(cast.ToInt(workFlowId), adminUserId, true)
				if subErr == nil {
					for libraryId := range children.Librarys {
						robotCsl.Librarys[libraryId] = nil
					}
					for formId := range children.Forms {
						robotCsl.Forms[formId] = nil
					}
					children.Librarys = nil
					children.Forms = nil
					children.Workflows = nil
					robotCsl.Workflows = append(robotCsl.Workflows, children)
				} else {
					logs.Error(subErr.Error())
				}
			}
		}
		if len(robot[`library_ids`]) > 0 {
			for _, libraryId := range strings.Split(robot[`library_ids`], `,`) {
				robotCsl.Librarys[cast.ToInt(libraryId)] = nil
			}
		}
		if len(robot[`form_ids`]) > 0 {
			for _, formId := range strings.Split(robot[`form_ids`], `,`) {
				robotCsl.Forms[cast.ToInt(formId)] = nil
			}
		}
	case define.ApplicationTypeFlow:
		nodes, sqlErr := msql.Model(`work_flow_node`, define.Postgres).Where(`robot_id`, cast.ToString(id)).
			Where(`data_type`, cast.ToString(define.DataTypeDraft)).Select()
		if sqlErr != nil {
			err = sqlErr
			return
		}
		robotCsl.Nodes = nodes
		for _, node := range nodes {
			nodeParams := work_flow.NodeParams{}
			_ = tool.JsonDecodeUseNumber(node[`node_params`], &nodeParams)
			switch cast.ToInt(node[`node_type`]) {
			case work_flow.NodeTypeLibs:
				if len(nodeParams.Libs.LibraryIds) > 0 { //存在关联的知识库
					for _, libraryId := range strings.Split(nodeParams.Libs.LibraryIds, `,`) {
						robotCsl.Librarys[cast.ToInt(libraryId)] = nil
					}
				}
			case work_flow.NodeTypeFormInsert:
				if nodeParams.FormInsert.FormId > 0 {
					robotCsl.Forms[nodeParams.FormInsert.FormId.Int()] = nil
				}
			case work_flow.NodeTypeFormDelete:
				if nodeParams.FormDelete.FormId > 0 {
					robotCsl.Forms[nodeParams.FormDelete.FormId.Int()] = nil
				}
			case work_flow.NodeTypeFormUpdate:
				if nodeParams.FormUpdate.FormId > 0 {
					robotCsl.Forms[nodeParams.FormUpdate.FormId.Int()] = nil
				}
			case work_flow.NodeTypeFormSelect:
				if nodeParams.FormSelect.FormId > 0 {
					robotCsl.Forms[nodeParams.FormSelect.FormId.Int()] = nil
				}
			}
		}
	}
	if len(simple) > 0 && simple[0] {
		return
	}
	//处理知识库
	for libraryId := range robotCsl.Librarys {
		libraryCsl, sqlErr := common.BuildLibraryCsl(libraryId, adminUserId)
		if sqlErr == nil {
			robotCsl.Librarys[libraryId] = libraryCsl
		} else {
			logs.Error(sqlErr.Error())
		}
	}
	//处理表单
	for formId := range robotCsl.Forms {
		formCsl, sqlErr := common.BuildFormCsl(formId, adminUserId)
		if sqlErr == nil {
			robotCsl.Forms[formId] = formCsl
		} else {
			logs.Error(sqlErr.Error())
		}
	}
	return
}

func ApplyRobotCsl(adminUserId, userId int, token string, robotCsl *common.RobotCsl) (any, error) {
	if robotCsl == nil || len(robotCsl.Robot) == 0 {
		return nil, errors.New(`机器人信息不存在`)
	}
	//获取默认的模型配置
	models, err := common.GetDefaultModelParams(adminUserId)
	if err != nil {
		return nil, err
	}
	logs.Debug(`models:%s`, tool.JsonEncodeNoError(models))
	//初始化新旧id的maps
	cslIdMaps := common.NewCslIdMaps()
	//分段分类(精选)id(旧=>新)
	category, err := msql.Model(`category`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).ColumnObj(`id`, `type`)
	if err != nil {
		return nil, err
	}
	for oldId, oldType := range robotCsl.Category {
		cslIdMaps.Category[cast.ToInt(oldId)] = cast.ToInt(category[oldType])
	}
	//开始导入知识库
	for _, libraryCsl := range robotCsl.Librarys {
		sqlErr := libraryCsl.Import(adminUserId, userId, cslIdMaps, models, token)
		if sqlErr != nil {
			logs.Error(sqlErr.Error())
		}
	}
	//开始导入数据表
	for _, formCsl := range robotCsl.Forms {
		sqlErr := formCsl.Import(adminUserId, cslIdMaps)
		if sqlErr != nil {
			logs.Error(sqlErr.Error())
		}
	}
	//开始导入机器人
	var newRobot msql.Params
	switch cast.ToInt(robotCsl.Robot[`application_type`]) {
	case define.ApplicationTypeChat:
		newRobot, err = ApplyChatRobot(robotCsl.Robot, cslIdMaps, models, token)
		if err != nil {
			return nil, err
		}
		//导入关联工作流
		workFlowIds := make([]string, 0)
		for _, workflow := range robotCsl.Workflows {
			relationWorkflow, sqlErr := ApplyFlowRobot(adminUserId, workflow.Robot, workflow.Nodes, cslIdMaps, models, token)
			if sqlErr != nil {
				logs.Error(sqlErr.Error())
				continue
			}
			workFlowIds = append(workFlowIds, relationWorkflow[`id`])
		}
		//关联工作流操作
		if len(workFlowIds) > 0 {
			up := msql.Datas{`work_flow_ids`: strings.Join(workFlowIds, `,`), `update_time`: tool.Time2Int()}
			if _, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, newRobot[`id`]).Update(up); err != nil {
				return nil, fmt.Errorf(`应用关联工作流操作失败:%s`, err.Error())
			}
			lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: newRobot[`robot_key`]})
		}
	case define.ApplicationTypeFlow:
		newRobot, err = ApplyFlowRobot(adminUserId, robotCsl.Robot, robotCsl.Nodes, cslIdMaps, models, token)
		if err != nil {
			return nil, err
		}
	}
	//返回机器人信息
	return newRobot, nil
}

func ApplyChatRobot(robot msql.Params, cslIdMaps *common.CslIdMaps, models *common.DefaultModelParams, token string) (msql.Params, error) {
	chatRobot := make(map[string]string)
	for key, val := range robot {
		switch key {
		case `id`, `admin_user_id`, `robot_key`, `application_type`, `work_flow_ids`, `default_library_id`, `default_app_config`, `start_node_key`, `work_flow_model_config_ids`, `creator`:
		case `library_ids`:
			newVals := make([]string, 0)
			if len(val) > 0 {
				for _, oldLibraryId := range strings.Split(val, `,`) {
					if newLibraryId := cslIdMaps.Librarys[cast.ToInt(oldLibraryId)]; newLibraryId > 0 {
						newVals = append(newVals, cast.ToString(newLibraryId))
					}
				}
			}
			chatRobot[key] = strings.Join(newVals, `,`)
		case `form_ids`:
			newVals := make([]string, 0)
			if len(val) > 0 {
				for _, oldFormId := range strings.Split(val, `,`) {
					if newFormId := cslIdMaps.Forms[cast.ToInt(oldFormId)]; newFormId > 0 {
						newVals = append(newVals, cast.ToString(newFormId))
					}
				}
			}
			chatRobot[key] = strings.Join(newVals, `,`)
		case `model_config_id`, `optimize_question_model_config_id`, `intention_model_config_id`:
			if cast.ToInt(val) > 0 {
				chatRobot[key] = cast.ToString(models.LlmModelConfigId)
			}
		case `use_model`, `optimize_question_use_model`, `intention_use_model`:
			if len(val) > 0 {
				chatRobot[key] = models.LlmUseModel
			}
		case `rerank_model_config_id`:
			if cast.ToInt(val) > 0 {
				chatRobot[key] = cast.ToString(models.RerankModelConfigId)
			}
		case `rerank_use_model`:
			if len(val) > 0 {
				chatRobot[key] = models.RerankUseModel
			}
		case `unknown_summary_model_config_id`:
			if cast.ToInt(val) > 0 {
				chatRobot[key] = cast.ToString(models.VectorModelConfigId)
			}
		case `unknown_summary_use_model`:
			if len(val) > 0 {
				chatRobot[key] = models.VectorUseModel
			}
		case `robot_avatar`:
			chatRobot[`avatar_from_template`] = val
		default:
			chatRobot[key] = val
		}
	}
	code, err := common.RequestChatWiki(`/manage/saveRobot`, http.MethodPost, token, chatRobot)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapString(code.Data), nil
}

func ApplyFlowRobot(adminUserId int, robot msql.Params, nodes []msql.Params, cslIdMaps *common.CslIdMaps, models *common.DefaultModelParams, token string) (msql.Params, error) {
	flowRobot := map[string]string{
		`robot_name`:           robot[`robot_name`],
		`robot_intro`:          robot[`robot_intro`],
		`avatar_from_template`: robot[`robot_avatar`],
	}
	code, err := common.RequestChatWiki(`/manage/addFlowRobot`, http.MethodPost, token, flowRobot)
	if err != nil {
		return nil, err
	}
	newRobot := cast.ToStringMapString(code.Data)
	//工作流节点
	for _, node := range nodes {
		nodeType := cast.ToInt(node[`node_type`])
		if !tool.InArrayInt(nodeType, work_flow.NodeTypes[:]) {
			logs.Error(`节点类型非法:%s`, node[`node_type`])
			continue
		}
		insData := msql.Datas{`admin_user_id`: adminUserId, `robot_id`: newRobot[`id`]}
		for key, val := range node {
			switch key {
			case `id`, `admin_user_id`, `robot_id`, `node_info_json`:
			case `node_params`:
				insData[key], insData[`node_info_json`] = ReplaceNodeParams(adminUserId, nodeType, val, node[`node_info_json`], cslIdMaps, models)
			default:
				insData[key] = val
			}
		}
		_, err = msql.Model(`work_flow_node`, define.Postgres).Insert(insData)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return newRobot, nil
}

func ReplaceNodeParams(adminUserId int, nodeType int, nodeParamsStr, nodeInfoStr string, cslIdMaps *common.CslIdMaps, models *common.DefaultModelParams) (string, string) {
	nodeParams := work_flow.NodeParams{}
	if err := tool.JsonDecodeUseNumber(nodeParamsStr, &nodeParams); err != nil {
		logs.Error(err.Error())
	}
	replace := make(map[string]any) //修正node_info_json.dataRaw数据
	switch nodeType {
	case work_flow.NodeTypeCate:
		nodeParams.Cate.ModelConfigId = common.MixedInt(models.LlmModelConfigId)
		nodeParams.Cate.UseModel = models.LlmUseModel
		replace[`model_config_id`] = models.LlmModelConfigId
		replace[`use_model`] = models.LlmUseModel
	case work_flow.NodeTypeLibs:
		newLibraryIds := make([]string, 0)
		if len(nodeParams.Libs.LibraryIds) > 0 {
			for _, oldLibraryId := range strings.Split(nodeParams.Libs.LibraryIds, `,`) {
				if newLibraryId := cslIdMaps.Librarys[cast.ToInt(oldLibraryId)]; newLibraryId > 0 {
					newLibraryIds = append(newLibraryIds, cast.ToString(newLibraryId))
				}
			}
		}
		nodeParams.Libs.LibraryIds = strings.Join(newLibraryIds, `,`)
		nodeParams.Libs.RerankModelConfigId = common.MixedInt(models.RerankModelConfigId)
		nodeParams.Libs.RerankUseModel = models.RerankUseModel
		replace[`library_ids`] = nodeParams.Libs.LibraryIds
		replace[`model_config_id`] = models.RerankModelConfigId
		replace[`use_model`] = models.RerankUseModel
	case work_flow.NodeTypeLlm:
		nodeParams.Llm.ModelConfigId = common.MixedInt(models.LlmModelConfigId)
		nodeParams.Llm.UseModel = models.LlmUseModel
		replace[`model_config_id`] = models.LlmModelConfigId
		replace[`use_model`] = models.LlmUseModel
	case work_flow.NodeTypeManual:
		nodeParams.Manual.SwitchType = work_flow.StaffAll
		nodeParams.Manual.SwitchStaff = ``
		replace[`switch_type`] = work_flow.StaffAll
		replace[`switch_staff`] = ``
	case work_flow.NodeTypeQuestionOptimize:
		nodeParams.QuestionOptimize.ModelConfigId = common.MixedInt(models.LlmModelConfigId)
		nodeParams.QuestionOptimize.UseModel = models.LlmUseModel
		replace[`model_config_id`] = models.LlmModelConfigId
		replace[`use_model`] = models.LlmUseModel
	case work_flow.NodeTypeParamsExtractor:
		nodeParams.ParamsExtractor.ModelConfigId = common.MixedInt(models.LlmModelConfigId)
		nodeParams.ParamsExtractor.UseModel = models.LlmUseModel
		replace[`model_config_id`] = models.LlmModelConfigId
		replace[`use_model`] = models.LlmUseModel
	case work_flow.NodeTypeFormInsert:
		nodeParams.FormInsert.FormId = common.MixedInt(cslIdMaps.Forms[nodeParams.FormInsert.FormId.Int()])
		replace[`admin_user_id`] = adminUserId
		replace[`form_id`] = nodeParams.FormInsert.FormId
	case work_flow.NodeTypeFormDelete:
		nodeParams.FormDelete.FormId = common.MixedInt(cslIdMaps.Forms[nodeParams.FormDelete.FormId.Int()])
		replace[`admin_user_id`] = adminUserId
		replace[`form_id`] = nodeParams.FormDelete.FormId
		nodeParams.FormDelete.Where = ReplaceFormOpWhere(nodeParams.FormDelete.Where, cslIdMaps, &replace)
	case work_flow.NodeTypeFormUpdate:
		nodeParams.FormUpdate.FormId = common.MixedInt(cslIdMaps.Forms[nodeParams.FormUpdate.FormId.Int()])
		replace[`admin_user_id`] = adminUserId
		replace[`form_id`] = nodeParams.FormUpdate.FormId
		nodeParams.FormUpdate.Where = ReplaceFormOpWhere(nodeParams.FormUpdate.Where, cslIdMaps, &replace)
	case work_flow.NodeTypeFormSelect:
		nodeParams.FormSelect.FormId = common.MixedInt(cslIdMaps.Forms[nodeParams.FormSelect.FormId.Int()])
		replace[`admin_user_id`] = adminUserId
		replace[`form_id`] = nodeParams.FormSelect.FormId
		nodeParams.FormSelect.Where = ReplaceFormOpWhere(nodeParams.FormSelect.Where, cslIdMaps, &replace)
	}
	return tool.JsonEncodeNoError(nodeParams), AmendNodeinfojson(nodeInfoStr, replace)
}

func ReplaceFormOpWhere(where []define.FormFilterCondition, cslIdMaps *common.CslIdMaps, replace *map[string]any) []define.FormFilterCondition {
	newWhere := make([]define.FormFilterCondition, 0)
	for _, condition := range where {
		oldFormFieldId := condition.FormFieldId
		newFormFieldId := cslIdMaps.FormFields[oldFormFieldId]
		if oldFormFieldId == 0 || newFormFieldId > 0 {
			condition.FormFieldId = newFormFieldId
			newWhere = append(newWhere, condition)
		}
		(*replace)[fmt.Sprintf(`form_field_id#%d`, oldFormFieldId)] = newFormFieldId
	}
	return newWhere
}

func AmendNodeinfojson(nodeInfoStr string, replace map[string]any) string {
	if len(replace) == 0 {
		return nodeInfoStr
	}
	nodeInfoJson := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(nodeInfoStr, &nodeInfoJson); err != nil {
		logs.Error(err.Error())
	}
	if _, ok := nodeInfoJson[`dataRaw`]; !ok {
		return nodeInfoStr
	}
	var dataRaw any
	if err := tool.JsonDecodeUseNumber(cast.ToString(nodeInfoJson[`dataRaw`]), &dataRaw); err != nil {
		logs.Error(err.Error())
	}
	dataRaw = AmendDataRaw(dataRaw, replace)
	nodeInfoJson[`dataRaw`] = tool.JsonEncodeNoError(dataRaw)
	return tool.JsonEncodeNoError(nodeInfoJson)
}

func AmendDataRaw(dataRaw any, replace map[string]any, keys ...string) any {
	switch realData := dataRaw.(type) {
	case []any:
		newData := make([]interface{}, len(realData))
		for idx, val := range realData {
			newData[idx] = AmendDataRaw(val, replace)
		}
		return newData
	case map[string]any:
		newData := make(map[string]interface{})
		for key, val := range realData {
			newData[key] = AmendDataRaw(val, replace, key)
		}
		return newData
	default:
		if len(keys) > 0 {
			if newVal, ok := replace[fmt.Sprintf(`%s#%v`, keys[0], realData)]; ok {
				return newVal //替换成新值
			}
			if newVal, ok := replace[keys[0]]; ok {
				return newVal //替换成新值
			}
		}
		return realData
	}
}
