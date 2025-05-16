// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func SaveLibrarySearch(c *gin.Context) {
	var userId int
	if userId = getLoginUserId(c); userId == 0 {
		return
	}
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	temperature := cast.ToFloat32(c.DefaultPostForm(`temperature`, `0.5`))
	maxToken := cast.ToInt(c.DefaultPostForm(`max_token`, `2000`))
	contextPair := cast.ToInt(c.DefaultPostForm(`context_pair`, `6`))
	topK := cast.ToInt(c.DefaultPostForm(`size`, `200`))
	similarity := cast.ToFloat32(c.DefaultPostForm(`similarity`, `0.6`))
	searchType := cast.ToInt(c.DefaultPostForm(`search_type`, `1`))
	rerankStatus := cast.ToInt(c.PostForm(`rerank_status`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	promptType := cast.ToString(c.PostForm(`prompt_type`))
	prompt := cast.ToString(c.PostForm(`prompt`))
	if len(prompt) == 0 && cast.ToInt(promptType) != 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `prompt`))))
		return
	}
	m := msql.Model(`user_search_config`, define.Postgres)
	data, err := m.Where("user_id", cast.ToString(userId)).Field(`id`).Find()
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	// 准备保存的数据
	saveData := msql.Datas{
		"user_id":                userId,
		"model_config_id":        modelConfigId,
		"use_model":              useModel,
		"temperature":            temperature,
		"max_token":              maxToken,
		"context_pair":           contextPair,
		"size":                   topK,
		"similarity":             similarity,
		"search_type":            searchType,
		"rerank_status":          rerankStatus,
		"rerank_model_config_id": rerankModelConfigID,
		"rerank_use_model":       rerankUseModel,
		"update_time":            tool.Time2Int(),
		"prompt_type":            promptType,
		"prompt":                 prompt,
	}
	if len(data) == 0 {
		saveData[`create_time`] = tool.Time2Int()
		_, err = m.Insert(saveData)
	} else {
		_, err = m.Where("id", cast.ToString(data[`id`])).Update(saveData)
	}
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, nil)
}

// GetLibrarySearch
func GetLibrarySearch(c *gin.Context) {
	var userId int
	if userId = getLoginUserId(c); userId == 0 {
		return
	}

	m := msql.Model(`user_search_config`, define.Postgres)
	data, err := m.Where("user_id", cast.ToString(userId)).Find()
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, data)
}

func LibraryAiSummary(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	c.Header(`Content-Type`, `text/event-stream`)
	c.Header(`Cache-Control`, `no-cache`)
	c.Header(`Connection`, `keep-alive`)
	libraryIds := cast.ToString(c.PostForm(`id`))
	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	question := strings.TrimSpace(c.PostForm(`question`))
	size := cast.ToInt(c.DefaultPostForm(`size`, "10"))
	similarity := cast.ToFloat64(c.PostForm(`similarity`))
	searchType := cast.ToInt(c.PostForm(`search_type`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	rerankStatus := strings.TrimSpace(c.DefaultPostForm(`rerank_status`, `1`))
	temperature := cast.ToFloat64(c.DefaultPostForm(`temperature`, `0.5`))
	maxToken := cast.ToInt(c.DefaultPostForm(`max_token`, `2000`))
	recallType := cast.ToString(c.PostForm(`recall_type`))
	prompt := cast.ToString(c.PostForm(`prompt`))
	if modelConfigId <= 0 || useModel == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `llm_model_err`))))
		return
	}
	if len(libraryIds) <= 0 || len(question) == 0 || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if len(prompt) == 0 {
		prompt = define.PromptLibAiSummary
	}
	if size > 10 {
		size = 10
	}
	if searchType != define.SearchTypeMixed && searchType != define.SearchTypeVector && searchType != define.SearchTypeFullText && searchType != define.SearchTypeGraph {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}

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
		}
		robot[`admin_user_id`] = info[`admin_user_id`]
		robot[`model_config_id`] = cast.ToString(modelConfigId)
		robot[`use_model`] = useModel
		robot[`id`] = strconv.Itoa(0)
	}
	var (
		wg = sync.WaitGroup{}
	)
	chanStream := make(chan sse.Event)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := common.LibraryAiSummary(common.GetLang(c), question, prompt, []string{}, libraryIds, size, maxToken, similarity, temperature, searchType, robot, chanStream); err != nil {
			common.FmtError(c, `sys_err`)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Stream(func(_ io.Writer) bool {
			if event, ok := <-chanStream; ok {
				if data, ok := event.Data.(string); ok {
					event.Data = strings.ReplaceAll(data, "\r", ``)
				}
				c.SSEvent(event.Event, event.Data)
				return true
			}
			return false
		})
	}()
	wg.Wait()
}
