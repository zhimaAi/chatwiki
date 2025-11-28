-- +goose Up

CREATE TABLE "public"."plugin_config" (
    "id" serial NOT NULL primary key,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(255) NOT NULL DEFAULT 0,
    "type" varchar(255) NOT NULL DEFAULT 'notice',
    "has_loaded" bool NOT NULL DEFAULT false,
    "data" jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX ON "public"."plugin_config" ("admin_user_id");

COMMENT ON TABLE "public"."plugin_config" IS '插件配置表';
COMMENT ON COLUMN "public"."plugin_config"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."plugin_config"."name" IS '插件名称';
COMMENT ON COLUMN "public"."plugin_config"."type" IS '类别';
COMMENT ON COLUMN "public"."plugin_config"."has_loaded" IS '是否已加载';
COMMENT ON COLUMN "public"."plugin_config"."data" IS '插件数据';
