-- +goose Up

CREATE TABLE "form_field"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "form_id"           int4 NOT NULL DEFAULT 0,
    "name"              varchar(100) NOT NULL DEFAULT '',
    "description"       varchar(1000) NOT NULL DEFAULT '',
    "type"              varchar(100) NOT NULL DEFAULT '',
    "required"          bool NOT NULL DEFAULT true,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "form_field" ("admin_user_id", "form_id", "name");

COMMENT ON TABLE "form_field" IS '表单字段';

COMMENT ON COLUMN "form_field"."id" IS 'ID';
COMMENT ON COLUMN "form_field"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "form_field"."form_id" IS '关联的表单ID';
COMMENT ON COLUMN "form_field"."name" IS '字段名';
COMMENT ON COLUMN "form_field"."description" IS '字段描述';
COMMENT ON COLUMN "form_field"."type" IS '字段类型';
COMMENT ON COLUMN "form_field"."required" IS '是否必填';
COMMENT ON COLUMN "form_field"."create_time" IS '创建时间';
COMMENT ON COLUMN "form_field"."update_time" IS '更新时间';
