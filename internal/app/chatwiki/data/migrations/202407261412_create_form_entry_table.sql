-- +goose Up

CREATE TABLE "form_entry"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "form_id"           int4 NOT NULL DEFAULT 0,
    "delete_time"       int4 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "form_entry" ("admin_user_id", "form_id", "delete_time");

COMMENT ON TABLE "form_entry" IS '数据条目';

COMMENT ON COLUMN "form_entry"."id" IS 'ID';
COMMENT ON COLUMN "form_entry"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "form_entry"."form_id" IS '关联的表单ID';
COMMENT ON COLUMN "form_entry"."delete_time" IS '删除时间';
COMMENT ON COLUMN "form_entry"."create_time" IS '创建时间';
COMMENT ON COLUMN "form_entry"."update_time" IS '更新时间';
