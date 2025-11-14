// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/define"
	"net/http"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

// MCP Server 缓存
var (
	mcpServerCache     = make(map[int]*CachedMCPServer)
	mcpServerCacheLock sync.RWMutex
)

// CachedMCPServer MCP 服务器缓存结构
type CachedMCPServer struct {
	Server      *server.MCPServer
	HTTPHandler http.Handler
	CreatedAt   time.Time
	LastUsedAt  time.Time
	AdminUserID int
	ServerID    int
	ToolsHash   string
}

// GetMCPServerCache 获取缓存的 MCP 服务器
func GetMCPServerCache(serverId int) (*CachedMCPServer, bool) {
	mcpServerCacheLock.RLock()
	defer mcpServerCacheLock.RUnlock()
	cached, exists := mcpServerCache[serverId]
	return cached, exists
}

// SetMCPServerCache 设置 MCP 服务器缓存
func SetMCPServerCache(serverId int, cached *CachedMCPServer) {
	mcpServerCacheLock.Lock()
	defer mcpServerCacheLock.Unlock()
	mcpServerCache[serverId] = cached
}

// UpdateMCPServerLastUsed 更新最后使用时间
func UpdateMCPServerLastUsed(serverId int) {
	mcpServerCacheLock.Lock()
	defer mcpServerCacheLock.Unlock()
	if cached, exists := mcpServerCache[serverId]; exists {
		cached.LastUsedAt = time.Now()
	}
}

// ClearMCPServerCache 清除 MCP Server 缓存（当工具配置变更时调用）
func ClearMCPServerCache(adminUserId int) {
	// Query server_id by admin_user_id
	mcpServerInfo, err := msql.Model(`mcp_server`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil {
		logs.Error("query mcp_server error: %v", err)
		return
	}
	if len(mcpServerInfo) == 0 {
		logs.Info("No mcp_server found for admin_user_id: %d", adminUserId)
		return
	}

	serverId := cast.ToInt(mcpServerInfo[`id`])

	mcpServerCacheLock.Lock()
	defer mcpServerCacheLock.Unlock()

	if _, exists := mcpServerCache[serverId]; exists {
		delete(mcpServerCache, serverId)
		logs.Info("Cleared MCP Server cache for admin_user_id: %d, server_id: %d", adminUserId, serverId)
	}
}
