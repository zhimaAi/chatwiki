-- +goose Up

ALTER TABLE "chat_ai_library_file"
    ADD COLUMN "doc_type" int2          NOT NULL DEFAULT 1,
    ADD COLUMN "doc_url"  varchar(2000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_library_file"."doc_type" IS '文档类型:1本地文档,2在线文档,3自定义文档';
COMMENT ON COLUMN "chat_ai_library_file"."doc_url" IS '在线文档的url链接';