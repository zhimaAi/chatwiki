-- +goose Up

ALTER TABLE "public"."chat_ai_faq_files"
ALTER COLUMN "chunk_prompt" TYPE varchar(5000) COLLATE "pg_catalog"."default";
