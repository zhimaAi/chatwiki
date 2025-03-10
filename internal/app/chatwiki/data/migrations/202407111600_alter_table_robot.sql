-- +goose Up

ALTER TABLE "chat_ai_robot"
    ALTER COLUMN "max_token" SET DEFAULT 2000;
ALTER TABLE "chat_ai_robot"
    ALTER COLUMN "context_pair" SET DEFAULT 6;
ALTER TABLE "chat_ai_robot"
    ALTER COLUMN "top_k" SET DEFAULT 8;