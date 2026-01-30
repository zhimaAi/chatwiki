// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetMcpServerList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	info, err := msql.Model(`mcp_server`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	toolList := make([]msql.Params, 0)
	if len(info) > 0 {
		toolList, err = msql.Model(`mcp_tool`, define.Postgres).
			Alias(`m`).
			Where(`m.admin_user_id`, cast.ToString(adminUserId)).
			Where(`server_id`, cast.ToString(info[`id`])).
			Join(`chat_ai_robot r`, `m.robot_id = r.id`, `inner`).
			Join(`work_flow_node_version w`, `r.id = w.robot_id and r.start_node_Key = w.node_key and w.work_flow_version_id = (SELECT id FROM work_flow_version WHERE robot_id = r.id ORDER BY create_time DESC LIMIT 1)`, `left`).
			Field(`m.*,r.robot_name,r.robot_intro,r.robot_avatar,w.node_params`).
			Select()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		// 简化 node_params 结构，只返回必要的参数信息
		for i := range toolList {
			simplifiedParams := simplifyNodeParams(toolList[i])
			// 将简化后的参数序列化为 JSON 字符串
			simplifiedParamsJSON, err := json.Marshal(simplifiedParams)
			if err != nil {
				logs.Error("Failed to marshal simplified_params: %v", err)
				continue
			}
			// 创建新的 map 来存储简化后的参数
			newTool := make(msql.Params)
			for k, v := range toolList[i] {
				newTool[k] = v
			}
			newTool["params"] = string(simplifiedParamsJSON)
			delete(newTool, `node_params`)
			toolList[i] = newTool
		}
	}

	response := map[string]interface{}{
		"info":      info,
		"tool_list": toolList,
	}

	c.String(http.StatusOK, lib_web.FmtJson(response, nil))
}

func SaveMcpServer(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	serverId := cast.ToInt(c.PostForm(`server_id`))
	name := strings.TrimSpace(c.PostForm(`name`))
	description := strings.TrimSpace(c.PostForm(`description`))

	// check params
	if len(name) == 0 || len(description) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if utf8.RuneCountInString(name) > 100 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `name`))))
		return
	}
	if utf8.RuneCountInString(description) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `description`))))
		return
	}

	// upload file
	avatar := ``
	if serverId == 0 {
		avatar = define.LocalUploadPrefix + `default/mcp_avatar.svg`
	}
	fileHeader, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, adminUserId, `mcp_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}

	if serverId > 0 { // edit
		mcpServer, err := msql.Model(`mcp_server`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(serverId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(mcpServer) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}

		updateData := msql.Datas{
			`name`:        name,
			`description`: description,
			`update_time`: tool.Time2Int(),
		}
		if len(avatar) > 0 {
			updateData[`avatar`] = avatar
		}

		_, err = msql.Model(`mcp_server`, define.Postgres).
			Where(`id`, cast.ToString(serverId)).
			Update(updateData)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
		return
	} else { // add
		old, err := msql.Model(`mcp_server`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(old) > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exist_relation_mcp_server`))))
			return
		}

		generatePrefixedKey := func(prefix string, userID int, length int) (string, error) {
			b := make([]byte, length)
			if _, err := rand.Read(b); err != nil {
				return "", err
			}
			data := append(b, []byte(fmt.Sprintf("%d", userID))...)
			hash := sha256.Sum256(data)
			keyPart := hex.EncodeToString(hash[:])[:length*2]

			return fmt.Sprintf("%s_%s", prefix, keyPart), nil
		}
		apiKey, err := generatePrefixedKey(`chatwiki`, adminUserId, 16)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		_, err = msql.Model(`mcp_server`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`:  adminUserId,
			`name`:           name,
			`description`:    description,
			`create_time`:    tool.Time2Int(),
			`update_time`:    tool.Time2Int(),
			`avatar`:         avatar,
			`api_key`:        apiKey,
			`publish_status`: define.McpServerDraft,
		}, `id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
	}
}

func UpdateMcpServerPublishStatus(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	serverId := cast.ToInt(c.PostForm(`server_id`))
	publishStatus := cast.ToInt(c.PostForm(`publish_status`))
	if serverId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if publishStatus != 0 && publishStatus != 1 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `publish_status`, `publish_status`))))
		return
	}
	mcpServer, err := msql.Model(`mcp_server`, define.Postgres).Where(`id`, cast.ToString(serverId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(mcpServer) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	toolCount, err := msql.Model(`mcp_tool`, define.Postgres).Where(`server_id`, cast.ToString(serverId)).Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if toolCount == 0 && publishStatus == 1 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_workflow`))))
		return
	}
	_, err = msql.Model(`mcp_server`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(serverId)).
		Update(msql.Datas{
			`publish_status`: publishStatus,
			`update_time`:    tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DeleteMcpServer(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	serverId := strings.TrimSpace(c.PostForm("server_id"))
	if len(serverId) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	_, err := msql.Model(`mcp_server`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(serverId)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	_, err = msql.Model(`mcp_tool`, define.Postgres).
		Where(`server_id`, cast.ToString(serverId)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	common.ClearMCPServerCache(adminUserId)

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func EditMcpTool(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	toolId := strings.TrimSpace(c.PostForm("tool_id"))
	name := strings.TrimSpace(c.PostForm("name"))
	if len(toolId) == 0 || len(name) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	// Validate function name format
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if !matched {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `invalid_function_name`))))
		return
	}
	mcpTool, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, toolId).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(mcpTool) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	duplicate, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, `!=`, toolId).
		Where(`name`, name).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(duplicate) > 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `duplicated_mcp_tool_name`))))
		return
	}

	_, err = msql.Model(`mcp_tool`, define.Postgres).
		Where(`id`, cast.ToString(mcpTool[`id`])).
		Update(msql.Datas{
			`name`:        name,
			`update_time`: tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	common.ClearMCPServerCache(adminUserId)

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func SaveMcpTool(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	serverId := cast.ToInt(c.PostForm("server_id"))
	robotIdsStr := strings.TrimSpace(c.PostForm("robot_id"))

	if serverId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	// 验证 MCP Server 是否存在
	mcpServer, err := msql.Model(`mcp_server`, define.Postgres).
		Where(`id`, cast.ToString(serverId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(mcpServer) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	// 解析传入的 robot_id 列表
	newRobotIds := make(map[int]bool)
	if len(robotIdsStr) > 0 {
		for _, idStr := range strings.Split(robotIdsStr, ",") {
			id := cast.ToInt(strings.TrimSpace(idStr))
			if id > 0 {
				newRobotIds[id] = true
			}
		}
	}

	// 获取当前 server_id 下的所有 mcp_tool
	existingTools, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`server_id`, cast.ToString(serverId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// 构建现有的 robot_id 映射
	existingRobotIds := make(map[int]int) // robot_id -> tool_id
	for _, t := range existingTools {
		robotId := cast.ToInt(t[`robot_id`])
		toolId := cast.ToInt(t[`id`])
		existingRobotIds[robotId] = toolId
	}

	// 找出需要添加的 robot_id（新列表中有，旧列表中没有）
	for robotId := range newRobotIds {
		if _, exists := existingRobotIds[robotId]; !exists {
			// 验证 robot 是否存在且类型正确
			robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
				Where(`admin_user_id`, cast.ToString(adminUserId)).
				Where(`application_type`, cast.ToString(define.ApplicationTypeFlow)).
				Where(`id`, cast.ToString(robotId)).
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
			if cast.ToInt(robot[`application_type`]) == define.ApplicationTypeFlow && cast.ToUint(robot[`data_type`]) == define.DataTypeDraft {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `exist_relation_mcp_tool`))))
				return
			}

			// 生成随机名称
			getRandomStr := func() (string, error) {
				b := make([]byte, 10)
				if _, err := rand.Read(b); err != nil {
					return "", err
				}
				return "chatwiki_" + hex.EncodeToString(b), nil
			}
			name, err := getRandomStr()
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}

			// 添加新的 mcp_tool
			_, err = msql.Model(`mcp_tool`, define.Postgres).Insert(msql.Datas{
				`admin_user_id`: adminUserId,
				`server_id`:     serverId,
				`robot_id`:      robotId,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
				`name`:          name,
			}, `id`)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
		}
	}

	// 找出需要删除的 robot_id（旧列表中有，新列表中没有）
	for robotId, toolId := range existingRobotIds {
		if !newRobotIds[robotId] {
			_, err = msql.Model(`mcp_tool`, define.Postgres).
				Where(`id`, cast.ToString(toolId)).
				Where(`admin_user_id`, cast.ToString(adminUserId)).
				Delete()
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
		}
	}

	// 检查删除后是否还有 tool，如果没有则将 server 的 publish_status 设为 0
	count, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`server_id`, cast.ToString(serverId)).
		Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if count == 0 {
		_, err = msql.Model(`mcp_server`, define.Postgres).
			Where(`id`, cast.ToString(serverId)).
			Update(msql.Datas{
				`update_time`:    tool.Time2Int(),
				`publish_status`: 0,
			})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

	common.ClearMCPServerCache(adminUserId)

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DeleteMcpTool(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	toolId := strings.TrimSpace(c.PostForm("tool_id"))
	if len(toolId) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	mcpTool, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(toolId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(mcpTool) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	serverId := cast.ToInt(mcpTool[`server_id`])

	_, err = msql.Model(`mcp_tool`, define.Postgres).
		Where(`id`, cast.ToString(mcpTool[`id`])).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	count, err := msql.Model(`mcp_tool`, define.Postgres).
		Where(`server_id`, cast.ToString(serverId)).
		Count()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if count == 0 {
		_, err = msql.Model(`mcp_server`, define.Postgres).
			Where(`id`, cast.ToString(serverId)).
			Update(msql.Datas{
				`update_time`:    tool.Time2Int(),
				`publish_status`: 0,
			})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

	common.ClearMCPServerCache(adminUserId)

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// simplifyNodeParams 简化 node_params 结构，只返回必要的参数信息
// 默认返回 content 和 open_id，如果 node_params 不为空则还包含 diy_global 中的参数
func simplifyNodeParams(toolData msql.Params) []map[string]interface{} {
	params := make([]map[string]interface{}, 0)

	// 默认参数：content 和 open_id
	params = append(params, map[string]interface{}{
		"key":      "content",
		"desc":     lib_define.DialogueContent,
		"type":     "string",
		"required": true,
	})
	params = append(params, map[string]interface{}{
		"key":      "open_id",
		"desc":     lib_define.UserRequestIdentifier,
		"type":     "string",
		"required": true,
	})

	// 如果 node_params 不为空，提取 diy_global 中的参数
	nodeParamsStr := cast.ToString(toolData["node_params"])
	if nodeParamsStr != "" {
		var nodeParams work_flow.NodeParams
		if err := json.Unmarshal([]byte(nodeParamsStr), &nodeParams); err != nil {
			logs.Error("Failed to parse node_params: %v", err)
			return params
		}

		// 添加 diy_global 中的参数
		for _, param := range nodeParams.Start.DiyGlobal {
			params = append(params, map[string]interface{}{
				"key":      param.Key,
				"desc":     param.Desc,
				"type":     param.Typ,
				"required": param.Required,
			})
		}
	}

	return params
}
