-- +goose Up

ALTER TABLE "public"."chat_ai_library"
ADD COLUMN "qa_index_type" int4 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_library"."qa_index_type" IS 'QA文档-索引方式,1问题答案一起生成索引,2仅对问题生成索引';