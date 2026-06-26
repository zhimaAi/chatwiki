-- +goose Up

ALTER TABLE "public"."work_flow_logs"
    ADD COLUMN "status" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."work_flow_logs"."status" IS '状态：0=运行中，1=已完成，2=已停止';

CREATE INDEX IF NOT EXISTS "idx_work_flow_logs_admin_status"
    ON "public"."work_flow_logs" ("admin_user_id", "status", "create_time" DESC);
