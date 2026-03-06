-- +goose Up

ALTER TABLE "chat_ai_message" ADD COLUMN "debug_log" text NOT NULL default '';

COMMENT ON COLUMN "chat_ai_message"."debug_log" IS '运行日志';
