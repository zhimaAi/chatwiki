-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "think_switch"   int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_robot"."think_switch" IS '思考过程开关 1：开 0：关';