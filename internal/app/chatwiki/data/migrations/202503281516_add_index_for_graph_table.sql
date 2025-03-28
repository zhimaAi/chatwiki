-- +goose Up

CREATE INDEX graphrag_entity_name_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"name"'::agtype]));

CREATE INDEX graphrag_entity_file_id_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"file_id"'::agtype]));

CREATE INDEX graphrag_entity_library_id_idx ON "graphrag"."Entity"
    USING btree (agtype_access_operator(VARIADIC ARRAY[properties, '"library_id"'::agtype]));