-- +goose Up
ALTER TABLE "public"."work_flow_storage_cache"
    ALTER COLUMN "dialog_id" TYPE bigint,
    ALTER COLUMN "session_id" TYPE bigint;
