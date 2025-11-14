-- +goose Up

ALTER TABLE "public"."work_flow_logs"
    ADD COLUMN "version_id" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "node_logs" text NOT NULL DEFAULT '[]',
    ADD COLUMN "question" varchar(1000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."work_flow_logs"."version_id" IS '版本id';
COMMENT ON COLUMN "public"."work_flow_logs"."node_logs" IS '节点执行日志';
COMMENT ON COLUMN "public"."work_flow_logs"."question" IS '问题';

CREATE INDEX ON "work_flow_logs" ("admin_user_id" , "create_time");