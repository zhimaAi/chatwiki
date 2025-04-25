-- +goose Up

CREATE TABLE "category"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "name"          varchar(8)    NOT NULL DEFAULT '',
    "type"          varchar(1)    NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);
CREATE INDEX ON "category" ("admin_user_id");

COMMENT ON TABLE "category" IS '分类标记';

COMMENT ON COLUMN "category"."id" IS 'ID';
COMMENT ON COLUMN "category"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "category"."name" IS '名称';
COMMENT ON COLUMN "category"."type" IS '类别 a、b、c、d、e等';
COMMENT ON COLUMN "category"."create_time" IS '创建时间';
COMMENT ON COLUMN "category"."update_time" IS '更新时间';