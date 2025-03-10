-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "work_flow_model_config_ids" int4[] NOT NULL DEFAULT '{}';

COMMENT ON COLUMN "chat_ai_robot"."work_flow_model_config_ids" IS '工作流模型配置ID集';
