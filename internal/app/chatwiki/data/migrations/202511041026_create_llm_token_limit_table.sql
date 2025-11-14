-- +goose Up

CREATE TABLE "llm_token_app_limit"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "robot_id"          int4 NOT NULL DEFAULT 0,
    "token_app_type"    varchar(32) NOT NULL DEFAULT '',
    "max_token"         int8 NOT NULL DEFAULT 0,
    "use_token"         int8 NOT NULL DEFAULT 0,
    "description"       varchar(500) NOT NULL DEFAULT '',
    "switch_status"     int2 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX on "llm_token_app_limit"("admin_user_id", "token_app_type", "robot_id");

COMMENT ON TABLE "llm_token_app_limit" IS 'token每日限额配置';

COMMENT ON COLUMN "llm_token_app_limit"."id" IS '自增ID';
COMMENT ON COLUMN "llm_token_app_limit"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "llm_token_app_limit"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "llm_token_app_limit"."token_app_type" IS '应用名：workflow工作流，chatwiki_robot机器人，other其他';
COMMENT ON COLUMN "llm_token_app_limit"."max_token" IS '限额token量';
COMMENT ON COLUMN "llm_token_app_limit"."use_token" IS '已消耗token量';
COMMENT ON COLUMN "llm_token_app_limit"."description" IS '备注';
COMMENT ON COLUMN "llm_token_app_limit"."switch_status" IS '开关状态，1开启，0关闭（默认）';
COMMENT ON COLUMN "llm_token_app_limit"."create_time" IS '创建时间';
COMMENT ON COLUMN "llm_token_app_limit"."update_time" IS '更新时间';
