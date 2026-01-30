// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func getParamsFromUri(uri, lang string) (string, string, error) {
	findKey := common.GetUrlPath(uri)
	if findKey == `` {
		return ``, ``, errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	if !strings.HasPrefix(findKey, `/open/workflow/webhook/`) {
		return ``, ``, errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	findKeyParams := strings.Split(findKey, `/`)
	if len(findKeyParams) != 6 {
		return ``, ``, errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	if findKeyParams[4] == `` || len(findKeyParams[5]) != 10 {
		return ``, ``, errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	return findKeyParams[5], findKeyParams[4], nil
}

func SaveTriggerWebhookConfig(adminUserId string, trigger TriggerConfig, robot msql.Params, lang string) error {
	findKey, robotKey, err := getParamsFromUri(trigger.TriggerWebHookConfig.Url, lang)
	if findKey == `` {
		return errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	if robotKey != robot[`robot_key`] {
		return errors.New(i18n.Show(lang, `param_err`, `url`))
	}
	_, err = msql.Model(`work_flow_trigger`, define.Postgres).Insert(msql.Datas{
		`admin_user_id`: cast.ToInt(robot[`admin_user_id`]),
		`robot_id`:      cast.ToString(robot[`id`]),
		`trigger_type`:  TriggerTypeWebHook,
		`trigger_json`:  tool.JsonEncodeNoError(trigger),
		`find_key`:      findKey,
		`create_time`:   time.Now().Unix(),
		`update_time`:   time.Now().Unix(),
	}, `id`)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	lib_redis.DelCacheData(define.Redis, common.TriggerFindKeyCacheBuildHandler{
		AdminUserId: adminUserId,
		FindKey:     findKey,
	})
	return nil
}

func TriggerWebhookVerifyStart(findKey, robotKey string, c *gin.Context) (robotInfo msql.Params, webhookConfig TriggerWebHookConfig, trigger msql.Params, err error) {
	robotInfo, err = common.GetRobotInfo(robotKey)
	if err != nil {
		return
	}
	if len(robotInfo) == 0 {
		err = errors.New(`robot not exist`)
		return
	}
	err = lib_redis.GetCacheWithBuild(define.Redis, common.TriggerFindKeyCacheBuildHandler{
		AdminUserId: robotInfo[`admin_user_id`],
		FindKey:     findKey,
	}, &trigger, time.Hour*24)
	if err != nil {
		return
	}
	if len(trigger) == 0 {
		err = errors.New(`trigger not exist`)
		return
	}
	triggerConfig := TriggerConfig{}
	err = tool.JsonDecode(trigger[`trigger_json`], &triggerConfig)
	if err != nil {
		return
	}
	webhookConfig = triggerConfig.TriggerWebHookConfig
	if cast.ToInt(triggerConfig.TriggerWebHookConfig.SwitchVerify) > 0 {
		_, err = common.ParseAuthorizationToken(c)
		if err != nil {
			return
		}
	}
	nodeParam, err := msql.Model(`work_flow_node`, define.Postgres).
		Where(`robot_id`, robotInfo[`id`]).
		Where(`node_type`, cast.ToString(NodeTypeStart)).Value(`node_params`)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	isOk := false
	nodeParams := NodeParams{}
	err = tool.JsonDecode(nodeParam, &nodeParams)
	if err != nil {
		return
	}
	if len(nodeParams.Start.TriggerList) == 0 {
		err = errors.New(`workflow triggers is empty`)
		return
	}
	realIp := lib_web.GetClientIP(c)
	for _, triggerVal := range nodeParams.Start.TriggerList {
		if triggerVal.TriggerSwitch == false {
			continue
		}
		if triggerVal.TriggerType != triggerConfig.TriggerType {
			continue
		}
		if triggerVal.TriggerWebHookConfig.Url != triggerConfig.TriggerWebHookConfig.Url {
			continue
		}
		if cast.ToInt(triggerVal.TriggerWebHookConfig.SwitchAllowIp) > 0 && !tool.InArray(realIp, strings.Split(triggerVal.TriggerWebHookConfig.AllowIps, "\n")) {
			logs.Debug(LogTriggerPrefix + ` ip (` + realIp + `) not allow`)
			continue
		}
		isOk = true
	}
	if isOk == false {
		err = errors.New(`not find trigger`)
		return
	}
	return
}

func StartWebhook(c *gin.Context) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf(`%v`, r))
		}
	}()

	robotKey := c.Param(`robot_key`)
	findKey := c.Param(`find_key`)
	if findKey == `` || robotKey == `` {
		err = errors.New(`missing parameter`)
		return
	}
	robotInfo, webhookConfig, trigger, err := TriggerWebhookVerifyStart(findKey, robotKey, c)
	if err != nil {
		return
	}
	// testParams extract passed parameters, triggerOutputs used for trigger outputs assignment
	testParams, triggerOutputs := getRunParams(robotInfo, webhookConfig, c)
	if webhookConfig.ResponseType == WebHookResponseTypeNow {
		ret = webhookConfig.ResponseNow
		go func() {
			_, err := runWorkFlow(robotInfo, trigger, testParams, triggerOutputs)
			if err != nil {
				logs.Error(LogTriggerPrefix + err.Error())
			}
		}()
		return
	}
	resData, err := runWorkFlow(robotInfo, trigger, testParams, triggerOutputs)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	ret = tool.JsonEncodeNoError(resData)
	return
}

func getRunParams(robotInfo msql.Params, webhookConfig TriggerWebHookConfig, c *gin.Context) (map[string]any, []TriggerOutputParam) {
	testParams := make(map[string]any)
	triggerOutputs := make([]TriggerOutputParam, 0)
	getParams(webhookConfig, testParams, &triggerOutputs, c)
	if webhookConfig.Method == http.MethodPost {
		if webhookConfig.RequestContentType == define.ContentTypeMultipart {
			getFormTestParams(webhookConfig, robotInfo, testParams, &triggerOutputs, c)
		} else if webhookConfig.RequestContentType == define.ContentTypeJson {
			getJsonTestParams(webhookConfig, testParams, &triggerOutputs, c)
		} else if webhookConfig.RequestContentType == define.ContentTypeFormUrl {
			getXFormTestParams(webhookConfig, testParams, &triggerOutputs, c)
		}
	}
	return testParams, triggerOutputs
}

func runWorkFlow(robotInfo msql.Params, trigger msql.Params, testParams map[string]any, triggerOutputs []TriggerOutputParam) (resData map[string]any, err error) {
	resData = make(map[string]any)

	startTime := time.Now().Unix()
	workFlowParams := &WorkFlowParams{
		ChatRequestParam: &define.ChatRequestParam{
			ChatBaseParam: &define.ChatBaseParam{
				AdminUserId: cast.ToInt(robotInfo[`admin_user_id`]),
				Robot:       robotInfo,
			},
		},
		TriggerParams: TriggerParams{
			TriggerType:    TriggerTypeWebHook,
			TestParams:     testParams,
			TriggerOutputs: triggerOutputs,
		},
	}
	flow, _, err := BaseCallWorkFlow(workFlowParams)
	setRunResult(trigger[`id`], startTime, err)
	if err != nil {
		return
	}
	content, _ := TakeOutputReply(flow)
	if variables, ok := flow.output[`special.finish_variables`]; ok {
		if len(variables.Vals) > 0 {
			decodeErr := tool.JsonDecode(cast.ToString(*variables.Vals[0].String), &resData)
			if decodeErr != nil {
				logs.Error(LogTriggerPrefix + decodeErr.Error())
			}
		}
	}
	resData[`content`] = content
	return
}

func getParams(webhookConfig TriggerWebHookConfig, testParams map[string]any, triggerOutputs *[]TriggerOutputParam, c *gin.Context) {
	if len(webhookConfig.Params) > 0 {
		for _, param := range webhookConfig.Params {
			testParams[param.Key] = c.Query(param.Key)
			variable := cast.ToString(param.Desc)
			if variable != `` {
				*triggerOutputs = append(*triggerOutputs, TriggerOutputParam{
					StartNodeParam: StartNodeParam{
						Key: param.Key,
						Typ: param.Typ,
					},
					Variable: variable,
				})
			}
		}
	}
}

func getXFormTestParams(webhookConfig TriggerWebHookConfig, testParams map[string]any, triggerOutputs *[]TriggerOutputParam, c *gin.Context) {
	if len(webhookConfig.XForm) == 0 {
		return
	}
	for _, param := range webhookConfig.XForm {
		if param.Typ == common.TypString {
			testParams[param.Key] = c.PostForm(param.Key)
			variable := cast.ToString(param.Desc)
			if variable != `` {
				*triggerOutputs = append(*triggerOutputs, TriggerOutputParam{
					StartNodeParam: StartNodeParam{
						Key: param.Key,
						Typ: param.Typ,
					},
					Variable: variable,
				})
			}
		}
	}
}

func getFormTestParams(webhookConfig TriggerWebHookConfig, robotInfo msql.Params, testParams map[string]any, triggerOutputs *[]TriggerOutputParam, c *gin.Context) {
	if len(webhookConfig.Form) == 0 {
		return
	}
	for _, param := range webhookConfig.Form {
		if param.Typ == `file` {
			file, fileErr := c.FormFile(param.Key)
			if fileErr != nil {
				logs.Error(LogTriggerPrefix + fileErr.Error())
				continue
			}
			uploadInfo, err := common.SaveUploadedFile(file, define.LibFileLimitSize, robotInfo[`admin_user_id`],
				`work_flow_webhook`, define.AllExt)
			if err != nil {
				logs.Error(LogTriggerPrefix + err.Error())
				continue
			}
			if uploadInfo == nil {
				logs.Warning(LogTriggerPrefix + `upload file is empty`)
				continue
			}
			if !common.IsUrl(uploadInfo.Link) {
				testParams[param.Key] = define.Config.WebService[`image_domain`] + uploadInfo.Link
			} else {
				testParams[param.Key] = uploadInfo.Link
			}
		} else if param.Typ == common.TypString {
			testParams[param.Key] = c.PostForm(param.Key)
		}
		variable := cast.ToString(param.Desc)
		if variable != `` {
			*triggerOutputs = append(*triggerOutputs, TriggerOutputParam{
				StartNodeParam: StartNodeParam{
					Key: param.Key,
					Typ: common.TypString,
				},
				Variable: variable,
			})
		}
	}
}

func getJsonTestParams(webhookConfig TriggerWebHookConfig, testParams map[string]any, triggerOutputs *[]TriggerOutputParam, c *gin.Context) {
	if len(webhookConfig.Json) == 0 {
		return
	}
	bodyBs, err := c.GetRawData()
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	if len(bodyBs) == 0 {
		logs.Warning(LogTriggerPrefix + `body is empty`)
		return
	}
	bodyMap := make(map[string]any)
	err = tool.JsonDecode(cast.ToString(bodyBs), &bodyMap)
	if err != nil {
		logs.Error(LogTriggerPrefix + err.Error())
		return
	}
	jsonOutputs := common.SimpleFields{}
	webhookConfig.Json.SimplifyFieldsDeepExtract(&jsonOutputs, ``, bodyMap)
	for _, param := range jsonOutputs {
		var val any
		var vals = make([]any, 0)
		switch param.Typ {
		case common.TypString:
			if len(param.Vals) > 0 {
				val = *param.Vals[0].String
				testParams[param.Key] = val
			}
		case common.TypNumber:
			if len(param.Vals) > 0 {
				val = *param.Vals[0].Number
				testParams[param.Key] = val
			}
		case common.TypBoole:
			if len(param.Vals) > 0 {
				val = *param.Vals[0].Boole
				testParams[param.Key] = val
			}
		case common.TypArrString:
			if len(param.Vals) > 0 {
				for _, outputVal := range param.Vals {
					vals = append(vals, *outputVal.String)
				}
				testParams[param.Key] = vals
			}
		case common.TypArrNumber:
			if len(param.Vals) > 0 {
				for _, outputVal := range param.Vals {
					vals = append(vals, *outputVal.Number)
				}
				testParams[param.Key] = vals
			}
		case common.TypArrObject:
			if len(param.Vals) > 0 {
				for _, outputVal := range param.Vals {
					vals = append(vals, outputVal.Object)
				}
				testParams[param.Key] = vals
			}
		}
		variable := cast.ToString(param.Desc)
		if variable != `` {
			*triggerOutputs = append(*triggerOutputs, TriggerOutputParam{
				StartNodeParam: StartNodeParam{
					Key: param.Key,
					Typ: param.Key,
				},
				Variable: variable,
			})
		}
	}
}
