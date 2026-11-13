-- +goose Up

ALTER TABLE "public"."chat_ai_prompt_optimize_record"
    ADD COLUMN "reply" text NOT NULL DEFAULT '',
    ADD COLUMN "reasoning_content" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."reply" IS '提示词助手的回复，包含改动说明或普通问答正文';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."reasoning_content" IS '模型返回的可展示思考过程，不参与后续模型上下文';

-- +goose Down

ALTER TABLE "public"."chat_ai_prompt_optimize_record"
    DROP COLUMN "reasoning_content",
    DROP COLUMN "reply";
