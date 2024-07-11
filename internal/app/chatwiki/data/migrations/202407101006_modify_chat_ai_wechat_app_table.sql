-- +goose Up

ALTER TABLE chat_ai_robot
    ADD COLUMN fast_command_switch int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_robot"."fast_command_switch" IS '快捷指令1:开';
