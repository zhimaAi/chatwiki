// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zhimaAi/go_tools/logs"
	"strings"
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
