-- +goose Up

ALTER TABLE "chat_ai_message"
    ADD COLUMN recall_time int4 NOT NULL DEFAULT 0,
    ADD COLUMN request_time int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "chat_ai_message"."recall_time" IS '召回时间';
COMMENT ON COLUMN "chat_ai_message"."request_time" IS '请求时间';
