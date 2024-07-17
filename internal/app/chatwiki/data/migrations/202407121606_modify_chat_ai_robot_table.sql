-- +goose Up

ALTER TABLE chat_ai_robot
    ADD COLUMN yunpc_fast_command_switch int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_robot"."yunpc_fast_command_switch" IS '嵌入网站-快捷指令1:开';

ALTER TABLE fast_command
    ADD COLUMN app_id int4 NOT NULL DEFAULT -1;

COMMENT ON COLUMN "fast_command"."app_id" IS 'app_id,-1:yun_h5,-2:yun_pc';

CREATE INDEX idx_robot_app_id ON public.fast_command (robot_id,app_id);
