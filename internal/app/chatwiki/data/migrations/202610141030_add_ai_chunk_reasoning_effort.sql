-- +goose Up

ALTER TABLE "public"."chat_ai_library"
    ADD COLUMN "ai_chunk_reasoning_effort" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library"."ai_chunk_reasoning_effort" IS 'AI分段深度思考等级，空字符串兼容旧开关配置';

ALTER TABLE "public"."chat_ai_library_file"
    ADD COLUMN "ai_chunk_reasoning_effort" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file"."ai_chunk_reasoning_effort" IS 'AI分段深度思考等级，空字符串兼容旧开关配置';

-- +goose Down

ALTER TABLE "public"."chat_ai_library_file"
    DROP COLUMN "ai_chunk_reasoning_effort";

ALTER TABLE "public"."chat_ai_library"
    DROP COLUMN "ai_chunk_reasoning_effort";
