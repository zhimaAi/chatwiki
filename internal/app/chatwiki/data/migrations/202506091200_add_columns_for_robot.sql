-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
ADD COLUMN "default_library_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_robot"."default_library_id" IS '默认关联知识库';