-- +goose Up

ALTER TABLE "public"."chat_ai_export_task"
ADD COLUMN "library_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_export_task"."library_id" IS '知识库id';