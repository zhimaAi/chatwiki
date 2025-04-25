-- +goose Up

ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "category_id" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file_data"."category_id" IS '所属分类id';