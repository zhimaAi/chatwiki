-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "optimize_question_reasoning_effort" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_robot"."optimize_question_reasoning_effort" IS '问题优化深度思考等级，空字符串兼容旧开关配置';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN "optimize_question_reasoning_effort";
