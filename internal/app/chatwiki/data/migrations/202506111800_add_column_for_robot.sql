-- +goose Up

ALTER TABLE "chat_ai_robot"
ADD COLUMN "cache_config" varchar(100) NOT NULL DEFAULT '{"cache_switch":0,"valid_time":86400}';

COMMENT ON COLUMN "chat_ai_robot"."cache_config" IS '缓存配置:{"cache_switch":0,"valid_time":86400}';