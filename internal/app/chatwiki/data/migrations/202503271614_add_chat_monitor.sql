-- +goose Up

CREATE TABLE "public"."chat_ai_receiver"
(
    "id"                bigserial     NOT NULL primary key,
    "admin_user_id"     int4          NOT NULL DEFAULT 0,
    "robot_id"          int4          NOT NULL DEFAULT 0,
    "robot_key"         varchar(15)   NOT NULL DEFAULT '',
    "openid"            varchar(100)  NOT NULL DEFAULT '',
    "nickname"          varchar(100)  NOT NULL DEFAULT '',
    "name"              varchar(100)  NOT NULL DEFAULT '',
    "avatar"            varchar(500)  NOT NULL DEFAULT '',
    "session_id"        int4          NOT NULL DEFAULT 0,
    "last_chat_time"    int4          NOT NULL DEFAULT 0,
    "last_chat_message" varchar(1000) NOT NULL DEFAULT '',
    "app_type"          varchar(100)  NOT NULL DEFAULT '',
    "app_id"            varchar(100)  NOT NULL DEFAULT '',
    "come_from"         varchar(1000) NOT NULL DEFAULT '',
    "unread"            int4          NOT NULL DEFAULT 0,
    "create_time"       int4          NOT NULL DEFAULT 0,
    "update_time"       int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_receiver" ("admin_user_id");
CREATE INDEX ON "public"."chat_ai_receiver" ("robot_id");
CREATE INDEX ON "public"."chat_ai_receiver" ("openid");
CREATE UNIQUE INDEX ON "public"."chat_ai_receiver" ("session_id");
CREATE INDEX ON "public"."chat_ai_receiver" ("last_chat_time");

COMMENT ON TABLE "public"."chat_ai_receiver" IS '文档问答机器人-会话表';

COMMENT ON COLUMN "public"."chat_ai_receiver"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."robot_key" IS '对外的key';
COMMENT ON COLUMN "public"."chat_ai_receiver"."openid" IS '客户ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."nickname" IS '客户昵称';
COMMENT ON COLUMN "public"."chat_ai_receiver"."name" IS '客户名称';
COMMENT ON COLUMN "public"."chat_ai_receiver"."avatar" IS '客户头像';
COMMENT ON COLUMN "public"."chat_ai_receiver"."session_id" IS '会话ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."last_chat_time" IS '最近聊天时间';
COMMENT ON COLUMN "public"."chat_ai_receiver"."last_chat_message" IS '最近聊天内容';
COMMENT ON COLUMN "public"."chat_ai_receiver"."app_type" IS '应用类型';
COMMENT ON COLUMN "public"."chat_ai_receiver"."app_id" IS '应用ID';
COMMENT ON COLUMN "public"."chat_ai_receiver"."come_from" IS '来自';
COMMENT ON COLUMN "public"."chat_ai_receiver"."unread" IS '未读消息数';
COMMENT ON COLUMN "public"."chat_ai_receiver"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_receiver"."update_time" IS '更新时间';

ALTER TABLE "public"."chat_ai_message"
    ADD COLUMN "nickname" varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "name"     varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "avatar"   varchar(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_message"."nickname" IS '发送人昵称';
COMMENT ON COLUMN "public"."chat_ai_message"."name" IS '发送人名称';
COMMENT ON COLUMN "public"."chat_ai_message"."avatar" IS '发送人头像';