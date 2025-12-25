-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "question_multiple_switch" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."question_multiple_switch" IS '多模态输入开关:0关,1开';
