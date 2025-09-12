-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ALTER COLUMN "prompt" TYPE varchar(10000);

ALTER TABLE "public"."chat_ai_prompt_library_items"
    ALTER COLUMN "prompt" TYPE varchar(10000);