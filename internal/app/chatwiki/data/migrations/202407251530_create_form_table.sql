-- +goose Up

CREATE TABLE "form"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "name"              varchar(100) NOT NULL DEFAULT '',
    "description"       varchar(1000) NOT NULL DEFAULT '',
    "delete_time"       int4 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);
CREATE INDEX ON "form" ("admin_user_id", "name");

COMMENT ON TABLE "form" IS '表单';

COMMENT ON COLUMN "form"."id" IS 'ID';
COMMENT ON COLUMN "form"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "form"."name" IS '表单名称';
COMMENT ON COLUMN "form"."description" IS '表单描述';
COMMENT ON COLUMN "form"."delete_time" IS '删除时间';
COMMENT ON COLUMN "form"."create_time" IS '创建时间';
COMMENT ON COLUMN "form"."update_time" IS '更新时间';
