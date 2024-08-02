-- +goose Up

ALTER TABLE "chat_ai_robot" ALTER COLUMN "answer_source_switch" SET DEFAULT true;
