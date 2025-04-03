-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "semantic_chunk_use_model" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library_file" ADD COLUMN "semantic_chunk_model_config_id" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."semantic_chunk_use_model" IS '语义模式下分段选择的嵌入模型名';
COMMENT ON COLUMN "chat_ai_library_file"."semantic_chunk_model_config_id" IS '语义模式下分段选择的嵌入模型配置id';