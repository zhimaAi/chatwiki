-- +goose Up

CREATE TABLE "public"."chat_ai_unknown_issue_stats"
(
    "id"            serial NOT NULL primary key,
    "admin_user_id" int4   NOT NULL DEFAULT 0,
    "robot_id"      int4   NOT NULL DEFAULT 0,
    "stats_day"     int4   NOT NULL DEFAULT 0,
    "question"      text   NOT NULL DEFAULT '',
    "trigger_total" int4   NOT NULL DEFAULT 0,
    "create_time"   int4   NOT NULL DEFAULT 0,
    "update_time"   int4   NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "public"."chat_ai_unknown_issue_stats"
    ("admin_user_id", "robot_id", "stats_day", "question");
CREATE INDEX ON "public"."chat_ai_unknown_issue_stats" ("robot_id", "stats_day");

COMMENT ON TABLE "public"."chat_ai_unknown_issue_stats" IS '未知问题统计';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."stats_day" IS '日期(yyyymmdd)';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."question" IS '未知问题';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."update_time" IS '更新时间';