-- +goose Up

SET search_path = ag_catalog, "$user", public;

LOAD 'age';

select * from cypher('graphrag', $$
    MERGE (s:Entity {name: 'chatwiki', library_id: 0, file_id: 0, data_id: 0})
    MERGE (o:Entity {name: '属于', library_id: 0, file_id: 0, data_id: 0})
    CREATE (s)-[r:RELATES_TO {confidence: '芝麻小事', library_id: 0, file_id: 0, data_id: 0}]->(o)
    RETURN r
$$) as (r agtype);

CREATE INDEX graphrag_entity_name_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"name"'::agtype]));

CREATE INDEX graphrag_entity_file_id_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"file_id"'::agtype]));

CREATE INDEX graphrag_entity_library_id_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"library_id"'::agtype]));