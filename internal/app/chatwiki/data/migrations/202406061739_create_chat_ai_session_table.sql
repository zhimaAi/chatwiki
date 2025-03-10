-- +goose Up

CREATE TABLE "chat_ai_session"
(
    "id"                serial        NOT NULL primary key,
    "admin_user_id"     int4          NOT NULL DEFAULT 0,
    "dialogue_id"       int4          NOT NULL DEFAULT 0,
    "robot_id"          int4          NOT NULL DEFAULT 0,
    "openid"            varchar(100)  NOT NULL DEFAULT '',
    "last_chat_time"    int4          NOT NULL DEFAULT 0,
    "last_chat_message" varchar(1000) NOT NULL DEFAULT '',
    "create_time"       int4          NOT NULL DEFAULT 0,
    "update_time"       int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_session" ("admin_user_id");
CREATE INDEX ON "chat_ai_session" ("dialogue_id");
CREATE INDEX ON "chat_ai_session" ("robot_id");
CREATE INDEX ON "chat_ai_session" ("openid");

COMMENT ON TABLE "chat_ai_session" IS '文档问答机器人-会话表';

COMMENT ON COLUMN "chat_ai_session"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_session"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_session"."dialogue_id" IS '对话ID';
COMMENT ON COLUMN "chat_ai_session"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_session"."openid" IS '客户ID';
COMMENT ON COLUMN "chat_ai_session"."last_chat_time" IS '最近聊天时间';
COMMENT ON COLUMN "chat_ai_session"."last_chat_message" IS '最近聊天内容';
COMMENT ON COLUMN "chat_ai_session"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_session"."update_time" IS '更新时间';
