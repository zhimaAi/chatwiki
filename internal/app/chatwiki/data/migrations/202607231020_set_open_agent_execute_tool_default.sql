-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ALTER COLUMN "open_agent_execute_tool" SET DEFAULT 1;

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    ALTER COLUMN "open_agent_execute_tool" SET DEFAULT 0;
