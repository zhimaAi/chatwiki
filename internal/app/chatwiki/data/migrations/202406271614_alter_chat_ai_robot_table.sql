-- +goose Up

ALTER TABLE "chat_ai_robot"
    ALTER COLUMN "chat_type" SET DEFAULT 1;