// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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

// MCP Server cache
var (
	mcpServerCache     = make(map[int]*CachedMCPServer)
	mcpServerCacheLock sync.RWMutex
)

// CachedMCPServer is the cached MCP server struct
type CachedMCPServer struct {
	Server      *server.MCPServer
	HTTPHandler http.Handler
	CreatedAt   time.Time
	LastUsedAt  time.Time
	AdminUserID int
	ServerID    int
	ToolsHash   string
}

// GetMCPServerCache gets a cached MCP server
func GetMCPServerCache(serverId int) (*CachedMCPServer, bool) {
	mcpServerCacheLock.RLock()
	defer mcpServerCacheLock.RUnlock()
	cached, exists := mcpServerCache[serverId]
	return cached, exists
}

// SetMCPServerCache sets a cached MCP server
func SetMCPServerCache(serverId int, cached *CachedMCPServer) {
	mcpServerCacheLock.Lock()
	defer mcpServerCacheLock.Unlock()
	mcpServerCache[serverId] = cached
}

// UpdateMCPServerLastUsed updates the last used time
func UpdateMCPServerLastUsed(serverId int) {
	mcpServerCacheLock.Lock()
	defer mcpServerCacheLock.Unlock()
	if cached, exists := mcpServerCache[serverId]; exists {
		cached.LastUsedAt = time.Now()
	}
}

// ClearMCPServerCache clears MCP Server cache (called when tool config changes)
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
