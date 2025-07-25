-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "enable_thinking" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."enable_thinking" IS '深度思考开关:0关,1开';
