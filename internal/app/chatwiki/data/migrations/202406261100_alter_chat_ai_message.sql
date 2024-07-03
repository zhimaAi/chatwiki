-- +goose Up

ALTER TABLE "chat_ai_message" ALTER COLUMN "content" TYPE text;
