-- +goose Up

CREATE TABLE "chat_robot_ability"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "robot_id"          int4 NOT NULL DEFAULT 0,
    "ability_type"      varchar(100) NOT NULL DEFAULT '',
    "switch_status"     int2 NOT NULL DEFAULT 0,
    "fixed_menu"        int2 NOT NULL DEFAULT 0,
    "ai_reply_status"   int2 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "chat_robot_ability"("robot_id","ability_type");

COMMENT ON TABLE "chat_robot_ability" IS '聊天机器人能力表';

COMMENT ON COLUMN "chat_robot_ability"."id" IS '自增ID';
COMMENT ON COLUMN "chat_robot_ability"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_robot_ability"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "chat_robot_ability"."ability_type" IS '功能类型：keyword_reply关键词回复';
COMMENT ON COLUMN "chat_robot_ability"."switch_status" IS '开启状态：0关1开';
COMMENT ON COLUMN "chat_robot_ability"."fixed_menu" IS '固定菜单：0关1开';
COMMENT ON COLUMN "chat_robot_ability"."ai_reply_status" IS 'AI回复：0关1开';
COMMENT ON COLUMN "chat_robot_ability"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_robot_ability"."update_time" IS '更新时间';
