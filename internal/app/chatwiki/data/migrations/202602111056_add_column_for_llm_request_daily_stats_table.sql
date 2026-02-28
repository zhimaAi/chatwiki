-- +goose Up
ALTER TABLE "public"."llm_request_daily_stats"
    ADD COLUMN "app_id" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."llm_request_daily_stats"."app_id" IS '应用ID，对应chat_ai_wechat_app';
