-- +goose Up

CREATE TABLE "company"
(
    "id" serial NOT NULL primary key,
    "name" varchar(100) NOT NULL DEFAULT '' ,
    "avatar" varchar(255) NOT NULL DEFAULT '',
    "parent_id" int4 NOT NULL DEFAULT 0,
    "is_deleted" int2 NOT NULL DEFAULT 0,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

COMMENT ON COLUMN "company"."id" IS 'ID';
COMMENT ON COLUMN "company"."name" IS '名称';
COMMENT ON COLUMN "company"."avatar" IS 'con';
COMMENT ON COLUMN "company"."is_deleted" IS '是否删除 1:删除';
COMMENT ON COLUMN "company"."create_time" IS '创建时间';
COMMENT ON COLUMN "company"."update_time" IS '更新时间';
COMMENT ON COLUMN "company"."parent_id" IS '所属人id';
