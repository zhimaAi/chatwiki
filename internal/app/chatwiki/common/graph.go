// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// GraphDB 图数据库操作封装
type GraphDB struct {
	schema string // 图数据库schema名称
}

// NewGraphDB 创建图数据库操作实例
func NewGraphDB(schema string) *GraphDB {
	return &GraphDB{
		schema: schema,
	}
}

// ExecuteCypher 执行Cypher查询
func (g *GraphDB) ExecuteCypher(query string, args ...interface{}) ([]msql.Params, error) {
	// 构建完整的查询
	cypherQuery := fmt.Sprintf(query, args...)

	sql := `LOAD 'age'; SET search_path = ag_catalog, "$user", public;` + msql.Model(cypherQuery, define.Postgres).Field("*").BuildSql()

	return msql.RawValues(define.Postgres, sql, nil)
}

// CreateTriple 创建三元组
func (g *GraphDB) CreateTriple(subject, predicate, object string, attributes map[string]interface{}) error {
	// 构建属性字符串
	attrStr := ""
	if len(attributes) > 0 {
		attrs := make([]string, 0, len(attributes))
		for k, v := range attributes {
			switch v := v.(type) {
			case string:
				attrs = append(attrs, fmt.Sprintf("%s: '%s'", k, strings.ReplaceAll(v, "'", "''")))
			case time.Time:
				attrs = append(attrs, fmt.Sprintf("%s: %d", k, tool.Time2Int()))
			default:
				attrs = append(attrs, fmt.Sprintf("%s: %v", k, v))
			}
		}
		attrStr = ", " + strings.Join(attrs, ", ")
	}

	query := fmt.Sprintf(`
		cypher('%s', $$
			MERGE (s:Entity {name: '%s'})
			MERGE (o:Entity {name: '%s'})
			CREATE (s)-[r:%s {%s}]->(o)
			RETURN r
		$$) as (r agtype)
	`, g.schema,
		strings.ReplaceAll(subject, "'", "''"),
		strings.ReplaceAll(object, "'", "''"),
		predicate,
		strings.TrimPrefix(attrStr, ", "))

	_, err := g.ExecuteCypher(query)
	return err
}

// FindByEntity 根据实体查找关联三元组
func (g *GraphDB) FindByEntity(subject, relation, object string, libraryIds []string, limit int) ([]msql.Params, error) {
	// 构建库ID条件
	libraryCondition := ""
	if len(libraryIds) > 0 {
		libraryConditions := make([]string, 0)
		for _, id := range libraryIds {
			libraryConditions = append(libraryConditions, fmt.Sprintf("r.library_id = %s", id))
		}
		if len(libraryConditions) > 0 {
			libraryCondition = "AND (" + strings.Join(libraryConditions, " OR ") + ")"
		}
	}

	query := fmt.Sprintf(`
		cypher('%s', $$
			MATCH (s:Entity)-[r]->(o:Entity)
			WHERE (s.name =~ '(?i).*%s.*') %s
			WITH s, r, o,
         		CASE
					WHEN type(r) =~ '(?i).*%s.*' THEN 3
             		WHEN o.name =~ '(?i).*%s.*' THEN 2
             		ELSE 1
         		END * r.confidence as relevance
			RETURN s.name as subject, type(r) as relation, o.name as object, 
				   r.confidence as confidence, r.library_id as library_id, 
				   r.file_id as file_id, r.data_id as data_id, relevance, 
				   1 as depth
			ORDER BY relevance DESC
			LIMIT %d
		$$) as (subject agtype, relation agtype, object agtype, confidence agtype, 
			   library_id agtype, file_id agtype, data_id agtype, relevance agtype,
			   depth agtype)
	`, g.schema, subject, libraryCondition, relation, object, limit)

	return g.ExecuteCypher(query)
}

// DeleteByLibrary 删除库相关的所有图数据
func (g *GraphDB) DeleteByLibrary(libraryId int) error {
	query := fmt.Sprintf(`
		cypher('%s', $$
			MATCH (n)
			WHERE n.library_id = %d
			DETACH DELETE n
			RETURN count(*) as deleted
		$$) as (deleted agtype)
	`, g.schema, libraryId)

	_, err := g.ExecuteCypher(query)
	return err
}

// DeleteByFile 删除文件相关的所有图数据
func (g *GraphDB) DeleteByFile(fileId int) error {
	query := fmt.Sprintf(`
		cypher('%s', $$
			MATCH (n)
			WHERE n.file_id = %d
			DETACH DELETE n
			RETURN count(*) as deleted
		$$) as (deleted agtype)
	`, g.schema, fileId)

	_, err := g.ExecuteCypher(query)
	return err
}

// DeleteByData 删除数据相关的所有图数据
func (g *GraphDB) DeleteByData(dataId int) error {
	query := fmt.Sprintf(`
		cypher('%s', $$
			MATCH (n)
			WHERE n.data_id = %d
			DETACH DELETE n
			RETURN count(*) as deleted
		$$) as (deleted agtype)
	`, g.schema, dataId)

	_, err := g.ExecuteCypher(query)
	return err
}

// FindRelatedEntities Find related triples in the graph database for an entity, supporting multi-hop queries
func (g *GraphDB) FindRelatedEntities(entity string, libraryIds []string, limit int, maxDepth int) ([]msql.Params, error) {
	// For safety, limit maximum recursion depth to 2
	if maxDepth > 2 {
		maxDepth = 2
	}

	// Process entity name, remove special characters
	entity = strings.ReplaceAll(entity, "'", "''")

	// Store all results
	allResults := make([]msql.Params, 0)
	visitedEntities := make(map[string]bool)
	visitedEntities[entity] = true

	// Build library ID conditions
	libraryCondition := ""
	if len(libraryIds) > 0 {
		libraryConditions := make([]string, 0)
		for _, id := range libraryIds {
			libraryConditions = append(libraryConditions, fmt.Sprintf("r.library_id = %s", id))
		}
		if len(libraryConditions) > 0 {
			libraryCondition = "AND (" + strings.Join(libraryConditions, " OR ") + ")"
		}
	}

	// First level query - combined subject and object query (exact match)
	firstLevelQuery := fmt.Sprintf(`
		cypher('%s', $$
			MATCH (n1:Entity)-[r]->(n2:Entity)
			WHERE (n1.name = '%s' OR n2.name = '%s') %s
			RETURN n1.name as subject, type(r) as relation, n2.name as object, 
				   r.confidence as confidence, r.library_id as library_id, 
				   r.file_id as file_id, r.data_id as data_id,
				   1 as depth
			LIMIT %d
		$$) as (subject agtype, relation agtype, object agtype, confidence agtype, 
			   library_id agtype, file_id agtype, data_id agtype, depth agtype)
	`, g.schema, entity, entity, libraryCondition, limit)

	// Execute first level query
	firstLevelResults, err := g.ExecuteCypher(firstLevelQuery)
	if err != nil {
		logs.Error("First level graph query failed: %s", err.Error())
	} else {
		allResults = append(allResults, firstLevelResults...)
	}

	// Collect entities for second level query
	if maxDepth > 1 {
		secondLevelEntities := make([]string, 0)

		// Extract entities from first level results
		for _, result := range firstLevelResults {
			subject := cast.ToString(result["subject"])
			object := cast.ToString(result["object"])

			if !visitedEntities[subject] {
				visitedEntities[subject] = true
				secondLevelEntities = append(secondLevelEntities, subject)
			}

			if !visitedEntities[object] {
				visitedEntities[object] = true
				secondLevelEntities = append(secondLevelEntities, object)
			}
		}

		// Limit second level entity count to avoid query explosion
		if len(secondLevelEntities) > 8 {
			secondLevelEntities = secondLevelEntities[:8]
		}

		// Execute batched queries for second level entities
		batchSize := 3 // Number of entities to include in a single query
		for i := 0; i < len(secondLevelEntities); i += batchSize {
			endIdx := i + batchSize
			if endIdx > len(secondLevelEntities) {
				endIdx = len(secondLevelEntities)
			}

			// Prepare batch of entities for this query
			batch := secondLevelEntities[i:endIdx]
			if len(batch) == 0 {
				continue
			}

			// Build WHERE condition for entity batch
			entityConditions := make([]string, 0, len(batch)*2)
			for _, e := range batch {
				safeEntity := strings.ReplaceAll(e, "'", "''")
				// Use fuzzy match at second level
				entityConditions = append(entityConditions,
					fmt.Sprintf("n1.name = '%s'", safeEntity),
					fmt.Sprintf("n2.name = '%s'", safeEntity))
			}
			entityWhere := strings.Join(entityConditions, " OR ")

			// Second level query with batch of entities
			secondLevelQuery := fmt.Sprintf(`
				cypher('%s', $$
					MATCH (n1:Entity)-[r]->(n2:Entity)
					WHERE (%s) %s
					RETURN n1.name as subject, type(r) as relation, n2.name as object, 
						   r.confidence as confidence, r.library_id as library_id, 
						   r.file_id as file_id, r.data_id as data_id,
						   2 as depth
					LIMIT %d
				$$) as (subject agtype, relation agtype, object agtype, confidence agtype, 
					   library_id agtype, file_id agtype, data_id agtype, depth agtype)
			`, g.schema, entityWhere, libraryCondition, limit/len(secondLevelEntities)*len(batch))

			// Execute second level batch query
			secondLevelResults, err := g.ExecuteCypher(secondLevelQuery)
			if err != nil {
				logs.Error("Second level graph query failed for batch: %s", err.Error())
				continue
			}

			// Add to total results
			allResults = append(allResults, secondLevelResults...)
		}
	}

	return allResults, nil
}
