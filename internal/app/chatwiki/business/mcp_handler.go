// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

// MCPAuthMiddleware validates Bearer Token and loads user info
func MCPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// Query server info by api_key (only published servers allowed)
		mcpServerInfo, err := msql.Model(`mcp_server`, define.Postgres).
			Where(`api_key`, apiKey).
			Where(`publish_status`, cast.ToString(define.McpServerPublished)).
			Find()
		if err != nil {
			logs.Error("query mcp_server error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			c.Abort()
			return
		}

		if len(mcpServerInfo) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid api key",
			})
			c.Abort()
			return
		}

		// Store user info in context
		adminUserId := cast.ToInt(mcpServerInfo[`admin_user_id`])
		serverId := cast.ToInt(mcpServerInfo[`id`])
		serverName := cast.ToString(mcpServerInfo[`name`])
		serverDescription := cast.ToString(mcpServerInfo[`description`])

		c.Set("mcp_admin_user_id", adminUserId)
		c.Set("mcp_server_id", serverId)
		c.Set("mcp_server_name", serverName)
		c.Set("mcp_server_description", serverDescription)

		c.Next()
	}
}

// HandleMCPRequest handles MCP requests with cached server instances
func HandleMCPRequest(c *gin.Context) {
	adminUserId, exists := c.Get("mcp_admin_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	serverId := cast.ToInt(c.MustGet("mcp_server_id"))
	serverName := cast.ToString(c.MustGet("mcp_server_name"))
	serverDescription := cast.ToString(c.MustGet("mcp_server_description"))

	cached, err := getOrCreateMCPServer(
		cast.ToInt(adminUserId),
		serverId,
		serverName,
		serverDescription,
	)
	if err != nil {
		logs.Error("get or create mcp server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to initialize mcp server",
		})
		return
	}

	common.UpdateMCPServerLastUsed(serverId)

	cached.HTTPHandler.ServeHTTP(c.Writer, c.Request)
}

// getOrCreateMCPServer retrieves or creates an MCP Server instance
func getOrCreateMCPServer(adminUserId, serverId int, serverName, serverDescription string) (*common.CachedMCPServer, error) {
	cached, exists := common.GetMCPServerCache(serverId)

	tools, toolsHash, err := loadUserToolsWithHash(adminUserId, serverId)
	if err != nil {
		return nil, fmt.Errorf("load user tools error: %w", err)
	}

	// Return cached instance if exists and tools config unchanged
	if exists && cached.ToolsHash == toolsHash {
		return cached, nil
	}

	logs.Info("Creating MCP Server for user %d, server %d (tools hash: %s)", adminUserId, serverId, toolsHash)

	mcpServer := server.NewMCPServer(
		serverName,
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	for _, tool := range tools {
		mcpServer.AddTool(tool.Definition, tool.Handler)
		logs.Info("Registered tool: %s for user: %d", tool.Definition.Name, adminUserId)
	}

	httpHandler := server.NewStreamableHTTPServer(mcpServer)

	newCached := &common.CachedMCPServer{
		Server:      mcpServer,
		HTTPHandler: httpHandler,
		CreatedAt:   time.Now(),
		LastUsedAt:  time.Now(),
		AdminUserID: adminUserId,
		ServerID:    serverId,
		ToolsHash:   toolsHash,
	}

	common.SetMCPServerCache(serverId, newCached)

	logs.Info("MCP Server cached for server_id: %d, tools count: %d", serverId, len(tools))

	return newCached, nil
}

// ToolInfo contains tool definition and handler
type ToolInfo struct {
	Definition mcp.Tool
	Handler    func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// loadUserToolsWithHash loads user tools and computes config hash
func loadUserToolsWithHash(adminUserId int, serverId int) ([]ToolInfo, string, error) {
	tools, err := loadUserTools(adminUserId, serverId)
	if err != nil {
		return nil, "", err
	}

	toolsHash := fmt.Sprintf("%d", time.Now().Unix())
	if len(tools) > 0 {
		var names []string
		for _, t := range tools {
			names = append(names, t.Definition.Name)
		}
		toolsHash = fmt.Sprintf("%v", names)
	}

	return tools, toolsHash, nil
}

// loadUserTools loads tools from database for a user
func loadUserTools(adminUserId int, serverId int) ([]ToolInfo, error) {
	toolList, err := msql.Model(`mcp_tool`, define.Postgres).
		Alias(`m`).
		Where(`m.admin_user_id`, cast.ToString(adminUserId)).
		Where(`server_id`, cast.ToString(serverId)).
		Join(`chat_ai_robot r`, `m.robot_id = r.id`, `inner`).
		Join(`work_flow_node w`, `r.id = w.robot_id and r.start_node_Key = w.node_key and w.data_type = `+cast.ToString(define.DataTypeRelease), `left`).
		Field(`m.*,r.id as robot_id,r.robot_key,r.robot_name,r.robot_intro,w.node_params`).
		Select()
	if err != nil {
		return nil, fmt.Errorf("query mcp_tool error: %w", err)
	}

	var tools []ToolInfo

	for _, toolData := range toolList {
		nodeParamsStr := cast.ToString(toolData[`node_params`])

		// Build tool definition with dynamic parameters
		toolOptions := []mcp.ToolOption{
			mcp.WithDescription(cast.ToString(toolData[`robot_intro`])),
		}
		toolOptions = append(toolOptions, mcp.WithString(`content`, mcp.Description(lib_define.DialogueContent), mcp.Required()))
		toolOptions = append(toolOptions, mcp.WithString(`open_id`, mcp.Description(lib_define.UserRequestIdentifier), mcp.Required()))

		if toolData[`node_params`] != "" {
			// Parse node_params JSON using work_flow.NodeParams
			var nodeParams work_flow.NodeParams
			var allParams []work_flow.StartNodeParam
			if err := json.Unmarshal([]byte(nodeParamsStr), &nodeParams); err != nil {
				logs.Error("Failed to parse node_params for tool %s: %v", toolData[`name`], err)
			} else {
				allParams = append(allParams, nodeParams.Start.DiyGlobal...)
			}

			// Add parameters from node_params
			for _, param := range allParams {
				paramOptions := []mcp.PropertyOption{
					mcp.Description(param.Desc),
				}
				if param.Required {
					paramOptions = append(paramOptions, mcp.Required())
				}

				// Add parameter based on type
				switch strings.ToLower(param.Typ) {
				case "string":
					toolOptions = append(toolOptions, mcp.WithString(param.Key, paramOptions...))
				case "number", "int", "integer":
					toolOptions = append(toolOptions, mcp.WithNumber(param.Key, paramOptions...))
				case "boolean", "bool":
					toolOptions = append(toolOptions, mcp.WithBoolean(param.Key, paramOptions...))
				default:
					// Default to string for unknown types
					toolOptions = append(toolOptions, mcp.WithString(param.Key, paramOptions...))
				}
			}
		}

		toolDef := mcp.NewTool(cast.ToString(toolData[`name`]), toolOptions...)

		handler := chatMessage(
			adminUserId,
			cast.ToString(toolData[`robot_key`]),
		)

		tools = append(tools, ToolInfo{
			Definition: toolDef,
			Handler:    handler,
		})
	}

	return tools, nil
}

func chatMessage(adminUserId int, robotKey string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		args := req.GetArguments()

		content, ok := args["content"].(string)
		if !ok || content == "" {
			return nil, fmt.Errorf("content parameter is required")
		}

		openID, ok := args["open_id"].(string)
		if !ok || openID == "" {
			return nil, fmt.Errorf("open_id parameter is required")
		}

		// Extract optional global parameters
		var global map[string]any
		if globalVal, exists := args["global"]; exists {
			if globalMap, ok := globalVal.(map[string]any); ok {
				global = globalMap
			}
		}

		// Validate robot key
		if !common.CheckRobotKey(robotKey) {
			return nil, fmt.Errorf("invalid robot_key")
		}

		// Get robot info
		robot, err := common.GetRobotInfo(robotKey)
		if err != nil {
			logs.Error("get robot info error: %v", err)
			return nil, fmt.Errorf("failed to get robot info")
		}
		if len(robot) == 0 {
			return nil, fmt.Errorf("robot not found")
		}

		// Verify admin user id
		if cast.ToInt(robot["admin_user_id"]) != adminUserId {
			return nil, fmt.Errorf("admin_user_id mismatch")
		}

		// Validate open id
		if !common.IsChatOpenid(openID) {
			return nil, fmt.Errorf("invalid open_id")
		}

		// Get customer info
		customer, err := common.GetCustomerInfo(openID, adminUserId)
		if err != nil {
			logs.Error("get customer info error: %v", err)
			return nil, fmt.Errorf("failed to get customer info")
		}

		// Build chat request parameters
		chatBaseParam := &define.ChatBaseParam{
			AppType:     lib_define.AppOpenApi,
			Openid:      openID,
			AdminUserId: adminUserId,
			Robot:       robot,
			Customer:    customer,
		}

		isClose := false
		chatParams := &define.ChatRequestParam{
			ChatBaseParam:  chatBaseParam,
			Lang:           define.LangZhCn,
			Question:       strings.TrimSpace(content),
			IsClose:        &isClose,
			WorkFlowGlobal: global,
		}

		// Execute non-streaming chat request
		chanStream := make(chan sse.Event)
		go func() {
			for range chanStream {
				// Consume stream events to prevent blocking
			}
		}()

		message, err := DoChatRequest(chatParams, false, chanStream)
		if err != nil {
			logs.Error("chat request error: %v", err)
			return nil, fmt.Errorf("chat request failed: %w", err)
		}

		return mcp.NewToolResultText(message[`content`]), nil
	}
}
