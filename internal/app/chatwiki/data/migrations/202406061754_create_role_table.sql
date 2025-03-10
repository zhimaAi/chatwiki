-- +goose Up

CREATE TABLE "role"
(
    "id"            serial          NOT NULL primary key,
    "name"          varchar(100)    NOT NULL DEFAULT '',
    "mark"          varchar(500)    NOT NULL DEFAULT 0,
    "is_deleted"    int2            NOT NULL DEFAULT 0,
    "create_time"   int4            NOT NULL DEFAULT 0,
    "update_time"   int4            NOT NULL DEFAULT 0,
    "operate_id"    int4            NOT NULL DEFAULT 0,
    "role_type"     int4            NOT NULL DEFAULT 0,
    "operate_name"  varchar(128)    NOT NULL DEFAULT '',
    "create_name"   varchar(128)    NOT NULL DEFAULT '',
    "parent_id"     int4            DEFAULT 0 NOT NULL
);

COMMENT ON TABLE role IS '文档问答机器人-角色表';

COMMENT ON COLUMN "role"."id" IS 'ID';
COMMENT ON COLUMN "role"."name" IS '角色名';
COMMENT ON COLUMN "role"."mark" IS '备注';
COMMENT ON COLUMN "role"."is_deleted" IS '是否删除 1:删除';
COMMENT ON COLUMN "role"."create_time" IS '创建时间';
COMMENT ON COLUMN "role"."update_time" IS '更新时间';
COMMENT ON COLUMN "role"."operate_id" IS '操作人id';
COMMENT ON COLUMN "role"."operate_name" IS '操作人';
COMMENT ON COLUMN "role"."create_name" IS '创建人';
COMMENT ON COLUMN "role".parent_id IS '所属人';
COMMENT ON COLUMN "role"."role_type" IS '角色类型1：所有者,2:管理员 3：成员';
