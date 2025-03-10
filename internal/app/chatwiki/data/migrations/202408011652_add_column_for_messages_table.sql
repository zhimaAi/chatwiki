-- +goose Up

ALTER TABLE "chat_ai_message" ADD COLUMN "is_valid_function_call" boolean NOT NULL DEFAULT false;

COMMENT ON COLUMN "chat_ai_message"."is_valid_function_call" IS '是否为有效的函数调用';
