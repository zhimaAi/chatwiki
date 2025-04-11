// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"fmt"
	"strings"
	"time"

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
	if len(args) > 0 {
		query = fmt.Sprintf("select * from "+query, args...)
	} else {
		query = "select * from " + query
	}
	queries := []string{
		`SELECT * FROM ag_catalog.get_cypher_keywords() limit 0`,
		`SET search_path = ag_catalog, "$user", public`,
		query,
	}
	return g.ExecuteLastQuery(queries)
}

// ExecuteLastQuery 执行最后一条查询
func (g *GraphDB) ExecuteLastQuery(queries []string) ([]msql.Params, error) {
	tx, err := msql.Begin(define.Postgres)
	if err != nil {
		return nil, err
	}
	result := make([]msql.Params, 0)
	for _, query := range queries {
		result, err = msql.RawValues(define.Postgres, query, tx)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, err
		}
	}

	return result, err
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
       		MATCH (n {library_id: %d})
       		OPTIONAL MATCH (n)-[r]-()
       		DELETE r, n
			RETURN count(distinct n) as nodes, count(r) as relations
   		$$) as (nodes agtype, relations agtype)
	`, g.schema, libraryId)

	_, err := g.ExecuteCypher(query)
	return err
}

// DeleteByFile 删除文件相关的所有图数据
func (g *GraphDB) DeleteByFile(fileId int) error {
	query := fmt.Sprintf(`
		cypher('%s', $$
       		MATCH (n {file_id: %d})
       		OPTIONAL MATCH (n)-[r]-()
       		DELETE r, n
       		RETURN count(distinct n) as nodes, count(r) as relations
   		$$) as (nodes agtype, relations agtype)
	`, g.schema, fileId)

	_, err := g.ExecuteCypher(query)
	return err
}

// DeleteByData 删除数据相关的所有图数据
func (g *GraphDB) DeleteByData(dataId int) error {
	query := fmt.Sprintf(`
		cypher('%s', $$
       		MATCH (n {data_id: %d})
       		OPTIONAL MATCH (n)-[r]-()
       		DELETE r, n
       		RETURN count(distinct n) as nodes, count(r) as relations
   		$$) as (nodes agtype, relations agtype)
	`, g.schema, dataId)

	_, err := g.ExecuteCypher(query)
	return err
}

// FindRelatedEntities Find related entities in the graph database, up to 3 levels deep
func (g *GraphDB) FindRelatedEntities(entity string, libraryIds []string, limit int, maxDepth int) ([]msql.Params, error) {
	// Process entity name, remove special characters
	entity = strings.ReplaceAll(entity, "'", "''")

	// Store all results
	allResults := make([]msql.Params, 0)

	// 构建IN条件
	inCondition := ""
	for i, id := range libraryIds {
		if i > 0 {
			inCondition += ","
		}
		inCondition += id
	}

	// 带库ID过滤条件的查询，将过滤移到WITH子句中
	query := fmt.Sprintf(`
			cypher('%s', $$
				// 查找具有指定深度的路径
				MATCH path = (start:Entity)-[*1..%d]-(connected:Entity)
				WHERE start.name = '%s'
				
				// 解开路径中的关系并直接过滤
				WITH path, length(path) AS depth
				UNWIND relationships(path) AS rel
				WITH rel, startNode(rel) AS start_node, endNode(rel) AS end_node, depth
				WHERE rel.library_id IN [%s]
				
				// 返回结果
				RETURN start_node.name AS subject, 
					   type(rel) AS relation, 
					   end_node.name AS object,
					   rel.confidence AS confidence, 
					   rel.library_id AS library_id,
					   rel.file_id AS file_id, 
					   rel.data_id AS data_id,
					   depth AS depth
				LIMIT %d
			$$) as (subject agtype, relation agtype, object agtype, confidence agtype, 
				   library_id agtype, file_id agtype, data_id agtype, depth agtype)
		`, g.schema, maxDepth, entity, inCondition, limit)

	results, err := g.ExecuteCypher(query)
	if err != nil {
		logs.Error("Graph query failed: %s", err.Error())
		return allResults, err
	}
	return results, nil
}
