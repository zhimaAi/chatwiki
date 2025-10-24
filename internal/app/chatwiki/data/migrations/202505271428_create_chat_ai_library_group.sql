-- +goose Up

CREATE TABLE "public"."chat_ai_library_group"
(
    "id"            serial      NOT NULL primary key,
    "admin_user_id" int4        NOT NULL DEFAULT 0,
    "library_id"    int4        NOT NULL DEFAULT 0,
    "group_name"    varchar(50) NOT NULL DEFAULT '',
    "create_time"   int4        NOT NULL DEFAULT 0,
    "update_time"   int4        NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_library_group" ("admin_user_id");
CREATE INDEX ON "public"."chat_ai_library_group" ("library_id");

COMMENT ON TABLE "public"."chat_ai_library_group" IS '文档问答机器人-知识库分组';
COMMENT ON COLUMN "public"."chat_ai_library_group"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_library_group"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_library_group"."library_id" IS '知识库ID';
COMMENT ON COLUMN "public"."chat_ai_library_group"."group_name" IS '分组名称';
COMMENT ON COLUMN "public"."chat_ai_library_group"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_library_group"."update_time" IS '更新时间';


ALTER TABLE "public"."chat_ai_library_file_data"
    ADD COLUMN "group_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_library_file_data"."group_id" IS '分组ID';