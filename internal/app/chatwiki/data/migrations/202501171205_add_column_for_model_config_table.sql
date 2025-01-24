-- +goose Up

ALTER TABLE "chat_ai_model_config"
    ADD COLUMN "show_model_name" varchar(100) default '';

COMMENT ON COLUMN "chat_ai_model_config"."show_model_name" IS '用于显示的模型名称';
