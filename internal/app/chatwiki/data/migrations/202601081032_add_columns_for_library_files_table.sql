-- +goose Up

ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "feishu_app_id" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "feishu_app_secret" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "feishu_document_id" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file"."feishu_app_id" IS '飞书app_id';
COMMENT ON COLUMN "public"."chat_ai_library_file"."feishu_app_secret" IS '飞书app_secret';
COMMENT ON COLUMN "public"."chat_ai_library_file"."feishu_document_id" IS '飞书文档id';


