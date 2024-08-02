-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "answer_source_switch" bool NOT NULL DEFAULT false;
COMMENT ON COLUMN "chat_ai_robot"."answer_source_switch" IS '显示引文开关';