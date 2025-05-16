-- +goose Up

CREATE TABLE "public"."chat_ai_prompt_library_group"
(
    "id"            serial       NOT NULL primary key,
    "admin_user_id" int4         NOT NULL DEFAULT 0,
    "group_name"    varchar(100) NOT NULL DEFAULT '',
    "group_desc"    varchar(100) NOT NULL DEFAULT '',
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_prompt_library_group" ("admin_user_id");

COMMENT ON TABLE "public"."chat_ai_prompt_library_group" IS '提示词库-分组表';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."group_name" IS '分组名称';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."group_desc" IS '分组描述';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_group"."update_time" IS '更新时间';


CREATE TABLE "public"."chat_ai_prompt_library_items"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "title"         varchar(30)   NOT NULL DEFAULT '',
    "group_id"      int4          NOT NULL DEFAULT 0,
    "prompt_type"   int2          NOT NULL DEFAULT 0,
    "prompt"        varchar(2000) NOT NULL DEFAULT '',
    "prompt_struct" text          NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_prompt_library_items" ("admin_user_id", "group_id");
CREATE INDEX ON "public"."chat_ai_prompt_library_items" ("group_id");

COMMENT ON TABLE "public"."chat_ai_prompt_library_items" IS '提示词库-内容表';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."title" IS '标题';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."group_id" IS '分组ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."prompt_type" IS '提示词类型:0自定义提示词,1结构化提示词';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."prompt" IS '自定义提示词';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."prompt_struct" IS '结构化提示词配置json';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_prompt_library_items"."update_time" IS '更新时间';
