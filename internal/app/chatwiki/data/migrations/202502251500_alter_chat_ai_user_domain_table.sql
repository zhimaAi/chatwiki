-- +goose Up

ALTER TABLE "chat_ai_user_domain"
    ADD COLUMN "is_upload"  int2 NOT NULL DEFAULT 0;


