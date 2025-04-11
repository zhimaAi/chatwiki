-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "feedback_switch"   int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_robot"."feedback_switch" IS '反馈开关 1：开 0：关';