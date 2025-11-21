-- +goose Up

ALTER TABLE "public"."chat_ai_robot" 
  ALTER COLUMN "answer_source_switch" SET DEFAULT false;