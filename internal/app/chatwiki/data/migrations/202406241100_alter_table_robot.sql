-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "show_type" int2          NOT NULL DEFAULT 1;
COMMENT ON COLUMN "chat_ai_robot"."show_type" IS '展示类型:1:markdown';
