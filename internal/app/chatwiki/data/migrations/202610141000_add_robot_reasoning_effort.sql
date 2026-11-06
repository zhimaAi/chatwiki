-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "reasoning_effort" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_robot"."reasoning_effort" IS '机器人主模型深度思考等级，空字符串兼容旧开关配置';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN "reasoning_effort";
