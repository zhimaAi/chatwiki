-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "open_agent_write_file_tool" int2 NOT NULL DEFAULT 0,
    ADD COLUMN "open_agent_execute_tool"    int2 NOT NULL DEFAULT 0,
    ADD COLUMN "open_agent_edit_file_tool"  int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."open_agent_write_file_tool" IS '启用Agent写文件工具:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot"."open_agent_execute_tool" IS '启用Agent执行命令工具:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot"."open_agent_edit_file_tool" IS '启用Agent编辑文件工具:0关,1开';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN IF EXISTS "open_agent_write_file_tool",
    DROP COLUMN IF EXISTS "open_agent_execute_tool",
    DROP COLUMN IF EXISTS "open_agent_edit_file_tool";
