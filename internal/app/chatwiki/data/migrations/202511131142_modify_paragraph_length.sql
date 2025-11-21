-- +goose Up

ALTER TABLE "public"."chat_ai_answer_source"
    ALTER COLUMN content TYPE VARCHAR(10000),
    ALTER COLUMN question TYPE VARCHAR(10000),
    ALTER COLUMN answer TYPE VARCHAR(10000);
