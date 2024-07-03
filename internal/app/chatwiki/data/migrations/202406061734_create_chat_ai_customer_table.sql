-- +goose Up

CREATE TABLE "chat_ai_customer"
(
    "id"                serial        NOT NULL primary key,
    "admin_user_id"     int4          NOT NULL DEFAULT 0,
    "openid"            varchar(100)  NOT NULL DEFAULT '',
    "nickname"          varchar(100)  NOT NULL DEFAULT '',
    "name"              varchar(100)  NOT NULL DEFAULT '',
    "avatar"            varchar(500)  NOT NULL DEFAULT '',
    "is_background"     int2          NOT NULL DEFAULT 0,
    "create_time"       int4          NOT NULL DEFAULT 0,
    "update_time"       int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_customer" ("admin_user_id");
CREATE UNIQUE INDEX ON "chat_ai_customer" ("openid", "admin_user_id");

COMMENT ON TABLE "chat_ai_customer" IS '文档问答机器人-客户表';

COMMENT ON COLUMN "chat_ai_customer"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_customer"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_customer"."openid" IS '客户ID';
COMMENT ON COLUMN "chat_ai_customer"."nickname" IS '客户昵称';
COMMENT ON COLUMN "chat_ai_customer"."name" IS '客户名称';
COMMENT ON COLUMN "chat_ai_customer"."avatar" IS '客户头像';
COMMENT ON COLUMN "chat_ai_customer".is_background IS '后台创建的';
COMMENT ON COLUMN "chat_ai_customer"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_customer"."update_time" IS '更新时间';
