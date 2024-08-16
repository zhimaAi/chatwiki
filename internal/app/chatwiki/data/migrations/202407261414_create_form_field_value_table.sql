-- +goose Up

CREATE TABLE "form_field_value"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "form_entry_id"     int4 NOT NULL DEFAULT 0,
    "form_field_id"     int4 NOT NULL DEFAULT 0,
    "type"              varchar(32) NOT NULL DEFAULT '',
    "string_content"    varchar(512) NOT NULL DEFAULT '',
    "integer_content"   int4 NOT NULL DEFAULT 0,
    "number_content"     float4 NOT NULL DEFAULT 0,
    "boolean_content"   boolean NOT NULL DEFAULT false,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "form_field_value" ("admin_user_id", "form_entry_id", "form_field_id", "string_content", "integer_content", "number_content", "boolean_content");

COMMENT ON TABLE "form_field_value" IS '字段内容';

COMMENT ON COLUMN "form_field_value"."id" IS 'ID';
COMMENT ON COLUMN "form_field_value"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "form_field_value"."form_entry_id" IS '表单条目的id';
COMMENT ON COLUMN "form_field_value"."form_field_id" IS '关联form_field表的id';
COMMENT ON COLUMN "form_field_value"."type" IS '数据类型 string、integer、float、boolean 等';
COMMENT ON COLUMN "form_field_value"."string_content" IS '字符串内容';
COMMENT ON COLUMN "form_field_value"."integer_content" IS '整数内容';
COMMENT ON COLUMN "form_field_value"."number_content" IS '浮点数内容';
COMMENT ON COLUMN "form_field_value"."boolean_content" IS '布尔内容';
COMMENT ON COLUMN "form_field_value"."create_time" IS '创建时间';
COMMENT ON COLUMN "form_field_value"."update_time" IS '更新时间';
