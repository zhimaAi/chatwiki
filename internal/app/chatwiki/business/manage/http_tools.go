// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/app/chatwiki/work_flow"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
)

// HttpToolRequest HTTP工具请求结构
type HttpToolRequest struct {
	ID          int    `form:"id" json:"id"`
	Name        string `form:"name" json:"name" binding:"required"`
	NameEn      string `form:"name_en" json:"name_en"`
	ToolKey     string `form:"tool_key" json:"tool_key"`
	Avatar      string `form:"avatar" json:"avatar"`
	Description string `form:"description" json:"description"`
}

// HttpToolNodeRequest HTTP工具节点请求结构
type HttpToolNodeRequest struct {
	ID              int    `form:"id" json:"id"`
	HttpToolID      int    `form:"http_tool_id" json:"http_tool_id" binding:"required"`
	NodeKey         string `form:"node_key" json:"node_key"`
	NodeNameEn      string `form:"node_name_en" json:"node_name_en"`
	NodeName        string `form:"node_name" json:"node_name" binding:"required"`
	NodeDescription string `form:"node_description" json:"node_description"`
	NodeRemark      string `form:"node_remark" json:"node_remark"`
	DataRaw         string `form:"data_raw" json:"data_raw"` // 原始数据，用于生成node_params
}

// HttpToolListFilterRequest HTTP工具列表过滤请求结构
type HttpToolListFilterRequest struct {
	Name      string `json:"name" form:"name"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	WithNodes bool   `json:"with_nodes" form:"with_nodes"` // 是否同时获取节点详情
}

// HttpToolNodeListFilterRequest HTTP工具节点列表过滤请求结构
type HttpToolNodeListFilterRequest struct {
	HttpToolID int    `json:"http_tool_id" form:"http_tool_id" binding:"required"`
	NodeName   string `json:"node_name" form:"node_name"`
	Page       int    `json:"page" form:"page"`
	Size       int    `json:"size" form:"size"`
}

// SaveHttpTool 保存HTTP工具（创建或更新）
func SaveHttpTool(c *gin.Context) {
	var req HttpToolRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 如果是更新操作，检查工具是否存在且属于当前用户
	if req.ID > 0 {
		toolInfo, err := common.GetHttpTool(req.ID)
		if err != nil {
			logs.Error("Get http tool error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}

		if len(toolInfo) == 0 {
			common.FmtError(c, `http_tool_not_exist`)
			return
		}

		if cast.ToInt(toolInfo["admin_user_id"]) != adminUserId {
			common.FmtError(c, `auth_no_permission`)
			return
		}

		// 如果是更新操作，使用数据库中的tool_key，不允许更改
		req.ToolKey = cast.ToString(toolInfo["tool_key"])
	}

	// 保存工具（创建或更新）
	id, err := common.SaveHttpTool(
		req.ID,
		adminUserId,
		req.Name,
		req.NameEn,
		req.ToolKey,
		req.Avatar,
		req.Description,
	)

	if err != nil {
		logs.Error("SaveHttpTool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteHttpTool 删除HTTP工具
func DeleteHttpTool(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查工具是否存在且属于当前用户
	toolInfo, err := common.GetHttpTool(id)
	if err != nil {
		logs.Error("Get http tool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(toolInfo) == 0 {
		common.FmtError(c, `http_tool_not_exist`)
		return
	}

	if cast.ToInt(toolInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 删除工具
	err = common.DeleteHttpTool(id)
	if err != nil {
		logs.Error("DeleteHttpTool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetHttpTool 获取单个HTTP工具
func GetHttpTool(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取工具信息及节点列表
	toolInfo, err := common.GetHttpToolWithNodes(id)
	if err != nil {
		logs.Error("Get http tool with nodes error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(toolInfo) == 0 {
		common.FmtError(c, `http_tool_not_exist`)
		return
	}

	// 检查权限
	if cast.ToInt(toolInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 如果toolInfo包含节点详情，需要将data_raw转换为node_params对象
	if nodes, ok := toolInfo["nodes"].([]map[string]any); ok {
		for j := range nodes {
			// 从data_raw解析NodeParams
			if dataRaw, ok := nodes[j]["data_raw"].(string); ok && dataRaw != "" {
				nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
				nodeParams.Curl.ToolInfo = GetHttpToolInfo(nodes[j])
				toolInfo["nodes"].([]map[string]any)[j]["http_tool_info"] = nodeParams.Curl.ToolInfo
				toolInfo["nodes"].([]map[string]any)[j]["node_params"] = nodeParams

			}
		}
	}

	// 返回数据
	common.FmtOk(c, toolInfo)
}

// GetHttpToolList 获取HTTP工具列表
func GetHttpToolList(c *gin.Context) {
	var req HttpToolListFilterRequest

	// 获取参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取工具列表，支持仅返回节点数量而不需要节点详情
	var list []map[string]any
	var total int
	var err error

	// 使用新的方法，支持仅获取节点数量
	list, total, err = common.GetHttpToolListWithFilterAndNodeCount(adminUserId, req.Name, req.Page, req.Size, req.WithNodes, true)
	if err != nil {
		logs.Error("GetHttpToolListWithFilterAndNodeCount error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// 如果需要获取节点详情，需要将data_raw转换为node_params对象
	if req.WithNodes {
		for i := range list {
			// 处理节点列表，先尝试不同的类型断言
			if nodesSlice, ok := list[i]["nodes"].([]map[string]any); ok {
				for j := range nodesSlice {
					// 从data_raw解析NodeParams
					if dataRaw, ok := nodesSlice[j]["data_raw"].(string); ok && dataRaw != "" {
						nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
						nodeParams.Curl.ToolInfo = GetHttpToolInfo(nodesSlice[j])
						nodesSlice[j]["http_tool_info"] = nodeParams.Curl.ToolInfo
						nodesSlice[j]["node_params"] = nodeParams
					}
				}
			}
		}
	}

	// 返回分页数据
	response := map[string]interface{}{
		"list":  list,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}

// GetHttpToolInfo 获取HTTP工具信息
func GetHttpToolInfo(nodes map[string]any) work_flow.HttpToolInfo {
	return work_flow.HttpToolInfo{
		HttpToolName:            cast.ToString(nodes["http_tool_name"]),
		HttpToolNameEn:          cast.ToString(nodes["http_tool_name_en"]),
		HttpToolKey:             cast.ToString(nodes["http_tool_key"]),
		HttpToolAvatar:          cast.ToString(nodes["http_tool_avatar"]),
		HttpToolDescription:     cast.ToString(nodes["http_tool_description"]),
		HttpToolNodeKey:         cast.ToString(nodes["node_key"]),
		HttpToolNodeName:        cast.ToString(nodes["node_name"]),
		HttpToolNodeNameEn:      cast.ToString(nodes["node_name_en"]),
		HttpToolNodeDescription: cast.ToString(nodes["node_description"]),
	}
}

// SaveHttpToolNode 保存HTTP工具节点（创建或更新）
func SaveHttpToolNode(c *gin.Context) {
	var req HttpToolNodeRequest

	// 获取参数
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查HTTP工具是否存在且属于当前用户
	httpToolInfo, err := common.GetHttpTool(req.HttpToolID)
	if err != nil {
		logs.Error("Get http tool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(httpToolInfo) == 0 {
		common.FmtError(c, `http_tool_not_exist`)
		return
	}

	if cast.ToInt(httpToolInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 如果是更新操作，检查节点是否存在且属于当前用户
	if req.ID > 0 {
		nodeInfo, err := common.GetHttpToolNode(req.ID)
		if err != nil {
			logs.Error("Get http tool node error: %s", err.Error())
			common.FmtError(c, `sys_err`)
			return
		}

		if len(nodeInfo) == 0 {
			common.FmtError(c, `http_tool_node_not_exist`)
			return
		}

		if cast.ToInt(nodeInfo["admin_user_id"]) != adminUserId {
			common.FmtError(c, `auth_no_permission`)
			return
		}
	}

	// 保存节点（创建或更新）- 包含node_params参数
	id, err := common.SaveHttpToolNode(
		req.ID,
		adminUserId,
		req.HttpToolID,
		req.NodeKey,
		req.NodeName,
		req.NodeNameEn,
		req.NodeDescription,
		req.NodeRemark,
		req.DataRaw,
	)

	if err != nil {
		logs.Error("SaveHttpToolNode error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, map[string]interface{}{"id": id})
}

// DeleteHttpToolNode 删除HTTP工具节点
func DeleteHttpToolNode(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查节点是否存在且属于当前用户
	nodeInfo, err := common.GetHttpToolNode(id)
	if err != nil {
		logs.Error("Get http tool node error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(nodeInfo) == 0 {
		common.FmtError(c, `http_tool_node_not_exist`)
		return
	}

	if cast.ToInt(nodeInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 删除节点
	err = common.DeleteHttpToolNode(id)
	if err != nil {
		logs.Error("DeleteHttpToolNode error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetHttpToolNode 获取单个HTTP工具节点
func GetHttpToolNode(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))

	// 检查参数
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 获取节点信息
	nodeInfo, err := common.GetHttpToolNode(id)
	if err != nil {
		logs.Error("Get http tool node error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(nodeInfo) == 0 {
		common.FmtError(c, `http_tool_node_not_exist`)
		return
	}

	// 检查权限
	if cast.ToInt(nodeInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 从data_raw解析NodeParams
	if dataRaw, ok := nodeInfo["data_raw"].(string); ok && dataRaw != "" {
		nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
		nodeParams.Curl.ToolInfo = GetHttpToolInfo(nodeInfo)
		nodeInfo["http_tool_info"] = nodeParams.Curl.ToolInfo
		nodeInfo["node_params"] = nodeParams

	}

	// 返回数据
	common.FmtOk(c, nodeInfo)
}

// GetHttpToolNodeList 获取HTTP工具节点列表
func GetHttpToolNodeList(c *gin.Context) {
	var req HttpToolNodeListFilterRequest

	// 获取参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取登录用户信息
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// 检查HTTP工具是否存在且属于当前用户
	httpToolInfo, err := common.GetHttpTool(req.HttpToolID)
	if err != nil {
		logs.Error("Get http tool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	if len(httpToolInfo) == 0 {
		common.FmtError(c, `http_tool_not_exist`)
		return
	}

	if cast.ToInt(httpToolInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// 获取节点列表
	list, total, err := common.GetHttpToolNodeListWithFilter(req.HttpToolID, req.NodeName, req.Page, req.Size)
	if err != nil {
		logs.Error("GetHttpToolNodeListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// 对每个节点，从data_raw解析NodeParams并替换node_params字段
	for i := range list {
		if dataRaw, ok := list[i]["data_raw"].(string); ok && dataRaw != "" {
			nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
			nodeParams.Curl.ToolInfo = GetHttpToolInfo(list[i])
			list[i]["http_tool_info"] = nodeParams.Curl.ToolInfo
			list[i]["node_params"] = nodeParams

		}
	}

	// 返回分页数据
	response := map[string]interface{}{
		"list":  list,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}
