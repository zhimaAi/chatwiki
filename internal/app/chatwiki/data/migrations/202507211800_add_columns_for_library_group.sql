-- +goose Up

ALTER TABLE "public"."chat_ai_library_group"
ADD COLUMN "sort" int4 NOT NULL DEFAULT 0,
ADD COLUMN "group_type" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_group"."group_type" IS '0:QA分段分组 1:知识库文件分组';

ALTER TABLE "public"."chat_ai_library_file"
ADD COLUMN "group_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."group_id" IS '分组id';