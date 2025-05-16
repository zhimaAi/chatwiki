// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zhimaAi/go_tools/logs"
)

// GraphDB 图数据库操作封装
type GraphDB struct {
	adminUserId int
}

// NewGraphDB 创建图数据库操作实例
func NewGraphDB(adminUserId int) *GraphDB {
	return &GraphDB{
		adminUserId: adminUserId,
	}
}

// Execute 执行Cypher查询
func (g *GraphDB) Execute(query string) (*neo4j.EagerResult, error) {
	logs.Debug("execute neo4j query: %s", query)
	ctx := context.Background()
	result, err := neo4j.ExecuteQuery(
		ctx,
		define.Neo4jDriver,
		query,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(define.Config.Neo4j["database"]),
	)
	return result, err
}

func (g *GraphDB) ConstructEntity(subject, object, sanitizedPredicate string, libraryId, fileId, dataId int, confidence float64) (*neo4j.EagerResult, error) {
	createGraphSQL := fmt.Sprintf(`
		MERGE (s:Entity_%d {name: '%s', library_id: %d, file_id: %d, data_id: %d})
		MERGE (o:Entity_%d {name: '%s', library_id: %d, file_id: %d, data_id: %d})
		CREATE (s)-[r:%s {confidence: %f, library_id: %d, file_id: %d, data_id: %d}]->(o)
	`, g.adminUserId, subject, libraryId, fileId, dataId,
		g.adminUserId, object, libraryId, fileId, dataId,
		sanitizedPredicate, confidence, libraryId, fileId, dataId)
	r, err := g.Execute(createGraphSQL)
	if err != nil {
		logs.Error(`create graph error: %s, cypher is %s`, err.Error(), createGraphSQL)
	}
	return r, err
}

func (g *GraphDB) GetFileEntity(fileId int) (*neo4j.EagerResult, error) {
	query := fmt.Sprintf(`
			MATCH (s:Entity_%d {file_id: %d})
			RETURN s.file_id as file_id, s.data_id as data_id, s.name as name, s.library_id as library_id
		`, g.adminUserId, fileId)
	return g.Execute(query)
}

func (g *GraphDB) GetEntityCount(idList []string) (*neo4j.EagerResult, error) {
	query := fmt.Sprintf(`
			MATCH (s:Entity_%d)
			WHERE s.file_id in [%s]
			RETURN s.file_id as file_id, count(s) as count
		`, g.adminUserId, strings.Join(idList, `,`))
	return g.Execute(query)
}

// DeleteByLibrary 删除库相关的所有图数据
func (g *GraphDB) DeleteByLibrary(libraryId int) error {
	query := fmt.Sprintf(`
		MATCH (n:Entity_%d {library_id: %d})
		DETACH DELETE n
	`, g.adminUserId, libraryId)
	_, err := g.Execute(query)
	return err
}

// DeleteByFile 删除文件相关的所有图数据
func (g *GraphDB) DeleteByFile(fileId int) error {
	query := fmt.Sprintf(`
		MATCH (n:Entity_%d {file_id: %d})
		DETACH DELETE n
	`, g.adminUserId, fileId)
	_, err := g.Execute(query)
	return err
}

// DeleteByData 删除数据相关的所有图数据
func (g *GraphDB) DeleteByData(dataId int) error {
	query := fmt.Sprintf(`
		MATCH (n:Entity_%d {data_id: %d})
		DETACH DELETE n
	`, g.adminUserId, dataId)
	_, err := g.Execute(query)
	return err
}

// FindRelatedEntities Find related entities in the graph database, up to 3 levels deep
func (g *GraphDB) FindRelatedEntities(entity string, libraryIds []string, limit int, maxDepth int) (*neo4j.EagerResult, error) {
	entity = strings.ReplaceAll(entity, "'", "''")
	inCondition := ""
	for i, id := range libraryIds {
		if i > 0 {
			inCondition += ","
		}
		inCondition += id
	}
	query := fmt.Sprintf(`
		MATCH path = (start:Entity_%d)-[*1..%d]-(connected:Entity_%d)
		WHERE start.name contains '%s' and start.library_id IN [%s]
			WITH path, length(path) AS depth
			UNWIND relationships(path) AS rel
			WITH rel, startNode(rel) AS start_node, endNode(rel) AS end_node, depth
		RETURN start_node.name AS subject, 
			type(rel) AS relation, 
			end_node.name AS object,
			rel.confidence AS confidence, 
			rel.library_id AS library_id,
			rel.file_id AS file_id, 
			rel.data_id AS data_id,
		depth AS depth
		LIMIT %d
	`, g.adminUserId, maxDepth, g.adminUserId, entity, inCondition, limit)

	results, err := g.Execute(query)
	if err != nil {
		logs.Error("Graph query failed: %s", err.Error())
		return nil, err
	}
	return results, nil
}

// GetFileRelationships 获取指定文件的所有关系边
func (g *GraphDB) GetFileRelationships(fileId int, dataId int, searchTerm ...string) (*neo4j.EagerResult, error) {
	query := fmt.Sprintf(`
		MATCH (s:Entity_%d {file_id: %d})-[r]-(o:Entity_%d {file_id: %d})
		RETURN 
			id(r) AS id,
			id(s) AS from_id,
			id(o) AS to_id,
			type(r) AS label,
			type(r) AS type
	`, g.adminUserId, fileId, g.adminUserId, fileId)

	if dataId > 0 {
		query = fmt.Sprintf(`
		MATCH (s:Entity_%d {file_id: %d, data_id: %d})-[r]-(o:Entity_%d {file_id: %d, data_id: %d})
		RETURN 
			id(r) AS id,
			id(s) AS from_id,
			id(o) AS to_id,
			type(r) AS label,
			type(r) AS type
	`, g.adminUserId, fileId, dataId, g.adminUserId, fileId, dataId)
	}

	// 如果指定了搜索词，尝试作为关系类型或节点名称过滤
	if len(searchTerm) > 0 && searchTerm[0] != "" {
		// 转义单引号
		escapedTerm := strings.ReplaceAll(searchTerm[0], "'", "''")

		// 修改查询，同时过滤关系类型和节点名称
		query = fmt.Sprintf(`
		MATCH (s:Entity_%d {file_id: %d})-[r]-(o:Entity_%d {file_id: %d})
		WHERE type(r) CONTAINS '%s' OR s.name CONTAINS '%s' OR o.name CONTAINS '%s'
		RETURN 
			id(r) AS id,
			id(s) AS from_id,
			id(o) AS to_id,
			type(r) AS label,
			type(r) AS type
	`, g.adminUserId, fileId, g.adminUserId, fileId, escapedTerm, escapedTerm, escapedTerm)

		if dataId > 0 {
			query = fmt.Sprintf(`
			MATCH (s:Entity_%d {file_id: %d, data_id: %d})-[r]-(o:Entity_%d {file_id: %d, data_id: %d})
			WHERE type(r) CONTAINS '%s' OR s.name CONTAINS '%s' OR o.name CONTAINS '%s'
			RETURN 
				id(r) AS id,
				id(s) AS from_id,
				id(o) AS to_id,
				type(r) AS label,
				type(r) AS type
		`, g.adminUserId, fileId, dataId, g.adminUserId, fileId, dataId, escapedTerm, escapedTerm, escapedTerm)
		}
	}

	return g.Execute(query)
}

// GetFileNodes 获取文件中所有唯一节点
func (g *GraphDB) GetFileNodes(fileId int, dataId int, searchTerm ...string) (*neo4j.EagerResult, error) {
	query := fmt.Sprintf(`
		MATCH (n:Entity_%d {file_id: %d})
		RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
		limit 500
	`, g.adminUserId, fileId)

	if dataId > 0 {
		query = fmt.Sprintf(`
		MATCH (n:Entity_%d {file_id: %d, data_id: %d})
		RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
		limit 500
	`, g.adminUserId, fileId, dataId)
	}

	// 如果指定了搜索词，添加过滤条件
	if len(searchTerm) > 0 && searchTerm[0] != "" {
		// 转义单引号
		escapedTerm := strings.ReplaceAll(searchTerm[0], "'", "''")

		if dataId > 0 {
			query = fmt.Sprintf(`
			MATCH (n:Entity_%d {file_id: %d, data_id: %d})
			WHERE n.name CONTAINS '%s'
			RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
			UNION
			MATCH (n:Entity_%d {file_id: %d, data_id: %d})-[r]-(m:Entity_%d {file_id: %d, data_id: %d})
			WHERE type(r) CONTAINS '%s' OR m.name CONTAINS '%s'
			limit 500
			RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
		`, g.adminUserId, fileId, dataId, escapedTerm,
				g.adminUserId, fileId, dataId, g.adminUserId, fileId, dataId, escapedTerm, escapedTerm)
		} else {
			query = fmt.Sprintf(`
			MATCH (n:Entity_%d {file_id: %d})
			WHERE n.name CONTAINS '%s'
			RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
			UNION
			MATCH (n:Entity_%d {file_id: %d})-[r]-(m:Entity_%d {file_id: %d})
			WHERE type(r) CONTAINS '%s' OR m.name CONTAINS '%s'
			limit 500
			RETURN DISTINCT id(n) AS id, n.name AS name, n.data_id as data_id
		`, g.adminUserId, fileId, escapedTerm,
				g.adminUserId, fileId, g.adminUserId, fileId, escapedTerm, escapedTerm)
		}
	}

	return g.Execute(query)
}
