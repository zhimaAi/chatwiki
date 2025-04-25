-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "similar_label" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_library_file" ADD COLUMN "similar_column" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_library_file"."similar_label" IS 'QA文档-相似问题标识符';
COMMENT ON COLUMN "chat_ai_library_file"."similar_column" IS 'QA文档-相似问题所在列';
