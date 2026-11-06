-- +goose Up

ALTER TABLE "public"."chat_ai_faq_files"
    ADD COLUMN "chunk_reasoning_effort" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."chunk_reasoning_effort" IS 'FAQ提取深度思考等级，空字符串兼容旧开关配置';

-- +goose Down

ALTER TABLE "public"."chat_ai_faq_files"
    DROP COLUMN "chunk_reasoning_effort";
