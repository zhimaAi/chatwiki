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

// HttpToolsCacheBuildHandler HTTP工具缓存构建处理器
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
	// 转换
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

// HttpToolsNodeCacheBuildHandler HTTP工具节点缓存构建处理器
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
	// 转换
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

// GetHttpToolInfo 获取HTTP工具信息
func GetHttpToolInfo(id int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &HttpToolsCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// GetHttpToolNodeInfo 获取HTTP工具节点信息
func GetHttpToolNodeInfo(id int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := lib_redis.GetCacheWithBuild(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: id}, &result, time.Hour)
	return result, err
}

// SaveHttpTool 保存HTTP工具（创建或更新）
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
		// 创建新记录
		if toolKey == `` {
			toolKey = GenerateToolKey(adminUserID)
		}
		data["tool_key"] = toolKey
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`http_tools`, define.Postgres).Insert(data, "id")
	} else {
		// 更新现有记录，但不更新tool_key
		_, err = msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &HttpToolsCacheBuildHandler{ID: int(newId)})
	}
	return newId, err
}

// DeleteHttpTool 删除HTTP工具
func DeleteHttpTool(id int) error {
	// 先删除相关的节点
	_, err := msql.Model(`http_tools_node`, define.Postgres).Where("http_tool_id", cast.ToString(id)).Delete()
	if err != nil {
		return err
	}

	// 删除主表记录
	_, err = msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &HttpToolsCacheBuildHandler{ID: id})
	}
	return err
}

// GetHttpToolWithNodes 获取单个HTTP工具及其节点信息
func GetHttpToolWithNodes(id int) (map[string]any, error) {
	return GetSingleHttpTool(id, true)
}

// GetHttpTool 获取单个HTTP工具（兼容旧版）
func GetHttpTool(id int) (msql.Params, error) {
	return msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Find()
}

// GenerateToolKey 生成HTTP工具Key
func GenerateToolKey(adminUserId int) string {
	// 生成节点key HttpTool +uuid v7 + adminUserId%5
	key, _ := uuid.NewV7()                           // 1. 生成 V7
	raw := strings.ReplaceAll(key.String(), "-", "") // 2. 去连字符
	nodeKey := fmt.Sprintf("HttpTool%s%d", raw, adminUserId%5)
	return nodeKey
}

// GenerateNodeKey 生成HTTP工具节点Key
func GenerateNodeKey(adminUserId int) string {
	// 生成节点key HttpTool +uuid v7 + adminUserId%5
	key, _ := uuid.NewV7()                           // 1. 生成 V7
	raw := strings.ReplaceAll(key.String(), "-", "") // 2. 去连字符
	nodeKey := fmt.Sprintf("HttpToolNode%s%d", raw, adminUserId%5)
	return nodeKey
}

// GetHttpToolWithNodeCount 获取单个HTTP工具及其节点数量
func GetHttpToolWithNodeCount(id int) (map[string]any, error) {
	return GetSingleHttpTool(id, false)
}

// GetHttpToolListWithFilter 获取HTTP工具列表（带过滤条件和分页）
func GetHttpToolListWithFilter(adminUserID int, name string, page, size int, withNodes bool) ([]map[string]any, int, error) {
	return GetHttpToolListWithFilterAndNodeCount(adminUserID, name, page, size, withNodes, true)
}

// GetHttpToolListWithFilterAndNodeCount 获取HTTP工具列表（带过滤条件和分页），可选择是否获取节点数量
func GetHttpToolListWithFilterAndNodeCount(adminUserID int, name string, page, size int, withNodes, withNodeCount bool) ([]map[string]any, int, error) {
	model := msql.Model(`http_tools`, define.Postgres).Where("admin_user_id", cast.ToString(adminUserID))

	// 根据工具名称模糊查询
	if name != "" {
		model.Where("name LIKE ?", "%"+name+"%")
	}

	// 添加分页
	list, total, err := model.Order("id DESC").Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为[]map[string]any格式
	//获取http_tool_info
	toolInfoList := make(map[string]msql.Params)
	resultList := make([]map[string]any, len(list))
	for i, item := range list {
		id := cast.ToString(item["id"])
		//为节点统一 返回值格式使用
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
		// 获取所有主表ID
		ids := make([]string, len(resultList))
		for i, item := range resultList {
			ids[i] = cast.ToString(item["id"])
		}

		// 如果需要获取节点数量
		if withNodeCount {
			// 一次性查询所有主表ID对应的节点数量
			nodesList, err := msql.Model(`http_tools_node`, define.Postgres).
				Where("http_tool_id", "in", strings.Join(ids, ",")).
				Group("http_tool_id").
				Field("http_tool_id, COUNT(*) as count").
				Select()
			if err == nil {
				// 创建ID到数量的映射
				countMap := make(map[string]int64)
				for _, node := range nodesList {
					httpToolId := cast.ToString(node["http_tool_id"])
					countMap[httpToolId] = cast.ToInt64(node["count"])
				}

				// 为每个结果设置节点数量
				for i := range resultList {
					idStr := cast.ToString(resultList[i]["id"])
					if count, exists := countMap[idStr]; exists {
						resultList[i]["node_count"] = count
					} else {
						resultList[i]["node_count"] = int64(0)
					}
				}
			} else {
				// 如果查询失败，设置默认值
				for i := range resultList {
					resultList[i]["node_count"] = int64(0)
				}
			}
		}

		// 如果需要同时获取节点详情
		if withNodes {
			// 一次性查询所有主表ID对应的节点列表
			httpToolNodes, err := msql.Model(`http_tools_node`, define.Postgres).
				Where("http_tool_id", "in", strings.Join(ids, ",")).
				Order("http_tool_id ASC, id ASC").
				Select()
			if err == nil {

				// 创建ID到节点列表的映射
				nodesMap := make(map[string][]map[string]any)
				for i := range httpToolNodes {
					httpToolId := cast.ToString(httpToolNodes[i]["http_tool_id"])
					//统一返回值格式
					httpTool := toolInfoList[httpToolId]
					httpToolNodes[i] = NodeAddHttpToolInfo(httpToolNodes[i], httpTool)
					// 转换为[]map[string]any格式
					nodeInfo := make(map[string]any)
					for k, v := range httpToolNodes[i] {
						nodeInfo[k] = v
					}
					nodesMap[httpToolId] = append(nodesMap[httpToolId], nodeInfo)
				}

				// 为每个结果设置节点列表
				for i := range resultList {
					idStr := cast.ToString(resultList[i]["id"])

					if nodes, exists := nodesMap[idStr]; exists {
						resultList[i]["nodes"] = nodes
					} else {
						resultList[i]["nodes"] = make([]map[string]any, 0)
					}
				}
			} else {
				// 如果查询失败，设置空列表
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

// GetSingleHttpTool 获取单个HTTP工具信息（复用列表查询逻辑）
func GetSingleHttpTool(id int, withNodes bool) (map[string]any, error) {
	// 直接查询指定ID的工具
	httpTool, err := msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(id)).Find()
	if err != nil || len(httpTool) == 0 {
		return nil, err
	}

	result := make(map[string]any)
	for k, v := range httpTool {
		result[k] = v
	}

	// 添加节点数量
	nodesCount, err := msql.Model(`http_tools_node`, define.Postgres).
		Where("http_tool_id", cast.ToString(id)).
		Count()
	if err == nil {
		result["node_count"] = nodesCount
	} else {
		result["node_count"] = int64(0)
	}

	// 如果需要节点详情
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
			//统一返回值格式
			httpToolNodes[i] = NodeAddHttpToolInfo(httpToolNodes[i], httpTool)
			// 转换为[]map[string]any格式
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

// SaveHttpToolNode 保存HTTP工具节点（创建或更新）
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
		// 创建新记录
		if nodeKey == `` {
			nodeKey = GenerateNodeKey(adminUserID)
		}
		data["node_key"] = nodeKey
		data["create_time"] = tool.Time2Int()
		newId, err = msql.Model(`http_tools_node`, define.Postgres).Insert(data, "id")
	} else {
		// 更新现有记录
		_, err = msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Update(data)
		newId = int64(id)
	}

	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: int(newId)})
	}
	return newId, err
}

// DeleteHttpToolNode 删除HTTP工具节点
func DeleteHttpToolNode(id int) error {
	_, err := msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Delete()
	if err == nil {
		// 清除缓存
		lib_redis.DelCacheData(define.Redis, &HttpToolsNodeCacheBuildHandler{ID: id})
	}
	return err
}

// GetHttpToolNode 获取单个HTTP工具节点
func GetHttpToolNode(id int) (map[string]any, error) {
	node, err := msql.Model(`http_tools_node`, define.Postgres).Where("id", cast.ToString(id)).Find()
	if err != nil {
		return nil, err
	}

	// 关联获取主表信息
	if len(node) > 0 {
		httpTool, err := msql.Model(`http_tools`, define.Postgres).Where("id", cast.ToString(node["http_tool_id"])).Find()
		if err != nil {
			return nil, err
		}

		if len(httpTool) > 0 {
			node = NodeAddHttpToolInfo(node, httpTool)
		}
	}

	// 转换为[]map[string]any格式
	result := make(map[string]any)
	for k, v := range node {
		result[k] = v
	}

	return result, nil
}

// GetHttpToolNodeListWithFilter 获取HTTP工具节点列表（带过滤条件和分页）
func GetHttpToolNodeListWithFilter(httpToolID int, nodeName string, page, size int) ([]map[string]any, int, error) {
	model := msql.Model(`http_tools_node`, define.Postgres).Where("http_tool_id", cast.ToString(httpToolID))

	// 根据节点名称模糊查询
	if nodeName != "" {
		model.Where("node_name LIKE ?", "%"+nodeName+"%")
	}

	// 添加分页
	list, total, err := model.Order("id DESC").Paginate(page, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为[]map[string]any格式
	resultList := make([]map[string]any, len(list))
	for i, item := range list {
		resultMap := make(map[string]any)
		for k, v := range item {
			resultMap[k] = v
		}
		resultList[i] = resultMap
	}

	if len(resultList) > 0 {
		// 获取所有节点关联的工具ID
		toolIds := make([]string, len(resultList))
		for i, item := range resultList {
			toolIds[i] = cast.ToString(item["http_tool_id"])
		}

		// 一次性查询所有工具信息
		httpTools, err := msql.Model(`http_tools`, define.Postgres).
			Where("id", "in", strings.Join(toolIds, ",")).
			Select()
		if err == nil {
			// 创建ID到工具信息的映射
			toolMap := make(map[string]msql.Params)
			for _, httpTool := range httpTools {
				toolMap[cast.ToString(httpTool["id"])] = httpTool
			}

			// 为每个节点设置工具信息
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
