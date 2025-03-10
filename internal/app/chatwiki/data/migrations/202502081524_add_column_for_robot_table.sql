-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "application_type" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_robot"."application_type" IS '应用类型:0聊天机器人,1工作流';
