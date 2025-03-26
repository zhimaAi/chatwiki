-- +goose Up

ALTER TABLE "chat_ai_library" ADD COLUMN "graph_switch" int2 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library" ADD COLUMN "graph_model_config_id" int2 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library" ADD COLUMN "graph_use_model" varchar(100) NOT NULL DEFAULT '';
COMMENT ON COLUMN "chat_ai_library"."graph_switch" IS '知识图谱开关 1开启 0关闭';
COMMENT ON COLUMN "chat_ai_library"."graph_model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "chat_ai_library"."graph_use_model" IS '构建知识图谱的对话大模型(枚举值)';

