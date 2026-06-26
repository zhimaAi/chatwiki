-- +goose Up

ALTER TABLE "public"."work_flow_storage_cache"
    ADD COLUMN "log_id" int8 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."work_flow_storage_cache"."log_id" IS '关联 work_flow_logs.id，维系一次调用的稳定身份';

CREATE INDEX IF NOT EXISTS "idx_work_flow_storage_cache_log_id"
    ON "public"."work_flow_storage_cache" ("log_id");
