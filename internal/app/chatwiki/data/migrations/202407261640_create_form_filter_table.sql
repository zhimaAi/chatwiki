-- +goose Up

CREATE TABLE "form_filter"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "form_id"           int4 NOT NULL DEFAULT 0,
    "name"              varchar(100) NOT NULL DEFAULT '',
    "type"              int2 NOT NULL DEFAULT 1,
    "enabled"           bool NOT NULL DEFAULT true,
    "sort"             int4 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "form_filter" ("form_id", "sort", "enabled");

COMMENT ON TABLE "form_filter" IS '表单过滤器';

COMMENT ON COLUMN "form_filter"."id" IS 'ID';
COMMENT ON COLUMN "form_filter"."form_id" IS '表单ID';
COMMENT ON COLUMN "form_filter"."name" IS '名称';
COMMENT ON COLUMN "form_filter"."type" IS '条件之间的关系类型 1与 2或';
COMMENT ON COLUMN "form_filter"."enabled" IS '是否启用';
COMMENT ON COLUMN "form_filter"."sort" IS '排序 数值越大越靠前';
COMMENT ON COLUMN "form_filter"."create_time" IS '创建时间';
COMMENT ON COLUMN "form_filter"."update_time" IS '更新时间';
