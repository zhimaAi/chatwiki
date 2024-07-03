-- +goose Up

CREATE TABLE "chat_ai_dialogue"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "robot_id"      int4          NOT NULL DEFAULT 0,
    "openid"        varchar(100)  NOT NULL DEFAULT '',
    "is_background" int2          NOT NULL DEFAULT 0,
    "subject"       varchar(1000) NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_dialogue" ("admin_user_id");
CREATE INDEX ON "chat_ai_dialogue" ("robot_id", "is_background");
CREATE INDEX ON "chat_ai_dialogue" ("openid");

COMMENT ON TABLE "chat_ai_dialogue" IS '文档问答机器人-对话表';

COMMENT ON COLUMN "chat_ai_dialogue"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_dialogue"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_dialogue"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_dialogue"."openid" IS '客户ID';
COMMENT ON COLUMN "chat_ai_dialogue".is_background IS '后台创建的';
COMMENT ON COLUMN "chat_ai_dialogue"."subject" IS '对话主题';
COMMENT ON COLUMN "chat_ai_dialogue"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_dialogue"."update_time" IS '更新时间';
