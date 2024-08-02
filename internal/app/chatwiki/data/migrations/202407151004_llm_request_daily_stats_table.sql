-- +goose Up

CREATE TABLE "llm_request_daily_stats"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "robot_id"          int4 NOT NULL DEFAULT 0,
    "date"              date NOT NULL,
    "app_type"          varchar(100) NOT NULL DEFAULT '',
    "type"              int2 NOT NULL DEFAULT 1,
    "amount"            int4 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "llm_request_daily_stats"("admin_user_id", "robot_id", "date", "app_type", "type");

COMMENT ON TABLE "llm_request_daily_stats" IS '每日统计分析';

COMMENT ON COLUMN "llm_request_daily_stats"."id" IS '自增ID';
COMMENT ON COLUMN "llm_request_daily_stats"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "llm_request_daily_stats"."robot_id" IS '机器人id';
COMMENT ON COLUMN "llm_request_daily_stats"."date" IS '日期';
COMMENT ON COLUMN "llm_request_daily_stats"."app_type" IS '应用类别';
COMMENT ON COLUMN "llm_request_daily_stats"."type" IS '类别 1日活用户数 2日新增用户数 3总消息数 4token消耗数';
COMMENT ON COLUMN "llm_request_daily_stats"."amount" IS '数量';
COMMENT ON COLUMN "llm_request_daily_stats"."create_time" IS '创建时间';
COMMENT ON COLUMN "llm_request_daily_stats"."update_time" IS '更新时间';
