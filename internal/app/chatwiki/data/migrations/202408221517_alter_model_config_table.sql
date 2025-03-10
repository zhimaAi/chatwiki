-- +goose Up

ALTER TABLE "chat_ai_model_config" ALTER COLUMN "api_key" TYPE varchar(1024);