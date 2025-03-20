-- +goose Up

ALTER TABLE "chat_ai_library_file_doc"
    ADD COLUMN "draft_content"  text          NOT NULL DEFAULT '',
    ADD COLUMN "is_pub"      int2          NOT NULL DEFAULT 0,
    ADD COLUMN "delete_time"      int4          NOT NULL DEFAULT 0;
COMMENT ON COLUMN "chat_ai_library_file_doc"."draft_content" IS '文档-草稿';
COMMENT ON COLUMN "chat_ai_library_file_doc"."create_time" IS '1：已发布';
COMMENT ON COLUMN "chat_ai_library_file_doc"."delete_time" IS '删除时间';


