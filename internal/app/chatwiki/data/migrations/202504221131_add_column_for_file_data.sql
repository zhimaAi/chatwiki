-- +goose Up

ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "similar_questions" jsonb NOT NULL DEFAULT '[]';

COMMENT ON COLUMN "chat_ai_library_file_data"."similar_questions" IS '相似问题';