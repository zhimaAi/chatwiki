-- +goose Up

ALTER TABLE "chat_ai_robot"
    ALTER COLUMN "prompt" TYPE varchar(2000);