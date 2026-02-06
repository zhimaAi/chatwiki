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

// HttpToolRequest HTTP tool request struct
type HttpToolRequest struct {
	ID          int    `form:"id" json:"id"`
	Name        string `form:"name" json:"name" binding:"required"`
	NameEn      string `form:"name_en" json:"name_en"`
	ToolKey     string `form:"tool_key" json:"tool_key"`
	Avatar      string `form:"avatar" json:"avatar"`
	Description string `form:"description" json:"description"`
}

// HttpToolNodeRequest HTTP tool node request struct
type HttpToolNodeRequest struct {
	ID              int    `form:"id" json:"id"`
	HttpToolID      int    `form:"http_tool_id" json:"http_tool_id" binding:"required"`
	NodeKey         string `form:"node_key" json:"node_key"`
	NodeNameEn      string `form:"node_name_en" json:"node_name_en"`
	NodeName        string `form:"node_name" json:"node_name" binding:"required"`
	NodeDescription string `form:"node_description" json:"node_description"`
	NodeRemark      string `form:"node_remark" json:"node_remark"`
	DataRaw         string `form:"data_raw" json:"data_raw"` // raw data, used to generate node_params
}

// HttpToolListFilterRequest HTTP tool list filter request struct
type HttpToolListFilterRequest struct {
	Name      string `json:"name" form:"name"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	WithNodes bool   `json:"with_nodes" form:"with_nodes"` // whether to include node details
}

// HttpToolNodeListFilterRequest HTTP tool node list filter request struct
type HttpToolNodeListFilterRequest struct {
	HttpToolID int    `json:"http_tool_id" form:"http_tool_id" binding:"required"`
	NodeName   string `json:"node_name" form:"node_name"`
	Page       int    `json:"page" form:"page"`
	Size       int    `json:"size" form:"size"`
}

// SaveHttpTool Save HTTP tool (create or update)
func SaveHttpTool(c *gin.Context) {
	var req HttpToolRequest

	// Get params
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// If updating, check whether the tool exists and belongs to current user
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

		// If updating, use tool_key from DB and forbid changes
		req.ToolKey = cast.ToString(toolInfo["tool_key"])
	}

	// Save tool (create or update)
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

// DeleteHttpTool Delete HTTP tool
func DeleteHttpTool(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the tool exists and belongs to current user
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

	// Delete tool
	err = common.DeleteHttpTool(id)
	if err != nil {
		logs.Error("DeleteHttpTool error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetHttpTool Get a single HTTP tool
func GetHttpTool(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get tool info and node list
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

	// Check permissions
	if cast.ToInt(toolInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// If toolInfo includes node details, convert data_raw to node_params
	if nodes, ok := toolInfo["nodes"].([]map[string]any); ok {
		for j := range nodes {
			// Parse NodeParams from data_raw
			if dataRaw, ok := nodes[j]["data_raw"].(string); ok && dataRaw != "" {
				nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
				nodeParams.Curl.ToolInfo = GetHttpToolInfo(nodes[j])
				toolInfo["nodes"].([]map[string]any)[j]["http_tool_info"] = nodeParams.Curl.ToolInfo
				toolInfo["nodes"].([]map[string]any)[j]["node_params"] = nodeParams

			}
		}
	}

	// Return data
	common.FmtOk(c, toolInfo)
}

// GetHttpToolList Get HTTP tool list
func GetHttpToolList(c *gin.Context) {
	var req HttpToolListFilterRequest

	// Get params
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get tool list; supports returning only node count without node details
	var list []map[string]any
	var total int
	var err error

	// Use new method; supports fetching only node count
	list, total, err = common.GetHttpToolListWithFilterAndNodeCount(adminUserId, req.Name, req.Page, req.Size, req.WithNodes, true)
	if err != nil {
		logs.Error("GetHttpToolListWithFilterAndNodeCount error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// If node details are required, convert data_raw to node_params
	if req.WithNodes {
		for i := range list {
			// Process node list; try different type assertions
			if nodesSlice, ok := list[i]["nodes"].([]map[string]any); ok {
				for j := range nodesSlice {
					// Parse NodeParams from data_raw
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

	// Return paginated data
	response := map[string]interface{}{
		"list":  list,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}

// GetHttpToolInfo Get HTTP tool info
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

// SaveHttpToolNode Save HTTP tool node (create or update)
func SaveHttpToolNode(c *gin.Context) {
	var req HttpToolNodeRequest

	// Get params
	if err := c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the HTTP tool exists and belongs to current user
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

	// If updating, check whether the node exists and belongs to current user
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

	// Save node (create or update) - includes node_params param
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

// DeleteHttpToolNode Delete HTTP tool node
func DeleteHttpToolNode(c *gin.Context) {
	id := cast.ToInt(c.PostForm("id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the node exists and belongs to current user
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

	// Delete node
	err = common.DeleteHttpToolNode(id)
	if err != nil {
		logs.Error("DeleteHttpToolNode error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	common.FmtOk(c, nil)
}

// GetHttpToolNode Get a single HTTP tool node
func GetHttpToolNode(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))

	// Validate params
	if id <= 0 {
		common.FmtError(c, `param_lack`, `id`)
		return
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Get node info
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

	// Check permissions
	if cast.ToInt(nodeInfo["admin_user_id"]) != adminUserId {
		common.FmtError(c, `auth_no_permission`)
		return
	}

	// Parse NodeParams from data_raw
	if dataRaw, ok := nodeInfo["data_raw"].(string); ok && dataRaw != "" {
		nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
		nodeParams.Curl.ToolInfo = GetHttpToolInfo(nodeInfo)
		nodeInfo["http_tool_info"] = nodeParams.Curl.ToolInfo
		nodeInfo["node_params"] = nodeParams

	}

	// Return data
	common.FmtOk(c, nodeInfo)
}

// GetHttpToolNodeList Get HTTP tool node list
func GetHttpToolNodeList(c *gin.Context) {
	var req HttpToolNodeListFilterRequest

	// Get params
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}

	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Get logged-in user info
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	// Check whether the HTTP tool exists and belongs to current user
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

	// Get node list
	list, total, err := common.GetHttpToolNodeListWithFilter(req.HttpToolID, req.NodeName, req.Page, req.Size)
	if err != nil {
		logs.Error("GetHttpToolNodeListWithFilter error: %s", err.Error())
		common.FmtError(c, `sys_err`)
		return
	}

	// For each node, parse NodeParams from data_raw and set node_params
	for i := range list {
		if dataRaw, ok := list[i]["data_raw"].(string); ok && dataRaw != "" {
			nodeParams := work_flow.DisposeNodeParams(work_flow.NodeTypeCurl, dataRaw, common.GetLang(c))
			nodeParams.Curl.ToolInfo = GetHttpToolInfo(list[i])
			list[i]["http_tool_info"] = nodeParams.Curl.ToolInfo
			list[i]["node_params"] = nodeParams

		}
	}

	// Return paginated data
	response := map[string]interface{}{
		"list":  list,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}

	common.FmtOk(c, response)
}
