-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "ai_chunk_model_config_id" int4 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library_file" ADD COLUMN "ai_chunk_model" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library_file" ADD COLUMN "ai_chunk_prumpt" varchar(500) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library_file" ADD COLUMN "ai_chunk_task_id" varchar(500) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library_file" ADD COLUMN "ai_chunk_size" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."ai_chunk_model" IS 'AI分段模型';
COMMENT ON COLUMN "chat_ai_library_file"."ai_chunk_task_id" IS 'AI分段异步任务ID';

ALTER TABLE "chat_ai_library" ADD COLUMN "ai_chunk_model_config_id" int4 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library" ADD COLUMN "ai_chunk_model" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library" ADD COLUMN "ai_chunk_prumpt" varchar(500) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library" ADD COLUMN "ai_chunk_size" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library"."ai_chunk_model" IS 'AI分段模型';