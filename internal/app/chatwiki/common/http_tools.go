// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// HttpToolsCacheBuildHandler HTTP tool cache build handler
type HttpToolsCacheBuildHandler struct {
	ID int
}

func (h *HttpToolsCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki_http_tools_%d`, h.ID)
}

func (h *HttpToolsCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`http_tools`, define.Postgres).Where(`id`, cast.ToString(h.ID)).Find()
	if err != nil {
		return nil, err
	}
	// Convert
	if len(data) > 0 {
		result := map[string]interface{}{
			"id":            cast.ToInt(data[`id`]),
			"admin_user_id": cast.ToInt(data[`admin_user_id`]),
			"name":          data[`name`],
			"tool_key":      data[`tool_key`],
			"avatar":        data[`avatar`],
			"description":   data[`description`],
			"create_time":   cast.ToInt(data[`create_time`]),
			"update_time":   cast.ToInt(data[`update_time`]),
		}
		return result, err
	} else {
		return nil, nil
	}
}

// HttpToolsNodeCacheBuildHandler HTTP tool node cache build handler
type HttpToolsNodeCacheBuildHandler struct {
	ID int
}

func (h *HttpToolsNodeCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`http_tools_node_%d`, h.ID)
}

func (h *HttpToolsNodeCacheBuildHandler) GetCacheData() (any, error) {
	data, err := msql.Model(`http_tools_node`, define.Postgres).Where(`id`, cast.ToString(h.ID)).Find()
	if err != nil {
		return nil, err
	}
	// Convert
	if len(data) > 0 {
		result := map[string]interface{}{
			"id":               cast.ToInt(data[`id`]),
			"admin_user_id":    cast.ToInt(data[`admin_user_id`]),
			"http_tool_id":     cast.ToInt(data[`http_tool_id`]),
			"node_key":         data[`node_key`],
			"node_name":        data[`node_name`],
			"node_description": data[`node_description`],
			"node_remark":      data[`node_remark`],
			"data_raw":         data[`data_raw`],
			"create_time":      cast.ToInt(data[`create_time`]),
			"update_time":      cast.ToInt(data[`update_time`]),
		}
		return result, err
	} else {
		return nil, nil
	}
}

// GetHttpToolInfo gets HTTP tool info
func GetHttpToolInfo(id int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &HttpToolsCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// GetHttpToolNodeInfo gets HTTP tool node info
func GetHttpToolNodeInfo(id int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// SaveHttpTool saves HTTP tool (create or update)
func SaveHttpTool(id, adminUserID int, name, nameEn, toolKey, avatar, description string) (int64, error) {
	data := msql.Datas{
		"admin_user_id": adminUserID,
		"name":          name,
		"name_en":       nameEn,
		"avatar":        avatar,
		"description":   description,
		"update_time":   tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// Create new record
		if toolKey == `` {
			toolKey = GenerateToolKey(adminUserID)
		}
		data["tool_key"] = toolKey
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`http_tools`, define.Postgres).Insert(data, "id")
	} else {
		// Update existing record, but don't update tool_key
		_, err = msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &HttpToolsCacheBuildHandler{ID: int(newId)})
	}
	return newId, err
}

// DeleteHttpTool deletes HTTP tool
func DeleteHttpTool(id int) error {
	// First delete related nodes
	_, err := msql.Model(`http_tools_node`, define.Postgres).Where("http_tool_id", cast.ToString(id)).Delete()
	if err != nil {
		return err
	}

	// Delete main table record
	_, err = msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &HttpToolsCacheBuildHandler{ID: id})
	}
	return err
}

// GetHttpToolWithNodes gets single HTTP tool and its node information
func GetHttpToolWithNodes(id int) (map[string]any, error) {
	return GetSingleHttpTool(id, true)
}

// GetHttpTool gets single HTTP tool (compatible with older versions)
func GetHttpTool(id int) (msql.Params, error) {
	return msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Find()
}

// GenerateToolKey generates HTTP tool key
func GenerateToolKey(adminUserId int) string {
	// Generate node key HttpTool +uuid v7 + adminUserId%5
	key, _ := uuid.NewV7()                           // 1. Generate V7
	raw := strings.ReplaceAll(key.String(), "-", "") // 2. Remove hyphens
	nodeKey := fmt.Sprintf("HttpTool%s%d", raw, adminUserId%5)
	return nodeKey
}

// GenerateNodeKey generates HTTP tool node key
func GenerateNodeKey(adminUserId int) string {
	// Generate node key HttpTool +uuid v7 + adminUserId%5
	key, _ := uuid.NewV7()                           // 1. Generate V7
	raw := strings.ReplaceAll(key.String(), "-", "") // 2. Remove hyphens
	nodeKey := fmt.Sprintf("HttpToolNode%s%d", raw, adminUserId%5)
	return nodeKey
}

// GetHttpToolWithNodeCount gets single HTTP tool and its node count
func GetHttpToolWithNodeCount(id int) (map[string]any, error) {
	return GetSingleHttpTool(id, false)
}

// GetHttpToolListWithFilter gets HTTP tool list (with filters and pagination)
func GetHttpToolListWithFilter(adminUserID int, name string, page, size int, withNodes bool) ([]map[string]any, int, error) {
	return GetHttpToolListWithFilterAndNodeCount(adminUserID, name, page, size, withNodes, true)
}

// GetHttpToolListWithFilterAndNodeCount gets HTTP tool list (with filters and pagination), option to get node count
func GetHttpToolListWithFilterAndNodeCount(adminUserID int, name string, page, size int, withNodes, withNodeCount bool) ([]map[string]any, int, error) {
	model := msql.Model(`http_tools`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID))

	// Fuzzy query by tool name
	if name != "" {
		model.Where("name LIKE ?", "%"+name+"%")
	}

	// Add pagination
	list, total, err := model.Order("id DESC").Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}

	// Convert to []map[string]any format
	// Get http_tool_info
	toolInfoList := make(map[string]msql.Params)
	resultList := make([]map[string]any, len(list))
	for i, item := range list {
		id := cast.ToString(item["id"])
		// Unify return value format for nodes
		toolInfoList[id] = msql.Params{
			`name`:        cast.ToString(item["name"]),
			`avatar`:      cast.ToString(item["avatar"]),
			`tool_key`:    cast.ToString(item["tool_key"]),
			`description`: cast.ToString(item["description"]),
		}
		resultMap := make(map[string]any)
		for k, v := range item {
			resultMap[k] = v
		}
		resultList[i] = resultMap
	}

	if len(resultList) > 0 {
		// Get all main table IDs
		ids := make([]string, len(resultList))
		for i, item := range resultList {
			ids[i] = cast.ToString(item["id"])
		}

		// If node count is needed
		if withNodeCount {
			// Query all node counts for main IDs at once
			nodesList, err := msql.Model(`http_tools_node`, define.Postgres).
				Where("http_tool_id", "in", strings.Join(ids, ",")).
				Group("http_tool_id").
				Field("http_tool_id, COUNT(*) as count").
				Select()
			if err == nil {
				// Create mapping from ID to count
				countMap := make(map[string]int64)
				for _, node := range nodesList {
					httpToolId := cast.ToString(node["http_tool_id"])
					countMap[httpToolId] = cast.ToInt64(node["count"])
				}

				// Set node count for each result
				for i := range resultList {
					idStr := cast.ToString(resultList[i]["id"])
					if count, exists := countMap[idStr]; exists {
						resultList[i]["node_count"] = count
					} else {
						resultList[i]["node_count"] = int64(0)
					}
				}
			} else {
				// If query fails, set default values
				for i := range resultList {
					resultList[i]["node_count"] = int64(0)
				}
			}
		}

		// If node details are also needed
		if withNodes {
			// Query all node lists for main IDs at once
			httpToolNodes, err := msql.Model(`http_tools_node`, define.Postgres).
				Where("http_tool_id", "in", strings.Join(ids, ",")).
				Order("http_tool_id ASC, id ASC").
				Select()
			if err == nil {

				// Create mapping from ID to node list
				nodesMap := make(map[string][]map[string]any)
				for i := range httpToolNodes {
					httpToolId := cast.ToString(httpToolNodes[i]["http_tool_id"])
					// Unify return value format
					httpTool := toolInfoList[httpToolId]
					httpToolNodes[i] = NodeAddHttpToolInfo(httpToolNodes[i], httpTool)
					// Convert to []map[string]any format
					nodeInfo := make(map[string]any)
					for k, v := range httpToolNodes[i] {
						nodeInfo[k] = v
					}
					nodesMap[httpToolId] = append(nodesMap[httpToolId], nodeInfo)
				}

				// Set node list for each result
				for i := range resultList {
					idStr := cast.ToString(resultList[i]["id"])

					if nodes, exists := nodesMap[idStr]; exists {
						resultList[i]["nodes"] = nodes
					} else {
						resultList[i]["nodes"] = make([]map[string]any, 0)
					}
				}
			} else {
				// If query fails, set empty list
				for i := range resultList {
					resultList[i]["nodes"] = make([]map[string]any, 0)
				}
			}
		}
	}

	return resultList, total, nil
}

func NodeAddHttpToolInfo(node msql.Params, httpTool msql.Params) msql.Params {
	node["http_tool_name"] = cast.ToString(httpTool["name"])
	node["http_tool_name_en"] = cast.ToString(httpTool["name_en"])
	node["http_tool_avatar"] = cast.ToString(httpTool["avatar"])
	node["http_tool_key"] = cast.ToString(httpTool["tool_key"])
	node["http_tool_description"] = cast.ToString(httpTool["description"])
	return node
}

// GetSingleHttpTool gets single HTTP tool information (reuse list query logic)
func GetSingleHttpTool(id int, withNodes bool) (map[string]any, error) {
	// Directly query tool with specified ID
	httpTool, err := msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Find()
	if err != nil || len(httpTool) == 0 {
		return nil, err
	}

	result := make(map[string]any)
	for k, v := range httpTool {
		result[k] = v
	}

	// Add node count
	nodesCount, err := msql.Model(`http_tools_node`, define.Postgres).
		Where("http_tool_id", cast.ToString(id)).
		Count()
	if err == nil {
		result["node_count"] = nodesCount
	} else {
		result["node_count"] = int64(0)
	}

	// If node details are needed
	if withNodes {
		httpToolNodes, err := msql.Model(`http_tools_node`, define.Postgres).
			Where("http_tool_id", cast.ToString(id)).
			Order("id ASC").
			Select()
		if err != nil {
			return nil, err
		}
		nodes := make([]map[string]any, 0)
		for i := range httpToolNodes {
			// Unify return value format
			httpToolNodes[i] = NodeAddHttpToolInfo(httpToolNodes[i], httpTool)
			// Convert to []map[string]any format
			nodeInfo := make(map[string]any)
			for k, v := range httpToolNodes[i] {
				nodeInfo[k] = v
			}
			nodes = append(nodes, nodeInfo)
		}
		result["nodes"] = nodes
	}

	return result, nil
}

// SaveHttpToolNode saves HTTP tool node (create or update)
func SaveHttpToolNode(id, adminUserID, httpToolID int, nodeKey, nodeName, nodeNameEn, nodeDescription, nodeRemark, dataRaw string) (int64, error) {
	data := msql.Datas{
		"admin_user_id":    adminUserID,
		"http_tool_id":     httpToolID,
		"node_name":        nodeName,
		"node_name_en":     nodeNameEn,
		"node_description": nodeDescription,
		"node_remark":      nodeRemark,
		"data_raw":         dataRaw,
		"update_time":      tool.Time2Int(),
	}

	var err error
	var newId int64

	if id <= 0 {
		// Create new record
		if nodeKey == `` {
			nodeKey = GenerateNodeKey(adminUserID)
		}
		data["node_key"] = nodeKey
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`http_tools_node`, define.Postgres).Insert(data, "id")
	} else {
		// Update existing record
		_, err = msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: int(newId)})
	}
	return newId, err
}

// DeleteHttpToolNode deletes HTTP tool node
func DeleteHttpToolNode(id int) error {
	_, err := msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// Clear cache
		lib_redis.DelCacheData(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: id})
	}
	return err
}

// GetHttpToolNode gets single HTTP tool node
func GetHttpToolNode(id int) (map[string]any, error) {
	node, err := msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Find()
	if err != nil {
		return nil, err
	}

	// Associate and get main table information
	if len(node) > 0 {
		httpTool, err := msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(node["http_tool_id"])).Find()
		if err != nil {
			return nil, err
		}

		if len(httpTool) > 0 {
			node = NodeAddHttpToolInfo(node, httpTool)
		}
	}

	// Convert to []map[string]any format
	result := make(map[string]any)
	for k, v := range node {
		result[k] = v
	}

	return result, nil
}

// GetHttpToolNodeListWithFilter gets HTTP tool node list (with filters and pagination)
func GetHttpToolNodeListWithFilter(httpToolID int, nodeName string, page, size int) ([]map[string]any, int, error) {
	model := msql.Model(`http_tools_node`, define.Postgres).Where("http_tool_id", cast.ToString(httpToolID))

	// Fuzzy query by node name
	if nodeName != "" {
		model.Where("node_name LIKE ?", "%"+nodeName+"%")
	}

	// Add pagination
	list, total, err := model.Order("id DESC").Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}

	// Convert to []map[string]any format
	resultList := make([]map[string]any, len(list))
	for i, item := range list {
		resultMap := make(map[string]any)
		for k, v := range item {
			resultMap[k] = v
		}
		resultList[i] = resultMap
	}

	if len(resultList) > 0 {
		// Get all tool IDs associated with nodes
		toolIds := make([]string, len(resultList))
		for i, item := range resultList {
			toolIds[i] = cast.ToString(item["http_tool_id"])
		}

		// Query all tool information at once
		httpTools, err := msql.Model(`http_tools`, define.Postgres).
			Where("id", "in", strings.Join(toolIds, ",")).
			Select()
		if err == nil {
			// Create mapping from ID to tool information
			toolMap := make(map[string]msql.Params)
			for _, httpTool := range httpTools {
				toolMap[cast.ToString(httpTool["id"])] = httpTool
			}

			// Set tool information for each node
			for i := range resultList {
				toolId := cast.ToString(resultList[i]["http_tool_id"])
				if httpTool, exists := toolMap[toolId]; exists {
					resultList[i]["http_tool_name"] = httpTool["name"]
					resultList[i]["http_tool_name_en"] = httpTool["name_en"]
					resultList[i]["http_tool_key"] = httpTool["tool_key"]
					resultList[i]["http_tool_avatar"] = httpTool["avatar"]
					resultList[i]["http_tool_description"] = httpTool["description"]
				}
			}
		}
	}

	return resultList, total, nil
}
