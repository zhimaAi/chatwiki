-- +goose Up

CREATE TABLE message_feedback
(
    "id"                    serial          NOT NULL primary key,
    "admin_user_id"         int4            NOT NULL DEFAULT 0,
    "robot_id"              int4            NOT NULL DEFAULT 0,
    "customer_message_id"   int4            NOT NULL DEFAULT 0,
    "ai_message_id"         int4            NOT NULL DEFAULT 0,
    "type"                  int2            NOT NULL DEFAULT 1,
    "content"               varchar(1000)   NOT NULL DEFAULT '',
    "robot"                 json            NOT NULL DEFAULT '{}',
    "create_time"           int4            NOT NULL DEFAULT 0,
    "update_time"           int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "message_feedback" ("admin_user_id", "robot_id", "ai_message_id", "type", "create_time");

COMMENT ON TABLE "message_feedback" IS '聊天消息反馈记录';

COMMENT ON COLUMN "message_feedback"."id" IS 'ID';
COMMENT ON COLUMN "message_feedback"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "message_feedback"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "message_feedback"."customer_message_id" IS '客户消息ID';
COMMENT ON COLUMN "message_feedback"."ai_message_id" IS '机器人消息ID';
COMMENT ON COLUMN "message_feedback"."type" IS '反馈类别 1点赞 2点踩';
COMMENT ON COLUMN "message_feedback"."content" IS '反馈内容';
COMMENT ON COLUMN "message_feedback"."robot" IS '当前机器人配置信息';
COMMENT ON COLUMN "message_feedback"."create_time" IS '创建时间';
COMMENT ON COLUMN "message_feedback"."update_time" IS '更新时间';
