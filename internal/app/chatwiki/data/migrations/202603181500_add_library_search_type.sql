-- +goose Up

ALTER TABLE "public"."user_search_config"
ADD COLUMN "library_search_type" varchar(50) NOT NULL DEFAULT '';
COMMENT ON COLUMN "public"."user_search_config"."library_search_type" IS '知识库搜索类型 fullTextSearch全文检索 keywordSearch关键词匹配';

ALTER TABLE "public"."chat_ai_robot"
ADD COLUMN "library_search_type" varchar(50) NOT NULL DEFAULT '';
COMMENT ON COLUMN "public"."chat_ai_robot"."library_search_type" IS '知识库搜索类型 fullTextSearch全文检索 keywordSearch关键词匹配';

CREATE INDEX idx_content_bigm ON chat_ai_library_file_data_index
USING gin (content gin_bigm_ops);
