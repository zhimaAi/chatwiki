-- +goose Up

ALTER TABLE "chat_ai_library"
    ADD COLUMN "use_model_switch" int2          NOT NULL DEFAULT 0,
    ADD COLUMN "summary_model_config_id" int4          NOT NULL DEFAULT 0,
    ADD COLUMN "statistics_set" varchar(1000)          NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_library"."use_model_switch" IS '嵌入开关 0:不开启 1：开启';
COMMENT ON COLUMN "chat_ai_library"."summary_model_config_id" IS 'AI总结配置id';
COMMENT ON COLUMN "chat_ai_library"."statistics_set" IS '统计设置';

