-- +goose Up

CREATE TABLE "llm_token_app_daily_stats"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "robot_id"          int4 NOT NULL DEFAULT 0,
    "token_app_type"    varchar(32) NOT NULL DEFAULT '',
    "corp"              varchar(32) NOT NULL DEFAULT '',
    "model"             varchar(32) NOT NULL DEFAULT '',
    "type"              varchar(16) NOT NULL DEFAULT '',
    "prompt_token"      int8 NOT NULL DEFAULT 0,
    "completion_token"  int8 NOT NULL DEFAULT 0,
    "request_num"       int8 NOT NULL DEFAULT 0,
    "date"              date NOT NULL,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "llm_token_app_daily_stats"("admin_user_id", "token_app_type", "date" , "robot_id");

COMMENT ON TABLE "llm_token_app_daily_stats" IS 'token消耗按应用按日统计';

COMMENT ON COLUMN "llm_token_app_daily_stats"."id" IS '自增ID';
COMMENT ON COLUMN "llm_token_app_daily_stats"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "llm_token_app_daily_stats"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "llm_token_app_daily_stats"."token_app_type" IS '应用名：workflow工作流，chatwiki_robot机器人，other其他';
COMMENT ON COLUMN "llm_token_app_daily_stats"."corp" IS '模型的企业';
COMMENT ON COLUMN "llm_token_app_daily_stats"."model" IS '模型名';
COMMENT ON COLUMN "llm_token_app_daily_stats"."type" IS '类型。llm、embedding';
COMMENT ON COLUMN "llm_token_app_daily_stats"."prompt_token" IS '输入token量';
COMMENT ON COLUMN "llm_token_app_daily_stats"."completion_token" IS '输出token量';
COMMENT ON COLUMN "llm_token_app_daily_stats"."request_num" IS '请求次数';
COMMENT ON COLUMN "llm_token_app_daily_stats"."date" IS '日期';
COMMENT ON COLUMN "llm_token_app_daily_stats"."create_time" IS '创建时间';
COMMENT ON COLUMN "llm_token_app_daily_stats"."update_time" IS '更新时间';
