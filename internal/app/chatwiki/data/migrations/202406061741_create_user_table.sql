-- +goose Up

CREATE TABLE "user"
(
    "id"            serial       NOT NULL primary key,
    "user_name"     varchar(128) NOT NULL DEFAULT '',
    "password"      varchar(128) NOT NULL DEFAULT '',
    "salt"          varchar(100) NOT NULL DEFAULT '',
    "user_type"     int4         NOT NULL DEFAULT 0,
    "login_ip"      varchar(128) DEFAULT '',
    "nick_name"     varchar(128) NOT NULL DEFAULT '',
    "avatar"        varchar(255) NOT NULL DEFAULT '',
    "is_deleted"    int2         NOT NULL DEFAULT 0,
    "operate_name"  varchar(128) NOT NULL DEFAULT '',
    "operate_id"    int4         NOT NULL DEFAULT 0,
    "user_roles"    varchar(64)  NOT NULL DEFAULT '0',
    "parent_id"     int4         DEFAULT 0 NOT NULL,
    "login_time"    int4         DEFAULT 0 NOT NULL,
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "user" ("user_name");

COMMENT ON TABLE "user" IS '文档问答机器人-用户表';

COMMENT ON COLUMN "user"."id" IS 'ID';
COMMENT ON COLUMN "user"."user_name" IS '用户名';
COMMENT ON COLUMN "user"."password" IS '密码';
COMMENT ON COLUMN "user"."login_ip" IS '登录IP';
COMMENT ON COLUMN "user"."nick_name" IS '昵称';
COMMENT ON COLUMN "user"."avatar" IS '头像';
COMMENT ON COLUMN "user"."is_deleted" IS '是否删除 1:删除';
COMMENT ON COLUMN "user"."operate_name" IS '操作人';
COMMENT ON COLUMN "user"."operate_id" IS '操作人id';
COMMENT ON COLUMN "user"."user_roles" IS '用户角色组';
COMMENT ON COLUMN "user".parent_id IS '所属人id';
COMMENT ON COLUMN "user".login_time IS '登录时间';
COMMENT ON COLUMN "user"."create_time" IS '创建时间';
COMMENT ON COLUMN "user"."update_time" IS '更新时间';
COMMENT ON COLUMN "user"."salt" IS '密码加盐';
COMMENT ON COLUMN "user"."user_type" IS '用户类型:1管理员';
