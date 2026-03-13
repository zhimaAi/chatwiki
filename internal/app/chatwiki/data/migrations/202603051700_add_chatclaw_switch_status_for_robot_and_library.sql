-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
ADD COLUMN "chat_claw_switch_status" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_robot"."chat_claw_switch_status" IS 'ChatClaw展示开关:0关,1开';

ALTER TABLE "public"."chat_ai_library"
ADD COLUMN "chat_claw_switch_status" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_library"."chat_claw_switch_status" IS 'ChatClaw展示开关:0关,1开';
